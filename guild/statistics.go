package guild

import (
	"sort"

	"github.com/dstotijn/go-notion"
	dbSchema "github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/guild/schema"
)

// StatContent for:
//   - hits
//     map[int]int(key: number of hits, value: number of articles)
//   - frontPages
//     number of articles
func (g *Guild) StatContent(dateStr string) (hits map[int]int, frontPages int, err error) {
	end, err := notion.ParseDateTime(dateStr)
	if err != nil {
		return
	}
	start := end.AddDate(0, 0, -7)

	datas, err := g.db.GetContentStats(&notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			notion.DatabaseQueryFilter{
				Property: "Count Time",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						Before: &end.Time,
					},
				},
			},
			notion.DatabaseQueryFilter{
				Property: "Count Time",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						After: &start,
					},
				},
			},
		},
	})
	if err != nil {
		return
	}

	// stat
	hits = map[int]int{
		500:  0,
		1000: 0,
		5000: 0,
	}
	for _, data := range datas {
		switch true {
		case data.Hits >= 5000:
			hits[5000]++
		case data.Hits >= 1000:
			hits[1000]++
		case data.Hits >= 500:
			hits[500]++
		}

		if len(data.FrontPage) >= 3 {
			frontPages++
		}
	}
	return
}

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
	if err != nil {
		return
	}

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

func (g *Guild) StatBeforeFinance(targetToken, nid, dateStr string) (totalAmount float64, contributors map[string]float64, rankOfContributor []schema.Contributor, err error) {
	date, err := notion.ParseDateTime(dateStr)
	if err != nil {
		return
	}
	end := date.AddDate(0, 0, -3)

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
						Before: &end,
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

func (g *Guild) StatBeforeFinanceByAmount(nid, dateStr string) (totalAmount float64, contributors map[string]float64, rankOfContributor []schema.Contributor, err error) {
	end, err := notion.ParseDateTime(dateStr)
	if err != nil {
		return
	}

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
						Before: &end.Time,
					},
				},
			},
		},
	})
	if err != nil {
		return
	}

	totalAmount, contributors, rankOfContributor = g.statFinanceByAmount(fins)
	return
}

func (g *Guild) StatBetweenFinance(targetToken, nid, startDate, endDate string) (totalAmount float64, contributors map[string]float64, rankOfContributor []schema.Contributor, err error) {
	start, err := notion.ParseDateTime(startDate)
	if err != nil {
		return
	}
	end, err := notion.ParseDateTime(endDate)
	if err != nil {
		return
	}

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
						OnOrAfter: &start.Time,
					},
				},
			},
			notion.DatabaseQueryFilter{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						OnOrBefore: &end.Time,
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

func (g *Guild) statFinanceByAmount(fins []dbSchema.FinData) (totalAmount float64, contributors map[string]float64, rankOfContributor []schema.Contributor) {
	// stat
	contributors = map[string]float64{}
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
