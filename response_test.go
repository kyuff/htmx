package htmx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyuff/htmx/internal/assert"
	"github.com/kyuff/htmx"
)

func TestResponse(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		t.Run("write 204 no content", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.Empty()

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, http.StatusNoContent, w.Code)
		})
	})

	t.Run("ClientRedirect", func(t *testing.T) {
		t.Run("set HX-Redirect header", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.ClientRedirect("/login")

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "/login", w.Header().Get("HX-Redirect"))
		})
	})

	t.Run("ClientLocation", func(t *testing.T) {
		t.Run("set HX-Location header", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.ClientLocation("/new-page")

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "/new-page", w.Header().Get("HX-Location"))
		})
	})

	t.Run("StopPoll", func(t *testing.T) {
		t.Run("write 286 stop polling", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.StopPoll()

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, htmx.StatusStopPolling, w.Code)
		})
	})

	t.Run("WithTrigger", func(t *testing.T) {
		t.Run("set HX-Trigger header", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithTrigger(htmx.Empty(), "showMessage")

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "showMessage", w.Header().Get("HX-Trigger"))
		})
	})

	t.Run("WithTriggerAfterSettle", func(t *testing.T) {
		t.Run("set HX-Trigger-After-Settle header", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithTriggerAfterSettle(htmx.Empty(), "settled")

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "settled", w.Header().Get("HX-Trigger-After-Settle"))
		})
	})

	t.Run("WithTriggerAfterSwap", func(t *testing.T) {
		t.Run("set HX-Trigger-After-Swap header", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithTriggerAfterSwap(htmx.Empty(), "swapped")

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "swapped", w.Header().Get("HX-Trigger-After-Swap"))
		})
	})

	t.Run("WithPushURL", func(t *testing.T) {
		t.Run("set HX-Push-Url header", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithPushURL(htmx.Empty(), "/pushed")

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "/pushed", w.Header().Get("HX-Push-Url"))
		})
	})

	t.Run("WithReplaceURL", func(t *testing.T) {
		t.Run("set HX-Replace-Url header", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithReplaceURL(htmx.Empty(), "/replaced")

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "/replaced", w.Header().Get("HX-Replace-Url"))
		})
	})

	t.Run("WithReswap", func(t *testing.T) {
		t.Run("set HX-Reswap header", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithReswap(htmx.Empty(), htmx.SwapOuterHTML)

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "outerHTML", w.Header().Get("HX-Reswap"))
		})
	})

	t.Run("WithRetarget", func(t *testing.T) {
		t.Run("set HX-Retarget header", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithRetarget(htmx.Empty(), "#other-div")

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "#other-div", w.Header().Get("HX-Retarget"))
		})
	})

	t.Run("WithReselect", func(t *testing.T) {
		t.Run("set HX-Reselect header", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithReselect(htmx.Empty(), ".content")

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, ".content", w.Header().Get("HX-Reselect"))
		})
	})

	t.Run("WithRefresh", func(t *testing.T) {
		t.Run("set HX-Refresh header to true", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithRefresh(htmx.Empty())

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "true", w.Header().Get("HX-Refresh"))
		})
	})

	t.Run("WithStatus", func(t *testing.T) {
		t.Run("override HTTP status code on a render response", func(t *testing.T) {
			// arrange
			var (
				view = testView[any]("msg.html", `{{ define "msg" }}ok{{ end }}`)
				w    = httptest.NewRecorder()
				r    = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithStatus(view.OK(nil), http.StatusTeapot)

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, http.StatusTeapot, w.Code)
		})
	})

	t.Run("WithTrigger chained on WithTrigger", func(t *testing.T) {
		t.Run("reuse existing withResponse wrapper", func(t *testing.T) {
			// arrange
			var (
				w = httptest.NewRecorder()
				r = httptest.NewRequest(http.MethodGet, "/", nil)
			)

			// act
			resp := htmx.WithTrigger(htmx.WithTriggerAfterSettle(htmx.Empty(), "settled"), "fired")

			// assert
			htmx.ExportRespond(resp, w, r)
			assert.Equal(t, "fired", w.Header().Get("HX-Trigger"))
			assert.Equal(t, "settled", w.Header().Get("HX-Trigger-After-Settle"))
		})
	})
}
