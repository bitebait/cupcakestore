# Repository Guidelines

## Project Structure and Architecture

Cupcake Store is a Go/Fiber application using GORM and server-rendered HTML. `main.go` starts the application through `bootstrap/`. Dependencies are composed once in `bootstrap/routes.go`; `routers/` only registers endpoints with injected controllers. Keep HTTP concerns in `controllers/` and `middlewares/`, business operations in `services/`, persistence in `repositories/`, and domain rules in `models/`.

`config/`, `database/`, and `session/` manage infrastructure; `helpers/` and `messages/` provide shared utilities. Templates live in `views/`, with layouts and reusable snippets. Application CSS and JavaScript live in `web/css/` and `web/js/`. The remaining `web/dist/` and `web/plugins/` files are required third-party CSS, fonts, images and licenses; do not restore unused vendor bundles.

Current documentation is indexed in `docs/README.md`. Historical academic documents and screenshots live in `docs/archive/`; current interface screenshots live in `docs/screenshots/`.

## Build, Test and Development

Use Go 1.26 or later, matching `go.mod`, with current security patches. SQLite requires CGO and a C compiler. Run commands from the repository root:

- `go mod download`: download the pinned dependencies.
- `go run .`: start local development at `http://localhost:8080` by default.
- `go build ./...`: compile all packages.
- `go test ./...`: run tests and default vet checks.
- `go test -race ./...`: verify concurrent behavior.
- `go vet ./...`: run static analysis.
- `go mod tidy`: reconcile module files after dependency changes; review the diff.
- `node --check web/js/ui.js` (and the other files in `web/js/`): check browser JavaScript syntax.
- `docker compose up --build -d --wait`: start with the official Compose plugin, persistent SQLite and image volumes.

The `.env` file is optional; existing environment variables take precedence. Copy `.env.example` only if a local `.env` does not already exist. No default administrative credentials exist. Templates and static assets are loaded from disk; no frontend build is required. Read `docs/DEPLOYMENT.md` for HTTPS, backups and importing an existing database. Startup still migrates and seeds the configured database, so use isolated databases for tests and reviews.

## Coding and Interface Style

Format changed Go files with `gofmt`. Use lowercase package names, exported `PascalCase` identifiers and unexported `camelCase` identifiers. Follow existing filenames such as `product_controller.go`, `product_service.go` and `product_repository.go`. Prefer explicit composition and focused functions over new frameworks or abstraction layers.

Preserve Portuguese user-facing text. Match the four-space indentation in HTML, CSS and JavaScript. Functional text uses at least `1rem` (16 px by default); secondary text uses at least `.875rem` (14 px). Adjust layout instead of shrinking text on mobile. Use shared snippets and native forms; every state-changing form includes the session-bound CSRF token. Scripts and styles remain external to comply with CSP. Render prices with `money`, while keeping numeric input values machine-readable. See `docs/FRONTEND.md`.

## Testing and Data Integrity

Tests use Go's `testing` package beside the relevant code and isolated databases or local HTTP servers. They do not require a `.env` file. Run checks appropriate to the change; add meaningful regressions for authorization, business invariants and concurrent writes. For UI changes, verify the affected storefront/admin flows at desktop and mobile sizes.

Preserve stock and order invariants: checkout and cancellation are transactional, finalized carts cannot be edited, and status-only order updates do not reprice delivery or overwrite payment terms. Pix instructions are generated locally and confirmation is manual; local cancellation does not revoke a saved QR or issue a refund. Use the explicit payment-confirmation operation to record the responsible administrator and time. Do not use real customer databases or payment services in automated tests.

## Repository Hygiene, Commits and Pull Requests

Keep secrets, databases, uploads, editor settings and generated artifacts out of Git. Preserve historical documents, vendor attribution and licenses. Follow `.editorconfig` and existing conventions; avoid unrelated formatting churn.

History uses `feat:`, `fix:`, `refactor:`, `test:`, `docs:` and `chore:`, optionally scoped. Write concise, focused subjects. PRs should explain behavior changes, relevant validation and known limitations; include screenshots for visual changes. CI validates Go, JavaScript syntax, vulnerabilities and Docker Compose startup.
