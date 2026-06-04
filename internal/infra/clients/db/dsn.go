package db

import (
	"fmt"
	"net/url"

	"github.com/viictormotorhead/financial-app/internal/infra/config"
)

func DSN() string {
	dbCfg := config.Config().DB

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbCfg.Host,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Name,
		dbCfg.Port,
		dbCfg.SSLMode,
	)
}

func MigrationURL() string {
	dbCfg := config.Config().DB

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(dbCfg.User),
		url.QueryEscape(dbCfg.Password),
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
		url.QueryEscape(dbCfg.SSLMode),
	)
}
