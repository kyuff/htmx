package htmx_test

import (
	"fmt"
	"io/fs"
	"net/http"
	"testing/fstest"

	"github.com/kyuff/htmx"
)

func staticHandler(resp htmx.Response) htmx.HandlerFunc {
	return func(_ *http.Request) (htmx.Response, error) {
		return resp, nil
	}
}

func errorHandler(msg string) htmx.HandlerFunc {
	return func(_ *http.Request) (htmx.Response, error) {
		return nil, fmt.Errorf("%s", msg)
	}
}

func testView[T any](name string, tmplBody string) *htmx.View[T] {
	fsys := fstest.MapFS{
		name: &fstest.MapFile{Data: []byte(tmplBody)},
	}
	return htmx.NewView[T](fsys, name)
}

// failingFS wraps an fs.FS and returns an error when opening a specific file.
type failingFS struct {
	inner    fs.FS
	failName string
}

func (f failingFS) Open(name string) (fs.File, error) {
	if name == f.failName {
		return nil, fmt.Errorf("injected read error for %q", name)
	}
	return f.inner.Open(name)
}

// succeedOnceThenFailFS succeeds on the first Open of failName, then fails.
// Used to pass setup-time validation but fail on the first request.
type succeedOnceThenFailFS struct {
	inner    fs.FS
	failName string
	opened   bool
}

func (f *succeedOnceThenFailFS) Open(name string) (fs.File, error) {
	if name == f.failName {
		if !f.opened {
			f.opened = true
			return f.inner.Open(name)
		}
		return nil, fmt.Errorf("injected read error for %q", name)
	}
	return f.inner.Open(name)
}
