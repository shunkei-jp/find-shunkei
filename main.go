package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ebiyu/zeroconf"
)

func main() {
	fmt.Printf("IPv4 \thostname \tservice\n")

	var queryes = []string{"_shunkei_vtx_tx._tcp", "_shunkei_vtx_rx._tcp", "_momo_tx._udp"}

	var wg sync.WaitGroup
	for _, query := range queryes {
		wg.Add(1)
		go func(query string) {
			defer wg.Done()
			lookup(query)
		}(query)
	}
	wg.Wait()
}

func lookup(query string) {
	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		log.Fatalln("Failed to initialize resolver:", err.Error())
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			for _, addr := range entry.AddrIPv4 {
                fmt.Printf("%v \t%v \t%v\thttp://%v/\n", addr, entry.HostName, query, addr)
			}
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	err = resolver.BrowseWithStrategy(ctx, query, "local.", zeroconf.ForceIPv4, entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
}
