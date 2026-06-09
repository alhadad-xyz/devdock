# DevDock

DevDock is a local development launcher that detects your project, generates a lightweight environment config, and starts your app together with Dockerized services like MySQL, PostgreSQL, and Redis.

DevDock removes repetitive local setup work: no more writing the same Compose boilerplate, manually checking ports, or remembering service commands for every project.

## ✨ Features

- **🪄 Intelligent Auto-Detection:** Automatically detects your supported framework (Laravel, Next.js, or existing docker-compose projects).
- **🚀 One-Command Spin Up:** Run your local framework process (e.g., `php artisan serve` or `npm run dev`) *and* your Docker databases concurrently.
- **🛡️ Port Conflict Detection:** Dynamically scans your host machine before starting. If another process is hogging port `3306` or `5432`, DevDock stops safely and tells you exactly what PID is causing the issue.
- **👻 Detached Mode:** Run `devdock up --detach` to background your application. DevDock safely manages the PID and tears it down gracefully when you run `devdock down`.
- **📊 Unified Observability:** Run `devdock status` for a clean dashboard of your Docker services and your host application process, complete with database connection strings.

## 🛠 Supported Stacks

| Framework | Supported Databases | Notes |
|-----------|--------------------|-------|
| **Laravel** | MySQL, PostgreSQL, Redis | Detects `composer.json` and `artisan`. Generates `.env`. |
| **Next.js** | PostgreSQL, MySQL, Redis | Detects `package.json` and `next.config.*`. |
| **Docker Compose** | *Any* | Existing `compose.yml`. DevDock acts as a friendly wrapper. |

## 📦 Installation

Currently, DevDock can be installed from source via Go:

```bash
# Clone the repository
git clone https://github.com/yourusername/devdock.git
cd devdock

# Install the binary
go install ./cmd/devdock
```

*Prerequisites: You must have [Go](https://go.dev/) 1.23+ and [Docker](https://www.docker.com/) installed on your machine. Ensure your `~/go/bin` is in your `$PATH`.*

## ⚡ Quick Start

Navigate to your existing application (e.g., a Laravel or Next.js project) and run:

### 1. Initialize
```bash
devdock init
```
DevDock will detect your framework and ask which databases you'd like to use (MySQL, Postgres, Redis). For Laravel and Next.js projects, DevDock generates `.devdock.yml` and `compose.yml`. For existing Docker Compose projects, DevDock keeps your current Compose file and acts as a friendly wrapper around it. For supported frameworks, it even sets up your local `.env` database connection strings automatically!

### 2. Diagnose Your Environment
```bash
devdock doctor
```
DevDock checks whether Docker is installed and running, Docker Compose v2 is available, required runtimes are installed, `.devdock.yml` is valid, and configured ports are available. 

If something is wrong, DevDock prints a clear fix:

```txt
✗ Port 5432 is already in use

  PID 8821 (postgres) is using this port.

  Fix: Change services.postgres.port in .devdock.yml to 5433,
       then run `devdock up` again.
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
Get a clean overview of your app's local URL, running Docker health states, and database connection strings.

### 5. View Logs
```bash
devdock logs app
devdock logs mysql
```
Stream logs from your backgrounded app process or from any of your Docker services.

### 6. Tear Down
```bash
devdock down
```
Gracefully terminates your host application process and spins down your Docker containers. Pass `--volumes` if you want to wipe your database data.

## ⚙️ Configuration (`.devdock.yml`)

DevDock saves its state in `.devdock.yml` at the root of your project. You can manually edit this file to change ports, database versions, or your framework's start command.

```yaml
version: "1"
project:
  name: my-app
  type: laravel
app:
  command: php artisan serve --host=127.0.0.1 --port=8000
  port: 8000
services:
  mysql:
    enabled: true
    version: "8.0"
    port: 3306
  redis:
    enabled: true
    version: "7"
    port: 6379
```
If you make changes, simply run `devdock up` again.

## 🚑 Troubleshooting

- **`devdock up` hangs on "Starting app process...":** Ensure the port specified in `.devdock.yml` matches your framework's actual start port.
- **Port Conflict Errors:** Run `devdock doctor` to see what PID is using the port. Either kill that process or edit `.devdock.yml` to map to a different port.
- **App Process Logs Not Showing:** If you run detached mode, tail `~/.devdock/logs/<project>.app.log` or simply run `devdock logs app`.
- **Docker Daemon Not Running:** `devdock doctor` will catch this. Start Docker Desktop or your Docker daemon and try again.

## 🤝 Contributing

Contributions are welcome! DevDock is built in Go. Framework recipes are defined cleanly in `recipes/*.yml` making it incredibly easy to add support for new frameworks like Django, Rails, or SvelteKit.
