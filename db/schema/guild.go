package schema

import (
	"time"

	"github.com/dstotijn/go-notion"
)

type GuildData struct {
	NID                string // notion id for update
	Name               string
	Info               string
	Link               string // link to guild notion id
	Tags               []string
	TotalContributors  float64
	WeeklyContributors float64
	NewContributors    float64
	TotalIncentive     float64
	WeeklyIncentive    float64
	Date               string
	Rank               float64
}

func (a *GuildData) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	a.NID = nid
	if len(props["Name"].Title) > 0 {
		a.Name = props["Name"].Title[0].Text.Content
	}
	if len(props["Info"].RichText) > 0 {
		a.Info = props["Info"].RichText[0].Text.Content
	}
	if len(props["Link"].RichText) > 0 {
		a.Link = props["Link"].RichText[0].Mention.Page.ID
	}
	if props["Tags"].MultiSelect != nil {
		tags := make([]string, len(props["Tags"].MultiSelect))
		for i, tag := range props["Tags"].MultiSelect {
			tags[i] = tag.Name
		}
	}
	if props["Total Contributors"].Number != nil {
		a.TotalContributors = *props["Total Contributors"].Number
	}
	if props["Weekly Contributors"].Number != nil {
		a.WeeklyContributors = *props["Weekly Contributors"].Number
	}
	if props["New Contributors"].Number != nil {
		a.NewContributors = *props["New Contributors"].Number
	}
	if props["Total Incentive"].Number != nil {
		a.NewContributors = *props["Total Incentive"].Number
	}
	if props["Weekly Incentive"].Number != nil {
		a.NewContributors = *props["Weekly Incentive"].Number
	}
	if props["Date"].Date != nil {
		a.Date = props["Date"].Date.Start.Format("2006-01-02")
	}
	if props["Rank"].Number != nil {
		a.Rank = *props["Rank"].Number
	}
}

func (a *GuildData) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if a.Name != "" {
		props["Name"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: a.Name,
					},
				},
			},
		}
	}
	if a.Info != "" {
		props["Info"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: a.Info,
					},
				},
			},
		}
	}
	if a.Link != "" {
		props["Link"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					// Type:      notion.RichTextTypeMention,
					// PlainText: a.Name,
					Mention: &notion.Mention{
						Page: &notion.ID{a.Link},
					},
				},
			},
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
	props["Total Contributors"] = notion.DatabasePageProperty{
		Number: &a.TotalContributors,
	}
	props["Weekly Contributors"] = notion.DatabasePageProperty{
		Number: &a.WeeklyContributors,
	}
	props["New Contributors"] = notion.DatabasePageProperty{
		Number: &a.NewContributors,
	}
	props["Total Incentive"] = notion.DatabasePageProperty{
		Number: &a.TotalIncentive,
	}
	props["Weekly Incentive"] = notion.DatabasePageProperty{
		Number: &a.WeeklyIncentive,
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
	if a.Rank != 0 {
		props["Rank"] = notion.DatabasePageProperty{
			Number: &a.Rank,
		}
	}
	return a.NID, &props
}
