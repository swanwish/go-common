package config

import (
	"errors"
	"strconv"

	"github.com/swanwish/go-common/utils"
)

type Configuration interface {
	Load(filePath string) error
	LoadContent(content string) error
	Get(key string) string
	Unmarshal(dest interface{}) error
}

var (
	Config Configuration = &IniConfiguration{}
)

var (
	ErrConfigurationTypeNotSpecified = errors.New("Configuration type not specified")
	ErrConfigurationFileNotExists    = errors.New("Configuration file not exists")
)

func Load(filePath string) error {
	if Config == nil {
		return ErrConfigurationTypeNotSpecified
	}
	if !utils.FileExists(filePath) {
		return ErrConfigurationFileNotExists
	}
	return Config.Load(filePath)
}

func LoadContent(content string) error {
	if Config == nil {
		return ErrConfigurationTypeNotSpecified
	}
	return Config.LoadContent(content)
}

func Get(key string) (string, error) {
	if Config == nil {
		return "", ErrConfigurationTypeNotSpecified
	}
	return Config.Get(key), nil
}

func GetInt(key string) (int64, error) {
	strValue, err := Get(key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(strValue, 10, 64)
}

func Unmarshal(dest interface{}) error {
	return Config.Unmarshal(dest)
}

//func NewConfig(configType, filePath string) (Configuration, error) {
//	var c Configuration
//	switch configType {
//	case "ini":
//		c = &IniConfiguration{}
//	}
//	if c != nil {
//		err := c.Load(filePath)
//		if err != nil {
//			return c, err
//		}
//		return c, nil
//	}
//	return nil, ErrConfigurationFileNotExists
//}
