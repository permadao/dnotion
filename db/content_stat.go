package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetContentStats(filter *notion.DatabaseQueryFilter) ([]schema.ConentStatData, error) {
	pages, err := d.GetPages(d.ContentStatDB, filter)
	if err != nil {
		return nil, err
	}

	datas := []schema.ConentStatData{}
	for _, page := range pages {
		data := NewContentStatDataFromPage(page)
		datas = append(datas, *data)
	}

	return datas, nil
}

func NewContentStatDataFromPage(page notion.Page) *schema.ConentStatData {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewContentStatDataFromProps(page.ID, props)
}

func NewContentStatDataFromProps(nid string, props notion.DatabasePageProperties) *schema.ConentStatData {
	c := &schema.ConentStatData{}
	c.DeserializePropertys(nid, props)
	return c
}
