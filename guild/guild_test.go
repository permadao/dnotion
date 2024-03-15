package guild

import (
	"fmt"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"testing"
)

func TestGrade(t *testing.T) {
	// c := config.New("config_temp")
	// d := db.New(c)
	// g := New(c, d)
	// start := time.Now().AddDate(0, 0, -28).Format("2006-01-02")
	// end := time.Now().Format("2006-01-02")
	// err := g.GenGrade("e8d79c55c0394cba83664f3e5737b0bd", "d8c270f68a8f44aaa6b24e17c927df2b", start, end)
	// fmt.Println(err)
}

func TestDevGrade(t *testing.T) {
	//c := config.New("config_temp")
	//d := db.New(c)
	//g := New(c, d)
	//err := g.GenDevGrade("146e1f661ed943e3a460b8cf12334b7b", "623ccfc9fb1443279decf90fb752215d", "2024-01-18", "2024-01-25")
	//fmt.Println(err)
}

func TestGuild_GenPromotionSettlement(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	err := g.GenPromotionSettlement("14debb08a4e8416e9b0de7ce46821506", "2ea3ff42b3b84d5cbc9a575d4c436878", "2024-03-15")
	if err != nil {
		fmt.Println(err)
	}
}
