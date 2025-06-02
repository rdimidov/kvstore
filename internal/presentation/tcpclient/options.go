package tcpclient

import (
	"time"
)

type Option func(*Client)

func WithTimeout(timeout time.Duration) Option {
	return func(s *Client) {
		s.timeout = timeout
	}
}

func WithBufferSize(size int) Option {
	return func(s *Client) {
		s.bufferSize = size
	}
}
