# Step 15 — Web UI Foundation and Application Shell

**Knowledge depth: 6/10**

This step creates the React application shell, gateway clients, authentication state, and agent settings page.

## Task 1 — Define the UI boundary

### Theory

Use [04 — Gateway and Protocol](../04-gateway-protocol.md), [18 — HTTP REST API](../18-http-api.md), and [19 — WebSocket RPC Methods](../19-websocket-rpc.md) as the UI contracts.

The browser owns presentation, routing, temporary interaction state, and connection recovery. The gateway owns sessions, agent settings, tools, memory, and runs.

```text
React route → query/mutation hook → HTTP/WS client → gateway
```

### Practice guide

Do not copy agent logic into TypeScript. Define typed API/protocol models that match Steps 12–13 and keep them in one client layer.

## Task 2 — Create the React project

### Practice guide

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

### Practice guide

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

Server state is authoritative and refreshable. UI state is temporary and local.

### Practice guide

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

### Practice guide

Create one HTTP client that adds the configured demo token and maps gateway error codes. Create one WebSocket client that:

1. Connects with authentication.
2. Tracks `connecting`, `connected`, `reconnecting`, and `failed` states.
3. Correlates request IDs with responses.
4. Dispatches typed events.
5. Reconnects with bounded exponential backoff.

For a local course project, the token may be entered by the user and kept in session storage. Explain that a production browser app should prefer a secure server-managed session over long-lived bearer tokens in JavaScript storage.

## Task 6 — Build agent settings

### Practice guide

The settings page edits:

- Model name from the server-approved choices.
- Custom system prompt.
- Enabled local skills.
- Enabled tools supported by the server.

Load settings from the gateway, validate before submit, save through an HTTP endpoint, and show server field errors. After success, invalidate/reload the agent query. Do not make browser-local settings authoritative.

## Task 7 — Verify the shell

### Practice guide

Check:

1. A page reload restores server-backed settings.
2. Invalid authentication shows a clear signed-out state.
3. Connection state is visible but not distracting.
4. Keyboard navigation reaches all shell controls.
5. Mobile layout keeps primary navigation and content usable.
6. No provider secret is present in the browser bundle or network requests.

This step is complete when the UI can connect to the gateway and manage the agent profile. Chat rendering is completed in Step 16.
