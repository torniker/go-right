package db

import "github.com/torniker/go-right/app/country/response"

type (
	// Countries list of Country(s)
	Countries []Country

	// Country is struct for countrყ item
	Country struct {
		Code      string `db:"code"`
		Name      string `db:"name"`
		Region    string `db:"region"`
		Subregion string `db:"subregion"`
	}
)

// Response converts Country to response.Country
func (c Country) Response() *response.Country {
	return &response.Country{
		Name:      c.Name,
		Code:      c.Code,
		Region:    c.Region,
		Subregion: c.Subregion,
	}
}

// Response converts Countries to []response.Country
func (cs *Countries) Response() *[]*response.Country {
	var ret []*response.Country
	if cs == nil {
		return &ret
	}
	for _, country := range *cs {
		ret = append(ret, country.Response())
	}
	return &ret
}
