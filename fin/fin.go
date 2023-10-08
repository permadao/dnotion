package fin

import (
	"github.com/dstotijn/go-notion"
	"github.com/everFinance/go-everpay/account"
	"github.com/everFinance/go-everpay/sdk"
	"github.com/everFinance/goether"
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	log "github.com/sirupsen/logrus"
)

var Fin *Finance

type Finance struct {
	// everpay sdk
	everpay *sdk.SDK

	// contributors
	uidToNid    map[string]string //  userid -> contributors page notion id
	nidToWallet map[string]string //  contributors page notion id -> wallet
}

func Init() {
	// create finance
	Fin = create(config.Config.Everpay.PrivKey, config.Config.Everpay.Url)
	// init contributors
	Fin.initContributors()
}

func create(prv, everpayURL string) *Finance {
	signer, err := goether.NewSigner(prv)
	if err != nil {
		panic(err)
	}
	sdk, err := sdk.New(signer, everpayURL)
	if err != nil {
		panic(err)
	}
	log.Info("wallet address:", sdk.AccId, "everpay network:", everpayURL)

	return &Finance{
		everpay:     sdk,
		uidToNid:    make(map[string]string),
		nidToWallet: make(map[string]string),
	}
}

func (f *Finance) initContributors() {
	pages, _ := db.DB.GetAllPagesFromDB(db.DB.ContributorsDB, nil)
	for _, page := range pages {
		p := page.Properties.(notion.DatabasePageProperties)
		people := p["Notion"].People
		if len(people) == 0 {
			continue
		}
		uid := people[0].BaseUser.ID
		f.uidToNid[uid] = page.ID

		if len(p["Wallet"].RichText) == 0 {
			continue
		}
		wallet := p["Wallet"].RichText[0].PlainText
		if _, _, err := account.IDCheck(wallet); err != nil {
			continue
		}

		f.nidToWallet[page.ID] = wallet
	}
}
