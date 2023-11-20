package errs

import "errors"

type Errors struct {
	list []error
}

func (self Errors) Add(err error) {
	if err != nil {
		self.list = append(self.list, err)
	}
}

func (self Errors) Combined() error {
	return errors.Join(self.list...)
}
