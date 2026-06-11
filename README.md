# DevDock ![Version](https://img.shields.io/badge/version-v0.2.0%20beta-blue)

![DevDock Terminal Demo](assets/demo.gif)

DevDock is a local development launcher that detects your project, provides a guided configuration, and starts your app alongside Dockerized services like MySQL, PostgreSQL, Redis, Mailpit, and MinIO.

DevDock removes repetitive local setup work: no more writing the same Compose boilerplate, manually checking ports, or remembering service commands for every project.

## ✨ Features

- **🪄 Intelligent Auto-Detection:** Automatically detects your supported framework (Laravel, Next.js, Express, Go Fiber, or existing docker-compose projects).
- **🚀 One-Command Spin Up:** Run your local framework process (e.g., `npm run dev`) *and* your Docker databases concurrently with a single command.
- **🛡️ Port Conflict Detection:** Dynamically scans your host machine before starting. If another process is hogging a port, DevDock safely stops and tells you exactly what to fix.
- **👻 Detached Mode:** Run `devdock up --detach` to background your application. DevDock safely manages the PID and tears it down gracefully when you run `devdock down`.
- **📊 Unified Observability:** Run `devdock status` for a clean dashboard of your Docker services and your host application process, complete with database connection strings.
- **🔧 Direct UI Access:** Run `devdock open mailpit` to instantly launch the Web UI in your browser.
- **🏃 Task Runner:** Define and run project-specific scripts via `devdock run`.

## 🛠 Supported Stacks

| Framework | Supported Services | Notes |
|-----------|--------------------|-------|
| **Laravel** | MySQL, PostgreSQL, Redis, Mailpit, MinIO | Detects `composer.json` and `artisan`. Generates `.env`. |
| **Next.js** | PostgreSQL, MySQL, Redis, Mailpit, MinIO | Detects `package.json` and `next.config.*`. |
| **Express** | PostgreSQL, MySQL, Redis, Mailpit, MinIO | Detects `express` in `package.json`. Sets up `node index.js` fallback. |
| **Go Fiber** | PostgreSQL, MySQL, Redis, Mailpit, MinIO | Detects `github.com/gofiber/fiber` in `go.mod`. Sets up `go run .` by default. |
| **Docker Compose** | *Any* | Existing `compose.yml`. DevDock acts as a friendly wrapper. |

## 📦 Installation

**Using Homebrew (Recommended for macOS)**
```bash
brew install alhadad-xyz/tap/devdock
```
*(Note: DevDock is a standalone binary. No Go installation is required when installing via Homebrew!)*

**Building from Source**
If you prefer to build from source or are on Linux:
```bash
go install github.com/alhadad-xyz/devdock/cmd/devdock@latest
```
*Prerequisite: Ensure `~/go/bin` is in your `$PATH`. Requires Go 1.23+.*

**Troubleshooting:**
- **`devdock: command not found` after Homebrew install**: Ensure Homebrew's `bin` directory (`/opt/homebrew/bin` on Apple Silicon or `/usr/local/bin` on Intel) is added to your PATH in `~/.zshrc` or `~/.bashrc`.
- **Manual Download**: You can also download pre-compiled binaries directly from the [GitHub Releases page](https://github.com/alhadad-xyz/devdock/releases).

## ⚡ Quick Start

Navigate to your existing application and run:

### 1. Initialize
```bash
devdock init
```
DevDock detects your framework and asks which services you want. For supported frameworks, it generates `.devdock.yml` and a managed `compose.yml`, and maps local `.env` database connection strings automatically! For existing Docker Compose projects, DevDock acts as a wrapper without modifying your Compose files.

### 2. Diagnose Your Environment
```bash
devdock doctor
```
One of our most powerful commands. DevDock runs a full diagnostic: Docker health, Compose v2 availability, Node/Go runtime checks, config validation, and port conflict detection. 

If something is wrong, DevDock gives you a clear What/Why/Fix:
```txt
✗ Port 5432 is already in use
  PID 8821 (postgres) is using this port.
  Fix: Change services.postgres.port in .devdock.yml to 5433.
```

### 3. Start the Environment
```bash
devdock up
```
This boots up your Docker containers in the background and runs your local application process in the foreground. 
*Want your terminal back? Run `devdock up --detach`.*

### 4. Check Status
```bash
devdock status
```
Get a clean overview of your app's local URL, running Docker health states, and connection credentials.

### 5. Access Web UIs
```bash
devdock open
devdock open mailpit
devdock open minio
```
Instantly opens the web interface for your app or the specified service in your default browser.

### 6. View Logs
```bash
devdock logs app
devdock logs mysql
devdock service logs postgres --tail 50
```
Stream logs from your backgrounded app process or from any of your Docker services.

### 7. Tear Down
```bash
devdock down
```
Gracefully terminates your host application process and spins down your Docker containers. Pass `--volumes` if you want to wipe your database data.

## 🧰 Managing Services

You can easily manage services via the CLI after initialization:

```bash
devdock service add redis
devdock service remove minio
devdock service status
devdock service logs postgres
```
*(Note: After adding a service, run `devdock up` to start it. After removing a service, run `devdock down && devdock up` to fully apply the change.)*

For existing Docker Compose projects, DevDock does not modify your `compose.yml`; add/remove services manually in your Compose file.

## 🏃 Running Commands

Define shortcuts in `.devdock.yml`:
```yaml
commands:
  migrate: "npx prisma migrate dev"
  test: "go test ./..."
```
Then run them instantly:
```bash
devdock run migrate
```

## ⚙️ Configuration (`.devdock.yml`)

DevDock saves its state in `.devdock.yml` at the root of your project.

```yaml
version: "1"
project:
  name: my-app
  type: express
app:
  command: npm run dev
  port: 3000
services:
  postgres:
    enabled: true
    version: "15"
    port: 5432
  mailpit:
    enabled: true
    version: "v1.21"
    smtp_port: 1025
    ui_port: 8025
```

## 🚑 Troubleshooting

- **`devdock up` hangs on "Starting app process...":** Ensure the port specified in `.devdock.yml` matches your framework's actual start port.
- **Command Not Found:** If you see `command not found: devdock`, ensure `~/go/bin` is added to your `PATH` in `~/.zshrc` or `~/.bashrc`.
- **Port Conflict Errors:** Run `devdock doctor` to see what PID is using the port. Either kill that process or edit `.devdock.yml` to map to a different port.
- **Docker Daemon Not Running:** `devdock doctor` will catch this. Start Docker Desktop (or OrbStack/Colima) and try again.
- **Missing PHP/Node/Go:** DevDock doctor will verify runtimes required by your project type.
- **App Process Logs Not Showing:** If you ran in detached mode, tail `~/.devdock/logs/<project>.app.log` or simply run `devdock logs app`.
- **MinIO `local` bucket not found:** Run `devdock down && devdock up` to recreate the automated initialization hooks.

## 🤝 Contributing

Contributions are welcome! DevDock is built in Go. Framework recipes are defined cleanly in `recipes/*.yml` making it incredibly easy to add support for new frameworks.
