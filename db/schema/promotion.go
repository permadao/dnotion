package schema

import (
	"github.com/dstotijn/go-notion"
	"time"
)

type PromotionStat struct {
	NID  string
	ID   string
	Date string
}

type PromotionPoints struct {
	NID         string
	Contributor string
	BasePoints  string
}

type PromotionSettlement struct {
	NID           string
	Contributor   string
	TotalScore    float64
	PersonalScore float64
	Rewards       float64
	Date          string
}

func (p *PromotionStat) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	p.NID = nid
	if len(props["ID"].Title) > 0 {
		p.ID = props["ID"].Title[0].Text.Content
	}
	if props["Date"].Date != nil {
		p.Date = props["Date"].Date.Start.Format("2006-01-02")
	}
}

func (p *PromotionStat) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if p.ID != "" {
		props["ID"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: p.ID,
					},
				},
			},
		}
	}
	if p.Date != "" {
		date, err := notion.ParseDateTime(p.Date)
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
	return p.NID, &props
}

func (p *PromotionPoints) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	p.NID = nid
	if len(props["Contributor"].Title) > 0 {
		p.Contributor = props["Contributor"].Title[0].Text.Content
	}
	if len(props["BasePoints"].Title) > 0 {
		p.Contributor = props["BasePoints"].Title[0].Text.Content
	}
}

func (p *PromotionPoints) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if p.Contributor != "" {
		props["Contributor"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: p.Contributor,
					},
				},
			},
		}
	}
	if p.BasePoints != "" {
		props["BasePoints"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: p.BasePoints,
					},
				},
			},
		}
	}
	return p.NID, &props
}
func (p *PromotionSettlement) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	p.NID = nid
	if len(props["Contributor"].Title) > 0 {
		p.Contributor = props["Contributor"].Title[0].Text.Content
	}
	if props["TotalScore"].Number != nil {
		p.TotalScore = *props["TotalScore"].Number
	}
	if props["PersonalScore"].Number != nil {
		p.PersonalScore = *props["PersonalScore"].Number
	}
	if props["Rewards"].Number != nil {
		p.Rewards = *props["Rewards"].Number
	}
	if props["Date"].Date != nil {
		p.Date = props["Date"].Date.Start.Format("2006-01-02")
	}
}

func (p *PromotionSettlement) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if p.Contributor != "" {
		props["Contributor"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: p.Contributor,
					},
				},
			},
		}
	}
	if p.TotalScore != 0 {
		props["TotalScore"] = notion.DatabasePageProperty{
			Number: &p.TotalScore,
		}
	}
	if p.PersonalScore != 0 {
		props["PersonalScore"] = notion.DatabasePageProperty{
			Number: &p.PersonalScore,
		}
	}
	if p.Rewards != 0 {
		props["Income"] = notion.DatabasePageProperty{
			Number: &p.Rewards,
		}
	}
	if p.Date != "" {
		date, err := notion.ParseDateTime(p.Date)
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
	return p.NID, &props
}
