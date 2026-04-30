package htmx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyuff/htmx/internal/assert"
	"github.com/kyuff/htmx"
)

func TestIsRequest(t *testing.T) {
	t.Run("return false when header is absent", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)

		// act
		got := htmx.IsRequest(r)

		// assert
		assert.False(t, got)
	})

	t.Run("return true when HX-Request is true", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)
		r.Header.Set("HX-Request", "true")

		// act
		got := htmx.IsRequest(r)

		// assert
		assert.Truef(t, got, "expected IsRequest to return true")
	})
}

func TestIsBoosted(t *testing.T) {
	t.Run("return false when header is absent", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)

		// act
		got := htmx.IsBoosted(r)

		// assert
		assert.False(t, got)
	})

	t.Run("return true when HX-Boosted is true", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)
		r.Header.Set("HX-Boosted", "true")

		// act
		got := htmx.IsBoosted(r)

		// assert
		assert.Truef(t, got, "expected IsBoosted to return true")
	})
}

func TestIsHistoryRestoreRequest(t *testing.T) {
	t.Run("return false when header is absent", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)

		// act
		got := htmx.IsHistoryRestoreRequest(r)

		// assert
		assert.False(t, got)
	})

	t.Run("return true when HX-History-Restore-Request is true", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)
		r.Header.Set("HX-History-Restore-Request", "true")

		// act
		got := htmx.IsHistoryRestoreRequest(r)

		// assert
		assert.Truef(t, got, "expected IsHistoryRestoreRequest to return true")
	})
}

func TestTarget(t *testing.T) {
	t.Run("return empty string when header is absent", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)

		// act
		got := htmx.Target(r)

		// assert
		assert.Equal(t, "", got)
	})

	t.Run("return the target element id", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)
		r.Header.Set("HX-Target", "result-div")

		// act
		got := htmx.Target(r)

		// assert
		assert.Equal(t, "result-div", got)
	})
}

func TestTrigger(t *testing.T) {
	t.Run("return empty string when header is absent", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)

		// act
		got := htmx.Trigger(r)

		// assert
		assert.Equal(t, "", got)
	})

	t.Run("return the trigger element id", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)
		r.Header.Set("HX-Trigger", "my-button")

		// act
		got := htmx.Trigger(r)

		// assert
		assert.Equal(t, "my-button", got)
	})
}

func TestTriggerName(t *testing.T) {
	t.Run("return empty string when header is absent", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)

		// act
		got := htmx.TriggerName(r)

		// assert
		assert.Equal(t, "", got)
	})

	t.Run("return the trigger element name", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)
		r.Header.Set("HX-Trigger-Name", "search")

		// act
		got := htmx.TriggerName(r)

		// assert
		assert.Equal(t, "search", got)
	})
}

func TestCurrentURL(t *testing.T) {
	t.Run("return empty string when header is absent", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)

		// act
		got := htmx.CurrentURL(r)

		// assert
		assert.Equal(t, "", got)
	})

	t.Run("return the current browser URL", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)
		r.Header.Set("HX-Current-URL", "https://example.com/page")

		// act
		got := htmx.CurrentURL(r)

		// assert
		assert.Equal(t, "https://example.com/page", got)
	})
}

func TestPrompt(t *testing.T) {
	t.Run("return empty string when header is absent", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)

		// act
		got := htmx.Prompt(r)

		// assert
		assert.Equal(t, "", got)
	})

	t.Run("return the user prompt response", func(t *testing.T) {
		// arrange
		var (
			r = httptest.NewRequest(http.MethodGet, "/", nil)
		)
		r.Header.Set("HX-Prompt", "confirmed")

		// act
		got := htmx.Prompt(r)

		// assert
		assert.Equal(t, "confirmed", got)
	})
}
