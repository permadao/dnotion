package schema

import (
	"github.com/dstotijn/go-notion"
)

type ContributorData struct {
	NID         string // notion id for update
	ID          *float64
	NotionName  string
	NotionID    string
	Wallet      string
	BuddyNotion string
	CreatedTime string
	// TODO: more fields
}

func (f *ContributorData) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	f.NID = nid
	if props["ID"].Number != nil {
		f.ID = props["ID"].Number
	}
	if len(props["Notion Name"].RichText) > 0 {
		f.NotionName = props["Notion Name"].RichText[0].Text.Content
	}
	if len(props["Notion ID"].RichText) > 0 {
		f.NotionID = props["Notion ID"].RichText[0].Text.Content
	}
	if len(props["Wallet"].RichText) > 0 {
		f.Wallet = props["Wallet"].RichText[0].Text.Content
	}
	if len(props["Buddy Notion"].People) > 0 {
		f.BuddyNotion = props["Buddy Notion"].People[0].BaseUser.ID
	}
	if props["Created time"].CreatedTime != nil {
		f.CreatedTime = props["Created time"].CreatedTime.Format("2006-01-02")
	}
}

func (f *ContributorData) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if f.ID != nil {
		props["ID"] = notion.DatabasePageProperty{
			Number: f.ID,
		}
	}
	if f.NotionName != "" {
		props["Notion Name"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: f.NotionName,
					},
				},
			},
		}
	}
	if f.NotionID != "" {
		props["Notion ID"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: f.NotionID,
					},
				},
			},
		}
	}
	if f.Wallet != "" {
		props["Wallet"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: f.Wallet,
					},
				},
			},
		}
	}
	return f.NID, &props
}
