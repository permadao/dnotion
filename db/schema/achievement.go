package schema

import (
	"time"

	"github.com/dstotijn/go-notion"
)

type AchievementData struct {
	NID   string // notion id for update
	ID    string
	Guild string
	Tags  []string
	Date  string
}

func (a *AchievementData) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	a.NID = nid
	if len(props["ID"].Title) > 0 {
		a.ID = props["ID"].Title[0].Text.Content
	}
	if props["Guild"].Select != nil {
		a.Guild = props["Guild"].Select.Name
	}
	if props["Tags"].MultiSelect != nil {
		tags := make([]string, len(props["Tags"].MultiSelect))
		for i, tag := range props["Tags"].MultiSelect {
			tags[i] = tag.Name
		}
	}
	if props["Date"].Date != nil {
		a.Date = props["Date"].Date.Start.Format("2006-01-02")
	}
}

func (a *AchievementData) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if a.ID != "" {
		props["ID"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: a.ID,
					},
				},
			},
		}
	}
	if a.Guild != "" {
		props["Guild"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: a.Guild},
		}
	}
	if a.Tags != nil {
		selects := make([]notion.SelectOptions, len(a.Tags))
		for i, tagName := range a.Tags {
			selects[i] = notion.SelectOptions{Name: tagName}
		}
		props["Tags"] = notion.DatabasePageProperty{
			MultiSelect: selects,
		}
	}
	if a.Date != "" {
		date, err := notion.ParseDateTime(a.Date)
		if err != nil {
			curTime := time.Now()
			formattedDate := curTime.Format("20060102")
			date, err = notion.ParseDateTime(formattedDate)
		}
		if err == nil {
			props["Date"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	return a.NID, &props
}
