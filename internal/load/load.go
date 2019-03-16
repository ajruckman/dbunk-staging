// Copyright (c) A.J. Ruckman 2019

// Package load contains code to read configuration data. 
package load

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/mholt/caddy"

	"github.com/ajruckman/dbunk-staging/internal/config"
	"github.com/ajruckman/dbunk-staging/internal/genrules"
	"github.com/ajruckman/dbunk-staging/internal/log"
	"github.com/ajruckman/dbunk-staging/pkg/stack"
)

func init() {
	flag.BoolVar(&config.TriggerGenDB, "dbunk-gen-db", false, "Run database generator")
	flag.BoolVar(&config.TriggerGenRules, "dbunk-gen-rules", false, "Run rule generator")
	flag.BoolVar(&config.OneShot, "dbunk-oneshot", false, "Exit after running database generator and/or rule generator")
	flag.Parse()
}

func ParseConfig(cont *caddy.Controller) {
	cont.Next()

	go func(){
		fmt.Println(stack.Stack())
	}()

	var last string

	// Iterate over words in the Corefile.
	for cont.Next() {
		var (
			v = strings.TrimPrefix(strings.ToLower(cont.Val()), "!")
			b = !strings.HasPrefix(cont.Val(), "!")
		)

		// Handle string fields. These fields come before string values.
		switch last {
		case "spoofeda":
			config.Conf.SpoofedA = v
			last = ""

		case "spoofedaaaa":
			config.Conf.SpoofedAAAA = v
			last = ""

		case "spoofedcname":
			config.Conf.SpoofedCNAME = v
			last = ""

		case "source":
			config.Conf.Sources = append(config.Conf.Sources, v)
			last = ""

		case "exception":
			config.Conf.Exceptions = append(config.Conf.Exceptions, regexp.MustCompile(v))
			last = ""

		// Handle boolean fields. If a '!' becomes before a field, it is set to
		// false; otherwise, it is true.
		default:
			switch v {
			case "verbose":
				config.Conf.Verbose = b

			case "printstatistics":
				config.Conf.PrintStatistics = b

			case "logunmatchedqueries":
				config.Conf.LogUnmatchedQueries = b

			case "domainneeded":
				config.Conf.DomainNeeded = b

			case "fastreturns":
				config.Conf.FastReturns = b

			case "spoofresponseaddress":
				config.Conf.SpoofResponseAddress = b

			case "interceptforwards":
				config.Conf.InterceptForwards = b

			case "killswitch":
				config.Conf.KillSwitch = b

			case "spoofeda", "spoofedaaaa", "spoofedcname", "source", "exception":
				last = v

			case "{", "}":

			default:
				log.Info("Ignoring unkonwn config parameter: " + v)
			}
		}
	}

	spew.Dump(config.Conf)
}

func HandleFlags() {
	if config.TriggerGenRules {
		genrules.GenRules()
	}

	if config.OneShot {
		log.Info("Oneshot was specified; exiting")
		os.Exit(0)
	}
}
