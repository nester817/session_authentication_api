package server

import (
	"net/http"
	"time"

	router "github.com/nester817/session_authentication_api.git/pkg/api/gin_router"
)

type Server struct {
	httpServer http.Server
}

func NewServer(addr string, h *router.Handler) *Server {
	return &Server{
		httpServer: http.Server{
			Addr:    addr,
			Handler: h.InitRouter(),

			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,

			IdleTimeout: 120 * time.Second,
		},
	}
}

func (s *Server) Run() {

}
