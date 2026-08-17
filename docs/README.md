# Build a GoClaw-Like Agent in Go

This mini-course leads you from an empty Go project to a complete AI-agent product with a web interface. The numbered GoClaw documents in this directory are supporting source knowledge: each course step introduces the few author documents that are useful at that moment.

## What this course covers

| Module | Steps | Focus |
|---|---|---|
| 01 — Foundation | 1–3 | Structure, contracts, configuration, and durable data |
| 02 — Core engine | 4–7 | Providers, prompts, orchestration, and conversation context |
| 03 — Memory and tools | 8–10 | Memory design, retrieval, and action boundaries |
| 04 — Runtime and interfaces | 11–15 | Scheduling, gateway, API, live events, and channels |
| 05 — Platform | 16–19 | Isolation, observability, teams, skills, and evolution |
| 06 — User interface | 20–22 | Application shell, chat experience, and management UI |

## Learning order

1. [Step1.FolderStructureInitialization.md](01-foundation/Step1.FolderStructureInitialization.md)
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
12. [Step12.GatewayFoundationAndIdentity.md](04-runtime-and-interfaces/Step12.GatewayFoundationAndIdentity.md)
13. [Step13.CLIAndHTTPAPI.md](04-runtime-and-interfaces/Step13.CLIAndHTTPAPI.md)
14. [Step14.WebSocketRPCAndLiveEvents.md](04-runtime-and-interfaces/Step14.WebSocketRPCAndLiveEvents.md)
15. [Step15.ChannelConnectorsAndMessageDelivery.md](04-runtime-and-interfaces/Step15.ChannelConnectorsAndMessageDelivery.md)
16. [Step16.TenantIsolationAndSecurityControls.md](05-production/Step16.TenantIsolationAndSecurityControls.md)
17. [Step17.ObservabilityAndVerification.md](05-production/Step17.ObservabilityAndVerification.md)
18. [Step18.MultiAgentTeamsAndDelegation.md](05-production/Step18.MultiAgentTeamsAndDelegation.md)
19. [Step19.SkillsPublishingAndAgentEvolution.md](05-production/Step19.SkillsPublishingAndAgentEvolution.md)
20. [Step20.WebUIFoundationAndApplicationShell.md](06-user-interface/Step20.WebUIFoundationAndApplicationShell.md)
21. [Step21.ChatInterfaceAndRealtimeStreaming.md](06-user-interface/Step21.ChatInterfaceAndRealtimeStreaming.md)
22. [Step22.ManagementUIAndProductPolish.md](06-user-interface/Step22.ManagementUIAndProductPolish.md)

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

The author reference documents stay in the root `docs/` directory so a concept is always one click away when it appears in a lesson.
