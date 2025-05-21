package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/rdimidov/kvstore/internal/domain"
	"github.com/rdimidov/kvstore/internal/presentation/interpreter"
)

// contextKey is used for context values
// txIDKey stores a transaction ID
type contextKey string

const (
	txIDKey   contextKey = "tx"
	inputMark string     = "> "
)

// app abstracts business logic methods used by CLI
type app interface {
	Get(ctx context.Context, key domain.Key) (*domain.Entry, error)
	Set(ctx context.Context, key domain.Key, value domain.Value) error
	Delete(ctx context.Context, key domain.Key) error
}

// interpr processes raw input and executes commands
type interpr interface {
	Execute(ctx context.Context, raw string) (*domain.Entry, error)
}

// Cli runs the command loop
type Cli struct {
	interpreter interpr
}

// NewCli returns a CLI bound to the given app
func NewCli(app app) (*Cli, error) {
	if app == nil {
		return nil, errors.New("app is nil")
	}
	interp, err := interpreter.New(app)
	if err != nil {
		return nil, err
	}
	return &Cli{interpreter: interp}, nil
}

// Run reads stdin, executes commands, and prints results
func (c *Cli) Run(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	c.prompt()

	for scanner.Scan() {
		input := scanner.Text()
		if input == "" {
			c.prompt()
			continue
		}

		ctx = context.WithValue(ctx, txIDKey, uuid.NewString())
		entry, err := c.interpreter.Execute(ctx, input)

		switch {
		case err != nil:
			fmt.Println("Error:", err)
		case entry != nil:
			fmt.Println(entry.Value)
		}
		c.prompt()
	}
	return scanner.Err()
}

// prompt displays the CLI prompt
func (c *Cli) prompt() {
	fmt.Print(inputMark)
}
