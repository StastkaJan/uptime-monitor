# Executable

Create `main.go` here in milestone 1. Its job is configuration, application construction, starting the scheduler and HTTP server, and graceful shutdown.

Keep business rules and route handlers in `internal/uptime`. Prefer a small `run() error` called by `main()` so deferred cleanup runs before a fatal exit.

Read configuration once at startup. On shutdown, cancel background checks, stop accepting requests, wait for in-flight work within a deadline, and close the database last.
