package finance

import (
	"github.com/everFinance/goether"
	"github.com/everVision/everpay-kits/sdk"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/logger"
)

var log = logger.New("finance")

type Finance struct {
	// everpay sdk
	everpay *sdk.SDK

	// db
	db *db.DB

	// contributors
	uidToNid    map[string]string //  userid -> contributors page notion id
	nidToWallet map[string]string //  contributors page notion id -> wallet
	tokenTagMap map[string]string //  token tags
}

func New(conf *config.Config, db *db.DB) *Finance {
	signer, err := goether.NewSigner(conf.Everpay.PrivKey)
	if err != nil {
		panic(err)
	}
	sdk, err := sdk.New(signer, conf.Everpay.Url)
	if err != nil {
		panic(err)
	}
	log.Info("wallet address:", sdk.AccId, "everpay network:", conf.Everpay.Url)

	fin := &Finance{
		everpay: sdk,

		db: db,

		uidToNid:    make(map[string]string),
		nidToWallet: make(map[string]string),
		tokenTagMap: map[string]string{
			"AR":   "arweave,ethereum-ar-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,0x4fadc7a98f2dc96510e42dd1a74141eeae0c1543",
			"BP":   "aostest-bp-_HbnZH5blAZH0CNT1k_dpRrGXWCzBg34hjMUkoDrXr0",
			"USDC": "ethereum-usdc-0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		},
	}

	fin.initContributors()
	return fin
}

func (f *Finance) initContributors() {
	contributors, err := f.db.GetContributors(nil)
	if err != nil {
		panic(err)
	}

	for _, c := range contributors {
		f.uidToNid[c.NotionID] = c.NID
		f.nidToWallet[c.NID] = c.Wallet
	}
}
