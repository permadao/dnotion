package translation

import "github.com/permadao/dnotion/db"

type Translator struct {
	db        *db.DB
	nidToName map[string]string
}

type Tier struct {
	level    Level
	interval [2]int
	val      [2]int
}

type Level struct {
	name  string
	color string
}

var (
	SUPREME = Level{name: "Supreme-至尊", color: "purple"}
	GLORY   = Level{name: "Glory-荣耀", color: "green"}
	DIAMOND = Level{"Diamond-钻石", "green"}
	GOLD    = Level{"Gold-黄金", "yellow"}
	SILVER  = Level{"Silver-白银", "grey"}
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

func (t *Translator) GetTierSlice() []Tier {
	ts := []Tier{}
	ts = append(ts, []Tier{
		{
			level:    SUPREME,
			interval: [2]int{0, 1},
		},
		{
			level:    GLORY,
			interval: [2]int{1, 2},
		},
		{
			level:    DIAMOND,
			interval: [2]int{2, 4},
		},
		{
			level:    GOLD,
			interval: [2]int{4, 7},
		},
		{
			level:    SILVER,
			interval: [2]int{7, 10},
		},
	}...)
	return ts
}
