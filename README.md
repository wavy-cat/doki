# doki

A minimalistic and fast port knocker written in Go.

---

## Getting started

Precompiled binaries are available on the [releases page](https://github.com/wavy-cat/doki/releases).

To see all available flags, run `./doki -h`:

```
Usage of doki:
  -4    Force use IPv4 protocol
  -6    Force use IPv6 protocol
  -address string
        Target IP address
  -domain string
        Target domain name
  -ignore-errors
        Ignore errors when establishing a connection
  -ports value
        Comma-separated list of ports (0-65535 range)
  -timeout duration
        Maximum time to establish a connection (default 10ms)
```

## Building from Source

1. Download and install Go from the [official website](https://go.dev/dl/).
2. Clone the repository: `git clone https://github.com/wavy-cat/doki`
3. In the project directory, run: `go build -trimpath -ldflags="-s -w" -o doki github.com/wavy-cat/doki/cmd/main`

> [!NOTE]
> If you're using Windows, the output file should be named `doki.exe` instead of `doki`.