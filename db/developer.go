package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetDeveloper(nid string, filter *notion.DatabaseQueryFilter) ([]schema.Developer, error) {
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	DevDatas := []schema.Developer{}
	for _, page := range pages {
		DevData := NewDevDataFromPage(page)
		DevDatas = append(DevDatas, *DevData)
	}

	return DevDatas, nil
}

func NewDevDataFromPage(page notion.Page) *schema.Developer {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewDevDataFromProps(page.ID, props)
}

func NewDevDataFromProps(nid string, props notion.DatabasePageProperties) *schema.Developer {
	data := &schema.Developer{}
	data.DeserializePropertys(nid, props)
	return data
}
