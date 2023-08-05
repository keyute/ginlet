package gimlet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterGroup_Apply(t *testing.T) {
	router := gin.Default()
	handler := func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
	path := "/test"
	group := RouterGroup{
		BasePath: path,
		Routes: map[string][]Route{
			http.MethodGet:    {{Handler: handler}},
			http.MethodPost:   {{Handler: handler}},
			http.MethodPatch:  {{Handler: handler}},
			http.MethodDelete: {{Handler: handler}},
		},
	}
	_, err := group.Apply(router.Group(""))
	assert.NoError(t, err)
	assert.NoError(t, testRequest(router, path, http.MethodGet, http.StatusOK))
	assert.NoError(t, testRequest(router, path, http.MethodPost, http.StatusOK))
	assert.NoError(t, testRequest(router, path, http.MethodPatch, http.StatusOK))
	assert.NoError(t, testRequest(router, path, http.MethodDelete, http.StatusOK))
}

func TestRouterGroup_ApplyError(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		Routes: map[string][]Route{
			http.MethodGet: {{}},
		},
	}
	_, err := group.Apply(router.Group(""))
	assert.Errorf(t, err, "handler is nil for route /")
}

func TestRouterGroup_Nested(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		BasePath: "/test",
		SubGroups: []RouterGroup{{
			BasePath: "/nested",
			Routes: map[string][]Route{
				http.MethodGet: {{Handler: func(c *gin.Context) {
					c.Status(http.StatusOK)
				}}},
			},
		}},
	}
	_, err := group.Apply(router.Group(""))
	assert.NoError(t, err)
	assert.NoError(t, testRequest(router, "/test/nested", http.MethodGet, http.StatusOK))
}

func TestRouterGroup_NestedWithError(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		SubGroups: []RouterGroup{
			{
				PreFunc: func() error {
					return fmt.Errorf("error")
				},
			},
		},
	}
	_, err := group.Apply(router.Group(""))
	assert.Errorf(t, err, "error")
}

func TestRouterGroup_PreFunc(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		PreFunc: func() error {
			return fmt.Errorf("error")
		},
	}
	_, err := group.Apply(router.Group(""))
	assert.Errorf(t, err, "error")
}

func TestRouterGroup_PostFunc(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		PostFunc: func() error {
			return fmt.Errorf("error")
		},
	}
	_, err := group.Apply(router.Group(""))
	assert.Errorf(t, err, "error")
}

func TestRouterGroup_Middleware(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
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
	_, err := group.Apply(router.Group(""))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test", w.Body.String())
}

func TestRestRouterGroup(t *testing.T) {
	router := gin.Default()
	group := RestRouterGroup{
		GetRoute: Route{
			Handler: func(c *gin.Context) {
				c.Status(http.StatusOK)
			},
			Path: "/test",
		},
	}
	_, err := group.Apply(router.Group(""))
	assert.NoError(t, err)
	assert.NoError(t, testRequest(router, "/test", http.MethodGet, http.StatusOK))
}
