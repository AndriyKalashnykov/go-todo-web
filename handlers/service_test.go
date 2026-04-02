package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AndriyKalashnykov/go-todo-web/handlers"
	"github.com/AndriyKalashnykov/go-todo-web/models"
	"github.com/labstack/echo/v5"
)

func setupRouter() *echo.Echo {
	e := echo.New()
	e.POST("/create", handlers.CreateTodo)
	e.GET("/get/:id", handlers.GetTodo)
	e.GET("/all", handlers.Todos)
	e.DELETE("/delete/:id", handlers.DeleteTodo)
	e.PATCH("/update/:id", handlers.UpdateTodo)
	return e
}

func createTodoViaAPI(t *testing.T, e *echo.Echo, task string) {
	t.Helper()
	body := `{"task":"` + task + `"}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("setup: create todo failed with status %d: %s", rec.Code, rec.Body.String())
	}
}

func TestCreateTodo(t *testing.T) {
	models.Reset()
	e := setupRouter()

	body := `{"task":"Write tests","completed":false}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var todo models.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &todo); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if todo.Task != "Write tests" {
		t.Errorf("expected task 'Write tests', got %q", todo.Task)
	}
}

func TestCreateTodoInvalidJSON(t *testing.T) {
	models.Reset()
	e := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetAllTodos(t *testing.T) {
	models.Reset()
	e := setupRouter()

	createTodoViaAPI(t, e, "first")
	createTodoViaAPI(t, e, "second")

	req := httptest.NewRequest(http.MethodGet, "/all", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var todos []models.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &todos); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(todos) != 2 {
		t.Errorf("expected 2 todos, got %d", len(todos))
	}
}

func TestGetAllTodosEmpty(t *testing.T) {
	models.Reset()
	e := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/all", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestGetTodo(t *testing.T) {
	models.Reset()
	e := setupRouter()
	createTodoViaAPI(t, e, "my task")

	req := httptest.NewRequest(http.MethodGet, "/get/1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var todo models.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &todo); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if todo.Task != "my task" {
		t.Errorf("expected task 'my task', got %q", todo.Task)
	}
	if todo.ID != 1 {
		t.Errorf("expected ID 1, got %d", todo.ID)
	}
}

func TestGetTodoNotFound(t *testing.T) {
	models.Reset()
	e := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/get/999", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetTodoInvalidID(t *testing.T) {
	models.Reset()
	e := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/get/abc", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// strconv.Atoi("abc") returns 0, and ID 0 doesn't exist
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d for invalid ID, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestDeleteTodo(t *testing.T) {
	models.Reset()
	e := setupRouter()
	createTodoViaAPI(t, e, "to delete")

	req := httptest.NewRequest(http.MethodDelete, "/delete/1", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	// Verify it's actually gone
	req = httptest.NewRequest(http.MethodGet, "/get/1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected deleted todo to return %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestDeleteTodoNotFound(t *testing.T) {
	models.Reset()
	e := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/delete/999", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUpdateTodo(t *testing.T) {
	models.Reset()
	e := setupRouter()
	createTodoViaAPI(t, e, "original task")

	body := `{"task":"updated task","completed":true}`
	req := httptest.NewRequest(http.MethodPatch, "/update/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var todo models.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &todo); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if todo.ID != 1 {
		t.Errorf("expected ID 1, got %d", todo.ID)
	}
}

func TestUpdateTodoInvalidJSON(t *testing.T) {
	models.Reset()
	e := setupRouter()
	createTodoViaAPI(t, e, "exists")

	req := httptest.NewRequest(http.MethodPatch, "/update/1", strings.NewReader(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for invalid JSON, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateTodoNotFound(t *testing.T) {
	models.Reset()
	e := setupRouter()

	body := `{"task":"nope"}`
	req := httptest.NewRequest(http.MethodPatch, "/update/999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestCRUDWorkflow(t *testing.T) {
	models.Reset()
	e := setupRouter()

	// Create
	body := `{"task":"integration test"}`
	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: expected %d, got %d", http.StatusCreated, rec.Code)
	}

	// Read
	req = httptest.NewRequest(http.MethodGet, "/get/1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("get: expected %d, got %d", http.StatusOK, rec.Code)
	}
	var todo models.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &todo); err != nil {
		t.Fatalf("get: failed to parse: %v", err)
	}
	if todo.Task != "integration test" {
		t.Errorf("get: expected task 'integration test', got %q", todo.Task)
	}

	// List
	req = httptest.NewRequest(http.MethodGet, "/all", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("all: expected %d, got %d", http.StatusOK, rec.Code)
	}

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/delete/1", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete: expected %d, got %d", http.StatusNoContent, rec.Code)
	}

	// Verify deleted
	req = httptest.NewRequest(http.MethodGet, "/all", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	var remaining []models.Todo
	if err := json.Unmarshal(rec.Body.Bytes(), &remaining); err != nil {
		t.Fatalf("all after delete: failed to parse: %v", err)
	}
	if len(remaining) != 0 {
		t.Errorf("expected 0 todos after delete, got %d", len(remaining))
	}
}
