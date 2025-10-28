package vak

import (
	"sort"
)

// Completer provides autocompletion functionality.
type Completer struct {
	registry *Registry
}

// NewCompleter creates a new completer.
func NewCompleter(registry *Registry) *Completer {
	return &Completer{
		registry: registry,
	}
}

// Complete provides completion suggestions for the given input.
func (c *Completer) Complete(input string) []string {
	suggestions := c.registry.Complete(input)
	sort.Strings(suggestions)
	return suggestions
}

// CompleteCommand provides completion for command names.
func (c *Completer) CompleteCommand(prefix string) []string {
	commands := c.registry.List()
	suggestions := make([]string, 0)

	for _, cmd := range commands {
		if startsWith(cmd.Name, prefix) {
			suggestions = append(suggestions, cmd.Name)
		}
	}

	sort.Strings(suggestions)
	return suggestions
}

// CompleteFlag provides completion for flag names.
func (c *Completer) CompleteFlag(commandName, prefix string) []string {
	cmd, exists := c.registry.Get(commandName)
	if !exists {
		return nil
	}

	suggestions := make([]string, 0)
	for _, flag := range cmd.Flags {
		flagName := "--" + flag.Name
		if startsWith(flagName, prefix) {
			suggestions = append(suggestions, flagName)
		}
		if flag.Short != "" {
			shortName := "-" + flag.Short
			if startsWith(shortName, prefix) {
				suggestions = append(suggestions, shortName)
			}
		}
	}

	sort.Strings(suggestions)
	return suggestions
}
