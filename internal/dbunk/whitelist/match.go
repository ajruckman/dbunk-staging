// Copyright (c) A.J. Ruckman 2019

package whitelist

import (
    "github.com/ajruckman/dbunk-staging/internal/dbunkdb"
    "github.com/ajruckman/dbunk-staging/internal/log"
)

func Match(key, qname, ip string) bool {
    match, err := dbunkdb.MatchWhitelist(key, qname, ip, "", "")
    log.Error(err, log.F{
       "key":   key,
       "qname": qname,
       "ip":    ip,
    })
    return match
}
