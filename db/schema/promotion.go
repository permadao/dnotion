package schema

import (
	"github.com/dstotijn/go-notion"
	"time"
)

type PromotionStat struct {
	NID   string
	ID    string
	OutDB string
	Date  string
}

type PromotionPoints struct {
	NID         string
	Contributor string
	BasePoints  float64
	Task        string
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
	if len(props["Contributor"].People) > 0 {
		p.Contributor = props["Contributor"].People[0].BaseUser.ID
	}
	if props["Base Points"].Rollup != nil {
		p.BasePoints = *props["Base Points"].Rollup.Number
	}
	if len(props["Task"].Relation) > 0 {
		p.Task = props["Task"].Relation[0].ID
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
	return p.NID, &props
}
func (p *PromotionSettlement) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	p.NID = nid
	if len(props["Contributor"].People) > 0 {
		p.Contributor = props["Contributor"].People[0].ID
	} else if len(props["Contributor Name"].RichText) > 0 {
		p.Contributor = props["Contributor Name"].RichText[0].PlainText
	}
	if props["Total Score"].Number != nil {
		p.TotalScore = *props["Total Score"].Number
	}
	if props["Personal Score"].Number != nil {
		p.PersonalScore = *props["Personal Score"].Number
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
	//if p.Contributor != "" {
	//	props["Contributor"] = notion.DatabasePageProperty{
	//		People: []notion.User{
	//			{
	//				BaseUser: notion.BaseUser{ID: p.Contributor},
	//			},
	//		},
	//	}
	//}
	if p.Contributor != "" {
		props["Contributor Name"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: p.Contributor,
					},
				},
			},
		}
	}
	if p.TotalScore != 0 {
		props["Total Score"] = notion.DatabasePageProperty{
			Number: &p.TotalScore,
		}
	}
	if p.PersonalScore != 0 {
		props["Personal Score"] = notion.DatabasePageProperty{
			Number: &p.PersonalScore,
		}
	}
	if p.Rewards != 0 {
		props["Rewards"] = notion.DatabasePageProperty{
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
