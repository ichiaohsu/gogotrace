package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type Config struct {
	source  string        // Source string, in either Domain Name string or IP string
	maxTTL  int           // Max hop limit
	timeout time.Duration // Socket will timeout after denoted period, default 5s
	ipv6    bool          // Set to ipv6 mode
}

// usage overwrite the default help message coming with flag
func usage() {
	fmt.Fprintf(os.Stderr,
		`Usage: enroute [-m maxTTL] [-w waittime] [-h] [-6] host
Options:
`)
	flag.PrintDefaults()
}

func main() {

	var (
		showHelp bool
		config   Config
	)
	flag.BoolVar(&showHelp, "h", false, "help for traceroute")
	flag.IntVar(&config.maxTTL, "m", 30, "Set the max time-to-live (max number of hops) used in outgoing probe packets")
	flag.DurationVar(&config.timeout, "w", time.Second*5, "Set the time (e.g. 10s, 3m, 1h) to wait for a response to a probe")
	flag.BoolVar(&config.ipv6, "6", false, "Set this up when designated to use ipv6")

	flag.Usage = usage

	flag.Parse()

	// Show help message if -h argument are set or there is no host info
	if flag.NArg() == 0 || showHelp {
		flag.Usage()
		return
	}

	config.source = flag.Arg(0)

	var (
		t   trace
		err error
	)
	if !config.ipv6 {
		t, err = NewTraceV4(config)
	} else {
		t, err = NewTraceV6(config)
	}
	if err != nil {
		fmt.Println(err)
	}
	// Send probes and get all the results
	results, err := t.Send()
	if err != nil {
		fmt.Println(err)
	}
	// Show results statistics
	results.MaxHopTime()
}
