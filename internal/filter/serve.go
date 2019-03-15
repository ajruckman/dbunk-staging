// Copyright (c) A.J. Ruckman 2019

package filter

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
)

func Serve(handler plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	spew.Dump(r)

	return plugin.NextOrFailure(handler.Name(), handler, ctx, w, r)
}