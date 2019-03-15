// Copyright (c) A.J. Ruckman 2019

package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/ajruckman/dbunk-staging/internal/dbunkdb"
	"github.com/ajruckman/dbunk-staging/internal/log"
)

func main() {
	blacklist, err := dbunkdb.Blacklist()
	log.Err(err)

	spew.Dump(blacklist)

	whitelist, err := dbunkdb.Whitelist()
	log.Err(err)

	fmt.Println()
	fmt.Println()
	fmt.Println()

	spew.Dump(whitelist)

	fmt.Println((*whitelist[0].IPs)[0].String())
	fmt.Println((*whitelist[0].Subnets)[0].String())
}
