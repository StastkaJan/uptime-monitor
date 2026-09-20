# Application package

Start with one package named `uptime`. Add files as the corresponding milestone needs them:

| File | Responsibility |
| --- | --- |
| `app.go` | Application construction and owned dependencies |
| `handlers.go` | Routes, form parsing, HTTP responses |
| `views.go` | Embedded assets, template parsing, page/fragment rendering |
| `monitor.go` | Monitor and check types, validation, status rules |
| `store.go` | Parameterized queries, transactions, schema initialization |
| `checker.go` | One bounded outbound HTTP check |
| `scheduler.go` | Due work, concurrency limits, cancellation |
| `demo.go` | Explicit local demo mode and controlled failure endpoint |

These are suggested ownership boundaries, not files to generate empty. Place `*_test.go` files beside the implementation. Create an embedded `schema.sql` when persistence is added; introduce versioned migrations when the schema first changes after data exists.

Use `templates/` for pages and named fragments. Start with a shared layout, monitor list, monitor form, and detail view. Use `static/` for CSS and a locally served, version-pinned HTMX distribution with its license.

Keep shared render, validation, and database code here. Do not create separate frontend and API models for the same server-rendered feature.
