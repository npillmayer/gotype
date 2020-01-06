package layout

import (
	"fmt"
	"strings"

	"github.com/npillmayer/gotype/engine/dom"
	"golang.org/x/net/html"
)

// Invaluable:
// https://developer.mozilla.org/en-US/docs/Web/CSS/Visual_formatting_model
//
// Regions:
// http://cna.mamk.fi/Public/FJAK/MOAC_MTA_HTML5_App_Dev/c06.pdf

// DisplayModesForDOMNode returns outer and inner display mode for a given DOM node.
func DisplayModesForDOMNode(domnode *dom.W3CNode) (outerMode DisplayMode, innerMode DisplayMode) {
	if domnode == nil || domnode.HTMLNode() == nil {
		return NoMode, NoMode
	}
	if domnode.NodeType() == html.TextNode {
		return InlineMode, InlineMode
	}
	display := domnode.ComputedStyles().GetPropertyValue("display")
	T().Infof("property display=%v", display)
	if display.String() == "initial" {
		outerMode, innerMode = DefaultDisplayModeForHTMLNode(domnode.HTMLNode())
	} else {
		var err error
		outerMode, innerMode, err = ParseDisplay(display.String())
		if err != nil {
			T().Errorf("unrecognized display property: %s", display)
			outerMode, innerMode = BlockMode, BlockMode
		} else if outerMode == NoMode {
			outerMode, innerMode = DefaultDisplayModeForHTMLNode(domnode.HTMLNode())
		}
	}
	T().Infof("display modes = %s | %s", outerMode.String(), innerMode.String())
	return
}

// DefaultDisplayModeForHTMLNode returns the default display mode for a HTML node type,
// as described by the CSS specification.
//
// TODO possibly move this to package style (= part of browser defaults)
// If, then return a string.
func DefaultDisplayModeForHTMLNode(h *html.Node) (DisplayMode, DisplayMode) {
	if h == nil {
		return NoMode, NoMode
	}
	switch h.Type {
	case html.DocumentNode:
		return BlockMode, BlockMode
	case html.TextNode:
		return InlineMode, InlineMode
	case html.ElementNode:
		switch h.Data {
		case "table":
			return BlockMode, TableMode
		case "ul", "ol":
			return BlockMode, ListItemMode
		case "li":
			return ListItemMode, FlowMode | BlockMode
		case "html", "body", "div", "section", "article", "nav":
			return BlockMode, BlockMode
		case "p":
			return BlockMode, InlineMode
		case "span", "i", "b", "strong", "em":
			return InlineMode, InlineMode
		case "h1", "h2", "h3", "h4", "h5", "h6":
			return BlockMode, InlineMode
		default:
			return BlockMode, BlockMode
		}
	default:
		T().Errorf("Have styled node for non-element ?!?")
		T().Errorf(" type of node = %d", h.Type)
		T().Errorf(" name of node = %s", h.Data)
		T().Infof("unknown HTML element will stack children vertically")
		return BlockMode, BlockMode
	}
}

// ---------------------------------------------------------------------------

func dbgNodeString(domnode *dom.W3CNode) string {
	if domnode == nil {
		return "DOM(null)"
	}
	return fmt.Sprintf("DOM(%s/%s)", domnode.NodeName(), shortText(domnode))
}

func shortText(n *dom.W3CNode) string {
	h := n.HTMLNode()
	s := "\""
	if len(h.Data) > 10 {
		s += h.Data[:10] + "...\""
	} else {
		s += h.Data + "\""
	}
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\t", `\t`, -1)
	s = strings.Replace(s, " ", "\u2423", -1)
	return s
}
