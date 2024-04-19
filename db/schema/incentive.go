package schema

import (
	"github.com/dstotijn/go-notion"
	"time"
)

type Incentive struct {
	NID               string
	ID                string
	AccountingDate    string
	NotionID          string
	Guild             string
	NotionName        string
	BuddyNotion       string
	TotalIncentive    float64
	WeeklyIncentive   float64
	PaymentDate       string
	OnboardDate       string
	FirstContribution string
}

type TotalIncentive struct {
	NID               string
	ID                string
	AccountingDate    string
	NotionID          string
	NotionName        string
	BuddyNotion       string
	TotalIncentive    float64
	WeeklyIncentive   float64
	PaymentDate       string
	OnboardDate       string
	FirstContribution string
	Medal             string
}

func (i *Incentive) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	i.NID = nid
	if len(props["ID"].Title) > 0 {
		i.ID = props["ID"].Title[0].Text.Content
	}
	if props["Accounting Date"].Date != nil {
		i.AccountingDate = props["Accounting Date"].Date.Start.Format("2006-01-02")
	}
	if len(props["Notion ID"].Title) > 0 {
		i.NotionID = props["Notion ID"].Title[0].Text.Content
	}
	if props["Guild"].Select != nil {
		i.Guild = props["Guild"].Select.Name
	}
	if len(props["Notion Name"].RichText) > 0 {
		i.NotionName = props["Notion Name"].RichText[0].PlainText
	}
	if len(props["Buddy Notion"].RichText) > 0 {
		i.BuddyNotion = props["Buddy Notion"].RichText[0].PlainText
	}
	if props["Total Incentive"].Number != nil {
		i.TotalIncentive = *props["Total Incentive"].Number
	}
	if props["Weekly Incentive"].Number != nil {
		i.WeeklyIncentive = *props["Weekly Incentive"].Number
	}
	if props["Payment Date"].Date != nil {
		i.PaymentDate = props["Payment Date"].Date.Start.Format("2006-01-02")
	}
	if props["Onboard Date"].Date != nil {
		i.OnboardDate = props["Onboard Date"].Date.Start.Format("2006-01-02")
	}
	if props["First contribution"].Select != nil {
		i.FirstContribution = props["First contribution"].Select.Name
	}
}

func (i *Incentive) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if i.ID != "" {
		props["ID"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: i.ID,
					},
				},
			},
		}
	}
	if i.AccountingDate != "" {
		date, err := notion.ParseDateTime(i.AccountingDate)
		if err != nil {
			curTime := time.Now()
			formattedDate := curTime.Format("2006-01-02")
			date, err = notion.ParseDateTime(formattedDate)
		}
		if err == nil {
			props["Accounting Date"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	if i.Guild != "" {
		props["Guild"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: i.Guild},
		}
	}
	if i.NotionID != "" {
		props["Notion ID"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: i.NotionID,
					},
				},
			},
		}
	}
	if i.NotionName != "" {
		props["Notion Name"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: i.NotionName,
					},
				},
			},
		}
	}
	if i.BuddyNotion != "" {
		props["Buddy Notion"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: i.BuddyNotion,
					},
				},
			},
		}
	}
	if i.TotalIncentive != 0 {
		props["Total Incentive"] = notion.DatabasePageProperty{
			Number: &i.TotalIncentive,
		}
	}
	if i.WeeklyIncentive != 0 {
		props["Weekly Incentive"] = notion.DatabasePageProperty{
			Number: &i.WeeklyIncentive,
		}
	}
	if i.PaymentDate != "" {
		date, err := notion.ParseDateTime(i.PaymentDate)
		if err != nil {
			curTime := time.Now()
			formattedDate := curTime.Format("2006-01-02")
			date, err = notion.ParseDateTime(formattedDate)
		}
		if err == nil {
			props["Payment Date"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	if i.OnboardDate != "" {
		date, err := notion.ParseDateTime(i.OnboardDate)
		if err != nil {
			curTime := time.Now()
			formattedDate := curTime.Format("2006-01-02")
			date, err = notion.ParseDateTime(formattedDate)
		}
		if err == nil {
			props["Onboard Date"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	if i.FirstContribution != "" {
		props["First contribution"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: i.FirstContribution},
		}
	}
	return i.NID, &props
}

func (t *TotalIncentive) DeserializePropertys(nid string, props notion.DatabasePageProperties) {

}

func (t *TotalIncentive) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if t.ID != "" {
		props["ID"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: t.ID,
					},
				},
			},
		}
	}
	if t.AccountingDate != "" {
		date, err := notion.ParseDateTime(t.AccountingDate)
		if err != nil {
			curTime := time.Now()
			formattedDate := curTime.Format("2006-01-02")
			date, err = notion.ParseDateTime(formattedDate)
		}
		if err == nil {
			props["Accounting Date"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	if t.NotionID != "" {
		props["Notion ID"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: t.NotionID,
					},
				},
			},
		}
	}
	if t.NotionName != "" {
		props["Notion Name"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: t.NotionName,
					},
				},
			},
		}
	}
	if t.BuddyNotion != "" {
		props["Buddy Notion"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: t.BuddyNotion,
					},
				},
			},
		}
	}
	if t.TotalIncentive != 0 {
		props["Total Incentive"] = notion.DatabasePageProperty{
			Number: &t.TotalIncentive,
		}
	}
	if t.WeeklyIncentive != 0 {
		props["Weekly Incentive"] = notion.DatabasePageProperty{
			Number: &t.WeeklyIncentive,
		}
	}
	if t.PaymentDate != "" {
		date, err := notion.ParseDateTime(t.PaymentDate)
		if err != nil {
			curTime := time.Now()
			formattedDate := curTime.Format("2006-01-02")
			date, err = notion.ParseDateTime(formattedDate)
		}
		if err == nil {
			props["Payment Date"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	if t.OnboardDate != "" {
		date, err := notion.ParseDateTime(t.OnboardDate)
		if err != nil {
			curTime := time.Now()
			formattedDate := curTime.Format("2006-01-02")
			date, err = notion.ParseDateTime(formattedDate)
		}
		if err == nil {
			props["Onboard Date"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	if t.FirstContribution != "" {
		props["First contribution"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: t.FirstContribution},
		}
	}

	if t.Medal != "" {
		props["Medal"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: t.Medal},
		}
	}
	return t.NID, &props
}
