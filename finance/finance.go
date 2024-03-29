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
