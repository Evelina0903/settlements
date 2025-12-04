package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterCreation(t *testing.T) {
	r := New()

	if r == nil {
		t.Fatalf("Expected router instance, got nil")
	}

	if r.root == nil {
		t.Errorf("Expected root node to be initialized")
	}

	if len(r.middlewares) != 0 {
		t.Errorf("Expected 0 middlewares initially, got %d", len(r.middlewares))
	}
}

func TestRouterHandleStaticRoute(t *testing.T) {
	r := New()
	handlerCalled := false

	handler := func(w http.ResponseWriter, req *http.Request, params Params) {
		handlerCalled = true
	}

	r.GET("/test", handler)

	// Test matching
	h, params := r.match("GET", "/test")

	if h == nil {
		t.Fatalf("Expected handler, got nil")
	}

	if params == nil {
		t.Errorf("Expected params map, got nil")
	}

	h(nil, nil, params)
	if !handlerCalled {
		t.Errorf("Expected handler to be called")
	}
}

func TestRouterDynamicParameter(t *testing.T) {
	r := New()
	handler := func(w http.ResponseWriter, req *http.Request, params Params) {}

	r.GET("/users/:id", handler)

	h, params := r.match("GET", "/users/123")

	if h == nil {
		t.Fatalf("Expected handler for /users/123, got nil")
	}

	if params == nil {
		t.Fatalf("Expected params, got nil")
	}

	if params["id"] != "123" {
		t.Errorf("Expected id='123', got '%s'", params["id"])
	}
}

func TestRouterMultipleDynamicParameters(t *testing.T) {
	r := New()
	handler := func(w http.ResponseWriter, req *http.Request, params Params) {}

	r.GET("/users/:id/posts/:postid", handler)

	h, params := r.match("GET", "/users/42/posts/99")

	if h == nil {
		t.Fatalf("Expected handler, got nil")
	}

	if params["id"] != "42" {
		t.Errorf("Expected id='42', got '%s'", params["id"])
	}

	if params["postid"] != "99" {
		t.Errorf("Expected postid='99', got '%s'", params["postid"])
	}
}

func TestRouterNotFound(t *testing.T) {
	r := New()
	handler := func(w http.ResponseWriter, req *http.Request, params Params) {}

	r.GET("/existing", handler)

	h, _ := r.match("GET", "/nonexistent")

	if h != nil {
		t.Errorf("Expected nil handler for nonexistent route, got handler")
	}
}

func TestRouterMethodNotAllowed(t *testing.T) {
	r := New()
	handler := func(w http.ResponseWriter, req *http.Request, params Params) {}

	r.GET("/test", handler)

	// Try different method
	h, _ := r.match("POST", "/test")

	if h != nil {
		t.Errorf("Expected nil handler for POST on GET route, got handler")
	}
}

func TestRouterMultipleMethods(t *testing.T) {
	r := New()
	getHandler := func(w http.ResponseWriter, req *http.Request, params Params) {}
	postHandler := func(w http.ResponseWriter, req *http.Request, params Params) {}

	r.GET("/test", getHandler)
	r.POST("/test", postHandler)

	h1, _ := r.match("GET", "/test")
	h2, _ := r.match("POST", "/test")

	if h1 == nil {
		t.Errorf("Expected GET handler, got nil")
	}

	if h2 == nil {
		t.Errorf("Expected POST handler, got nil")
	}
}

func TestRouterServeHTTP(t *testing.T) {
	r := New()
	handlerCalled := false

	handler := func(w http.ResponseWriter, req *http.Request, params Params) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}

	r.GET("/test", handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if !handlerCalled {
		t.Errorf("Expected handler to be called via ServeHTTP")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRouterMiddleware(t *testing.T) {
	r := New()
	middlewareCalled := false

	middleware := func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request, params Params) {
			middlewareCalled = true
			next(w, req, params)
		}
	}

	handler := func(w http.ResponseWriter, req *http.Request, params Params) {
		w.WriteHeader(http.StatusOK)
	}

	r.Use(middleware)
	r.GET("/test", handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if !middlewareCalled {
		t.Errorf("Expected middleware to be called")
	}
}

func TestRouterMultipleMiddleware(t *testing.T) {
	r := New()
	order := []string{}

	middleware1 := func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request, params Params) {
			order = append(order, "m1")
			next(w, req, params)
		}
	}

	middleware2 := func(next HandlerFunc) HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request, params Params) {
			order = append(order, "m2")
			next(w, req, params)
		}
	}

	handler := func(w http.ResponseWriter, req *http.Request, params Params) {
		order = append(order, "handler")
	}

	r.Use(middleware1)
	r.Use(middleware2)
	r.GET("/test", handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Middleware applied in reverse order (wrapping in reverse)
	if len(order) != 3 {
		t.Errorf("Expected 3 calls, got %d", len(order))
	}

	if order[0] != "m1" {
		t.Errorf("Expected m1 first, got %s", order[0])
	}

	if order[1] != "m2" {
		t.Errorf("Expected m2 second, got %s", order[1])
	}

	if order[2] != "handler" {
		t.Errorf("Expected handler last, got %s", order[2])
	}
}

func TestRouterCatchAll(t *testing.T) {
	r := New()
	handler := func(w http.ResponseWriter, req *http.Request, params Params) {}

	r.GET("/api/*rest", handler)

	h, params := r.match("GET", "/api/v1/users/123/posts")

	if h == nil {
		t.Fatalf("Expected handler for catch-all, got nil")
	}

	if params["rest"] != "v1/users/123/posts" {
		t.Errorf("Expected rest='v1/users/123/posts', got '%s'", params["rest"])
	}
}

func TestRouterStaticChildrenPriority(t *testing.T) {
	r := New()
	staticHandler := func(w http.ResponseWriter, req *http.Request, params Params) {}
	paramHandler := func(w http.ResponseWriter, req *http.Request, params Params) {}

	r.GET("/users/me", staticHandler)
	r.GET("/users/:id", paramHandler)

	h1, _ := r.match("GET", "/users/me")
	h2, _ := r.match("GET", "/users/123")

	if h1 == nil {
		t.Errorf("Expected handler for /users/me")
	}

	if h2 == nil {
		t.Errorf("Expected handler for /users/123")
	}
}

func TestRouterNotFoundStatus(t *testing.T) {
	r := New()
	handler := func(w http.ResponseWriter, req *http.Request, params Params) {}
	r.GET("/test", handler)

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404 status, got %d", w.Code)
	}
}

func TestRouterMethodNotAllowedStatus(t *testing.T) {
	r := New()
	handler := func(w http.ResponseWriter, req *http.Request, params Params) {}
	r.GET("/test", handler)

	req := httptest.NewRequest("POST", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 status, got %d", w.Code)
	}
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		input  string
		expect []string
	}{
		{"/", []string{}},
		{"", []string{}},
		{"/test", []string{"test"}},
		{"/test/path", []string{"test", "path"}},
		{"/test/path/", []string{"test", "path"}},
	}

	for _, test := range tests {
		result := splitPath(test.input)
		if len(result) != len(test.expect) {
			t.Errorf("splitPath(%q) expected %d segments, got %d", test.input, len(test.expect), len(result))
		}
		for i, seg := range result {
			if seg != test.expect[i] {
				t.Errorf("splitPath(%q) segment %d: expected %q, got %q", test.input, i, test.expect[i], seg)
			}
		}
	}
}
