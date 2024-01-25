package guild

import (
	"fmt"
	dbSchema "github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/guild/schema"
	"time"
)

func GRankToGrade(rankOfContributor []schema.Contributor, PageID int) (translators []dbSchema.Translator) {
	n := len(rankOfContributor)
	translators = make([]dbSchema.Translator, n)

	date := time.Now().Format("2006-01-02")

	for i, r := range rankOfContributor {
		per := (float64(i) + 1) / float64(n)
		PageID++
		translators[i] = dbSchema.Translator{
			ID:          fmt.Sprintf("%d", PageID),
			Contributor: r.Name,
			Level:       GTranslatorRank(per),
			Title:       "普通译员",
			Date:        date,
		}
	}
	return
}

func GRankToGradeForDev(rankOfContributor []schema.Contributor, developers []dbSchema.Developer, PageID int) (insert []dbSchema.Developer) {
	n := len(developers)
	date := time.Now().Format("2006-01-02")
	devMap := make(map[string]*dbSchema.Developer, n)
	for i := 0; i < len(developers); i++ {
		devMap[developers[i].Contributor] = &developers[i]
	}
	for _, r := range rankOfContributor {
		if d, ok := devMap[r.Name]; ok {
			if d.Income < 500.0 && r.Amount >= 500.0 {
				d.Date = date
				d.Level = GDeveloperLevel(r.Amount)
			}
			d.Income = r.Amount
		} else {
			PageID++
			insert = append(insert, dbSchema.Developer{
				ID:          fmt.Sprintf("%d", PageID),
				Contributor: r.Name,
				Level:       GDeveloperLevel(r.Amount),
				Income:      r.Amount,
				Date:        date,
			})
		}
	}
	return
}

func GTranslatorRank(per float64) (res string) {
	switch true {
	case per <= 0.1:
		res = "Supreme-至尊"
	case per <= 0.2:
		res = "Glory-荣耀"
	case per <= 0.4:
		res = "Diamond-钻石"
	case per <= 0.7:
		res = "Gold-黄金"
	default:
		res = "Silver-白银"
	}
	return
}

func GDeveloperLevel(income float64) (res string) {
	switch true {
	case income >= 500.0:
		res = "高级开发者"
	default:
		res = "普通开发者"
	}
	return
}
