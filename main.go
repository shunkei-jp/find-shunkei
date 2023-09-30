package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/grandcat/zeroconf"
)

func main() {
	fmt.Printf("IPv4 \thostname \tservice\n")
	lookup("_shunkei_rx._udp")
	lookup("_shunkei_tx._udp")

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
				fmt.Printf("%v \t%v \t%v\n", addr, entry.HostName, query)
			}
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = resolver.BrowseWithStrategy(ctx, query, "local.", zeroconf.ForceIPv4, entries)
	if err != nil {
		log.Fatalln("Failed to browse:", err.Error())
	}

	<-ctx.Done()
}
