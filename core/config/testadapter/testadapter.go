/*
Package testadapter is for application configuration during tests.

All configuration is started explicitely with a call to
config.Initialize(testadapter.New()).
There is no init() call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.


BSD License

Copyright (c) 2017–18, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of Norbert Pillmayer nor the names of its contributors
may be used to endorse or promote products derived from this software
without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

*/
package testadapter

import (
	"strconv"
	"strings"

	"github.com/npillmayer/gotype/core/config"
)

type Conf struct {
	conf map[string]string
}

func New() *Conf {
	return &Conf{conf: make(map[string]string)}
}

func (c *Conf) Init() {
	InitDefaults(c.conf)
}

// InitDefaults is usually called by Init().
func InitDefaults(m map[string]string) {
	m["tracing"] = "go"
	m["tracingonline"] = "true"
	m["tracingequations"] = "Error"
	m["tracingsyntax"] = "Error"
	m["tracingcommands"] = "Error"
	m["tracinginterpreter"] = "Error"
	m["tracinggraphics"] = "Error"
	m["tracingscripting"] = "Error"
	m["tracingcore"] = "Error"
	m["tracingengine"] = "Error"

	m["tracingcapsules"] = "Error"
	m["tracingrestores"] = "Error"
	m["tracingchoices"] = "true"
}

func (c *Conf) Set(key string, value string) (oldval string) {
	oldval = c.conf[key]
	c.conf[key] = value
	return
}

// IsSet is a predicate wether a configuration flag is set to true.
func (c *Conf) IsSet(key string) bool {
	_, found := c.conf[key]
	return found
}

func (c *Conf) GetString(key string) string {
	v := c.conf[key]
	return v
}

func (c *Conf) GetInt(key string) int {
	v, found := c.conf[key]
	if !found {
		return 0
	}
	n, _ := strconv.Atoi(v)
	return n
}

func (c *Conf) GetBool(key string) bool {
	v, found := c.conf[key]
	if !found {
		return false
	}
	return strings.EqualFold(v, "true")
}

var _ config.Configuration = &Conf{}
