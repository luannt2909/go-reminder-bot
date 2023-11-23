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
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		MaxAge:           86400,
	}))
	{
		r.GET("/tasks", h.findTask)
		r.GET("/tasks/:id", h.getOneTask)
		r.POST("/tasks", h.createTask)
		r.PUT("/tasks/:id", h.updateTask)
		r.DELETE("/tasks/:id", h.deleteTask)
	}
	r.Run(":3001")
}
