package guild

import (
	"fmt"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"testing"
	"time"
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
	success, paydateMap, _ := g.GenIncentiveStat("4c19704d927f4d52b2f030ebd1648ef3", "2023-03-04")
	if success {
		err := g.GenTotalIncentiveStat("46dc65c61cd94e43bdaeee7ed22f15d2", paydateMap)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestGuild_GenTotalIncentiveStatRecords(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	startDate, _ := time.Parse("2006-01-02", "2023-03-11")
	endDate, _ := time.Parse("2006-01-02", "2024-04-06")
	for !startDate.After(endDate) {
		startDateStr := startDate.Format("2006-01-02")
		success, paydateMap, err := g.GenIncentiveStat("4c19704d927f4d52b2f030ebd1648ef3", startDateStr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(paydateMap)
		if success {
			err := g.GenTotalIncentiveStat("46dc65c61cd94e43bdaeee7ed22f15d2", paydateMap)
			if err != nil {
				return
			}
		}
		fmt.Println("完成" + startDateStr)
		startDate = startDate.AddDate(0, 0, 7)
		fmt.Println("下一时间", startDate)
	}
}
