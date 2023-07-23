package dnotion

import (
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
)

func (n *DNotion) CheckAllDbsCountAndID() {
	// n.CheckDbsCountAndID(n.taskDBs)
	n.CheckDbsCountAndID(n.workloadDBs)
	n.CheckDbsCountAndID(n.financeDBs)
}

func (n *DNotion) CheckDbsCountAndID(dbs []string) {
	for _, nid := range dbs {
		count := n.GetCountFromDB(nid)
		lastid := n.GetLastIDFromDB(nid)

		if count != lastid {
			fmt.Printf("nid: %s with wrong count and last id: %d, %d\n", nid, count, lastid)
		}
	}
}

func (n *DNotion) CheckAllWorkloadAndAcutal() {
	for i, fnid := range n.financeDBs {
		t := time.Now()
		fmt.Println("Checking fnid", fnid)

		n.CheckWorkloadAndAcutal(fnid, n.workloadDBs[i])

		fmt.Println("Check done fnid", fnid, time.Since(t))
		fmt.Println()
	}
}

func (n *DNotion) CheckWorkloadAndAcutal(fnid, wnid string) {
	fins := n.GetAllPagesFromDB(fnid, nil)
	works := n.GetAllPagesFromDB(wnid, nil)

	workToUSD := map[string]float64{} // id -> usd
	for _, page := range works {
		wp := page.Properties.(notion.DatabasePageProperties)
		if wp == nil {
			fmt.Printf("Find work page is nil, fnid/wid:%s/%s\n", fnid, page.ID)
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

		actualUSD := 0.0
		if p["Actual USD"].Number != nil {
			actualUSD = *p["Actual USD"].Number
		}

		wid := p["Workload ID"].Relation[0].ID

		if _, ok := workToUSD[wid]; !ok {
			fmt.Printf("fnid/id: %v/%v check workload is failed, actual: %v, workload: %v\n",
				fnid, p["ID"].Title[0].PlainText, actualUSD, nil)
		}
		if actualUSD != workToUSD[wid] {
			fmt.Printf("fnid/id: %v/%v check workload is failed, actual: %v, workload: %v\n",
				fnid, p["ID"].Title[0].PlainText, actualUSD, workToUSD[wid])
		}
	}
}
