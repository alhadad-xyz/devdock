# DevDock

**DevDock** is a zero-config, intelligent CLI tool designed to magically spin up local development environments. It automatically detects your framework (Laravel, Next.js, etc.) and seamlessly orchestrates your local application process alongside Dockerized backing services like MySQL, PostgreSQL, and Redis.

Say goodbye to manual `docker-compose.yml` boilerplate and port conflict headaches. 

## ✨ Features

- **🪄 Intelligent Auto-Detection:** Automatically detects your framework (Laravel, Next.js, or existing docker-compose projects).
- **🚀 One-Command Spin Up:** Run your local framework process (e.g., `php artisan serve` or `npm run dev`) *and* your Docker databases concurrently.
- **🛡️ Port Conflict Detection:** Dynamically scans your host machine before starting. If another process is hogging port `3306` or `5432`, DevDock stops safely and tells you exactly what PID is causing the issue.
- **👻 Detached Mode:** Run `devdock up --detach` to background your application. DevDock safely manages the PID and tears it down gracefully when you run `devdock down`.
- **📊 Unified Observability:** Run `devdock status` for a clean dashboard of your Docker services and your host application process, complete with database connection strings.

## 📦 Installation

Currently, DevDock can be installed from source via Go:

```bash
# Clone the repository
git clone https://github.com/yourusername/devdock.git
cd devdock

# Build and install the binary
go build -o devdock cmd/devdock/main.go
sudo mv devdock /usr/local/bin/
```

*Prerequisites: You must have [Go](https://go.dev/) 1.23+ and [Docker](https://www.docker.com/) installed on your machine.*

## ⚡ Quick Start

Navigate to your existing application (e.g., a Laravel or Next.js project) and run:

### 1. Initialize
```bash
devdock init
```
DevDock will detect your framework and ask which databases you'd like to use (MySQL, Postgres, Redis). It will generate a lightweight `.devdock.yml` configuration and write out a transparent `compose.yml`. For supported frameworks, it even sets up your local `.env` database connection strings automatically!

### 2. Start the Environment
```bash
devdock up
```
This boots up your Docker containers in the background and runs your local application process in the foreground. 
*Want your terminal back? Run `devdock up --detach`.*

### 3. Check Status
```bash
devdock status
```
Get a clean overview of your app's local URL, running Docker health states, and database connection strings.

### 4. View Logs
```bash
devdock logs app
devdock logs mysql
```
Stream logs from your backgrounded app process or from any of your Docker services.

### 5. Tear Down
```bash
devdock down
```
Gracefully terminates your host application process and spins down your Docker containers. Pass `--volumes` if you want to wipe your database data.

## 🛠️ Configuration (`.devdock.yml`)

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

## 🤝 Contributing

Contributions are welcome! DevDock is built in Go. Framework recipes are defined cleanly in `recipes/*.yml` making it incredibly easy to add support for new frameworks like Django, Rails, or SvelteKit.
