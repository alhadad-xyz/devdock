# DevDock v0.3.0 Task Checklist

**Version:** v0.3.0
**Release Theme:** Distribution & Template Infrastructure
**Base:** DevDock v0.2.0 complete
**Ticket prefix:** DD-066 onward (continues from v0.2 DD-040–065)
**Reference:** DevDock v0.3.0 PRD Final

---

## Priority Reference

| Priority | Features |
|---|---|
| P0 — Must ship | Homebrew, release automation, self-update, manifest + cache, template commands, checksum + path traversal protection, hook review, `devdock config`, backward compat tests |
| P1 — Should ship | Telemetry config + event client, first-run prompt, local template warning, README privacy section, `--no-hooks` |
| P2 — Ship if capacity | Direct template URL (`--template-url`), template aliases, template lint command |

P2 does not block release.

---

## Dependency Map

```
DD-066 (release workflow)
├── DD-067 (checksums)
└── DD-068 (homebrew tap)

DD-069 (version metadata)
└── DD-071 (update --check — needs Version variable)

DD-071 (update --check)
└── DD-072 (update download + checksum)
    └── DD-073 (atomic binary replacement)
        └── DD-074 (homebrew detection)

DD-075 (manifest schema + Go structs)
├── DD-076 (registry HTTP client + TTL)
├── DD-077 (cache layout + metadata)
│   └── DD-078 (archive downloader)
│       └── DD-079 (SHA-256 verifier + path traversal)
├── DD-080 (template list)
├── DD-081 (template search)
├── DD-082 (template info)
└── DD-083 (template update)

DD-079 (verifier)
└── DD-084 (create registry integration)

DD-086 (hook review prompt — must exist before --no-hooks)
└── DD-085 (--no-hooks flag)

DD-087 (devdock config command — needed for telemetry)
├── DD-088 (telemetry config fields)
│   └── DD-089 (telemetry event client)
│       └── DD-090 (first-run telemetry prompt)
└── DD-091 (local template warning — uses TTY detection pattern from config)

DD-092–095 (regression + release — depend on all above)
```

---

# Week 1 — Release Distribution

## DD-066 — GitHub Actions Release Workflow

- `[ ]` **DD-066** Create `.github/workflows/release.yml`
- `[ ]` **DD-066** Trigger on tag push matching `v*.*.*`
- `[ ]` **DD-066** Steps: `go test ./...` → build `darwin/arm64` → build `darwin/amd64`
- `[ ]` **DD-066** Inject build variables via `-ldflags`: `Version`, `Commit` (short SHA), `BuildDate`
- `[ ]` **DD-066** Upload both binaries as release assets
- `[ ]` **DD-066** Workflow fails if `go test` fails — no binaries uploaded on test failure

**Done when:**
- `[ ]` **DD-066** Pushing tag `v0.3.0-test` builds both binaries and attaches them to a draft release
- `[ ]` **DD-066** `devdock --version` on each binary prints the injected version, not `dev`
- `[ ]` **DD-066** Pushing a commit that breaks a test: workflow fails before build step

**Unblocks:** DD-067, DD-068.

---

## DD-067 — Release Checksums

- `[ ]` **DD-067** After both binaries are built, generate: `shasum -a 256 devdock-darwin-arm64 devdock-darwin-amd64 > checksums.txt`
- `[ ]` **DD-067** Attach `checksums.txt` to the GitHub Release as a third artifact
- `[ ]` **DD-067** Verify locally: `shasum -a 256 -c checksums.txt` passes after downloading all three files

**Done when:**
- `[ ]` **DD-067** GitHub Release contains exactly three files: `devdock-darwin-arm64`, `devdock-darwin-amd64`, `checksums.txt`
- `[ ]` **DD-067** `shasum -a 256 -c checksums.txt` passes when run in the same directory as the downloaded binaries

**Unblocks:** DD-068 (formula needs SHA-256 values), DD-072 (update download needs to verify against this file).

---

## DD-068 — Homebrew Tap and Formula

- `[ ]` **DD-068** Create GitHub repo `alhadad-xyz/homebrew-tap`
- `[ ]` **DD-068** Add `Formula/devdock.rb`
- `[ ]` **DD-068** Formula uses architecture-specific `url` and `sha256` fields (one block for ARM, one for Intel)
- `[ ]` **DD-068** Formula installs binary as `devdock`
- `[ ]` **DD-068** Formula test block: `system "#{bin}/devdock", "--version"`
- `[ ]` **DD-068** After formula is committed, verify by running `brew tap alhadad-xyz/tap` on a test machine

**Done when:**
- `[ ]` **DD-068** On a clean Apple Silicon Mac with no prior DevDock: `brew install alhadad-xyz/tap/devdock && devdock --version` works
- `[ ]` **DD-068** On a clean Intel Mac (or Rosetta): same test passes
- `[ ]` **DD-068** `brew test devdock` passes

**Unblocks:** Nothing — standalone distribution milestone.

---

## DD-069 — Version Build Metadata

- `[ ]` **DD-069** Add package-level variables in `cmd/devdock/main.go`:
  ```go
  var (
      Version   = "dev"
      Commit    = "unknown"
      BuildDate = "unknown"
  )
  ```
- `[ ]` **DD-069** Update `devdock --version` output:
  - Release binary: `devdock version 0.3.0 (commit: abc1234, built: 2026-06-11)`
  - Dev build: `devdock version dev`
- `[ ]` **DD-069** Update Makefile `make build` to inject a `dev` version for local builds

**Done when:**
- `[ ]` **DD-069** Local `go build` produces binary that prints `devdock version dev`
- `[ ]` **DD-069** Release workflow binary prints `devdock version 0.3.0 (commit: <sha>, built: <date>)`
- `[ ]` **DD-069** `Version` variable is accessible by the update command implementation

**Unblocks:** DD-071.

---

## DD-070 — README Installation Update

- `[ ]` **DD-070** Replace existing installation section with Homebrew as the primary method:
  ```bash
  brew install alhadad-xyz/tap/devdock
  ```
- `[ ]` **DD-070** Add note: "No Go installation required"
- `[ ]` **DD-070** Keep source install as a secondary section:
  ```bash
  go install github.com/alhadad-xyz/devdock/cmd/devdock@latest
  ```
- `[ ]` **DD-070** Add manual download fallback (GitHub Releases page link)
- `[ ]` **DD-070** Add troubleshooting: `devdock: command not found` → PATH instructions for Homebrew

**Done when:**
- `[ ]` **DD-070** A developer who has never used Go can follow the README and successfully install DevDock
- `[ ]` **DD-070** Source install is still documented for contributors

---

## Week 1 End-to-End Validation

```bash
brew tap alhadad-xyz/tap
brew install devdock
devdock --version          # prints correct release version
devdock doctor             # all prerequisite checks pass
devdock --help             # prints help
```

Pass on both Apple Silicon and Intel. No Go installation on the test machine.

---

# Week 2 — Self-Update

## DD-071 — `devdock update --check`

- `[ ]` **DD-071** Add `update` command to Cobra with `--check` flag
- `[ ]` **DD-071** Fetch `https://api.github.com/repos/alhadad-xyz/devdock/releases/latest` with 10-second timeout
- `[ ]` **DD-071** Parse `tag_name`, strip leading `v`, compare with `Version` using semver (`golang.org/x/mod/semver`)
- `[ ]` **DD-071** Exit codes: 0 = up to date, 1 = update available, 2 = network/API error
- `[ ]` **DD-071** Print format when update available:
  ```
  DevDock v0.2.0 is installed. v0.3.0 is available.
  Run `devdock update` to upgrade.
  ```

**Done when:**
- `[ ]` **DD-071** On a binary built with `Version=0.2.0`, running against a mock GitHub response showing `0.3.0`: prints update message, exits 1
- `[ ]` **DD-071** On a binary with same version as mock latest: prints "up to date", exits 0
- `[ ]` **DD-071** Mock network failure: prints what/why/fix error, exits 2
- `[ ]` **DD-071** Unit test covers all three paths using a mock HTTP server

**Unblocks:** DD-072.

---

## DD-072 — Update Download and Checksum Verification

- `[ ]` **DD-072** Detect `runtime.GOOS` + `runtime.GOARCH`; resolve binary name: `devdock-darwin-arm64` or `devdock-darwin-amd64`
- `[ ]` **DD-072** Download binary asset to `~/.devdock/tmp/devdock-new` (120-second timeout)
- `[ ]` **DD-072** Download `checksums.txt` (10-second timeout)
- `[ ]` **DD-072** Parse `checksums.txt`: find the line where the second field matches the binary name; extract the hex digest
- `[ ]` **DD-072** Compute SHA-256 of the downloaded binary; compare against extracted digest
- `[ ]` **DD-072** If mismatch: delete `~/.devdock/tmp/devdock-new`, print error showing both expected and actual hex, exit 1
- `[ ]` **DD-072** If match: proceed to DD-073

**Done when:**
- `[ ]` **DD-072** Apple Silicon: `devdock-darwin-arm64` is downloaded
- `[ ]` **DD-072** Intel: `devdock-darwin-amd64` is downloaded
- `[ ]` **DD-072** Deliberate checksum mismatch (inject wrong expected hash in test): download deleted, error printed, old binary untouched
- `[ ]` **DD-072** Network timeout after 120s: temp file cleaned up, clear error message

**Unblocks:** DD-073.

---

## DD-073 — Atomic Binary Replacement

- `[ ]` **DD-073** Detect current binary path: `os.Executable()`, resolve symlinks
- `[ ]` **DD-073** Check write permission on the binary's parent directory
- `[ ]` **DD-073** Backup: `os.Rename(currentPath, currentPath+".backup")`
- `[ ]` **DD-073** Move new binary: `os.Rename(tempPath, currentPath)`
- `[ ]` **DD-073** Set executable bit: `os.Chmod(currentPath, 0755)`
- `[ ]` **DD-073** If move or chmod fails: `os.Rename(currentPath+".backup", currentPath)` to restore
- `[ ]` **DD-073** On success: delete backup, print `Updated to DevDock v<new version>.`
- `[ ]` **DD-073** If write permission denied: print error with `sudo` or permission fix suggestion

**Done when:**
- `[ ]` **DD-073** Manual install (binary in `/usr/local/bin`): update replaces it, new `devdock --version` prints new version
- `[ ]` **DD-073** Simulated failure (mock `os.Rename` failure): old binary is restored, exit 1, no temp files left
- `[ ]` **DD-073** After successful update: `*.backup` file is deleted

**Unblocks:** DD-074.

---

## DD-074 — Homebrew-Managed Install Detection

- `[ ]` **DD-074** If `brew` is in PATH: run `brew --prefix devdock 2>/dev/null`
- `[ ]` **DD-074** If output is non-empty and current binary path starts with that prefix: binary is Homebrew-managed
- `[ ]` **DD-074** Fallback (no `brew` command): check if path starts with `/opt/homebrew/` or `/usr/local/Cellar/`
- `[ ]` **DD-074** Do NOT treat `/usr/local/bin/` alone as a Homebrew indicator
- `[ ]` **DD-074** If Homebrew-managed: print the redirect message and exit 0 — do not download or modify anything

**Done when:**
- `[ ]` **DD-074** Homebrew-installed binary: `devdock update` prints redirect, no files touched
- `[ ]` **DD-074** Manually installed binary in `/usr/local/bin` (not Homebrew): `devdock update` proceeds with download
- `[ ]` **DD-074** Test: mock `brew --prefix devdock` returning a path prefix; verify detection is positive

---

## Week 2 End-to-End Validation

```bash
# Test 1: manual install update
# (On a machine with manually installed v0.2 binary)
devdock update --check       # exits 1, shows available version
devdock update               # downloads, verifies, replaces
devdock --version            # shows new version

# Test 2: homebrew install redirect
# (On a machine with Homebrew-installed devdock)
devdock update               # prints brew upgrade message, exits 0
which devdock                # path unchanged
```

---

# Week 3 — Template Registry Core

## DD-075 — Registry Manifest Schema and Go Structs

- `[ ]` **DD-075** Define `RegistryManifest`, `TemplateEntry` Go structs in `internal/registry/schema.go`
- `[ ]` **DD-075** Fields per PRD Section 7.2
- `[ ]` **DD-075** Validation function: `ValidateManifest(m RegistryManifest) error`
  - `schema_version` must be `"1"`
  - Each template's `checksum_sha256` must be exactly 64 lowercase hex characters (regex: `^[a-f0-9]{64}$`)
  - Duplicate `id` values: return error
- `[ ]` **DD-075** Add test fixtures: valid manifest, missing checksum, wrong checksum format, unsupported schema version, duplicate ID

**Done when:**
- `[ ]` **DD-075** Valid manifest fixture: `ValidateManifest` returns nil
- `[ ]` **DD-075** Missing checksum: returns error naming the template ID
- `[ ]` **DD-075** Checksum placeholder `"sha256-value-here"`: returns error (fails the 64-char hex regex)
- `[ ]` **DD-075** Wrong schema version: returns error
- `[ ]` **DD-075** Duplicate ID: returns error

**Unblocks:** DD-076, DD-077, DD-080–083.

---

## DD-076 — Registry HTTP Client with TTL

- `[ ]` **DD-076** Implement `registry.FetchManifest(url string) (*RegistryManifest, error)` with 10-second timeout
- `[ ]` **DD-076** Read `registry.url` from global config; read `DEVDOCK_REGISTRY_URL` env var (env var takes priority)
- `[ ]` **DD-076** Implement TTL check: read `manifest_cached_at` from `cache.json`; if age < 24h, return cached manifest without HTTP call
- `[ ]` **DD-076** If online and cache is stale (≥ 24h): fetch fresh manifest, save to `~/.devdock/templates/manifest.json`, update `cache.json`
- `[ ]` **DD-076** If offline (network error) and cache exists (even stale): return cached manifest with `CacheStale: true` flag
- `[ ]` **DD-076** If offline and no cache: return error

**Done when:**
- `[ ]` **DD-076** Fresh manifest fetched and cached: second call within 24h makes no HTTP request (verify with mock server call count)
- `[ ]` **DD-076** After 24h (mock `manifest_cached_at` to a 25h-old timestamp): fresh fetch occurs
- `[ ]` **DD-076** Network blocked, cache present: cached manifest returned with stale flag, no panic
- `[ ]` **DD-076** Network blocked, no cache: structured error returned

**Unblocks:** DD-080–083, DD-084.

---

## DD-077 — Template Cache Layout and Metadata

- `[ ]` **DD-077** Create cache directories on first use: `~/.devdock/templates/`, `~/.devdock/templates/archives/`
- `[ ]` **DD-077** Implement `cache.Load() (*CacheMetadata, error)` — reads `cache.json`; if file missing, returns empty struct (not an error)
- `[ ]` **DD-077** Implement `cache.Save(m *CacheMetadata) error` — writes atomically (temp file + rename)
- `[ ]` **DD-077** Implement `cache.GetTemplate(id, version string) (*CachedTemplate, bool)` — returns nil if not found
- `[ ]` **DD-077** Implement `cache.PutTemplate(entry CachedTemplate) error` — adds or updates entry, saves
- `[ ]` **DD-077** All `path` fields stored as absolute paths using `filepath.Abs()` — never store `~`
- `[ ]` **DD-077** If `cache.json` is corrupt (invalid JSON): log warning, treat as empty cache, do not error

**Done when:**
- `[ ]` **DD-077** Empty cache: `Load()` returns empty struct without error
- `[ ]` **DD-077** After `PutTemplate()`: `GetTemplate()` returns the entry with correct absolute path
- `[ ]` **DD-077** Simulate corrupt `cache.json`: `Load()` returns empty struct and logs warning (verify with test)
- `[ ]` **DD-077** Paths in saved entries: never start with `~`, always start with `/`

**Unblocks:** DD-078.

---

## DD-078 — Template Archive Downloader

- `[ ]` **DD-078** Implement `downloader.Fetch(archiveURL, destPath string) error`
- `[ ]` **DD-078** Download to `destPath + ".tmp"` (temp file, same directory)
- `[ ]` **DD-078** On success: `os.Rename(tempPath, destPath)` (atomic)
- `[ ]` **DD-078** On failure or cancellation: `os.Remove(tempPath)` — no partial archive left
- `[ ]` **DD-078** 120-second timeout for download
- `[ ]` **DD-078** Support `.tar.gz` only; reject other formats with a clear error

**Done when:**
- `[ ]` **DD-078** Successful download: archive file exists at `destPath`, no `.tmp` file remains
- `[ ]` **DD-078** Simulated network failure mid-download: temp file deleted, `destPath` does not exist
- `[ ]` **DD-078** Non-`.tar.gz` URL: error returned before download attempt

**Unblocks:** DD-079.

---

## DD-079 — SHA-256 Verifier and Path Traversal Protection

**(Part 1 — Checksum Verifier):**
- `[ ]` **DD-079** Implement `verifier.CheckArchive(archivePath, expectedHex string) error`
- `[ ]` **DD-079** Read archive file, compute SHA-256
- `[ ]` **DD-079** Compare against `expectedHex` (must be exactly 64 hex chars — pre-validate)
- `[ ]` **DD-079** On mismatch: delete the archive file, return error showing both expected and actual values
- `[ ]` **DD-079** Update `cache.json` entry: `"verified": true` after successful verification

**(Part 2 — Path Traversal Protection):**
- `[ ]` **DD-079** Implement `extractor.Extract(archivePath, destDir string) error`
- `[ ]` **DD-079** Before writing any file: resolve the full path of the entry
- `[ ]` **DD-079** If resolved path does not start with `filepath.Clean(destDir) + string(os.PathSeparator)`: reject the entire archive
- `[ ]` **DD-079** On rejection: remove any partially extracted files, return error
- `[ ]` **DD-079** On success: all files extracted under `destDir`

**Done when (checksum):**
- `[ ]` **DD-079** Correct checksum: `CheckArchive` returns nil, `verified: true` in cache
- `[ ]` **DD-079** Wrong checksum: archive file deleted, error message includes both hex values
- `[ ]` **DD-079** Empty/missing archive file: error returned, no panic

**Done when (path traversal):**
- `[ ]` **DD-079** Archive with normal paths: extracted successfully to `destDir`
- `[ ]` **DD-079** Archive containing `../../etc/passwd` entry: entire extraction rejected, zero files written to `destDir`, error message names the offending path
- `[ ]` **DD-079** Unit test using a crafted test archive with a traversal entry

**Unblocks:** DD-084 (create integration), DD-085 (--no-hooks), DD-086 (hook review).

---

## Week 3 Validation

Three separate manual checks — all must pass before moving to Week 4:

```bash
# Check 1: Online flow
DEVDOCK_REGISTRY_URL=<test-registry> devdock template update
# → manifest fetched and cached
# → ls ~/.devdock/templates/ shows manifest.json, cache.json

# Check 2: Checksum verification
# Manually corrupt the archive after download; re-run create
devdock create next-postgres test-corrupt
# → error: checksum mismatch; corrupted archive deleted from cache

# Check 3: Offline flow
# Disconnect network (or use DEVDOCK_REGISTRY_URL pointing to unreachable host)
devdock template list
# → shows cached manifest with "(may be outdated)" notice
devdock create next-postgres offline-test --offline
# → succeeds from cache (if template was cached in Check 1)
```

---

# Week 4 — Template Commands and Create Integration

## DD-080 — `devdock template list`

- `[ ]` **DD-080** Add `template` command group to Cobra
- `[ ]` **DD-080** Add `list` subcommand
- `[ ]` **DD-080** Call `registry.FetchManifest()` (respects TTL — no network if cache fresh)
- `[ ]` **DD-080** Print table: ID | Version | Category | Description
- `[ ]` **DD-080** Header shows cache age: `(cached Xh ago)` or `(cached manifest — may be outdated)` for stale
- `[ ]` **DD-080** Do not download archives
- `[ ]` **DD-080** If offline with no cache: what/why/fix error

**Done when:**
- `[ ]` **DD-080** `devdock template list` prints formatted table
- `[ ]` **DD-080** Second call within 24h: no HTTP request made (verify with mock server)
- `[ ]` **DD-080** Offline + stale cache: table shown with "may be outdated" notice
- `[ ]` **DD-080** Offline + no cache: helpful error shown

---

## DD-081 — `devdock template search <query>`

- `[ ]` **DD-081** Add `search` subcommand
- `[ ]` **DD-081** Load manifest from cache (same TTL logic as list)
- `[ ]` **DD-081** Filter: substring match against ID, name, description, category, tags, runtime values, service names
- `[ ]` **DD-081** Case-insensitive
- `[ ]` **DD-081** Print same table format as `list`
- `[ ]` **DD-081** Zero results: `No templates found for '<query>'. Run 'devdock template list' to see all templates.`

**Done when:**
- `[ ]` **DD-081** `devdock template search postgres` returns all templates with postgres in any field
- `[ ]` **DD-081** `devdock template search NEXTJS` (uppercase): finds nextjs templates
- `[ ]` **DD-081** `devdock template search zzznotfound`: prints no-results message

---

## DD-082 — `devdock template info <id>`

- `[ ]` **DD-082** Add `info` subcommand
- `[ ]` **DD-082** Look up template in manifest
- `[ ]` **DD-082** Check cache metadata for cached/verified status
- `[ ]` **DD-082** Print all fields per PRD Section 6.6 format
- `[ ]` **DD-082** For uncached template: `Cached: no`
- `[ ]` **DD-082** For cached but unverified: `Cached: yes (unverified — run 'devdock template update')`
- `[ ]` **DD-082** Unknown ID: helpful error with `devdock template list` as fix

**Done when:**
- `[ ]` **DD-082** `devdock template info next-postgres` on a cached+verified template shows all fields with `Cached: yes (verified)`
- `[ ]` **DD-082** Same command on an uncached template: `Cached: no`
- `[ ]` **DD-082** `devdock template info unknownid`: helpful error

---

## DD-083 — `devdock template update`

- `[ ]` **DD-083** Fetch manifest fresh (bypasses TTL — always re-fetches on explicit `update`)
- `[ ]` **DD-083** Save to `~/.devdock/templates/manifest.json`, update TTL
- `[ ]` **DD-083** Compare each template in manifest against `cache.json` entries
- `[ ]` **DD-083** Print summary using symbols: `✔` (cached up-to-date), `↑` (new version), `+` (new template), `=` (unchanged, not cached)
- `[ ]` **DD-083** `--all` flag: after manifest update, download + verify all templates into cache

**Done when:**
- `[ ]` **DD-083** `devdock template update` always re-fetches even if cache is fresh (no TTL bypass — explicit command)
- `[ ]` **DD-083** Summary output matches the spec format with correct symbols
- `[ ]` **DD-083** `devdock template update --all`: all templates cached and verified after completion
- `[ ]` **DD-083** Network failure: cache not corrupted, error message shown

---

## DD-084 — `devdock create` Registry Template Resolution

- `[ ]` **DD-084** Update resolution order in create command:
  1. Local path (argument starts with `/`, `./`, or `../`)
  2. Cached verified template (from `cache.json`)
  3. Bundled template (embedded in binary)
  4. Remote registry template (fetch, verify, cache, then use)
  5. Error: "Template '<id>' not found in cache, bundled templates, or registry"
- `[ ]` **DD-084** Add `--offline` flag: skip step 4; error if not in cache or bundle
- `[ ]` **DD-084** Preserve atomic project creation from v0.1/v0.2

**Done when:**
- `[ ]` **DD-084** `devdock create next-postgres test` with template not cached: downloads, verifies, creates project
- `[ ]` **DD-084** Same command repeated: uses cache (no download)
- `[ ]` **DD-084** `devdock create next-postgres test --offline` with cached template: succeeds
- `[ ]` **DD-084** `devdock create next-postgres test --offline` with no cache: clear error
- `[ ]` **DD-084** `devdock create laravel-api test` (bundled): works without internet even in v0.3
- `[ ]` **DD-084** Resolution order test: if template is in both cache and bundle, cache is used (verify with mock)

---

## DD-085 — `--no-hooks` Flag on `devdock create`

**Note:** This ticket depends on DD-086 (hook review prompt) being implemented first. If Week 4 runs long, implement DD-086 before DD-085.

- `[ ]` **DD-085** Add `--no-hooks` flag to `devdock create`
- `[ ]` **DD-085** When set: skip post-create hooks, skip hook review prompt
- `[ ]` **DD-085** Scaffold command still runs
- `[ ]` **DD-085** Print after create: "Post-create hooks were skipped. You may need to run the following manually:" followed by each skipped hook command

**Done when:**
- `[ ]` **DD-085** `devdock create next-postgres test --no-hooks`: creates project, no prompt shown, hooks not run
- `[ ]` **DD-085** Output lists the skipped hook commands
- `[ ]` **DD-085** Scaffold command (`pnpm create next-app .`) still runs
- `[ ]` **DD-085** Without `--no-hooks`: hook review prompt appears (DD-086 behavior)

---

## DD-086 — Hook Review Prompt

- `[ ]` **DD-086** Before running scaffold or post-create hooks for any template (bundled OR registry), print all commands grouped by type (scaffold vs post-create)
- `[ ]` **DD-086** Show source label: `Source: official (registry.devdock.dev)` or `Source: bundled`
- `[ ]` **DD-086** Prompt: `Run these commands? (Y/n)` — default Y
- `[ ]` **DD-086** If user says no:
  - No files written to final project directory (temp dir cleaned up)
  - Exit 0: "Create cancelled. No files were written."
- `[ ]` **DD-086** If stdin is not a TTY (CI/piped): auto-proceed with scaffold + skip hooks (same as `--no-hooks`), print warning: "Non-interactive mode: hooks skipped automatically."

**Done when:**
- `[ ]` **DD-086** Interactive terminal: hook review prompt shown, user can cancel
- `[ ]` **DD-086** Cancelled create: no final project directory exists, no temp files left
- `[ ]` **DD-086** Non-TTY: auto-proceeds without prompt, hooks skipped, warning printed
- `[ ]` **DD-086** `--no-hooks` flag: prompt not shown at all (handled in DD-085)

**Unblocks:** DD-085.

---

## Week 4 End-to-End Validation

```bash
devdock template list
devdock template search next
devdock template info next-postgres
devdock template update --all
devdock create next-postgres my-saas       # shows hook review prompt
# answer Y to proceed
cd my-saas && devdock up && curl localhost:3000
cd .. && rm -rf my-saas

devdock create next-postgres no-hook-test --no-hooks
# no prompt, hooks skipped, list of skipped hooks shown

devdock create next-postgres offline-test --offline
# uses cache
```

---

# Week 5 — `devdock config`, Telemetry, Security

## DD-087 — `devdock config` Command

- `[ ]` **DD-087** Add `config` command group to Cobra
- `[ ]` **DD-087** Add `get`, `set`, `list` subcommands
- `[ ]` **DD-087** Implement key registry: supported keys, types, validation rules, defaults (per PRD Section 12.2)
- `[ ]` **DD-087** `get <key>`: print value (from config or default if not set); unknown key → error with valid keys list
- `[ ]` **DD-087** `set <key> <value>`: validate type, write `~/.devdock/config.yml` atomically, print confirmation
- `[ ]` **DD-087** `list`: print all keys with current values and defaults in a table
- `[ ]` **DD-087** Unknown key for `set` or `get`: what/why/fix error naming the key and listing valid ones

**Done when:**
- `[ ]` **DD-087** `devdock config list`: prints all 4 keys with current values
- `[ ]` **DD-087** `devdock config get telemetry`: prints `false` on fresh install
- `[ ]` **DD-087** `devdock config set telemetry true`: updates config file, prints confirmation
- `[ ]` **DD-087** `devdock config set telemetry invalidvalue`: type validation error
- `[ ]` **DD-087** `devdock config set unknown.key value`: unknown key error
- `[ ]` **DD-087** Config writes are atomic (temp file + rename)

**Unblocks:** DD-088, DD-089, DD-090.

---

## DD-088 — Telemetry Config Fields

- `[ ]` **DD-088** Add `telemetry.enabled` and `telemetry.prompted` to global config schema
- `[ ]` **DD-088** Wire `devdock config set/get telemetry` to `telemetry.enabled`
- `[ ]` **DD-088** Read `DEVDOCK_NO_TELEMETRY` env var: if set to `1` or `true`, disable telemetry regardless of config
- `[ ]` **DD-088** Add `telemetry.IsEnabled()` function: returns false if `DEVDOCK_NO_TELEMETRY` is set, otherwise reads config

**Done when:**
- `[ ]` **DD-088** Fresh install: `telemetry.enabled` is `false`, `prompted` is `false`
- `[ ]` **DD-088** `devdock config set telemetry true` → `telemetry.enabled: true` in config file
- `[ ]` **DD-088** `DEVDOCK_NO_TELEMETRY=1 devdock config get telemetry` → reports disabled (env var override)
- `[ ]` **DD-088** `telemetry.IsEnabled()` returns `false` when env var is set, even if config is `true`

**Unblocks:** DD-089, DD-090.

---

## DD-089 — Telemetry Event Client

- `[ ]` **DD-089** Implement `telemetry.Send(event string, properties map[string]string) error`
- `[ ]` **DD-089** POST to `https://telemetry.devdock.dev/v1/events` with JSON payload:
  ```json
  {"event": "...", "properties": {"command": "...", "devdock_version": "...", "os": "...", "arch": "..."}}
  ```
- `[ ]` **DD-089** No session ID, no user ID, no device ID — every event is stateless
- `[ ]` **DD-089** Property whitelist: each event type has a fixed allowed set (per PRD Section 11.3); extra properties silently dropped
- `[ ]` **DD-089** Failed send: silently discard, never surface error to user
- `[ ]` **DD-089** If `DEVDOCK_DEBUG=1`: log failed sends to stderr (debug only)
- `[ ]` **DD-089** 5-second timeout on POST; fire-and-forget (non-blocking — use goroutine)
- `[ ]` **DD-089** Only call if `telemetry.IsEnabled()` returns true

**Done when:**
- `[ ]` **DD-089** `telemetry.enabled: false`: no HTTP calls made (verify with mock server: zero requests)
- `[ ]` **DD-089** `telemetry.enabled: true`: correct payload structure sent
- `[ ]` **DD-089** Mock server returns 500: no error surfaced to user, command proceeds normally
- `[ ]` **DD-089** Payload never contains `command_args`, `project_path`, `project_name`, or any field not in the whitelist

**Unblocks:** DD-090.

---

## DD-090 — First-Run Telemetry Prompt

- `[ ]` **DD-090** On the first `devdock` command where `telemetry.prompted` is `false`:
  - Print the prompt from PRD Section 11.1 with collected/never-collected summary
  - Default is `N`
  - If non-interactive (no TTY): skip prompt entirely, leave telemetry disabled, set `prompted: true`
  - On user answer: set `telemetry.enabled` accordingly, set `telemetry.prompted: true`, save config
- `[ ]` **DD-090** Prompt shown only once per machine — never repeated

**Done when:**
- `[ ]` **DD-090** Fresh install interactive: prompt shown on first command; not shown on second command
- `[ ]` **DD-090** Fresh install non-interactive: prompt never shown; `prompted: true` set; telemetry stays disabled
- `[ ]` **DD-090** After prompt: `devdock config get telemetry` reflects the user's choice

---

## DD-091 — Local Template Warning

- `[ ]` **DD-091** If `devdock create` argument starts with `/`, `./`, or `../`: treat as local path
- `[ ]` **DD-091** Print warning per PRD Section 8.4
- `[ ]` **DD-091** Prompt: `Continue? (y/N)` — default N
- `[ ]` **DD-091** Non-interactive (no TTY): refuse to proceed unless `--force-local` flag is passed
  - Error: "Local template requires explicit consent in non-interactive mode. Use --force-local to proceed."
- `[ ]` **DD-091** Still apply path traversal protection during extraction even for local templates

**Done when:**
- `[ ]` **DD-091** `devdock create ./my-template app`: warning shown, default N, proceeds on Y
- `[ ]` **DD-091** Non-TTY without `--force-local`: error with `--force-local` suggestion
- `[ ]` **DD-091** Non-TTY with `--force-local`: proceeds with warning printed but no prompt
- `[ ]` **DD-091** Path traversal protection still applies to local templates (verify with crafted test archive)

---

## DD-090b — Security and Privacy README Section

- `[ ]` **DD-090b** Add README section: **Template Security**
  - Explains checksum verification
  - Explains hook review prompt and `--no-hooks`
  - Explains path traversal protection
  - Explains local template warning
- `[ ]` **DD-090b** Add README section: **Telemetry & Privacy**
  - Lists what is collected (verbatim from PRD Section 11.3)
  - Lists what is never collected
  - Shows `devdock config set telemetry false` and `DEVDOCK_NO_TELEMETRY=1`

**Done when:**
- `[ ]` **DD-090b** README clearly explains both sections
- `[ ]` **DD-090b** New user can understand template risks and telemetry opt-out without reading the PRD

---

## Week 5 Validation

```bash
# Security
devdock create next-postgres hook-test
# → hook review prompt shown; answer N; no directory created

devdock create next-postgres hook-test --no-hooks
# → no prompt; hooks skipped; list of skipped commands shown

devdock create ./local-template local-test
# → local warning shown; default N; proceeds on Y

# Config
devdock config list
devdock config get telemetry        # false
devdock config set telemetry true
devdock config get telemetry        # true
DEVDOCK_NO_TELEMETRY=1 devdock config get telemetry   # still reports disabled
devdock config set telemetry false
devdock config set badkey value     # error: unknown key
```

---

# Week 6 — Regression Testing and Release

## DD-092 — Backward Compatibility Regression Suite

- `[ ]` **DD-092** Run automated fixture tests (from v0.2 DD-061 suite) — must still pass
- `[ ]` **DD-092** Manual regression tests for all previous stable flows:
  - v0.1 Laravel: `devdock init && devdock up && devdock down`
  - v0.1 Next.js: `devdock init && devdock up && devdock down`
  - v0.2 Express: `devdock init && devdock up && devdock down`
  - v0.2 Go Fiber: `devdock init && devdock up && devdock down`
  - v0.2 service commands: `devdock service add mailpit && devdock up && devdock open mailpit && devdock down`
- `[ ]` **DD-092** Test bundled template create with registry unreachable: `devdock create laravel-api test` with `DEVDOCK_REGISTRY_URL=http://localhost:9999` (unreachable) — must succeed from bundled

**Done when:**
- `[ ]` **DD-092** All automated fixture tests pass
- `[ ]` **DD-092** All 5 manual regression flows pass
- `[ ]` **DD-092** Bundled template create succeeds with registry unreachable

---

## DD-093 — Template Registry End-to-End Test

- `[ ]` **DD-093** Test full registry flow:
  ```bash
  devdock template update           # fetches and caches manifest
  devdock template list             # reads from cache (verify no second HTTP call)
  devdock template info next-postgres
  devdock template update --all     # downloads and verifies all templates
  devdock create next-postgres registry-test   # uses cached registry template
  devdock create next-postgres offline-test --offline  # uses cache
  ```
- `[ ]` **DD-093** Test checksum rejection:
  - Manually corrupt the cached archive (flip one byte)
  - `devdock create next-postgres corrupt-test` → error, corrupt archive deleted, no project created
- `[ ]` **DD-093** Test path traversal rejection:
  - Use a crafted test archive with a `../` entry
  - `devdock create ./test-traversal.tar.gz test` → error, no files written outside target

**Done when:** All steps above pass. Corruption and traversal attacks are both blocked.

---

## DD-094 — Homebrew and Update Integration Test

- `[ ]` **DD-094** On a clean macOS machine with no prior DevDock:
  ```bash
  brew install alhadad-xyz/tap/devdock
  devdock --version
  devdock update --check
  devdock doctor
  ```
- `[ ]` **DD-094** On a machine with manually installed v0.2 binary (not Homebrew):
  ```bash
  devdock update --check     # should show v0.3 available
  devdock update             # downloads, verifies, replaces
  devdock --version          # shows v0.3
  ```
- `[ ]` **DD-094** On a Homebrew-installed machine:
  ```bash
  devdock update             # prints 'brew upgrade devdock', exits 0
  devdock --version          # unchanged (Homebrew binary not modified)
  ```

**Done when:** All three scenarios pass.

---

## DD-095 — README and Release Notes

- `[ ]` **DD-095** Final README review: all v0.3 features documented (installation, template commands, update command, offline mode, security, telemetry)
- `[ ]` **DD-095** Add `CHANGELOG.md` entry for v0.3.0 with "New / Improved / Still Not Included" sections
- `[ ]` **DD-095** Add migration notes: "No migration required — v0.1/v0.2 projects work without changes"

**Done when:**
- `[ ]` **DD-095** A developer who has never used DevDock can follow README from installation through `devdock create` using only the README
- `[ ]` **DD-095** Changelog entry is present and accurate

---

## DD-096 — v0.3.0 Release

- `[ ]` **DD-096** Bump `Version` default from `dev` to `0.3.0` (release workflow will override via `-ldflags`)
- `[ ]` **DD-096** Create and push tag `v0.3.0`
- `[ ]` **DD-096** Verify GitHub Actions release workflow completes successfully
- `[ ]` **DD-096** Verify all three release artifacts are attached: both binaries + checksums
- `[ ]` **DD-096** Verify Homebrew formula update commit is pushed to `homebrew-tap`
- `[ ]` **DD-096** Verify `brew upgrade devdock` works on a machine with v0.2 installed
- `[ ]` **DD-096** Publish release notes from changelog
- `[ ]` **DD-096** Ask v0.2 testers to validate:
  ```bash
  brew install alhadad-xyz/tap/devdock
  devdock template list
  devdock create next-postgres my-test
  devdock update --check
  ```

**Done when:**
- `[ ]` **DD-096** v0.3.0 GitHub Release is published with all artifacts
- `[ ]` **DD-096** Homebrew formula installs v0.3.0
- `[ ]` **DD-096** At least 2 external testers complete the 4-command validation above successfully

---

# P2 — Ship If Capacity Allows

## DD-097 — Direct Template URL (`--template-url`)

- `[ ]` **DD-097** Add `--template-url <url>` flag to `devdock create`
- `[ ]` **DD-097** Requires `--checksum <sha256hex>` to be provided — reject otherwise
- `[ ]` **DD-097** Download, verify checksum, extract with path traversal protection
- `[ ]` **DD-097** Show local-template-style warning: "Creating from direct URL — not registry-verified"

**Done when:** Direct URL create works only when `--checksum` is provided and verified. Without `--checksum`, clear error.

---

## DD-098 — Template Aliases

- `[ ]` **DD-098** Add `"aliases": ["next", "nextjs-postgres"]` field to manifest template entry schema
- `[ ]` **DD-098** `devdock template search` and `devdock create` resolve aliases to canonical ID
- `[ ]` **DD-098** `devdock template info <alias>` redirects to canonical template

**Done when:** `devdock create next my-app` resolves to `next-postgres` if alias `next` is defined in manifest.

---

## DD-099 — Template Lint Command

- `[ ]` **DD-099** Add `devdock template lint <path>`
- `[ ]` **DD-099** Validates `template.yml`: required fields, variable syntax, scaffold command presence
- `[ ]` **DD-099** Validates overlay files exist
- `[ ]` **DD-099** Validates hook file references
- `[ ]` **DD-099** Reports pass/fail per check

**Done when:** Valid template passes all checks. Template missing `template.yml` fails with clear error listing what's missing.

---

# v0.3 Definition of Done Sign-Off

Must all be checked before tagging v0.3.0:

## Distribution
- [ ] GitHub Actions builds both binaries with correct version string
- [ ] `checksums.txt` generated and attached to release
- [ ] `brew install alhadad-xyz/tap/devdock` works on Apple Silicon
- [ ] `brew install alhadad-xyz/tap/devdock` works on Intel
- [ ] `devdock --version` shows correct version after Homebrew install
- [ ] README uses Homebrew as primary install method

## Update System
- [ ] `devdock update --check` exits 0/1/2 correctly per PRD spec
- [ ] `devdock update` downloads correct architecture binary
- [ ] Checksum verified before any file modification
- [ ] Checksum mismatch: original binary unchanged
- [ ] Homebrew install: `brew upgrade devdock` redirect, no files modified
- [ ] Successful update: `devdock --version` shows new version

## Template Registry
- [ ] Manifest fetched and cached with 24h TTL
- [ ] Stale cache used offline with notice; fresh fetch when online + stale
- [ ] No cache + offline: clear error with fix
- [ ] `devdock template list` format matches spec
- [ ] `devdock template search` case-insensitive, covers all fields
- [ ] `devdock template info` shows cached/verified status correctly
- [ ] `devdock template update` prints new/updated/unchanged per symbol format
- [ ] `devdock template update --all` downloads and verifies all templates
- [ ] SHA-256 mismatch: archive deleted, error with both hex values shown
- [ ] Missing checksum: template rejected
- [ ] Path traversal in archive: entire archive rejected, zero files written
- [ ] `--offline` works from cache; fails clearly if not cached

## Create Integration
- [ ] Resolution order: cached → bundled → remote (tested)
- [ ] Hook review prompt shown for all template sources
- [ ] Cancel at hook review: no final directory, exit 0
- [ ] `--no-hooks`: hooks skipped, prompt skipped, list of skipped commands printed
- [ ] Non-TTY: auto-skips hooks with warning
- [ ] Bundled templates work when registry unreachable
- [ ] Failed create leaves no partial final directory

## Security & Privacy
- [ ] Path traversal protection: malicious archive test passes
- [ ] Local template warning shown; default N; blocked in non-TTY without `--force-local`
- [ ] Telemetry disabled by default
- [ ] First-run prompt default is N; non-TTY skips prompt silently
- [ ] `devdock config set telemetry true/false` works
- [ ] `DEVDOCK_NO_TELEMETRY=1` overrides config
- [ ] Telemetry failure silently discarded, never fails user command
- [ ] Telemetry payload never contains paths, names, args, or secrets (test with mock server)
- [ ] README security section documents all protections

## `devdock config`
- [ ] `config list` shows all supported keys
- [ ] `config get <key>` returns current or default value
- [ ] `config set <key> <value>` validates and writes atomically
- [ ] `config set <unknown>` prints error with valid keys
- [ ] Config writes are atomic

## Backward Compatibility
- [ ] Automated fixture tests: all pass (v0.1, v0.2 configs)
- [ ] v0.1 Laravel project: `devdock up && devdock down` passes
- [ ] v0.1 Next.js project: `devdock up && devdock down` passes
- [ ] v0.2 Express + Go Fiber + Mailpit + MinIO: all flows pass
- [ ] Bundled template create succeeds with registry unreachable
