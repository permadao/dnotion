package fin

import (
	"context"
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db"
	log "github.com/sirupsen/logrus"
)

func (f *Finance) UpdateAllWorkToFin() (errlogs []string) {
	for i, v := range db.DB.WorkloadDBs {
		t := time.Now()
		log.Info("Update workload to finance, wid", v)

		errs := f.UpdateWorkToFin(v, db.DB.FinanceDBs[i])
		errlogs = append(errlogs, errs...)

		log.Infof("Workload to Finance, %s/%s updated, since:%v\n\n", v, db.DB.FinanceDBs[i], time.Since(t))
	}
	return
}

func (f *Finance) UpdateWorkToFin(workNid, finNid string) (errlogs []string) {
	// get last Page id
	wPageID := db.DB.GetLastIDFromDB(workNid)
	fPageID := db.DB.GetLastIDFromDB(finNid)
	// wPageID := 422
	// fPageID := 420

	for fPageID < wPageID {
		fPageID++

		fPageIDStr := fmt.Sprintf("%d", fPageID)
		wpage, err := db.DB.GetPageFromDBByID(workNid, fPageIDStr)
		if err != nil {
			msg := fmt.Sprintf("error getting page from workload DB, wid:%s fpid:%s, err:%s", workNid, fPageIDStr, err)
			log.Error(msg)
			errlogs = append(errlogs, msg)
			continue
		}
		wpagep := wpage.Properties.(notion.DatabasePageProperties)
		wusd := 0.0
		if wpagep["USD"].Number != nil {
			wusd = *wpagep["USD"].Number
		} else if wpagep["USD"].Formula != nil {
			wusd = *wpagep["USD"].Formula.Number
		}

		// generate DatabasePageProperties
		dpp := notion.DatabasePageProperties{
			"ID": notion.DatabasePageProperty{
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: fPageIDStr,
						},
					},
				},
			},
			"Workload ID": notion.DatabasePageProperty{
				Relation: []notion.Relation{
					{ID: wpage.ID},
				},
			},
			// "Actual Token": notion.DatabasePageProperty{
			// 	Select: &notion.SelectOptions{Name: "USDC"},
			// },
			"Amount": notion.DatabasePageProperty{
				Number: &wusd,
			},
		}
		// add contributors to properties
		if len(wpagep["Contributor"].People) > 0 {
			if v, ok := f.uidToNid[wpagep["Contributor"].People[0].BaseUser.ID]; ok {
				dpp["Contributor"] = notion.DatabasePageProperty{
					Relation: []notion.Relation{
						{ID: v},
					}}
			}
		}

		// create finance expense
		if _, err := db.DB.DBClient.CreatePage(context.Background(), notion.CreatePageParams{
			ParentType:             notion.ParentTypeDatabase,
			ParentID:               finNid,
			DatabasePageProperties: &dpp,
		}); err != nil {
			msg := "create page failed:" + err.Error()
			log.Error(msg)
			errlogs = append(errlogs, msg)
		}
	}
	return

}
