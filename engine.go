package gimlet

import "github.com/gin-gonic/gin"

// Engine is a struct that contains all the information needed to create a gin.Engine.
type Engine struct {
	// Middlewares are the global middlewares that will be applied to all routes.
	Middlewares []gin.HandlerFunc
	// RouterGroups contains all the RouterGroups to be used with the gin engine.
	RouterGroups []RouterGroup
	// PreFunc is a function that is executed before New is called
	PreFunc func() error
	// PostFunc is a function that is executed after New is called
	PostFunc func() error
}

// New creates a gin.Engine from the information in the Engine struct.
func (e *Engine) New() (*gin.Engine, error) {
	if e.PreFunc != nil {
		if err := e.PreFunc(); err != nil {
			return nil, err
		}
	}

	r := gin.New()
	r.Use(e.Middlewares...)
	controller := r.Group("")
	for _, rg := range e.RouterGroups {
		if _, err := rg.Apply(controller); err != nil {
			return nil, err
		}
	}

	if e.PostFunc != nil {
		if err := e.PostFunc(); err != nil {
			return nil, err
		}
	}
	return r, nil
}
