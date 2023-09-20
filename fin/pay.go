package fin

import (
	"context"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/utils"
	log "github.com/sirupsen/logrus"
)

func (n *Finance) PayAll() {
	for _, v := range n.NotionDB.FinanceDBs {
		t := time.Now()
		log.Info("Paying, fid", v)

		n.Pay(v)

		log.Infof("Finance payment, %s updated, since: %v", v, time.Since(t))
	}
}

func (n *Finance) Pay(fnid string) {
	// get Status is In progress
	pages := n.NotionDB.GetAllPagesFromDB(fnid, &notion.DatabaseQueryFilter{
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
		if _, err := n.NotionDB.Client.UpdatePage(context.Background(), page.ID,
			notion.UpdatePageParams{
				DatabasePageProperties: notion.DatabasePageProperties{
					"Status": notion.DatabasePageProperty{
						Status: &notion.SelectOptions{Name: "Done"},
					},
				},
			}); err != nil {
			log.Errorf("Update nid/id: %v/%v to `done` failed. %v", fnid, page.ID, err)
			continue
		}

		// payment
		tx, err := n.everpay.Transfer(
			"arweave,ethereum-ar-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,0x4fadc7a98f2dc96510e42dd1a74141eeae0c1543",
			utils.FloatToBigInt(token), wallet,
			`{"appName": "permadao-payroll", "permadaoUrl": "`+page.URL+`"}`)
		if err != nil {
			log.Errorf("Payment failed nid/id: %v/%v. %v", fnid, page.ID, err)
			// rollback
			if _, err := n.NotionDB.Client.UpdatePage(context.Background(), page.ID,
				notion.UpdatePageParams{
					DatabasePageProperties: notion.DatabasePageProperties{
						"Status": notion.DatabasePageProperty{
							Status: &notion.SelectOptions{Name: "In progress"},
						},
					},
				}); err != nil {
				log.Errorf("rollback nid/id: %v/%v to `In progress` failed. %v", fnid, page.ID, err)
			}
			continue
		}

		// update receipt
		receipt := "https://scan.everpay.io/tx/" + tx.HexHash()
		if _, err := n.NotionDB.Client.UpdatePage(context.Background(), page.ID,
			notion.UpdatePageParams{
				DatabasePageProperties: notion.DatabasePageProperties{
					"Receipt(url)": notion.DatabasePageProperty{
						URL: &receipt,
					},
				},
			}); err != nil {
			log.Errorf("Update nid/id: %v/%v receipt failed. %v", fnid, page.ID, err)
		}

	}
}
