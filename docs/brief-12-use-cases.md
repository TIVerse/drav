# Section 12: Use Cases & Application Patterns

[← Back to Index](brief-index.md) | [Previous: Developer Experience](brief-11-developer-experience.md) | [Next: Research & Innovation →](brief-13-research-innovation.md)

---

## 12.1 DevOps Dashboards

### Real-Time Monitoring Dashboard

**Use Case**: Live system metrics visualization for SRE teams

```go
package main

import (
    "time"
    "github.com/abhineesh/drav"
)

type SystemDashboard struct {
    cpuUsage    drav.Observable[float64]
    memUsage    drav.Observable[float64]
    diskIO      drav.Observable[float64]
    networkIO   drav.Observable[float64]
    processes   drav.Observable[[]Process]
    logs        drav.Observable[[]LogEntry]
}

func (d *SystemDashboard) Render() drav.View {
    return drav.Column(
        // Header
        drav.Panel("System Dashboard", 
            drav.Text(time.Now().Format("15:04:05")),
        ),
        
        // Metrics row
        drav.Row(
            drav.Panel("CPU", drav.LineChart(d.cpuHistory)),
            drav.Panel("Memory", drav.LineChart(d.memHistory)),
            drav.Panel("Disk I/O", drav.BarChart(d.diskIO)),
            drav.Panel("Network", drav.BarChart(d.networkIO)),
        ),
        
        // Process list
        drav.Panel("Top Processes",
            drav.Table(
                []string{"PID", "Name", "CPU%", "Memory"},
                d.processTableData(),
            ),
        ),
        
        // Live logs
        drav.Panel("System Logs",
            drav.LogViewer(d.logs),
        ),
    )
}

func main() {
    app := drav.NewApp(drav.WithTheme(drav.DarkTheme))
    dashboard := NewSystemDashboard()
    
    // Update metrics every second
    app.OnTick(time.Second, dashboard.UpdateMetrics)
    
    // Register commands
    app.RegisterCommand("refresh", drav.Command{
        Description: "Force refresh metrics",
        Handler:     dashboard.ForceRefresh,
    })
    
    app.RegisterCommand("filter", drav.Command{
        Description: "Filter logs by level",
        Usage:       "filter <level>",
        Handler:     dashboard.FilterLogs,
    })
    
    app.SetRoot(dashboard)
    app.Run()
}
```

**Key Features**:
- Real-time metric updates
- Interactive charts
- Command-driven filtering
- Live log streaming

---

## 12.2 Kubernetes Dashboard

### Pod Manager TUI

**Use Case**: Manage Kubernetes pods from terminal

```go
type K8sDashboard struct {
    namespace   drav.Observable[string]
    pods        drav.Observable[[]Pod]
    selectedPod drav.Observable[*Pod]
    logs        drav.Observable[string]
}

func (k *K8sDashboard) Render() drav.View {
    return drav.Row(
        // Left panel: Pod list
        drav.Panel("Pods",
            drav.List(
                k.podNames(),
                k.selectedIndex,
            ).OnSelect(k.onPodSelect),
        ),
        
        // Right panels: Pod details and logs
        drav.Column(
            drav.Panel("Pod Details",
                k.renderPodDetails(),
            ),
            
            drav.Panel("Logs",
                drav.Text(k.logs.Get()),
            ),
        ),
    )
}

func main() {
    app := drav.NewApp()
    dashboard := NewK8sDashboard()
    
    // Commands
    app.RegisterCommand("ns", drav.Command{
        Description: "Switch namespace",
        Usage:       "ns <namespace>",
        Complete:    dashboard.CompleteNamespace,
        Handler:     dashboard.SwitchNamespace,
    })
    
    app.RegisterCommand("delete", drav.Command{
        Description: "Delete selected pod",
        Handler:     dashboard.DeletePod,
        Confirm:     true,  // Ask for confirmation
    })
    
    app.RegisterCommand("scale", drav.Command{
        Description: "Scale deployment",
        Usage:       "scale <deployment> <replicas>",
        Handler:     dashboard.ScaleDeployment,
    })
    
    app.SetRoot(dashboard)
    app.Run()
}
```

**Key Features**:
- Namespace switching
- Pod selection and details
- Live log tailing
- Interactive commands (delete, scale, etc.)

---

## 12.3 Git TUI

### Interactive Git Client

**Use Case**: Visual git workflow management

```go
type GitTUI struct {
    repoPath    string
    status      drav.Observable[*GitStatus]
    branches    drav.Observable[[]Branch]
    commits     drav.Observable[[]Commit]
    selectedFile drav.Observable[string]
    diff        drav.Observable[string]
}

func (g *GitTUI) Render() drav.View {
    return drav.Row(
        // Left: Changed files
        drav.Panel("Changed Files",
            drav.List(
                g.status.Get().ChangedFiles,
                g.selectedFileIndex,
            ).OnSelect(g.showDiff),
        ),
        
        // Center: Diff view
        drav.Panel("Diff",
            drav.Text(g.diff.Get()).Syntax("diff"),
        ),
        
        // Right: Commit history
        drav.Panel("Commits",
            g.renderCommitHistory(),
        ),
    )
}

func main() {
    app := drav.NewApp()
    git := NewGitTUI(".")
    
    // Git commands
    app.RegisterCommand("stage", drav.Command{
        Description: "Stage selected file",
        Shortcut:    "s",
        Handler:     git.StageFile,
    })
    
    app.RegisterCommand("commit", drav.Command{
        Description: "Commit staged changes",
        Usage:       "commit <message>",
        Handler:     git.Commit,
    })
    
    app.RegisterCommand("push", drav.Command{
        Description: "Push to remote",
        Handler:     git.Push,
    })
    
    app.RegisterCommand("branch", drav.Command{
        Description: "Create or switch branch",
        Usage:       "branch <name>",
        Complete:    git.CompleteBranch,
        Handler:     git.SwitchBranch,
    })
    
    app.SetRoot(git)
    app.Run()
}
```

**Key Features**:
- File staging
- Diff visualization
- Commit history
- Branch management
- Command palette for git operations

---

## 12.4 Database Client

### SQL Query Interface

**Use Case**: Interactive database exploration

```go
type DBClient struct {
    connection  *sql.DB
    databases   drav.Observable[[]string]
    tables      drav.Observable[[]string]
    query       drav.Observable[string]
    results     drav.Observable[*QueryResult]
    history     drav.Observable[[]Query]
}

func (db *DBClient) Render() drav.View {
    return drav.Column(
        // Query editor
        drav.Panel("Query",
            drav.TextArea(db.query).
                Syntax("sql").
                OnSubmit(db.ExecuteQuery),
        ),
        
        // Results table
        drav.Panel("Results",
            drav.Table(
                db.results.Get().Columns,
                db.results.Get().Rows,
            ).Paginate(50),
        ),
        
        // Status bar
        drav.Row(
            drav.Text(fmt.Sprintf("Rows: %d", db.results.Get().RowCount)),
            drav.Text(fmt.Sprintf("Time: %s", db.results.Get().Duration)),
        ),
    )
}

func main() {
    app := drav.NewApp()
    client := NewDBClient("postgres://localhost/mydb")
    
    app.RegisterCommand("connect", drav.Command{
        Description: "Connect to database",
        Usage:       "connect <connection-string>",
        Handler:     client.Connect,
    })
    
    app.RegisterCommand("export", drav.Command{
        Description: "Export results to CSV",
        Usage:       "export <filename>",
        Handler:     client.ExportResults,
    })
    
    app.SetRoot(client)
    app.Run()
}
```

**Key Features**:
- SQL editor with syntax highlighting
- Query execution
- Result pagination
- Export to CSV
- Connection management

---

## 12.5 Log Viewer

### Real-Time Log Analysis

**Use Case**: Monitor and analyze application logs

```go
type LogViewer struct {
    logs        drav.Observable[[]LogEntry]
    filter      drav.Observable[LogFilter]
    following   drav.Observable[bool]
    selectedLog drav.Observable[*LogEntry]
}

func (lv *LogViewer) Render() drav.View {
    return drav.Row(
        // Log stream
        drav.Panel("Logs",
            drav.LogList(lv.filteredLogs()).
                Follow(lv.following.Get()).
                Highlight(lv.filter.Get().Pattern),
        ),
        
        // Log details
        drav.Panel("Details",
            lv.renderLogDetails(),
        ),
    )
}

func main() {
    app := drav.NewApp()
    viewer := NewLogViewer("/var/log/app.log")
    
    app.RegisterCommand("filter", drav.Command{
        Description: "Filter logs",
        Usage:       "filter <pattern>",
        Handler:     viewer.SetFilter,
    })
    
    app.RegisterCommand("level", drav.Command{
        Description: "Filter by log level",
        Usage:       "level <error|warn|info|debug>",
        Complete:    drav.Complete([]string{"error", "warn", "info", "debug"}),
        Handler:     viewer.FilterByLevel,
    })
    
    app.RegisterCommand("follow", drav.Command{
        Description: "Follow log tail",
        Handler:     viewer.ToggleFollow,
    })
    
    app.SetRoot(viewer)
    app.Run()
}
```

**Key Features**:
- Real-time log streaming
- Pattern matching and highlighting
- Level filtering
- Follow mode
- Log details panel

---

## 12.6 API Testing Tool

### HTTP Client TUI

**Use Case**: Test and debug REST APIs

```go
type APITester struct {
    request     drav.Observable[*HTTPRequest]
    response    drav.Observable[*HTTPResponse]
    collections drav.Observable[[]RequestCollection]
    history     drav.Observable[[]Request]
}

func (at *APITester) Render() drav.View {
    return drav.Column(
        // Request builder
        drav.Panel("Request",
            drav.Column(
                drav.Input(at.request.Get().URL, "URL"),
                drav.Select(
                    []string{"GET", "POST", "PUT", "DELETE"},
                    at.methodIndex,
                ),
                drav.TextArea(at.request.Get().Body).
                    Syntax("json"),
            ),
        ),
        
        // Response viewer
        drav.Panel("Response",
            drav.Tabs([]drav.Tab{
                {Name: "Body", Content: at.responseBody()},
                {Name: "Headers", Content: at.responseHeaders()},
                {Name: "Cookies", Content: at.responseCookies()},
            }),
        ),
    )
}

func main() {
    app := drav.NewApp()
    tester := NewAPITester()
    
    app.RegisterCommand("send", drav.Command{
        Description: "Send HTTP request",
        Shortcut:    "Ctrl+Enter",
        Handler:     tester.SendRequest,
    })
    
    app.RegisterCommand("save", drav.Command{
        Description: "Save request to collection",
        Usage:       "save <name>",
        Handler:     tester.SaveRequest,
    })
    
    app.SetRoot(tester)
    app.Run()
}
```

**Key Features**:
- Request builder
- Response visualization
- Request collections
- History tracking

---

## 12.7 File Manager

### Terminal File Browser

**Use Case**: Navigate and manage files

```go
type FileManager struct {
    currentDir  drav.Observable[string]
    files       drav.Observable[[]FileInfo]
    selected    drav.Observable[int]
    preview     drav.Observable[string]
}

func (fm *FileManager) Render() drav.View {
    return drav.Row(
        // File list
        drav.Panel(fm.currentDir.Get(),
            drav.List(fm.fileNames(), fm.selected).
                OnSelect(fm.updatePreview),
        ),
        
        // File preview
        drav.Panel("Preview",
            drav.Text(fm.preview.Get()).
                Syntax(fm.detectSyntax()),
        ),
    )
}

func main() {
    app := drav.NewApp()
    fm := NewFileManager(".")
    
    // Vim-like navigation
    app.OnKey(drav.Key('j'), fm.MoveDown)
    app.OnKey(drav.Key('k'), fm.MoveUp)
    app.OnKey(drav.KeyEnter, fm.OpenFile)
    
    app.RegisterCommand("cd", drav.Command{
        Description: "Change directory",
        Usage:       "cd <path>",
        Complete:    drav.DirectoryCompleter(),
        Handler:     fm.ChangeDirectory,
    })
    
    app.RegisterCommand("mkdir", drav.Command{
        Description: "Create directory",
        Handler:     fm.MakeDirectory,
    })
    
    app.SetRoot(fm)
    app.Run()
}
```

---

## 12.8 Data Pipeline Monitor

### ETL Job Dashboard

**Use Case**: Monitor data pipeline execution

```go
type PipelineMonitor struct {
    pipelines   drav.Observable[[]Pipeline]
    jobs        drav.Observable[[]Job]
    selectedJob drav.Observable[*Job]
    metrics     drav.Observable[*JobMetrics]
}

func (pm *PipelineMonitor) Render() drav.View {
    return drav.Column(
        // Pipeline overview
        drav.Row(
            drav.Panel("Active", pm.activePipelinesCount()),
            drav.Panel("Failed", pm.failedPipelinesCount()),
            drav.Panel("Success Rate", pm.successRate()),
        ),
        
        // Job list
        drav.Panel("Jobs",
            drav.Table(
                []string{"Job", "Status", "Progress", "Duration"},
                pm.jobTableData(),
            ),
        ),
        
        // Job details
        drav.Panel("Details",
            pm.renderJobDetails(),
        ),
    )
}
```

---

## 12.9 Common Patterns

### Pattern: Master-Detail View

```go
type MasterDetail struct {
    items    drav.Observable[[]Item]
    selected drav.Observable[*Item]
}

func (md *MasterDetail) Render() drav.View {
    return drav.Row(
        drav.Panel("Items", 
            drav.List(md.itemNames(), md.selectedIndex),
        ),
        drav.Panel("Details",
            md.renderDetails(),
        ),
    )
}
```

### Pattern: Tabbed Interface

```go
func (c *Component) Render() drav.View {
    return drav.Tabs([]drav.Tab{
        {Name: "Overview", Content: c.renderOverview()},
        {Name: "Details", Content: c.renderDetails()},
        {Name: "Logs", Content: c.renderLogs()},
    }, c.activeTab)
}
```

### Pattern: Modal Dialogs

```go
func (c *Component) Render() drav.View {
    view := c.renderMain()
    
    if c.showModal.Get() {
        view = drav.Modal("Confirm", 
            drav.Text("Are you sure?"),
            drav.Row(
                drav.Button("Yes", c.onConfirm),
                drav.Button("No", c.onCancel),
            ),
        ).Over(view)
    }
    
    return view
}
```

---

## Summary

**DRAV excels at**:
- Real-time monitoring dashboards
- Interactive CLI tools
- Data visualization
- System administration
- Developer tools
- Log analysis
- API testing
- File management

**Common Patterns**:
- Master-detail layouts
- Tabbed interfaces
- Command palettes
- Real-time updates
- Modal dialogs
- List selection

**Target Users**:
- DevOps engineers
- SREs
- Platform engineers
- Data engineers
- Developers

---

[← Back to Index](brief-index.md) | [Previous: Developer Experience](brief-11-developer-experience.md) | [Next: Research & Innovation →](brief-13-research-innovation.md)
