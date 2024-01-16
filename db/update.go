package db

import (
	"context"

	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) UpdatePage(idb schema.IDB) error {
	nid, props := idb.SerializePropertys()
	_, err := d.DBWriteClient.UpdatePage(context.Background(), nid, notion.UpdatePageParams{
		DatabasePageProperties: *props,
	})
	return err
}

func (d *DB) CreatePage(parentID string, idb schema.IDB) error {
	_, props := idb.SerializePropertys()
	_, err := d.DBWriteClient.CreatePage(context.Background(), notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               parentID,
		DatabasePageProperties: props,
	})
	return err
}
