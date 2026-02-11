package main

import (
	"log"

	"github.com/AndriyKalashnykov/go-todo-web/handlers"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {

	// Create a new instance of Echo
	e := echo.New()

	// Middleware
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.Use(middleware.Recover())

	// Routes
	e.POST("/create", handlers.CreateTodo)
	e.GET("/get/:id", handlers.GetTodo)
	e.GET("/all", handlers.Todos)
	e.DELETE("/delete/:id", handlers.DeleteTodo)
	e.PATCH("/update/:id", handlers.UpdateTodo)
	// Start the server
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
