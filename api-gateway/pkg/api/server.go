package api

import (
	"api-gateway/pkg/config"

	"github.com/labstack/echo/v4"
)

type Server struct {
	port   string
	engine *echo.Echo
}

// NewServerHTTP creates a new server with given handler functions
func NewServerHTTP(cfg config.Config) *Server {

	engine := echo.New()

	return &Server{
		engine: engine,
		port:   cfg.ApiPort,
	}
}

func (c *Server) Start() {
	c.engine.Start((":" + c.port))
}
