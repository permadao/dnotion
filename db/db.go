package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/config"
)

type DB struct {
	// notion client
	DBClient       *notion.Client
	RedirectClient *notion.Client

	// db nid
	TaskDBs        []string // notion id
	WorkloadDBs    []string // notion id
	FinanceDBs     []string // notion id
	ContributorsDB string   // notion id
	GuildDB        string   // notion id
	ContentStatDB  string   // notion id
}

func New(conf *config.Config) *DB {
	return &DB{
		DBClient:       notion.NewClient(conf.NotionDB.DBSecret),
		RedirectClient: notion.NewClient(conf.NotionDB.ClientSecret),

		TaskDBs:        conf.NotionDB.TaskDBs,
		WorkloadDBs:    conf.NotionDB.WorkloadDBs,
		FinanceDBs:     conf.NotionDB.FinDBs,
		ContributorsDB: conf.NotionDB.ContributorsDB,
		GuildDB:        conf.NotionDB.GuildDB,
		ContentStatDB:  conf.NotionDB.ContentStatDB,
	}
}
