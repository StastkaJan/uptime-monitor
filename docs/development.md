# Development directions

## Starting implementation

This is a documentation scaffold, so there is currently no Go module or executable. Go was not found on PATH when it was prepared.

1. Install a supported stable Go release and confirm `go version` works in a new terminal.
2. From this folder, initialize the module with the intended repository import path. For local-only development, `go mod init example.com/uptime-monitor` is an acceptable temporary choice.
3. Implement milestone 1 from [the roadmap](roadmap.md). Add the SQLite driver only in the storage milestone and vendor a pinned HTMX asset when interactions are introduced.
4. Commit dependency lock information and update this guide with the exact versions and any new setup steps.

After milestone 1 exists, use:

```powershell
go run ./cmd/uptime
```

After configuration and storage exist, the intended optional overrides are:

```powershell
$env:UPTIME_ADDR = '127.0.0.1:8080'
$env:UPTIME_DB = 'data/uptime.db'
$env:UPTIME_DEMO = 'true'
go run ./cmd/uptime
```

Create the data directory safely on startup if missing. Keep runtime databases and credentials out of version control. Build output goes to `bin/`; create that directory before building.

## How to split work

Work in small, end-to-end slices: one user behavior, its persistence if needed, its HTML, and its tests. The first usable slice is adding a monitor and seeing it after a restart. Finish each roadmap milestone's acceptance checks before advancing.

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
