# OmniConfig

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/dl/)

> **One CLI tool to read, write, and manage configuration files across *any* operating system — including iOS and Android.**

Auto-detects the OS and config format so you never have to think about it — just run one command, anywhere.

## Features

- **Cross-platform** — Linux, macOS, Windows, iOS, Android (amd64 & arm64)
- **Multi-format** — JSON, YAML, TOML, INI, XML, .env, plist, and more
- **Auto-detection** — Detects your OS and config file format automatically
- **Single binary** — No dependencies, static builds
- **Mobile aware** — Understands iOS and Android config locations

## Quick Start

```bash
# Read a config value
omniconfig get database.host

# Write a config value
omniconfig set server.port 8080

# Detect OS and config format
omniconfig detect

# Get help
omniconfig --help
```

## Install

### Linux / macOS / Android (Termux)
```bash
curl -fsSL https://omnicofig.sh/install | sh
```
Or with wget:
```bash
wget -qO- https://omnicofig.sh/install | sh
```

### Windows (PowerShell)
```powershell
irm https://omnicofig.sh/install.ps1 | iex
```

### Build from source
```bash
git clone https://github.com/omnicofig/cli.git
cd cli
make build
```

## Commands

| Command | Description |
|---------|-------------|
| `get <key>` | Read a configuration value |
| `set <key> <value>` | Write a configuration value (creates backup) |
| `detect` | Detect OS and config file format |
| `version` | Print version information |
| `--help` | Show help for any command |
| `--config` | Specify a config file path |
| `--list-platforms` | List all supported target platforms |

## Supported Platforms

| Platform | amd64 | arm64 | Notes |
|----------|-------|-------|-------|
| Linux    | ✅    | ✅    | |
| macOS    | ✅    | ✅    | Intel & Apple Silicon |
| Windows  | ✅    | ✅    | |
| iOS      | —     | ✅    | Needs Xcode |
| Android  | ✅    | ✅    | Needs NDK / Termux |

## Supported Formats

| Format | Extensions | Status |
|--------|-----------|--------|
| JSON | `.json` | ✅ Implemented |
| YAML | `.yml`, `.yaml` | ✅ Implemented |
| TOML | `.toml` | ✅ Implemented |
| INI | `.ini`, `.cfg` | ✅ Implemented |
| .env | `.env` | ✅ Implemented |
| XML | `.xml` | ✅ Implemented |
| plist | `.plist` | 🔜 Coming |

## Building from Source

```bash
# Build for your current platform
make build

# Build for all desktop platforms
make build-desktop

# Build for mobile (iOS/Android) — requires Xcode/NDK
make build-mobile

# Build everything
make build-all

# Run tests
make test
```

## License

MIT — see [LICENSE](LICENSE).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.