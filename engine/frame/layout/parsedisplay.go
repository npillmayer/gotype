package layout

import "fmt"

// ParseDisplay returns mode flags from a display property string (outer and inner).
func ParseDisplay(display string) (DisplayMode, DisplayMode, error) {
	// TODO
	if display == "" {
		return NoMode, NoMode, nil
	}
	switch display {
	case "block":
		return BlockMode, BlockMode, nil
	case "inline":
		return InlineMode, InlineMode, nil
	case "list-item":
		return ListItemMode, FlowMode, nil
	case "inline-block":
		return InlineMode, BlockMode, nil
	case "table":
		return BlockMode, TableMode, nil
	case "inline-table":
		return InlineMode, TableMode, nil
	}
	return NoMode, NoMode, fmt.Errorf("Unknown display mode: %s", display)
}
