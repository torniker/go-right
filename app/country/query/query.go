package query

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/torniker/go-right/app/country/db"
)

type Query struct {
	db sqlx.ExtContext
}

func New(db sqlx.ExtContext) *Query {
	return &Query{db: db}
}

// List returns all countries
func (q *Query) List(c context.Context) (db.Countries, error) {
	var cs db.Countries
	err := sqlx.SelectContext(c, q.db, &cs, `
		SELECT 
			code,
			name,
			region,
			subregion
		FROM countries
		ORDER BY code`)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

// ByCode returns country by 2-letter code
func (q *Query) ByCode(c context.Context, code string) (*db.Country, error) {
	var cs db.Country
	err := sqlx.GetContext(c, q.db, &cs, `
		SELECT 
			code,
			name,
			region,
			subregion
		FROM countries
		WHERE code=$1`, code)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cs, nil
}