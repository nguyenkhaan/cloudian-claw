# Build a GoClaw-Like Agent in Go

This mini-course leads you from an empty Go project to a stable student-scale AI agent with a web interface. It follows GoClaw's architecture while deliberately implementing only the components that matter most: one provider, durable sessions and memory, safe tools, a gateway, real-time chat, and a usable UI. The numbered GoClaw documents in this directory are supporting source knowledge; each step introduces the few author documents that are useful at that moment.

## What this course covers

| Module | Steps | Focus |
|---|---|---|
| 01 — Foundation | 1–3 | Structure, contracts, configuration, and durable data |
| 02 — Core engine | 4–7 | Providers, prompts, orchestration, and conversation context |
| 03 — Memory and tools | 8–10 | Memory design, retrieval, and action boundaries |
| 04 — Runtime and interfaces | 11–13 | Scheduling, gateway, API, and live events |
| 05 — Secure project | 14 | Single-tenant access and safe boundaries |
| 06 — User interface | 15–16 | Application shell, real-time chat, and run timeline |

## Learning order

1. [Step1.ProjectScopeAndFolderStructure.md](01-foundation/Step1.ProjectScopeAndFolderStructure.md)
2. [Step2.CoreContractsAndConfiguration.md](01-foundation/Step2.CoreContractsAndConfiguration.md)
3. [Step3.DataModelAndStoreLayer.md](01-foundation/Step3.DataModelAndStoreLayer.md)
4. [Step4.LLMProviderAdapter.md](02-core-engine/Step4.LLMProviderAdapter.md)
5. [Step5.BootstrapPromptAndCoreSkills.md](02-core-engine/Step5.BootstrapPromptAndCoreSkills.md)
6. [Step6.AgentOrchestrationLoop.md](02-core-engine/Step6.AgentOrchestrationLoop.md)
7. [Step7.SessionHistoryAndContextBudget.md](02-core-engine/Step7.SessionHistoryAndContextBudget.md)
8. [Step8.MemoryArchitectureAndPolicy.md](03-memory-and-tools/Step8.MemoryArchitectureAndPolicy.md)
9. [Step9.HybridMemoryRetrievalAndConsolidation.md](03-memory-and-tools/Step9.HybridMemoryRetrievalAndConsolidation.md)
10. [Step10.ToolRegistryAndAuthorization.md](03-memory-and-tools/Step10.ToolRegistryAndAuthorization.md)
11. [Step11.SessionSchedulerAndBackgroundEvents.md](04-runtime-and-interfaces/Step11.SessionSchedulerAndBackgroundEvents.md)
12. [Step12.GatewayHTTPAPIAndCLI.md](04-runtime-and-interfaces/Step12.GatewayHTTPAPIAndCLI.md)
13. [Step13.WebSocketEventsAndUIContract.md](04-runtime-and-interfaces/Step13.WebSocketEventsAndUIContract.md)
14. [Step14.SingleTenantSecurityAndAccessControl.md](05-production/Step14.SingleTenantSecurityAndAccessControl.md)
15. [Step15.WebUIFoundationAndApplicationShell.md](06-user-interface/Step15.WebUIFoundationAndApplicationShell.md)
16. [Step16.ChatUIRunTimelineAndProjectReview.md](06-user-interface/Step16.ChatUIRunTimelineAndProjectReview.md)

## Source map

Read a supporting author document only when a step introduces it. The course extracts the concept you need first, then links to GoClaw for the full architecture, implementation details, and adjacent features. Treat GoClaw as the reference implementation, not as a directory structure to copy blindly.

```text
main.go → cmd/ → gateway
                   ├── internal/agent + internal/pipeline
                   ├── internal/providers + internal/providerresolve
                   ├── internal/tools + internal/mcp
                   ├── internal/store/{pg,sqlitestore}
                   ├── internal/memory + internal/consolidation
                   └── internal/http + pkg/protocol + ui/
```

The author reference documents stay in the root `docs/` directory so a concept is always one click away when it appears in a lesson. Channels, teams, ACP, credentialed execution, skill publishing, self-evolution, heartbeat, and platform-scale multi-tenancy are retained as architecture study, not compulsory implementation work.
