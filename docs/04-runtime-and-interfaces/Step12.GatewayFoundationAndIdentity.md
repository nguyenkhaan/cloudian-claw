# Step 12 — Gateway Foundation and Identity

**Knowledge depth: 7/10**

Read [04 — Gateway and Protocol](../04-gateway-protocol.md) before creating a server. It explains how GoClaw joins HTTP, WebSocket, health checks, routing, and identity propagation at one boundary. Read [20 — API Keys & Authentication](../20-api-keys-auth.md) alongside it so the gateway does not treat raw request parameters as authority.

## A gateway is the edge of the agent system

The gateway has a narrower purpose than the agent loop:

```text
accept connection → authenticate → resolve scope → route request → schedule run → deliver response
```

It does not decide how to answer the user. That remains inside the shared application service and agent runtime built in earlier steps.

## Identity travels with the request

At the gateway edge, resolve the values that every downstream operation needs:

```text
request ID · tenant ID · user ID · agent ID · role · locale · session key
```

The system passes this context to the scheduler, stores, tools, and tracing. GoClaw's gateway/router code demonstrates this progression. Understanding it now makes the tenant and permission model in Step 16 much easier.

## Keep the server composition readable

```go
type Server struct {
	chat      ChatService
	auth      Authenticator
	scheduler Scheduler
	router    http.Handler
}
```

Compose concrete stores, providers, registries, and channels in one application entry point. The handlers should only depend on the interfaces they call. This prevents HTTP behavior from becoming a second agent implementation.

## Start with a small route surface

`/health` says whether the process can serve; `/v1/chat/completions` exposes a normal request/response interaction; `/ws` is added in Step 14 for persistent event-driven interaction. Other management routes can grow with the platform.

