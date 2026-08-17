# Step 24 — Web UI Foundation and Agent Settings

**Knowledge depth: 7/10**

Build the React application shell, typed clients, authentication state, and agent settings workflow.

## Step outcome

The responsive UI connects to the gateway, survives reload, and treats server-backed settings as authoritative.

## Task 1 — Define the UI boundary

### Theory

- **Architectural role:** Route → feature hook → typed client forms the inbound presentation adapter. The gateway is the source of truth for agents, sessions, runs, tools, and memory.
- **Why:** Copying agent logic or policy into TypeScript creates two implementations, and the browser is not a trust boundary. A typed client centralizes contract drift in one place.
- **GoClaw reference:** Inspect root routes in [`goclaw/ui/web/src/App.tsx`](../../goclaw/ui/web/src/App.tsx), protocol types in [`goclaw/ui/web/src/api/protocol.ts`](../../goclaw/ui/web/src/api/protocol.ts), the HTTP client in [`goclaw/ui/web/src/api/http-client.ts`](../../goclaw/ui/web/src/api/http-client.ts), and the WebSocket client in [`goclaw/ui/web/src/api/ws-client.ts`](../../goclaw/ui/web/src/api/ws-client.ts).

Use [04 — Gateway and Protocol](../04-gateway-protocol.md), [18 — HTTP REST API](../18-http-api.md), and [19 — WebSocket RPC Methods](../19-websocket-rpc.md) as the UI contracts.

The browser owns presentation, routing, temporary interaction state, and connection recovery. The gateway owns sessions, agent settings, tools, memory, and runs.

```text
React route → query/mutation hook → HTTP/WS client → gateway
```

### Goal

Decide which presentation and interaction concerns belong to the browser and which data and behavior always belong to the gateway.

### Guide to implement

Do not copy agent logic into TypeScript. Define typed API and protocol models that match Steps 20–21 and keep them in one client layer.

## Task 2 — Create the React project

### Theory

- **Architectural role:** `app` composes cross-cutting providers and routes, `features` owns UI use cases, `components` contains reusable presentation, `lib/api` is the gateway adapter, and `stores` contains appropriate client state.
- **Why:** Feature-based organization keeps chat or settings changes local while clients and protocols remain shared. Provider keys must not enter the browser because build-time environment variables are still shipped to users.
- **GoClaw reference:** Compare [`goclaw/ui/web/src/App.tsx`](../../goclaw/ui/web/src/App.tsx), pages under [`goclaw/ui/web/src/pages`](../../goclaw/ui/web/src/pages), shared components in [`goclaw/ui/web/src/components`](../../goclaw/ui/web/src/components), and app providers in [`goclaw/ui/web/src/components/providers/app-providers.tsx`](../../goclaw/ui/web/src/components/providers/app-providers.tsx).

### Goal

Create a frontend skeleton organized by features and boundaries instead of placing components, APIs, and state in one flat directory.

### Guide to implement

Create a React + TypeScript application under `web/`. A suitable first structure is:

```text
web/src/
├── app/          routes and providers
├── features/
│   ├── chat/
│   ├── sessions/
│   ├── agent-settings/
│   └── runs/
├── components/   reusable UI pieces
├── lib/          HTTP/WS clients and formatters
├── stores/       connection and local UI state
└── styles/
```

Set development configuration for gateway HTTP and WebSocket URLs. Do not place provider keys in browser environment variables.

## Task 3 — Build the application shell

### Theory

- **Architectural role:** The shell owns the route outlet, responsive navigation, and global feedback. An error boundary prevents a render failure from spreading across the app but does not replace API error handling.
- **Why:** Designing accessibility and responsiveness in the shell avoids reworking every page. Empty, loading, and fatal states are product states, not decoration added at the end.
- **GoClaw reference:** Find layout components with `rg --files goclaw/ui/web/src/components/layout`, inspect root composition in [`goclaw/ui/web/src/App.tsx`](../../goclaw/ui/web/src/App.tsx), and review shared feedback components under [`goclaw/ui/web/src/components/ui`](../../goclaw/ui/web/src/components/ui).

### Goal

Build a stable navigation, layout, error, and loading shell so feature pages can focus on their content.

### Guide to implement

Add:

- Root providers and error boundary.
- Routes for `/chat/:sessionKey?` and `/settings/agent`.
- Desktop sidebar and mobile drawer navigation.
- Main scrollable content region.
- Notification/toast layer.
- Loading, empty, and fatal-error states.

Use accessible buttons, labels, focus styles, and keyboard navigation. Test the shell at a narrow mobile width from the beginning.

## Task 4 — Separate server and UI state

### Theory

- **Architectural role:** Server state has query keys and fetch, invalidation, and reconciliation behavior. UI state has a local lifecycle and may disappear on reload. A streaming draft can be a temporary projection reconciled with a durable message later.
- **Why:** Copying server entities into a global store creates manual cache invalidation and stale data. Putting drawer or tab state into a query cache creates unnecessary network coupling.
- **GoClaw reference:** Compare chat server-state hooks in [`goclaw/ui/web/src/pages/chat/hooks/use-chat-messages.ts`](../../goclaw/ui/web/src/pages/chat/hooks/use-chat-messages.ts) with the UI store in [`goclaw/ui/web/src/stores/use-ui-store.ts`](../../goclaw/ui/web/src/stores/use-ui-store.ts) and streaming projection in [`goclaw/ui/web/src/stores/use-chat-messages-store.ts`](../../goclaw/ui/web/src/stores/use-chat-messages-store.ts).

Server state is authoritative and refreshable. UI state is temporary and local.

### Goal

Classify state by source of truth and recovery behavior before choosing hooks or stores.

### Guide to implement

Treat these as server state:

```text
agent settings, sessions, messages, memories, runs
```

Treat these as UI state:

```text
active drawer, selected tab, draft input, expanded tool card
```

Use a query/cache library or small typed hooks for server state. Do not make a global client store the only copy of server data.

## Task 5 — Implement authentication and clients

### Theory

- **Architectural role:** The HTTP client handles requests and responses. The WebSocket client is a connection state machine with a pending-RPC map and event dispatcher. A React provider exposes the connection to the component tree.
- **Why:** Components that call `fetch` or open sockets directly repeat authentication and retry logic and are harder to clean up. Session storage is only a demo compromise; production should prefer server-managed sessions to reduce token exposure.
- **GoClaw reference:** Inspect [`goclaw/ui/web/src/api/http-client.ts`](../../goclaw/ui/web/src/api/http-client.ts), [`goclaw/ui/web/src/api/ws-client.ts`](../../goclaw/ui/web/src/api/ws-client.ts), [`goclaw/ui/web/src/components/providers/ws-provider.tsx`](../../goclaw/ui/web/src/components/providers/ws-provider.tsx), and authentication state in [`goclaw/ui/web/src/stores/use-auth-store.ts`](../../goclaw/ui/web/src/stores/use-auth-store.ts).

### Goal

Centralize credential attachment, error mapping, RPC correlation, and reconnect policy in the client layer.

### Guide to implement

Create one HTTP client that adds the configured demo token and maps gateway error codes. Create one WebSocket client that:

1. Connects with authentication.
2. Tracks `connecting`, `connected`, `reconnecting`, and `failed` states.
3. Correlates request IDs with responses.
4. Dispatches typed events.
5. Reconnects with bounded exponential backoff.

For a local course project, the token may be entered by the user and kept in session storage. Explain that a production browser app should prefer a secure server-managed session over long-lived bearer tokens in JavaScript storage.

## Task 6 — Build agent settings

### Theory

- **Architectural role:** The server provides allowed options and current state, the form keeps a draft, the mutation sends a patch, and success invalidates the authoritative query. Browser validation improves UX, while server validation protects security and invariants.
- **Why:** Browser-only settings make gateway behavior differ from what the UI shows. Server-approved model and tool lists prevent the UI from enabling options the runtime does not support.
- **GoClaw reference:** Inspect agent schemas and types in [`goclaw/ui/web/src/schemas/agent.schema.ts`](../../goclaw/ui/web/src/schemas/agent.schema.ts) and [`goclaw/ui/web/src/types/agent.ts`](../../goclaw/ui/web/src/types/agent.ts), find agent pages with `rg --files goclaw/ui/web/src/pages/agents`, and review backend methods in [`goclaw/internal/gateway/methods/agents_update.go`](../../goclaw/internal/gateway/methods/agents_update.go).

### Goal

Implement one complete CRUD feature to verify typed forms, validation, mutations, and cache invalidation.

### Guide to implement

The settings page edits:

- Model name from the server-approved choices.
- Custom system prompt.
- Enabled local skills.
- Enabled tools supported by the server.

Load settings from the gateway, validate before submit, save through an HTTP endpoint, and show server field errors. After success, invalidate/reload the agent query. Do not make browser-local settings authoritative.

## Task 7 — Verify the shell

### Theory

- **Architectural role:** Browser tests check composition between routes, providers, clients, and state. Secret scans and build inspection verify that provider credentials do not cross the frontend boundary.
- **Why:** Many shell bugs appear only after a hard reload, deep link, or socket disconnect. Early accessibility and mobile checks prevent component APIs from becoming tied to desktop interactions.
- **GoClaw reference:** Review HTTP client tests in [`goclaw/ui/web/src/api/http-client.test.ts`](../../goclaw/ui/web/src/api/http-client.test.ts), the authentication guard in [`goclaw/ui/web/src/components/shared/require-auth.tsx`](../../goclaw/ui/web/src/components/shared/require-auth.tsx), and real pages and layouts under [`goclaw/ui/web/src`](../../goclaw/ui/web/src).

### Goal

Verify that the shell works across reloads, authentication failures, responsive layouts, and keyboard use, not only the happy path.

### Guide to implement

Check:

1. A page reload restores server-backed settings.
2. Invalid authentication shows a clear signed-out state.
3. Connection state is visible but not distracting.
4. Keyboard navigation reaches all shell controls.
5. Mobile layout keeps primary navigation and content usable.
6. No provider secret is present in the browser bundle or network requests.

This step is complete when the UI can connect to the gateway and manage the agent profile. Chat rendering is completed in Step 25.
