package guild

import (
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

func (g *Guild) GenDevGrade(guidNid, gradeNid, endDate string) (err error) {
	_, _, rankOfContributor, err := g.StatBeforeFinanceByAmount(guidNid, endDate)
	if err != nil {
		return
	}
	id, err := g.db.GetLastID(gradeNid)
	if err != nil {
		log.Error("get last id from page failed", "finance nid", gradeNid, "err", err)
		return
	}
	developers, err := g.db.GetDeveloper(gradeNid, nil)
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
