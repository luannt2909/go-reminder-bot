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
	g.Static("/admin", "./admin/reminder-admin/dist")
	g.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400,
	}))
	r := g.Group("/api/v1")
	{
		r.GET("/reminders", h.findReminders)
		r.GET("/reminders/:id", h.getOneReminder)
		r.POST("/reminders", h.createReminder)
		r.PUT("/reminders/:id", h.updateReminder)
		r.DELETE("/reminders/:id", h.deleteReminder)
	}
	g.Run(":2909")
}
