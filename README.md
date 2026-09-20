# Uptime Monitor

A Go + HTMX showcase: add a website, watch its health, and inspect outages as they happen.

**Status: planning scaffold.** The folders and implementation directions are prepared; application code, dependencies, and tests have not been created yet.

## The idea

Build a small, local-first dashboard for a single operator monitoring up to 50 HTTP endpoints. Demonstrate Go through concurrent checks, cancellation, and durable results; demonstrate HTMX through live updates, forms, and searchable history.

The main demo is a complete story: add a monitor, see it become healthy, make a controlled demo endpoint fail, watch an incident open, then restore it and see the incident close.

## First version

- Add and edit a monitor's name, URL, and check interval; pause or resume it.
- Show healthy, down, pending, stale, and paused states with text and color.
- Refresh dashboard results without refreshing the whole page.
- Show recent response times, check results, and incident history per monitor.
- Preserve monitors and history across restarts.
- Provide an explicit local demo mode with a controllable healthy/failing endpoint.

Keep email alerts, accounts, billing, distributed workers, and elaborate charts outside the first version.

## Default stack

| Concern | Choice |
| --- | --- |
| Server and routing | Go standard library, `net/http` |
| HTML | `html/template`, embedded templates and assets |
| Interactions | HTMX, initially polling every five seconds |
| Styling | Plain CSS, responsive HTML, native form controls |
| Storage | SQLite through `database/sql` and `modernc.org/sqlite` |
| Checks | Shared HTTP client, bounded goroutines, contexts |
| Tests | `testing`, `net/http/httptest`, temporary SQLite files |
| Delivery | One executable plus a writable data directory |

## Folder map

```text
uptime-monitor/
  AGENTS.md                 Working rules for future implementation
  cmd/uptime/               Executable entry point
  internal/uptime/          Application code, initially one package
    templates/             Full pages and reusable HTML fragments
    static/                CSS and a pinned local HTMX asset
  docs/
    architecture.md        Boundaries, data, routes, operational defaults
    development.md         Setup, splitting work, reuse, dependencies
    testing.md             Verification strategy and acceptance checks
    roadmap.md             Ordered implementation slices
```

Read [architecture](docs/architecture.md), then follow the [roadmap](docs/roadmap.md). Use [development directions](docs/development.md) and the [testing guide](docs/testing.md) while implementing.

The first task is a runnable page and health endpoint, as described in milestone 1. No application run command works until that milestone is implemented.
