package finance

import (
	"fmt"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/utils"
)

func (f *Finance) PayAll() (errlogs []string) {
	for _, v := range f.db.FinanceDBs {
		t := time.Now()
		log.Info("Paying, fid", v)

		errs := f.Pay(v)
		errlogs = append(errlogs, errs...)

		log.Info("Finance payment", "updated", v, "since", time.Since(t))
	}
	return
}

func (f *Finance) Pay(fnid string) (errs []string) {
	// 初始化分发金额记录
	distributedAR := 0.0
	distributedBP := 0.0

	// get Status is In progress
	pages, err := f.db.GetPages(fnid, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Status: &notion.StatusDatabaseQueryFilter{
						Equals: schema.StatusInProgress,
					},
				},
			},
		},
	})
	if err != nil {
		errs = append(errs, err.Error())
		return
	}

	// for payments
	for _, page := range pages {
		// get wallet and amount
		finData := db.NewFinDataFromPage(page)

		wallet := ""
		if finData.Contributor != "" {
			if v, ok := f.nidToWallet[finData.Contributor]; ok {
				wallet = v
			}
		}
		if wallet == "" {
			msg := fmt.Sprintf("Contributor not found, nid/id: %v/%v", fnid, page.ID)
			log.Error(msg)
			errs = append(errs, msg)
			continue
		}

		// get token info
		tokenInfo, ok := f.tokens[finData.TargetToken]
		if !ok {
			msg := fmt.Sprintf("Target token not exist: %s ; nid/id: %v/%v", finData.TargetToken, fnid, page.ID)
			log.Error(msg)
			errs = append(errs, msg)
			continue
		}

		// update to done
		finData.Status = schema.StatusDone
		if err := f.db.UpdatePage(finData); err != nil {
			msg := fmt.Sprintf("Update nid/id: %v/%v to `done` failed. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
			continue
		}

		// payment
		tx, err := f.everpay.Transfer(
			tokenInfo.Tag,
			utils.FloatToBigInt(finData.TargetAmount, tokenInfo.Decimals),
			wallet,
			`{"appName": "`+"dnotion"+`", "permadaoUrl": "`+page.URL+`"}`)
		if err != nil {
			msg := fmt.Sprintf("Payment failed nid/id: %v/%v. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
			// rollback
			finData.Status = schema.StatusInProgress
			if err := f.db.UpdatePage(finData); err != nil {
				msg := fmt.Sprintf("rollback nid/id: %v/%v to `In progress` failed. %v", fnid, page.ID, err)
				log.Error(msg)
				errs = append(errs, msg)
			}
			continue
		}

		// update receipt
		receipt := "https://scan.everpay.io/tx/" + tx.HexHash()
		finData.ReceiptUrl = receipt
		if err := f.db.UpdatePage(finData); err != nil {
			msg := fmt.Sprintf("Update nid/id: %v/%v receipt failed. %v", fnid, page.ID, err)
			log.Error(msg)
			errs = append(errs, msg)
		}

		// 记录分发金额
		if finData.TargetToken == schema.TokenAR {
			distributedAR += finData.TargetAmount
		} else if finData.TargetToken == schema.TokenBP {
			distributedBP += finData.TargetAmount
		}
	}

	// 更新 WeeklyDBs
	err = f.PAyUpdateWeeklyDbs(fnid, distributedAR, distributedBP)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Failed to update WeeklyDBs: %v", err))
	}

	return
}

func (f *Finance) ExistGuildAndStatus(guildName, status string) (*schema.WeeklyData, error) {
	// 查询数据库，筛选出符合条件的页面
	pages, err := f.db.GetPages(f.db.WeeklyDBs, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Guild",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Select: &notion.SelectDatabaseQueryFilter{
						Equals: guildName, // 使用 SelectDatabaseQueryFilter 匹配 Guild
					},
				},
			},
			{
				Property: "Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Status: &notion.StatusDatabaseQueryFilter{
						Equals: schema.StatusNotStarted,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query WeeklyDBs for guild: %s, status: %s, error: %s", guildName, status, err.Error())
	}

	if len(pages) == 0 {
		return nil, nil // 没有找到符合条件的页面
	}

	// 返回第一个匹配的页面
	page := pages[0]
	props := page.Properties.(notion.DatabasePageProperties)

	// 提取数据到 WeeklyData 结构
	weeklyData := schema.WeeklyData{
		NID:   page.ID,
		Guild: guildName,
	}
	if arProp, ok := props["Actual Distributed AR"]; ok && arProp.Number != nil {
		weeklyData.ActualDistributedAR = *arProp.Number
	}
	if bpProp, ok := props["Actual Distributed BP"]; ok && bpProp.Number != nil {
		weeklyData.ActualDistributedBP = *bpProp.Number
	}

	return &weeklyData, nil
}

// 更新 WeeklyDBs 中的分发列
func (f *Finance) PAyUpdateWeeklyDbs(fnid string, distributedAR, distributedBP float64) error {
	// 从配置中获取对应的 Guild 名称
	guildName, exists := f.db.FinanceDBsComments[fnid]
	if !exists {
		return fmt.Errorf("fin_nid: %s not found in finance_dbs_comments", fnid)
	}

	// 查询目标记录：根据 Guild 和 Status 为 NotStarted 的条件
	existingPage, err := f.ExistGuildAndStatus(guildName, schema.StatusNotStarted)
	if err != nil {
		return fmt.Errorf("failed to find existing WeeklyDBs record for guild: %s, error: %w", guildName, err)
	}

	if existingPage != nil {
		// 更新已有页面
		existingPage.ActualDistributedAR += distributedAR
		existingPage.ActualDistributedBP += distributedBP
		// 更新状态为 Done
		existingPage.Status = schema.StatusDone
		if err := f.db.UpdatePage(existingPage); err != nil {
			return fmt.Errorf("failed to update WeeklyDBs page: %w", err)
		}
		log.Info("Updated WeeklyDBs with distributed amounts", "Guild", guildName, "AR", distributedAR, "BP", distributedBP)
	}

	return nil
}
