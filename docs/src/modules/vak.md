# Vak - Command Engine

**Package:** `github.com/TIVerse/drav/pkg/vak`

## Overview

Vak (Sanskrit: वाक्, "speech, command") is DRAV's command engine. It provides command registration, parsing, execution, history, undo/redo, and completion.

## Key Concepts

### Commands

Commands are executable actions with metadata:

```go
type Command struct {
    Name        string
    Summary     string
    Description string
    Usage       string
    Flags       []Flag
    Execute     ExecuteFunc
}
```

### Command Registry

Central registry for all commands:

```go
registry := vak.NewRegistry()
registry.Register(myCommand)
```

### Command Palette

UI component for command discovery and execution (to be implemented with widgets).

## Core API

### Creating a Command

```go
cmd := vak.Command{
    Name:    "save",
    Summary: "Save current file",
    Usage:   "save [filename]",
    Execute: func(ctx context.Context, args []string) (vak.Result, error) {
        filename := "untitled.txt"
        if len(args) > 0 {
            filename = args[0]
        }
        
        if err := saveFile(filename); err != nil {
            return vak.ErrorResult(err.Error()), err
        }
        
        return vak.SuccessResult(fmt.Sprintf("Saved to %s", filename)), nil
    },
}
```

### Registering Commands

```go
registry := vak.NewRegistry()
if err := registry.Register(cmd); err != nil {
    log.Fatal(err)
}
```

### Executing Commands

```go
result, err := registry.Execute(ctx, "save myfile.txt")
if err != nil {
    log.Printf("Command failed: %v", err)
}

fmt.Println(result.Message())
```

### Command with Flags

```go
cmd := vak.Command{
    Name:    "grep",
    Summary: "Search for pattern",
    Flags: []vak.Flag{
        {Name: "case-sensitive", Short: "c", Type: vak.FlagBool},
        {Name: "count", Short: "n", Type: vak.FlagInt},
    },
    Execute: func(ctx context.Context, args []string) (vak.Result, error) {
        flags := vak.ParseFlags(args)
        caseSensitive := flags.Bool("case-sensitive")
        count := flags.Int("count")
        
        // Execute search
        results := search(args[0], caseSensitive, count)
        return vak.SuccessResult(fmt.Sprintf("Found %d matches", len(results))), nil
    },
}
```

## History

### Command History

Track executed commands:

```go
history := vak.NewHistory(100)  // Keep last 100 commands

// Add command
history.Add("save file.txt")

// Navigate history
prev := history.Previous()
next := history.Next()

// Search history
matches := history.Search("save")
```

### Persistence

Save/load history:

```go
// Save
if err := history.SaveToFile("~/.drav_history"); err != nil {
    log.Printf("Failed to save history: %v", err)
}

// Load
if err := history.LoadFromFile("~/.drav_history"); err != nil {
    log.Printf("Failed to load history: %v", err)
}
```

## Undo/Redo

### Undoable Commands

Implement undo functionality:

```go
type UndoableCommand struct {
    vak.Command
    Undo func(ctx context.Context) error
}

cmd := UndoableCommand{
    Command: vak.Command{
        Name: "delete-line",
        Execute: func(ctx context.Context, args []string) (vak.Result, error) {
            lineNum := parseInt(args[0])
            deletedLine := editor.DeleteLine(lineNum)
            
            // Store for undo
            ctx = context.WithValue(ctx, "deleted", deletedLine)
            ctx = context.WithValue(ctx, "lineNum", lineNum)
            
            return vak.SuccessResult("Line deleted"), nil
        },
    },
    Undo: func(ctx context.Context) error {
        line := ctx.Value("deleted").(string)
        lineNum := ctx.Value("lineNum").(int)
        editor.InsertLine(lineNum, line)
        return nil
    },
}
```

### Undo Stack

Manage undo/redo:

```go
stack := vak.NewUndoStack()

// Execute with undo support
result, err := stack.ExecuteWithUndo(ctx, cmd, args)

// Undo last command
if err := stack.Undo(); err != nil {
    log.Printf("Cannot undo: %v", err)
}

// Redo
if err := stack.Redo(); err != nil {
    log.Printf("Cannot redo: %v", err)
}
```

## Completion

### Command Completion

Auto-complete command names:

```go
completions := registry.Complete("sa")
// Returns: ["save", "save-as", "save-all"]
```

### Argument Completion

Custom completion for command arguments:

```go
cmd := vak.Command{
    Name: "open",
    Complete: func(ctx context.Context, args []string) []string {
        // Complete file names
        if len(args) == 0 {
            return listFiles(".")
        }
        prefix := args[len(args)-1]
        return listFiles(filepath.Dir(prefix))
    },
}
```

## Parser

### Command Parsing

Parse command strings:

```go
parser := vak.NewParser()
parsed, err := parser.Parse("save --force myfile.txt")
if err != nil {
    log.Fatal(err)
}

fmt.Println(parsed.Command)  // "save"
fmt.Println(parsed.Flags)    // map["force": true]
fmt.Println(parsed.Args)     // ["myfile.txt"]
```

### Quoted Arguments

Handle quoted strings:

```go
// Supports both single and double quotes
parsed, _ := parser.Parse(`echo "Hello World"`)
// Args: ["Hello World"]

parsed, _ = parser.Parse(`echo 'It\'s working'`)
// Args: ["It's working"]
```

## Patterns

### Application Commands

```go
func registerAppCommands(registry *vak.Registry, app *dravya.App) {
    registry.Register(vak.Command{
        Name:    "quit",
        Summary: "Exit application",
        Execute: func(ctx context.Context, args []string) (vak.Result, error) {
            app.Lifecycle().Shutdown(nil)
            return vak.SuccessResult("Goodbye!"), nil
        },
    })
    
    registry.Register(vak.Command{
        Name:    "help",
        Summary: "Show available commands",
        Execute: func(ctx context.Context, args []string) (vak.Result, error) {
            commands := registry.List()
            var help strings.Builder
            for _, cmd := range commands {
                help.WriteString(fmt.Sprintf("%s - %s\n", cmd.Name, cmd.Summary))
            }
            return vak.SuccessResult(help.String()), nil
        },
    })
}
```

### Editor Commands

```go
type Editor struct {
    buffer   *Buffer
    registry *vak.Registry
}

func (e *Editor) registerCommands() {
    e.registry.Register(vak.Command{
        Name:    "goto",
        Summary: "Go to line",
        Usage:   "goto <line>",
        Execute: func(ctx context.Context, args []string) (vak.Result, error) {
            if len(args) == 0 {
                return vak.ErrorResult("Line number required"), nil
            }
            
            line, err := strconv.Atoi(args[0])
            if err != nil {
                return vak.ErrorResult("Invalid line number"), err
            }
            
            e.buffer.GoToLine(line)
            return vak.SuccessResult(fmt.Sprintf("Moved to line %d", line)), nil
        },
    })
}
```

### Command Aliases

```go
type AliasRegistry struct {
    *vak.Registry
    aliases map[string]string
}

func (a *AliasRegistry) RegisterAlias(alias, command string) {
    a.aliases[alias] = command
}

func (a *AliasRegistry) Execute(ctx context.Context, input string) (vak.Result, error) {
    // Expand alias
    if expanded, ok := a.aliases[input]; ok {
        input = expanded
    }
    
    return a.Registry.Execute(ctx, input)
}
```

## Best Practices

### 1. Descriptive Names

Use clear, verb-based names:

```go
// Good
"save", "open", "close", "search", "replace"

// Bad
"s", "o", "c", "find-text", "do-replace"
```

### 2. Provide Summaries

Always include summary and usage:

```go
cmd := vak.Command{
    Name:        "search",
    Summary:     "Search for text in current file",
    Usage:       "search <pattern> [--case-sensitive]",
    Description: "Searches for the specified pattern...",
}
```

### 3. Validate Arguments

Check arguments early:

```go
Execute: func(ctx context.Context, args []string) (vak.Result, error) {
    if len(args) < 1 {
        return vak.ErrorResult("Missing required argument: filename"), 
               fmt.Errorf("filename required")
    }
    
    filename := args[0]
    if !fileExists(filename) {
        return vak.ErrorResult(fmt.Sprintf("File not found: %s", filename)),
               fmt.Errorf("file not found")
    }
    
    // Execute command
    return vak.SuccessResult("Success"), nil
}
```

### 4. Consistent Return Values

Always return appropriate results:

```go
// Success
return vak.SuccessResult("Operation completed"), nil

// Error
return vak.ErrorResult("Operation failed"), err

// With data
return vak.ResultWithData(data), nil
```

### 5. Command Categories

Organize commands by category:

```go
cmd := vak.Command{
    Name:     "save",
    Category: "File",
    // ...
}
```

## UI Integration

### Command Palette Widget

```go
type CommandPalette struct {
    registry   *vak.Registry
    input      *prana.Observable[string]
    filtered   *prana.Observable[[]vak.Command]
    selected   int
}

func (p *CommandPalette) Render(ctx maya.RenderContext) maya.View {
    commands := p.filtered.Get()
    
    items := make([]maya.ListItem, len(commands))
    for i, cmd := range commands {
        items[i] = maya.ListItem{
            Label:    fmt.Sprintf("%s - %s", cmd.Name, cmd.Summary),
            Value:    cmd.Name,
            Selected: i == p.selected,
        }
    }
    
    return maya.Column(
        maya.Input(p.input.Get(),
            maya.WithPlaceholder("Type a command..."),
        ),
        maya.List(items),
    )
}

func (p *CommandPalette) filter(query string) {
    commands := p.registry.Complete(query)
    p.filtered.Set(commands)
}
```

## Examples

### Simple Command

```go
registry := vak.NewRegistry()

registry.Register(vak.Command{
    Name:    "echo",
    Summary: "Echo text",
    Execute: func(ctx context.Context, args []string) (vak.Result, error) {
        text := strings.Join(args, " ")
        return vak.SuccessResult(text), nil
    },
})

result, _ := registry.Execute(context.Background(), "echo Hello World")
fmt.Println(result.Message())  // "Hello World"
```

### Command with Undo

```go
type TextEditor struct {
    content string
    undoStack *vak.UndoStack
}

func (e *TextEditor) registerCommands() {
    e.undoStack.Register(vak.UndoableCommand{
        Command: vak.Command{
            Name: "insert",
            Execute: func(ctx context.Context, args []string) (vak.Result, error) {
                text := strings.Join(args, " ")
                oldContent := e.content
                e.content += text
                
                // Store for undo
                ctx = context.WithValue(ctx, "oldContent", oldContent)
                return vak.SuccessResult("Text inserted"), nil
            },
        },
        Undo: func(ctx context.Context) error {
            e.content = ctx.Value("oldContent").(string)
            return nil
        },
    })
}
```

## Performance Considerations

### Command Lookup

Use hash map for O(1) lookup:

```go
// Good - already implemented in Registry
registry.Get("save")  // Fast

// Avoid linear search
for _, cmd := range allCommands {
    if cmd.Name == "save" { ... }  // Slow
}
```

### History Size

Limit history to reasonable size:

```go
history := vak.NewHistory(1000)  // Good
history := vak.NewHistory(1000000)  // Too large
```

## Related Modules

- **[Dravya](dravya.md)**: Command registry integration
- **[Maya](maya.md)**: Command palette UI

## See Also

- [Command Examples](../../examples/04-commands/)
- [Keyboard Shortcuts](../concepts.md#shortcuts)
