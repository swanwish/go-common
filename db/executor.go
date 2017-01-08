package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/swanwish/go-common/logs"
)

const (
	ErrorMessage_GetConnectionFailed  = "Failed to get database connection"
	ErrorMessageNoConnectionProvider  = "Connection provider not specified"
	ErrorMessageNoTransactionFunction = "Transaction function not specified"
)

var (
	ErrNoConnectionProvider  = errors.New(ErrorMessageNoConnectionProvider)
	ErrNoTransactionFunction = errors.New(ErrorMessageNoTransactionFunction)
)

type SqlExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type DB interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	GetInt(query string, args ...interface{}) (int64, error)
	GetConnection() (*sqlx.DB, error)
	InTransaction(func(*sql.Tx) error) error
}

var (
	connectionMutex = &sync.Mutex{}
)

type DefaultDB struct {
	ConnectionGetterFunc func() (*sqlx.DB, error)
	pool                 *sqlx.DB
}

func (d DefaultDB) Select(dest interface{}, query string, args ...interface{}) error {
	if d.ConnectionGetterFunc == nil {
		logs.Error(ErrorMessageNoConnectionProvider)
		return ErrNoConnectionProvider
	}
	logSql(logs.LOG_LEVEL_DEBUG, query, nil, args...)
	dbConnection, err := d.ConnectionGetterFunc()
	if err != nil {
		logs.Error(ErrorMessage_GetConnectionFailed, err)
		return err
	}
	defer dbConnection.Close()

	err = dbConnection.Select(dest, query, args...)
	if err != nil {
		logSql(logs.LOG_LEVEL_ERROR, query, err, args...)
	}
	return err
}

func (d DefaultDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	dbConnection, err := d.GetConnection()
	if err != nil {
		logs.Error(ErrorMessage_GetConnectionFailed, err)
		return nil, err
	}
	defer dbConnection.Close()

	logSql(logs.LOG_LEVEL_DEBUG, query, nil, args...)
	result, err := dbConnection.Exec(query, args...)
	if err != nil {
		logSql(logs.LOG_LEVEL_ERROR, query, err, args...)
	}
	return result, err
}

func (d DefaultDB) Get(dest interface{}, query string, args ...interface{}) error {
	dbConnection, err := d.GetConnection()
	if err != nil {
		logs.Error(ErrorMessage_GetConnectionFailed, err)
		return err
	}
	defer dbConnection.Close()

	logSql(logs.LOG_LEVEL_DEBUG, query, nil, args...)
	err = dbConnection.Get(dest, query, args...)
	if err != nil {
		logSql(logs.LOG_LEVEL_ERROR, query, err, args...)
	}
	return err
}

func (d DefaultDB) GetInt(query string, args ...interface{}) (int64, error) {
	logSql(logs.LOG_LEVEL_DEBUG, query, nil, args)
	var maxValue int64
	err := d.Get(&maxValue, query, args...)
	if err != nil {
		logSql(logs.LOG_LEVEL_ERROR, query, err, args...)
	}
	return maxValue, err
}

func (d DefaultDB) GetConnection() (*sqlx.DB, error) {
	connectionMutex.Lock()
	defer connectionMutex.Unlock()

	if d.pool != nil {
		return d.pool, nil
	}

	if d.ConnectionGetterFunc == nil {
		logs.Errorf(ErrorMessageNoConnectionProvider)
		return nil, ErrNoConnectionProvider
	}
	pool, err := d.ConnectionGetterFunc()
	if err != nil {
		logs.Errorf("Failed to get database connection, the error is %v", err)
		return nil, err
	}
	d.pool = pool
	return d.pool, nil
}

func (d DefaultDB) InTransaction(TransactionFunc func(*sql.Tx) error) error {
	if TransactionFunc == nil {
		logs.Errorf(ErrorMessageNoTransactionFunction)
		return ErrNoTransactionFunction
	}

	if d.ConnectionGetterFunc == nil {
		logs.Errorf(ErrorMessageNoConnectionProvider)
		return ErrNoConnectionProvider
	}

	dbConnection, err := d.GetConnection()
	if err != nil {
		logs.Errorf(ErrorMessage_GetConnectionFailed, err)
		return err
	}
	defer dbConnection.Close()

	tx, err := dbConnection.Begin()
	err = TransactionFunc(tx)

	if err != nil {
		logs.Errorf("Failed to execute transaction function, the error is %v", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()

}

func logSql(level int64, sql string, err error, params ...interface{}) {
	if logs.MaxLogLevel <= level {
		cons := []string{}
		for _, param := range params {
			cons = append(cons, fmt.Sprintf("%v", param))
		}
		format := "Execute sql:%s\nParams:`%s`"
		logParams := []interface{}{sql, strings.Join(cons, "`, `")}
		if err != nil {
			format = format + " failed, the error is %v"
			logParams = append(logParams, err)
		}
		logs.Logf(level, logs.EMPTY_KEY, format, logParams...)
	}
}
