package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/ebiyu/zeroconf"
)

func main() {
	timeout := flag.Int("t", 1, "timeout")
	flag.Parse()

	fmt.Printf("IPv4 \thostname \tservice\n")

	var queryes = []string{"_shunkei_vtx_tx._tcp", "_shunkei_vtx_rx._tcp"}

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

	select {
	case <-done:
		fmt.Println("done")
	case result := <-resultsChan:
		fmt.Printf("%v \t%v \t%v\thttp://%v/\n", result.ipv4, result.hostname, result.service, result.ipv4)
	}

	<-done
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
