package ginlet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpMethod string

const (
	MethodGet     HttpMethod = http.MethodGet
	MethodPost    HttpMethod = http.MethodPost
	MethodPut     HttpMethod = http.MethodPut
	MethodPatch   HttpMethod = http.MethodPatch
	MethodDelete  HttpMethod = http.MethodDelete
	MethodOptions HttpMethod = http.MethodOptions
	MethodHead    HttpMethod = http.MethodHead
	MethodConnect HttpMethod = http.MethodConnect
	MethodTrace   HttpMethod = http.MethodTrace
)

// Router represents a group of routes that share a common base path and middleware functions.
type Router struct {
	// BasePath is the common prefix for all routes in this group
	BasePath string
	// Middlewares are executed for each route in this group before the route's own middlewares.
	Middlewares []gin.HandlerFunc
	// PersistentMiddlewares are executed for each route before the route's own middlewares and Middlewares.
	PersistentMiddlewares []gin.HandlerFunc
	// Routes is a map where the key is the HTTP method and the value is a slice of Route.
	Routes map[HttpMethod][]Route
	// PreFunc is a function that is executed before Apply is called
	PreFunc func(rg *gin.RouterGroup) error
	// PostFunc is a function that is executed after Apply is called
	PostFunc func(rg *gin.RouterGroup) error
	// SubRouters are nested router groups. Each subgroup inherits BasePath and PersistentMiddlewares from its parent.
	SubRouters []Routable
}

// RestRouter is a specialized version of Router for creating RESTful API endpoints.
type RestRouter struct {
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
	// Router is the base group of routes for the RESTful API. If Router.Routes is not empty, it will be ignored.
	Router
}

// Routable is an interface that all route types must implement to be used with Router and Engine.
type Routable interface {
	Apply(parent *gin.RouterGroup) (*gin.RouterGroup, error)
}

// Route represents a single route
type Route struct {
	// Path is the path for this route
	Path string
	// Handler is the function to execute for this route.
	Handler gin.HandlerFunc
	// Middlewares are functions executed for this route after Router.PersistentMiddlewares and Router.Middlewares.
	Middlewares []gin.HandlerFunc
}

func (r *Route) apply(group *gin.RouterGroup, httpMethod HttpMethod, middlewares []gin.HandlerFunc) error {
	if r.Handler == nil {
		return fmt.Errorf("handler is nil for route %s", r.Path)
	}
	if r.Middlewares != nil {
		middlewares = append(middlewares, r.Middlewares...)
	}
	group.Handle(string(httpMethod), r.Path, append(middlewares, r.Handler)...)
	return nil
}

// Apply adds routes to a provided gin router group based on the Router's configuration.
func (rg *Router) Apply(parent *gin.RouterGroup) (*gin.RouterGroup, error) {
	group := parent.Group(rg.BasePath)
	if rg.PreFunc != nil {
		if err := rg.PreFunc(group); err != nil {
			return nil, err
		}
	}
	group.Use(rg.PersistentMiddlewares...)

	for method, routes := range rg.Routes {
		for _, r := range routes {
			if err := r.apply(group, method, rg.Middlewares); err != nil {
				return nil, err
			}
		}
	}

	for _, rg := range rg.SubRouters {
		if _, err := rg.Apply(group); err != nil {
			return nil, err
		}
	}

	if rg.PostFunc != nil {
		if err := rg.PostFunc(group); err != nil {
			return nil, err
		}
	}
	return group, nil
}

// Apply adds routes to a provided gin router group based on the RestRouter's configuration.
func (rrg *RestRouter) Apply(parent *gin.RouterGroup) (*gin.RouterGroup, error) {
	rrg.Routes = make(map[HttpMethod][]Route)
	routes := map[HttpMethod]*Route{
		MethodGet:    &rrg.GetRoute,
		MethodPost:   &rrg.PostRoute,
		MethodPatch:  &rrg.PatchRoute,
		MethodPut:    &rrg.PutRoute,
		MethodDelete: &rrg.DeleteRoute,
	}

	for method, route := range routes {
		if route.Handler != nil {
			rrg.Routes[method] = []Route{*route}
		}
	}

	return rrg.Router.Apply(parent)
}

// NewEngine creates a new gin.Engine with the provided Routable.
// Keep in mind that the engine is created with gin.New not gin.Default, so you will need to add your own middleware.
func NewEngine(route Routable) (*gin.Engine, error) {
	engine := gin.New()
	if _, err := route.Apply(engine.Group("")); err != nil {
		return nil, err
	}
	return engine, nil
}
