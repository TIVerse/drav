package vak

import (
	"context"
	"fmt"
	"sync"
)

// UndoManager manages undo/redo operations.
type UndoManager struct {
	mu        sync.RWMutex
	undoStack []UndoEntry
	redoStack []UndoEntry
	maxSize   int
}

// UndoEntry represents an undoable command.
type UndoEntry struct {
	Command *Command
	Context context.Context
}

// NewUndoManager creates a new undo manager.
func NewUndoManager(maxSize int) *UndoManager {
	return &UndoManager{
		undoStack: make([]UndoEntry, 0, maxSize),
		redoStack: make([]UndoEntry, 0, maxSize),
		maxSize:   maxSize,
	}
}

// Push pushes a command onto the undo stack.
func (um *UndoManager) Push(entry UndoEntry) {
	um.mu.Lock()
	defer um.mu.Unlock()

	if len(um.undoStack) >= um.maxSize {
		// Remove oldest entry
		um.undoStack = um.undoStack[1:]
	}

	um.undoStack = append(um.undoStack, entry)

	// Clear redo stack when a new action is performed
	um.redoStack = make([]UndoEntry, 0, um.maxSize)
}

// Undo undoes the last command.
func (um *UndoManager) Undo(ctx context.Context) error {
	um.mu.Lock()
	if len(um.undoStack) == 0 {
		um.mu.Unlock()
		return fmt.Errorf("nothing to undo")
	}

	// Pop from undo stack
	entry := um.undoStack[len(um.undoStack)-1]
	um.undoStack = um.undoStack[:len(um.undoStack)-1]

	// Push to redo stack
	um.redoStack = append(um.redoStack, entry)
	um.mu.Unlock()

	// Execute undo
	if entry.Command.Undo == nil {
		return fmt.Errorf("command does not support undo")
	}

	return entry.Command.Undo(ctx)
}

// Redo redoes the last undone command.
func (um *UndoManager) Redo(ctx context.Context) error {
	um.mu.Lock()
	if len(um.redoStack) == 0 {
		um.mu.Unlock()
		return fmt.Errorf("nothing to redo")
	}

	// Pop from redo stack
	entry := um.redoStack[len(um.redoStack)-1]
	um.redoStack = um.redoStack[:len(um.redoStack)-1]

	// Push back to undo stack
	um.undoStack = append(um.undoStack, entry)
	um.mu.Unlock()

	// Re-execute command
	_, err := entry.Command.Execute(ctx, []string{})
	return err
}

// CanUndo returns whether there are commands to undo.
func (um *UndoManager) CanUndo() bool {
	um.mu.RLock()
	defer um.mu.RUnlock()
	return len(um.undoStack) > 0
}

// CanRedo returns whether there are commands to redo.
func (um *UndoManager) CanRedo() bool {
	um.mu.RLock()
	defer um.mu.RUnlock()
	return len(um.redoStack) > 0
}

// Clear clears both undo and redo stacks.
func (um *UndoManager) Clear() {
	um.mu.Lock()
	defer um.mu.Unlock()
	um.undoStack = make([]UndoEntry, 0, um.maxSize)
	um.redoStack = make([]UndoEntry, 0, um.maxSize)
}

// UndoStackSize returns the size of the undo stack.
func (um *UndoManager) UndoStackSize() int {
	um.mu.RLock()
	defer um.mu.RUnlock()
	return len(um.undoStack)
}

// RedoStackSize returns the size of the redo stack.
func (um *UndoManager) RedoStackSize() int {
	um.mu.RLock()
	defer um.mu.RUnlock()
	return len(um.redoStack)
}
