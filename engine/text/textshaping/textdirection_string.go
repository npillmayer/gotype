// Code generated by "stringer -type=TextDirection"; DO NOT EDIT.

package textshaping

import "strconv"

const _TextDirection_name = "LeftToRight"

var _TextDirection_index = [...]uint8{0, 11}

func (i TextDirection) String() string {
	if i < 0 || i >= TextDirection(len(_TextDirection_index)-1) {
		return "TextDirection(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TextDirection_name[_TextDirection_index[i]:_TextDirection_index[i+1]]
}
