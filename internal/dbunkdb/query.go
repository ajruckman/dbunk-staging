// Copyright (c) A.J. Ruckman 2019

package dbunkdb

import (
	"github.com/ajruckman/dbunk-staging/pkg/schema"
)

func Blacklist() (res []schema.Blacklist, err error) {
	err = DB.DB.Select(&res, `SELECT * FROM blacklist`)
	return
}

func Whitelist() (res []schema.Whitelist, err error) {
	err = DB.DB.Select(&res, `SELECT * FROM whitelist`)
	return
}
