package dnsseed

import (
	"testing"
	"time"

	"github.com/caddyserver/caddy"
	"github.com/oleksandrblack/dnsseeder/safecoin/network"
)

// TestSetup tests the various things that should be parsed by setup.
func TestSetup(t *testing.T) {
	tt := []struct {
		config string
		validConfig bool
		magic network.Network
		interval time.Duration
		bootstrap []string
		ttl uint32
	}{
		{`dnsseed`, false, 0, 0, []string{}, 0},
		{`dnsseed mainnet`, false, 0, 0, []string{}, 0},
		{`dnsseed { }`, false, 0, 0, []string{}, 0},
		{`dnsseed { network }`, false, 0, 0, []string{}, 0},
		{`dnsseed { network mainnet }`, true, network.Mainnet, defaultUpdateInterval, []string{}, defaultTTL},
		{`dnsseed { 
			network testnet
			crawl_interval 15s
			bootstrap_peers
			}`, 
			false, 0, 0, []string{}, 0,
		},
		{`dnsseed { 
			network testnet
			crawl_interval
			bootstrap_peers 127.0.0.1:8770
			}`, 
			false, 0, 0, []string{}, 0,
		},
		{`dnsseed { 
			network testnet
			crawl_interval 15s
			bootstrap_peers 127.0.0.1:8770
			}`, 
			true, network.Testnet, time.Duration(15) * time.Second, []string{"127.0.0.1:8770"}, defaultTTL,
		},
		{`dnsseed { 
			network testnet
			crawl_interval 15s
			bootstrap_peers 127.0.0.1:8770
			boop snoot every 15s
			}`, 
			false, 0, 0, []string{}, 0,
		},
		{`dnsseed { 
			network mainnet
			crawl_interval 30m
			bootstrap_peers 127.0.0.1:8770 127.0.0.2:8770
			record_ttl 300
			}`, 
			true, network.Mainnet, time.Duration(30) * time.Minute, []string{"127.0.0.1:8770", "127.0.0.2:8770"}, 300,
		},
	}

	for _, test := range tt {
		c := caddy.NewTestController("dns", test.config)
		opts, err := parseConfig(c)
		if (err == nil) != test.validConfig {
			t.Errorf("Unexpected error in test case `%s`: %v", test.config, err)
			t.FailNow()
		}

		if err != nil && !test.validConfig {
			// bad parse, as expected
			continue
		}
		
		if opts.networkMagic != test.magic {
			t.Errorf("Input: %s wrong network magic", test.config)
		}

		if opts.updateInterval != test.interval {
			t.Errorf("Input: %s wrong update interval", test.config)
		}

		for i, s := range opts.bootstrapPeers {
			if s != test.bootstrap[i] {
				t.Errorf("Input: %s wrong bootstrap peer", test.config)
			}
		}

		if opts.recordTTL != test.ttl {
			t.Errorf("Input: %s wrong TTL", test.config)
		}
	}
}