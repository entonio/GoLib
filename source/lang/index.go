package lang

type index struct {
	current int
}

func NewIndex() *index {
	return &index{current: -1}
}

func (self *index) Set(value int) int {
	self.current = value
	return self.current
}

func (self *index) Same() int {
	return self.current
}

func (self *index) Next() int {
	self.current = self.current + 1
	return self.current
}

func (self *index) Skip(offset int) int {
	self.current = self.current + 1 + offset
	return self.current
}
