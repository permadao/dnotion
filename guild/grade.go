package guild

import (
	"fmt"
	dbSchema "github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/guild/schema"
	"math"
	"strings"
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
	for contributor, points := range contributors {
		name := notionidToName[contributor]
		promotionSettlement = append(promotionSettlement, dbSchema.PromotionSettlement{
			Contributor:         name,
			ContributorNotionID: contributor,
			TotalScore:          totalPoints,
			PersonalScore:       points,
			Rewards:             points,
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

func GetGuildFinMap() (guildFinMap map[string]string) {
	guildFinMap = make(map[string]string)
	guildFinMap["Activity"] = "f2160eae42e9483882f01d3daa7090fa"
	guildFinMap["Promotion"] = "990db3313e42412b8c6ab07e399a2635"
	guildFinMap["Admin"] = "caac7a1aefcc4ed0b02b8adbc106f021"
	guildFinMap["Content"] = "328f2bfbfdbe4f9581af37f393893e36"
	guildFinMap["Translation"] = "e8d79c55c0394cba83664f3e5737b0bd"
	guildFinMap["Submission"] = "a815dcd96395424a93d9854b4418ab03"
	guildFinMap["Develop"] = "146e1f661ed943e3a460b8cf12334b7b"
	guildFinMap["PSP Market"] = "a9ce0c5902b14e4891ed0fb6333a9e92"
	guildFinMap["PSP Product"] = "27555aec8d734b6889ae1836d7a67b4a"
	guildFinMap["Meeting"] = "0c8f5483e1344d919e5e5a49d6d8dabb"
	return
}

func GenStatRecords(completeStatResults map[string]*schema.StatResult, weekStatResults map[string]*schema.StatResult, hisRecords map[string]dbSchema.CIncentiveGuild, guild string, acDate string, paymentDate string, pageID *float64, g *Guild) (insertRecords []dbSchema.CIncentiveGuild, updateRecords []dbSchema.CIncentiveGuild) {
	for token, statResult := range completeStatResults {
		if _, ok := weekStatResults[token]; !ok {
			continue
		}
		completeContributors := statResult.Contributors
		weekContributors := weekStatResults[token].Contributors
		for contributor, amount := range completeContributors {
			notionName := contributor
			if _, ok := g.nidToName[contributor]; ok {
				notionName = g.nidToName[contributor]
			} else if _, ok := g.notionidToName[contributor]; ok {
				notionName = g.notionidToName[contributor]
			}
			if _, ok := weekContributors[contributor]; !ok {
				continue
			}
			incentiveRecord := dbSchema.CIncentiveGuild{
				AccountingDate:  acDate,
				Guild:           guild,
				NotionName:      notionName,
				NotionID:        contributor,
				BuddyNotion:     g.notionidToName[g.nidToContributor[contributor].BuddyNotion],
				TotalIncentive:  amount,
				WeeklyIncentive: weekContributors[contributor],
				PaymentDate:     paymentDate,
				OnboardDate:     g.nidToContributor[contributor].CreatedTime,
				Token:           token,
			}
			//所有币别全时期总金额
			var totalIncentive float64
			for _, result := range completeStatResults {
				totalIncentive += result.Contributors[contributor]
			}
			var weekIncentive float64
			for _, result := range weekStatResults {
				weekIncentive += result.Contributors[contributor]
			}
			if totalIncentive == weekIncentive {
				incentiveRecord.FirstContribution = "Yes"
			}

			if record, ok := hisRecords[GetKey(incentiveRecord)]; ok {
				record.TotalIncentive = incentiveRecord.TotalIncentive
				record.WeeklyIncentive = incentiveRecord.WeeklyIncentive
				record.FirstContribution = incentiveRecord.FirstContribution
				updateRecords = append(updateRecords, record)
			} else if record, ok := hisRecords[GetKeyNoToken(incentiveRecord)]; ok {
				record.TotalIncentive = incentiveRecord.TotalIncentive
				record.WeeklyIncentive = incentiveRecord.WeeklyIncentive
				record.Token = incentiveRecord.Token
				record.FirstContribution = incentiveRecord.FirstContribution
				delete(hisRecords, GetKeyNoToken(incentiveRecord))
				updateRecords = append(updateRecords, record)
			} else {
				*pageID++
				incentiveRecord.ID = *pageID
				insertRecords = append(insertRecords, incentiveRecord)
			}
		}
	}
	return
}

func CalTotalIncentive(data []dbSchema.CIncentiveGuild, hisData map[string]schema.ResultSepToken, records map[string]dbSchema.CIncentive, pageID *float64) (insertRecords []dbSchema.CIncentive, updateRecords []dbSchema.CIncentive) {
	contributorMap := map[string]*dbSchema.CIncentive{}
	//不区分币别的激励
	allIncentive := map[string]float64{}
	weekIncentive := map[string]float64{}
	for _, incentive := range data {
		token := incentive.Token
		key := strings.Join([]string{incentive.NotionID, incentive.PaymentDate, incentive.Token}, "-")
		if _, ok := contributorMap[key]; ok {
			totalIncentive := contributorMap[key]
			(*totalIncentive).TotalIncentive += incentive.WeeklyIncentive
			(*totalIncentive).WeeklyIncentive += incentive.WeeklyIncentive
			allIncentive[incentive.NotionID] += incentive.WeeklyIncentive
			weekIncentive[incentive.NotionID] += incentive.WeeklyIncentive
		} else {
			totalIncentive := 0.0
			if ti, ok := hisData[incentive.PaymentDate]; ok {
				if rst, ok := ti[incentive.Token]; ok {
					totalIncentive = rst[incentive.NotionID]
				}
			}
			contributorMap[key] = &dbSchema.CIncentive{
				ID:              *pageID,
				AccountingDate:  incentive.AccountingDate,
				NotionID:        incentive.NotionID,
				NotionName:      incentive.NotionName,
				BuddyNotion:     incentive.BuddyNotion,
				TotalIncentive:  totalIncentive + incentive.WeeklyIncentive,
				WeeklyIncentive: incentive.WeeklyIncentive,
				PaymentDate:     incentive.PaymentDate,
				OnboardDate:     incentive.OnboardDate,
				Token:           token,
			}
			*pageID++
			allIncentive[incentive.NotionID] += contributorMap[key].TotalIncentive
			weekIncentive[incentive.NotionID] += contributorMap[key].WeeklyIncentive
		}
	}

	for _, totalIncentiveRecord := range contributorMap {
		if allIncentive[totalIncentiveRecord.NotionID] == weekIncentive[totalIncentiveRecord.NotionID] {
			(*totalIncentiveRecord).FirstContribution = "Yes"
		}
		(*totalIncentiveRecord).Medal = GDMedal(totalIncentiveRecord.TotalIncentive-totalIncentiveRecord.WeeklyIncentive, totalIncentiveRecord.TotalIncentive)
		//upsert
		if record, ok := records[GetTIKey(*totalIncentiveRecord)]; ok {
			record.TotalIncentive = totalIncentiveRecord.TotalIncentive
			record.WeeklyIncentive = totalIncentiveRecord.WeeklyIncentive
			record.FirstContribution = totalIncentiveRecord.FirstContribution
			record.Medal = totalIncentiveRecord.Medal
			updateRecords = append(updateRecords, record)
		} else if record, ok := records[GetTIKeyNoToken(*totalIncentiveRecord)]; ok {
			record.TotalIncentive = totalIncentiveRecord.TotalIncentive
			record.WeeklyIncentive = totalIncentiveRecord.WeeklyIncentive
			record.FirstContribution = totalIncentiveRecord.FirstContribution
			record.Medal = totalIncentiveRecord.Medal
			record.Token = totalIncentiveRecord.Token
			delete(records, GetTIKeyNoToken(*totalIncentiveRecord))
			updateRecords = append(updateRecords, record)
		} else {
			insertRecords = append(insertRecords, *totalIncentiveRecord)
		}
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

func GDMedal(hisIncentive, curIncentive float64) (medal string) {
	switch true {
	case hisIncentive < 200 && curIncentive >= 200:
		{
			medal = "达成200"
		}
	case hisIncentive < 50 && curIncentive >= 50:
		{
			medal = "达成50"
		}
	default:
		{
			medal = ""
		}
	}
	return
}

func GetKey(incentive dbSchema.CIncentiveGuild) string {
	return strings.Join([]string{incentive.Guild, incentive.NotionID, incentive.PaymentDate, incentive.Token}, "-")
}

func GetKeyNoToken(incentive dbSchema.CIncentiveGuild) string {
	return strings.Join([]string{incentive.Guild, incentive.NotionID, incentive.PaymentDate, ""}, "-")
}

func GetTIKey(incentive dbSchema.CIncentive) string {
	return strings.Join([]string{incentive.NotionID, incentive.PaymentDate, incentive.Token}, "-")
}

func GetTIKeyNoToken(incentive dbSchema.CIncentive) string {
	return strings.Join([]string{incentive.NotionID, incentive.PaymentDate, ""}, "-")
}
