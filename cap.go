package main

import (
	"strings"
)

var Caps = map[string]string{}

var (
	capls    string
	capls302 string
)

// Can't be in init() because Caps will be registered with init in the future
// and init()s are executed by filename alphabetical order
func setupCapls() {
	capls = ""
	capls302 = ""
	for k, v := range Caps {
		capls += k
		capls302 += k
		if v != "" {
			capls302 += "=" + v
		}
		capls += " "
		capls302 += " "
	}
	capls = strings.TrimSuffix(capls, " ")
	capls302 = strings.TrimSuffix(capls302, " ")
}
