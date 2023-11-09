# dnotion

Below is a demonstration code for PermaDAO financial settlement:

## Example

### Finance

```golang
package main

import (
	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/finance"
)

func main() {
	config := config.New("../config.toml")
	db := db.New(config)
	fin := finance.New(config, db)

	// 1. check counts
	fin.CheckAllDbsCountAndID()

	// 2. check actual usd equality with workload usd
	// dao.CheckAllWorkloadAndAcutal()

	// 3. Update workload tx to finance table
	// dao.UpdateAllWorkToFin()

	// 4. Update all Finance transactions to progress
	// dao.UpdateAllFinToProgress("2023-07-21", "AR", 6.11)
	// dao.UpdateFinToProgress("f5d84582e9d8471f8f903563cce72567", "2022-07-22", "AR", 5.9)

	// 5. process payment by everpay
	// dao.PayAll()
	// dao.Pay("caac7a1aefcc4ed0b02b8adbc106f021")
}
```

### Achievement

```golang
package main

import (
	"fmt"

	"github.com/permadao/dnotion/config"
	"github.com/permadao/dnotion/db"
	"github.com/permadao/dnotion/guild"
)

func main() {
	config := config.New("./config.toml")
	db := db.New(config)
	g := guild.New(config, db)

	totalAmount, contributors, rank, _ := g.StatFinance("AR", "a815dcd96395424a93d9854b4418ab03")
	fmt.Println(totalAmount, contributors, rank)
	// g.GenGuilds("AR", "2023-10-27")
	// g.GenGuilds("AR", "2023-11-03")
}
```