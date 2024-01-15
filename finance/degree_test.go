package finance

import (
	"encoding/json"
	"fmt"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/guild"
	"testing"
)

func TestF(t *testing.T) {
	myconfig := config.New("config_temp")
	d := db.New(myconfig)
	g := guild.New(myconfig, d)
	amount, contributors, contributor, err := g.StatWeeklyFinance("AR", "e8d79c55c0394cba83664f3e5737b0bd", "2024-01-12")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(fmt.Sprintf("总金额：%f", amount))
	fmt.Println(json.Marshal(contributors))
	fmt.Println(json.Marshal(contributor))
}
