package easyui

type TreeItem struct {
	Id         string      `json:"id"`
	ParentId   string      `json:"-"`
	Text       string      `json:"text"`
	IconClass  string      `json:"iconCls"`
	Children   []TreeItem  `json:"children,omitempty"`
	Attributes interface{} `json:"attributes,omitempty"`
}

type GridData struct {
	Rows  []interface{} `json:"rows"`
	Total int64         `json:"total"`
}

func CreateGridData() GridData {
	gridData := GridData{}
	gridData.Rows = make([]interface{}, 0)
	return gridData
}
