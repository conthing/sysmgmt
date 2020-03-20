package services

import (
	"context"
	"testing"
	"time"

	"github.com/grandcat/zeroconf"

	"github.com/conthing/utils/common"
	"github.com/stretchr/testify/assert"
)

// MDNS 服务
func TestMDNS(t *testing.T) {
	common.Log.Debug("Discovering MDNS...")
	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if assert.NoError(t, err) {

		entries := make(chan *zeroconf.ServiceEntry)
		go func(results <-chan *zeroconf.ServiceEntry) {
			for entry := range results {
				common.Log.Debug(entry)
			}
			common.Log.Debug("No more entries.")
		}(entries)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		err = resolver.Browse(ctx, "_workstation._tcp", "local.", entries)
		if assert.NoError(t, err) {
			<-ctx.Done()
		}

	}
}
