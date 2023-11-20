package xlsrange

import (
	"strings"

	"github.com/tealeg/xlsx"
)

func (self *Range) HasRow(rowOffset int) bool {
	row := int(self.topRowIndex) + rowOffset
	if row >= 0 {
		rows := self.sheet.Rows
		if row < len(rows) {
			return true
		}
	}
	return false
}

func (self *Range) HasRC(rc RC) bool {
	return self.Has(int(rc.Row), int(rc.Col))
}

func (self *Range) GetRC(rc RC) string {
	return self.Get(int(rc.Row), int(rc.Col))
}

func (self *Range) ReadRC(rc RC) (value string, exists bool) {
	return self.Read(int(rc.Row), int(rc.Col))
}

func (self *Range) Has(rowOffset int, columnOffset int) bool {
	return self.findRC(rowOffset, columnOffset) != nil
}

func (self *Range) RowIsBlank(rowOffset int) bool {
	var hasContents bool
	self.EachCell(uint(rowOffset), 1, 0, 0, false, func(row int, column int, cell *xlsx.Cell, stop *bool) {
		if len(cellVal(cell)) > 0 {
			hasContents = true
			*stop = true
		}
	})
	return !hasContents
}

func (self *Range) GetAllColumns(rowOffset int) (values []string) {
	self.EachCell(uint(rowOffset), 1, 0, 0, false, func(row int, column int, cell *xlsx.Cell, stop *bool) {
		values = append(values, cellVal(cell))
	})
	return
}

func (self *Range) Get(rowOffset int, columnOffset int) string {
	return cellVal(self.findRC(rowOffset, columnOffset))
}

func (self *Range) Read(rowOffset int, columnOffset int) (value string, exists bool) {
	c := self.findRC(rowOffset, columnOffset)
	return cellVal(c), c != nil
}

func (self *Range) findRC(rowOffset int, columnOffset int) *xlsx.Cell {
	row := int(self.topRowIndex) + rowOffset
	col := int(self.leftColumnIndex) + columnOffset
	if row >= 0 {
		rows := self.sheet.Rows
		if row < len(rows) {
			if col >= 0 {
				cols := rows[row].Cells
				if col < len(cols) {
					return cols[col]
				}
			}
		}
	}
	return nil
}

func cellVal(cell *xlsx.Cell) string {
	if cell != nil {
		return strings.TrimSpace(strings.ReplaceAll(cell.Value, textMarker, ""))
	} else {
		return ""
	}
}

func (self *Range) RowWithColumnValue(columnOffset int, value string) int {
	for i, row := range self.sheet.Rows {
		columns := row.Cells
		if len(columns) > columnOffset {
			if cellVal(columns[columnOffset]) == value {
				return i
			}
		}
	}
	return -1
}

/*
func xfVal(file *xlsx.File, sheet int, row int, column int) string {
	if file != nil {
		sheets := file.Sheets
		if len(sheets) > sheet {
			rows := sheets[sheet].Rows
			if len(rows) > row {
				columns := rows[row].Cells
				if len(columns) > column {
					return val(columns[column])
				}
			}
		}
	}
	return ""
}

func xfRowWithColumnValue(file *xlsx.File, sheet int, column int, value string) int {
	if file != nil {
		sheets := file.Sheets
		if len(sheets) > sheet {
			for i, row := range sheets[sheet].Rows {
				columns := row.Cells
				if len(columns) > column {
					if val(columns[column]) == value {
						return i
					}
				}
			}
		}
	}
	return -1
}

func xfCell(file *xlsx.File, sheet int, row int, column int) *xlsx.Cell {
	if file != nil && len(file.Sheets) > sheet {
		rows := file.Sheets[sheet].Rows
		if len(rows) > row {
			columns := rows[row].Cells
			if len(columns) > column {
				return columns[column]
			}
		}
	}
	return nil
}

func rowColumnVal(row *xlsx.Row, column int) string {
	if c := rowCell(row, column); c != nil {
		return val(c)
	} else {
		return ""
	}
}

func rowCell(row *xlsx.Row, column int) *xlsx.Cell {
	if column < len(row.Cells) {
		return row.Cells[column]
	} else {
		return nil
	}
}
*/
