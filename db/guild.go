package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetAllGuilds() ([]schema.GuildData, error) {
	pages, err := d.GetAllPagesFromDB(d.GuildDB, nil)
	if err != nil {
		return nil, err
	}

	achiementDatas := []schema.GuildData{}
	for _, page := range pages {
		achiementData := NewGuildDataFromPage(page)
		achiementDatas = append(achiementDatas, *achiementData)
	}

	return achiementDatas, nil
}

func NewGuildDataFromPage(page notion.Page) *schema.GuildData {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewGuildDataFromProps(page.ID, props)
}

func NewGuildDataFromProps(nid string, props notion.DatabasePageProperties) *schema.GuildData {
	contributorsData := &schema.GuildData{}
	contributorsData.DeserializePropertys(nid, props)
	return contributorsData
}
