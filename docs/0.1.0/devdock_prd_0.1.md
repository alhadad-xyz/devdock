# DevDock v0.1 Revised Task Checklist

**Type:** Engineering task checklist  
**Scope:** DevDock v0.1 solo developer build  
**Target:** 8 weeks  
**Source:** DevDock v0.1 Engineering Spec Rev 1  
**Status:** Ready for Week 1 implementation  

---

## Implementation Priorities

### P0 — Engine Works

These tasks prove the core DevDock engine can read config, detect projects, generate services, and run the app.

- `[ ]` Config loader
- `[ ]` Project detector
- `[ ]` Docker Compose generator
- `[ ]` `devdock init`
- `[ ]` `devdock up`
- `[ ]` `devdock down`

### P1 — Usable Developer Tool

These tasks make DevDock useful enough for real local development.

- `[ ]` `devdock doctor`
- `[ ]` `devdock logs`
- `[ ]` `devdock status`
- `[ ]` Laravel recipe polish
- `[ ]` Docker Compose project support
- `[ ]` Clear error formatting

### P2 — Product Demo & Release Polish

These tasks make DevDock demoable and ready for external users.

- `[ ]` `devdock create next-postgres`
- `[ ]` `devdock create laravel-api`
- `[ ]` Install script
- `[ ]` Release README
- `[ ]` GitHub Release
- `[ ]` External tester validation

---

# Week 1 — CLI Skeleton & Configuration

**Goal:** `devdock --help` works, `.devdock.yml` loads and validates, and the project has the basic internal structure needed for later weeks.

## Core CLI

- `[ ]` **DD-001** Initialize Go module.
- `[ ]` **DD-002** Add Cobra CLI structure.
- `[ ]` **DD-003** Create `cmd/devdock/main.go` entry point.
- `[ ]` **DD-004** Create `internal/cli/root.go` root command.
- `[ ]` **DD-005** Implement `devdock` with no arguments welcome output.
- `[ ]` **DD-006** Implement `--help` support.
- `[ ]` **DD-007** Implement `--version` support.
- `[ ]` **DD-008** Implement internal version variable.
- `[ ]` **DD-009** Inject version using build flags.
- `[ ]` **DD-010** Implement `--quiet` global flag.
- `[ ]` **DD-011** Implement `--json` global flag placeholder.
- `[ ]` **DD-012** Implement `--project` global flag or project root resolver.
- `[ ]` **DD-013** Resolve project path from `--project` or current directory.

## Friendly Unknown Command Handling

- `[ ]` **DD-014** Implement unknown command handler.
- `[ ]` **DD-015** Unknown command should print:

```txt
'devdock <command>' is not available in this version.

Run `devdock --help` to see available commands.
```

- `[ ]` **DD-016** Unknown command exits with code `1`.

## DevDock Home Directory

- `[ ]` **DD-017** Create DevDock home directory initializer.
- `[ ]` **DD-018** Create `~/.devdock/` if missing.
- `[ ]` **DD-019** Create `~/.devdock/pids/` if missing.
- `[ ]` **DD-020** Create `~/.devdock/logs/` if missing.
- `[ ]` **DD-021** Create default `~/.devdock/config.yml` if missing.
- `[ ]` **DD-022** Define `~/.devdock/config.yml` schema.

Example:

```yaml
version: "1"
defaults:
  package_manager: pnpm
  editor: code
```

## Centralized Error Handling

- `[ ]` **DD-023** Define centralized app error type.
- `[ ]` **DD-024** Add error category constants.
- `[ ]` **DD-025** Add exit code mapping.
- `[ ]` **DD-026** Implement what/why/fix error formatter.
- `[ ]` **DD-027** Ensure user-facing errors never show stack traces.
- `[ ]` **DD-028** Write stack traces or debug details to `~/.devdock/logs/error.log`.

Required error format:

```txt
✗ <What went wrong>

  <Why it happened, if known>

  Fix: <One specific action to take>
```

## `.devdock.yml` Schema

- `[ ]` **DD-029** Create `internal/config/schema.go`.
- `[ ]` **DD-030** Define `Config` struct.
- `[ ]` **DD-031** Define `ProjectConfig` struct.
- `[ ]` **DD-032** Define `RuntimeConfig` struct.
- `[ ]` **DD-033** Define `AppConfig` struct.
- `[ ]` **DD-034** Define `ServicesConfig` struct.
- `[ ]` **DD-035** Define `PostgresConfig` struct.
- `[ ]` **DD-036** Define `MySQLConfig` struct.
- `[ ]` **DD-037** Define `RedisConfig` struct.
- `[ ]` **DD-038** Define `CommandsConfig` map.

## Config Loader

- `[ ]` **DD-039** Create `internal/config/loader.go`.
- `[ ]` **DD-040** Load `.devdock.yml` from project root.
- `[ ]` **DD-041** Parse YAML into config struct.
- `[ ]` **DD-042** Validate `version` is exactly `"1"`.
- `[ ]` **DD-043** Validate `project.name` exists.
- `[ ]` **DD-044** Validate `project.name` matches `^[a-z0-9-]+$`.
- `[ ]` **DD-045** Validate `project.type` is one of `laravel`, `nextjs`, `docker-compose`.
- `[ ]` **DD-046** Validate `app.command` is required for `laravel` and `nextjs`.
- `[ ]` **DD-047** Validate `app.port` is required for `laravel` and `nextjs`.
- `[ ]` **DD-048** Validate `app.command` and `app.port` are optional and ignored for `docker-compose`.
- `[ ]` **DD-049** Validate `app.run_mode` supports only `host` or empty in v0.1.
- `[ ]` **DD-050** If `app.run_mode: container`, return clear “not supported yet” error.
- `[ ]` **DD-051** Warn on unknown fields without failing.
- `[ ]` **DD-052** Normalize service enabled flags.
- `[ ]` **DD-053** Validate service ports are valid TCP ports.
- `[ ]` **DD-054** Validate no duplicate configured ports inside the same config.

## Basic `doctor` Placeholder

- `[ ]` **DD-055** Create `internal/cli/doctor.go`.
- `[ ]` **DD-056** Implement basic `devdock doctor` command.
- `[ ]` **DD-057** Check Docker installed using `docker --version`.
- `[ ]` **DD-058** Check Docker daemon running using `docker info`.
- `[ ]` **DD-059** Check Docker Compose v2 using `docker compose version`.
- `[ ]` **DD-060** If `.devdock.yml` exists, parse and validate it.
- `[ ]` **DD-061** If `.devdock.yml` does not exist, show helpful message but do not crash.

## Tests

- `[ ]` **DD-062** Write config loader unit tests for valid `nextjs` config.
- `[ ]` **DD-063** Write config loader unit tests for valid `laravel` config.
- `[ ]` **DD-064** Write config loader unit tests for valid `docker-compose` config.
- `[ ]` **DD-065** Write config loader unit tests for invalid version.
- `[ ]` **DD-066** Write config loader unit tests for invalid project name.
- `[ ]` **DD-067** Write config loader unit tests for invalid project type.
- `[ ]` **DD-068** Write config loader unit tests for missing `app.command` on `nextjs`.
- `[ ]` **DD-069** Write config loader unit tests for missing `app.port` on `laravel`.
- `[ ]` **DD-070** Write config loader unit tests proving `docker-compose` does not require `app.command` or `app.port`.

## Week 1 Done When

- `[ ]` `devdock --version` prints the current version.
- `[ ]` `devdock --help` prints available v0.1 commands only.
- `[ ]` `devdock` with no arguments prints the welcome message.
- `[ ]` `devdock doctor` runs without panic.
- `[ ]` `devdock doctor` can parse a valid `.devdock.yml`.
- `[ ]` `devdock doctor` reports schema errors for invalid `.devdock.yml`.
- `[ ]` Config loader tests pass.

---

# Week 2 — Project Detector

**Goal:** `devdock detect` correctly identifies Laravel, Next.js, and existing Docker Compose projects.

## File Scanner

- `[ ]` **DD-071** Create `internal/detector/detector.go`.
- `[ ]` **DD-072** Implement current-directory file scanner.
- `[ ]` **DD-073** Do not traverse subdirectories in v0.1.
- `[ ]` **DD-074** Detect file presence by exact filename.
- `[ ]` **DD-075** Support `docker-compose.yml`.
- `[ ]` **DD-076** Support `compose.yml`.
- `[ ]` **DD-077** Support `next.config.js`.
- `[ ]` **DD-078** Support `next.config.ts`.
- `[ ]` **DD-079** Support `composer.json`.
- `[ ]` **DD-080** Support `artisan`.
- `[ ]` **DD-081** Support `package.json`.

## Detection Rules

- `[ ]` **DD-082** Detect `composer.json` + `artisan` as `laravel` with high confidence.
- `[ ]` **DD-083** Detect `package.json` + `next.config.js` as `nextjs` with high confidence.
- `[ ]` **DD-084** Detect `package.json` + `next.config.ts` as `nextjs` with high confidence.
- `[ ]` **DD-085** Detect `docker-compose.yml` as `docker-compose` with high confidence.
- `[ ]` **DD-086** Detect `compose.yml` as `docker-compose` with high confidence.
- `[ ]` **DD-087** Detect `composer.json` only as `laravel` with low confidence.
- `[ ]` **DD-088** Detect `package.json` only as `nextjs` with low confidence.
- `[ ]` **DD-089** Return unknown result if no supported files are found.

## Detection Output

- `[ ]` **DD-090** Create detection result struct: type, confidence, reasons.
- `[ ]` **DD-091** Implement `devdock detect` command.
- `[ ]` **DD-092** Print detected project type.
- `[ ]` **DD-093** Print confidence level.
- `[ ]` **DD-094** Print reasons or matched files.
- `[ ]` **DD-095** Add `--json` output support for `devdock detect`.

## Low-Confidence Flow

- `[ ]` **DD-096** Implement low-confidence confirmation prompt.
- `[ ]` **DD-097** If user accepts, continue with guessed type.
- `[ ]` **DD-098** If user rejects, show supported type picker.
- `[ ]` **DD-099** Supported picker options: `laravel`, `nextjs`, `docker-compose`.
- `[ ]` **DD-100** Persist selected type into later `init` flow.

## Unknown Stack Error

- `[ ]` **DD-101** Implement unknown stack error message.
- `[ ]` **DD-102** Message should list supported types.
- `[ ]` **DD-103** Message should show manual init examples.

Required output:

```txt
DevDock could not detect a supported project type.

Supported types in this version: laravel, nextjs, docker-compose

To initialize manually, run:
  devdock init --type=laravel
  devdock init --type=nextjs
```

## Test Fixtures

- `[ ]` **DD-104** Create `testdata/fixtures/laravel/`.
- `[ ]` **DD-105** Add minimal Laravel fixture files: `composer.json`, `artisan`.
- `[ ]` **DD-106** Create `testdata/fixtures/nextjs/`.
- `[ ]` **DD-107** Add minimal Next.js fixture files: `package.json`, `next.config.js`.
- `[ ]` **DD-108** Create `testdata/fixtures/docker-compose/`.
- `[ ]` **DD-109** Add minimal Docker Compose fixture file: `compose.yml`.
- `[ ]` **DD-110** Create low-confidence Composer fixture.
- `[ ]` **DD-111** Create low-confidence package fixture.
- `[ ]` **DD-112** Create unknown project fixture.

## Tests

- `[ ]` **DD-113** Test Laravel high-confidence detection.
- `[ ]` **DD-114** Test Next.js high-confidence detection via `next.config.js`.
- `[ ]` **DD-115** Test Next.js high-confidence detection via `next.config.ts`.
- `[ ]` **DD-116** Test Docker Compose detection via `docker-compose.yml`.
- `[ ]` **DD-117** Test Docker Compose detection via `compose.yml`.
- `[ ]` **DD-118** Test Laravel low-confidence detection.
- `[ ]` **DD-119** Test Next.js low-confidence detection.
- `[ ]` **DD-120** Test unknown project behavior.

## Week 2 Done When

- `[ ]` `devdock detect` correctly identifies Laravel fixture.
- `[ ]` `devdock detect` correctly identifies Next.js fixture.
- `[ ]` `devdock detect` correctly identifies Docker Compose fixture.
- `[ ]` Low-confidence project asks for confirmation.
- `[ ]` Unknown project prints useful manual init guidance.
- `[ ]` Detector tests pass with 100% pass rate on fixture suite.

---

# Week 3 — Docker Compose Generator & `devdock init`

**Goal:** `devdock init` creates `.devdock.yml`, `compose.yml`, and `.env` for Laravel and Next.js, while Docker Compose projects get a minimal `.devdock.yml` only.

## Naming & Normalization

- `[ ]` **DD-121** Add project name normalizer.
- `[ ]` **DD-122** Convert `My App` to `my-app`.
- `[ ]` **DD-123** Convert `my_app` to `my-app` for project names.
- `[ ]` **DD-124** Add database name normalizer.
- `[ ]` **DD-125** Convert `my-app` to `my_app` for database names.
- `[ ]` **DD-126** Add Compose project name strategy.
- `[ ]` **DD-127** Ensure Compose project names avoid conflicts across projects.
- `[ ]` **DD-128** Ensure volume names are deterministic.
- `[ ]` **DD-129** Ensure network names are deterministic.

## Service Definitions

- `[ ]` **DD-130** Create `internal/services/service.go` shared service interface.
- `[ ]` **DD-131** Implement PostgreSQL service definition.
- `[ ]` **DD-132** Implement PostgreSQL image tag support.
- `[ ]` **DD-133** Implement PostgreSQL environment variables.
- `[ ]` **DD-134** Implement PostgreSQL health check using `pg_isready`.
- `[ ]` **DD-135** Implement PostgreSQL volume mapping.
- `[ ]` **DD-136** Implement MySQL service definition.
- `[ ]` **DD-137** Implement MySQL image tag support.
- `[ ]` **DD-138** Implement MySQL environment variables.
- `[ ]` **DD-139** Implement MySQL health check using `mysqladmin ping`.
- `[ ]` **DD-140** Implement MySQL volume mapping.
- `[ ]` **DD-141** Implement Redis service definition.
- `[ ]` **DD-142** Implement Redis image tag support.
- `[ ]` **DD-143** Implement Redis health check using `redis-cli ping`.
- `[ ]` **DD-144** Implement Redis volume mapping if needed.

## Compose Generator

- `[ ]` **DD-145** Create `internal/compose/generator.go`.
- `[ ]` **DD-146** Generate `compose.yml` from `.devdock.yml`.
- `[ ]` **DD-147** Generate Compose `name` field.
- `[ ]` **DD-148** Generate enabled services only.
- `[ ]` **DD-149** Generate isolated bridge network.
- `[ ]` **DD-150** Generate isolated named volumes.
- `[ ]` **DD-151** Generate service health checks.
- `[ ]` **DD-152** Generate configured host port mappings.
- `[ ]` **DD-153** Generate updated header:

```yaml
# Generated by DevDock from .devdock.yml
# devdock version: 0.1.0
# Manual edits to this file may be overwritten. Edit .devdock.yml instead.
```

## Generated File Ownership

- `[ ]` **DD-154** Add generated-file ownership detection for `compose.yml`.
- `[ ]` **DD-155** Detect DevDock-owned `compose.yml` by header.
- `[ ]` **DD-156** If `compose.yml` has DevDock header, allow regeneration without prompt.
- `[ ]` **DD-157** If `compose.yml` exists without DevDock header, show diff before overwrite.
- `[ ]` **DD-158** If user declines overwrite, exit safely.
- `[ ]` **DD-159** Write files atomically using temp file + rename.

## Environment Generator

- `[ ]` **DD-160** Create `internal/env/generator.go`.
- `[ ]` **DD-161** Generate Laravel MySQL environment values.
- `[ ]` **DD-162** Generate Laravel Redis environment values.
- `[ ]` **DD-163** Generate Next.js PostgreSQL `DATABASE_URL`.
- `[ ]` **DD-164** Generate Next.js Redis URL if enabled.
- `[ ]` **DD-165** Generate `.env.example` with placeholder values.
- `[ ]` **DD-166** Add `.env` overwrite protection.
- `[ ]` **DD-167** Add `.env.example` overwrite protection.
- `[ ]` **DD-168** Implement `.gitignore` update helper for `.env` and `.devdock.local.yml` future compatibility.

## `devdock init` Flow

- `[ ]` **DD-169** Create `internal/cli/init.go`.
- `[ ]` **DD-170** Implement `devdock init` command.
- `[ ]` **DD-171** Add `--type` flag.
- `[ ]` **DD-172** Add `--db` flag.
- `[ ]` **DD-173** Add `--redis` flag.
- `[ ]` **DD-174** Add `--force` flag.
- `[ ]` **DD-175** If `.devdock.yml` exists and `--force` is false, exit with helpful message.
- `[ ]` **DD-176** Run project detector.
- `[ ]` **DD-177** Display detection result.
- `[ ]` **DD-178** Display suggested config.
- `[ ]` **DD-179** Ask confirmation for Laravel and Next.js.
- `[ ]` **DD-180** Implement service selection prompt.
- `[ ]` **DD-181** Default Laravel services: MySQL enabled, Redis optional.
- `[ ]` **DD-182** Default Next.js services: PostgreSQL enabled, Redis optional.
- `[ ]` **DD-183** Write `.devdock.yml`.
- `[ ]` **DD-184** Generate `compose.yml`.
- `[ ]` **DD-185** Generate `.env`.
- `[ ]` **DD-186** Generate `.env.example`.
- `[ ]` **DD-187** Print next step: `devdock up`.
- `[ ]` **DD-188** Ensure `init` does not start services.

## Docker Compose Project Init Flow

- `[ ]` **DD-189** Detect Docker Compose project during `init`.
- `[ ]` **DD-190** Write minimal `.devdock.yml` only.
- `[ ]` **DD-191** Do not generate `compose.yml`.
- `[ ]` **DD-192** Do not generate `.env`.
- `[ ]` **DD-193** Do not write `app.command`.
- `[ ]` **DD-194** Do not write `app.port`.
- `[ ]` **DD-195** Print message explaining DevDock will proxy existing Compose file.

Minimal Docker Compose `.devdock.yml`:

```yaml
version: "1"

project:
  name: my-app
  type: docker-compose
```

## Tests

- `[ ]` **DD-196** Test Compose generation for PostgreSQL.
- `[ ]` **DD-197** Test Compose generation for MySQL.
- `[ ]` **DD-198** Test Compose generation for Redis.
- `[ ]` **DD-199** Test Compose generation for PostgreSQL + Redis.
- `[ ]` **DD-200** Test volume name generation.
- `[ ]` **DD-201** Test network name generation.
- `[ ]` **DD-202** Test Compose header generation.
- `[ ]` **DD-203** Test DevDock-owned Compose detection.
- `[ ]` **DD-204** Test project name normalization.
- `[ ]` **DD-205** Test database name normalization.
- `[ ]` **DD-206** Test Laravel `.env` generation.
- `[ ]` **DD-207** Test Next.js `.env` generation.
- `[ ]` **DD-208** Test Docker Compose project minimal config generation.

## Week 3 Done When

- `[ ]` Running `devdock init` in a Next.js fixture creates `.devdock.yml`.
- `[ ]` Running `devdock init` in a Next.js fixture creates valid `compose.yml`.
- `[ ]` Generated `compose.yml` can run with `docker compose up`.
- `[ ]` Running `devdock init` in a Laravel fixture creates expected MySQL/Redis config.
- `[ ]` Running `devdock init` in a Docker Compose fixture creates minimal `.devdock.yml` only.
- `[ ]` Generator tests pass.

---

# Week 4 — Core Lifecycle: `devdock up` and `devdock down`

**Goal:** DevDock can start services, run the app process, handle port conflicts, support foreground/detached modes, and cleanly stop everything.

## Docker Command Wrapper

- `[ ]` **DD-209** Create `internal/docker/client.go`.
- `[ ]` **DD-210** Implement wrapper for `docker compose up -d`.
- `[ ]` **DD-211** Implement wrapper for `docker compose down`.
- `[ ]` **DD-212** Implement wrapper for `docker compose down --volumes`.
- `[ ]` **DD-213** Implement wrapper for `docker compose ps`.
- `[ ]` **DD-214** Implement wrapper for `docker compose logs`.
- `[ ]` **DD-215** Implement wrapper for `docker compose pull` if needed.
- `[ ]` **DD-216** Ensure Docker commands run from project root.
- `[ ]` **DD-217** Capture stdout/stderr safely.

## Prerequisite Checks for `up`

- `[ ]` **DD-218** Check Docker installed.
- `[ ]` **DD-219** Check Docker daemon running.
- `[ ]` **DD-220** Check Docker Compose v2.
- `[ ]` **DD-221** Check required runtime exists for Laravel.
- `[ ]` **DD-222** Check required runtime exists for Next.js.
- `[ ]` **DD-223** Skip runtime checks for `docker-compose` type.
- `[ ]` **DD-224** Implement `--skip-checks` flag.

## Port Conflict Detection

- `[ ]` **DD-225** Create `internal/ports/checker.go`.
- `[ ]` **DD-226** Check app port availability.
- `[ ]` **DD-227** Check PostgreSQL port availability.
- `[ ]` **DD-228** Check MySQL port availability.
- `[ ]` **DD-229** Check Redis port availability.
- `[ ]` **DD-230** Identify PID holding port using `lsof` on macOS.
- `[ ]` **DD-231** Identify process name holding port.
- `[ ]` **DD-232** Print structured what/why/fix error for port conflicts.
- `[ ]` **DD-233** Name the exact `.devdock.yml` field to edit.
- `[ ]` **DD-234** Use exit code `4` for port conflicts.

Example:

```txt
✗ Port 5432 is already in use

  PID 8821 (postgres) is using this port.

  Fix: Change services.postgres.port in .devdock.yml to 5433,
       then run `devdock up` again.
```

## Service Startup

- `[ ]` **DD-235** Create `internal/cli/up.go`.
- `[ ]` **DD-236** Implement `devdock up` command.
- `[ ]` **DD-237** Add `--detach` flag.
- `[ ]` **DD-238** Add `--build` flag.
- `[ ]` **DD-239** Add `--skip-checks` flag.
- `[ ]` **DD-240** Load and validate config.
- `[ ]` **DD-241** Regenerate `compose.yml` if needed for Laravel/Next.js.
- `[ ]` **DD-242** Do not regenerate `compose.yml` for `docker-compose` type.
- `[ ]` **DD-243** Run service startup using `docker compose up -d`.
- `[ ]` **DD-244** Stream startup progress.
- `[ ]` **DD-245** Poll service health checks.
- `[ ]` **DD-246** Implement 30-second health timeout.
- `[ ]` **DD-247** On unhealthy service, print error and suggest `devdock logs <service>`.
- `[ ]` **DD-248** Use exit code `5` for service startup failure.

## Host App Process Runner

- `[ ]` **DD-249** Create `internal/process/runner.go`.
- `[ ]` **DD-250** Start app command on host.
- `[ ]` **DD-251** Ensure app process uses project root as working directory.
- `[ ]` **DD-252** Stream app stdout/stderr in foreground mode.
- `[ ]` **DD-253** Pass environment variables from `.env` where appropriate.
- `[ ]` **DD-254** Do not start host app process for `docker-compose` type.
- `[ ]` **DD-255** Detect if app port becomes reachable.
- `[ ]` **DD-256** Print app URL after startup.

## Foreground Mode

- `[ ]` **DD-257** Implement foreground mode as default.
- `[ ]` **DD-258** Capture Ctrl+C / SIGINT.
- `[ ]` **DD-259** Ctrl+C sends SIGINT to app process only.
- `[ ]` **DD-260** Docker services remain running after Ctrl+C.
- `[ ]` **DD-261** Print message reminding user to run `devdock down` to stop services.

## Detached Mode

- `[ ]` **DD-262** Create `internal/process/pid.go`.
- `[ ]` **DD-263** Start app process in background with `--detach`.
- `[ ]` **DD-264** Write PID to `~/.devdock/pids/<project-name>.pid`.
- `[ ]` **DD-265** Detect stale PID files.
- `[ ]` **DD-266** Remove stale PID files automatically.
- `[ ]` **DD-267** Prevent duplicate app process when `devdock up --detach` is run twice.
- `[ ]` **DD-268** Write detached app stdout/stderr to `~/.devdock/logs/<project>.app.log`.
- `[ ]` **DD-269** Print detached app log path.

## `devdock down`

- `[ ]` **DD-270** Create `internal/cli/down.go`.
- `[ ]` **DD-271** Implement `devdock down` command.
- `[ ]` **DD-272** If PID file exists, send SIGTERM to app process.
- `[ ]` **DD-273** Remove PID file after app process stops.
- `[ ]` **DD-274** If PID file is stale, remove it without error.
- `[ ]` **DD-275** If no PID file exists, continue without error.
- `[ ]` **DD-276** Run `docker compose down`.
- `[ ]` **DD-277** Add `--volumes` flag.
- `[ ]` **DD-278** Prompt before `docker compose down --volumes`.
- `[ ]` **DD-279** Default destructive confirmation answer to `No`.
- `[ ]` **DD-280** Ensure no orphan containers remain after down.

## Docker Compose Project Lifecycle

- `[ ]` **DD-281** For `docker-compose` type, `devdock up` runs `docker compose up -d` only.
- `[ ]` **DD-282** For `docker-compose` type, `devdock down` runs `docker compose down` only.
- `[ ]` **DD-283** For `docker-compose` type, no host app process is started.
- `[ ]` **DD-284** For `docker-compose` type, no PID file is written.

## Week 4 Done When

- `[ ]` `devdock up` in Next.js project starts PostgreSQL and app.
- `[ ]` App is reachable at `http://localhost:3000`.
- `[ ]` Ctrl+C stops only the app process.
- `[ ]` PostgreSQL remains accessible after Ctrl+C.
- `[ ]` `devdock down` stops Docker services.
- `[ ]` `devdock up --detach` writes PID file.
- `[ ]` `devdock down` stops detached app process.
- `[ ]` Port conflicts produce structured what/why/fix errors.
- `[ ]` `devdock down --volumes` prompts before deleting data.

---

# Week 5 — Framework Recipes, Logs & Status

**Goal:** Laravel setup works end-to-end, and DevDock can show useful logs and status for Laravel, Next.js, and Docker Compose projects.

## Laravel Recipe

- `[ ]` **DD-285** Create `recipes/laravel.yml`.
- `[ ]` **DD-286** Embed Laravel recipe via `go:embed`.
- `[ ]` **DD-287** Suggest PHP `8.3`.
- `[ ]` **DD-288** Suggest command `php artisan serve --host=127.0.0.1 --port=8000`.
- `[ ]` **DD-289** Suggest port `8000`.
- `[ ]` **DD-290** Suggest MySQL enabled by default.
- `[ ]` **DD-291** Suggest Redis optional.
- `[ ]` **DD-292** Generate Laravel `.env` database values.
- `[ ]` **DD-293** Generate Laravel `.env` Redis values if enabled.
- `[ ]` **DD-294** Ensure `APP_URL=http://localhost:8000`.

## Next.js Recipe Polish

- `[ ]` **DD-295** Create or finalize `recipes/next.yml`.
- `[ ]` **DD-296** Embed Next.js recipe via `go:embed`.
- `[ ]` **DD-297** Suggest Node `22`.
- `[ ]` **DD-298** Suggest command `pnpm dev`.
- `[ ]` **DD-299** Suggest port `3000`.
- `[ ]` **DD-300** Suggest PostgreSQL enabled by default.
- `[ ]` **DD-301** Suggest Redis optional.
- `[ ]` **DD-302** Generate `DATABASE_URL`.

## App Log Handling

- `[ ]` **DD-303** Ensure foreground app output is visible during `devdock up`.
- `[ ]` **DD-304** For detached mode, tail `~/.devdock/logs/<project>.app.log`.
- `[ ]` **DD-305** Implement fallback when app logs are unavailable.
- `[ ]` **DD-306** Print helpful message if app process is not running.

## `devdock logs`

- `[ ]` **DD-307** Create `internal/cli/logs.go`.
- `[ ]` **DD-308** Implement `devdock logs` command.
- `[ ]` **DD-309** Support `devdock logs` for all Docker service logs.
- `[ ]` **DD-310** Support `devdock logs postgres`.
- `[ ]` **DD-311** Support `devdock logs mysql`.
- `[ ]` **DD-312** Support `devdock logs redis`.
- `[ ]` **DD-313** Support `devdock logs app`.
- `[ ]` **DD-314** Support `--tail <n>`.
- `[ ]` **DD-315** Support `--since <duration>`.
- `[ ]` **DD-316** Prefix logs with service labels.
- `[ ]` **DD-317** Add service label coloring.
- `[ ]` **DD-318** Default to follow mode.

## Docker Compose Logs Behavior

- `[ ]` **DD-319** For `docker-compose` type, proxy `docker compose logs -f`.
- `[ ]` **DD-320** For `docker-compose` type, proxy `docker compose logs <service> -f`.
- `[ ]` **DD-321** For `docker-compose` type, `devdock logs app` should explain no host app process exists.

## `devdock status`

- `[ ]` **DD-322** Create `internal/cli/status.go`.
- `[ ]` **DD-323** Implement `devdock status` command.
- `[ ]` **DD-324** Load project config.
- `[ ]` **DD-325** Query Docker service status.
- `[ ]` **DD-326** Query Docker health status.
- `[ ]` **DD-327** Query app process status from PID file for detached app.
- `[ ]` **DD-328** For foreground app, show app status as unknown or stopped if no PID exists.
- `[ ]` **DD-329** Print project name, path, type, and mode.
- `[ ]` **DD-330** Print service table.
- `[ ]` **DD-331** Print app URL.
- `[ ]` **DD-332** Print database connection string for enabled DB.
- `[ ]` **DD-333** Implement Docker Compose project status separately.
- `[ ]` **DD-334** For `docker-compose` type, query `docker compose ps` and format output.

## Manual Laravel Test

- `[ ]` **DD-335** Clone official Laravel project.
- `[ ]` **DD-336** Run `devdock init`.
- `[ ]` **DD-337** Confirm Laravel detection.
- `[ ]` **DD-338** Confirm generated `.devdock.yml`.
- `[ ]` **DD-339** Confirm generated `compose.yml`.
- `[ ]` **DD-340** Run `devdock up`.
- `[ ]` **DD-341** Confirm MySQL health.
- `[ ]` **DD-342** Confirm Redis health if enabled.
- `[ ]` **DD-343** Confirm `curl localhost:8000` returns HTTP 200.
- `[ ]` **DD-344** Run `devdock down`.
- `[ ]` **DD-345** Confirm no orphan containers remain.

## Week 5 Done When

- `[ ]` Real Laravel app passes `devdock init → devdock up → curl localhost:8000 → devdock down`.
- `[ ]` `devdock logs postgres` streams PostgreSQL logs.
- `[ ]` `devdock logs app` works for detached app logs.
- `[ ]` `devdock status` shows service status and health.
- `[ ]` Docker Compose project status/logs behavior works separately from Laravel/Next.js.

---

# Week 6 — Scaffolding: `devdock create`

**Goal:** DevDock can create new Laravel and Next.js projects using official scaffolding tools, then add DevDock-specific files.

## Embedded Templates

- `[ ]` **DD-346** Create `templates/laravel-api/`.
- `[ ]` **DD-347** Create `templates/next-postgres/`.
- `[ ]` **DD-348** Add `template.yml` for Laravel template.
- `[ ]` **DD-349** Add `template.yml` for Next.js template.
- `[ ]` **DD-350** Add `.devdock.yml.tpl` for Laravel template.
- `[ ]` **DD-351** Add `.env.tpl` for Laravel template.
- `[ ]` **DD-352** Add `.env.example.tpl` for Laravel template.
- `[ ]` **DD-353** Add `.devdock.yml.tpl` for Next.js template.
- `[ ]` **DD-354** Add `prisma/schema.prisma.tpl` for Next.js template.
- `[ ]` **DD-355** Add `.env.tpl` for Next.js template.
- `[ ]` **DD-356** Add `.env.example.tpl` for Next.js template.
- `[ ]` **DD-357** Embed templates using Go `embed` package.

## Template Renderer

- `[ ]` **DD-358** Create `internal/template/renderer.go`.
- `[ ]` **DD-359** Implement variable substitution for `{{project_name}}`.
- `[ ]` **DD-360** Implement variable substitution for `{{project_name_snake}}`.
- `[ ]` **DD-361** Implement variable substitution for `{{project_name_pascal}}`.
- `[ ]` **DD-362** Implement variable substitution for `{{db_name}}`.
- `[ ]` **DD-363** Implement `{{random_hex_32}}`.
- `[ ]` **DD-364** Implement `{{random_hex_64}}`.
- `[ ]` **DD-365** Add renderer tests.

## Scaffold Runner

- `[ ]` **DD-366** Create `internal/template/scaffold.go`.
- `[ ]` **DD-367** Read `scaffold_command` from `template.yml`.
- `[ ]` **DD-368** Check prerequisite command exists before scaffold.
- `[ ]` **DD-369** Check `pnpm` before Next.js scaffold.
- `[ ]` **DD-370** Check `node` before Next.js scaffold.
- `[ ]` **DD-371** Check `composer` before Laravel scaffold.
- `[ ]` **DD-372** Check `php` before Laravel scaffold.
- `[ ]` **DD-373** Run `pnpm create next-app .` in temp directory.
- `[ ]` **DD-374** Run `composer create-project laravel/laravel .` in temp directory.
- `[ ]` **DD-375** Stream scaffold command output.
- `[ ]` **DD-376** Capture scaffold errors.
- `[ ]` **DD-377** Delete temp directory on scaffold failure.

## Atomic Project Creation

- `[ ]` **DD-378** Validate final project directory does not exist.
- `[ ]` **DD-379** Create temp directory outside final path.
- `[ ]` **DD-380** Run scaffold in temp directory.
- `[ ]` **DD-381** Render DevDock files into temp directory.
- `[ ]` **DD-382** Generate `compose.yml` from rendered `.devdock.yml`.
- `[ ]` **DD-383** Update `.gitignore`.
- `[ ]` **DD-384** Move temp directory to final project path only after all steps succeed.
- `[ ]` **DD-385** Ensure no partial final project directory is left on failure.
- `[ ]` **DD-386** Clean up temp directory on any error.

## Post-Scaffold Hooks

- `[ ]` **DD-387** Add Laravel `post_scaffold.sh`.
- `[ ]` **DD-388** Laravel hook runs `php artisan key:generate` or equivalent.
- `[ ]` **DD-389** Add Next.js `post_scaffold.sh`.
- `[ ]` **DD-390** Next.js hook runs `pnpm add prisma @prisma/client`.
- `[ ]` **DD-391** Next.js hook runs `pnpm prisma generate`.
- `[ ]` **DD-392** Stream hook output.
- `[ ]` **DD-393** Stop create flow on hook failure.
- `[ ]` **DD-394** Delete temp directory on hook failure.
- `[ ]` **DD-395** Do not add hook confirmation UI in v0.1.

## `devdock create` Command

- `[ ]` **DD-396** Create `internal/cli/create.go`.
- `[ ]` **DD-397** Implement `devdock create` command.
- `[ ]` **DD-398** Support non-interactive `devdock create next-postgres my-saas`.
- `[ ]` **DD-399** Support non-interactive `devdock create laravel-api my-api`.
- `[ ]` **DD-400** Support interactive template picker.
- `[ ]` **DD-401** Ask project name in interactive mode.
- `[ ]` **DD-402** Ask Redis option in interactive mode.
- `[ ]` **DD-403** Print progress for each create step.
- `[ ]` **DD-404** Print final next steps.
- `[ ]` **DD-405** Print final app URL.

## Optional Flags

- `[ ]` **DD-406** Decide whether `--no-install` is included in v0.1 or deferred.
- `[ ]` **DD-407** If included, implement `--no-install`.
- `[ ]` **DD-408** If deferred, return friendly “not available in v0.1” message.

## Template Tests

- `[ ]` **DD-409** Add renderer tests for all variables.
- `[ ]` **DD-410** Add template metadata parse tests.
- `[ ]` **DD-411** Add smoke test for Next.js template file rendering.
- `[ ]` **DD-412** Add smoke test for Laravel template file rendering.

## Manual Create Tests

- `[ ]` **DD-413** Run `devdock create next-postgres my-test`.
- `[ ]` **DD-414** Run `cd my-test && devdock up`.
- `[ ]` **DD-415** Confirm PostgreSQL healthy.
- `[ ]` **DD-416** Confirm app reachable at `http://localhost:3000`.
- `[ ]` **DD-417** Run `devdock down`.
- `[ ]` **DD-418** Run `devdock create laravel-api my-api`.
- `[ ]` **DD-419** Run `cd my-api && devdock up`.
- `[ ]` **DD-420** Confirm MySQL healthy.
- `[ ]` **DD-421** Confirm app reachable at `http://localhost:8000`.
- `[ ]` **DD-422** Run `devdock down`.

## Week 6 Done When

- `[ ]` `devdock create next-postgres my-test && cd my-test && devdock up && curl localhost:3000` succeeds.
- `[ ]` `devdock create laravel-api my-api && cd my-api && devdock up && curl localhost:8000` succeeds.
- `[ ]` Failed scaffold leaves no partial final directory.
- `[ ]` Template renderer tests pass.
- `[ ]` Template smoke tests pass.

---

# Week 7 — Diagnostics, Error Polish & Clean Machine Testing

**Goal:** `devdock doctor` is complete, all errors are actionable, and the two core flows work on another clean macOS machine.

## Full Doctor Checks

- `[ ]` **DD-423** Complete Docker installed check.
- `[ ]` **DD-424** Complete Docker daemon running check.
- `[ ]` **DD-425** Complete Docker Compose v2 check.
- `[ ]` **DD-426** Complete `.devdock.yml` exists check.
- `[ ]` **DD-427** Complete `.devdock.yml` valid YAML check.
- `[ ]` **DD-428** Complete `.devdock.yml` schema validation check.
- `[ ]` **DD-429** Complete Node version check.
- `[ ]` **DD-430** Complete PHP version check.
- `[ ]` **DD-431** Complete port availability check.
- `[ ]` **DD-432** Complete `.env` exists check.
- `[ ]` **DD-433** Skip irrelevant runtime checks based on project type.
- `[ ]` **DD-434** Skip app-specific checks for `docker-compose` project type.
- `[ ]` **DD-435** Ensure doctor is read-only in v0.1.
- `[ ]` **DD-436** Ensure doctor has no auto-fix behavior in v0.1.

## Doctor Output

- `[ ]` **DD-437** Print doctor title with project name.
- `[ ]` **DD-438** Print pass checks with `✔`.
- `[ ]` **DD-439** Print failed checks with `✗`.
- `[ ]` **DD-440** Include concrete fix for every failed check.
- `[ ]` **DD-441** Print issue count summary.
- `[ ]` **DD-442** Exit code `0` if all checks pass.
- `[ ]` **DD-443** Exit code `1` if any check fails.

## Error Message Review

- `[ ]` **DD-444** Review config errors.
- `[ ]` **DD-445** Review detector errors.
- `[ ]` **DD-446** Review init errors.
- `[ ]` **DD-447** Review Compose generation errors.
- `[ ]` **DD-448** Review Docker errors.
- `[ ]` **DD-449** Review port conflict errors.
- `[ ]` **DD-450** Review process runner errors.
- `[ ]` **DD-451** Review logs errors.
- `[ ]` **DD-452** Review status errors.
- `[ ]` **DD-453** Review create/scaffold errors.
- `[ ]` **DD-454** Verify every error follows what/why/fix format.
- `[ ]` **DD-455** Verify no raw stack traces are shown to user.
- `[ ]` **DD-456** Verify stack traces or debug details go to `~/.devdock/logs/error.log`.

## Common Failure Scenarios

- `[ ]` **DD-457** Test Docker not installed.
- `[ ]` **DD-458** Test Docker not running.
- `[ ]` **DD-459** Test Docker Compose v2 missing or invalid.
- `[ ]` **DD-460** Test `.devdock.yml` missing.
- `[ ]` **DD-461** Test invalid YAML.
- `[ ]` **DD-462** Test invalid schema.
- `[ ]` **DD-463** Test Node missing.
- `[ ]` **DD-464** Test PHP missing.
- `[ ]` **DD-465** Test port conflict.
- `[ ]` **DD-466** Test `.env` missing.
- `[ ]` **DD-467** Test scaffold failure.
- `[ ]` **DD-468** Test service health failure.

## Second Clean macOS Machine Test

- `[ ]` **DD-469** Install DevDock binary on second Apple Silicon Mac or fresh macOS user profile.
- `[ ]` **DD-470** Run `devdock doctor` before setup.
- `[ ]` **DD-471** Follow doctor instructions until prerequisites pass.
- `[ ]` **DD-472** Run existing Laravel flow.
- `[ ]` **DD-473** Run new Next.js project flow.
- `[ ]` **DD-474** Run Docker Compose project flow.
- `[ ]` **DD-475** Record every error encountered.
- `[ ]` **DD-476** Fix confusing errors.
- `[ ]` **DD-477** Re-run flows after fixes.

## Week 7 Done When

- `[ ]` `devdock doctor` covers all 9 v0.1 checks.
- `[ ]` `devdock doctor` is read-only.
- `[ ]` Every failed doctor check includes a concrete fix.
- `[ ]` Port conflict errors in `devdock up` match doctor format.
- `[ ]` No user-facing stack traces exist.
- `[ ]` A developer can resolve the 5 most common failures from error messages alone.
- `[ ]` Two core flows pass on a second clean macOS machine or fresh user profile.

---

# Week 8 — Integration Testing & Release

**Goal:** v0.1.0 is packaged, documented, released, and validated by external testers.

## End-to-End Test Matrix

- `[ ]` **DD-478** Test existing Laravel project on Apple Silicon Mac.
- `[ ]` **DD-479** Test new Next.js project on Apple Silicon Mac.
- `[ ]` **DD-480** Test existing Docker Compose project on Apple Silicon Mac.
- `[ ]` **DD-481** Test existing Laravel project on Intel Mac.
- `[ ]` **DD-482** Test new Next.js project on Intel Mac.
- `[ ]` **DD-483** Test existing Docker Compose project on Intel Mac.
- `[ ]` **DD-484** Test `devdock down --volumes` prompt.
- `[ ]` **DD-485** Test running `devdock up` twice.
- `[ ]` **DD-486** Test Ctrl+C behavior in foreground mode.
- `[ ]` **DD-487** Test detached mode app PID cleanup.
- `[ ]` **DD-488** Test stale PID cleanup.
- `[ ]` **DD-489** Test no orphan containers after `devdock down`.

## Build System

- `[ ]` **DD-490** Create `Makefile`.
- `[ ]` **DD-491** Add `make test`.
- `[ ]` **DD-492** Add `make build`.
- `[ ]` **DD-493** Add `make build-all`.
- `[ ]` **DD-494** Build `devdock-darwin-arm64`.
- `[ ]` **DD-495** Build `devdock-darwin-amd64`.
- `[ ]` **DD-496** Inject version into binaries.
- `[ ]` **DD-497** Verify binary runs on Apple Silicon.
- `[ ]` **DD-498** Verify binary runs on Intel Mac.

## Signing & Release Packaging

- `[ ]` **DD-499** Decide whether v0.1 binary signing is required or deferred.
- `[ ]` **DD-500** If signing is included, sign release binaries.
- `[ ]` **DD-501** Generate checksums for release binaries.
- `[ ]` **DD-502** Create release archive or direct binaries.
- `[ ]` **DD-503** Write `install.sh`.
- `[ ]` **DD-504** `install.sh` detects Apple Silicon vs Intel.
- `[ ]` **DD-505** `install.sh` downloads correct binary.
- `[ ]` **DD-506** `install.sh` installs to `/usr/local/bin/devdock` or explains permission issue.
- `[ ]` **DD-507** Test install script on fresh machine.

## README

- `[ ]` **DD-508** Write project summary.
- `[ ]` **DD-509** Write installation instructions.
- `[ ]` **DD-510** Write prerequisite section.
- `[ ]` **DD-511** Write quick start for existing Laravel project.
- `[ ]` **DD-512** Write quick start for new Next.js project.
- `[ ]` **DD-513** Write Docker Compose project usage section.
- `[ ]` **DD-514** Write command reference for v0.1 commands.
- `[ ]` **DD-515** Write troubleshooting section.
- `[ ]` **DD-516** Mention explicitly what is not in v0.1.
- `[ ]` **DD-517** Add screenshots or terminal output examples if available.

## GitHub Release

- `[ ]` **DD-518** Tag `v0.1.0`.
- `[ ]` **DD-519** Create GitHub Release.
- `[ ]` **DD-520** Upload Apple Silicon binary.
- `[ ]` **DD-521** Upload Intel binary.
- `[ ]` **DD-522** Upload checksums.
- `[ ]` **DD-523** Upload `install.sh`.
- `[ ]` **DD-524** Add release notes.
- `[ ]` **DD-525** Include known limitations.

## External Tester Validation

- `[ ]` **DD-526** Recruit 3 external developers.
- `[ ]` **DD-527** Ask tester 1 to run existing Laravel flow.
- `[ ]` **DD-528** Ask tester 2 to run new Next.js flow.
- `[ ]` **DD-529** Ask tester 3 to run both flows.
- `[ ]` **DD-530** Collect setup time from each tester.
- `[ ]` **DD-531** Collect error/confusion points.
- `[ ]` **DD-532** Fix critical blockers.
- `[ ]` **DD-533** Record non-critical issues for v0.1.1 or v0.2.

## Community Launch

- `[ ]` **DD-534** Prepare short launch post.
- `[ ]` **DD-535** Post on relevant developer communities.
- `[ ]` **DD-536** Post on IndieHackers.
- `[ ]` **DD-537** Prepare Hacker News post but avoid overhyping.
- `[ ]` **DD-538** Share in Laravel, Next.js, and Go/PHP developer groups where allowed.
- `[ ]` **DD-539** Create issue templates for bug reports.
- `[ ]` **DD-540** Create roadmap issue for v0.2.

## Week 8 Done When

- `[ ]` v0.1.0 is tagged.
- `[ ]` GitHub Release is live.
- `[ ]` Install script works on macOS.
- `[ ]` Apple Silicon binary works.
- `[ ]` Intel binary works.
- `[ ]` README explains both core flows clearly.
- `[ ]` At least 3 external developers complete the core flow on their own machine.
- `[ ]` Critical tester blockers are fixed or documented.

---

# Final v0.1 Definition of Done

v0.1.0 is complete only when all items below pass on a clean macOS machine.

## Flow 1 — Existing Laravel Project

```bash
git clone https://github.com/laravel/laravel my-laravel
cd my-laravel
devdock init
devdock up
curl localhost:8000
devdock down
```

- `[ ]` Laravel is detected.
- `[ ]` `.devdock.yml` is created.
- `[ ]` `compose.yml` is generated.
- `[ ]` MySQL starts and becomes healthy.
- `[ ]` Redis starts and becomes healthy if enabled.
- `[ ]` `php artisan serve` starts.
- `[ ]` `curl localhost:8000` returns HTTP 200.
- `[ ]` `devdock down` stops services.
- `[ ]` No orphan containers remain.

## Flow 2 — New Next.js Project

```bash
devdock create next-postgres my-saas
cd my-saas
devdock up
curl localhost:3000
devdock down
```

- `[ ]` Project directory is created atomically.
- `[ ]` Next.js scaffold is generated.
- `[ ]` DevDock files are generated.
- `[ ]` PostgreSQL starts and becomes healthy.
- `[ ]` `pnpm dev` starts.
- `[ ]` `curl localhost:3000` returns HTTP 200.
- `[ ]` `devdock down` stops services.
- `[ ]` No orphan containers remain.

## Flow 3 — Existing Docker Compose Project

```bash
cd existing-compose-project
devdock init
devdock up
devdock status
devdock logs
devdock down
```

- `[ ]` Docker Compose project is detected.
- `[ ]` Minimal `.devdock.yml` is created.
- `[ ]` Existing `compose.yml` is not overwritten.
- `[ ]` `devdock up` proxies `docker compose up -d`.
- `[ ]` `devdock status` formats `docker compose ps`.
- `[ ]` `devdock logs` proxies `docker compose logs -f`.
- `[ ]` `devdock down` proxies `docker compose down`.

## Quality Gates

- `[ ]` `devdock doctor` passes all checks on a correctly configured machine.
- `[ ]` `devdock doctor` fails correctly when Docker is not running.
- `[ ]` Running `devdock up` twice does not create duplicate app processes.
- `[ ]` Running `devdock up` twice does not create orphan containers.
- `[ ]` `devdock down --volumes` prompts before deleting data.
- `[ ]` Every error message follows the what/why/fix format.
- `[ ]` Port conflict errors name the exact `.devdock.yml` field to edit.
- `[ ]` No silent file overwrites occur in any flow.
- `[ ]` Ctrl+C during `devdock up` stops only the app process.
- `[ ]` Docker services remain accessible after Ctrl+C.
- `[ ]` `compose.yml` header clearly warns that manual edits may be overwritten.
- `[ ]` Failed `devdock create` leaves no partial final project directory.
- `[ ]` User-facing output never shows raw stack traces.

---

# Explicitly Deferred to v0.2+

Do not implement these in v0.1 unless all Definition of Done items are already complete.

## v0.2

- `[ ]` `devdock service add/remove`
- `[ ]` Mailpit service
- `[ ]` MinIO service
- `[ ]` MongoDB service
- `[ ]` Express stack
- `[ ]` NestJS stack
- `[ ]` Go Fiber stack
- `[ ]` FastAPI stack
- `[ ]` `devdock run <command>`
- `[ ]` `devdock open`
- `[ ]` `devdock doctor --fix`
- `[ ]` `devdock restart [service]`

## v0.3

- `[ ]` Template registry / remote fetch
- `[ ]` Self-update: `devdock update`
- `[ ]` Telemetry
- `[ ]` Homebrew tap
- `[ ]` Container app run mode

## v1.0+

- `[ ]` Config migration: `devdock migrate-config`
- `[ ]` SQLite state database
- `[ ]` `.test` domain support
- `[ ]` Desktop GUI
- `[ ]` Windows support
- `[ ]` Linux support

---

# Recommended First Coding Checkpoints

## Checkpoint 1 — CLI Exists

```bash
devdock --version
devdock --help
devdock doctor
```

- `[ ]` Version prints.
- `[ ]` Help prints v0.1 commands only.
- `[ ]` Doctor runs without panic.

## Checkpoint 2 — Config Works

```bash
devdock doctor
```

- `[ ]` Valid `.devdock.yml` parses.
- `[ ]` Invalid `.devdock.yml` gives helpful schema error.

## Checkpoint 3 — Init Works

```bash
cd sample-next-app
devdock init
```

- `[ ]` Next.js is detected.
- `[ ]` `.devdock.yml` is generated.
- `[ ]` `compose.yml` is generated.

## Checkpoint 4 — Lifecycle Works

```bash
cd sample-next-app
devdock up
curl localhost:3000
devdock down
```

- `[ ]` PostgreSQL starts.
- `[ ]` App starts.
- `[ ]` App responds HTTP 200.
- `[ ]` Services stop cleanly.

## Checkpoint 5 — Create Works

```bash
devdock create next-postgres my-test
cd my-test
devdock up
curl localhost:3000
devdock down
```

- `[ ]` Project scaffolds.
- `[ ]` DevDock files are generated.
- `[ ]` Full flow succeeds.
