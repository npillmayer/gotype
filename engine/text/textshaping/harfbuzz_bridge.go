package textshaping

//#cgo CPPFLAGS: -I/usr/local/include/harfbuzz
//#cgo LDFLAGS: -L/usr/local/lib -lharfbuzz
/*
#include <stdlib.h>
#include <stdio.h>
#include <math.h>
#include <hb.h>
#include <hb-ot.h>

char *get_codepoint_from_glyph_info(hb_font_t *,hb_glyph_info_t *, int);
hb_glyph_info_t *get_glyph_info_at(hb_glyph_info_t *, int);
hb_glyph_position_t *get_glyph_position_at(hb_glyph_position_t *, int);
*/
import "C"
import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/npillmayer/gotype/core/font"
)

/*
----------------------------------------------------------------------

BSD License

Copyright (c) 2017-2018, Norbert Pillmayer

All rights reserved.

----------------------------------------------------------------------

Helper CGo module as a wrapper around the Harfbuzz text shaping library.

Extract from the Harfbuzz manual page:

Create a buffer and put your text in it.

  #include <hb.h>
  hb_buffer_t *buf;
  buf = hb_buffer_create();
  hb_buffer_add_utf8(buf, text, strlen(text), 0, strlen(text));

Create a face and a font.

  #include <hb-ft.h>
  hb_font_t *font = hb_ft_font_create(...);

Shape!

  hb_shape(font, buf, NULL, 0);

Get the glyph and position information.

  hb_glyph_info_t *glyph_info    = hb_buffer_get_glyph_infos(buf, &glyph_count);
  hb_glyph_position_t *glyph_pos = hb_buffer_get_glyph_positions(buf, &glyph_count);

TODO: make the HB buffer behave as a good GC citizen.
*/

// Allocate the central Harfbuzz data structure and return a (hidden)
// pointer to it.
func allocHBBuffer() uintptr {
	hbbuf := C.hb_buffer_create()
	return uintptr(unsafe.Pointer(hbbuf))
}

func freeHBBuffer(buf uintptr) {
	hbbuf := (*C.struct_hb_buffer_t)(unsafe.Pointer(buf))
	C.hb_buffer_destroy(hbbuf)
}

// Helper: convert a Textdirection enum into a flag suited for Harfbuzz
func dir2hbdir(textdir TextDirection) int32 {
	switch textdir {
	case LeftToRight:
		return 4
	case RightToLeft:
		return 5
	case TopToBottom:
		return 6
	case BottomToTop:
		return 7
	}
	return 4
}

// Set the text direction flag for a Harfbuzz buffer.
func setHBBufferDirection(hbbuf uintptr, dir TextDirection) {
	ptr := (*C.struct_hb_buffer_t)(unsafe.Pointer(hbbuf))
	C.hb_buffer_set_direction(ptr, C.hb_direction_t(dir2hbdir(dir)))
}

// Set the script info for a Harfbuzz buffer.
func setHBBufferScript(hbbuf uintptr, script ScriptID) {
	ptr := (*C.struct_hb_buffer_t)(unsafe.Pointer(hbbuf))
	C.hb_buffer_set_script(ptr, C.hb_script_t(script))
}

// This is Harfbuzz's main function: share a piece of text, using
// a given font. The result of a call to this function will be
// attached to the buffer and may be received by a successive call
// to 'getHBGlyphInfo()'.
func harfbuzzShape(hbbuf uintptr, text string, hbfont uintptr) {
	ptr := (*C.struct_hb_buffer_t)(unsafe.Pointer(hbbuf))
	fptr := (*C.struct_hb_font_t)(unsafe.Pointer(hbfont))
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))
	C.hb_buffer_add_utf8(ptr, cstr, -1, 0, -1)
	C.hb_shape(fptr, ptr, nil, 0)
}

// Harfbuzz uses a different font structure, created from the same
// font binaries we're using in Go.
func makeHBFont(tc *font.TypeCase) uintptr {
	len := len(tc.ScalableFontParent().Binary)
	bytez := C.CBytes(tc.ScalableFontParent().Binary)
	blob := C.hb_blob_create((*C.char)(bytez), (C.uint)(len), C.HB_MEMORY_MODE_WRITABLE, bytez, nil)
	face := C.hb_face_create(blob, 0)
	f := C.hb_font_create(face) // f is of type *C.hb_font_t
	C.hb_ot_font_set_funcs(f)
	sz := (C.int(int(tc.PtSize() * 64.0)))
	C.hb_font_set_scale(f, sz, sz)
	return uintptr(unsafe.Pointer(f))
}

// Retrieve the glyph information from a previous shaper-run.
func getHBGlyphInfo(hbbuf uintptr) *hbGlyphSequence {
	ptr := (*C.struct_hb_buffer_t)(unsafe.Pointer(hbbuf))
	buflen := C.hb_buffer_get_length(ptr)
	info := C.hb_buffer_get_glyph_infos(ptr, nil)
	pos := C.hb_buffer_get_glyph_positions(ptr, nil)
	if info == nil || pos == nil {
		fmt.Printf("*** error: info or pos is zero\n")
	}
	seq := &hbGlyphSequence{
		length: int(buflen),
		info:   info,
		pos:    pos,
	}
	return seq
}

// Go-Container for a Harfbuzz-result.
type hbGlyphSequence struct {
	length int
	info   *C.hb_glyph_info_t
	pos    *C.hb_glyph_position_t
}

// Implement the GlyphSequence interface
func (seq *hbGlyphSequence) GlyphCount() int {
	return seq.length
}

type hbGlyphInfo struct {
	glyph    rune
	cluster  int
	xadvance float64
	yadvance float64
	x        float64
	y        float64
}

// Implement the GlyphSequence interface
func (seq *hbGlyphSequence) GetGlyphInfoAt(i int) GlyphInfo {
	gi := &hbGlyphInfo{}
	info := C.get_glyph_info_at(seq.info, C.int(i))
	pos := C.get_glyph_position_at(seq.pos, C.int(i))
	gi.glyph = rune(info.codepoint)
	gi.cluster = int(info.cluster)
	gi.xadvance = float64(pos.x_advance) / 64.0
	gi.yadvance = float64(pos.y_advance) / 64.0
	gi.x = float64(pos.x_offset) / 64.0
	gi.y = float64(pos.y_offset) / 64.0
	return gi
}

// Implement the GlyphInfo interface
func (gi *hbGlyphInfo) Glyph() rune {
	return gi.glyph
}

// Implement the GlyphInfo interface
func (gi *hbGlyphInfo) Cluster() int {
	return gi.cluster
}

// Implement the GlyphInfo interface
func (gi *hbGlyphInfo) XAdvance() float64 {
	return gi.xadvance
}

// Implement the GlyphInfo interface
func (gi *hbGlyphInfo) YAdvance() float64 {
	return gi.yadvance
}

// Implement the GlyphInfo interface
func (gi *hbGlyphInfo) XPosition() float64 {
	return gi.x
}

// Implement the GlyphInfo interface
func (gi *hbGlyphInfo) YPosition() float64 {
	return gi.y
}

// For debugging purposes: string representation of a glyph sequence,
// displaying code-points/glyph-IDs.
func (seq *hbGlyphSequence) String() string {
	var sb strings.Builder
	for i := 0; i < seq.length; i++ {
		s := C.get_codepoint_from_glyph_info(nil, seq.info, C.int(i))
		if i > 0 {
			sb.WriteString("|")
		}
		sb.WriteString(C.GoString(s))
	}
	return sb.String()
}

// For debugging purposes: string representation of a glyph sequence.
// Includes the glyphs' names as provided by the font.
func hbGlyphString(hbfont uintptr, seq *hbGlyphSequence) string {
	fptr := (*C.struct_hb_font_t)(unsafe.Pointer(hbfont))
	var sb strings.Builder
	for i := 0; i < seq.length; i++ {
		s := C.get_codepoint_from_glyph_info(fptr, seq.info, C.int(i))
		if i > 0 {
			sb.WriteString("|")
		}
		sb.WriteString(C.GoString(s))
	}
	return sb.String()
}
