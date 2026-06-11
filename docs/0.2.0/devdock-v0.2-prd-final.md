# DevDock v0.2 Product Requirements Document

**Product:** DevDock
**Version:** v0.2.0
**Release Theme:** Workflow Expansion
**Type:** Product Requirements Document
**Status:** Approved for Engineering Planning
**Base Version:** DevDock v0.1.0 complete
**Target Platform:** macOS
**Last Updated:** June 2026

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Product Goal](#2-product-goal)
3. [Scope](#3-scope)
4. [Problems to Solve](#4-problems-to-solve)
5. [User Stories & Acceptance Criteria](#5-user-stories--acceptance-criteria)
6. [Functional Requirements](#6-functional-requirements)
7. [New Services Specification](#7-new-services-specification)
8. [New Stack Support](#8-new-stack-support)
9. [Service Registry Interface](#9-service-registry-interface)
10. [`.devdock.yml` Schema Changes](#10-devdockyml-schema-changes)
11. [Docker Compose Project Rules](#11-docker-compose-project-rules)
12. [Non-Functional Requirements](#12-non-functional-requirements)
13. [Backward Compatibility Contract](#13-backward-compatibility-contract)
14. [6-Week Build Plan](#14-6-week-build-plan)
15. [Definition of Done](#15-definition-of-done)
16. [Risks & Mitigations](#16-risks--mitigations)

---

## 1. Executive Summary

DevDock v0.1 proved the core local development lifecycle — from zero to running project in under 5 minutes. Every feature in v0.1 served one goal: getting the project started.

v0.2 serves a different goal: making DevDock useful **after** the project is already running, during the hours of actual development that follow.

v0.2 introduces:

- `devdock open` — open any project URL without remembering ports
- `devdock run <command>` — consistent command wrapper across all stacks
- `devdock service add/remove` — modify the service list without touching Compose files
- Mailpit — local email inbox for testing
- MinIO — local S3-compatible storage for file upload testing
- Express and Go Fiber stack support — proving DevDock is multi-stack

The demo workflow that motivates v0.2:

```bash
devdock up --detach
devdock open
devdock run migrate
devdock service add mailpit
devdock up
devdock open mailpit
devdock logs mailpit
devdock service remove redis
devdock down
```

---

## 2. Product Goal

> Make DevDock useful during everyday development, not only during first-time setup.

v0.1 answered: "How do I start this project?"
v0.2 answers: "How do I work inside this project every day?"

---

## 3. Scope

### 3.1 In Scope

| Category | v0.2 Contents |
|---|---|
| Platform | macOS only (no change from v0.1) |
| Existing stacks | Laravel, Next.js, existing Docker Compose |
| New stacks | Express, Go Fiber |
| Existing services | PostgreSQL, MySQL, Redis |
| New services | Mailpit, MinIO |
| New commands | `open`, `run`, `service add`, `service remove`, `service status`, `service logs` |
| Deferred command | `restart` — P2, see Section 3.3 |
| Config | Extend `.devdock.yml` v1 schema, fully backward compatible |
| App run mode | Host app + container services only (no change) |
| Domains | `localhost:PORT` only (no change) |
| Distribution | GitHub Releases (no change) |

### 3.2 Explicitly Out of Scope

| Feature | Deferred To |
|---|---|
| Template registry / remote fetching | v0.3 |
| Self-update (`devdock update`) | v0.3 |
| Telemetry | v0.3 |
| Homebrew tap | v0.3 |
| `.test` domain support | v1.1 |
| SQLite project state database | v1.0 |
| Desktop GUI | v2.0 |
| Windows/Linux support | v2.0 |
| Plugin system | Not planned |

### 3.3 `devdock restart` Priority

`devdock restart` is **P2 — include only if Week 5 has capacity after all P0/P1 work is done and tested.**

It does not appear in the Definition of Done. Shipping v0.2 without `restart` is an acceptable outcome. Do not hold the release for it.

If implemented, the scope is: `devdock restart <service>` only (restart a Docker service by name). `devdock restart app` is explicitly excluded — restarting a foreground app process from a separate terminal session is ambiguous and can be done with Ctrl+C + `devdock up`.

---

## 4. Problems to Solve

### 4.1 Opening Local URLs Is Manual and Error-Prone

After `devdock up`, the developer must remember which port their app, mail UI, and storage console are on. With multiple projects open, this causes constant confusion.

**Solution:** `devdock open`, `devdock open mailpit`, `devdock open minio` — reads port from `.devdock.yml` and opens the right URL with one command.

### 4.2 Project Commands Are Inconsistent Across Stacks

```bash
# On Laravel:
php artisan migrate
php artisan test
composer install

# On Next.js:
pnpm prisma migrate dev
pnpm test
pnpm install

# On Go:
go test ./...
go build -o bin/app .
```

Every stack has different commands. When switching between projects, developers constantly forget the exact invocation.

**Solution:** `devdock run migrate`, `devdock run test`, `devdock run build` — consistent interface, stack-specific implementation.

### 4.3 Services Change After Project Initialization

A project starts with PostgreSQL. Three weeks later it needs email testing. Today that means: Google "docker mailpit", copy a Compose snippet, merge it manually, figure out env vars, remember what port Mailpit's SMTP is on.

**Solution:** `devdock service add mailpit` handles all of this in one command.

### 4.4 DevDock Supports Only Two Stacks

v0.1 supports Laravel and Next.js. Adding only these two risks the perception that DevDock is a PHP/JS tool.

**Solution:** Add Express and Go Fiber to demonstrate the multi-stack promise from the Vision PRD.

---

## 5. User Stories & Acceptance Criteria

### US-01: Open App in Browser

**As a** developer,
**I want to** open my local app URL from DevDock,
**So that** I don't need to remember which port is in use.

**Acceptance Criteria:**

- `devdock open` reads `app.port` from `.devdock.yml` and runs `open http://localhost:<port>` (macOS)
- `devdock open app` is identical to `devdock open`
- If `.devdock.yml` is missing: error with fix pointing to `devdock init`
- If the app process is not running (no PID file or PID not alive): print a warning, then open the URL anyway — the developer may want to check why it's not running
- Exit code is always 0 after the browser open call (the browser handles connectivity)

---

### US-02: Open Service Web UIs

**As a** developer,
**I want to** open the Mailpit inbox or MinIO console with one command,
**So that** I can inspect emails and uploaded files without looking up the port.

**Acceptance Criteria:**

- `devdock open mailpit` opens `http://localhost:<mailpit.ui_port>` (default 8025)
- `devdock open minio` opens `http://localhost:<minio.console_port>` (default 9001)
- If the named service is not configured in `.devdock.yml`: error with fix pointing to `devdock service add <name>`
- If the named service has no web UI (e.g., `devdock open postgres`): print a specific message — "postgres does not have a web interface. Connect via your database client at localhost:5432."
- `devdock open <unknown>`: lists all openable targets for the current project

---

### US-03: Run Named Project Commands

**As a** developer,
**I want to** run project-specific commands through a consistent interface,
**So that** I don't need to remember the exact invocation for each stack.

**Acceptance Criteria:**

- `devdock run migrate` runs the command string defined at `commands.migrate` in `.devdock.yml`
- Command runs in the project root directory, via shell (`/bin/sh -c "<command>"`) to support compound commands
- stdout and stderr stream in real time, not buffered
- The exit code of `devdock run <cmd>` matches the exit code of the underlying command exactly
- `devdock run` with no arguments lists all defined commands from `.devdock.yml`
- `devdock run <undefined>` prints the undefined command name, lists available commands, and exits 1

---

### US-04: Add a Service After Initialization

**As a** developer,
**I want to** add a new service to an existing project without editing Compose files,
**So that** I can add capabilities as the project grows.

**Acceptance Criteria:**

- `devdock service add mailpit` adds Mailpit to `.devdock.yml` and regenerates `compose.yml`
- Relevant environment variables are printed to the terminal in a clearly labeled block
- User is asked whether to append those vars to `.env`; if yes, only missing keys are appended (existing keys are never overwritten)
- If the service is already enabled in `.devdock.yml`: print "mailpit is already configured. Run `devdock up` to start it." and exit 0
- If the project type is `docker-compose`: print the Docker Compose project protection error and exit 1
- If services are currently running: print "Run `devdock up` to start the new service. Existing running services will not be affected." and exit 0

**Service add does not start services.** It only modifies config. The developer runs `devdock up` separately.

---

### US-05: Remove a Service

**As a** developer,
**I want to** remove a service I no longer need, without accidentally destroying data,
**So that** I can keep the project config clean.

**Acceptance Criteria:**

- `devdock service remove redis` prompts: "Remove redis from this project? (y/N)" — default is N
- After confirmation: removes `redis` from `.devdock.yml`, regenerates `compose.yml`
- If the service has a named Docker volume: prints a separate warning before the first confirmation:
  ```
  ⚠ redis has persistent data in volume: devdock-my-app-redis-data
  This data will not be deleted unless you run `devdock down --volumes` separately.
  ```
- Volumes are never automatically deleted by `devdock service remove`
- If the service is not configured: print "redis is not configured in this project." and exit 0
- Service remove does not stop a currently-running container — developer must `devdock down` and `devdock up` to apply

---

### US-06: Use Mailpit for Email Testing

**As a** developer,
**I want to** catch all outgoing emails from my local app in a web inbox,
**So that** I can test email flows without sending to real addresses.

**Acceptance Criteria:**

- `devdock service add mailpit` adds Mailpit with correct defaults
- `devdock up` starts Mailpit; health check passes before app starts
- `devdock open mailpit` opens `http://localhost:8025`
- `devdock logs mailpit` streams Mailpit container logs
- `devdock status` shows Mailpit with SMTP port and web URL

---

### US-07: Use MinIO for Local Object Storage

**As a** developer,
**I want to** test S3-compatible file uploads against a local MinIO instance,
**So that** I can develop file storage features without using a real AWS bucket.

**Acceptance Criteria:**

- `devdock service add minio` adds MinIO with correct defaults
- `devdock up` starts MinIO; health check passes before app starts
- The `local` bucket is created automatically on first start (via MinIO init container or `mc` post-start script)
- `devdock open minio` opens the MinIO console at `http://localhost:9001`
- MinIO data in the named volume persists across `devdock down` (but not across `devdock down --volumes`)
- `devdock status` shows MinIO with API port and console URL

---

### US-08: Initialize and Run an Express Project

**As a** Node.js developer,
**I want** DevDock to detect and run Express projects,
**So that** I can use DevDock on my existing Express APIs.

**Acceptance Criteria:**

- `devdock init` in an Express project produces a working `.devdock.yml`
- Detection checks `package.json` `dependencies` field for `express` key (high confidence)
- Low-confidence detection (package.json with no express dependency, but common server patterns) asks for confirmation
- Recipe detects `dev` script in `package.json`; falls back to `start`; falls back to `node index.js`
- `devdock up`, `open`, `status`, `logs`, `doctor`, and `run` all work for Express projects

---

### US-09: Initialize and Run a Go Fiber Project

**As a** Go developer,
**I want** DevDock to detect and run Go Fiber projects,
**So that** I can use DevDock on my Go APIs.

**Acceptance Criteria:**

- `devdock init` in a Go Fiber project produces a working `.devdock.yml`
- Detection parses `go.mod` and looks for a `require` block entry containing `github.com/gofiber/fiber` (high confidence)
- `go.mod` present but no Fiber import: low confidence, asks for confirmation
- `devdock doctor` validates that `go` runtime is installed and matches the version in config
- `devdock up` runs `go run .` (assumes `package main` in project root — documented assumption)
- `devdock run test` runs `go test ./...`

---

## 6. Functional Requirements

### 6.1 `devdock open [target]`

#### Command Table

| Command | Action |
|---|---|
| `devdock open` | Opens `http://localhost:<app.port>` |
| `devdock open app` | Same as `devdock open` |
| `devdock open mailpit` | Opens `http://localhost:<mailpit.ui_port>` |
| `devdock open minio` | Opens `http://localhost:<minio.console_port>` |
| `devdock open postgres` | "postgres has no web interface" message (exit 0) |
| `devdock open mysql` | "mysql has no web interface" message (exit 0) |
| `devdock open redis` | "redis has no web interface" message (exit 0) |
| `devdock open <unknown>` | Lists openable targets for this project (exit 1) |

#### Implementation

Uses macOS `open` command: `open http://localhost:<port>`.

The open target registry is defined as a map in the service registry (see Section 9). Each service definition declares whether it has a web UI and what port to use.

#### Error Format

```
✗ Cannot open mailpit

  Mailpit is not configured for this project.

  Fix: Run `devdock service add mailpit` then `devdock up`.
```

---

### 6.2 `devdock run <command>`

#### Behavior

1. Load `.devdock.yml`
2. Look up `commands.<name>`
3. If not found: print error with available command list and exit 1
4. Execute via shell: `/bin/sh -c "<command string>"` in the project root directory
5. Pipe stdout and stderr directly to the terminal (no buffering)
6. Exit with the same code as the underlying command

Shell execution is required because recipe commands include compound expressions (`composer install && npm install`) and pipes. Direct `exec` would break these.

#### Output When No Args

```
Available commands for my-app:

  install    composer install && npm install
  dev        php artisan serve
  migrate    php artisan migrate
  seed       php artisan db:seed
  test       php artisan test

Run `devdock run <command>` to execute.
```

#### Error Format

```
✗ Command 'deploy' is not defined

  Available commands: install, dev, migrate, seed, test

  Fix: Add 'commands.deploy' to .devdock.yml, or run one of the commands above.
```

---

### 6.3 `devdock service add <name>`

#### Supported Services

PostgreSQL, MySQL, Redis, Mailpit, MinIO.

Attempting to add any other name:

```
✗ 'mongodb' is not a supported service

  Supported services: postgres, mysql, redis, mailpit, minio

  Fix: Add the service manually to compose.yml, then run `devdock up`.
```

#### Full Behavior

1. Load and validate `.devdock.yml`
2. Reject if `project.type: docker-compose` (see Section 11)
3. If service already enabled: print "already configured" message, exit 0
4. Merge service defaults from service registry into `.devdock.yml`
5. Write `.devdock.yml` atomically (temp file → rename)
6. Regenerate `compose.yml`
7. Print environment variable block (see format below)
8. Prompt: "Append these variables to .env? (Y/n)"
   - If yes: append only missing keys via `env.MergeSafe()`
   - If no: print "You can add them manually when needed."
9. Print next step: "Run `devdock up` to start the new service."

#### Environment Variable Output Format

```
Environment variables for mailpit:

  MAIL_MAILER=smtp
  MAIL_HOST=127.0.0.1
  MAIL_PORT=1025
  MAIL_USERNAME=null
  MAIL_PASSWORD=null
  MAIL_ENCRYPTION=null

Append these to .env? (Y/n)
```

#### Running Services Behavior

If `devdock up` is currently running (Docker services are healthy), `devdock service add` still modifies the config files. After printing env vars, it adds:

```
⚠ Services are currently running.

  The new service will not start until you run `devdock up` again.
  Existing running services will not be interrupted.
```

`devdock service add` never touches running containers.

---

### 6.4 `devdock service remove <name>`

#### Full Behavior

1. Load `.devdock.yml`
2. Reject if `project.type: docker-compose`
3. If service not found: print "not configured" message, exit 0
4. If service has a named volume: print volume warning (non-blocking)
5. Prompt: "Remove <name> from this project? (y/N)" — default N
6. If confirmed: set `services.<name>.enabled: false` and remove service block
7. Write `.devdock.yml` atomically
8. Regenerate `compose.yml`
9. Print: "Run `devdock down && devdock up` to apply the change."

Volume warning format:

```
⚠ redis stores data in a named volume: devdock-my-app-redis-data

  This data will remain until you run:
    devdock down --volumes

  Removing the service from DevDock does not delete existing data.
```

Confirmation prompt:

```
Remove redis from this project? (y/N):
```

**Volumes are never touched by `devdock service remove`.** Only config files are modified.

---

### 6.5 `devdock service status`

Shows services only (not app process). This is distinct from `devdock status` which shows the full project overview.

```
Services — my-app

Service     Status    Port    Healthy   Web URL
───────────────────────────────────────────────────────
postgres    running   5432    ✔         —
redis       running   6379    ✔         —
mailpit     running   1025    ✔         http://localhost:8025
minio       running   9000    ✔         http://localhost:9001
```

For services with two ports (Mailpit, MinIO): the primary port shown is the application port (SMTP for Mailpit, API for MinIO). The web UI port is shown in the Web URL column.

If no services are running: "No services are running. Run `devdock up` to start."

---

### 6.6 `devdock service logs <name>`

Thin wrapper around the existing `devdock logs` implementation.

```bash
devdock service logs postgres
devdock service logs mailpit --tail 50
devdock service logs minio --since 5m
```

If the named service is not configured in `.devdock.yml`:

```
✗ mailpit is not configured in this project.

  Fix: Run `devdock service add mailpit` then `devdock up`.
```

---

### 6.7 `devdock restart <service>` (P2 — Optional)

**Only implement if Week 5 has capacity after all P0/P1 work is tested.**

Scope if implemented:

```bash
devdock restart postgres   # stop and start a Docker service
devdock restart redis
```

Not in scope for v0.2:

- `devdock restart app` — excluded; foreground app restart is Ctrl+C + `devdock up`
- `devdock restart` (no args) — excluded

Implementation: `docker compose restart <service>` in the project directory.

---

## 7. New Services Specification

### 7.1 Image Versioning Policy

Both new services must use pinned Docker image versions, not `latest`. The `:latest` tag makes builds non-reproducible and can silently break on `docker pull`.

| Service | Pinned Image | Notes |
|---|---|---|
| Mailpit | `axllent/mailpit:v1.21` | Update quarterly; test before updating |
| MinIO | `minio/minio:RELEASE.2024-11-07T00-52-20Z` | Update quarterly; test before updating |

The version is stored in the service definition Go struct and in the service registry. Upgrading service versions requires a DevDock release.

---

### 7.2 Mailpit

#### Docker Compose Output

```yaml
mailpit:
  image: axllent/mailpit:v1.21
  restart: unless-stopped
  ports:
    - "1025:1025"    # SMTP
    - "8025:8025"    # Web UI
  networks:
    - devdock-my-app
  healthcheck:
    test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8025/api/v1/info"]
    interval: 5s
    timeout: 5s
    retries: 6
```

**Health check:** `GET /api/v1/info` — Mailpit's own health endpoint. Returns JSON on healthy.

Mailpit has no persistent volume (email data is in-memory only). This is intentional — emails clear on restart.

#### `.devdock.yml` Config

```yaml
services:
  mailpit:
    enabled: true
    version: "v1.21"
    smtp_port: 1025
    ui_port: 8025
```

#### Environment Variables by Stack

**Laravel:**
```env
MAIL_MAILER=smtp
MAIL_HOST=127.0.0.1
MAIL_PORT=1025
MAIL_USERNAME=null
MAIL_PASSWORD=null
MAIL_ENCRYPTION=null
```

**Express / Node.js:**
```env
SMTP_HOST=127.0.0.1
SMTP_PORT=1025
```

---

### 7.3 MinIO

#### Docker Compose Output

```yaml
minio:
  image: minio/minio:RELEASE.2024-11-07T00-52-20Z
  restart: unless-stopped
  command: server /data --console-address ":9001"
  environment:
    MINIO_ROOT_USER: devdock
    MINIO_ROOT_PASSWORD: devdock-secret
  ports:
    - "9000:9000"    # API
    - "9001:9001"    # Console
  volumes:
    - devdock-my-app-minio-data:/data
  networks:
    - devdock-my-app
  healthcheck:
    test: ["CMD", "mc", "ready", "local"]
    interval: 5s
    timeout: 5s
    retries: 10

minio-init:
  image: minio/mc:latest
  depends_on:
    minio:
      condition: service_healthy
  entrypoint: >
    /bin/sh -c "
    mc alias set local http://minio:9000 devdock devdock-secret;
    mc mb --ignore-existing local/local;
    exit 0;
    "
  networks:
    - devdock-my-app
```

**MinIO bucket auto-creation:** A `minio-init` one-shot container runs after MinIO is healthy. It creates the `local` bucket using `mc` (MinIO Client). This runs on every `devdock up` but `--ignore-existing` makes it idempotent. Without this, the env vars referencing `AWS_BUCKET=local` would fail silently until the user manually created the bucket via the console.

**Health check:** Uses `mc ready local` to verify MinIO is ready to accept requests.

#### `.devdock.yml` Config

```yaml
services:
  minio:
    enabled: true
    version: "RELEASE.2024-11-07T00-52-20Z"
    api_port: 9000
    console_port: 9001
    access_key: devdock
    secret_key: devdock-secret
    bucket: local
```

#### Environment Variables by Stack

**Laravel:**
```env
FILESYSTEM_DISK=s3
AWS_ACCESS_KEY_ID=devdock
AWS_SECRET_ACCESS_KEY=devdock-secret
AWS_DEFAULT_REGION=us-east-1
AWS_BUCKET=local
AWS_ENDPOINT=http://127.0.0.1:9000
AWS_USE_PATH_STYLE_ENDPOINT=true
```

**Express / Node.js:**
```env
S3_ENDPOINT=http://127.0.0.1:9000
S3_ACCESS_KEY_ID=devdock
S3_SECRET_ACCESS_KEY=devdock-secret
S3_BUCKET=local
S3_REGION=us-east-1
```

---

## 8. New Stack Support

### 8.1 Express

#### Detection Logic (Priority Order)

| Condition | Confidence | Result |
|---|---|---|
| `package.json` exists AND `dependencies.express` key present | High | `express` |
| `package.json` exists AND no `express` dependency but `devDependencies.express` present | Medium | `express` (warn: unusual) |
| `package.json` exists AND no express in any dependencies | Low | Ask user |
| `package.json` missing | — | Not Express |

Detection reads `package.json` and parses it as JSON. Checks the `dependencies` field first. Does not string-search — must parse the JSON structure.

If Next.js was already detected (higher priority rule), the Express rule is skipped.

#### App Command Detection (Recipe)

The Express recipe checks `scripts` in `package.json` in this priority order:

1. If `scripts.dev` exists → use `npm run dev`
2. If `scripts.start` exists → use `npm start`
3. Fallback → `node index.js` and warn: "Could not detect start command. Defaulting to `node index.js`. Update `app.command` in `.devdock.yml` if needed."

#### Recipe (`recipes/express.yml`)

```yaml
id: express
name: Express

runtime:
  node: "22"

app:
  port: 3000
  run_mode: host

services:
  suggested:
    - postgres
    - redis

env_templates:
  postgres:
    DATABASE_URL: "postgresql://devdock:secret@localhost:5432/{{project_name}}"
    DB_HOST: "127.0.0.1"
    DB_PORT: "5432"
    DB_NAME: "{{project_name}}"
    DB_USER: "devdock"
    DB_PASSWORD: "secret"
  redis:
    REDIS_URL: "redis://localhost:6379"

commands:
  install: npm install
  dev: npm run dev
  start: npm start
  test: npm test
```

**Note on commands:** The recipe provides defaults. `devdock init` reads the actual `scripts` block from `package.json` and merges them into the generated `.devdock.yml`. Commands found in `package.json` take precedence over recipe defaults.

#### Doctor Check for Express

`devdock doctor` adds: validate `node --version` matches `runtime.node` in config.

---

### 8.2 Go Fiber

#### Detection Logic

| Condition | Confidence | Result |
|---|---|---|
| `go.mod` exists AND file contains `github.com/gofiber/fiber` in require block | High | `go-fiber` |
| `go.mod` exists AND no fiber dependency | Low | Ask user (could be generic Go) |
| `go.mod` missing | — | Not Go |

**Parsing `go.mod`:** Read the file as text. Look for any line (after the `require` keyword block) that contains `github.com/gofiber/fiber`. Do not use a full Go module parser — a string search for `github.com/gofiber/fiber` after the `module` line is sufficient for detection purposes.

#### App Command and Root Assumption

`go run .` requires a `package main` in the project root directory. This is the most common Go project layout for simple APIs.

If `main.go` is not found in the project root: `devdock init` warns:

```
⚠ No main.go found in project root.

  DevDock defaults to `go run .` which requires package main in the root directory.
  If your main package is in a subdirectory (e.g., cmd/server/), update app.command
  in .devdock.yml:

    app:
      command: go run ./cmd/server
```

This is a warning, not an error. The `.devdock.yml` is still generated.

#### Recipe (`recipes/go-fiber.yml`)

```yaml
id: go-fiber
name: Go Fiber

runtime:
  go: "1.22"

app:
  command: go run .
  port: 8080
  run_mode: host

services:
  suggested:
    - postgres
    - redis

env_templates:
  postgres:
    DATABASE_URL: "postgresql://devdock:secret@localhost:5432/{{project_name}}?sslmode=disable"
    DB_HOST: "127.0.0.1"
    DB_PORT: "5432"
    DB_NAME: "{{project_name}}"
    DB_USER: "devdock"
    DB_PASSWORD: "secret"
  redis:
    REDIS_URL: "redis://localhost:6379"
    REDIS_HOST: "127.0.0.1"
    REDIS_PORT: "6379"

commands:
  dev: go run .
  test: go test ./...
  build: go build -o bin/app .
```

#### Doctor Check for Go Fiber

`devdock doctor` adds:
- Validate `go version` output matches `runtime.go` in config
- If `go` is not installed: print install link (https://go.dev/dl/)

---

## 9. Service Registry Interface

The service registry is a Go map (`map[string]ServiceDefinition`) that all service-related commands use as the source of truth. This prevents hardcoded service logic scattered across commands.

**`ServiceDefinition` Go struct:**

```go
type ServiceDefinition struct {
    // Identity
    Name        string   // "mailpit"
    DisplayName string   // "Mailpit"

    // Docker
    Image       string   // "axllent/mailpit:v1.21"
    HasVolume   bool     // true for stateful services

    // Ports
    Ports []PortDefinition

    // Web UI
    WebUI *WebUIDefinition // nil if no web interface

    // Health check
    HealthCheck HealthCheckDefinition

    // Env var templates by stack type
    EnvTemplates map[string]map[string]string // map[stackType]map[varName]varValue

    // Default .devdock.yml service config (merged on service add)
    DefaultConfig map[string]interface{}
}

type PortDefinition struct {
    Name      string // "smtp", "api", "ui", "console"
    Container int
    Host      int    // default host port (user can override in .devdock.yml)
    Primary   bool   // true = shown in `devdock status` port column
}

type WebUIDefinition struct {
    PortName string // references a PortDefinition.Name
    Label    string // "Mailpit Inbox", "MinIO Console"
}

type HealthCheckDefinition struct {
    Type    string // "http" | "tcp" | "exec"
    Path    string // for type=http
    Port    int    // for type=http or type=tcp
    Command []string // for type=exec
}
```

All five services (postgres, mysql, redis, mailpit, minio) are registered at startup. This makes adding a new service in v0.3+ a single-file change.

---

## 10. `.devdock.yml` Schema Changes

### 10.1 Backward Compatibility

v0.2 is fully backward compatible with v0.1 configs. All new fields are optional and have defaults. A v0.1 config loaded by DevDock v0.2 must work without modification.

### 10.2 New Project Types

```
express
go-fiber
```

### 10.3 New Service Fields

All new fields are optional in v0.2. When absent, service registry defaults apply.

```yaml
services:
  mailpit:
    enabled: true
    version: "v1.21"          # optional — defaults to registry value
    smtp_port: 1025           # optional — defaults to 1025
    ui_port: 8025             # optional — defaults to 8025

  minio:
    enabled: true
    version: "RELEASE.2024-11-07T00-52-20Z"  # optional — defaults to registry value
    api_port: 9000            # optional — defaults to 9000
    console_port: 9001        # optional — defaults to 9001
    access_key: devdock       # optional — defaults to "devdock"
    secret_key: devdock-secret  # optional — defaults to "devdock-secret"
    bucket: local             # optional — defaults to "local"
```

### 10.4 No Schema Version Bump

The `version: "1"` field remains unchanged. v0.2 additions are purely additive. A `devdock migrate-config` run is not required.

---

## 11. Docker Compose Project Rules

For `project.type: docker-compose`, DevDock operates in passthrough mode (see v0.1 spec). The following commands must be blocked for this project type:

- `devdock service add <name>`
- `devdock service remove <name>`

Error format:

```
✗ Service management is not available for docker-compose projects

  This project uses an existing compose.yml that DevDock does not own.
  DevDock will not modify it automatically.

  Fix: Add the service manually to compose.yml, then run `devdock up`.
```

The following commands work normally in passthrough mode:

- `devdock service status` — proxies to `docker compose ps`
- `devdock service logs <name>` — proxies to `docker compose logs <name>`
- `devdock open` — reads app port from `.devdock.yml` (user-specified)
- `devdock run <command>` — runs named commands from `.devdock.yml`

---

## 12. Non-Functional Requirements

### 12.1 Performance

| Operation | Target |
|---|---|
| `devdock open` | < 500ms from command to browser open |
| `devdock run <cmd>` (first output) | < 200ms to first output byte |
| `devdock service add` | < 2 seconds (config write + compose regen) |
| `devdock service remove` | < 2 seconds (config write + compose regen) |

### 12.2 Reliability

- All config writes are atomic (temp file → rename)
- All destructive operations require explicit confirmation
- `.env` is never overwritten — only appended to with missing keys
- `compose.yml` is only regenerated for DevDock-owned projects

### 12.3 Error Quality

Every error introduced in v0.2 must follow the what/why/fix format established in v0.1. No exceptions.

---

## 13. Backward Compatibility Contract

| Test | Method |
|---|---|
| v0.1 Laravel `.devdock.yml` loads without error in v0.2 | Automated fixture test |
| v0.1 Next.js `.devdock.yml` loads without error in v0.2 | Automated fixture test |
| v0.1 Docker Compose `.devdock.yml` loads without error in v0.2 | Automated fixture test |
| v0.1 configs do not trigger schema migration prompt | Automated fixture test |
| Unknown fields in `.devdock.yml` still log warning, not error | Automated fixture test |

These tests run in the existing test suite. A failing backward compatibility test blocks release.

---

## 14. 6-Week Build Plan

### Week 1 — `devdock open` + `devdock run`

**Goal:** The two most commonly-used new commands are done and tested.

Build: `devdock open` (all targets), `devdock run` (all behaviors), update Laravel and Next.js recipe command blocks, unit tests for both commands.

**Done when:** `devdock up --detach && devdock open` opens the browser. `devdock run migrate` on a Laravel project runs the migration with real-time output and correct exit code.

---

### Week 2 — Service Registry + Service Management Commands

**Goal:** Service management CLI is fully built; tested against existing services.

Build: Service registry abstraction (all 5 services), `devdock service add`, `devdock service remove`, `devdock service status`, `devdock service logs`, all unit tests.

**Done when:** `devdock service add redis` on a Next.js project (that doesn't have Redis) updates `.devdock.yml`, regenerates `compose.yml`, prints env vars, and correctly declines to append to `.env` if user says no. `devdock service remove redis` prompts and removes cleanly without deleting data.

---

### Week 3 — Mailpit and MinIO Services

**Goal:** Both new services start correctly, are healthy, and have correct env vars.

Build: Mailpit service definition (pinned image, health check), MinIO service definition (pinned image, health check, `minio-init` container for bucket creation), Compose generator updates for multi-port and init containers, all env var templates.

**Done when:** `devdock service add mailpit && devdock up && devdock open mailpit` → Mailpit inbox opens. `devdock service add minio && devdock up && devdock open minio` → MinIO console opens. After `devdock down && devdock up`, MinIO data (a test file uploaded via the console) is still present.

---

### Week 4 — Express and Go Fiber Support

**Goal:** Both new stacks detect correctly and run through DevDock.

Build: Express detection (JSON parser, script detection fallback chain), Express recipe, Go Fiber detection (`go.mod` parser), Go Fiber recipe, doctor updates for both runtimes, fixture projects for detector tests.

**Done when:** An existing Express app passes `devdock init && devdock up && curl localhost:3000`. An existing Go Fiber app passes `devdock init && devdock up && curl localhost:8080`. Both pass on a clean macOS machine.

---

### Week 5 — Polish + Doctor Updates + Backward Compatibility

**Goal:** All errors are correct. All new stacks are recognized by doctor. Backward compat is automated.

Build: Error message audit (all new errors against what/why/fix format), doctor checks for Express and Go Fiber runtimes, Mailpit/MinIO port checks in doctor, automated backward compatibility fixture tests, `devdock restart <service>` if capacity allows.

**Done when:** All errors from Weeks 1–4 pass the error format audit. Doctor passes for all 4 stacks. All 5 backward compatibility fixture tests pass.

---

### Week 6 — Integration Testing + Release

**Goal:** v0.2.0 shipped.

Build: All 9 end-to-end flows (below), README updates, release binary, GitHub Release.

**Done when:** All Definition of Done items checked. Release tagged and published.

---

## 15. Definition of Done

v0.2.0 is not released until every item below is checked. `devdock restart` is excluded from this list.

### Workflow Commands

- [ ] `devdock open` opens app URL in browser
- [ ] `devdock open mailpit` opens Mailpit UI
- [ ] `devdock open minio` opens MinIO console
- [ ] `devdock open postgres` prints "no web interface" message (does not error)
- [ ] `devdock open <unknown>` lists available targets
- [ ] `devdock run <command>` runs with real-time output and correct exit code
- [ ] `devdock run` (no args) lists available commands
- [ ] `devdock run <undefined>` prints helpful error with available commands

### Service Management

- [ ] `devdock service add mailpit` updates config and prints env vars
- [ ] `devdock service add minio` updates config and prints env vars
- [ ] `devdock service add <already-configured>` is a no-op with clear message
- [ ] `devdock service add` while running prints restart-required notice, does not touch containers
- [ ] `devdock service remove <name>` prompts and removes without deleting volume data
- [ ] `devdock service remove <stateful-service>` shows volume warning before confirmation
- [ ] `devdock service add/remove` blocked for `docker-compose` project type
- [ ] `devdock service status` shows correct output with dual-port services
- [ ] `devdock service logs <name>` works with `--tail` and `--since`

### New Services

- [ ] Mailpit starts and health check passes
- [ ] Mailpit web UI is accessible at `http://localhost:8025`
- [ ] Mailpit SMTP is accessible on port 1025
- [ ] MinIO starts and health check passes
- [ ] MinIO `local` bucket is auto-created on `devdock up`
- [ ] MinIO console is accessible at `http://localhost:9001`
- [ ] MinIO data persists across `devdock down` + `devdock up`
- [ ] MinIO data is deleted by `devdock down --volumes`

### New Stacks

- [ ] Express detected from `package.json` dependencies (high confidence)
- [ ] Express command fallback chain works (dev → start → `node index.js`)
- [ ] Express `devdock init && devdock up && curl localhost:3000` passes
- [ ] Go Fiber detected from `go.mod` Fiber import (high confidence)
- [ ] Go Fiber main.go warning shown when no root main.go
- [ ] Go Fiber `devdock init && devdock up && curl localhost:8080` passes
- [ ] Go runtime validation in `devdock doctor`

### Backward Compatibility

- [ ] v0.1 Laravel config loads without error or migration prompt
- [ ] v0.1 Next.js config loads without error or migration prompt
- [ ] v0.1 Docker Compose config loads without error or migration prompt
- [ ] Automated backward compat fixture tests pass (5 tests)

---

## 16. Risks & Mitigations

| Risk | Likelihood | Mitigation |
|---|---|---|
| MinIO bucket init container complicates Compose output | Medium | Test early in Week 3; document the init container pattern for future services |
| Express command detection wrong for custom project setups | High | Make fallback chain visible in init output; user can edit `.devdock.yml` |
| Go Fiber detection false-positives on projects that import Fiber transitively | Low | Require the `require` block in go.mod, not just any occurrence of the string |
| Pinned image versions become outdated | Medium | Add a quarterly image review to the release checklist |
| `devdock service add` while running causes confusion | Medium | The "running services will not be affected" message mitigates this; covered in Definition of Done |
| `restart` scope creep delays release | Medium | `restart` is not in Definition of Done; cut it cleanly if not done by Week 5 end |
