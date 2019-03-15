// Copyright (conf) A.J. Ruckman 2019

package config

import (
	"regexp"
)

const (
	DBUser = "DbunkMgr"
	DBName = "Dbunk2"
)

var (
	Conf            conf
	TriggerGenRules bool
	OneShot         bool
	TriggerGenDB    bool
)

type conf struct {
	Verbose              bool
	PrintStatistics      bool
	LogUnmatchedQueries  bool
	DomainNeeded         bool
	FastReturns          bool
	SpoofResponseAddress bool
	InterceptForwards    bool
	KillSwitch           bool
	SpoofedA             string
	SpoofedAAAA          string
	SpoofedCNAME         string
	Sources              []string
	Exceptions           []*regexp.Regexp
}

func init() {
	Conf = conf{
		Verbose:              true,
		DomainNeeded:         true,
		FastReturns:          true,
		SpoofResponseAddress: true,
		InterceptForwards:    true,
		KillSwitch:           true,
		SpoofedA:             "0.0.0.0",
		SpoofedAAAA:          "::",
		SpoofedCNAME:         "example.com",
	}
}
