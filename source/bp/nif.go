package bp

import (
	"fmt"
	"regexp"
	"strings"
)

type NIF struct {
	original string
	digits   string
}

var regexpNonDigits = regexp.MustCompile(`[^\d]+`)

func NewNIF(nif string) NIF {
	return NIF{
		original: nif,
		digits:   regexpNonDigits.ReplaceAllString(nif, ""),
	}
}

type TipoNIF string

const (
	NIF_pessoa_singular                   TipoNIF = "Pessoa Singular"
	NIF_pessoa_singular_nao_residente     TipoNIF = "Pessoa Singular Não Residente"
	NIF_pessoa_colectiva                  TipoNIF = "Pessoa Colectiva"
	NIF_administracao_publica             TipoNIF = "Administração Pública"
	NIF_heranca_indivisa                  TipoNIF = "Herança Indivisa"
	NIF_pessoa_colectiva_nao_residente    TipoNIF = "Pessoa Colectiva Não Residente"
	NIF_fundo_de_investimento             TipoNIF = "Fundo de Investimento"
	NIF_atribuicao_oficiosa               TipoNIF = "Atribuição Oficiosa"
	NIF_regime_excepcional                TipoNIF = "Regime Excepcional"
	NIF_empresario_em_nome_individual     TipoNIF = "Empresário em Nome Individual (Extinto)"
	NIF_condominio_ou_sociedade_irregular TipoNIF = "Condomínio ou Sociedade Irregular"
	NIF_nao_residente                     TipoNIF = "Não Residente"
	NIF_sociedade_civil                   TipoNIF = "Sociedade Civil"
)

var TipoNIFByPrefix = map[uint]TipoNIF{
	1:  NIF_pessoa_singular,
	2:  NIF_pessoa_singular,
	3:  NIF_pessoa_singular,
	45: NIF_pessoa_singular_nao_residente,
	5:  NIF_pessoa_colectiva,
	6:  NIF_administracao_publica,
	70: NIF_heranca_indivisa,
	74: NIF_heranca_indivisa,
	71: NIF_pessoa_colectiva_nao_residente,
	72: NIF_fundo_de_investimento,
	77: NIF_atribuicao_oficiosa,
	79: NIF_regime_excepcional,
	8:  NIF_empresario_em_nome_individual,
	90: NIF_condominio_ou_sociedade_irregular,
	91: NIF_condominio_ou_sociedade_irregular,
	98: NIF_nao_residente,
	99: NIF_sociedade_civil,
}

func (self NIF) String() string {
	return self.digits
}

func (self NIF) Tipo() TipoNIF {
	for prefix, tipo := range TipoNIFByPrefix {
		if strings.HasPrefix(self.digits, fmt.Sprint(prefix)) {
			return tipo
		}
	}
	var tipo TipoNIF
	return tipo
}

func (self NIF) IsValid() bool {
	return ChecksumIsValid(self.digits, 9, 2)
}

/*
Multiplique o 8.º dígito por 2, o 7.º dígito por 3, o 6.º dígito por 4, o 5.º dígito por 5, o 4.º dígito por 6, o 3.º dígito por 7, o 2.º dígito por 8, e o 1.º digito por 9
Adicione os resultados
Calcule o Módulo 11 do resultado, isto é, o resto da divisão do número por 11.
Se o resto for 0 ou 1, o dígito de controle será 0
Se for outro algarismo x, o dígito de controle será o resultado de 11 - x
*/
func ChecksumIsValid(s string, n int, t int) bool {
	if len(s) != n {
		//log.Debug("len(s) != n: %d %d", len(s), n)
		return false
	}

	sum := 0
	for i, r := range s {
		d := int(r - '0')
		if d < 0 || d > 9 {
			//log.Debug("d < 0 || d > 9: %d", d)
			return false
		}
		if i < n-1 {
			//log.Debug("%d: %d * %d", i+1, d, n-i)
			sum += d * (n - i)
		} else {
			mod := sum % (n + t)
			if mod < t {
				//log.Debug("d == 0?: %d", d)
				return d == 0
			} else {
				x := 11 - mod
				//log.Debug("d == x?: %d %d", d, x)
				return d == x
			}
		}
	}

	panic("unreachable")
}

func (self NIF) IsIndividual() bool {
	switch self.Tipo() {
	case NIF_pessoa_singular,
		NIF_pessoa_singular_nao_residente,
		NIF_heranca_indivisa:
		return true
	default:
		return false
	}
}

func (self NIF) IsHeranca() bool {
	switch self.Tipo() {
	case NIF_heranca_indivisa:
		return true
	default:
		return false
	}
}

func (self NIF) IsFundo() bool {
	switch self.Tipo() {
	case NIF_fundo_de_investimento:
		return true
	default:
		return false
	}
}

func (self NIF) IsEstatal() bool {
	switch self.Tipo() {
	case NIF_administracao_publica:
		return true
	default:
		return false
	}
}

func (self NIF) IsColectivo() bool {
	return len(self.digits) > 0 && !self.IsIndividual()
}
