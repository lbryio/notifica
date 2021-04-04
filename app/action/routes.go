package action

import (
	"net/http"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/orderedmap"
	"github.com/lbryio/notifica/app/config"
)

// GetRoutes returns the routes for the API Server
func GetRoutes() *Routes {
	routes := Routes{}

	routes.Set("/status", Status)
	routes.Set("/", Root)

	if config.IsDebugMode {
		routes.Set("/test", Test)
	}

	return &routes
}

// Routes container for the route map between path and handler
type Routes struct {
	m *orderedmap.Map
}

// Set sets the map entry for the route
func (r *Routes) Set(key string, h api.Handler) {
	if r.m == nil {
		r.m = orderedmap.New()
	}

	r.m.Set(key, h)
}

// Each applies a function wrapper middleware to each route
func (r *Routes) Each(f func(string, http.Handler)) {
	if r.m == nil {
		return
	}
	for _, k := range r.m.Keys() {
		a, _ := r.m.Get(k)
		f(k, a.(http.Handler))
	}
}

// Walk applies a function to each route
func (r *Routes) Walk(f func(string, http.Handler) http.Handler) {
	if r.m == nil {
		return
	}
	for _, k := range r.m.Keys() {
		a, _ := r.m.Get(k)
		r.m.Set(k, f(k, a.(http.Handler)))
	}
}
