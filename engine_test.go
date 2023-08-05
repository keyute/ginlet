package ginlet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(r *gin.Engine, path string, method string, code int) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	if w.Code != code {
		return fmt.Errorf("expected %d, got %d: %s", code, w.Code, w.Body.String())
	}
	return nil
}

func TestEngine_New(t *testing.T) {
	e := Engine{
		Middlewares: []gin.HandlerFunc{
			func(c *gin.Context) {
				if c.Request.Method == http.MethodDelete {
					c.AbortWithStatusJSON(http.StatusBadRequest, `{"error": "DELETE not allowed"}`)
				}
			},
		},
		RouterGroups: []BaseRoute{
			&RouterGroup{
				Routes: map[string][]Route{
					http.MethodGet: {{Handler: func(c *gin.Context) {
						c.String(http.StatusOK, "GET")
					}},
					}},
			},
			&RouterGroup{
				BasePath: "/test",
				Routes: map[string][]Route{
					http.MethodPost: {{Handler: func(c *gin.Context) {
						c.String(http.StatusOK, "POST")
					}}},
				},
				SubGroups: []BaseRoute{
					&RouterGroup{
						Routes: map[string][]Route{
							http.MethodPatch: {{
								Handler: func(c *gin.Context) {
									c.String(http.StatusOK, "PATCH")
								},
								Path: "/patch",
							}},
						},
					},
				},
			},
		},
	}
	r, err := e.New()
	assert.NotNil(t, r, err)
	assert.NoError(t, testRequest(r, "/", http.MethodGet, http.StatusOK))
	assert.NoError(t, testRequest(r, "/test", http.MethodPost, http.StatusOK))
	assert.NoError(t, testRequest(r, "/test/patch", http.MethodPatch, http.StatusOK))
	assert.NoError(t, testRequest(r, "/test/patch", http.MethodDelete, http.StatusBadRequest))
}
