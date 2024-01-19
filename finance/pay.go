package finance

import (
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/utils"
)

func (f *Finance) PayAll() (errlogs []string) {
	for _, v := range f.db.FinanceDBs {
		t := time.Now()
		log.Info("Paying, fid", v)

		errs := f.Pay(v)
		errlogs = append(errlogs, errs...)

		log.Info("Finance payment", "updated", v, "since", time.Since(t))
	}
	return
}

func (f *Finance) Pay(fnid string) (errs []string) {
	// get Status is In progress
	pages, err := f.db.GetPages(fnid, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			notion.DatabaseQueryFilter{
				Property: "Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Status: &notion.StatusDatabaseQueryFilter{
						Equals: schema.StatusInProgress,
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
		finData := db.NewFinDataFromPage(page)
		token := finData.TargetAmount

		wallet := ""
		if finData.Contributor != "" {
			if v, ok := f.nidToWallet[finData.Contributor]; ok {
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
		finData.Status = schema.StatusDone
		if err := f.db.UpdatePage(finData); err != nil {
			msg := fmt.Sprintf("Update nid/id: %v/%v to `done` failed. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
			continue
		}

		// payment
		tx, err := f.everpay.Transfer(
			"arweave,ethereum-ar-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,0x4fadc7a98f2dc96510e42dd1a74141eeae0c1543",
			utils.FloatToBigInt(token), wallet,
			`{"appName": "`+"dnotion"+`", "permadaoUrl": "`+page.URL+`"}`)
		if err != nil {
			msg := fmt.Sprintf("Payment failed nid/id: %v/%v. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
			// rollback
			finData.Status = schema.StatusInProgress
			if err := f.db.UpdatePage(finData); err != nil {
				msg := fmt.Sprintf("rollback nid/id: %v/%v to `In progress` failed. %v", fnid, page.ID, err)
				log.Error(msg)
				errs = append(errs, msg)
			}
			continue
		}

		// update receipt
		receipt := "https://scan.everpay.io/tx/" + tx.HexHash()
		finData.ReceiptUrl = receipt
		if err := f.db.UpdatePage(finData); err != nil {
			msg := fmt.Sprintf("Update nid/id: %v/%v receipt failed. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
		}
	}
	return
}
