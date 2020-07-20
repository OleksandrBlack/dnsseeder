package dnsseed

import (
	"context"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/OleksandrBlack/dnsseeder/safecoin"
)

// SafecoinSeeder discovers IP addresses by asking Safecoin peers for them.
type SafecoinSeeder struct {
	Next   plugin.Handler
	Zones  []string
	seeder *safecoin.Seeder
	opts   *options
}

// Name satisfies the Handler interface.
func (zs SafecoinSeeder) Name() string { return "dnsseed" }

// Ready implements the ready.Readiness interface, once this flips to true CoreDNS
// assumes this plugin is ready for queries; it is not checked again.
func (zs SafecoinSeeder) Ready() bool {
	// setup() has attempted an initial connection to the backing peer already.
	return zs.seeder.Ready()
}

func (zs SafecoinSeeder) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// Check if it's a question for us
	state := request.Request{W: w, Req: r}
	zone := plugin.Zones(zs.Zones).Matches(state.Name())
	if zone == "" {
		return plugin.NextOrFailure(zs.Name(), zs.Next, ctx, w, r)
	}

	var peerIPs []net.IP
	switch state.QType() {
	case dns.TypeA:
		peerIPs = zs.seeder.Addresses(25)
	case dns.TypeAAAA:
		peerIPs = zs.seeder.AddressesV6(25)
	default:
		return dns.RcodeNotImplemented, nil
	}

	a := new(dns.Msg)
	a.SetReply(r)
	a.Authoritative = true
	a.Answer = make([]dns.RR, 0, 25)

	for i := 0; i < len(peerIPs); i++ {
		var rr dns.RR

		if peerIPs[i].To4() == nil {
			rr = new(dns.AAAA)
			rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Ttl: zs.opts.recordTTL, Class: state.QClass()}
			rr.(*dns.AAAA).AAAA = peerIPs[i]
		} else {
			rr = new(dns.A)
			rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Ttl: zs.opts.recordTTL, Class: state.QClass()}
			rr.(*dns.A).A = peerIPs[i]
		}

		a.Answer = append(a.Answer, rr)
	}

	w.WriteMsg(a)
	return dns.RcodeSuccess, nil
}
