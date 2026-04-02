package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewApp(t *testing.T) {
	app := NewApp()
	if app == nil {
		t.Fatal("expected non-nil Echo instance")
	}

	req := httptest.NewRequest(http.MethodGet, "/all", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d for /all, got %d", http.StatusOK, rec.Code)
	}
}
