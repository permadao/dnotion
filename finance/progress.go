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

func (f *Finance) findExistingWeeklyPage(paymentDateStr, finNid string) (*notion.Page, error) {
	// 从配置中获取 Guild 名称
	guildName, exists := f.db.FinanceDBsComments[finNid]
	if !exists {
		return nil, fmt.Errorf("fin_nid: %s not found in finance_dbs_comments", finNid)
	}

	// 将日期字符串转换为 time.Time
	paymentDate, err := time.Parse("2006-01-02", paymentDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid paymentDateStr format, expected 'YYYY-MM-DD': %s", err.Error())
	}

	// 查询符合条件的页面
	pages, err := f.db.GetPages(f.db.WeeklyDBs, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						Equals: &paymentDate,
					},
				},
			},
			{
				Property: "Guild",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Select: &notion.SelectDatabaseQueryFilter{
						Equals: guildName, // 使用 SelectDatabaseQueryFilter 匹配 Guild
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query WeeklyDBs for date: %s, guild: %s, error: %s", paymentDateStr, guildName, err.Error())
	}

	// 如果有匹配的页面，返回第一个
	if len(pages) > 0 {
		return &pages[0], nil
	}

	// 如果没有匹配的页面，返回 nil
	return nil, nil
}

func (f *Finance) UpdateWeeklyDBs(finNid string, totalAmount float64, paymentDateStr string, actualToken string, actualPrice float64) error {
	log.Info("Updating WeeklyDBs", "fin_nid", finNid, "totalAmount", totalAmount)

	// 查询是否存在相同时间和公会的记录
	existingPage, err := f.findExistingWeeklyPage(paymentDateStr, finNid)
	if err != nil {
		return fmt.Errorf("failed to check existing WeeklyDBs for fin_nid: %s, error: %s", finNid, err.Error())
	}

	// 从配置中获取对应的 Guild 名称
	guildName, exists := f.db.FinanceDBsComments[finNid]
	if !exists {
		return fmt.Errorf("fin_nid: %s not found in finance_dbs_comments", finNid)
	}

	if existingPage != nil {
		// 如果已存在相同时间和公会的记录，则修改该记录
		props := existingPage.Properties.(notion.DatabasePageProperties)

		// 获取现有 USD 和 BP 的值
		currentUSD := 0.0
		if usdProp, ok := props["USD"]; ok && usdProp.Number != nil {
			currentUSD = *usdProp.Number
		}

		currentBP := 0.0
		if bpProp, ok := props["BP"]; ok && bpProp.Number != nil {
			currentBP = *bpProp.Number
		}
		updateData := schema.WeeklyData{
			NID:   existingPage.ID, // 设置 NID，用于更新页面
			Date:  paymentDateStr,  // 时间
			Guild: guildName,       // Guild 列
		}

		// 动态更新 USD 或 BP
		if actualToken == schema.TokenAR {
			updateData.USD = currentUSD + totalAmount
			updateData.ARCurrencyRate = actualPrice
		} else if actualToken == schema.TokenBP {
			updateData.BP = currentBP + totalAmount
		} else {
			return fmt.Errorf("unsupported actualToken: %s", actualToken)
		}
		// 打印待创建的页面数据
		log.Info("Attempting to create page in WeeklyDBs", "WeeklyDBs", f.db.WeeklyDBs, "updateData", updateData)

		if err := f.db.UpdatePage(&updateData); err != nil {
			return fmt.Errorf("failed to update existing page for fin_nid: %s, error: %s", finNid, err.Error())
		}

		log.Info("Updated existing record in WeeklyDBs", "page_id", existingPage.ID, "guild", guildName, "totalAmount", totalAmount)
		return nil
	}

	// 如果不存在相同时间和公会的记录，则新增记录
	weekLastID, err := f.db.GetLastID(f.db.WeeklyDBs)
	if err != nil {
		return fmt.Errorf("failed to fetch the last ID for fin_nid: %s, error: %s", f.db.WeeklyDBs, err.Error())
	}

	// 解析出下一个 ID (假设 ID 是数字格式)
	newID := weekLastID + 1

	// 构造新的记录数据结构
	newPageData := schema.WeeklyData{
		ID:    fmt.Sprintf("%d", newID), // 假设 ID 为字符串格式
		Guild: guildName,                // Guild 列
		Date:  paymentDateStr,           // 时间列
	}

	// 动态设置 USD 或 BP
	if actualToken == "AR" {
		newPageData.USD = totalAmount
		newPageData.ARCurrencyRate = actualPrice
	} else if actualToken == "BP" {
		newPageData.BP = totalAmount
	} else {
		return fmt.Errorf("unsupported actualToken: %s", actualToken)
	}
	// 打印待创建的页面数据
	log.Info("Attempting to create page in WeeklyDBs", "WeeklyDBs", f.db.WeeklyDBs, "newPageData", newPageData)

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

	// 构造基础查询条件（前两个条件）
	queryFilters := []notion.DatabaseQueryFilter{
		{
			Property: "Status",
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				Status: &notion.StatusDatabaseQueryFilter{
					Equals: schema.StatusNotStarted,
				},
			},
		},
		{
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
	}

	// 动态添加第三个条件
	if actualToken == schema.TokenAR {
		// 如果 actualToken 是 "AR"，构造一个 Or 条件，筛选 USD 或空值
		queryFilters = append(queryFilters, notion.DatabaseQueryFilter{
			Or: []notion.DatabaseQueryFilter{
				{
					Property: "Workload Token",
					DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
						Rollup: &notion.RollupDatabaseQueryFilter{
							Any: &notion.DatabaseQueryPropertyFilter{
								Select: &notion.SelectDatabaseQueryFilter{
									Equals: schema.TokenUSD, // 等于 "USD"
								},
							},
						},
					},
				},
				{
					Property: "Workload Token",
					DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
						Rollup: &notion.RollupDatabaseQueryFilter{
							Any: &notion.DatabaseQueryPropertyFilter{
								Select: &notion.SelectDatabaseQueryFilter{
									IsEmpty: true, // 正确
								},
							},
						},
					},
				},
			},
		})
	} else if actualToken == schema.TokenBP {
		// 如果 actualToken 是 "BP"，仅筛选 BP
		queryFilters = append(queryFilters, notion.DatabaseQueryFilter{
			Property: "Workload Token",
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				Rollup: &notion.RollupDatabaseQueryFilter{
					Any: &notion.DatabaseQueryPropertyFilter{
						Select: &notion.SelectDatabaseQueryFilter{
							Equals: schema.TokenBP, // 等于 "BP"
						},
					},
				},
			},
		})
	}

	// 构造查询
	pages, err := f.db.GetPages(finNid, &notion.DatabaseQueryFilter{
		And: queryFilters, // 使用动态构建的查询条件
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

	// 将 totalAmount 填入 WeeklyDBs 表的 USD 列
	if err := f.UpdateWeeklyDBs(finNid, totalAmount, paymentDateStr, actualToken, actualPrice); err != nil {
		msg := fmt.Sprintf("Failed to update WeeklyDBs for fin_nid: %s, error: %s", finNid, err.Error())
		log.Error(msg)
		errs = append(errs, msg)
	}

	log.Info("Update done", "fin_nid", finNid, "totalAmount", totalAmount, "time", time.Since(t))
	return
}
