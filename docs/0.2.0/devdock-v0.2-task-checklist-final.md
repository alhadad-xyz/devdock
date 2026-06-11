# DevDock v0.2 Task Checklist

**Version:** v0.2.0
**Release Theme:** Workflow Expansion
**Base:** DevDock v0.1.0 complete
**Ticket prefix:** DD-040 onward (continues from v0.1 DD-001–039)
**Reference:** DevDock v0.2 PRD Final

Ticket format matches v0.1: each has a build description and **Done when** condition.
P2 (`devdock restart`) tickets are at the end; they do not block release.

---

## Priority Reference

| Priority | Features |
|---|---|
| P0 — Must ship | `devdock open`, `devdock run`, `devdock service add/remove/status/logs`, Mailpit, MinIO |
| P1 — Should ship | Express recipe, Go Fiber recipe, doctor updates, backward compat tests |
| P2 — Ship if capacity | `devdock restart <service>` only |

**P2 is not in the Definition of Done. Do not delay release for it.**

---

## Dependency Map

```
DD-040 (service registry)
├── DD-041 (open — reads registry for web UI targets)
├── DD-044 (service add — uses registry defaults)
├── DD-045 (service remove — uses registry for volume flag)
├── DD-046 (service status — uses registry for port/URL display)
├── DD-050 (Mailpit definition — extends registry)
└── DD-053 (MinIO definition — extends registry)

DD-042 (run) — standalone, no registry dependency
DD-044 → DD-043 (env var prompt — shared logic)
DD-050 → DD-056 (Compose generator multi-port updates)
DD-053 → DD-056

DD-057 (Express detection) → DD-058 (Express recipe)
DD-059 (Go Fiber detection) → DD-060 (Go Fiber recipe)

DD-061 (doctor updates) — depends on DD-058, DD-060 being done
DD-062 (error audit) — depends on all Week 1–4 tickets
DD-063 (backward compat tests) — standalone
```

---

# Week 1 — `devdock open` + `devdock run`

## Service Registry Abstraction

- `[ ]` **DD-040** Define `ServiceDefinition`, `PortDefinition`, `WebUIDefinition`, `HealthCheckDefinition` Go structs in `internal/services/registry.go`
- `[ ]` **DD-040** Implement `registry.Get(name string) (ServiceDefinition, bool)`
- `[ ]` **DD-040** Implement `registry.All() []ServiceDefinition`
- `[ ]` **DD-040** Register all five existing services: postgres, mysql, redis (with basic definitions); Mailpit and MinIO definitions will be added in Week 3 but stubs should be registered now
- `[ ]` **DD-040** Each service definition must include: ports, web UI (if any), whether it has a persistent volume

**Done when:** `registry.Get("postgres")` returns the correct definition. `registry.Get("mongodb")` returns `false`. All 5 services are registered. Unit test covers each registered service.

**Unblocks:** DD-041, DD-044, DD-045, DD-046, DD-050, DD-053.

---

## `devdock open` Command

- `[ ]` **DD-041** Add `open` command to Cobra CLI
- `[ ]` **DD-041** Implement target resolution:
  - No arg → read `app.port` from `.devdock.yml`
  - `app` arg → same as no arg
  - Service name → look up `WebUIDefinition` from registry; if nil (no web UI), print service-specific message: "postgres does not have a web interface. Connect via your database client at localhost:5432."
  - Unknown arg → list all openable targets for current project and exit 1
- `[ ]` **DD-041** Check if named service is enabled in `.devdock.yml`; if not, print error with `devdock service add <name>` fix
- `[ ]` **DD-041** If app process is not running (no PID or PID dead): print warning, open URL anyway
- `[ ]` **DD-041** Implement browser open: `exec.Command("open", url).Run()` (macOS)

**Done when:**
- `[ ]` **DD-041** `devdock open` in a running Next.js project opens `http://localhost:3000` in the browser
- `[ ]` **DD-041** `devdock open postgres` prints the "no web interface" message (not an error, exits 0)
- `[ ]` **DD-041** `devdock open mailpit` when mailpit is not in `.devdock.yml` prints the "not configured" error
- `[ ]` **DD-041** `devdock open unknownthing` lists openable targets and exits 1
- `[ ]` **DD-041** App not running: warning is printed, browser still opens

**Unblocks:** DD-047 (open mailpit), DD-048 (open minio) — those work automatically once their web UI definitions are registered.

---

## `devdock run` Command

- `[ ]` **DD-042** Add `run` command to Cobra CLI
- `[ ]` **DD-042** Parse `commands` block from `.devdock.yml` into `map[string]string`
- `[ ]` **DD-042** Implement command lookup by name
- `[ ]` **DD-042** Execute via shell: `exec.Command("/bin/sh", "-c", commandString)` in project root directory
  - Shell required for compound commands (`composer install && npm install`, pipes)
- `[ ]` **DD-042** Stream stdout and stderr directly (no buffering) — use `cmd.Stdout = os.Stdout`, `cmd.Stderr = os.Stderr`
- `[ ]` **DD-042** Return the underlying command's exit code via `cmd.ProcessState.ExitCode()`
- `[ ]` **DD-042** No args → list all commands in formatted table (name + command string)
- `[ ]` **DD-042** Unknown command name → print name, list available, exit 1

**Done when:**
- `[ ]` **DD-042** `devdock run migrate` on a Laravel project with `commands.migrate: php artisan migrate` runs the migration, output streams in real time, exit code matches
- `[ ]` **DD-042** `devdock run` (no args) lists all defined commands
- `[ ]` **DD-042** `devdock run deploy` when deploy is not defined shows helpful error and available commands, exits 1
- `[ ]` **DD-042** A command that exits 1 causes `devdock run` to exit 1
- `[ ]` **DD-042** A compound command (`echo hello && echo world`) works correctly via shell execution

**Unblocks:** Nothing — standalone.

---

## Recipe Command Block Updates

- `[ ]` **DD-043** Update `recipes/laravel.yml` to add full command block:
  - `install: composer install && npm install`
  - `dev: php artisan serve --host=127.0.0.1 --port=8000`
  - `migrate: php artisan migrate`
  - `seed: php artisan db:seed`
  - `test: php artisan test`
- `[ ]` **DD-043** Update `recipes/next.yml`:
  - `install: pnpm install`
  - `dev: pnpm dev`
  - `build: pnpm build`
  - `lint: pnpm lint`
  - `test: pnpm test`
- `[ ]` **DD-043** Verify `devdock init` for both stacks now includes command block in generated `.devdock.yml`

**Done when:** `devdock init` in a fresh Laravel project produces a `.devdock.yml` with all 5 Laravel commands. Same for Next.js. `devdock run` (no args) lists them correctly.

---

## Week 1 End-to-End Validation

Manual test — run before moving to Week 2:

```bash
# In an existing Laravel project (from v0.1 tests):
devdock up --detach
devdock open                 # browser opens localhost:8000
devdock run migrate          # migration runs with real-time output, exits 0
devdock run deploy           # "command not defined" error, exits 1
devdock run                  # lists all commands
devdock down
```

**Done when:** All steps above pass without errors or unexpected behavior.

---

# Week 2 — Service Management Commands

## `devdock service add` Command

- `[ ]` **DD-044** Add `service` command group to Cobra; add `add` subcommand
- `[ ]` **DD-044** Load and validate `.devdock.yml`
- `[ ]` **DD-044** Reject if `project.type: docker-compose` (use standard error message from PRD Section 11)
- `[ ]` **DD-044** Look up service in registry; reject unknown names with list of supported services
- `[ ]` **DD-044** Check if service already enabled; if so, print "already configured" message and exit 0
- `[ ]` **DD-044** Merge service registry defaults into `.devdock.yml` service block
- `[ ]` **DD-044** Write `.devdock.yml` atomically (temp file + `os.Rename`)
- `[ ]` **DD-044** Regenerate `compose.yml` via existing Compose generator
- `[ ]` **DD-044** Print env var block (formatted, clearly labeled by stack type if multi-stack project)
- `[ ]` **DD-044** Prompt: "Append these to .env? (Y/n)" — default Y
- `[ ]` **DD-044** If yes: call `env.MergeSafe()` (only adds missing keys, never overwrites)
- `[ ]` **DD-044** Detect if services are running (via PID file or `docker compose ps`); if so, add running-services notice
- `[ ]` **DD-044** Print final next step

**Done when:**
- `[ ]` **DD-044** `devdock service add redis` on a Next.js project (no Redis) updates `.devdock.yml`, regenerates `compose.yml`, prints env vars, appends to `.env` on Y
- `[ ]` **DD-044** Running the same command again: "redis is already configured" and exit 0
- `[ ]` **DD-044** `devdock service add mongodb`: "not a supported service" error with list
- `[ ]` **DD-044** `devdock service add redis` when Docker services are running: adds config + prints restart notice, does NOT touch containers
- `[ ]` **DD-044** For `docker-compose` project type: blocked with correct error

**Unblocks:** DD-050, DD-053 (Mailpit and MinIO need this working first).

---

## `devdock service remove` Command

- `[ ]` **DD-045** Add `remove` subcommand to `service` command group
- `[ ]` **DD-045** Load `.devdock.yml`; reject for `docker-compose` project type
- `[ ]` **DD-045** If service not in config: "not configured" message, exit 0
- `[ ]` **DD-045** If service has `HasVolume: true` in registry: print volume warning (non-blocking, shown before confirmation)
- `[ ]` **DD-045** Prompt: "Remove <name> from this project? (y/N)" — default N
- `[ ]` **DD-045** If confirmed: remove service block from `.devdock.yml`, write atomically
- `[ ]` **DD-045** Regenerate `compose.yml`
- `[ ]` **DD-045** Print: "Run `devdock down && devdock up` to apply."

**Done when:**
- `[ ]` **DD-045** `devdock service remove redis` prompts, shows volume warning for Redis (stateful), removes on confirm
- `[ ]` **DD-045** Answering N aborts with "No changes made." and exit 0
- `[ ]` **DD-045** Service removed: `.devdock.yml` no longer has the redis block; `compose.yml` no longer has redis service
- `[ ]` **DD-045** Running containers are not affected (verify with `docker ps` after removal)

---

## `devdock service status` Command

- `[ ]` **DD-046** Add `status` subcommand to `service` command group
- `[ ]` **DD-046** Query Docker state for all enabled services via `docker.ComposePs()`
- `[ ]` **DD-046** Look up web UI definition from registry for each service
- `[ ]` **DD-046** Print table: Service | Status | Port (primary port) | Healthy | Web URL
- `[ ]` **DD-046** For services with no web UI: show `—` in Web URL column
- `[ ]` **DD-046** "No services configured" if no services in `.devdock.yml`
- `[ ]` **DD-046** "Services are configured but not running. Run `devdock up`." if configured but stopped

**Done when:** `devdock service status` in a project with postgres + mailpit + minio shows all three with correct ports, health, and URLs. Mailpit shows SMTP port (1025) as primary, web URL in last column. MinIO shows API port (9000) as primary, console URL in last column.

---

## `devdock service logs` Command

- `[ ]` **DD-047** Add `logs` subcommand to `service` command group
- `[ ]` **DD-047** Validate service is configured in `.devdock.yml`; if not, print error with fix
- `[ ]` **DD-047** Delegate to existing `docker.ComposeLogs(projectDir, serviceName, tail, since)` implementation
- `[ ]` **DD-047** Pass through `--tail` and `--since` flags

**Done when:** `devdock service logs postgres --tail 20` streams last 20 lines. `devdock service logs mailpit` (not configured) prints the "not configured" error with fix.

---

## Week 2 End-to-End Validation

```bash
# In existing Next.js project with PostgreSQL:
devdock service status              # postgres: running, no mailpit/minio yet
devdock service add redis           # adds redis, prints env vars
devdock service status              # redis: stopped (not started yet)
devdock up                          # both postgres and redis start
devdock service status              # both healthy
devdock service remove redis        # prompts, shows volume warning, removes
# reject prompt: no changes
devdock service remove redis        # accept this time
devdock service status              # only postgres
devdock service logs postgres --tail 5
devdock down
```

---

# Week 3 — Mailpit and MinIO Services

## Mailpit Service Definition

- `[x]` **DD-048** Create `internal/services/mailpit.go`
- `[x]` **DD-048** Define Mailpit `ServiceDefinition` with:
  - Image: `axllent/mailpit:v1.21` (pinned, not latest)
  - Ports: SMTP 1025 (primary), Web 8025
  - `HasVolume: false`
  - `WebUI: &WebUIDefinition{PortName: "ui", Label: "Mailpit Inbox"}`
  - Health check: HTTP GET `http://localhost:8025/api/v1/info` (type=http)
  - Env templates for `laravel` and `express` stack types
  - Default `.devdock.yml` config: `enabled: true, version: "v1.21", smtp_port: 1025, ui_port: 8025`
- `[x]` **DD-048** Register in service registry

**Done when:** `registry.Get("mailpit")` returns correct definition. All fields populated correctly per PRD Section 7.2.

---

## MinIO Service Definition

- `[x]` **DD-049** Create `internal/services/minio.go`
- `[x]` **DD-049** Define MinIO `ServiceDefinition` with:
  - Image: `minio/minio:RELEASE.2024-11-07T00-52-20Z` (pinned)
  - Ports: API 9000 (primary), Console 9001
  - `HasVolume: true`
  - `WebUI: &WebUIDefinition{PortName: "console", Label: "MinIO Console"}`
  - Health check: exec `mc ready local` (type=exec)
  - Env templates for `laravel` and `express` stack types
  - Default `.devdock.yml` config: all fields from PRD Section 7.3
- `[x]` **DD-049** Register in service registry

**Done when:** `registry.Get("minio")` returns correct definition. `HasVolume: true` is set (triggers volume warning on remove). All env var templates match PRD Section 7.3 exactly.

---

## Compose Generator: Multi-Port and Init Container Support

- `[x]` **DD-050** Update `compose.Generate()` to handle services with multiple port mappings (one port per `- "host:container"` line)
- `[x]` **DD-050** Add init container support: if a `ServiceDefinition` has an `InitContainer` field, generate the init service in `compose.yml` with correct `depends_on: condition: service_healthy`
- `[x]` **DD-050** Implement MinIO init container (see PRD Section 7.3 for exact Compose YAML)
- `[x]` **DD-050** Add unit tests: Mailpit compose output has 2 port mappings; MinIO compose output has init container and volume

**Done when:** Generated `compose.yml` for Mailpit matches the spec output exactly (2 port bindings, health check using wget). Generated `compose.yml` for MinIO matches spec output exactly (2 ports, volume, command, init container).

---

## Mailpit Integration Test

Manual test:

```bash
cd existing-next-postgres-project   # from v0.1 tests
devdock service add mailpit         # updates config, prints env vars
devdock up                          # postgres and mailpit start
devdock service status              # both healthy, mailpit shows http://localhost:8025
devdock open mailpit                # Mailpit inbox opens in browser
devdock service logs mailpit        # streams logs
devdock down
```

**Done when:** All steps pass without errors. Mailpit inbox loads correctly in browser.

---

## MinIO Integration Test

Manual test:

```bash
cd existing-next-postgres-project
devdock service add minio           # updates config, prints env vars
devdock up                          # postgres and minio start; init bucket created
devdock service status              # minio healthy, shows http://localhost:9001
devdock open minio                  # MinIO console opens
# In console: upload a test file to the 'local' bucket
devdock down
devdock up                          # services restart
devdock open minio                  # test file still present in 'local' bucket
devdock down --volumes
devdock up
# Verify: 'local' bucket is empty (data was deleted)
```

**Done when:** All steps pass. Data persistence across down/up confirmed. Data deletion after `--volumes` confirmed. Bucket `local` exists immediately after `devdock up` (without manual creation).

---

## Week 3 End-to-End Validation

```bash
devdock service add mailpit
devdock service add minio
devdock up
devdock service status
devdock open mailpit
devdock open minio
devdock service logs mailpit --tail 10
devdock service logs minio --tail 10
devdock service remove mailpit      # prompts; mailpit has no volume warning
devdock service remove minio        # prompts; minio shows volume warning
devdock down
```

---

# Week 4 — Express and Go Fiber Support

## Express Detection

- `[x]` **DD-053** Add `express` to project type enum
- `[x]` **DD-053** Detection rule: read and JSON-parse `package.json`; check `dependencies` field for key `"express"` (high confidence)
- `[x]` **DD-053** Also check `devDependencies` if not in `dependencies` (medium confidence — show warning: "express is unusual as a devDependency")
- `[x]` **DD-053** If `package.json` exists but no express in either field: low confidence, ask user
- `[x]` **DD-053** Next.js detection (higher priority) must still take precedence — Express detector only fires if Next.js detector did not match
- `[x]` **DD-053** Add `testdata/fixtures/express/` fixture: `package.json` with `"dependencies": {"express": "^4.18.0"}`
- `[x]` **DD-053** Add `testdata/fixtures/express-ambiguous/` fixture: `package.json` with no express

**Done when:** Detector tests pass for all Express fixtures. `detect("fixtures/express")` returns `express`, high confidence. `detect("fixtures/express-ambiguous")` returns low confidence. `detect("fixtures/nextjs")` still returns `nextjs`, not `express`.

---

## Express Recipe and Command Fallback

- `[x]` **DD-054** Create `recipes/express.yml` matching PRD Section 8.1
- `[x]` **DD-054** Implement app command detection in `devdock init` for Express:
  1. Check `package.json` `scripts.dev` → use `npm run dev`
  2. Check `package.json` `scripts.start` → use `npm start`
  3. Fallback → `node index.js` + warn user
- `[x]` **DD-054** Generate correct env vars for PostgreSQL and Redis based on project stack
- `[x]` **DD-054** `devdock doctor` for Express: validate Node.js runtime version

**Done when:** `devdock init` in a real Express project (e.g., a simple hello-world Express app with `npm run dev` script) produces a correct `.devdock.yml` with the right command, port 3000, and suggested services. `devdock up` runs the app successfully.

---

## Express End-to-End Test

Manual test (use a real Express project, not a fixture):

```bash
git clone https://github.com/expressjs/express express-test
cd express-test
devdock init                    # detects Express, suggests postgres + redis
devdock up                      # starts services and app
curl localhost:3000             # HTTP 200
devdock status
devdock open
devdock down
```

**Done when:** All steps pass on a clean macOS machine.

---

## Go Fiber Detection

- `[x]` **DD-056** Add `go-fiber` to project type enum
- `[x]` **DD-056** Detection: if `go.mod` exists, read file as text, search for a line containing `github.com/gofiber/fiber` within the `require (...)` block (not just anywhere in the file)
- `[x]` **DD-056** High confidence if found; low confidence if `go.mod` exists but no Fiber import
- `[x]` **DD-056** Add `testdata/fixtures/gofiber/` fixture: `go.mod` with `require (github.com/gofiber/fiber/v2 v2.52.5)`
- `[x]` **DD-056** Add `testdata/fixtures/gofiber-generic/` fixture: `go.mod` with only stdlib

**Done when:** Detector tests pass for all Go fixtures. Fiber fixture returns `go-fiber`, high confidence. Generic Go fixture returns low confidence and prompts user.

---

## Go Fiber Recipe

- `[x]` **DD-057** Create `recipes/go-fiber.yml` matching PRD Section 8.2
- `[x]` **DD-057** `devdock init` in Go Fiber project: check if `main.go` exists in project root; if not, print the `cmd/` subdirectory warning
- `[x]` **DD-057** `devdock doctor` for Go Fiber: validate `go version` matches `runtime.go` from config; if Go not installed, print https://go.dev/dl/ link

**Done when:** `devdock init` in a Go Fiber project with `main.go` in root produces a correct `.devdock.yml`. Running `devdock init` in a project with `cmd/server/main.go` (no root main.go) shows the warning and still generates the config.

---

## Go Fiber End-to-End Test

Manual test (use a real Go Fiber project):

```bash
# Create or clone a minimal Go Fiber hello world
mkdir gofiber-test && cd gofiber-test
go mod init example.com/myapp
go get github.com/gofiber/fiber/v2
# Add main.go with fiber hello world
devdock init                    # detects Go Fiber, suggests postgres
devdock up                      # starts postgres and go run .
curl localhost:8080             # HTTP 200
devdock run test               # runs go test ./...
devdock down
```

**Done when:** All steps pass on a clean macOS machine with Go installed.

---

# Week 5 — Polish, Doctor, Backward Compatibility

## Doctor Updates for v0.2

- `[x]` **DD-059** Add Express doctor check: if `project.type: express`, validate `node --version` against `runtime.node`
- `[x]` **DD-059** Add Go Fiber doctor check: if `project.type: go-fiber`, validate `go version` against `runtime.go`; if Go not found: link to https://go.dev/dl/
- `[x]` **DD-059** Add Mailpit port checks: if mailpit enabled, check ports 1025 and 8025 are available (before `devdock up`)
- `[x]` **DD-059** Add MinIO port checks: if minio enabled, check ports 9000 and 9001 are available

**Done when:** `devdock doctor` in a Go Fiber project with Go not installed: Go check fails with install link. In a project with mailpit, port 8025 already in use: doctor reports the conflict with correct port number and service name.

---

## Error Message Audit

- `[x]` **DD-060** Go through every `return error` and `ui.ErrorWithFix` added in Weeks 1–4
- `[x]` **DD-060** For each: verify it has: what happened, why it happened, what to do
- `[x]` **DD-060** Create a checklist of all v0.2 error paths and mark each as pass/fail
- `[x]` **DD-060** Fix all that fail

Target error paths to audit:
- `[x]` **DD-060** `devdock open <unknown>`
- `[x]` **DD-060** `devdock open <service-not-configured>`
- `[x]` **DD-060** `devdock open <service-no-webui>`
- `[x]` **DD-060** `devdock run <undefined>`
- `[x]` **DD-060** `devdock service add <unknown-service>`
- `[x]` **DD-060** `devdock service add` on docker-compose project
- `[x]` **DD-060** `devdock service remove <not-configured>`
- `[x]` **DD-060** `devdock service logs <not-configured>`
- `[x]` **DD-060** `devdock service add` when running
- `[x]` **DD-060** Express command fallback warning
- `[x]` **DD-060** Go Fiber main.go warning

**Done when:** All error paths in the checklist above pass the what/why/fix format review.

---

## Backward Compatibility Fixture Tests

- `[x]` **DD-061** Create `testdata/backward-compat/` directory with 3 fixture configs:
  - `v0.1-laravel.yml` — a v0.1 Laravel `.devdock.yml` (no commands block, no Mailpit, no MinIO)
  - `v0.1-nextjs.yml` — a v0.1 Next.js `.devdock.yml`
  - `v0.1-docker-compose.yml` — a v0.1 Docker Compose `.devdock.yml`
- `[x]` **DD-061** Add automated tests that load each fixture via `config.Load()` and assert:
  - No error is returned
  - No schema migration prompt is triggered
  - Unknown fields produce a warning, not an error
  - All v0.1 commands (`devdock up`, `devdock down`, `devdock status`) resolve correctly against the loaded config

**Done when:** All 3 × 4 = 12 assertions pass. These tests run in `go test ./...` with no manual intervention.

---

## `devdock restart <service>` (P2 — If Capacity Allows)

Only start this ticket after DD-059, DD-060, and DD-061 are complete.

- `[x]` **DD-062** Add `restart` command to Cobra CLI (service only — no app restart)
- `[x]` **DD-062** `devdock restart postgres` → `docker compose restart postgres` in project directory
- `[x]` **DD-062** `devdock restart app` → print: "Restarting the app is not supported. Use Ctrl+C then `devdock up`." and exit 0
- `[x]` **DD-062** `devdock restart` (no arg) → print: "Specify a service to restart. Example: `devdock restart postgres`" and exit 1
- `[x]` **DD-062** Unknown service name → list configured services

**Done when:** `devdock restart postgres` in a running project restarts the container and it returns to healthy. `devdock restart app` prints the helpful message without error.

---

# Week 6 — Integration Testing and Release

## Full End-to-End Flow Tests

Run all 9 flows. Document result (pass/fail + notes) for each. Do not proceed to release until all pass.

| Flow | Command Sequence |
|---|---|
| 1 — open and run | `devdock up --detach && devdock open && devdock run migrate` |
| 2 — add mailpit | `devdock service add mailpit && devdock up && devdock open mailpit` |
| 3 — add minio | `devdock service add minio && devdock up && devdock open minio && devdock down && devdock up` (verify persistence) |
| 4 — remove service | `devdock service add redis && devdock service remove redis` (confirm prompt, verify config clean) |
| 5 — Express | `devdock init && devdock up && curl localhost:3000` (in Express project) |
| 6 — Go Fiber | `devdock init && devdock up && curl localhost:8080 && devdock run test` (in Go Fiber project) |
| 7 — v0.1 Laravel backward compat | v0.1 Laravel project: `devdock up && devdock status && devdock down` |
| 8 — v0.1 Next.js backward compat | v0.1 Next.js project: `devdock up && devdock status && devdock down` |
| 9 — docker-compose protection | `devdock service add postgres` in docker-compose project → correct error |

Run all flows on both Apple Silicon and Intel.

---

## README and Documentation Update
## v0.2.0 Release

- `[ ]` **DD-065** Bump version to `0.2.0` in `cmd/devdock/main.go`
- `[ ]` **DD-065** Build release binaries for `darwin/arm64` and `darwin/amd64`
- `[ ]` **DD-065** Run `devdock --version` on each binary to confirm version string
- `[ ]` **DD-065** Test `install.sh` on a clean macOS machine (no prior DevDock install)
- `[ ]` **DD-065** Create GitHub Release with both binaries and updated `install.sh`
- `[ ]` **DD-065** Write release notes (from PRD Section 15)
- `[ ]` **DD-065** Ask v0.1 beta testers to validate v0.2 on their machines

**Done when:** GitHub Release is published. At least 2 external testers confirm flows 1 and 2 work on their machines.

---

## v0.2 Definition of Done Sign-Off

**Must all be checked before tagging v0.2.0:**

### Workflow Commands
- `[ ]` `devdock open` opens app URL (Flow 1)
- `[ ]` `devdock open mailpit` opens Mailpit UI (Flow 2)
- `[ ]` `devdock open minio` opens MinIO console (Flow 3)
- `[ ]` `devdock open postgres` prints "no web interface" message, exits 0
- `[ ]` `devdock open <unknown>` lists targets, exits 1
- `[ ]` `devdock run <command>` runs with real-time output, correct exit code (Flow 1)
- `[ ]` `devdock run` lists available commands
- `[ ]` `devdock run <undefined>` shows available commands, exits 1

### Service Management
- `[ ]` `devdock service add mailpit` updates config, prints env vars (Flow 2)
- `[ ]` `devdock service add minio` updates config, prints env vars (Flow 3)
- `[ ]` `devdock service add <already-configured>` is a no-op
- `[ ]` `devdock service add` while running: config updated, notice printed, containers untouched
- `[ ]` `devdock service remove` prompts before removing (Flow 4)
- `[ ]` `devdock service remove <stateful-service>` shows volume warning before prompt
- `[ ]` `devdock service add/remove` blocked for docker-compose projects (Flow 9)
- `[ ]` `devdock service status` shows correct dual-port output for Mailpit and MinIO
- `[ ]` `devdock service logs <name>` works with `--tail` and `--since`

### New Services
- `[ ]` Mailpit starts and health check passes (Flow 2)
- `[ ]` Mailpit SMTP accessible on 1025
- `[ ]` Mailpit web UI accessible at http://localhost:8025
- `[ ]` MinIO starts and health check passes (Flow 3)
- `[ ]` MinIO `local` bucket exists immediately after `devdock up` without manual creation
- `[ ]` MinIO console accessible at http://localhost:9001
- `[ ]` MinIO data persists across `devdock down` + `devdock up`
- `[ ]` MinIO data deleted by `devdock down --volumes`

### New Stacks
- `[ ]` Express detected from `package.json` dependencies (Flow 5)
- `[ ]` Express command fallback: `dev` → `start` → `node index.js`
- `[ ]` Express `devdock init && devdock up && curl localhost:3000` passes on clean machine
- `[ ]` Go Fiber detected from `go.mod` Fiber require (Flow 6)
- `[ ]` Go Fiber main.go warning shown when no root main.go
- `[ ]` Go Fiber `devdock init && devdock up && curl localhost:8080` passes on clean machine
- `[ ]` Go runtime check in `devdock doctor`

### Backward Compatibility
- `[ ]` Automated backward compat fixture tests all pass (DD-061)
- `[ ]` v0.1 Laravel project: `devdock up && devdock down` passes (Flow 7)
- `[ ]` v0.1 Next.js project: `devdock up && devdock down` passes (Flow 8)
- `[ ]` v0.1 configs load without schema migration prompt

### Quality
- `[ ]` All v0.2 errors pass the what/why/fix audit (DD-060)
- `[ ]` `install.sh` tested on clean macOS machine
- `[ ]` Both binaries tested on Apple Silicon and Intel
- `[ ]` README updated (DD-064)
- `[ ]` Version string is `0.2.0`
