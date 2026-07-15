package services

import (
	"merger/utility"
	"slices"
)

type Setting struct{}

func (s *Setting) DeleteColsAndRows(model utility.Setting) {
	if len(model.Rows) == 0 && len(model.Cols) == 0 {
		utility.ShowWarning("先选择需要删除的行或列再操作")
		return
	}

	if _, sheet := s.check(model.Workbook, model.Sheet); sheet != nil {
		var message string
		if len(model.Rows) == sheet.Rows && len(model.Cols) == sheet.Columns {
			message = "确定删除全部的行和全部的列吗?"
		} else if len(model.Rows) == sheet.Rows {
			message = "确定删除全部的行吗?"
		} else if len(model.Cols) == sheet.Columns {
			message = "确定删除全部的列吗?"
		}

		if message != "" {
			utility.Confirm(message, "", func(ok bool) {
				if ok {
					sheet.Rows = 0
					sheet.Columns = 0
					sheet.Data = nil
					utility.State().App.Event.Emit("setting:deleted:row_col", model)
				}
			})
			return
		}

		var colums int
		var rows []*utility.Row
		for idx, row := range sheet.Data {
			if slices.Contains(model.Rows, idx) {
				continue
			}

			var cells []utility.Cell
			for col, cell := range row.Data {
				if !slices.Contains(model.Cols, col) {
					cells = append(cells, cell)
				}
			}
			row.Data = cells
			row.Columns = len(cells)
			if row.Columns > colums {
				colums = row.Columns
			}
			rows = append(rows, row)
		}
		sheet.Data = rows
		sheet.Rows = len(rows)
		sheet.Columns = colums
		utility.State().App.Event.Emit("setting:deleted:row_col", model)
	}
}

func (s *Setting) GetMain() utility.Main {
	return utility.State().Main
}

func (s *Setting) SetMain(model utility.Main) {
	if _, sheet := s.check(model.Workbook, model.Sheet); sheet == nil {
		return
	}
	utility.State().Main = model
	utility.State().App.Event.Emit("workbooks:updated")
}

func (s *Setting) check(workbook, sheetName string) (*utility.Workbook, *utility.Sheet) {
	if workbook == "" {
		utility.ShowWarning("请选择正确的工作薄")
		return nil, nil
	}
	if sheetName == "" {
		utility.ShowWarning("请选择正确的工作表")
		return nil, nil
	}
	wk, ok := utility.State().Workbooks[workbook]
	if !ok {
		utility.ShowWarning("工作薄未上传")
		return nil, nil
	}
	sheet, ok := wk.Sheets[sheetName]
	if !ok {
		utility.ShowWarning("工作薄中找不到对应的表")
		return nil, nil
	}
	return wk, sheet
}

func (s *Setting) SetHeader(model utility.Setting) {
	if len(model.Rows) != 1 {
		utility.ShowWarning("设置表头需要选择对应的行且只能选择一行")
		return
	}

	if _, sheet := s.check(model.Workbook, model.Sheet); sheet != nil {
		sheet.Header = model.Rows[0]
		utility.State().App.Event.Emit("workbook:header:setting", model)
	}
}
