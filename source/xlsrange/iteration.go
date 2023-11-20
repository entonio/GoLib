package xlsrange

import "github.com/tealeg/xlsx"

type RowFn func(row int, xr *xlsx.Row, stop *bool)
type CellFn func(row int, column int, cell *xlsx.Cell, stop *bool)

func (self *Range) EachRow(fromRowIndex uint, rowCount uint, addIfMissing bool, f RowFn) *Range {
	toRowIndex := self.toRowIndex(fromRowIndex, rowCount)
	if addIfMissing {
		for len(self.sheet.Rows) <= int(toRowIndex) {
			//log.Debug("addRow")
			self.sheet.AddRow()
		}
	} else {
		toRowIndex = self.bottomRowIndex()
	}
	//log.Debug("rows: %d", len(self.sheet.Rows))
	stop := false
	for row := fromRowIndex; row <= toRowIndex; row += 1 {
		f(int(row), self.sheet.Rows[row], &stop)
		if stop {
			break
		}
	}
	return self
}

func (self *Range) EachCell(fromRowIndex uint, rowCount uint, fromColumnIndex uint, columnCount uint, addIfMissing bool, f CellFn) *Range {
	globalToColumnIndex := self.toColumnIndex(fromColumnIndex, columnCount)
	//log.Debug("ri %d, rc %d, ci %d, cc %d", fromRowIndex, rowCount, fromColumnIndex, columnCount)
	return self.EachRow(fromRowIndex, rowCount, addIfMissing, func(row int, xr *xlsx.Row, stop *bool) {
		toColumnIndex := globalToColumnIndex
		if addIfMissing {
			for len(xr.Cells) <= int(toColumnIndex) {
				//log.Debug("addCell")
				xr.AddCell()
			}
		} else {
			toColumnIndex = self.rightColumnIndex()
		}
		//log.Debug("cells: %d", len(row.Cells))
		for column := fromColumnIndex; column <= toColumnIndex; column += 1 {
			//log.Debug("c %d", c)
			f(row, int(column), xr.Cells[column], stop)
			if *stop {
				break
			}
		}
	})
}

func (self *Range) AllRows(f RowFn) *Range {
	return self.EachRow(self.topRowIndex, self.rowCount, true, f)
}

func (self *Range) ensureAllCells() *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {})
}

func (self *Range) AllCells(f CellFn) *Range {
	return self.EachCell(self.topRowIndex, self.rowCount, self.leftColumnIndex, self.columnCount, true, f)
}

func (self *Range) LeftCells(f CellFn) *Range {
	return self.EachCell(self.topRowIndex, self.rowCount, self.leftColumnIndex, 1, true, f)
}

func (self *Range) RightCells(f CellFn) *Range {
	return self.EachCell(self.topRowIndex, self.rowCount, self.rightColumnIndex(), 1, true, f)
}

func (self *Range) topCells(f CellFn) *Range {
	return self.EachCell(self.topRowIndex, 1, self.leftColumnIndex, self.columnCount, true, f)
}

func (self *Range) bottomCells(f CellFn) *Range {
	return self.EachCell(self.bottomRowIndex(), 1, self.leftColumnIndex, self.columnCount, true, f)
}
