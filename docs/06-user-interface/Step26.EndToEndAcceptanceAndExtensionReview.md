# Step 26 — End-to-End Acceptance and Extension Review

**Knowledge depth: 8/10**

Verify the complete product slice and document how future modules attach to the core.

## Step outcome

The full acceptance scenario passes across restart and negative security paths, and every deferred feature has a named extension point.

## Task 1 — Run the end-to-end acceptance check

### Theory

- **Architectural role:** The scenario deliberately crosses every boundary and one restart. Each step provides evidence for an invariant instead of only checking that an answer exists.
- **Why:** Individual layers can pass tests while wiring is wrong, event IDs do not match, or state is not durable. The final negative paths confirm that new features did not break isolation or path safety.
- **GoClaw reference:** Use [`goclaw/cmd/gateway_setup.go`](../../goclaw/cmd/gateway_setup.go) as a wiring map, [`goclaw/tests`](../../goclaw/tests) as examples of layered acceptance and invariant tests, and the chat UI in [`goclaw/ui/web/src/pages/chat`](../../goclaw/ui/web/src/pages/chat) to compare the real flow.

### Goal

Prove the complete vertical slice from configuration through model, tools, memory, persistence, UI, and security controls.

### Guide to implement

Complete this scenario:

1. Start PostgreSQL, apply migrations, and start the gateway/UI.
2. Edit the agent's system prompt and enable one local skill.
3. Create a session and send a message over WebSocket.
4. Observe text and one safe tool round-trip.
5. Complete another session so episodic memory is consolidated.
6. Start a new session and recall the scoped fact.
7. Restart the gateway and reload history/settings.
8. Open the run timeline and inspect usage, tools, and retrieval facts.
9. Attempt cross-scope access and a workspace traversal; both must fail.

## Task 2 — Review extension boundaries

### Theory

- **Architectural role:** A channel is an inbound adapter to `ChatService`, a provider is an outbound adapter, cron produces scheduler jobs, a managed skill store sits behind the loader and selector, and team orchestration is a layer above runs and scheduling.
- **Why:** An extension point is meaningful only when its interface, owner, and dependency direction are clear. “Refactor later” is not an extension architecture.
- **GoClaw reference:** Compare channel dispatch in [`goclaw/internal/channels/dispatch.go`](../../goclaw/internal/channels/dispatch.go), the provider registry in [`goclaw/internal/providers/registry.go`](../../goclaw/internal/providers/registry.go), cron wiring in [`goclaw/cmd/gateway_cron.go`](../../goclaw/cmd/gateway_cron.go), the skill loader and store in [`goclaw/internal/skills/loader.go`](../../goclaw/internal/skills/loader.go) and [`goclaw/internal/store/skill_store.go`](../../goclaw/internal/store/skill_store.go), and team and delegation code under [`goclaw/internal/orchestration`](../../goclaw/internal/orchestration).

The remaining author documents explain platform features that should extend the core rather than rewrite it:

- [05 — Channels and Messaging](../05-channels-messaging.md)
- [11 — Agent Teams](../11-agent-teams.md)
- [14 — Skills Runtime](../14-skills-runtime.md)
- [15 — Core Skills System](../15-core-skills-system.md)
- [16 — Skill Publishing](../16-skill-publishing.md)
- [21 — Agent Evolution and Skill Management](../21-agent-evolution-and-skill-management.md)
- [22 — Heartbeat System](../22-heartbeat-system.md)

### Goal

Prove that the current architecture has concrete extension seams without partially implementing future features.

### Guide to implement

For each future feature, name its extension point in the current project. Examples: a channel calls `ChatService`, a new provider implements `model.Provider`, cron submits a scheduler job, and publishing adds a managed skill store behind the existing skill loader.

The course is complete when the acceptance scenario passes and future features have clear extension points without partial implementations.
