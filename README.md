# Rogo 🚀

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.20-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**A Redis-compatible, in-memory key-value store written in pure Go**

Rogo is a lightweight, high-performance key-value database that implements Redis-like commands with Go's efficiency and concurrency model. Perfect for embedding in Go applications or learning database internals.

## Features

- 📌 Redis-compatible commands (`SET`, `GET`, `DEL`, `EXPIRE`, etc.)
- ⚡ Blazing fast in-memory storage
- 🧵 Thread-safe with goroutine support
- 💾 Optional persistence (AOF/RDB style)
- 📡 Simple TCP server interface
- 🚫 Zero external dependencies

## Quick Start

```bash
# Install and run
go install github.com/yourname/rogo@latest
rogo --port 6380
```

Connect using any Redis client:

```bash
redis-cli -p 6380
> SET foo bar
OK
> GET foo
"bar"
```

## Benmarks

## Dockerizing Rogo

```bash
docker build -t rogo .
docker run -p 6380:6380 rogo
```

## Roadmap

- [ ] Add persistence (AOF/RDB)
- [ ] Add more Redis-like commands
- [ ] Add more tests
- [ ] Add more documentation

## Contributing

PRs welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT