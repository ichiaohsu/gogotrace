package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv6"
)

type TraceV6 struct {
	dst *net.IPAddr // Destination IP Address
	// dst     net.Addr // Destination IP Address
	network string // Network listen proto
	address string
	Config  // Set identical argement parsed from CLI
}

func NewTraceV6(config Config) (*TraceV6, error) {
	t := &TraceV6{
		network: "ip6:58",
		address: "::",
		Config:  config,
	}
	ipAddr, err := net.ResolveIPAddr("ip6", config.source)
	fmt.Printf("tracerout to %s(%v), %d hop max\n", config.source, ipAddr, config.maxTTL)
	if err != nil {
		return nil, err
	}
	t.dst = ipAddr
	return t, nil
}

func (t *TraceV6) Send() (results reports, err error) {

	conn, err := net.ListenPacket(t.network, t.address) // ICMP for IPv6
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	p := ipv6.NewPacketConn(conn)

	if err := p.SetControlMessage(ipv6.FlagHopLimit|ipv6.FlagSrc|ipv6.FlagDst|ipv6.FlagInterface, true); err != nil {
		log.Fatal(err)
	}
	var f ipv6.ICMPFilter
	f.SetAll(true)
	f.Accept(ipv6.ICMPTypeTimeExceeded)
	f.Accept(ipv6.ICMPTypeEchoReply)
	if err := p.SetICMPFilter(&f); err != nil {
		log.Fatal(err)
	}

	var wcm ipv6.ControlMessage
	// received bytes
	rb := make([]byte, 1500)

TraceLoop:
	for i := 1; i <= t.maxTTL; i++ {

		var hop Hop

		b, err := createICMP(ipv6.ICMPTypeEchoRequest, i)
		if err != nil {
			return results, err
		}

		begin := time.Now()
		// ipv6 need to setup TTL through control message
		wcm.HopLimit = i

		if _, err := p.WriteTo(b, &wcm, t.dst); err != nil {
			return results, err
		}
		if err := p.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
			return results, err
		}

		n, _, peer, err := p.ReadFrom(rb)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				hop = Hop{ID: i, Addr: "*"}
				hop.formatPrint()
				results = append(results, hop)
				continue TraceLoop
			}
			log.Fatal(err)
		}
		rm, err := icmp.ParseMessage(58, rb[:n])
		if err != nil {
			log.Fatal(err)
		}
		rtt := time.Since(begin)

		var hopTime time.Duration
		if len(results) == 0 {
			hopTime = rtt
		} else {
			hopTime = rtt - results[len(results)-1].RTT
		}

		switch rm.Type {
		case ipv6.ICMPTypeTimeExceeded:
			hop = Hop{ID: i, Addr: peer.String(), RTT: rtt, HopTime: hopTime}
			hop.formatPrint()
		case ipv6.ICMPTypeEchoReply:
			hop = Hop{ID: i, Addr: peer.String(), RTT: rtt, HopTime: hopTime}
			hop.formatPrint()
			break TraceLoop
		}
		results = append(results, hop)
	}
	return results, nil
}
