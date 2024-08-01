package schema

type Contributor struct {
	Name   string
	Amount float64
	Wallet string
}

type StatResult struct {
	TotalAmount       float64
	Contributors      map[string]float64
	RankOfContributor []Contributor
}

// ResultSepToken 按token分开的激励结果
type ResultSepToken map[string]map[string]float64
