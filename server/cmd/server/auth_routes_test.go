package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProtectedEndpointsRequireAuthorization(t *testing.T) {
	srv := NewServer(":0", nil)
	handler := srv.router()

	protectedCases := []struct {
		name   string
		method string
		path   string
		body   []byte
	}{
		{name: "swagger ui", method: http.MethodGet, path: "/swagger/index.html"},
		{name: "get profile", method: http.MethodGet, path: "/api/v1/profile"},
		{name: "update profile", method: http.MethodPut, path: "/api/v1/profile", body: []byte(`{}`)},
		{name: "update password", method: http.MethodPut, path: "/api/v1/profile/password", body: []byte(`{}`)},
		{name: "list users", method: http.MethodGet, path: "/api/v1/users"},
		{name: "get user by id", method: http.MethodGet, path: "/api/v1/users/1"},
		{name: "user lookup", method: http.MethodGet, path: "/api/v1/users/lookup"},
		{name: "list set types", method: http.MethodGet, path: "/api/v1/set_types"},
		{name: "list project types", method: http.MethodGet, path: "/api/v1/project_types"},
		{name: "list warehouses", method: http.MethodGet, path: "/api/v1/warehouse"},
		{name: "list equipment sets", method: http.MethodGet, path: "/api/v1/equipment_set"},
		{name: "list equipment", method: http.MethodGet, path: "/api/v1/equipment"},
		{name: "list projects", method: http.MethodGet, path: "/api/v1/projects"},
		{name: "list drafts", method: http.MethodGet, path: "/api/v1/drafts"},
		{name: "equipment in project", method: http.MethodGet, path: "/api/v1/equipment_in_project/1"},
		{name: "equipment in draft", method: http.MethodGet, path: "/api/v1/equipment_in_draft/1"},
	}

	for _, tc := range protectedCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusForbidden {
				t.Fatalf("expected %d, got %d", http.StatusForbidden, rr.Code)
			}
		})
	}
}

func TestLoginAndRegisterArePublic(t *testing.T) {
	srv := NewServer(":0", nil)
	handler := srv.router()

	cases := []struct {
		name string
		path string
	}{
		{name: "login", path: "/api/v1/login"},
		{name: "register", path: "/api/v1/register"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.path, bytes.NewBufferString(`{}`))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code == http.StatusForbidden {
				t.Fatalf("expected non-%d for public endpoint, got %d", http.StatusForbidden, rr.Code)
			}
		})
	}
}
