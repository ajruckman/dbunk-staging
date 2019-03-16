// Copyright (c) A.J. Ruckman 2019

package stack

import (
    "regexp"
    "runtime/debug"
    "strconv"
    "strings"
)

var (
    matchGoroutineNum         = regexp.MustCompile(`^goroutine (\d+)`)
    matchGoroutineCreatorFunc = regexp.MustCompile(`^created by (\S.*?)(?:\((?:0x\w+(?:, \.\.\.\)|, |\))|\)|\.\.\.\))+|$)`)
    matchFunc                 = regexp.MustCompile(`^(\S.*?)(?:\((?:0x\w+(?:, \.\.\.\)|, |\))|\)|\.\.\.\))+|$)`)
    matchFuncPath             = regexp.MustCompile(`^\s+(.*)\:(\d+)(?: \+0x\w+|$)`)
)

type tracedFunc struct {
    Func string
    Path string
    Line int
}

type trace struct {
    GoroutineID      int
    GoroutineCreator tracedFunc
    Functions        []tracedFunc
}

func Stack() []trace {
    var (
        in        bool
        inCreator bool
        cur       trace
        curFunc   tracedFunc
        traces    []trace
    )

    s := string(debug.Stack())

    for _, line := range strings.Split(s, "\n") {
        if matchGoroutineNum.MatchString(line) {
            if ! in {
                gid_s := matchGoroutineNum.FindStringSubmatch(line)[1]
                gid, _ := strconv.Atoi(gid_s)
                cur.GoroutineID = gid

                in = true

            } else {
                traces = append(traces, cur)
                in = false
            }

        } else if matchGoroutineCreatorFunc.MatchString(line) {
            cur.GoroutineCreator.Func = matchGoroutineCreatorFunc.FindStringSubmatch(line)[1]

            inCreator = true

        } else if inCreator {
            cur.GoroutineCreator.Path = matchFuncPath.FindStringSubmatch(line)[1]
            cur.GoroutineCreator.Line, _ = strconv.Atoi(matchFuncPath.FindStringSubmatch(line)[2])
            traces = append(traces, cur)

            inCreator = false
            in = false

        } else if matchFunc.MatchString(line) {
            curFunc.Func = matchFunc.FindStringSubmatch(line)[1]

        } else if matchFuncPath.MatchString(line) {
            curFunc.Path = matchFuncPath.FindStringSubmatch(line)[1]
            curFunc.Line, _ = strconv.Atoi(matchFuncPath.FindStringSubmatch(line)[2])
            cur.Functions = append(cur.Functions, curFunc)
        }
    }

    if in {
        traces = append(traces, cur)
    }

    return traces
}
