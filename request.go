package htmx

import "net/http"

const (
	headerRequest               = "HX-Request"
	headerBoosted               = "HX-Boosted"
	headerCurrentURL            = "HX-Current-URL"
	headerHistoryRestoreRequest = "HX-History-Restore-Request"
	headerPrompt                = "HX-Prompt"
	headerTarget                = "HX-Target"
	headerTrigger               = "HX-Trigger"
	headerTriggerName           = "HX-Trigger-Name"
)

// IsRequest reports whether the request was made by htmx.
func IsRequest(r *http.Request) bool {
	return r.Header.Get(headerRequest) == "true"
}

// IsBoosted reports whether the request came from an hx-boost enabled element.
func IsBoosted(r *http.Request) bool {
	return r.Header.Get(headerBoosted) == "true"
}

// IsHistoryRestoreRequest reports whether the request is for history restoration
// after a miss in the local history cache.
func IsHistoryRestoreRequest(r *http.Request) bool {
	return r.Header.Get(headerHistoryRestoreRequest) == "true"
}

// Target returns the id of the target element, or empty string if not set.
func Target(r *http.Request) string {
	return r.Header.Get(headerTarget)
}

// Trigger returns the id of the triggered element, or empty string if not set.
func Trigger(r *http.Request) string {
	return r.Header.Get(headerTrigger)
}

// TriggerName returns the name of the triggered element, or empty string if not set.
func TriggerName(r *http.Request) string {
	return r.Header.Get(headerTriggerName)
}

// CurrentURL returns the current URL of the browser, or empty string if not set.
func CurrentURL(r *http.Request) string {
	return r.Header.Get(headerCurrentURL)
}

// Prompt returns the user response to an hx-prompt, or empty string if not set.
func Prompt(r *http.Request) string {
	return r.Header.Get(headerPrompt)
}
