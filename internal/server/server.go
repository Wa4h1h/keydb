package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"slices"
	"time"

	"github.com/Wa4h1h/memdb/internal/evaluator"

	"github.com/Wa4h1h/memdb/internal/store"
	"go.uber.org/zap"
)

type Server struct {
	Logger           *zap.SugaredLogger
	evaluator        *evaluator.Evaluator
	l                net.Listener
	Port             string
	ReadTimeout      time.Duration
	ReadBackOffLimit int
}

func NewServer(port string,
	logger *zap.Logger, readTimeout time.Duration,
	readBackOffLimit int,
	ttlBackgroundWorkerInterval int,
) *Server {
	s := &Server{
		Port:             port,
		Logger:           logger.Sugar().Named("server"),
		ReadTimeout:      readTimeout,
		ReadBackOffLimit: readBackOffLimit,
		evaluator: evaluator.NewEvaluator(store.NewMemStore(ttlBackgroundWorkerInterval),
			logger.Sugar().Named("evaluator"),
		),
	}

	return s
}

func (s *Server) Close() error {
	if err := s.l.Close(); err != nil {
		return fmt.Errorf("error: close server: %w", err)
	}

	return nil
}

func (s *Server) ListenAndAccept() error {
	var err error

	s.l, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", s.Port))
	if err != nil {
		return fmt.Errorf("error: listen on port %s: %w", s.Port, err)
	}

	for {
		conn, err := s.l.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Println(err.Error())

				continue
			}

			return fmt.Errorf("error: accept connection: %w", err)
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			s.Logger.Errorf("error: close remote connection %s: %s",
				conn.RemoteAddr().String(), err)
		}
	}()

	cmdBytes := make([]byte, 0)
	tries := 0

	for {
		if err := conn.SetReadDeadline(time.Now().Add(s.ReadTimeout)); err != nil {
			s.Logger.Errorf("error: set read timeout: %s", err)

			tries++

			if tries >= s.ReadBackOffLimit {
				s.simpleWrite(conn, "unexpected error")

				return
			}

			continue
		}

		readBytes := make([]byte, 100)

		n, err := conn.Read(readBytes)
		if err != nil {
			tries++

			switch {
			case os.IsTimeout(err):
				s.simpleWrite(conn, "connection timed out")
			case !errors.Is(err, io.EOF):
				if tries < s.ReadBackOffLimit {
					continue
				}

				s.Logger.Errorf("error: read bytes: %s", err)
				s.simpleWrite(conn, "unexpected error")

				return
			}

			return
		}

		if n > 0 {
			cmdBytes = append(cmdBytes, readBytes[:n]...)

			if slices.Contains(cmdBytes, '\n') {
				var response string

				response, err := s.evaluator.Evaluate(string(cmdBytes))
				if err != nil {
					response = s.checkError(err)

					s.Logger.Debug(response)
				}

				_, err = conn.Write([]byte(response))
				if err != nil {
					s.Logger.Errorf("error: write to conn: %s", err)

					return
				}

				cmdBytes = cmdBytes[:0]
			}
		}
	}
}
