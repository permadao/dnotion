package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetTranslator(nid string, filter *notion.DatabaseQueryFilter) ([]schema.Translator, error) {
	pages, err := d.GetPages(nid, filter)
	if err != nil {
		return nil, err
	}

	TRDatas := []schema.Translator{}
	for _, page := range pages {
		TRData := NewTRDataFromPage(page)
		TRDatas = append(TRDatas, *TRData)
	}

	return TRDatas, nil
}

func NewTRDataFromPage(page notion.Page) *schema.Translator {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewTRDataFromProps(page.ID, props)
}

func NewTRDataFromProps(nid string, props notion.DatabasePageProperties) *schema.Translator {
	data := &schema.Translator{}
	data.DeserializePropertys(nid, props)
	return data
}
