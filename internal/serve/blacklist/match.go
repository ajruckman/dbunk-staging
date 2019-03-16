// Copyright (c) A.J. Ruckman 2019

package blacklist

func Match(key, qname string) bool {
    for _, v := range ruleCache[key] {
       if v.Rule.MatchString(qname) {
           return true
       }
    }
    return false
}
