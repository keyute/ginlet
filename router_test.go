package gimlet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterGroup_Get(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		BasePath: "/test",
		GetHandler: func(c *gin.Context) {
			c.Status(http.StatusOK)
		},
	}
	_, err := group.New(router.Group(""))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterGroup_Post(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		BasePath: "/test",
		PostHandler: func(c *gin.Context) {
			c.Status(http.StatusOK)
		},
	}
	_, err := group.New(router.Group(""))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterGroup_Patch(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		BasePath: "/test",
		PatchHandler: func(c *gin.Context) {
			c.Status(http.StatusOK)
		},
	}
	_, err := group.New(router.Group(""))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterGroup_Delete(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		BasePath: "/test",
		DeleteHandler: func(c *gin.Context) {
			c.Status(http.StatusOK)
		},
	}
	_, err := group.New(router.Group(""))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/test", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterGroup_Nested(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		BasePath: "/test",
		RouterGroups: []RouterGroup{
			{
				BasePath: "/nested",
				GetHandler: func(c *gin.Context) {
					c.Status(http.StatusOK)
				},
			},
		},
	}
	_, err := group.New(router.Group(""))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test/nested", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouterGroup_NestedWithError(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		RouterGroups: []RouterGroup{
			{
				PreFunc: func() error {
					return fmt.Errorf("error")
				},
			},
		},
	}
	_, err := group.New(router.Group(""))
	assert.Errorf(t, err, "error")
}

func TestRouterGroup_PreFunc(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		PreFunc: func() error {
			return fmt.Errorf("error")
		},
	}
	_, err := group.New(router.Group(""))
	assert.Errorf(t, err, "error")
}

func TestRouterGroup_PostFunc(t *testing.T) {
	router := gin.Default()
	group := RouterGroup{
		PostFunc: func() error {
			return fmt.Errorf("error")
		},
	}
	_, err := group.New(router.Group(""))
	assert.Errorf(t, err, "error")
}
