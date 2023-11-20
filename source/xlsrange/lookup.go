package xlsrange

import "fmt"

func (self *Range) Matches(top uint, left uint, bottom uint, right uint, values []string) bool {
	i := -1
	for row := top; row <= bottom; row += 1 {
		for col := left; col <= right; col += 1 {
			i += 1
			required := ""
			if i < len(values) {
				required = values[i]
			}
			present := self.Get(int(row), int(col))
			fmt.Printf("%d,%d: present [%s], required [%s]\n", row, col, present, required)
			if present != required {
				return false
			}
		}
	}
	return true
}
