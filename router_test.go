package ginlet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterGroup_Apply(t *testing.T) {
	handler := func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
	path := "/test"
	group := Router{
		BasePath: path,
		Routes: map[string][]Route{
			http.MethodGet:    {{Handler: handler}},
			http.MethodPost:   {{Handler: handler}},
			http.MethodPatch:  {{Handler: handler}},
			http.MethodDelete: {{Handler: handler}},
		},
	}
	router, err := NewEngine(&group)
	assert.NotNil(t, router, err)
	assert.NoError(t, testRequest(router, path, http.MethodGet, http.StatusOK))
	assert.NoError(t, testRequest(router, path, http.MethodPost, http.StatusOK))
	assert.NoError(t, testRequest(router, path, http.MethodPatch, http.StatusOK))
	assert.NoError(t, testRequest(router, path, http.MethodDelete, http.StatusOK))
}

func TestRouterGroup_ApplyError(t *testing.T) {
	group := Router{
		Routes: map[string][]Route{
			http.MethodGet: {{}},
		},
	}
	_, err := NewEngine(&group)
	assert.Errorf(t, err, "handler is nil for route /")
}

func TestRouterGroup_Nested(t *testing.T) {
	group := Router{
		BasePath: "/test",
		SubRouters: []Routable{&Router{
			BasePath: "/nested",
			Routes: map[string][]Route{
				http.MethodGet: {{Handler: func(c *gin.Context) {
					c.Status(http.StatusOK)
				}}},
			},
		}},
	}
	router, err := NewEngine(&group)
	assert.NotNil(t, router, err)
	assert.NoError(t, testRequest(router, "/test/nested", http.MethodGet, http.StatusOK))
}

func TestRouterGroup_NestedWithError(t *testing.T) {
	group := Router{
		SubRouters: []Routable{
			&Router{
				PreFunc: func(rg *gin.RouterGroup) error {
					return fmt.Errorf("error")
				},
			},
		},
	}
	_, err := NewEngine(&group)
	assert.Errorf(t, err, "error")
}

func TestRouterGroup_PreFunc(t *testing.T) {
	group := Router{
		PreFunc: func(rg *gin.RouterGroup) error {
			return fmt.Errorf("error")
		},
	}
	_, err := NewEngine(&group)
	assert.Errorf(t, err, "error")
}

func TestRouterGroup_PostFunc(t *testing.T) {
	group := Router{
		PostFunc: func(rg *gin.RouterGroup) error {
			return fmt.Errorf("error")
		},
	}
	_, err := NewEngine(&group)
	assert.Errorf(t, err, "error")
}

func TestRouterGroup_Middleware(t *testing.T) {
	group := Router{
		Routes: map[string][]Route{
			http.MethodGet: {{
				Handler: func(c *gin.Context) {
					c.String(http.StatusOK, c.GetString("test"))
				},
				Middlewares: []gin.HandlerFunc{
					func(c *gin.Context) {
						c.Set("test", "test")
					},
				},
			}},
		},
	}
	router, err := NewEngine(&group)
	assert.NotNil(t, router, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test", w.Body.String())
}

func TestRestRouterGroup(t *testing.T) {
	group := RestRouter{
		GetRoute: Route{
			Handler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			Path: "/test",
		},
	}
	router, err := NewEngine(&group)
	assert.NotNil(t, router, err)
	assert.NoError(t, testRequest(router, "/test", http.MethodGet, http.StatusOK))
}

func testRequest(r *gin.Engine, path string, method string, code int) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	if w.Code != code {
		return fmt.Errorf("expected %d, got %d: %s", code, w.Code, w.Body.String())
	}
	return nil
}
