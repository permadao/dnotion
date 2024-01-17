package translation

import (
	"github.com/permadao/dnotion/db"
	trschema "github.com/permadao/dnotion/translation/schema"
)

type Translator struct {
	db        *db.DB
	nidToName map[string]string
}

var (
	SUPREME = trschema.Level{Name: "Supreme-至尊", Color: "purple"}
	GLORY   = trschema.Level{Name: "Glory-荣耀", Color: "green"}
	DIAMOND = trschema.Level{Name: "Diamond-钻石", Color: "green"}
	GOLD    = trschema.Level{Name: "Gold-黄金", Color: "yellow"}
	SILVER  = trschema.Level{Name: "Silver-白银", Color: "grey"}
)

func New(db *db.DB) *Translator {
	tr := &Translator{
		db:        db,
		nidToName: map[string]string{},
	}
	tr.initContributors()
	return tr
}

func (t *Translator) initContributors() {
	contributors, err := t.db.GetContributors(nil)
	if err != nil {
		panic(err)
	}

	for _, c := range contributors {
		if c.NotionName != "" {
			t.nidToName[c.NID] = c.NotionName
		}
	}
}

func (t *Translator) GetTierSlice() []trschema.Tier {
	ts := []trschema.Tier{}
	ts = append(ts, []trschema.Tier{
		{
			Level:    SUPREME,
			Interval: [2]int{0, 1},
		},
		{
			Level:    GLORY,
			Interval: [2]int{1, 2},
		},
		{
			Level:    DIAMOND,
			Interval: [2]int{2, 4},
		},
		{
			Level:    GOLD,
			Interval: [2]int{4, 7},
		},
		{
			Level:    SILVER,
			Interval: [2]int{7, 10},
		},
	}...)
	return ts
}
