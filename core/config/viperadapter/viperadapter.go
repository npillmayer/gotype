/*
Package viperadapter is for application configuration.

All configuration is started explicitely with a call to
config.Initialize(viperadapter.New()).
There is no init() call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.


BSD License

Copyright (c) 2017–20, Norbert Pillmayer

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

3. Neither the name of this software nor the names of its contributors
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
package viperadapter

import (
	"fmt"

	"github.com/npillmayer/gotype/core/config"
	"github.com/spf13/viper"
)

type VConf struct{}

func New() *VConf {
	return &VConf{}
}

func (c *VConf) Init() {
	c.InitDefaults()
	c.InitConfigPath()
}

// InitDefaults is usually called by Init().
func (c *VConf) InitDefaults() {
	/*
		viper.SetDefault("ContentDir", "content")
		viper.SetDefault("LayoutDir", "layouts")
		viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
		viper.Set("Verbose", true)
	*/

	viper.SetDefault("tracing", "go")
	viper.SetDefault("tracingonline", true)
	viper.SetDefault("tracingequations", "Error")
	viper.SetDefault("tracingsyntax", "Error")
	viper.SetDefault("tracingcommands", "Error")
	viper.SetDefault("tracinginterpreter", "Error")
	viper.SetDefault("tracinggraphics", "Error")

	viper.SetDefault("tracingcapsules", "Error")
	viper.SetDefault("tracingrestores", "Error")
	viper.SetDefault("tracingchoices", true)
}

// InitConfigPath is usually called by Init().
func (c *VConf) InitConfigPath() {
	viper.SetConfigName("gotype")        // name of config file (without extension)
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	viper.AddConfigPath("$GOPATH/etc/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.gotype") // call multiple times to add many search paths
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {                      // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

// Set overrides any configuration values set from the environment.
func (c *VConf) Set(key string, value interface{}) {
	viper.Set(key, value)
}

// IsSet is a predicate wether a configuration flag is set to true.
func (c *VConf) IsSet(key string) bool {
	return viper.IsSet(key)
}

func (c *VConf) GetString(key string) string {
	return viper.GetString(key)
}

func (c *VConf) GetInt(key string) int {
	return viper.GetInt(key)
}

func (c *VConf) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (c *VConf) IsInteractive() bool {
	return viper.GetBool("tracingonline")
}

var _ config.Configuration = &VConf{}
