# htmx

A Go package for building htmx-powered web applications. It provides an HTTP handler that manages routing, template compilation, and rendering — distinguishing between full-page requests and htmx fragment requests automatically.

## Installation

```sh
go get github.com/kyuff/htmx
```

Requires Go 1.22 or later. No external dependencies.

## Usage

### Setup

```go
h := htmx.New()

// Set the base layout template (wraps pages for full-page requests)
h.Layout(webFS, "layouts/base.html")

// Register global components (e.g. navbar) included in all renders
h.AddInclude(webFS, "partials/navbar.html")

// Serve static files
h.FileServer(webFS, "assets", "assets")

http.ListenAndServe(":8080", h)
```

### Pages and Partials

```go
// Define a typed view bound to a template file
type CounterVM struct{ Count int }
var CounterPage = htmx.NewView[CounterVM](webFS, "counter.html")

// Register a page — renders full layout for normal requests,
// content fragment only for htmx requests
h.Page("GET /counter", htmx.HandlerFunc(func(r *http.Request) (htmx.Response, error) {
    return CounterPage.OK(CounterVM{Count: 0}), nil
}))

// Register a partial — always renders as an HTML fragment
h.Partial("POST /counter/increment", htmx.HandlerFunc(func(r *http.Request) (htmx.Response, error) {
    return CounterValue.OK(CounterVM{Count: 1}), nil
}))
```

### Response helpers

```go
// Redirect, location, stop polling, empty response
htmx.ClientRedirect("/login")
htmx.ClientLocation("/new-page")
htmx.StopPoll()
htmx.Empty()

// Response modifiers (chainable)
htmx.WithTrigger(resp, "counterUpdated")
htmx.WithPushURL(resp, "/counter")
htmx.WithReswap(resp, htmx.SwapOuterHTML)
htmx.WithRefresh(resp)
```

### Request helpers

```go
htmx.IsRequest(r)              // HX-Request: true
htmx.IsBoosted(r)              // HX-Boosted: true
htmx.Target(r)                 // HX-Target
htmx.Trigger(r)                // HX-Trigger
htmx.TriggerName(r)            // HX-Trigger-Name
htmx.CurrentURL(r)             // HX-Current-URL
htmx.Prompt(r)                 // HX-Prompt
```

### Testing

```go
// Render a view standalone (no layout) and assert on the HTML
html := htmx.RenderTest(t, CounterPage, CounterVM{Count: 42})

// Extract typed data from a Response in handler tests
htmx.AssertData(t, resp, func(t *testing.T, got CounterVM) {
    assert.Equal(t, 42, got.Count)
})
```

## License

MIT
