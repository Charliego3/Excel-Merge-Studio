package services

import (
	"merger/utility"
	"slices"
)

type Setting struct{}

func (s *Setting) Result(setting utility.Setting, sheet *utility.Sheet) utility.SettingWithHeader {
	return utility.SettingWithHeader{
		Workbook:   setting.Workbook,
		Sheet:      setting.Sheet,
		Cols:       setting.Cols,
		Rows:       setting.Rows,
		Header:     sheet.Header,
		PrimaryKey: sheet.PrimaryKey,
	}
}

// TODO: 注意
// 删除行和列的时候要修正合并单元格
// 删除行时要考虑设置表头之后，如果删除表头上面的行，需要向上移动，如果删除表头下面的行，需要向下移动
// 删除列时要考虑主键列，同上删除行的表头
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
					sheet.Reset()
					utility.State.App.Event.Emit("setting:deleted:row_col", s.Result(model, sheet))
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

			var cells []*utility.Cell
			for col, cell := range row.Data {
				if !slices.Contains(model.Cols, col) {
					cells = append(cells, cell)
				} else if sheet.PrimaryKey == col {
					// 清除主键的位置
					sheet.PrimaryKey = -1
				} else { // 跳过被删除的列：处理删除列时的合并单元格

				}
			}
			row.Data = cells
			length := len(cells)
			if length > colums {
				colums = length
			}
			rows = append(rows, row)
		}
		sheet.Data = rows
		sheet.Rows = len(rows)
		sheet.Columns = colums
		utility.State.App.Event.Emit("setting:deleted:row_col", s.Result(model, sheet))
	}
}

func (s *Setting) GetMain() utility.Main {
	return utility.State.Main
}

func (s *Setting) SetMain(model utility.Main) {
	if _, sheet := s.check(model.Workbook, model.Sheet); sheet == nil {
		return
	}
	utility.State.Main = model
	utility.State.App.Event.Emit("workbooks:updated")
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
	wk, ok := utility.State.Workbooks[workbook]
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
		utility.State.App.Event.Emit("workbook:header:setting", model)
	}
}
