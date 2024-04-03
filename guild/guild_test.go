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
	//c := config.New("config_temp")
	//d := db.New(c)
	//g := New(c, d)
	//err := g.GenPromotionSettlement("", "", "2024-03-15")
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func TestGuild_GenStatRecords(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	success, _, err := g.GenIncentiveStat("45305636a546442ab5fb36fc5446b035", "2024-04-03")
	if err != nil {
		fmt.Println(err)
	}
	if success {
		fmt.Println(success)
	}
}

func TestGuild_GenTotalIncentiveStatRecords(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	err := g.GenTotalIncentiveStat("3b2c67a9b20e42208f5f3f24b8cec52c", "2024-04-03")
	if err != nil {
		fmt.Println(err)
	}
}
