package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wa4h1h/memdb/internal/server"
	"github.com/Wa4h1h/memdb/pkg"
)

//nolint:gochecknoglobals
var (
	port        = pkg.GetEnv[string]("8000", false, "PORT")
	logLevel    = pkg.GetEnv[string]("debug", false, "LOG_LEVEL")
	readTimeout = pkg.GetEnv[uint]("30", false, "READ_TIMEOUT")
)

func main() {
	logger := pkg.NewLogger(logLevel)
	s := server.NewServer(port, logger, time.Duration(readTimeout))

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

		logger.Info("server closed")
	}(s)

	// listen shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
}
