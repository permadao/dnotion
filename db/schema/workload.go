package schema

import "github.com/dstotijn/go-notion"

type WorkloadData struct {
	NID        string
	ID         string
	Status     string
	TaskStatus string
	Name       string
	Note       string
	Amount     float64 // usd amount except translation guide
	// TranslationAmount float64 // !! usd amount for translation guide -- Deprecated, remove at 2024.6.27 by webb
	Contributor string
	Token       string
}

func (f *WorkloadData) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	f.NID = nid
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
	if props["Amount"].Number != nil {
		f.Amount = *props["Amount"].Number
	}
	// -- remove at 2024.6.27 by webb
	// else if props["Amount"].Formula != nil {
	// 	f.TranslationAmount = *props["Amount"].Formula.Number
	// }
	if len(props["Contributor"].People) > 0 {
		f.Contributor = props["Contributor"].People[0].BaseUser.ID
	}
	if props["Token"].Select != nil {
		f.Token = props["Token"].Select.Name
	}
}

func (f *WorkloadData) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
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
	if f.Amount != 0 {
		props["Amount"] = notion.DatabasePageProperty{
			Number: &f.Amount,
		}
	}
	// -- remove at 2024.6.27 by webb
	// else if f.TranslationAmount != 0 {
	// 	props["Amount"] = notion.DatabasePageProperty{
	// 		Formula: &notion.FormulaResult{
	// 			Number: &f.TranslationAmount,
	// 		},
	// 	}
	// }
	if f.Contributor != "" {
		props["Contributor"] = notion.DatabasePageProperty{
			People: []notion.User{
				{
					BaseUser: notion.BaseUser{ID: f.Contributor},
				},
			},
		}
	}
	if f.Token != "" {
		props["Token"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: f.Token},
		}
	}
	return f.NID, &props
}
