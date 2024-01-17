package translation

import (
	"fmt"
	"github.com/dstotijn/go-notion"
	dbSchema "github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/guild/schema"
	log "github.com/sirupsen/logrus"
	"math"
	"sort"
	"time"
)

const TRDB = "e8d79c55c0394cba83664f3e5737b0bd"

func (t *Translator) Grade() (err error) {
	//Summary Finance of last 4 weeks
	rankOfContributor, err := t.SummaryFinance()
	if err != nil {
		return
	}

	PageID, err := t.db.GetLastID(TRDB)
	if err != nil {
		msg := fmt.Sprintf("get last id from page failed:%s, finance nid: %s", err.Error(), TRDB)
		log.Error(msg)
		return
	}

	translators, err := t.RankToGrade(rankOfContributor, PageID)
	if err != nil {
		return
	}

	//Update
	for _, tr := range translators {
		if err = t.db.CreatePage(TRDB, &tr); err != nil {
			msg := "create page failed:" + err.Error()
			log.Error(msg)
			return
		}
	}
	return
}

func (t *Translator) RankToGrade(rankOfContributor []schema.Contributor, PageID int) (translators []dbSchema.Translator, err error) {
	n := len(rankOfContributor)
	gradedNum := 0
	ts := t.GetTierSlice()
	translators = make([]dbSchema.Translator, n)

	date := time.Now().Format("20060102")
	l, r := 0, 0
	for i := 0; i < len(ts) && gradedNum < n; i++ {
		gap := int(math.Ceil(float64((ts[i].Interval[1]-ts[i].Interval[0])*n) / 10.0))
		l, r = r, r+gap
		gradedNum += r - l
		for j := l; j < r; j++ {
			PageID++
			translators[j] = dbSchema.Translator{
				ID:          fmt.Sprintf("%d", PageID),
				Contributor: rankOfContributor[j].Name,
				Level:       ts[i].Level.Name,
				Title:       "普通译员",
				Date:        date,
			}
		}
	}
	return
}

func (t *Translator) SummaryFinance() (rankOfContributor []schema.Contributor, err error) {
	start := notion.NewDateTime(time.Now().AddDate(0, 0, -28), false)
	end := notion.NewDateTime(time.Now(), false)
	fins, err := t.db.GetFinances(TRDB, &notion.DatabaseQueryFilter{
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
						OnOrBefore: &end.Time,
						OnOrAfter:  &start.Time,
					},
				},
			},
		},
	})
	if err != nil {
		return
	}
	rankOfContributor = t.statFinance("AR", fins)
	return
}

func (t *Translator) statFinance(targetToken string, fins []dbSchema.FinData) (rankOfContributor []schema.Contributor) {
	// stat
	contributors := map[string]float64{}
	for _, fin := range fins {
		if fin.TargetToken != targetToken {
			continue
		}

		_, ok := t.nidToName[fin.Contributor]
		if !ok {
			continue
		}
		if c, ok := contributors[fin.Contributor]; ok {
			contributors[fin.Contributor] = c + fin.TargetAmount
		} else {
			contributors[fin.Contributor] = fin.TargetAmount
		}
	}

	// rank
	for k, v := range contributors {
		rankOfContributor = append(rankOfContributor, schema.Contributor{Name: k, Amount: v})
	}
	sort.Slice(rankOfContributor, func(i, j int) bool {
		return rankOfContributor[i].Amount > rankOfContributor[j].Amount
	})

	return
}
