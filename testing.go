package htmx

import (
	"bytes"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// RenderTest compiles the view's template (standalone, without layout)
// and renders it with the given data. It returns the rendered HTML string.
// The test fails immediately if the template cannot be parsed or executed.
func RenderTest[T any](t testing.TB, v *View[T], data T) string {
	t.Helper()

	tmpl := template.New("")

	// Parse all .html files from the view's filesystem for cross-references.
	err := fs.WalkDir(v.fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.ToLower(filepath.Ext(path)) != ".html" {
			return nil
		}
		content, err := fs.ReadFile(v.fsys, path)
		if err != nil {
			return err
		}
		_, err = tmpl.Parse(string(content))
		return err
	})
	if err != nil {
		t.Fatalf("RenderTest: parse templates: %v", err)
	}

	// Try the view's derived name first, fall back to "content"
	// (pages define "content" while partials use the file-based name).
	name := v.name
	if tmpl.Lookup(name) == nil {
		name = "content"
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		t.Fatalf("RenderTest: execute template %q: %v", name, err)
	}
	return buf.String()
}

// AssertData extracts the typed view model data from a Response and passes it
// to the assertion function. The test fails immediately if the response is nil,
// not a render response, or the data does not match the expected type.
func AssertData[T any](t *testing.T, resp Response, fn func(t *testing.T, got T)) {
	t.Helper()
	if resp == nil {
		t.Fatal("AssertData: response is nil")
		return
	}
	rr, ok := asRender(resp)
	if !ok {
		t.Fatalf("AssertData: expected render response, got %T", resp)
		return
	}
	data, ok := rr.data.(T)
	if !ok {
		t.Fatalf("AssertData: expected data type %T, got %T", *new(T), rr.data)
		return
	}
	fn(t, data)
}
