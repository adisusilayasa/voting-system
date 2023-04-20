package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 10
)

type Server struct {
	gin    *gin.Engine
	logger *zap.Logger
}

func NewServer() *Server {
	logger, _ := zap.NewProduction()
	return &Server{
		gin:    gin.New(),
		logger: logger,
	}
}

func (s *Server) RunServer() error {
	gin.DefaultWriter = colorable.NewColorableStderr()

	gin.SetMode(gin.ReleaseMode)
	server := &http.Server{
		Handler:      s.gin,
		Addr:         ":8080",
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		// MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		s.logger.Sugar().Infof("Server is listening on PORT: %s", "8080")
		if err := server.ListenAndServe(); err != nil {
			s.logger.Sugar().Fatal("Error starting Server: ", err)
		}
	}()

	if err := s.MapHandlers(s.gin); err != nil {
		log.Println(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return server.Shutdown(ctx)
}
