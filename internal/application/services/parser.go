package services

import (
	"strings"

	"github.com/rdimidov/kvstore/internal/domain"
)

// List of supported command names
const (
	getCommand = "GET"
	setCommand = "SET"
	delCommand = "DEL"
)

// Expected number of arguments for commands
const (
	minArgsLen      = 2
	getDelArgsLen   = 2
	setArgsLen      = 3
	commandNameIdx  = 0
	commandKeyIdx   = 1
	commandValueIdx = 2
)

// Command represents a parsed command structure extracted from a raw input line.
// Value is only set for SET commands.
type Command struct {
	Cmd   string
	Key   domain.Key
	Value *domain.Value
}

// Parser is responsible for converting raw strings into structured Command objects.
type Parser struct{}

// Parse analyzes the input string and builds a Command object.
// Supported commands are: GET <key>, SET <key> <value>, DEL <key>.
func (Parser) Parse(raw string) (*Command, error) {
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
	case getCommand, delCommand:
		if len(tokens) != getDelArgsLen {
			return nil, ErrInvalidCmd
		}
		return &Command{Cmd: cmd, Key: key}, nil
	case setCommand:
		if len(tokens) != setArgsLen {
			return nil, ErrInvalidCmd
		}
		value, err := domain.NewValue(tokens[commandValueIdx])
		if err != nil {
			return nil, err
		}
		return &Command{Cmd: cmd, Key: key, Value: &value}, nil
	}
	return nil, ErrInvalidCmd
}
