package prettier_go

import (
	"testing"
)

func TestPrettierFormat(t *testing.T) {
	// General
	if ret, _ := FormatTypeScript("if (1) {\n1}", PrettierFormatOptions{TabWidth: 4}); ret != "if (1) {\n    1\n}\n" {
		t.Error("FormatTypeScript error", ret)
	}
	if ret, _ := FormatTypeScript("if (1) {\n1}", PrettierFormatOptions{UseTabs: true}); ret != "if (1) {\n\t1\n}\n" {
		t.Error("FormatTypeScript error", ret)
	}
	if ret, _ := FormatTypeScript("if ('1') {\n1}", PrettierFormatOptions{SingleQuote: false}); ret != "if (\"1\") {\n    1\n}\n" {
		t.Error("FormatTypeScript error", ret)
	}
	if ret, _ := FormatTypeScript("1", PrettierFormatOptions{Semi: true}); ret != "1;\n" {
		t.Error("FormatTypeScript error", ret)
	}

	// js
	if ret, _ := FormatTypeScript("var a=1", PrettierFormatOptions{}); ret != "var a = 1\n" {
		t.Error("FormatTypeScript error", ret)
	}
	// ts
	if ret, _ := FormatTypeScript("var a: number = 1", PrettierFormatOptions{}); ret != "var a: number = 1\n" {
		t.Error("FormatTypeScript error", ret)
	}
	// jsx
	if ret, _ := FormatTypeScript("var a=<b>1</b>", PrettierFormatOptions{}); ret != "var a = <b>1</b>\n" {
		t.Error("FormatTypeScript error", ret)
	}

	// JSON
	if ret, err := FormatJSON("{a:1}", PrettierFormatOptions{}); ret != "{ \"a\": 1 }\n" {
		t.Error("FormatJSON error", ret, err)
	}
	// Markdown
	if ret, err := FormatMarkdown("#  1", PrettierFormatOptions{}); ret != "# 1\n" {
		t.Error("FormatMarkdown error", ret, err)
	}
}
