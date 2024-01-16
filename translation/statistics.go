package translation

import (
	"fmt"
	"github.com/dstotijn/go-notion"
	dbSchema "github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/guild/schema"
	log "github.com/sirupsen/logrus"
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
	n := len(rankOfContributor)
	gradedNum := 0
	ts := t.GetTierSlice()
	for i := 0; i < len(ts) && gradedNum < n; i++ {
		l, r := ts[i].interval[0]*n/10, ts[i].interval[1]*n/10
		ts[i].val = [2]int{l, r}
		gradedNum += r - l
	}
	trSlice := []dbSchema.Translator{}
	PageID, err := t.db.GetLastID(TRDB)
	if err != nil {
		msg := fmt.Sprintf("get last id from page failed:%s, finance nid: %s", err.Error(), TRDB)
		log.Error(msg)
		return
	}
	date := time.Now().Format("20060102")
	for _, tier := range ts {
		start, end := tier.val[0], tier.val[1]
		for i := start; i < end; i++ {
			PageID++
			tr := dbSchema.Translator{
				ID:          fmt.Sprintf("%d", PageID),
				Contributor: rankOfContributor[i].Name,
				Level:       tier.level.name,
				Title:       "普通译员",
				Date:        date,
			}
			trSlice = append(trSlice, tr)
		}
	}

	//Update
	for _, tr := range trSlice {
		if err = t.db.CreatePage(TRDB, &tr); err != nil {
			msg := "create page failed:" + err.Error()
			log.Error(msg)
			return
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
