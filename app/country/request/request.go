package request

import (
	"github.com/torniker/go-right/app/country/db"
	"github.com/torniker/go-right/pkg/server"
)

type (
	// CountrySave request to store country's main info
	CountrySave struct {
		Code      string `json:"alpha2Code"`
		Name      string `json:"name"`
		Region    string `json:"region"`
		Subregion string `json:"subregion"`
	}

	validSaveRequest struct {
		CountrySave
	}
)

func (cr *CountrySave) Validate() (*validSaveRequest, error) {
	if cr.Name == "" || len(cr.Name) < 2 || len(cr.Name) > 100 {
		return nil, server.ErrBadRequest("bad country name")
	}
	return &validSaveRequest{
		CountrySave: *cr,
	}, nil
}

// DB converts request to db object
func (sr CountrySave) DB() db.Country {
	return db.Country{
		Code:      sr.Code,
		Name:      sr.Name,
		Region:    sr.Region,
		Subregion: sr.Subregion,
	}
}
