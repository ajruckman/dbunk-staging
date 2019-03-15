// Copyright (c) A.J. Ruckman 2019

package log

import (
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var (
	log = clog.NewWithPlugin("dbunk")
)

func Info(v ...interface{}) {
	log.Info(v...)
}

func Err(err error, v ...interface{}) {
	if err == nil {
		return
	}
	log.Error(err, v)
	panic(err)
}
