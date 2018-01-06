// Generated from Gallery.g4 by ANTLR 4.7.

package grammar // Gallery
import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = reflect.Copy
var _ = strconv.Itoa

var parserATN = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 48, 274,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	3, 2, 3, 2, 3, 2, 7, 2, 50, 10, 2, 12, 2, 14, 2, 53, 11, 2, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 5, 3, 60, 10, 3, 3, 4, 3, 4, 3, 4, 3, 4, 7, 4, 66, 10,
	4, 12, 4, 14, 4, 69, 11, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3,
	6, 3, 7, 3, 7, 5, 7, 81, 10, 7, 3, 8, 3, 8, 3, 8, 6, 8, 86, 10, 8, 13,
	8, 14, 8, 87, 3, 9, 3, 9, 3, 9, 6, 9, 93, 10, 9, 13, 9, 14, 9, 94, 3, 10,
	3, 10, 3, 10, 3, 10, 7, 10, 101, 10, 10, 12, 10, 14, 10, 104, 11, 10, 3,
	10, 3, 10, 3, 10, 3, 10, 7, 10, 110, 10, 10, 12, 10, 14, 10, 113, 11, 10,
	3, 10, 3, 10, 5, 10, 117, 10, 10, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3,
	11, 7, 11, 125, 10, 11, 12, 11, 14, 11, 128, 11, 11, 3, 12, 3, 12, 3, 12,
	5, 12, 133, 10, 12, 3, 12, 3, 12, 3, 12, 7, 12, 138, 10, 12, 12, 12, 14,
	12, 141, 11, 12, 3, 13, 3, 13, 3, 13, 6, 13, 146, 10, 13, 13, 13, 14, 13,
	147, 3, 13, 5, 13, 151, 10, 13, 3, 14, 3, 14, 3, 14, 3, 15, 3, 15, 3, 15,
	3, 15, 3, 15, 3, 15, 3, 15, 3, 15, 7, 15, 164, 10, 15, 12, 15, 14, 15,
	167, 11, 15, 3, 16, 3, 16, 6, 16, 171, 10, 16, 13, 16, 14, 16, 172, 3,
	17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17,
	3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3,
	17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17,
	3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 5, 17, 213, 10, 17, 3, 18, 3,
	18, 5, 18, 217, 10, 18, 3, 19, 3, 19, 3, 19, 3, 19, 5, 19, 223, 10, 19,
	3, 20, 3, 20, 5, 20, 227, 10, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3,
	20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20, 3, 20,
	5, 20, 245, 10, 20, 3, 21, 3, 21, 3, 21, 7, 21, 250, 10, 21, 12, 21, 14,
	21, 253, 11, 21, 3, 21, 3, 21, 3, 21, 7, 21, 258, 10, 21, 12, 21, 14, 21,
	261, 11, 21, 5, 21, 263, 10, 21, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 5,
	22, 270, 10, 22, 3, 23, 3, 23, 3, 23, 2, 5, 20, 22, 28, 24, 2, 4, 6, 8,
	10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44,
	2, 7, 3, 2, 17, 19, 3, 2, 13, 14, 3, 2, 15, 16, 3, 2, 27, 28, 3, 2, 43,
	44, 2, 295, 2, 51, 3, 2, 2, 2, 4, 59, 3, 2, 2, 2, 6, 61, 3, 2, 2, 2, 8,
	70, 3, 2, 2, 2, 10, 74, 3, 2, 2, 2, 12, 80, 3, 2, 2, 2, 14, 82, 3, 2, 2,
	2, 16, 89, 3, 2, 2, 2, 18, 116, 3, 2, 2, 2, 20, 118, 3, 2, 2, 2, 22, 132,
	3, 2, 2, 2, 24, 142, 3, 2, 2, 2, 26, 152, 3, 2, 2, 2, 28, 155, 3, 2, 2,
	2, 30, 170, 3, 2, 2, 2, 32, 212, 3, 2, 2, 2, 34, 216, 3, 2, 2, 2, 36, 222,
	3, 2, 2, 2, 38, 244, 3, 2, 2, 2, 40, 262, 3, 2, 2, 2, 42, 269, 3, 2, 2,
	2, 44, 271, 3, 2, 2, 2, 46, 47, 5, 4, 3, 2, 47, 48, 7, 6, 2, 2, 48, 50,
	3, 2, 2, 2, 49, 46, 3, 2, 2, 2, 50, 53, 3, 2, 2, 2, 51, 49, 3, 2, 2, 2,
	51, 52, 3, 2, 2, 2, 52, 3, 3, 2, 2, 2, 53, 51, 3, 2, 2, 2, 54, 60, 5, 6,
	4, 2, 55, 60, 5, 8, 5, 2, 56, 60, 5, 10, 6, 2, 57, 60, 5, 12, 7, 2, 58,
	60, 5, 18, 10, 2, 59, 54, 3, 2, 2, 2, 59, 55, 3, 2, 2, 2, 59, 56, 3, 2,
	2, 2, 59, 57, 3, 2, 2, 2, 59, 58, 3, 2, 2, 2, 60, 5, 3, 2, 2, 2, 61, 62,
	7, 22, 2, 2, 62, 67, 7, 43, 2, 2, 63, 64, 7, 7, 2, 2, 64, 66, 7, 43, 2,
	2, 65, 63, 3, 2, 2, 2, 66, 69, 3, 2, 2, 2, 67, 65, 3, 2, 2, 2, 67, 68,
	3, 2, 2, 2, 68, 7, 3, 2, 2, 2, 69, 67, 3, 2, 2, 2, 70, 71, 5, 40, 21, 2,
	71, 72, 7, 3, 2, 2, 72, 73, 5, 20, 11, 2, 73, 9, 3, 2, 2, 2, 74, 75, 7,
	38, 2, 2, 75, 76, 5, 2, 2, 2, 76, 77, 7, 39, 2, 2, 77, 11, 3, 2, 2, 2,
	78, 81, 5, 14, 8, 2, 79, 81, 5, 16, 9, 2, 80, 78, 3, 2, 2, 2, 80, 79, 3,
	2, 2, 2, 81, 13, 3, 2, 2, 2, 82, 85, 5, 20, 11, 2, 83, 84, 7, 4, 2, 2,
	84, 86, 5, 20, 11, 2, 85, 83, 3, 2, 2, 2, 86, 87, 3, 2, 2, 2, 87, 85, 3,
	2, 2, 2, 87, 88, 3, 2, 2, 2, 88, 15, 3, 2, 2, 2, 89, 92, 5, 22, 12, 2,
	90, 91, 9, 2, 2, 2, 91, 93, 5, 22, 12, 2, 92, 90, 3, 2, 2, 2, 93, 94, 3,
	2, 2, 2, 94, 92, 3, 2, 2, 2, 94, 95, 3, 2, 2, 2, 95, 17, 3, 2, 2, 2, 96,
	97, 7, 41, 2, 2, 97, 102, 7, 43, 2, 2, 98, 99, 7, 7, 2, 2, 99, 101, 7,
	43, 2, 2, 100, 98, 3, 2, 2, 2, 101, 104, 3, 2, 2, 2, 102, 100, 3, 2, 2,
	2, 102, 103, 3, 2, 2, 2, 103, 117, 3, 2, 2, 2, 104, 102, 3, 2, 2, 2, 105,
	106, 7, 42, 2, 2, 106, 111, 7, 43, 2, 2, 107, 108, 7, 7, 2, 2, 108, 110,
	7, 43, 2, 2, 109, 107, 3, 2, 2, 2, 110, 113, 3, 2, 2, 2, 111, 109, 3, 2,
	2, 2, 111, 112, 3, 2, 2, 2, 112, 117, 3, 2, 2, 2, 113, 111, 3, 2, 2, 2,
	114, 115, 7, 40, 2, 2, 115, 117, 7, 47, 2, 2, 116, 96, 3, 2, 2, 2, 116,
	105, 3, 2, 2, 2, 116, 114, 3, 2, 2, 2, 117, 19, 3, 2, 2, 2, 118, 119, 8,
	11, 1, 2, 119, 120, 5, 22, 12, 2, 120, 126, 3, 2, 2, 2, 121, 122, 12, 3,
	2, 2, 122, 123, 7, 37, 2, 2, 123, 125, 5, 22, 12, 2, 124, 121, 3, 2, 2,
	2, 125, 128, 3, 2, 2, 2, 126, 124, 3, 2, 2, 2, 126, 127, 3, 2, 2, 2, 127,
	21, 3, 2, 2, 2, 128, 126, 3, 2, 2, 2, 129, 130, 8, 12, 1, 2, 130, 133,
	5, 28, 15, 2, 131, 133, 5, 24, 13, 2, 132, 129, 3, 2, 2, 2, 132, 131, 3,
	2, 2, 2, 133, 139, 3, 2, 2, 2, 134, 135, 12, 4, 2, 2, 135, 136, 9, 3, 2,
	2, 136, 138, 5, 28, 15, 2, 137, 134, 3, 2, 2, 2, 138, 141, 3, 2, 2, 2,
	139, 137, 3, 2, 2, 2, 139, 140, 3, 2, 2, 2, 140, 23, 3, 2, 2, 2, 141, 139,
	3, 2, 2, 2, 142, 145, 5, 28, 15, 2, 143, 144, 7, 12, 2, 2, 144, 146, 5,
	28, 15, 2, 145, 143, 3, 2, 2, 2, 146, 147, 3, 2, 2, 2, 147, 145, 3, 2,
	2, 2, 147, 148, 3, 2, 2, 2, 148, 150, 3, 2, 2, 2, 149, 151, 5, 26, 14,
	2, 150, 149, 3, 2, 2, 2, 150, 151, 3, 2, 2, 2, 151, 25, 3, 2, 2, 2, 152,
	153, 7, 12, 2, 2, 153, 154, 7, 36, 2, 2, 154, 27, 3, 2, 2, 2, 155, 156,
	8, 15, 1, 2, 156, 157, 5, 32, 17, 2, 157, 165, 3, 2, 2, 2, 158, 159, 12,
	4, 2, 2, 159, 160, 9, 4, 2, 2, 160, 164, 5, 32, 17, 2, 161, 162, 12, 3,
	2, 2, 162, 164, 5, 30, 16, 2, 163, 158, 3, 2, 2, 2, 163, 161, 3, 2, 2,
	2, 164, 167, 3, 2, 2, 2, 165, 163, 3, 2, 2, 2, 165, 166, 3, 2, 2, 2, 166,
	29, 3, 2, 2, 2, 167, 165, 3, 2, 2, 2, 168, 169, 7, 35, 2, 2, 169, 171,
	5, 32, 17, 2, 170, 168, 3, 2, 2, 2, 171, 172, 3, 2, 2, 2, 172, 170, 3,
	2, 2, 2, 172, 173, 3, 2, 2, 2, 173, 31, 3, 2, 2, 2, 174, 175, 7, 29, 2,
	2, 175, 213, 5, 38, 20, 2, 176, 177, 5, 34, 18, 2, 177, 178, 5, 38, 20,
	2, 178, 213, 3, 2, 2, 2, 179, 180, 5, 36, 19, 2, 180, 181, 7, 10, 2, 2,
	181, 182, 5, 22, 12, 2, 182, 183, 7, 7, 2, 2, 183, 184, 5, 22, 12, 2, 184,
	185, 7, 11, 2, 2, 185, 213, 3, 2, 2, 2, 186, 187, 5, 38, 20, 2, 187, 188,
	7, 10, 2, 2, 188, 189, 5, 22, 12, 2, 189, 190, 7, 7, 2, 2, 190, 191, 5,
	22, 12, 2, 191, 192, 7, 11, 2, 2, 192, 213, 3, 2, 2, 2, 193, 213, 5, 38,
	20, 2, 194, 195, 7, 24, 2, 2, 195, 213, 5, 32, 17, 2, 196, 197, 7, 33,
	2, 2, 197, 198, 5, 22, 12, 2, 198, 199, 7, 34, 2, 2, 199, 200, 5, 32, 17,
	2, 200, 213, 3, 2, 2, 2, 201, 202, 7, 31, 2, 2, 202, 213, 5, 32, 17, 2,
	203, 204, 7, 30, 2, 2, 204, 205, 5, 22, 12, 2, 205, 206, 7, 34, 2, 2, 206,
	207, 5, 32, 17, 2, 207, 213, 3, 2, 2, 2, 208, 209, 7, 25, 2, 2, 209, 213,
	5, 32, 17, 2, 210, 211, 9, 5, 2, 2, 211, 213, 5, 40, 21, 2, 212, 174, 3,
	2, 2, 2, 212, 176, 3, 2, 2, 2, 212, 179, 3, 2, 2, 2, 212, 186, 3, 2, 2,
	2, 212, 193, 3, 2, 2, 2, 212, 194, 3, 2, 2, 2, 212, 196, 3, 2, 2, 2, 212,
	201, 3, 2, 2, 2, 212, 203, 3, 2, 2, 2, 212, 208, 3, 2, 2, 2, 212, 210,
	3, 2, 2, 2, 213, 33, 3, 2, 2, 2, 214, 217, 9, 3, 2, 2, 215, 217, 5, 36,
	19, 2, 216, 214, 3, 2, 2, 2, 216, 215, 3, 2, 2, 2, 217, 35, 3, 2, 2, 2,
	218, 219, 7, 45, 2, 2, 219, 220, 7, 16, 2, 2, 220, 223, 7, 45, 2, 2, 221,
	223, 7, 45, 2, 2, 222, 218, 3, 2, 2, 2, 222, 221, 3, 2, 2, 2, 223, 37,
	3, 2, 2, 2, 224, 226, 7, 45, 2, 2, 225, 227, 7, 21, 2, 2, 226, 225, 3,
	2, 2, 2, 226, 227, 3, 2, 2, 2, 227, 245, 3, 2, 2, 2, 228, 245, 5, 40, 21,
	2, 229, 230, 7, 8, 2, 2, 230, 231, 5, 22, 12, 2, 231, 232, 7, 7, 2, 2,
	232, 233, 5, 22, 12, 2, 233, 234, 7, 9, 2, 2, 234, 245, 3, 2, 2, 2, 235,
	236, 7, 8, 2, 2, 236, 237, 5, 22, 12, 2, 237, 238, 7, 9, 2, 2, 238, 245,
	3, 2, 2, 2, 239, 240, 7, 38, 2, 2, 240, 241, 5, 2, 2, 2, 241, 242, 5, 22,
	12, 2, 242, 243, 7, 39, 2, 2, 243, 245, 3, 2, 2, 2, 244, 224, 3, 2, 2,
	2, 244, 228, 3, 2, 2, 2, 244, 229, 3, 2, 2, 2, 244, 235, 3, 2, 2, 2, 244,
	239, 3, 2, 2, 2, 245, 39, 3, 2, 2, 2, 246, 251, 7, 44, 2, 2, 247, 250,
	5, 42, 22, 2, 248, 250, 5, 44, 23, 2, 249, 247, 3, 2, 2, 2, 249, 248, 3,
	2, 2, 2, 250, 253, 3, 2, 2, 2, 251, 249, 3, 2, 2, 2, 251, 252, 3, 2, 2,
	2, 252, 263, 3, 2, 2, 2, 253, 251, 3, 2, 2, 2, 254, 259, 7, 43, 2, 2, 255,
	258, 5, 42, 22, 2, 256, 258, 5, 44, 23, 2, 257, 255, 3, 2, 2, 2, 257, 256,
	3, 2, 2, 2, 258, 261, 3, 2, 2, 2, 259, 257, 3, 2, 2, 2, 259, 260, 3, 2,
	2, 2, 260, 263, 3, 2, 2, 2, 261, 259, 3, 2, 2, 2, 262, 246, 3, 2, 2, 2,
	262, 254, 3, 2, 2, 2, 263, 41, 3, 2, 2, 2, 264, 270, 7, 45, 2, 2, 265,
	266, 7, 10, 2, 2, 266, 267, 5, 22, 12, 2, 267, 268, 7, 11, 2, 2, 268, 270,
	3, 2, 2, 2, 269, 264, 3, 2, 2, 2, 269, 265, 3, 2, 2, 2, 270, 43, 3, 2,
	2, 2, 271, 272, 9, 6, 2, 2, 272, 45, 3, 2, 2, 2, 30, 51, 59, 67, 80, 87,
	94, 102, 111, 116, 126, 132, 139, 147, 150, 163, 165, 172, 212, 216, 222,
	226, 244, 249, 251, 257, 259, 262, 269,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "':='", "'='", "':'", "';'", "','", "'('", "')'", "'['", "']'", "",
	"'+'", "'-'", "'*'", "'/'", "'||'", "'|-'", "'~'", "", "", "", "", "",
	"", "'edge'", "'frame'", "'box'", "", "'subpath'", "'reverse'", "'with'",
	"'point'", "'of'", "", "'cycle'", "", "'begingroup'", "'endgroup'", "'proof'",
	"'save'", "'show'", "", "", "", "'.'",
}
var symbolicNames = []string{
	"", "ASSIGN", "EQUALS", "COLON", "SEMIC", "COMMA", "LPAREN", "RPAREN",
	"LBRACKET", "RBRACKET", "PATHJOIN", "PLUS", "MINUS", "TIMES", "OVER", "PARALLEL",
	"PERPENDIC", "CONGRUENT", "WHATEVER", "UNIT", "TYPE", "INTERNAL", "PAIRPART",
	"EDGECONSTR", "EDGE", "FRAME", "BOX", "MATHFUNC", "SUBPATH", "REVERSE",
	"WITH", "POINT", "OF", "TRANSFORM", "CYCLE", "PATHCLIPOP", "BEGINGROUP",
	"ENDGROUP", "PROOF", "SAVE", "SHOW", "TAG", "MIXEDTAG", "DECIMALTOKEN",
	"DOT", "LABEL", "WS",
}

var ruleNames = []string{
	"statementlist", "statement", "declaration", "assignment", "compound",
	"constraint", "equation", "orientation", "command", "expression", "tertiary",
	"path", "cycle", "secondary", "transformer", "primary", "scalarmulop",
	"numtokenatom", "atom", "variable", "subscript", "anytag",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type GalleryParser struct {
	*antlr.BaseParser
}

func NewGalleryParser(input antlr.TokenStream) *GalleryParser {
	this := new(GalleryParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "Gallery.g4"

	return this
}

// GalleryParser tokens.
const (
	GalleryParserEOF          = antlr.TokenEOF
	GalleryParserASSIGN       = 1
	GalleryParserEQUALS       = 2
	GalleryParserCOLON        = 3
	GalleryParserSEMIC        = 4
	GalleryParserCOMMA        = 5
	GalleryParserLPAREN       = 6
	GalleryParserRPAREN       = 7
	GalleryParserLBRACKET     = 8
	GalleryParserRBRACKET     = 9
	GalleryParserPATHJOIN     = 10
	GalleryParserPLUS         = 11
	GalleryParserMINUS        = 12
	GalleryParserTIMES        = 13
	GalleryParserOVER         = 14
	GalleryParserPARALLEL     = 15
	GalleryParserPERPENDIC    = 16
	GalleryParserCONGRUENT    = 17
	GalleryParserWHATEVER     = 18
	GalleryParserUNIT         = 19
	GalleryParserTYPE         = 20
	GalleryParserINTERNAL     = 21
	GalleryParserPAIRPART     = 22
	GalleryParserEDGECONSTR   = 23
	GalleryParserEDGE         = 24
	GalleryParserFRAME        = 25
	GalleryParserBOX          = 26
	GalleryParserMATHFUNC     = 27
	GalleryParserSUBPATH      = 28
	GalleryParserREVERSE      = 29
	GalleryParserWITH         = 30
	GalleryParserPOINT        = 31
	GalleryParserOF           = 32
	GalleryParserTRANSFORM    = 33
	GalleryParserCYCLE        = 34
	GalleryParserPATHCLIPOP   = 35
	GalleryParserBEGINGROUP   = 36
	GalleryParserENDGROUP     = 37
	GalleryParserPROOF        = 38
	GalleryParserSAVE         = 39
	GalleryParserSHOW         = 40
	GalleryParserTAG          = 41
	GalleryParserMIXEDTAG     = 42
	GalleryParserDECIMALTOKEN = 43
	GalleryParserDOT          = 44
	GalleryParserLABEL        = 45
	GalleryParserWS           = 46
)

// GalleryParser rules.
const (
	GalleryParserRULE_statementlist = 0
	GalleryParserRULE_statement     = 1
	GalleryParserRULE_declaration   = 2
	GalleryParserRULE_assignment    = 3
	GalleryParserRULE_compound      = 4
	GalleryParserRULE_constraint    = 5
	GalleryParserRULE_equation      = 6
	GalleryParserRULE_orientation   = 7
	GalleryParserRULE_command       = 8
	GalleryParserRULE_expression    = 9
	GalleryParserRULE_tertiary      = 10
	GalleryParserRULE_path          = 11
	GalleryParserRULE_cycle         = 12
	GalleryParserRULE_secondary     = 13
	GalleryParserRULE_transformer   = 14
	GalleryParserRULE_primary       = 15
	GalleryParserRULE_scalarmulop   = 16
	GalleryParserRULE_numtokenatom  = 17
	GalleryParserRULE_atom          = 18
	GalleryParserRULE_variable      = 19
	GalleryParserRULE_subscript     = 20
	GalleryParserRULE_anytag        = 21
)

// IStatementlistContext is an interface to support dynamic dispatch.
type IStatementlistContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatementlistContext differentiates from other interfaces.
	IsStatementlistContext()
}

type StatementlistContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementlistContext() *StatementlistContext {
	var p = new(StatementlistContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_statementlist
	return p
}

func (*StatementlistContext) IsStatementlistContext() {}

func NewStatementlistContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementlistContext {
	var p = new(StatementlistContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_statementlist

	return p
}

func (s *StatementlistContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementlistContext) AllStatement() []IStatementContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IStatementContext)(nil)).Elem())
	var tst = make([]IStatementContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IStatementContext)
		}
	}

	return tst
}

func (s *StatementlistContext) Statement(i int) IStatementContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *StatementlistContext) AllSEMIC() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserSEMIC)
}

func (s *StatementlistContext) SEMIC(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserSEMIC, i)
}

func (s *StatementlistContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementlistContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementlistContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterStatementlist(s)
	}
}

func (s *StatementlistContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitStatementlist(s)
	}
}

func (p *GalleryParser) Statementlist() (localctx IStatementlistContext) {
	localctx = NewStatementlistContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, GalleryParserRULE_statementlist)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(49)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(44)
				p.Statement()
			}
			{
				p.SetState(45)
				p.Match(GalleryParserSEMIC)
			}

		}
		p.SetState(51)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext())
	}

	return localctx
}

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_statement
	return p
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) Declaration() IDeclarationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDeclarationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDeclarationContext)
}

func (s *StatementContext) Assignment() IAssignmentContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAssignmentContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAssignmentContext)
}

func (s *StatementContext) Compound() ICompoundContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICompoundContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICompoundContext)
}

func (s *StatementContext) Constraint() IConstraintContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IConstraintContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IConstraintContext)
}

func (s *StatementContext) Command() ICommandContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICommandContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICommandContext)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (p *GalleryParser) Statement() (localctx IStatementContext) {
	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, GalleryParserRULE_statement)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(57)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(52)
			p.Declaration()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(53)
			p.Assignment()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(54)
			p.Compound()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(55)
			p.Constraint()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(56)
			p.Command()
		}

	}

	return localctx
}

// IDeclarationContext is an interface to support dynamic dispatch.
type IDeclarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDeclarationContext differentiates from other interfaces.
	IsDeclarationContext()
}

type DeclarationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDeclarationContext() *DeclarationContext {
	var p = new(DeclarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_declaration
	return p
}

func (*DeclarationContext) IsDeclarationContext() {}

func NewDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclarationContext {
	var p = new(DeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_declaration

	return p
}

func (s *DeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *DeclarationContext) TYPE() antlr.TerminalNode {
	return s.GetToken(GalleryParserTYPE, 0)
}

func (s *DeclarationContext) AllTAG() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserTAG)
}

func (s *DeclarationContext) TAG(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserTAG, i)
}

func (s *DeclarationContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserCOMMA)
}

func (s *DeclarationContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserCOMMA, i)
}

func (s *DeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterDeclaration(s)
	}
}

func (s *DeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitDeclaration(s)
	}
}

func (p *GalleryParser) Declaration() (localctx IDeclarationContext) {
	localctx = NewDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, GalleryParserRULE_declaration)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(59)
		p.Match(GalleryParserTYPE)
	}
	{
		p.SetState(60)
		p.Match(GalleryParserTAG)
	}
	p.SetState(65)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == GalleryParserCOMMA {
		{
			p.SetState(61)
			p.Match(GalleryParserCOMMA)
		}
		{
			p.SetState(62)
			p.Match(GalleryParserTAG)
		}

		p.SetState(67)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IAssignmentContext is an interface to support dynamic dispatch.
type IAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAssignmentContext differentiates from other interfaces.
	IsAssignmentContext()
}

type AssignmentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAssignmentContext() *AssignmentContext {
	var p = new(AssignmentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_assignment
	return p
}

func (*AssignmentContext) IsAssignmentContext() {}

func NewAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentContext {
	var p = new(AssignmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_assignment

	return p
}

func (s *AssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *AssignmentContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *AssignmentContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(GalleryParserASSIGN, 0)
}

func (s *AssignmentContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *AssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterAssignment(s)
	}
}

func (s *AssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitAssignment(s)
	}
}

func (p *GalleryParser) Assignment() (localctx IAssignmentContext) {
	localctx = NewAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, GalleryParserRULE_assignment)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(68)
		p.Variable()
	}
	{
		p.SetState(69)
		p.Match(GalleryParserASSIGN)
	}
	{
		p.SetState(70)
		p.expression(0)
	}

	return localctx
}

// ICompoundContext is an interface to support dynamic dispatch.
type ICompoundContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCompoundContext differentiates from other interfaces.
	IsCompoundContext()
}

type CompoundContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompoundContext() *CompoundContext {
	var p = new(CompoundContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_compound
	return p
}

func (*CompoundContext) IsCompoundContext() {}

func NewCompoundContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompoundContext {
	var p = new(CompoundContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_compound

	return p
}

func (s *CompoundContext) GetParser() antlr.Parser { return s.parser }

func (s *CompoundContext) BEGINGROUP() antlr.TerminalNode {
	return s.GetToken(GalleryParserBEGINGROUP, 0)
}

func (s *CompoundContext) Statementlist() IStatementlistContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementlistContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStatementlistContext)
}

func (s *CompoundContext) ENDGROUP() antlr.TerminalNode {
	return s.GetToken(GalleryParserENDGROUP, 0)
}

func (s *CompoundContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompoundContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompoundContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterCompound(s)
	}
}

func (s *CompoundContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitCompound(s)
	}
}

func (p *GalleryParser) Compound() (localctx ICompoundContext) {
	localctx = NewCompoundContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, GalleryParserRULE_compound)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(72)
		p.Match(GalleryParserBEGINGROUP)
	}
	{
		p.SetState(73)
		p.Statementlist()
	}
	{
		p.SetState(74)
		p.Match(GalleryParserENDGROUP)
	}

	return localctx
}

// IConstraintContext is an interface to support dynamic dispatch.
type IConstraintContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsConstraintContext differentiates from other interfaces.
	IsConstraintContext()
}

type ConstraintContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConstraintContext() *ConstraintContext {
	var p = new(ConstraintContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_constraint
	return p
}

func (*ConstraintContext) IsConstraintContext() {}

func NewConstraintContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstraintContext {
	var p = new(ConstraintContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_constraint

	return p
}

func (s *ConstraintContext) GetParser() antlr.Parser { return s.parser }

func (s *ConstraintContext) Equation() IEquationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEquationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEquationContext)
}

func (s *ConstraintContext) Orientation() IOrientationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IOrientationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IOrientationContext)
}

func (s *ConstraintContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConstraintContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConstraintContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterConstraint(s)
	}
}

func (s *ConstraintContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitConstraint(s)
	}
}

func (p *GalleryParser) Constraint() (localctx IConstraintContext) {
	localctx = NewConstraintContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, GalleryParserRULE_constraint)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(78)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(76)
			p.Equation()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(77)
			p.Orientation()
		}

	}

	return localctx
}

// IEquationContext is an interface to support dynamic dispatch.
type IEquationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEquationContext differentiates from other interfaces.
	IsEquationContext()
}

type EquationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEquationContext() *EquationContext {
	var p = new(EquationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_equation
	return p
}

func (*EquationContext) IsEquationContext() {}

func NewEquationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EquationContext {
	var p = new(EquationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_equation

	return p
}

func (s *EquationContext) GetParser() antlr.Parser { return s.parser }

func (s *EquationContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *EquationContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *EquationContext) AllEQUALS() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserEQUALS)
}

func (s *EquationContext) EQUALS(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserEQUALS, i)
}

func (s *EquationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EquationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EquationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterEquation(s)
	}
}

func (s *EquationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitEquation(s)
	}
}

func (p *GalleryParser) Equation() (localctx IEquationContext) {
	localctx = NewEquationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, GalleryParserRULE_equation)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(80)
		p.expression(0)
	}
	p.SetState(83)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == GalleryParserEQUALS {
		{
			p.SetState(81)
			p.Match(GalleryParserEQUALS)
		}
		{
			p.SetState(82)
			p.expression(0)
		}

		p.SetState(85)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IOrientationContext is an interface to support dynamic dispatch.
type IOrientationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsOrientationContext differentiates from other interfaces.
	IsOrientationContext()
}

type OrientationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrientationContext() *OrientationContext {
	var p = new(OrientationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_orientation
	return p
}

func (*OrientationContext) IsOrientationContext() {}

func NewOrientationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrientationContext {
	var p = new(OrientationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_orientation

	return p
}

func (s *OrientationContext) GetParser() antlr.Parser { return s.parser }

func (s *OrientationContext) AllTertiary() []ITertiaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITertiaryContext)(nil)).Elem())
	var tst = make([]ITertiaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITertiaryContext)
		}
	}

	return tst
}

func (s *OrientationContext) Tertiary(i int) ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *OrientationContext) AllPARALLEL() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserPARALLEL)
}

func (s *OrientationContext) PARALLEL(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserPARALLEL, i)
}

func (s *OrientationContext) AllPERPENDIC() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserPERPENDIC)
}

func (s *OrientationContext) PERPENDIC(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserPERPENDIC, i)
}

func (s *OrientationContext) AllCONGRUENT() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserCONGRUENT)
}

func (s *OrientationContext) CONGRUENT(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserCONGRUENT, i)
}

func (s *OrientationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrientationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OrientationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterOrientation(s)
	}
}

func (s *OrientationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitOrientation(s)
	}
}

func (p *GalleryParser) Orientation() (localctx IOrientationContext) {
	localctx = NewOrientationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, GalleryParserRULE_orientation)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(87)
		p.tertiary(0)
	}
	p.SetState(90)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<GalleryParserPARALLEL)|(1<<GalleryParserPERPENDIC)|(1<<GalleryParserCONGRUENT))) != 0) {
		p.SetState(88)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<GalleryParserPARALLEL)|(1<<GalleryParserPERPENDIC)|(1<<GalleryParserCONGRUENT))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
		{
			p.SetState(89)
			p.tertiary(0)
		}

		p.SetState(92)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ICommandContext is an interface to support dynamic dispatch.
type ICommandContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCommandContext differentiates from other interfaces.
	IsCommandContext()
}

type CommandContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCommandContext() *CommandContext {
	var p = new(CommandContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_command
	return p
}

func (*CommandContext) IsCommandContext() {}

func NewCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandContext {
	var p = new(CommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_command

	return p
}

func (s *CommandContext) GetParser() antlr.Parser { return s.parser }

func (s *CommandContext) CopyFrom(ctx *CommandContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *CommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ShowcmdContext struct {
	*CommandContext
}

func NewShowcmdContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ShowcmdContext {
	var p = new(ShowcmdContext)

	p.CommandContext = NewEmptyCommandContext()
	p.parser = parser
	p.CopyFrom(ctx.(*CommandContext))

	return p
}

func (s *ShowcmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowcmdContext) SHOW() antlr.TerminalNode {
	return s.GetToken(GalleryParserSHOW, 0)
}

func (s *ShowcmdContext) AllTAG() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserTAG)
}

func (s *ShowcmdContext) TAG(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserTAG, i)
}

func (s *ShowcmdContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserCOMMA)
}

func (s *ShowcmdContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserCOMMA, i)
}

func (s *ShowcmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterShowcmd(s)
	}
}

func (s *ShowcmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitShowcmd(s)
	}
}

type ProofcmdContext struct {
	*CommandContext
}

func NewProofcmdContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ProofcmdContext {
	var p = new(ProofcmdContext)

	p.CommandContext = NewEmptyCommandContext()
	p.parser = parser
	p.CopyFrom(ctx.(*CommandContext))

	return p
}

func (s *ProofcmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProofcmdContext) PROOF() antlr.TerminalNode {
	return s.GetToken(GalleryParserPROOF, 0)
}

func (s *ProofcmdContext) LABEL() antlr.TerminalNode {
	return s.GetToken(GalleryParserLABEL, 0)
}

func (s *ProofcmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterProofcmd(s)
	}
}

func (s *ProofcmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitProofcmd(s)
	}
}

type SavecmdContext struct {
	*CommandContext
}

func NewSavecmdContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SavecmdContext {
	var p = new(SavecmdContext)

	p.CommandContext = NewEmptyCommandContext()
	p.parser = parser
	p.CopyFrom(ctx.(*CommandContext))

	return p
}

func (s *SavecmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SavecmdContext) SAVE() antlr.TerminalNode {
	return s.GetToken(GalleryParserSAVE, 0)
}

func (s *SavecmdContext) AllTAG() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserTAG)
}

func (s *SavecmdContext) TAG(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserTAG, i)
}

func (s *SavecmdContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserCOMMA)
}

func (s *SavecmdContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserCOMMA, i)
}

func (s *SavecmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterSavecmd(s)
	}
}

func (s *SavecmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitSavecmd(s)
	}
}

func (p *GalleryParser) Command() (localctx ICommandContext) {
	localctx = NewCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, GalleryParserRULE_command)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(114)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case GalleryParserSAVE:
		localctx = NewSavecmdContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(94)
			p.Match(GalleryParserSAVE)
		}
		{
			p.SetState(95)
			p.Match(GalleryParserTAG)
		}
		p.SetState(100)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == GalleryParserCOMMA {
			{
				p.SetState(96)
				p.Match(GalleryParserCOMMA)
			}
			{
				p.SetState(97)
				p.Match(GalleryParserTAG)
			}

			p.SetState(102)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case GalleryParserSHOW:
		localctx = NewShowcmdContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(103)
			p.Match(GalleryParserSHOW)
		}
		{
			p.SetState(104)
			p.Match(GalleryParserTAG)
		}
		p.SetState(109)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == GalleryParserCOMMA {
			{
				p.SetState(105)
				p.Match(GalleryParserCOMMA)
			}
			{
				p.SetState(106)
				p.Match(GalleryParserTAG)
			}

			p.SetState(111)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case GalleryParserPROOF:
		localctx = NewProofcmdContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(112)
			p.Match(GalleryParserPROOF)
		}
		{
			p.SetState(113)
			p.Match(GalleryParserLABEL)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *ExpressionContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionContext) PATHCLIPOP() antlr.TerminalNode {
	return s.GetToken(GalleryParserPATHCLIPOP, 0)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *GalleryParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *GalleryParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 18
	p.EnterRecursionRule(localctx, 18, GalleryParserRULE_expression, _p)

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(117)
		p.tertiary(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(124)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewExpressionContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, GalleryParserRULE_expression)
			p.SetState(119)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(120)
				p.Match(GalleryParserPATHCLIPOP)
			}
			{
				p.SetState(121)
				p.tertiary(0)
			}

		}
		p.SetState(126)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
	}

	return localctx
}

// ITertiaryContext is an interface to support dynamic dispatch.
type ITertiaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTertiaryContext differentiates from other interfaces.
	IsTertiaryContext()
}

type TertiaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTertiaryContext() *TertiaryContext {
	var p = new(TertiaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_tertiary
	return p
}

func (*TertiaryContext) IsTertiaryContext() {}

func NewTertiaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TertiaryContext {
	var p = new(TertiaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_tertiary

	return p
}

func (s *TertiaryContext) GetParser() antlr.Parser { return s.parser }

func (s *TertiaryContext) CopyFrom(ctx *TertiaryContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *TertiaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TertiaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type PathtertiaryContext struct {
	*TertiaryContext
}

func NewPathtertiaryContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PathtertiaryContext {
	var p = new(PathtertiaryContext)

	p.TertiaryContext = NewEmptyTertiaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*TertiaryContext))

	return p
}

func (s *PathtertiaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathtertiaryContext) Path() IPathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathContext)
}

func (s *PathtertiaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterPathtertiary(s)
	}
}

func (s *PathtertiaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitPathtertiary(s)
	}
}

type LonesecondaryContext struct {
	*TertiaryContext
}

func NewLonesecondaryContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LonesecondaryContext {
	var p = new(LonesecondaryContext)

	p.TertiaryContext = NewEmptyTertiaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*TertiaryContext))

	return p
}

func (s *LonesecondaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LonesecondaryContext) Secondary() ISecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISecondaryContext)
}

func (s *LonesecondaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterLonesecondary(s)
	}
}

func (s *LonesecondaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitLonesecondary(s)
	}
}

type TermContext struct {
	*TertiaryContext
}

func NewTermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TermContext {
	var p = new(TermContext)

	p.TertiaryContext = NewEmptyTertiaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*TertiaryContext))

	return p
}

func (s *TermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TermContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *TermContext) Secondary() ISecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISecondaryContext)
}

func (s *TermContext) PLUS() antlr.TerminalNode {
	return s.GetToken(GalleryParserPLUS, 0)
}

func (s *TermContext) MINUS() antlr.TerminalNode {
	return s.GetToken(GalleryParserMINUS, 0)
}

func (s *TermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterTerm(s)
	}
}

func (s *TermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitTerm(s)
	}
}

func (p *GalleryParser) Tertiary() (localctx ITertiaryContext) {
	return p.tertiary(0)
}

func (p *GalleryParser) tertiary(_p int) (localctx ITertiaryContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewTertiaryContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ITertiaryContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 20
	p.EnterRecursionRule(localctx, 20, GalleryParserRULE_tertiary, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(130)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		localctx = NewLonesecondaryContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(128)
			p.secondary(0)
		}

	case 2:
		localctx = NewPathtertiaryContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(129)
			p.Path()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(137)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewTermContext(p, NewTertiaryContext(p, _parentctx, _parentState))
			p.PushNewRecursionContext(localctx, _startState, GalleryParserRULE_tertiary)
			p.SetState(132)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			p.SetState(133)
			_la = p.GetTokenStream().LA(1)

			if !(_la == GalleryParserPLUS || _la == GalleryParserMINUS) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
			{
				p.SetState(134)
				p.secondary(0)
			}

		}
		p.SetState(139)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext())
	}

	return localctx
}

// IPathContext is an interface to support dynamic dispatch.
type IPathContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPathContext differentiates from other interfaces.
	IsPathContext()
}

type PathContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathContext() *PathContext {
	var p = new(PathContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_path
	return p
}

func (*PathContext) IsPathContext() {}

func NewPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathContext {
	var p = new(PathContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_path

	return p
}

func (s *PathContext) GetParser() antlr.Parser { return s.parser }

func (s *PathContext) AllSecondary() []ISecondaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISecondaryContext)(nil)).Elem())
	var tst = make([]ISecondaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISecondaryContext)
		}
	}

	return tst
}

func (s *PathContext) Secondary(i int) ISecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISecondaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISecondaryContext)
}

func (s *PathContext) AllPATHJOIN() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserPATHJOIN)
}

func (s *PathContext) PATHJOIN(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserPATHJOIN, i)
}

func (s *PathContext) Cycle() ICycleContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICycleContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICycleContext)
}

func (s *PathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterPath(s)
	}
}

func (s *PathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitPath(s)
	}
}

func (p *GalleryParser) Path() (localctx IPathContext) {
	localctx = NewPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, GalleryParserRULE_path)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(140)
		p.secondary(0)
	}
	p.SetState(143)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(141)
				p.Match(GalleryParserPATHJOIN)
			}
			{
				p.SetState(142)
				p.secondary(0)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(145)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext())
	}
	p.SetState(148)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 13, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(147)
			p.Cycle()
		}

	}

	return localctx
}

// ICycleContext is an interface to support dynamic dispatch.
type ICycleContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCycleContext differentiates from other interfaces.
	IsCycleContext()
}

type CycleContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCycleContext() *CycleContext {
	var p = new(CycleContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_cycle
	return p
}

func (*CycleContext) IsCycleContext() {}

func NewCycleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CycleContext {
	var p = new(CycleContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_cycle

	return p
}

func (s *CycleContext) GetParser() antlr.Parser { return s.parser }

func (s *CycleContext) PATHJOIN() antlr.TerminalNode {
	return s.GetToken(GalleryParserPATHJOIN, 0)
}

func (s *CycleContext) CYCLE() antlr.TerminalNode {
	return s.GetToken(GalleryParserCYCLE, 0)
}

func (s *CycleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CycleContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CycleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterCycle(s)
	}
}

func (s *CycleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitCycle(s)
	}
}

func (p *GalleryParser) Cycle() (localctx ICycleContext) {
	localctx = NewCycleContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, GalleryParserRULE_cycle)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(150)
		p.Match(GalleryParserPATHJOIN)
	}
	{
		p.SetState(151)
		p.Match(GalleryParserCYCLE)
	}

	return localctx
}

// ISecondaryContext is an interface to support dynamic dispatch.
type ISecondaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSecondaryContext differentiates from other interfaces.
	IsSecondaryContext()
}

type SecondaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySecondaryContext() *SecondaryContext {
	var p = new(SecondaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_secondary
	return p
}

func (*SecondaryContext) IsSecondaryContext() {}

func NewSecondaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SecondaryContext {
	var p = new(SecondaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_secondary

	return p
}

func (s *SecondaryContext) GetParser() antlr.Parser { return s.parser }

func (s *SecondaryContext) CopyFrom(ctx *SecondaryContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *SecondaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SecondaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type TransformContext struct {
	*SecondaryContext
}

func NewTransformContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TransformContext {
	var p = new(TransformContext)

	p.SecondaryContext = NewEmptySecondaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SecondaryContext))

	return p
}

func (s *TransformContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TransformContext) Secondary() ISecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISecondaryContext)
}

func (s *TransformContext) Transformer() ITransformerContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITransformerContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITransformerContext)
}

func (s *TransformContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterTransform(s)
	}
}

func (s *TransformContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitTransform(s)
	}
}

type LoneprimaryContext struct {
	*SecondaryContext
}

func NewLoneprimaryContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LoneprimaryContext {
	var p = new(LoneprimaryContext)

	p.SecondaryContext = NewEmptySecondaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SecondaryContext))

	return p
}

func (s *LoneprimaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LoneprimaryContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *LoneprimaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterLoneprimary(s)
	}
}

func (s *LoneprimaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitLoneprimary(s)
	}
}

type FactorContext struct {
	*SecondaryContext
}

func NewFactorContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FactorContext {
	var p = new(FactorContext)

	p.SecondaryContext = NewEmptySecondaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SecondaryContext))

	return p
}

func (s *FactorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FactorContext) Secondary() ISecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISecondaryContext)
}

func (s *FactorContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *FactorContext) TIMES() antlr.TerminalNode {
	return s.GetToken(GalleryParserTIMES, 0)
}

func (s *FactorContext) OVER() antlr.TerminalNode {
	return s.GetToken(GalleryParserOVER, 0)
}

func (s *FactorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterFactor(s)
	}
}

func (s *FactorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitFactor(s)
	}
}

func (p *GalleryParser) Secondary() (localctx ISecondaryContext) {
	return p.secondary(0)
}

func (p *GalleryParser) secondary(_p int) (localctx ISecondaryContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewSecondaryContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ISecondaryContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 26
	p.EnterRecursionRule(localctx, 26, GalleryParserRULE_secondary, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	localctx = NewLoneprimaryContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(154)
		p.Primary()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(163)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(161)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext()) {
			case 1:
				localctx = NewFactorContext(p, NewSecondaryContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, GalleryParserRULE_secondary)
				p.SetState(156)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				p.SetState(157)
				_la = p.GetTokenStream().LA(1)

				if !(_la == GalleryParserTIMES || _la == GalleryParserOVER) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
				{
					p.SetState(158)
					p.Primary()
				}

			case 2:
				localctx = NewTransformContext(p, NewSecondaryContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, GalleryParserRULE_secondary)
				p.SetState(159)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				{
					p.SetState(160)
					p.Transformer()
				}

			}

		}
		p.SetState(165)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext())
	}

	return localctx
}

// ITransformerContext is an interface to support dynamic dispatch.
type ITransformerContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTransformerContext differentiates from other interfaces.
	IsTransformerContext()
}

type TransformerContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTransformerContext() *TransformerContext {
	var p = new(TransformerContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_transformer
	return p
}

func (*TransformerContext) IsTransformerContext() {}

func NewTransformerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TransformerContext {
	var p = new(TransformerContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_transformer

	return p
}

func (s *TransformerContext) GetParser() antlr.Parser { return s.parser }

func (s *TransformerContext) AllTRANSFORM() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserTRANSFORM)
}

func (s *TransformerContext) TRANSFORM(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserTRANSFORM, i)
}

func (s *TransformerContext) AllPrimary() []IPrimaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPrimaryContext)(nil)).Elem())
	var tst = make([]IPrimaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPrimaryContext)
		}
	}

	return tst
}

func (s *TransformerContext) Primary(i int) IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *TransformerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TransformerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TransformerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterTransformer(s)
	}
}

func (s *TransformerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitTransformer(s)
	}
}

func (p *GalleryParser) Transformer() (localctx ITransformerContext) {
	localctx = NewTransformerContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, GalleryParserRULE_transformer)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(168)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(166)
				p.Match(GalleryParserTRANSFORM)
			}
			{
				p.SetState(167)
				p.Primary()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(170)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 16, p.GetParserRuleContext())
	}

	return localctx
}

// IPrimaryContext is an interface to support dynamic dispatch.
type IPrimaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPrimaryContext differentiates from other interfaces.
	IsPrimaryContext()
}

type PrimaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryContext() *PrimaryContext {
	var p = new(PrimaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_primary
	return p
}

func (*PrimaryContext) IsPrimaryContext() {}

func NewPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryContext {
	var p = new(PrimaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_primary

	return p
}

func (s *PrimaryContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryContext) CopyFrom(ctx *PrimaryContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *PrimaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type PathpointContext struct {
	*PrimaryContext
}

func NewPathpointContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PathpointContext {
	var p = new(PathpointContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *PathpointContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathpointContext) POINT() antlr.TerminalNode {
	return s.GetToken(GalleryParserPOINT, 0)
}

func (s *PathpointContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *PathpointContext) OF() antlr.TerminalNode {
	return s.GetToken(GalleryParserOF, 0)
}

func (s *PathpointContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *PathpointContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterPathpoint(s)
	}
}

func (s *PathpointContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitPathpoint(s)
	}
}

type EdgeconstraintContext struct {
	*PrimaryContext
}

func NewEdgeconstraintContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EdgeconstraintContext {
	var p = new(EdgeconstraintContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *EdgeconstraintContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EdgeconstraintContext) EDGECONSTR() antlr.TerminalNode {
	return s.GetToken(GalleryParserEDGECONSTR, 0)
}

func (s *EdgeconstraintContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *EdgeconstraintContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterEdgeconstraint(s)
	}
}

func (s *EdgeconstraintContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitEdgeconstraint(s)
	}
}

type InterpolationContext struct {
	*PrimaryContext
}

func NewInterpolationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InterpolationContext {
	var p = new(InterpolationContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *InterpolationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InterpolationContext) Numtokenatom() INumtokenatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtokenatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtokenatomContext)
}

func (s *InterpolationContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(GalleryParserLBRACKET, 0)
}

func (s *InterpolationContext) AllTertiary() []ITertiaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITertiaryContext)(nil)).Elem())
	var tst = make([]ITertiaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITertiaryContext)
		}
	}

	return tst
}

func (s *InterpolationContext) Tertiary(i int) ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *InterpolationContext) COMMA() antlr.TerminalNode {
	return s.GetToken(GalleryParserCOMMA, 0)
}

func (s *InterpolationContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(GalleryParserRBRACKET, 0)
}

func (s *InterpolationContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *InterpolationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterInterpolation(s)
	}
}

func (s *InterpolationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitInterpolation(s)
	}
}

type ScalarnumatomContext struct {
	*PrimaryContext
}

func NewScalarnumatomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ScalarnumatomContext {
	var p = new(ScalarnumatomContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *ScalarnumatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScalarnumatomContext) Scalarmulop() IScalarmulopContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IScalarmulopContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IScalarmulopContext)
}

func (s *ScalarnumatomContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *ScalarnumatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterScalarnumatom(s)
	}
}

func (s *ScalarnumatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitScalarnumatom(s)
	}
}

type SimplenumatomContext struct {
	*PrimaryContext
}

func NewSimplenumatomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimplenumatomContext {
	var p = new(SimplenumatomContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *SimplenumatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimplenumatomContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *SimplenumatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterSimplenumatom(s)
	}
}

func (s *SimplenumatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitSimplenumatom(s)
	}
}

type SubpathContext struct {
	*PrimaryContext
}

func NewSubpathContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SubpathContext {
	var p = new(SubpathContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *SubpathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubpathContext) SUBPATH() antlr.TerminalNode {
	return s.GetToken(GalleryParserSUBPATH, 0)
}

func (s *SubpathContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *SubpathContext) OF() antlr.TerminalNode {
	return s.GetToken(GalleryParserOF, 0)
}

func (s *SubpathContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *SubpathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterSubpath(s)
	}
}

func (s *SubpathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitSubpath(s)
	}
}

type PairpartContext struct {
	*PrimaryContext
}

func NewPairpartContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PairpartContext {
	var p = new(PairpartContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *PairpartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairpartContext) PAIRPART() antlr.TerminalNode {
	return s.GetToken(GalleryParserPAIRPART, 0)
}

func (s *PairpartContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *PairpartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterPairpart(s)
	}
}

func (s *PairpartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitPairpart(s)
	}
}

type BoxContext struct {
	*PrimaryContext
}

func NewBoxContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BoxContext {
	var p = new(BoxContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *BoxContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoxContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *BoxContext) FRAME() antlr.TerminalNode {
	return s.GetToken(GalleryParserFRAME, 0)
}

func (s *BoxContext) BOX() antlr.TerminalNode {
	return s.GetToken(GalleryParserBOX, 0)
}

func (s *BoxContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterBox(s)
	}
}

func (s *BoxContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitBox(s)
	}
}

type FuncnumatomContext struct {
	*PrimaryContext
}

func NewFuncnumatomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FuncnumatomContext {
	var p = new(FuncnumatomContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *FuncnumatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncnumatomContext) MATHFUNC() antlr.TerminalNode {
	return s.GetToken(GalleryParserMATHFUNC, 0)
}

func (s *FuncnumatomContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *FuncnumatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterFuncnumatom(s)
	}
}

func (s *FuncnumatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitFuncnumatom(s)
	}
}

type ReversepathContext struct {
	*PrimaryContext
}

func NewReversepathContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ReversepathContext {
	var p = new(ReversepathContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *ReversepathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReversepathContext) REVERSE() antlr.TerminalNode {
	return s.GetToken(GalleryParserREVERSE, 0)
}

func (s *ReversepathContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *ReversepathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterReversepath(s)
	}
}

func (s *ReversepathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitReversepath(s)
	}
}

func (p *GalleryParser) Primary() (localctx IPrimaryContext) {
	localctx = NewPrimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, GalleryParserRULE_primary)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(210)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext()) {
	case 1:
		localctx = NewFuncnumatomContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(172)
			p.Match(GalleryParserMATHFUNC)
		}
		{
			p.SetState(173)
			p.Atom()
		}

	case 2:
		localctx = NewScalarnumatomContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(174)
			p.Scalarmulop()
		}
		{
			p.SetState(175)
			p.Atom()
		}

	case 3:
		localctx = NewInterpolationContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(177)
			p.Numtokenatom()
		}
		{
			p.SetState(178)
			p.Match(GalleryParserLBRACKET)
		}
		{
			p.SetState(179)
			p.tertiary(0)
		}
		{
			p.SetState(180)
			p.Match(GalleryParserCOMMA)
		}
		{
			p.SetState(181)
			p.tertiary(0)
		}
		{
			p.SetState(182)
			p.Match(GalleryParserRBRACKET)
		}

	case 4:
		localctx = NewInterpolationContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(184)
			p.Atom()
		}
		{
			p.SetState(185)
			p.Match(GalleryParserLBRACKET)
		}
		{
			p.SetState(186)
			p.tertiary(0)
		}
		{
			p.SetState(187)
			p.Match(GalleryParserCOMMA)
		}
		{
			p.SetState(188)
			p.tertiary(0)
		}
		{
			p.SetState(189)
			p.Match(GalleryParserRBRACKET)
		}

	case 5:
		localctx = NewSimplenumatomContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(191)
			p.Atom()
		}

	case 6:
		localctx = NewPairpartContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(192)
			p.Match(GalleryParserPAIRPART)
		}
		{
			p.SetState(193)
			p.Primary()
		}

	case 7:
		localctx = NewPathpointContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(194)
			p.Match(GalleryParserPOINT)
		}
		{
			p.SetState(195)
			p.tertiary(0)
		}
		{
			p.SetState(196)
			p.Match(GalleryParserOF)
		}
		{
			p.SetState(197)
			p.Primary()
		}

	case 8:
		localctx = NewReversepathContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(199)
			p.Match(GalleryParserREVERSE)
		}
		{
			p.SetState(200)
			p.Primary()
		}

	case 9:
		localctx = NewSubpathContext(p, localctx)
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(201)
			p.Match(GalleryParserSUBPATH)
		}
		{
			p.SetState(202)
			p.tertiary(0)
		}
		{
			p.SetState(203)
			p.Match(GalleryParserOF)
		}
		{
			p.SetState(204)
			p.Primary()
		}

	case 10:
		localctx = NewEdgeconstraintContext(p, localctx)
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(206)
			p.Match(GalleryParserEDGECONSTR)
		}
		{
			p.SetState(207)
			p.Primary()
		}

	case 11:
		localctx = NewBoxContext(p, localctx)
		p.EnterOuterAlt(localctx, 11)
		p.SetState(208)
		_la = p.GetTokenStream().LA(1)

		if !(_la == GalleryParserFRAME || _la == GalleryParserBOX) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
		{
			p.SetState(209)
			p.Variable()
		}

	}

	return localctx
}

// IScalarmulopContext is an interface to support dynamic dispatch.
type IScalarmulopContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsScalarmulopContext differentiates from other interfaces.
	IsScalarmulopContext()
}

type ScalarmulopContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyScalarmulopContext() *ScalarmulopContext {
	var p = new(ScalarmulopContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_scalarmulop
	return p
}

func (*ScalarmulopContext) IsScalarmulopContext() {}

func NewScalarmulopContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ScalarmulopContext {
	var p = new(ScalarmulopContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_scalarmulop

	return p
}

func (s *ScalarmulopContext) GetParser() antlr.Parser { return s.parser }

func (s *ScalarmulopContext) PLUS() antlr.TerminalNode {
	return s.GetToken(GalleryParserPLUS, 0)
}

func (s *ScalarmulopContext) MINUS() antlr.TerminalNode {
	return s.GetToken(GalleryParserMINUS, 0)
}

func (s *ScalarmulopContext) Numtokenatom() INumtokenatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtokenatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtokenatomContext)
}

func (s *ScalarmulopContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScalarmulopContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ScalarmulopContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterScalarmulop(s)
	}
}

func (s *ScalarmulopContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitScalarmulop(s)
	}
}

func (p *GalleryParser) Scalarmulop() (localctx IScalarmulopContext) {
	localctx = NewScalarmulopContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, GalleryParserRULE_scalarmulop)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(214)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case GalleryParserPLUS, GalleryParserMINUS:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(212)
		_la = p.GetTokenStream().LA(1)

		if !(_la == GalleryParserPLUS || _la == GalleryParserMINUS) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}

	case GalleryParserDECIMALTOKEN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(213)
			p.Numtokenatom()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// INumtokenatomContext is an interface to support dynamic dispatch.
type INumtokenatomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNumtokenatomContext differentiates from other interfaces.
	IsNumtokenatomContext()
}

type NumtokenatomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumtokenatomContext() *NumtokenatomContext {
	var p = new(NumtokenatomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_numtokenatom
	return p
}

func (*NumtokenatomContext) IsNumtokenatomContext() {}

func NewNumtokenatomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumtokenatomContext {
	var p = new(NumtokenatomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_numtokenatom

	return p
}

func (s *NumtokenatomContext) GetParser() antlr.Parser { return s.parser }

func (s *NumtokenatomContext) AllDECIMALTOKEN() []antlr.TerminalNode {
	return s.GetTokens(GalleryParserDECIMALTOKEN)
}

func (s *NumtokenatomContext) DECIMALTOKEN(i int) antlr.TerminalNode {
	return s.GetToken(GalleryParserDECIMALTOKEN, i)
}

func (s *NumtokenatomContext) OVER() antlr.TerminalNode {
	return s.GetToken(GalleryParserOVER, 0)
}

func (s *NumtokenatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumtokenatomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NumtokenatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterNumtokenatom(s)
	}
}

func (s *NumtokenatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitNumtokenatom(s)
	}
}

func (p *GalleryParser) Numtokenatom() (localctx INumtokenatomContext) {
	localctx = NewNumtokenatomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, GalleryParserRULE_numtokenatom)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(220)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 19, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(216)
			p.Match(GalleryParserDECIMALTOKEN)
		}
		{
			p.SetState(217)
			p.Match(GalleryParserOVER)
		}
		{
			p.SetState(218)
			p.Match(GalleryParserDECIMALTOKEN)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(219)
			p.Match(GalleryParserDECIMALTOKEN)
		}

	}

	return localctx
}

// IAtomContext is an interface to support dynamic dispatch.
type IAtomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAtomContext differentiates from other interfaces.
	IsAtomContext()
}

type AtomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAtomContext() *AtomContext {
	var p = new(AtomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_atom
	return p
}

func (*AtomContext) IsAtomContext() {}

func NewAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtomContext {
	var p = new(AtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_atom

	return p
}

func (s *AtomContext) GetParser() antlr.Parser { return s.parser }

func (s *AtomContext) CopyFrom(ctx *AtomContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *AtomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type VaratomContext struct {
	*AtomContext
}

func NewVaratomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *VaratomContext {
	var p = new(VaratomContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *VaratomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VaratomContext) Variable() IVariableContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IVariableContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IVariableContext)
}

func (s *VaratomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterVaratom(s)
	}
}

func (s *VaratomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitVaratom(s)
	}
}

type ExprgroupContext struct {
	*AtomContext
}

func NewExprgroupContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprgroupContext {
	var p = new(ExprgroupContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *ExprgroupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprgroupContext) BEGINGROUP() antlr.TerminalNode {
	return s.GetToken(GalleryParserBEGINGROUP, 0)
}

func (s *ExprgroupContext) Statementlist() IStatementlistContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementlistContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStatementlistContext)
}

func (s *ExprgroupContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *ExprgroupContext) ENDGROUP() antlr.TerminalNode {
	return s.GetToken(GalleryParserENDGROUP, 0)
}

func (s *ExprgroupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterExprgroup(s)
	}
}

func (s *ExprgroupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitExprgroup(s)
	}
}

type DecimalContext struct {
	*AtomContext
}

func NewDecimalContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DecimalContext {
	var p = new(DecimalContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *DecimalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DecimalContext) DECIMALTOKEN() antlr.TerminalNode {
	return s.GetToken(GalleryParserDECIMALTOKEN, 0)
}

func (s *DecimalContext) UNIT() antlr.TerminalNode {
	return s.GetToken(GalleryParserUNIT, 0)
}

func (s *DecimalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterDecimal(s)
	}
}

func (s *DecimalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitDecimal(s)
	}
}

type SubexpressionContext struct {
	*AtomContext
}

func NewSubexpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SubexpressionContext {
	var p = new(SubexpressionContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *SubexpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubexpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(GalleryParserLPAREN, 0)
}

func (s *SubexpressionContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *SubexpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(GalleryParserRPAREN, 0)
}

func (s *SubexpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterSubexpression(s)
	}
}

func (s *SubexpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitSubexpression(s)
	}
}

type LiteralpairContext struct {
	*AtomContext
}

func NewLiteralpairContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LiteralpairContext {
	var p = new(LiteralpairContext)

	p.AtomContext = NewEmptyAtomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AtomContext))

	return p
}

func (s *LiteralpairContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralpairContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(GalleryParserLPAREN, 0)
}

func (s *LiteralpairContext) AllTertiary() []ITertiaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITertiaryContext)(nil)).Elem())
	var tst = make([]ITertiaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITertiaryContext)
		}
	}

	return tst
}

func (s *LiteralpairContext) Tertiary(i int) ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *LiteralpairContext) COMMA() antlr.TerminalNode {
	return s.GetToken(GalleryParserCOMMA, 0)
}

func (s *LiteralpairContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(GalleryParserRPAREN, 0)
}

func (s *LiteralpairContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterLiteralpair(s)
	}
}

func (s *LiteralpairContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitLiteralpair(s)
	}
}

func (p *GalleryParser) Atom() (localctx IAtomContext) {
	localctx = NewAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, GalleryParserRULE_atom)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(242)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		localctx = NewDecimalContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(222)
			p.Match(GalleryParserDECIMALTOKEN)
		}
		p.SetState(224)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(223)
				p.Match(GalleryParserUNIT)
			}

		}

	case 2:
		localctx = NewVaratomContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(226)
			p.Variable()
		}

	case 3:
		localctx = NewLiteralpairContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(227)
			p.Match(GalleryParserLPAREN)
		}
		{
			p.SetState(228)
			p.tertiary(0)
		}
		{
			p.SetState(229)
			p.Match(GalleryParserCOMMA)
		}
		{
			p.SetState(230)
			p.tertiary(0)
		}
		{
			p.SetState(231)
			p.Match(GalleryParserRPAREN)
		}

	case 4:
		localctx = NewSubexpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(233)
			p.Match(GalleryParserLPAREN)
		}
		{
			p.SetState(234)
			p.tertiary(0)
		}
		{
			p.SetState(235)
			p.Match(GalleryParserRPAREN)
		}

	case 5:
		localctx = NewExprgroupContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(237)
			p.Match(GalleryParserBEGINGROUP)
		}
		{
			p.SetState(238)
			p.Statementlist()
		}
		{
			p.SetState(239)
			p.tertiary(0)
		}
		{
			p.SetState(240)
			p.Match(GalleryParserENDGROUP)
		}

	}

	return localctx
}

// IVariableContext is an interface to support dynamic dispatch.
type IVariableContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsVariableContext differentiates from other interfaces.
	IsVariableContext()
}

type VariableContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVariableContext() *VariableContext {
	var p = new(VariableContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_variable
	return p
}

func (*VariableContext) IsVariableContext() {}

func NewVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableContext {
	var p = new(VariableContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_variable

	return p
}

func (s *VariableContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableContext) MIXEDTAG() antlr.TerminalNode {
	return s.GetToken(GalleryParserMIXEDTAG, 0)
}

func (s *VariableContext) AllSubscript() []ISubscriptContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISubscriptContext)(nil)).Elem())
	var tst = make([]ISubscriptContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISubscriptContext)
		}
	}

	return tst
}

func (s *VariableContext) Subscript(i int) ISubscriptContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISubscriptContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISubscriptContext)
}

func (s *VariableContext) AllAnytag() []IAnytagContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAnytagContext)(nil)).Elem())
	var tst = make([]IAnytagContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAnytagContext)
		}
	}

	return tst
}

func (s *VariableContext) Anytag(i int) IAnytagContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnytagContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAnytagContext)
}

func (s *VariableContext) TAG() antlr.TerminalNode {
	return s.GetToken(GalleryParserTAG, 0)
}

func (s *VariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterVariable(s)
	}
}

func (s *VariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitVariable(s)
	}
}

func (p *GalleryParser) Variable() (localctx IVariableContext) {
	localctx = NewVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, GalleryParserRULE_variable)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.SetState(260)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case GalleryParserMIXEDTAG:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(244)
			p.Match(GalleryParserMIXEDTAG)
		}
		p.SetState(249)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 23, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(247)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case GalleryParserLBRACKET, GalleryParserDECIMALTOKEN:
					{
						p.SetState(245)
						p.Subscript()
					}

				case GalleryParserTAG, GalleryParserMIXEDTAG:
					{
						p.SetState(246)
						p.Anytag()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}

			}
			p.SetState(251)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 23, p.GetParserRuleContext())
		}

	case GalleryParserTAG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(252)
			p.Match(GalleryParserTAG)
		}
		p.SetState(257)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(255)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case GalleryParserLBRACKET, GalleryParserDECIMALTOKEN:
					{
						p.SetState(253)
						p.Subscript()
					}

				case GalleryParserTAG, GalleryParserMIXEDTAG:
					{
						p.SetState(254)
						p.Anytag()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}

			}
			p.SetState(259)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ISubscriptContext is an interface to support dynamic dispatch.
type ISubscriptContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSubscriptContext differentiates from other interfaces.
	IsSubscriptContext()
}

type SubscriptContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySubscriptContext() *SubscriptContext {
	var p = new(SubscriptContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_subscript
	return p
}

func (*SubscriptContext) IsSubscriptContext() {}

func NewSubscriptContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubscriptContext {
	var p = new(SubscriptContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_subscript

	return p
}

func (s *SubscriptContext) GetParser() antlr.Parser { return s.parser }

func (s *SubscriptContext) DECIMALTOKEN() antlr.TerminalNode {
	return s.GetToken(GalleryParserDECIMALTOKEN, 0)
}

func (s *SubscriptContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(GalleryParserLBRACKET, 0)
}

func (s *SubscriptContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *SubscriptContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(GalleryParserRBRACKET, 0)
}

func (s *SubscriptContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubscriptContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubscriptContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterSubscript(s)
	}
}

func (s *SubscriptContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitSubscript(s)
	}
}

func (p *GalleryParser) Subscript() (localctx ISubscriptContext) {
	localctx = NewSubscriptContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, GalleryParserRULE_subscript)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(267)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case GalleryParserDECIMALTOKEN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(262)
			p.Match(GalleryParserDECIMALTOKEN)
		}

	case GalleryParserLBRACKET:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(263)
			p.Match(GalleryParserLBRACKET)
		}
		{
			p.SetState(264)
			p.tertiary(0)
		}
		{
			p.SetState(265)
			p.Match(GalleryParserRBRACKET)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IAnytagContext is an interface to support dynamic dispatch.
type IAnytagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAnytagContext differentiates from other interfaces.
	IsAnytagContext()
}

type AnytagContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAnytagContext() *AnytagContext {
	var p = new(AnytagContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = GalleryParserRULE_anytag
	return p
}

func (*AnytagContext) IsAnytagContext() {}

func NewAnytagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnytagContext {
	var p = new(AnytagContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = GalleryParserRULE_anytag

	return p
}

func (s *AnytagContext) GetParser() antlr.Parser { return s.parser }

func (s *AnytagContext) TAG() antlr.TerminalNode {
	return s.GetToken(GalleryParserTAG, 0)
}

func (s *AnytagContext) MIXEDTAG() antlr.TerminalNode {
	return s.GetToken(GalleryParserMIXEDTAG, 0)
}

func (s *AnytagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnytagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AnytagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.EnterAnytag(s)
	}
}

func (s *AnytagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(GalleryListener); ok {
		listenerT.ExitAnytag(s)
	}
}

func (p *GalleryParser) Anytag() (localctx IAnytagContext) {
	localctx = NewAnytagContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, GalleryParserRULE_anytag)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	p.SetState(269)
	_la = p.GetTokenStream().LA(1)

	if !(_la == GalleryParserTAG || _la == GalleryParserMIXEDTAG) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}

	return localctx
}

func (p *GalleryParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 9:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	case 10:
		var t *TertiaryContext = nil
		if localctx != nil {
			t = localctx.(*TertiaryContext)
		}
		return p.Tertiary_Sempred(t, predIndex)

	case 13:
		var t *SecondaryContext = nil
		if localctx != nil {
			t = localctx.(*SecondaryContext)
		}
		return p.Secondary_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *GalleryParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *GalleryParser) Tertiary_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 1:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *GalleryParser) Secondary_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
