package server

import (
	"log"
	"net/http"

	"github.com/isaulin-svg/sprint6/internal/handlers"
)

type Server struct {
	HTTP   *http.Server
	logger *log.Logger
}

func NewServer(logger *log.Logger) *Server {
	mux := http.NewServeMux()

	srv := &Server{
		HTTP: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
		logger: logger,
	}

	// Обработчики
	mux.HandleFunc("/", handlers.HandlerIndex)
	mux.HandleFunc("/upload", handlers.HandlerUpload)

	return srv
}
