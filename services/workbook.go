package services

import (
	"fmt"
	"merger/utility"
)

type Workbook struct{}

func (w *Workbook) GetWorkbook(id string) *utility.Workbook {
	if wk, ok := utility.State().Workbooks[id]; ok {
		return wk
	}
	utility.State().App.Dialog.Info().
		AttachToWindow(utility.State().MainWindow).
		SetTitle("Workbook not found").
		SetMessage("The workbook with the specified ID was not found.").
		Show()
	return nil
}

func (w *Workbook) Workbooks() []*utility.Workbook {
	fmt.Println(len(utility.State().Workbooks))
	return utility.State().AllWorkbooks()
}
