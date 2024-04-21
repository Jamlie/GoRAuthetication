package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload"

	"server/internal/database"
)

type Server struct {
	port  int
	db    database.Service
	store *sessions.CookieStore
	app   *http.ServeMux
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:  port,
		db:    database.New(),
		store: sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET"))),
		app:   http.NewServeMux(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) Use(pattern string, handler http.Handler) {
	s.app.Handle(pattern, handler)
}
