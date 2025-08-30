package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/handler"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/logger"
	"github.com/FIAP-SOAT-G20/tc4-order-service/internal/infrastructure/route"
	"github.com/gin-gonic/gin/binding"

	"github.com/go-playground/validator/v10"
)

type Server struct {
	router *route.Router
	config *config.Config
	logger *logger.Logger
}

func NewServer(cfg *config.Config, logger *logger.Logger, handlers *route.Handlers) *Server {
	router := route.NewRouter(logger, cfg)

	RegisterCustomValidation()
	router.RegisterRoutes(handlers)

	return &Server{
		router: router,
		config: cfg,
		logger: logger,
	}
}

func (s *Server) Start() error {
	server := &http.Server{
		Addr:         ":" + s.config.ServerPort,
		Handler:      s.router.Engine(),
		ReadTimeout:  s.config.ServerReadTimeout,
		WriteTimeout: s.config.ServerWriteTimeout,
		IdleTimeout:  s.config.ServerIdleTimeout,
	}

	var serverErr error

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		s.logger.Info("server is running", "port", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("server failed to start", "error", err)
			serverErr = err
			sigChannel <- syscall.SIGINT
		}
	}()

	<-sigChannel

	gracefullyShutdown(server, s)

	return serverErr
}

func gracefullyShutdown(server *http.Server, s *Server) {
	s.logger.Info("shutting down server...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), s.config.ServerGracefulShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctxTimeout); err != nil {
		s.logger.Error("server failed to shutdown", "error", err)
		panic(err)
	}

	s.logger.Info("server exited gracefully")
	fmt.Println("Bye! 👋")
}

func RegisterCustomValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("order_status_exists", handler.OrderStatusValidator)
		if err != nil {
			panic(err)
		}
	}
}
