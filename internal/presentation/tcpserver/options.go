package tcpserver

import (
	"time"
)

type Option func(*Server)

func WithBufferSize(size int) Option {
	return func(s *Server) {
		s.bufferSize = size
	}
}

func WithTimeouts(read, write time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = read
		s.writeTimeout = write
	}
}
