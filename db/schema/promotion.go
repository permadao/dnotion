package schema

import "github.com/dstotijn/go-notion"

type PromotionStat struct {
	NID string
}

type PromotionPoints struct {
	NID        string
	BasePoints string
}

type Settlement struct {
	NID           string
	Contributor   string
	TotalScore    float64
	PersonalScore float64
	Rewards       float64
	Date          string
}

func (f *PromotionStat) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	f.NID = nid
}

func (f *PromotionStat) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	return f.NID, &props
}

func (f *PromotionPoints) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	f.NID = nid
}

func (f *PromotionPoints) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	return f.NID, &props
}
func (f *Settlement) DeserializePropertys(nid string, props notion.DatabasePageProperties) {
	f.NID = nid
}

func (f *Settlement) SerializePropertys() (nid string, nprops *notion.DatabasePageProperties) {
	props := notion.DatabasePageProperties{}
	return f.NID, &props
}
