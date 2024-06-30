package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ebiyu/zeroconf"
)

func main() {
	timeout := flag.Int("t", 2, "timeout")
	receiver := flag.Bool("rx", false, "Print only receiver")
	transmitter := flag.Bool("tx", false, "Print only transmitter")
	first := flag.Bool("1", false, "Exit after first device found")
	ipOnly := flag.Bool("ip-only", false, "Print only IP address")
	printHost := flag.Bool("host", false, "Print hostname")
	flag.Parse()

	// In first-exit mode, timeout may be at least 5 seconds
	if *first && *timeout < 5 {
		*timeout = 5
	}

	var queryes = []string{}

	var lite = false
	if !*receiver && !*transmitter && !lite {
		*receiver = true
		*transmitter = true
		lite = true
	}

	if *receiver {
		queryes = append(queryes, "_shunkei_vtx_rx._tcp")
	}
	if *transmitter {
		queryes = append(queryes, "_shunkei_vtx_tx._tcp")
	}
	if lite {
		queryes = append(queryes, "_shunkei_vtxlite_tx._tcp")
	}

	resultsChan := make(chan LookupResult)

	var wg sync.WaitGroup
	for _, query := range queryes {
		wg.Add(1)
		go func(query string) {
			defer wg.Done()
			err := Lookup(resultsChan, query, *timeout)
			if err != nil {
				log.Fatalln("Failed to lookup:", err.Error())
			}
		}(query)
	}
	done := make(chan bool)
	go func() {
		wg.Wait()
		close(done)
	}()

	found := 0

	for {

		select {
		case <-done:
			if found == 0 {
				fmt.Fprintf(os.Stderr, "No device found\n")
				os.Exit(9)
			} else {
				fmt.Fprintf(os.Stderr, "Found %d device(s)\n", found)
				os.Exit(0)
			}
		case result := <-resultsChan:
			// print header
			if found == 0 && !*ipOnly {
				if *printHost {
					fmt.Printf("IPv4 Address \tHostname \tDevice Type \tWeb UI\n")
				} else {
					fmt.Printf("IPv4 Address \tDevice Type \tWeb UI\n")
				}
			}

			deviceType := "Unknown"
			switch result.service {
			case "_shunkei_vtx_rx._tcp":
				deviceType = "Shunkei VTX Receiver"
			case "_shunkei_vtx_tx._tcp":
				deviceType = "Shunkei VTX Transmitter"
			case "_shunkei_vtxlite_tx._tcp":
				deviceType = "Shunkei VTX Lite Transmitter"
			}

			if *ipOnly {
				if *first {
					fmt.Printf("%v", result.ipv4)
				} else {
					fmt.Printf("%v\n", result.ipv4)
				}
			} else {
				if *printHost {
					fmt.Printf("%v \t%v \t%v\thttp://%v/\n", result.ipv4, result.hostname, deviceType, result.ipv4)
				} else {
					fmt.Printf("%v \t%v\thttp://%v/\n", result.ipv4, deviceType, result.ipv4)
				}
			}

			found++

			if *first {
				os.Exit(0)
			}
		}
	}
}

type LookupResult struct {
	ipv4     net.IP
	hostname string
	service  string
}

func Lookup(resultsChan chan<- LookupResult, query string, timeout int) error {
	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			for _, addr := range entry.AddrIPv4 {
				resultsChan <- LookupResult{ipv4: addr, hostname: entry.HostName, service: query}
			}
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	err = resolver.BrowseWithStrategy(ctx, query, "local.", zeroconf.ForceIPv4, entries)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}
