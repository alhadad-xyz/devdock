# DevDock v0.3.0 Product Requirements Document

**Product:** DevDock
**Version:** v0.3.0
**Release Theme:** Distribution & Template Infrastructure
**Type:** Product Requirements Document
**Status:** Approved for Engineering Planning
**Base Version:** DevDock v0.2.0 complete
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
7. [Template Registry Specification](#7-template-registry-specification)
8. [Template Security Model](#8-template-security-model)
9. [Update System Specification](#9-update-system-specification)
10. [Homebrew Distribution Specification](#10-homebrew-distribution-specification)
11. [Telemetry & Privacy](#11-telemetry--privacy)
12. [`devdock config` Command Specification](#12-devdock-config-command-specification)
13. [Configuration Schema Changes](#13-configuration-schema-changes)
14. [Backward Compatibility Contract](#14-backward-compatibility-contract)
15. [Non-Functional Requirements](#15-non-functional-requirements)
16. [6-Week Build Plan](#16-6-week-build-plan)
17. [Definition of Done](#17-definition-of-done)
18. [Risks & Mitigations](#18-risks--mitigations)

---

## 1. Executive Summary

DevDock v0.1 proved the core local development lifecycle. DevDock v0.2 expanded DevDock into a daily workflow tool with `devdock open`, `devdock run`, service management, Mailpit, MinIO, Express, and Go Fiber.

DevDock v0.3 shifts focus entirely away from local workflow features and toward **distribution, self-update, and template infrastructure**. No new stacks or services ship in v0.3.

v0.3 delivers:

- **Homebrew installation** — the primary install method for all new users
- **Self-update command** — keep DevDock current without manual binary management
- **GitHub Actions release automation** — consistent, checksum-verified releases
- **Template registry** — official remote templates, no longer bundled inside the binary
- **Template commands** — list, search, inspect, and cache templates
- **Template security** — checksum verification, hook review, path traversal protection
- **Offline mode** — cached templates work without internet
- **Opt-in telemetry** — privacy-safe, anonymous usage data with a hard opt-in default

The v0.3 core promise:

```bash
brew install alhadad-xyz/tap/devdock
devdock template list
devdock create next-postgres my-saas
devdock update --check
```

---

## 2. Product Goal

> Make DevDock easy to install, keep up to date, and extend through trusted remote templates.

| Version | Question Answered |
|---|---|
| v0.1 | Can DevDock start projects? |
| v0.2 | Can DevDock help during everyday development? |
| v0.3 | Can DevDock become easy to adopt and grow beyond bundled templates? |

---

## 3. Scope

### 3.1 In Scope

| Category | v0.3 Contents |
|---|---|
| Platform | macOS only |
| Distribution | Homebrew tap, GitHub Releases automation, install script polish |
| Update system | `devdock update`, `devdock update --check` |
| Template registry | Static official manifest hosted over HTTPS |
| Template commands | `template list`, `template search`, `template info`, `template update` |
| Template create | `devdock create <template> <project>` from remote registry |
| Template cache | `~/.devdock/templates/` with TTL-based manifest cache |
| Offline mode | `devdock create --offline` uses cached templates only |
| Template security | SHA-256 checksum, path traversal protection, hook review, `--no-hooks` |
| Telemetry | Opt-in only, privacy-safe, with `devdock config` controls |
| `devdock config` | New command group for global config management |
| Release automation | GitHub Actions build, checksums, release artifacts |

### 3.2 Explicitly Out of Scope

| Feature | Deferred To |
|---|---|
| New stacks (NestJS, FastAPI, Django, Rails) | v0.4 (registry templates only) |
| New services | v0.4 |
| Desktop GUI | v2.0 |
| Windows/Linux support | v2.0 |
| `.test` domain support | v1.1 |
| Community template publishing | v0.4+ |
| Template signing with asymmetric keys | v1.0+ |
| Team accounts or paid templates | v2.x |
| Plugin system | Not planned |
| AI environment fixer | v3.x |

### 3.3 Scope Discipline

v0.3 must not add framework or service features. If a new stack or service seems tempting, it belongs as a registry template in v0.4, not as compiled-in code in v0.3.

---

## 4. Problems to Solve

### 4.1 Installation Requires Go Knowledge

v0.1 and v0.2 require cloning the repo and building from source. This blocks adoption from developers who don't have Go installed.

**Solution:** `brew install alhadad-xyz/tap/devdock` — no Go, no source, no `~/go/bin`.

### 4.2 Updating Is Manual

Users currently re-clone or download binaries manually.

**Solution:** `devdock update` — detects architecture, downloads the correct binary, verifies the checksum, replaces atomically.

### 4.3 Bundled Templates Don't Scale

Every template update requires a new DevDock release. The binary grows with every template added. Template improvements can't ship independently.

**Solution:** A static official registry at `registry.devdock.dev`. Templates are downloaded on demand, cached locally, and usable offline after first fetch.

### 4.4 Remote Templates Require Trust Controls

Remote templates run scaffold commands and hooks. Without controls, a compromised or malicious template could execute arbitrary code.

**Solution:** SHA-256 checksum verification for every downloaded archive, path traversal protection during extraction, hook review prompt before any execution, and `--no-hooks` to skip hooks entirely. Official registry only in v0.3 — no community publishing yet.

### 4.5 Offline Work Must Keep Working

Developers with unreliable internet or air-gapped environments need templates to work after first download.

**Solution:** Template archives are cached to `~/.devdock/templates/`. `--offline` forces cache-only resolution.

### 4.6 Product Data Is Needed, But Privacy Matters

To improve DevDock, aggregate usage data on which commands, templates, and services are most used would be valuable. But developer tools must never collect project data, paths, or secrets.

**Solution:** Opt-in telemetry (default off), anonymous events only, no session ID, configurable via `devdock config`.

---

## 5. User Stories & Acceptance Criteria

### US-01: Install via Homebrew

**Acceptance Criteria:**
- `brew install alhadad-xyz/tap/devdock` installs the latest stable release on both Apple Silicon and Intel
- `devdock --version` works immediately after installation with no PATH changes required
- README lists Homebrew as the primary install method
- Source install (`go install`) remains documented as an alternative

---

### US-02: Check for Updates

**Acceptance Criteria:**
- `devdock update --check` fetches latest GitHub Release metadata
- Up to date: prints `DevDock v0.3.0 is up to date.`
- Behind: prints current version, latest version, and `Run 'devdock update' to upgrade.`
- Network unreachable: prints error with what/why/fix format, exits 2
- Never downloads or modifies anything when `--check` is used

---

### US-03: Self-Update

**Acceptance Criteria:**
- `devdock update` downloads the binary matching current OS and architecture
- Downloads `checksums.txt`, verifies SHA-256 before touching the current binary
- If Homebrew-managed: prints `brew upgrade devdock` and exits 0 without modifying any files
- Atomic replacement: old binary is backed up, new binary is moved in place; backup restored on failure
- If the download or checksum fails, the original binary is never modified

---

### US-04: List Templates

**Acceptance Criteria:**
- `devdock template list` fetches the official manifest and prints ID, version, category, description
- If cached manifest is valid (not stale): uses cache without network call
- If offline and manifest is cached (even if stale): uses cache, prints a notice that it may be outdated
- If offline and no manifest cached: prints a clear error with `devdock template update` as the fix
- Never downloads template archives during list

---

### US-05: Search Templates

**Acceptance Criteria:**
- `devdock template search <query>` filters by ID, name, description, category, tags, runtime, and services
- Case-insensitive
- Zero results: `No templates found for '<query>'.`
- Works with cached manifest if offline

---

### US-06: Inspect a Template

**Acceptance Criteria:**
- `devdock template info <id>` shows all metadata without executing anything: ID, name, version, category, description, runtime requirements, services, scaffold command, hooks, source (official/local), cached status, checksum status
- Unknown template ID: helpful error with `devdock template list` as the fix

---

### US-07: Create Project from Registry Template

**Acceptance Criteria:**
- `devdock create next-postgres my-saas` fetches from registry if not cached, verifies checksum, extracts, and creates the project
- Scaffold command and hooks are printed and confirmed before execution
- `--no-hooks` skips post-create hooks and the hook confirmation prompt; scaffold still runs
- Project creation remains atomic: failure leaves no partial final directory

---

### US-08: Offline Template Create

**Acceptance Criteria:**
- `devdock create next-postgres my-app --offline` uses cache only — no network calls
- If template is cached and checksum-verified: succeeds
- If template not cached: error with instructions to run `devdock template update` while online

---

### US-09: Refresh Template Cache

**Acceptance Criteria:**
- `devdock template update` fetches the latest manifest and prints new/updated/unchanged summary
- `devdock template update --all` downloads and verifies all official templates into cache
- Old cached versions are not deleted automatically
- Network failure: error with fix; cached manifest is not corrupted

---

### US-10: Opt In/Out of Telemetry

**Acceptance Criteria:**
- Telemetry is disabled by default; first-run prompt has `N` as the default
- User controls: `devdock config set telemetry true/false`, `devdock config get telemetry`
- `DEVDOCK_NO_TELEMETRY=1` disables telemetry regardless of config
- README documents what is collected and what is never collected
- No project names, paths, file contents, `.env` values, secrets, or arguments that could contain paths are ever sent

---

## 6. Functional Requirements

### 6.1 Homebrew Installation

**Tap repository:** `github.com/alhadad-xyz/homebrew-tap`
**Formula path:** `Formula/devdock.rb`

The formula:
- Downloads the official GitHub Release binary for the correct architecture
- Verifies SHA-256 before installing
- Installs as `devdock`
- Formula test: `system "#{bin}/devdock", "--version"`

Installation:
```bash
brew install alhadad-xyz/tap/devdock
```

Source alternative (remains documented):
```bash
go install github.com/alhadad-xyz/devdock/cmd/devdock@latest
```

---

### 6.2 GitHub Actions Release Automation

The release workflow triggers on Git tag `v*.*.*` and must:

1. Run `go test ./...` — fail the release if tests fail
2. Build `darwin/arm64` with `-ldflags "-X main.Version=${TAG} -X main.Commit=${SHA} -X main.BuildDate=${DATE}"`
3. Build `darwin/amd64` with the same flags
4. Generate SHA-256 checksums: `shasum -a 256 devdock-darwin-arm64 devdock-darwin-amd64 > checksums.txt`
5. Attach all three files to the GitHub Release
6. (Optional, post-release) Update Homebrew formula checksums via a follow-up commit to the tap repo

**Binary naming:**
```
devdock-darwin-arm64
devdock-darwin-amd64
checksums.txt
```

---

### 6.3 `devdock update`

#### Behavior: `--check`

1. Read `Version` from build-time variable
2. Fetch `https://api.github.com/repos/alhadad-xyz/devdock/releases/latest`
3. Parse `tag_name`, strip leading `v`
4. Compare using semver
5. Print status
6. Exit 0 if up to date, 1 if update available, 2 on network/API error

Output when behind:
```
DevDock v0.2.0 is installed. v0.3.0 is available.

Run `devdock update` to upgrade.
```

#### Behavior: `devdock update`

1. Detect `runtime.GOOS` and `runtime.GOARCH`
2. Fetch latest release metadata
3. Resolve binary asset name (`devdock-darwin-arm64` or `devdock-darwin-amd64`)
4. Detect Homebrew management (see Section 9.3)
5. If Homebrew-managed: print message and exit 0 — do not proceed
6. Download binary to `~/.devdock/tmp/devdock-new`
7. Download `checksums.txt`
8. Verify SHA-256 against `checksums.txt` entry for the binary name
9. If mismatch: delete temp file, print error with checksum values, exit 1
10. Backup current binary: `os.Rename(currentPath, currentPath+".backup")`
11. Move new binary: `os.Rename(tempPath, currentPath)`
12. Set executable permission: `os.Chmod(currentPath, 0755)`
13. If step 11–12 fails: restore backup, exit 1
14. Print: `Updated to DevDock v0.3.0.`
15. Delete backup

**Timeout:** 10 seconds for metadata fetch; 120 seconds for binary download.

---

### 6.4 `devdock template list`

Output format:
```
Official Templates  (cached 2h ago)

ID                  Version   Category    Description
───────────────────────────────────────────────────────────────
next-postgres       0.3.0     fullstack   Next.js with PostgreSQL and Prisma
laravel-api         0.3.0     backend     Laravel API with MySQL
express-postgres    0.3.0     backend     Express API with PostgreSQL
go-fiber-api        0.3.0     backend     Go Fiber API with PostgreSQL
```

When using stale cache (offline with expired manifest):
```
Official Templates  (cached manifest — may be outdated)
```

---

### 6.5 `devdock template search <query>`

Searches: ID, name, description, category, tags, runtime, services. Case-insensitive substring match.

No results:
```
No templates found for 'django'.

Run `devdock template list` to see all available templates.
```

---

### 6.6 `devdock template info <id>`

Output format:
```
Template: next-postgres
Name:     Next.js + PostgreSQL
Version:  0.3.0
Source:   official
Cached:   yes (verified)

Description:
  Next.js app with PostgreSQL and Prisma.

Runtime:
  node: 22

Services:
  postgres (required)
  redis    (optional)

Scaffold command:
  pnpm create next-app .

Post-create hooks:
  pnpm add prisma @prisma/client
  pnpm prisma generate
```

If template is not cached: `Cached: no`. If checksum unverified: `Cached: yes (unverified — run 'devdock template update' to refresh)`.

---

### 6.7 `devdock template update`

Behavior:
1. Fetch manifest from registry
2. Save to `~/.devdock/templates/manifest.json`
3. Update `~/.devdock/templates/cache.json` manifest TTL
4. Compare each template's cached version vs manifest version
5. Print summary

Output format:
```
Template registry updated.

  ✔ next-postgres     0.3.0  (cached, up to date)
  ↑ laravel-api       0.2.0 → 0.3.0  (new version available — run 'devdock template update --all' to download)
  + go-fiber-api      0.3.0  (new)
  = express-postgres  0.3.0  (no change)
```

`devdock template update --all`: after manifest update, downloads and verifies every template archive.

---

### 6.8 `devdock create` Template Resolution Order

When resolving a template ID or path:

1. **Local path** — if the argument is a file path (starts with `/`, `./`, or `../`)
2. **Cached registry template** — if ID matches a verified cached template
3. **Bundled template** — if ID matches a template compiled into the binary
4. **Remote registry template** — fetch, verify, cache, then use
5. **Error** — template not found anywhere

**Why cached before bundled:** A cached template may be a newer version than the bundled one. The registry is the source of truth; the bundle is a fallback for offline scenarios.

**With `--offline`:** Steps 4 (remote) is skipped. Error if not in cache or bundle.

**With `--no-hooks`:** Post-create hooks are skipped. The scaffold command still runs. The hook confirmation prompt is also skipped — `--no-hooks` implies user consent to run the scaffold without review. A warning is printed: "Post-create hooks were skipped. You may need to run them manually: `<list of skipped commands>`."

---

## 7. Template Registry Specification

### 7.1 Registry Type

Static official registry. HTTPS-only. No authentication. No community publishing in v0.3.

**Default URL:** `https://registry.devdock.dev/manifest.json`

Override priority (highest to lowest):
1. `DEVDOCK_REGISTRY_URL` environment variable
2. `registry.url` in `~/.devdock/config.yml`
3. Default URL

The env var name is `DEVDOCK_REGISTRY_URL` — this exact spelling must be used in all documentation and code.

---

### 7.2 Manifest Schema

```json
{
  "schema_version": "1",
  "updated_at": "2026-06-11T00:00:00Z",
  "templates": [
    {
      "id": "next-postgres",
      "name": "Next.js + PostgreSQL",
      "version": "0.3.0",
      "category": "fullstack",
      "description": "Next.js app with PostgreSQL and Prisma.",
      "tags": ["nextjs", "postgres", "prisma"],
      "runtime": {
        "node": "22"
      },
      "services": ["postgres", "redis"],
      "archive_url": "https://registry.devdock.dev/templates/next-postgres-0.3.0.tar.gz",
      "checksum_sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "official": true
    }
  ]
}
```

**`checksum_sha256` format:** lowercase hex-encoded SHA-256 digest, exactly 64 characters. Any other format is a validation error. An empty string or the string `"sha256-value-here"` is also a validation error — templates with placeholder checksums must not be used.

**Validation rules:**
- `schema_version` must be `"1"` — any other value rejects the manifest
- Each template must have a non-empty `checksum_sha256` of exactly 64 hex characters
- Duplicate `id` values within one manifest reject the manifest

---

### 7.3 Template Archive Structure

Remote templates use the same overlay structure as bundled templates:

```
template/
├── template.yml
├── overlay/
│   ├── .devdock.yml.tpl
│   ├── .env.tpl
│   └── .env.example.tpl
└── hooks/
    └── post_scaffold.sh
```

---

### 7.4 Manifest Cache TTL

The manifest is cached at `~/.devdock/templates/manifest.json`. The cache is valid for **24 hours** from the last successful fetch.

| Condition | Behavior |
|---|---|
| Cache age < 24h, online | Use cache; no network call |
| Cache age < 24h, offline | Use cache; no network call |
| Cache age ≥ 24h, online | Refetch manifest; update cache |
| Cache age ≥ 24h, offline | Use stale cache; print notice "Manifest may be outdated" |
| No cache, online | Fetch manifest; create cache |
| No cache, offline | Error: "No cached manifest. Run `devdock template update` while online." |

The `cached_at` timestamp is stored in `cache.json`.

---

### 7.5 Cache Layout

```
~/.devdock/
└── templates/
    ├── manifest.json          # cached registry manifest
    ├── cache.json             # cache metadata (TTL, versions, checksums)
    └── archives/
        ├── next-postgres/
        │   └── 0.3.0/
        │       ├── template.tar.gz
        │       ├── checksum.txt       # single-line: <sha256hex>  <filename>
        │       └── extracted/         # unpacked template contents
        └── laravel-api/
            └── 0.3.0/
```

---

### 7.6 Cache Metadata (`cache.json`)

```json
{
  "manifest_cached_at": "2026-06-11T10:00:00Z",
  "templates": [
    {
      "id": "next-postgres",
      "version": "0.3.0",
      "cached_at": "2026-06-11T10:05:00Z",
      "checksum_sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "path": "/Users/arief/.devdock/templates/archives/next-postgres/0.3.0/extracted",
      "verified": true
    }
  ]
}
```

**Important:** `path` must be stored as an absolute path resolved at write time. Never store `~` — it will not be expanded by Go's `os` package when reading.

---

## 8. Template Security Model

### 8.1 Checksum Verification

Every downloaded template archive must have a matching SHA-256 checksum in the manifest.

| Condition | Action |
|---|---|
| Manifest has no checksum for template | Reject — print error |
| Checksum is not 64 hex chars | Reject — print error |
| Download succeeds, checksum matches | Extract archive, mark `verified: true` in cache |
| Download succeeds, checksum mismatch | Delete downloaded archive, print error with both expected and actual values, do not extract |
| Archive already in cache, verified | Use without re-downloading |
| Archive in cache, not verified | Re-download and re-verify |

### 8.2 Path Traversal Protection

Before extracting any `.tar.gz` archive, DevDock must validate every entry path:

- Resolve the entry path against the target extraction directory
- Reject the entire archive if any entry path resolves outside the target directory
- This prevents malicious archives from writing to `../../` paths

Rejection error:
```
✗ Template archive rejected — path traversal detected

  Archive entry '../../../etc/passwd' would escape the extraction directory.
  This template archive may be malicious.

  Fix: Report this issue at https://github.com/alhadad-xyz/devdock/issues
```

### 8.3 Hook Review Prompt

Before running scaffold or post-create hooks for **remote or bundled** templates, DevDock prints all commands and asks for confirmation.

```
This template will run the following commands:

Scaffold:
  pnpm create next-app .

Post-create hooks:
  pnpm add prisma @prisma/client
  pnpm prisma generate

Source: official (registry.devdock.dev)

Run these commands? (Y/n)
```

**Bundled templates:** Hook review IS shown for bundled templates too. The rationale: a user may not know what a bundled template executes; transparency is the right default regardless of source trust level.

If user says no:
- No files are written to the final project directory
- Temp directory is cleaned up
- Exit 0 with message: "Create cancelled. No files were written."

**`--no-hooks` behavior:**
- Skips post-create hooks entirely
- Also skips the hook review confirmation prompt (the flag itself is the user's consent to skip review)
- Scaffold command still runs (required to produce the project structure)
- Prints: "Post-create hooks were skipped. You may need to run the following manually:" followed by the list of skipped hooks

**Non-interactive / CI mode:** If stdin is not a TTY (e.g., CI pipeline), DevDock automatically applies `--no-hooks` behavior and prints a warning. Hooks requiring interactive confirmation are never blocked in CI.

### 8.4 Local Template Warning

When the template argument is a file path:

```
⚠ Creating from a local template path.

  DevDock cannot verify this template against a registry checksum.
  Path: /home/user/my-templates/custom-api

  Review template.yml and any hooks before continuing.

Continue? (y/N)
```

Default is N. In non-interactive mode (no TTY): refuses to proceed unless `--force-local` flag is passed.

### 8.5 Trust Model Summary

| Source | Checksum Required | Hook Review Shown | Path Traversal Check |
|---|---|---|---|
| Bundled | No (shipped in binary) | Yes | N/A (pre-extracted) |
| Official registry | Yes (mandatory) | Yes | Yes |
| Local path | No | Yes + warning | Yes |
| Community registry | Not supported in v0.3 | — | — |

---

## 9. Update System Specification

### 9.1 Version Build Variables

```go
var (
    Version   = "dev"
    Commit    = "unknown"
    BuildDate = "unknown"
)
```

Release workflow injects these using `-ldflags`:
```
-X 'main.Version=0.3.0'
-X 'main.Commit=abc1234'
-X 'main.BuildDate=2026-06-11'
```

Local dev builds use the defaults (`dev`, `unknown`, `unknown`). `devdock --version` output:
- Release: `devdock version 0.3.0 (commit: abc1234, built: 2026-06-11)`
- Dev build: `devdock version dev`

---

### 9.2 Version Comparison

Use semantic versioning (semver). A Go semver library (`golang.org/x/mod/semver`) is acceptable.

Pre-release version ordering: `0.3.0-beta < 0.3.0 < 0.3.1 < 0.4.0`

---

### 9.3 Homebrew Detection

**Authoritative method:** Run `brew --prefix devdock 2>/dev/null`. If this returns a non-empty path and the current binary path starts with that prefix, the binary is Homebrew-managed.

Fallback for machines without `brew`: check if the binary path starts with `/opt/homebrew/` (Apple Silicon) or `/usr/local/Cellar/` (Intel).

Do not treat `/usr/local/bin/` alone as a Homebrew indicator — it is also used for manual installs on Intel Macs.

When Homebrew-managed:
```
DevDock appears to be installed via Homebrew.

  Fix: Run:
    brew upgrade devdock
```

Exit 0 — this is not an error.

---

## 10. Homebrew Distribution Specification

### 10.1 Tap Repository

```
github.com/alhadad-xyz/homebrew-tap
Formula/devdock.rb
```

### 10.2 Formula Requirements

- Uses GitHub Release asset URL (not the source archive)
- Uses architecture-specific URL and SHA-256
- Installs binary as `devdock`
- Formula `test` block: `system "#{bin}/devdock", "--version"`

### 10.3 Formula Update Process

After each GitHub Release:
1. Automated step in release workflow computes SHA-256 for both binaries
2. Updates `devdock.rb` with new version, URLs, and checksums
3. Commits and pushes to `homebrew-tap` repo
4. Manual verification: `brew upgrade devdock` on a test machine

---

## 11. Telemetry & Privacy

### 11.1 Default Behavior

Telemetry is disabled by default. On the first `devdock` command after installing v0.3.0 on a machine with no existing `~/.devdock/config.yml`, DevDock shows:

```
Help improve DevDock by sharing anonymous usage data? (y/N)

What's collected:
  command names, template IDs, service names, error codes,
  DevDock version, macOS architecture

What's never collected:
  project names, file paths, file contents, .env values,
  source code, secrets, shell arguments

Change anytime: devdock config set telemetry false
```

Default is `N`. If no input (non-interactive): defaults to N, no prompt shown.

The `prompted` field in config is set to `true` after this prompt runs, so it never shows again.

### 11.2 Telemetry Endpoint

Events are sent as HTTP POST to:
```
https://telemetry.devdock.dev/v1/events
```

Payload format:
```json
{
  "event": "command_used",
  "properties": {
    "command": "up",
    "devdock_version": "0.3.0",
    "os": "darwin",
    "arch": "arm64"
  }
}
```

**No session ID, no user ID, no device ID.** Each event is stateless and anonymous. There is no way to correlate multiple events from the same user or machine.

Failed sends are silently discarded — a telemetry failure must never fail the user's command, show an error, or log anything visible to the user. In `DEVDOCK_DEBUG=1` mode, failed sends may log a debug message.

### 11.3 Allowed Events

| Event | Properties Sent |
|---|---|
| `command_used` | `command` (name only, no arguments) |
| `template_created` | `template_id` |
| `service_added` | `service_name` |
| `doctor_failed` | `error_code` (numeric, no message text) |
| `up_failed` | `error_code` |

**Forbidden in all events:** project name, project path, any file path, file contents, `.env` values, secrets, command arguments that could contain paths or secrets.

### 11.4 Environment Controls

| Variable | Effect |
|---|---|
| `DEVDOCK_NO_TELEMETRY=1` | Force-disables telemetry regardless of config |
| `DEVDOCK_DEBUG=1` | Enables debug logging (does not affect telemetry on/off) |

---

## 12. `devdock config` Command Specification

`devdock config` is a new command group in v0.3 for reading and writing `~/.devdock/config.yml`.

### 12.1 Commands

```bash
devdock config get <key>
devdock config set <key> <value>
devdock config list
```

### 12.2 Supported Keys

| Key | Type | Values | Default |
|---|---|---|---|
| `telemetry` | bool | `true`, `false` | `false` |
| `defaults.package_manager` | string | `npm`, `pnpm`, `yarn` | `pnpm` |
| `defaults.editor` | string | any string | `code` |
| `registry.url` | string | valid HTTPS URL | `https://registry.devdock.dev/manifest.json` |

### 12.3 Behavior

`devdock config get <key>`:
- Prints the current value
- If key does not exist in config, prints the default value
- Unknown key: prints error with `devdock config list` as the fix

`devdock config set <key> <value>`:
- Validates the value for the key's type
- Writes to `~/.devdock/config.yml` atomically
- Prints: `Set <key> = <value>`
- Unknown key: error with list of valid keys
- Invalid value: error with valid options

`devdock config list`:
- Prints all supported keys, current values, and defaults

### 12.4 Error Format

```
✗ Unknown config key 'theme'

  Supported keys:
    telemetry               (bool)
    defaults.package_manager (string: npm|pnpm|yarn)
    defaults.editor         (string)
    registry.url            (string: HTTPS URL)

  Fix: Run `devdock config list` to see all supported keys.
```

---

## 13. Configuration Schema Changes

`~/.devdock/config.yml` gains two new top-level sections in v0.3. All new fields are optional with defaults.

```yaml
version: "1"

defaults:
  package_manager: pnpm
  editor: code

registry:
  url: "https://registry.devdock.dev/manifest.json"   # optional — defaults to official registry

telemetry:
  enabled: false
  prompted: false        # true after first-run prompt has been shown
```

**No breaking change:** A v0.1/v0.2 config with only `defaults` still loads without error. New sections are additive.

---

## 14. Backward Compatibility Contract

v0.3 must not break any v0.1 or v0.2 behavior. Required regression tests (automated):

| Test | Type |
|---|---|
| v0.1 Laravel `.devdock.yml` loads without error | Automated |
| v0.1 Next.js `.devdock.yml` loads without error | Automated |
| v0.2 Express `.devdock.yml` loads without error | Automated |
| v0.2 Go Fiber `.devdock.yml` loads without error | Automated |
| v0.2 Mailpit/MinIO config generates valid Compose | Automated |
| `devdock create <bundled>` works when registry unreachable | Manual |
| All v0.2 commands work with v0.2 `.devdock.yml` | Manual end-to-end |

---

## 15. Non-Functional Requirements

### 15.1 Performance

| Operation | Target |
|---|---|
| `devdock template list` (cached manifest) | < 300ms |
| `devdock template search` (cached) | < 300ms |
| `devdock template info` (cached) | < 300ms |
| `devdock update --check` | < 3s on normal network |
| Manifest fetch (online) | < 5s |
| Template archive download (100KB archive) | < 30s on normal network |
| SHA-256 verification (100KB archive) | < 1s |

### 15.2 Reliability

- Failed template download: no partial archive left in cache
- Failed self-update: original binary fully intact
- Failed telemetry send: silent discard, command proceeds normally
- Failed manifest fetch: cached manifest used if available
- `devdock create` atomicity: maintained from v0.1/v0.2

### 15.3 Security

- All remote template archives: SHA-256 verification before extraction
- All archive extractions: path traversal validation before any file write
- Homebrew-managed binary: never overwritten by `devdock update`
- Telemetry: opt-in default, no PII, no correlatable identifiers

---

## 16. 6-Week Build Plan

### Week 1 — Release Distribution

**Goal:** New users can install DevDock without cloning the repo or having Go installed.

Build: GitHub Actions release workflow, release checksums, Homebrew tap and formula, version metadata build variables, README installation update.

**Done when:** `brew install alhadad-xyz/tap/devdock && devdock --version` works on a clean macOS machine with no Go installed.

---

### Week 2 — Self-Update

**Goal:** Existing users can update DevDock from the CLI.

Build: `devdock update --check`, `devdock update` (download + checksum + atomic replace), Homebrew detection using `brew --prefix devdock`.

**Done when:** A manually installed v0.2 binary can update itself to v0.3. A Homebrew-installed binary refuses direct update and prints `brew upgrade devdock`.

---

### Week 3 — Template Registry Core

**Goal:** DevDock can fetch, verify, cache, and serve remote template archives.

Build: Registry manifest schema and Go structs, registry HTTP client (with TTL, timeout, env override), template cache layout and metadata, archive downloader (with atomic write), SHA-256 verifier with path traversal protection, cache TTL logic.

**Done when:** DevDock fetches the official manifest, downloads `next-postgres`, verifies the checksum, extracts safely, and successfully uses the extracted template from cache when the registry is blocked (offline simulation).

---

### Week 4 — Template Commands + Create Integration

**Goal:** Users can discover templates and create projects from registry templates.

Build: `devdock template list/search/info/update` commands, `devdock create` resolution order update (cached before bundled), `--offline` flag, `--no-hooks` flag (skips hooks and confirmation), hook review prompt.

**Done when:** `devdock template list`, `devdock template info next-postgres`, and `devdock create next-postgres my-test` all work using a live registry template on a clean machine.

---

### Week 5 — `devdock config`, Telemetry, and Security

**Goal:** Security model is complete, telemetry is opt-in and documented, config is manageable.

Build: `devdock config get/set/list` command, telemetry config fields, first-run prompt, telemetry event client (with endpoint), local template warning with TTY/non-TTY detection, README security and privacy sections.

**Done when:** `devdock create next-postgres test` shows hook review prompt and allows cancellation. `devdock config set telemetry true` enables telemetry. `DEVDOCK_NO_TELEMETRY=1` overrides config. Local template path shows warning with default N.

---

### Week 6 — Regression Testing and Release

**Goal:** v0.3.0 shipped, all previous functionality verified.

Build: Full regression suite (v0.1/v0.2 projects), registry integration test (online + offline), Homebrew install test, self-update test, README and changelog, GitHub Release.

**Done when:** v0.3.0 is tagged and published. At least 2 external testers complete: `brew install` → `template list` → `devdock create next-postgres my-test`.

---

## 17. Definition of Done

v0.3.0 is not released until every required item passes.

### Distribution
- [ ] GitHub Actions builds `darwin/arm64` binary with correct version string
- [ ] GitHub Actions builds `darwin/amd64` binary with correct version string
- [ ] `checksums.txt` attached to GitHub Release
- [ ] `brew install alhadad-xyz/tap/devdock` installs DevDock on Apple Silicon
- [ ] `brew install alhadad-xyz/tap/devdock` installs DevDock on Intel
- [ ] `devdock --version` shows correct version after Homebrew install
- [ ] README uses Homebrew as primary install method

### Update System
- [ ] `devdock update --check` exits 0 when up to date
- [ ] `devdock update --check` exits 1 when update available, prints both versions
- [ ] `devdock update --check` exits 2 on network failure with helpful error
- [ ] `devdock update` downloads correct binary for architecture
- [ ] Downloaded binary checksum is verified before any file modification
- [ ] Checksum mismatch: original binary unchanged, error printed
- [ ] Homebrew-managed binary: `brew upgrade devdock` printed, no files modified
- [ ] Successful update: new version confirmed by `devdock --version`

### Template Registry
- [ ] Official manifest fetched and cached with 24h TTL
- [ ] Stale cache used when offline with notice
- [ ] No cache + offline: clear error with fix
- [ ] `devdock template list` output matches spec format
- [ ] `devdock template search postgres` returns postgres-related templates
- [ ] `devdock template search xyznotfound` prints no-results message
- [ ] `devdock template info next-postgres` shows all fields including cached/verified status
- [ ] `devdock template update` prints new/updated/unchanged summary correctly
- [ ] `devdock template update --all` downloads and verifies all templates
- [ ] Old cached versions not deleted on update
- [ ] Template archives are SHA-256 verified before extraction
- [ ] Missing checksum rejects template
- [ ] Checksum mismatch rejects template and deletes bad archive
- [ ] Path traversal in archive: entire archive rejected, no files written
- [ ] Cached verified template used with `--offline`

### Create Integration
- [ ] Resolution order: cached → bundled → remote (verified with tests)
- [ ] `devdock create next-postgres my-test` works from registry template
- [ ] Bundled templates still work when registry unreachable
- [ ] `--offline` uses cache only; fails clearly if template not cached
- [ ] Hook review prompt shown before scaffold/hooks for all source types
- [ ] User can cancel at hook review; no files written to final directory
- [ ] `--no-hooks` skips hooks and prompt; prints list of skipped commands
- [ ] Non-TTY/CI: auto-applies `--no-hooks` behavior with warning
- [ ] Failed create leaves no partial final directory

### Security & Privacy
- [ ] Path traversal protection test: malicious archive rejected
- [ ] Local template path shows warning, default N, blocked in non-TTY without `--force-local`
- [ ] Telemetry disabled by default on fresh install
- [ ] First-run prompt default is `N`
- [ ] `devdock config set telemetry true` enables telemetry
- [ ] `devdock config set telemetry false` disables it
- [ ] `DEVDOCK_NO_TELEMETRY=1` disables regardless of config
- [ ] Telemetry failure does not fail or surface an error to the user
- [ ] README documents collected and forbidden telemetry data

### `devdock config`
- [ ] `devdock config get telemetry` prints current value
- [ ] `devdock config set telemetry false` writes to config file
- [ ] `devdock config set <unknown>` prints error with valid keys
- [ ] `devdock config list` shows all keys and values
- [ ] Config writes are atomic

### Backward Compatibility
- [ ] v0.1 Laravel project: `devdock up && devdock down` passes
- [ ] v0.1 Next.js project: `devdock up && devdock down` passes
- [ ] v0.2 Express project: `devdock up && devdock down` passes
- [ ] v0.2 Go Fiber project: `devdock up && devdock down` passes
- [ ] v0.2 Mailpit/MinIO configs still generate valid Compose output
- [ ] Automated backward compat fixture tests: all pass

---

## 18. Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| `registry.devdock.dev` unavailable | Medium | Medium | Bundled templates always work; cached templates work offline; clear network error messages |
| Compromised template archive | Low | High | SHA-256 mandatory; path traversal protection; official-only registry in v0.3 |
| Self-update corrupts binary | Low | High | Checksum verified before any file modification; backup-then-replace; restore on failure |
| Homebrew detection false positive on Intel `/usr/local/bin` | Medium | Medium | Use `brew --prefix devdock` as authoritative method before falling back to path checks |
| Telemetry builds user distrust | Medium | Medium | Opt-in default; no correlatable IDs; fully documented; `DEVDOCK_NO_TELEMETRY` escape hatch |
| Formula checksum mismatch after release | Low | High | Release automation generates and commits checksums atomically; manual verification step in release checklist |
| Archive path traversal exploit in official templates | Very Low | High | Path traversal protection blocks this even for official archives; defense-in-depth |
