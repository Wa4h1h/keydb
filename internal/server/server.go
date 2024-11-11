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
	Logger      *zap.Logger
	evaluator   *evaluator.Evaluator
	Port        string
	ReadTimeout time.Duration
	l           net.Listener
}

func (s *Server) Close() error {
	if err := s.l.Close(); err != nil {
		return fmt.Errorf("error occured while trying to close server: %w", err)
	}

	return nil
}

func (s *Server) ListenAndAccept() error {
	var err error

	s.l, err = net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", s.Port))
	if err != nil {
		return fmt.Errorf("failed to start server on port %s: %w", s.Port, err)
	}

	for {
		conn, err := s.l.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Println(err.Error())

				return fmt.Errorf("failed to accept connections: %w", err)
			}

			return nil
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	cmdBytes := make([]byte, 0)

	for {
		if err := conn.SetReadDeadline(time.Now().Add(s.ReadTimeout * time.Second)); err != nil {
			s.Logger.Error(fmt.Sprintf("error while setting read timeout"))
			s.simpleWrite(conn, "unexpected error")

			return
		}

		readBytes := make([]byte, 50)

		n, err := conn.Read(readBytes)
		if err != nil {
			if os.IsTimeout(err) {
				s.simpleWrite(conn, "connection timed out")
			} else if !errors.Is(err, io.EOF) {
				s.simpleWrite(conn, "unexpected error")
			}

			return
		}

		if n > 0 {
			cmdBytes = append(cmdBytes, readBytes[:n]...)

			if slices.Contains(cmdBytes, '\n') {
				var response string
				response, err := s.evaluator.Execute(string(cmdBytes))
				if err != nil {
					response = s.checkError(err)
					s.Logger.Debug(response)
				}
				_, err = conn.Write([]byte(response))
				if err != nil {
					s.Logger.Error(fmt.Sprintf("error while writing to conn:%s", err.Error()))

					return
				}
				cmdBytes = make([]byte, 0)
			}
		}
	}
}

func NewServer(port string,
	logger *zap.Logger, readTimeout time.Duration,
) *Server {
	s := &Server{
		Port:        port,
		Logger:      logger,
		ReadTimeout: readTimeout,
		evaluator: evaluator.NewEvaluator(store.NewMemStore(),
			logger.Sugar().Named("evaluator"),
		),
	}

	return s
}
