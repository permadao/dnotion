package main

import (
	"time"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/guild"
	"fmt"
)
func GetPreviousDate(days int) (date string) {
	now := time.Now()
	last := now.AddDate(0, 0, -days)
	date = last.Format("2006-01-02")
	return
}

func main() {
	config := config.New("./config.toml")
	db := db.New(config)
	g := guild.New(config, db)
	start := GetPreviousDate(0)
	startDateOfNews := GetPreviousDate(15 * 7)
	fmt.Println(startDateOfNews)
	if err := g.GenNewsGrade("ad2cf585b08843fea7cf40a682bc4529", "d5f9fc70910b45d4ab8811f37716637d", startDateOfNews, start); err != nil {
		fmt.Sprintf("Contributor not found, nid/id:", "11")
	}
}