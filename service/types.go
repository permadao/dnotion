package service

type UpdateFinParams struct {
	FinNid         string  `json:"nid,omitempty"`
	PaymentDateStr string  `json:"date"`
	ActualToken    string  `json:"actualtoken"`
	ActualPrice    float64 `json:"actualprice"`
	TargetToken    string  `json:"targettoken"`
	TargetPrice    float64 `json:"targetprice"`
}
