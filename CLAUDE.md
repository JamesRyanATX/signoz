# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

SigNoz is an open-source observability platform (logs, metrics, traces, APM) built on OpenTelemetry with ClickHouse as the primary telemetry datastore. There are two editions: **community** (SQLite for metadata) and **enterprise** (PostgreSQL + licensing). Entry points are `cmd/community/` and `cmd/enterprise/`.

## Development Environment

Start infrastructure (ClickHouse + OTel Collector) before running the backend:

```bash
make devenv-up          # starts ClickHouse + signoz-otel-collector via Docker Compose
make devenv-clickhouse  # ClickHouse only
make devenv-postgres    # PostgreSQL only (enterprise)
```

## Backend (Go 1.24)

```bash
make go-run-enterprise  # run enterprise server with race detector (connects to local ClickHouse)
make go-run-community   # run community server
make go-test            # go test -race ./...
```

To run a single Go test package:
```bash
go test -race ./pkg/sqlstore/...
```

Formatting and linting (checked in CI):
```bash
gofmt -w .
golangci-lint run
```

**Linting rules enforced by `.golangci.yml`:**
- Use `slog` — `go.uber.org/zap` is banned
- Use `github.com/SigNoz/signoz/pkg/errors` — stdlib `errors` package is banned
- `pkg/query-service/` and `ee/query-service/` are excluded from linting (legacy code)
- `slog` keys must be `snake_case`, messages must be lowercase

## Frontend (React/TypeScript)

```bash
cd frontend
yarn install
yarn dev           # development server
yarn jest          # run all tests
yarn jest path/to/test.test.tsx   # run single test file
yarn lint          # ESLint
yarn fmt           # Prettier check
yarn prettify      # Prettier fix
yarn build         # production build
```

## Architecture

```
cmd/
  community/       # community edition entry point
  enterprise/      # enterprise edition entry point
pkg/               # shared Go packages
  query-service/   # main API server + query engine (legacy, excluded from linting)
  modules/         # feature modules: dashboard, alerts, org, user, role, etc.
  sqlstore/        # metadata DB abstraction (SQLite / PostgreSQL)
  alertmanager/    # alert management
  apis/            # API definitions
  http/            # HTTP handlers
  cache/           # caching layer
  factory/         # pluggable provider pattern (alerts, licensing, web, cache, email)
ee/                # enterprise-only code
  query-service/   # enterprise extensions (legacy)
  licensing/       # license system
  sqlstore/        # enterprise SQL store (PostgreSQL)
  zeus/            # license verification service
  authz/           # authorization (OpenFGA)
frontend/
  src/
    api/           # API client methods
    components/    # reusable UI components
    container/     # smart/page-level components
    hooks/         # custom React hooks
    lib/           # utilities
    modules/       # frontend feature modules
    mocks-server/  # MSW mock server
    tests/         # test utilities (test-utils harness)
```

Community uses SQLite for metadata; enterprise uses PostgreSQL. ClickHouse stores all telemetry (logs, metrics, traces) in both editions. The `pkg/factory/` package provides pluggable providers wired up at startup.

## Frontend Testing Conventions

Always import from the internal harness — never directly from `@testing-library/react`:

```ts
import { render, screen, userEvent, waitFor } from 'tests/test-utils';
import { server, rest } from 'mocks-server/server';
```

Use MSW to override API responses per test:
```ts
server.use(
  rest.get('*/api/v1/foo', (_req, res, ctx) => res(ctx.status(200), ctx.json({ ok: true })))
);
```

Use `userEvent` (not `fireEvent`) for user interactions; always `await` them:
```ts
const user = userEvent.setup({ pointerEventsCheck: 0 });
await user.click(screen.getByRole('button', { name: /save/i }));
```

Mock decision: global mocks (in `__mocks__/`) for infra used in 20+ files (react-router-dom, react-query, antd, browser APIs). Local mocks for business logic, specific API responses, and error scenarios.

No global fake timers — use per-test fake timers only for debounce/throttle, then restore with `jest.useRealTimers()`.

All mock functions must be typed with `jest.MockedFunction<T>`. No `any` in test files.

## Integration Tests (Python/Poetry)

```bash
make py-test   # cd tests/integration && poetry run pytest --basetemp=./tmp/ -vv --capture=no src/
```

## Docker / Production Builds

```bash
make docker-build-community   # build community Docker image
make docker-build-enterprise  # build enterprise Docker image
```
