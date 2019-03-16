// Copyright (c) A.J. Ruckman 2019

package serve

import (
    "context"
    "fmt"
    "strings"

    "github.com/coredns/coredns/plugin"
    "github.com/miekg/dns"

    "github.com/ajruckman/dbunk-staging/internal/common"
    "github.com/ajruckman/dbunk-staging/internal/config"
    "github.com/ajruckman/dbunk-staging/internal/log"
    "github.com/ajruckman/dbunk-staging/internal/serve/blacklist"
    "github.com/ajruckman/dbunk-staging/internal/serve/respond"
    "github.com/ajruckman/dbunk-staging/internal/serve/whitelist"
)

func Serve(handler plugin.Handler, ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (rcode int, err error) {
    ip := strings.Split(w.RemoteAddr().String(), ":")[0]

    for _, v := range r.Question {
       qname := strings.TrimSuffix(v.Name, ".")
       _ = qname

       if strings.Count(qname, ".") == 0 {
          if config.Conf.DomainNeeded {
              fmt.Println("!!! domain needed")
              fmt.Println()

              err = respond.WithCode(w, r, dns.RcodeBadName)
              log.Error(err)

              return dns.RcodeBadName, nil
          }
       }

       key := common.HostToKey(qname)
       if len(key) == 0 {
          fmt.Println("!!! key too short")
          fmt.Println()

          err = respond.WithCode(w, r, dns.RcodeBadName)
          log.Error(err)

          return dns.RcodeBadName, nil
       }

       if whitelist.Match(key, qname, ip) {
          fmt.Println("!!! whitelist")
          fmt.Println()

          return plugin.NextOrFailure(handler.Name(), handler, ctx, w, r)
       }

       if blacklist.Match(key, qname) {
          fmt.Println("!!! blacklist")
          fmt.Println()

          err = respond.WithSpoofedAddr(w, r, v.Qtype)
          log.Error(err)

          return dns.RcodeRefused, nil
       }
    }

    fmt.Println("!!! pass")
    fmt.Println()

    return plugin.NextOrFailure(handler.Name(), handler, ctx, w, r)
}
