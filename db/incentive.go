package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetIncentiveData(filter *notion.DatabaseQueryFilter) ([]schema.Incentive, error) {
	nid := "531e2bb0eebd421988457692648c6953"
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	data := []schema.Incentive{}
	for _, page := range pages {
		d := NewIncentiveDataFromPage(page)
		data = append(data, *d)
	}

	return data, nil
}

func NewIncentiveDataFromPage(page notion.Page) *schema.Incentive {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewIncentiveDataFromProps(page.ID, props)
}

func NewIncentiveDataFromProps(nid string, props notion.DatabasePageProperties) *schema.Incentive {
	data := &schema.Incentive{}
	data.DeserializePropertys(nid, props)
	return data
}
