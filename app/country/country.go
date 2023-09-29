package country

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/torniker/go-right/app/country/command"
	"github.com/torniker/go-right/app/country/migration"
	"github.com/torniker/go-right/app/country/query"
	"github.com/torniker/go-right/app/country/request"
	"github.com/torniker/go-right/app/country/response"
	"github.com/torniker/go-right/pkg/server"
)

type Service struct {
	cmd   *command.Command
	query *query.Query
}

func New(db sqlx.ExtContext) *Service {
	migration.Run(db)
	return &Service{
		cmd:   command.New(db),
		query: query.New(db),
	}
}

// List returns countries for page
func (s *Service) List(c context.Context) (*[]*response.Country, error) {
	countries, err := s.query.List(c)
	if err != nil {
		return nil, server.ErrInternal(err)
	}
	return countries.Response(), nil
}

// ByCode returns country by code
func (s *Service) ByCode(c context.Context, code string) (*response.Country, error) {
	cdb, err := s.query.ByCode(c, code)
	if err != nil {
		return nil, server.ErrInternal(err)
	}
	if cdb == nil {
		return nil, server.ErrBadRequest("bad country id")
	}
	return cdb.Response(), nil
}

// Save inserts or updates countries
func (s *Service) Save(c context.Context, csr request.CountrySave) error {
	err := s.cmd.Upsert(c, csr.DB())
	if err != nil {
		return server.ErrInternal(err)
	}
	return err
}
