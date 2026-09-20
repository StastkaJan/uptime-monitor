# Implementation directions

This folder contains the Uptime Monitor planning scaffold. Read README.md and docs/roadmap.md before implementing. Follow the requested milestone; do not implement the entire backlog unless asked.

- Use the Go standard library first. Default dependencies are HTMX and the SQLite driver only.
- Keep one executable and one internal application package initially.
- Split files by responsibility; split packages only when a concrete boundary makes the code easier to understand.
- Keep SQL in store files, HTTP handling in handler files, and outbound checks in checker files.
- Use the same templates and validation for ordinary browser requests and HTMX requests.
- Keep business state on the server. Add custom browser JavaScript only for a demonstrated interaction that needs it.
- Prefer concrete types. Introduce small interfaces at their consumer only when actual substitution is needed.
- Search existing code before creating helpers. Avoid generic repositories, service containers, base handlers, and miscellaneous utils packages.
- Bound background concurrency and HTTP timeouts. Propagate cancellation and close response bodies.
- Use parameterized SQL, escaped templates, request size limits, and same-origin protection for mutations.
- Bind to loopback by default. Follow the public-demo constraints in docs/architecture.md before exposing the service.
- Write behavior tests beside the code. Use controlled HTTP servers and real temporary databases rather than live websites or database mocks.
- Run formatting, tests, vet, and build after meaningful code changes; run the race detector for concurrent behavior where supported.
- Update the roadmap and setup instructions with each completed slice. Report exactly what was verified and any unavailable checks.
- Add no deployment, external notifications, or unrelated features as part of a local implementation task.

Detailed acceptance criteria and commands are in docs/testing.md.
