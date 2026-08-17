# Step 20 — Web UI Foundation and Application Shell

**Knowledge depth: 6/10**

There is no separate GoClaw author document devoted only to the React shell. Use [04 — Gateway and Protocol](../04-gateway-protocol.md), [18 — HTTP REST API](../18-http-api.md), and [19 — WebSocket RPC Methods](../19-websocket-rpc.md) as the contract knowledge behind the UI. Then study the implementation in `goclaw/ui/web/` as the code companion.

## The UI is a client of the gateway

The browser must not recreate agent logic. It owns presentation, local interaction state, routing, and connection recovery; the gateway remains the source of truth for sessions, configuration, tools, and long-running work.

```text
React route → query/mutation hook → HTTP or WebSocket client → gateway contract
```

## Establish the application shell

Create a React + TypeScript application with a root layout, authenticated route boundary, navigation, a responsive content area, a notification layer, and a central connection/auth state. GoClaw's web application uses React, Vite, Tailwind, Radix primitives, Zustand, and React Router; use the same separation even if you choose different libraries.

```text
src/
├── app/          route tree and providers
├── features/     chat, agents, settings, skills, traces
├── components/   reusable visual building blocks
├── lib/          API/WS clients and formatters
├── stores/       connection and UI state
└── i18n/         user-facing translations
```

## Separate server state from UI state

Sessions, agents, tasks, and configuration come from the server and must be refreshable. Drawer state, active tab, keyboard visibility, and optimistic input state are browser-local. Mixing these two kinds of state is a common source of stale interfaces.

## Design for small screens from the start

The agent UI is often used in a narrow browser. Use responsive layouts, accessible controls, and scroll regions that do not hide the chat input behind mobile browser chrome. These details make the course result feel like a product rather than a developer demo.

