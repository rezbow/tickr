package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rezbow/tickr/internal/auth"
	"github.com/rezbow/tickr/internal/database"
	"github.com/rezbow/tickr/internal/entities"
	"github.com/rezbow/tickr/internal/events"
	"github.com/rezbow/tickr/internal/payment"
	"github.com/rezbow/tickr/internal/tickets"
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
	ticketService := tickets.NewTicketsService(db, logger)
	paymentService := payment.NewPaymentService(db, logger)
	jwtService := auth.NewJWTService()

	engine := gin.Default()

	// Public routes (no authentication required)
	engine.POST("/auth/login", userService.LoginHandler)
	engine.POST("/auth/refresh", userService.RefreshTokenHandler)
	engine.POST("/users", userService.CreateUserHandler)
	engine.GET("/events", eventsService.GetEventsHandler)
	engine.GET("/events/:id", eventsService.GetEventHandler)
	engine.GET("/events/:id/tickets", ticketService.GetEventTicketsHandler)
	engine.GET("/tickets/:id", ticketService.GetTicket)

	// Protected routes (authentication required)
	protected := engine.Group("/")
	protected.Use(auth.AuthMiddleware(jwtService))
	{
		// Auth routes
		protected.POST("/auth/logout", userService.LogoutHandler)
		protected.GET("/auth/profile", userService.GetProfileHandler)

		// User management (admin only)
		protected.GET("/users", auth.RequireRole("admin"), userService.GetUsersHandler)
		protected.GET("/users/:id", auth.RequireRole("admin"), userService.GetUserHandler)
		protected.DELETE("/users/:id", auth.RequireRole("admin"), userService.DeleteUserHandler)
		protected.PUT("/users/:id", auth.RequireOwnershipOrRole("admin"), userService.UpdateUserHander)

		// Event management (organizers and admins)
		protected.POST("/events", auth.RequireRoles([]string{"organizer", "admin"}), eventsService.CreateEventHandler)
		protected.DELETE("/events/:id", auth.RequireEntityOwnershipOrRole(db, entities.Event{}, "admin"), eventsService.DeleteEventHandler)
		protected.POST("/events/:id/tickets", auth.RequireEntityOwnershipOrRole(db, entities.Event{}, "admin"), ticketService.CreateTicketHandler)

		// Ticket management (organizers and admins)
		protected.DELETE("/tickets/:id", auth.RequireEntityOwnershipOrRole(db, entities.Ticket{}, "admin"), ticketService.DeleteTicket)

		// Payment management (authenticated users)
		protected.POST("/payments", paymentService.BuyTicketHandler)
		protected.GET("/payments/:id", paymentService.GetPaymentHandler)

	}

	engine.Run(":8080")
}
