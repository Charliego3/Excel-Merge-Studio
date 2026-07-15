package services

import (
	"fmt"
	"merger/utility"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/xuri/excelize/v2"
)

type MergeIndex struct {
	// 所有属于合并区域的单元格
	cells map[string]*MergeInfo

	// 合并区域左上角
	starts map[string]*MergeInfo
}

type Reader struct {
	file     *excelize.File
	workbook *utility.Workbook
	sheet    *utility.Sheet
	merges   *MergeIndex

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
			utility.State().App.Dialog.Error().
				SetTitle("读取文件出错").
				SetMessage(fmt.Sprintf("%s", message)).
				AttachToWindow(utility.State().MainWindow).
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

	var sheets []*utility.Sheet
	for _, sheetName := range r.workbook.SheetNames {
		r.sheet = &utility.Sheet{WorkbookId: id, Name: sheetName}
		rows, err := r.file.Rows(sheetName)
		if err != nil {
			return nil
		}

		rowIdx := 1
		for rows.Next() {
			columns, err := rows.Columns()
			if err != nil {
				return nil
			}

			if len(columns) > r.sheet.Columns {
				r.sheet.Columns = len(columns)
			}

			row := &utility.Row{}
			// var colStart int
			for col := 1; col <= len(columns); {
				merge, ok := r.merged(col, rowIdx)
				if ok {
					cell := newCell(merge)
					cell.IsMerged = true
					cell.Skip = !r.isMergeFirstCell(col, rowIdx)
					row.Data = append(row.Data, cell)

					for i := 1; i < cell.ColSpan; i++ {
						cell := newCell(merge)
						cell.IsMerged = true
						cell.Skip = true
						row.Data = append(row.Data, cell)
					}
					col += cell.ColSpan
				} else {
					var cell utility.Cell
					cellName, _ := excelize.CoordinatesToCellName(col, rowIdx)
					value, _ := r.file.CalcCellValue(sheetName, cellName)
					cell.Value = strings.TrimSpace(value)
					cell.StartCol = col
					cell.StartRow = rowIdx
					row.Data = append(row.Data, cell)
					col++
				}
			}

			// for index, _ := range columns {
			// 	var cell utility.Cell
			// 	col := index + 1
			// 	merge, ok := r.merged(col, rowIdx)
			// 	if ok {
			// 		cell.Value = merge.Value
			// 		cell.IsMerged = true
			// 		cell.Skip = !r.isMergeFirstCell(col, rowIdx)
			// 		cell.ColSpan = merge.ColSpan
			// 		cell.RowSpan = merge.RowSpan
			// 		cell.StartCol = merge.startCol
			// 		cell.EndCol = merge.endCol
			// 		cell.StartRow = merge.startRow
			// 		cell.EndRow = merge.endRow
			// 		colStart += cell.ColSpan
			// 	} else {
			// 		cellName, _ := excelize.CoordinatesToCellName(col, rowIdx)
			// 		cell.Value, _ = r.file.CalcCellValue(sheetName, cellName)
			// 		colStart++
			// 		cell.StartCol = colStart
			// 		cell.StartRow = rowIdx
			// 		row.Columns += 1
			// 	}

			// 	cell.Value = strings.TrimSpace(cell.Value)
			// 	row.Data = append(row.Data, cell)
			// }
			rowIdx++
			r.sheet.Data = append(r.sheet.Data, row)
		}

		var truncated bool
	ROW:
		for i := len(r.sheet.Data) - 1; i >= 0; i-- {
			length := len(r.sheet.Data[i].Data)
			if r.sheet.Columns > length {
				for j := length + 1; j <= r.sheet.Columns; j++ {
					var value string
					merge, skip := r.merged(j, i+1)
					if skip {
						axis, _ := excelize.CoordinatesToCellName(merge.startCol, merge.startRow)
						if m, ok := r.merges.starts[axis]; ok {
							value = m.Value
						}
					}
					r.sheet.Data[i].Data = append(r.sheet.Data[i].Data, utility.Cell{
						Value:    value,
						Skip:     skip,
						StartRow: i + 1,
						EndRow:   i + 1,
						StartCol: j,
						EndCol:   j,
						ColSpan:  1,
						RowSpan:  1,
					})
				}
			}

			if !truncated {
				for _, row := range r.sheet.Data[i].Data {
					if row.Value != "" {
						truncated = true
						continue ROW
					}
				}
				r.sheet.Data = r.sheet.Data[:i]
			}
		}

		r.sheet.Rows = len(r.sheet.Data)
		r.merges = nil
		err = rows.Close()
		if err != nil {
			return nil
		}
		r.workbook.Sheets[sheetName] = r.sheet
		sheets = append(sheets, r.sheet)
	}

	(&Workbook{}).AddWorkbook(r.workbook)
	utility.State().App.Event.Emit("workbooks:updated")
	return sheets
}

func (r *Reader) isMerged(col, row int) bool {
	if _, ok := r.merged(col, row); ok {
		return true
	}
	return false
}

func (r *Reader) isMergeFirstCell(col, row int) bool {
	_, ok := r.mergeIndex(col, row, func(axis string) (*MergeInfo, bool) {
		info, ok := r.merges.starts[axis]
		return info, ok
	})
	return ok
}

func (r *Reader) merged(col, row int) (*MergeInfo, bool) {
	return r.mergeIndex(col, row, func(axis string) (*MergeInfo, bool) {
		info, ok := r.merges.cells[axis]
		return info, ok
	})
}

func (r *Reader) mergeIndex(col, row int, fn func(string) (*MergeInfo, bool)) (*MergeInfo, bool) {
	if r.merges == nil {
		r.initCurrentSheetMerges()
	}

	axis, _ := excelize.CoordinatesToCellName(col, row)
	if info, ok := fn(axis); ok {
		return info, true
	}
	return nil, false
}

func (r *Reader) initCurrentSheetMerges() {
	merges, err := r.file.GetMergeCells(r.sheet.Name)
	if err != nil {
		panic(err)
	}

	r.merges = &MergeIndex{
		cells:  make(map[string]*MergeInfo),
		starts: make(map[string]*MergeInfo),
	}
	for _, m := range merges {
		sc, sr, _ := excelize.CellNameToCoordinates(m.GetStartAxis())
		ec, er, _ := excelize.CellNameToCoordinates(m.GetEndAxis())

		info := &MergeInfo{
			Start:   m.GetStartAxis(),
			End:     m.GetEndAxis(),
			RowSpan: er - sr + 1,
			ColSpan: ec - sc + 1,
			Value:   m.GetCellValue(),

			startRow: sr,
			startCol: sc,
			endRow:   er,
			endCol:   ec,
		}

		r.merges.starts[m.GetStartAxis()] = info

		for i := sr; i <= er; i++ {
			for c := sc; c <= ec; c++ {
				axis, _ := excelize.CoordinatesToCellName(c, i)
				r.merges.cells[axis] = info
			}
		}
	}
}

func (r *Reader) reset() error {
	r.merges = nil
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
	path, err := utility.State().App.Dialog.OpenFile().
		SetTitle("请选择表格文件").
		SetMessage("选择要合并的表格文件").
		SetButtonText("确定").
		CanCreateDirectories(false).
		AttachToWindow(utility.State().MainWindow).
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
	utility.State().App.Event.Emit("workbook:read:start")
	return map[string]any{
		"id":     id,
		"file":   path,
		"sheets": r.read(id, path),
	}
}

func newCell(merge *MergeInfo) utility.Cell {
	var cell utility.Cell
	cell.Value = strings.TrimSpace(merge.Value)
	cell.ColSpan = merge.ColSpan
	cell.RowSpan = merge.RowSpan
	cell.StartCol = merge.startCol
	cell.EndCol = merge.endCol
	cell.StartRow = merge.startRow
	cell.EndRow = merge.endRow
	return cell
}

type MergeInfo struct {
	Start   string
	End     string
	RowSpan int
	ColSpan int

	startRow, startCol int
	endRow, endCol     int

	Value string
}
