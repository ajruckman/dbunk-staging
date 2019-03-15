// Copyright (c) A.J. Ruckman 2019

package pqext

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// Based on: github.com/lib/pq/array.go > StringArray

/////

type Regexp struct {
	*regexp.Regexp
}

func (a *Regexp) Scan(src interface{}) (err error) {
	r, err := regexp.Compile(src.(string))
	a.Regexp = r

	return
}

func (a *Regexp) Value() (driver.Value, error) {
	return a, nil
}

/////

type HardwareAddrArray []net.HardwareAddr

func (a *HardwareAddrArray) Scan(src interface{}) (err error) {
	elems, err := parseArrayTyped(src, "HardwareAddrArray")
	if err != nil {
		return
	}

	b := make([]net.HardwareAddr, len(elems))
	for i, v := range elems {
		b[i], err = net.ParseMAC(string(v))
		if err != nil {
			return
		}
	}
	*a = b

	return
}

func (a *HardwareAddrArray) Value() (driver.Value, error) {
	return a, nil
}

/////

type InetArray []net.IP

func (a *InetArray) Scan(src interface{}) (err error) {
	elems, err := parseArrayTyped(src, "InetArray")
	if err != nil {
		return
	}

	b := make([]net.IP, len(elems))
	for i, v := range elems {
		b[i] = net.ParseIP(string(v))
	}
	*a = b

	return
}

func (a *InetArray) Value() (driver.Value, error) {
	return a, nil
}

/////

type CidrArray []*net.IPNet

func (a *CidrArray) Scan(src interface{}) (err error) {
	elems, err := parseArrayTyped(src, "CidrArray")
	if err != nil {
		return
	}

	b := make([]*net.IPNet, len(elems))
	for i, v := range elems {
		_, b[i], err = net.ParseCIDR(string(v))
		if err != nil {
			return
		}
	}
	*a = b

	return nil
}

func (a *CidrArray) Value() (driver.Value, error) {
	return a, nil
}

/////

func parseArrayTyped(src interface{}, typ string) (elems [][]byte, err error) {
	var dims []int
	dims, elems, err = parseArray(src.([]byte), []byte{','})
	if err != nil {
		return
	} else if len(dims) > 1 {
		err = fmt.Errorf("cannot convert ARRAY%s to %s", strings.Replace(fmt.Sprint(dims), " ", "][", -1), typ)
	}

	return
}

// From: github.com/lib/pq/array.go
// Copyright (c) 2011-2013, 'pq' Contributors
// -----
// parseArray extracts the dimensions and elements of an array represented in
// text format. Only representations emitted by the backend are supported.
// Notably, whitespace around brackets and delimiters is significant, and NULL
// is case-sensitive.
//
// See http://www.postgresql.org/docs/current/static/arrays.html#ARRAYS-IO
func parseArray(src, del []byte) (dims []int, elems [][]byte, err error) {
	var depth, i int

	if len(src) < 1 || src[0] != '{' {
		return nil, nil, fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '{', 0)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '{':
			depth++
			i++
		case '}':
			elems = make([][]byte, 0)
			goto Close
		default:
			break Open
		}
	}
	dims = make([]int, i)

Element:
	for i < len(src) {
		switch src[i] {
		case '{':
			if depth == len(dims) {
				break Element
			}
			depth++
			dims[depth-1] = 0
			i++
		case '"':
			var elem []byte
			var escape bool
			for i++; i < len(src); i++ {
				if escape {
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						elems = append(elems, elem)
						i++
						break Element
					}
				}
			}
		default:
			for start := i; i < len(src); i++ {
				if bytes.HasPrefix(src[i:], del) || src[i] == '}' {
					elem := src[start:i]
					if len(elem) == 0 {
						return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
					}
					if bytes.Equal(elem, []byte("NULL")) {
						elem = nil
					}
					elems = append(elems, elem)
					break Element
				}
			}
		}
	}

	for i < len(src) {
		if bytes.HasPrefix(src[i:], del) && depth > 0 {
			dims[depth-1]++
			i += len(del)
			goto Element
		} else if src[i] == '}' && depth > 0 {
			dims[depth-1]++
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	for i < len(src) {
		if src[i] == '}' && depth > 0 {
			depth--
			i++
		} else {
			return nil, nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}
	if depth > 0 {
		err = fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '}', i)
	}
	if err == nil {
		for _, d := range dims {
			if (len(elems) % d) != 0 {
				err = fmt.Errorf("pq: multidimensional arrays must have elements with matching dimensions")
			}
		}
	}
	return
}
