package config

import (
	"io/ioutil"
	"os"
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

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		logs.Errorf("Failed to read configuration file.", err)
		return err
	}

	c.configurationKVM = make(map[string]string, 0)
	lines := strings.Split(string(fileContent), "\n")
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
