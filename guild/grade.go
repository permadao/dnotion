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

func CalculatePromotionRewards(promotionPoints []dbSchema.PromotionPoints, date string) (promotionSettlement []dbSchema.PromotionSettlement) {
	totalPoints := 0.0
	contributors := map[string]float64{}
	promotionNum := map[string]struct{}{}
	for _, p := range promotionPoints {
		totalPoints += p.BasePoints
		contributors[p.Contributor] += p.BasePoints
		if _, ok := promotionNum[p.Task]; !ok {
			promotionNum[p.Task] = struct{}{}
		}
	}
	pool := CalculateRewardPool(float64(len(contributors)), float64(len(promotionNum)))
	for contributorID, points := range contributors {
		promotionSettlement = append(promotionSettlement, dbSchema.PromotionSettlement{
			Contributor:   contributorID,
			TotalScore:    totalPoints,
			PersonalScore: points,
			Rewards:       pool * points / totalPoints,
			Date:          date,
		})
	}
	return
}

// CalculateRewardPool Calculation method of total weekly funds
// entry : total number of weekly contributors
// promotions: total number of weekly tasks
// constant : constant
// entryW : weight for entry
// pw : weight for promotion
func CalculateRewardPool(entry, promotions float64) float64 {
	constant := 50.0
	entryW := 0.8
	pW := 0.7
	return constant * math.Pow(entry, entryW) / (1 + math.Pow(entry/promotions, pW))
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
