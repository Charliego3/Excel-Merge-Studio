package utility

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

type AppState struct {
	App        *application.App
	MainWindow *application.WebviewWindow

	WorkbookIds []string
	Workbooks   map[string]*Workbook
}

func (s *AppState) Sheets(id string) []*Sheet {
	if !s.ContainsWorkbook(id) {
		return nil
	}
	wk := s.Workbooks[id]
	if wk == nil {
		return nil
	}

	var sheets []*Sheet
	for _, sheetName := range wk.SheetNames {
		sheets = append(sheets, wk.Sheets[sheetName])
	}
	return sheets
}

func (s *AppState) AddWorkbook(workbook *Workbook) {
	s.Workbooks[workbook.ID] = workbook
	if s.ContainsWorkbook(workbook.ID) {
		return
	}
	s.WorkbookIds = append(s.WorkbookIds, workbook.ID)
}

func (s *AppState) ContainsWorkbook(id string) bool {
	for _, workbookId := range s.WorkbookIds {
		if workbookId == id {
			return true
		}
	}
	return false
}

func (s *AppState) AllWorkbooks() []*Workbook {
	var workbooks []*Workbook
	for _, id := range s.WorkbookIds {
		workbooks = append(workbooks, s.Workbooks[id])
	}
	return workbooks
}

type WorkbookInfo struct {
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
	Name    string
	Columns int
	Rows    int
	Data    [][]Cell
}

type Cell struct {
	Value string

	// 是否属于某个合并区域（包括左上角）
	IsMerged bool

	ColSpan int
	RowSpan int

	Skip bool
}

var state *AppState = &AppState{
	Workbooks: make(map[string]*Workbook),
}

func State() *AppState {
	return state
}
