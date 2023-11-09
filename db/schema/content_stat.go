package schema

import (
	"time"

	"github.com/dstotijn/go-notion"
)

type ConentStatData struct {
	NID       string // notion id for update
	Name      string // title
	Platform  string
	Hits      float64
	FrontPage []string
	CountTime string
}

func (c *ConentStatData) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	c.NID = nid
	if len(props["Name"].Title) > 0 {
		c.Name = props["Name"].Title[0].Text.Content
	}
	if props["Platform"].Select != nil {
		c.Platform = props["Platform"].Select.Name
	}
	if props["Hits"].Number != nil {
		c.Hits = *props["Hits"].Number
	}
	if props["Front Page"].MultiSelect != nil {
		fp := make([]string, len(props["Front Page"].MultiSelect))
		for i, tag := range props["Front Page"].MultiSelect {
			fp[i] = tag.Name
		}
		c.FrontPage = fp
	}
	if props["Count Time"].Date != nil {
		c.CountTime = props["Count Time"].Date.Start.Format("2006-01-02")
	}
}

func (c *ConentStatData) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if c.Name != "" {
		props["Name"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: c.Name,
					},
				},
			},
		}
	}
	if c.Platform != "" {
		props["Platform"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: c.Platform},
		}
	}
	props["Hits"] = notion.DatabasePageProperty{
		Number: &c.Hits,
	}
	if c.FrontPage != nil {
		selects := make([]notion.SelectOptions, len(c.FrontPage))
		for i, tagName := range c.FrontPage {
			selects[i] = notion.SelectOptions{Name: tagName}
		}
		props["Front Page"] = notion.DatabasePageProperty{
			MultiSelect: selects,
		}
	}
	if c.CountTime != "" {
		date, err := notion.ParseDateTime(c.CountTime)
		if err != nil {
			curTime := time.Now()
			formattedDate := curTime.Format("20060102")
			date, err = notion.ParseDateTime(formattedDate)
		}
		if err == nil {
			props["Count Time"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	return c.NID, &props
}
