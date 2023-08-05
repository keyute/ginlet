package gimlet

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// RouterGroup is a struct that represents a group of routes
type RouterGroup struct {
	// BasePath is the common prefix for all routes in this group
	BasePath string

	// Middlewares are functions that are executed for each route in this group before the route itself
	Middlewares []gin.HandlerFunc

	// PersistentMiddlewares are functions that are executed for each route in this group before the route itself
	// including the ones that might be added to the group later
	PersistentMiddlewares []gin.HandlerFunc

	// GetHandler is the handler for GET requests
	GetHandler func(c *gin.Context)

	// GetPath is the path for GET requests
	GetPath string

	// PostHandler is the handler for POST requests
	PostHandler func(c *gin.Context)

	// PostPath is the path for POST requests
	PostPath string

	// PatchHandler is the handler for PATCH requests
	PatchHandler func(c *gin.Context)

	// PatchPath is the path for PATCH requests
	PatchPath string

	// PutHandler is the handler for PUT requests
	PutHandler func(c *gin.Context)

	// PutPath is the path for PUT requests
	PutPath string

	// DeleteHandler is the handler for DELETE requests
	DeleteHandler func(c *gin.Context)

	// DeletePath is the path for DELETE requests
	DeletePath string

	// PreFunc is a function that is executed before New is called
	PreFunc func() error

	// PostFunc is a function that is executed after New is called
	PostFunc func() error

	// RouterGroups are the subgroups of this group
	RouterGroups []RouterGroup
}

func createRoute(controller *gin.RouterGroup, httpMethod string, path string, handler func(c *gin.Context), middlewares []gin.HandlerFunc) {
	if handler != nil {
		middlewares = append(middlewares, handler)
	}
	controller.Handle(httpMethod, path, middlewares...)
}

// New adds routes to a provided gin router group based on the RouterGroup's configuration.
func (rg *RouterGroup) New(parent *gin.RouterGroup) (*gin.RouterGroup, error) {
	if rg.PreFunc != nil {
		if err := rg.PreFunc(); err != nil {
			return nil, err
		}
	}
	controller := parent.Group(rg.BasePath)
	controller.Use(rg.PersistentMiddlewares...)

	createRoute(controller, http.MethodGet, rg.GetPath, rg.GetHandler, rg.Middlewares)
	createRoute(controller, http.MethodPost, rg.PostPath, rg.PostHandler, rg.Middlewares)
	createRoute(controller, http.MethodPatch, rg.PatchPath, rg.PatchHandler, rg.Middlewares)
	createRoute(controller, http.MethodPut, rg.PutPath, rg.PutHandler, rg.Middlewares)
	createRoute(controller, http.MethodDelete, rg.DeletePath, rg.DeleteHandler, rg.Middlewares)

	for _, rg := range rg.RouterGroups {
		if _, err := rg.New(controller); err != nil {
			return nil, err
		}
	}

	if rg.PostFunc != nil {
		if err := rg.PostFunc(); err != nil {
			return nil, err
		}
	}
	return controller, nil
}
