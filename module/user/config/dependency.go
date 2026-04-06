package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func initializeDatabase(cfg ServiceConfig) (*sql.DB, error) {
	dbCfg := cfg.DatabaseConfig
	query := dbCfg.QueryString
	if query == "" {
		query = "parseTime=true&loc=UTC"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Database,
		query,
	)

	db, err := sql.Open(dbCfg.Driver, dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbCfg.MaxOpenConns)
	db.SetMaxIdleConns(dbCfg.MaxIdleConns)
	db.SetConnMaxIdleTime(dbCfg.MaxIdleTime)
	db.SetConnMaxLifetime(dbCfg.MaxConnLifetime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
