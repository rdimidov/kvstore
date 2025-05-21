package tcpclient

import (
	"net"
	"testing"
	"time"
)

func startTestTCPServer(t *testing.T, response string) (addr string, stop func()) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		buf := make([]byte, 1024)
		_, _ = conn.Read(buf)
		_, _ = conn.Write([]byte(response))
	}()

	return ln.Addr().String(), func() {
		_ = ln.Close()
	}
}

func TestClient_SendSuccess(t *testing.T) {
	t.Parallel()
	expectedResp := "hello client"
	addr, stop := startTestTCPServer(t, expectedResp)
	defer stop()

	client, err := New(addr,
		WithTimeout(2*time.Second),
		WithBufferSize(1024),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	resp, err := client.Send([]byte("hello server"))
	if err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	if string(resp) != expectedResp {
		t.Errorf("expected %q, got %q", expectedResp, string(resp))
	}
}

func TestClient_SendResponseTooBig(t *testing.T) {
	t.Parallel()
	// simulate response longer than buffer size
	expectedResp := "response-too-long"
	addr, stop := startTestTCPServer(t, expectedResp)
	defer stop()

	client, err := New(addr,
		WithTimeout(2*time.Second),
		WithBufferSize(5),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	_, err = client.Send([]byte("hello"))
	if err == nil || err.Error() != "received response is too big" {
		t.Errorf("expected 'received response is too big' error, got %v", err)
	}
}

func TestClient_Timeout(t *testing.T) {
	t.Parallel()
	// start a server that doesn't respond
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, _ := ln.Accept()
		defer conn.Close()
		time.Sleep(5 * time.Second)
	}()

	client, err := New(ln.Addr().String(),
		WithTimeout(1*time.Second),
		WithBufferSize(1024),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	_, err = client.Send([]byte("hello"))
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}
