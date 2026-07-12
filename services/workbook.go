package services

import (
	"fmt"
	"merger/utility"
)

type Workbook struct{}

func (w *Workbook) GetWorkbook(id string) utility.WorkbookInfo {
	if wk, ok := utility.State().Workbooks[id]; ok {
		info := utility.WorkbookInfo{
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
	return utility.State().AllWorkbooks()
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
