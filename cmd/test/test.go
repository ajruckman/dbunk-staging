// Copyright (c) A.J. Ruckman 2019

package main

import (
	"fmt"
	"runtime/debug"
)

func main() {

	fmt.Println(string(debug.Stack()))

	//blacklist, err := dbunkdb.Blacklist()
	//log.Err(err)
	//
	//spew.Dump(blacklist)
	//
	//whitelist, err := dbunkdb.Whitelist()
	//log.Err(err)
	//
	//fmt.Println()
	//fmt.Println()
	//fmt.Println()
	//
	//spew.Dump(whitelist)
	//
	//fmt.Println((*whitelist[0].IPs)[0].String())
	//fmt.Println((*whitelist[0].Subnets)[0].String())
}
