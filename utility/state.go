package utility

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

type AppState struct {
	App        *application.App
	MainWindow *application.WebviewWindow

	WorkbookIds []string
	Workbooks   map[string]*Workbook
	MergeCells  map[string]map[string]SheetMergeCells

	Main Main
}

var State *AppState = &AppState{
	Workbooks:  make(map[string]*Workbook),
	MergeCells: make(map[string]map[string]SheetMergeCells),
}
