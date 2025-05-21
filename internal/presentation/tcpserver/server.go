package tcpserver

import (
	"context"
	"errors"
	"net"
	"time"

	"go.uber.org/zap"
)

const defaultBufferSize int = 4096

type handler interface {
	Execute(context.Context, []byte) []byte
}

type Server struct {
	listener     net.Listener
	handler      handler
	bufferSize   int
	readTimeout  time.Duration
	writeTimeout time.Duration
	logger       *zap.SugaredLogger
}

func New(address string, handler handler, logger *zap.SugaredLogger, options ...Option) (*Server, error) {
	if handler == nil {
		return nil, errors.New("handler is required")
	}

	ln, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	s := &Server{
		listener:   ln,
		handler:    handler,
		logger:     logger,
		bufferSize: defaultBufferSize,
	}

	for _, opt := range options {
		opt(s)
	}
	return s, nil
}

func (s *Server) Start(ctx context.Context) {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					s.logger.Infow("listener already closed", "error", err)
					return
				}
				s.logger.Infow("failed to accept connection", "error", err)
				continue
			}
			go s.handleConnection(ctx, conn)
		}
	}()
	s.logger.Infof("server listening on %v", s.listener.Addr())
	<-ctx.Done()
	s.shutdown()
}

func (s *Server) shutdown() {
	if err := s.listener.Close(); err != nil {
		s.logger.Infow("could not close listener correctly", "error", err)
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		if v := recover(); v != nil {
			s.logger.Errorw("panic in connection handler", "panic", v)
		}
		if err := conn.Close(); err != nil {
			s.logger.Info("could not close connection", "error", err)
		}
	}()

	buf := make([]byte, s.bufferSize)
	for {
		if s.readTimeout != 0 {
			if err := conn.SetReadDeadline(time.Now().Add(s.readTimeout)); err != nil {
				s.logger.Infow("failed to set read timeout", "error", err)
				return
			}
		}

		count, err := conn.Read(buf)
		if err != nil {
			s.logger.Infow("failed to read data", "error", err)
			return
		}

		if count == s.bufferSize {
			s.logger.Info("buffer is full")
			return
		}

		if s.writeTimeout != 0 {
			if err := conn.SetWriteDeadline(time.Now().Add(s.writeTimeout)); err != nil {
				s.logger.Infow("failed to set write timeout", "error", err)
				return
			}
		}

		resp := s.handler.Execute(ctx, buf[:count])
		if _, err := conn.Write(resp); err != nil {
			s.logger.Infow("failed to write data", "error", err)
			return
		}
	}
}
