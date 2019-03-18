// Copyright (c) A.J. Ruckman 2019

package blacklist

import (
    "math/rand"
    "testing"
)

func BenchmarkMatch(b *testing.B) {
    s := rand.NewSource(1)
    r := rand.New(s)
    c := ruleCache["2o7.net"]
    l := len(c)

    for n := 0; n < b.N; n++ {
        v := c[r.Intn(l-1)]
        Match("2o7.net", *v.Host)
    }
}
