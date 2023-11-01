package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetFinancesByParentID(parentID string) {}

func NewFinDataFromProps(nid string, props *notion.DatabasePageProperties) *schema.FinData {
	finData := &schema.FinData{}
	finData.DeserializePropertys(nid, *props)
	return finData
}
