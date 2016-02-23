package config

type XmlConfiguration struct {
}

func (c *XmlConfiguration) Load(filePath string) error {
	return nil
}

func (c *XmlConfiguration) Get(key string) (string, error) {
	return "", nil
}
