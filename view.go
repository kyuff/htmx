package htmx

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// viewProvider is a sealed interface that allows HTMX to access
// a View's template content for compilation. Only types in this
// package can implement it.
type viewProvider interface {
	templateContent() ([]byte, error)
	templateName() string
	templateFS() fs.FS
}

// View binds a template file to a typed view model.
// It is created at the feature level and used to produce
// type-safe render responses.
type View[T any] struct {
	fsys fs.FS
	path string
	name string
}

// NewView creates a View that binds a template file to a model type.
// The fsys and path identify the template file. Template compilation
// is deferred until the first request.
func NewView[T any](fsys fs.FS, path string) *View[T] {
	return &View[T]{
		fsys: fsys,
		path: path,
		name: viewName(path),
	}
}

// OK creates a render response carrying the given data.
// The returned Response is used by HTMX to lazily compile and
// execute the View's template.
func (v *View[T]) OK(data T) Response {
	return &renderResponse{
		data:   data,
		view:   v,
		status: http.StatusOK,
	}
}

func (v *View[T]) templateContent() ([]byte, error) {
	return fs.ReadFile(v.fsys, v.path)
}

func (v *View[T]) templateName() string {
	return v.name
}

func (v *View[T]) templateFS() fs.FS {
	return v.fsys
}

// viewName derives the Go template name from a file path.
func viewName(name string) string {
	return strings.TrimSuffix(path.Base(name), ".html")
}
