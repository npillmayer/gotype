package style

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
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"errors"
	"fmt"
	"strings"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/engine/tree"
)

// T returns a global tracer. We trace to the EngineTracer
func T() tracing.Trace {
	return gtrace.EngineTracer
}

// Property is a raw value for a CSS property. For example, with
//
//     color: black
//
// a property value of "black" is set. The main purpose of wrapping
// the raw string value into type Property is to provide a set of
// convenient type conversion functions and other helpers.
type Property string

// NullStyle is an empty property value.
const NullStyle Property = ""

func (p Property) String() string {
	return string(p)
}

// IsInitial denotes if a property is of inheritence-type "initial"
func (p Property) IsInitial() bool {
	return p == "initial"
}

// IsInherited denotes if a property is of inheritence-type "inheritet"
func (p Property) IsInherited() bool {
	return p == "inherited"
}

// IsEmpty checks wether a property is empty, i.e. the null-string.
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

// PropertyGroup is a collection of propertes sharing a common topic.
// CSS knows a whole lot of properties. We split them up into organisatorial
// groups.
//
// The mapping of property into groups is documented with
// GroupNameFromPropertyKey[...].
type PropertyGroup struct {
	name      string
	Parent    *PropertyGroup
	propsDict map[string]Property
}

// NewPropertyGroup creates a new empty property group, given its name.
func NewPropertyGroup(groupname string) *PropertyGroup {
	pg := &PropertyGroup{}
	pg.name = groupname
	return pg
}

// Name returns the name of the property group. Once named (during
// construction, property groups may not be renamed.
func (pg *PropertyGroup) Name() string {
	return pg.name
}

// Stringer for property groups; used for debugging.
func (pg *PropertyGroup) String() string {
	s := "[" + pg.name + "] =\n"
	for k, v := range pg.propsDict {
		s += fmt.Sprintf("  %s = %s\n", k, v)
	}
	return s
}

// Properties returns all properties of a group.
func (pg *PropertyGroup) Properties() []KeyValue {
	i := 0
	r := make([]KeyValue, len(pg.propsDict))
	for k, v := range pg.propsDict {
		r[i] = KeyValue{k, v}
		i++
	}
	return r
}

// IsSet is a predicated wether a property is set within this group.
func (pg *PropertyGroup) IsSet(key string) bool {
	if pg.propsDict == nil {
		return false
	}
	v, ok := pg.propsDict[key]
	return ok && !v.IsEmpty()
}

// Get a property's value.
//
// Style property values are always converted to lower case.
func (pg *PropertyGroup) Get(key string) (Property, bool) {
	if pg.propsDict == nil {
		return NullStyle, false
	}
	p, ok := pg.propsDict[key]
	return p, ok
}

// Set a property's value. Overwrites an existing value, if present.
//
// Style property values are always converted to lower case.
func (pg *PropertyGroup) Set(key string, p Property) {
	p = Property(strings.ToLower(string(p)))
	if pg.propsDict == nil {
		pg.propsDict = make(map[string]Property)
	}
	pg.propsDict[key] = p
}

// Add a property's value. Does not overwrite an existing value, i.e., does nothing
// if a value is already set.
func (pg *PropertyGroup) Add(key string, p Property) {
	if pg.propsDict == nil {
		pg.propsDict = make(map[string]Property)
	}
	_, exists := pg.propsDict[key]
	if !exists {
		pg.propsDict[key] = p
	}
}

// ForkOnProperty creates a new PropertyGroup, pre-filled with a given property.
// If 'cascade' is true, the new PropertyGroup will be
// linking to the ancesting PropertyGroup containing this property.
func (pg *PropertyGroup) ForkOnProperty(key string, p Property, cascade bool) (*PropertyGroup, bool) {
	var ancestor *PropertyGroup
	if cascade {
		ancestor = pg.Cascade(key)
		if ancestor != nil {
			p2, _ := ancestor.Get(key)
			if p2 == p {
				return pg, false
			}
		}
	}
	npg := NewPropertyGroup(pg.name)
	npg.Parent = ancestor
	//npg.signature = pg.signature
	npg.Set(key, p)
	return npg, true
}

// Cascade finds the ancesiting PropertyGroup containing the given property-key.
func (pg *PropertyGroup) Cascade(key string) *PropertyGroup {
	it := pg
	for it != nil && !it.IsSet(key) { // stopper is default partial
		it = it.Parent
	}
	return it
}

// GroupNameFromPropertyKey returns the style property group name for a
// style property.
// Example:
//    GroupNameFromPropertyKey("margin-top") => "Margins"
//
// Unknown style property keys will return a group name of "X".
func GroupNameFromPropertyKey(key string) string {
	groupname, found := groupNameFromPropertyKey[key]
	if !found {
		groupname = "X"
	}
	return groupname
}

// Symbolic names for string literals, denoting PropertyGroups.
const (
	PGMargins   = "Margins"
	PGPadding   = "Padding"
	PGBorder    = "Border"
	PGDimension = "Dimension"
	PGDisplay   = "Display"
	PGRegion    = "Region"
	PGX         = "X"
)

var groupNameFromPropertyKey = map[string]string{
	"margin-top":                 PGMargins, // Margins
	"margin-left":                PGMargins,
	"margin-right":               PGMargins,
	"margin-bottom":              PGMargins,
	"padding-top":                PGPadding, // Padding
	"padding-left":               PGPadding,
	"padding-right":              PGPadding,
	"padding-bottom":             PGPadding,
	"border-top-color":           PGBorder, // Border
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
	"width":      "Dimension", // Dimension
	"height":     "Dimension",
	"min-width":  "Dimension",
	"min-height": "Dimension",
	"max-width":  "Dimension",
	"max-height": "Dimension",
	"display":    PGDisplay, // Display
	"float":      PGDisplay,
	"visibility": PGDisplay,
	"position":   PGDisplay,
	"flow-into":  PGRegion,
	"flow-from":  PGRegion,
}

// isCascading returns wether the standard behaviour for a propery is to be
// inherited or not, i.e., a call to retrieve its value will cascade.
func isCascading(key string) bool {
	if strings.HasPrefix(key, "list-style") {
		return true
	}
	switch key {
	case "color", "cursor", "direction", "position", "flow-into", "flow-from":
		return true
	case "letter-spacing", "line-height", "quotes", "visibility", "white-space":
		return true
	case "word-spacing", "word-break", "word-wrap":
		return true
	}
	return false
}

// SplitCompoundProperty splits up a shortcut property into its individual
// components. Returns a slice of key-value pairs representing the
// individual (fine grained) style properties.
// Example:
//    SplitCompountProperty("padding", "3px")
// will return
//    "padding-top"    => "3px"
//    "padding-right"  => "3px"
//    "padding-bottom" => "3px"
//    "padding-left  " => "3px"
// For the logic behind this, refer to e.g.
// https://www.w3schools.com/css/css_padding.asp .
func SplitCompoundProperty(key string, value Property) ([]KeyValue, error) {
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
	return nil, fmt.Errorf("Not recognized as compound property: %s", key)
}

// CSS logic to distribute individual values from compound shortcuts is as
// follows: https://www.w3schools.com/css/css_border.asp
func feazeCompound4(pre string, suf string, dirs [4]string, fields []string) ([]KeyValue, error) {
	l := len(fields)
	if l == 0 || l > 4 {
		return nil, fmt.Errorf("Expecting 1-3 values for %s-%s", pre, suf)
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

// --- Property Map -----------------------------------------------------

// PropertyMap holds CSS properties. nil is a legal (empty) property map.
// A property map is the entity styling a DOM node: a DOM node links to a property map,
// which contains zero or more property groups. Property maps may share property groups.
type PropertyMap struct {
	// As CSS defines a whole lot of properties, we segment them into logical groups.
	m map[string]*PropertyGroup // into struct to make it opaque for clients
}

// NewPropertyMap returns a new empty property map.
func NewPropertyMap() *PropertyMap {
	return &PropertyMap{}
}

func (pmap *PropertyMap) String() string {
	s := "Property Map = {\n"
	for _, v := range pmap.m {
		s += v.String()
	}
	s += "}"
	return s
}

// Size returns the number of property groups.
func (pmap *PropertyMap) Size() int {
	if pmap == nil {
		return 0
	}
	return len(pmap.m)
}

// Group returns the property group for a group name or nil.
func (pmap *PropertyMap) Group(groupname string) *PropertyGroup {
	if pmap == nil {
		return nil
	}
	group, _ := pmap.m[groupname]
	return group
}

// Property returns a style property value, together with an indicator
// wether it has been found in the properties map.
// No cascading is performed
func (pmap *PropertyMap) Property(key string) (Property, bool) {
	groupname := GroupNameFromPropertyKey(key)
	group := pmap.Group(groupname)
	if group == nil {
		return NullStyle, false
	}
	return group.Get(key)
}

// GetPropertyValue returns the property value for a given key.
// If the property is inherited, it may cascade.
func (pmap *PropertyMap) GetPropertyValue(key string, node *tree.Node, styler Interf) Property {
	p, ok := pmap.Property(key)
	if ok {
		if p != "inherit" {
			return p
		}
	}
	// not found in local dicts => cascade, if allowed
	if p == "inherit" || isCascading(key) {
		groupname := GroupNameFromPropertyKey(key)
		var group *PropertyGroup
		//func (pg *PropertyGroup) Cascade(key string) *PropertyGroup {
		for node != nil && group == nil {
			styler := styler(node)
			group = styler.ComputedStyles().Group(groupname)
			node = node.Parent()
		}
		if group == nil {
			return NullStyle
		}
		p, _ := group.Cascade(key).Get(key)
		return p
	}
	return NullStyle
}

// AddAllFromGroup transfers all style properties from a property group
// to a property map. If overwrite is set, existing style property values
// will be overwritten, otherwise only new values are set.
//
// If the property map does not yet contain a group of this kind, it will
// simply set this group (instead of copying values).
func (pmap *PropertyMap) AddAllFromGroup(group *PropertyGroup, overwrite bool) *PropertyMap {
	if pmap == nil {
		pmap = NewPropertyMap()
	}
	if pmap.m == nil {
		pmap.m = make(map[string]*PropertyGroup)
	}
	g := pmap.Group(group.name)
	if g == nil {
		pmap.m[group.name] = group
	} else {
		for k, v := range group.propsDict {
			if overwrite {
				g.Set(k, v)
			} else {
				g.Add(k, v)
			}
		}
	}
	return pmap
}

// Add adds a property to this property map, e.g.,
//    pm.Add("funny-margin", "big")
func (pmap *PropertyMap) Add(key string, value Property) {
	if pmap == nil {
		return
	}
	groupname := GroupNameFromPropertyKey(key)
	group, found := pmap.m[groupname]
	if !found {
		group = NewPropertyGroup(groupname)
		pmap.m[groupname] = group
	}
	group.Set(key, value)
}

// InitializeDefaultPropertyValues creates an internal data structure to
// hold all the default values for CSS properties.
// In real-world browsers these are the user-agent CSS values.
func InitializeDefaultPropertyValues(additionalProps []KeyValue) *PropertyMap {
	m := make(map[string]*PropertyGroup, 15)
	root := NewPropertyGroup("Root")

	x := NewPropertyGroup(PGX) // special group for extension properties
	for _, kv := range additionalProps {
		x.Set(kv.Key, kv.Value)
	}
	m[PGX] = x

	margins := NewPropertyGroup(PGMargins)
	margins.Set("margin-top", "0")
	margins.Set("margin-left", "0")
	margins.Set("margin-right", "0")
	margins.Set("margin-bottom", "0")
	margins.Parent = root
	m[PGMargins] = margins

	padding := NewPropertyGroup(PGPadding)
	padding.Set("padding-top", "0")
	padding.Set("padding-left", "0")
	padding.Set("padding-right", "0")
	padding.Set("padding-bottom", "0")
	padding.Parent = root
	m[PGPadding] = padding

	border := NewPropertyGroup(PGBorder)
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
	m[PGBorder] = border

	dimension := NewPropertyGroup(PGDimension)
	dimension.Set("width", "10%")
	dimension.Set("width", "100pt")
	dimension.Set("min-width", "0")
	dimension.Set("min-height", "0")
	dimension.Set("max-width", "10000pt")
	dimension.Set("max-height", "10000pt")
	dimension.Parent = root
	m[PGDimension] = dimension

	region := NewPropertyGroup(PGRegion)
	region.Set("flow-from", "")
	region.Set("flow-into", "")
	m[PGRegion] = region

	display := NewPropertyGroup(PGDisplay)
	display.Set("display", "inline")
	display.Set("float", "none")
	display.Set("visibility", "visible")
	display.Set("position", "static")
	m[PGDisplay] = display

	/*
	   type DisplayStyle struct {
	   	Display    uint8 // https://www.tutorialrepublic.com/css-reference/css-display-property.php
	   	Top        dimen.Dimen
	   	Left       dimen.Dimen
	   	Right      dimen.Dimen
	   	Bottom     dimen.Dimen
	   	ZIndex     int
	   	Overflow   uint8
	   	OverflowX  uint8
	   	OverflowY  uint8
	   	Clip       string // geometric shape
	   }

	   type ColorModel string

	   type Color struct {
	   	Color   color.Color
	   	Model   ColorModel
	   	Opacity uint8
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

// --- Package Level Functions ------------------------------------------

/*
func GetCascadedProperty(sn *styledtree.StyNode, key string) (Property, error) {
	groupname := GroupNameFromPropertyKey(key)
	var group *PropertyGroup
	for sn != nil && group == nil {
		group = sn.ComputedStyles().Group(groupname)
		sn = sn.ParentNode()
	}
	if group == nil {
		errmsg := fmt.Sprintf("Cannot find ancestor with prop-group %s -- did you create global properties?", groupname)
		return Property(""), errors.New(errmsg)
	}
	return group.Cascade(key).Get(key), nil // must succeed
}
*/

// GetCascadedProperty gets the value of a property. The search cascades to
// parent property maps, if available.
//
// Clients will usually call GetProperty(…) instead as this will respect
// CSS semantics for inherited properties.
//
// The call to GetCascadedProperty will flag an error if the style property
// isn't found (which should not happen, as every property should be included
// in the 'user-agent' default style properties).
func GetCascadedProperty(n *tree.Node, key string, sty Interf) (Property, error) {
	//
	groupname := GroupNameFromPropertyKey(key)
	var group *PropertyGroup
	for n != nil && group == nil {
		styler := sty(n)
		group = styler.ComputedStyles().Group(groupname)
		n = n.Parent()
	}
	if group == nil {
		errmsg := fmt.Sprintf("Cannot find ancestor with prop-group %s -- did you create global properties?", groupname)
		return NullStyle, errors.New(errmsg)
	}
	p, _ := group.Cascade(key).Get(key)
	return p, nil // must succeed
}

// GetProperty gets the value of a property. If the property is not set
// locally on the style node and the property is inheritable, he search
// cascades to parent property maps, if available.
//
// The call to GetProperty will flag an error if the style property isn't found
// (which should not happen, as every property should be included in the
// 'user-agent' default style properties).
func GetProperty(n *tree.Node, key string, sty Interf) (Property, error) {
	if nonInherited[key] {
		T().Debugf("Property %s is not inherited", key)
		styler := sty(n)
		p, ok := GetLocalProperty(styler.ComputedStyles(), key)
		if !ok {
			p = GetDefaultProperty(styler, key)
		}
		return p, nil
	}
	return GetCascadedProperty(n, key, sty)
}

/*
func GetLocalProperty(sn *styledtree.StyNode, key string) (Property, bool) {
	groupname := GroupNameFromPropertyKey(key)
	var group *PropertyGroup
	group = sn.ComputedStyles().Group(groupname)
	if group == nil {
		return "", false
	}
	return group.Get(key), true
}
*/

// GetLocalProperty returns a style property value, if it is set locally
// for a styled node's property map. No cascading is performed.
func GetLocalProperty(pmap *PropertyMap, key string) (Property, bool) {
	groupname := GroupNameFromPropertyKey(key)
	var group *PropertyGroup
	group = pmap.Group(groupname)
	if group == nil {
		return "", false
	}
	return group.Get(key)
}

// --------------------------------------------------------------------------

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
