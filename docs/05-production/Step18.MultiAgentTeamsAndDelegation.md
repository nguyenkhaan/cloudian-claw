# Step 18 — Multi-Agent Teams and Delegation

**Knowledge depth: 8/10**

Read [11 — Agent Teams](../11-agent-teams.md) before adding a second agent. It explains GoClaw's teams, task board, mailbox, task ownership, and result aggregation. Revisit [13 — WebSocket Team & Delegation Events](../13-ws-team-events.md) to see how that work becomes visible to clients.

## Delegation is a workflow, not a recursive prompt

Spawning another model call is easy. Coordinating useful independent work requires explicit ownership, scope, task state, and a way to return artifacts to the parent.

```text
parent agent → creates scoped task → child agent run
child agent → records progress/artifacts → task result
parent agent → reads result → produces final answer
```

## Decide what a child may inherit

A child usually needs a focused prompt, a limited tool set, a work budget, and a workspace boundary. It should not automatically receive every credential, unrestricted shell ability, or all parent conversation history. GoClaw's delegation modes and agent links show why this decision deserves a real data model.

## Coordinate through durable state

Use task records and events rather than only in-memory channels. A task board makes progress inspectable, lets a user intervene, and lets a UI reconnect after a process or browser restart.

## Keep the first team small

Begin with one coordinator and one specialist. The conceptual goal is to learn result handoff and task ownership; broad fan-out can wait until scheduler limits, trace visibility, and cancellation are already clear.

