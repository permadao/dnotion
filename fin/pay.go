package fin

import (
	"context"
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/utils"
	log "github.com/sirupsen/logrus"
)

func (n *Finance) PayAll() (errlogs []string) {
	for _, v := range db.DB.FinanceDBs {
		t := time.Now()
		log.Info("Paying, fid", v)

		errs := n.Pay(v)
		errlogs = append(errlogs, errs...)

		log.Infof("Finance payment, %s updated, since: %v", v, time.Since(t))
	}
	return
}

func (n *Finance) Pay(fnid string) (errs []string) {
	// get Status is In progress
	pages := db.DB.GetAllPagesFromDB(fnid, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			notion.DatabaseQueryFilter{
				Property: "Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Status: &notion.StatusDatabaseQueryFilter{
						Equals: "In progress",
					},
				},
			},
		},
	})

	// for payments
	for _, page := range pages {
		// get wallet and amount
		p := page.Properties.(notion.DatabasePageProperties)
		token := 0.0
		if p["Target Amount"].Formula != nil {
			token = *p["Target Amount"].Formula.Number
		}

		wallet := ""
		if len(p["Contributor"].Relation) > 0 {
			if v, ok := n.nidToWallet[p["Contributor"].Relation[0].ID]; ok {
				wallet = v
			}
		}

		// update to done
		if _, err := db.DB.DBClient.UpdatePage(context.Background(), page.ID,
			notion.UpdatePageParams{
				DatabasePageProperties: notion.DatabasePageProperties{
					"Status": notion.DatabasePageProperty{
						Status: &notion.SelectOptions{Name: "Done"},
					},
				},
			}); err != nil {
			msg := fmt.Sprintf("Update nid/id: %v/%v to `done` failed. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
			continue
		}

		// payment
		tx, err := n.everpay.Transfer(
			config.Config.Everpay.TokenTag,
			utils.FloatToBigInt(token), wallet,
			`{"appName": "`+config.Config.Everpay.AppName+`", "permadaoUrl": "`+page.URL+`"}`)
		if err != nil {
			msg := fmt.Sprintf("Payment failed nid/id: %v/%v. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
			// rollback
			if _, err := db.DB.DBClient.UpdatePage(context.Background(), page.ID,
				notion.UpdatePageParams{
					DatabasePageProperties: notion.DatabasePageProperties{
						"Status": notion.DatabasePageProperty{
							Status: &notion.SelectOptions{Name: "In progress"},
						},
					},
				}); err != nil {
				msg := fmt.Sprintf("rollback nid/id: %v/%v to `In progress` failed. %v", fnid, page.ID, err)
				log.Error(msg)
				errs = append(errs, msg)
			}
			continue
		}

		// update receipt
		receipt := config.Config.Everpay.ScanUrl + "/tx/" + tx.HexHash()
		if _, err := db.DB.DBClient.UpdatePage(context.Background(), page.ID,
			notion.UpdatePageParams{
				DatabasePageProperties: notion.DatabasePageProperties{
					"Receipt(url)": notion.DatabasePageProperty{
						URL: &receipt,
					},
				},
			}); err != nil {
			msg := fmt.Sprintf("Update nid/id: %v/%v receipt failed. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
		}
	}
	return
}
