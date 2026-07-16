package utility

type Setting struct {
	Workbook string
	Sheet    string
	Cols     []int
	Rows     []int
}

type SettingWithHeader struct {
	Workbook string
	Sheet    string
	Cols     []int
	Rows     []int

	Header     int
	PrimaryKey int
}

type Main struct {
	Workbook string
	Sheet    string
}

type MergeInfo struct {
	Start   string
	End     string
	RowSpan int
	ColSpan int

	StartRow, StartCol int
	EndRow, EndCol     int

	Value string
}

type SheetMergeCells struct {
	// 所有属于合并区域的单元格
	Cells map[string]MergeInfo

	// 合并区域左上角
	Starts map[string]MergeInfo
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

	IsMain    bool
	MainSheet string
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

	Header     int
	PrimaryKey int
}

func (s *Sheet) Reset() {
	s.Header = 0
	s.PrimaryKey = -1
	s.Rows = 0
	s.Columns = 0
	s.Data = nil
}

type Row struct {
	// Columns int
	Data []*Cell
}

type Cell struct {
	Value string

	// 是否属于某个合并区域（包括左上角）
	IsMerged bool

	ColSpan int
	RowSpan int

	Skip bool

	// Columns  int
	StartRow int
	StartCol int
	EndRow   int
	EndCol   int
}
