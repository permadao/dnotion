package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/config"
)

type DB struct {
	// notion client
	DBReadClient  *notion.Client
	DBWriteClient *notion.Client

	// db nid
	TaskDBs                 []string // notion id
	WorkloadDBs             []string // notion id
	FinanceDBs              []string // notion id
	ContributorsDB          string   // notion id
	GuildDB                 string   // notion id
	ContentStatDB           string   // notion id
	CincentiveWeeklyDB      string   // notion id
	CincentiveWeeklyGuildDB string   // notion id
}

func New(conf *config.Config) *DB {
	return &DB{
		DBWriteClient: notion.NewClient(conf.NotionDB.WriteSecret),
		DBReadClient:  notion.NewClient(conf.NotionDB.ReadSecret),

		TaskDBs:                 conf.NotionDB.TaskDBs,
		WorkloadDBs:             conf.NotionDB.WorkloadDBs,
		FinanceDBs:              conf.NotionDB.FinDBs,
		ContributorsDB:          conf.NotionDB.ContributorsDB,
		GuildDB:                 conf.NotionDB.GuildDB,
		ContentStatDB:           conf.NotionDB.ContentStatDB,
		CincentiveWeeklyDB:      conf.NotionDB.CincentiveWeeklyDB,
		CincentiveWeeklyGuildDB: conf.NotionDB.CincentiveWeeklyGuildDB,
	}
}
