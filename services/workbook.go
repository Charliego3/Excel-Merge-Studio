package services

import (
	"fmt"
	"merger/utility"
)

type Workbook struct{}

func (w *Workbook) RemoveWorkbook(id string) bool {
	if _, ok := utility.State().Workbooks[id]; ok {
		delete(utility.State().Workbooks, id)
		for index, wkId := range utility.State().WorkbookIds {
			if wkId == id {
				utility.State().WorkbookIds = append(utility.State().WorkbookIds[:index], utility.State().WorkbookIds[index+1:]...)
				break
			}
		}
		utility.State().App.Event.Emit("workbooks:updated")
		return true
	}
	return false
}

func (w *Workbook) GetWorkbook(id string) utility.WorkbookInfo {
	if wk, ok := utility.State().Workbooks[id]; ok {
		info := utility.WorkbookInfo{
			ID:       id,
			FilePath: wk.FilePath,
			Name:     wk.Name,
		}
		for _, sheetName := range wk.SheetNames {
			info.Sheets = append(info.Sheets, wk.Sheets[sheetName])
		}
		return info
	}
	utility.State().App.Dialog.Info().
		AttachToWindow(utility.State().MainWindow).
		SetTitle("Workbook not found").
		SetMessage("The workbook with the specified ID was not found.").
		Show()
	return utility.WorkbookInfo{}
}

func (w *Workbook) Workbooks() []*utility.Workbook {
	fmt.Println(len(utility.State().Workbooks))
	return w.AllWorkbooks()
}

func (w *Workbook) WorkbooksMeta() []utility.WorkbookMeta {
	var metas []utility.WorkbookMeta
	for _, wk := range utility.State().Workbooks {
		metas = append(metas, utility.WorkbookMeta{
			ID:         wk.ID,
			FilePath:   wk.FilePath,
			Name:       wk.Name,
			SheetNames: wk.SheetNames,
		})
	}
	return metas
}

func (w *Workbook) Sheets(id string) []*utility.Sheet {
	if !w.ContainsWorkbook(id) {
		return nil
	}
	wk := utility.State().Workbooks[id]
	if wk == nil {
		return nil
	}

	var sheets []*utility.Sheet
	for _, sheetName := range wk.SheetNames {
		sheets = append(sheets, wk.Sheets[sheetName])
	}
	return sheets
}

func (w *Workbook) AddWorkbook(workbook *utility.Workbook) {
	utility.State().Workbooks[workbook.ID] = workbook
	if w.ContainsWorkbook(workbook.ID) {
		return
	}
	utility.State().WorkbookIds = append(utility.State().WorkbookIds, workbook.ID)
}

func (w *Workbook) ContainsWorkbook(id string) bool {
	for _, workbookId := range utility.State().WorkbookIds {
		if workbookId == id {
			return true
		}
	}
	return false
}

func (w *Workbook) AllWorkbooks() []*utility.Workbook {
	var workbooks []*utility.Workbook
	for _, id := range utility.State().WorkbookIds {
		workbooks = append(workbooks, utility.State().Workbooks[id])
	}
	return workbooks
}
