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
