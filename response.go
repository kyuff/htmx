package htmx

import "net/http"

const (
	headerLocation                   = "HX-Location"
	headerPushURL                    = "HX-Push-Url"
	headerRedirect                   = "HX-Redirect"
	headerRefresh                    = "HX-Refresh"
	headerReplaceURL                 = "HX-Replace-Url"
	headerReswap                     = "HX-Reswap"
	headerRetarget                   = "HX-Retarget"
	headerReselect                   = "HX-Reselect"
	headerResponseTrigger            = "HX-Trigger"
	headerResponseTriggerAfterSettle = "HX-Trigger-After-Settle"
	headerResponseTriggerAfterSwap   = "HX-Trigger-After-Swap"
)

// StatusStopPolling is the HTTP status code (286) that tells htmx to stop polling.
const StatusStopPolling = 286

// Response represents the outcome of an htmx handler invocation.
// The interface is sealed: only types in this package can implement it.
type Response interface {
	// respond writes the response to the http.ResponseWriter.
	// It is unexported to seal the interface.
	respond(w http.ResponseWriter, r *http.Request)
}

// Handler handles an htmx request and returns a Response.
type Handler interface {
	Handle(r *http.Request) (Response, error)
}

// HandlerFunc is an adapter to use a plain function as a Handler.
type HandlerFunc func(r *http.Request) (Response, error)

func (f HandlerFunc) Handle(r *http.Request) (Response, error) {
	return f(r)
}

// --- Render response (created by View[T].OK) ---

// renderResponse holds template data and the originating view.
// HTMX uses the view reference for lazy template compilation.
type renderResponse struct {
	data    any
	view    viewProvider
	status  int
	headers http.Header
}

func (r *renderResponse) respond(w http.ResponseWriter, _ *http.Request) {
	writeHeaders(w, r.headers)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.status != 0 && r.status != http.StatusOK {
		w.WriteHeader(r.status)
	}
}

// --- Non-render responses ---

// headerResponse writes status + headers with no body.
type headerResponse struct {
	status  int
	headers http.Header
}

func (r *headerResponse) respond(w http.ResponseWriter, _ *http.Request) {
	writeHeaders(w, r.headers)
	w.WriteHeader(r.status)
}

// Empty creates a 204 No Content response. htmx will do nothing.
func Empty() Response {
	return &headerResponse{status: http.StatusNoContent}
}

// ClientRedirect creates a response with the HX-Redirect header,
// causing htmx to perform a client-side redirect.
func ClientRedirect(url string) Response {
	h := make(http.Header)
	h.Set(headerRedirect, url)
	return &headerResponse{status: http.StatusOK, headers: h}
}

// ClientLocation creates a response with the HX-Location header,
// causing htmx to navigate without a full page reload.
func ClientLocation(url string) Response {
	h := make(http.Header)
	h.Set(headerLocation, url)
	return &headerResponse{status: http.StatusOK, headers: h}
}

// StopPoll creates a 286 response that tells htmx to stop polling.
func StopPoll() Response {
	return &headerResponse{status: StatusStopPolling}
}

// --- With* modifiers ---

// withResponse wraps a Response to add headers or override status.
type withResponse struct {
	inner   Response
	headers http.Header
	status  int
}

func (r *withResponse) respond(w http.ResponseWriter, req *http.Request) {
	r.inner.respond(w, req)
	writeHeaders(w, r.headers)
	if r.status != 0 {
		w.WriteHeader(r.status)
	}
}

// withModifier returns a withResponse wrapping resp.
// If resp is already a renderResponse, it modifies it in place for efficiency.
func withModifier(resp Response) *withResponse {
	if wr, ok := resp.(*withResponse); ok {
		return wr
	}
	return &withResponse{inner: resp, headers: make(http.Header)}
}

// WithTrigger sets the HX-Trigger response header.
func WithTrigger(resp Response, event string) Response {
	w := withModifier(resp)
	w.headers.Set(headerResponseTrigger, event)
	return w
}

// WithTriggerAfterSettle sets the HX-Trigger-After-Settle response header.
func WithTriggerAfterSettle(resp Response, event string) Response {
	w := withModifier(resp)
	w.headers.Set(headerResponseTriggerAfterSettle, event)
	return w
}

// WithTriggerAfterSwap sets the HX-Trigger-After-Swap response header.
func WithTriggerAfterSwap(resp Response, event string) Response {
	w := withModifier(resp)
	w.headers.Set(headerResponseTriggerAfterSwap, event)
	return w
}

// WithPushURL sets the HX-Push-Url response header.
func WithPushURL(resp Response, url string) Response {
	w := withModifier(resp)
	w.headers.Set(headerPushURL, url)
	return w
}

// WithReplaceURL sets the HX-Replace-Url response header.
func WithReplaceURL(resp Response, url string) Response {
	w := withModifier(resp)
	w.headers.Set(headerReplaceURL, url)
	return w
}

// WithReswap sets the HX-Reswap response header.
func WithReswap(resp Response, swap Swap) Response {
	w := withModifier(resp)
	w.headers.Set(headerReswap, string(swap))
	return w
}

// WithRetarget sets the HX-Retarget response header.
func WithRetarget(resp Response, selector string) Response {
	w := withModifier(resp)
	w.headers.Set(headerRetarget, selector)
	return w
}

// WithReselect sets the HX-Reselect response header.
func WithReselect(resp Response, selector string) Response {
	w := withModifier(resp)
	w.headers.Set(headerReselect, selector)
	return w
}

// WithRefresh sets HX-Refresh: true, triggering a full page refresh.
func WithRefresh(resp Response) Response {
	w := withModifier(resp)
	w.headers.Set(headerRefresh, "true")
	return w
}

// WithStatus overrides the HTTP status code.
func WithStatus(resp Response, code int) Response {
	w := withModifier(resp)
	w.status = code
	return w
}

// --- helpers ---

func writeHeaders(w http.ResponseWriter, headers http.Header) {
	for key, values := range headers {
		for _, v := range values {
			w.Header().Set(key, v)
		}
	}
}

// asRender extracts the render response if resp is a render type.
func asRender(resp Response) (*renderResponse, bool) {
	switch r := resp.(type) {
	case *renderResponse:
		return r, true
	case *withResponse:
		return asRender(r.inner)
	default:
		return nil, false
	}
}
