package macro

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/itdesigndev/conv"
)

func testReplace(t *testing.T, expected string, s string, data interface{}) {

	replaced, err := Replace(s, data)
	if err != nil {
		t.Fatal(err)
	}
	if replaced != expected {
		t.Fatalf("Replace failed:\n\texpected = '%s'\n\t  actual = '%s'", expected, replaced)
	} else {
		//t.Logf("Replace suceed:\n\texpected = '%s'\n\t  actual = '%s'", expected, replaced)
	}

}

func TestReplaceMap(t *testing.T) {
	testReplace(t, "Hello Bar!", "Hello {{.Foo}}!", map[string]string{"Foo": "Bar"})
}

func TestReplaceStruct(t *testing.T) {
	testReplace(t, "Hello Bar!", "Hello {{.Foo}}!", struct{ Foo string }{"Bar"})
}

func TestRemoveNonExisting(t *testing.T) {
	testReplace(t, "Hello !", "Hello {{.Foo2}}!", map[string]string{"Foo": "Bar"})
}

func TestB64Dec(t *testing.T) {
	testReplace(t, "Hello Bar!", "Hello {{.Foo | B64Dec }}!", map[string]string{"Foo": conv.B64Enc("Bar")})
}

func TestB64GenDec(t *testing.T) {
	for _, v := range []string{"b64", "base64"} {
		testReplace(t, "Hello Bar!", "Hello {{.Foo | Dec .Format }}!", map[string]string{"Foo": conv.B64Enc("Bar"), "Format": v})
	}
}

func TestB64Enc(t *testing.T) {
	testReplace(t, "Hello "+conv.B64Enc("Bar")+"!", "Hello {{.Foo | B64Enc }}!", map[string]string{"Foo": "Bar"})
}

func TestB64GenEnc(t *testing.T) {
	for _, v := range []string{"b64", "base64"} {
		testReplace(t, "Hello Bar!", "Hello {{.Foo | Dec .Format }}!", map[string]string{"Foo": conv.B64Enc("Bar"), "Format": v})
	}
}

func TestJsonDec(t *testing.T) {
	var exp map[string]string
	json := `{"foo": "bar", "bar": "foo"}`
	conv.JsonMustDec(json, &exp)
	testReplace(t, "Hello "+fmt.Sprintf("%v", exp)+"!", "Hello {{.Foo | JsonDec }}!", map[string]string{"Foo": json})
}

func TestJsonGenDec(t *testing.T) {
	var exp map[string]string
	json := `{"foo": "bar", "bar": "foo"}`
	conv.JsonMustDec(json, &exp)
	for _, v := range []string{"json", "application/json", "application/x-json", "text/json", "text/x-json"} {
		testReplace(t, "Hello "+fmt.Sprintf("%v", exp)+"!", "Hello {{.Foo | Dec .Format }}!", map[string]string{"Foo": json, "Format": v})
	}
}

func TestJsonEnc(t *testing.T) {
	data := map[string]string{
		"foo": "bar",
		"bar": "foo",
	}
	testReplace(t, "Hello "+conv.JsonMustEnc(data)+"!", "Hello {{.Foo | JsonEnc }}!", map[string]interface{}{"Foo": data})
}

func TestJsonGenEnc(t *testing.T) {
	data := map[string]string{
		"foo": "bar",
		"bar": "foo",
	}
	for _, v := range []string{"json", "application/json", "application/x-json", "text/json", "text/x-json"} {
		testReplace(t, "Hello "+conv.JsonMustEnc(data)+"!", "Hello {{.Foo | Enc .Format }}!", map[string]interface{}{"Foo": data, "Format": v})
	}
}

func TestYamlDec(t *testing.T) {
	var exp = make(map[string]string)
	yaml := `
foo: bar
bar: foo
`
	conv.YamlMustDec(yaml, &exp)
	testReplace(t, "Hello "+fmt.Sprintf("%v", exp)+"!", "Hello {{.Foo | YamlDec }}!", map[string]string{"Foo": yaml})
}

func TestYamlGenDec(t *testing.T) {
	var exp = make(map[string]string)
	yaml := `
foo: bar
bar: foo
`
	conv.YamlMustDec(yaml, &exp)
	for _, v := range []string{"yaml", "application/yaml", "application/x-yaml", "text/yaml", "text/x-yaml"} {
		testReplace(t, "Hello "+fmt.Sprintf("%v", exp)+"!", "Hello {{.Foo | Dec .Format }}!", map[string]string{"Foo": yaml, "Format": v})
	}
}

func TestYamlEnc(t *testing.T) {
	data := map[string]string{
		"foo": "bar",
		"bar": "foo",
	}
	testReplace(t, "Hello "+conv.YamlMustEnc(data)+"!", "Hello {{.Foo | YamlEnc }}!", map[string]interface{}{"Foo": data})
}

func TestYamlGenEnc(t *testing.T) {
	data := map[string]string{
		"foo": "bar",
		"bar": "foo",
	}
	for _, v := range []string{"yaml", "application/yaml", "application/x-yaml", "text/yaml", "text/x-yaml"} {
		testReplace(t, "Hello "+conv.YamlMustEnc(data)+"!", "Hello {{.Foo | Enc .Format }}!", map[string]interface{}{"Foo": data, "Format": v})
	}
}

func TestFunc(t *testing.T) {
	Func("custom", func(s string) string {
		return s + " World!"
	})
	testReplace(t, "Hello World!", `{{ custom "Hello" }}`, nil)
}

func TestFuncs(t *testing.T) {
	Funcs(template.FuncMap{
		"custom1": func(s string) string {
			return s + " World"
		},
		"custom2": func(s string) string {
			return s + "!"
		},
	})
	testReplace(t, "Hello World!", `{{ custom1 "Hello" | custom2 }}`, nil)
}
