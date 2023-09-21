package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/config"
)

var DB *NotionDB

type NotionDB struct {
	// notion client
	DBClient       *notion.Client
	RedirectClient *notion.Client

	// db nid
	TaskDBs        []string // notion id
	WorkloadDBs    []string // notion id
	FinanceDBs     []string // notion id
	ContributorsDB string   // notion id
}

func Init() {
	DB = &NotionDB{
		DBClient:       notion.NewClient(config.Config.NotionDB.DBSecret),
		RedirectClient: notion.NewClient(config.Config.NotionDB.ClientSecret),
		TaskDBs:        config.Config.NotionDB.TaskDBs,
		WorkloadDBs:    config.Config.NotionDB.WorkloadDBs,
		FinanceDBs:     config.Config.NotionDB.FinDBs,
		ContributorsDB: config.Config.NotionDB.ContributorsDB,
	}
}
