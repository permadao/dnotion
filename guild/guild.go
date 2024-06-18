package guild

import (
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

func (g *Guild) GenIncentiveStat(outNid, now string) (success bool, paymentDateMap map[string]int, err error) {
	gfm := GetGuildFinMap()
	startDate := "1970-01-01"
	insert := []dbSchema.Incentive{}
	//pageId, err := g.db.GetLastID(outNid)
	//if err != nil {
	//	pageId = 1
	//}
	pageId := 0
	paymentDateMap = map[string]int{}
	for guild, nid := range gfm {
		_, contributorsAllTime, _, err := g.StatBetweenFinanceGroupByCNID("", nid, startDate, now)
		if err != nil {
			log.Error("statistic the incentive of various guild failed", "err", err)
			return false, paymentDateMap, err
		}
		_, contributorsThisWeek, _, paymentDate, err := g.StatWeeklyFinanceGroupByCNID("", nid, now)
		if err != nil {
			log.Error("statistic the incentive of various guild failed", "err", err)
			return false, paymentDateMap, err
		}
		insert = append(insert, GenStatRecords(contributorsAllTime, contributorsThisWeek, guild, now, paymentDate, pageId, g)...)
		paymentDateMap[paymentDate]++
	}

	for _, tr := range insert {
		if err = g.db.CreatePage(outNid, &tr); err != nil {
			log.Error("create the records of incentive's statistic page failed", "err", err)
			return
		}
	}
	success = true
	return
}

func (g *Guild) GenTotalIncentiveStat(outNid string, paymentDateMap map[string]int) (err error) {
	incentiveData := []dbSchema.Incentive{}
	historyIncentiveByDate := map[string]map[string]float64{}
	for paymentDateStr, _ := range paymentDateMap {
		paymentDate, err1 := notion.ParseDateTime(paymentDateStr)
		if err1 != nil {
			return err1
		}
		data, err2 := g.db.GetIncentiveData(&notion.DatabaseQueryFilter{
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
		if len(data) != 0 {
			historyIncentive := map[string]float64{}
			hisdata, err3 := g.db.GetIncentiveData(&notion.DatabaseQueryFilter{
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
			if err3 != nil {
				return err3
			}
			for _, d := range hisdata {
				historyIncentive[d.NotionID] += d.WeeklyIncentive
			}
			historyIncentiveByDate[paymentDateStr] = historyIncentive
		}

		incentiveData = append(incentiveData, data...)
	}
	//pageId, err := g.db.GetLastID(outNid)
	//if err != nil {
	//	pageId = 1
	//}

	pageId := 0
	insert := CalTotalIncentive(incentiveData, historyIncentiveByDate, pageId)
	for _, tr := range insert {
		if err = g.db.CreatePage(outNid, &tr); err != nil {
			log.Error("create the incentive_weekly_guild page failed", "err", err)
			return
		}
	}

	return
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
	data, err := g.db.GetIncentiveData(&notion.DatabaseQueryFilter{
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
