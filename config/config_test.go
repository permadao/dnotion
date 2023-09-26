package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	Init("config_temp")

	a := assert.New(t)
	a.Equal(Config.AppName, "dnotion")
	a.Equal(Config.Everpay.Url, "https://api.everpay.io")
	a.Equal(Config.NotionDB.BaseUrl, "https://api.notion.com")
	a.Equal(len(Config.NotionDB.FinDBs), 9)
	a.Equal(len(Config.NotionDB.TaskDBs), 7)
	a.Equal(Config.Everpay.TokenTag, "arweave,ethereum-ar-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,0x4fadc7a98f2dc96510e42dd1a74141eeae0c1543")
}
