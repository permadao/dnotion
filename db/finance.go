package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetFinances(nid string, filter *notion.DatabaseQueryFilter) ([]schema.FinData, error) {
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	finDatas := []schema.FinData{}
	for _, page := range pages {
		finData := NewFinDataFromPage(page)
		finDatas = append(finDatas, *finData)
	}

	return finDatas, nil
}

func NewFinDataFromPage(page notion.Page) *schema.FinData {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewFinDataFromProps(page.ID, props)
}

func NewFinDataFromProps(nid string, props notion.DatabasePageProperties) *schema.FinData {
	finData := &schema.FinData{}
	finData.DeserializePropertys(nid, props)
	return finData
}
