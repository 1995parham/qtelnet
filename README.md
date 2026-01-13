<h1 align="center"> qtelnet </h1>
<h6 align="center"> Your telnet in the future </h6>

## Introduction

QUIC is an application layer protocol based on UDP which tries to simulate TCP features in the application layer.
`qtelnet` works like telnet but for QUIC. You can connect to a QUIC server and send or receive information bidirectionally.

## Installation

```bash
go install github.com/1995parham/qtelnet/cmd/client@latest
```

Or build from source:

```bash
git clone https://github.com/1995parham/qtelnet.git
cd qtelnet
go build -o qtelnet ./cmd/client
```

## Usage

```bash
qtelnet [flags] <address> <port>
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--alpn` | `-a` | `quic-echo` | ALPN protocol to use |
| `--insecure` | `-K` | `true` | Skip certificate verification |
| `--help` | `-h` | | Show help |

### Examples

Connect to a local QUIC server:

```bash
./qtelnet 127.0.0.1 8080
```

Connect to an HTTP/3 server (connection only, see [Notes](#notes)):

```bash
./qtelnet -a h3 cloudflare-quic.com 443
```

## Testing

The project includes a simple echo server for testing:

```bash
# Build both binaries
go build -o qtelnet ./cmd/client
go build -o qtelnet-server ./cmd/server

# Terminal 1: Start the server
./qtelnet-server

# Terminal 2: Connect with the client
./qtelnet 127.0.0.1 8080
```

Type messages in the client terminal, and the server will echo them back.

## Project Structure

```
qtelnet/
├── cmd/
│   ├── client/          # Client entry point
│   │   └── main.go
│   └── server/          # Test echo server
│       └── main.go
├── internal/
│   ├── cmd/             # CLI command handling
│   │   └── root.go
│   └── handler/         # QUIC stream handlers
│       ├── accept.go    # Incoming stream handler
│       └── prompt.go    # User input handler
├── go.mod
└── README.md
```

## How It Works

- **Accepter**: Runs in a goroutine, listening for incoming QUIC streams from the server and printing received data to stdout
- **Prompt**: Reads user input from stdin, opens a new QUIC stream for each line, and sends data to the server

## Notes

- **HTTP/3 servers won't work** with qtelnet for making HTTP requests. HTTP/3 uses binary framing (like HTTP/2), not plain text like HTTP/1.1. You can establish a QUIC connection, but you can't send readable HTTP requests.
- For raw telnet-like communication, use the included test server or deploy your own QUIC server that expects plain text over QUIC streams.
- The `-a` flag lets you specify different ALPN protocols to connect to various QUIC servers, but meaningful communication requires a server that understands plain text.
