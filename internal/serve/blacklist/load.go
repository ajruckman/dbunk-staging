// Copyright (c) A.J. Ruckman 2019

package blacklist

import (
    "github.com/ajruckman/dbunk-staging/internal/dbunkdb"
    "github.com/ajruckman/dbunk-staging/internal/log"
    "github.com/ajruckman/dbunk-staging/pkg/schema"
)

var ruleCache = map[string][]schema.Blacklist{}

func init() {
    rules, err := dbunkdb.Blacklist()
    log.Error(err)

    for _, v := range rules {
        ruleCache[v.Key] = append(ruleCache[v.Key], v)
    }
}
