package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	server *http.Server
}

func New(port string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHealthCheck)

	return &Server{
		server: &http.Server{
			Addr:    ":" + port,
			Handler: mux,
		},
	}
}

func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Error en servidor HTTP: %v", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Error cerrando servidor HTTP: %v", err)
	}
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bot está funcionando!")
}
