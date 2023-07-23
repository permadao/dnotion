package dnotion

import (
	"context"
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
)

func (n *DNotion) UpdateAllWorkToFin() {
	for i, v := range n.workloadDBs {
		t := time.Now()
		fmt.Println("Update workload to finance, wid", v)

		n.UpdateWorkToFin(v, n.financeDBs[i])

		fmt.Printf("Workload to Finance, %s/%s updated, since:%v\n\n", v, n.financeDBs[i], time.Since(t))
	}
}

func (n *DNotion) UpdateWorkToFin(workNid, finNid string) {
	// get last Page id
	wPageID := n.GetLastIDFromDB(workNid)
	fPageID := n.GetLastIDFromDB(finNid)
	// wPageID := 422
	// fPageID := 420

	for fPageID < wPageID {
		fPageID++

		fPageIDStr := fmt.Sprintf("%d", fPageID)
		wpage := n.GetPageFromDBByID(workNid, fPageIDStr)
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
			"Actual USD": notion.DatabasePageProperty{
				Number: &wusd,
			},
		}
		// add contributors to properties
		if len(wpagep["Contributor"].People) > 0 {
			if v, ok := n.uidToNid[wpagep["Contributor"].People[0].BaseUser.ID]; ok {
				dpp["Contributor"] = notion.DatabasePageProperty{
					Relation: []notion.Relation{
						{ID: v},
					}}
			}
		}

		// create finance expense
		if _, err := n.Client.CreatePage(context.Background(), notion.CreatePageParams{
			ParentType:             notion.ParentTypeDatabase,
			ParentID:               finNid,
			DatabasePageProperties: &dpp,
		}); err != nil {
			fmt.Println("Create page failed:", err)
		}
	}

}
