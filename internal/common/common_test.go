// Copyright (c) A.J. Ruckman 2019

package common

import (
    "testing"

    "github.com/magiconair/properties/assert"
)

type testCase struct {
    In, Out string
}

func TestHostToRawSLD(t *testing.T) {
    var testCases = []testCase{
        {"google", ""},
        {"google.com", "google.com"},
        {"photos.google.com", "google.com"},
        {".google.com", "google.com"},
        {".photos.google.com", "google.com"},
    }

    for _, v := range testCases {
        assert.Equal(t, HostToRawSLD(v.In), v.Out)
    }
}

func TestHostToKey(t *testing.T) {
    var testCases = []testCase{
        {"google", ""},
        {"google.com", "google.com"},
        {"photos.google.com", "google.com"},
        {".google.com", "google.com"},
        {".photos.google.com", "google.com"},
        
        {"google", ""},
        {"google.co.uk", "google.co.uk"},
        {"photos.google.co.uk", "google.co.uk"},
        {".google.co.uk", "google.co.uk"},
        {".photos.google.co.uk", "google.co.uk"},
    }

    for _, v := range testCases {
        assert.Equal(t, HostToKey(v.In), v.Out)
    }
}
