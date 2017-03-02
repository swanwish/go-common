package cache

import (
	"sync"
	"time"

	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/utils"
)

var (
	LastValueTimingCache       = lastValueTimingCache{cacheLock: &sync.Mutex{}}
	CacheExpireSeconds   int64 = 2
)

type lastValueTimingCache struct {
	cacheMap  map[string]*lastValueTimingCacheItem
	cacheLock *sync.Mutex
}

type lastValueTimingCacheItem struct {
	LastAccessTime int64
	LastValue      int64
}

func NewLastValueTimingCacheItem(lastValue int64) *lastValueTimingCacheItem {
	return &lastValueTimingCacheItem{LastAccessTime: time.Now().Unix(), LastValue: lastValue}
}

func (item *lastValueTimingCacheItem) IsExpired() bool {
	if time.Now().Unix()-item.LastAccessTime > CacheExpireSeconds {
		return true
	}
	return false
}

func AddTimingCacheItem(key string, item *lastValueTimingCacheItem) {
	LastValueTimingCache.cacheLock.Lock()
	defer LastValueTimingCache.cacheLock.Unlock()
	mapKey := utils.GetMD5Hash(key)
	logs.Debugf("Add cache item for key %s, change to map key %s", key, mapKey)
	if LastValueTimingCache.cacheMap == nil {
		LastValueTimingCache.cacheMap = make(map[string]*lastValueTimingCacheItem, 0)
	}
	LastValueTimingCache.cacheMap[mapKey] = item
}

func GetLastValue(key string) (int64, bool) {
	LastValueTimingCache.cacheLock.Lock()
	defer LastValueTimingCache.cacheLock.Unlock()
	mapKey := utils.GetMD5Hash(key)
	if item, exists := LastValueTimingCache.cacheMap[mapKey]; exists {
		if item.IsExpired() {
			delete(LastValueTimingCache.cacheMap, mapKey)
			return 0, false
		}
		return item.LastValue, true
	}
	return 0, false
}

func StartCleanExpiredItemsTask() {
	go func() {
		logs.Debugf("check duration %d", CacheExpireSeconds)
		tickChan := time.NewTicker(time.Second * time.Duration(CacheExpireSeconds)).C
		for {
			select {
			case <-tickChan:
				cleanExpiredItems()
			}
		}
	}()
}

func cleanExpiredItems() {
	LastValueTimingCache.cacheLock.Lock()
	defer LastValueTimingCache.cacheLock.Unlock()
	for key, value := range LastValueTimingCache.cacheMap {
		if value.IsExpired() {
			delete(LastValueTimingCache.cacheMap, key)
			logs.Debugf("Item with key %s has expired", key)
		}
	}
}
