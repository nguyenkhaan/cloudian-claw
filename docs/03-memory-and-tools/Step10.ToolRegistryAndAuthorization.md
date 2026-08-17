# Step 10 — Tool Registry and Authorization

**Knowledge depth: 9/10**

This step gives the agent a small, safe action surface: time, memory search, and workspace file access.

## Task 1 — Define the tool boundary

### Theory

Read [03 — Tools System](../03-tools-system.md) and [09 — Security](../09-security.md). For privileged command execution, study [19 — Credentialed Exec](../19-credentialed-exec.md); do not add a raw shell in this project.

The model proposes a tool call. The runtime validates, authorizes, executes, scrubs, truncates, and records it.

### Practice guide

Create `internal/tools/types.go`:

```go
type Capability string

const (
	Read    Capability = "read"
	Write   Capability = "write"
	Network Capability = "network"
	Exec    Capability = "exec"
)

type Result struct {
	Content string
	IsError bool
}

type Tool interface {
	Name() string
	Description() string
	Parameters() map[string]any
	Capabilities() []Capability
	Execute(context.Context, map[string]any) Result
}

type Policy interface {
	Allow(context.Context, Tool, map[string]any) error
}
```

Every `Parameters` value must be a JSON Schema object.

## Task 2 — Implement the registry

### Theory

The registry is an execution boundary, not only a map. It owns duplicate detection, deterministic definitions, validation, policy, timeout, panic recovery, and result conversion.

### Practice guide

Implement these operations:

```go
type Registry struct {
	tools  map[string]Tool
	policy Policy
}

func (r *Registry) Register(Tool) error
func (r *Registry) Definitions(context.Context) []model.ToolDefinition
func (r *Registry) Execute(context.Context, model.ToolCall) model.Message
```

`Execute` must:

1. Reject unknown names and missing call IDs.
2. Validate arguments against JSON Schema.
3. Ask the policy for permission.
4. Run with a tool-specific deadline.
5. Recover panics at this boundary.
6. Scrub secrets and truncate output.
7. Return `role="tool"` with the original call ID for success and failure.

Sort definitions by name so provider requests and tests are deterministic.

## Task 3 — Implement capability policy

### Theory

Use intersection, not “last setting wins”:

```text
registered tool
∩ agent allow-list
∩ authenticated scope
∩ resource boundary
∩ optional approval
```

### Practice guide

Start with this policy:

| Tool | Default | Required check |
|---|---|---|
| `datetime` | Allow | None beyond validation. |
| `memory_search` | Allow | Same agent/user scope. |
| `read_file` | Allow | Path stays under workspace. |
| `write_file` | Approval | Path stays under workspace and approval matches arguments. |
| `exec` | Not registered | Outside course scope. |

Pass user, agent, session, workspace, and approval through immutable request context or explicit execution input. Do not mutate singleton tools for each request.

## Task 4 — Implement workspace file tools

### Practice guide

For every requested path:

1. Reject empty and NUL-containing values.
2. Join relative paths to the configured workspace root.
3. Clean and resolve the path.
4. Resolve existing symlinks where applicable.
5. Confirm the result is the root or a descendant of it.
6. Apply byte limits to reads and writes.

Return a model-visible error for traversal or denied access. Do not reveal sensitive host paths in that error.

`write_file` should create or replace one explicit file only. Do not support globs, recursive deletion, or shell syntax.

## Task 5 — Connect the first tools

### Practice guide

Register:

- `datetime`: returns current time for an allowed timezone.
- `memory_search`: uses Step 9 scoped retrieval.
- `read_file`: reads a bounded UTF-8/text file from the workspace.
- `write_file`: optional, approval-gated bounded write.

Pass the registry definitions to each provider request. Pass returned tool messages back into the Step 6 loop.

## Task 6 — Add loop and output protections

### Practice guide

Add:

- Per-tool timeouts.
- Maximum result bytes with an explicit truncation marker.
- Credential-pattern redaction.
- Rate limits for expensive tools.
- Detection of repeated identical calls/results.

Parallelize only independent read-only I/O. Persist observations in the model's original call order.

## Task 7 — Verify authorization

### Practice guide

Test:

1. Duplicate registration.
2. Unknown tool.
3. Invalid JSON argument type.
4. Denied write without approval.
5. `../` traversal and symlink escape.
6. Successful bounded read.
7. Panic recovery and output truncation.
8. Correct tool-call ID in every result.

This step is complete when no tool can run without independent runtime validation and authorization.
