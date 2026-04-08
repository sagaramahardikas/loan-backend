package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	loanConfig "example.com/loan/module/loan/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"

	_ "github.com/go-sql-driver/mysql"
)

type ServiceConfig struct {
	RPCHost        string         `envconfig:"HOST" required:"true"`
	DatabaseConfig DatabaseConfig `envconfig:"DB"`
}

type DatabaseConfig struct {
	Driver          string        `envconfig:"DRIVER" required:"true"`
	Host            string        `envconfig:"HOST" required:"true"`
	Port            int           `envconfig:"PORT" required:"true"`
	Username        string        `envconfig:"USERNAME" required:"true"`
	Password        string        `envconfig:"PASSWORD" required:"true"`
	Database        string        `envconfig:"DATABASE" required:"true"`
	QueryString     string        `envconfig:"QUERYSTRING" required:"true"`
	MaxOpenConns    int           `envconfig:"MAX_OPEN_CONNS" default:"30"`
	MaxIdleConns    int           `envconfig:"MAX_IDLE_CONNS" default:"15"`
	MaxIdleTime     time.Duration `envconfig:"MAX_IDLE_TIME" default:"2m"`
	MaxConnLifetime time.Duration `envconfig:"MAX_CONNECTION_LIFETIME" default:"5m"`
}

func LoadConfig() (ServiceConfig, error) {
	var cfg ServiceConfig

	// load from .env if exists
	if _, err := os.Stat(".env"); err == nil {
		if err := gotenv.Load(); err != nil {
			return cfg, err
		}
	}

	// parse environment variable to config struct using "service" namespace
	// to prevent conflict with another modules
	err := envconfig.Process("service", &cfg)
	return cfg, err
}

func LoadLoanConfig(cfg *loanConfig.LoanConfig) error {
	return envconfig.Process("loan", cfg)
}

func InitializeDatabase(cfg ServiceConfig) (*sql.DB, error) {
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
