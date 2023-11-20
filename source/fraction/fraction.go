package fraction

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type Fraction struct {
	num int64
	den int64
}

func Indeterminate() Fraction {
	return Fraction{num: 0, den: 0}
}

func Zero() Fraction {
	return Fraction{num: 0, den: 1}
}

func NewFraction(num int64, den int64) Fraction {
	return Fraction{num: num, den: den}
}

func fractionFromFloat(f float64) Fraction {
	r := big.NewRat(0, 1).SetFloat64(f)
	return fractionFromRational(r)
}

func fractionFromRational(rational *big.Rat) Fraction {
	return Fraction{num: rational.Num().Int64(), den: rational.Denom().Int64()}
}

func (self Fraction) IsValid() bool {
	return self.den != 0
}

func (self Fraction) Float() float64 {
	switch {
	case self.den != 0:
		return float64(self.num) / float64(self.den)
	case self.num == 0:
		return math.NaN()
	case self.num > 0:
		return math.Inf(+1)
	case self.num < 0:
		return math.Inf(-1)
	}
	panic("logic cannot reach here")
}

func (self Fraction) String() string {
	return fmt.Sprint(self.num) + "/" + fmt.Sprint(self.den)
}

func (self Fraction) rational() *big.Rat {
	return big.NewRat(self.num, self.den)
}

func (self Fraction) Plus(f Fraction) Fraction {
	if self.Float() == 0 {
		return NewFraction(f.num, f.den)
	} else if f.Float() == 0 {
		return f.Plus(self)
	} else if self.den == f.den {
		return NewFraction(self.num+f.num, self.den)
	} else {
		sum := big.NewRat(0, 1).Add(self.rational(), f.rational())
		return fractionFromRational(sum)
	}
}

func Parse(s string) (Fraction, error) {
	split := strings.Split(s, "/")
	if len(split) != 2 ||
		len(strings.TrimSpace(split[0])) == 0 ||
		len(strings.TrimSpace(split[1])) == 0 {
		return Indeterminate(), fmt.Errorf("Fraction not well formed: %s", s)
	}
	if strings.Contains(s, ".") {
		num, err := strconv.ParseFloat(split[0], 64)
		if err != nil {
			return Indeterminate(), err
		}
		den, err := strconv.ParseFloat(split[1], 64)
		if err != nil {
			return Indeterminate(), err
		}
		return fractionFromFloat(num / den), nil
	} else {
		num, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			return Indeterminate(), err
		}
		den, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return Indeterminate(), err
		}
		return NewFraction(num, den), nil
	}
}
