package htmx

import "net/http"

// ExportRespond exposes the sealed respond method for testing.
var ExportRespond = func(resp Response, w http.ResponseWriter, r *http.Request) {
	resp.respond(w, r)
}
