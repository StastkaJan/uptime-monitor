# Implementation tickets

Repository-local backlog for the [roadmap](roadmap.md). These are issue-ready descriptions, not published tracker issues. UPT-001 through UPT-003 are complete and G1 passed. Remaining tickets are planned. See [verification evidence](verification.md). Check off a ticket only after acceptance passes; [integration stages](integration-stages.md) define milestone gates.

## Shared definition of done

- Implement the ticket's behavior and necessary failure handling in the existing application package, using the standard library first.
- Add behavior tests beside the implementation with controlled HTTP servers and real temporary SQLite databases as relevant. Never test against live websites or the development database.
- Run formatting, tests, vet, and build through Compose as specified in [testing.md](testing.md). Run race checks for concurrent behavior; CI also runs them on a supported platform. Record exact results and unavailable checks.
- Update setup instructions for new runnable behavior, keep ticket/roadmap checkboxes accurate, and supply evidence for any gate closed by the ticket.
- Preserve ordinary browser behavior, escaped output, existing data, and the named volume. Add no deployment, external notifications, or unrelated dependencies.

Dependencies below are prerequisites within each milestone. Every ticket also requires the preceding milestone's integration gate. Acceptance criteria supplement this definition of done and the rules in [architecture.md](architecture.md).

## M1: Runnable shell

- [x] UPT-001
- [x] UPT-002
- [x] UPT-003

### UPT-001: Serve the first embedded page

**Depends on:** none. **Scope:** module, entry point, routes, embedded templates and CSS.

Initialize the module through Compose and serve `GET /` and `GET /healthz`. Keep wiring in the executable and handlers/rendering in the application package.

Acceptance:

- Dashboard requests return complete HTML with local embedded CSS; health requests return success independently of monitored endpoints.
- The module toolchain aligns with Compose; no SQLite or HTMX dependency is added yet.
- Handler tests verify both routes; a binary run from another working directory still serves its assets.
- `docker compose up` starts the application using the existing service.

### UPT-002: Configure and stop the process reliably

**Depends on:** UPT-001. **Scope:** configuration and lifecycle.

Read and validate runtime configuration, set server timeouts, and implement signal-driven shutdown with a deadline shorter than Compose's 15-second grace period.

Acceptance:

- Defaults/overrides follow architecture; native runs bind loopback and Compose publishes only to host loopback.
- Invalid configuration fails startup with an actionable message; header and idle timeouts are nonzero.
- SIGTERM stops the real Compose app within its grace period; verification leaves no stray process.
- Tests cover valid overrides, invalid values, and orderly server shutdown.

### UPT-003: Establish the verification pipeline

**Depends on:** UPT-002. **Scope:** CI and verification documentation. **Closes:** G1.

Add CI with the aligned pinned toolchain and existing Compose workflow; keep commands usable locally.

Acceptance:

- CI validates Compose, checks formatting without rewriting source, and runs tests, vet, race checks, and build.
- Race verification uses a supported platform; tests never touch the persistent application volume.
- README/setup instructions accurately describe the runnable shell.
- Record G1 startup, assets, shutdown, and automated-check evidence; verify actual CI execution before closing G1.

## M2: Persistent monitors

- [ ] UPT-004
- [ ] UPT-005
- [ ] UPT-006
- [ ] UPT-007

### UPT-004: Store and validate monitors

**Depends on:** none within M2. **Scope:** SQLite initialization, monitor model, validation, store operations.

Add the pinned SQLite driver, initial monitor schema, parameterized queries, and shared validation. Initialize the configured directory/database without destroying existing data.

Acceptance:

- Fresh initialization and reopening succeed; monitors survive reopening a temporary database.
- Configure one connection, foreign keys, a bounded busy timeout, and consistent UTC timestamps.
- Reject unsupported schemes, missing hosts, credentials, oversized fields, and intervals outside 10-3600 seconds; default to 60 seconds.
- Document concrete name/URL length limits and enforce the 50-monitor cap even under concurrent creates.
- Database failures return errors; dependency versions and setup steps are documented.

### UPT-005: Create and list monitors through ordinary forms

**Depends on:** UPT-004. **Scope:** dashboard list, new form, create handler, shared mutation protection.

Implement `GET /monitors/new` and `POST /monitors`, then list stored monitors. Introduce protection with this first mutation route and reuse it thereafter.

Acceptance:

- With JavaScript disabled, a valid submission creates one monitor and redirects to a full page showing it.
- Invalid input returns 422 with field errors and entered values retained; failed writes produce visible generic 5xx errors and detailed server logs.
- A documented body limit rejects oversized submissions; cross-origin mutations are rejected.
- Tests cover creation, validation, cap enforcement, write failure, escaping, parameterized values, and mutation protection.

### UPT-006: View and edit a monitor

**Depends on:** UPT-005. **Scope:** detail page and edit routes.

Implement `GET /monitors/{id}`, `GET /monitors/{id}/edit`, and `POST /monitors/{id}` using the shared form and validator. Initially show settings on the detail page; history arrives in M3.

Acceptance:

- Name, URL, and interval edits persist without JavaScript; URL editing is available before any check exists.
- Invalid edits retain input and stored values. Unknown IDs return 404; malformed IDs produce controlled client errors.
- Tests cover validation, persistence, failed writes, and the reused request limits/protection.
- UPT-009 will lock URL changes after the first persisted check; do not invent placeholder history for this ticket.

### UPT-007: Pause, resume, and prove durable settings

**Depends on:** UPT-006. **Scope:** pause/resume and restart acceptance. **Closes:** G2.

Implement `POST /monitors/{id}/pause` with an explicit desired paused state, reusing protection and ordinary form redirects.

Acceptance:

- Repeating a pause/resume request preserves the requested state instead of toggling it again.
- Dashboard/detail views display the persisted setting; failed writes do not show false success.
- The complete create/edit/pause/resume flow works with JavaScript disabled.
- Temporary-database reopen tests and the real Compose down/up check preserve settings without deleting volumes; record G2 evidence.

## M3: Checks and incidents

- [ ] UPT-008
- [ ] UPT-009
- [ ] UPT-010
- [ ] UPT-011
- [ ] UPT-012

### UPT-008: Perform one bounded HTTP check

**Depends on:** none within M3. **Scope:** checker and tests.

Use a shared HTTP client/transport for GET checks with a five-second deadline. Measure response-header latency, close response bodies, and bound error descriptions.

Acceptance:

- Controlled 2xx responses succeed; 503 and redirects fail, and redirects are not followed.
- A blocked target reaches its deadline; parent cancellation ends promptly and is distinguishable from endpoint failure.
- Tests cover transport errors, statuses, deadlines, cancellation, and body closure without live websites.
- No page content is stored and no database transaction spans an outbound request.

### UPT-009: Commit checks and incidents together

**Depends on:** UPT-008. **Scope:** schema migration, transactional persistence, URL edit lock.

Add indexed checks/incidents with foreign keys and a partial unique index allowing one open incident per monitor. Preserve M2 data through a numbered migration.

Acceptance:

- First failure opens an incident, repeated failures retain it, and recovery resolves it; successful checks alone create no incident.
- Result and incident changes commit atomically; forced write failure leaves no partial result or inconsistent transition.
- Shutdown cancellation never persists a failed check or opens an incident.
- After the first persisted check, reject URL changes even by direct POST; name/interval remain editable. Enforce the rule at the write boundary against concurrent result recording.
- Migration tests start with populated M2 storage and preserve monitors across reopening; document backup and recovery.

### UPT-010: Schedule checks and integrate shutdown

**Depends on:** UPT-009. **Scope:** scheduler, worker lifecycle, process wiring.

Use a one-second tick, four concurrent checks maximum, and one in-flight check per monitor. Track due times in memory, load current settings at dispatch, and schedule the next run from completion.

Acceptance:

- Startup makes monitors due once without replaying downtime; new and resumed monitors become eligible.
- Pause prevents new dispatch while allowing existing work to finish; edits apply to subsequent dispatches.
- Channel-synchronized tests and explicit timestamps prove concurrency bounds, no duplicates, and checking without a browser open.
- Shutdown cancels requests, waits within its deadline, and closes SQLite last without recording cancellation as an outage.
- Failed persistence is logged and does not kill scheduling or invent stored observations; race checks pass.

### UPT-011: Show current state and history

**Depends on:** UPT-010. **Scope:** status rules, store queries, dashboard/detail views.

Render status, recent checks with response-header latency, and incident history from durable observations.

Acceptance:

- No checks means pending; otherwise show healthy/down unless stale or paused. Paused takes precedence; stale starts after twice the interval plus the five-second timeout.
- Staleness never opens an incident; open incidents survive pauses and detail pages explain the observation gap.
- History queries are bounded and distinguish HTTP failures from transport errors; provide simple server-side pagination where necessary to access older retained history.
- Tests cover status boundaries, empty history, pause/recovery, and reopening; statuses have text labels.
- A controlled 200 -> 503 -> 503 -> 200 sequence appears as one resolved incident.

### UPT-012: Bound retention without losing observation state

**Depends on:** UPT-011. **Scope:** daily cleanup and lifecycle wiring. **Closes:** G3.

Delete samples older than seven days in bounded daily batches; retain incidents and enough observation state for stale status and the URL edit lock. The simplest option is retaining the latest check row per monitor and documenting that exception.

Acceptance:

- Timestamp-driven tests prove the cutoff, a documented per-run deletion cap, and preservation of recent rows and incidents.
- A monitor with only old checks stays stale/paused after cleanup; it cannot revert to pending or regain URL editing.
- Cleanup observes cancellation, uses short transactions, and finishes before database shutdown.
- Document the selected retention approach; pass routine/race checks and record all G3 evidence.

## M4: HTMX interactions

- [ ] UPT-013
- [ ] UPT-014
- [ ] UPT-015

### UPT-013: Poll shared dashboard results

**Depends on:** none within M4. **Scope:** local HTMX asset/license, results route and fragment.

Vendor pinned HTMX and implement `GET /monitors/results` using the same named template included by `GET /`. Poll every five seconds.

Acceptance:

- Status updates appear without full reloads; dashboard requests still return complete documents.
- Replace only the results region, preserve focused controls, and keep form fields outside it.
- Handler tests verify the fragment; browser checks verify polling and usable behavior when refresh fails.
- No CDN or frontend build is required; the dashboard remains usable without JavaScript.

### UPT-014: Enhance forms with shared server behavior

**Depends on:** UPT-013. **Scope:** enhanced creation, editing, pause/resume, response selection.

Add HTMX attributes and fragment rendering to existing handlers while retaining shared validation and ordinary form redirects.

Acceptance:

- Enhanced create/edit and pause/resume update the relevant view without full reloads or duplicate mutations.
- Ordinary invalid forms return full HTML with 422; enhanced invalid forms may return the shared fragment with 200, preserving values/errors.
- Unexpected failures remain 5xx and are visible; all mutation protection remains effective.
- Responses varying by `HX-Request` send `Vary: HX-Request`; direct loads and history restoration receive complete documents.
- Test both request modes and verify navigation, validation, focus, and error visibility in the browser.

### UPT-015: Filter results during live updates

**Depends on:** UPT-014. **Scope:** name/URL filtering and request coordination. **Closes:** G4.

Add submit-based search using a query parameter and the shared results template. Carry the applied filter into subsequent polling requests.

Acceptance:

- Name/URL matching and empty queries behave consistently in full-page and enhanced requests; filter values remain SQL data.
- Filters/forms stay outside polling replacements; input, unsaved form values, and focus survive refreshes.
- Delayed search/poll responses cannot overwrite newer filtered results; prefer HTMX synchronization to custom JavaScript.
- Empty matches have a useful state; clearing the filter restores all monitors.
- Filtering works without JavaScript, back/forward behavior is correct, and all G4 checks pass.

## M5: Showcase finish

- [ ] UPT-016
- [ ] UPT-017
- [ ] UPT-018

### UPT-016: Provide a controlled local outage demo

**Depends on:** none within M5. **Scope:** demo seeding, target, protected controls.

When `UPTIME_DEMO=true`, seed one known monitor and mount a same-process endpoint whose synchronized state switches between 200 and 503 through explicit POST actions.

Acceptance:

- Disabled mode exposes no demo routes and creates no demo monitor; enabled startup is idempotent across restarts.
- Healthy -> failing -> healthy follows the real scheduler/checker/store path and yields one resolved incident; repeated failure does not duplicate it.
- Controls reuse same-origin protection and request limits; incidental GETs and unrelated monitors cannot alter demo state.
- Concurrent demo requests pass race checks; the known target works under Compose and native defaults.
- Document enable/disable commands, state reset behavior, and treatment of an existing demo monitor when disabled.

### UPT-017: Finish accessible and responsive states

**Depends on:** UPT-016. **Scope:** templates and plain CSS polish.

Complete empty, loading, stale, validation, and unexpected-error states on the dashboard, forms, and detail pages.

Acceptance:

- Controls have labels, keyboard access and visible focus work, and status meaning does not rely on color alone.
- Narrow/mobile layouts remain usable without clipped actions or unreadable history.
- Loading does not erase input; refresh/mutation failures are visible with a usable retry path.
- Browser checks cover JavaScript enabled/disabled and the demo flow; add behavior regression tests where logic changed.

### UPT-018: Verify and document the finished showcase

**Depends on:** UPT-017. **Scope:** final acceptance, setup guide, walkthrough/screenshots. **Closes:** G5.

Run complete acceptance and write a short reproducible outage/recovery walkthrough with screenshots of the actual application.

Acceptance:

- Record outcomes for all eleven browser checks, including shutdown during work, embedded assets, and volume persistence.
- A fresh setup following the README reaches the demo with exact commands; dependency/toolchain pins and instructions agree.
- Formatting, tests, vet, race verification, build, Compose validation, and final CI pass; unavailable required checks leave the gate open.
- README describes implemented behavior and links the walkthrough/screenshots; ticket/roadmap completion reflects evidence.
- Document public-demo restrictions from architecture without adding deployment, accounts, alerts, or deferred features. Stop after G5.
