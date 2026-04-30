package htmx

// Swap defines how htmx swaps HTML content into the DOM.
type Swap string

const (
	// SwapInnerHTML replaces the inner html of the target element.
	SwapInnerHTML Swap = "innerHTML"
	// SwapOuterHTML replaces the entire target element with the response.
	SwapOuterHTML Swap = "outerHTML"
	// SwapAfterBegin prepends the content before the first child of the target.
	SwapAfterBegin Swap = "afterbegin"
	// SwapBeforeBegin inserts the content before the target element.
	SwapBeforeBegin Swap = "beforebegin"
	// SwapBeforeEnd appends the content after the last child of the target.
	SwapBeforeEnd Swap = "beforeend"
	// SwapAfterEnd inserts the content after the target element.
	SwapAfterEnd Swap = "afterend"
	// SwapDelete deletes the target element regardless of the response.
	SwapDelete Swap = "delete"
	// SwapNone does not append content from the response.
	SwapNone Swap = "none"
)
