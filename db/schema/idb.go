package schema

import "github.com/dstotijn/go-notion"

type IDB interface {
	DeserializePropertys(nid string, props notion.DatabasePageProperties)
	SerializePropertys() (nid string, nprops *notion.DatabasePageProperties)
}
