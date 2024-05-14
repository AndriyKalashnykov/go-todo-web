package handlers

import (
	"github.com/AndriyKalashnykov/go-todo-web/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func CreateTodo(c echo.Context) error {
	todo := new(models.Todo)
	if err := c.Bind(todo); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request"})
	}
	models.AddTodo(*todo)
	return c.JSON(http.StatusCreated, todo)
}
func Todos(c echo.Context) error {
	return c.JSON(http.StatusOK, models.GetAllTodo())

}
func GetTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	todo := models.GetTodo(id)
	if todo == nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: "Invalid id"})
	}
	return c.JSON(http.StatusOK, todo)
}
func UpdateTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	todo := models.GetTodo(id)
	if todo == nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: "Invalid id"})
	}
	update := new(models.Todo)
	if err := c.Bind(update); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request"})
	}
	update.Task = todo.Task
	update.Completed = todo.Completed
	return c.JSON(http.StatusOK, todo)
}
func DeleteTodo(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := models.DeleteTodoByID(id); err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: "Todo not found"})
	}
	return c.NoContent(http.StatusNoContent)
}
