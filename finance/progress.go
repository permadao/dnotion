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

func (f *Finance) UpdateWeeklyDBs(finNid string, totalAmount float64) error {
	log.Info("Updating WeeklyDBs", "fin_nid", finNid, "totalAmount", totalAmount)
	log.Info("WeeklyDBs: %v", f.db.WeeklyDBs)
	log.Info("WeeklyDBs: %v", f.db.WeeklyDBs)

	// 获取最新的 ID
	weekLastID, err := f.db.GetLastID(f.db.WeeklyDBs)
	if err != nil {
		return fmt.Errorf("failed to fetch the last ID for fin_nid: %s, error: %s", f.db.WeeklyDBs, err.Error())
	}

	// 解析出下一个 ID (假设 ID 是数字格式)
	newID := weekLastID + 1

	// 从配置中获取对应的 Guild 名称
	guildName, exists := f.db.FinanceDBsComments[finNid]
	if !exists {
		return fmt.Errorf("fin_nid: %s not found in finance_dbs_comments", finNid)
	}

	// 构造新的记录数据结构
	newPageData := schema.WeeklyData{
		ID:    fmt.Sprintf("%d", newID), // 假设 ID 为字符串格式
		Guild: guildName,                // Guild 列
		AR:    totalAmount,              // AR 列
	}

	// 新增记录到 WeeklyDBs 表
	if err := f.db.CreatePage(f.db.WeeklyDBs, &newPageData); err != nil {
		return fmt.Errorf("failed to add new page for fin_nid: %s, error: %s", finNid, err.Error())
	}

	log.Info("Added new record to WeeklyDBs successfully", "new_id", newID, "fin_nid", finNid, "guild", guildName, "totalAmount", totalAmount)
	return nil
}

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
		wpagep, ok := page.Properties.(notion.DatabasePageProperties)
		if !ok {
			msg := fmt.Sprintf("Failed to assert properties as DatabasePageProperties in page id: %s", page.ID)
			log.Error(msg)
			errs = append(errs, msg)
			continue
		}

		workloadData := db.NewWrokloadDataFromProps(page.ID, &wpagep)
		wusd := workloadData.Amount
		totalAmount += wusd

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

	// 将 totalAmount 填入 WeeklyDBs 表的 AR 列
	if err := f.UpdateWeeklyDBs(finNid, totalAmount); err != nil {
		msg := fmt.Sprintf("Failed to update WeeklyDBs for fin_nid: %s, error: %s", finNid, err.Error())
		log.Error(msg)
		errs = append(errs, msg)
	}

	log.Info("Update done", "fin_nid", finNid, "totalAmount", totalAmount, "time", time.Since(t))
	return
}
