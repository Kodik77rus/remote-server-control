package server

import (
	"context"
	"log"
	"net/http"
	"remote-server-control/internal/handlers"
	"remote-server-control/internal/middleware"
	"time"
)

const _apiPrifix = "/api/v1"

type Server struct {
	config *Config

	server *http.Server
	ctx    context.Context
}

//Server constructor
func New(c *Config, ctx context.Context) *Server {
	return &Server{
		ctx:    ctx,
		config: c,
		server: configirateServer(c),
	}
}

//Start server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.Handle(_apiPrifix+"/remote-execution", http.HandlerFunc(handlers.ExecuteRemoteCommand))

	wrappedMux := middleware.NewLogger(
		middleware.NewResponseHeader(mux, "Content-Type", "application/json"),
	)

	s.setMux(wrappedMux)

	log.Printf("About to listen on 8443. Go to https://127.0.0.1:8443/")

	log.Fatalln(s.server.ListenAndServeTLS("", ""))

	return nil
}

//Server shutdown
func (s *Server) Shutdown() {
	log.Printf("Server stopped")

	ctxShutDownn, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	var err error

	if err = s.server.Shutdown(ctxShutDownn); err != nil {
		log.Fatalln("Server exited properly")
	}

	if err == http.ErrServerClosed {
		err = nil
	}
}

func (s *Server) setMux(wrappredMux *middleware.Logger) {
	s.server.Handler = wrappredMux
}

//Server configurator
func configirateServer(c *Config) *http.Server {
	return &http.Server{
		Addr:         c.port,
		TLSConfig:    c.tlsCfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
}
