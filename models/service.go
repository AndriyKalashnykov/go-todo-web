package models

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

var todos []*Todo
var currentID int

func AddTodo(todo Todo) {
	currentID++
	todo.ID = currentID
	todos = append(todos, &todo)

}
func UpdateTodo(id int, updated *Todo) (*Todo, error) {
	for _, todo := range todos {
		if todo.ID == id {
			todo.Task = updated.Task
			todo.Completed = updated.Completed
			return todo, nil
		}
	}
	return nil, &NotFoundError{Message: "Todo not found"}
}
func GetTodo(id int) *Todo {
	for _, todo := range todos {
		if todo.ID == id {
			return todo
		}
	}
	return nil
}

func GetAllTodo() []*Todo {
	return todos
}
func DeleteTodoByID(id int) error {
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return &NotFoundError{Message: "Todo not found"}
}
