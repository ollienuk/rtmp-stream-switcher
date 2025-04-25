# RTMP Stream Switcher

A lightweight Go program that acts as an RTMP server and forwards an incoming stream to an external RTMP destination. It supports primary and backup stream sources, automatically switching to the backup if the primary disconnects, and back to primary when it reconnects.

## Features

- Accepts RTMP streams via custom stream keys  
- Forwards incoming stream to a configured output RTMP URL  
- Automatically switches between a primary and backup stream source  
- Uses the [`joy4`](https://github.com/nareix/joy4) library for RTMP handling  

## How It Works

- The server listens on port `:1935` for incoming RTMP connections.  
- It checks the stream key against `FirstStreamKey` and `SecondStreamKey` from the config.  
- Only one stream (primary or backup) can be active at a time:  
  - If the primary connects, the backup is dropped.  
  - If the backup connects and the primary is inactive, it is forwarded.  
- Streams are pushed to the `OutputStreamURL`.  

## Configuration

The program uses a `config.json` file to load settings. Here's a sample format:

```json
{
  "FirstStreamKey": "/primary",
  "SecondStreamKey": "/backup",
  "OutputStreamURL": "rtmp://example.com/live/stream"
}

````
 Installation via Makefile

This project includes a `Makefile` that automates the installation of the RTMP Stream Switcher as a `systemd` service on Linux.

## 📦 Prerequisites

- Go (installed and added to your `PATH`)
- `make` utility
- Linux system with `systemd`
- Root or sudo privileges

---

## 🚀 Quick Start

To build and install the RTMP Stream Switcher as a system service:

```bash
make
```

This will:

1. Build the Go application.
2. Create the `/opt/switch` install directory.
3. Copy the binary and `config.json`.
4. Generate the `switch.service` systemd file.
5. Enable and start the systemd service.

---

## 🔧 Makefile Targets

### `make`

Builds the application and performs the full install process.

### `make build`

Compiles `main.go` into a binary named `switch`.

### `make copy-files`

Creates the install directory and copies the binary + config file.

### `make service`

Generates the systemd unit file at `/etc/systemd/system/switch.service`.

### `make start`

Reloads systemd, enables the service, and starts it.

### `make stop`

Stops the service.

### `make status`

Shows the current status of the systemd service.

### `make clean`

Removes the local `switch` binary from the project directory.

### `make uninstall`

Stops the service, disables it, removes the systemd file, and deletes the `/opt/switch` install directory.

---

## 📝 Example Usage

```bash
# Install and start the service
make

# Check service status
make status

# Stop the service
make stop

# Clean up the binary
make clean

# Uninstall everything
make uninstall
```

---

## 📁 File Structure

```
.
├── main.go
├── config.json
├── Makefile
└── README.md
```

---

## ❗ Notes

- The service runs as `root`. Adjust the `User` and `Group` in the Makefile or systemd unit if needed.
- Make sure `config.json` exists before running `make`.

---

## ✅ Done!

Your RTMP Stream Switcher is now running as a background service and will auto-start on boot!
