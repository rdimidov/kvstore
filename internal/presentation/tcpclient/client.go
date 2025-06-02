package tcpclient

import (
	"errors"
	"net"
	"time"
)

type Client struct {
	conn       net.Conn
	timeout    time.Duration
	bufferSize int
}

func New(serverAddr string, options ...Option) (*Client, error) {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return nil, err
	}

	c := &Client{conn: conn}
	for _, opt := range options {
		opt(c)
	}

	if c.timeout != 0 {
		if err := c.conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Client) Send(message []byte) ([]byte, error) {
	if _, err := c.conn.Write(message); err != nil {
		return nil, err
	}
	resp := make([]byte, c.bufferSize)
	count, err := c.conn.Read(resp)
	if err != nil {
		return nil, err
	}
	if count == c.bufferSize {
		return nil, errors.New("received response is too big")
	}

	return resp[:count], nil
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
