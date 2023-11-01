package db

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dstotijn/go-notion"
)

func (d *DB) GetPageFromDBByID(nid, id string) (*notion.Page, error) {
	res, err := d.DBClient.QueryDatabase(context.Background(), nid, &notion.DatabaseQuery{
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

func (d *DB) GetLastPageFromDB(nid string) (page notion.Page, err error) {
	res, err := d.DBClient.QueryDatabase(context.Background(), nid, &notion.DatabaseQuery{
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
		fmt.Println("get last page error: ", err)
		return
	}
	if len(res.Results) < 1 {
		err = fmt.Errorf("get last page failed, nid: %s", nid)
		return
	}

	return res.Results[0], nil
}

func (d *DB) GetLastIDFromDB(nid string) (id int, err error) {
	page, err := d.GetLastPageFromDB(nid)
	if err != nil {
		return 0, err
	}

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

func (d *DB) GetAllPagesFromDB(nid string, filter *notion.DatabaseQueryFilter) (pages []notion.Page, err error) {
	hasMore := true
	nextCursor := ""

	for hasMore {
		var res notion.DatabaseQueryResponse
		res, err = d.DBClient.QueryDatabase(context.Background(), nid, &notion.DatabaseQuery{
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
			fmt.Println("query database error: ", err)
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

func (d *DB) GetCountFromDB(nid string) (int, error) {
	pages, err := d.GetAllPagesFromDB(nid, nil)
	if err != nil {
		return 0, err
	}
	return len(pages), nil
}
