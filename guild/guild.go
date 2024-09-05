package guild

import (
	"github.com/permadao/dnotion/utils"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	dbSchema "github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/guild/schema"
	"github.com/permadao/dnotion/logger"
)

var log = logger.New("guild")

type Guild struct {
	db *db.DB

	// contributors
	nidToName        map[string]string   //  contributor data nid -> contributors name
	nidToWallet      map[string]string   // contributor data nid -> contributors wallet
	notionidToName   map[string]string   // contributor data notion id -> contributors name
	notionidToID     map[string]*float64 // contributor data notion id -> contributors sort id
	nidToContributor map[string]dbSchema.ContributorData
}

func New(conf *config.Config, db *db.DB) *Guild {
	g := &Guild{
		db: db,

		nidToName:        map[string]string{},
		nidToWallet:      map[string]string{},
		notionidToName:   map[string]string{},
		notionidToID:     map[string]*float64{},
		nidToContributor: map[string]dbSchema.ContributorData{},
	}

	g.LoadContributors()
	return g
}

func (g *Guild) LoadContributors() {
	contributors, err := g.db.GetContributors(nil)
	if err != nil {
		panic(err)
	}

	for _, c := range contributors {
		if c.NotionName != "" {
			g.nidToName[c.NID] = c.NotionName
		}
		if c.Wallet != "" {
			g.nidToWallet[c.NID] = c.Wallet
		}
		if c.NotionID != "" {
			g.notionidToName[c.NotionID] = c.NotionName
			g.notionidToID[c.NotionID] = c.ID
		}
		g.nidToContributor[c.NID] = c
	}
}

func (g *Guild) GenGuilds(targetToken, date string) {
	// content stat
	hits, frontPages, _ := g.StatContent(date)

	for guildName, info := range schema.Guilds {
		// guild stat
		totalAmount, contributors, rankOfContributor, _ := g.StatWeeklyFinance(targetToken, info.FinNID, date)
		_, beforeContributors, _, _ := g.StatBeforeFinance(targetToken, info.FinNID, date)
		allTotalAmount, _, allRankOfContributor, _ := g.StatFinance(targetToken, info.FinNID)
		news := float64(0)
		for name, _ := range contributors {
			if _, ok := beforeContributors[name]; !ok {
				news++
			}
		}

		// achievements
		tags := []string{}
		if a := AGuildActiviy(len(rankOfContributor)); a != "" {
			tags = append(tags, a)
		}
		if len(rankOfContributor) > 0 && totalAmount != 0 {
			if a := AFairDistribution(rankOfContributor[0].Amount / totalAmount); a != "" {
				tags = append(tags, a)
			}
		}
		if a := AReadership(hits); a != "" && (guildName == "内容公会 - 投稿" || guildName == "品宣公会") {
			tags = append(tags, a)
		}
		if a := AMediaPicks(frontPages); a != "" && (guildName == "内容公会 - 投稿" || guildName == "品宣公会") {
			tags = append(tags, a)
		}

		guild := dbSchema.GuildData{
			Name:               guildName,
			Info:               info.Info,
			Link:               info.NID,
			Tags:               tags,
			TotalContributors:  float64(len(allRankOfContributor)),
			WeeklyContributors: float64(len(rankOfContributor)),
			NewContributors:    news,
			TotalIncentive:     allTotalAmount,
			WeeklyIncentive:    totalAmount,
			Date:               date,
			Rank:               info.Rank,
		}
		if err := g.db.CreatePage(g.db.GuildDB, &guild); err != nil {
			log.Error("create guild failed", "err", err)
		}
	}
}

func (g *Guild) GenGrade(guidNid, gradeNid, startDate, endDate string) (err error) {
	_, _, rankOfContributor, err := g.StatBetweenFinance("AR", guidNid, startDate, endDate)
	if err != nil {
		return
	}
	id, err := g.db.GetLastID(gradeNid)
	if err != nil {
		log.Error("get last id from page failed", "finance nid", gradeNid, "err", err)
		return
	}
	grades := GRankToGrade(rankOfContributor, id)

	for _, tr := range grades {
		if err = g.db.CreatePage(gradeNid, &tr); err != nil {
			log.Error("create grade page failed failed", "err", err)
			return
		}
	}

	return
}

func (g *Guild) GenDevGrade(guidNid, gradeNid, lastDate, endDate string) (err error) {
	_, _, rankOfContributor, err := g.StatBeforeFinanceByAmount(guidNid, endDate)
	if err != nil {
		return
	}
	lastD, err := notion.ParseDateTime(lastDate)
	if err != nil {
		return
	}
	id, err := g.db.GetLastID(gradeNid)
	if err != nil {
		log.Error("get last id from page failed", "finance nid", gradeNid, "err", err)
		return
	}
	developers, err := g.db.GetDeveloper(gradeNid, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						Equals: &lastD.Time,
					},
				},
			},
		},
	})
	if err != nil {
		return
	}
	insert := GRankToGradeForDev(rankOfContributor, developers, id, endDate)

	for _, tr := range insert {
		if err = g.db.CreatePage(gradeNid, &tr); err != nil {
			log.Error("create grade page failed", "err", err)
			return
		}
	}

	return
}

// 新闻工会等级计算
func (g *Guild) GenNewsGrade(guidNid, gradeNid, lastDate, endDate string) (err error) {
	_, aggrContributorsFor15weeks, aggrContributorsForAllDay, err := g.StatNewsFinance(guidNid, lastDate)
	if err != nil {
		return
	}
	// fmt.Println(aggrContributorsFor15weeks)
	news, err := g.db.GetNews(gradeNid, nil)
	if err != nil {
		return
	}
	id, err := g.db.GetLastID(gradeNid)
	if err != nil {
		log.Error("get last id from page failed", "finance nid", gradeNid, "err", err)
		return
	}
	insert := GRankToGradeForNews(aggrContributorsFor15weeks, aggrContributorsForAllDay, news, id, endDate)
	for _, tr := range insert {
		if err = g.db.UpdatePage(&tr); err != nil {
			log.Error("update grade page failed", "err", err)
			return
		}

	}

	return
}

func (g *Guild) GenPromotionSettlement(guidNid, outNid, endDate string) (err error) {
	//1 query the weekly table of promotion
	endD, err := notion.ParseDateTime(endDate)
	if err != nil {
		return
	}
	ps, err := g.db.GetPromotionStat(guidNid, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						Equals: &endD.Time,
					},
				},
			},
		},
	})
	if err != nil || len(ps) == 0 {
		return
	}
	//2 to statistics the points of promotions
	promotionPoints, err := g.db.GetPromotionPoints(ps[0].ID, nil)

	//3 output
	insert := CalculatePromotionRewards(promotionPoints, g.notionidToName, g.notionidToID, endDate)
	for _, tr := range insert {
		if err = g.db.CreatePage(outNid, &tr); err != nil {
			log.Error("create the reward of promotion's page failed", "err", err)
			return
		}
	}
	return
}

func (g *Guild) GenIncentiveStat(outNid, now string) (success bool, paymentDateCount map[string]int, err error) {
	guildFinMap := GetGuildFinMap()
	startDate := "1970-01-01"
	insert := []dbSchema.CIncentiveGuild{}
	update := []dbSchema.CIncentiveGuild{}
	paymentDateCount = map[string]int{}
	hisRecords, err := g.GetHisIncentiveRecords(now)
	if err != nil {
		log.Error("GetHisIncentiveRecords failed", "err", err)
		return
	}
	//获取总数
	counts, err := g.db.GetCount(outNid)
	pageId := float64(counts)
	for guild, nid := range guildFinMap {
		weekStatResults, paymentDate, err := g.StatWeeklyFinanceGroupByCNID(nid, now)
		if err != nil {
			log.Error("statistic the incentive of various guild failed", "err", err)
			return false, paymentDateCount, err
		}
		if len(weekStatResults) == 0 {
			continue
		}
		completeStatResults, err2 := g.StatBetweenFinanceGroupByCNidToken(nid, startDate, now)
		if err2 != nil {
			log.Error("statistic the total incentive of various guild failed", "err", err2)
			return false, paymentDateCount, err2
		}
		insertRecords, updateRecords := GenStatRecords(completeStatResults, weekStatResults, hisRecords, guild, now, paymentDate, &pageId, g)
		insert = append(insert, insertRecords...)
		update = append(update, updateRecords...)
		paymentDateCount[paymentDate]++
	}
	for _, datum := range insert {
		if err = g.db.CreatePage(outNid, &datum); err != nil {
			log.Error("create the records of incentive's statistic page failed", "err", err)
			return
		}
	}
	for _, datum := range update {
		if err = g.db.UpdatePage(&datum); err != nil {
			log.Error("update the records of incentive's statistic page failed", "err", err)
			return
		}
	}
	success = true
	return
}

func (g *Guild) GenTotalIncentiveStat(outNid string, paymentDateMap map[string]int) (err error) {
	incentiveData := []dbSchema.CIncentiveGuild{}
	historyIncentiveByDate := map[string]schema.ResultSepToken{}
	records := map[string]dbSchema.CIncentive{}
	for paymentDateStr, _ := range paymentDateMap {
		totalIncentiveRecords, err := g.GetHisTotalIncentiveRecords(paymentDateStr)
		if err != nil {
			log.Error("GetHisTotalIncentiveRecords failed", "err", err)
			return err
		}
		records = utils.MergeMaps(records, totalIncentiveRecords)
		paymentDate, err1 := notion.ParseDateTime(paymentDateStr)
		if err1 != nil {
			return err1
		}
		data, err2 := g.db.GetIncentiveGuildData(utils.CincentiveWeeklyGuildRs, &notion.DatabaseQueryFilter{
			And: []notion.DatabaseQueryFilter{
				{
					Property: "Payment Date",
					DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
						Date: &notion.DatePropertyFilter{
							Equals: &paymentDate.Time,
						},
					},
				},
			},
		})
		if err2 != nil {
			return err2
		}
		if len(data) > 0 {
			historyIncentive, err3 := g.GetHisTotalIncentiveRecordsSepToke(paymentDate)
			if err3 != nil {
				return err3
			}
			historyIncentiveByDate[paymentDateStr] = historyIncentive
			incentiveData = append(incentiveData, data...)
		}
	}
	if len(incentiveData) == 0 {
		return
	}
	//获取总数
	counts, err := g.db.GetCount(outNid)
	pageId := float64(counts)
	insert, update := CalTotalIncentive(incentiveData, historyIncentiveByDate, records, &pageId)
	for _, tr := range insert {
		if err = g.db.CreatePage(outNid, &tr); err != nil {
			log.Error("create the incentive_weekly page failed", "err", err)
			return
		}
	}
	for _, tr := range update {
		if err = g.db.UpdatePage(&tr); err != nil {
			log.Error("update the incentive_weekly page failed", "err", err)
			return
		}
	}
	return
}

func (g *Guild) GetHisTotalIncentiveRecordsSepToke(paymentDate notion.DateTime) (schema.ResultSepToken, error) {
	result := schema.ResultSepToken{}
	hisData, err := g.db.GetIncentiveGuildData(utils.CincentiveWeeklyGuildRs, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						Before: &paymentDate.Time,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	for _, datum := range hisData {
		token := datum.Token
		if historyIncentive, ok := result[token]; ok {
			historyIncentive[datum.NotionID] += datum.WeeklyIncentive
			result[token] = historyIncentive
		} else {
			historyIncentive := map[string]float64{}
			historyIncentive[datum.NotionID] += datum.WeeklyIncentive
			result[token] = historyIncentive
		}
	}
	return result, nil
}

func (g *Guild) GetHisTotalIncentiveRecords(paymentDate string) (map[string]dbSchema.CIncentive, error) {
	end, _ := notion.ParseDateTime(paymentDate)
	//本周已经生成的数据
	records, err := g.db.GetCIncentiveData(utils.CincentiveWeeklyRs, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						Equals: &end.Time,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	hisRecords := make(map[string]dbSchema.CIncentive)
	for _, record := range records {
		key := GetTIKey(record)
		hisRecords[key] = record
	}
	return hisRecords, nil
}

func (g *Guild) GetHisIncentiveRecords(acDate string) (map[string]dbSchema.CIncentiveGuild, error) {
	//这周的时间范围
	endDateTime, _ := time.Parse("2006-01-02", acDate)
	startDate := endDateTime.AddDate(0, 0, -6).Format("2006-01-02")
	start, _ := notion.ParseDateTime(startDate)
	end, _ := notion.ParseDateTime(acDate)
	//本周已经生成的数据
	records, err := g.db.GetIncentiveGuildData(utils.CincentiveWeeklyGuildRs, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrAfter: &start.Time,
					},
				},
			},
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrBefore: &end.Time,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	hisRecords := make(map[string]dbSchema.CIncentiveGuild)
	for _, record := range records {
		key := GetKey(record)
		hisRecords[key] = record
	}
	return hisRecords, nil
}

func (g *Guild) IsExistRecord(endDate string) (isExist bool, err error) {
	parse, _ := time.Parse("2006-01-02", endDate)
	startDate := parse.AddDate(0, 0, -6).Format("2006-01-02")
	start, err := notion.ParseDateTime(startDate)
	if err != nil {
		return
	}
	end, err := notion.ParseDateTime(endDate)
	if err != nil {
		return
	}

	fins, err := g.db.GetFinances("caac7a1aefcc4ed0b02b8adbc106f021", &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			notion.DatabaseQueryFilter{
				Property: "Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Status: &notion.StatusDatabaseQueryFilter{
						Equals: "Done",
					},
				},
			},
			notion.DatabaseQueryFilter{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrAfter: &start.Time,
					},
				},
			},
			notion.DatabaseQueryFilter{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrBefore: &end.Time,
					},
				},
			},
		},
	})
	if err != nil {
		return
	}
	return len(fins) > 0 && fins[0].PaymentDate != "", err
}

func (g *Guild) IsExistIncentiveStatRecord(endDateStr string) bool {
	endDate, _ := notion.ParseDateTime(endDateStr)
	data, err := g.db.GetIncentiveGuildData(utils.CincentiveWeeklyGuildRs, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrAfter: &endDate.Time,
					},
				},
			},
		},
	})
	if err != nil {
		return true
	}
	return len(data) > 0
}
