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
	nidToName   map[string]string //  contributor data nid -> contributors name
	nidToUserID map[string]string
}

func New(conf *config.Config, db *db.DB) *Guild {
	g := &Guild{
		db: db,

		nidToName:   map[string]string{},
		nidToUserID: map[string]string{},
	}

	g.initContributors()
	return g
}

func (g *Guild) initContributors() {
	contributors, err := g.db.GetContributors(nil)
	if err != nil {
		panic(err)
	}

	for _, c := range contributors {
		if c.NotionName != "" {
			g.nidToName[c.NID] = c.NotionName
		}
		if c.NotionID != "" {
			g.nidToUserID[c.NID] = c.NotionID
		}
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
			fmt.Println(err)
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
		msg := fmt.Sprintf("get last id from page failed:%s, finance nid: %s", err.Error(), gradeNid)
		fmt.Printf(msg)
		return
	}
	grades := GRankToGrade(rankOfContributor, id)
	fmt.Printf("%+v\n", grades)

	for _, tr := range grades {
		if err = g.db.CreatePage(gradeNid, &tr); err != nil {
			msg := "create page failed:" + err.Error()
			fmt.Println(msg)
			return
		}
	}

	return
}
