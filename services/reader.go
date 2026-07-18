package services

import (
	"fmt"
	"merger/utility"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Reader struct {
	file     *excelize.File
	workbook *utility.Workbook
	sheet    *utility.Sheet

	reading bool
}

func (r *Reader) read(id, file string) []*utility.Sheet {
	var err error
	defer func() {
		message := recover()
		if message == nil {
			message = err
		}
		if message != nil {
			debug.PrintStack()
			utility.State.App.Dialog.Error().
				SetTitle("读取文件出错").
				SetMessage(fmt.Sprintf("%s", message)).
				AttachToWindow(utility.State.MainWindow).
				Show()
		}
		r.reset()
	}()

	r.file, err = excelize.OpenFile(file)
	if err != nil {
		return nil
	}
	defer r.file.Close()

	r.workbook = &utility.Workbook{
		ID:         id,
		FilePath:   file,
		Name:       filepath.Base(file),
		SheetNames: r.file.GetSheetList(),
		Sheets:     make(map[string]*utility.Sheet),
	}

	// TODO: 可以先将所有单元格空白的行缓存，
	// 如果后面出现非空白的，再将这些空白行插入Sheet.Data 中，
	// 否则丢弃这些空白行
	var sheets []*utility.Sheet
	for _, sheetName := range r.workbook.SheetNames {
		r.sheet = &utility.Sheet{WorkbookId: id, Name: sheetName}
		rows, err := r.file.Rows(sheetName)
		if err != nil {
			return nil
		}

		merger := utility.NewMerger(id, sheetName)
		merger.InitCurrentSheetMerges(r.file)

		row_idx := 1
		for rows.Next() {
			columns, err := rows.Columns()
			if err != nil {
				return nil
			}

			if len(columns) > r.sheet.Columns {
				r.sheet.Columns = len(columns)
			}

			row := &utility.Row{}
			for col_idx := 1; col_idx <= len(columns); {
				merge, ok := merger.MergedCell(col_idx, row_idx)
				if ok {
					cell := newCell(merge)
					cell.IsMerged = true
					cell.Skip = !merger.IsMergeFirstCell(col_idx, row_idx)
					row.Data = append(row.Data, cell)

					for i := 1; i < cell.ColSpan; i++ {
						cell := newCell(merge)
						cell.IsMerged = true
						cell.Skip = true
						row.Data = append(row.Data, cell)
					}
					col_idx += cell.ColSpan
				} else {
					var cell utility.Cell
					cellName, _ := excelize.CoordinatesToCellName(col_idx, row_idx)
					value, _ := r.file.CalcCellValue(sheetName, cellName)
					cell.Value = strings.TrimSpace(value)
					cell.StartCol = col_idx
					cell.StartRow = row_idx
					row.Data = append(row.Data, &cell)
					col_idx++
				}
			}
			row_idx++
			r.sheet.Data = append(r.sheet.Data, row)
		}

		var truncated bool
		slices.Backward(r.sheet.Data)(func(row_idx int, row *utility.Row) bool {
			length := len(row.Data)
			if r.sheet.Columns > length {
				for col_idx := length + 1; col_idx <= r.sheet.Columns; col_idx++ {
					var value string
					merge, skip := merger.MergedCell(col_idx, row_idx+1)
					if skip {
						if m, ok := merger.MergedFirstCell(merge.StartCol, merge.StartRow); ok {
							value = m.Value
						}
					}
					row.Data = append(row.Data, &utility.Cell{
						Value:    value,
						Skip:     skip,
						StartRow: row_idx + 1,
						EndRow:   row_idx + 1,
						StartCol: col_idx,
						EndCol:   col_idx,
						ColSpan:  1,
						RowSpan:  1,
					})
				}
			}

			if !truncated {
				for _, row := range row.Data {
					if row.Value != "" {
						truncated = true
						return true
					}
				}
				r.sheet.Data = r.sheet.Data[:row_idx]
			}
			return true
		})

		r.sheet.Rows = len(r.sheet.Data)
		err = rows.Close()
		if err != nil {
			return nil
		}
		r.workbook.Sheets[sheetName] = r.sheet
		sheets = append(sheets, r.sheet)
	}

	(&Workbook{}).AddWorkbook(r.workbook)
	utility.State.App.Event.Emit("workbooks:updated")
	return sheets
}

func (r *Reader) reset() error {
	r.sheet = nil
	r.workbook = nil
	r.file = nil
	return nil
}

func (r *Reader) ShowFilePicker() map[string]any {
	if r.reading {
		utility.ShowWarning("正在读取中，请稍候...")
		return nil
	}

	r.reading = true
	defer func() {
		r.reading = false
	}()
	path, err := utility.State.App.Dialog.OpenFile().
		SetTitle("请选择表格文件").
		SetMessage("选择要合并的表格文件").
		SetButtonText("确定").
		CanCreateDirectories(false).
		AttachToWindow(utility.State.MainWindow).
		AddFilter("Excel (*.xlsx;*.xls;*.csv)", "*.xlsx;*.xls;*.csv").
		PromptForSingleSelection()

	if err != nil || path == "" {
		return nil
	}

	id := utility.Sha256Hash(path)
	ch := make(chan bool, 1)
	wk := &Workbook{}
	if wk.ContainsWorkbook(id) {
		utility.Confirm(
			"该文件已经存在，是否要覆盖?",
			fmt.Sprintf("Workbook %s is already open. Do you want to override it?", path),
			func(b bool) {
				ch <- b
			},
		)
	} else {
		ch <- true
	}

	if !<-ch {
		return map[string]any{
			"id":     id,
			"file":   path,
			"sheets": wk.Sheets(id),
		}
	}
	fmt.Println(path, id)
	utility.State.App.Event.Emit("workbook:read:start")
	return map[string]any{
		"id":     id,
		"file":   path,
		"sheets": r.read(id, path),
	}
}

func newCell(merge utility.MergeInfo) *utility.Cell {
	var cell utility.Cell
	cell.Value = strings.TrimSpace(merge.Value)
	cell.ColSpan = merge.ColSpan
	cell.RowSpan = merge.RowSpan
	cell.StartCol = merge.StartCol
	cell.EndCol = merge.EndCol
	cell.StartRow = merge.StartRow
	cell.EndRow = merge.EndRow
	return &cell
}
