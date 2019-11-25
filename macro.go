package macro

import (
	"bytes"
	"sync"
	"text/template"

	"github.com/itdesigndev/conv"
)

var funcMap = template.FuncMap{
	"B64Dec": conv.B64Dec,
	"B64Enc": conv.B64Enc,
	"JsonDec": func(s string) (v interface{}, err error) {
		err = conv.JsonDec(s, &v)
		return
	},
	"JsonEnc": conv.JsonEnc,
	"YamlDec": func(s string) (v interface{}, err error) {
		err = conv.YamlDec(s, &v)
		return
	},
	"YamlEnc": conv.YamlEnc,
	"Enc":     conv.Enc,
	"Dec": func(format string, src string) (dst interface{}, err error) {
		err = conv.Dec(format, src, &dst)
		return
	},
}
var muFunc sync.RWMutex

// Replace replaces macros in string with data
// Yoy can use full go template functionality, see https://golang.org/pkg/text/template/
func Replace(s string, data interface{}) (out string, err error) {
	var tmpl = template.New("m")
	tmpl.Option("missingkey=zero") // remove non existing macros

	muFunc.RLock()
	tmpl.Funcs(funcMap) // register funcMap
	muFunc.RUnlock()

	if _, err = tmpl.Parse(s); err == nil {
		var b bytes.Buffer
		err = tmpl.Execute(&b, data)
		out = b.String()
	}
	return
}

// Funcs allows to register custom macro functions
// see https://golang.org/pkg/text/template/ for detailed explanation
func Funcs(fm template.FuncMap) {
	muFunc.Lock()
	defer muFunc.Unlock()
	for k, v := range fm {
		funcMap[k] = v
	}
}

// Func same as Funcs but only register one function
func Func(name string, f interface{}) {
	muFunc.Lock()
	defer muFunc.Unlock()
	funcMap[name] = f
}
