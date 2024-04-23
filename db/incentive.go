package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetIncentiveData(filter *notion.DatabaseQueryFilter) ([]schema.Incentive, error) {
	nid := "4c19704d927f4d52b2f030ebd1648ef3"
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

func (d *DB) GetTotalIncentiveData(filter *notion.DatabaseQueryFilter) ([]schema.TotalIncentive, error) {
	nid := "04c301f8dc5448759c5919e618822854"
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	data := []schema.TotalIncentive{}
	for _, page := range pages {
		d := NewTotalIncentiveDataFromPage(page)
		data = append(data, *d)
	}

	return data, nil
}

func NewTotalIncentiveDataFromPage(page notion.Page) *schema.TotalIncentive {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewTotalIncentiveDataFromProps(page.ID, props)
}

func NewTotalIncentiveDataFromProps(nid string, props notion.DatabasePageProperties) *schema.TotalIncentive {
	data := &schema.TotalIncentive{}
	data.DeserializePropertys(nid, props)
	return data
}
