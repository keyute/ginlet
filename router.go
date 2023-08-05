package gimlet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RouterGroup represents a group of routes that share a common base path and middleware functions.
type RouterGroup struct {
	// BasePath is the common prefix for all routes in this group
	BasePath string
	// Middlewares are executed for each route in this group before the route's own middlewares.
	Middlewares []gin.HandlerFunc
	// PersistentMiddlewares are executed for each route before the route's own middlewares and Middlewares.
	PersistentMiddlewares []gin.HandlerFunc
	// Routes is a map where the key is the HTTP method and the value is a slice of Route.
	Routes map[string][]Route
	// PreFunc is a function that is executed before Apply is called
	PreFunc func() error
	// PostFunc is a function that is executed after Apply is called
	PostFunc func() error
	// SubGroups are nested router groups. Each subgroup inherits BasePath and PersistentMiddlewares from its parent.
	SubGroups []RouterGroup
}

// RestRouterGroup is a specialized version of RouterGroup for creating RESTful API endpoints.
// If RouterGroup.Routes is not empty, it will be ignored.
type RestRouterGroup struct {
	//GetRoute is the route for GET requests
	GetRoute Route
	//PostRoute is the route for POST requests
	PostRoute Route
	//PatchRoute is the route for PATCH requests
	PatchRoute Route
	//PutRoute is the route for PUT requests
	PutRoute Route
	//DeleteRoute is the route for DELETE requests
	DeleteRoute Route
	// RouterGroup is the base group of routes for the RESTful API.
	RouterGroup
}

// Route represents a single route
type Route struct {
	// Path is the path for this route
	Path string
	// Handler is the function to execute for this route.
	Handler func(c *gin.Context)
	// Middlewares are functions executed for this route after RouterGroup.PersistentMiddlewares and RouterGroup.Middlewares.
	Middlewares []gin.HandlerFunc
}

func (r *Route) apply(group *gin.RouterGroup, httpMethod string, middlewares []gin.HandlerFunc) error {
	if r.Handler == nil {
		return fmt.Errorf("handler is nil for route %s", r.Path)
	}
	if r.Middlewares != nil {
		middlewares = append(middlewares, r.Middlewares...)
	}
	group.Handle(httpMethod, r.Path, append(middlewares, r.Handler)...)
	return nil
}

// Apply adds routes to a provided gin router group based on the RouterGroup's configuration.
func (rg *RouterGroup) Apply(parent *gin.RouterGroup) (*gin.RouterGroup, error) {
	if rg.PreFunc != nil {
		if err := rg.PreFunc(); err != nil {
			return nil, err
		}
	}
	group := parent.Group(rg.BasePath)
	group.Use(rg.PersistentMiddlewares...)

	for method, routes := range rg.Routes {
		for _, r := range routes {
			if err := r.apply(group, method, rg.Middlewares); err != nil {
				return nil, err
			}
		}
	}

	for _, rg := range rg.SubGroups {
		if _, err := rg.Apply(group); err != nil {
			return nil, err
		}
	}

	if rg.PostFunc != nil {
		if err := rg.PostFunc(); err != nil {
			return nil, err
		}
	}
	return group, nil
}

// Apply adds routes to a provided gin router group based on the RestRouterGroup's configuration.
func (rrg *RestRouterGroup) Apply(parent *gin.RouterGroup) (*gin.RouterGroup, error) {
	rrg.Routes = make(map[string][]Route)
	routes := map[string]*Route{
		http.MethodGet:    &rrg.GetRoute,
		http.MethodPost:   &rrg.PostRoute,
		http.MethodPatch:  &rrg.PatchRoute,
		http.MethodPut:    &rrg.PutRoute,
		http.MethodDelete: &rrg.DeleteRoute,
	}

	for method, route := range routes {
		if route.Handler != nil {
			rrg.Routes[method] = []Route{*route}
		}
	}

	return rrg.RouterGroup.Apply(parent)
}
