package gimlet

import "github.com/gin-gonic/gin"

// Engine is a struct that contains all the information needed to create a gin.Engine.
type Engine struct {
	Middlewares  []gin.HandlerFunc
	RouterGroups []RouterGroup
}

// New creates a gin.Engine from the information in the Engine struct.
func (e *Engine) New() (*gin.Engine, error) {
	r := gin.New()
	r.Use(e.Middlewares...)
	controller := r.Group("")
	for _, rg := range e.RouterGroups {
		if _, err := rg.Apply(controller); err != nil {
			return nil, err
		}
	}
	return r, nil
}
