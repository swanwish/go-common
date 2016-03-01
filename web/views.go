package web

import (
	"encoding/json"
	"errors"
	"html/template"
	"io/ioutil"
	"path/filepath"

	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/utils"
)

type ViewConfiguration struct {
	CommonTemplates []string         `json:"commonTemplates"`
	CommonView      string           `json:"commonView"`
	DefaultViewId   string           `json:"defaultViewId"`
	LoginViewId     string           `json:"loginViewId"`
	Views           []ViewDefinition `json:"views"`
	ImportViews     []string         `json:"importViews"`
}

type Style struct {
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	Media string `json:"media"`
}

type ViewDefinition struct {
	Id              string   `json:"id"`
	SelectViewId    string   `json:"selectViewId"`
	AnonymousAccess bool     `json:"anonymousAccess"`
	View            string   `json:"view"`
	PageTitle       string   `json:"pageTitle"`
	PageHeader      string   `json:"pageHeader"`
	Templates       []string `json:"templates"`
	Scripts         []string `json:"scripts"`
	Styles          []Style  `json:"styles"`
	MainRequire     string   `json:"mainRequire"`
	ModuleScripts   []string `json:"moduleScripts"`
	StyleSheets     []string `json:"stylesheets"`
}

var (
	viewMap           = map[string]ViewDefinition{}
	viewTemplateCache = make(map[string]*template.Template, 0)
	DefaultViewId     string
	LoginViewId       string

	ErrConfigurationFileNotExist = errors.New("Configuration file not exist")
	ErrNoSuchView                = errors.New("No such view")
	ErrInvalidConfiguration      = errors.New("Invalid configuration")

	ProductMode bool
)

func LoadViews(configureFileName string) error {
	if configureFileName == "" {
		logs.Errorf("The configuration file name is not specified.")
		return ErrConfigurationFileNotExist
	}
	logs.Debugf("Load view from file `%s`", configureFileName)
	if utils.FileExists(configureFileName) {
		content, err := ioutil.ReadFile(configureFileName)
		if err != nil {
			logs.Errorf("Failed to read configuration file %s, the error is %v", configureFileName, err)
			return err
		}
		viewConfiguration := ViewConfiguration{}
		err = json.Unmarshal(content, &viewConfiguration)
		if err != nil {
			logs.Errorf("Failed to unmarshal json data %s, the error is %v", string(content), err)
			return err
		}
		DefaultViewId = viewConfiguration.DefaultViewId
		LoginViewId = viewConfiguration.LoginViewId
		cacheViewDefinitions(viewConfiguration, viewConfiguration.Views)
		importExternalViews(viewConfiguration, filepath.Dir(configureFileName))
		return nil
	}
	logs.Errorf("The configuration file %s does not exist", configureFileName)
	return ErrConfigurationFileNotExist
}

func importExternalViews(viewConfiguration ViewConfiguration, dir string) error {
	for _, importView := range viewConfiguration.ImportViews {
		filePath := filepath.Join(dir, importView)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			logs.Errorf("Failed to read configuration file %s, the error is %v", filePath, err)
			return err
		}
		config := ViewConfiguration{}
		err = json.Unmarshal(content, &config)
		if err != nil {
			logs.Errorf("Failed to unmarshal json data %s, the error is %v", string(content), err)
			return err
		}
		err = cacheViewDefinitions(viewConfiguration, config.Views)
		if err != nil {
			logs.Errorf("Faile to cache view definitions, the error is %v", err)
			return err
		}
	}
	return nil
}

func cacheViewDefinitions(viewConfiguration ViewConfiguration, views []ViewDefinition) error {
	for _, view := range views {
		view.Templates = append(view.Templates, viewConfiguration.CommonTemplates...)
		if view.View == "" {
			if viewConfiguration.CommonView != "" {
				view.View = viewConfiguration.CommonView
			} else {
				logs.Errorf("The view property is not set for view %s", view.Id)
				return ErrInvalidConfiguration
			}
		}
		viewMap[view.Id] = view
	}
	return nil
}

func GetView(id string) (ViewDefinition, error) {
	if id == "" {
		logs.Errorf("The view id is empty.")
		return ViewDefinition{}, ErrNoSuchView
	}
	if view, ok := viewMap[id]; ok {
		return view, nil
	}
	return ViewDefinition{}, ErrNoSuchView
}

func GetTemplate(id string) (*template.Template, error) {
	if t, ok := viewTemplateCache[id]; ok {
		return t, nil
	}
	if view, ok := viewMap[id]; ok {
		t, err := template.ParseFiles(view.Templates...)
		if err != nil {
			logs.Errorf("Failed to parse templates, the error is %v", err)
			return nil, err
		}
		if ProductMode {
			viewTemplateCache[id] = t
		}
		return t, err
	}
	return nil, ErrNoSuchView
}
