package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config := New("config_temp")

	a := assert.New(t)
	a.Equal(config.AppName, "dnotion")
	a.Equal(config.Everpay.Url, "https://api.everpay.io")
	a.Equal(config.NotionDB.BaseUrl, "https://api.notion.com")
	a.Equal(len(config.NotionDB.FinDBs), 9)
	a.Equal(len(config.NotionDB.TaskDBs), 7)
	a.Equal(config.Everpay.TokenTag, "arweave,ethereum-ar-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,0x4fadc7a98f2dc96510e42dd1a74141eeae0c1543")
}
