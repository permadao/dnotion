package db

import (
	"github.com/dstotijn/go-notion"
)

type NotionDB struct {
	// notion client
	Client *notion.Client

	// db nid
	TaskDBs        []string // notion id
	WorkloadDBs    []string // notion id
	FinanceDBs     []string // notion id
	ContributorsDB string   // notion id
}

func New(secret string, taskDBs, workloadDBs, financeDBs []string, contributorsDB string) *NotionDB {
	return &NotionDB{
		Client:         notion.NewClient(secret),
		TaskDBs:        taskDBs,
		WorkloadDBs:    workloadDBs,
		FinanceDBs:     financeDBs,
		ContributorsDB: contributorsDB,
	}
}
