package prettier_go

import (
	"dario.cat/mergo"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/dop251/goja"
	"log"
)

//go:embed js/standalone.js
var standalone string

//go:embed js/babel.js
var babel string

//go:embed js/typescript.js
var typescript string

//go:embed js/estree.js
var estree string

//go:embed js/markdown.js
var markdown string

const (
	TypeScript FormatType = "ts"
	JSON       FormatType = "json"
	Markdown   FormatType = "md"
)

type FormatType string

type PrettierFormatOptions struct {
	UseTabs       bool   `json:"useTabs"`
	TabWidth      int    `json:"tabWidth"`
	PrintWidth    int    `json:"printWidth"`
	SingleQuote   bool   `json:"singleQuote"`
	TrailingComma string `json:"trailingComma"`
	Semi          bool   `json:"semi"`
}

var vm *goja.Runtime
var format goja.Callable

var defaults = PrettierFormatOptions{
	UseTabs:       false,
	TabWidth:      4,
	PrintWidth:    80,
	SingleQuote:   false,
	TrailingComma: "none",
	Semi:          false,
}

func ToJsonString[T any](value T) string {
	data, err := json.Marshal(value)
	if err != nil {
		log.Println("ToJsonString error:", err)
		return ""
	}
	return string(data)
}

func handlePromise(value goja.Value, code string) (string, error) {
	if p, ok := value.Export().(*goja.Promise); ok {
		switch p.State() {
		case goja.PromiseStateRejected:
			return code, errors.New(p.Result().String())
		case goja.PromiseStateFulfilled:
			return p.Result().Export().(string), nil
		default:
			return code, errors.New("unexpected promise state pending")
		}
	}

	return "", nil
}

func FormatTypeScript(code string, opts PrettierFormatOptions) (string, error) {
	if err := mergo.Map(&opts, defaults); err != nil {
		log.Println("mergo defaults error:", err)
		return code, err
	}

	log.Println("opts: ", ToJsonString(opts))

	value, err := format(goja.Undefined(), vm.ToValue("typescript"), vm.ToValue(code), vm.ToValue(ToJsonString(opts)))
	if err != nil {
		return code, err
	}

	return handlePromise(value, code)
}
func FormatJSON(code string, opts PrettierFormatOptions) (string, error) {
	if err := mergo.Map(&opts, defaults); err != nil {
		log.Println("mergo defaults error:", err)
		return code, err
	}

	log.Println("opts: ", ToJsonString(opts))

	value, err := format(goja.Undefined(), vm.ToValue("json"), vm.ToValue(code), vm.ToValue(ToJsonString(opts)))
	if err != nil {
		return code, err
	}

	return handlePromise(value, code)
}
func FormatMarkdown(code string, opts PrettierFormatOptions) (string, error) {
	if err := mergo.Map(&opts, defaults); err != nil {
		log.Println("mergo defaults error:", err)
		return code, err
	}

	log.Println("opts: ", ToJsonString(opts))

	value, err := format(goja.Undefined(), vm.ToValue("markdown"), vm.ToValue(code), vm.ToValue(ToJsonString(opts)))
	if err != nil {
		return code, err
	}

	return handlePromise(value, code)
}

var exported = `
async function format(parser, code, optsJsonString) {
	return prettier.format(code, Object.assign(JSON.parse(optsJsonString), {
		parser: parser,
		plugins: prettierPlugins
	}));
}
`

func init() {
	vm = goja.New()
	_, err := vm.RunString(standalone + babel + typescript + estree + markdown + exported)
	if err != nil {
		log.Fatalln("prettier error:", err)
	}

	ret, ok := goja.AssertFunction(vm.Get("format"))
	if !ok {
		panic("not a function")
	} else {
		format = ret
	}

	//err = vm.ExportTo(vm.Get("format"), &format)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("=====")
	//for k, _ := range prettier {
	//	fmt.Println(k)
	//}
	//fmt.Println("=====")
}
