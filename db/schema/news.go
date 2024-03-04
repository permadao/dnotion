package schema
import (
	"github.com/dstotijn/go-notion"
	"time"
	// "fmt"
)
type News struct {
	NID         string // notion id for update
	ID          string
	Executor string
	ExecutorWorkload       string
	Rank      string
	WorkloadOf15Weeds        float64
	RankCode string
	CumulativeWorkload  float64
	Date string
}

func (d *News) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	d.NID = nid
	// fmt.Println(props["ID"].RichText[0].PlainText)
	// fmt.Println(props["Executor Workload"])
	// fmt.Println(props["Rank"])
	// fmt.Println(props["Executor"])
	// fmt.Println(props["Cumulative Workload"])
	// fmt.Println(props["Rank EN"])
	// fmt.Println(props["Workload of the last 15 Weeks"])
	// fmt.Println(props["Date"])

	if  len(props["ID"].RichText)>0 {
		d.ID = props["ID"].RichText[0].PlainText
	}
	if len(props["Executor"].People) > 0 {
		d.Executor = props["Executor"].People[0].BaseUser.ID
	}
	if props["Executor Workload"].Title != nil {
		d.ExecutorWorkload = props["Executor Workload"].Title[0].Text.Content
	}
	if props["Rank"].Select != nil {
		d.Rank = props["Rank"].Select.Name
	}
	if props["Workload of the last 15 Weeks"].Number != nil {
		// fmt.Println(*props["Workload of the last 15 Weeks"].Number)
		d.WorkloadOf15Weeds = *props["Workload of the last 15 Weeks"].Number
	}
	if len(props["Rank EN"].RichText)>0  {
		d.RankCode = props["Rank EN"].RichText[0].PlainText
	}
	if props["Cumulative Workload"].Number != nil {
		d.CumulativeWorkload = *props["Cumulative Workload"].Number
	}
	if props["Date"].Date != nil {
		d.Date = props["Date"].Date.Start.Format("2006-01-02")
	}
}

func (d *News) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	if d.ID != "" {
		props["ID"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: d.ID,
					},
				},
			},
	    }
    }
	if d.Executor != "" {
		props["Executor"] = notion.DatabasePageProperty{
			People: []notion.User{
				{
					BaseUser: notion.BaseUser{ID: d.Executor},
				},
			},
		}
	}
	if d.ExecutorWorkload != "" {
		props["Executor Workload"] = notion.DatabasePageProperty{
			Title: []notion.RichText{
				{
					Text: &notion.Text{
						Content: d.ExecutorWorkload,
					},
				},
			},
		}
	}
	if d.Rank != "" {
		props["Rank"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name:  d.Rank},
		}
	}
	if d.WorkloadOf15Weeds != 0 {
		props["Workload of the last 15 Weeks"] = notion.DatabasePageProperty{
			Number: &d.WorkloadOf15Weeds,
		}
	}
	if d.RankCode != "" {
		props["Rank EN"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Text: &notion.Text{
						Content: d.RankCode,
					},
				},
			},
	    }
    }
	if d.CumulativeWorkload != 0 {
		props["Cumulative Workload"] = notion.DatabasePageProperty{
			Number: &d.CumulativeWorkload,
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
