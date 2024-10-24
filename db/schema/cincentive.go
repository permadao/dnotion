package schema

import (
	"github.com/dstotijn/go-notion"
	"time"
)

type CIncentive struct {
	NID               string
	ID                float64
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
	Token             string
}

func (t *CIncentive) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	t.NID = nid
	if props["ID"].Number != nil {
		t.ID = *props["ID"].Number
	}
	if props["Accounting Date"].Date != nil {
		t.AccountingDate = props["Accounting Date"].Date.Start.Format("2006-01-02")
	}
	if len(props["Notion ID"].Title) > 0 {
		t.NotionID = props["Notion ID"].Title[0].Text.Content
	}
	if len(props["Notion Name"].RichText) > 0 {
		t.NotionName = props["Notion Name"].RichText[0].PlainText
	}
	if len(props["Buddy Notion"].RichText) > 0 {
		t.BuddyNotion = props["Buddy Notion"].RichText[0].PlainText
	}
	if props["Total Incentive"].Number != nil {
		t.TotalIncentive = *props["Total Incentive"].Number
	}
	if props["Weekly Incentive"].Number != nil {
		t.WeeklyIncentive = *props["Weekly Incentive"].Number
	}
	if props["Payment Date"].Date != nil {
		t.PaymentDate = props["Payment Date"].Date.Start.Format("2006-01-02")
	}
	if props["Onboard Date"].Date != nil {
		t.OnboardDate = props["Onboard Date"].Date.Start.Format("2006-01-02")
	}
	if props["First contribution"].Select != nil {
		t.FirstContribution = props["First contribution"].Select.Name
	}
	if props["Medal"].Select != nil {
		t.Medal = props["Medal"].Select.Name
	}
	if props["Token"].Select != nil {
		t.Token = props["Token"].Select.Name
	}
}

func (t *CIncentive) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if t.ID != 0 {
		props["ID"] = notion.DatabasePageProperty{
			Number: &t.ID,
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
	if t.Token != "" {
		props["Token"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: t.Token},
		}
	}
	return t.NID, &props
}
