package dnotion

import (
	"context"
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
)

func (n *DNotion) UpdateAllFinToProgress(paymentDateStr, targetToken string, tokenPrice float64) {
	for _, v := range n.financeDBs {
		t := time.Now()
		fmt.Println("Update Finance to progress, fid", v)

		n.UpdateFinToProgress(v, paymentDateStr, targetToken, tokenPrice)

		fmt.Printf("Finance to progress, %s updated, since: %v\n\n", v, time.Since(t))
	}
}

func (n *DNotion) UpdateFinToProgress(finNid, paymentDateStr, targetToken string, tokenPrice float64) {
	paymentDate, err := notion.ParseDateTime(paymentDateStr)
	if err != nil {
		fmt.Println("invalid payment date", paymentDateStr)
		return
	}

	// get Status is Not started & Workload Status is Acctual txs
	pages := n.GetAllPagesFromDB(finNid, &notion.DatabaseQueryFilter{
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
		if _, err := n.Client.UpdatePage(context.Background(), page.ID, notion.UpdatePageParams{
			DatabasePageProperties: notion.DatabasePageProperties{
				"Target Token": notion.DatabasePageProperty{
					Select: &notion.SelectOptions{Name: targetToken},
				},
				"Token Price": notion.DatabasePageProperty{
					Number: &tokenPrice,
				},
				"Status": notion.DatabasePageProperty{
					Status: &notion.SelectOptions{Name: "In progress"},
				},
				"Payment Date": notion.DatabasePageProperty{
					Date: &notion.Date{Start: paymentDate},
				},
			},
		}); err != nil {
			fmt.Printf("Update nid/id: %v/%v failed. %v\n", finNid, page.ID, err)
		}
	}
}
