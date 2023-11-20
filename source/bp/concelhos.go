package bp

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"golib/fsx"
	"golib/lang"
	"golib/maps"
)

type Concelhos struct {
	distritosByConcelho map[key]string
	concelhosByLocation map[key]string
	sinonimosByLocation map[key]string
	keysSortedByLength  []key
}

func NewConcelhos(file string) *Concelhos {
	self := &Concelhos{
		distritosByConcelho: make(map[key]string),
		concelhosByLocation: make(map[key]string),
		sinonimosByLocation: make(map[key]string),
	}
	contents := fsx.MustRead(file)
	regexpLineDistrito := regexp.MustCompile(`^DISTRITO\s*([^\s].*)$`)
	regexpLineConcelho := regexp.MustCompile(`^CONCELHO\s*([^\s].*)$`)
	regexpLineSinonimo := regexp.MustCompile(`^SINONIMO\s*([^\s].*)$`)
	regexpLineFlags := regexp.MustCompile(`^(.+)\s*_(.*)_$`)
	var distrito string
	var concelho string
	var sinonimo string
	for index, line := range strings.Split(contents, "\n") {
		lineNumber := fmt.Sprintf("%s:%d: ", file, index+1)

		line := strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "--") {
			continue
		}

		match := regexpLineDistrito.FindStringSubmatch(line)
		if len(match) == 2 {
			distrito = match[1]
			concelho = ""
			sinonimo = ""
			continue
		}

		match = regexpLineConcelho.FindStringSubmatch(line)
		if len(match) == 2 {
			distrito = ""
			concelho = match[1]
			sinonimo = ""
			continue
		}

		match = regexpLineSinonimo.FindStringSubmatch(line)
		if len(match) == 2 {
			distrito = ""
			concelho = ""
			sinonimo = match[1]
			continue
		}

		switch {
		case len(distrito) > 0:
			line, flags := self.splitLineFromFlags(regexpLineFlags, line)
			if !strings.Contains(flags, "N") {
				self.addLocationToConcelho(line, line)
			}
			self.addConcelhoToDistrito(line, distrito)
		case len(concelho) > 0:
			self.addLocationToConcelho(line, concelho)
		case len(sinonimo) > 0:
			self.addNameToSinonimo(line, sinonimo)
		default:
			assert(false, lineNumber+"Line belongs neither to a distrito nor a concelho nor a sinonimo: "+line)
		}
	}

	keys := maps.Keys(self.concelhosByLocation)
	sort.Slice(keys, func(i, j int) bool {
		return len(string(keys[i])) > len(string(keys[j]))
	})
	self.keysSortedByLength = keys

	return self
}

func (self *Concelhos) splitLineFromFlags(regexpLineFlags *regexp.Regexp, lineAndFlags string) (line string, flags string) {
	match := regexpLineFlags.FindStringSubmatch(lineAndFlags)
	if len(match) == 3 {
		return match[1], strings.ToUpper(match[2])
	} else {
		return lineAndFlags, ""
	}
}

func (self *Concelhos) addNameToSinonimo(name string, sinonimo string) {
	key := keyFromLocation(name)
	existing := self.sinonimosByLocation[key]
	assert(len(existing) == 0, name+" already mapped to "+existing)
	self.sinonimosByLocation[key] = sinonimo
}

func (self *Concelhos) addLocationToConcelho(location string, concelho string) {
	key := keyFromLocation(location)
	existing := self.concelhosByLocation[key]
	assert(len(existing) == 0, location+" already mapped to "+existing)
	self.concelhosByLocation[key] = concelho
}

func (self *Concelhos) addConcelhoToDistrito(concelho string, distrito string) {
	key := keyFromLocation(concelho)
	existing := self.distritosByConcelho[key]
	assert(len(existing) == 0, concelho+" already mapped to "+existing)
	self.distritosByConcelho[key] = distrito
}

func (self *Concelhos) SinonimoForName(name string) string {
	key := keyFromLocation(name)
	return self.sinonimosByLocation[key]
}

func (self *Concelhos) DistritoForConcelho(concelho string) string {
	key := keyFromLocation(concelho)
	return self.distritosByConcelho[key]
}

func (self *Concelhos) ConcelhoForExactLocation(location string) string {
	key := keyFromLocation(location)
	return self.concelhosByLocation[key]
}

func (self *Concelhos) ConcelhoForPartialLocation(location string) string {
	key := keyFromLocation(location)
	/*
		if strings.Contains(location, "ogueira") {
			fmt.Printf("LOCATION [%s] KEY [%s]\n", location, key)
		}
		if strings.HasSuffix(location, "Lisboa") {
		}
	*/
	for _, candidate := range self.keysSortedByLength {
		if strings.HasSuffix(string(key), string(candidate)) {
			return self.concelhosByLocation[candidate]
		}
	}
	return ""
}

func (self *Concelhos) ConcelhoForLocation(location string) (concelho string) {
	concelho = self.ConcelhoForExactLocation(location)
	if len(concelho) == 0 {
		concelho = self.ConcelhoForPartialLocation(location)
	}
	return
}

type key string

func keyFromLocation(location string) key {
	result := location
	result = lang.KeyFromString(result)
	/*
		if strings.Contains(location, "zeit") {
			fmt.Println("LOCATION [" + location + "] -> [" + result + "]")
		}
	*/
	return key(result)
}

func CorrectConcelho(concelho string) string {
	switch concelho {
	case "Vila Velha de Rodão":
		return "Vila Velha de Ródão"
	}
	return concelho
}
