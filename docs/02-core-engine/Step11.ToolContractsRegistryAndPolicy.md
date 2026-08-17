# Step 11 — Tool Contracts, Registry, and Policy

**Knowledge depth: 8/10**

Create the runtime boundary that advertises, validates, and authorizes tools.

## Step outcome

No tool can execute outside the registry, and effective permission is the intersection of all restrictions.

## Task 1 — Define the tool boundary

### Theory

- **Architectural role:** `Tool` describes schema, capabilities, and execution, `Policy` decides permission from request context, and `Result` is the canonical observation returned to the agent loop.
- **Why:** A tool call from the model is untrusted input. Keeping capability metadata separate from implementation lets policy evaluate read, write, network, and exec access without parsing tool names or prompts.
- **GoClaw reference:** Compare `Tool` and optional capability interfaces in [`goclaw/internal/tools/types.go`](../../goclaw/internal/tools/types.go), `ToolCapability` and `ToolMetadata` in [`goclaw/internal/tools/capability.go`](../../goclaw/internal/tools/capability.go), and provider conversion through `ToProviderDef` in `types.go`.

Read [03 — Tools System](../03-tools-system.md) and [09 — Security](../09-security.md). For privileged command execution, study [19 — Credentialed Exec](../19-credentialed-exec.md); do not add a raw shell in this project.

The model proposes a tool call. The runtime validates, authorizes, executes, scrubs, truncates, and records it.

### Goal

Define an action surface that the model may propose while the runtime keeps full control.

### Guide to implement

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

- **Architectural role:** The registry is both the catalogue for provider definitions and the execution pipeline. Policy, validation, timeout, recovery, and scrubbing wrap the concrete tool implementation.
- **Why:** If a caller can invoke tools directly, a new path can bypass authorization. Deterministic definitions stabilize prompt caching and tests, while preserving the call ID keeps protocol valid even on failure.
- **GoClaw reference:** Follow `Register`, `ExecuteWithContext`, and `ProviderDefs` in [`goclaw/internal/tools/registry.go`](../../goclaw/internal/tools/registry.go), then inspect boundary tests in [`goclaw/internal/tools/registry_test.go`](../../goclaw/internal/tools/registry_test.go) and [`goclaw/internal/tools/boundary_test.go`](../../goclaw/internal/tools/boundary_test.go).

The registry is an execution boundary, not only a map. It owns duplicate detection, deterministic definitions, validation, policy, timeout, panic recovery, and result conversion.

### Goal

Create one choke point for every tool invocation, from lookup to canonical result.

### Guide to implement

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

- **Architectural role:** Registration says a tool exists, the agent allow-list narrows the surface, principal and resource scope narrow the data, and approval adds authorization for one exact mutation. No layer may broaden another layer.
- **Why:** A “last setting wins” rule can turn a less-trusted configuration into an override. Request metadata must be immutable and per-call so a singleton tool never keeps identity from an earlier request.
- **GoClaw reference:** Read the policy engine in [`goclaw/internal/tools/policy.go`](../../goclaw/internal/tools/policy.go), capability metadata in [`goclaw/internal/tools/capability.go`](../../goclaw/internal/tools/capability.go), and context propagation in [`goclaw/internal/tools/context_keys.go`](../../goclaw/internal/tools/context_keys.go). Use `rg 'Allow|Authorize' goclaw/internal/tools` to follow the exact call path in the current revision.

Use intersection, not “last setting wins”:

```text
registered tool
∩ agent allow-list
∩ authenticated scope
∩ resource boundary
∩ optional approval
```

### Goal

Calculate effective permission as the intersection of several independent restrictions.

### Guide to implement

Start with this policy:

| Tool | Default | Required check |
|---|---|---|
| `datetime` | Allow | None beyond validation. |
| `memory_search` | Allow | Same agent/user scope. |
| `read_file` | Allow | Path stays under workspace. |
| `write_file` | Approval | Path stays under workspace and approval matches arguments. |
| `exec` | Not registered | Outside course scope. |

Pass user, agent, session, workspace, and approval through immutable request context or explicit execution input. Do not mutate singleton tools for each request.
