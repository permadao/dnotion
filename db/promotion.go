package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetPromotionStat(nid string, filter *notion.DatabaseQueryFilter) ([]schema.PromotionStat, error) {
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	datas := []schema.PromotionStat{}
	for _, page := range pages {
		data := NewPromotionStatDataFromPage(page)
		datas = append(datas, *data)
	}

	return datas, nil
}

func NewPromotionStatDataFromPage(page notion.Page) *schema.PromotionStat {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewPromotionStatDataFromProps(page.ID, props)
}

func NewPromotionStatDataFromProps(nid string, props notion.DatabasePageProperties) *schema.PromotionStat {
	data := &schema.PromotionStat{}
	data.DeserializePropertys(nid, props)
	return data
}

func (d *DB) GetPromotionPoints(nid string, filter *notion.DatabaseQueryFilter) ([]schema.PromotionPoints, error) {
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	datas := []schema.PromotionPoints{}
	for _, page := range pages {
		data := NewPromotionPointsDataFromPage(page)
		datas = append(datas, *data)
	}

	return datas, nil
}

func NewPromotionPointsDataFromPage(page notion.Page) *schema.PromotionPoints {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewPromotionPointsDataFromProps(page.ID, props)
}

func NewPromotionPointsDataFromProps(nid string, props notion.DatabasePageProperties) *schema.PromotionPoints {
	data := &schema.PromotionPoints{}
	data.DeserializePropertys(nid, props)
	return data
}

func (d *DB) GetPromotionSettlement(nid string, filter *notion.DatabaseQueryFilter) ([]schema.PromotionSettlement, error) {
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	datas := []schema.PromotionSettlement{}
	for _, page := range pages {
		data := NewSettlementDataFromPage(page)
		datas = append(datas, *data)
	}

	return datas, nil
}

func NewSettlementDataFromPage(page notion.Page) *schema.PromotionSettlement {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewSettlementDataFromProps(page.ID, props)
}

func NewSettlementDataFromProps(nid string, props notion.DatabasePageProperties) *schema.PromotionSettlement {
	data := &schema.PromotionSettlement{}
	data.DeserializePropertys(nid, props)
	return data
}
