# AI Agent Module Map and Implementation Roadmap

Follow the Steps in order. The roadmap starts with compile-time structure and small contracts, then adds infrastructure, agent behavior, memory, runtime interfaces, security, and the UI. Every Step produces a result that can be tested before the next module depends on it.

## Project boundary

The first release is intentionally smaller than GoClaw:

- One Go gateway and one React web application.
- One OpenAI-compatible LLM provider.
- PostgreSQL with pgvector.
- One agent profile with per-user sessions.
- Durable history, summaries, and episodic memory.
- A small tool set: `datetime`, `memory_search`, `read_file`, and approval-gated `write_file`.
- HTTP, WebSocket streaming, and a thin CLI.
- Structured logs and a durable run timeline.

The runtime path is:

```text
Web UI / CLI
      ↓
HTTP + WebSocket gateway
      ↓
Chat application service
      ↓
Per-session scheduler
      ↓
Agent loop
  ↙      ↓       ↘
LLM    Tools    Memory
provider          ↓
      PostgreSQL + pgvector
```

## Modules to implement

| Module | Main location | Responsibilities | First delivered |
|---|---|---|---|
| Process and composition | `cmd/agentkit`, `internal/app` | Commands, configuration, dependency wiring, startup, and shutdown. | Steps 1–2 |
| Canonical model | `internal/model` | Provider-neutral messages, tool calls, requests, responses, usage, and interfaces. | Step 3 |
| Durable data | `migrations`, `internal/store/postgres` | Agent settings, scoped sessions, ordered messages, summaries, approvals, memory, and run facts. | Steps 4–5 |
| LLM provider | `internal/providers/openai` | Wire translation, HTTP calls, SSE streaming, tool-call fragments, usage, and safe errors. | Steps 6–7 |
| Prompt composition | `internal/agent` | Stable identity, editable instructions, memory, summary, history, and current-message ordering. | Step 8 |
| Local skills | `skills`, `internal/skills` | Validate, load, select, and inject local `SKILL.md` instruction packages. | Step 9 |
| Agent execution | `internal/agent` | Per-run state, bounded think/act/observe transitions, persistence, usage, failure, and cancellation. | Steps 10 and 13 |
| Tool runtime | `internal/tools` | Registry, JSON Schema, policy, path safety, approvals, limits, redaction, and canonical results. | Steps 11–12 |
| Session context | `internal/session` | Provider-valid history, token budgets, pruning, tool-pair preservation, and durable compaction. | Steps 14–15 |
| Memory | `internal/memory` | Memory policy, embeddings, scoped hybrid retrieval, prompt injection, and explicit search. | Steps 16–18 |
| Background runtime | `internal/runtime` | Domain events, per-session FIFO scheduling, global concurrency, backpressure, and draining. | Steps 18–19 |
| Gateway interfaces | `internal/transport/http`, `internal/transport/websocket` | Authentication, validation, HTTP APIs, WebSocket protocol, streaming events, and recovery. | Steps 20–21 |
| CLI adapter | `cmd/agentkit` | `serve`, `chat`, migration, and provider verification commands using gateway contracts. | Step 22 |
| Security and audit | Cross-cutting | Trusted identity, scoped queries, exact-action approval, transport hardening, redaction, and audit. | Step 23 |
| Web application | `web` | Responsive shell, typed clients, settings, chat, reconnect recovery, approvals, and timeline. | Steps 24–25 |
| Acceptance and extensions | Tests and documentation | End-to-end verification and explicit extension points for deferred modules. | Step 26 |

## Structure of every Step

Each Step has one measurable `Step outcome` and several Tasks. Every Task uses the same learning structure:

1. `Theory` explains the component, its architectural role, why the boundary exists, and where to research it in GoClaw.
2. `Goal` states the capability produced by that Task.
3. `Guide to implement` provides the concrete files, contracts, code shape, tests, and completion checks.

Do not begin the next Step until the current `Step outcome` is demonstrated by a build, test, fixture, or runnable behavior.

## Phase 1 — Foundation and durable boundaries

| Step | Module work | Completion goal | Depends on |
|---|---|---|---|
| [Step 1 — Product Scope and Runtime Map](01-foundation/Step1.ProductScopeAndRuntimeMap.md) | Fix the first-release scope and assign every runtime responsibility to one package. | The runtime path and build-now/deferred scope are explicit. | None |
| [Step 2 — Project Skeleton and Configuration](01-foundation/Step2.ProjectSkeletonAndConfiguration.md) | Create the Go module, package skeleton, entry point, and validated configuration. | The empty application builds and fails fast on invalid config. | Step 1 |
| [Step 3 — Core Contracts and Dependency Inversion](01-foundation/Step3.CoreContractsAndDependencyInversion.md) | Define canonical model, provider, run, and session contracts. | The core compiles with fakes and imports no SDK or SQL driver. | Step 2 |
| [Step 4 — Durable Data Model and Migrations](01-foundation/Step4.DurableDataModelAndMigrations.md) | Model ownership and create the first PostgreSQL migration. | Agent, session, and ordered-message invariants exist in the schema. | Step 3 |
| [Step 5 — PostgreSQL Stores and Scope](01-foundation/Step5.PostgreSQLStoresAndScope.md) | Implement session and agent stores with transaction and scope guarantees. | State survives restart and cannot be read across user scope. | Step 4 |

## Phase 2 — Provider, prompt, and basic agent execution

| Step | Module work | Completion goal | Depends on |
|---|---|---|---|
| [Step 6 — OpenAI-Compatible Provider Basics](02-core-engine/Step6.OpenAICompatibleProviderBasics.md) | Implement capabilities, wire translation, and non-streaming chat. | A canonical request reaches one provider and returns a canonical response. | Steps 3 and 5 |
| [Step 7 — Provider Streaming and Adapter Tests](02-core-engine/Step7.ProviderStreamingAndAdapterTests.md) | Add bounded SSE parsing, fragmented tool-call assembly, and adapter tests. | Streaming, cancellation, and provider errors behave safely. | Step 6 |
| [Step 8 — Prompt Builder and Context Composition](02-core-engine/Step8.PromptBuilderAndContextComposition.md) | Compose ordered prompt sections from explicit inputs. | Prompt order, trust labels, and exactly-once current message are tested. | Steps 3 and 5 |
| [Step 9 — Local Skills Loading and Selection](02-core-engine/Step9.LocalSkillsLoadingAndSelection.md) | Load and select safe local `SKILL.md` packages. | Only enabled skills appear in deterministic prompt order. | Step 8 |
| [Step 10 — Basic Text Agent Loop](02-core-engine/Step10.BasicTextAgentLoop.md) | Add per-run state, preparation, one bounded think loop, and final persistence. | A text-only request completes without shared mutable run state. | Steps 5–9 |

## Phase 3 — Tools and long-session behavior

| Step | Module work | Completion goal | Depends on |
|---|---|---|---|
| [Step 11 — Tool Contracts, Registry, and Policy](02-core-engine/Step11.ToolContractsRegistryAndPolicy.md) | Define the tool boundary, registry pipeline, and effective capability policy. | No tool can execute outside validation and authorization. | Steps 3 and 10 |
| [Step 12 — Workspace Tools and Runtime Protections](02-core-engine/Step12.WorkspaceToolsAndRuntimeProtections.md) | Add the initial tools, workspace containment, limits, and approval checks. | Safe reads work and unauthorized or unsafe writes fail. | Step 11 |
| [Step 13 — Tool-Aware Agent Loop](02-core-engine/Step13.ToolAwareAgentLoop.md) | Extend the loop with act-observe iterations, tool messages, and failure semantics. | Text and tool round trips preserve order, IDs, budgets, and cancellation. | Steps 10–12 |
| [Step 14 — Session History and Context Budget](02-core-engine/Step14.SessionHistoryAndContextBudget.md) | Build provider-valid history, count tokens, and prune complete groups. | A long prompt fits its context window without breaking tool pairs. | Step 13 |
| [Step 15 — Durable Summary Compaction](02-core-engine/Step15.DurableSummaryCompaction.md) | Compact old turns into a durable summary checkpoint. | Repeated compaction survives restart and does not reprocess the same prefix. | Step 14 |

## Phase 4 — Memory and asynchronous processing

| Step | Module work | Completion goal | Depends on |
|---|---|---|---|
| [Step 16 — Memory Architecture and Policy](03-memory-and-tools/Step16.MemoryArchitectureAndPolicy.md) | Define layers, provenance, retention, scope, retrieval flow, and trust rules. | The memory contract and policy are explicit before storage exists. | Steps 8 and 15 |
| [Step 17 — Hybrid Memory Storage and Retrieval](03-memory-and-tools/Step17.HybridMemoryStorageAndRetrieval.md) | Add pgvector schema, embeddings, hybrid ranking, and token-capped injection. | A scoped fact can be retrieved by exact term or paraphrase. | Step 16 |
| [Step 18 — Domain Events and Memory Consolidation](03-memory-and-tools/Step18.DomainEventsAndMemoryConsolidation.md) | Add an in-process event bus and retry-safe session consolidation. | Completed sessions produce episodic records outside the response path. | Steps 15 and 17 |

## Phase 5 — Runtime and external interfaces

| Step | Module work | Completion goal | Depends on |
|---|---|---|---|
| [Step 19 — Session Scheduler and Shutdown](04-runtime-and-interfaces/Step19.SessionSchedulerAndShutdown.md) | Add keyed FIFO queues, global concurrency, cancellation, backpressure, and draining. | Same-session work is serialized while unrelated sessions run concurrently. | Steps 13 and 18 |
| [Step 20 — Gateway and HTTP API](04-runtime-and-interfaces/Step20.GatewayAndHTTPAPI.md) | Add the chat application service, edge authentication, health, chat, and request logs. | HTTP exposes the same durable scheduler-backed behavior. | Step 19 |
| [Step 21 — WebSocket Protocol and Runtime Events](04-runtime-and-interfaces/Step21.WebSocketProtocolAndRuntimeEvents.md) | Add typed frames, connection lifecycle, chat events, and reconnect recovery. | A client can start, observe, disconnect from, and recover a run. | Step 20 |
| [Step 22 — CLI Adapter](04-runtime-and-interfaces/Step22.CLIAdapter.md) | Add gateway-backed `serve`, `chat`, migration, and provider commands. | CLI and HTTP continue the same durable session behavior. | Step 20 |

## Phase 6 — Security, UI, and product acceptance

| Step | Module work | Completion goal | Depends on |
|---|---|---|---|
| [Step 23 — Security Hardening, Approvals, and Audit](05-production/Step23.SecurityHardeningApprovalsAndAudit.md) | Apply trusted identity, scoped stores, transport controls, exact-action approval, and audit. | Adversarial tests show that every external path has an independent enforcement point. | Steps 12 and 19–22 |
| [Step 24 — Web UI Foundation and Agent Settings](06-user-interface/Step24.WebUIFoundationAndAgentSettings.md) | Build the React shell, typed clients, auth state, state boundaries, and settings. | The responsive UI connects to the gateway and persists server-backed settings. | Steps 20–23 |
| [Step 25 — Chat UI, Reconnection, and Run Timeline](06-user-interface/Step25.ChatUIReconnectionAndRunTimeline.md) | Add session navigation, streaming messages, tool approval UX, recovery, and timeline. | Real-time UI state always reconciles with durable server state. | Steps 21, 23, and 24 |
| [Step 26 — End-to-End Acceptance and Extension Review](06-user-interface/Step26.EndToEndAcceptanceAndExtensionReview.md) | Run the full product scenario and identify extension points. | The product survives restart and negative security paths, with deferred modules documented. | All previous Steps |

## Deferred modules

The first implementation does not include:

- Multi-tenancy, tenant administration, or full RBAC.
- Agent teams, subagents, and delegation.
- Telegram, Discord, WhatsApp, or other messaging channels.
- Cron, heartbeat, and autonomous schedules.
- Multi-provider routing, load balancing, and fallback.
- MCP, ACP, unrestricted shell execution, or browser automation.
- Skill publishing, marketplace management, runtime installation, or self-evolution.
- Document-ingestion pipelines, Knowledge Vault, or a knowledge graph.
- Voice, TTS, media generation, desktop, or distributed queues.

These features should extend the completed contracts rather than change the dependency direction of the core.
