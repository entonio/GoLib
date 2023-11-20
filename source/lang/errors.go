package lang

import (
	"fmt"

	"github.com/go-errors/errors"
)

type Errors struct {
	Map map[any]error
}

func (self Errors) Found() bool {
	return self.Len() > 0
}

func (self Errors) Len() int {
	return len(self.Map)
}

func (self Errors) String() string {
	return fmt.Sprint(self.Map)
}

func (self Errors) AsError(s string) error {
	return errors.New(s)
}

func (self Errors) Add(e error) bool {
	return self.Set(fmt.Sprintf("%d", len(self.Map)), e)
}

func (self Errors) Set(key any, e error) bool {
	if e == nil {
		return false
	}
	if self.Map == nil {
		self.Map = make(map[any]error)
	}
	self.Map[key] = e
	return true
}

func (self Errors) AddAll(errs Errors) {
	if self.Map == nil {
		self.Map = make(map[any]error)
	}
	for key, e := range errs.Map {
		self.Map[key] = e
	}
}
