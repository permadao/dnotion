package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetNewsFinances(nid string, filter *notion.DatabaseQueryFilter) ([]schema.NewsFinData, error) {
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	finDatas := []schema.NewsFinData{}
	for _, page := range pages {
		finData := NewsFinDataFromPage(page)
		finDatas = append(finDatas, *finData)
	}

	return finDatas, nil
}

func NewsFinDataFromPage(page notion.Page) *schema.NewsFinData {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewsFinDataFromProps(page.ID, props)
}

func NewsFinDataFromProps(nid string, props notion.DatabasePageProperties) *schema.NewsFinData {
	finData := &schema.NewsFinData{}
	finData.DeserializePropertys(nid, props)
	return finData
}
