package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dstotijn/go-notion"
)

func (db *NotionDB) GetPageFromDBByID(nid, id string) (*notion.Page, error) {
	res, err := db.Client.QueryDatabase(context.Background(), nid, &notion.DatabaseQuery{
		Filter: &notion.DatabaseQueryFilter{
			Property: "ID",
			DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
				Title: &notion.TextPropertyFilter{
					Equals: id,
				},
			},
		},
		PageSize: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("get page with id %s error: %s", id, err)
	}
	if len(res.Results) < 1 {
		return nil, fmt.Errorf("no page found with id %s", id)
	}

	return &res.Results[0], nil
}

func (db *NotionDB) GetLastPageFromDB(nid string) (page notion.Page) {
	res, err := db.Client.QueryDatabase(context.Background(), nid, &notion.DatabaseQuery{
		Sorts: []notion.DatabaseQuerySort{
			notion.DatabaseQuerySort{
				Property:  "Sort ID",
				Direction: notion.SortDirDesc,
				Timestamp: notion.SortTimeStampCreatedTime,
			},
		},
		PageSize: 1,
	})
	if err != nil {
		fmt.Println("err", err)
		return
	}
	if len(res.Results) < 1 {
		return
	}

	return res.Results[0]
}

func (db *NotionDB) GetLastIDFromDB(nid string) (id int) {
	page := db.GetLastPageFromDB(nid)

	var err error
	idd := page.Properties.(notion.DatabasePageProperties)["ID"]
	switch idd.Type {
	case notion.DBPropTypeTitle:
		id, err = strconv.Atoi(idd.Title[0].PlainText)
		if err != nil {
			fmt.Println("err", err)
		}
	case notion.DBPropTypeNumber:
		id = int(*idd.Number)
	}

	return
}

func (n *NotionDB) GetAllPagesFromDB(nid string, filter *notion.DatabaseQueryFilter) (pages []notion.Page) {
	hasMore := true
	nextCursor := ""

	for hasMore {
		res, err := n.Client.QueryDatabase(context.Background(), nid, &notion.DatabaseQuery{
			Filter: filter,
			Sorts: []notion.DatabaseQuerySort{
				notion.DatabaseQuerySort{
					Property:  "Sort ID",
					Direction: notion.SortDirDesc,
					Timestamp: notion.SortTimeStampLastEditedTime,
				},
			},
			StartCursor: nextCursor,
		})
		if err != nil {
			fmt.Println("err", err)
			return
		}

		// append
		pages = append(pages, res.Results...)

		hasMore = res.HasMore
		if hasMore {
			nextCursor = *res.NextCursor
		}
	}
	return
}

func (n *NotionDB) GetCountFromDB(nid string) int {
	return len(n.GetAllPagesFromDB(nid, nil))
}
