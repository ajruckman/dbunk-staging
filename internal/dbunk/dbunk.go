// Copyright (c) A.J. Ruckman 2019

package dbunk

import (
    "context"
    "fmt"
    "strings"

    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/dbunk-staging/internal/common"
    "github.com/ajruckman/dbunk-staging/internal/config"
    "github.com/ajruckman/dbunk-staging/internal/dbunk/blacklist"
    "github.com/ajruckman/dbunk-staging/internal/dbunk/respond"
    "github.com/ajruckman/dbunk-staging/internal/dbunk/whitelist"
    "github.com/ajruckman/dbunk-staging/internal/log"
)

func Serve(handler plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (rcode int, err error) {
    ip := strings.Split(w.RemoteAddr().String(), ":")[0]

    for _, v := range r.Question {
        qname := strings.TrimSuffix(v.Name, ".")

        if strings.Count(qname, ".") == 0 {
            if config.Conf.DomainNeeded {
                return domainNeeded(w, r)
            }
        }

        key := common.HostToKey(qname)

        if len(key) == 0 {
            return keyFailure(w, r)
        }

        if whitelist.Match(key, qname, ip) {
            return allow(handler, ctx, w, r)
        }

        if blacklist.Match(key, qname) {
            return block(w, r, v.Qtype)
        }
    }

    return pass(handler, ctx, w, r)
}

func domainNeeded(w dns.ResponseWriter, r *dns.Msg) (rcode int, err error) {
    fmt.Println("!!! domain needed")
    fmt.Println()

    err = respond.WithCode(w, r, dns.RcodeBadName)
    log.Error(err)

    return dns.RcodeBadName, nil
}

func keyFailure(w dns.ResponseWriter, r *dns.Msg) (rcode int, err error) {
    fmt.Println("!!! key too short")
    fmt.Println()

    err = respond.WithCode(w, r, dns.RcodeBadName)
    log.Error(err)

    return dns.RcodeBadName, nil
}

func allow(handler plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (rcode int, err error) {
    fmt.Println("!!! whitelist")
    fmt.Println()

    return plugin.NextOrFailure(handler.Name(), handler, ctx, w, r)
}

func block(w dns.ResponseWriter, r *dns.Msg, qtype uint16) (rcode int, err error) {
    fmt.Println("!!! blacklist")
    fmt.Println()

    err = respond.WithSpoofedAddr(w, r, qtype)
    log.Error(err)

    return dns.RcodeRefused, nil
}

func pass(handler plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (rcode int, err error) {
    fmt.Println("!!! pass")
    fmt.Println()

    return plugin.NextOrFailure(handler.Name(), handler, ctx, w, r)
}

