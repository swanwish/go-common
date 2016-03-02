package easyui

func (model *GridData) AddRow(row interface{}) {
	model.Rows = append(model.Rows, row)
	model.Total = int64(len(model.Rows))
}
