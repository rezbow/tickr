package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rezbow/tickr/internal/database"
	"github.com/rezbow/tickr/internal/events"
	"github.com/rezbow/tickr/internal/users"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		panic("missing DB_URL env")
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db := database.SetupDatabase(dsn)
	userService := users.NewUserService(db, logger)
	eventsService := events.NewEventsService(db, logger)

	engine := gin.Default()
	// users
	engine.POST("/users", userService.CreateUserHandler)
	engine.DELETE("/users/:id", userService.DeleteUserHandler)
	engine.GET("/users/:id", userService.GetUserHandler)
	engine.PUT("/users/:id", userService.UpdateUserHander)
	engine.GET("/users", userService.GetUsersHandler)
	// events
	engine.POST("/events", eventsService.CreateEventHandler)
	engine.GET("/events/:id", eventsService.GetEventHandler)
	engine.GET("/events", eventsService.GetEventsHandler)
	engine.DELETE("/events/:id", eventsService.DeleteEventHandler)
	engine.PUT("/events/:id", eventsService.UpdateEventHandler)

	engine.GET("/events/:id/tickets", eventsService.GetEventTicketsHandler)

	engine.Run(":8080")
}
