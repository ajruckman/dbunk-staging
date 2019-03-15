// Copyright (c) A.J. Ruckman 2019

package common

import "regexp"

// HostToKey returns a key unique to a hostname and all of its subdomains. This
// key will always be the same for domains with the same second- or third-level
// domain.
// This function is one of the most important components of dbunk. It allows us
// to quickly narrow down hundreds of thousands of possible rules for a request
// to just a few.
func HostToKey(host string) string {
	sldMatch := MatchSLD.FindStringSubmatch(host)
	if len(sldMatch) == 0 {
		return ""
	}
	sld := sldMatch[0]

	if tldsToParticularize[sld] {
		tldMatch := MatchTLD.FindStringSubmatch(host)
		// Fall back to second-level domain if third-level domain is blank.
		if len(tldMatch) == 0 {
			return sld
		} else {
			return tldMatch[0]
		}
	} else {
		return sld
	}
}

// HostToRawSLD always returns the second-level domain. This function is used
// when importing rules from host lists and when de-duping rules to ensure that,
// if a general rule for a SLD exists (e.g. '.*\.2o7\.net'), we do not import or
// keep subdomains of that domain (e.g. '10.122.207.net).
func HostToRawSLD(host string) string {
	sldMatch := MatchSLD.FindStringSubmatch(host)
	if len(sldMatch) == 0 {
		return ""
	}
	sld := sldMatch[0]
	return sld
}

var (
	// Match second-level domain
	MatchSLD = regexp.MustCompile(`[^\s.]+\.[^.]+$`)

	// Match third-level do main
	MatchTLD = regexp.MustCompile(`[^\s.]+\.[^\s.]+\.[^.]+$`)
)

// Return third-level domains instead of second-level domains for hostnames with
// these suffixes. This is so that hostnames with two-part TLDs (like .co.uk)
// are filtered effectively with appropriate keys. Otherwise, all hostnames with
// two-part TLDs would have the same keys, potentially resulting in thousands of
// rules with the same key, like all domains with a .co.uk TLD. I have found
// that these keys show up the most in imported rules; however, it should be
// safe to add any two-part TLD here.
var tldsToParticularize = map[string]bool{
	"com.br":            true,
	"000webhostapp.com": true,
	"co.uk":             true,
	"com.au":            true,
	"neliver.com":       true,
}
