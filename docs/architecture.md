# Architecture

## Shape

One process owns an HTTP server, a scheduler, and a SQLite connection pool. HTML is the browser contract. The scheduler records checks independently of whether anyone has the dashboard open.

```text
Browser -- forms / HTMX --> Go handlers --> SQLite
Browser <-- HTML fragments -- templates <-- query results

Scheduler --> bounded HTTP checks --> monitored endpoints
          --> result + incident transaction --> SQLite
```

Begin with `cmd/uptime` for wiring and `internal/uptime` for application code. File boundaries provide enough separation for this size of project. Do not add a JSON API or additional service processes for the browser.

## Technologies

- Use Go 1.27.1, as pinned in `compose.yaml`; align the module and CI toolchain when code is added, and commit `go.mod` and `go.sum`. Update these pins together when upgrading. Standard-library [net/http](https://pkg.go.dev/net/http) supplies routing, clients, and server lifecycle support.
- Use [html/template](https://pkg.go.dev/html/template) for escaped output. Embed templates and static assets so the executable can run from any directory. Never turn user input into trusted `template.HTML`.
- Use a pinned stable HTMX release. [HTMX supports polling and HTML fragment responses](https://htmx.org/docs/); poll the results container every five seconds initially. Add SSE only if a later requirement needs lower latency.
- Use `database/sql` with [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite), a CGo-free SQLite driver. Choose a compatible stable version when implementing storage and commit resolved dependency versions.
- Use plain CSS and the Go standard library's structured logger. A frontend build pipeline, ORM, router framework, and message broker are unnecessary for the initial scope.

## Docker Compose development

The root `compose.yaml` defines one `app` service using the [official Go image](https://github.com/docker-library/official-images/blob/master/library/golang), pinned to `golang:1.27.1-bookworm`. The Debian-based toolchain supports Go builds and race tests. `GOTOOLCHAIN=local` prevents silent toolchain downloads; update the image when dependencies require a newer Go version.

Bind-mount source at `/workspace`; store SQLite under `/data` in the `uptime-data` named volume. Separate volumes cache downloaded modules and build output. SQLite remains inside the application, with no database service. The same service runs one-off setup and verification commands via `docker compose run --rm --no-deps app ...`.

The startup command builds a temporary executable and replaces the shell with it so shutdown reaches the application. Compose provides an init process and a 15-second stop grace period; implement an application shutdown deadline shorter than that. Restart the service after source or embedded asset changes to rebuild. There is no automatic file watcher.

The server binds `0.0.0.0:8080` inside the container; the published host port is restricted to `127.0.0.1:8080`. This follows the [Compose service configuration](https://docs.docker.com/reference/compose-file/services/). Container loopback refers to the container itself, so the built-in demo endpoint stays in the same process.

This is a development environment with source and compiler available. Add a multi-stage Dockerfile with a non-root runtime only when preparing a deployable image; it is not required to use the current Compose toolchain.

## Data and rules

| Record | Initial fields |
| --- | --- |
| Monitor | ID, name, URL, interval seconds, paused, created time |
| Check | ID, monitor ID, start time, duration milliseconds, HTTP status if available, success, bounded error description |
| Incident | ID, monitor ID, opened time, optional resolved time |

Store timestamps in UTC using one consistent representation. Index checks by `(monitor_id, started_at)` and incidents by monitor and opened time. Enforce foreign keys and at most one open incident per monitor with a partial unique index.

A successful check receives HTTP 200-299. Redirects count as failures initially: this keeps the target explicit and the behavior easy to demonstrate. Use GET, measure elapsed time until response headers arrive, and close the response body without storing page content. Label the measurement as response-header latency.

One failed check opens an incident; repeated failures keep it open. The next success closes it. A first-ever failure also opens an incident. Save the check and incident transition in one transaction so a partial write cannot disagree with the dashboard. Shutdown cancellation is not an endpoint outage and must not create a failed check.

Show pending until the first check, paused when disabled, and stale if the latest result is older than twice the configured interval plus the check timeout. Stale means the observation is old; it does not create an outage. Retain an open incident while paused and identify that observation gap on the detail page. For MVP, changing a URL is allowed only before the first check; afterward create another monitor so histories are not mixed.

Keep recent checks for seven days with a bounded daily cleanup; retain incidents. Do not call a fraction of successful samples time-based uptime. Defer uptime percentages until the treatment of gaps and pauses is explicitly defined.

## Scheduling and storage

Start with these documented constants: 50 monitors maximum, 60-second default interval, 10-3600-second allowed interval, five-second check timeout, four concurrent checks, and a one-second scheduler tick.

Track next-due times and in-flight monitor IDs in memory. Claim only available worker capacity, skip a monitor already being checked, and schedule the next run from completion time. On startup, monitors become due once; do not replay checks missed during downtime. On pause, prevent new checks; an existing check may finish. Load current settings when dispatching a check.

Reuse one HTTP client and transport across checks. Give each request a context deadline and propagate process shutdown. Do not hold a database transaction while making an outbound request.

For this small deployment, start with one SQLite connection (`SetMaxOpenConns(1)`), foreign keys enabled, and a bounded busy timeout. Keep transactions short and close query rows before issuing another query. Serialized database access is an intentional ceiling; change it only after measuring contention.

## Browser routes and rendering

| Route | Purpose |
| --- | --- |
| `GET /` | Full dashboard with embedded results fragment |
| `GET /monitors/results` | Filtered results fragment, also used for polling |
| `GET /monitors/new` | Add form |
| `POST /monitors` | Create monitor |
| `GET /monitors/{id}` | Detail page with checks and incidents |
| `GET /monitors/{id}/edit` | Edit form |
| `POST /monitors/{id}` | Update monitor |
| `POST /monitors/{id}/pause` | Explicit desired paused state |
| `GET /healthz` | Process liveness, independent of monitored endpoints |

Use normal HTML forms and redirects first. Enhance with HTMX using the same validation and named templates. A full-page request must always receive a complete document. If the same URL varies by `HX-Request`, set `Vary: HX-Request` and keep history restoration working.

Keep filters and forms outside the polling region; include current filter values in refresh requests. Start with interval polling and a submit-based search so competing search requests do not overwrite newer results. Add typing-triggered search only with explicit request synchronization.

Return form field errors without losing entered values. Ordinary requests can receive a full form with 422; HTMX requests may receive the form fragment with 200 for a validation failure, avoiding version-dependent error-swap configuration. Unexpected failures remain 5xx and show a visible generic error; log details on the server.

Use semantic HTML, explicit labels, keyboard-accessible actions, visible focus, and text alongside status colors. Do not replace a focused control during polling. Basic navigation and forms must work without JavaScript.

## Runtime and demo defaults

Native defaults: `UPTIME_ADDR=127.0.0.1:8080`, `UPTIME_DB=data/uptime.db`, and `UPTIME_DEMO=false`. Compose explicitly overrides the address to `0.0.0.0:8080` and database path to `/data/uptime.db`. It interpolates `UPTIME_DEMO` from the host environment or a local Compose `.env` file. The Go application reads environment variables directly and does not load `.env` files itself. Validate configuration at startup and log actionable errors. Set HTTP server header and idle timeouts and a shutdown deadline below Compose's 15-second grace period.

In explicit local demo mode, seed one demo monitor and mount a controlled endpoint plus a POST action to switch it between 200 and 503. Keep its state synchronized across requests. Other monitors and incidental GET requests must not change demo state.

Validate URL scheme, host, lengths, and intervals; allow only HTTP/HTTPS and reject embedded credentials. Use parameterized SQL, limited form bodies, and same-origin/CSRF protection for mutations even on localhost.

A public portfolio demo should use fixed developer-controlled targets and read-only visitor access. Do not expose arbitrary URL creation from this local scaffold: user-selected outbound targets require authentication and connection-time SSRF protection covering resolved addresses and redirects. Demo-only loopback access must be narrowly scoped to the known controlled target. These are release conditions for a public deployment, not extra infrastructure for the local MVP.
