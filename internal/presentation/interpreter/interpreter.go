package interpreter

import (
	"context"
	"errors"
	"strings"

	"github.com/rdimidov/kvstore/internal/domain"
)

// List of supported command names
const (
	getCommand = "GET"
	setCommand = "SET"
	delCommand = "DEL"
)

// Expected number of arguments for each command
const (
	minArgsLen      = 2
	getArgsLen      = 2
	delArgsLen      = 2
	setArgsLen      = 3
	commandNameIdx  = 0
	commandKeyIdx   = 1
	commandValueIdx = 2
)

// ErrInvalidCmd is returned when the input does not match any supported command format.
var ErrInvalidCmd = errors.New("invalid command")

// application defines the set of operations supported by the business logic layer.
type handler interface {
	Get(ctx context.Context, key domain.Key) (*domain.Entry, error)
	Set(ctx context.Context, key domain.Key, value domain.Value) error
	Delete(ctx context.Context, key domain.Key) error
}

// Interpreter handles parsing raw input strings and executing corresponding application commands.
type Interpreter struct {
	handler handler
}

// New creates a new Interpreter with the given application implementation.
// Returns an error if the application is nil.
func New(handler handler) (*Interpreter, error) {
	if handler == nil {
		return nil, errors.New("handler is nil")
	}
	return &Interpreter{handler: handler}, nil
}

// Execute parses the raw input string and invokes the matching application method.
// Supported formats:
//
//	GET <key>
//	DEL <key>
//	SET <key> <value>
func (i *Interpreter) Execute(ctx context.Context, raw string) (*domain.Entry, error) {
	tokens := strings.Fields(raw)
	if len(tokens) < minArgsLen {
		return nil, ErrInvalidCmd
	}

	key, err := domain.NewKey(tokens[commandKeyIdx])
	if err != nil {
		return nil, err
	}

	switch tokens[commandNameIdx] {
	case getCommand:
		if len(tokens) != getArgsLen {
			return nil, ErrInvalidCmd
		}
		return i.handler.Get(ctx, key)

	case delCommand:
		if len(tokens) != delArgsLen {
			return nil, ErrInvalidCmd
		}
		return nil, i.handler.Delete(ctx, key)

	case setCommand:
		if len(tokens) != setArgsLen {
			return nil, ErrInvalidCmd
		}
		value, err := domain.NewValue(tokens[commandValueIdx])
		if err != nil {
			return nil, err
		}
		return nil, i.handler.Set(ctx, key, value)
	}

	return nil, ErrInvalidCmd
}

type RawInterpreter struct {
	Interpreter
}

func NewRaw(handler handler) (*RawInterpreter, error) {
	if handler == nil {
		return nil, errors.New("application is nil")
	}
	return &RawInterpreter{
		Interpreter: Interpreter{handler: handler},
	}, nil
}

func (r *RawInterpreter) Execute(ctx context.Context, data []byte) []byte {
	raw := strings.TrimSpace(string(data))
	result, err := r.Interpreter.Execute(ctx, raw)
	if err != nil {
		return []byte("ERR " + err.Error() + "\n")
	}

	if result == nil {
		return []byte("OK\n")
	}

	return []byte(result.Value.String() + "\n")
}
