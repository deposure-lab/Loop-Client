# AGG Loop Client - Technical Documentation

## 1. Overview

**AGG Loop Client** is a high-performance network agent engineered in Go, designed to establish secure tunnels within a Zero Trust Network Access (ZTNA) architecture.

The agent enables secure mapping of local network services to the AGG Labs Edge infrastructure using **WebSocket Secure (WSS)** and native HTTP request processing.

The client is optimized for:

* minimal resource footprint
* maximum throughput
* native binary execution (no interpreted runtime)

---

## 2. Architecture and Performance

### 🔹 Static Compilation

The agent is compiled into a single standalone binary, eliminating external runtime dependencies.

### 🔹 Concurrency (Goroutines)

Each incoming request is handled within a lightweight goroutine, allowing:

* thousands of concurrent connections
* minimal overhead

### 🔹 Dynamic Routing

The agent automatically connects to the datacenter with the lowest latency to optimize packet routing.

---

## 3. Pre-compiled Binaries (Releases)

Pre-built binaries are available on the official GitHub repository.

### 📥 Accessing Releases

1. Navigate to the **Releases** section (right sidebar on GitHub)
2. Select the latest version (**marked as "Latest"**)
3. Download the appropriate binary:

| Platform                     | File                        |
| ---------------------------- | --------------------------- |
| Linux (x64)                  | `aggloop-linux-amd64`       |
| Linux (ARM64 / Raspberry Pi) | `aggloop-linux-arm64`       |
| Windows                      | `aggloop-windows-amd64.exe` |
| macOS (Apple Silicon)        | `aggloop-macos-arm64`       |

---

## 4. Installation and Build

### ⚙️ Requirements

* Go **1.21+**

### 🔨 Build from Source

```bash
go mod tidy
go build -o aggloop main.go
```

### 🔁 Cross-Compilation

#### Linux (x64 / ARM64)

```bash
GOOS=linux GOARCH=amd64 go build -o aggloop
```

#### Windows

```bash
GOOS=windows GOARCH=amd64 go build -o aggloop.exe
```

#### macOS (Apple Silicon)

```bash
GOOS=darwin GOARCH=arm64 go build -o aggloop
```

---

## 5. Configuration

### 📍 Default Config Location

| System        | Path                    |
| ------------- | ----------------------- |
| Linux / macOS | `/etc/agglabs/loop.yml` |
| Windows       | `C:\agglabs\loop.yml`   |

### 🧾 YAML Structure

```yaml
version: "3"
token: "YOUR_AUTH_TOKEN"
region: "default"

applications:
  app-name:
    appId: "APPLICATION_UUID"
    addr: "8080"
    scheme: "http"
    inspect: true
```

---

## 6. Command Line Interface (CLI)

### 🔐 Authentication

```bash
aggloop add-token <ACCESS_TOKEN>
```

### 🧙 Service Initialization

```bash
aggloop create
```

### 🚀 Tunnel Operations

```bash
aggloop start <app-name>
```

Runs in foreground with real-time traffic inspection.

---

## 7. Process Management (Daemon Mode)

Supported on:

* Linux
* macOS

❌ Not supported on Windows

### ▶️ Run in Background

```bash
sudo aggloop start <app-name> -b
```

### ⚙️ Flags

| Flag               | Description                               |
| ------------------ | ----------------------------------------- |
| `-b, --background` | Run process in background and save PID    |
| `-r, --auto-retry` | Auto reconnect on failure (default: true) |

### 🛑 Stop Process

```bash
sudo aggloop stop <app-name>
```

---

## 8. Security Specifications

### 🔒 Network Isolation

* Outbound-only TLS connections
* No inbound firewall ports required

### 🛡️ Least Privilege Principle

* Root required only for:

  * PID file handling
  * system config access

### 🔐 Encryption

* TLS 1.3
* WebSocket Secure (WSS)

All traffic within the tunnel is fully encrypted.

---

## 📌 Summary

AGG Loop Client provides:

* secure ZTNA tunneling
* high concurrency via Go
* lightweight deployment
* strong encryption standards
