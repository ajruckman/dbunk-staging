// Copyright (c) A.J. Ruckman 2019

package log

import (
    "fmt"
    "strings"

    clog "github.com/coredns/coredns/plugin/pkg/log"
)

var (
    log = clog.NewWithPlugin("dbunk")
)

type F map[string]interface{}

func logFields(t, msg string, fields ...F) {
    callLogger(t, msg)
    if len(fields) > 0 {
        for k, v := range fields[0] {
            if strings.Contains(fmt.Sprintf("%v", v), "\n") {
                callLogger(t, "─> ", k, ":")
                for _, line := range strings.Split(fmt.Sprintf("%v", v), "\n") {
                    callLogger(t, "───> ", line)
                }
            } else {
                callLogger(t, "─> ", k, ": ", fmt.Sprintf("%v", v))
            }
        }
    }
}

func callLogger(t string, v ...interface{}) {
    switch t {
    case "debug":
        log.Debug(v...)
    case "info":
        log.Info(v...)
    case "warning":
        log.Warning(v...)
    case "error":
        log.Error(v...)
    case "fatal":
        log.Fatal(v...)
    }
}

func Debug(msg string, fields ...F) {
    logFields("debug", msg, fields...)
}

func Info(msg string, fields ...F) {
    logFields("info", msg, fields...)
}

func Warning(msg string, fields ...F) {
    logFields("warning", msg, fields...)
}

func Error(err error, fields ...F) {
    if err == nil {
        return
    }
    logFields("error", err.Error(), fields...)
}

func Fatal(err error, fields ...F) {
    if err == nil {
        return
    }
    logFields("fatal", err.Error(), fields...)
}
