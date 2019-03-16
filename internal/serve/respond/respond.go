// Copyright (c) A.J. Ruckman 2019

package respond

import (
    "net"

    "github.com/coredns/coredns/request"
    "github.com/miekg/dns"

    "github.com/ajruckman/dbunk-staging/internal/config"
)

// WithCode writes an rcode to the client.
func WithCode(w dns.ResponseWriter, r *dns.Msg, code int) error {
    m := new(dns.Msg)
    m.SetReply(r)
    m.Authoritative = true
    m.Compress = true
    m.RecursionAvailable = true
    m.Rcode = code

    return w.WriteMsg(m)
}

// WithSpoofedAddr writes an IP to the client.
// Code from: https://coredns.io/2017/03/01/how-to-add-plugins-to-coredns/
func WithSpoofedAddr(w dns.ResponseWriter, r *dns.Msg, qtype uint16) error {
    m := new(dns.Msg)
    m.SetReply(r)
    m.Authoritative = true
    m.Compress = true
    m.RecursionAvailable = true
    m.Rcode = dns.RcodeSuccess

    var rr dns.RR
    state := request.Request{W: w, Req: r}

    switch qtype {
    case dns.TypeA:
        rr = &dns.A{
            Hdr: dns.RR_Header{
                Name:   state.QName(),
                Rrtype: dns.TypeA,
                Class:  state.QClass(),
            },
            A: net.ParseIP(config.Conf.SpoofedA).To4(),
        }
    case dns.TypeAAAA:
        rr = &dns.AAAA{
            Hdr: dns.RR_Header{
                Name:   state.QName(),
                Rrtype: dns.TypeAAAA,
                Class:  state.QClass(),
            },
            AAAA: net.ParseIP(config.Conf.SpoofedAAAA),
        }
    case dns.TypeCNAME:
        rr = &dns.CNAME{
            Hdr: dns.RR_Header{
                Name:   state.QName(),
                Rrtype: dns.TypeCNAME,
                Class:  state.QClass(),
            },
            Target: dns.Fqdn(config.Conf.SpoofedCNAME),
        }
    default:
        rr = &dns.ANY{
            Hdr: dns.RR_Header{
                Name:   state.QName(),
                Rrtype: dns.TypeCNAME,
                Class:  state.QClass(),
            },
        }
    }
    m.Answer = append(m.Answer, rr)

    return w.WriteMsg(m)
}
