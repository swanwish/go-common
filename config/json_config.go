package config

type JsonConfiguration struct {
}

type ConfigItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ConfigWrapper struct {
	Items []ConfigItem `json:"list"`
}

func (c *JsonConfiguration) Load(filePath string) error {
	return nil
}

func (c *JsonConfiguration) Get(key string) (string, error) {
	return "", nil
}
