package settings

import (
	"github.com/zieckey/goini"
)

type ini struct {
	path string
	ini  *goini.INI
}

func newIni(path string) *ini {
	goini := goini.New()
	err := goini.ParseFile(path)
	assertNil(err)
	return &ini{
		path: path,
		ini:  goini,
	}
}

func (self *ini) B(section string, key string) bool {
	v, ok := self.ini.SectionGetBool(section, key)
	assert(ok, self.path+": could not get "+section+"."+key+" bool")
	return v
}

func (self *ini) I(section string, key string) int {
	v, ok := self.ini.SectionGetInt(section, key)
	assert(ok, self.path+": could not get "+section+"."+key+" int")
	return v
}

func (self *ini) F(section string, key string) float64 {
	v, ok := self.ini.SectionGetFloat(section, key)
	assert(ok, self.path+": could not get "+section+"."+key+" float")
	return v
}

func (self *ini) S(section string, key string) string {
	v, ok := self.ini.SectionGet(section, key)
	assert(ok, self.path+": could not get "+section+"."+key+" string")
	return v
}

func (self *ini) SL(section string, key string) []string {
	return SL(self.S(section, key))
}

func (self *ini) C(section string, key string) string {
	return C(self.S(section, key))
}
