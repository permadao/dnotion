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

	totalAmount, contributors, rankOfContributor = g.statFinance("", fins)
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

// if targetToken == "", stat fin.Amount
func (g *Guild) statFinance(targetToken string, fins []dbSchema.FinData) (totalAmount float64, contributors map[string]float64, rankOfContributor []schema.Contributor) {
	// stat
	aggrContributors := map[string]float64{} // contributors id -> sum of rewards
	for _, fin := range fins {
		switch targetToken {
		case "":
			totalAmount += fin.Amount
		case fin.TargetToken:
			totalAmount += fin.TargetAmount
		default:
			continue
		}

		if fin.Contributor == "" {
			continue
		}
		if c, ok := aggrContributors[fin.Contributor]; ok {
			aggrContributors[fin.Contributor] = c + fin.TargetAmount
		} else {
			aggrContributors[fin.Contributor] = fin.TargetAmount
		}
	}

	// gen contributors & rank
	contributors = map[string]float64{}
	for k, v := range aggrContributors {
		name, ok := g.nidToName[k]
		if !ok {
			continue
		}
		contributors[name] = v

		wallet := g.nidToWallet[k]
		rankOfContributor = append(rankOfContributor, schema.Contributor{name, v, wallet})
	}

	// sort and rank
	sort.Slice(rankOfContributor, func(i, j int) bool {
		return rankOfContributor[i].Amount > rankOfContributor[j].Amount
	})

	return
}

func (g *Guild) statNewsFinance(fins []dbSchema.NewsFinData, startDate string) (totalAmount float64, aggrContributorsFor15weeks map[string]float64, aggrContributorsForAllDay map[string]float64) {
	// stat
	aggrContributorsFor15weeks = map[string]float64{} // contributors id -> sum of rewards
	aggrContributorsForAllDay = map[string]float64{}  // contributors id -> sum of rewards

	for _, fin := range fins {
		totalAmount += fin.USD

		if fin.Contributor == "" {
			continue
		}
		createDate := fin.CreatedTime
		if createDate >= startDate {
			if c, ok := aggrContributorsFor15weeks[fin.Contributor]; ok {
				aggrContributorsFor15weeks[fin.Contributor] = c + fin.USD
			} else {
				aggrContributorsFor15weeks[fin.Contributor] = fin.USD
			}
		}
		if c, ok := aggrContributorsForAllDay[fin.Contributor]; ok {
			aggrContributorsForAllDay[fin.Contributor] = c + fin.USD
		} else {
			aggrContributorsForAllDay[fin.Contributor] = fin.USD
		}
	}

	return
}

// 统计新闻组最近15周的激励
func (g *Guild) StatNewsFinance(nid, startDate string) (totalAmount float64, aggrContributorsFor15weeks map[string]float64, aggrContributorsForAllDay map[string]float64, err error) {
	fins, err := g.db.GetNewsFinances(nid, nil)
	// fmt.Println(fins[0])
	if err != nil {
		return
	}
	totalAmount, aggrContributorsFor15weeks, aggrContributorsForAllDay = g.statNewsFinance(fins, startDate)
	return
}
