package finance

import (
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
	log "github.com/sirupsen/logrus"
)

func (f *Finance) UpdateAllFinToProgress(
	paymentDateStr,
	actualToken string, actualPrice float64,
	targetToken string, targetPrice float64,
) (errs []string) {
	for _, v := range f.db.FinanceDBs {
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
	t := time.Now()
	log.Info("update fin to progress, fin_nid: ", finNid)

	// get Status is Not started & Workload Status is Acctual txs
	pages, err := f.db.GetPages(finNid, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			notion.DatabaseQueryFilter{
				Property: "Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Status: &notion.StatusDatabaseQueryFilter{
						Equals: schema.StatusNotStarted,
					},
				},
			},
			notion.DatabaseQueryFilter{
				Property: "Workload Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Rollup: &notion.RollupDatabaseQueryFilter{
						Any: &notion.DatabaseQueryPropertyFilter{
							Status: &notion.StatusDatabaseQueryFilter{
								Equals: schema.StatusAccrual,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		msg := fmt.Sprintf("get pages failed, fin_nid:%s, error: %s", finNid, err.Error())
		log.Error(msg)
		errs = append(errs, msg)
		return
	}
	// update page
	for _, page := range pages {
		finData := schema.FinData{
			NID:         page.ID,
			ActualToken: actualToken,
			ActualPrice: actualPrice,
			TargetToken: targetToken,
			TargetPrice: targetPrice,
			Status:      schema.StatusInProgress,
			PaymentDate: paymentDateStr,
		}
		if err := f.db.UpdatePage(&finData); err != nil {
			msg := fmt.Sprintf("Update nid/id: %v/%v failed. %v", finNid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
		}
	}
	log.Infof("Update done, fin_nid: %s, time: %s", finNid, time.Since(t).String())
	return
}
