package main

import (
	"strings"
)

var Caps = map[string]string{}

var capls string

// Can't be in init() because Caps will be registered with init in the future
// and init()s are executed by filename alphabetical order
func setupCapls() {
	capls = ""
	for k, v := range Caps {
		capls += k
		if v != "" {
			capls += "=" + v
		}
		capls += " "
	}
	capls = strings.TrimSuffix(capls, " ")
}
