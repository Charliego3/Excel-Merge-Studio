package services

import (
	"fmt"
	"merger/utility"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Reader struct {
	file     *excelize.File
	workbook *utility.Workbook
	sheet    *utility.Sheet
	merges   []merge

	reading bool
}

func (r *Reader) read(id, file string) map[string]any {
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
		id := utility.Sha256Hash(sheetName)
		r.sheet = &utility.Sheet{
			ID:   id,
			Name: sheetName,
		}

		rows, err := r.file.Rows(sheetName)
		if err != nil {
			return nil
		}

		rowIndex := 1
		for rows.Next() {
			columns, err := rows.Columns()
			if err != nil {
				return nil
			}

			if len(columns) > r.sheet.Columns {
				r.sheet.Columns = len(columns)
			}

			// fmt.Println("\n\n\nCells: ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
			var rowValues []utility.Cell
			for index, _ := range columns {
				var cell utility.Cell
				columnIndex := index + 1
				if r.onMerge(rowIndex, columnIndex, func(merge merge) {
					cell.Value = merge.value
					cell.IsMerged = true
					cell.Skip = !(rowIndex == merge.startRow && columnIndex == merge.startCol)
					cell.ColSpan = merge.endedCol - merge.startCol + 1
					cell.RowSpan = merge.endedRow - merge.startRow + 1
				}) {
					cellName, _ := excelize.CoordinatesToCellName(columnIndex, rowIndex)
					cell.Value, _ = r.file.CalcCellValue(sheetName, cellName)
					cell.ColSpan = 1
					cell.RowSpan = 1
				}

				cell.Value = strings.TrimSpace(cell.Value)
				// fmt.Printf("%#v\n", cell)
				rowValues = append(rowValues, cell)
			}
			// fmt.Println("Cells: ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑\n\n\n")
			rowIndex++
			r.sheet.Data = append(r.sheet.Data, rowValues)
		}

		// fmt.Println("Sheet:", r.sheet.Name, ", Columns:", r.sheet.Columns)
		for index, columns := range r.sheet.Data {
			var columnCount int
			for _, column := range columns {
				if column.IsMerged && !column.Skip {
					columnCount += column.ColSpan
				} else if !column.IsMerged {
					columnCount++
				}
			}
			if columnCount < r.sheet.Columns {
				// fmt.Println("\n\n\n")
				// fmt.Println("第", index+1, "行", ", 实际列数:", columnCount, ", 期望列数:", r.sheet.Columns)
				// for j, cell := range columns {
				// fmt.Printf("第 %d 行 第 %+d 列: %s\n", index+1, j+1, cell.Value)
				// }
				for i := columnCount; i < r.sheet.Columns; i++ {
					columns = append(columns, utility.Cell{})
					columnCount++
				}
				r.sheet.Data[index] = columns
				// fmt.Printf("%#v\n", columnCount)
				// fmt.Println("\n\n\n")
			} else {
				// fmt.Println("\n\n\n========", "第", index+1, "行", ", 实际列数:", columnCount, ", 期望列数:", r.sheet.Columns)
				// for _, cell := range columns {
				// fmt.Printf("%+v\n", cell)
				// }
				// fmt.Println("========\n\n\n")
			}
		}

	ROW:
		for i := len(r.sheet.Data) - 1; i >= 0; i-- {
			for _, v := range r.sheet.Data[i] {
				if v.Value != "" {
					break ROW
				}
			}
			r.sheet.Data = r.sheet.Data[:i]
		}

		r.sheet.Rows = rowIndex
		r.merges = nil
		err = rows.Close()
		if err != nil {
			return nil
		}
		r.workbook.Sheets[id] = r.sheet
		sheets = append(sheets, r.sheet)
	}

	utility.State().AddWorkbook(r.workbook)
	return map[string]any{
		"file":   file,
		"sheets": sheets,
	}
}

func (r *Reader) isMerged(row, column int) bool {
	if r.merges == nil {
		r.initCurrentSheetMerges()
	}

	for _, merge := range r.merges {
		if row >= merge.startRow && row <= merge.endedRow && column >= merge.startCol && column <= merge.endedCol {
			return true
		}
	}
	return false
}

func (r *Reader) onMerge(row, column int, fn func(merge merge)) bool {
	if r.merges == nil {
		r.initCurrentSheetMerges()
	}

	for _, merge := range r.merges {
		if row >= merge.startRow && row <= merge.endedRow && column >= merge.startCol && column <= merge.endedCol {
			fn(merge)
			return false
		}
	}
	return true
}

func (r *Reader) initCurrentSheetMerges() {
	mergeCells, err := r.file.GetMergeCells(r.sheet.Name)
	if err != nil {
		panic(err)
	}

	var merges []merge
	for _, cell := range mergeCells {
		merge := merge{}
		merge.startCol, merge.startRow, err = excelize.CellNameToCoordinates(cell.GetStartAxis())
		if err != nil {
			panic(err)
		}

		merge.endedCol, merge.endedRow, err = excelize.CellNameToCoordinates(cell.GetEndAxis())
		if err != nil {
			panic(err)
		}
		merge.value = cell.GetCellValue()
		merges = append(merges, merge)
	}
	r.merges = merges
	// fmt.Printf("Merges: %#v\n\n", merges)
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
		utility.State().App.Dialog.Warning().
			SetTitle("正在读取中").
			SetMessage("请稍候...").
			AttachToWindow(utility.State().MainWindow).
			Show()
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
	if utility.State().ContainsWorkbook(id) {
		dialog := utility.State().App.Dialog.Question().
			SetTitle("该文件已经存在，是否要覆盖?").
			AttachToWindow(utility.State().MainWindow).
			SetMessage(fmt.Sprintf("Workbook %s is already open. Do you want to override it?", path))

		cancelBtn := dialog.AddButton("取消")
		cancelBtn.OnClick(func() {
			ch <- false
		})
		overrideBtn := dialog.AddButton("覆盖")
		overrideBtn.OnClick(func() {
			ch <- true
		})
		dialog.SetDefaultButton(overrideBtn)
		dialog.SetCancelButton(cancelBtn)
		dialog.Show()
	} else {
		ch <- true
	}

	if !<-ch {
		return nil
	}
	fmt.Println(path)
	return r.read(id, path)
}

type merge struct {
	startCol, startRow int
	endedCol, endedRow int

	value string
}
