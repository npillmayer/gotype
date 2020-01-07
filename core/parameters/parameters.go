/*
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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE. */
package parameters

import (
	"golang.org/x/text/unicode/bidi"

	"github.com/npillmayer/gotype/core/dimen"
)

type TypesettingParameter int

//go:generate stringer -type=TypesettingParameter
const (
	none TypesettingParameter = iota
	P_LANGUAGE
	P_SCRIPT
	P_TEXTDIRECTION
	P_BASELINESKIP
	P_LINESKIP
	P_LINESKIPLIMIT
	P_HYPHENCHAR
	P_HYPHENPENALTY
	P_MINHYPHENLENGTH
	P_STOPPER
)

type ParameterGroup struct {
	params map[TypesettingParameter]interface{}
	level  int
	next   *ParameterGroup
}

type TypesettingRegisters struct {
	base       [P_STOPPER]interface{}
	groups     *ParameterGroup
	grouplevel int
}

// ----------------------------------------------------------------------

func NewTypesettingRegisters() *TypesettingRegisters {
	regs := &TypesettingRegisters{}
	initParameters(&regs.base)
	return regs
}

func initParameters(p *[P_STOPPER]interface{}) {
	p[P_LANGUAGE] = "en_EN"               // a string
	p[P_SCRIPT] = "Latin"                 // a string
	p[P_TEXTDIRECTION] = bidi.LeftToRight //
	p[P_BASELINESKIP] = 12 * dimen.PT     // dimension
	p[P_LINESKIP] = 0                     // dimension
	p[P_LINESKIPLIMIT] = 0                // dimension
	p[P_HYPHENCHAR] = int('-')            // a rune
	p[P_HYPHENPENALTY] = 0                // a numeric penalty (int)
	p[P_MINHYPHENLENGTH] = dimen.Infty    // a numeric quantitiv (int) = # of runes
}

func (regs *TypesettingRegisters) Begingroup() {
	regs.grouplevel++
}

func (regs *TypesettingRegisters) Endgroup() {
	if regs.grouplevel > 0 {
		if regs.groups != nil && regs.groups.level == regs.grouplevel {
			regs.groups = regs.groups.next
			regs.grouplevel--
		}
	}
}

func (regs *TypesettingRegisters) Push(key TypesettingParameter, value interface{}) {
	if regs.grouplevel > 0 {
		var g *ParameterGroup
		if regs.groups == nil {
			g = &ParameterGroup{}
			g.params = make(map[TypesettingParameter]interface{})
			g.level = regs.grouplevel
			regs.groups = g
		} else {
			if regs.groups.level < regs.grouplevel {
				g = &ParameterGroup{}
				g.params = make(map[TypesettingParameter]interface{})
				g.level = regs.grouplevel
				g.next = regs.groups
				regs.groups = g
			} else {
				g = regs.groups
			}
		}
		g.params[key] = value
	} else {
		regs.base[key] = value
	}
}

func (regs *TypesettingRegisters) Get(key TypesettingParameter) interface{} {
	if key <= 0 || key == P_STOPPER {
		panic("parameter key outside range of typesetting parameters")
	}
	var value interface{}
	if regs.grouplevel > 0 {
		for g := regs.groups; g != nil; g = g.next {
			value = g.params[key]
			if value != nil {
				break
			}
		}
	}
	if value == nil {
		value = regs.base[key]
	}
	return value
}

func (regs *TypesettingRegisters) S(key TypesettingParameter) string {
	return regs.Get(key).(string)
}

func (regs *TypesettingRegisters) N(key TypesettingParameter) int {
	return regs.Get(key).(int)
}

func (regs *TypesettingRegisters) D(key TypesettingParameter) dimen.Dimen {
	return regs.Get(key).(dimen.Dimen)
}
