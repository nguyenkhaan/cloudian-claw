# Step 12 — Workspace Tools and Runtime Protections

**Knowledge depth: 8/10**

Add a small safe tool set and protect filesystem and output boundaries.

## Step outcome

Time and bounded workspace tools run with path containment, limits, redaction, and approval where required.

## Task 1 — Implement workspace file tools

### Theory

- **Architectural role:** A resolver canonicalizes paths before I/O, tools apply size and type limits, and policy decides read, write, and approval access. Symlink resolution is part of authorization, not only path formatting.
- **Why:** Prefix checks on raw strings do not stop `..`, symlinks, or sibling paths with the same prefix. Excluding globs, shell syntax, and recursive deletion keeps each call's blast radius small and predictable.
- **GoClaw reference:** Study the read and path boundary in [`goclaw/internal/tools/filesystem.go`](../../goclaw/internal/tools/filesystem.go), writes in [`goclaw/internal/tools/filesystem_write.go`](../../goclaw/internal/tools/filesystem_write.go), other path helpers under [`goclaw/internal/tools`](../../goclaw/internal/tools), and security tests with `rg 'traversal|symlink' goclaw/internal/tools goclaw/internal/http`.

### Goal

Provide useful file access while making the workspace root a clear sandbox boundary.

### Guide to implement

For every requested path:

1. Reject empty and NUL-containing values.
2. Join relative paths to the configured workspace root.
3. Clean and resolve the path.
4. Resolve existing symlinks where applicable.
5. Confirm the result is the root or a descendant of it.
6. Apply byte limits to reads and writes.

Return a model-visible error for traversal or denied access. Do not reveal sensitive host paths in that error.

`write_file` should create or replace one explicit file only. Do not support globs, recursive deletion, or shell syntax.

## Task 2 — Connect the first tools

### Theory

- **Architectural role:** The composition root registers instances, the registry exposes definitions to the provider, the agent loop turns calls into executions and results, and concrete tools know only their minimum dependencies.
- **Why:** A small tool set with different risk classes makes the architecture prove discovery, scope, and approval without adding raw exec or network complexity.
- **GoClaw reference:** Inspect built-in wiring in [`goclaw/cmd/gateway_builtin_tools.go`](../../goclaw/cmd/gateway_builtin_tools.go) and [`goclaw/cmd/gateway_tools_wiring.go`](../../goclaw/cmd/gateway_tools_wiring.go). Concrete examples include [`goclaw/internal/tools/datetime.go`](../../goclaw/internal/tools/datetime.go) and the filesystem tools under [`goclaw/internal/tools`](../../goclaw/internal/tools).

### Goal

Complete one vertical tool loop containing a pure utility, scoped reads, and an approval-gated mutation.

### Guide to implement

Register:

- `datetime`: returns current time for an allowed timezone.
- `read_file`: reads a bounded UTF-8/text file from the workspace.
- `write_file`: optional, approval-gated bounded write.

Pass the registry definitions to each provider request. Step 13 connects returned tool messages to the agent loop.

Do not register a fake `memory_search` implementation here. Step 16 defines its stable contract, and Step 17 connects real scoped retrieval.

## Task 3 — Add loop and output protections

### Theory

- **Architectural role:** Cross-cutting protections live in the registry and agent loop instead of being implemented differently by every tool. A tool-specific override may only narrow the shared limits.
- **Why:** A valid tool can still hang, return megabytes, expose credentials, or be called repeatedly by the model. Preserving result order separates execution scheduling from conversation semantics.
- **GoClaw reference:** Inspect rate limiting and scrubbing in [`goclaw/internal/tools/registry.go`](../../goclaw/internal/tools/registry.go), result truncation in [`goclaw/internal/agent/tool_result_truncation.go`](../../goclaw/internal/agent/tool_result_truncation.go), repeated-call controls in [`goclaw/internal/agent/loop_tools.go`](../../goclaw/internal/agent/loop_tools.go), and timing in [`goclaw/internal/agent/tool_timing.go`](../../goclaw/internal/agent/tool_timing.go).

### Goal

Limit the time, size, rate, and repeated cycles of the tool subsystem.

### Guide to implement

Add:

- Per-tool timeouts.
- Maximum result bytes with an explicit truncation marker.
- Credential-pattern redaction.
- Rate limits for expensive tools.
- Detection of repeated identical calls/results.

Parallelize only independent read-only I/O. Persist observations in the model's original call order.

## Task 4 — Verify authorization

### Theory

- **Architectural role:** Registry-boundary tests use malicious inputs and failing tools, filesystem tests use a temporary root and symlinks, and every outcome must become a canonical tool message with the original call ID.
- **Why:** Happy-path unit tests for individual tools do not verify defense in depth. Denials and panics must also remain provider-valid so the model can recover without crashing the runtime.
- **GoClaw reference:** Review [`goclaw/internal/tools/registry_test.go`](../../goclaw/internal/tools/registry_test.go), [`goclaw/internal/tools/boundary_test.go`](../../goclaw/internal/tools/boundary_test.go), [`goclaw/internal/tools/policy_prefix_test.go`](../../goclaw/internal/tools/policy_prefix_test.go), and security tests found with `rg 'symlink|traversal|approval' goclaw/internal/tools`.

### Goal

Prove that no execution path bypasses lookup, validation, scope, approval, or result correlation.

### Guide to implement

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
