// Copyright (c) A.J. Ruckman 2019

package blacklist

import (
    "testing"
)

func BenchmarkMatch(b *testing.B) {
    for n := 0; n < b.N; n++ {
        Match("2o7.net", "ziffdavisenterpriseglobal.112.2o7.net")
    }
}
