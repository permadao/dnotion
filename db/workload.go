package db

import (
	"context"

	"github.com/dstotijn/go-notion"
)

type WorkloadData struct {
	ID             string
	Status         string
	TaskStatus     string
	Name           string
	Note           string
	Usd            float64 // usd amount except translation guide
	TranslattonUsd float64 // !! usd amount for translation guide
	Contributor    string
}

func NewWrokloadDataFromProps(props *notion.DatabasePageProperties) *WorkloadData {
	workloadData := &WorkloadData{}
	workloadData.DeserializePropertys(*props)
	return workloadData
}

func (f *WorkloadData) UpdatePage(pageId string) error {
	props := f.SerializePropertys()
	_, err := DB.DBClient.UpdatePage(context.Background(), pageId, notion.UpdatePageParams{
		DatabasePageProperties: *props,
	})
	return err
}

func (f *WorkloadData) CreatePage(parentID string) error {
	_, err := DB.DBClient.CreatePage(context.Background(), notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               parentID,
		DatabasePageProperties: f.SerializePropertys(),
	})
	return err
}

func (f *WorkloadData) SerializePropertys() *notion.DatabasePageProperties {
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
	if f.Status != "" {
		props["Status"] = notion.DatabasePageProperty{
			Status: &notion.SelectOptions{Name: f.Status},
		}
	}
	if f.TaskStatus != "" {
		props["Task Status"] = notion.DatabasePageProperty{
			Status: &notion.SelectOptions{Name: f.TaskStatus},
		}
	}
	if f.Name != "" {
		props["Name"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Type: notion.RichTextTypeText,
					Text: &notion.Text{
						Content: f.Name,
					},
					PlainText: f.Name,
				},
			},
		}
	}
	if f.Note != "" {
		props["Note"] = notion.DatabasePageProperty{
			RichText: []notion.RichText{
				{
					Type: notion.RichTextTypeText,
					Text: &notion.Text{
						Content: f.Note,
					},
					PlainText: f.Note,
				},
			},
		}
	}
	if f.Usd != 0 {
		props["USD"] = notion.DatabasePageProperty{
			Number: &f.Usd,
		}
	} else if f.TranslattonUsd != 0 {
		props["USD"] = notion.DatabasePageProperty{
			Formula: &notion.FormulaResult{
				Number: &f.TranslattonUsd,
			},
		}
	}
	if f.Contributor != "" {
		props["Contributor"] = notion.DatabasePageProperty{
			People: []notion.User{
				{
					BaseUser: notion.BaseUser{ID: f.Contributor},
				},
			},
		}
	}
	return &props
}

func (f *WorkloadData) DeserializePropertys(props notion.DatabasePageProperties) {
	if len(props["ID"].Title) > 0 {
		f.ID = props["ID"].Title[0].Text.Content
	}
	if props["Status"].Select != nil {
		f.Status = props["Status"].Select.Name
	}
	if props["Task Status"].Select != nil {
		f.TaskStatus = props["Task Status"].Select.Name
	}
	if len(props["Name"].RichText) > 0 {
		f.Name = props["Name"].RichText[0].Text.Content
	}
	if len(props["Note"].RichText) > 0 {
		f.Note = props["Note"].RichText[0].Text.Content
	}
	if props["USD"].Number != nil {
		f.Usd = *props["USD"].Number
	} else if props["USD"].Formula != nil {
		f.TranslattonUsd = *props["USD"].Formula.Number
	}
	if len(props["Contributor"].People) > 0 {
		f.Contributor = props["Contributor"].People[0].BaseUser.ID
	}
}
