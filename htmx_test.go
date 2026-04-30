package htmx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/kyuff/htmx/internal/assert"
	"github.com/kyuff/htmx"
)

func TestHTMX(t *testing.T) {
	t.Run("ServeHTTP", func(t *testing.T) {
		t.Run("delegate to registered handler", func(t *testing.T) {
			// arrange
			var (
				sut  = htmx.New()
				w    = httptest.NewRecorder()
				r    = httptest.NewRequest(http.MethodGet, "/test", nil)
				body = "hello from handler"
			)
			sut.HandleFunc("GET /test", func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte(body))
			})

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, body, w.Body.String())
		})

		t.Run("return 404 for unregistered route", func(t *testing.T) {
			// arrange
			var (
				sut = htmx.New()
				w   = httptest.NewRecorder()
				r   = httptest.NewRequest(http.MethodGet, "/missing", nil)
			)

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	})

	t.Run("Handle", func(t *testing.T) {
		t.Run("register an http.Handler", func(t *testing.T) {
			// arrange
			var (
				sut  = htmx.New()
				w    = httptest.NewRecorder()
				r    = httptest.NewRequest(http.MethodPost, "/submit", nil)
				body = "handled"
			)
			sut.Handle("POST /submit", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte(body))
			}))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, body, w.Body.String())
		})
	})

	t.Run("FileServer", func(t *testing.T) {
		t.Run("serve static files under url prefix", func(t *testing.T) {
			// arrange
			var (
				fsys = fstest.MapFS{
					"static/css/app.css": &fstest.MapFile{Data: []byte("body{}")},
				}
				sut = htmx.New()
				w   = httptest.NewRecorder()
				r   = httptest.NewRequest(http.MethodGet, "/assets/css/app.css", nil)
			)
			assert.NoError(t, sut.FileServer(fsys, "static", "assets"))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), "body{}")
		})
	})

	t.Run("Page", func(t *testing.T) {
		t.Run("render full page for normal request", func(t *testing.T) {
			// arrange
			type vm struct{ Title string }
			var (
				layoutFS = fstest.MapFS{
					"base.html": &fstest.MapFile{
						Data: []byte(`{{ define "base" }}<html><body>{{ block "content" . }}{{ end }}</body></html>{{ end }}`),
					},
				}
				view = testView[vm]("home.html", `{{ define "content" }}<h1>{{ .Title }}</h1>{{ end }}`)
				sut  = htmx.New()
				w    = httptest.NewRecorder()
				r    = httptest.NewRequest(http.MethodGet, "/", nil)
			)
			assert.NoError(t, sut.Layout(layoutFS, "base.html"))
			sut.Page("GET /", staticHandler(view.OK(vm{Title: "Hello"})))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), "<html>")
			assert.Contains(t, w.Body.String(), "<h1>Hello</h1>")
		})

		t.Run("render content fragment for htmx request", func(t *testing.T) {
			// arrange
			type vm struct{ Title string }
			var (
				layoutFS = fstest.MapFS{
					"base.html": &fstest.MapFile{
						Data: []byte(`{{ define "base" }}<html><body>{{ block "content" . }}{{ end }}</body></html>{{ end }}`),
					},
				}
				view = testView[vm]("home.html", `{{ define "content" }}<h1>{{ .Title }}</h1>{{ end }}`)
				sut  = htmx.New()
				w    = httptest.NewRecorder()
				r    = httptest.NewRequest(http.MethodGet, "/", nil)
			)
			assert.NoError(t, sut.Layout(layoutFS, "base.html"))
			r.Header.Set("HX-Request", "true")
			sut.Page("GET /", staticHandler(view.OK(vm{Title: "Fragment"})))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "<h1>Fragment</h1>", w.Body.String())
		})

		t.Run("set content-type to text/html", func(t *testing.T) {
			// arrange
			var (
				layoutFS = fstest.MapFS{
					"base.html": &fstest.MapFile{
						Data: []byte(`{{ define "base" }}<html>{{ block "content" . }}{{ end }}</html>{{ end }}`),
					},
				}
				view = testView[any]("ct.html", `{{ define "content" }}ok{{ end }}`)
				sut  = htmx.New()
				w    = httptest.NewRecorder()
				r    = httptest.NewRequest(http.MethodGet, "/ct", nil)
			)
			assert.NoError(t, sut.Layout(layoutFS, "base.html"))
			sut.Page("GET /ct", staticHandler(view.OK(nil)))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
		})

		t.Run("return 500 when handler returns error", func(t *testing.T) {
			// arrange
			var (
				sut = htmx.New()
				w   = httptest.NewRecorder()
				r   = httptest.NewRequest(http.MethodGet, "/fail", nil)
			)
			sut.Page("GET /fail", errorHandler("loader error"))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("handle non-render response from page handler", func(t *testing.T) {
			// arrange
			var (
				sut = htmx.New()
				w   = httptest.NewRecorder()
				r   = httptest.NewRequest(http.MethodGet, "/redirect", nil)
			)
			sut.Page("GET /redirect", staticHandler(htmx.ClientRedirect("/login")))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "/login", w.Header().Get("HX-Redirect"))
		})

		t.Run("include global includes in page rendering", func(t *testing.T) {
			// arrange
			var (
				layoutFS = fstest.MapFS{
					"base.html": &fstest.MapFile{
						Data: []byte(`{{ define "base" }}<html>{{ template "nav" . }}{{ block "content" . }}{{ end }}</html>{{ end }}`),
					},
					"nav.html": &fstest.MapFile{
						Data: []byte(`{{ define "nav" }}<nav>Menu</nav>{{ end }}`),
					},
				}
				view = testView[any]("home.html", `{{ define "content" }}<main>Page</main>{{ end }}`)
				sut  = htmx.New()
				w    = httptest.NewRecorder()
				r    = httptest.NewRequest(http.MethodGet, "/", nil)
			)
			assert.NoError(t, sut.Layout(layoutFS, "base.html"))
			assert.NoError(t, sut.AddInclude(layoutFS, "nav.html"))
			sut.Page("GET /", staticHandler(view.OK(nil)))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), "<nav>Menu</nav>")
			assert.Contains(t, w.Body.String(), "<main>Page</main>")
		})
	})

	t.Run("Partial", func(t *testing.T) {
		t.Run("render partial fragment", func(t *testing.T) {
			// arrange
			type vm struct{ Name string }
			var (
				view = testView[vm]("greeting.html", `{{ define "greeting" }}<span>Hello {{ .Name }}</span>{{ end }}`)
				sut  = htmx.New()
				w    = httptest.NewRecorder()
				r    = httptest.NewRequest(http.MethodGet, "/partials/greeting", nil)
			)
			sut.Partial("GET /partials/greeting", staticHandler(view.OK(vm{Name: "World"})))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "<span>Hello World</span>", w.Body.String())
		})

		t.Run("set content-type to text/html", func(t *testing.T) {
			// arrange
			var (
				view = testView[any]("ct.html", `{{ define "ct" }}ok{{ end }}`)
				sut  = htmx.New()
				w    = httptest.NewRecorder()
				r    = httptest.NewRequest(http.MethodGet, "/partials/ct", nil)
			)
			sut.Partial("GET /partials/ct", staticHandler(view.OK(nil)))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
		})

		t.Run("return 500 when handler returns error", func(t *testing.T) {
			// arrange
			var (
				sut = htmx.New()
				w   = httptest.NewRecorder()
				r   = httptest.NewRequest(http.MethodGet, "/partials/fail", nil)
			)
			sut.Partial("GET /partials/fail", errorHandler("loader error"))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})

		t.Run("handle non-render response from partial handler", func(t *testing.T) {
			// arrange
			var (
				sut = htmx.New()
				w   = httptest.NewRecorder()
				r   = httptest.NewRequest(http.MethodGet, "/partials/empty", nil)
			)
			sut.Partial("GET /partials/empty", staticHandler(htmx.Empty()))

			// act
			sut.ServeHTTP(w, r)

			// assert
			assert.Equal(t, http.StatusNoContent, w.Code)
		})
	})
}
