# gogotrace
gogotrace is a simple network inspect tool written in Go. It exploited the power of Internet Control Message Protocol(ICMP) to performs similar probes like `traceroute` on MacOS and Linux.

## Prerequisites
gogotrace was developed on MacOS 10.13 and Golang 1.12. The `GO Modules` function was enabled on this machine. So installing a Go with module support will be recommended. The requirements will be as follows:

- Golang 1.11+
- Set ENV `GO111MODULE=on`
## Build from source
It's easy to build cross-platform Go executives. Here provided corresponding commands to build 64bits executive in both platform:
### for MacOS
```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o gogotrace *.go
```

### for Linux
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gogotrace *.go
```

## How to use
Current version uses the IP level socket. This requires the root privilege to query. Therefore the use method, whether with `go run` or executive, should always be done with `sudo`.

All the example below will use **www.google.com** as host.

### Run with `go run`
To trace route to host simply:
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
#### Max TTL
To limit the max hops and stop packets from wandering around, use `-m` to denote the value.
For example, if you only want `gogotrace` to trace 15 hops at max:
```bash
./gogotrace -m 15 www.google.com
```

#### Socket timeout
It is possible to set the socket timeout, for example, 10 seconds, when `gogotrace` sends out probe packets. It will wait for remote to respond for 10s each time.
```bash
./gogotrace -w 10s www.google.com
```

#### Help
As any other friendly commands, `gogotrace` provide help at any time!
```bash
./gogotrace -h
```

## Limitations
1. Must use `sudo`
2. It only sends a packet for each hop