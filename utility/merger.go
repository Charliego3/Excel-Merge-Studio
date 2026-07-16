package utility

import "github.com/xuri/excelize/v2"

type Merger struct {
	WorkbookId string
	Sheet      string
}

func NewMerger(workbookId, sheet string) *Merger {
	return &Merger{
		WorkbookId: workbookId,
		Sheet:      sheet,
	}
}

func (m *Merger) DeleteSheetMergeCells() {
	if wbMerges, ok := State.MergeCells[m.WorkbookId]; ok {
		delete(wbMerges, m.Sheet)
	}
}

func (m *Merger) DeleteWorkbookMergeCells() {
	delete(State.MergeCells, m.WorkbookId)
}

func (m *Merger) IsMerged(col, row int) bool {
	_, ok := m.MergedCell(col, row)
	return ok
}

func (m *Merger) IsMergeFirstCell(col, row int) bool {
	_, ok := m.MergedFirstCell(col, row)
	return ok
}

func (m *Merger) MergedFirstCell(col, row int) (MergeInfo, bool) {
	info, ok := m.mergeIndex(col, row, func(axis string) (MergeInfo, bool) {
		info, ok := State.MergeCells[m.WorkbookId][m.Sheet].Starts[axis]
		return info, ok
	})
	return info, ok
}

func (m *Merger) MergedCell(col, row int) (MergeInfo, bool) {
	return m.mergeIndex(col, row, func(axis string) (MergeInfo, bool) {
		info, ok := State.MergeCells[m.WorkbookId][m.Sheet].Cells[axis]
		return info, ok
	})
}

func (m *Merger) mergeIndex(col, row int, fn func(string) (MergeInfo, bool)) (MergeInfo, bool) {
	axis, _ := excelize.CoordinatesToCellName(col, row)
	if info, ok := fn(axis); ok {
		return info, true
	}
	return MergeInfo{}, false
}

func (m *Merger) InitCurrentSheetMerges(file *excelize.File) {
	if wbMerges, ok := State.MergeCells[m.WorkbookId]; ok {
		if _, ok := wbMerges[m.Sheet]; ok {
			return
		}
	}

	mergeCells, err := file.GetMergeCells(m.Sheet)
	if err != nil {
		panic(err)
	}

	sheetMergeCells := SheetMergeCells{
		Cells:  make(map[string]MergeInfo),
		Starts: make(map[string]MergeInfo),
	}
	for _, m := range mergeCells {
		sc, sr, _ := excelize.CellNameToCoordinates(m.GetStartAxis())
		ec, er, _ := excelize.CellNameToCoordinates(m.GetEndAxis())

		info := MergeInfo{
			Start:   m.GetStartAxis(),
			End:     m.GetEndAxis(),
			RowSpan: er - sr + 1,
			ColSpan: ec - sc + 1,
			Value:   m.GetCellValue(),

			StartRow: sr,
			StartCol: sc,
			EndRow:   er,
			EndCol:   ec,
		}

		sheetMergeCells.Starts[m.GetStartAxis()] = info

		for i := sr; i <= er; i++ {
			for c := sc; c <= ec; c++ {
				axis, _ := excelize.CoordinatesToCellName(c, i)
				sheetMergeCells.Cells[axis] = info
			}
		}
	}

	workbookMergeCells := make(map[string]SheetMergeCells)
	workbookMergeCells[m.Sheet] = sheetMergeCells
	State.MergeCells[m.WorkbookId] = workbookMergeCells
}
