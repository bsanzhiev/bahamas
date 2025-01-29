package server

import "net/http"

// Router wraps for http.Handler with routes registration ability
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new Router instance
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Handle registers a new route with a matcher for URL path and a handler
func (r *Router) Handle(path string, handler http.Handler) {
	r.mux.Handle(path, handler)
}

// ServeHTTP implements http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
