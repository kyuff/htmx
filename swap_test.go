package htmx_test

import (
	"testing"

	"github.com/kyuff/htmx/internal/assert"
	"github.com/kyuff/htmx"
)

func TestSwap(t *testing.T) {
	t.Run("have expected string values", func(t *testing.T) {
		// arrange
		var (
			cases = []struct {
				swap htmx.Swap
				want string
			}{
				{htmx.SwapInnerHTML, "innerHTML"},
				{htmx.SwapOuterHTML, "outerHTML"},
				{htmx.SwapAfterBegin, "afterbegin"},
				{htmx.SwapBeforeBegin, "beforebegin"},
				{htmx.SwapBeforeEnd, "beforeend"},
				{htmx.SwapAfterEnd, "afterend"},
				{htmx.SwapDelete, "delete"},
				{htmx.SwapNone, "none"},
			}
		)

		for _, tc := range cases {
			t.Run(tc.want, func(t *testing.T) {
				// act
				got := string(tc.swap)

				// assert
				assert.Equal(t, tc.want, got)
			})
		}
	})
}
