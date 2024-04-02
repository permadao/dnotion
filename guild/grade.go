package guild

import (
	"fmt"
	dbSchema "github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/guild/schema"
	"math"
	"time"
)

func GRankToGrade(rankOfContributor []schema.Contributor, PageID int) (translators []dbSchema.Translator) {
	n := len(rankOfContributor)
	translators = make([]dbSchema.Translator, n)

	date := time.Now().Format("2006-01-02")

	for i, r := range rankOfContributor {
		per := (float64(i) + 1) / float64(n)
		PageID++
		translators[i] = dbSchema.Translator{
			ID:          fmt.Sprintf("%d", PageID),
			Contributor: r.Name,
			Level:       GTranslatorRank(per),
			Title:       "普通译员",
			Date:        date,
		}
	}
	return
}

func GRankToGradeForDev(rankOfContributor []schema.Contributor, developers []dbSchema.Developer, PageID int, endDate string) (insert []dbSchema.Developer) {
	n := len(developers)
	coreDev := make(map[string]struct{}, n)
	for i := 0; i < len(developers); i++ {
		if developers[i].Level == "核心开发者" {
			coreDev[developers[i].Contributor] = struct{}{}
		}
	}
	for _, r := range rankOfContributor {
		PageID++
		level := GDeveloperLevel(r.Amount)
		if _, ok := coreDev[r.Name]; ok {
			level = "核心开发者"
		}
		insert = append(insert, dbSchema.Developer{
			ID:          fmt.Sprintf("%d", PageID),
			Contributor: r.Name,
			Level:       level,
			Income:      r.Amount,
			Date:        endDate,
		})
	}
	return
}

func GTranslatorRank(per float64) (res string) {
	switch true {
	case per <= 0.1:
		res = "Supreme-至尊"
	case per <= 0.2:
		res = "Glory-荣耀"
	case per <= 0.4:
		res = "Diamond-钻石"
	case per <= 0.7:
		res = "Gold-黄金"
	default:
		res = "Silver-白银"
	}
	return
}

func GDeveloperLevel(income float64) (res string) {
	switch true {
	case income >= 500.0:
		res = "高级开发者"
	default:
		res = "普通开发者"
	}
	return
}
func GRankToGradeForNews(aggrContributorsFor15weeks map[string]float64, aggrContributorsForAllDay map[string]float64, news []dbSchema.News, PageID int, endDate string) (insert []dbSchema.News) {
	n := len(news)
	// fmt.Println(n)
	for i := 0; i < n; i++ {
		contributor := news[i].Executor
		fmt.Println(news[i].ID)
		// curRewards := news[i].CumulativeWorkload
		// fmt.Println(curRewards)
		// curRewardsOf15weeks := news[i].WorkloadOf15Weeds
		// fmt.Println(curRewardsOf15weeks)

		reward, _ := aggrContributorsForAllDay[contributor]

		rewardOf15weeks, _ := aggrContributorsFor15weeks[contributor]

		PageID++
		rank, rankCode := GDNewsLevel(rewardOf15weeks)
		insert = append(insert, dbSchema.News{
			NID: news[i].NID,
			ID:  news[i].ID,
			// Executor: news[i].Executor,
			// ExecutorWorkload: news[i].ExecutorWorkload,
			Rank:               rank,
			WorkloadOf15Weeds:  rewardOf15weeks,
			RankCode:           rankCode,
			CumulativeWorkload: reward,
			Date:               endDate,
		})
	}
	return
}

func CalculatePromotionRewards(promotionPoints []dbSchema.PromotionPoints, notionidToName map[string]string, notionidToID map[string]*float64, date string) (promotionSettlement []dbSchema.PromotionSettlement) {
	totalPoints := 0.0
	contributors := map[string]float64{}
	for _, p := range promotionPoints {
		totalPoints += p.BasePoints
		contributors[p.Contributor] += p.BasePoints
	}
	pool := CalculateRewardPool(float64(len(contributors)), float64(len(promotionPoints)))
	for contributor, points := range contributors {
		name := notionidToName[contributor]
		promotionSettlement = append(promotionSettlement, dbSchema.PromotionSettlement{
			Contributor:         name,
			ContributorNotionID: contributor,
			TotalScore:          totalPoints,
			PersonalScore:       points,
			Rewards:             pool * points / totalPoints,
			Date:                date,
		})
	}
	return
}

// CalculateRewardPool Calculation method of total weekly funds
// entry : total number of weekly contributors
// promotions: total number of weekly promotion record
// constant : constant
// entryW : weight for entry
// pw : weight for promotion
func CalculateRewardPool(entry, recordNum float64) float64 {
	constant := 50.0
	entryW := 0.8
	pW := 0.7
	return constant * math.Pow(entry, entryW) / (1 + math.Pow(entry/recordNum, pW))
}

func GetGuildFinMap() (gfm map[string]string) {
	gfm["Content"] = "328f2bfbfdbe4f9581af37f393893e36"
	gfm["Translation"] = "e8d79c55c0394cba83664f3e5737b0bd"
	gfm["Submission"] = "a815dcd96395424a93d9854b4418ab03"
	gfm["Promotion"] = "990db3313e42412b8c6ab07e399a2635"
	gfm["Activity"] = "f2160eae42e9483882f01d3daa7090fa"
	gfm["admin"] = "caac7a1aefcc4ed0b02b8adbc106f021"
	gfm["dev"] = "146e1f661ed943e3a460b8cf12334b7b"
	gfm["Psp Market"] = "a9ce0c5902b14e4891ed0fb6333a9e92"
	gfm["Psp Product"] = "27555aec8d734b6889ae1836d7a67b4a"
	gfm["MapDAO"] = "d52d6f8994504b89bad2b9dd8ff5d586"
	return
}

func GenStatRecords(contributorsAllTime map[string]float64, contributorsThisWeek map[string]float64, guild string, acDate string, paymentDate string, pageID int) (incentiveRecords []dbSchema.Incentive) {
	for contributor, amount := range contributorsAllTime {
		incentiveRecord := dbSchema.Incentive{
			AccountingDate: acDate,
			Guild:          guild,
			NotionName:     contributor,
			BuddyNotion:    "",
			TotalIncentive: amount,
			PaymentDate:    paymentDate,
			OnboardDate:    "",
		}
		if weekAmount, ok := contributorsThisWeek[contributor]; ok {
			incentiveRecord.WeeklyIncentive = weekAmount
		}
		if incentiveRecord.TotalIncentive == incentiveRecord.WeeklyIncentive {
			incentiveRecord.FirstContribution = "Yes"
		}
		pageID++
		incentiveRecord.ID = fmt.Sprintf("%d", pageID)
		incentiveRecords = append(incentiveRecords, incentiveRecord)
	}
	return
}

func CalTotalIncentive(data []dbSchema.Incentive, pageID int) (totalIncentiveRecords []dbSchema.TotalIncentive) {
	contributorMap := map[string]*dbSchema.TotalIncentive{}
	for _, incentive := range data {
		if _, ok := contributorMap[incentive.NotionName]; ok {
			totalIncentive := contributorMap[incentive.NotionName]
			(*totalIncentive).TotalIncentive += incentive.TotalIncentive
			(*totalIncentive).WeeklyIncentive += incentive.WeeklyIncentive
		} else {
			pageID++
			contributorMap[incentive.NotionName] = &dbSchema.TotalIncentive{
				ID:                fmt.Sprintf("%d", pageID),
				AccountingDate:    incentive.AccountingDate,
				NotionName:        incentive.NotionName,
				BuddyNotion:       incentive.BuddyNotion,
				TotalIncentive:    incentive.TotalIncentive,
				WeeklyIncentive:   incentive.WeeklyIncentive,
				PaymentDate:       incentive.PaymentDate,
				OnboardDate:       incentive.OnboardDate,
				FirstContribution: incentive.FirstContribution,
			}
		}
	}
	for _, totalIncentiveRecord := range contributorMap {
		totalIncentiveRecords = append(totalIncentiveRecords, *totalIncentiveRecord)
	}
	return
}

func GDNewsLevel(income float64) (res string, code string) {
	switch true {
	case income >= 1600.0:
		{
			res = "氪石作者"
			code = "Kryptonite"
		}
	case income >= 800.0:
		{
			res = "钻石作者"
			code = "Diamond"
		}
	case income >= 400.0:
		{
			res = "黄金作者"
			code = "Gold"
		}
	case income >= 200.0:
		{
			res = "白银作者"
			code = "Silver"
		}
	default:
		{
			res = "普通作者"
			code = "normal"
		}
	}
	return
}
