package vak

import (
	"context"
	"fmt"
	"sync"
)

// Registry manages command registration and execution.
type Registry struct {
	mu       sync.RWMutex
	commands map[string]*Command
	parser   *Parser
	history  *History
	undoMgr  *UndoManager
}

// NewRegistry creates a new command registry.
func NewRegistry() *Registry {
	return &Registry{
		commands: make(map[string]*Command),
		parser:   NewParser(),
		history:  NewHistory(100),
		undoMgr:  NewUndoManager(50),
	}
}

// Register registers a command.
func (r *Registry) Register(cmd Command) error {
	if err := cmd.Validate(); err != nil {
		return fmt.Errorf("invalid command: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.commands[cmd.Name]; exists {
		return fmt.Errorf("command already registered: %s", cmd.Name)
	}

	r.commands[cmd.Name] = &cmd
	return nil
}

// Unregister removes a command.
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.commands[name]; !exists {
		return fmt.Errorf("command not found: %s", name)
	}

	delete(r.commands, name)
	return nil
}

// Get retrieves a command by name.
func (r *Registry) Get(name string) (*Command, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cmd, exists := r.commands[name]
	return cmd, exists
}

// List returns all registered commands.
func (r *Registry) List() []Command {
	r.mu.RLock()
	defer r.mu.RUnlock()

	commands := make([]Command, 0, len(r.commands))
	for _, cmd := range r.commands {
		commands = append(commands, *cmd)
	}
	return commands
}

// Execute parses and executes a command string.
func (r *Registry) Execute(ctx context.Context, input string) (Result, error) {
	// Parse input
	parsed, err := r.parser.Parse(input)
	if err != nil {
		return ErrorResult(fmt.Sprintf("parse error: %v", err)), err
	}

	// Get command
	r.mu.RLock()
	cmd, exists := r.commands[parsed.Name]
	r.mu.RUnlock()

	if !exists {
		err := fmt.Errorf("command not found: %s", parsed.Name)
		return ErrorResult(err.Error()), err
	}

	// Execute command
	result, err := cmd.Execute(ctx, parsed.Args)

	// Add to history
	r.history.Add(HistoryEntry{
		Input:   input,
		Command: parsed.Name,
		Success: result.Success(),
	})

	// Add to undo stack if undoable
	if err == nil && result.Success() && cmd.CanUndo() {
		r.undoMgr.Push(UndoEntry{
			Command: cmd,
			Context: ctx,
		})
	}

	return result, err
}

// Complete provides autocompletion suggestions.
func (r *Registry) Complete(prefix string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	suggestions := make([]string, 0)

	// Try to parse partial input
	parsed, _ := r.parser.Parse(prefix)

	if parsed.Name == "" || len(parsed.Args) == 0 {
		// Complete command names
		for name := range r.commands {
			if len(prefix) == 0 || startsWith(name, prefix) {
				suggestions = append(suggestions, name)
			}
		}
	} else {
		// Complete arguments using command's Complete function
		if cmd, exists := r.commands[parsed.Name]; exists && cmd.Complete != nil {
			lastArg := ""
			if len(parsed.Args) > 0 {
				lastArg = parsed.Args[len(parsed.Args)-1]
			}
			suggestions = cmd.Complete(lastArg)
		}
	}

	return suggestions
}

// History returns the command history.
func (r *Registry) History() *History {
	return r.history
}

// UndoManager returns the undo manager.
func (r *Registry) UndoManager() *UndoManager {
	return r.undoMgr
}

// Undo undoes the last command.
func (r *Registry) Undo(ctx context.Context) error {
	return r.undoMgr.Undo(ctx)
}

// Redo redoes the last undone command.
func (r *Registry) Redo(ctx context.Context) error {
	return r.undoMgr.Redo(ctx)
}

// startsWith checks if a string starts with a prefix (case-insensitive).
func startsWith(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if toLower(s[i]) != toLower(prefix[i]) {
			return false
		}
	}
	return true
}

// toLower converts a byte to lowercase.
func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}
