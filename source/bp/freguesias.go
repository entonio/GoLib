package bp

import (
	"fmt"
	"sort"
	"strings"

	"golib/arrays"
	"golib/lang"
	"golib/log"
	"golib/maps"
	"golib/text"
	"golib/xlsrange"

	"github.com/tealeg/xlsx"
)

type CodigoFreguesia struct {
	codigo string
}

func NewCodigoFreguesia(codigo string) CodigoFreguesia {
	codigo = strings.ToUpper(strings.ReplaceAll(codigo, " ", ""))
	if len(codigo) == 5 {
		codigo = "0" + codigo
	}
	return CodigoFreguesia{
		codigo: codigo,
	}
}

func (self CodigoFreguesia) String() string {
	return self.codigo
}

func (self CodigoFreguesia) SeemsValid() bool {
	return len(self.codigo) == 6
}

type Freguesia struct {
	concelho        string
	nome            string
	codigo          CodigoFreguesia
	nomeAlternativo string
	extinta         bool
	uniao           bool

	codigoNova  CodigoFreguesia
	substitutas []*Freguesia
	nomeSD      string
	substituta  *Freguesia
}

func (self *Freguesia) actual() *Freguesia {
	if self.substituta != nil {
		return self.substituta
	}
	return self
}

func (self *Freguesia) codigoConcelho() string {
	return self.codigo.codigo[0:4]
}

func (self *Freguesia) concelhoIs(s string) bool {
	c := lang.RemoveDiacritics(self.concelho)
	s = lang.RemoveDiacritics(s)
	return strings.EqualFold(c, s)
}

func (self *Freguesia) NomeSD() string {
	if len(self.nomeSD) == 0 {
		self.nomeSD = lang.RemoveDiacritics(self.nome)
	}
	return self.nomeSD
}

func (self *Freguesia) nomes() string {
	return self.nomeReplacingUFBy("")
}

func (self *Freguesia) nomeReplacingUFBy(subst string) string {
	s := self.nome
	s = strings.Replace(s, "União de Freguesias", "União das Freguesias", 1)
	s = strings.Replace(s, "União das Freguesias dos ", subst, 1)
	s = strings.Replace(s, "União das Freguesias das ", subst, 1)
	s = strings.Replace(s, "União das Freguesias do ", subst, 1)
	s = strings.Replace(s, "União das Freguesias da ", subst, 1)
	s = strings.Replace(s, "União das Freguesias de ", subst, 1)
	return s
}

func (self *Freguesia) nomeSDP() string {
	n := self.NomeSD()
	if strings.HasSuffix(n, ")") && strings.Contains(n, "(") {
		split := strings.SplitN(n, "(", 2)
		p1 := strings.TrimSpace(split[0])
		p2 := strings.TrimSpace(split[1][:len(split[1])-1])
		if self.concelhoIs(p1) {
			return p2
		}
		if self.concelhoIs(p2) {
			return p1
		}
		log.Debug("[%s] [%s] [%s]", self.concelho, p1, p2)
	}
	return n
}

type Freguesias struct {
	freguesiasByCodigo map[CodigoFreguesia]*Freguesia
}

const (
	XL_Distrito = iota
	XL_Concelho
	XL_Freguesia
	XL_Codigo
	XL_NomeAlternativo
	XL_CodSF
	XL_NomeSF
	XL_CodigoNova
)

/*
func val(cells []*xlsx.Cell, index int) string {
	if len(cells) > index {
		return strings.TrimSpace(cells[index].Value)
	} else {
		return ""
	}
}
*/

func NewFreguesias(path string) *Freguesias {
	self := &Freguesias{
		freguesiasByCodigo: make(map[CodigoFreguesia]*Freguesia),
	}

	file, err := xlsx.OpenFile(path)
	assertNil(err)

	xls := xlsrange.NewRangeFromSheet(file.Sheets[0])

	xls.EachRow(1, 0, true, func(row int, xr *xlsx.Row, stop *bool) {

		freguesia := &Freguesia{
			concelho:        xls.Get(row, XL_Concelho),
			nome:            xls.Get(row, XL_Freguesia),
			codigo:          NewCodigoFreguesia(xls.Get(row, XL_Codigo)),
			nomeAlternativo: xls.Get(row, XL_NomeAlternativo),
			extinta:         len(xls.Get(row, XL_CodSF)) == 0,
			codigoNova:      NewCodigoFreguesia(xls.Get(row, XL_CodigoNova)),
		}

		freguesia.nome = normaliseNome(freguesia.nome)
		freguesia.nomeAlternativo = normaliseNome(freguesia.nomeAlternativo)

		freguesia.uniao = strings.HasPrefix(freguesia.nome, "Uni")

		self.freguesiasByCodigo[freguesia.codigo] = freguesia
	})

	for _, c := range self.ListCodigos() {
		f := self.freguesiasByCodigo[c]
		if len(f.codigoNova.codigo) > 0 {
			f.substituta = self.freguesiasByCodigo[f.codigoNova]
		}
	}

	log.Debug("%s: lidas %d freguesias", path, len(self.freguesiasByCodigo))
	return self
}

func normaliseNome(nome string) string {
	nome = strings.Replace(nome, "(Freguesia Extinta)", " ", -1)
	nome = strings.Replace(nome, "(Extinta)", " ", -1)
	nome = strings.ReplaceAll(nome, "São ", "S. ")
	nome = strings.ReplaceAll(nome, "Setubal", "Setúbal")
	return strings.TrimSpace(nome)
}

func (self *Freguesias) byCodigo(codigo string) *Freguesia {
	c := NewCodigoFreguesia(codigo)
	return self.freguesiasByCodigo[c]
}

func (self *Freguesias) substitutas(f *Freguesia) []*Freguesia {
	if !f.extinta {
		return nil
	}
	if len(f.substitutas) == 0 {
		n := f.nomeSDP()
		c := f.codigoConcelho()
		var substitutas []*Freguesia
		for _, s := range self.freguesiasByCodigo {
			if s != f && !s.extinta && s.codigoConcelho() == c {
				sn := s.NomeSD()
				if sn == f.NomeSD() || sn == n ||
					strings.Contains(sn, n+",") ||
					strings.Contains(sn, n+" e ") ||
					strings.HasSuffix(sn, ", "+n+")") ||
					strings.HasSuffix(sn, " e "+n+")") ||
					strings.HasSuffix(sn, ", "+n) ||
					strings.HasSuffix(sn, " e "+n) {
					substitutas = append(substitutas, s)
				}
			}
		}
		sort.Slice(substitutas, func(i1, i2 int) bool {
			s1 := substitutas[i1]
			s2 := substitutas[i2]
			if s1.uniao && !s2.uniao {
				return true
			}
			if !s1.uniao && s2.uniao {
				return false
			}
			return s1.codigo.codigo < s2.codigo.codigo
		})
		f.substitutas = substitutas
		if len(substitutas) == 0 {
			log.Trace("[%s:%s] -> [%s]", f.concelho, f.nome, n)
		}
	}
	return f.substitutas
}

func (self *Freguesias) CSVCalcularSubstitutas() string {
	var csv []string
	for _, c := range self.ListCodigos() {
		f := self.freguesiasByCodigo[c]
		line := fmt.Sprintf("%s\t%s",
			f.codigo.codigo,
			f.nome,
		)
		for _, s := range self.substitutas(f) {
			line += fmt.Sprintf("\t%s\t%s",
				s.codigo.codigo,
				s.nome,
			)
		}
		csv = append(csv, line)
	}
	return strings.Join(csv, "\n")
}

func (self *Freguesias) ListCodigos() []CodigoFreguesia {
	codigos := maps.Keys(self.freguesiasByCodigo)
	sort.Slice(codigos, func(i1, i2 int) bool {
		return codigos[i1].codigo < codigos[i2].codigo
	})
	return codigos
}

func (self *Freguesias) CodigoValido(codigo string) bool {
	c := NewCodigoFreguesia(codigo)
	return c.SeemsValid()
}

/*
func isConsonant(c rune) bool {
	switch c {
	case 'b', 'c', 'ç', 'd', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'p', 'q', 'r', 's', 't', 'v', 'w', 'x', 'z':
		return true
	case 'B', 'C', 'Ç', 'D', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'V', 'W', 'X', 'Z':
		return true
	default:
		return false
	}
}
*/

func (self *Freguesias) Abreviatura(codigo string) (abreviatura string) {
	if codigo == "151210" {
		return "União"
	}
	if f := self.byCodigo(codigo); f != nil {
		n := f.nome
		n = strings.ReplaceAll(n, "-", " ")
		n = strings.ReplaceAll(n, " de ", " ")
		n = strings.ReplaceAll(n, " da ", " ")
		n = strings.ReplaceAll(n, " das ", " ")
		n = strings.ReplaceAll(n, " do ", " ")
		n = strings.ReplaceAll(n, " dos ", " ")
		n = strings.ReplaceAll(n, "União Freguesias", "UF")
		n = strings.ReplaceAll(n, "Nossa", "N.")
		n = strings.ReplaceAll(n, "Sra", "S.")
		n = strings.ReplaceAll(n, "Srª", "S.")
		n = strings.ReplaceAll(n, "Senhora", "S.")
		n = strings.ReplaceAll(n, "N. S.", "NS")
		n = strings.ReplaceAll(n, "Santa", "St.")
		n = strings.ReplaceAll(n, "Santo", "St.")
		n = strings.ReplaceAll(n, "São", "S.")

		words := strings.Split(n, " ")
		if len(words) == 1 {
			abreviatura = words[0]
		} else {
			keep := arrays.New("Setúbal", "Azeitão")
			a := arrays.Map(words, func(s string) string {
				var prefix string
				var suffix string
				if s[0] == '(' {
					prefix = "("
					s = s[1:]
				}
				if s[len(s)-1] == ')' {
					suffix = ")"
					s = s[:len(s)-1]
				}
				if s != strings.ToUpper(s) && !arrays.Contains(keep, s) {
					silabas := text.Syllables(s)
					if len(silabas) > 2 || (len(silabas) == 2 && len(silabas[1]) >= 3) {
						s1c := text.ConsonantalPrefix(silabas[1])
						//if s == "Graça" {
						//	fmt.Println(silabas)
						//}
						if len(silabas) > 2 || len(s1c) <= 2 {
							sp := silabas[0] + s1c + "."
							if len(sp) < len(s) {
								if suffix == "" {
									suffix = " "
								}
								return prefix + sp + suffix
							}
							/*
								if len(a) == 2 {
									return a
								} else {

								}
							*/
						}
					}
				}
				return prefix + s + suffix + " "
			})
			abreviatura = strings.TrimRight(strings.Join(a, ""), " ")
		}
		/*
			C := false
			V := false
			VC := false
			VCV := false
			for _, c := range n {
				if unicode.IsLetter(c) {
					if VCV {
						continue
					}
					if !isConsonant(c) {
						if VC {
							VCV = true
							if !C {
								continue
							}
						} else if !C {
							V = true
						}
					} else {
						if V {
							VC = true
						} else {
							C = true
						}
					}
					sigla += string(c)
					continue
				}
				C = false
				V = false
				VC = false
				VCV = false
				if c == '(' {
					sigla += " ("
				} else if c == ')' {
					sigla += string(c)
				} else if c == ' ' {
					sigla += string(c)
				}
			}
			sigla = strings.ReplaceAll(sigla, "  ", " ")
			sigla = strings.ReplaceAll(sigla, "N ", "N")
			sigla = strings.ReplaceAll(sigla, "S ", "S")
			if len(sigla) < 2 {
				sigla = f.nome
			}
		*/
	}
	return
}

/*
func (self *Freguesias) Sigla(codigo string) (sigla string) {
	if f := self.byCodigo(codigo); f != nil {
		n := f.nome
		for _, c := range n {
			if c == '(' {
				sigla += " ("
			} else if c == ')' {
				sigla += string(c)
			} else if unicode.IsUpper(c) {
				sigla += string(c)
			}
		}
		/*
			if f.uniao {
				sigla = strings.TrimPrefix(sigla, "UF")
			}
		* /
		if len(sigla) < 2 {
			sigla = n
			/*
				max := 5
				if len(n) > max {
					sigla = n[:max]
				} else {
					sigla = n
				}
			* /
		}
	}
	return
}
*/

func (self *Freguesias) NomesCompactos(codigo string) string {
	if f := self.byCodigo(codigo); f != nil {
		result := f.nomeReplacingUFBy("UF ")
		if len(f.nomeAlternativo) > 0 {
			result = f.nomeAlternativo + " = " + result
		}
		return result
	}
	return ""
}

func (self *Freguesias) NomeIncludesAlternativo(nome string) bool {
	return strings.Contains(nome, " = ")
}

func (self *Freguesias) IsConhecida(codigo string) bool {
	return self.byCodigo(codigo) != nil
}

func (self *Freguesias) IsExtinta(codigo string) bool {
	if f := self.byCodigo(codigo); f != nil {
		return f.extinta
	}
	panic("No freguesia found for codigo [" + codigo + "]")
}

func (self *Freguesias) IsUniao(codigo string) bool {
	if f := self.byCodigo(codigo); f != nil {
		return f.uniao
	}
	panic("No freguesia found for codigo [" + codigo + "]")
}

func (self *Freguesias) CodigoActual(codigo string) string {
	if f := self.byCodigo(codigo); f != nil {
		return f.actual().codigo.String()
	}
	return codigo
}

func (self *Freguesias) TestAbreviaturas() {
	test := func(codigo string) {
		fmt.Println(codigo + " " + self.Abreviatura(codigo))
		fmt.Println("       " + self.NomesCompactos(codigo))
	}
	for f := 1; f < 40; f += 1 {
		test(fmt.Sprintf("0603%02d", f))
	}
	for f := 1; f < 10; f += 1 {
		test(fmt.Sprintf("1512%02d", f))
	}
}
