package command

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/torniker/go-right/app/country/db"
)

type Command struct {
	db sqlx.ExtContext
}

func New(db sqlx.ExtContext) *Command {
	return &Command{db: db}
}

// Upsert creates or updates country
func (cmd *Command) Upsert(c context.Context, country db.Country) error {
	_, err := sqlx.NamedExecContext(c, cmd.db, `
	INSERT INTO countries
		(code
		,name
		,region
		,subregion)
	VALUES
		(:code
		,:name
		,:region
		,:subregion)
	ON CONFLICT (code) DO
	UPDATE SET
		name = :name
		,region = :region
		,subregion = :subregion
	`, country)
	if err != nil {
		return err
	}
	return nil
}
