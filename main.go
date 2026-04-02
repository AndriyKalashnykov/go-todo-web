package main

import (
	"log"

	"github.com/AndriyKalashnykov/go-todo-web/handlers"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func NewApp() *echo.Echo {
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

	return e
}

func main() {
	if err := NewApp().Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
