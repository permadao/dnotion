package db

import (
	"context"
	"time"

	"github.com/dstotijn/go-notion"
)

type FinData struct {
	ID             string
	Amount         float64
	ActualToken    string
	ActualPrice    float64
	TargetToken    string
	TargetPrice    float64
	TargetAmount   float64
	Status         string
	WorkloadID     string
	WorkloadStatus string
	PaymentDate    string
	Contributor    string
	ReceiptUrl     string
	AffiliatedDAO  string
}

func NewFinDataFromProps(props *notion.DatabasePageProperties) *FinData {
	finData := &FinData{}
	finData.DeserializePropertys(*props)
	return finData
}

func (f *FinData) UpdatePage(pageId string) error {
	props := f.SerializePropertys()
	_, err := DB.DBClient.UpdatePage(context.Background(), pageId, notion.UpdatePageParams{
		DatabasePageProperties: *props,
	})
	return err
}

func (f *FinData) CreatePage(parentID string) error {
	_, err := DB.DBClient.CreatePage(context.Background(), notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               parentID,
		DatabasePageProperties: f.SerializePropertys(),
	})
	return err
}

func (f *FinData) SerializePropertys() *notion.DatabasePageProperties {
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
	if f.Amount != 0 {
		props["Amount"] = notion.DatabasePageProperty{
			Number: &f.Amount,
		}
	}
	if f.ActualToken != "" {
		props["Actual Token"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: f.ActualToken},
		}
	}
	if f.ActualPrice != 0 {
		props["Actual Price"] = notion.DatabasePageProperty{
			Number: &f.ActualPrice,
		}
	}
	if f.TargetToken != "" {
		props["Target Token"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: f.TargetToken},
		}
	}
	if f.TargetPrice != 0 {
		props["Target Price"] = notion.DatabasePageProperty{
			Number: &f.TargetPrice,
		}
	}
	if f.Status != "" {
		props["Status"] = notion.DatabasePageProperty{
			Status: &notion.SelectOptions{Name: f.Status},
		}
	}
	if f.WorkloadStatus != "" {
		props["Workload Status"] = notion.DatabasePageProperty{
			Status: &notion.SelectOptions{Name: f.WorkloadStatus},
		}
	}
	if f.WorkloadID != "" {
		props["Workload ID"] = notion.DatabasePageProperty{
			Relation: []notion.Relation{
				{ID: f.WorkloadID},
			},
		}
	}
	if f.PaymentDate != "" {
		date, err := notion.ParseDateTime(f.PaymentDate)
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
	if f.ReceiptUrl != "" {
		props["Receipt(url)"] = notion.DatabasePageProperty{
			URL: &f.ReceiptUrl,
		}
	}
	if f.Contributor != "" {
		props["Contributor"] = notion.DatabasePageProperty{
			Relation: []notion.Relation{
				{ID: f.Contributor},
			},
		}
	}
	if f.AffiliatedDAO != "" {
		props["Affiliated DAO"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: f.AffiliatedDAO},
		}
	}
	return &props
}

func (f *FinData) DeserializePropertys(props notion.DatabasePageProperties) {
	if len(props["ID"].Title) > 0 {
		f.ID = props["ID"].Title[0].Text.Content
	}
	if props["Amount"].Number != nil {
		f.Amount = *props["Amount"].Number
	}
	if props["Actual Token"].Select != nil {
		f.ActualToken = props["Actual Token"].Select.Name
	}
	if props["Actual Price"].Number != nil {
		f.ActualPrice = *props["Actual Price"].Number
	}
	if props["Target Token"].Select != nil {
		f.TargetToken = props["Target Token"].Select.Name
	}
	if props["Target Price"].Number != nil {
		f.TargetPrice = *props["Target Price"].Number
	}
	if props["Target Amount"].Formula.Number != nil {
		f.TargetAmount = *props["Target Amount"].Formula.Number
	}
	if props["Status"].Select != nil {
		f.Status = props["Status"].Select.Name
	}
	if len(props["Workload ID"].Relation) > 0 {
		f.WorkloadID = props["Workload ID"].Relation[0].ID
	}
	if props["Workload Status"].Select != nil {
		f.WorkloadStatus = props["Workload Status"].Select.Name
	}
	if props["Payment Date"].Date != nil {
		f.PaymentDate = props["Payment Date"].Date.Start.Format("2006-01-02")
	}
	if props["Receipt(url)"].URL != nil {
		f.ReceiptUrl = *props["Receipt(url)"].URL
	}
	if len(props["Contributor"].Relation) > 0 {
		f.Contributor = props["Contributor"].Relation[0].ID
	}
	if props["Affiliated DAO"].Select != nil {
		f.AffiliatedDAO = props["Affiliated DAO"].Select.Name
	}
}
