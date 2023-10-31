package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/everFinance/go-everpay/account"
)

// TODO: need refactor to contractor schema
func (d *DB) GetAllContributors() (uidToNid, nidToWallet map[string]string) {
	uidToNid = map[string]string{}
	nidToWallet = map[string]string{}

	pages, _ := d.GetAllPagesFromDB(d.ContributorsDB, nil)
	for _, page := range pages {
		p := page.Properties.(notion.DatabasePageProperties)
		people := p["Notion"].People
		if len(people) == 0 {
			continue
		}
		uid := people[0].BaseUser.ID
		uidToNid[uid] = page.ID

		if len(p["Wallet"].RichText) == 0 {
			continue
		}
		wallet := p["Wallet"].RichText[0].PlainText
		if _, _, err := account.IDCheck(wallet); err != nil {
			continue
		}

		nidToWallet[page.ID] = wallet
	}

	return
}
