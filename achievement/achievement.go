package achievement

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/achievement/schema"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
)

type Achievement struct {
	db *db.DB

	// contributors
	nidToName map[string]string //  contributor data nid -> contributors name
}

func New(conf *config.Config, db *db.DB) *Achievement {
	a := &Achievement{
		db: db,

		nidToName: map[string]string{},
	}
	a.initContributors()
	return a
}

func (a *Achievement) initContributors() {
	contributors, err := a.db.GetAllContributors()
	if err != nil {
		panic(err)
	}

	for _, c := range contributors {
		if c.NotionName != "" {
			a.nidToName[c.NID] = c.NotionName
		}
	}
}

// Statistics for:
// - totalAmount
// - rankOfContributor
func (a *Achievement) StatWeeklyFinances(nid string, yyyymmdd string) (
	totalAmount float64, rankOfContributor []schema.Contributor,
	err error) {

	date, err := YYYYMMDDtoTime(yyyymmdd)
	if err != nil {
		return
	}

	fins, err := a.db.GetFinancesByNID(nid, &notion.DatabaseQueryFilter{
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
						Equals: date,
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

		name, ok := a.nidToName[fin.Contributor]
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

func YYYYMMDDtoTime(yyyymmdd string) (t *time.Time, err error) {
	if len(yyyymmdd) != 8 {
		return nil, errors.New("Invalid YYYYMMDD")
	}

	year, err := strconv.Atoi(yyyymmdd[:4])
	if err != nil {
		return
	}
	month, err := strconv.Atoi(yyyymmdd[4:6])
	if err != nil {
		return
	}
	day, err := strconv.Atoi(yyyymmdd[6:])
	if err != nil {
		return
	}

	tim := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	t = &tim
	return
}
