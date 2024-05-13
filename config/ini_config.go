package config

import (
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/swanwish/go-common/logs"
)

type IniConfiguration struct {
	configurationKVM map[string]string
}

func (c *IniConfiguration) Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		logs.Errorf("Failed to read configuration file.", err)
		return err
	}
	defer file.Close()
	logs.Infof("Read configure from file %s", filePath)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		logs.Errorf("Failed to read configuration file.", err)
		return err
	}

	return c.LoadContent(string(fileContent))
}

func (c *IniConfiguration) LoadContent(content string) error {
	c.configurationKVM = make(map[string]string, 0)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if strings.Index(line, "#") == 0 {
			continue
		}
		if strings.Index(line, ";") == 0 {
			continue
		}
		pair := strings.Split(line, "=")
		if len(pair) < 2 {
			logs.Debugf("Skip line %s.", line)
			continue
		} else {
			c.configurationKVM[pair[0]] = line[len(pair[0])+1:]
		}
	}
	return nil
}

func (c *IniConfiguration) Get(key string) string {
	if c.configurationKVM[key] == "" {
		return ""
	}
	return c.configurationKVM[key]
}

func (c *IniConfiguration) Unmarshal(dest interface{}) error {
	v := reflect.ValueOf(dest).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if !field.CanSet() {
			continue
		}

		tag := t.Field(i).Tag.Get("ini")
		if tag == "" {
			continue
		}

		if value := c.Get(tag); value != "" {
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
					field.SetInt(intValue)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
					field.SetUint(uintValue)
				}
			case reflect.Bool:
				if boolValue, err := strconv.ParseBool(value); err == nil {
					field.SetBool(boolValue)
				}
			case reflect.Float32, reflect.Float64:
				if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
					field.SetFloat(floatValue)
				}
				// Add more cases as needed
			}
		}
	}
	return nil
}
