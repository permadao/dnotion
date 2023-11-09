package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetGuilds(filter *notion.DatabaseQueryFilter) ([]schema.GuildData, error) {
	pages, err := d.GetPages(d.GuildDB, nil)
	if err != nil {
		return nil, err
	}

	guilds := []schema.GuildData{}
	for _, page := range pages {
		guild := NewGuildDataFromPage(page)
		guilds = append(guilds, *guild)
	}

	return guilds, nil
}

func NewGuildDataFromPage(page notion.Page) *schema.GuildData {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewGuildDataFromProps(page.ID, props)
}

func NewGuildDataFromProps(nid string, props notion.DatabasePageProperties) *schema.GuildData {
	guild := &schema.GuildData{}
	guild.DeserializePropertys(nid, props)
	return guild
}
