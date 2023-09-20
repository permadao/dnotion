package fin

import (
	"context"
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
	log "github.com/sirupsen/logrus"
)

func (f *Finance) UpdateAllFinToProgress(
	paymentDateStr,
	actualToken string, actualPrice float64,
	targetToken string, targetPrice float64,
) (errs []string) {
	for _, v := range f.NotionDB.FinanceDBs {
		t := time.Now()
		log.Info("Update Finance to progress, fid", v)

		e := f.UpdateFinToProgress(v, paymentDateStr, actualToken, actualPrice, targetToken, targetPrice)
		errs = append(errs, e...)
		log.Infof("Finance to progress, %s updated, since: %v\n\n", v, time.Since(t))
	}
	return
}

func (f *Finance) UpdateFinToProgress(
	finNid, paymentDateStr,
	actualToken string, actualPrice float64,
	targetToken string, targetPrice float64,
) (errs []string) {
	paymentDate, err := notion.ParseDateTime(paymentDateStr)
	if err != nil {
		msg := fmt.Sprintf("invalid payment date %s", paymentDateStr)
		log.Error(msg)
		errs = append(errs, msg)
		return
	}

	// get Status is Not started & Workload Status is Acctual txs
	pages := f.NotionDB.GetAllPagesFromDB(finNid, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			notion.DatabaseQueryFilter{
				Property: "Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Status: &notion.StatusDatabaseQueryFilter{
						Equals: "Not started",
					},
				},
			},
			notion.DatabaseQueryFilter{
				Property: "Workload Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Rollup: &notion.RollupDatabaseQueryFilter{
						Any: &notion.DatabaseQueryPropertyFilter{
							Status: &notion.StatusDatabaseQueryFilter{
								Equals: "Accrual",
							},
						},
					},
				},
			},
		},
	})
	for _, page := range pages {
		if _, err := f.NotionDB.Client.UpdatePage(context.Background(), page.ID, notion.UpdatePageParams{
			DatabasePageProperties: notion.DatabasePageProperties{
				"Actual Token": notion.DatabasePageProperty{
					Select: &notion.SelectOptions{Name: actualToken},
				},
				"Actual Price": notion.DatabasePageProperty{
					Number: &actualPrice,
				},
				"Target Token": notion.DatabasePageProperty{
					Select: &notion.SelectOptions{Name: targetToken},
				},
				"Target Price": notion.DatabasePageProperty{
					Number: &targetPrice,
				},
				"Status": notion.DatabasePageProperty{
					Status: &notion.SelectOptions{Name: "In progress"},
				},
				"Payment Date": notion.DatabasePageProperty{
					Date: &notion.Date{Start: paymentDate},
				},
			},
		}); err != nil {
			msg := fmt.Sprintf("Update nid/id: %v/%v failed. %v", finNid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
		}
	}
	return
}
