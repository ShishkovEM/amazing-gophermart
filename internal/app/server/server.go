package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ShishkovEM/amazing-gophermart/internal/app/config"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/controllers"
	"github.com/ShishkovEM/amazing-gophermart/internal/app/storage"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, storage *storage.Storage) *Server {
	serverReadTimeout, _ := time.ParseDuration(cfg.ServerReadTimeout)
	serverWriteTimeout, _ := time.ParseDuration(cfg.ServerWriteTimeout)

	routes := controllers.Routes(storage, []byte(cfg.SecretKey), cfg.CookieLifeTime)
	server := http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      routes,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
	}
	return &Server{
		httpServer: &server,
	}
}

func (a *Server) Run() error {
	addr := a.httpServer.Addr
	log.Printf("Web-server started at http://%s", addr)
	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil {
			log.Printf("Something wrong with server: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
