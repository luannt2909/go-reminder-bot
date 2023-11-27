package server

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}
type server struct {
	handler Handler
}

func (s server) Stop(ctx context.Context) {
	//TODO implement me
}

func NewServer(handler Handler) Server {
	return &server{handler: handler}
}

func (s server) Start(ctx context.Context) {
	h := s.handler
	g := gin.Default()
	//g.Static("/admin", "./admin/reminder-admin/dist")
	g.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400,
	}))
	r := g.Group("/api/v1")
	webhookG := r.Group("/webhook")
	{
		webhookG.POST("/send", h.SendMessage)
	}
	authG := r.Group("/auth")
	{
		authG.POST("/authenticate", h.Login)
		authG.POST("/logout", AuthenticateUser)
	}

	authGroup := r.Group("", AuthenticateUser)
	reminderG := authGroup.Group("/reminders")
	{
		reminderG.GET("", h.findReminders)
		reminderG.GET("/:id", h.getOneReminder)
		reminderG.POST("", h.createReminder)
		reminderG.PUT("/:id", h.updateReminder)
		reminderG.DELETE("/:id", h.deleteReminder)
	}
	userG := authGroup.Group("/users")
	{
		userG.GET("", h.findUsers)
		userG.GET("/:id", h.getOneUser)
		userG.POST("", h.createUser)
		userG.PUT("/:id", h.updateUser)
		userG.DELETE("/:id", h.deleteUser)
	}

	g.Run()
}
