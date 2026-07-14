package utility

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Setting struct {
	Workbook string
	Sheet    string
	Cols     []int
	Rows     []int
}

type Main struct {
	Workbook string
	Sheet    string
}

type AppState struct {
	App        *application.App
	MainWindow *application.WebviewWindow

	WorkbookIds []string
	Workbooks   map[string]*Workbook

	Main Main
}

type WorkbookInfo struct {
	ID       string
	FilePath string
	Name     string
	Sheets   []*Sheet
}

type WorkbookMeta struct {
	ID         string
	FilePath   string
	Name       string
	SheetNames []string
}

type Workbook struct {
	ID         string
	FilePath   string
	Name       string
	SheetNames []string
	Sheets     map[string]*Sheet
}

type Sheet struct {
	WorkbookId string

	Name    string
	Columns int
	Rows    int
	Data    []*Row

	Header int
}

type Row struct {
	Columns int
	Data    []Cell
}

type Cell struct {
	Value string

	// 是否属于某个合并区域（包括左上角）
	IsMerged bool

	ColSpan int
	RowSpan int

	Skip bool

	Columns  int
	StartRow int
	StartCol int
	EndRow   int
	EndCol   int
}

var state *AppState = &AppState{
	Workbooks: make(map[string]*Workbook),
}

func State() *AppState {
	return state
}
