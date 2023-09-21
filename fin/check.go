package fin

import (
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db"
	log "github.com/sirupsen/logrus"
)

// check db count and return nid from failed dbs
func (f *Finance) CheckAllDbsCountAndID() (faileddbs []string) {
	faileddbs = append(faileddbs, f.CheckDbsCountAndID(db.DB.WorkloadDBs)...)
	faileddbs = append(faileddbs, f.CheckDbsCountAndID(db.DB.FinanceDBs)...)
	return
}

func (f *Finance) CheckDbsCountAndID(dbs []string) (faileddbs []string) {
	for _, nid := range dbs {
		count := db.DB.GetCountFromDB(nid)
		lastid := db.DB.GetLastIDFromDB(nid)

		if count != lastid {
			log.Errorf("nid: %s with wrong count and last id: %d, %d\n", nid, count, lastid)
			faileddbs = append(faileddbs, nid)
		}
	}
	return
}

func (f *Finance) CheckAllWorkloadAndAmount() (errLogs []string) {
	for i, fnid := range db.DB.FinanceDBs {
		t := time.Now()
		log.Info("Checking fnid", fnid)

		errs := f.CheckWorkloadAndAmount(fnid, db.DB.WorkloadDBs[i])
		errLogs = append(errLogs, errs...)

		log.Info("Check done fnid", fnid, time.Since(t))
	}
	return
}

func (f *Finance) CheckWorkloadAndAmount(fnid, wnid string) (errlogs []string) {
	fins := db.DB.GetAllPagesFromDB(fnid, nil)
	works := db.DB.GetAllPagesFromDB(wnid, nil)

	workToUSD := map[string]float64{} // id -> usd
	for _, page := range works {
		wp := page.Properties.(notion.DatabasePageProperties)
		if wp == nil {
			err := fmt.Sprintf("Find work page is nil, fnid/wid:%s/%s\n", fnid, page.ID)
			log.Error(err)
			errlogs = append(errlogs, err)
			continue
		}

		workloadUSD := 0.0
		switch wp["USD"].Type {
		case notion.DBPropTypeFormula:
			if wp["USD"].Formula != nil {
				workloadUSD = *wp["USD"].Formula.Number
			}
		case notion.DBPropTypeNumber:
			if wp["USD"].Number != nil {
				workloadUSD = *wp["USD"].Number
			}
		}
		workToUSD[page.ID] = workloadUSD
	}

	for _, page := range fins {
		p := page.Properties.(notion.DatabasePageProperties)

		actualAmount := 0.0
		if p["Amount"].Number != nil {
			actualAmount = *p["Amount"].Number
		}

		wid := p["Workload ID"].Relation[0].ID

		if _, ok := workToUSD[wid]; !ok {
			err := fmt.Sprintf("fnid/id: %v/%v check workload is failed, actual: %v, workload: %v\n",
				fnid, p["ID"].Title[0].PlainText, actualAmount, nil)
			log.Error(err)
			errlogs = append(errlogs, err)
		}
		if actualAmount != workToUSD[wid] {
			err := fmt.Sprintf("fnid/id: %v/%v check workload is failed, actual: %v, workload: %v\n",
				fnid, p["ID"].Title[0].PlainText, actualAmount, workToUSD[wid])
			log.Error(err)
			errlogs = append(errlogs, err)
		}
	}
	return
}
