# Step 22 — Management UI and Product Polish

**Knowledge depth: 7/10**

Use [18 — HTTP REST API](../18-http-api.md), [19 — WebSocket RPC Methods](../19-websocket-rpc.md), and [20 — API Keys & Authentication](../20-api-keys-auth.md) when connecting management screens to the gateway. [23 — AI Agent Permission Matrix](../23-ai-agent-permission-matrix.md) tells the UI which controls should be visible or disabled for each role.

## Complete the product surface

The chat page proves the agent works. A usable agent system also needs interfaces for the concepts introduced throughout this course:

```text
Agents and provider selection
Sessions and run history
Tools, approvals, and credentials
Memory and knowledge sources
Skills and published versions
Schedules, teams, tasks, and traces
Tenant, user, and permission administration
```

These do not need to become a single crowded settings page. Group them by the mental model of the person using the system: conversation work, agent setup, operational visibility, and administration.

## Let roles shape the interface

The UI should reflect the server's authority decisions, not invent its own. A viewer can inspect allowed data, an operator can run normal workflows, and an administrator can manage configuration. The server still enforces the action; the UI simply makes capability boundaries understandable.

## Make complex work visible

For long-running tasks, show queued, running, waiting, complete, and failed states. Link a user-facing event to its relevant session, task, trace, or approval detail. This gives the platform a coherent experience across chat, background work, teams, and management.

## Finish with consistency

Apply a stable language system, error style, loading treatment, empty states, keyboard navigation, and responsive behavior across every screen. At this stage, the UI is no longer a wrapper around a chat box: it is the operational surface of the complete agent system.

