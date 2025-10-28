package vak

import (
	"fmt"
	"strconv"
)

// FlagSet manages command flags.
type FlagSet struct {
	values map[string]any
}

// NewFlagSet creates a new flag set.
func NewFlagSet() *FlagSet {
	return &FlagSet{
		values: make(map[string]any),
	}
}

// Set sets a flag value.
func (fs *FlagSet) Set(name string, value any) {
	fs.values[name] = value
}

// Get retrieves a flag value.
func (fs *FlagSet) Get(name string) (any, bool) {
	val, exists := fs.values[name]
	return val, exists
}

// String retrieves a string flag.
func (fs *FlagSet) String(name string) (string, error) {
	val, exists := fs.values[name]
	if !exists {
		return "", fmt.Errorf("flag not found: %s", name)
	}

	str, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("flag %s is not a string", name)
	}
	return str, nil
}

// Int retrieves an int flag.
func (fs *FlagSet) Int(name string) (int, error) {
	val, exists := fs.values[name]
	if !exists {
		return 0, fmt.Errorf("flag not found: %s", name)
	}

	// Try direct int
	if i, ok := val.(int); ok {
		return i, nil
	}

	// Try parsing string
	if str, ok := val.(string); ok {
		return strconv.Atoi(str)
	}

	return 0, fmt.Errorf("flag %s is not an int", name)
}

// Bool retrieves a bool flag.
func (fs *FlagSet) Bool(name string) (bool, error) {
	val, exists := fs.values[name]
	if !exists {
		return false, fmt.Errorf("flag not found: %s", name)
	}

	// Try direct bool
	if b, ok := val.(bool); ok {
		return b, nil
	}

	// Try parsing string
	if str, ok := val.(string); ok {
		return strconv.ParseBool(str)
	}

	return false, fmt.Errorf("flag %s is not a bool", name)
}

// Float retrieves a float64 flag.
func (fs *FlagSet) Float(name string) (float64, error) {
	val, exists := fs.values[name]
	if !exists {
		return 0, fmt.Errorf("flag not found: %s", name)
	}

	// Try direct float
	if f, ok := val.(float64); ok {
		return f, nil
	}

	// Try parsing string
	if str, ok := val.(string); ok {
		return strconv.ParseFloat(str, 64)
	}

	return 0, fmt.Errorf("flag %s is not a float", name)
}

// Has checks if a flag exists.
func (fs *FlagSet) Has(name string) bool {
	_, exists := fs.values[name]
	return exists
}
