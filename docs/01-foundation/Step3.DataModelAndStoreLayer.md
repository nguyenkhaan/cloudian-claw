# Step 3 — Data Model and Store Layer

**Knowledge depth: 7/10**

Read [06 — Store Layer and Data Model](../06-store-data-model.md) before designing tables. It is the GoClaw source for interface-first stores, PostgreSQL/SQLite variants, migrations, and the data that must survive process restarts. For tenant-scoped data, keep [23 — Multi-Tenant Architecture](../23-multi-tenant-architecture.md) open while you read.

## Why persistence comes early

An agent cannot become reliable if its conversation, identity, configuration, and audit data live only in process memory. The store layer turns runtime concepts into durable records while keeping SQL out of the agent loop.

```text
agent loop → SessionStore / AgentStore / ProviderStore interfaces → database implementation
```

The interface belongs close to the package that consumes it. The database implementation belongs behind that interface. This is the same separation that lets GoClaw support PostgreSQL for the gateway and SQLite for its desktop edition.

## Start with four aggregates

```text
tenants       who owns data and policy
agents        model choice, prompt identity, tool policy
sessions      one ongoing conversation and its summary
messages      ordered canonical conversation records
```

Later steps add memories, traces, API keys, tasks, skills, and channels. Do not flatten all of them into one JSON blob: a session needs ordered messages, an agent needs configuration, and a tenant defines the scope for both.

## Scope before schema details

Every durable record needs a clear owner. For an initial system, use this mental key:

```text
tenant → agent → user → session → message
```

The order matters. A session key is a convenient identifier, not proof that the caller is allowed to read it. This principle becomes the foundation of Step 16.

## A small store boundary

```go
type SessionStore interface {
	Load(ctx context.Context, key string) ([]model.Message, string, error)
	Append(ctx context.Context, key string, messages ...model.Message) error
	SetSummary(ctx context.Context, key, summary string) error
}
```

The agent loop sees this contract, never table names or SQL dialect details. That makes it possible to reason about conversation behavior separately from migration work.

## How to study the GoClaw source

In `06-store-data-model.md`, follow one entity from its store interface to its PostgreSQL implementation, then to the matching migration. Notice how GoClaw keeps relationships, indices, and query scope explicit. Repeat that exercise later for memories and traces.

