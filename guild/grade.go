package guild

import (
	dbSchema "github.com/permadao/dnotion/db/schema"
	"github.com/permadao/dnotion/guild/schema"
)

func GRankToGrade(rankOfContributor []schema.Contributor) (translators []dbSchema.Translator) {
	n := len(rankOfContributor)
	translators = make([]dbSchema.Translator, n)

	for i, r := range rankOfContributor {
		per := (float64(i) + 1) / float64(n)
		translators[i] = dbSchema.Translator{
			Contributor: r.Name,
			Level:       GTranslatorRank(per),
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
