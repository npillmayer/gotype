package style

import (
	"errors"
	"fmt"
	"strings"
)

/* -----------------------------------------------------------------
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
----------------------------------------------------------------- */

// --- CSS Properties ---------------------------------------------------
//
// https://www.tutorialrepublic.com/css-reference/css3-properties.php
// https://www.mediaevent.de/xhtml/kernattribute.html

// Property is a raw value for a CSS property. For example, with
//     color: black
// a property value of "black" is set.
type Property string

func (p Property) String() string {
	return string(p)
}

// property = "initial"
func (p Property) IsInitial() bool {
	return p == "initial"
}

// property = "inherited"
func (p Property) IsInherited() bool {
	return p == "inherited"
}

func (p Property) IsEmpty() bool {
	return p == ""
}

// KeyValue is a container for a style property.
type KeyValue struct {
	Key   string
	Value Property
}

// --- CSS Property Groups ----------------------------------------------
//
// Caching is currently not implemented.

// propertyGroup is a collected of propertes sharing a common topic.
// CSS knows a whole lot of properties. We split them up into organisatorial
// groups.
//
// The mapping of property into groups is documented with
// GroupNameFromPropertyKey[...].
type propertyGroup struct {
	name          string
	Parent        *propertyGroup
	propertiesMap map[string]Property
	//signature       uint32 // signature of IDs and classes, used for caching
}

func newPropertyGroup(ofType string) *propertyGroup {
	pg := &propertyGroup{}
	pg.name = ofType
	return pg
}

// Stringer for property groups; used for debugging.
func (pg *propertyGroup) String() string {
	s := "[" + pg.name + "] =\n"
	for k, v := range pg.propertiesMap {
		s += fmt.Sprintf("  %s = %s\n", k, v)
	}
	return s
}

// IsSet is a predicated wether a property is set within this group.
func (pg *propertyGroup) IsSet(key string) bool {
	if pg.propertiesMap == nil {
		return false
	}
	v, ok := pg.propertiesMap[key]
	return ok && !v.IsEmpty()
}

// Get a property's value.
func (pg *propertyGroup) Get(key string) Property {
	if pg.propertiesMap == nil {
		return ""
	}
	return pg.propertiesMap[key]
}

// Set a property's value.
func (pg *propertyGroup) Set(key string, p Property) {
	if pg.propertiesMap == nil {
		pg.propertiesMap = make(map[string]Property)
	}
	pg.propertiesMap[key] = p
}

func (pg *propertyGroup) SpawnOn(key string, p Property, cascade bool) (*propertyGroup, bool) {
	if cascade {
		found := pg.Cascade(key)
		if found.Get(key) == p {
			return pg, false
		}
	}
	npg := newPropertyGroup(pg.name)
	//npg.signature = pg.signature
	npg.Set(key, p)
	return npg, true
}

func (pg *propertyGroup) Cascade(key string) *propertyGroup {
	it := pg
	for !it.IsSet(key) { // stopper is default partial
		it = it.Parent
	}
	return it
}

// Signature for being able to cache a fragment.
// Caching is currently not implemented.
//
// Class values have to be sorted.
//
// Returns hash and number of ID+class attributes.
//
// TODO
/*
func HashSignatureAttributes(htmlNode *html.Node) (uint32, uint8) {
	var hash uint32 = 0
	var count uint8 = 0
	signature := ""
	for _, a := range htmlNode.Attr {
		if a.Key == "id" {
			signature += a.Key
			count += 1
		}
	}
	for _, a := range htmlNode.Attr {
		if a.Key == "class" {
			signature += a.Key
			count += 1
		}
	}
	hash = xxhash.Checksum32([]byte(signature))
	return hash, count
}
*/

var groupNameFromPropertyKey = map[string]string{
	"margin-top":                 "Margins", // Margins
	"margin-left":                "Margins",
	"margin-right":               "Margins",
	"margin-bottom":              "Margins",
	"padding-top":                "Padding", // Padding
	"padding-left":               "Padding",
	"padding-right":              "Padding",
	"padding-bottom":             "Padding",
	"border-top-color":           "Border", // Border
	"border-left-color":          "Border",
	"border-right-color":         "Border",
	"border-bottom-color":        "Border",
	"border-top-width":           "Border",
	"border-left-width":          "Border",
	"border-right-width":         "Border",
	"border-bottom-width":        "Border",
	"border-top-style":           "Border",
	"border-left-style":          "Border",
	"border-right-style":         "Border",
	"border-bottom-style":        "Border",
	"border-top-left-radius":     "Border",
	"border-top-right-radius":    "Border",
	"border-bottom-left-radius":  "Border",
	"border-bottom-right-radius": "Border",
	"width":                      "Dimension", // Dimension
	"height":                     "Dimension",
	"min-width":                  "Dimension",
	"min-height":                 "Dimension",
	"max-width":                  "Dimension",
	"max-height":                 "Dimension",
}

// SplitCompountProperty splits up a shortcut property into its individual
// components.
func splitCompoundProperty(key string, value Property) ([]KeyValue, error) {
	fields := strings.Fields(value.String())
	switch key {
	case "margins":
		return feazeCompound4("margin", "", fourDirs, fields)
	case "padding":
		return feazeCompound4("padding", "", fourDirs, fields)
	case "border-color":
		return feazeCompound4("border", "color", fourDirs, fields)
	case "border-width":
		return feazeCompound4("border", "width", fourDirs, fields)
	case "border-style":
		return feazeCompound4("border", "style", fourDirs, fields)
	case "border-radius":
		return feazeCompound4("border", "style", fourCorners, fields)
	}
	return nil, errors.New(fmt.Sprintf("Not recognized as compound property: %s", key))
}

// CSS logic to distribute individual values from compound shortcuts is as
// follows: https://www.w3schools.com/css/css_border.asp
func feazeCompound4(pre string, suf string, dirs [4]string, fields []string) ([]KeyValue, error) {
	l := len(fields)
	if l == 0 || l > 4 {
		return nil, errors.New(fmt.Sprintf("Expecting 1-3 values for %s-%s", pre, suf))
	}
	r := make([]KeyValue, 4, 4)
	r[0] = KeyValue{p(pre, suf, dirs[0]), Property(fields[0])}
	if l >= 2 {
		r[1] = KeyValue{p(pre, suf, dirs[1]), Property(fields[1])}
		if l >= 3 {
			r[2] = KeyValue{p(pre, suf, dirs[2]), Property(fields[2])}
			if l == 4 {
				r[3] = KeyValue{p(pre, suf, dirs[3]), Property(fields[3])}
			} else {
				r[3] = KeyValue{p(pre, suf, dirs[3]), Property(fields[1])}
			}
		} else {
			r[2] = KeyValue{p(pre, suf, dirs[2]), Property(fields[0])}
			r[3] = KeyValue{p(pre, suf, dirs[3]), Property(fields[1])}
		}
	} else {
		r[1] = KeyValue{p(pre, suf, dirs[1]), Property(fields[0])}
		r[2] = KeyValue{p(pre, suf, dirs[2]), Property(fields[0])}
		r[3] = KeyValue{p(pre, suf, dirs[3]), Property(fields[0])}
	}
	return r, nil
}

var fourDirs = [4]string{"top", "right", "bottom", "left"}
var fourCorners = [4]string{"top-right", "bottom-right", "bottom-left", "top-left"}

func p(prefix string, suffix string, tag string) string {
	if suffix == "" {
		return prefix + "-" + tag
	}
	if prefix == "" {
		return tag + "-" + suffix
	}
	return prefix + "-" + tag + "-" + suffix
}

// PropertyMap holds CSS properties. As CSS defines a whole lot of properties,
// we segment them into logical group.
type PropertyMap struct {
	m map[string]*propertyGroup // into struct to make it opaque for clients
}

// Add adds a property to this property map, e.g.,
//    pm.Add("funny-margin", "big")
func (pm *PropertyMap) Add(key string, value string) {
	if pm != nil {
		group, found := pm.m["X"]
		if !found {
			group = newPropertyGroup("X")
			pm.m["X"] = group
		}
		group.Set(key, Property(value))
	}
}

// InitializeDefaultPropertyValues creates an internal data structure to
// hold all the default values for CSS properties.
// In real-world browsers these are the user-agent CSS values.
func initializeDefaultPropertyValues(additionalProps []KeyValue) *PropertyMap {
	m := make(map[string]*propertyGroup, 15)
	root := newPropertyGroup("Root")

	x := newPropertyGroup("X") // special group for extension properties
	for _, kv := range additionalProps {
		x.Set(kv.Key, kv.Value)
	}
	m["X"] = x

	margins := newPropertyGroup("Margins")
	margins.Set("margin-top", "0")
	margins.Set("margin-left", "0")
	margins.Set("margin-right", "0")
	margins.Set("margin-bottom", "0")
	margins.Parent = root
	m["Margins"] = margins

	padding := newPropertyGroup("Padding")
	padding.Set("padding-top", "0")
	padding.Set("padding-left", "0")
	padding.Set("padding-right", "0")
	padding.Set("padding-bottom", "0")
	padding.Parent = root
	m["Padding"] = padding

	border := newPropertyGroup("Border")
	border.Set("border-top-color", "black")
	border.Set("border-left-color", "black")
	border.Set("border-right-color", "black")
	border.Set("border-bottom-color", "black")
	border.Set("border-top-width", "medium")
	border.Set("border-left-width", "medium")
	border.Set("border-right-width", "medium")
	border.Set("border-bottom-width", "medium")
	border.Set("border-top-style", "solid")
	border.Set("border-left-style", "solid")
	border.Set("border-right-style", "solid")
	border.Set("border-bottom-style", "solid")
	border.Set("border-top-left-radius", "0")
	border.Set("border-top-right-radius", "0")
	border.Set("border-bottom-left-radius", "0")
	border.Set("border-bottom-right-radius", "0")
	border.Parent = root
	m["Border"] = border

	dimension := newPropertyGroup("Dimension")
	dimension.Set("width", "10%")
	dimension.Set("width", "100pt")
	dimension.Set("min-width", "0")
	dimension.Set("min-height", "0")
	dimension.Set("max-width", "10000pt")
	dimension.Set("max-height", "10000pt")
	dimension.Parent = root
	m["Dimension"] = dimension

	/*
	   type ColorModel string

	   type Color struct {
	   	Color   color.Color
	   	Model   ColorModel
	   	Opacity uint8
	   }

	   type DisplayStyle struct {
	   	Display    uint8 // https://www.tutorialrepublic.com/css-reference/css-display-property.php
	   	Position   uint8
	   	Top        dimen.Dimen
	   	Left       dimen.Dimen
	   	Right      dimen.Dimen
	   	Bottom     dimen.Dimen
	   	Float      uint8
	   	ZIndex     int
	   	Overflow   uint8
	   	OverflowX  uint8
	   	OverflowY  uint8
	   	Clip       string // geometric shape
	   	Visibility bool
	   }

	   type Background struct {
	   	Color color.Color
	   	//Position TODO
	   	Image  image.Image
	   	Origin dimen.Point
	   	Size   dimen.Point
	   	Clip   uint8
	   }

	   type Font struct {
	   	Family     string
	   	Style      string
	   	Variant    uint16
	   	Stretch    uint8
	   	Size       dimen.Dimen
	   	SizeAdjust dimen.Dimen
	   }

	   type TextProperties struct {
	   	Direction          uint8
	   	WordSpacing        uint8
	   	LetterSpacing      uint8
	   	VerticalAlignment  uint8
	   	TextAlignment      uint8 // + TextJustify
	   	TextAlignLast      uint8
	   	TextIndentation    dimen.Dimen // first line
	   	TabSize            dimen.Dimen
	   	LineHeight         uint8
	   	TextDecoration     uint8
	   	TextTransformation uint8
	   	WordWrap           uint8
	   	WordBreak          uint8
	   	Whitespace         uint8
	   	TextOverflow       uint8
	   }


	   type GeneratedContent struct {
	   	Content          string
	   	Quotes           string
	   	CounterReset     uint8
	   	CounterIncrement uint8
	   }

	   type Print struct {
	   	PageBreakAfter  uint8
	   	PageBreakBefore uint8
	   	PageBreakInside uint8
	   }

	   type Outline struct {
	   	Color  color.Color
	   	Offset dimen.Dimen
	   	Style  uint8
	   	Width  dimen.Dimen
	   }

	   //list-style-type:
	   //	disc | circle | square | decimal | decimal-leading-zero | lower-roman |
	   //  upper-roman | lower-greek | lower-latin | upper-latin | armenian |
	   //  georgian | lower-alpha | upper-alpha | none | initial | inherit
	   type List struct {
	   	StyleImage    image.Image
	   	StylePosition uint8 // inside, outside
	   	StyleType     uint8
	   }

	*/

	return &PropertyMap{m}
}
