package schema

import (
	"time"
	// "fmt"
	"github.com/dstotijn/go-notion"
)

type NewsFinData struct {
	NID             string // notion id for update
	ID              string
	CreatedTime     string
	TaskStatus      string
	Contributor     string
	ContributorRank string
	Amount          float64
}

func (f *NewsFinData) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	f.NID = nid

	// fmt.Println( props["Task Status"] )
	// fmt.Println( props["Contributor_rank"] )
	if len(props["ID"].Title) > 0 {
		f.ID = props["ID"].Title[0].Text.Content
	}
	if props["Created time"].CreatedTime != nil {
		f.CreatedTime = props["Created time"].CreatedTime.Format("2006-01-02")
	}
	if props["Task Status"].Select != nil {
		f.TaskStatus = props["Task Status"].Select.Name
	}
	if len(props["Contributor"].People) > 0 {
		f.Contributor = props["Contributor"].People[0].BaseUser.ID
	}
	if len(props["Contributor_rank"].Relation) > 0 {
		f.ContributorRank = props["Contributor_rank"].Relation[0].ID
	}
	if props["Amount"].Number != nil {
		f.Amount = *props["Amount"].Number
	}
	// fmt.Println(f)
}

func (f *NewsFinData) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if f.ID != "" {
		props["ID"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: f.ID,
					},
				},
			},
		}
	}
	if f.CreatedTime != "" {
		date, err := notion.ParseDateTime(f.CreatedTime)
		if err != nil {
			curTime := time.Now()
			formattedDate := curTime.Format("2006-01-02")
			date, err = notion.ParseDateTime(formattedDate)
		}
		if err == nil {
			props["Created time"] = notion.DatabasePageProperty{
				Date: &notion.Date{Start: date},
			}
		}
	}
	if f.TaskStatus != "" {
		props["Task Status"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: f.TaskStatus},
		}
	}
	if f.Contributor != "" {
		props["Contributor"] = notion.DatabasePageProperty{
			People: []notion.User{
				{
					BaseUser: notion.BaseUser{
						ID: f.Contributor,
					},
				},
			},
		}
	}

	if f.ContributorRank != "" {
		props["Contributor_rank"] = notion.DatabasePageProperty{
			Relation: []notion.Relation{
				{ID: f.ContributorRank},
			},
		}
	}

	if f.Amount != 0 {
		props["Amount"] = notion.DatabasePageProperty{
			Number: &f.Amount,
		}
	}
	return f.NID, &props
}
