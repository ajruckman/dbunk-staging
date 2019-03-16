// Copyright (c) A.J. Ruckman 2019

// Package plugin implements Dbunk's CoreDNS plugin. CoreDNS prohibits loading
// plugins from the '/internal' directory, so we place it in '/pkg' instead.
// Based on: https://github.com/coredns/example/blob/master/example.go
package plugin

import (
    "context"

    "github.com/coredns/coredns/core/dnsserver"
    "github.com/coredns/coredns/plugin"
    "github.com/mholt/caddy"
    "github.com/miekg/dns"

    "github.com/ajruckman/dbunk-staging/internal/load"
    "github.com/ajruckman/dbunk-staging/internal/serve"
)

func init() {
    caddy.RegisterPlugin("dbunk", caddy.Plugin{
        ServerType: "dns",
        Action:     setup,
    })
}

func setup(c *caddy.Controller) error {
    load.ParseConfig(c)
    load.HandleFlags()

    // Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
    dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
        return Dbunk{Next: next}
    })

    // All OK, return a nil error.
    return nil
}

type Dbunk struct {
    Next plugin.Handler
}

func (e Dbunk) Name() string {
    return "dbunk"
}

func (e Dbunk) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
    return serve.Serve(e.Next, ctx, w, r)
}
