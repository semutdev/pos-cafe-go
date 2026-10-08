package main

import (
	"bytes"
	"testing"

	"github.com/semutdev/goigniter/system/core"
	"github.com/semutdev/goigniter/system/helpers"
)

func TestAllTranspiledViewsParse(t *testing.T) {
	helpers.Init("http://localhost:8080")

	// Custom template funcs to cover CI3 helpers if any
	funcs := helpers.AllTemplateFuncs()
	funcs["currency"] = func(val any) string { return "" }
	funcs["format_date"] = func(val any) string { return "" }
	funcs["date"] = func(format string, val any) string { return "" }

	engine, err := core.NewTemplateEngine(core.TemplateConfig{
		Dir:     "./app/views",
		Ext:     ".html",
		FuncMap: funcs,
		Reload:  false,
	})
	if err != nil {
		t.Fatalf("Template engine failed to load: %v", err)
	}
	viewsToTest := []string{
		"admin/dashboard/index",
		"pos/index",
		"admin/user/index",
		"welcome_message",
		"inc/header",
		"inc/footer",
	}

	for _, viewName := range viewsToTest {
		var buf bytes.Buffer
		if err := engine.Render(&buf, viewName, core.Map{
			"Title": "Test",
			"Users": []any{},
			"ListKategori": []any{},
		}); err != nil {
			t.Fatalf("Render %s failed: %v", viewName, err)
		}
		t.Logf("Rendered %s: %d bytes", viewName, buf.Len())
	}
}
