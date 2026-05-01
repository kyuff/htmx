package htmx_test

import (
	"sync"
	"testing"

	"github.com/kyuff/htmx/internal/assert"
	"github.com/kyuff/htmx"
)

func TestRenderTest(t *testing.T) {
	t.Run("render page template content", func(t *testing.T) {
		// arrange
		type vm struct{ Name string }
		var (
			view = testView[vm]("greeting.html", `{{ define "content" }}<h1>Hello {{ .Name }}</h1>{{ end }}`)
		)

		// act
		got := htmx.RenderTest(t, view, vm{Name: "World"})

		// assert
		assert.Contains(t, got, "Hello World")
	})

	t.Run("render partial template by file name", func(t *testing.T) {
		// arrange
		type vm struct{ Count int }
		var (
			view = testView[vm]("counter.html", `{{ define "counter" }}<span>{{ .Count }}</span>{{ end }}`)
		)

		// act
		got := htmx.RenderTest(t, view, vm{Count: 42})

		// assert
		assert.Contains(t, got, "42")
	})

	t.Run("fail when template has invalid syntax", func(t *testing.T) {
		// arrange
		var (
			x    = &testing.T{}
			wg   sync.WaitGroup
			view = testView[any]("bad.html", `{{ define "content" }}{{ .`)
		)

		// act
		wg.Add(1)
		go func() {
			defer wg.Done()
			htmx.RenderTest(x, view, nil)
		}()
		wg.Wait()

		// assert
		assert.Equal(t, true, x.Failed())
	})

	t.Run("fail when template name cannot be resolved", func(t *testing.T) {
		// arrange
		var (
			x    = &testing.T{}
			wg   sync.WaitGroup
			// Template defines "other" — neither the file name "unknown" nor "content".
			view = testView[any]("unknown.html", `{{ define "other" }}ok{{ end }}`)
		)

		// act
		wg.Add(1)
		go func() {
			defer wg.Done()
			htmx.RenderTest(x, view, nil)
		}()
		wg.Wait()

		// assert
		assert.Equal(t, true, x.Failed())
	})
}

func TestAssertData(t *testing.T) {
	t.Run("extract typed data from render response", func(t *testing.T) {
		// arrange
		type vm struct{ Value string }
		var (
			view = testView[vm]("item.html", `{{ define "item" }}{{ .Value }}{{ end }}`)
			resp = view.OK(vm{Value: "hello"})
		)

		// act + assert
		htmx.AssertData(t, resp, func(t *testing.T, got vm) {
			assert.Equal(t, "hello", got.Value)
		})
	})

	t.Run("fail when response is nil", func(t *testing.T) {
		// arrange
		var (
			x  = &testing.T{}
			wg sync.WaitGroup
		)

		// act
		wg.Add(1)
		go func() {
			defer wg.Done()
			htmx.AssertData[any](x, nil, func(t *testing.T, got any) {})
		}()
		wg.Wait()

		// assert
		assert.Equal(t, true, x.Failed())
	})

	t.Run("fail when response is not a render response", func(t *testing.T) {
		// arrange
		var (
			x  = &testing.T{}
			wg sync.WaitGroup
		)

		// act
		wg.Add(1)
		go func() {
			defer wg.Done()
			htmx.AssertData[any](x, htmx.Empty(), func(t *testing.T, got any) {})
		}()
		wg.Wait()

		// assert
		assert.Equal(t, true, x.Failed())
	})

	t.Run("fail when data type does not match", func(t *testing.T) {
		// arrange
		type vm struct{ Value string }
		var (
			x    = &testing.T{}
			wg   sync.WaitGroup
			view = testView[vm]("item.html", `{{ define "item" }}{{ .Value }}{{ end }}`)
			resp = view.OK(vm{Value: "hello"})
		)

		// act
		wg.Add(1)
		go func() {
			defer wg.Done()
			htmx.AssertData[int](x, resp, func(t *testing.T, got int) {})
		}()
		wg.Wait()

		// assert
		assert.Equal(t, true, x.Failed())
	})
}
