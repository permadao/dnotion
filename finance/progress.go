package finance

import (
	"fmt"
	"github.com/permadao/dnotion/db"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (f *Finance) UpdateAllFinToProgress(
	paymentDateStr,
	actualToken string, actualPrice float64,
	targetToken string, targetPrice float64,
) (errs []string) {
	for _, v := range f.db.FinanceDBs {
		t := time.Now()
		log.Info("Update Finance to progress", "fid", v)
		e, totalAmount := f.UpdateFinToProgress(v, paymentDateStr, actualToken, actualPrice, targetToken, targetPrice)

		errs = append(errs, e...)
		// 输出总金额
		fmt.Printf("Total Amount: %f\n", totalAmount)
		log.Info("Finance to progress", "updated", v, "since", time.Since(t))
	}
	return
}

//func (f *Finance) UpdateFinToProgress(
//	finNid, paymentDateStr,
//	actualToken string, actualPrice float64,
//	targetToken string, targetPrice float64,
//) (errs []string) {
//	t := time.Now()
//	log.Info("update fin to progress", "fin_nid", finNid)
//
//	// get Status is Not started & Workload Status is Acctual txs
//	pages, err := f.db.GetPages(finNid, &notion.DatabaseQueryFilter{
//		And: []notion.DatabaseQueryFilter{
//			notion.DatabaseQueryFilter{
//				Property: "Status",
//				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
//					Status: &notion.StatusDatabaseQueryFilter{
//						Equals: schema.StatusNotStarted,
//					},
//				},
//			},
//			notion.DatabaseQueryFilter{
//				Property: "Workload Status",
//				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
//					Rollup: &notion.RollupDatabaseQueryFilter{
//						Any: &notion.DatabaseQueryPropertyFilter{
//							Status: &notion.StatusDatabaseQueryFilter{
//								Equals: schema.StatusAccrual,
//							},
//						},
//					},
//				},
//			},
//		},
//	})
//	if err != nil {
//		msg := fmt.Sprintf("get pages failed, fin_nid:%s, error: %s", finNid, err.Error())
//		log.Error(msg)
//		errs = append(errs, msg)
//		return
//	}
//	// update page
//	for _, page := range pages {
//		finData := schema.FinData{
//			NID:         page.ID,
//			ActualToken: actualToken,
//			ActualPrice: actualPrice,
//			TargetToken: targetToken,
//			TargetPrice: targetPrice,
//			Status:      schema.StatusInProgress,
//			PaymentDate: paymentDateStr,
//		}
//		if err := f.db.UpdatePage(&finData); err != nil {
//			msg := fmt.Sprintf("Update nid/id: %v/%v failed. %v", finNid, page.ID, err)
//			log.Error(msg)
//			errs = append(errs, msg)
//		}
//	}
//	log.Info("Update done", "fin_nid", finNid, "time", time.Since(t))
//	return
//}

func (f *Finance) UpdateFinToProgress(
	finNid, paymentDateStr,
	actualToken string, actualPrice float64,
	targetToken string, targetPrice float64,
) (errs []string, totalAmount float64) {
	t := time.Now()
	log.Info("update fin to progress", "fin_nid", finNid)

	// 获取符合条件的页面
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

	// 遍历符合条件的页面，提取并累加amount值
	for _, page := range pages {
		// 将 page.Properties 转换为 notion.DatabasePageProperties

		wpagep, ok := page.Properties.(notion.DatabasePageProperties)
		if !ok {
			msg := fmt.Sprintf("Failed to assert properties as DatabasePageProperties in page id: %s", page.ID)
			log.Error(msg)
			errs = append(errs, msg)
			continue
		}

		// 使用 NewWrokloadDataFromProps 创建 workloadData 实例
		workloadData := db.NewWrokloadDataFromProps(page.ID, &wpagep)

		// 从 workloadData 获取 Amount 值
		wusd := workloadData.Amount
		totalAmount += wusd
		// 更新页面信息
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

	log.Info("Update done", "fin_nid", finNid, "totalAmount", totalAmount, "time", time.Since(t))
	return
}
