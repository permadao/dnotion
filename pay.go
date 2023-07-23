package dnotion

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/dstotijn/go-notion"
)

func (n *DNotion) PayAll() {
	for _, v := range n.financeDBs {
		t := time.Now()
		fmt.Println("Paying, fid", v)

		n.Pay(v)

		fmt.Printf("Finance payment, %s updated, since: %v\n\n", v, time.Since(t))
	}
}

func (n *DNotion) Pay(fnid string) {
	// get Status is In progress
	pages := n.GetAllPagesFromDB(fnid, &notion.DatabaseQueryFilter{
		And: []notion.DatabaseQueryFilter{
			notion.DatabaseQueryFilter{
				Property: "Status",
				DatabaseQueryPropertyFilter: notion.DatabaseQueryPropertyFilter{
					Status: &notion.StatusDatabaseQueryFilter{
						Equals: "In progress",
					},
				},
			},
		},
	})

	// for payments
	for _, page := range pages {
		// get wallet and amount
		p := page.Properties.(notion.DatabasePageProperties)
		token := 0.0
		if p["Token"].Formula != nil {
			token = *p["Token"].Formula.Number
		}

		wallet := ""
		if len(p["Contributor"].Relation) > 0 {
			if v, ok := n.nidToWallet[p["Contributor"].Relation[0].ID]; ok {
				wallet = v
			}
		}

		// update to done
		if _, err := n.Client.UpdatePage(context.Background(), page.ID,
			notion.UpdatePageParams{
				DatabasePageProperties: notion.DatabasePageProperties{
					"Status": notion.DatabasePageProperty{
						Status: &notion.SelectOptions{Name: "Done"},
					},
				},
			}); err != nil {
			fmt.Printf("Update nid/id: %v/%v to `done` failed. %v\n", fnid, page.ID, err)
			continue
		}

		// payment
		tx, err := n.everpay.Transfer(
			"arweave,ethereum-ar-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,0x4fadc7a98f2dc96510e42dd1a74141eeae0c1543",
			FloatToBigInt(token), wallet,
			`{"appName": "permadao-payroll", "permadaoUrl": "`+page.URL+`"}`)
		if err != nil {
			fmt.Printf("Payment failed nid/id: %v/%v. %v\n", fnid, page.ID, err)
			// rollback
			if _, err := n.Client.UpdatePage(context.Background(), page.ID,
				notion.UpdatePageParams{
					DatabasePageProperties: notion.DatabasePageProperties{
						"Status": notion.DatabasePageProperty{
							Status: &notion.SelectOptions{Name: "In progress"},
						},
					},
				}); err != nil {
				fmt.Printf("Update nid/id: %v/%v to `done` failed. %v\n", fnid, page.ID, err)
			}
			continue
		}

		// update receipt
		receipt := "https://scan.everpay.io/tx/" + tx.HexHash()
		if _, err := n.Client.UpdatePage(context.Background(), page.ID,
			notion.UpdatePageParams{
				DatabasePageProperties: notion.DatabasePageProperties{
					"Receipt(url)": notion.DatabasePageProperty{
						URL: &receipt,
					},
				},
			}); err != nil {
			fmt.Printf("Update nid/id: %v/%v to `done` failed. %v\n", fnid, page.ID, err)
		}

	}
}

func FloatToBigInt(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000))

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result)

	return result
}
