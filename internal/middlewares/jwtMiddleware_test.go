package middlewares_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/minhmannh2001/authconnecthub/internal/middlewares"
	"github.com/minhmannh2001/authconnecthub/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/undefinedlabs/go-mpatch"
)

func TestIsHtmxRequest_HtmxHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{URL: &url.URL{Path: "/users"}, Method: http.MethodGet}
	c.Request.Header = http.Header{}
	c.Request.Header.Set("HX-Request", "true")

	mockConfig := &config.Config{App: config.App{SwaggerPath: "test/swagger.yaml"}}
	c.Set("config", mockConfig)

	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}}
	patch, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	middlewares.IsHtmxRequest(c)

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, c.Writer.Status())
}

func TestIsHtmxRequest_NoHtmxHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{App: config.App{SwaggerPath: "test/swagger.yaml"}}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsHtmxRequest)

	engine.LoadHTMLGlob("../../tests/templates/*")

	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}}
	patch, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	engine.ServeHTTP(w, req)

	expectedHtml := "\n    <div hx-get=\"/users\" hx-swap=\"outerHTML\" hx-trigger=\"load\"></div>\n"

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedHtml, w.Body.String())

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsHtmxRequest_PathMethodIsNotInSwagger(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{App: config.App{SwaggerPath: "test/swagger.yaml"}}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsHtmxRequest)

	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}}
	patch, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsHtmxRequest_GetSwaggerInfoFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{App: config.App{SwaggerPath: "test/swagger.yaml"}}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsHtmxRequest)

	patch, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return nil, errors.New("swagger error")
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	engine.ServeHTTP(w, req)

	expected_response := "{\"message\":\"error loading swagger information\"}"
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, expected_response, w.Body.String())

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

// https://intellij-support.jetbrains.com/hc/en-us/community/posts/360009685279-Go-test-working-directory-keeps-changing-to-dir-of-the-test-file-instead-of-value-in-template

func TestTriggerHtmxReload_WithQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{App: config.App{SwaggerPath: "test/swagger.yaml"}}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsHtmxRequest)

	engine.LoadHTMLGlob("../../tests/templates/*")

	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}}
	patch, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users?param1=value1&param2=value2", nil)
	engine.ServeHTTP(w, req)

	expectedHtml := "\n    <div hx-get=\"/users?param1=value1&amp;param2=value2\" hx-swap=\"outerHTML\" hx-trigger=\"load\"></div>\n"

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedHtml, w.Body.String())

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsAuthorized_PublicRouteNoToken(t *testing.T) {
	mockAuth := mocks.NewIAuthUC(t)

	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{App: config.App{SwaggerPath: "test/swagger.yaml"}}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsAuthorized(mockAuth))
	engine.GET("/users", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{}}}}}
	patch, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsAuthorized_PublicRouteWithToken(t *testing.T) {
	mockAuth := mocks.NewIAuthUC(t)
	mockAuth.On("ValidateToken", mock.Anything).Return("minhmannh2001", nil)

	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{
		App: config.App{SwaggerPath: "test/swagger.yaml"},
	}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsAuthorized(mockAuth))
	engine.GET("/users", func(c *gin.Context) {
		username, _ := c.Get("username")
		c.Data(http.StatusOK, "text/plain", []byte(username.(string)))
		c.Status(http.StatusOK)
	})

	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{}}}}}
	patch, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "accessToken"))
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "minhmannh2001", w.Body.String())

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsAuthorized_GetSwaggerInfoFails(t *testing.T) {
	mockAuth := mocks.NewIAuthUC(t)

	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{
		App: config.App{SwaggerPath: "test/swagger.yaml"},
	}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsAuthorized(mockAuth))
	patch, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return nil, errors.New("swagger error")
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	engine.ServeHTTP(w, req)

	assert.Equal(t, "/500.html", w.Result().Header.Get("HX-Redirect"))

	err = patch.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsAuthorized_PrivateRouteMissingToken(t *testing.T) {
	mockAuth := mocks.NewIAuthUC(t)

	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{
		App: config.App{SwaggerPath: "test/swagger.yaml"},
	}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsAuthorized(mockAuth))
	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}}
	patchGetSwaggerInfo, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	patchHashMap, err := mpatch.PatchMethod(helper.HashMap, func(data map[string]interface{}) (string, error) {
		return "abc", nil
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	engine.ServeHTTP(w, req)

	expected := "/v1/auth/login?toast-message=login-is-required-for-this-action.-sign-in-or-create-an-account-to-continue.&toast-type=danger&hash-value=abc"
	assert.Equal(t, expected, w.Result().Header.Get("HX-Redirect"))

	err = patchGetSwaggerInfo.Unpatch()
	if err != nil {
		t.Fatal(err)
	}

	err = patchHashMap.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsAuthorized_PrivateRouteBlacklistedToken(t *testing.T) {
	mockAuth := mocks.NewIAuthUC(t)
	mockAuth.On("IsTokenBlacklisted", mock.Anything).Return(true, nil)
	mockAuth.On("RetrieveFieldFromJwtToken", mock.Anything, mock.Anything, mock.Anything).Return(interface{}(true), nil)

	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{
		App: config.App{SwaggerPath: "test/swagger.yaml"},
	}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsAuthorized(mockAuth))

	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}}
	patchGetSwaggerInfo, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	patchHashMap, err := mpatch.PatchMethod(helper.HashMap, func(data map[string]interface{}) (string, error) {
		return "abc", nil
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "accessToken"))
	engine.ServeHTTP(w, req)

	expectedHXRedirect := "/v1/auth/login?toast-message=your-token-is-invalid.-please-log-in-to-continue.&toast-type=danger&hash-value=abc"
	expectedHXTrigger := "{\"deleteToken\":{\"deleteFrom\":\"local\"}}"
	assert.Equal(t, expectedHXRedirect, w.Result().Header.Get("HX-Redirect"))
	assert.Equal(t, expectedHXTrigger, w.Result().Header.Get("HX-Trigger"))

	err = patchGetSwaggerInfo.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	err = patchHashMap.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsAuthorized_PrivateRouteValidToken(t *testing.T) {
	mockAuth := mocks.NewIAuthUC(t)
	mockAuth.On("IsTokenBlacklisted", mock.Anything).Return(false, nil)
	mockAuth.On("ValidateToken", mock.Anything).Return("minhmannh2001", nil)

	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{
		App: config.App{SwaggerPath: "test/swagger.yaml"},
	}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsAuthorized(mockAuth))
	engine.GET("/users", func(c *gin.Context) {
		username, _ := c.Get("username")
		c.Data(http.StatusOK, "text/plain", []byte(username.(string)))
		c.Status(http.StatusOK)
	})

	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}}
	patchGetSwaggerInfo, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "accessToken"))
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "minhmannh2001", w.Body.String())

	err = patchGetSwaggerInfo.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsAuthorized_PrivateRouteExpiredTokenLogout(t *testing.T) {
	mockAuth := mocks.NewIAuthUC(t)
	mockAuth.On("IsTokenBlacklisted", mock.Anything).Return(false, nil)
	mockAuth.On("ValidateToken", mock.Anything).Return("", errors.New("token expired"))

	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{
		App: config.App{SwaggerPath: "test/swagger.yaml"},
	}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsAuthorized(mockAuth))
	engine.GET("/v1/auth/logout", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/v1/auth/logout": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}}
	patchGetSwaggerInfo, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/v1/auth/logout", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "accessToken"))
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = patchGetSwaggerInfo.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}

func TestIsAuthorized_PrivateRouteExpiredToken(t *testing.T) {
	mockAuth := mocks.NewIAuthUC(t)
	mockAuth.On("IsTokenBlacklisted", mock.Anything).Return(false, nil)
	mockAuth.On("ValidateToken", mock.Anything).Return("", errors.New("token expired"))
	mockAuth.On("RetrieveFieldFromJwtToken", mock.Anything, mock.Anything, mock.Anything).Return(interface{}(true), nil)

	gin.SetMode(gin.TestMode)
	_, engine := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{
		App: config.App{SwaggerPath: "test/swagger.yaml"},
	}
	engine.Use(func(c *gin.Context) {
		c.Set("config", mockConfig)
		c.Next()
	})
	engine.Use(middlewares.IsAuthorized(mockAuth))

	mockSwaggerInfo := &entity.SwaggerInfo{Paths: map[string]entity.PathItem{"/users": {Get: &entity.Operation{Security: []interface{}{map[string]interface{}{"JWT": nil}}}}}}
	patchGetSwaggerInfo, err := mpatch.PatchMethod(helper.GetSwaggerInfo, func(filePath string) (*entity.SwaggerInfo, error) {
		return mockSwaggerInfo, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	patchHashMap, err := mpatch.PatchMethod(helper.HashMap, func(data map[string]interface{}) (string, error) {
		return "abc", nil
	})
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "accessToken"))
	engine.ServeHTTP(w, req)

	expectedHXRedirect := "/v1/auth/login?toast-message=your-session-has-expired.-please-log-in-to-continue.&toast-type=danger&hash-value=abc"
	assert.Equal(t, expectedHXRedirect, w.Header().Get("HX-Redirect"))

	err = patchGetSwaggerInfo.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
	err = patchHashMap.Unpatch()
	if err != nil {
		t.Fatal(err)
	}
}
