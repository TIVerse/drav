package vak

import (
	"fmt"
	"strings"
)

// Parser parses command strings.
type Parser struct{}

// NewParser creates a new parser.
func NewParser() *Parser {
	return &Parser{}
}

// ParsedCommand represents a parsed command.
type ParsedCommand struct {
	Name  string
	Args  []string
	Flags map[string]any
}

// Parse parses a command string.
func (p *Parser) Parse(input string) (*ParsedCommand, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	tokens := p.tokenize(input)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens")
	}

	parsed := &ParsedCommand{
		Name:  tokens[0],
		Args:  make([]string, 0),
		Flags: make(map[string]any),
	}

	// Parse remaining tokens
	i := 1
	for i < len(tokens) {
		token := tokens[i]

		if strings.HasPrefix(token, "--") {
			// Long flag
			flagName := strings.TrimPrefix(token, "--")
			if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "-") {
				// Flag with value
				parsed.Flags[flagName] = tokens[i+1]
				i += 2
			} else {
				// Boolean flag
				parsed.Flags[flagName] = true
				i++
			}
		} else if strings.HasPrefix(token, "-") && len(token) == 2 {
			// Short flag
			flagName := strings.TrimPrefix(token, "-")
			if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "-") {
				// Flag with value
				parsed.Flags[flagName] = tokens[i+1]
				i += 2
			} else {
				// Boolean flag
				parsed.Flags[flagName] = true
				i++
			}
		} else {
			// Regular argument
			parsed.Args = append(parsed.Args, token)
			i++
		}
	}

	return parsed, nil
}

// tokenize splits input into tokens, respecting quotes.
func (p *Parser) tokenize(input string) []string {
	tokens := make([]string, 0)
	current := ""
	inQuote := false
	quoteChar := rune(0)

	for _, ch := range input {
		if ch == '"' || ch == '\'' {
			if !inQuote {
				inQuote = true
				quoteChar = ch
			} else if ch == quoteChar {
				inQuote = false
				quoteChar = 0
			} else {
				current += string(ch)
			}
		} else if ch == ' ' && !inQuote {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}

	if current != "" {
		tokens = append(tokens, current)
	}

	return tokens
}
