package db

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/swanwish/go-common/config"
	"github.com/swanwish/go-common/logs"
)

/**
 * Get connection get function by configure id
 */
func GetOrmConnectionGetFunc(id string) func() (*xorm.Engine, error) {
	return func() (*xorm.Engine, error) {
		driver, _ := config.Get(fmt.Sprintf("db_driver_%s", id))
		dsn, _ := config.Get(fmt.Sprintf("db_dsn_%s", id))
		if driver == "" || dsn == "" {
			logs.Errorf("The driver or dsn is not specified for id %s.", id)
			return nil, ErrConfigurationMissing
		}
		db, err := xorm.NewEngine(driver, dsn)
		return db, err
	}
}

func GetOrmConnectionWithDriverAndDsn(driver, dsn string) func() (*xorm.Engine, error) {
	return func() (*xorm.Engine, error) {
		if driver == "" || dsn == "" {
			logs.Errorf("The driver or dsn is not specified.")
			return nil, ErrConfigurationMissing
		}
		db, err := xorm.NewEngine(driver, dsn)
		return db, err
	}
}

func GetOrmDBConnection(id string) DB {
	return DefaultDB{ConnectionGetterFunc: GetConnectionGetFunc(id)}
}

func GetOrmDBConnectionWithDriverAndDsn(driver, dsn string) DB {
	return DefaultDB{ConnectionGetterFunc: GetConnectionWithDriverAndDsn(driver, dsn)}
}
