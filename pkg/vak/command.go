package vak

import (
	"context"
	"fmt"
)

// Command represents an executable command.
type Command struct {
	Name     string
	Summary  string
	Usage    string
	Flags    []Flag
	Examples []string
	Execute  ExecuteFunc
	Undo     UndoFunc
	Complete CompleteFunc
}

// ExecuteFunc is the function that executes a command.
type ExecuteFunc func(ctx context.Context, args []string) (Result, error)

// UndoFunc is an optional function to undo a command.
type UndoFunc func(ctx context.Context) error

// CompleteFunc provides autocompletion suggestions.
type CompleteFunc func(prefix string) []string

// Flag represents a command flag.
type Flag struct {
	Name        string
	Short       string
	Type        FlagType
	Description string
	Required    bool
	Default     any
}

// FlagType represents the type of a flag.
type FlagType int

const (
	FlagTypeString FlagType = iota
	FlagTypeInt
	FlagTypeBool
	FlagTypeFloat
)

// String returns the string representation of flag type.
func (ft FlagType) String() string {
	switch ft {
	case FlagTypeString:
		return "string"
	case FlagTypeInt:
		return "int"
	case FlagTypeBool:
		return "bool"
	case FlagTypeFloat:
		return "float"
	default:
		return "unknown"
	}
}

// Result represents the result of command execution.
type Result interface {
	Success() bool
	Message() string
	Data() any
}

// BaseResult provides a basic result implementation.
type BaseResult struct {
	SuccessFlag bool
	Msg         string
	ResultData  any
}

// Success returns whether the command succeeded.
func (r *BaseResult) Success() bool {
	return r.SuccessFlag
}

// Message returns the result message.
func (r *BaseResult) Message() string {
	return r.Msg
}

// Data returns the result data.
func (r *BaseResult) Data() any {
	return r.ResultData
}

// SuccessResult creates a success result.
func SuccessResult(message string) Result {
	return &BaseResult{
		SuccessFlag: true,
		Msg:         message,
	}
}

// SuccessWithData creates a success result with data.
func SuccessWithData(message string, data any) Result {
	return &BaseResult{
		SuccessFlag: true,
		Msg:         message,
		ResultData:  data,
	}
}

// ErrorResult creates an error result.
func ErrorResult(message string) Result {
	return &BaseResult{
		SuccessFlag: false,
		Msg:         message,
	}
}

// Validate validates the command definition.
func (c *Command) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("command name cannot be empty")
	}
	if c.Execute == nil {
		return fmt.Errorf("command must have an Execute function")
	}
	return nil
}

// Help generates help text for the command.
func (c *Command) Help() string {
	help := fmt.Sprintf("Command: %s\n", c.Name)
	help += fmt.Sprintf("Summary: %s\n", c.Summary)
	help += fmt.Sprintf("Usage: %s\n", c.Usage)

	if len(c.Flags) > 0 {
		help += "\nFlags:\n"
		for _, flag := range c.Flags {
			flagStr := fmt.Sprintf("  --%s", flag.Name)
			if flag.Short != "" {
				flagStr += fmt.Sprintf(", -%s", flag.Short)
			}
			flagStr += fmt.Sprintf(" (%s)", flag.Type.String())
			if flag.Required {
				flagStr += " [required]"
			}
			help += flagStr + "\n"
			if flag.Description != "" {
				help += fmt.Sprintf("    %s\n", flag.Description)
			}
		}
	}

	if len(c.Examples) > 0 {
		help += "\nExamples:\n"
		for _, example := range c.Examples {
			help += fmt.Sprintf("  %s\n", example)
		}
	}

	return help
}

// CanUndo returns whether the command can be undone.
func (c *Command) CanUndo() bool {
	return c.Undo != nil
}
