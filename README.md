# CloudianClaw

[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15%2B-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![React](https://img.shields.io/badge/React-18-%2361DAFB?logo=react&logoColor=white)](https://react.dev/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5-%233178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Realtime](https://img.shields.io/badge/Realtime-WebSocket%20%2F%20SSE-9cf)](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API)
[![CLI](https://img.shields.io/badge/CLI-Cobra%20%2F%20BubbleTea-cc0000)](#)
[![LLM](https://img.shields.io/badge/LLM-OpenAI--compatible-412991)](https://platform.openai.com/docs/guides/text-generation)
[![Local-first](https://img.shields.io/badge/Deployment-Local--first-brightgreen)](#)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](#)

**CloudianClaw** is a local-first, single-owner **AI Agent Gateway**. It acts as an
*Agent Harness* that turns Large Language Models into independent AI Agents for
coding assistance and troubleshooting. The system is built around a single
**Agent Loop** orchestrator, surrounded by specialized capabilities (Provider
abstraction, Tool Registry, Skill/Rule loading, Memory, Session store, Realtime
events, and Tracing).

> **Design principle:** *The Agent Loop owns execution orchestration. Specialized
> components own capabilities and state.* No Subagent, Delegation, or
> multi-agent orchestration is included.

---

## Overview

CloudianClaw lets you run, configure, and observe AI Agents from a **React
Dashboard**, a **Terminal CLI**, or external clients via API keys. It is local,
offline-friendly, and portable:

- **Local-first launch** — no complex environment setup; migrations run
  automatically at startup.
- **Flexible agent management** — configure agents through `SOUL.md` /
  `IDENTITY.md`, plus behavior Rules and capability Skills.
- **Bounded Agent Loop** — LLMs may call system Tools securely, with hard limits
  that prevent infinite loops.
- **Real-time interaction** — stream responses and watch reasoning, skill
  activation, and tool calls live through the Web Dashboard or Terminal CLI.
- **Portable vector memory** — embeddings are stored as `double precision[]`
  (or `jsonb`) arrays in PostgreSQL; **Cosine Similarity is computed in Go**
  (no `pgvector` extension required).
- **Graceful degradation** — failures in memory, tracing, persistence, or
  realtime delivery never crash an in-flight Agent execution.

### Scope

| In scope | Out of scope |
| --- | --- |
| Gateway auth (Gateway API Token), Provider discovery/validation | Channels (Slack/Discord/Telegram), MCP |
| Agent lifecycle (CRUD, clone), SOUL/IDENTITY config | Subagent orchestration, Delegation |
| Skill/Rules loading, sandbox testing | Teams / multi-tenancy |
| Tool Registry + permission policy, file/CLI tools | TTS, Webhook dispatcher, Cron scheduler |
| Session persistence & restoration | Skill marketplace |
| Memory + Go Cosine vector search, compaction | Multi-provider routing/fallback |
| Realtime events (WS/SSE), streaming | |
| Tracing/observability, API keys + global config | |
| Terminal CLI + React Dashboard | |

See [`docs/documents/prd_product_overview.md`](docs/documents/prd_product_overview.md)
and [`docs/documents/architecture.md`](docs/documents/architecture.md) for the
full requirements and architecture.

---

## Folder structure

The codebase follows a **Domain-Driven Design (DDD)** layout: domain types and
contracts live under `internal/domain`, concrete implementations under
`internal/impl`, transport/gateway under `internal/transport`, and the TUI under
`internal/tui`.

```text
CloudianClaw/
├── cmd/
│   └── agentkit/            # CLI binary entrypoint (Cobra commands: serve, chat, ...)
│
├── internal/
│   ├── model/               # Core domain entities & runtime contracts (no SQL/HTTP deps)
│   │                         #   Agent, Provider, Model, Skill, Rule, Tool, Session,
│   │                         #   Message, Memory, Usage, Execution, ExecutionRequest,
│   │                         #   RuntimeContext, ExecutionState, ToolCall/Result/Error
│   │
│   ├── domain/              # Pure domain layer (no framework/infra imports)
│   │   ├── entity/          # Entity definitions & domain validation
│   │   ├── event/           # Domain events (execution.*, tool.*, skill.*, ...)
│   │   ├── interface/       # Port interfaces (Provider, Tool, Store, Registry, ...)
│   │   └── service/         # Agent Loop orchestrator & domain services
│   │                         #   (runtime context builder, bounded execution loop)
│   │
│   ├── impl/                # Concrete implementations (adapters) of the ports
│   │   ├── config/          # Config loading (.env) + application container wiring
│   │   ├── database/        # DB connection & migration runner
│   │   ├── memory/          # Memory write/embed + Go Cosine similarity search
│   │   ├── provider/        # Provider abstraction + OpenAI-compatible adapter
│   │   ├── session/         # Session lifecycle & history persistence
│   │   ├── store/           # Store interfaces
│   │   │   └── postgres/    # PostgreSQL implementations (database/sql)
│   │   └── tool/            # Tool Registry + built-in file/CLI tools + policy
│   │
│   ├── transport/           # Gateway / routing layer (request admission)
│   │   ├── http/            # HTTP server, auth middleware, REST + SSE, Agent Router
│   │   └── socket/          # WebSocket realtime event transport
│   │
│   ├── tui/                 # Terminal UI (BubbleTea)
│   │   ├── command/         # CLI command tree (Cobra)
│   │   └── formatter/       # Terminal output formatting
│   │
│   └── exception/           # Centralized error types & wrapping
│
├── agent/                   # Local agent assets (filesystem, not code execution)
│   ├── skills/              # SKILL.md packages (Markdown + YAML frontmatter)
│   └── rules/               # Rule definitions (Markdown)
│
├── web/                     # React + TypeScript + TailwindCSS Dashboard
│
├── docs/                    # Project documentation
│   ├── documents/           # PRD, architecture, idea (CloudianClaw-specific specs)
│   ├── references/          # Upstream GoClaw study material (reference only)
│   ├── guides/              # Setup/usage guides
│   ├── plans/               # Implementation plan (plan.md)
│   └── ui/                  # UI wireframes & images
│
├── migrations/              # SQL migrations (created at build time)
├── main.go                  # Temporary placeholder entry (CLI wired via cmd/agentkit)
├── go.mod                   # Go module: cloudian/cloudian-claw
└── .env.example             # Configuration template
```

> **Note on references:** Files under `docs/references/` (incl. `STEP_INDEX.md`,
> `README.md`) describe the upstream **GoClaw** project and are kept as study
> material. They are *reference knowledge*, not the directory structure to copy.
> The authoritative CloudianClaw layout is this scaffold and `docs/plans/plan.md`.

---

## How to setup

### 1. Prerequisites

| Tool | Version | Purpose |
| --- | --- | --- |
| Go | 1.22+ | Build & run the backend / CLI |
| PostgreSQL | 15+ | Durable storage (no `pgvector` needed) |
| Node.js + npm | 18+ | Build & run the React Dashboard (`web/`) |
| Git | any | Clone the repository |

### 2. Clone & install Go dependencies

```bash
git clone <your-repo-url> CloudianClaw
cd CloudianClaw

# Download Go module dependencies
go mod download

# (Optional) Tidy modules after adding/removing deps
go mod tidy
```

### 3. Configure environment

Copy the example env file and fill in your values:

```bash
cp .env.example .env
```

Edit `.env` with at least the following (all values are local / static):

```dotenv
# --- Gateway ---
# Static token checked against the Authorization: Bearer <token> header.
# Clients also persist this in LocalStorage after onboarding.
GATEWAY_API_TOKEN=change-me-to-a-long-random-string

# --- Server ---
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# --- PostgreSQL ---
# The migration runner creates all tables automatically at startup.
POSTGRES_URL=postgres://CloudianClaw:CloudianClaw@localhost:5432/CloudianClaw?sslmode=disable

# --- AI defaults (overridable per agent) ---
DEFAULT_PROVIDER=openrouter
DEFAULT_MODEL=openai/gpt-4o-mini
AI_TEMPERATURE=0.7
AI_MAX_TOKENS=4096

# --- Quota (per agent / per user, single-owner) ---
QUOTA_MAX_REQUESTS_PER_DAY=0        # 0 = unlimited

# --- Global tool toggle ---
GLOBAL_TOOLS_ENABLED=true

# --- Bounded execution limits ---
MAX_TOOL_CALLS=20
MAX_PROVIDER_RETRIES=3
MAX_EXECUTION_DURATION=300s
MAX_CONTEXT_TOKENS=200000
MAX_CONTINUATION_DEPTH=10

# --- Provider credential (example: OpenRouter) ---
OPENROUTER_API_KEY=sk-or-...
```

> The schema uses `double precision[]` (or `jsonb`) for embeddings so the system
> runs on a stock PostgreSQL without the `pgvector` extension. See
> [`docs/DATABASE.txt`](docs/DATABASE.txt) for the data model.

### 4. Prepare the database

Create the database (if it does not exist) and let the server apply migrations:

```bash
# Create the database role + database (one-time, using psql)
createdb CloudianClaw
# or: psql -U postgres -c "CREATE DATABASE CloudianClaw;"
```

Migrations run automatically when the server boots (see `internal/impl/database`).
To run them explicitly (once the migration CLI exists):

```bash
go run ./cmd/agentkit migrate up
```

### 5. Run the server (Gateway)

```bash
# Build
go build ./...

# Start the Gateway (HTTP + WebSocket) on SERVER_PORT
go run ./cmd/agentkit serve
```

On startup the server will:
1. Load `.env` configuration (fails fast on missing required vars).
2. Run pending SQL migrations against PostgreSQL.
3. Listen for requests, enforcing `GATEWAY_API_TOKEN` auth.

Verify health:

```bash
curl -i http://localhost:8080/health
```

### 6. Run the Terminal CLI

```bash
# Interactive chat with the default agent
go run ./cmd/agentkit chat --agent default

# Other commands (structure defined in internal/tui/command)
go run ./cmd/agentkit --help
go run ./cmd/agentkit agent --help
go run ./cmd/agentkit skill --help
go run ./cmd/agentkit session --help
go run ./cmd/agentkit config --help
go run ./cmd/agentkit api-key --help
```

The CLI streams tokens to the terminal and prints tool-call sequences clearly.

### 7. Run the React Dashboard

```bash
cd web
npm install
npm run dev          # Vite dev server (default http://localhost:5173)
# Production build:
npm run build && npm run preview
```

Open the Dashboard, complete the **Onboarding** wizard
(enter Gateway Token → configure Provider/Model → create first Agent via
`SOUL.md` / `IDENTITY.md`), then chat with your agent in real time.

### 8. Common commands

```bash
# Format & lint (per AGENTS.md / docs)
gofmt -l -w .
goimports -w .
golangci-lint run

# Test
go test ./...
go test -race -cover ./...
```

---

## Documentation map

- `docs/documents/prd_product_overview.md` — product requirements (US / FR).
- `docs/documents/architecture.md` — full system & component architecture.
- `docs/documents/idea.md` — scope, use cases, system use cases.
- `docs/plans/plan.md` — 13-phase implementation plan mapped to this scaffold.
- `docs/DATABASE.txt` — data model (single-owner, Go-Cosine embeddings).
- `docs/references/` — upstream GoClaw study material (reference only).


Build with Cloudian Love CLoud 
