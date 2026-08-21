# Project Implementation Plan

## 1. Implementation Overview
- **Overall Implementation Strategy:** This project implements **Cloudclaw** (AI Agent Gateway local-first, single-owner). The implementation follows a bottom-up approach starting from the core domain types, application bootstrapping, and storage persistence layers. Once the foundation is laid, we implement the transport gateway, provider interfaces, tool registry, and agent loop execution orchestrator. Finally, memory similarity search, event tracing, the CLI, and the React Dashboard UI are integrated.
- **Major Modules (DDD layout — see `README.md` > Folder structure):**
  - Core Domain & Model Types (`internal/domain/entity`, `internal/domain/{entity,event,interface,service}`)
  - Config & Application Container (`internal/impl/config`)
  - PostgreSQL Persistence Store (`internal/impl/store`, `internal/impl/store/postgres`, `internal/impl/database`)
  - Gateway & Authentication Transport (`internal/transport/http`, `internal/transport/socket`)
  - Provider Abstraction & Model Execution (`internal/impl/provider`)
  - Skill and Rule Systems (`internal/impl/...` + `agent/skills`, `agent/rules`)
  - Tool Registry & Permissions (`internal/impl/tool`)
  - Agent Loop & State Machine (`internal/domain/service`)
  - Session & Local Memory (`internal/impl/session`, `internal/impl/memory`)
  - Observability & Realtime Events (`internal/transport/socket`, `internal/impl/...`)
  - CLI & React Dashboard Frontends (`cmd/agentkit`, `internal/tui`, `web/`)
- **Implementation Order & Dependencies:**
  1. Foundation & Bootstrapping (Phase 1)
  2. PostgreSQL Schema, Migrations & Stores (Phase 2)
  3. Gateway Transport & Authentication Middleware (Phase 3)
  4. Provider Abstraction & OpenAI-Compatible Layer (Phase 4)
  5. Skills System & Rules Sandbox (Phase 5)
  6. Tool Registry & Built-in File/CMD Tools (Phase 6)
  7. Runtime Agent Loop & Compaction (Phase 7)
  8. Message Persistence, Local Files Uploads, & Go Memory similarity search (Phase 8)
  9. Tracing, Metrics & Realtime event streams (Phase 9)
  10. API Keys & Global Config (Phase 10)
  11. Cobra CLI & BubbleTea Chat UI (Phase 11)
  12. React Dashboard with Onboarding Setup (Phase 12)
  13. Reliability, Fault Isolation & E2E Validation (Phase 13)

---

## 2. Architecture → Implementation Mapping
Below is the mapping of components defined in [architecture.md](file:///home/cloud/workspace/project/cloudclaw/docs/documents/architecture.md) and [idea.md](file:///home/cloud/workspace/project/cloudclaw/docs/documents/idea.md) to implementation Phases:
* **Gateway & Authentication:** Phase 3 (`internal/transport/http`)
* **Agent Router:** Phase 3 (`internal/transport/http`)
* **Agent Loop & Context Builder:** Phase 7 (`internal/domain/service`)
* **Provider Abstraction:** Phase 4 (`internal/impl/provider`)
* **Tool Registry:** Phase 6 (`internal/impl/tool`)
* **Skill Loader & Rule Engine:** Phase 5 (`agent/skills`, `agent/rules` + `internal/impl/...`)
* **Session Store:** Phase 8 (`internal/impl/session`, `internal/impl/store/postgres`)
* **Memory & Vector Search (Go Cosine):** Phase 8 (`internal/impl/memory`)
* **Realtime Events System:** Phase 9 (`internal/transport/socket`)
* **Tracing Store:** Phase 9 (`internal/impl/store/postgres`)
* **API Key & Config Store:** Phase 10 (`internal/impl/store/postgres`, `internal/impl/config`)
* **Cobra CLI / React UI:** Phase 11 & Phase 12 (`cmd/agentkit`, `internal/tui`, `web/`)

---

## 3. Phase 1 — Foundation & Bootstrapping
### Objective
Establish the compile-time contracts, base domain models, and core configuration loading so that the project can build cleanly without external drivers.

### Modules Covered
* `internal/domain/entity`
* `internal/domain/{entity,event,interface,service}`
* `internal/impl/config`

### Task 1.1 — Domain Entities Standardization
- **Description:** Define Go structures representing the core components of the system.
- **Objective:** Create compilation-stable domain models.
- **Requirements:**
  - Define `Agent`, `Provider`, `Model`, `Skill`, `Rule`, `Tool`, `Session`, `Message`, `Memory`, `Execution`, `APIKey`, and `Config` structs in `internal/domain/entity`.
  - Define `Usage` for calculate token per user'request 
  - Include basic validation functions for each entity.
  - Do not import any SQL driver or third-party HTTP transport libraries in `internal/domain/entity`.
- **Implementation Guidance:** Use Go standard types. Map to the database schemas defined later. Refer to GoClaw's representation of domain entities in [06-store-data-model.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/06-store-data-model.md).
- **Dependencies:** None.
- **Validation:** Run `go test ./internal/domain/entity/...` to verify they compile and validate.

### Task 1.2 — Runtime Contracts & State Machine Definitions
- **Description:** Define the structures that carry request and status state through the Agent Runtime.
- **Objective:** Establish the communication model of the Agent Loop.
- **Requirements:**
  - Define `ExecutionRequest`, `RuntimeContext`, `ExecutionState`, `ToolCall`, `ToolResult`, `ToolError`, `ModelResponse`, and `ExecutionResult`.
  - Implement `ExecutionState` containing Execution ID, Agent ID, Session ID, Current Model Turn, Tool Call Count, Retry Count, Status, and Error State.
  - State boundaries must include intermediate states (`INITIALIZING`, `BUILDING_CONTEXT`, `CALLING_PROVIDER`, `EXECUTING_TOOL`, `STREAMING_RESPONSE`, `FINALIZING`) and terminal states (`COMPLETED`, `FAILED`, `CANCELLED`, `LIMIT_REACHED`).
- **Implementation Guidance:** See [01-agent-loop.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/01-agent-loop.md) for loop state configurations.
- **Dependencies:** Task 1.1.
- **Validation:** Write unit tests demonstrating state transitions, ensuring terminal states cannot transition back.

### Task 1.3 — Configuration Loader & App Container
- **Description:** Implement static configuration parsing and application dependency injection setup.
- **Objective:** Initialize the Go server dependency graph cleanly.
- **Requirements:**
  - Load environment variables from `.env` using standard configuration patterns (e.g. `caarlos0/env` or custom loader).
  - Config must include: Server Host/Port, PostgreSQL URL, Default AI settings, Quota settings, Global Tools Toggle, and Execution limits (max tool calls, max retries, max duration, max context, max depth).
  - Implement the application `Container` in `internal/impl/config` to wire up configuration, database, stores, services, and gateways.
- **Implementation Guidance:** Read the configuration patterns in GoClaw's startup in [README.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/README.md).
- **Dependencies:** Task 1.2.
- **Validation:** Run configuration load tests with invalid/missing variables to verify fail-fast behavior.

### Phase Completion Criteria
* `go build ./...` compiles cleanly.
* Configuration parses successfully from environment variables.

---

## 4. Phase 2 — Persistence Layer
### Objective
Set up the PostgreSQL relational schema and store implementations to allow durable persistence of agent and session states.

### Modules Covered
* `migrations/`
* `internal/impl/store`
* `internal/impl/store/postgres`
* `internal/impl/database`

### Task 2.1 — Schema Design & Migrations
- **Description:** Write SQL migrations for all tables in PostgreSQL.
- **Objective:** Create the physical database schema.
- **Requirements:**
  - Create table definitions for `agents`, `providers`, `provider_models`, `skills`, `skill_versions`, `skill_assignments`, `rules`, `rule_versions`, `rule_assignments`, `tools`, `sessions`, `messages`, `memories`, `executions`, `execution_events`, `traces`, `spans`, `api_keys`, `api_key_usages`, and `global_configs`.
  - **Memory Table Rule:** Embeddings must be stored in a column of type `double precision[]` or `jsonb` instead of requiring pgvector to ensure local portability.
  - Implement a migration runner (e.g. `golang-migrate` or built-in runner) triggered at server startup or via CLI.
- **Implementation Guidance:** Reference GoClaw's DB architecture in [06-store-data-model.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/06-store-data-model.md).
- **Dependencies:** Phase 1.
- **Validation:** Execute migrations against a fresh PostgreSQL instance, verifying all constraints and tables exist.

### Task 2.2 — Store Interface Implementations
- **Description:** Implement SQL-specific logic for all store interfaces defined in `internal/store`.
- **Objective:** Provide abstraction over SQL execution.
- **Requirements:**
  - Build implementations for `AgentStore`, `ProviderStore`, `SessionStore`, `SkillStore`, `RuleStore`, `ToolStore`, `MemoryStore`, `TracingStore`, `APIKeyStore`, and `ConfigStore`.
  - Wrap database operations in Go standard `database/sql` transactions.
  - Implement uniform error wrapping to map SQL errors to domain errors.
- **Implementation Guidance:** Refer to GoClaw's Store Layer abstractions. Wire implementation behind interfaces using dependency inversion.
- **Dependencies:** Task 2.1.
- **Validation:** Write database integration tests (using test DB or containerized PostgreSQL) to verify CRUD transactions, rollbacks, and foreign key integrity.

### Phase Completion Criteria
* Database migrations run successfully up and down.
* Store tests verify complete persistence cycles of all domain models.

---

## 5. Phase 3 — Gateway Transport & Request Routing
### Objective
Expose the server entry points via HTTP, validate incoming Gateway API Tokens, and route requests to target Agent loops.

### Modules Covered
* `internal/transport/http`

### Task 3.1 — HTTP Server & Token Authentication Middleware
- **Description:** Launch the HTTP engine and protect endpoints with Gateway API Token verification.
- **Objective:** Handle request admission securely.
- **Requirements:**
  - Start an HTTP server (e.g. using standard `net/http` or `go-chi` router) on the configured port.
  - Implement `/health` endpoint returning database health status.
  - Write an authentication middleware that validates requests against `GATEWAY_API_TOKEN` configured in `.env`.
  - Require the `Authorization: Bearer <token>` HTTP header.
  - Reject unauthorized requests with HTTP 401. Secret tokens must never be logged.
- **Implementation Guidance:** Refer to GoClaw's token path in [20-api-keys-auth.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/20-api-keys-auth.md).
- **Dependencies:** Phase 2.
- **Validation:** Send requests to the server without a token, with an invalid token, and with a valid token, verifying 401 vs 200 responses.

### Task 3.2 — Agent Router & ExecutionRequest Construction
- **Description:** Implement resolution of incoming chat HTTP/WebSocket requests to actual agent profiles and sessions.
- **Objective:** Direct execution commands to the runtime engine.
- **Requirements:**
  - Implement the `AgentRouter` package.
  - Verify that the target Agent exists in `AgentStore` and is active.
  - Resolve the session ID in `SessionStore`, verifying it belongs to the target Agent.
  - Build the immutable `ExecutionRequest` structure from the client request payload.
- **Implementation Guidance:** Refer to GoClaw's request validation flow in [04-gateway-protocol.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/04-gateway-protocol.md).
- **Dependencies:** Task 3.1.
- **Validation:** Mock store responses and test routing under agent/session mismatches or non-existent agent errors.

### Phase Completion Criteria
* Health checks and token auth validation are covered.
* Incoming requests are successfully translated into `ExecutionRequest` containers.

---

## 6. Phase 4 — Provider Abstraction & Model Execution
### Objective
Establish the LLM provider interfaces and implement the OpenAI-Compatible HTTP/SSE streaming connection.

### Modules Covered
* `internal/impl/provider`
* `internal/impl/provider/openai`

### Task 4.1 — Provider Abstraction Interface
- **Description:** Define the Go interface for model completion, listing, and streaming.
- **Objective:** Decouple the Agent Loop from specific provider APIs.
- **Requirements:**
  - Create the `Provider` interface.
  - Include methods for: `GenerateCompletion`, `StreamCompletion`, `ListModels`, and `ValidateKey`.
  - Support return types for tool calls, usage tokens, and mapped Go standard errors.
- **Implementation Guidance:** Keep provider implementations isolated. See GoClaw's Provider Architecture in [02-providers.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/02-providers.md).
- **Dependencies:** Phase 1.
- **Validation:** Compile mock providers verifying no dependency on HTTP-specific SDKs.

### Task 4.2 — OpenAI-Compatible Adapter Implementation
- **Description:** Write the concrete provider implementation supporting OpenAI compatible APIs.
- **Objective:** Allow execution on OpenRouter, Gemini, OpenAI, DeepSeek, and Vercel endpoints.
- **Requirements:**
  - Implement the `Provider` interface for OpenAI standard endpoints.
  - Support Server-Sent Events (SSE) streaming.
  - Parse tool call fragments from streaming chunks and assemble them into a cohesive `ToolCall` slice.
  - Implement model listing (`ListModels`) by parsing response schemas.
  - Translate provider HTTP errors into canonical system errors.
- **Implementation Guidance:** Refer to GoClaw's provider stream parsing in [02-providers.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/02-providers.md).
- **Dependencies:** Task 4.1.
- **Validation:** Write mock HTTP server integration tests to verify successful chunk parsing, SSE EOF handling, and tool call reconstruction.

### Phase Completion Criteria
* Abstract LLM provider logic compiles.
* Local unit tests successfully parse streaming completion chunks and lists.

---

## 7. Phase 5 — Skills & Rules Systems
### Objective
Build the local Markdown YAML frontmatter parser for Skills and rule verification sandboxes.

### Modules Covered
* `agent/skills`
* `agent/rules`
* `internal/impl/...` (skill/rule loading + sandbox)

### Task 5.1 — Local Skills Filesystem Loader & Parser
- **Description:** Build the system that scans, reads, and parses local Skill definitions.
- **Objective:** Expose natural language instruction packages.
- **Requirements:**
  - Read files from the configured local filesystem path (e.g. `agent/skills/`).
  - Parse Skill metadata from the YAML frontmatter (name, description, version, dependencies).
  - Extract the main instruction body (Markdown) to inject directly into the Agent's System Prompt.
  - Reject skills with syntax errors or invalid frontmatter without crashing.
- **Implementation Guidance:** Refer to GoClaw's core skill resolution in [15-core-skills-system.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/15-core-skills-system.md).
- **Dependencies:** Phase 2.
- **Validation:** Put custom mock markdown files with frontmatter inside a test directory, scan, and verify metadata maps correctly to `internal/domain/entity/Skill`.

### Task 5.2 — Rules CRUD & Sandbox Executor
- **Description:** Implement CRUD actions for system Rules and a temporary sandbox to test Rules.
- **Objective:** Manage behavioral constraints and test them on the fly.
- **Requirements:**
  - Build Rule management API (create, edit, delete, assign to Agent).
  - Implement Rule Sandbox: Create a mock session containing the rule, construct a mock Runtime Context, send a prompt to the LLM Provider, and stream back the response without saving the history or events to the main DB.
- **Implementation Guidance:** Reference GoClaw's sandbox pattern. Ensure isolation from production sessions.
- **Dependencies:** Phase 4, Task 5.1.
- **Validation:** Run sandbox testing calls, verifying no production databases/sessions are touched or modified.

### Phase Completion Criteria
* Skill Markdown files are successfully discovered and parsed.
* Sandbox executes live LLM provider calls with isolated context.

---

## 8. Phase 6 — Tools System & Security Boundaries
### Objective
Define the Tool interface and build the centralized Tool Registry that enforces permissions and workspace directory safety constraints.

### Modules Covered
* `internal/impl/tool`

### Task 6.1 — Tool Registry & Contracts
- **Description:** Define the programmatic interface for all tools and the registry manager.
- **Objective:** Maintain a catalog of execution tools.
- **Requirements:**
  - Define `Tool` interface containing `Metadata` (name, description, JSON input schema) and `Execute(ctx, args) (ToolResult, error)`.
  - Implement a central `Registry` to handle registration, listing, schema generation for LLM ingestion, and lookup.
- **Implementation Guidance:** Refer to GoClaw's tools structure in [03-tools-system.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/03-tools-system.md).
- **Dependencies:** Phase 1.
- **Validation:** Write unit tests registering tools and verifying duplicate registration is rejected.

### Task 6.2 — Path Containment & Built-in File / CLI Tools
- **Description:** Build the core execution tools with strict security constraints.
- **Objective:** Execute filesystem commands safely.
- **Requirements:**
  - Build `read_file` and `write_file` tools.
  - **Path Safety:** Enforce that all file operations occur strictly within the workspace directory. Resolve absolute paths and block path traversal attempts (`..`, symlinks pointing outside).
  - Implement confirmation checks for modifying operations.
- **Implementation Guidance:** See GoClaw's Path safety checks and Direct Exec mode in [19-credentialed-exec.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/19-credentialed-exec.md).
- **Dependencies:** Task 6.1.
- **Validation:** Write unit tests verifying that trying to read or write a file outside the workspace root returns a permission error.

### Task 6.3 — Centralized Tool Authorization Policy
- **Description:** Enforce policies inside the registry before permitting execution.
- **Objective:** Guard execution limits at the registry layer.
- **Requirements:**
  - Implement registry lookup and permission check: Is the tool globally enabled? Is the tool assigned/allowed for the target Agent? Are arguments valid against the JSON Schema?
  - Return `ToolError` on unauthorized or invalid calls to feed back to the LLM loop.
- **Implementation Guidance:** Check authorization checks in [23-ai-agent-permission-matrix.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/23-ai-agent-permission-matrix.md).
- **Dependencies:** Task 6.2, Phase 4.
- **Validation:** Test registry pipeline execution with unauthorized agents, invalid schema arguments, and disabled tools.

### Phase Completion Criteria
* Built-in file tools block path traversal.
* Registry successfully validates schemas and rejects unauthorized calls.

---

## 9. Phase 7 — Agent Loop Execution Engine
### Objective
Implement the main orchestrator (Agent Loop), context snapshot builder, and token-based conversation pruning/compaction.

### Modules Covered
* `internal/domain/service` (Agent Loop)
* `internal/domain/service/compaction`

### Task 7.1 — Runtime Context Snapshot Builder
- **Description:** Build the snapshot assembly service that collects all state at the start of an execution.
- **Objective:** Ensure execution determinism.
- **Requirements:**
  - Assemble: Agent identity/SOUL, active Rules, active Skills, Session history, memory vectors, tool permissions, and provider model config.
  - Return a completely immutable `RuntimeContext` snapshot.
- **Implementation Guidance:** Snapshot all variables to verify changes to DB records mid-run do not mutate active context.
- **Dependencies:** Phase 3, Phase 5, Phase 6.
- **Validation:** Assert that modifying rule database tables during runtime does not alter the current Execution's context.

### Task 7.2 — Agent Loop State Machine & Execution Loop
- **Description:** Build the main orchestrator loop coordinating model completions and tool executions.
- **Objective:** Drive the Agent's reasoning loop.
- **Requirements:**
  - Implement a bounded execution loop: Prepare context -> Call LLM -> Parse Output (Final Text or ToolCall) -> If ToolCall, send to Tool Registry -> Obtain ToolResult/ToolError -> Add to history context -> Call LLM again.
  - Enforce bounds: max tool calls (budget), max durations, and max depth to prevent runaway execution.
  - Support cancellation contexts for live cancellation.
- **Implementation Guidance:** Refer to V2 execution patterns in [01-agent-loop.md](file:///home/cloud/workspace/project/01-agent-loop.md).
- **Dependencies:** Task 7.1, Phase 4.
- **Validation:** Write mock tests verifying that the loop stops when max tool call limits are hit, or when final text is returned.

### Task 7.3 — Pre-Execution Compaction & Context Budgeting
- **Description:** Check token bounds before model calls and prune or compact history.
- **Objective:** Fit context windows correctly.
- **Requirements:**
  - Estimate tokens for history, system prompt, and context.
  - If token counts exceed limits, compact the history: preserve recent messages, calculate summaries for old messages, and update the compaction checkpoint.
  - Keep related tool calls and results grouped together to avoid orphan calls.
- **Implementation Guidance:** Refer to GoClaw compaction and history context budgeting in [14-skills-runtime.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/14-skills-runtime.md).
- **Dependencies:** Task 7.2.
- **Validation:** Write unit tests feeding excessive message histories to verify they are compacted down before provider calls.

### Phase Completion Criteria
* Agent Loop compiles and correctly coordinates provider calls and tool runs.
* Context bounds are enforced and compacted when threshold limits are exceeded.

---

## 10. Phase 8 — Sessions & Vector Memory
### Objective
Implement database session state, local file uploads, and Go-based Memory vector similarity search.

### Modules Covered
* `internal/impl/session`
* `internal/impl/memory`

### Task 8.1 — Session Management & History Persistence
- **Description:** Write conversation turn logs to PostgreSQL at the end of executions.
- **Objective:** Save conversation state.
- **Requirements:**
  - Persist messages: user text, assistant text, tool calls, and tool results.
  - Save history sequentially.
  - Ensure persistence failures do not rollback already streamed responses.
- **Implementation Guidance:** See GoClaw's store models.
- **Dependencies:** Phase 2, Phase 7.
- **Validation:** Verify session loading returns identical message sequences.

### Task 8.2 — Local File Uploads Storage
- **Description:** Implement attachment handling for incoming client requests.
- **Objective:** Store binary files locally.
- **Requirements:**
  - Build upload endpoint saving files to a local directory (e.g. `data/attachments/`).
  - Store metadata and filepath in the database.
  - Expose static file server router endpoint to serve these files.
- **Implementation Guidance:** Abstract the file saver via Go interfaces so we can easily swap to Cloudinary later.
- **Dependencies:** Phase 3.
- **Validation:** Upload an image via HTTP, verify it is written to the local disk, and retrieve it via the static endpoint.

### Task 8.3 — Go-Based Cosine Similarity Vector Memory
- **Description:** Build vector embedding loading and similarity scoring in Go.
- **Objective:** Implement local-first vector search.
- **Requirements:**
  - Build embedding generation logic calling the configured embedding provider.
  - **Similarity Matching:** Read memory vector records from PostgreSQL (stored as float arrays/JSONB). Calculate Cosine Similarity directly in Go memory.
  - Sort candidate chunks and select the top K elements.
- **Implementation Guidance:** Refer to GoClaw's memory retrieval in [07-bootstrap-skills-memory.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/07-bootstrap-skills-memory.md). Instead of using PostgreSQL `pgvector` queries, perform the score matching inside Go.
- **Dependencies:** Phase 4, Phase 2.
- **Validation:** Inject sample vectors into PostgreSQL, search for similar vectors in Go, and assert that the closest vectors are correctly ranked.

### Phase Completion Criteria
* Message histories are stored.
* Files are uploaded and accessible.
* Cosine Similarity matches vector coordinates inside Go correctly.

---

## 11. Phase 9 — Observability & Realtime Events
### Objective
Publish live execution events over WebSocket/SSE and log telemetry spans.

### Modules Covered
* `internal/transport/socket` (realtime events)
* `internal/impl/store/postgres` (tracing store)

### Task 9.1 — Realtime Event Publisher
- **Description:** Stream agent runtime events to clients.
- **Objective:** Enable interactive visualization of reasoning steps.
- **Requirements:**
  - Build events interface supporting: `execution.started`, `skill.activated`, `reasoning.step`, `tool.requested`, `tool.started`, `tool.completed`, `tool.failed`, `response.delta`, `execution.completed`, `execution.failed`.
  - Expose WebSocket or SSE endpoint.
  - Stream events asynchronously. Event delivery failures must not crash the Agent Loop.
- **Implementation Guidance:** Refer to GoClaw websocket lifecycle in [04-gateway-protocol.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/04-gateway-protocol.md).
- **Dependencies:** Phase 3, Phase 7.
- **Validation:** Connect a websocket client, run an agent loop, and verify receipt of standard lifecycle events.

### Task 9.2 — Tracing & Metrics Telemetry
- **Description:** Implement trace spans for debugging and latency tracking.
- **Objective:** Measure system performance.
- **Requirements:**
  - Create root execution trace containing spans for Gateway, Router, Context Build, Model Call, Tool execution, and Finalization.
  - Measure token usage, duration, and error occurrences.
  - Save tracing metadata asynchronously to `TracingStore`.
- **Implementation Guidance:** See GoClaw tracing architecture in [10-tracing-observability.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/10-tracing-observability.md).
- **Dependencies:** Phase 2, Phase 7.
- **Validation:** Query trace logs via store after a run and verify all expected child spans exist and match.

### Phase Completion Criteria
* Realtime client receives progressive event chunks.
* Database traces record spans of executed sessions.

---

## 12. Phase 10 — API Keys & Global Configurations
### Objective
Implement multi-scope API keys and global system configuration endpoints.

### Modules Covered
* `internal/impl/store/postgres` (api_key, global_config)
* `internal/impl/config`

### Task 10.1 — API Keys Authentication
- **Description:** Implement CRUD and verification for client integration API keys.
- **Objective:** Allow external applications to call the gateway.
- **Requirements:**
  - Implement API Key generation (show plaintext key only once on creation, save salted hash in DB).
  - Define scopes (e.g. `chat`, `read:agents`, `write:agents`).
  - Write authentication check checking scopes and key expiration.
- **Implementation Guidance:** Map auth paths to match details in [20-api-keys-auth.md](file:///home/cloud/workspace/project/cloudclaw/docs/references/20-api-keys-auth.md).
- **Dependencies:** Phase 3.
- **Validation:** Test invoking agent chat HTTP endpoints using a generated API Key, verifying scope denials.

### Task 10.2 — Global Configuration Management
- **Description:** Build global settings storage and runtime update propagation.
- **Objective:** Manage server-wide default parameters.
- **Requirements:**
  - Support settings: Host, Port, default Provider, default Model, and Quota parameters.
  - Ensure that updates to global settings do not alter running execution contexts (snapshot boundary).
- **Dependencies:** Phase 2, Phase 7.
- **Validation:** Save configuration parameters, verify agents utilize updated default models on the next run.

### Phase Completion Criteria
* API keys validate client requests with specific scope checks.
* Global configurations persist and propagate correctly.

---

## 13. Phase 11 — Terminal CLI
### Objective
Build the Cobra command tree and BubbleTea interactive chat console for command-line control.

### Modules Covered
* `cmd/agentkit`
* `internal/tui`

### Task 11.1 — Cobra Command Tree & Server Boot
- **Description:** Implement the main CLI entry points.
- **Objective:** Run CLI actions.
- **Requirements:**
  - Create commands: `serve`, `chat`, `agent`, `skill`, `session`, `config`, and `api-key`.
  - Bind CLI commands to backend container/API clients.
- **Implementation Guidance:** Keep CLI clean and delegate to API layer.
- **Dependencies:** Phase 1, Phase 3.
- **Validation:** Run `./agentkit --help` and verify all commands display correctly.

### Task 11.2 — Interactive CLI Chat
- **Description:** Implement console chat interface using Go BubbleTea.
- **Objective:** Provide a fast CLI chat client.
- **Requirements:**
  - Implement interactive prompt supporting command flags (`/agent`, `/session`, `/tool`).
  - Stream tokens directly to the terminal.
  - Print tool execution sequences clearly.
- **Implementation Guidance:** Use Go standard BubbleTea library for rendering state.
- **Dependencies:** Task 11.1, Phase 9.
- **Validation:** Run `./agentkit chat --agent default`, send a message, and check streaming output.

### Phase Completion Criteria
* CLI builds as a single binary.
* Console chat communicates with LLM provider.

---

## 14. Phase 12 — React Dashboard Frontend
### Objective
Implement the React + TS Dashboard application, onboarding screens, and agent management controls.

### Modules Covered
* `web/`

### Task 12.1 — React Foundation & Shared State
- **Description:** Set up the Web application shell, state management, and HTTP/WebSocket clients.
- **Objective:** Initialize the frontend app.
- **Requirements:**
  - Scaffold React with TypeScript.
  - Implement Gateway token check saving token in LocalStorage.
  - Create API client and WebSocket connection lifecycle wrapper.
- **Dependencies:** Phase 3, Phase 9.
- **Validation:** Load React app in browser, log in with Gateway Token, and check WebSocket connection status.

### Task 12.2 — Onboarding Flow Setup
- **Description:** Implement the onboarding wizard.
- **Objective:** Guide first-time user setups.
- **Requirements:**
  - Create 3-step wizard screen:
    1. Enter Gateway Token.
    2. Add Provider details (API Key) & choose Model.
    3. Configure first Agent (edit `SOUL.md` / `IDENTITY.md`).
- **Dependencies:** Task 12.1, Phase 4, Phase 2.
- **Validation:** Walk through steps, verify successful first agent generation in DB, and redirect to main Dashboard.

### Task 12.3 — Chat & Management Panels
- **Description:** Build the main chat dashboard and configurations panels.
- **Objective:** Interact with agents visually.
- **Requirements:**
  - Implement chat layout (session history list, message input, tool execution timeline).
  - Implement Agent settings panels (create, edit, clone agent, select provider/model, toggle skills, rules, and allowed tools).
  - Implement Skill & Rule view panels, Rule Sandbox test form, API Key manager, and Config panels.
- **Dependencies:** Task 12.2, Phase 5, Phase 6, Phase 8.
- **Validation:** Verify streaming responses render, tool calls are visible, and Rule Sandbox runs successfully in browser.

### Phase Completion Criteria
* React Dashboard builds and operates without runtime errors.
* Onboarding flow successfully creates agents from empty database state.

---

## 15. Phase 13 — Reliability, Fault Isolation & E2E Validation
### Objective
Enforce robust failure handling, secure workspace containment, and execute full end-to-end scenarios.

### Modules Covered
* Cross-cutting

### Task 13.1 — Bounded Retries & Fault Isolation
- **Description:** Build retry wrappers and isolate execution crashes from side-effect components.
- **Objective:** Maximize runtime reliability.
- **Requirements:**
  - Implement transient error detection (retry HTTP calls for rate limits or server errors up to `max provider retries`).
  - Isolate memory, tracing, database write, and WebSocket failures. If they crash, log the error but do not disrupt message streaming to the client.
- **Dependencies:** Phase 7, Phase 8, Phase 9.
- **Validation:** Inject mock database write crashes during execution and verify that the client still receives the streamed response.

### Task 13.2 — End-to-End Product Validation
- **Description:** Run E2E scenarios verification.
- **Objective:** Validate the final product.
- **Requirements:**
  - Automate complete scenario: Start with empty DB -> Onboarding setup -> Create agent -> Chat with agent -> Agent invokes `read_file` -> Compaction triggered -> Verify history persists.
- **Dependencies:** Phase 11, Phase 12.
- **Validation:** Execute standard E2E test runs verifying all state invariants survive server restarts.

### Phase Completion Criteria
* Side-effect failures do not interrupt chat streams.
* Complete onboarding to chat sequence is validated successfully.

---

## 16. Final Integration
The final integration step wires the CLI and Dashboard applications to the single HTTP/WebSocket Gateway server. The system boots dynamically from database credentials. Config variables override default behaviors, and memory embeddings search operates locally via Cosine Similarity calculations inside the Go application engine.

---

## 17. Final Definition of Done
A feature or phase is considered **Done** when:
1. **Compilation:** `go build ./...` compiles cleanly and frontend React build runs without errors.
2. **Quality Checks:** `golangci-lint run` passes with zero findings, and files are formatted via `gofmt` and `goimports`.
3. **Automated Testing:** `go test -race -cover ./...` passes.
4. **Acceptance Criteria:** Every user story matches all acceptance criteria (including visual browser checks for frontend components).
5. **Security:** Path containment blocks file tool directory traversals. Secrets are redacted from traces and log logs.
6. **Documentation:** Local setup configurations and API endpoints are fully documented in `docs/`.
