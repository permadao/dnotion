package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
	// "fmt"
)

func (d *DB) GetNews(nid string, filter *notion.DatabaseQueryFilter) ([]schema.News, error) {
	// i1,_:=d.GetCount(nid)
	// fmt.Println(i1)
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}
	// fmt.Println(len(pages))
	NewsDatas := []schema.News{}
	for _, page := range pages {
		NewsData := NewNewsDataFromPage(page)
		NewsDatas = append(NewsDatas, *NewsData)
	}

	return NewsDatas, nil
}

func NewNewsDataFromPage(page notion.Page) *schema.News {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewNewsDataFromProps(page.ID, props)
}

func NewNewsDataFromProps(nid string, props notion.DatabasePageProperties) *schema.News {
	data := &schema.News{}
	data.DeserializePropertys(nid, props)
	return data
}
