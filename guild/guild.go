package guild

import (
	"fmt"

	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	dbSchema "github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/guild/schema"
)

type Guild struct {
	db *db.DB

	// contributors
	nidToName map[string]string //  contributor data nid -> contributors name
}

func New(conf *config.Config, db *db.DB) *Guild {
	g := &Guild{
		db: db,

		nidToName: map[string]string{},
	}

	g.initContributors()
	return g
}

func (g *Guild) initContributors() {
	contributors, err := g.db.GetAllContributors()
	if err != nil {
		panic(err)
	}

	for _, c := range contributors {
		if c.NotionName != "" {
			g.nidToName[c.NID] = c.NotionName
		}
	}
}

func (g *Guild) GenGuilds(date string) {
	for guildName, info := range schema.Guilds {
		totalAmount, rankOfContributor, _ := g.StatWeeklyFinances(info.FinNID, date)

		tags := []string{}
		if a := AGuildActiviy(len(rankOfContributor)); a != "" {
			tags = append(tags, a)
		}
		if len(rankOfContributor) > 0 && totalAmount != 0 {
			if a := AFairDistribution(rankOfContributor[0].Amount / totalAmount); a != "" {
				tags = append(tags, a)
			}
		}

		guild := dbSchema.GuildData{
			Name: guildName,
			Link: info.NID,
			Tags: tags,
			Info: info.Info,
			Date: date,
			Rank: info.Rank,
		}
		if err := g.db.CreatePage(g.db.GuildDB, &guild); err != nil {
			fmt.Println(err)
		}
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
