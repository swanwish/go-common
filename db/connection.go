package db

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/swanwish/go-common/config"
	"github.com/swanwish/go-common/logs"
)

var (
	ErrConfigurationMissing = errors.New("Configuration is missiong")
)

/**
 * Get connection get function by configure id
 */
func GetConnectionGetFunc(id string) func() (string, *sqlx.DB, error) {
	return func() (string, *sqlx.DB, error) {
		driver, _ := config.Get(fmt.Sprintf("db_driver_%s", id))
		dsn, _ := config.Get(fmt.Sprintf("db_dsn_%s", id))
		if driver == "" || dsn == "" {
			logs.Errorf("The driver or dsn is not specified for id %s.", id)
			return "", nil, ErrConfigurationMissing
		}
		db, err := sqlx.Open(driver, dsn)
		return driver, db, err
	}
}

func GetConnectionWithDriverAndDsn(driver, dsn string) func() (string, *sqlx.DB, error) {
	return func() (string, *sqlx.DB, error) {
		if driver == "" || dsn == "" {
			logs.Errorf("The driver or dsn is not specified.")
			return "", nil, ErrConfigurationMissing
		}
		db, err := sqlx.Open(driver, dsn)
		return driver, db, err
	}
}

func GetDBConnection(id string) DB {
	return DefaultDB{ConnectionGetterFunc: GetConnectionGetFunc(id)}
}

func GetDBConnectionWithDriverAndDsn(driver, dsn string) DB {
	return DefaultDB{ConnectionGetterFunc: GetConnectionWithDriverAndDsn(driver, dsn)}
}
