/*
Package config defines types for application configuration.

There is no init() call to set up configuration a priori. The reason
is to avoid coupling to a specific configuration framework, but rather
relay this decision to the client.

Please refer to package config.gconf for concrete usage.

BSD License

Copyright (c) 2017–21, Norbert Pillmayer

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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. */
package config

import (
	"sync"

	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gologadapter"
)

// Configuration is an interface to be implemented by every configuration
// adapter.
type Configuration interface {
	InitDefaults()               // initialize key/value pairs
	IsSet(key string) bool       // is a config key set?
	GetString(key string) string // get config value as string
	GetInt(key string) int       // get config value as integer
	GetBool(key string) bool     // get config value as boolean
	IsInteractive() bool         // Are we running in interactive mode?
}

var knownTraceAdapters = map[string]tracing.Adapter{
	"go": gologadapter.GetAdapter(),
	//"logrus": logrusadapter.GetAdapter(), // now to be set by AddTraceAdapter()
}
var adapterMutex = &sync.RWMutex{} // guard knownTraceAdapters[]

// AddTraceAdapter is an extension point for clients who want to use
// their own tracing adapter implementation.
// key will be used at configuration initialization time to identify
// this adapter, e.g. in configuration files.
//
// Clients will have to call this before any call to tracing-initialization,
// otherwise the adapter cannot be found.
func AddTraceAdapter(key string, adapter tracing.Adapter) {
	adapterMutex.Lock()
	defer adapterMutex.Unlock()
	knownTraceAdapters[key] = adapter
}

// GetAdapterFromConfiguration gets the concrete tracing implementation adapter
// from the appcation configuration. The configuration key name is "tracing".
//
// The value must be one of the known tracing adapter keys.
// Default is an adapter for the Go standard log package.
func GetAdapterFromConfiguration(conf Configuration) tracing.Adapter {
	adapterPackage := conf.GetString("tracing")
	adapterMutex.RLock()
	defer adapterMutex.RUnlock()
	adapter := knownTraceAdapters[adapterPackage]
	if adapter == nil {
		adapter = gologadapter.GetAdapter()
	}
	return adapter
}
