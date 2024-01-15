package schema

import "github.com/dstotijn/go-notion"

type Translator struct {
	NID         string // notion id for update
	ID          string
	Contributor string
	Level       string
	Title       string
}

func (f *Translator) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	f.NID = nid
	if len(props["ID"].Title) > 0 {
		f.ID = props["ID"].Title[0].Text.Content
	}
	if len(props["Contributor"].People) > 0 {
		f.Contributor = props["Contributor"].People[0].Name
	}
	if props["Level"].Select != nil {
		f.Level = props["Level"].Select.Name
	}
	if props["Title"].Select != nil {
		f.Title = props["Title"].Select.Name
	}
}

func (f *Translator) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
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
	if f.Contributor != "" {
		props["Contributor"] = notion.DatabasePageProperty{
			People: []notion.User{
				{
					Name: f.Contributor,
				},
			},
		}
	}
	if f.Level != "" {
		props["Level"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: f.Level},
		}
	}
	if f.Title != "" {
		props["Title"] = notion.DatabasePageProperty{
			Select: &notion.SelectOptions{Name: f.Title},
		}
	}
	return f.NID, &props
}
