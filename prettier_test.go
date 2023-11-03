package prettier_go

import (
	_ "embed"
	"testing"
)

//go:embed js/antd.min.js
var antd string

func TestPrettierFormat(t *testing.T) {
	// General
	if ret, _ := FormatTypeScript("if (1) {\n1}", PrettierOption{TabWidth: 4}); ret != "if (1) {\n    1\n}\n" {
		t.Error("FormatTypeScript error", ret)
	}
	if ret, _ := FormatTypeScript("if (1) {\n1}", PrettierOption{UseTabs: true}); ret != "if (1) {\n\t1\n}\n" {
		t.Error("FormatTypeScript error", ret)
	}
	if ret, _ := FormatTypeScript("if ('1') {\n1}", PrettierOption{SingleQuote: false}); ret != "if (\"1\") {\n    1\n}\n" {
		t.Error("FormatTypeScript error", ret)
	}
	if ret, _ := FormatTypeScript("1", PrettierOption{Semi: true}); ret != "1;\n" {
		t.Error("FormatTypeScript error", ret)
	}

	// js
	if ret, _ := FormatTypeScript("var a=1", PrettierOption{}); ret != "var a = 1\n" {
		t.Error("FormatTypeScript error", ret)
	}
	// ts
	if ret, _ := FormatTypeScript("var a: number = 1", PrettierOption{}); ret != "var a: number = 1\n" {
		t.Error("FormatTypeScript error", ret)
	}
	// jsx
	if ret, _ := FormatTypeScript("var a=<b>1</b>", PrettierOption{}); ret != "var a = <b>1</b>\n" {
		t.Error("FormatTypeScript error", ret)
	}

	// JSON
	if ret, err := FormatJSON("{a:1}", PrettierOption{}); ret != "{ \"a\": 1 }\n" {
		t.Error("FormatJSON error", ret, err)
	}
	// Markdown
	if ret, err := FormatMarkdown("#  1", PrettierOption{}); ret != "# 1\n" {
		t.Error("FormatMarkdown error", ret, err)
	}
	// HTML
	if ret, err := FormatHTML("<body  class='a' ></body>", PrettierOption{}); ret != "<body class=\"a\"></body>\n" {
		t.Error("FormatHTML error", ret, err)
	}
}

func TestPrettierFormat1(t *testing.T) {
	script, err := FormatTypeScript(antd, PrettierOption{})
	if err != nil {
		t.Log(err)
	} else {
		t.Log(len(script))
	}

}
