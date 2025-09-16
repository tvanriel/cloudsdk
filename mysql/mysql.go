package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
)

type NewMySQLParams struct {
	fx.In

	Configuration Configuration
}

func NewMySQL(opts NewMySQLParams) (*sql.DB, error) {
	setDefaultOptsIfUnset(&opts.Configuration)

	conn, err := sql.Open("mysql", makeDSN(opts.Configuration.Host, opts.Configuration.User, opts.Configuration.Password, opts.Configuration.DBName, opts.Configuration.Opts))
	if err != nil {
		return nil, fmt.Errorf("mysql: cannot open db: %w", err)
	}

	return conn, nil
}
