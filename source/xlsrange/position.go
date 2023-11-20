package xlsrange

import "fmt"

func (self *Range) NextRowSameColumn() *Range {
	return self.Row(self.topRowIndex + 1)
}

func (self *Range) NextRow() *Range {
	self.leftColumnIndex = 0
	return self.NextRowSameColumn()
}

func (self *Range) RC(rc RC) *Range {
	return self.Position(rc.Row, rc.Col)
}

func (self *Range) Position(rowIndex uint, columnIndex uint) *Range {
	return self.Row(rowIndex).Cell(columnIndex)
}

func (self *Range) Row(rowIndex uint) *Range {
	self.topRowIndex = rowIndex
	return self.SelectRow(self.topRowIndex)
}

func (self *Range) NextCell() *Range {
	return self.Cell(self.leftColumnIndex + 1)
}

func (self *Range) Cell(columnIndex uint) *Range {
	self.leftColumnIndex = columnIndex
	return self.SelectColumn(self.leftColumnIndex)
}

func (self *Range) CurrentCell() *Range {
	return self.SelectColumn(self.leftColumnIndex)
}

func (self *Range) CurrentRow() *Range {
	return self.SelectRow(self.topRowIndex)
}

func (self *Range) LeftLetter(offset int) string {
	return asLetter(int(self.leftColumnIndex) + offset)
}

func (self *Range) RightLetter(offset int) string {
	return asLetter(int(self.rightColumnIndex()) + offset)
}

func (self *Range) TopIndex(offset int) int {
	return int(self.topRowIndex) + offset
}

func (self *Range) TopNumber(offset int) string {
	return fmt.Sprint(self.TopIndex(offset) + 1)
}

func (self *Range) BottomIndex(offset int) int {
	return int(self.bottomRowIndex()) + offset
}

func (self *Range) BottomNumber(offset int) string {
	return fmt.Sprint(self.BottomIndex(offset) + 1)
}

func (self *Range) MoveToBottomRow() {
	self.MoveToRow(self.BottomRow())
}

func (self *Range) MoveToBottomRowSameColumn() {
	self.MoveToRowSameColumn(self.BottomRow())
}

func (self *Range) MoveToRow(row *Range) {
	self.leftColumnIndex = 0
	self.MoveToRowSameColumn(row)
}

func (self *Range) MoveToRowSameColumn(row *Range) {
	target := row.rowAtOffset(0)
	for rowIndex, candidate := range self.sheet.Rows {
		if candidate == target {
			self.topRowIndex = uint(rowIndex)
			return
		}
	}
	panic("Target row not found")
}

func (self *Range) MoveToCell(cell *Range) {
	target := cell.cellAtOffset(0, 0)
	for rowIndex, row := range self.sheet.Rows {
		for columnIndex, candidate := range row.Cells {
			if candidate == target {
				self.topRowIndex = uint(rowIndex)
				self.leftColumnIndex = uint(columnIndex)
				return
			}
		}
	}
	panic("Target cell not found")
}

func (self *Range) BottomRow() *Range {
	// could be 0
	return self.SelectRow(self.MaxRowIndex())
}

func (self *Range) MaxRowIndex() uint {
	return uint(max(0, len(self.sheet.Rows)-1))
}

func (self *Range) MaxColumnIndex() uint {
	maxColumnIndex := 0
	for _, row := range self.sheet.Rows {
		maxColumnIndex = max(maxColumnIndex, len(row.Cells)-1)
	}
	return uint(maxColumnIndex)
}

func (self *Range) EnsureSameColumnsInAllRows() {
	maxColumnIndex := self.MaxColumnIndex()
	for rowIndex, _ := range self.sheet.Rows {
		self.SelectPosition(uint(rowIndex), maxColumnIndex)
	}
}

func (self *Range) EnsureColumns() *Range {
	maxColumnIndex := self.MaxColumnIndex()
	self.SelectPosition(uint(self.topRowIndex), maxColumnIndex)
	return self
}

func (self *Range) ColumnWidthsCM(centimeters ...float64) {
	self.SelectPosition(0, uint(len(centimeters))-1).ensureAllCells()
	for columnIndex, cm := range centimeters {
		self.sheet.SetColWidth(columnIndex, columnIndex, cm*9/2.54)
	}
}

func asLetter(c int) string {
	if c < 0 {
		panic("Cannot convert negative column " + fmt.Sprint(c) + " to letter")
	}
	if c < 26 {
		return fmt.Sprintf("%c", 'A'+c)
	}
	return asLetter(c/26) + asLetter(c%26)
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
