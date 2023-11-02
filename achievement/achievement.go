package achievement

import (
	"fmt"

	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/db/schema"
)

type Achievement struct {
	db *db.DB

	// contributors
	nidToName         map[string]string //  contributor data nid -> contributors name
	finNidToGuildName map[string]string // finNid -> Guild Name in achievements
}

func New(conf *config.Config, db *db.DB) *Achievement {
	a := &Achievement{
		db: db,

		nidToName: map[string]string{},
		finNidToGuildName: map[string]string{
			"328f2bfbfdbe4f9581af37f393893e36": "内容公会 - 策划",    // content
			"e8d79c55c0394cba83664f3e5737b0bd": "内容公会 - 翻译",    // translation
			"a815dcd96395424a93d9854b4418ab03": "内容公会 - 投稿",    // submission
			"990db3313e42412b8c6ab07e399a2635": "品宣公会",         // promotion
			"f2160eae42e9483882f01d3daa7090fa": "活动公会",         // activity
			"caac7a1aefcc4ed0b02b8adbc106f021": "管理公会",         // admin
			"146e1f661ed943e3a460b8cf12334b7b": "开发公会",         // dev
			"a9ce0c5902b14e4891ed0fb6333a9e92": "PSPC Market",  // psp market
			"27555aec8d734b6889ae1836d7a67b4a": "PSPC Product", // psp prod
		},
	}

	a.initContributors()
	return a
}

func (a *Achievement) initContributors() {
	contributors, err := a.db.GetAllContributors()
	if err != nil {
		panic(err)
	}

	for _, c := range contributors {
		if c.NotionName != "" {
			a.nidToName[c.NID] = c.NotionName
		}
	}
}

func (a *Achievement) GenAchievements(date string) {
	lastID, _ := a.db.GetLastIDFromDB(a.db.AchievementDB)

	for finNid, guildName := range a.finNidToGuildName {
		totalAmount, rankOfContributor, _ := a.StatWeeklyFinances(finNid, date)

		tags := []string{}
		if a := AGuildActiviy(len(rankOfContributor)); a != "" {
			tags = append(tags, a)
		}
		if len(rankOfContributor) > 0 && totalAmount != 0 {
			if a := AFairDistribution(rankOfContributor[0].Amount / totalAmount); a != "" {
				tags = append(tags, a)
			}
		}

		lastID++
		ach := schema.AchievementData{
			ID:    fmt.Sprintf("%d", lastID),
			Guild: guildName,
			Tags:  tags,
			Date:  date,
		}
		a.db.CreatePage(a.db.AchievementDB, &ach)
	}
}

// 成就：公会热度
func AGuildActiviy(num int) (res string) {
	switch true {
	case num > 15:
		res = "热度/门庭若市"
	case num > 8:
		res = "热度/小有热闹"
	case num > 3:
		res = "热度/结伴而行"
	default:
		res = "热度/冷冷清清"
	}
	return
}

// 成就：激励分配分散度
func AFairDistribution(per float64) (res string) {
	switch true {
	case per > 0.8:
		res = "分配/超集中"
	case per > 0.5:
		res = "分配/个别集中"
	default:
		res = ""
	}
	return
}
