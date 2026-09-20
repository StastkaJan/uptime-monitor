# Testing and acceptance

There are no runnable tests in this scaffold yet. Add behavior tests with each implementation slice. Go's [built-in test tooling](https://go.dev/doc/tutorial/add-a-test) is the default; no assertion or mocking framework is needed.

## What to test

| Area | Meaningful proof |
| --- | --- |
| Validation | Reject unsupported schemes, missing hosts, embedded credentials, oversized names, and invalid intervals; allow valid inputs |
| Checker | A controlled 200 succeeds; 503 and redirects fail; a blocked response hits its deadline; shutdown cancellation does not record an outage |
| Scheduler | No duplicate check for one monitor; no more than four concurrent checks; paused monitors are not dispatched; shutdown releases workers |
| Storage | A fresh database initializes; monitors survive reopening; check and incident updates are atomic; repeated failures create one incident; recovery closes it |
| Status | Pending, stale, paused, failure, and recovery follow the documented rules |
| Handlers | Full requests get documents, fragment requests get correct fragments, invalid forms retain values, and failed writes are reported |
| Trust boundaries | User text is escaped; SQL values remain data; cross-origin mutations are rejected; invalid targets are rejected |
| Demo | Explicit demo mode gates its routes; a controlled outage and recovery produce the expected incident |

Use `httptest.Server` for outbound checks and `httptest.ResponseRecorder` for handlers. Use a SQLite file under `t.TempDir()` for storage tests, including reopen tests. Do not depend on real websites, network availability, or a shared development database.

Synchronize concurrency tests with channels and completion signals. Drive due-time calculations with explicit timestamps rather than long sleeps. Deadline tests should use generous outer bounds and a handler that waits for cancellation. Use a second database connection or reopening where necessary to prove persistence and rollback rather than inspecting implementation internals.

## Routine commands

Validate Compose configuration now with `docker compose config --quiet`. Once Go code exists, run from the project root using the same service as development:

```powershell
docker compose run --rm --no-deps app gofmt -w ./cmd ./internal
docker compose run --rm --no-deps app go test ./...
docker compose run --rm --no-deps app go vet ./...
docker compose run --rm --no-deps app go build -o /tmp/uptime ./cmd/uptime
```

For scheduler/concurrency changes, also run:

```powershell
docker compose run --rm --no-deps app go test -race ./...
```

These one-off containers override application startup and do not publish ports, so checks can run while the dashboard is open. The temporary build output is discarded after verification. Tests must use `t.TempDir()` databases instead of the configured application database.

The race detector needs a supported architecture and C toolchain even though the chosen SQLite driver is CGo-free. The pinned Debian-based Go image supplies the toolchain; use Linux amd64 or arm64 for race verification. In CI, validate Compose and check formatting without rewriting source, then run tests, vet, race checks, and build with the same pinned Go version. Add CI once the first runnable slice exists.

Do not set an arbitrary coverage target. Cover important failure paths and state transitions; avoid tests that only mirror trivial field assignments or assert whole-page snapshots.

## Browser acceptance before the showcase is done

1. Start with an empty database. Add a monitor, edit its name and interval, and restart; settings survive.
2. Verify a new monitor moves from pending to healthy without reloading the page.
3. In local demo mode, trigger a failure; exactly one incident opens. Restore the endpoint; that incident closes.
4. Pause a monitor; no new work is scheduled after in-flight work finishes. Resume it and observe checking restart.
5. Submit an invalid form; readable errors appear and entered values remain.
6. Filter results while polling continues; filters persist and the focused input is preserved.
7. Disable JavaScript; navigation and forms still work with full-page responses.
8. Check a narrow mobile viewport, keyboard navigation, visible focus, and status labels without relying on color.
9. Stop the process during active checks; shutdown completes and the next start opens the database successfully.
10. Start the built executable from another directory using an explicit database path; embedded pages and assets still load.
11. Run `docker compose down`, then `docker compose up`; saved monitors and history survive in the named volume. Verify shutdown completes within the configured grace period.

Manual browser checks are sufficient initially. Add a small automated browser smoke test only when regressions or repeated release checks justify its tooling.
