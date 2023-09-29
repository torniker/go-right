package migration

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func Run(db sqlx.ExtContext) {
	query := `CREATE TABLE IF NOT EXISTS countries (
			code         CHAR(2) PRIMARY KEY,
			name         VARCHAR(255) NOT NULL,
			region       VARCHAR(127) NOT NULL,
			subregion    VARCHAR(127) NOT NULL
		);`
	sqlx.MustExecContext(context.Background(), db, query)
}
