// Copyright (c) A.J. Ruckman 2019

package dbunkdb

import (
    "github.com/ajruckman/dbunk-staging/pkg/schema"
)

func Blacklist() (res []schema.Blacklist, err error) {
    err = DB.DB.Select(&res, `
      SELECT b.key, b.rule, b.host, b.source
      FROM blacklist b
             INNER JOIN sources s ON b.source = s.source
      WHERE s.enabled = TRUE
    `)
    return
}

func Whitelist() (res []schema.Whitelist, err error) {
    err = DB.DB.Select(&res, `SELECT * FROM whitelist`)
    return
}

func MatchBlacklist(key, qname string) (matched bool, err error) {
    err = DB.DB.Get(&matched, `      
      SELECT exists(
          SELECT b.key
          FROM blacklist b
                 INNER JOIN sources s ON b.source = s.source
          WHERE s.enabled = TRUE
            AND b.key = $1
            AND $2 ~ RULE
        );
    `, key, qname)

    return
}

func MatchWhitelist(key, qname, ip, hostname, mac string) (matched bool, err error) {
    err = DB.DB.Get(&matched, `      
      SELECT EXISTS(
          SELECT key
          FROM whitelist
          WHERE $1 = key
            AND $2 ~ rule
            AND now() < expires
            AND (
              $3 = ANY (ips)
              OR CASE
                   WHEN text($3) != ''
                     THEN $3 << ANY (subnets)
                   ELSE
                     FALSE
                END
              OR $4 = ANY (hostnames)
              OR CASE
                   WHEN text($5) != ''
                     THEN $5::MACADDR = ANY (macs)
                   ELSE
                     FALSE
                END
            )
        );
    `, key, qname, ip, hostname, mac)

    return
}
