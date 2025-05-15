package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rdimidov/kvstore/internal/application/services"
)

type contextKey string

const txIDKey contextKey = "tx"

type Cli struct {
	parser  parser
	handler handler
}

func NewCli(h handler, p parser) *Cli {
	if p == nil {
		p = services.Parser{} // JIT DI
	}
	return &Cli{
		parser:  p,
		handler: h,
	}
}

func (c *Cli) Run(ctx context.Context) error {
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for reader.Scan() {
		input := reader.Text()

		if input == "" {
			fmt.Print("> ")
			continue
		}

		// Create a new context with unique transaction ID
		txCtx := context.WithValue(ctx, txIDKey, uuid.NewString())

		cmd, err := c.parser.Parse(input)
		if err != nil {
			fmt.Println("Parse error:", err)
			fmt.Print("> ")
			continue
		}

		entry, err := c.handler.Handle(txCtx, cmd)
		if err != nil {
			fmt.Println("Error:", err)
		} else if entry != nil {
			fmt.Printf("%s = %s\n", entry.Key, entry.Value)
		}
		fmt.Print("> ")
	}
	return reader.Err()
}
