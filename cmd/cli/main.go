package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/rdimidov/kvstore/internal/presentation/tcpclient"
	"go.uber.org/zap"
)

const (
	defaultTimeout    = 10 * time.Minute
	defaultBufferSize = 4096
)

func main() {
	// Define and parse required CLI flag for server address
	addrFlag := flag.String("address", "", "server address (required), e.g. 127.0.0.1:3223")
	timeoutFlag := flag.Duration("timeout", defaultTimeout, "timeout for server connection, e.g. 5s, 1m")
	bufSizeFlag := flag.Int("buf", defaultBufferSize, "buffer size in bytes, e.g 1024, 4096")
	flag.Parse()

	// Exit if address is not provided
	if *addrFlag == "" {
		fmt.Fprintln(os.Stderr, "error: --address flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Initialize structured logger (Zap)
	loggerCore, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	logger := loggerCore.Sugar()
	defer logger.Sync() //nolint: all

	// Create TCP client with timeout
	client, err := tcpclient.New(*addrFlag,
		tcpclient.WithTimeout(*timeoutFlag),
		tcpclient.WithBufferSize(*bufSizeFlag),
	)
	if err != nil {
		logger.Fatalw("could not connect to server", "error", err)
	}
	defer client.Close()

	// Start REPL loop for user input
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ") // Prompt
		input, err := reader.ReadString('\n')
		if err != nil {
			// Exit cleanly on EOF (e.g. Ctrl+D)
			if err == io.EOF {
				logger.Infow("received EOF, exiting")
				return
			}
			logger.Errorw("could not read input", "error", err)
			continue
		}

		// Trim newline and skip empty input
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Send user input to server and print response
		resp, err := client.Send([]byte(input))
		if err != nil {
			logger.Fatalw("could not send message", "error", err)
		}

		fmt.Println(string(resp))
	}
}
