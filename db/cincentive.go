package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetCIncentiveData(nid string, filter *notion.DatabaseQueryFilter) ([]schema.CIncentive, error) {
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	data := []schema.CIncentive{}
	for _, page := range pages {
		d := NewCIncentiveDataFromPage(page)
		data = append(data, *d)
	}

	return data, nil
}

func NewCIncentiveDataFromPage(page notion.Page) *schema.CIncentive {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewTCIncentiveDataFromProps(page.ID, props)
}

func NewTCIncentiveDataFromProps(nid string, props notion.DatabasePageProperties) *schema.CIncentive {
	data := &schema.CIncentive{}
	data.DeserializePropertys(nid, props)
	return data
}
