package xlsrange

func (self *Range) SelectAll() *Range {
	return self.Select(0, 0, self.MaxRowIndex(), self.MaxColumnIndex())
}

func (self *Range) SelectColumn(columnIndex uint) *Range {
	return self.SelectPosition(self.topRowIndex, columnIndex)
}

func (self *Range) SelectRC(rc RC) *Range {
	return self.SelectPosition(
		rc.Row,
		rc.Col,
	)
}

func (self *Range) SelectPosition(rowIndex uint, columnIndex uint) *Range {
	row := self.SelectRow(rowIndex).rowAtOffset(0)
	for len(row.Cells) <= int(columnIndex) {
		row.AddCell()
	}
	return NewRangeFromCell(self.sheet, rowIndex, columnIndex)
}

func (self *Range) SelectRow(rowIndex uint) *Range {
	for len(self.sheet.Rows) <= int(rowIndex) {
		self.sheet.AddRow()
	}
	return NewRangeFromRow(self.sheet, rowIndex)
}

func (self *Range) SelectField(field Field) *Range {
	return self.Select(
		field.Label.Row,
		field.Label.Col,
		field.Value.Row,
		field.Value.Col,
	)
}

func (self *Range) Select(top uint, left uint, bottom uint, right uint) *Range {
	return newRange(self.sheet,
		top,
		bottom-top+1,
		left,
		right-left+1,
	)
}

func (self *Range) Cut(top uint, left uint, bottom uint, right uint) *Range {
	return newRange(self.sheet,
		self.topRowIndex+top,
		subtract(self.rowCount, top+bottom),
		self.leftColumnIndex+left,
		subtract(self.columnCount, left+right),
	)
}

func (self *Range) RowIndex() uint {
	return self.topRowIndex
}

func (self *Range) ColumnIndex() uint {
	return self.leftColumnIndex
}
