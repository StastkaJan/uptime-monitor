# Implementation roadmap

M1 is complete and G1 passed, including hosted CI. See [verification evidence](verification.md) for results. Later milestones remain planned.

Implement in this order. Each milestone contains focused, reviewable tickets. The Docker Compose development configuration is prepared; M1 is complete and later milestones remain planned.

Use this document for scope and progress, [integration stages](integration-stages.md) for component assembly and acceptance gates, and [implementation tickets](tickets.md) for dependencies and actionable acceptance criteria. Architecture owns behavior and defaults; the testing guide owns verification commands.

| Milestone | Outcome | Tickets | Integration gate |
| --- | --- | --- | --- |
| M1. Runnable shell | Local page and reliable process lifecycle | UPT-001 through UPT-003 | G1: startup, assets, shutdown, CI |
| M2. Persistent monitors | Ordinary browser monitor management | UPT-004 through UPT-007 | G2: protected mutations and restart persistence |
| M3. Checks and incidents | Bounded checks and durable outage history | UPT-008 through UPT-012 | G3: failure/recovery and race verification |
| M4. HTMX interactions | Live results, enhanced forms, filtering | UPT-013 through UPT-015 | G4: live interaction and browser fallback parity |
| M5. Showcase finish | Reproducible local demonstration | UPT-016 through UPT-018 | G5: complete browser acceptance and handoff |

Complete each gate before advancing. Dates and estimates remain unset until an implementer and capacity are known. The ticket checklist in [tickets.md](tickets.md) tracks individual completion; the items below summarize milestone scope.

## 1. Runnable shell

- [x] Initialize the Go module through Compose and add the executable entry point.
- [x] Serve an embedded dashboard shell and CSS, plus `GET /healthz`.
- [x] Add address configuration, server timeouts, and graceful shutdown.
- [x] Test the health route and full dashboard response.
- [x] Validate Compose configuration, startup, port access, and graceful shutdown with the real app.
- [x] Add CI for formatting, tests, vet, race checks, build, and Compose validation with aligned toolchain pins.

Done at G1 when `docker compose up` serves a page at `http://localhost:8080`, Compose-based verification and CI pass, shutdown fits the configured grace period, and the binary serves embedded assets from another working directory.

## 2. Persistent monitors

- [ ] Initialize SQLite and its schema; add monitor validation and store operations.
- [ ] Add ordinary HTML create/edit forms, list/detail pages, and pause/resume.
- [ ] Reuse the monitor form and validate on the server.
- [ ] Add same-origin protection, body limits, and readable error responses.
- [ ] Test validation, mutations, persistence across reopening, and failure handling.

Done at G2 when a user can create, edit, pause, and resume a monitor and find it after `docker compose down` followed by `docker compose up`, with JavaScript disabled and the named data volume preserved.

## 3. Checks and incidents

- [ ] Add one timed HTTP check and bounded scheduling with cancellation.
- [ ] Record checks and incident transitions atomically.
- [ ] Migrate existing monitor data safely and prevent URL changes after the first recorded check.
- [ ] Show pending, healthy, down, stale, and paused states.
- [ ] Show recent check latency and incident history; prune old check samples.
- [ ] Test deadlines, concurrency limits, pause/resume, shutdown, and incident transitions.

Done at G3 when a controlled target moving from 200 to 503 to 200 produces one resolved incident, repeated failures do not duplicate it, restart preserves history, retention preserves current observation state, and concurrent checks pass race verification.

## 4. HTMX interactions

- [ ] Add a locally served pinned HTMX asset with its license.
- [ ] Poll the results fragment every five seconds.
- [ ] Enhance forms and add name/URL filtering using shared templates.
- [ ] Preserve form input, filter state, and focus during refreshes.
- [ ] Verify fragment responses and normal browser fallbacks.

Done at G4 when status updates, forms, and filtering work without full reloads, while the core flow remains usable without JavaScript and polling preserves input and focus.

## 5. Showcase finish

- [ ] Add explicit local demo mode with a controlled failure/recovery action.
- [ ] Polish empty, loading, stale, and error states and responsive styling.
- [ ] Add a short demo walkthrough, screenshots, and exact setup commands; verify CI on the completed application.
- [ ] Complete the browser acceptance list in [testing.md](testing.md).
- [ ] Document public-demo restrictions before any separate deployment task.

Done at G5 when someone can run the project from its README and reproduce an outage/recovery demonstration in a few minutes, all browser acceptance checks pass, and every gate has recorded evidence.

## Stop condition

Mark work complete only after its acceptance criteria and the shared ticket definition of done pass. Record verification with the implementing change; an unavailable required check leaves its gate open. Update setup instructions as each slice becomes runnable.

The first showcase is complete after milestone 5 and its acceptance checks. Keep future ideas out of the implementation until explicitly selected: alerts, accounts, SSE, charts, distributed checking, and time-based uptime reporting.
