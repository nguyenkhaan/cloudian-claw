# PROJECT ARCHITECTURE DESIGN

## 1. System Design

### 1.1 Objective

Define the architectural blueprint, component responsibilities, runtime execution model, request flow, data flow, failure behavior, and implementation boundaries of the local AI Agent Gateway system.

The architecture is designed around a single core principle: 

> **Agent Loop is the central execution orchestrator of the Agent Runtime.**

The Agent Loop owns the execution lifecycle of an Agent invocation, while specialized components provide external capabilities and persistence:

* Provider Abstraction communicates with LLM providers.
* Tool Registry validates permissions and executes tools.
* Skill and Rule components provide agent behavior and capabilities.
* Memory provides context retrieval and post-execution memory processing.
* Session Store provides durable conversation state.
* Realtime Event component publishes execution events.
* Tracing Store records execution telemetry.
* PostgreSQL provides durable persistence.
* Gateway and Agent Router handle request admission and routing before execution.

The system is local-first and single-owner. It does not include multi-tenant organization management or multi-agent orchestration.

It is like an AI Agent harness. That can turn an LLM into an Agent can execute tasks 

---

## 2. Architectural Scope

### 2.1 In Scope

The system provides the following core capabilities:

* Gateway authentication using a Gateway API Token
* Provider registration and API key management
* Provider model discovery and validation
* Agent creation and lifecycle management
* Agent configuration through SOUL.md / IDENTITY.md
* Skill loading, validation, versioning, and assignment
* Rule creation, validation, versioning, assignment, and sandbox testing
* Tool registration and agent-specific access policies
* Agent execution through the Agent Loop
* LLM provider abstraction
* Tool calling and bounded tool execution loops
* Session creation, persistence, restoration, and history
* Memory retrieval
* Memory write and embedding
* Automatic and manual conversation compaction
* Realtime execution events
* Streaming model responses
* Execution tracing and telemetry
* API key creation and usage tracking
* Global server / AI / quota / tool configuration
* Terminal CLI interaction
* React Dashboard interaction

### 2.2 Out of Scope

The following are explicitly excluded:

* Channel adapters
* Slack / Discord / Telegram integrations
* MCP protocol handling
* Subagent orchestration
* Agent-to-agent delegation
* Teams / Organizations / multi-tenancy
* TTS
* Webhook dispatcher
* Cron scheduler
* Team-based quota management
* Agent hierarchy such as parent/child agents

---

# 3. System Architecture

## 3.1 Technology Stack

| Layer                    | Technology                       | Responsibility / Rationale                                                              |
| ------------------------ | -------------------------------- | --------------------------------------------------------------------------------------- |
| Backend                  | Go (Golang)                      | High performance, fast startup, concurrent request handling                             |
| Frontend                 | React + TypeScript + TailwindCSS | Component reusability, type safety, rapid UI development                                |
| Database                 | PostgreSQL                       | Durable relational state, data integrity, JSONB support                                 |
| Cache                    | In-Memory RAM + Redis            | Low-latency session/hot-data access                                                     |
| Realtime                 | WebSocket / SSE                  | Streaming responses and realtime execution events                                       |
| CLI                      | Cobra / BubbleTea                | Interactive terminal experience                                                         |
| LLM Providers            | OpenAI-Compatible HTTP + SSE     | Support providers such as OpenAI, Gemini, DeepSeek, DashScope, and compatible endpoints |
| Agent Protocol Providers | ACP / JSON-RPC 2.0 / stdio       | Support compatible local agent CLI providers                                            |
| Skills                   | Local Filesystem                 | Load and validate local skills                                                          |
| Agent Configuration      | Filesystem + persistent metadata | Store SOUL.md / IDENTITY.md and agent configuration                                     |

---

## 3.2 Architectural Overview

The architecture is divided into four logical layers:

1. **Client Layer**
2. **Gateway / Routing Layer**
3. **Agent Runtime Layer**
4. **Persistence / External Capability Layer**

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│                              CLIENT LAYER                                    │
│                                                                              │
│   ┌───────────────────────┐              ┌──────────────────────────────┐   │
│   │      Terminal CLI     │              │      React Dashboard         │   │
│   │    Cobra / BubbleTea  │              │ React + TypeScript + Tailwind│   │
│   └───────────┬───────────┘              └──────────────┬───────────────┘   │
│               │                                         │                    │
└───────────────┼─────────────────────────────────────────┼────────────────────┘
                │                                         │
                └───────────────────┬─────────────────────┘
                                    ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                           GATEWAY / ROUTING LAYER                            │
│                                                                              │
│   ┌────────────────┐       ┌────────────────┐       ┌────────────────────┐  │
│   │   HTTP Server  │──────▶│  Auth / Request│──────▶│    Agent Router    │  │
│   │ Socket Server  │       │   Validation   │       │                    │  │
│   └────────────────┘       └────────────────┘       └─────────┬──────────┘  │
│                                                               │             │
└───────────────────────────────────────────────────────────────┼─────────────┘
                                                                │
                                                                ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                              AGENT RUNTIME                                   │
│                                                                              │
│                     ┌──────────────────────────────┐                         │
│                     │         AGENT LOOP           │                         │
│                     │                              │                         │
│                     │  Execution Orchestrator      │                         │
│                     │                              │                         │
│                     │  1. Resolve Runtime Context  │                         │
│                     │  2. Build Model Context      │                         │
│                     │  3. Call Provider            │                         │
│                     │  4. Process Tool Calls       │                         │
│                     │  5. Continue Execution Loop  │                         │
│                     │  6. Stream Final Response    │                         │
│                     │  7. Finalize Execution       │                         │
│                     └─────────────┬────────────────┘                         │
│                                   │                                          │
│             ┌─────────────────────┼─────────────────────────────┐            │
│             │                     │                             │            │
│             ▼                     ▼                             ▼            │
│   ┌──────────────────┐  ┌─────────────────────┐     ┌───────────────────┐  │
│   │ Provider         │  │    Tool Registry    │     │ Runtime Context   │  │
│   │ Abstraction      │  │                     │     │ Sources           │  │
│   │                  │  │ Permission Policy   │     │                   │  │
│   │ OpenAI-compatible│  │ Validation          │     │ Agent Config      │  │
│   │ ACP              │  │ Execution           │     │ Rules             │  │
│   └────────┬─────────┘  └──────────┬──────────┘     │ Skills            │  │
│            │                       │                │ Memory            │  │
│            │                       │                │ Session           │  │
│            │                       │                │ Tool Permissions  │  │
│            │                       │                └───────────────────┘  │
│            │                       │                                          │
│            └───────────────┬───────┘                                          │
│                            │                                                  │
│                            ▼                                                  │
│                 ┌────────────────────────┐                                   │
│                 │   REALTIME EVENTS       │                                   │
│                 │                        │                                   │
│                 │ execution / reasoning  │                                   │
│                 │ skill / tool / error   │                                   │
│                 └────────────┬───────────┘                                   │
│                              │                                               │
└──────────────────────────────┼───────────────────────────────────────────────┘
                               │
                               ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                         PERSISTENCE / STATE LAYER                            │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐  │
│   │                         STORE LAYER                                  │  │
│   │                                                                     │  │
│   │ Session Store │ Agent Store │ Provider Store │ Skill / Rule Store   │  │
│   │ Memory Store  │ Tool Store  │ Tracing Store │ API Key Store         │  │
│   └───────────────────────────────────────┬─────────────────────────────┘  │
│                                           │                                  │
│                                           ▼                                  │
│                                  ┌────────────────┐                          │
│                                  │  PostgreSQL     │                          │
│                                  │   Database      │                          │
│                                  └────────────────┘                          │
│                                                                              │
│   Local Filesystem: Skills / SOUL.md / IDENTITY.md                          │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 3.3 Architectural Boundaries

### Gateway

The Gateway is responsible for:

* receiving HTTP/WebSocket requests
* authenticating the Gateway API Token
* validating request structure
* applying request-level protection such as rate limiting
* routing requests to the correct Agent execution entry point

The Gateway does **not** execute agent reasoning or tools.

---

### Agent Router

The Agent Router is responsible for:

* resolving the target Agent
* validating that the requested Agent exists and is executable
* resolving the execution target from request/session context
* constructing the `ExecutionRequest`
* forwarding the request to the Agent Loop

The Agent Router does **not** own the model/tool loop.

---

### Agent Loop

The Agent Loop is the **central execution orchestrator**.

It owns:

* execution lifecycle
* Runtime Context resolution
* model invocation
* tool-call processing
* bounded execution loop
* execution state
* realtime execution event emission
* final response streaming
* execution finalization

The Agent Loop does **not**:

* implement provider-specific APIs
* execute individual tools
* directly enforce tool permissions
* own persistent session storage
* implement memory retrieval algorithms
* own PostgreSQL persistence directly

Those responsibilities belong to specialized dependencies.

---

## 3.4 Agent Runtime Context

Every execution creates a **Runtime Context Snapshot**.

The snapshot represents all information that is authoritative for that execution.

```text
Agent Runtime Context
│
├── Agent Configuration
│   ├── SOUL.md
│   ├── IDENTITY.md
│   └── Provider / Model selection
│
├── Rules
│   └── Rules assigned to the Agent
│
├── Skills
│   └── Active and validated Skills
│
├── Session
│   └── Conversation history
│
├── Memory
│   └── Relevant retrieved memory/context
│
└── Tool Permissions
    └── Tools allowed for this Agent
```

Once the Runtime Context Snapshot has been created:

> Changes to Agent configuration, Skill, Rule, Memory configuration, or Tool policy do not mutate the execution already in progress.

Changes become effective from the **next execution**.

This provides deterministic execution boundaries.

---

# 4. Component Specifications

Each component follows:

**Purpose → Responsibilities → Communicates With → Input → Output → Failure Behavior → Boundary**

---

## 4.1 Component: Gateway

| Property              | Description                                                                      |
| --------------------- | -------------------------------------------------------------------------------- |
| **Purpose**           | Entry point for client and external API requests                                 |
| **Responsibilities**  | Authentication, request validation, routing, rate limiting, connection handling  |
| **Communicates With** | HTTP Server, Socket Server, Agent Router, API Key Store                          |
| **Input**             | HTTP / WebSocket request                                                         |
| **Output**            | Validated internal request or authentication/validation error                    |
| **Failure Behavior**  | Rejects invalid or unauthorized requests; emits/logs security and request errors |
| **Boundary**          | Does not execute Agent logic                                                     |

---

## 4.2 Component: Agent Router

| Property              | Description                                                         |
| --------------------- | ------------------------------------------------------------------- |
| **Purpose**           | Resolve the Agent execution target                                  |
| **Responsibilities**  | Agent lookup, request normalization, execution request construction |
| **Communicates With** | Gateway, Agent Store, Session Store, Agent Loop                     |
| **Input**             | User request, Agent ID/name, Session ID                             |
| **Output**            | `ExecutionRequest`                                                  |
| **Failure Behavior**  | Rejects missing Agent, invalid Session, or invalid request context  |
| **Boundary**          | Does not own execution loop                                         |

### ExecutionRequest

Conceptually:

```text
ExecutionRequest
├── Agent ID
├── Session ID
├── User Message
├── Attachments / Files
├── Request Metadata
└── Execution Options
```

---

## 4.3 Component: Agent Service

| Property              | Description                                                                           |
| --------------------- | ------------------------------------------------------------------------------------- |
| **Purpose**           | Manage Agent lifecycle and configuration                                              |
| **Responsibilities**  | Create, update, delete, clone, validate Agent configuration                           |
| **Communicates With** | Agent Store, Provider Store, Skill Store, Rule Store, Tool Store                      |
| **Input**             | Agent configuration                                                                   |
| **Output**            | Persisted Agent definition                                                            |
| **Failure Behavior**  | Reject invalid configuration and retain previous valid configuration where applicable |
| **Boundary**          | Control-plane component; does not execute the Agent                                   |

### Agent Configuration

```text
Agent
├── SOUL.md
├── IDENTITY.md
├── Provider
├── Model
├── Skills
├── Rules
├── Allowed Tools
└── Memory Configuration
```

---

## 4.4 Component: Agent Loop

| Property              | Description                                                                                                                          |
| --------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| **Purpose**           | Orchestrate one complete Agent execution                                                                                             |
| **Responsibilities**  | Build Runtime Context, invoke Provider, process tool calls, enforce execution bounds, stream output, emit events, finalize execution |
| **Communicates With** | Provider Abstraction, Tool Registry, Skill Store/Loader, Rule Store, Memory, Session Store, Realtime Events, Tracing Store           |
| **Input**             | `ExecutionRequest`                                                                                                                   |
| **Output**            | Streaming response + execution result                                                                                                |
| **Failure Behavior**  | Applies failure policy according to error class; logs every failure and emits user-visible error events when appropriate             |
| **Boundary**          | Owns orchestration, not implementation details of tools/providers/stores                                                             |

### Agent Loop Lifecycle

```text
ExecutionRequest
      │
      ▼
Resolve Runtime Context
      │
      ▼
Build Model Context
      │
      ▼
Provider Generate / Stream
      │
      ├─────────────── Final Response ───────────────┐
      │                                               │
      └──────── Tool Call ──▶ Tool Registry           │
                              │                       │
                              ▼                       │
                         Tool Result / Error          │
                              │                       │
                              └──────▶ Provider ──────┘
                                     
                          repeat until
                        terminal execution state
```

---

## 4.5 Component: Runtime Context Builder

Runtime Context Builder is an internal responsibility of Agent Loop rather than an independent runtime service.

### Purpose

Construct the immutable Runtime Context Snapshot used by one execution.

### Inputs

* Agent configuration
* SOUL.md
* IDENTITY.md
* Rules
* Skills
* Session history
* Memory retrieval result
* Tool permissions
* Provider/model configuration

### Processing

```text
1. Load Agent configuration
2. Load active Rules
3. Load validated Skills
4. Load Session history
5. Check context length policy
6. Perform automatic compaction when necessary
7. Retrieve relevant Memory
8. Resolve Tool permissions
9. Construct Runtime Context Snapshot
```

### Output

```text
RuntimeContext
├── Agent Identity
├── System Instructions
├── Rules
├── Skills
├── Conversation Context
├── Retrieved Memory
├── Available Tools
└── Provider / Model Configuration
```

### Boundary

The Runtime Context Builder prepares the execution context but does not execute the model or tools.

---

## 4.6 Component: Provider Abstraction

| Property              | Description                                                                                                      |
| --------------------- | ---------------------------------------------------------------------------------------------------------------- |
| **Purpose**           | Provide a provider-neutral interface for model execution                                                         |
| **Responsibilities**  | Provider selection, request construction, streaming, response normalization, provider-specific error translation |
| **Communicates With** | Agent Loop, Provider Store, external AI Providers                                                                |
| **Input**             | Model context, model configuration                                                                               |
| **Output**            | Model stream, final response, normalized ToolCall, provider error                                                |
| **Failure Behavior**  | Retry transient errors within configured limit; terminate execution for permanent/configuration failures         |
| **Boundary**          | Does not execute tools or manage sessions                                                                        |

### Supported Interfaces

```text
Provider Abstraction
│
├── OpenAI-Compatible
│   ├── HTTP
│   └── SSE
│
└── ACP
    └── JSON-RPC 2.0 / stdio
```

Provider-specific implementations must be hidden behind the abstraction.

---

## 4.7 Component: Tool Registry

| Property              | Description                                                                                                    |
| --------------------- | -------------------------------------------------------------------------------------------------------------- |
| **Purpose**           | Register, authorize, validate, and execute Agent tools                                                         |
| **Responsibilities**  | Tool registration, tool discovery, permission policy enforcement, argument validation, execution, tool logging |
| **Communicates With** | Agent Loop, Tool Store, filesystem/process/network tools                                                       |
| **Input**             | `ToolCall`                                                                                                     |
| **Output**            | `ToolResult` or `ToolError`                                                                                    |
| **Failure Behavior**  | Returns structured ToolError rather than silently failing; logs failure and emits realtime event               |
| **Boundary**          | The Tool Registry is the authority for tool permission and execution                                           |

### Tool Call Boundary

```text
Agent Loop
    │
    │ ToolCall
    ▼
Tool Registry
    │
    ├── Tool exists?
    ├── Agent has permission?
    ├── Arguments valid?
    ├── Policy allows execution?
    │
    ├── YES ──▶ Execute Tool ──▶ ToolResult
    │
    └── NO ───▶ ToolError
```

The Agent Loop does **not** independently evaluate tool permissions.

---

## 4.8 Component: Skill System

| Property              | Description                                                                                  |
| --------------------- | -------------------------------------------------------------------------------------------- |
| **Purpose**           | Load and validate local Agent Skills                                                         |
| **Responsibilities**  | Filesystem loading, syntax validation, dependency validation, version handling, activation   |
| **Communicates With** | Agent Loop, filesystem, Skill Store                                                          |
| **Input**             | Skill configuration / filesystem path                                                        |
| **Output**            | Validated Skill definition                                                                   |
| **Failure Behavior**  | Reject invalid Skill and retain the last valid active version where applicable               |
| **Boundary**          | Skill content defines capabilities/instructions but does not directly execute the Agent loop |

### Skill Lifecycle

```text
Filesystem
   │
   ▼
Load Skill
   │
   ▼
Validate
   │
   ├── Invalid → reject
   │
   └── Valid
        │
        ▼
Version / Dependency Check
        │
        ▼
Assign to Agent
        │
        ▼
Available in next execution
```

---

## 4.9 Component: Rule System

| Property              | Description                                                                  |
| --------------------- | ---------------------------------------------------------------------------- |
| **Purpose**           | Define behavioral constraints for an Agent                                   |
| **Responsibilities**  | CRUD, validation, versioning, assignment, sandbox testing                    |
| **Communicates With** | Agent Service, Agent Loop, Rule Store                                        |
| **Input**             | Rule definition and assignment                                               |
| **Output**            | Validated Rule set                                                           |
| **Failure Behavior**  | Invalid rules are rejected; failed updates do not replace the active version |
| **Boundary**          | Rules become part of Runtime Context and are applied before model execution  |

Rules are not an independent reasoning engine.

They are contextual constraints applied during Runtime Context construction.

---

## 4.10 Component: Memory

| Property              | Description                                                                                                                       |
| --------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| **Purpose**           | Retrieve, store, embed, and compact long-term conversational context                                                              |
| **Responsibilities**  | Retrieval, embedding, memory write, context compaction                                                                            |
| **Communicates With** | Agent Loop, Session Store, Memory Store, embedding provider                                                                       |
| **Input**             | Session history, query/context, Memory configuration                                                                              |
| **Output**            | Relevant memory context or post-execution memory records                                                                          |
| **Failure Behavior**  | Retrieval failure degrades execution when possible; memory persistence failure is logged without rolling back a streamed response |
| **Boundary**          | Memory does not control Agent execution                                                                                           |

### Memory Read Path

```text
Execution begins
      │
      ▼
Session History
      │
      ▼
Context Policy Check
      │
      ├── Over threshold
      │       │
      │       ▼
      │   Automatic Compaction
      │
      ▼
Memory Retrieval
      │
      ▼
Relevant Memory Context
      │
      ▼
Runtime Context
```

### Memory Write Path

```text
Execution complete
      │
      ▼
Persist Session
      │
      ▼
Memory Processing
      │
      ├── Store memory
      ├── Generate embeddings
      └── Apply compaction policy when required
```

### Manual Compaction

Users can explicitly request compaction for a Session through Dashboard or CLI.

```text
User
 │
 ▼
Compact Session Request
 │
 ▼
Memory Component
 │
 ▼
Compaction Policy
 │
 ▼
Compacted Session / Context
```

Automatic and manual compaction use the same Memory policy:

* threshold
* keep recent messages
* maximum token budget

---

## 4.11 Component: Session Store

| Property              | Description                                                                               |
| --------------------- | ----------------------------------------------------------------------------------------- |
| **Purpose**           | Persist and restore conversation state                                                    |
| **Responsibilities**  | Create, update, archive, delete, restore sessions; persist messages                       |
| **Communicates With** | Agent Router, Agent Loop, Memory                                                          |
| **Input**             | Session commands and message records                                                      |
| **Output**            | Session state and conversation history                                                    |
| **Failure Behavior**  | Persistence failures are logged and reported, but do not rollback already streamed output |
| **Boundary**          | Stores state; does not decide Agent behavior                                              |

---

## 4.12 Component: Realtime Event System

| Property              | Description                                                                                                 |
| --------------------- | ----------------------------------------------------------------------------------------------------------- |
| **Purpose**           | Publish execution activity to the client in real time                                                       |
| **Responsibilities**  | Event emission, serialization, WebSocket/SSE delivery                                                       |
| **Communicates With** | Agent Loop, Frontend                                                                                        |
| **Input**             | Runtime execution events                                                                                    |
| **Output**            | WebSocket/SSE events                                                                                        |
| **Failure Behavior**  | Realtime delivery failure is logged; execution should not fail solely because event delivery is unavailable |
| **Boundary**          | Observability/UI delivery only; no execution decisions                                                      |

### Event Categories

```text
execution.started
skill.activated
reasoning.step
tool.requested
tool.started
tool.completed
tool.failed
response.delta
execution.completed
execution.failed
persistence.failed
memory.failed
```

Reasoning events are intended for internal execution visualization and do not require exposing hidden chain-of-thought content.

---

## 4.13 Component: Tracing / Monitoring

| Property              | Description                                                                              |
| --------------------- | ---------------------------------------------------------------------------------------- |
| **Purpose**           | Record execution telemetry and diagnose failures                                         |
| **Responsibilities**  | Trace creation, spans, latency, token usage, error recording, execution metadata         |
| **Communicates With** | Agent Loop, Provider Abstraction, Tool Registry, Store Layer                             |
| **Input**             | Execution and component telemetry                                                        |
| **Output**            | Persisted traces/metrics                                                                 |
| **Failure Behavior**  | Trace persistence failure is logged and must not invalidate a successful Agent execution |
| **Boundary**          | Observability only                                                                       |

### Tracing Scope

```text
Execution Trace
│
├── Request
├── Runtime Context creation
├── Memory retrieval
├── Provider call
├── Tool call(s)
├── Provider continuation(s)
├── Final response
├── Session persistence
└── Memory processing
```

---

## 4.14 Component: Store Layer

The Store Layer provides persistence abstractions and prevents runtime components from depending directly on PostgreSQL.

```text
Store Layer
│
├── Agent Store
├── Provider Store
├── Session Store
├── Skill Store
├── Rule Store
├── Tool Store
├── Memory Store
├── Tracing Store
└── API Key Store
```

The runtime layer communicates with Store interfaces rather than directly issuing database operations.

---

# 5. Request Flow Diagram

## 5.1 Standard Agent Execution Request Flow

```text
┌──────────────┐
│ User / FE /  │
│     CLI      │
└──────┬───────┘
       │
       │ User Message
       ▼
┌──────────────────────┐
│      Gateway         │
│ Auth + Validation    │
└──────────┬───────────┘
           │
           │ Validated Request
           ▼
┌──────────────────────┐
│    Agent Router      │
│ Resolve Agent/Session│
└──────────┬───────────┘
           │
           │ ExecutionRequest
           ▼
┌─────────────────────────────────────┐
│            Agent Loop               │
│                                     │
│  ┌───────────────────────────────┐  │
│  │ Build Runtime Context         │  │
│  │                               │  │
│  │ Agent Config                  │  │
│  │ SOUL / IDENTITY               │  │
│  │ Rules                         │  │
│  │ Skills                        │  │
│  │ Session History               │  │
│  │ Memory Retrieval              │  │
│  │ Tool Permissions              │  │
│  └───────────────┬───────────────┘  │
│                  │                  │
│                  ▼                  │
│         Provider Abstraction        │
│                  │                  │
│                  ▼                  │
│            Model Response           │
│                  │                  │
│          ┌───────┴────────┐         │
│          │                │         │
│      Final Text        ToolCall     │
│          │                │         │
│          │                ▼         │
│          │         ┌─────────────┐  │
│          │         │Tool Registry│  │
│          │         └──────┬──────┘  │
│          │                │         │
│          │        ┌───────┴───────┐ │
│          │        │               │ │
│          │      Allowed        Denied│
│          │        │               │ │
│          │        ▼               ▼ │
│          │    ToolResult       ToolError
│          │        │               │ │
│          │        └───────┬───────┘ │
│          │                │         │
│          │                ▼         │
│          │      Continue Provider   │
│          │      execution loop      │
│          │                │         │
│          └────────────────┘         │
│                  │                  │
└──────────────────┼──────────────────┘
                   │
                   ▼
        ┌────────────────────────┐
        │ Final Response Stream  │
        │ WebSocket / SSE        │
        └───────────┬────────────┘
                    │
                    ▼
              User / FE / CLI
```

---

## 5.2 Realtime Event Flow

Realtime events are side effects of the Agent Loop rather than execution stages.

```text
                    Agent Loop
                        │
        ┌───────────────┼────────────────┐
        │               │                │
        ▼               ▼                ▼
 execution.*       skill.* / tool.*   response.delta
        │               │                │
        └───────────────┼────────────────┘
                        ▼
                Realtime Event System
                        │
                        ▼
                 WebSocket / SSE
                        │
                        ▼
                   Dashboard
```

The Realtime Event System does not control execution.

If realtime delivery fails:

```text
Realtime delivery failure
        │
        ├── Log error
        ├── Record trace event
        └── Continue Agent execution
```

---

# 6. Data Flow Diagram — Agent Execution

The Data Flow describes the data being transformed throughout one execution rather than only the control sequence.

```text
┌──────────────────────┐
│   User Message       │
│   + Attachments      │
└──────────┬───────────┘
           │
           ▼
┌──────────────────────┐
│      Gateway         │
│ Authenticated Request│
└──────────┬───────────┘
           │
           ▼
┌──────────────────────┐
│    Agent Router      │
│  ExecutionRequest    │
└──────────┬───────────┘
           │
           ▼
┌──────────────────────────────────────┐
│       Runtime Context Resolution     │
│                                      │
│ Agent Config ───────────┐            │
│ SOUL / IDENTITY ────────┤            │
│ Rules ──────────────────┤            │
│ Skills ─────────────────┤            │
│ Session History ────────┤            │
│ Memory Retrieval ───────┤            │
│ Tool Permissions ───────┘            │
└────────────────┬─────────────────────┘
                 │
                 ▼
       ┌──────────────────────┐
       │ Runtime Context      │
       │ Snapshot             │
       └──────────┬───────────┘
                  │
                  ▼
       ┌──────────────────────┐
       │ Model Input Context  │
       └──────────┬───────────┘
                  │
                  ▼
       ┌──────────────────────┐
       │ Provider Abstraction │
       └──────────┬───────────┘
                  │
                  ▼
       ┌──────────────────────┐
       │ Model Output         │
       │ Text / ToolCall      │
       └──────────┬───────────┘
                  │
          ┌───────┴────────┐
          │                │
          ▼                ▼
       Final Text        ToolCall
          │                │
          │                ▼
          │       ┌───────────────────┐
          │       │   Tool Registry   │
          │       │                   │
          │       │ Validate          │
          │       │ Authorize         │
          │       │ Execute           │
          │       └────────┬──────────┘
          │                │
          │       ┌────────┴────────┐
          │       │                 │
          │       ▼                 ▼
          │   ToolResult        ToolError
          │       │                 │
          │       └────────┬────────┘
          │                │
          │                ▼
          │       Provider continuation
          │                │
          │                └───────┐
          │                        │
          └────────────────────────┘
                                   │
                                   ▼
                         Final Response Stream
                                   │
                    ┌──────────────┼───────────────┐
                    │              │               │
                    ▼              ▼               ▼
                Frontend       Session Store    Tracing Store
                    │              │               │
                    │              ▼               ▼
                    │         Persisted        Execution
                    │          Session           Trace
                    │
                    ▼
             User-visible UI
```

---

# 7. Runtime Data Flow in Detail

## 7.1 Phase 1 — Request Admission

```text
Client Request
    │
    ▼
Gateway
    │
    ├── Authenticate
    ├── Validate request
    └── Apply request-level policies
    │
    ▼
Agent Router
```

Output:

```text
ExecutionRequest
```

---

## 7.2 Phase 2 — Runtime Context Resolution

Agent Loop starts an execution by resolving all context required for deterministic execution.

```text
ExecutionRequest
      │
      ├──────────────▶ Agent Configuration
      │
      ├──────────────▶ SOUL / IDENTITY
      │
      ├──────────────▶ Rules
      │
      ├──────────────▶ Skills
      │
      ├──────────────▶ Session History
      │
      ├──────────────▶ Memory
      │
      └──────────────▶ Tool Permissions
                           │
                           ▼
                  Runtime Context Snapshot
```

The Runtime Context becomes immutable for the lifetime of the execution.

Changes to Skill, Rule, Tool configuration, or Agent configuration made after this point apply only to subsequent executions.

---

# 8. Context Compaction Flow

## 8.1 Automatic Compaction

Automatic compaction occurs during Runtime Context preparation.

```text
Session History
      │
      ▼
Check Context Policy
      │
      ├── Within threshold ──────────────┐
      │                                  │
      └── Exceeds threshold              │
                 │                       │
                 ▼                       │
          Memory Compaction              │
                 │                       │
                 ▼                       │
        Compacted History                │
                 │                       │
                 └──────────────┬────────┘
                                ▼
                         Memory Retrieval
                                │
                                ▼
                        Runtime Context
```

Compaction policy includes:

* threshold
* keep recent messages
* maximum token budget

Compaction is performed before model execution rather than in the middle of an already-running Agent execution.

---

## 8.2 Manual Compaction

A user can explicitly request compaction.

```text
Dashboard / CLI
      │
      ▼
Compact Session Request
      │
      ▼
Session / Memory Component
      │
      ▼
Apply Compaction Policy
      │
      ▼
Persist Compacted Session
```

Manual compaction does not require an Agent model execution.

---

# 9. Agent Model / Tool Execution Loop

The core runtime behavior is:

```text
┌────────────────────────────────────┐
│          Agent Execution           │
└─────────────────┬──────────────────┘
                  │
                  ▼
         Provider.Generate()
                  │
                  ▼
          Model Output Event
                  │
         ┌────────┴────────┐
         │                 │
         ▼                 ▼
    Final Response       ToolCall
         │                 │
         │                 ▼
         │          Tool Registry
         │                 │
         │        ┌────────┴────────┐
         │        │                 │
         │        ▼                 ▼
         │    ToolResult         ToolError
         │        │                 │
         │        └────────┬────────┘
         │                 │
         │                 ▼
         │       Add result to execution
         │       context / conversation
         │                 │
         │                 ▼
         │        Provider.Generate()
         │                 │
         │                 └─────────┐
         │                           │
         └───────────────────────────┘
                  │
                  ▼
           Terminal State
```

The loop continues until one of the following occurs:

* model produces final response
* permanent provider failure
* execution limit is reached
* unrecoverable runtime failure occurs
* request is explicitly cancelled

There is no Subagent or Agent-to-Agent delegation loop.

---

# 10. Tool Error Flow

Tool failures do not automatically terminate the Agent execution.

```text
ToolCall
   │
   ▼
Tool Registry
   │
   ├── Permission valid
   ├── Arguments valid
   └── Tool available
   │
   ├────────────── YES ──────────────▶ Execute
   │                                    │
   │                                    ▼
   │                                ToolResult
   │
   └────────────── NO / ERROR ─────────┐
                                       ▼
                                   ToolError
                                       │
                            ┌──────────┴──────────┐
                            │                     │
                            ▼                     ▼
                         Log Error          Emit Event
                            │                     │
                            └──────────┬──────────┘
                                       ▼
                              Return ToolError
                                   to model
                                       │
                                       ▼
                                Model decides
                                next strategy
```

This allows the Agent to react to failed tool execution instead of immediately terminating.

---

# 11. Provider Failure Flow

Provider failures are classified into transient and permanent failures.

## 11.1 Transient Provider Failure

```text
Provider Request
      │
      ▼
Transient Failure
      │
      ▼
Retry Policy
      │
      ├── Retry available
      │       │
      │       ▼
      │   Retry Provider
      │
      └── Retry exhausted
              │
              ▼
        Execution Failure
```

Every failure is logged and included in tracing.

---

## 11.2 Permanent Provider Failure

Examples include invalid model configuration or non-recoverable provider errors.

```text
Provider Request
      │
      ▼
Permanent Failure
      │
      ├── Log Error
      ├── Record Trace
      ├── Emit Error Event
      └── Terminate Execution
```

The user must receive enough information to understand that the execution failed.

---

# 12. Memory Failure Flow

Memory is designed to degrade gracefully when possible.

```text
Memory Retrieval
      │
      ▼
Memory Failure
      │
      ├── Log Error
      ├── Record Trace
      ├── Emit diagnostic event
      │
      ▼
Continue execution
without unavailable memory context
```

The Agent can continue when the remaining Session context is sufficient.

A Memory failure must not automatically terminate the Agent execution.

---

# 13. Persistence Failure Flow

Persistence occurs after the response has already been streamed.

```text
Final Response
      │
      ▼
Response Stream → User
      │
      ▼
Execution Finalization
      │
      ├── Session Persistence
      ├── Memory Write
      └── Tracing Persistence
```

If persistence fails:

```text
Persistence Failure
      │
      ├── Log Error
      ├── Record diagnostic information
      ├── Emit user-visible event where appropriate
      └── Do NOT rollback streamed response
```

A persistence failure therefore affects durability/observability, not the already-delivered response.

---

# 14. Realtime Streaming Flow

The system has two related but distinct realtime mechanisms.

## 14.1 Execution Events

Execution events describe what the Agent is doing.

```text
Agent Loop
   │
   ├── execution.started
   ├── skill.activated
   ├── reasoning.step
   ├── tool.requested
   ├── tool.started
   ├── tool.completed
   ├── tool.failed
   ├── execution.failed
   └── execution.completed
   │
   ▼
Realtime Event System
   │
   ▼
WebSocket / SSE
```

## 14.2 Response Streaming

Model output is streamed separately as response data:

```text
Provider
   │
   ▼
Response Delta
   │
   ▼
Agent Loop
   │
   ▼
Realtime Transport
   │
   ▼
FE / CLI
```

The final response is delivered progressively rather than waiting until the entire model execution is complete.

---

# 15. Session Lifecycle

## 15.1 Session Creation

```text
User
 │
 ▼
Agent Selection
 │
 ▼
Create Session
 │
 ▼
Session Store
 │
 ▼
Persist Session
```

---

## 15.2 Session Execution

```text
Session
 │
 ▼
Load History
 │
 ▼
Context Policy
 │
 ├── Compaction if required
 │
 ▼
Runtime Context
 │
 ▼
Agent Execution
 │
 ▼
Final Response
 │
 ▼
Persist Messages
```

---

## 15.3 Session Restoration

```text
Stored Session
      │
      ▼
Session Store
      │
      ▼
Load Session Metadata
      │
      ▼
Load Conversation History
      │
      ▼
Optional Memory Retrieval
      │
      ▼
Runtime Context
```

---

# 16. Configuration Change Semantics

Configuration changes are **execution-boundary changes**.

For example:

```text
Execution N
    │
    ├── Skill A
    ├── Rule A
    └── Tool X
```

If the user changes configuration while Execution N is running:

```text
User changes Skill / Rule / Tool configuration
                │
                ▼
       Persist new configuration
                │
                ▼
          Execution N
        remains unchanged
```

The next execution uses the new configuration:

```text
Execution N + 1
    │
    ├── New Skill configuration
    ├── New Rule configuration
    └── New Tool policy
```

This ensures execution determinism and prevents configuration changes from mutating an already-running Runtime Context.

---

# 17. Store and Persistence Architecture

```text
┌─────────────────────────────────────┐
│             Domain Layer             │
│                                     │
│ Agent / Provider / Session / Skill  │
│ Rule / Tool / Memory / API Key      │
└─────────────────────┬───────────────┘
                      │
                      ▼
┌─────────────────────────────────────┐
│             Store Layer              │
│                                     │
│ Agent Store                         │
│ Provider Store                      │
│ Session Store                       │
│ Skill Store                         │
│ Rule Store                          │
│ Tool Store                          │
│ Memory Store                        │
│ Tracing Store                       │
│ API Key Store                       │
└─────────────────────┬───────────────┘
                      │
                      ▼
                ┌───────────┐
                │ PostgreSQL│
                └───────────┘
```

Stores hide persistence implementation details from runtime components.

---

# 18. Data Ownership

| Component                       | Primary Data Responsibility                                  |
| ------------------------------- | ------------------------------------------------------------ |
| **Agent Service / Agent Store** | Agent identity and configuration                             |
| **Provider Store**              | Provider configuration, model metadata, credentials metadata |
| **Skill Store / Filesystem**    | Skill definitions and versions                               |
| **Rule Store**                  | Rules, versions, assignments                                 |
| **Tool Store**                  | Tool registration and metadata                               |
| **Session Store**               | Sessions and conversation messages                           |
| **Memory Store**                | Memory records, embeddings, memory metadata                  |
| **Tracing Store**               | Execution traces, spans, telemetry                           |
| **API Key Store**               | External API key metadata and usage                          |
| **Realtime Event System**       | Runtime event delivery, not durable system-of-record data    |

---

# 19. API Key Flow

API keys are intended for clients or agents that need to call the local gateway.

```text
User
 │
 ▼
API Key Management
 │
 ├── Create
 ├── Scope
 ├── Expiration
 └── Revoke
 │
 ▼
API Key Store
```

At request time:

```text
External Client
      │
      ▼
Gateway
      │
      ▼
Validate API Key
      │
      ├── Scope
      ├── Expiration
      └── Status
      │
      ▼
Authorized Request
```

Usage is recorded for monitoring.

---

# 20. Global Configuration

Global configuration consists of:

```text
Config
│
├── Server
│   ├── Host
│   └── Port
│
├── AI Defaults
│   ├── Default Provider
│   ├── Default Model
│   └── Model parameters
│
├── Quota
│   └── Per-Agent / Per-User limits
│
└── Tools
    └── Global tool enable/disable configuration
```

Global configuration is separate from Agent-specific runtime configuration.

Agent configuration overrides or specializes global defaults where applicable.

---

# 21. Observability Model

Observability is built around three levels.

## 21.1 Logs

Logs capture failures and diagnostic information.

Every execution failure should include sufficient information to identify:

* operation
* component
* execution/session identifier
* failure category
* cause
* recovery behavior

---

## 21.2 Traces

Tracing records the lifecycle of the execution:

```text
Execution
│
├── Gateway
├── Router
├── Runtime Context
├── Memory Retrieval
├── Provider Call
├── Tool Call
├── Provider Continuation
├── Final Response
├── Session Persistence
└── Memory Processing
```

---

## 21.3 Realtime Events

Realtime events are intended for user-facing execution visibility:

```text
Chat
Reasoning Step
Skill Activation
Tool Call
Tool Result
Error
Completion
```

These are transient UI events and are not the primary durable source of execution history.

---

# 22. Execution State

The Agent Loop maintains explicit execution state.

Conceptually:

```text
ExecutionState
│
├── Execution ID
├── Agent ID
├── Session ID
├── Runtime Context Snapshot
├── Current Model Turn
├── Tool Call Count
├── Retry Count
├── Execution Status
└── Error State
```

Possible terminal states:

```text
COMPLETED
FAILED
CANCELLED
LIMIT_REACHED
```

Possible intermediate states:

```text
INITIALIZING
BUILDING_CONTEXT
CALLING_PROVIDER
EXECUTING_TOOL
STREAMING_RESPONSE
FINALIZING
```

---

# 23. Execution Limits

The Agent Loop must enforce bounded execution.

Limits can apply to:

* maximum tool calls
* maximum provider retries
* maximum execution duration
* maximum context size
* maximum model/tool continuation depth

The objective is to prevent unbounded execution and accidental infinite loops.

There is no recursive Agent-to-Agent execution because Subagent and Delegation are out of scope.

---

# 24. End-to-End Agent Execution

The complete execution can be summarized as:

```text
┌────────────────────────────────────────────────────────────┐
│                     1. REQUEST                            │
│                                                            │
│ User → Gateway → Agent Router → ExecutionRequest           │
└───────────────────────────┬────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────────┐
│                  2. CONTEXT RESOLUTION                     │
│                                                            │
│ Agent Config                                               │
│ + SOUL / IDENTITY                                          │
│ + Rules                                                    │
│ + Skills                                                   │
│ + Session History                                          │
│ + Memory Retrieval / Compaction                            │
│ + Tool Permissions                                         │
│                                                            │
│               ↓ Runtime Context Snapshot                   │
└───────────────────────────┬────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────────┐
│                    3. AGENT LOOP                           │
│                                                            │
│            Agent Loop → Provider                           │
│                         │                                  │
│                    Model Output                            │
│                         │                                  │
│                  ┌──────┴──────┐                           │
│                  │             │                            │
│               Response      ToolCall                       │
│                  │             │                            │
│                  │        Tool Registry                    │
│                  │             │                            │
│                  │       ToolResult/Error                  │
│                  │             │                            │
│                  │             └──────▶ Provider            │
│                  │                       │                  │
│                  │                       └──── repeat ────┐ │
│                  │                                        │ │
│                  └────────────────────────────────────────┘ │
└───────────────────────────┬────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────────┐
│                     4. RESPONSE                            │
│                                                            │
│ Provider → Agent Loop → WebSocket / SSE → User             │
└───────────────────────────┬────────────────────────────────┘
                            ▼
┌────────────────────────────────────────────────────────────┐
│                   5. FINALIZATION                           │
│                                                            │
│ Session Persistence                                        │
│ + Memory Write / Embedding                                 │
│ + Compaction if required                                   │
│ + Trace Persistence                                        │
│ + Final Execution Event                                    │
└────────────────────────────────────────────────────────────┘
```

---

# 25. Architecture Validation

The architecture has been validated against the functional scope defined for the core Agent Gateway.

## 25.1 Agent Execution

**Validated**

The Agent Loop is the single runtime orchestrator.

```text
Gateway
  ↓
Agent Router
  ↓
Agent Loop
  ↓
Provider ↔ Tool Registry loop
```

There is no second orchestration engine competing with Agent Loop.

---

## 25.2 Provider Boundary

**Validated**

Provider-specific behavior stays behind Provider Abstraction.

```text
Agent Loop
   ↓
Provider Abstraction
   ↓
OpenAI-Compatible / ACP
```

The Agent Loop does not contain provider-specific HTTP or protocol logic.

---

## 25.3 Tool Boundary

**Validated**

Tool access is centralized in Tool Registry.

```text
Agent Loop
   ↓
Tool Registry
   ↓
Permission Validation
   ↓
Tool Execution
```

Agent Loop does not independently bypass tool policy.

---

## 25.4 Skill / Rule Boundary

**Validated**

Skill and Rule are resolved before execution and become part of Runtime Context.

Configuration changes apply to the next execution.

```text
Skill / Rule update
        ↓
Persist configuration
        ↓
Next execution
        ↓
New Runtime Context
```

---

## 25.5 Memory Boundary

**Validated**

Memory has two distinct paths:

```text
READ:
Session → Memory Retrieval → Runtime Context

WRITE:
Execution Complete → Memory Write / Embedding / Compaction
```

Automatic compaction happens before Runtime Context creation when context limits require it.

Manual compaction is also supported.

---

## 25.6 Session Boundary

**Validated**

Session is durable state, not execution control.

```text
Session Store
     ↓
Conversation History
     ↓
Runtime Context

Execution
     ↓
Session Persistence
```

Session persistence failure does not rollback a response that has already been streamed.

---

## 25.7 Realtime Boundary

**Validated**

Realtime Events are side effects of Agent execution.

```text
Agent Loop
    ↓
Realtime Event System
    ↓
WebSocket / SSE
    ↓
Frontend
```

Realtime delivery failure does not terminate the Agent execution.

---

## 25.8 Observability Boundary

**Validated**

Failures are observable through:

```text
Log
 +
Trace
 +
Realtime Error Event
```

This allows the user and system operator to understand execution failures.

---

## 25.9 Failure Semantics

**Validated**

| Failure                     | Behavior                                                        |
| --------------------------- | --------------------------------------------------------------- |
| Tool execution error        | Return ToolError to model; model may recover or change strategy |
| Tool permission denied      | Return ToolError; log and emit event                            |
| Provider transient error    | Retry with bounded retry policy                                 |
| Provider permanent error    | Log + trace + error event + terminate                           |
| Memory retrieval failure    | Log + trace; continue without unavailable memory when possible  |
| Session persistence failure | Log + trace; do not rollback streamed response                  |
| Memory persistence failure  | Log + trace; do not rollback streamed response                  |
| Tracing failure             | Log locally; execution continues                                |
| Realtime delivery failure   | Log; execution continues                                        |

---

# 26. Design Principles Applied

These principles guide all architectural decisions.

| Principle                            | Application                                                                                        |
| ------------------------------------ | -------------------------------------------------------------------------------------------------- |
| **Write for the next reader**        | Explicit component boundaries, contracts, failure semantics, and execution stages                  |
| **One responsibility per component** | Provider, Tool Registry, Memory, Session, Realtime, and Tracing each own a distinct responsibility |
| **Explicit over implicit**           | Runtime Context and ExecutionRequest explicitly define the state entering execution                |
| **Fail visibly**                     | Every failure is logged; user-relevant failures generate realtime diagnostics                      |
| **Local reasoning**                  | Runtime components depend on abstractions rather than implementation-specific stores/providers     |
| **Fast by design**                   | Streaming response path avoids waiting for post-execution persistence                              |
| **Deterministic execution**          | Runtime Context is snapshotted at execution start                                                  |
| **Graceful degradation**             | Memory, persistence, tracing, and realtime failures do not unnecessarily destroy a valid execution |
| **Bounded execution**                | Tool calls, retries, duration, context size, and continuation depth are constrained                |
| **No agent loops**                   | No Subagent/Delegation; Agent execution is bounded to one Agent Runtime                            |

---

# 27. Final Architecture Summary

The system can be understood as five primary responsibilities:

```text
                    ┌─────────────────────┐
                    │      Gateway        │
                    │ Authenticate/Route  │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │    Agent Router     │
                    │ Resolve Execution   │
                    └──────────┬──────────┘
                               │
                               ▼
             ┌────────────────────────────────┐
             │          AGENT LOOP             │
             │                                │
             │  Runtime Context               │
             │       ↓                        │
             │  Provider                      │
             │       ↕                        │
             │  Tool Registry                 │
             │       ↕                        │
             │  Bounded Tool/Model Loop       │
             └───────┬───────────────┬────────┘
                     │               │
                     │               │
                     ▼               ▼
             ┌─────────────┐   ┌──────────────┐
             │ Realtime    │   │ Persistence  │
             │ Events      │   │ + Tracing    │
             └──────┬──────┘   └──────┬───────┘
                    │                 │
                    ▼                 ▼
                 Frontend          PostgreSQL
```

The central architectural rule is:

> **Agent Loop owns execution orchestration. Specialized components own capabilities and state.**

The execution lifecycle is therefore:

```text
Request
  ↓
Authenticate
  ↓
Route
  ↓
Resolve Runtime Context
  ↓
Compact if necessary
  ↓
Retrieve Memory
  ↓
Call Provider
  ↓
Process Tool Calls
  ↓
Repeat bounded model/tool loop
  ↓
Stream Final Response
  ↓
Persist Session
  ↓
Write Memory / Embedding
  ↓
Trace Execution
  ↓
Complete
```

This architecture remains aligned with the core product scope: a local, single-owner AI Agent Gateway centered around **Agent, Skill, Rules, Tool, Provider, Memory, Session, Monitoring, Realtime Events, API Key, and Configuration**, without introducing Channel, MCP, Subagent, Delegation, Teams, TTS, Webhook, or Cron capabilities.
