package guild

import (
	"sort"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/guild/schema"
)

// StatWeeklyFinance Statistics for:
// - totalAmount
// - rankOfContributor
func (g *Guild) StatWeeklyFinances(nid string, dateStr string) (
	totalAmount float64, rankOfContributor []schema.Contributor,
	err error) {

	date, err := notion.ParseDateTime(dateStr)

	fins, err := g.db.GetFinancesByNID(nid, &notion.DatabaseQueryFilter{
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

	// stat
	contributors := map[string]float64{}
	for _, fin := range fins {
		totalAmount += fin.Amount

		name, ok := g.nidToName[fin.Contributor]
		if !ok {
			continue
		}
		if c, ok := contributors[name]; ok {
			contributors[name] = c + fin.Amount
		} else {
			contributors[name] = fin.Amount
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
