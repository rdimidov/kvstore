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
	if len(tokens) < 2 {
		return nil, ErrInvalidCmd
	}

	cmd := strings.ToUpper(tokens[0])
	key, err := domain.NewKey(tokens[1])
	if err != nil {
		return nil, err
	}

	switch cmd {
	case getCommand, delCommand:
		if len(tokens) != 2 {
			return nil, ErrInvalidCmd
		}
		return &Command{Cmd: cmd, Key: key}, nil
	case setCommand:
		if len(tokens) != 3 {
			return nil, ErrInvalidCmd
		}
		value, err := domain.NewValue(tokens[2])
		if err != nil {
			return nil, err
		}
		return &Command{Cmd: cmd, Key: key, Value: &value}, nil
	}
	return nil, ErrInvalidCmd
}
