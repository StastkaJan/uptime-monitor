# Integration stages

These stages assemble the [roadmap](roadmap.md) into one executable and one internal application package. Integration here means connecting application components, not adding external services. G1 passed, including hosted CI; later gates remain pending. See [verification evidence](verification.md).

## Sequence and contracts

Implement the [tickets](tickets.md) in dependency order and integrate small, runnable changes. Every ticket inherits the preceding milestone's gate. Keep SQL in store files, HTTP handling in handler files, outbound requests in checker files, and process wiring in the executable. Share templates and validation across ordinary and HTMX requests.

Use [architecture.md](architecture.md) for routes, status rules, database records, limits, and response behavior. Extend concrete types only when the ticket needs them. Do not introduce additional services, generic repositories, or a separate frontend data model.

| Stage | Entry condition | Components connected | Exit gate |
| --- | --- | --- | --- |
| I1. Process and browser | Scaffold and available Compose toolchain | Compose -> executable -> HTTP -> embedded pages/assets -> CI | G1 |
| I2. Browser and persistence | G1 passed | Protected forms -> shared validation -> SQLite -> monitor views | G2 |
| I3. Background monitoring | G2 passed | Scheduler -> checker -> result/incident transaction -> status/history | G3 |
| I4. Progressive enhancement | G3 passed | HTMX -> existing handlers/templates -> filtered live results | G4 |
| I5. Demonstration and handoff | G4 passed | Demo controls -> controlled target -> checks -> incident UI -> operator guide | G5 |

## I1 / G1: Process and browser

**Tickets:** UPT-001 through UPT-003.

The executable owns startup/shutdown; the application owns routes and embedded views. `/healthz` reports process liveness independently of monitored targets. Native runs default to loopback; Compose publishes only to host loopback. Module, Compose, and CI use aligned Go versions.

Gate evidence:

- Validate Compose; run formatting verification, tests, vet, race checks on a supported platform, and build using the documented toolchain.
- Start via Compose and request the dashboard, stylesheet, and health endpoint.
- Run the built executable from another working directory and verify embedded assets.
- Stop the service and verify clean exit inside the 15-second Compose grace period.
- Verify CI executes the required checks on the runnable slice.

## I2 / G2: Browser and persistence

**Tickets:** UPT-004 through UPT-007.

Mutations share validation, limited request bodies, same-origin protection, and parameterized SQL. Ordinary successful forms redirect; invalid forms preserve values. SQLite uses one connection, foreign keys, and a bounded busy timeout. No outbound checks run yet.

Gate evidence:

- With JavaScript disabled, create a monitor, view it, edit it, pause it, and resume it.
- Reject invalid targets, intervals, over-limit creation, and cross-origin mutations; escape user text.
- Prove persistence across temporary-database reopening and report failed writes without false success.
- Run `docker compose down` then `docker compose up`; verify settings survive in the named volume. Never use `down --volumes` for this check.
- Pass the routine checks in [testing.md](testing.md).

## I3 / G3: Background monitoring

**Tickets:** UPT-008 through UPT-012.

The scheduler uses a shared HTTP client without holding a database transaction. One store operation commits each result and incident transition together. Checks run without a browser, with four workers maximum and one in-flight check per monitor. Shutdown cancels checks, waits for workers, and closes SQLite last; cancellation must not create an outage.

Gate evidence:

- Against a controlled target, persist 200 -> 503 -> 503 -> 200 and assert one resolved incident. A first-ever failure also opens an incident.
- Verify pending, healthy, down, stale, and paused behavior; retain open incidents during pauses and show the observation gap.
- Prove concurrency bounds, no duplicate dispatch, pause/resume, deadlines, and restart without replaying missed runs.
- Migrate a populated M2 database without losing monitors; prove rollback of a failed result/incident write.
- Verify history survives restart, cleanup is bounded, and retention preserves the last observation needed for status and URL edit protection.
- Stop during active work, verify clean reopen, and pass routine checks plus race verification.

## I4 / G4: Progressive enhancement

**Tickets:** UPT-013 through UPT-015.

The dashboard and refresh route use the same results template. Keep forms and filters outside the polling region. Responses that vary by `HX-Request` send `Vary: HX-Request`; full-page and history-restoration requests receive complete documents. Poll every five seconds and include the applied filter. Coordinate search and polling so delayed responses cannot overwrite newer results.

Gate evidence:

- Verify status updates and filtering without page reloads; filters remain applied across refreshes.
- Submit valid and invalid enhanced forms; preserve entered values and show unexpected failures while retaining 5xx status codes.
- Preserve focus, unsaved form values, and filter input over multiple polls.
- Exercise back/forward navigation and direct loads, then repeat the core form flow with JavaScript disabled.
- Pass handler/fragment checks and routine verification.

## I5 / G5: Demonstration and handoff

**Tickets:** UPT-016 through UPT-018.

Explicit demo mode seeds one controlled monitor without duplication and uses a same-process target. Only protected POST actions change it between 200 and 503. Demo routes are absent when disabled. The UI works on narrow screens and with a keyboard and explains status using text and color.

Gate evidence:

- Follow the README through outage/recovery, repeated failure without duplicate incidents, and restart without duplicate seeding.
- Verify demo isolation: incidental GETs and unrelated monitors cannot alter demo state; disabled mode exposes no demo actions.
- Complete all eleven browser checks in [testing.md](testing.md), recording environment and outcomes.
- Pass final CI and routine/race checks; capture screenshots of the actual application.
- Verify exact setup commands and document public-demo restrictions without implementing public deployment.

## Evidence and recovery

For each gate, record the tested commit, environment/toolchain, command results, manual scenarios, and unavailable checks in the implementing change or PR. A blocked check is not a passed gate. Update ticket and roadmap checkboxes only when acceptance passes.

If integration fails, fix or revert the smallest responsible change and repeat affected checks before advancing. Preserve the named SQLite volume. Before a schema change, stop the application and take the documented backup; test migrations on populated temporary databases. Do not assume an older executable can read a newer schema: use the tested recovery procedure or restore the backup before returning to an older version.
