# PROJECT ARCHITECTURE DESIGN 
## 1. System Design

**Objective**: Define the architectural blueprint, component responsibilities, data flows, and implementation approach.

### 3.1 Technology Stack

| Layer | Technology | Rationale |
|-------|------------|-----------|
| Backend | Go (Golang) | High performance, fast startup, concurrent request handling |
| Frontend | React + TypeScript + TailwindCSS | Component reusability, type safety, rapid UI development |
| Database | PostgreSQL | Relational data integrity, JSONB support, extensibility |
| Cache | In-Memory (RAM) + Redis | Low-latency session state, hot data caching |
| Real-time | WebSocket / SSE | Streaming responses, live event updates |
| CLI | Cobra / BubbleTea | Interactive terminal experience, cross-platform |

### 3.2 Architectural Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT LAYER                              │
├──────────────────────┬──────────────────────────────────────────┤
│   React Dashboard    │     Terminal CLI Interface                │
│   (localhost:3000)   │     (go run ./cmd/agent-cli)             │
└──────────┬───────────┴──────────────────────────────────────────┘
           │                    │
           │   HTTP / WebSocket │
           └────────┬───────────┘
                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                     GATEWAY LAYER                                │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐   │
│  │ Auth     │ │ Rate     │ │ Request  │ │ Router / Load    │   │
│  │ Middleware│ │ Limiter  │ │ Validator│ │ Balancer         │   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘   │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                      SERVICE LAYER                               │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────────┐   │
│  │Agent     │ │Session   │ │Skill     │ │Tool               │   │
│  │Service   │ │Service   │ │Service   │ │Registry Service   │   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘   │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐                        │
│  │Memory    │ │Provider  │ │Monitoring│ │Config               │   │
│  │Service   │ │Service   │ │Service   │ │Service              │   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────────┘   │
└────────────────────────┬────────────────────────────────────────┘
                         │
          ┌──────────────┼──────────────┐
          ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│  PostgreSQL  │ │   RAM/Cache   │ │  External    │
│  (Sessions,  │ │  (Hot Data,   │ │  Providers   │
│   Agents,    │ │   Sessions)   │ │  (OpenRouter,│
│   Skills,    │ │               │ │   Vercel...) │
│   Rules)     │ │               │ │              │
└──────────────┘ └──────────────┘ └──────────────┘
```

### 3.3 Component Specifications

Each component follows the specification format: **Purpose**, **Communicates With**, **Input/Output**, **Failure Behavior**.

#### Component: Agent Service
| Attribute | Description |
|-----------|-------------|
| **Purpose** | Orchestrates agent creation, configuration, and execution lifecycle |
| **Communicates With** | Provider Service, Skill Service, Tool Registry, Memory Service |
| **Input** | Agent definition (SOUL.md, model selection, skill/tool assignments), execution requests |
| **Output** | Agent instances, execution results, streaming responses |
| **Failure Behavior** | Returns validation errors for invalid configs; falls back to default agent if primary fails |

#### Component: Provider Service
| Attribute | Description |
|-----------|-------------|
| **Purpose** | Manages connections to external AI model providers and their available models |
| **Communicates With** | Agent Service, Gateway Layer |
| **Input** | Provider API keys, model selection criteria |
| **Output** | Validated provider connections, model catalogs, generated completions |
| **Failure Behavior** | Marks provider as unhealthy; routes requests to alternate provider if available |

#### Component: Skill Service
| Attribute | Description |
|-----------|-------------|
| **Purpose** | Loads, validates, manages lifecycle and assignment of agent skills |
| **Communicates With** | Agent Service, File System, Tool Registry |
| **Input** | Skill definitions (files), enable/disable commands, version info |
| **Output** | Validated skill objects, skill metadata, dependency graphs |
| **Failure Behavior** | Skips invalid skills with warnings; maintains previous working version |

#### Component: Tool Registry
| Attribute | Description |
|-----------|-------------|
| **Purpose** | Catalogs all available tools and controls agent tool access permissions |
| **Communicates With** | Agent Service, External Systems, Skill Service |
| **Input** | Tool definitions, permission policies, execution requests |
| **Output** | Approved tool invocations, execution results, audit logs |
| **Failure Behavior** | Rejects unauthorized tool calls; returns structured error with allowed tools |

#### Component: Memory Service
| Attribute | Description |
|-----------|-------------|
| **Purpose** | Manages short-term session context and long-term agent memory storage |
| **Communicates With** | Agent Service, PostgreSQL, Embedding Provider |
| **Input** | Conversation messages, embedding configuration, compaction thresholds |
| **Output** | Relevant context snippets, memory summaries, retrieval rankings |
| **Failure Behavior** | Falls back to raw message history if vector search fails; respects compaction limits |

#### Component: Session Service
| Attribute | Description |
|-----------|-------------|
| **Purpose** | Tracks conversation sessions, persists state, manages session lifecycle |
| **Communicates With** | Agent Service, PostgreSQL, Memory Service |
| **Input** | Session creation requests, message payloads, session metadata |
| **Output** | Active sessions, message history, session status |
| **Failure Behavior** | Persists partial messages; recovers session state from last known checkpoint |

#### Component: Monitoring Service
| Attribute | Description |
|-----------|-------------|
| **Purpose** | Collects telemetry, traces agent operations, exposes system health metrics |
| **Communicates With** | All Services, Realtime Event System |
| **Input** | Metrics, trace spans, error events, log entries |
| **Output** | Dashboards, alert triggers, trace visualizations, event cards |
| **Failure Behavior** | Continues operation with degraded observability; buffers events for later processing |

### 3.4 Request Flow Diagram

```
┌─────────┐     ┌──────────┐     ┌───────────┐     ┌──────────┐
│  User   │────▶│ Gateway  │────▶│  Router   │────▶│  Agent   │
│ (FE/CLI)│     │ Auth/Val │     │ Selection │     │ Service  │
└─────────┘     └──────────┘     └───────────┘     └──────────┘
                                           │              │
                                           │              ▼
                                           │     ┌──────────────────┐
                                           │     │ LLM Provider     │
                                           │     │ (OpenRouter/etc) │
                                           │     └──────────────────┘
                                           │              │
                                           │              ▼
                                           │     ┌──────────────────┐
                                           │     │ Decision Engine  │
                                           │     │ + Tool Calling   │
                                           │     └──────────────────┘
                                           │              │
                                           │        ┌─────┴─────┐
                                           │        ▼           ▼
                                           │  ┌──────────┐ ┌──────────┐
                                           │  │ Tools    │ │ Memory   │
                                           │  │ Execute  │ │ Retrieve │
                                           │  └──────────┘ └──────────┘
                                           │              │
                                           └──────────────┼──────────────┘
                                                          ▼
                                                   ┌──────────┐
                                                   │Response  │
                                                   │Stream    │
                                                   └──────────┘
```

### 3.5 Data Flow: Agent Execution

```
1. User sends message via Dashboard or CLI
        │
        ▼
2. Gateway validates request and authenticates user
        │
        ▼
3. Router selects target agent based on session context
        │
        ▼
4. Agent Service loads agent configuration (SOUL.md, model, skills, tools)
        │
        ▼
5. Memory Service retrieves relevant context for this session
        │
        ▼
6. Combined prompt sent to selected LLM Provider
        │
        ▼
7. LLM generates response (may include tool calls)
        │
        ▼
8. If tool calls present → Tool Registry validates and executes
        │
        ▼
9. Results fed back to LLM for final response
        │
        ▼
10. Streaming response sent back to user interface
        │
        ▼
11. Session Service persists message pair to PostgreSQL
        │
        ▼
12. Monitoring Service records metrics and emits realtime event
```

### 3.6 Design Principles Applied

These principles guide all architectural decisions:

| Principle | Application |
|-----------|-------------|
| **Write for the next reader** | Clear component boundaries, explicit contracts, documented failure modes |
| **One responsibility per component** | Each service owns one domain concept (agents, sessions, skills, etc.) |
| **Explicit over implicit** | All configuration comes from user input; no magic defaults |
| **Fail visibly** | Errors include context, cause, and suggested resolution |
| **Local reasoning** | Components communicate through well-defined interfaces; internal state is encapsulated |
| **Fast by design** | Go backend ensures low latency; in-memory caching for hot paths; PostgreSQL for durable state |
| **No agent loops** | Guard clauses prevent recursive tool calling; bounded recursion depth enforced |
