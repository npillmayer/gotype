#include <stdlib.h>
#include <stdio.h>
#include <math.h>
#include <hb.h>
#include <hb-ot.h>

/*
Support for low level interfacing to Harfbuzz structures.
Will be called from "harfbuzz_bridge.go".
*/

char strbuffer[1024];

/* Return the code-point (glyph ID) from a glyph_info struct as a string.
 * Please note that the returned string is shared (not threadsafe whatsoever!),
 * as this is used mainly for debugging reasons.
 */
char *get_codepoint_from_glyph_info(hb_font_t *hb_font, hb_glyph_info_t *info, int i) {
    hb_codepoint_t cp = info[i].codepoint;
    if (hb_font != NULL) {
        char glyphname[32];
        hb_font_get_glyph_name(hb_font, cp, glyphname, sizeof(glyphname));
        sprintf(strbuffer, "%04X:%s", cp, glyphname);
    }
    else {
        sprintf(strbuffer, "%04X", cp);
    }
    return strbuffer;
}


/* Helper: get a glyph info struct from an array.
 */
hb_glyph_info_t *get_glyph_info_at(hb_glyph_info_t *info, int i) {
    return &info[i];
}

/* Helper: get a glyph position struct from an array.
 */
hb_glyph_position_t *get_glyph_position_at(hb_glyph_position_t *pos, int i) {
    return &pos[i];
}
