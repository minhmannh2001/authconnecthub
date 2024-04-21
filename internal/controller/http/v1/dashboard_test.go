package v1_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	v1 "github.com/minhmannh2001/authconnecthub/internal/controller/http/v1"
	"github.com/minhmannh2001/authconnecthub/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func TestDashboardRoutes_getDashboardHandler(t *testing.T) {
	// Test case 1: Test with no HX-Reload header
	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	funcMap := template.FuncMap{
		"Capitalize": capitalize,
	}
	engine.SetFuncMap(funcMap)
	engine.LoadHTMLGlob("../../../../templates/*")
	v1.NewDashboardRoutes(engine.Group("/v1"), logger.NewLogger(), nil, nil, nil)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/dashboard", nil)
	engine.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/html; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "<!DOCTYPE html>")

	// Test case 2: Test with HX-Reload header
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/v1/dashboard", nil)
	req2.Header.Set("HX-Reload", "true")
	engine.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "text/html; charset=utf-8", w2.Header().Get("Content-Type"))
	assert.NotContains(t, w2.Body.String(), "<!DOCTYPE html>")
}
