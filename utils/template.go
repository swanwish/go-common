package utils

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/swanwish/go-common/logs"
)

var templateMap = map[string]*template.Template{}

func ExecuteTemplate(filenames []string, data interface{}) ([]byte, error) {
	//templateFiles := []string{"templates/container.xml"}
	key := GetMD5Hash(strings.Join(filenames, ","))
	var tpl *template.Template
	var exists bool
	var err error
	if tpl, exists = templateMap[key]; !exists {
		tpl, err = template.ParseFiles(filenames...)
		if err != nil {
			logs.Errorf("Failed to parse template file %s", strings.Join(filenames, ", "))
			return nil, err
		}
		templateMap[key] = tpl
	}
	buff := bytes.NewBufferString("")
	err = tpl.Execute(buff, data)
	if err != nil {
		logs.Errorf("Failed to generate contaienr.xml, the error is %v", err)
		return nil, err
	}
	return buff.Bytes(), nil
}
