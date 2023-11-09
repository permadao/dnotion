package guild

import (
	"sort"

	"github.com/dstotijn/go-notion"
	dbSchema "github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/guild/schema"
)

// StatFinance for:
// - totalAmount
// - rankOfContributor
func (g *Guild) StatFinance(targetToken, nid string) (totalAmount float64, contributors map[string]float64, rankOfContributor []schema.Contributor, err error) {
	fins, err := g.db.GetFinances(nid, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			notion.DatabaseQueryFilter{
				Property: "Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Status: &notion.StatusDatabaseQueryFilter{
						Equals: "Done",
					},
				},
			},
		},
	})
	if err != nil {
		return
	}

	totalAmount, contributors, rankOfContributor = g.statFinance(targetToken, fins)
	return
}

func (g *Guild) StatWeeklyFinance(targetToken, nid, dateStr string) (totalAmount float64, contributors map[string]float64, rankOfContributor []schema.Contributor, err error) {
	date, err := notion.ParseDateTime(dateStr)
	fins, err := g.db.GetFinances(nid, &notion.DatabaseQueryFilter{
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
						Equals: &date.Time,
					},
				},
			},
		},
	})
	if err != nil {
		return
	}

	totalAmount, contributors, rankOfContributor = g.statFinance(targetToken, fins)
	return
}

func (g *Guild) statFinance(targetToken string, fins []dbSchema.FinData) (totalAmount float64, contributors map[string]float64, rankOfContributor []schema.Contributor) {
	// stat
	contributors = map[string]float64{}
	for _, fin := range fins {
		if fin.TargetToken != targetToken {
			continue
		}
		totalAmount += fin.TargetAmount

		name, ok := g.nidToName[fin.Contributor]
		if !ok {
			continue
		}
		if c, ok := contributors[name]; ok {
			contributors[name] = c + fin.TargetAmount
		} else {
			contributors[name] = fin.TargetAmount
		}
	}

	// rank
	for k, v := range contributors {
		rankOfContributor = append(rankOfContributor, schema.Contributor{k, v})
	}
	sort.Slice(rankOfContributor, func(i, j int) bool {
		return rankOfContributor[i].Amount > rankOfContributor[j].Amount
	})

	return
}
