package tcpserver

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func startTestServer(t *testing.T, handler handler) (addr string, cancel context.CancelFunc) {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr = ln.Addr().String()
	_ = ln.Close()

	ctx, cancel := context.WithCancel(context.Background())
	logger := zap.NewNop().Sugar()

	server, err := New(addr, handler, logger)
	require.NoError(t, err)

	go server.Start(ctx)
	time.Sleep(100 * time.Millisecond) // Give time to start

	return addr, cancel
}

func TestServer_HandleRequest(t *testing.T) {
	// Create and set expectations on the mock handler
	mockHandler := newMockhandler(t)
	mockHandler.
		EXPECT().
		Execute(mock.Anything, []byte("test-input")).
		Return([]byte("echo-response")).
		Once()

	addr, cancel := startTestServer(t, mockHandler)
	defer cancel()

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer conn.Close()

	_, err = conn.Write([]byte("test-input"))
	require.NoError(t, err)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	require.NoError(t, err)

	require.Equal(t, "echo-response", string(buf[:n]))
}

func TestServer_GracefulShutdown(t *testing.T) {
	mockHandler := newMockhandler(t)

	addr, cancel := startTestServer(t, mockHandler)

	// Connect to the server (this should succeed)
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer conn.Close()

	// Trigger shutdown
	cancel()
	time.Sleep(100 * time.Millisecond) // Let the shutdown complete

	// Try to connect again (should fail)
	_, err = net.Dial("tcp", addr)
	require.Error(t, err, "expected connection to fail after shutdown")
}
