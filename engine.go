package ginlet

import "github.com/gin-gonic/gin"

// Engine is a struct that contains all the information needed to create a gin.Engine.
type Engine struct {
	// Middlewares are the global middlewares that will be applied to all routes.
	Middlewares []gin.HandlerFunc
	// RouterGroups contains all the RouterGroups to be used with the gin engine.
	RouterGroups []RouterGroup
	// PreFunc is a function that is executed before New is called
	PreFunc func(r *gin.Engine) error
	// PostFunc is a function that is executed after New is called
	PostFunc func(r *gin.Engine) error
}

// New creates a gin.Engine from the information in the Engine struct.
func (e *Engine) New() (*gin.Engine, error) {
	r := gin.New()
	if e.PreFunc != nil {
		if err := e.PreFunc(r); err != nil {
			return nil, err
		}
	}

	r.Use(e.Middlewares...)
	controller := r.Group("")
	for _, rg := range e.RouterGroups {
		if _, err := rg.Apply(controller); err != nil {
			return nil, err
		}
	}

	if e.PostFunc != nil {
		if err := e.PostFunc(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}
