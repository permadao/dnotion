package dnotion

import (
	"fmt"

	"github.com/dstotijn/go-notion"
	"github.com/everFinance/go-everpay/account"
	"github.com/everFinance/go-everpay/sdk"
	"github.com/everFinance/goether"
)

type DNotion struct {
	// notion client
	Client *notion.Client
	// everpay sdk
	everpay *sdk.SDK

	// db nid
	taskDBs        []string // notion id
	workloadDBs    []string // notion id
	financeDBs     []string // notion id
	contributorsDB string   // notion id

	// contributors
	uidToNid    map[string]string //  userid -> contributors page notion id
	nidToWallet map[string]string //  contributors page notion id -> wallet
}

func New(secret, prv, everpayURL string, taskDBs, workloadDBs, financeDBs []string, contributorsDB string) *DNotion {
	signer, err := goether.NewSigner(prv)
	if err != nil {
		panic(err)
	}
	sdk, err := sdk.New(signer, everpayURL)
	if err != nil {
		panic(err)
	}
	fmt.Println("wallet address:", sdk.AccId, "everpay network:", everpayURL)

	return &DNotion{
		Client:  notion.NewClient(secret),
		everpay: sdk,

		taskDBs:        taskDBs,
		workloadDBs:    workloadDBs,
		financeDBs:     financeDBs,
		contributorsDB: contributorsDB,

		uidToNid:    make(map[string]string),
		nidToWallet: make(map[string]string),
	}
}

func (n *DNotion) InitContributors() {
	pages := n.GetAllPagesFromDB(n.contributorsDB, nil)
	for _, page := range pages {
		p := page.Properties.(notion.DatabasePageProperties)
		people := p["Notion"].People
		if len(people) == 0 {
			continue
		}
		uid := people[0].BaseUser.ID
		n.uidToNid[uid] = page.ID

		if len(p["Wallet"].RichText) == 0 {
			continue
		}
		wallet := p["Wallet"].RichText[0].PlainText
		if _, _, err := account.IDCheck(wallet); err != nil {
			continue
		}

		n.nidToWallet[page.ID] = wallet
	}
}
