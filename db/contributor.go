package db

import (
	"github.com/dstotijn/go-notion"
	"github.com/permadao/dnotion/db/schema"
)

func (d *DB) GetAllContributors() ([]schema.ContributorData, error) {
	pages, err := d.GetAllPagesFromDB(d.ContributorsDB, nil)
	if err != nil {
		return nil, err
	}

	contributorDatas := []schema.ContributorData{}
	for _, page := range pages {
		contributorData := NewContributorDataFromPage(page)
		contributorDatas = append(contributorDatas, *contributorData)
	}

	return contributorDatas, nil
}

func NewContributorDataFromPage(page notion.Page) *schema.ContributorData {
	props := page.Properties.(notion.DatabasePageProperties)
	return NewContributorDataFromProps(page.ID, props)
}

func NewContributorDataFromProps(nid string, props notion.DatabasePageProperties) *schema.ContributorData {
	contributorsData := &schema.ContributorData{}
	contributorsData.DeserializePropertys(nid, props)
	return contributorsData
}
