package cache

import (
	"errors"
	"sync"

	"github.com/swanwish/go-common/logs"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	ErrParameterMissing    = errors.New("Parameter is missing")
	DefaultLevelDbFileName = "leveldb"
	mutex                  = &sync.Mutex{}
)

type CacheLevelDB struct {
	LevelDbFileName string
}

func (cache CacheLevelDB) GetLevelDB() (*leveldb.DB, error) {
	dbFileName := DefaultLevelDbFileName
	if cache.LevelDbFileName != "" {
		dbFileName = cache.LevelDbFileName
	}
	return leveldb.OpenFile(dbFileName, nil)
}

func (cache CacheLevelDB) SetValue(key, value []byte) error {
	if len(key) == 0 {
		logs.Error("The key is empty.")
		return ErrParameterMissing
	}
	mutex.Lock()
	defer mutex.Unlock()
	db, err := cache.GetLevelDB()
	if err != nil {
		logs.Errorf("Faieled to get level db connection, the error is %v\n", err)
		return err
	}
	defer db.Close()

	return db.Put(key, value, nil)
}

func (cache CacheLevelDB) GetValue(key []byte) ([]byte, error) {
	mutex.Lock()
	defer mutex.Unlock()
	db, err := cache.GetLevelDB()
	if err != nil {
		logs.Errorf("Failed to get level db connection, the error is %v\n", err)
		return []byte{}, err
	}
	defer db.Close()

	value, err := db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		err = ErrKeyNotFound
	}
	return value, err
}

func (c CacheLevelDB) Delete(key []byte) error {
	mutex.Lock()
	defer mutex.Unlock()
	db, err := c.GetLevelDB()
	if err != nil {
		logs.Errorf("Failed to get db connection, the error is %v\n", err)
		return err
	}
	defer db.Close()

	return db.Delete(key, nil)
}
