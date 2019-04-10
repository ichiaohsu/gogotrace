package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type trace interface {
	Send() (results reports, err error)
}

type TraceV4 struct {
	dst *net.IPAddr // Destination IP Address
	// dst     net.Addr // Destination IP Address
	network string // Network listen proto
	address string
	Config  // Set identical argement parsed from CLI
}

func createICMP(typ icmp.Type, seq int) (b []byte, err error) {

	m := icmp.Message{
		Type: typ,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	b, err = m.Marshal(nil)
	return b, err
}

// NewTrace return a new Trace Object
func NewTraceV4(config Config) (*TraceV4, error) {

	// Parse ip address
	// Accept both ip string or domain names
	ipAddr, err := net.ResolveIPAddr("ip4", config.source)
	if err != nil {
		return nil, err
	}

	fmt.Printf("tracerout to %s(%v), %d hop max\n", config.source, ipAddr, config.maxTTL)

	// var dst net.Addr = ipAddr
	t := &TraceV4{
		network: "ip4:icmp",
		address: "0.0.0.0",
		dst:     ipAddr,
		Config:  config,
	}
	return t, nil
}

func (t *TraceV4) Send() (results reports, err error) {

	// Listen to ICMP packet on ip4
	// conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	conn, err := net.ListenPacket(t.network, t.address)
	// conn, err := icmp.ListenPacket(t.proto, "0.0.0.0")
	if err != nil {
		return results, err
	}
	defer conn.Close()

	// Create ipv4 PacketConn to enable setTTL
	p := ipv4.NewPacketConn(conn)

	// received bytes
	rb := make([]byte, 1500)
	// Start the loop
TraceLoop:
	for i := 1; i < t.maxTTL; i++ {

		// hop is the result of single probe
		var hop Hop
		// Create ICMP bytes message
		b, err := createICMP(ipv4.ICMPTypeEcho, i)
		if err != nil {
			return results, err
		}
		// Set ICMP TTL info
		if err := p.SetTTL(i); err != nil {
			return results, err
		}

		begin := time.Now()
		// Start to send
		if _, err := p.WriteTo(b, nil, t.dst); err != nil {
			return results, err
		}
		if err := p.SetReadDeadline(time.Now().Add(t.timeout)); err != nil {
			return results, err
		}

		n, _, peer, err := p.ReadFrom(rb)
		if err != nil {
			hop = Hop{ID: i, Addr: "*"}
			hop.formatPrint()
			results = append(results, hop)
			// Continue next loop
			continue TraceLoop
		}
		// Parse receviced message
		// Protocol-numbers: ICMP = 1, IPv6-ICMP = 58
		rm, err := icmp.ParseMessage(1, rb[:n])
		if err != nil {
			return results, err
		}
		// Round Trip Time
		rtt := time.Since(begin)

		// Calculate the time spent between consecutive hop
		// If it's the first probe, it's the first round trip time
		// Else it is the rtt - previous rtt
		var hopTime time.Duration
		if len(results) == 0 {
			hopTime = rtt
		} else {
			hopTime = rtt - results[len(results)-1].RTT
		}

		// Parsing received ICMP message
		switch rm.Type {
		case ipv4.ICMPTypeTimeExceeded:
			hop = Hop{ID: i, Addr: peer.String(), RTT: rtt, HopTime: hopTime}
			hop.formatPrint()
		case ipv4.ICMPTypeEchoReply:
			hop = Hop{ID: i, Addr: peer.String(), RTT: rtt, HopTime: hopTime}
			hop.formatPrint()
			break TraceLoop
		default:
			hop = Hop{ID: i, Addr: "*"}
			hop.formatPrint()
		}
		results = append(results, hop)
	}
	return results, nil
}
