// Copyright (c) A.J. Ruckman 2019

package respond

import "github.com/miekg/dns"

func WithCode(w dns.ResponseWriter, r *dns.Msg, code int) error {
    m := new(dns.Msg)
    m.SetReply(r)
    m.Authoritative = true
    m.Compress = true
    m.RecursionAvailable = true
    m.Rcode = code

    return w.WriteMsg(m)
}
