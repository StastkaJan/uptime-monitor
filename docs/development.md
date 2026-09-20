# Development directions

## Starting implementation

The runnable shell includes the Go module and executable. Docker Compose is the default workflow; installing Go on the host is optional. Go 1.27.1 on Docker Desktop Linux amd64 has been verified; see [verification evidence](verification.md).

1. Install Docker with Compose v2 and start its engine. On Windows, use Docker Desktop in Linux-container mode. Confirm `docker compose version` and `docker info` work.
2. From this folder, validate the configuration and check the toolchain using the commands below. The service uses Go 1.27.1; keep CI aligned with that version.
3. Follow milestone 1 from [the roadmap](roadmap.md), starting with UPT-001 in [the ticket backlog](tickets.md), and verify its [integration gate](integration-stages.md). Add the SQLite driver only in the storage milestone and vendor a pinned HTMX asset when interactions are introduced.
4. Commit dependency lock information and update this guide with the exact versions and any new setup steps.

Check the prepared toolchain before starting the application. The first image download requires internet access:

```powershell
docker compose config --quiet
docker compose run --rm --no-deps app go version
```

The committed module uses only the standard library; no dependency download or `go.sum` is needed yet. One-off commands override the app startup command and do not publish its port. On native Linux, the default container user can create root-owned source files; use `docker compose run --rm --no-deps --user "$(id -u):$(id -g)" -e GOCACHE=/tmp/go-build app ...` for commands that write source files, such as module initialization or formatting.

Start the app and open `http://localhost:8080`:

```powershell
docker compose up
```

To run in the background, use `docker compose up -d`; inspect output with `docker compose logs -f app`. After changing source or embedded assets, run `docker compose restart app` to rebuild and restart. Initial builds can take longer while dependencies download.

After demo mode is implemented, enable it in PowerShell:

```powershell
$env:UPTIME_DEMO = 'true'
docker compose up -d
```

Compose recreates the service when its environment changes; `restart` alone does not apply changed environment values. Use `docker compose stop` to stop it or `docker compose down` to remove its containers and network. Both preserve named volumes. `docker compose down --volumes` also deletes the SQLite data and caches; use that only for an intentional reset.

The database lives at `/data/uptime.db` in the `uptime-data` volume, separate from source. Before backing it up, stop the application and copy the data directory from the volume; do not copy a live database without SQLite-aware backup handling. Tests use temporary databases and must never reset the application volume.

Keep runtime databases and credentials out of version control. The application should create its data directory safely if missing. A native run remains possible with Go 1.27.1 installed: `go run ./cmd/uptime`, using `data/uptime.db` by default. The Compose configuration uses the official toolchain image directly; a custom Dockerfile is deferred until an executable deployment image is needed.

## Runtime configuration

The executable reads environment variables directly; Compose can additionally
read its own `.env` file. Native defaults are `UPTIME_ADDR=127.0.0.1:8080`,
`UPTIME_DB=data/uptime.db`, and `UPTIME_DEMO=false`. Compose sets the address to
`0.0.0.0:8080` inside the container and publishes only `127.0.0.1:8080` on the host.

`UPTIME_ADDR` requires an explicit host and numeric port from 1 to 65535.
`UPTIME_DB` must be nonblank, and `UPTIME_DEMO` accepts Go boolean values such as
`true` and `false`. Invalid or explicitly empty values fail startup with the
variable name and corrective guidance. Database and demo settings are validated
now; their features are introduced in later milestones.

The HTTP server has a five-second header timeout and a 60-second idle timeout.
Ctrl+C or SIGTERM stops accepting requests and allows up to ten seconds for
active requests to finish, then closes remaining connections. This deadline is
shorter than Compose's 15-second stop grace period. `docker compose stop` keeps
the named data volume.

## How to split work

Work in small, end-to-end slices: one user behavior, its persistence if needed, its HTML, and its tests. The first usable slice is adding a monitor and seeing it after a restart. Use ticket dependencies for implementation order and finish each milestone's integration gate before advancing. Introduce CI in milestone 1 and extend it as behavior is added. Keep ticket and roadmap checkboxes aligned with verified completion.

Split source files by responsibility as described in [the application folder](../internal/uptime/README.md). Split a package only when it has a clear cohesive job and the import direction stays simple. A file's length alone is not a reason for a new abstraction.

If multiple contributors work at once, first agree on the monitor/check types, routes, and expected rendered fragments. Assign ownership of files; avoid simultaneous edits to shared schema, models, or templates. Merge foundational changes before dependent changes. Separate services and separate repositories are unnecessary.

## Reuse rules

- Search existing code and templates before writing another helper.
- Reuse the same monitor form for create/edit and the same list fragment for initial rendering/filtering/polling.
- Keep URL and interval validation in one place, called by all mutation paths.
- Keep the result-plus-incident transaction in one store operation.
- Reuse one HTTP client, one parsed template set, and the application database handle.
- Extract repeated code when shared behavior is clear. A short repeated line is cheaper than a speculative framework.
- Keep helpers close to their callers. No `utils` package or generic CRUD layer.
- Prefer real HTTP test servers and temporary databases to interfaces created only to support mocks.

Use explicit constructor parameters for dependencies; avoid globals. A concrete application struct can own the store, client, templates, and scheduler. Small function parameters for a clock or check execution are reasonable only where deterministic tests require them.

## Dependencies and changes

Use the standard library first. Each new dependency should solve a current problem and have a short rationale in the change description. Pin chosen versions; do not serve an unversioned CDN script. Preserve the HTMX license with its local distribution.

Start with plain SQL and a schema file. Once a schema change must preserve existing data, add a numbered migration, a migration test starting from the old schema, and backup/restore instructions. Do not silently recreate a user's database.

Describe each change by observable behavior and verification. Keep unrelated cleanup separate. Update completed roadmap checkboxes and distinguish implemented behavior from planned behavior.

## Deliberate limits

One process, one operator, at most 50 monitors, and one database connection define the initial scale. Add accounts when multiple people need private configurations; consider another database or worker arrangement only after measured contention or a deployment requirement demands it. Add charts, SSE, notifications, and browser-test infrastructure when a milestone actually calls for them.
