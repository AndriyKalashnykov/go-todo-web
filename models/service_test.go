package models

import (
	"testing"
)

func resetState() {
	todos = nil
	currentID = 0
}

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{Message: "item not found"}
	if err.Error() != "item not found" {
		t.Errorf("expected 'item not found', got %q", err.Error())
	}
}

func TestAddTodo(t *testing.T) {
	resetState()

	AddTodo(Todo{Task: "Buy groceries"})

	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}
	if todos[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", todos[0].ID)
	}
	if todos[0].Task != "Buy groceries" {
		t.Errorf("expected task 'Buy groceries', got %q", todos[0].Task)
	}
	if todos[0].Completed {
		t.Error("expected Completed to be false for new todo")
	}
}

func TestAddTodoIncrementsID(t *testing.T) {
	resetState()

	AddTodo(Todo{Task: "first"})
	AddTodo(Todo{Task: "second"})
	AddTodo(Todo{Task: "third"})

	if len(todos) != 3 {
		t.Fatalf("expected 3 todos, got %d", len(todos))
	}
	for i, want := range []int{1, 2, 3} {
		if todos[i].ID != want {
			t.Errorf("todos[%d].ID = %d, want %d", i, todos[i].ID, want)
		}
	}
}

func TestGetTodo(t *testing.T) {
	resetState()
	AddTodo(Todo{Task: "test task"})

	todo := GetTodo(1)
	if todo == nil {
		t.Fatal("expected todo, got nil")
	}
	if todo.Task != "test task" {
		t.Errorf("expected 'test task', got %q", todo.Task)
	}
	if todo.ID != 1 {
		t.Errorf("expected ID 1, got %d", todo.ID)
	}
}

func TestGetTodoNotFound(t *testing.T) {
	resetState()

	todo := GetTodo(999)
	if todo != nil {
		t.Errorf("expected nil for nonexistent ID, got %v", todo)
	}
}

func TestGetAllTodo(t *testing.T) {
	resetState()
	AddTodo(Todo{Task: "first"})
	AddTodo(Todo{Task: "second"})

	all := GetAllTodo()
	if len(all) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(all))
	}
	if all[0].Task != "first" {
		t.Errorf("expected first todo 'first', got %q", all[0].Task)
	}
	if all[1].Task != "second" {
		t.Errorf("expected second todo 'second', got %q", all[1].Task)
	}
}

func TestGetAllTodoEmpty(t *testing.T) {
	resetState()

	all := GetAllTodo()
	if len(all) != 0 {
		t.Errorf("expected empty list, got %d items", len(all))
	}
}

func TestUpdateTodo(t *testing.T) {
	resetState()
	AddTodo(Todo{Task: "original", Completed: false})

	updated, err := UpdateTodo(1, &Todo{Task: "modified", Completed: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Task != "modified" {
		t.Errorf("expected task 'modified', got %q", updated.Task)
	}
	if !updated.Completed {
		t.Error("expected Completed to be true after update")
	}
	if updated.ID != 1 {
		t.Errorf("expected ID to remain 1, got %d", updated.ID)
	}
}

func TestUpdateTodoNotFound(t *testing.T) {
	resetState()

	_, err := UpdateTodo(999, &Todo{Task: "nope"})
	if err == nil {
		t.Fatal("expected error for nonexistent ID, got nil")
	}
	nfe, ok := err.(*NotFoundError)
	if !ok {
		t.Fatalf("expected *NotFoundError, got %T", err)
	}
	if nfe.Message != "Todo not found" {
		t.Errorf("expected 'Todo not found', got %q", nfe.Message)
	}
}

func TestUpdateTodoPartialFields(t *testing.T) {
	resetState()
	AddTodo(Todo{Task: "original", Completed: false})

	updated, err := UpdateTodo(1, &Todo{Completed: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated.Completed {
		t.Error("expected Completed to be true")
	}
}

func TestDeleteTodoByID(t *testing.T) {
	resetState()
	AddTodo(Todo{Task: "to delete"})

	err := DeleteTodoByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 0 {
		t.Errorf("expected 0 todos after delete, got %d", len(todos))
	}
}

func TestDeleteTodoByIDNotFound(t *testing.T) {
	resetState()

	err := DeleteTodoByID(999)
	if err == nil {
		t.Fatal("expected error for nonexistent ID, got nil")
	}
	if _, ok := err.(*NotFoundError); !ok {
		t.Errorf("expected *NotFoundError, got %T", err)
	}
}

func TestDeleteTodoPreservesOthers(t *testing.T) {
	resetState()
	AddTodo(Todo{Task: "keep first"})
	AddTodo(Todo{Task: "delete me"})
	AddTodo(Todo{Task: "keep last"})

	if err := DeleteTodoByID(2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(todos))
	}
	if todos[0].Task != "keep first" {
		t.Errorf("expected first todo 'keep first', got %q", todos[0].Task)
	}
	if todos[1].Task != "keep last" {
		t.Errorf("expected second todo 'keep last', got %q", todos[1].Task)
	}
}

func TestResetClearsState(t *testing.T) {
	resetState()
	AddTodo(Todo{Task: "one"})
	AddTodo(Todo{Task: "two"})

	Reset()

	if len(todos) != 0 {
		t.Errorf("expected 0 todos after Reset, got %d", len(todos))
	}
	if currentID != 0 {
		t.Errorf("expected currentID 0 after Reset, got %d", currentID)
	}

	AddTodo(Todo{Task: "after reset"})
	if todos[0].ID != 1 {
		t.Errorf("expected ID to restart at 1 after Reset, got %d", todos[0].ID)
	}
}
