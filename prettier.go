package prettier_go

import (
	"dario.cat/mergo"
	_ "embed"
	"encoding/json"
	"errors"
	"github.com/dop251/goja"
	"log"
	"time"
)

//go:embed js/standalone.js
var standalone string

//go:embed js/babel.js
var babelPlugin string

//go:embed js/typescript.js
var typescriptPlugin string

//go:embed js/estree.js
var estreePlugin string

//go:embed js/markdown.js
var markdownPlugin string

//go:embed js/html.js
var htmlPlugin string

type FormatType string

type PrettierOption struct {
	UseTabs       bool   `json:"useTabs"`
	TabWidth      int    `json:"tabWidth"`
	PrintWidth    int    `json:"printWidth"`
	SingleQuote   bool   `json:"singleQuote"`
	TrailingComma string `json:"trailingComma"`
	Semi          bool   `json:"semi"`
}

var vm *goja.Runtime
var format goja.Callable

var defaults = PrettierOption{
	UseTabs:       false,
	TabWidth:      4,
	PrintWidth:    80,
	SingleQuote:   false,
	TrailingComma: "none",
	Semi:          false,
}

func toJsonString[T any](value T) string {
	data, err := json.Marshal(value)
	if err != nil {
		log.Println("toJsonString error:", err)
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

func FormatText(ext goja.Value, code goja.Value, opts goja.Value) (goja.Value, error) {
	time.AfterFunc(2*time.Second, func() {
		vm.Interrupt("timeout > 2s")
	})

	value, err := format(goja.Undefined(), ext, code, opts)

	if err != nil {
		return code, err
	}
	return value, nil
}
func FormatTypeScript(code string, opts PrettierOption) (string, error) {
	if err := mergo.Map(&opts, defaults); err != nil {
		log.Println("mergo defaults error:", err)
		return code, err
	}

	value, err := FormatText(vm.ToValue("typescript"), vm.ToValue(code), vm.ToValue(toJsonString(opts)))
	if err != nil {
		return code, err
	}

	return handlePromise(value, code)
}
func FormatJSON(code string, opts PrettierOption) (string, error) {
	if err := mergo.Map(&opts, defaults); err != nil {
		log.Println("mergo defaults error:", err)
		return code, err
	}

	value, err := FormatText(vm.ToValue("json"), vm.ToValue(code), vm.ToValue(toJsonString(opts)))
	if err != nil {
		return code, err
	}

	return handlePromise(value, code)
}
func FormatMarkdown(code string, opts PrettierOption) (string, error) {
	if err := mergo.Map(&opts, defaults); err != nil {
		log.Println("mergo defaults error:", err)
		return code, err
	}

	value, err := FormatText(vm.ToValue("markdown"), vm.ToValue(code), vm.ToValue(toJsonString(opts)))
	if err != nil {
		return code, err
	}

	return handlePromise(value, code)
}
func FormatHTML(code string, opts PrettierOption) (string, error) {
	if err := mergo.Map(&opts, defaults); err != nil {
		log.Println("mergo defaults error:", err)
		return code, err
	}

	value, err := FormatText(vm.ToValue("html"), vm.ToValue(code), vm.ToValue(toJsonString(opts)))
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
	_, err := vm.RunString(standalone + babelPlugin + typescriptPlugin + estreePlugin + htmlPlugin + markdownPlugin + exported)
	if err != nil {
		log.Fatalln("prettier error:", err)
	}

	ret, ok := goja.AssertFunction(vm.Get("format"))
	if !ok {
		panic("not a function")
	} else {
		format = ret
	}
}
