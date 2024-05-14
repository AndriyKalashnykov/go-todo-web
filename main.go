package main

import (
	"github.com/AndriyKalashnykov/go-todo-web/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// Create a new instance of Echo
	e := echo.New()

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/create", handlers.CreateTodo)
	e.GET("/get/:id", handlers.GetTodo)
	e.GET("/all", handlers.Todos)
	e.DELETE("/delete/:id", handlers.DeleteTodo)
	e.PATCH("/update/:id", handlers.UpdateTodo)
	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
