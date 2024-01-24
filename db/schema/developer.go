package schema

import (
	"github.com/dstotijn/go-notion"
	"time"
)

type Developer struct {
	NID         string // notion id for update
	ID          string
	Contributor string
	Level       string
	Income      float64
	Date        string
}

func (d *Developer) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	d.NID = nid
	if len(props["ID"].Title) > 0 {
		d.ID = props["ID"].Title[0].Text.Content
	}
	if len(props["Contributor"].RichText) > 0 {
		d.Contributor = props["Contributor"].RichText[0].PlainText
	}
	if props["Level"].Select != nil {
		d.Level = props["Level"].Select.Name
	}
	if props["Income"].Number != nil {
		d.Income = *props["Income"].Number
	}
	if props["Date"].Date != nil {
		d.Date = props["Date"].Date.Start.Format("2006-01-02")
	}
}

func (d *Developer) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if d.ID != "" {
		props["ID"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: d.ID,
					},
				},
			},
		}
	}
	if d.Contributor != "" {
		props["Contributor"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: d.Contributor,
					},
				},
			},
		}
	}
	if d.Level != "" {
		props["Level"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: d.Level},
		}
	}
	if d.Income != 0 {
		props["Income"] = notion.DatabasePageProperty{
			Number: &d.Income,
		}
	}
	if d.Date != "" {
		date, err := notion.ParseDateTime(d.Date)
		if err != nil {
			curTime := time.Now()
			formattedDate := curTime.Format("2006-01-02")
			date, err = notion.ParseDateTime(formattedDate)
		}
		if err == nil {
			props["Date"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	return d.NID, &props
}
