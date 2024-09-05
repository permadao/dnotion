package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetIncentiveGuildData(nid string, filter *notion.DatabaseQueryFilter) ([]schema.CIncentiveGuild, error) {
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	data := []schema.CIncentiveGuild{}
	for _, page := range pages {
		d := NewIncentiveGuildDataFromPage(page)
		data = append(data, *d)
	}

	return data, nil
}

func NewIncentiveGuildDataFromPage(page notion.Page) *schema.CIncentiveGuild {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewIncentiveGuildDataFromProps(page.ID, props)
}

func NewIncentiveGuildDataFromProps(nid string, props notion.DatabasePageProperties) *schema.CIncentiveGuild {
	data := &schema.CIncentiveGuild{}
	data.DeserializePropertys(nid, props)
	return data
}
