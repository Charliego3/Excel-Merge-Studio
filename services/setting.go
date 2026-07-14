package services

import "merger/utility"

type Setting struct{}

func (s *Setting) SetHeader(model utility.Setting) {
	if model.Workbook == "" {
		s.showWarning("请选择正确的工作薄")
		return
	}
	if model.Sheet == "" {
		s.showWarning("请选择正确的工作表")
		return
	}
	if len(model.Rows) != 1 {
		s.showWarning("设置表头需要选择对应的行且只能选择一行")
		return
	}

	workbook := utility.State().Workbooks[model.Workbook]
	sheet := workbook.Sheets[model.Sheet]
	sheet.Header = model.Rows[0]

	utility.State().App.Event.Emit("workbook:sheet:setting", model)
}

func (s *Setting) showWarning(message string) {
	utility.State().App.Dialog.Warning().
		AttachToWindow(utility.State().MainWindow).
		SetTitle("Warning").
		SetMessage(message).
		Show()
}
