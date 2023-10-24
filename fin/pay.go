package fin

import (
	"fmt"
	serverSchema "github.com/everFinance/go-everpay/sdk/schema"
	"time"

	"github.com/dstotijn/go-notion"
	paySchema "github.com/everFinance/go-everpay/pay/schema"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/utils"
	log "github.com/sirupsen/logrus"
)

func (n *Finance) PayAll() (errlogs []string) {
	TokenMap := n.everpay.GetTokens()

	for _, v := range db.DB.FinanceDBs {
		t := time.Now()
		log.Info("Paying, fid", v)

		errs := n.Pay(v, TokenMap)
		errlogs = append(errlogs, errs...)

		log.Infof("Finance payment, %s updated, since: %v", v, time.Since(t))
	}
	return
}

func (n *Finance) Pay(fnid string, tokenMap map[string]serverSchema.TokenInfo) (errs []string) {
	// get Status is In progress
	pages, err := db.DB.GetAllPagesFromDB(fnid, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			notion.DatabaseQueryFilter{
				Property: "Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Status: &notion.StatusDatabaseQueryFilter{
						Equals: db.StatusInProgress,
					},
				},
			},
		},
	})
	if err != nil {
		errs = append(errs, err.Error())
		return
	}

	// for payments
	for _, page := range pages {
		// get wallet and amount
		p := page.Properties.(notion.DatabasePageProperties)
		pageData := db.NewFinDataFromProps(&p)
		token := pageData.TargetAmount

		wallet := ""
		if pageData.Contributor != "" {
			if v, ok := n.nidToWallet[pageData.Contributor]; ok {
				wallet = v
			}
		}
		if wallet == "" {
			msg := fmt.Sprintf("Contributor not found, nid/id: %v/%v", fnid, page.ID)
			log.Error(msg)
			errs = append(errs, msg)
			continue
		}

		// update to done
		finData := db.FinData{}
		finData.Status = db.StatusDone
		if err := finData.UpdatePage(page.ID); err != nil {
			msg := fmt.Sprintf("Update nid/id: %v/%v to `done` failed. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
			continue
		}

		// payment
		var tx *paySchema.Transaction
		if pageInfo, ok := tokenMap[pageData.ActualToken]; !ok {
			msg := fmt.Sprintf("Unknown Token")
			log.Error(msg)
			continue
		} else {
			tx, err = n.everpay.Transfer(
				pageInfo.Tag,
				utils.FloatToBigInt(token), wallet,
				`{"appName": "`+config.Config.Everpay.AppName+`", "permadaoUrl": "`+page.URL+`"}`)
			if err != nil {
				log.WithError(err).Error("trasfer error")
			}
		}

		if err != nil {
			msg := fmt.Sprintf("Payment failed nid/id: %v/%v. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
			// rollback
			finData.Status = db.StatusInProgress
			if err := finData.UpdatePage(page.ID); err != nil {
				msg := fmt.Sprintf("rollback nid/id: %v/%v to `In progress` failed. %v", fnid, page.ID, err)
				log.Error(msg)
				errs = append(errs, msg)
			}
			continue
		}

		// update receipt
		receipt := config.Config.Everpay.ScanUrl + "/tx/" + tx.HexHash()
		receiptData := db.FinData{}
		receiptData.ReceiptUrl = receipt
		if err := receiptData.UpdatePage(page.ID); err != nil {
			msg := fmt.Sprintf("Update nid/id: %v/%v receipt failed. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
		}
	}
	return
}
