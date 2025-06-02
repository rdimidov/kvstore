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
type application interface {
	Get(ctx context.Context, key domain.Key) (*domain.Entry, error)
	Set(ctx context.Context, key domain.Key, value domain.Value) error
	Delete(ctx context.Context, key domain.Key) error
}

// Interpreter handles parsing raw input strings and executing corresponding application commands.
type Interpreter struct {
	app application
}

// New creates a new Interpreter with the given application implementation.
// Returns an error if the application is nil.
func New(app application) (*Interpreter, error) {
	if app == nil {
		return nil, errors.New("application is nil")
	}
	return &Interpreter{app: app}, nil
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

	cmd := strings.ToUpper(tokens[commandNameIdx])
	key, err := domain.NewKey(tokens[commandKeyIdx])
	if err != nil {
		return nil, err
	}

	switch cmd {
	case getCommand:
		if len(tokens) != getArgsLen {
			return nil, ErrInvalidCmd
		}
		return i.app.Get(ctx, key)

	case delCommand:
		if len(tokens) != delArgsLen {
			return nil, ErrInvalidCmd
		}
		return nil, i.app.Delete(ctx, key)

	case setCommand:
		if len(tokens) != setArgsLen {
			return nil, ErrInvalidCmd
		}
		value, err := domain.NewValue(tokens[commandValueIdx])
		if err != nil {
			return nil, err
		}
		return nil, i.app.Set(ctx, key, value)
	}

	return nil, ErrInvalidCmd
}

type RawInterpreter struct {
	Interpreter
}

func NewRaw(app application) (*RawInterpreter, error) {
	if app == nil {
		return nil, errors.New("application is nil")
	}
	return &RawInterpreter{
		Interpreter: Interpreter{app: app},
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
