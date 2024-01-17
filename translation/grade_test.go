package translation

import (
	"fmt"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/guild/schema"
	"testing"
)

func TestGrade(t *testing.T) {
	myconfig := config.New("config_temp")
	d := db.New(myconfig)
	tr := New(d)
	if err := tr.Grade(); err != nil {
		fmt.Println(err)
	}
}

func TestTranslator_RankToGrade(t *testing.T) {
	rankOfContributor := []schema.Contributor{
		{
			"001",
			100.00,
		},
		{
			"002",
			80.00,
		},
		{
			"003",
			60.00,
		},
		{
			"004",
			50.00,
		},
		{
			"005",
			40.00,
		},
	}

	rankOfContributor2 := []schema.Contributor{
		{
			"006",
			30.00,
		},
		{
			"007",
			30.00,
		},
		{
			"008",
			30.00,
		},
		{
			"009",
			30.00,
		},
		{
			"010",
			30.00,
		},
	}
	tr := Translator{}
	translators, _ := tr.RankToGrade(rankOfContributor, 0)
	fmt.Println(translators)
	rankOfContributor = append(rankOfContributor, rankOfContributor2...)
	translators, _ = tr.RankToGrade(rankOfContributor, 0)
	fmt.Println(translators)
}
