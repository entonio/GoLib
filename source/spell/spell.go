// Package spell is based on an Excel VBA formula meant to convert numbers into
// their textual representation. The formula was offered under free terms in an
// accountant's forum. Later on, it became clear it was a portuguese adaptation
// of Microsoft's SpellNumber formula, that can be downloaded from https://support.microsoft.com/en-us/office/convert-numbers-into-words-a0d166fb-e1ea-4090-95c8-69442cd55d98?ui=en-us&rs=en-us&ad=us. The formula, both in the original and the adaptation, isn't perfect,
// but works well enough to justify porting it to go rather than investing in a
// clean room alternative. Over the years, tweaks have been made to improve parts
// of its behaviour. The ported part is the SpellNum func, the others are convenience
// APIs that apply transforms to its result.
package spell

import (
	"fmt"
	"strings"

	"golib/cs"
	"golib/vba"
)

var Int = vba.Int
var Str = vba.Str

var INSTR = vba.INSTR
var LCASE = vba.LCASE
var LEFT = vba.LEFT
var LEN = vba.LEN
var MID = vba.MID
var REPLACE = vba.REPLACE
var RIGHT = vba.RIGHT
var TRIM = vba.TRIM

func SpellFree(MyNumber any) string {
	t := SpellDecimal(MyNumber)
	result := &t
	result = cs.ReplaceSuffix(result, "vírgula dez", "vírgula um")
	result = cs.ReplaceSuffix(result, "vírgula vinte", "vírgula dois")
	result = cs.ReplaceSuffix(result, "vírgula trinta", "vírgula três")
	result = cs.ReplaceSuffix(result, "vírgula quarenta", "vírgula quatro")
	result = cs.ReplaceSuffix(result, "vírgula cinquenta", "vírgula cinco")
	result = cs.ReplaceSuffix(result, "vírgula sessenta", "vírgula seis")
	result = cs.ReplaceSuffix(result, "vírgula setenta", "vírgula sete")
	result = cs.ReplaceSuffix(result, "vírgula oitenta", "vírgula oito")
	result = cs.ReplaceSuffix(result, "vírgula noventa", "vírgula nove")
	return *result
}

func SpellDecimal(MyNumber any) string {
	t := SpellEur(MyNumber)
	result := &t
	result = REPLACE(result, " euros", " euro")
	result = REPLACE(result, " euro e ", " vírgula ")

	result = REPLACE(result, "vírgula um ", "vírgula zero um ")
	result = REPLACE(result, "vírgula dois ", "vírgula zero dois ")
	result = REPLACE(result, "vírgula três ", "vírgula zero três ")
	result = REPLACE(result, "vírgula quatro ", "vírgula zero quatro ")
	result = REPLACE(result, "vírgula cinco ", "vírgula zero cinco ")
	result = REPLACE(result, "vírgula seis ", "vírgula zero seis ")
	result = REPLACE(result, "vírgula sete ", "vírgula zero sete ")
	result = REPLACE(result, "vírgula oito ", "vírgula zero oito ")
	result = REPLACE(result, "vírgula nove ", "vírgula zero nove ")

	result = REPLACE(result, " euro", "")
	result = REPLACE(result, " cêntimos", "")
	result = REPLACE(result, " cêntimo", "")
	return *result
}

func SpellEur(MyNumber any) string {
	t := SpellNum(MyNumber)
	result := &t
	result = LCASE(result)
	result = REPLACE(result, " e zero cêntimos", "")
	return *result
}

func SpellNum(MyNumber any) string {
	// non-excel preprocessing
	var s_MyNumber = fmt.Sprint(MyNumber)
	s_MyNumber = strings.ReplaceAll(s_MyNumber, " ", "")
	s_MyNumber = strings.ReplaceAll(s_MyNumber, ",", ".")
	s_MyNumber = strings.Trim(s_MyNumber, ".")
	MyNumber = strings.ToLower(s_MyNumber)
	if "" == MyNumber {
		return ""
	}

	// ported from excel

	var Euros = ""
	var Cents = ""
	var temp string
	var Place = make([]string, 8)
	Place[2] = "Mil "
	Place[3] = "Milhões "
	Place[4] = "Mil milhões "
	//Place[5] = "Mil ";
	//Place[6] = "Milhão ";
	//Place[7] = "Mil Milhão ";

	// string representation of amount;
	MyNumber = TRIM(Str(MyNumber))

	// Position of decimal place 0 if none;
	DecimalPlace := *INSTR(MyNumber, ".")

	// Convert cents and set MyNumber to euro amount.;
	if DecimalPlace > 0 {
		Cents = GetTens(LEFT(*MID(MyNumber, DecimalPlace+1, nil)+"00", 2))
		MyNumber = TRIM(LEFT(MyNumber, DecimalPlace-1))
	}

	count := 1
	for *LEN(MyNumber) > 0 {

		// TODO: add to excel
		if count > 4 {
			return ""
		}

		temp = GetHundreds(RIGHT(MyNumber, 3))

		if temp != "" {
			switch count {
			// 1st hundreds: 000.000.000.XXX,00 €
			case 1:
				Euros = temp + Place[count] + Euros
				break

			// 2nd hundreds: 000.000.XXX.000,00 €
			case 2:
				if *Int(RIGHT(MyNumber, 3)) == 1 {
					Euros = Place[count] + Euros
				} else {
					Euros = temp + Place[count] + Euros
				}
				break

			// 3rd hundreds: 000.XXX.000.000,00 €
			case 3:
				if *Int(RIGHT(MyNumber, 3)) == 1 {
					Euros = temp + Place[count+3] + Euros
				} else {
					Euros = temp + Place[count] + Euros
				}
				break

			// 4th hundreds: XXX.000.000.000,00 €
			case 4:
				if *Int(RIGHT(MyNumber, 3)) == 1 {
					Euros = temp + Place[count+3] + Euros
				} else {
					Euros = temp + Place[count] + Euros
				}
				break
			}
		}

		if *LEN(MyNumber) > 3 {
			MyNumber = LEFT(MyNumber, *LEN(MyNumber)-3)
		} else {
			MyNumber = ""
		}

		count = count + 1
	}

	//begin change
	if *LEFT(Euros, 2) == "e " {
		Euros = *RIGHT(Euros, *LEN(Euros)-2)
	}
	//end change

	switch Euros {
	case "":
		Euros = "Zero Euros"
		break
	case "Um ":
		Euros = "Um Euro"
		break
	default:
		Euros = Euros + "Euros"
		break
	}

	switch Cents {
	case "":
		Cents = " e Zero Cêntimos"
		break
	case "Um ":
		Cents = " e Um Cêntimo"
		break
	default:
		Cents = " e " + Cents + "Cêntimos"
		break
	}

	return Euros + Cents
}

// Converts a number from 100-999 into text;
func GetHundreds(MyNumber any) string {
	var result string
	var separator1 string // between hundreds and tens;

	var separator2 string // between tens and digits;

	separator1 = ""
	separator2 = ""

	if *Int(MyNumber) == 0 {
		return ""
	}

	MyNumber = RIGHT("000"+*Str(MyNumber), 3)

	switch *Int(LEFT(MyNumber, 1)) {
	case 1:
		if *Int(MyNumber) == 100 {
			result = "Cem "
		} else {
			result = "Cento "
		}
		break
	case 2:
		result = "Duzentos "
		break
	case 3:
		result = "Trezentos "
		break
	case 4:
		result = "Quatrocentos "
		break
	case 5:
		result = "Quinhentos "
		break
	case 6:
		result = "Seiscentos "
		break
	case 7:
		result = "Setecentos "
		break
	case 8:
		result = "Oitocentos "
		break
	case 9:
		result = "Novecentos "
		break
	default:
		result = ""
		break
	}

	one := 1
	two := 1

	if result != "" && *Int(MID(MyNumber, 2, &one)) > 0 {
		separator1 = "e "
	}

	if *Int(MID(MyNumber, 1, &two)) > 0 && *Int(MID(MyNumber, 3, &one)) > 0 {
		separator2 = "e "
	}

	// Convert the tens and ones place.
	if *MID(MyNumber, 2, &one) != "0" {
		result = result + separator1 + GetTens(MID(MyNumber, 2, nil))
	} else {
		result = result + separator2 + GetDigit(MID(MyNumber, 3, nil))
	}

	//begin change
	if result != "" {
		if (separator1 == "" && separator2 == "") || *Int(LEFT(MyNumber, 1)) == 0 {
			result = "e " + result
		}
	}
	//end change

	return result
}

// Converts a number from 10 to 99 into text.
func GetTens(TensText any) string {

	var result string
	var temp string

	result = "" // Null out the temporary function value.

	if *Int(LEFT(TensText, 1)) == 1 {
		// if ( value between 10-19...

		switch *Int(TensText) {
		case 10:
			result = "Dez "
			break
		case 11:
			result = "Onze "
			break
		case 12:
			result = "Doze "
			break
		case 13:
			result = "Treze "
			break
		case 14:
			result = "Catorze "
			break
		case 15:
			result = "Quinze "
			break
		case 16:
			result = "Dezasseis "
			break
		case 17:
			result = "Dezassete "
			break
		case 18:
			result = "Dezoito "
			break
		case 19:
			result = "Dezanove "
			break
		default:
			break
		}
	} else { // if ( value between 20-99...
		switch *Int(LEFT(TensText, 1)) {
		case 2:
			result = "Vinte "
			break
		case 3:
			result = "Trinta "
			break
		case 4:
			result = "Quarenta "
			break
		case 5:
			result = "Cinquenta "
			break
		case 6:
			result = "Sessenta "
			break
		case 7:
			result = "Setenta "
			break
		case 8:
			result = "Oitenta "
			break
		case 9:
			result = "Noventa "
			break
		default:
			break
		}
		temp = GetDigit(RIGHT(TensText, 1)) // Retrieve ones place

		// use "e" separator between tens and digits, if applicable
		if result != "" && temp != "" {
			result = result + "e " + temp
		} else {
			result = result + temp
		}
	}
	return result
}

// Converts a number from 1 to 9 into text.
func GetDigit(Digit any) string {
	switch *Int(Digit) {
	case 1:
		return "Um "
	case 2:
		return "Dois "
	case 3:
		return "Três "
	case 4:
		return "Quatro "
	case 5:
		return "Cinco "
	case 6:
		return "Seis "
	case 7:
		return "Sete "
	case 8:
		return "Oito "
	case 9:
		return "Nove "
	default:
		return ""
	}
}
