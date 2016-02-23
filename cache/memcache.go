package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strconv"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/swanwish/go-common/logs"
)

type MemcacheConfig struct {
	Servers                  string
	TimeoutMillisecond       int64
	DefaultExpirationSeconds int64
}

var (
	mc                    *memcache.Client
	Enabled               bool
	currentConfig         MemcacheConfig
	ErrMemcacheNotEnabled = errors.New("Memcached not enabled.")
)

func SetupMemcache(servers string, timeoutMillisecond, defaultExpirationSeconds int64) {
	currentConfig = MemcacheConfig{Servers: servers, TimeoutMillisecond: timeoutMillisecond, DefaultExpirationSeconds: defaultExpirationSeconds}
	if servers != "" {
		mc = memcache.New(servers)
		if timeoutMillisecond > 0 {
			mc.Timeout = time.Duration(timeoutMillisecond * int64(time.Millisecond))
		}
		Enabled = true
	} else {
		logs.Debugf("The memcache is not configured.")
		Enabled = false
	}
}

func Get(k string, data interface{}) error {
	if Enabled {
		item, err := mc.Get(k)
		if err != nil {
			if err != memcache.ErrCacheMiss {
				logs.Errorf("Failed to get memcached object for key %s with error %v\n", k, err)
			}
			return err
		}
		err = gobUnmarshal(item.Value, data)
		if err != nil {
			logs.Errorf("Failed to unmarshal data.", err)
			return err
		}
		return nil
	} else {
		return ErrMemcacheNotEnabled
	}
}

func Set(k string, v interface{}, timeout int32) error {
	if Enabled {
		value, err := gobMarshal(v)
		if err != nil {
			logs.Errorf("Failed to marshal value", v, err)
			return err
		}
		if timeout == 0 && currentConfig.DefaultExpirationSeconds > 0 {
			timeout = int32(currentConfig.DefaultExpirationSeconds)
		}
		mc.Set(&memcache.Item{Key: k, Value: value, Expiration: timeout})
		return nil
	} else {
		return ErrMemcacheNotEnabled
	}
}

func Delete(k string) error {
	if Enabled {
		err := mc.Delete(k)
		if err != nil && err != memcache.ErrCacheMiss {
			logs.Errorf("Failed to delete key %s from memcache, the eror is %v\n", k, err)
			return err
		}
		return nil
	} else {
		return ErrMemcacheNotEnabled
	}
}

func GetIntValue(key string) (int64, error) {
	if Enabled {
		item, err := mc.Get(key)
		if err != nil {
			logs.Errorf("Failed to get memcached value for key %s with error %v", key, err)
			return 0, err
		}
		value, err := strconv.ParseInt(string(item.Value), 10, 64)
		if err != nil {
			logs.Errorf("Failed to parse memcached value.", err)
			return 0, err
		}
		return value, nil
	} else {
		return 0, errors.New("Memcached not enabled.")
	}
}

func SetItem(item memcache.Item) {
	if Enabled {
		mc.Set(&item)
	}
}

func SetIntValue(key string, value int64) {
	if Enabled {
		mc.Set(&memcache.Item{Key: key, Value: []byte(strconv.FormatInt(value, 10))})
	}
}

func Increment(key string, delta uint64) (int64, error) {
	if Enabled {
		n, err := mc.Increment(key, delta)
		if err != nil {
			logs.Errorf("Failed to increment value. err: %v", err)
			return 0, err
		}
		return int64(n), nil
	}
	return 0, ErrMemcacheNotEnabled
}

func gobMarshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func gobUnmarshal(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
}
