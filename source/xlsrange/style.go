package xlsrange

import (
	"golib/xlsconstants"

	"github.com/tealeg/xlsx"
)

func (self *Range) HeightMulti(multiplier float64) *Range {
	return self.HeightPoints(multiplier * xlsconstants.DefaultHeightPoints)
}

func (self *Range) HeightCM(centimeters float64) *Range {
	return self.HeightIN(centimeters / 2.54)
}

func (self *Range) HeightIN(inches float64) *Range {
	return self.HeightPoints(inches * 72)
}

func (self *Range) HeightPoints(points float64) *Range {
	return self.AllRows(func(row int, xr *xlsx.Row, stop *bool) {
		xr.SetHeightCM(0) // indirect way to set row.isCustom = true
		xr.Height = points
	})
}

/*
func (self *Range) HeightReset() *Range {
	return self.allRows(func(row *xlsx.Row) {
		row.isCustom = false
	})
}
*/

func (self *Range) Border(style xlsconstants.BorderStyle, color xlsconstants.Color) *Range {
	self.BorderBottom(style, color)
	self.BorderTop(style, color)
	self.BorderRight(style, color)
	self.BorderLeft(style, color)
	return self
}

func (self *Range) BorderStyle(style xlsconstants.BorderStyle) *Range {
	self.BorderBottomStyle(style)
	self.BorderTopStyle(style)
	self.BorderRightStyle(style)
	self.BorderLeftStyle(style)
	return self
}

func (self *Range) BorderColor(color xlsconstants.Color) *Range {
	self.BorderBottomColor(color)
	self.BorderTopColor(color)
	self.BorderRightColor(color)
	self.BorderLeftColor(color)
	return self
}

func (self *Range) BorderLeft(style xlsconstants.BorderStyle, color xlsconstants.Color) *Range {
	return self.LeftCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		s := cell.GetStyle()
		s.Border.Left = string(style)
		s.Border.LeftColor = string(color)
	})
}

func (self *Range) BorderLeftStyle(style xlsconstants.BorderStyle) *Range {
	return self.LeftCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Border.Left = string(style)
	})
}

func (self *Range) BorderLeftColor(color xlsconstants.Color) *Range {
	return self.LeftCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Border.LeftColor = string(color)
	})
}

func (self *Range) BorderRight(style xlsconstants.BorderStyle, color xlsconstants.Color) *Range {
	return self.RightCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		s := cell.GetStyle()
		s.Border.Right = string(style)
		s.Border.RightColor = string(color)
	})
}

func (self *Range) BorderRightStyle(style xlsconstants.BorderStyle) *Range {
	return self.RightCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Border.Right = string(style)
	})
}

func (self *Range) BorderRightColor(color xlsconstants.Color) *Range {
	return self.RightCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Border.RightColor = string(color)
	})
}

func (self *Range) BorderTop(style xlsconstants.BorderStyle, color xlsconstants.Color) *Range {
	return self.topCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		s := cell.GetStyle()
		s.Border.Top = string(style)
		s.Border.TopColor = string(color)
	})
}

func (self *Range) BorderTopStyle(style xlsconstants.BorderStyle) *Range {
	return self.topCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Border.Top = string(style)
	})
}

func (self *Range) BorderTopColor(color xlsconstants.Color) *Range {
	return self.topCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Border.TopColor = string(color)
	})
}

func (self *Range) BorderBottom(style xlsconstants.BorderStyle, color xlsconstants.Color) *Range {
	return self.bottomCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		s := cell.GetStyle()
		s.Border.Bottom = string(style)
		s.Border.BottomColor = string(color)
	})
}

func (self *Range) BorderBottomStyle(style xlsconstants.BorderStyle) *Range {
	return self.bottomCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Border.Bottom = string(style)
	})
}

func (self *Range) BorderBottomColor(color xlsconstants.Color) *Range {
	return self.bottomCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Border.BottomColor = string(color)
	})
}

func (self *Range) FillFgColor(color xlsconstants.Color) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Fill.FgColor = string(color)
	})
}

func (self *Range) FillPatternType(value xlsconstants.Pattern) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Fill.PatternType = string(value)
	})
}

func (self *Range) FontColorGray() *Range {
	return self.FontColor(xlsconstants.Gray)
}

func (self *Range) FontColorRed() *Range {
	return self.FontColor(xlsconstants.Red)
}

func (self *Range) FontColorMagenta() *Range {
	return self.FontColor(xlsconstants.Magenta)
}

func (self *Range) FontColorLightBlue() *Range {
	return self.FontColor(xlsconstants.LightBlue)
}

func (self *Range) FontColorOrange() *Range {
	return self.FontColor(xlsconstants.Orange)
}

func (self *Range) FontColorGreen() *Range {
	return self.FontColor(xlsconstants.Green)
}

func (self *Range) FontColorWhite() *Range {
	return self.FontColor(xlsconstants.White)
}

func (self *Range) FontColor(color xlsconstants.Color) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Font.Color = string(color)
	})
}

func (self *Range) FontFamily(font xlsconstants.Font) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Font.Name = string(font)
	})
}

func (self *Range) FontBold() *Range {
	return self.setFontBold(true)
}

func (self *Range) FontNotBold() *Range {
	return self.setFontBold(false)
}

func (self *Range) setFontBold(b bool) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Font.Bold = b
	})
}

func (self *Range) AlignTop() *Range {
	return self.Align(xlsconstants.Top)
}

func (self *Range) AlignMiddle() *Range {
	return self.Align(xlsconstants.Middle)
}

func (self *Range) AlignBottom() *Range {
	return self.Align(xlsconstants.Bottom)
}

func (self *Range) AlignLeft() *Range {
	return self.Align(xlsconstants.Left)
}

func (self *Range) AlignCenter() *Range {
	return self.Align(xlsconstants.Center)
}

func (self *Range) AlignRight() *Range {
	return self.Align(xlsconstants.Right)
}

func (self *Range) AlignLeftMiddle() *Range {
	return self.Align(xlsconstants.Left, xlsconstants.Middle)
}

func (self *Range) AlignCenterMiddle() *Range {
	return self.Align(xlsconstants.Center, xlsconstants.Middle)
}

func (self *Range) AlignRightMiddle() *Range {
	return self.Align(xlsconstants.Right, xlsconstants.Middle)
}

func (self *Range) Wrap() *Range {
	return self.setWrap(true)
}

func (self *Range) setWrap(b bool) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.GetStyle().Alignment.WrapText = b
	})
}

func (self *Range) Align(alignments ...xlsconstants.Alignment) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		for _, alignment := range alignments {
			switch alignment {
			case xlsconstants.Left,
				xlsconstants.Center,
				xlsconstants.Right:
				cell.GetStyle().Alignment.Horizontal = string(alignment)

			case xlsconstants.Middle:
				cell.GetStyle().Alignment.Vertical = string(xlsconstants.Center)

			case xlsconstants.Top,
				xlsconstants.Bottom:
				cell.GetStyle().Alignment.Vertical = string(alignment)
			}
		}
	})
}
