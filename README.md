# gogotrace
gogotrace is a simple network inspect tool written in Go. It performs similar probes like `traceroute` on MacOS.

## Build from source
### for MacOS
```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o gogotrace *.go
```

### for Linux
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gogotrace *.go
```

## How to use
Current version uses the ipv4 ICMP socket. This requires the root privilege. So it must be run with `sudo`.

### Run with `go run`
To trace route to **www.google.com**:
```bash
sudo go run *.go www.google.com
```

### Run with built executive
Assumed the target is still **www.google.com**.
The built executive is named `gogotrace` like explained above:
```bash
sudo ./gogotrace www.google.com
```

### Supported arguments
1. Max TTL
To limit the max hop times, use `-m` to denote the value.
For example, if you only want `gogotrace` to trace 15 hops at max:
```bash
./gogotrace -m 15 www.google.com
```

2. Socket timeout
It is possible to set the socket timeout, for example, 10 seconds, each time `gogotrace` sends out probe packets.
```bash
./gogotrace -w 10s www.google.com
```

As any other friendly commands, `gogotrace` provide help at any time!
```bash
./gogotrace -h
```

## Limitations
1. Must use `sudo`
2. It only sends a packet for each hop
