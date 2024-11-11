package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wa4h1h/memdb/internal/config"
	"github.com/Wa4h1h/memdb/internal/server"
	"github.com/Wa4h1h/memdb/internal/utils"
)

func main() {
	cfg := config.LoadServerConfig()
	logger := utils.NewLogger(cfg.LogLevel)
	s := server.NewServer(cfg.Port,
		logger,
		time.Duration(cfg.TCPReadTimeout)*time.Second,
		cfg.BackOffLimit,
		cfg.TTLBackgroundWorkerInterval)

	go func() {
		if err := s.ListenAndAccept(); err != nil && !errors.Is(err, net.ErrClosed) {
			panic(err.Error())
		}
	}()
	logger.Info(fmt.Sprintf("started listening on port %s", s.Port))

	defer func(ser *server.Server) {
		if err := ser.Close(); err != nil {
			panic(err.Error())
		}

		logger.Info("keydb closed")
	}(s)

	// listen shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
