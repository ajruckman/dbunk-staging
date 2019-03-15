// Copyright (c) A.J. Ruckman 2019

package schema

import (
	"net"
	"time"

	"github.com/lib/pq"

	"github.com/ajruckman/dbunk-staging/pkg/pqext"
)

type Sources struct {
	Source  string `db:"source"`
	Enabled bool   `db:"enabled"`
}

type Blacklist struct {
	Key    string       `db:"key"`
	Rule   pqext.Regexp `db:"rule"`
	Host   *string      `db:"host"`
	Source string       `db:"source"`
}

type Whitelist struct {
	Key       string                   `db:"key"`
	Rule      pqext.Regexp             `db:"rule"`
	Expires   *time.Time               `db:"expires"`
	IPs       *pqext.InetArray         `db:"ips"`
	Subnets   *pqext.CidrArray         `db:"subnets"`
	Hostnames *pq.StringArray          `db:"hostnames"`
	MACs      *pqext.HardwareAddrArray `db:"macs"`
}

type Leases struct {
	Time     time.Time `db:"time"`
	MAC      string    `db:"mac"`
	Lease    string    `db:"lease"`
	Hostname *string   `db:"hostname"`
	Arg      string    `db:"arg"`
}

type Log struct {
	ID             int             `db:"id"`
	Time           time.Time       `db:"time"`
	Client         net.IP          `db:"client"`
	Action         string          `db:"action"`
	Rule           *string         `db:"rule"`
	Source         *string         `db:"source"`
	Question       *string         `db:"question"`
	QuestionType   *string         `db:"question_type"`
	Responses      *pq.StringArray `db:"responses"`
	ResponsesTypes *pq.StringArray `db:"responses_types"`
}
