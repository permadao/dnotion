package guild

import (
	"fmt"
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"sort"
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
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	now := time.Now()
	startDate := now.AddDate(0, 0, -15*7)
	startDateOfNews := startDate.Format("2006-01-02")
	end := time.Now().Format("2006-01-02")
	err := g.GenNewsGrade("ad2cf585b08843fea7cf40a682bc4529", "d5f9fc70910b45d4ab8811f37716637d", startDateOfNews, end)
	fmt.Println(err)
}

func TestGuild_GenPromotionSettlement(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	err := g.GenPromotionSettlement("14debb08a4e8416e9b0de7ce46821506", "2ea3ff42b3b84d5cbc9a575d4c436878", "2024-04-12")
	if err != nil {
		fmt.Println(err)
	}
}

func TestGuild_GenStatRecords(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	paymentDate, _ := notion.ParseDateTime("2023-12-15")
	data, _ := g.db.GetIncentiveData(&notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			{
				Property: "Payment Date",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Date: &notion.DatePropertyFilter{
						After: &paymentDate.Time,
					},
				},
			},
		},
	})
	paydateMap := map[string]int{}
	for _, d := range data {
		paydateMap[d.PaymentDate]++
	}
	paydateSlice := []string{}
	for k, _ := range paydateMap {
		paydateSlice = append(paydateSlice, k)
	}
	sort.Slice(paydateSlice, func(i, j int) bool {
		return paydateSlice[i] < paydateSlice[j]
	})

	for _, pd := range paydateSlice {
		fmt.Println("准备开始", pd)
		pm := map[string]int{}
		pm[pd] = 1
		err := g.GenTotalIncentiveStat("04c301f8dc5448759c5919e618822854", pm)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("完成", pd)
	}
}

func TestGuild_GenTotalIncentiveStatRecords(t *testing.T) {
	c := config.New("config_temp")
	d := db.New(c)
	g := New(c, d)
	startDate, _ := time.Parse("2006-01-02", "2023-12-08")
	endDate, _ := time.Parse("2006-01-02", "2024-04-06")
	for !startDate.After(endDate) {
		startDateStr := startDate.Format("2006-01-02")
		success, paydateMap, err := g.GenIncentiveStat("4c19704d927f4d52b2f030ebd1648ef3", startDateStr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(paydateMap)
		if success {
			err := g.GenTotalIncentiveStat("04c301f8dc5448759c5919e618822854", paydateMap)
			if err != nil {
				return
			}
		}
		fmt.Println("完成" + startDateStr)
		startDate = startDate.AddDate(0, 0, 7)
		fmt.Println("下一时间", startDate)
	}
}
