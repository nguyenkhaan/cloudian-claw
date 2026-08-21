# Build a GoClaw-Like Agent in Go (upstream reference material)

> **Cloudclaw note:** The documents in this `docs/references/` directory describe
> the upstream **GoClaw** project. They are retained as *supporting source
> knowledge* for Cloudclaw's design — not as Cloudclaw's implementation plan or
> directory structure. Cloudclaw's actual layout is the DDD scaffold in
> `README.md` > Folder structure, and its plan is `docs/plans/plan.md`. The
> concepts are authoritative; the module paths (e.g. `internal/providers`,
> `internal/skills`, `ui/web`) are GoClaw's and must be adapted to Cloudclaw's
> packages.

This mini-course leads you from an empty Go project to a stable student-scale AI agent with a web interface. It follows GoClaw's architecture while deliberately implementing only the components that matter most: one provider, durable sessions and memory, safe tools, a gateway, real-time chat, and a usable UI. The numbered GoClaw documents in this directory are supporting source knowledge; each step introduces the few author documents that are useful at that moment.

## What this course covers

| Module | Steps | Focus |
|---|---|---|
| 01 — Foundation | 1–5 | Scope, skeleton, contracts, migrations, and PostgreSQL stores |
| 02 — Core engine | 6–15 | Provider, prompt, skills, agent loops, tools, history, and compaction |
| 03 — Memory and tools | 16–18 | Memory policy, hybrid retrieval, domain events, and consolidation |
| 04 — Runtime and interfaces | 19–22 | Scheduling, HTTP, WebSocket, and CLI adapters |
| 05 — Secure project | 23 | Single-tenant hardening, approvals, and audit |
| 06 — User interface | 24–26 | Application shell, real-time chat, timeline, and acceptance review |

## Learning order

See [STEP_INDEX.md](STEP_INDEX.md) for the module map, dependency-ordered roadmap, and measurable outcome of every Step.

1. [Step1.ProductScopeAndRuntimeMap.md](01-foundation/Step1.ProductScopeAndRuntimeMap.md)
2. [Step2.ProjectSkeletonAndConfiguration.md](01-foundation/Step2.ProjectSkeletonAndConfiguration.md)
3. [Step3.CoreContractsAndDependencyInversion.md](01-foundation/Step3.CoreContractsAndDependencyInversion.md)
4. [Step4.DurableDataModelAndMigrations.md](01-foundation/Step4.DurableDataModelAndMigrations.md)
5. [Step5.PostgreSQLStoresAndScope.md](01-foundation/Step5.PostgreSQLStoresAndScope.md)
6. [Step6.OpenAICompatibleProviderBasics.md](02-core-engine/Step6.OpenAICompatibleProviderBasics.md)
7. [Step7.ProviderStreamingAndAdapterTests.md](02-core-engine/Step7.ProviderStreamingAndAdapterTests.md)
8. [Step8.PromptBuilderAndContextComposition.md](02-core-engine/Step8.PromptBuilderAndContextComposition.md)
9. [Step9.LocalSkillsLoadingAndSelection.md](02-core-engine/Step9.LocalSkillsLoadingAndSelection.md)
10. [Step10.BasicTextAgentLoop.md](02-core-engine/Step10.BasicTextAgentLoop.md)
11. [Step11.ToolContractsRegistryAndPolicy.md](02-core-engine/Step11.ToolContractsRegistryAndPolicy.md)
12. [Step12.WorkspaceToolsAndRuntimeProtections.md](02-core-engine/Step12.WorkspaceToolsAndRuntimeProtections.md)
13. [Step13.ToolAwareAgentLoop.md](02-core-engine/Step13.ToolAwareAgentLoop.md)
14. [Step14.SessionHistoryAndContextBudget.md](02-core-engine/Step14.SessionHistoryAndContextBudget.md)
15. [Step15.DurableSummaryCompaction.md](02-core-engine/Step15.DurableSummaryCompaction.md)
16. [Step16.MemoryArchitectureAndPolicy.md](03-memory-and-tools/Step16.MemoryArchitectureAndPolicy.md)
17. [Step17.HybridMemoryStorageAndRetrieval.md](03-memory-and-tools/Step17.HybridMemoryStorageAndRetrieval.md)
18. [Step18.DomainEventsAndMemoryConsolidation.md](03-memory-and-tools/Step18.DomainEventsAndMemoryConsolidation.md)
19. [Step19.SessionSchedulerAndShutdown.md](04-runtime-and-interfaces/Step19.SessionSchedulerAndShutdown.md)
20. [Step20.GatewayAndHTTPAPI.md](04-runtime-and-interfaces/Step20.GatewayAndHTTPAPI.md)
21. [Step21.WebSocketProtocolAndRuntimeEvents.md](04-runtime-and-interfaces/Step21.WebSocketProtocolAndRuntimeEvents.md)
22. [Step22.CLIAdapter.md](04-runtime-and-interfaces/Step22.CLIAdapter.md)
23. [Step23.SecurityHardeningApprovalsAndAudit.md](05-production/Step23.SecurityHardeningApprovalsAndAudit.md)
24. [Step24.WebUIFoundationAndAgentSettings.md](06-user-interface/Step24.WebUIFoundationAndAgentSettings.md)
25. [Step25.ChatUIReconnectionAndRunTimeline.md](06-user-interface/Step25.ChatUIReconnectionAndRunTimeline.md)
26. [Step26.EndToEndAcceptanceAndExtensionReview.md](06-user-interface/Step26.EndToEndAcceptanceAndExtensionReview.md)

## Source map

Read a supporting author document only when a step introduces it. The course extracts the concept you need first, then links to GoClaw for the full architecture, implementation details, and adjacent features. Treat GoClaw as the reference implementation, not as a directory structure to copy blindly.

```text
# GoClaw upstream structure (reference only — NOT Cloudclaw's layout)
main.go → cmd/ → gateway
                   ├── internal/agent + internal/pipeline
                   ├── internal/providers + internal/providerresolve
                   ├── internal/tools + internal/mcp
                   ├── internal/store/{pg,sqlitestore}
                   ├── internal/memory + internal/consolidation
                   └── internal/http + pkg/protocol + ui/
```

> The tree above is **GoClaw's** structure. Cloudclaw uses a different DDD
> scaffold — see `README.md` > Folder structure. The author reference documents
> stay in the root `docs/` directory so a concept is always one click away when
> it appears in a lesson. Channels, teams, ACP, credentialed execution, skill
> publishing, self-evolution, heartbeat, and platform-scale multi-tenancy are
> retained as architecture study, not compulsory implementation work.
