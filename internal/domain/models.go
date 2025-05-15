package domain

import "github.com/rdimidov/kvstore/internal/domain/validator"

type Key string

func NewKey(k string) (Key, error) {
	if !validator.IsValidString(k) {
		return "", ErrKeyIsNotValid
	}
	return Key(k), nil
}

func (k Key) String() string {
	return string(k)
}

type Value string

func NewValue(v string) (Value, error) {
	if !validator.IsValidString(v) {
		return "", ErrValueIsNotValid
	}
	return Value(v), nil
}

func (v Value) String() string {
	return string(v)
}

type Entry struct {
	Key   Key
	Value Value
}

func NewEntryFromKV(k Key, v Value) Entry {
	return Entry{k, v}
}
