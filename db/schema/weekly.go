package schema

import (
	"github.com/dstotijn/go-notion"
)

type WeeklyData struct {
	NID                 string  // Notion ID for update
	ID                  string  // ID column
	Guild               string  // Guild column
	USD                 float64 // USD column
	AR                  float64 // AR column
	ActualDistributedAR float64 // Actual Distributed AR column
	BP                  float64 // BP column
	ActualDistributedBP float64 // Actual Distributed BP column
	ARCurrencyRate      float64 // AR Currency Rate column
	Date                string  // Date column
	Status              string
}

func (w *WeeklyData) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	w.NID = nid
	if len(props["ID"].Title) > 0 {
		w.ID = props["ID"].Title[0].Text.Content
	}
	if props["Guild"].Select != nil {
		w.Guild = props["Guild"].Select.Name
	}
	if props["USD"].Number != nil {
		w.USD = *props["USD"].Number
	}
	if props["AR"].Number != nil {
		w.AR = *props["AR"].Number
	}
	if props["Actual Distributed AR"].Number != nil {
		w.ActualDistributedAR = *props["Actual Distributed AR"].Number
	}
	if props["BP"].Number != nil {
		w.BP = *props["BP"].Number
	}
	if props["Actual Distributed BP"].Number != nil {
		w.ActualDistributedBP = *props["Actual Distributed BP"].Number
	}
	if props["AR Currency Rate"].Number != nil {
		w.ARCurrencyRate = *props["AR Currency Rate"].Number
	}
	if props["Date"].Date != nil {
		w.Date = props["Date"].Date.Start.Format("2006-01-02")
	}
	if props["Status"].Select != nil {
		w.Status = props["Status"].Select.Name
	}
}

func (w *WeeklyData) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if w.ID != "" {
		props["ID"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: w.ID,
					},
				},
			},
		}
	}
	if w.Guild != "" {
		props["Guild"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: w.Guild},
		}
	}
	if w.USD != 0 {
		props["USD"] = notion.DatabasePageProperty{
			Number: &w.USD,
		}
	}
	if w.AR != 0 {
		props["AR"] = notion.DatabasePageProperty{
			Number: &w.AR,
		}
	}
	if w.ActualDistributedAR != 0 {
		props["Actual Distributed AR"] = notion.DatabasePageProperty{
			Number: &w.ActualDistributedAR,
		}
	}
	if w.BP != 0 {
		props["BP"] = notion.DatabasePageProperty{
			Number: &w.BP,
		}
	}
	if w.ActualDistributedBP != 0 {
		props["Actual Distributed BP"] = notion.DatabasePageProperty{
			Number: &w.ActualDistributedBP,
		}
	}
	if w.ARCurrencyRate != 0 {
		props["AR Currency Rate"] = notion.DatabasePageProperty{
			Number: &w.ARCurrencyRate,
		}
	}
	if w.Date != "" {
		date, err := notion.ParseDateTime(w.Date)
		if err == nil {
			props["Date"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	if w.Status != "" {
		props["Status"] = notion.DatabasePageProperty{
			Status: &notion.SelectOptions{Name: w.Status},
		}
	}
	return w.NID, &props
}
