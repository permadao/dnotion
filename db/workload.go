package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func NewWrokloadDataFromProps(nid string, props *notion.DatabasePageProperties) *schema.WorkloadData {
	workloadData := &schema.WorkloadData{}
	workloadData.DeserializePropertys(nid, *props)
	return workloadData
}
