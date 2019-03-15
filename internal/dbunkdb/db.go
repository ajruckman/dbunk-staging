// Copyright (c) A.J. Ruckman 2019

package dbunkdb

import (
	"github.com/ajruckman/lib/db"
	"github.com/ajruckman/lib/err"

	"github.com/ajruckman/dbunk-staging/internal/config"
)

var (
	DB  libdb.Database
)

func init() {
	DB = libdb.Database{
		User:         config.DBUser,
		Hostname:     "127.0.0.1",
		DatabaseName: config.DBName,
	}
	err := DB.Init()

	liberr.Err(err)
}
