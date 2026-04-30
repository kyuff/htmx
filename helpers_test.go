package htmx_test

import (
	"fmt"
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
