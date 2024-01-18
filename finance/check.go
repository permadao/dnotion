package finance

import (
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
)

// check db count and return nid from failed dbs
func (f *Finance) CheckAllDbsCountAndID() (faileddbs []string) {
	faileddbs = append(faileddbs, f.CheckDbsCountAndID(f.db.WorkloadDBs)...)
	faileddbs = append(faileddbs, f.CheckDbsCountAndID(f.db.FinanceDBs)...)
	return
}

func (f *Finance) CheckDbsCountAndID(dbs []string) (faileddbs []string) {
	for _, nid := range dbs {
		t := time.Now()
		log.Info("Checking count and id fnid", "nid", nid)
		count, err := f.db.GetCount(nid)
		if err != nil {
			msg := fmt.Sprintf("get count failed, nid: %s, err: %v\n", nid, err)
			log.Error(msg)
			faileddbs = append(faileddbs, msg)
			continue
		}
		lastid, err := f.db.GetLastID(nid)
		if err != nil {
			msg := fmt.Sprintf("get last id failed, nid: %s, err: %v\n", nid, err)
			log.Error(msg)
			faileddbs = append(faileddbs, msg)
			continue
		}
		if count != lastid {
			msg := fmt.Sprintf("nid: %s with wrong count and last id: %d, %d\n", nid, count, lastid)
			log.Error(msg)
			faileddbs = append(faileddbs, msg)
		}
		log.Info("Check done", "fin_nid", nid, "time", time.Since(t))
	}
	return
}

func (f *Finance) CheckAllWorkloadAndAmount() (errLogs []string) {
	for i, fnid := range f.db.FinanceDBs {
		t := time.Now()
		log.Info("Checking workload and amount fnid: ", fnid)

		errs := f.CheckWorkloadAndAmount(fnid, f.db.WorkloadDBs[i])
		errLogs = append(errLogs, errs...)

		log.Info("Check done", "fin_nid", fnid, "time", time.Since(t))
	}
	return
}

func (f *Finance) CheckWorkloadAndAmount(fnid, wnid string) (errlogs []string) {
	fins, err := f.db.GetPages(fnid, nil)
	if err != nil {
		msg := fmt.Sprintf("get fins failed, fnid: %s, err: %v\n", fnid, err)
		log.Error(msg)
		errlogs = append(errlogs, msg)
		return
	}
	works, err := f.db.GetPages(wnid, nil)
	if err != nil {
		msg := fmt.Sprintf("get workloads failed, wnid: %s, err: %v\n", wnid, err)
		log.Error(msg)
		errlogs = append(errlogs, msg)
		return
	}

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
