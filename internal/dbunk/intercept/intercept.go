// Copyright (c) A.J. Ruckman 2019

package intercept

import "github.com/miekg/dns"

// ResponseInterceptor wraps a dns.ResponseWriter to intercept messages written
// to it by the next plugins in the plugin chain. This allows us to log
// responses from the 'forward' plugin.
// Based on: https://github.com/coredns/example/blob/master/example.go
type ResponseInterceptor struct {
    dns.ResponseWriter
}

// WriteMsg calls a ResponseInterceptor's underlying ResponseWriter's WriteMsg
// method.
func (r *ResponseInterceptor) WriteMsg(res *dns.Msg) error {
    return r.ResponseWriter.WriteMsg(res)
}

// NewResponseInterceptor returns a new ResponseInterceptor.
func NewResponseInterceptor(w dns.ResponseWriter) *ResponseInterceptor {
    return &ResponseInterceptor{
        ResponseWriter: w,
    }
}
