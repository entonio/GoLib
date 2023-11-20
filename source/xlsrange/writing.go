package xlsrange

import (
	"fmt"
	"time"

	"github.com/tealeg/xlsx"
)

var textMarker = "\u200b"

func (self *Range) Printf(format string, values ...any) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.Value = fmt.Sprintf(format, values...)
	})
}

func (self *Range) Text(value any) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.Value = textMarker + fmt.Sprint(value)
	})
}

func (self *Range) Formula(value string) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.SetFormula(value)
	})
}

func (self *Range) Value(value any) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.Value = fmt.Sprint(value)
	})
}

func (self *Range) Date(value time.Time, format string) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.SetDateTime(value)
		cell.NumFmt = format
	})
}

func (self *Range) Float(value any) *Range {
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.SetFloat(0)
		cell.Value = fmt.Sprint(value)
	})
}

func (self *Range) FormatDecimal(decimals uint, separateThousands bool) *Range {
	format := numberFormat(decimals, separateThousands)
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.NumFmt = format
	})
}

func (self *Range) FormatPercent(decimals uint, separateThousands bool) *Range {
	format := numberFormat(decimals, separateThousands) + "%"
	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.NumFmt = format
	})
}

func (self *Range) FormatAccounting() *Range {
	return self.FormatAccountingWith(2, true, "")
}

func (self *Range) FormatAccountingWith(decimals uint, separateThousands bool, symbol string) *Range {
	number := numberFormat(decimals, separateThousands)
	currency := currencyFormat(symbol)
	format := "" +
		"_ * " + number + "_)\\ " + currency + "_ ;" +
		"_ * \\(" + number + "\\)\\ " + currency + "_ ;" +
		"_ * -??_)\\ " + currency + "_ ;" +
		"_ @_ "

	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.NumFmt = format
	})
}

func (self *Range) FormatCurrency() *Range {
	return self.FormatCurrencyWith(2, true, false, "")
}

func (self *Range) FormatCurrencyWith(decimals uint, separateThousands bool, useSymbol bool, symbol string) *Range {
	number := numberFormat(decimals, separateThousands)
	var format string
	if useSymbol {
		currency := currencyFormat(symbol)
		format = number + "\\ " + currency
	} else {
		format = number
	}

	return self.AllCells(func(row int, column int, cell *xlsx.Cell, stop *bool) {
		cell.NumFmt = format
	})
}

func numberFormat(decimals uint, separateThousands bool) string {
	number := "#"
	if separateThousands {
		number += ","
	}
	number += "##0"
	if decimals > 0 {
		number += "."
		for i := uint(0); i < decimals; i += 1 {
			number += "0"
		}
	}
	return number
}

func currencyFormat(symbol string) string {
	if len(symbol) > 0 {
		return "\\\"" + symbol + "\\\""
	} else {
		return "_€"
	}
}
