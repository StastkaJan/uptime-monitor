# Implementation roadmap

M1 is complete and G1 passed, including hosted CI. See [verification evidence](verification.md) for results. Later milestones remain planned.

Implement in this order. Each milestone should be a small reviewable change or a few focused commits. The Docker Compose development configuration is prepared; the unchecked items describe future work.

## 1. Runnable shell

- [x] Initialize the Go module through Compose and add the executable entry point.
- [x] Serve an embedded dashboard shell and CSS, plus `GET /healthz`.
- [x] Add address configuration, server timeouts, and graceful shutdown.
- [x] Test the health route and full dashboard response.
- [x] Validate Compose configuration, startup, port access, and graceful shutdown with the real app.
- [x] Add CI for formatting, tests, vet, race checks, build, and Compose validation with aligned toolchain pins.

Done when `docker compose up` serves a page at `http://localhost:8080`, Compose-based tests/vet/build pass, and the binary serves embedded assets from another working directory.

## 2. Persistent monitors

- [ ] Initialize SQLite and its schema; add monitor validation and store operations.
- [ ] Add ordinary HTML create/edit forms, list/detail pages, and pause/resume.
- [ ] Reuse the monitor form and validate on the server.
- [ ] Add same-origin protection, body limits, and readable error responses.
- [ ] Test validation, mutations, persistence across reopening, and failure handling.

Done when a user can add and edit a monitor and find it after `docker compose down` followed by `docker compose up`, with JavaScript disabled and the named data volume preserved.

## 3. Checks and incidents

- [ ] Add one timed HTTP check and bounded scheduling with cancellation.
- [ ] Record checks and incident transitions atomically.
- [ ] Show pending, healthy, down, stale, and paused states.
- [ ] Show recent check latency and incident history; prune old check samples.
- [ ] Test deadlines, concurrency limits, pause/resume, shutdown, and incident transitions.

Done when a controlled target moving from 200 to 503 to 200 produces one resolved incident, repeated failures do not duplicate it, and concurrent checks pass race verification.

## 4. HTMX interactions

- [ ] Add a locally served pinned HTMX asset with its license.
- [ ] Poll the results fragment every five seconds.
- [ ] Enhance forms and add name/URL filtering using shared templates.
- [ ] Preserve form input, filter state, and focus during refreshes.
- [ ] Verify fragment responses and normal browser fallbacks.

Done when status updates and forms work without full reloads, while the core flow remains usable without JavaScript.

## 5. Showcase finish

- [ ] Add explicit local demo mode with a controlled failure/recovery action.
- [ ] Polish empty, loading, stale, and error states and responsive styling.
- [ ] Add a short demo walkthrough, screenshots, exact setup commands, and CI.
- [ ] Complete the browser acceptance list in [testing.md](testing.md).
- [ ] Document public-demo restrictions before any separate deployment task.

Done when someone can run the project from its README and reproduce an outage/recovery demonstration in a few minutes.

## Stop condition

The first showcase is complete after milestone 5 and its acceptance checks. Keep future ideas out of the implementation until explicitly selected: alerts, accounts, SSE, charts, distributed checking, and time-based uptime reporting.
