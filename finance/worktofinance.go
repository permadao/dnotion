package finance

import (
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/db/schema"
	log "github.com/sirupsen/logrus"
)

func (f *Finance) UpdateAllWorkToFin() (errlogs []string) {
	for i, v := range f.db.WorkloadDBs {
		t := time.Now()
		log.Info("Update workload to finance, wid", v)

		errs := f.UpdateWorkToFin(v, f.db.FinanceDBs[i])
		errlogs = append(errlogs, errs...)

		log.Infof("Workload to Finance, %s/%s updated, since:%v\n\n", v, f.db.FinanceDBs[i], time.Since(t))
	}
	return
}

func (f *Finance) UpdateWorkToFin(workNid, finNid string) (errlogs []string) {
	// get last Page id
	wPageID, err := f.db.GetLastIDFromDB(workNid)
	if err != nil {
		msg := fmt.Sprintf("get last id from page failed:%s, workload nid: %s", err.Error(), workNid)
		log.Error(msg)
		errlogs = append(errlogs, msg)
		return
	}
	fPageID, err := f.db.GetLastIDFromDB(finNid)
	if err != nil {
		msg := fmt.Sprintf("get last id from page failed:%s, finance nid: %s", err.Error(), finNid)
		log.Error(msg)
		errlogs = append(errlogs, msg)
		return

	}

	for fPageID < wPageID {
		fPageID++

		fPageIDStr := fmt.Sprintf("%d", fPageID)
		wpage, err := f.db.GetPageFromDBByID(workNid, fPageIDStr)
		if err != nil {
			msg := fmt.Sprintf("error getting page from workload DB, wid:%s fpid:%s, err:%s", workNid, fPageIDStr, err)
			log.Error(msg)
			errlogs = append(errlogs, msg)
			continue
		}
		wpagep := wpage.Properties.(notion.DatabasePageProperties)
		workloadData := db.NewWrokloadDataFromProps(wpage.ID, &wpagep)
		wusd := workloadData.Usd

		// generate DatabasePageProperties
		dpp := schema.FinData{}
		dpp.ID = fPageIDStr
		dpp.WorkloadID = wpage.ID
		dpp.Amount = wusd
		//dpp.ActualToken = "USDC"

		// add contributors to properties
		if workloadData.Contributor != "" {
			if v, ok := f.uidToNid[workloadData.Contributor]; ok {
				dpp.Contributor = v
			}
		}

		// create finance expense
		if err := f.db.CreatePage(finNid, &dpp); err != nil {
			msg := "create page failed:" + err.Error()
			log.Error(msg)
			errlogs = append(errlogs, msg)
		}
	}
	return

}
