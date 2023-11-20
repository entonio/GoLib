package xlsrange

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

type Range struct {
	sheet           *xlsx.Sheet
	topRowIndex     uint
	rowCount        uint
	leftColumnIndex uint
	columnCount     uint
}

type Field struct {
	Label RC
	Value RC
}

type RC struct {
	Row uint
	Col uint
}

func HField(row uint, labelCol uint, valueCol uint) Field {
	return Field{
		Label: RC{Row: row, Col: labelCol},
		Value: RC{Row: row, Col: valueCol},
	}
}

func VField(labelRow uint, valueRow uint, col uint) Field {
	return Field{
		Label: RC{Row: labelRow, Col: col},
		Value: RC{Row: valueRow, Col: col},
	}
}

func NewRangeFromSheet(sheet *xlsx.Sheet) *Range {
	return newRange(sheet, 0, 0, 0, 0)
}

func NewRangeFromRow(sheet *xlsx.Sheet, topRowIndex uint) *Range {
	return newRange(sheet, topRowIndex, 1, 0, 0)
}

func NewRangeFromCell(sheet *xlsx.Sheet, topRowIndex uint, leftColumnIndex uint) *Range {
	return newRange(sheet, topRowIndex, 1, leftColumnIndex, 1)
}

func newRange(sheet *xlsx.Sheet, topRowIndex uint, rowCount uint, leftColumnIndex uint, columnCount uint) *Range {
	return &Range{
		sheet:           sheet,
		topRowIndex:     topRowIndex,
		rowCount:        rowCount,
		leftColumnIndex: leftColumnIndex,
		columnCount:     columnCount,
	}
}

func (self *Range) File() *xlsx.File {
	return self.sheet.File
}

func (self *Range) String() string {
	return fmt.Sprintf("Range<%d-%d, %d-%d of %dx%d>",
		self.topRowIndex,
		self.topRowIndex+self.rowCount,
		self.leftColumnIndex,
		self.leftColumnIndex+self.columnCount,
		len(self.sheet.Rows),
		self.MaxColumnIndex(),
	)
}

func (self *Range) rowAtOffset(offset uint) *xlsx.Row {
	return self.sheet.Rows[self.topRowIndex+offset]
}

func (self *Range) cellAtOffset(r uint, c uint) *xlsx.Cell {
	return self.rowAtOffset(r).Cells[self.leftColumnIndex+c]
}

func (self *Range) bottomRowIndex() uint {
	if self.rowCount > 0 {
		return self.topRowIndex + self.rowCount - 1
	} else {
		return self.MaxRowIndex()
	}
}

func (self *Range) rightColumnIndex() uint {
	if self.columnCount > 0 {
		return self.leftColumnIndex + self.columnCount - 1
	} else {
		return self.MaxColumnIndex()
	}
}

func (self *Range) toRowIndex(fromRowIndex uint, rowCount uint) uint {
	if rowCount > 0 {
		return fromRowIndex + rowCount - 1
	} else {
		return self.bottomRowIndex()
	}
}

func (self *Range) toColumnIndex(fromColumnIndex uint, columnCount uint) uint {
	if columnCount > 0 {
		return fromColumnIndex + columnCount - 1
	} else {
		return self.rightColumnIndex()
	}
}

func subtract(from uint, amount uint) uint {
	if amount > from {
		return 0
	} else {
		return from - amount
	}
}
