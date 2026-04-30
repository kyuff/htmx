package htmx

import (
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

// HTMX is the central handler for htmx-powered routes.
// It is used as a singleton passed to features that register their handlers.
type HTMX struct {
	mux      *http.ServeMux
	layout   *templateRef
	includes []templateRef

	mu        sync.RWMutex
	compiled  map[viewProvider]*template.Template
	pcompiled map[viewProvider]*template.Template
}

type templateRef struct {
	fsys fs.FS
	name string
}

// New creates a new HTMX instance.
func New() *HTMX {
	return &HTMX{
		mux:       http.NewServeMux(),
		compiled:  make(map[viewProvider]*template.Template),
		pcompiled: make(map[viewProvider]*template.Template),
	}
}

// Layout configures the base layout template used to wrap pages
// for non-htmx (full page) requests.
func (h *HTMX) Layout(fsys fs.FS, name string) error {
	if _, err := fs.ReadFile(fsys, name); err != nil {
		return fmt.Errorf("layout %q: %w", name, err)
	}
	h.layout = &templateRef{fsys: fsys, name: name}
	return nil
}

// AddInclude registers a template for inclusion in all page and partial
// rendering without an HTTP endpoint. Use for global components like navbars.
func (h *HTMX) AddInclude(fsys fs.FS, name string) error {
	if _, err := fs.ReadFile(fsys, name); err != nil {
		return fmt.Errorf("include %q: %w", name, err)
	}
	h.includes = append(h.includes, templateRef{fsys: fsys, name: name})
	return nil
}

// FileServer serves static files from the directory at path within fsys
// under the given URL prefix. It can be called multiple times to serve
// different directories.
func (h *HTMX) FileServer(fsys fs.FS, path string, urlPrefix string) error {
	sub, err := fs.Sub(fsys, path)
	if err != nil {
		return fmt.Errorf("file server %q: %w", path, err)
	}
	prefix := "/" + strings.Trim(urlPrefix, "/") + "/"
	h.mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.FS(sub))))
	return nil
}

// Page registers a handler for a page route. Pages render with the
// full layout for non-htmx requests and just the content block for
// htmx requests. Template compilation is deferred until the first request.
func (h *HTMX) Page(pattern string, handler Handler) {
	h.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		resp, err := handler.Handle(r)
		if err != nil {
			slog.Error("[htmx] handler error", "pattern", pattern, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rr, ok := asRender(resp)
		if !ok {
			resp.respond(w, r)
			return
		}

		tmpl, err := h.compilePage(rr.view)
		if err != nil {
			slog.Error("[htmx] template compile error", "pattern", pattern, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp.respond(w, r)
		if IsRequest(r) {
			tmpl.ExecuteTemplate(w, "content", rr.data)
		} else {
			tmpl.ExecuteTemplate(w, "base", rr.data)
		}
	})
}

// Partial registers a handler for a partial route. Partials always
// render as HTML fragments. Template compilation is deferred until the first request.
func (h *HTMX) Partial(pattern string, handler Handler) {
	h.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		resp, err := handler.Handle(r)
		if err != nil {
			slog.Error("[htmx] handler error", "pattern", pattern, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		rr, ok := asRender(resp)
		if !ok {
			resp.respond(w, r)
			return
		}

		tmpl, err := h.compilePartial(rr.view)
		if err != nil {
			slog.Error("[htmx] template compile error", "pattern", pattern, "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp.respond(w, r)
		tmpl.ExecuteTemplate(w, rr.view.templateName(), rr.data)
	})
}

// Handle registers the handler for the given pattern.
func (h *HTMX) Handle(pattern string, handler http.Handler) {
	h.mux.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern.
func (h *HTMX) HandleFunc(pattern string, handler http.HandlerFunc) {
	h.mux.HandleFunc(pattern, handler)
}

func (h *HTMX) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

// compilePage returns a compiled template for a page view, including
// layout and all registered includes. Results are cached by view identity.
func (h *HTMX) compilePage(vp viewProvider) (*template.Template, error) {
	h.mu.RLock()
	if tmpl, ok := h.compiled[vp]; ok {
		h.mu.RUnlock()
		return tmpl, nil
	}
	h.mu.RUnlock()

	h.mu.Lock()
	defer h.mu.Unlock()

	// Double-check after acquiring write lock.
	if tmpl, ok := h.compiled[vp]; ok {
		return tmpl, nil
	}

	tmpl := template.New("")

	if h.layout != nil {
		content, err := fs.ReadFile(h.layout.fsys, h.layout.name)
		if err != nil {
			return nil, err
		}
		if _, err = tmpl.Parse(string(content)); err != nil {
			return nil, err
		}
	}

	for _, inc := range h.includes {
		content, err := fs.ReadFile(inc.fsys, inc.name)
		if err != nil {
			return nil, err
		}
		if _, err = tmpl.Parse(string(content)); err != nil {
			return nil, err
		}
	}

	// Parse all HTML files from the view's filesystem so that
	// cross-references between templates in the same feature work.
	if err := parseAllHTML(tmpl, vp.templateFS()); err != nil {
		return nil, err
	}

	h.compiled[vp] = tmpl
	return tmpl, nil
}

// compilePartial returns a compiled template for a partial view,
// including all registered includes but not the layout.
func (h *HTMX) compilePartial(vp viewProvider) (*template.Template, error) {
	h.mu.RLock()
	if tmpl, ok := h.pcompiled[vp]; ok {
		h.mu.RUnlock()
		return tmpl, nil
	}
	h.mu.RUnlock()

	h.mu.Lock()
	defer h.mu.Unlock()

	if tmpl, ok := h.pcompiled[vp]; ok {
		return tmpl, nil
	}

	tmpl := template.New("")

	for _, inc := range h.includes {
		content, err := fs.ReadFile(inc.fsys, inc.name)
		if err != nil {
			return nil, err
		}
		if _, err = tmpl.Parse(string(content)); err != nil {
			return nil, err
		}
	}

	if err := parseAllHTML(tmpl, vp.templateFS()); err != nil {
		return nil, err
	}

	h.pcompiled[vp] = tmpl
	return tmpl, nil
}

// parseAllHTML walks the filesystem and parses all .html files into tmpl.
func parseAllHTML(tmpl *template.Template, fsys fs.FS) error {
	return fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || strings.ToLower(filepath.Ext(path)) != ".html" {
			return nil
		}
		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		_, err = tmpl.Parse(string(content))
		return err
	})
}
