# Implementation verification

## UPT-001: Embedded shell

Implemented on branch `implement/m1-runnable-shell` with Go 1.27.1 through
Docker Desktop's Linux amd64 engine and the existing Compose service.

- `docker compose config --quiet`: passed.
- Compose `gofmt -w ./cmd ./internal`, `go test ./...`, `go vet ./...`, and
  `go build -o /tmp/uptime ./cmd/uptime`: passed.
- `docker compose up -d`: dashboard, `/healthz`, and `/static/style.css`
  returned 200 with the expected HTML, text, and CSS content types.
- The compiled binary served its embedded dashboard and stylesheet while
  running from `/tmp`, outside the source directory.
- Temporary processes and the Compose application were stopped after checks;
  named volumes were preserved.

The ticket's Astra integrated a delegated handler-test implementation and
reviewed the result. The coordinating agent reviewed the module, routes,
embedding, templates, and tests before committing. There are no external Go
dependencies, so no `go.sum` is needed yet.

## UPT-002: Configuration and lifecycle

- Configuration tests cover native defaults, Compose-compatible overrides,
  IPv6, and actionable failures for empty or invalid environment values.
- Lifecycle tests prove active requests drain, overdue requests are canceled
  by forced connection closure, and listener failures are returned.
- Compose formatting, tests, vet, race checks, and build all passed on the
  combined UPT-001/UPT-002 implementation.
- The real Compose application returned 200 from `/healthz`; SIGTERM through
  `docker compose stop` finished in 0.836 seconds, below the 15-second grace
  period. Inspection showed exit code 0, running false, and PID 0; no project
  containers remained running.
- Verification containers shadowed `/data` with a disposable anonymous volume.

The ticket's Astra integrated delegated configuration tests and performed the
runtime checks. The coordinating agent reviewed configuration, signal wiring,
the ten-second shutdown deadline, and the lifecycle tests before committing.

## UPT-003 and combined G1 verification

The workflow runs Compose validation, nonmutating formatting verification,
tests, vet, race checks, and build using the Compose Go 1.27.1 pin on Linux.
Its Astra integrated a delegated workflow review with no findings. Live mount
inspection confirmed that `/data` uses an anonymous volume during checks.

All workflow-equivalent commands passed against the combined source. The
coordinating agent independently verified the final application:

- Host requests for health, dashboard, and stylesheet passed.
- Real Compose shutdown took 1.134 seconds and left exit code 0, running false,
  and PID 0. No application container remained running.
- The binary served embedded HTML/CSS from `/tmp`, used the native loopback
  default when `UPTIME_ADDR` was unset, and exited cleanly on SIGTERM.
- Invalid `UPTIME_ADDR` produced an actionable startup error and nonzero exit.
- An initial ten-second readiness window expired during the source build;
  retrying after the build completed passed. Allow for compilation before
  testing HTTP readiness.

The [first hosted run](https://github.com/StastkaJan/uptime-monitor/actions/runs/35524661259)
passed Compose validation, formatting, tests, vet, and race checks, but the
build failed obtaining VCS status (exit 128). CI now disables VCS stamping for
its disposable binary with `-buildvcs=false`, avoiding checkout ownership
differences between the hosted runner and container.

The [hosted rerun](https://github.com/StastkaJan/uptime-monitor/actions/runs/35524816647)
passed every step on commit `62b95da010c17e86c3bfd8660badce7d22fd5177`:
Compose validation, formatting, tests, vet, race checks, and build. The ticket's
Astra verified the run and the coordinating agent independently confirmed all
step conclusions. The subsequent gate-recording amendment changes documentation
only.

**G1 passed; UPT-001 through UPT-003 are complete.** Later milestones have not
been implemented.
