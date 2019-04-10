package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// usage overwrite the default help message coming with flag
func usage() {
	fmt.Fprintf(os.Stderr,
		`Usage: enroute [-m maxTTL] [-w waittime] [-h] host
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

	flag.Usage = usage

	flag.Parse()

	// Show help message if -h argument are set or there is no host info
	if flag.NArg() == 0 || showHelp {
		flag.Usage()
		return
	}

	config.source = flag.Arg(0)

	t, err := NewTrace(config)
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
