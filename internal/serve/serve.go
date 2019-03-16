// Copyright (c) A.J. Ruckman 2019

package filter

import (
    "context"
    "fmt"
    "strings"

    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/dbunk-staging/internal/common"
    "github.com/ajruckman/dbunk-staging/internal/config"
    "github.com/ajruckman/dbunk-staging/internal/filter/respond"
    "github.com/ajruckman/dbunk-staging/internal/log"
)

func Serve(handler plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
    var (
        err error
    )

    for _, v := range r.Question {
        qname := strings.TrimSuffix(v.Name, ".")
        fmt.Println(qname)

        if strings.Count(qname, ".") == 0 {
            if config.Conf.DomainNeeded {
                log.Info("")

                err = respond.WithCode(w, r, dns.RcodeBadName)
                log.Error(err)

                return dns.RcodeBadName, nil
            }

            key := common.HostToKey(qname)
            if len(key) == 0 {
                log.Warning("Derived key for " + qname + " is blank")
            }
        }
    }

    return plugin.NextOrFailure(handler.Name(), handler, ctx, w, r)
}

func handleZero
