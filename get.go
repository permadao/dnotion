package dnotion

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dstotijn/go-notion"
)

func (n *DNotion) GetPageFromDBByID(nid, id string) (page notion.Page) {
	res, err := n.Client.QueryDatabase(context.Background(), nid, &notion.DatabaseQuery{
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
		fmt.Println("err", err)
		return
	}
	if len(res.Results) < 1 {
		return
	}

	return res.Results[0]
}

func (n *DNotion) GetLastIDFromDB(nid string) (id int) {
	page := n.GetLastPageFromDB(nid)

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

func (n *DNotion) GetLastPageFromDB(nid string) (page notion.Page) {
	res, err := n.Client.QueryDatabase(context.Background(), nid, &notion.DatabaseQuery{
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

func (n *DNotion) GetCountFromDB(nid string) int {
	return len(n.GetAllPagesFromDB(nid, nil))
}

func (n *DNotion) GetAllPagesFromDB(nid string, filter *notion.DatabaseQueryFilter) (pages []notion.Page) {
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
