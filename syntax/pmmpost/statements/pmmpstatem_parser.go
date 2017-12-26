// Generated from PMMPStatem.g4 by ANTLR 4.7.

package statements // PMMPStatem
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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 57, 488,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9,
	39, 4, 40, 9, 40, 4, 41, 9, 41, 3, 2, 7, 2, 84, 10, 2, 12, 2, 14, 2, 87,
	11, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4,
	3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6,
	3, 6, 7, 6, 113, 10, 6, 12, 6, 14, 6, 116, 11, 6, 3, 7, 3, 7, 3, 7, 3,
	7, 3, 7, 5, 7, 123, 10, 7, 3, 8, 3, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 6,
	9, 132, 10, 9, 13, 9, 14, 9, 133, 3, 9, 3, 9, 3, 9, 6, 9, 139, 10, 9, 13,
	9, 14, 9, 140, 3, 9, 3, 9, 3, 9, 3, 9, 5, 9, 147, 10, 9, 3, 10, 3, 10,
	3, 10, 3, 10, 7, 10, 153, 10, 10, 12, 10, 14, 10, 156, 11, 10, 3, 11, 3,
	11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 13, 3, 13, 3, 13, 7, 13, 167, 10, 13,
	12, 13, 14, 13, 170, 11, 13, 3, 13, 3, 13, 3, 13, 7, 13, 175, 10, 13, 12,
	13, 14, 13, 178, 11, 13, 3, 13, 3, 13, 3, 13, 7, 13, 183, 10, 13, 12, 13,
	14, 13, 186, 11, 13, 3, 13, 3, 13, 3, 13, 7, 13, 191, 10, 13, 12, 13, 14,
	13, 194, 11, 13, 5, 13, 196, 10, 13, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14,
	5, 14, 203, 10, 14, 3, 15, 3, 15, 3, 15, 3, 15, 7, 15, 209, 10, 15, 12,
	15, 14, 15, 212, 11, 15, 3, 16, 3, 16, 3, 16, 3, 17, 3, 17, 3, 17, 3, 18,
	3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 3, 19, 5, 19, 227, 10, 19, 3, 19, 3,
	19, 5, 19, 231, 10, 19, 3, 20, 3, 20, 5, 20, 235, 10, 20, 3, 21, 3, 21,
	3, 21, 3, 21, 3, 21, 3, 21, 7, 21, 243, 10, 21, 12, 21, 14, 21, 246, 11,
	21, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 7, 22, 254, 10, 22, 12, 22,
	14, 22, 257, 11, 22, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3,
	23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23,
	3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 5, 23, 283, 10, 23, 3, 24, 3,
	24, 5, 24, 287, 10, 24, 3, 25, 3, 25, 3, 25, 3, 25, 5, 25, 293, 10, 25,
	3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 7, 26, 300, 10, 26, 12, 26, 14, 26,
	303, 11, 26, 3, 26, 3, 26, 3, 26, 7, 26, 308, 10, 26, 12, 26, 14, 26, 311,
	11, 26, 3, 26, 3, 26, 5, 26, 315, 10, 26, 3, 26, 3, 26, 3, 26, 3, 26, 3,
	26, 3, 26, 3, 26, 3, 26, 3, 26, 5, 26, 326, 10, 26, 3, 27, 3, 27, 3, 27,
	3, 27, 3, 27, 5, 27, 333, 10, 27, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3,
	28, 7, 28, 341, 10, 28, 12, 28, 14, 28, 344, 11, 28, 3, 29, 3, 29, 3, 29,
	3, 29, 3, 29, 3, 29, 5, 29, 352, 10, 29, 3, 29, 3, 29, 3, 29, 3, 29, 3,
	29, 7, 29, 359, 10, 29, 12, 29, 14, 29, 362, 11, 29, 3, 30, 3, 30, 3, 30,
	3, 30, 3, 30, 3, 30, 5, 30, 370, 10, 30, 3, 31, 3, 31, 3, 31, 3, 31, 3,
	31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31,
	3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 5, 31, 395, 10,
	31, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 7, 32,
	406, 10, 32, 12, 32, 14, 32, 409, 11, 32, 3, 32, 3, 32, 3, 32, 7, 32, 414,
	10, 32, 12, 32, 14, 32, 417, 11, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32,
	3, 32, 3, 32, 3, 32, 3, 32, 5, 32, 428, 10, 32, 3, 33, 3, 33, 3, 33, 7,
	33, 433, 10, 33, 12, 33, 14, 33, 436, 11, 33, 3, 34, 3, 34, 3, 34, 7, 34,
	441, 10, 34, 12, 34, 14, 34, 444, 11, 34, 3, 34, 5, 34, 447, 10, 34, 3,
	35, 3, 35, 5, 35, 451, 10, 35, 3, 36, 3, 36, 3, 36, 3, 37, 3, 37, 3, 37,
	3, 37, 3, 37, 7, 37, 461, 10, 37, 12, 37, 14, 37, 464, 11, 37, 3, 38, 3,
	38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 5, 38, 474, 10, 38, 3, 39,
	3, 39, 3, 39, 7, 39, 479, 10, 39, 12, 39, 14, 39, 482, 11, 39, 3, 40, 3,
	40, 3, 41, 3, 41, 3, 41, 2, 7, 40, 42, 54, 56, 72, 42, 2, 4, 6, 8, 10,
	12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46,
	48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 2,
	7, 4, 2, 27, 29, 45, 46, 3, 2, 21, 22, 3, 2, 23, 24, 3, 2, 49, 50, 4, 2,
	49, 50, 52, 53, 2, 526, 2, 85, 3, 2, 2, 2, 4, 90, 3, 2, 2, 2, 6, 94, 3,
	2, 2, 2, 8, 106, 3, 2, 2, 2, 10, 114, 3, 2, 2, 2, 12, 122, 3, 2, 2, 2,
	14, 124, 3, 2, 2, 2, 16, 146, 3, 2, 2, 2, 18, 148, 3, 2, 2, 2, 20, 157,
	3, 2, 2, 2, 22, 159, 3, 2, 2, 2, 24, 195, 3, 2, 2, 2, 26, 202, 3, 2, 2,
	2, 28, 204, 3, 2, 2, 2, 30, 213, 3, 2, 2, 2, 32, 216, 3, 2, 2, 2, 34, 219,
	3, 2, 2, 2, 36, 222, 3, 2, 2, 2, 38, 234, 3, 2, 2, 2, 40, 236, 3, 2, 2,
	2, 42, 247, 3, 2, 2, 2, 44, 282, 3, 2, 2, 2, 46, 286, 3, 2, 2, 2, 48, 292,
	3, 2, 2, 2, 50, 325, 3, 2, 2, 2, 52, 332, 3, 2, 2, 2, 54, 334, 3, 2, 2,
	2, 56, 351, 3, 2, 2, 2, 58, 369, 3, 2, 2, 2, 60, 394, 3, 2, 2, 2, 62, 427,
	3, 2, 2, 2, 64, 429, 3, 2, 2, 2, 66, 437, 3, 2, 2, 2, 68, 450, 3, 2, 2,
	2, 70, 452, 3, 2, 2, 2, 72, 455, 3, 2, 2, 2, 74, 473, 3, 2, 2, 2, 76, 475,
	3, 2, 2, 2, 78, 483, 3, 2, 2, 2, 80, 485, 3, 2, 2, 2, 82, 84, 5, 4, 3,
	2, 83, 82, 3, 2, 2, 2, 84, 87, 3, 2, 2, 2, 85, 83, 3, 2, 2, 2, 85, 86,
	3, 2, 2, 2, 86, 88, 3, 2, 2, 2, 87, 85, 3, 2, 2, 2, 88, 89, 7, 2, 2, 3,
	89, 3, 3, 2, 2, 2, 90, 91, 5, 6, 4, 2, 91, 92, 5, 10, 6, 2, 92, 93, 5,
	8, 5, 2, 93, 5, 3, 2, 2, 2, 94, 95, 7, 3, 2, 2, 95, 96, 7, 16, 2, 2, 96,
	97, 7, 56, 2, 2, 97, 98, 7, 15, 2, 2, 98, 99, 7, 54, 2, 2, 99, 100, 7,
	26, 2, 2, 100, 101, 7, 15, 2, 2, 101, 102, 7, 54, 2, 2, 102, 103, 7, 26,
	2, 2, 103, 104, 7, 17, 2, 2, 104, 105, 7, 14, 2, 2, 105, 7, 3, 2, 2, 2,
	106, 107, 7, 4, 2, 2, 107, 108, 7, 14, 2, 2, 108, 9, 3, 2, 2, 2, 109, 110,
	5, 12, 7, 2, 110, 111, 7, 14, 2, 2, 111, 113, 3, 2, 2, 2, 112, 109, 3,
	2, 2, 2, 113, 116, 3, 2, 2, 2, 114, 112, 3, 2, 2, 2, 114, 115, 3, 2, 2,
	2, 115, 11, 3, 2, 2, 2, 116, 114, 3, 2, 2, 2, 117, 123, 5, 14, 8, 2, 118,
	123, 5, 16, 9, 2, 119, 123, 5, 18, 10, 2, 120, 123, 5, 22, 12, 2, 121,
	123, 5, 26, 14, 2, 122, 117, 3, 2, 2, 2, 122, 118, 3, 2, 2, 2, 122, 119,
	3, 2, 2, 2, 122, 120, 3, 2, 2, 2, 122, 121, 3, 2, 2, 2, 123, 13, 3, 2,
	2, 2, 124, 125, 7, 47, 2, 2, 125, 126, 5, 10, 6, 2, 126, 127, 7, 48, 2,
	2, 127, 15, 3, 2, 2, 2, 128, 131, 5, 40, 21, 2, 129, 130, 7, 12, 2, 2,
	130, 132, 5, 40, 21, 2, 131, 129, 3, 2, 2, 2, 132, 133, 3, 2, 2, 2, 133,
	131, 3, 2, 2, 2, 133, 134, 3, 2, 2, 2, 134, 147, 3, 2, 2, 2, 135, 138,
	5, 54, 28, 2, 136, 137, 7, 12, 2, 2, 137, 139, 5, 54, 28, 2, 138, 136,
	3, 2, 2, 2, 139, 140, 3, 2, 2, 2, 140, 138, 3, 2, 2, 2, 140, 141, 3, 2,
	2, 2, 141, 147, 3, 2, 2, 2, 142, 143, 5, 76, 39, 2, 143, 144, 7, 12, 2,
	2, 144, 145, 5, 64, 33, 2, 145, 147, 3, 2, 2, 2, 146, 128, 3, 2, 2, 2,
	146, 135, 3, 2, 2, 2, 146, 142, 3, 2, 2, 2, 147, 17, 3, 2, 2, 2, 148, 149,
	5, 20, 11, 2, 149, 154, 5, 78, 40, 2, 150, 151, 7, 15, 2, 2, 151, 153,
	5, 78, 40, 2, 152, 150, 3, 2, 2, 2, 153, 156, 3, 2, 2, 2, 154, 152, 3,
	2, 2, 2, 154, 155, 3, 2, 2, 2, 155, 19, 3, 2, 2, 2, 156, 154, 3, 2, 2,
	2, 157, 158, 9, 2, 2, 2, 158, 21, 3, 2, 2, 2, 159, 160, 5, 24, 13, 2, 160,
	161, 7, 11, 2, 2, 161, 162, 5, 38, 20, 2, 162, 23, 3, 2, 2, 2, 163, 168,
	7, 53, 2, 2, 164, 167, 5, 52, 27, 2, 165, 167, 5, 80, 41, 2, 166, 164,
	3, 2, 2, 2, 166, 165, 3, 2, 2, 2, 167, 170, 3, 2, 2, 2, 168, 166, 3, 2,
	2, 2, 168, 169, 3, 2, 2, 2, 169, 196, 3, 2, 2, 2, 170, 168, 3, 2, 2, 2,
	171, 176, 7, 50, 2, 2, 172, 175, 5, 52, 27, 2, 173, 175, 5, 80, 41, 2,
	174, 172, 3, 2, 2, 2, 174, 173, 3, 2, 2, 2, 175, 178, 3, 2, 2, 2, 176,
	174, 3, 2, 2, 2, 176, 177, 3, 2, 2, 2, 177, 196, 3, 2, 2, 2, 178, 176,
	3, 2, 2, 2, 179, 184, 7, 52, 2, 2, 180, 183, 5, 52, 27, 2, 181, 183, 5,
	80, 41, 2, 182, 180, 3, 2, 2, 2, 182, 181, 3, 2, 2, 2, 183, 186, 3, 2,
	2, 2, 184, 182, 3, 2, 2, 2, 184, 185, 3, 2, 2, 2, 185, 196, 3, 2, 2, 2,
	186, 184, 3, 2, 2, 2, 187, 192, 7, 49, 2, 2, 188, 191, 5, 52, 27, 2, 189,
	191, 5, 80, 41, 2, 190, 188, 3, 2, 2, 2, 190, 189, 3, 2, 2, 2, 191, 194,
	3, 2, 2, 2, 192, 190, 3, 2, 2, 2, 192, 193, 3, 2, 2, 2, 193, 196, 3, 2,
	2, 2, 194, 192, 3, 2, 2, 2, 195, 163, 3, 2, 2, 2, 195, 171, 3, 2, 2, 2,
	195, 179, 3, 2, 2, 2, 195, 187, 3, 2, 2, 2, 196, 25, 3, 2, 2, 2, 197, 203,
	5, 28, 15, 2, 198, 203, 5, 30, 16, 2, 199, 203, 5, 32, 17, 2, 200, 203,
	5, 34, 18, 2, 201, 203, 5, 36, 19, 2, 202, 197, 3, 2, 2, 2, 202, 198, 3,
	2, 2, 2, 202, 199, 3, 2, 2, 2, 202, 200, 3, 2, 2, 2, 202, 201, 3, 2, 2,
	2, 203, 27, 3, 2, 2, 2, 204, 205, 7, 5, 2, 2, 205, 210, 5, 78, 40, 2, 206,
	207, 7, 15, 2, 2, 207, 209, 5, 78, 40, 2, 208, 206, 3, 2, 2, 2, 209, 212,
	3, 2, 2, 2, 210, 208, 3, 2, 2, 2, 210, 211, 3, 2, 2, 2, 211, 29, 3, 2,
	2, 2, 212, 210, 3, 2, 2, 2, 213, 214, 7, 6, 2, 2, 214, 215, 5, 78, 40,
	2, 215, 31, 3, 2, 2, 2, 216, 217, 7, 7, 2, 2, 217, 218, 5, 64, 33, 2, 218,
	33, 3, 2, 2, 2, 219, 220, 7, 8, 2, 2, 220, 221, 5, 64, 33, 2, 221, 35,
	3, 2, 2, 2, 222, 223, 7, 9, 2, 2, 223, 226, 7, 45, 2, 2, 224, 225, 7, 37,
	2, 2, 225, 227, 7, 54, 2, 2, 226, 224, 3, 2, 2, 2, 226, 227, 3, 2, 2, 2,
	227, 230, 3, 2, 2, 2, 228, 229, 7, 10, 2, 2, 229, 231, 7, 46, 2, 2, 230,
	228, 3, 2, 2, 2, 230, 231, 3, 2, 2, 2, 231, 37, 3, 2, 2, 2, 232, 235, 5,
	40, 21, 2, 233, 235, 5, 54, 28, 2, 234, 232, 3, 2, 2, 2, 234, 233, 3, 2,
	2, 2, 235, 39, 3, 2, 2, 2, 236, 237, 8, 21, 1, 2, 237, 238, 5, 42, 22,
	2, 238, 244, 3, 2, 2, 2, 239, 240, 12, 3, 2, 2, 240, 241, 9, 3, 2, 2, 241,
	243, 5, 42, 22, 2, 242, 239, 3, 2, 2, 2, 243, 246, 3, 2, 2, 2, 244, 242,
	3, 2, 2, 2, 244, 245, 3, 2, 2, 2, 245, 41, 3, 2, 2, 2, 246, 244, 3, 2,
	2, 2, 247, 248, 8, 22, 1, 2, 248, 249, 5, 44, 23, 2, 249, 255, 3, 2, 2,
	2, 250, 251, 12, 3, 2, 2, 251, 252, 9, 4, 2, 2, 252, 254, 5, 44, 23, 2,
	253, 250, 3, 2, 2, 2, 254, 257, 3, 2, 2, 2, 255, 253, 3, 2, 2, 2, 255,
	256, 3, 2, 2, 2, 256, 43, 3, 2, 2, 2, 257, 255, 3, 2, 2, 2, 258, 259, 7,
	33, 2, 2, 259, 283, 5, 50, 26, 2, 260, 261, 5, 46, 24, 2, 261, 262, 5,
	50, 26, 2, 262, 283, 3, 2, 2, 2, 263, 264, 5, 50, 26, 2, 264, 265, 7, 18,
	2, 2, 265, 266, 5, 40, 21, 2, 266, 267, 7, 15, 2, 2, 267, 268, 5, 40, 21,
	2, 268, 269, 7, 19, 2, 2, 269, 283, 3, 2, 2, 2, 270, 271, 5, 48, 25, 2,
	271, 272, 7, 18, 2, 2, 272, 273, 5, 40, 21, 2, 273, 274, 7, 15, 2, 2, 274,
	275, 5, 40, 21, 2, 275, 276, 7, 19, 2, 2, 276, 283, 3, 2, 2, 2, 277, 283,
	5, 50, 26, 2, 278, 279, 7, 32, 2, 2, 279, 283, 5, 60, 31, 2, 280, 281,
	7, 31, 2, 2, 281, 283, 5, 60, 31, 2, 282, 258, 3, 2, 2, 2, 282, 260, 3,
	2, 2, 2, 282, 263, 3, 2, 2, 2, 282, 270, 3, 2, 2, 2, 282, 277, 3, 2, 2,
	2, 282, 278, 3, 2, 2, 2, 282, 280, 3, 2, 2, 2, 283, 45, 3, 2, 2, 2, 284,
	287, 9, 3, 2, 2, 285, 287, 5, 48, 25, 2, 286, 284, 3, 2, 2, 2, 286, 285,
	3, 2, 2, 2, 287, 47, 3, 2, 2, 2, 288, 289, 7, 54, 2, 2, 289, 290, 7, 24,
	2, 2, 290, 293, 7, 54, 2, 2, 291, 293, 7, 54, 2, 2, 292, 288, 3, 2, 2,
	2, 292, 291, 3, 2, 2, 2, 293, 49, 3, 2, 2, 2, 294, 326, 7, 30, 2, 2, 295,
	326, 7, 25, 2, 2, 296, 301, 7, 53, 2, 2, 297, 300, 5, 52, 27, 2, 298, 300,
	5, 80, 41, 2, 299, 297, 3, 2, 2, 2, 299, 298, 3, 2, 2, 2, 300, 303, 3,
	2, 2, 2, 301, 299, 3, 2, 2, 2, 301, 302, 3, 2, 2, 2, 302, 326, 3, 2, 2,
	2, 303, 301, 3, 2, 2, 2, 304, 309, 7, 50, 2, 2, 305, 308, 5, 52, 27, 2,
	306, 308, 5, 80, 41, 2, 307, 305, 3, 2, 2, 2, 307, 306, 3, 2, 2, 2, 308,
	311, 3, 2, 2, 2, 309, 307, 3, 2, 2, 2, 309, 310, 3, 2, 2, 2, 310, 326,
	3, 2, 2, 2, 311, 309, 3, 2, 2, 2, 312, 314, 7, 54, 2, 2, 313, 315, 7, 26,
	2, 2, 314, 313, 3, 2, 2, 2, 314, 315, 3, 2, 2, 2, 315, 326, 3, 2, 2, 2,
	316, 317, 7, 16, 2, 2, 317, 318, 5, 40, 21, 2, 318, 319, 7, 17, 2, 2, 319,
	326, 3, 2, 2, 2, 320, 321, 7, 47, 2, 2, 321, 322, 5, 10, 6, 2, 322, 323,
	5, 40, 21, 2, 323, 324, 7, 48, 2, 2, 324, 326, 3, 2, 2, 2, 325, 294, 3,
	2, 2, 2, 325, 295, 3, 2, 2, 2, 325, 296, 3, 2, 2, 2, 325, 304, 3, 2, 2,
	2, 325, 312, 3, 2, 2, 2, 325, 316, 3, 2, 2, 2, 325, 320, 3, 2, 2, 2, 326,
	51, 3, 2, 2, 2, 327, 333, 7, 54, 2, 2, 328, 329, 7, 18, 2, 2, 329, 330,
	5, 40, 21, 2, 330, 331, 7, 19, 2, 2, 331, 333, 3, 2, 2, 2, 332, 327, 3,
	2, 2, 2, 332, 328, 3, 2, 2, 2, 333, 53, 3, 2, 2, 2, 334, 335, 8, 28, 1,
	2, 335, 336, 5, 56, 29, 2, 336, 342, 3, 2, 2, 2, 337, 338, 12, 3, 2, 2,
	338, 339, 9, 3, 2, 2, 339, 341, 5, 56, 29, 2, 340, 337, 3, 2, 2, 2, 341,
	344, 3, 2, 2, 2, 342, 340, 3, 2, 2, 2, 342, 343, 3, 2, 2, 2, 343, 55, 3,
	2, 2, 2, 344, 342, 3, 2, 2, 2, 345, 346, 8, 29, 1, 2, 346, 352, 5, 60,
	31, 2, 347, 348, 5, 42, 22, 2, 348, 349, 7, 23, 2, 2, 349, 350, 5, 60,
	31, 2, 350, 352, 3, 2, 2, 2, 351, 345, 3, 2, 2, 2, 351, 347, 3, 2, 2, 2,
	352, 360, 3, 2, 2, 2, 353, 354, 12, 5, 2, 2, 354, 355, 9, 4, 2, 2, 355,
	359, 5, 44, 23, 2, 356, 357, 12, 3, 2, 2, 357, 359, 5, 58, 30, 2, 358,
	353, 3, 2, 2, 2, 358, 356, 3, 2, 2, 2, 359, 362, 3, 2, 2, 2, 360, 358,
	3, 2, 2, 2, 360, 361, 3, 2, 2, 2, 361, 57, 3, 2, 2, 2, 362, 360, 3, 2,
	2, 2, 363, 364, 7, 37, 2, 2, 364, 370, 5, 44, 23, 2, 365, 366, 7, 39, 2,
	2, 366, 370, 5, 44, 23, 2, 367, 368, 7, 38, 2, 2, 368, 370, 5, 60, 31,
	2, 369, 363, 3, 2, 2, 2, 369, 365, 3, 2, 2, 2, 369, 367, 3, 2, 2, 2, 370,
	59, 3, 2, 2, 2, 371, 395, 5, 62, 32, 2, 372, 373, 5, 46, 24, 2, 373, 374,
	5, 62, 32, 2, 374, 395, 3, 2, 2, 2, 375, 376, 7, 35, 2, 2, 376, 377, 5,
	40, 21, 2, 377, 378, 7, 36, 2, 2, 378, 379, 5, 74, 38, 2, 379, 395, 3,
	2, 2, 2, 380, 381, 5, 50, 26, 2, 381, 382, 7, 18, 2, 2, 382, 383, 5, 54,
	28, 2, 383, 384, 7, 15, 2, 2, 384, 385, 5, 54, 28, 2, 385, 386, 7, 19,
	2, 2, 386, 395, 3, 2, 2, 2, 387, 388, 5, 48, 25, 2, 388, 389, 7, 18, 2,
	2, 389, 390, 5, 54, 28, 2, 390, 391, 7, 15, 2, 2, 391, 392, 5, 54, 28,
	2, 392, 393, 7, 19, 2, 2, 393, 395, 3, 2, 2, 2, 394, 371, 3, 2, 2, 2, 394,
	372, 3, 2, 2, 2, 394, 375, 3, 2, 2, 2, 394, 380, 3, 2, 2, 2, 394, 387,
	3, 2, 2, 2, 395, 61, 3, 2, 2, 2, 396, 397, 7, 16, 2, 2, 397, 398, 5, 40,
	21, 2, 398, 399, 7, 15, 2, 2, 399, 400, 5, 40, 21, 2, 400, 401, 7, 17,
	2, 2, 401, 428, 3, 2, 2, 2, 402, 407, 7, 52, 2, 2, 403, 406, 5, 52, 27,
	2, 404, 406, 5, 80, 41, 2, 405, 403, 3, 2, 2, 2, 405, 404, 3, 2, 2, 2,
	406, 409, 3, 2, 2, 2, 407, 405, 3, 2, 2, 2, 407, 408, 3, 2, 2, 2, 408,
	428, 3, 2, 2, 2, 409, 407, 3, 2, 2, 2, 410, 415, 7, 49, 2, 2, 411, 414,
	5, 52, 27, 2, 412, 414, 5, 80, 41, 2, 413, 411, 3, 2, 2, 2, 413, 412, 3,
	2, 2, 2, 414, 417, 3, 2, 2, 2, 415, 413, 3, 2, 2, 2, 415, 416, 3, 2, 2,
	2, 416, 428, 3, 2, 2, 2, 417, 415, 3, 2, 2, 2, 418, 419, 7, 16, 2, 2, 419,
	420, 5, 54, 28, 2, 420, 421, 7, 17, 2, 2, 421, 428, 3, 2, 2, 2, 422, 423,
	7, 47, 2, 2, 423, 424, 5, 10, 6, 2, 424, 425, 5, 54, 28, 2, 425, 426, 7,
	48, 2, 2, 426, 428, 3, 2, 2, 2, 427, 396, 3, 2, 2, 2, 427, 402, 3, 2, 2,
	2, 427, 410, 3, 2, 2, 2, 427, 418, 3, 2, 2, 2, 427, 422, 3, 2, 2, 2, 428,
	63, 3, 2, 2, 2, 429, 434, 5, 66, 34, 2, 430, 431, 7, 43, 2, 2, 431, 433,
	5, 66, 34, 2, 432, 430, 3, 2, 2, 2, 433, 436, 3, 2, 2, 2, 434, 432, 3,
	2, 2, 2, 434, 435, 3, 2, 2, 2, 435, 65, 3, 2, 2, 2, 436, 434, 3, 2, 2,
	2, 437, 442, 5, 68, 35, 2, 438, 439, 7, 20, 2, 2, 439, 441, 5, 68, 35,
	2, 440, 438, 3, 2, 2, 2, 441, 444, 3, 2, 2, 2, 442, 440, 3, 2, 2, 2, 442,
	443, 3, 2, 2, 2, 443, 446, 3, 2, 2, 2, 444, 442, 3, 2, 2, 2, 445, 447,
	5, 70, 36, 2, 446, 445, 3, 2, 2, 2, 446, 447, 3, 2, 2, 2, 447, 67, 3, 2,
	2, 2, 448, 451, 5, 72, 37, 2, 449, 451, 5, 54, 28, 2, 450, 448, 3, 2, 2,
	2, 450, 449, 3, 2, 2, 2, 451, 69, 3, 2, 2, 2, 452, 453, 7, 20, 2, 2, 453,
	454, 7, 42, 2, 2, 454, 71, 3, 2, 2, 2, 455, 456, 8, 37, 1, 2, 456, 457,
	5, 74, 38, 2, 457, 462, 3, 2, 2, 2, 458, 459, 12, 3, 2, 2, 459, 461, 5,
	58, 30, 2, 460, 458, 3, 2, 2, 2, 461, 464, 3, 2, 2, 2, 462, 460, 3, 2,
	2, 2, 462, 463, 3, 2, 2, 2, 463, 73, 3, 2, 2, 2, 464, 462, 3, 2, 2, 2,
	465, 474, 5, 76, 39, 2, 466, 467, 7, 41, 2, 2, 467, 474, 5, 74, 38, 2,
	468, 469, 7, 40, 2, 2, 469, 470, 5, 54, 28, 2, 470, 471, 7, 36, 2, 2, 471,
	472, 5, 74, 38, 2, 472, 474, 3, 2, 2, 2, 473, 465, 3, 2, 2, 2, 473, 466,
	3, 2, 2, 2, 473, 468, 3, 2, 2, 2, 474, 75, 3, 2, 2, 2, 475, 480, 7, 51,
	2, 2, 476, 479, 5, 52, 27, 2, 477, 479, 5, 80, 41, 2, 478, 476, 3, 2, 2,
	2, 478, 477, 3, 2, 2, 2, 479, 482, 3, 2, 2, 2, 480, 478, 3, 2, 2, 2, 480,
	481, 3, 2, 2, 2, 481, 77, 3, 2, 2, 2, 482, 480, 3, 2, 2, 2, 483, 484, 9,
	5, 2, 2, 484, 79, 3, 2, 2, 2, 485, 486, 9, 6, 2, 2, 486, 81, 3, 2, 2, 2,
	54, 85, 114, 122, 133, 140, 146, 154, 166, 168, 174, 176, 182, 184, 190,
	192, 195, 202, 210, 226, 230, 234, 244, 255, 282, 286, 292, 299, 301, 307,
	309, 314, 325, 332, 342, 351, 358, 360, 369, 394, 405, 407, 413, 415, 427,
	434, 442, 446, 450, 462, 473, 478, 480,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'beginfig'", "'endfig'", "'save'", "'showvariable'", "'draw'", "'fill'",
	"'pickup'", "'withcolor'", "':='", "'='", "':'", "';'", "','", "'('", "')'",
	"'['", "']'", "", "'+'", "'-'", "'*'", "'/'", "", "", "'numeric'", "'pair'",
	"'path'", "", "", "'length'", "", "'with'", "'point'", "'of'", "'scaled'",
	"'shifted'", "'rotated'", "'subpath'", "'reverse'", "'cycle'", "", "'intersectionpoint'",
	"", "", "'begingroup'", "'endgroup'", "", "", "", "", "", "", "'.'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "ASSIGN", "EQUALS", "COLON", "SEMIC",
	"COMMA", "LPAREN", "RPAREN", "LBRACKET", "RBRACKET", "PATHJOIN", "PLUS",
	"MINUS", "TIMES", "OVER", "WHATEVER", "UNIT", "NUMERIC", "PAIR", "PATH",
	"INTERNAL", "PAIRPART", "LENGTH", "MATHFUNC", "WITH", "POINT", "OF", "SCALED",
	"SHIFTED", "ROTATED", "SUBPATH", "REVERSE", "CYCLE", "PATHCLIPOP", "INTERSECTIONPOINT",
	"PEN", "COLOR", "BEGINGROUP", "ENDGROUP", "PTAG", "TAG", "PATHTAG", "MIXEDPTAG",
	"MIXEDTAG", "DECIMALTOKEN", "DOT", "LABEL", "WS",
}

var ruleNames = []string{
	"figures", "figure", "beginfig", "endfig", "statementlist", "statement",
	"compound", "equation", "declaration", "mptype", "assignment", "lvalue",
	"command", "saveStmt", "showvariableCmd", "drawCmd", "fillCmd", "pickupCmd",
	"expression", "numtertiary", "numsecondary", "numprimary", "scalarmulop",
	"numtokenatom", "numatom", "subscript", "pairtertiary", "pairsecondary",
	"transformer", "pairprimary", "pairatom", "pathexpression", "pathtertiary",
	"pathfragm", "cycle", "pathsecondary", "pathprimary", "pathatom", "tag",
	"anytag",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type PMMPStatemParser struct {
	*antlr.BaseParser
}

func NewPMMPStatemParser(input antlr.TokenStream) *PMMPStatemParser {
	this := new(PMMPStatemParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "PMMPStatem.g4"

	return this
}

// PMMPStatemParser tokens.
const (
	PMMPStatemParserEOF               = antlr.TokenEOF
	PMMPStatemParserT__0              = 1
	PMMPStatemParserT__1              = 2
	PMMPStatemParserT__2              = 3
	PMMPStatemParserT__3              = 4
	PMMPStatemParserT__4              = 5
	PMMPStatemParserT__5              = 6
	PMMPStatemParserT__6              = 7
	PMMPStatemParserT__7              = 8
	PMMPStatemParserASSIGN            = 9
	PMMPStatemParserEQUALS            = 10
	PMMPStatemParserCOLON             = 11
	PMMPStatemParserSEMIC             = 12
	PMMPStatemParserCOMMA             = 13
	PMMPStatemParserLPAREN            = 14
	PMMPStatemParserRPAREN            = 15
	PMMPStatemParserLBRACKET          = 16
	PMMPStatemParserRBRACKET          = 17
	PMMPStatemParserPATHJOIN          = 18
	PMMPStatemParserPLUS              = 19
	PMMPStatemParserMINUS             = 20
	PMMPStatemParserTIMES             = 21
	PMMPStatemParserOVER              = 22
	PMMPStatemParserWHATEVER          = 23
	PMMPStatemParserUNIT              = 24
	PMMPStatemParserNUMERIC           = 25
	PMMPStatemParserPAIR              = 26
	PMMPStatemParserPATH              = 27
	PMMPStatemParserINTERNAL          = 28
	PMMPStatemParserPAIRPART          = 29
	PMMPStatemParserLENGTH            = 30
	PMMPStatemParserMATHFUNC          = 31
	PMMPStatemParserWITH              = 32
	PMMPStatemParserPOINT             = 33
	PMMPStatemParserOF                = 34
	PMMPStatemParserSCALED            = 35
	PMMPStatemParserSHIFTED           = 36
	PMMPStatemParserROTATED           = 37
	PMMPStatemParserSUBPATH           = 38
	PMMPStatemParserREVERSE           = 39
	PMMPStatemParserCYCLE             = 40
	PMMPStatemParserPATHCLIPOP        = 41
	PMMPStatemParserINTERSECTIONPOINT = 42
	PMMPStatemParserPEN               = 43
	PMMPStatemParserCOLOR             = 44
	PMMPStatemParserBEGINGROUP        = 45
	PMMPStatemParserENDGROUP          = 46
	PMMPStatemParserPTAG              = 47
	PMMPStatemParserTAG               = 48
	PMMPStatemParserPATHTAG           = 49
	PMMPStatemParserMIXEDPTAG         = 50
	PMMPStatemParserMIXEDTAG          = 51
	PMMPStatemParserDECIMALTOKEN      = 52
	PMMPStatemParserDOT               = 53
	PMMPStatemParserLABEL             = 54
	PMMPStatemParserWS                = 55
)

// PMMPStatemParser rules.
const (
	PMMPStatemParserRULE_figures         = 0
	PMMPStatemParserRULE_figure          = 1
	PMMPStatemParserRULE_beginfig        = 2
	PMMPStatemParserRULE_endfig          = 3
	PMMPStatemParserRULE_statementlist   = 4
	PMMPStatemParserRULE_statement       = 5
	PMMPStatemParserRULE_compound        = 6
	PMMPStatemParserRULE_equation        = 7
	PMMPStatemParserRULE_declaration     = 8
	PMMPStatemParserRULE_mptype          = 9
	PMMPStatemParserRULE_assignment      = 10
	PMMPStatemParserRULE_lvalue          = 11
	PMMPStatemParserRULE_command         = 12
	PMMPStatemParserRULE_saveStmt        = 13
	PMMPStatemParserRULE_showvariableCmd = 14
	PMMPStatemParserRULE_drawCmd         = 15
	PMMPStatemParserRULE_fillCmd         = 16
	PMMPStatemParserRULE_pickupCmd       = 17
	PMMPStatemParserRULE_expression      = 18
	PMMPStatemParserRULE_numtertiary     = 19
	PMMPStatemParserRULE_numsecondary    = 20
	PMMPStatemParserRULE_numprimary      = 21
	PMMPStatemParserRULE_scalarmulop     = 22
	PMMPStatemParserRULE_numtokenatom    = 23
	PMMPStatemParserRULE_numatom         = 24
	PMMPStatemParserRULE_subscript       = 25
	PMMPStatemParserRULE_pairtertiary    = 26
	PMMPStatemParserRULE_pairsecondary   = 27
	PMMPStatemParserRULE_transformer     = 28
	PMMPStatemParserRULE_pairprimary     = 29
	PMMPStatemParserRULE_pairatom        = 30
	PMMPStatemParserRULE_pathexpression  = 31
	PMMPStatemParserRULE_pathtertiary    = 32
	PMMPStatemParserRULE_pathfragm       = 33
	PMMPStatemParserRULE_cycle           = 34
	PMMPStatemParserRULE_pathsecondary   = 35
	PMMPStatemParserRULE_pathprimary     = 36
	PMMPStatemParserRULE_pathatom        = 37
	PMMPStatemParserRULE_tag             = 38
	PMMPStatemParserRULE_anytag          = 39
)

// IFiguresContext is an interface to support dynamic dispatch.
type IFiguresContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFiguresContext differentiates from other interfaces.
	IsFiguresContext()
}

type FiguresContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFiguresContext() *FiguresContext {
	var p = new(FiguresContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_figures
	return p
}

func (*FiguresContext) IsFiguresContext() {}

func NewFiguresContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FiguresContext {
	var p = new(FiguresContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_figures

	return p
}

func (s *FiguresContext) GetParser() antlr.Parser { return s.parser }

func (s *FiguresContext) EOF() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserEOF, 0)
}

func (s *FiguresContext) AllFigure() []IFigureContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IFigureContext)(nil)).Elem())
	var tst = make([]IFigureContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IFigureContext)
		}
	}

	return tst
}

func (s *FiguresContext) Figure(i int) IFigureContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFigureContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IFigureContext)
}

func (s *FiguresContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FiguresContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FiguresContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterFigures(s)
	}
}

func (s *FiguresContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitFigures(s)
	}
}

func (p *PMMPStatemParser) Figures() (localctx IFiguresContext) {
	localctx = NewFiguresContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, PMMPStatemParserRULE_figures)
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
	p.SetState(83)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PMMPStatemParserT__0 {
		{
			p.SetState(80)
			p.Figure()
		}

		p.SetState(85)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(86)
		p.Match(PMMPStatemParserEOF)
	}

	return localctx
}

// IFigureContext is an interface to support dynamic dispatch.
type IFigureContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFigureContext differentiates from other interfaces.
	IsFigureContext()
}

type FigureContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFigureContext() *FigureContext {
	var p = new(FigureContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_figure
	return p
}

func (*FigureContext) IsFigureContext() {}

func NewFigureContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FigureContext {
	var p = new(FigureContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_figure

	return p
}

func (s *FigureContext) GetParser() antlr.Parser { return s.parser }

func (s *FigureContext) Beginfig() IBeginfigContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBeginfigContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBeginfigContext)
}

func (s *FigureContext) Statementlist() IStatementlistContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementlistContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStatementlistContext)
}

func (s *FigureContext) Endfig() IEndfigContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEndfigContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEndfigContext)
}

func (s *FigureContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FigureContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FigureContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterFigure(s)
	}
}

func (s *FigureContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitFigure(s)
	}
}

func (p *PMMPStatemParser) Figure() (localctx IFigureContext) {
	localctx = NewFigureContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, PMMPStatemParserRULE_figure)

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
		p.SetState(88)
		p.Beginfig()
	}
	{
		p.SetState(89)
		p.Statementlist()
	}
	{
		p.SetState(90)
		p.Endfig()
	}

	return localctx
}

// IBeginfigContext is an interface to support dynamic dispatch.
type IBeginfigContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBeginfigContext differentiates from other interfaces.
	IsBeginfigContext()
}

type BeginfigContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBeginfigContext() *BeginfigContext {
	var p = new(BeginfigContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_beginfig
	return p
}

func (*BeginfigContext) IsBeginfigContext() {}

func NewBeginfigContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BeginfigContext {
	var p = new(BeginfigContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_beginfig

	return p
}

func (s *BeginfigContext) GetParser() antlr.Parser { return s.parser }

func (s *BeginfigContext) LABEL() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserLABEL, 0)
}

func (s *BeginfigContext) AllDECIMALTOKEN() []antlr.TerminalNode {
	return s.GetTokens(PMMPStatemParserDECIMALTOKEN)
}

func (s *BeginfigContext) DECIMALTOKEN(i int) antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserDECIMALTOKEN, i)
}

func (s *BeginfigContext) AllUNIT() []antlr.TerminalNode {
	return s.GetTokens(PMMPStatemParserUNIT)
}

func (s *BeginfigContext) UNIT(i int) antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserUNIT, i)
}

func (s *BeginfigContext) SEMIC() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserSEMIC, 0)
}

func (s *BeginfigContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BeginfigContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BeginfigContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterBeginfig(s)
	}
}

func (s *BeginfigContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitBeginfig(s)
	}
}

func (p *PMMPStatemParser) Beginfig() (localctx IBeginfigContext) {
	localctx = NewBeginfigContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, PMMPStatemParserRULE_beginfig)

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
		p.SetState(92)
		p.Match(PMMPStatemParserT__0)
	}
	{
		p.SetState(93)
		p.Match(PMMPStatemParserLPAREN)
	}
	{
		p.SetState(94)
		p.Match(PMMPStatemParserLABEL)
	}
	{
		p.SetState(95)
		p.Match(PMMPStatemParserCOMMA)
	}
	{
		p.SetState(96)
		p.Match(PMMPStatemParserDECIMALTOKEN)
	}
	{
		p.SetState(97)
		p.Match(PMMPStatemParserUNIT)
	}
	{
		p.SetState(98)
		p.Match(PMMPStatemParserCOMMA)
	}
	{
		p.SetState(99)
		p.Match(PMMPStatemParserDECIMALTOKEN)
	}
	{
		p.SetState(100)
		p.Match(PMMPStatemParserUNIT)
	}
	{
		p.SetState(101)
		p.Match(PMMPStatemParserRPAREN)
	}
	{
		p.SetState(102)
		p.Match(PMMPStatemParserSEMIC)
	}

	return localctx
}

// IEndfigContext is an interface to support dynamic dispatch.
type IEndfigContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEndfigContext differentiates from other interfaces.
	IsEndfigContext()
}

type EndfigContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEndfigContext() *EndfigContext {
	var p = new(EndfigContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_endfig
	return p
}

func (*EndfigContext) IsEndfigContext() {}

func NewEndfigContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EndfigContext {
	var p = new(EndfigContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_endfig

	return p
}

func (s *EndfigContext) GetParser() antlr.Parser { return s.parser }

func (s *EndfigContext) SEMIC() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserSEMIC, 0)
}

func (s *EndfigContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EndfigContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EndfigContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterEndfig(s)
	}
}

func (s *EndfigContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitEndfig(s)
	}
}

func (p *PMMPStatemParser) Endfig() (localctx IEndfigContext) {
	localctx = NewEndfigContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, PMMPStatemParserRULE_endfig)

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
		p.SetState(104)
		p.Match(PMMPStatemParserT__1)
	}
	{
		p.SetState(105)
		p.Match(PMMPStatemParserSEMIC)
	}

	return localctx
}

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
	p.RuleIndex = PMMPStatemParserRULE_statementlist
	return p
}

func (*StatementlistContext) IsStatementlistContext() {}

func NewStatementlistContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementlistContext {
	var p = new(StatementlistContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_statementlist

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
	return s.GetTokens(PMMPStatemParserSEMIC)
}

func (s *StatementlistContext) SEMIC(i int) antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserSEMIC, i)
}

func (s *StatementlistContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementlistContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementlistContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterStatementlist(s)
	}
}

func (s *StatementlistContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitStatementlist(s)
	}
}

func (p *PMMPStatemParser) Statementlist() (localctx IStatementlistContext) {
	localctx = NewStatementlistContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, PMMPStatemParserRULE_statementlist)

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
	p.SetState(112)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(107)
				p.Statement()
			}
			{
				p.SetState(108)
				p.Match(PMMPStatemParserSEMIC)
			}

		}
		p.SetState(114)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext())
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
	p.RuleIndex = PMMPStatemParserRULE_statement
	return p
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) Compound() ICompoundContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICompoundContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICompoundContext)
}

func (s *StatementContext) Equation() IEquationContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEquationContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEquationContext)
}

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
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (p *PMMPStatemParser) Statement() (localctx IStatementContext) {
	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, PMMPStatemParserRULE_statement)

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

	p.SetState(120)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 2, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(115)
			p.Compound()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(116)
			p.Equation()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(117)
			p.Declaration()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(118)
			p.Assignment()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(119)
			p.Command()
		}

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
	p.RuleIndex = PMMPStatemParserRULE_compound
	return p
}

func (*CompoundContext) IsCompoundContext() {}

func NewCompoundContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompoundContext {
	var p = new(CompoundContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_compound

	return p
}

func (s *CompoundContext) GetParser() antlr.Parser { return s.parser }

func (s *CompoundContext) BEGINGROUP() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserBEGINGROUP, 0)
}

func (s *CompoundContext) Statementlist() IStatementlistContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementlistContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStatementlistContext)
}

func (s *CompoundContext) ENDGROUP() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserENDGROUP, 0)
}

func (s *CompoundContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompoundContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompoundContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterCompound(s)
	}
}

func (s *CompoundContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitCompound(s)
	}
}

func (p *PMMPStatemParser) Compound() (localctx ICompoundContext) {
	localctx = NewCompoundContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, PMMPStatemParserRULE_compound)

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
		p.SetState(122)
		p.Match(PMMPStatemParserBEGINGROUP)
	}
	{
		p.SetState(123)
		p.Statementlist()
	}
	{
		p.SetState(124)
		p.Match(PMMPStatemParserENDGROUP)
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
	p.RuleIndex = PMMPStatemParserRULE_equation
	return p
}

func (*EquationContext) IsEquationContext() {}

func NewEquationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EquationContext {
	var p = new(EquationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_equation

	return p
}

func (s *EquationContext) GetParser() antlr.Parser { return s.parser }

func (s *EquationContext) CopyFrom(ctx *EquationContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *EquationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EquationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type MultiequationContext struct {
	*EquationContext
}

func NewMultiequationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MultiequationContext {
	var p = new(MultiequationContext)

	p.EquationContext = NewEmptyEquationContext()
	p.parser = parser
	p.CopyFrom(ctx.(*EquationContext))

	return p
}

func (s *MultiequationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiequationContext) AllNumtertiary() []INumtertiaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem())
	var tst = make([]INumtertiaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(INumtertiaryContext)
		}
	}

	return tst
}

func (s *MultiequationContext) Numtertiary(i int) INumtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(INumtertiaryContext)
}

func (s *MultiequationContext) AllEQUALS() []antlr.TerminalNode {
	return s.GetTokens(PMMPStatemParserEQUALS)
}

func (s *MultiequationContext) EQUALS(i int) antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserEQUALS, i)
}

func (s *MultiequationContext) AllPairtertiary() []IPairtertiaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPairtertiaryContext)(nil)).Elem())
	var tst = make([]IPairtertiaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPairtertiaryContext)
		}
	}

	return tst
}

func (s *MultiequationContext) Pairtertiary(i int) IPairtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairtertiaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPairtertiaryContext)
}

func (s *MultiequationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterMultiequation(s)
	}
}

func (s *MultiequationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitMultiequation(s)
	}
}

type PathequationContext struct {
	*EquationContext
}

func NewPathequationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PathequationContext {
	var p = new(PathequationContext)

	p.EquationContext = NewEmptyEquationContext()
	p.parser = parser
	p.CopyFrom(ctx.(*EquationContext))

	return p
}

func (s *PathequationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathequationContext) Pathatom() IPathatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathatomContext)
}

func (s *PathequationContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserEQUALS, 0)
}

func (s *PathequationContext) Pathexpression() IPathexpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathexpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathexpressionContext)
}

func (s *PathequationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPathequation(s)
	}
}

func (s *PathequationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPathequation(s)
	}
}

func (p *PMMPStatemParser) Equation() (localctx IEquationContext) {
	localctx = NewEquationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, PMMPStatemParserRULE_equation)
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

	p.SetState(144)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 5, p.GetParserRuleContext()) {
	case 1:
		localctx = NewMultiequationContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(126)
			p.numtertiary(0)
		}
		p.SetState(129)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == PMMPStatemParserEQUALS {
			{
				p.SetState(127)
				p.Match(PMMPStatemParserEQUALS)
			}
			{
				p.SetState(128)
				p.numtertiary(0)
			}

			p.SetState(131)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case 2:
		localctx = NewMultiequationContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(133)
			p.pairtertiary(0)
		}
		p.SetState(136)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == PMMPStatemParserEQUALS {
			{
				p.SetState(134)
				p.Match(PMMPStatemParserEQUALS)
			}
			{
				p.SetState(135)
				p.pairtertiary(0)
			}

			p.SetState(138)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case 3:
		localctx = NewPathequationContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(140)
			p.Pathatom()
		}
		{
			p.SetState(141)
			p.Match(PMMPStatemParserEQUALS)
		}
		{
			p.SetState(142)
			p.Pathexpression()
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
	p.RuleIndex = PMMPStatemParserRULE_declaration
	return p
}

func (*DeclarationContext) IsDeclarationContext() {}

func NewDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclarationContext {
	var p = new(DeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_declaration

	return p
}

func (s *DeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *DeclarationContext) Mptype() IMptypeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IMptypeContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IMptypeContext)
}

func (s *DeclarationContext) AllTag() []ITagContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITagContext)(nil)).Elem())
	var tst = make([]ITagContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITagContext)
		}
	}

	return tst
}

func (s *DeclarationContext) Tag(i int) ITagContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITagContext)
}

func (s *DeclarationContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PMMPStatemParserCOMMA)
}

func (s *DeclarationContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserCOMMA, i)
}

func (s *DeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DeclarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterDeclaration(s)
	}
}

func (s *DeclarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitDeclaration(s)
	}
}

func (p *PMMPStatemParser) Declaration() (localctx IDeclarationContext) {
	localctx = NewDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, PMMPStatemParserRULE_declaration)
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
		p.SetState(146)
		p.Mptype()
	}
	{
		p.SetState(147)
		p.Tag()
	}
	p.SetState(152)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PMMPStatemParserCOMMA {
		{
			p.SetState(148)
			p.Match(PMMPStatemParserCOMMA)
		}
		{
			p.SetState(149)
			p.Tag()
		}

		p.SetState(154)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IMptypeContext is an interface to support dynamic dispatch.
type IMptypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsMptypeContext differentiates from other interfaces.
	IsMptypeContext()
}

type MptypeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMptypeContext() *MptypeContext {
	var p = new(MptypeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_mptype
	return p
}

func (*MptypeContext) IsMptypeContext() {}

func NewMptypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MptypeContext {
	var p = new(MptypeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_mptype

	return p
}

func (s *MptypeContext) GetParser() antlr.Parser { return s.parser }

func (s *MptypeContext) NUMERIC() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserNUMERIC, 0)
}

func (s *MptypeContext) PAIR() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPAIR, 0)
}

func (s *MptypeContext) PATH() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPATH, 0)
}

func (s *MptypeContext) PEN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPEN, 0)
}

func (s *MptypeContext) COLOR() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserCOLOR, 0)
}

func (s *MptypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MptypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MptypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterMptype(s)
	}
}

func (s *MptypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitMptype(s)
	}
}

func (p *PMMPStatemParser) Mptype() (localctx IMptypeContext) {
	localctx = NewMptypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, PMMPStatemParserRULE_mptype)
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
	p.SetState(155)
	_la = p.GetTokenStream().LA(1)

	if !(((_la-25)&-(0x1f+1)) == 0 && ((1<<uint((_la-25)))&((1<<(PMMPStatemParserNUMERIC-25))|(1<<(PMMPStatemParserPAIR-25))|(1<<(PMMPStatemParserPATH-25))|(1<<(PMMPStatemParserPEN-25))|(1<<(PMMPStatemParserCOLOR-25)))) != 0) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
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
	p.RuleIndex = PMMPStatemParserRULE_assignment
	return p
}

func (*AssignmentContext) IsAssignmentContext() {}

func NewAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentContext {
	var p = new(AssignmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_assignment

	return p
}

func (s *AssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *AssignmentContext) Lvalue() ILvalueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILvalueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILvalueContext)
}

func (s *AssignmentContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserASSIGN, 0)
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
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterAssignment(s)
	}
}

func (s *AssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitAssignment(s)
	}
}

func (p *PMMPStatemParser) Assignment() (localctx IAssignmentContext) {
	localctx = NewAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, PMMPStatemParserRULE_assignment)

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
		p.SetState(157)
		p.Lvalue()
	}
	{
		p.SetState(158)
		p.Match(PMMPStatemParserASSIGN)
	}
	{
		p.SetState(159)
		p.Expression()
	}

	return localctx
}

// ILvalueContext is an interface to support dynamic dispatch.
type ILvalueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLvalueContext differentiates from other interfaces.
	IsLvalueContext()
}

type LvalueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLvalueContext() *LvalueContext {
	var p = new(LvalueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_lvalue
	return p
}

func (*LvalueContext) IsLvalueContext() {}

func NewLvalueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LvalueContext {
	var p = new(LvalueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_lvalue

	return p
}

func (s *LvalueContext) GetParser() antlr.Parser { return s.parser }

func (s *LvalueContext) MIXEDTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserMIXEDTAG, 0)
}

func (s *LvalueContext) AllSubscript() []ISubscriptContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISubscriptContext)(nil)).Elem())
	var tst = make([]ISubscriptContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISubscriptContext)
		}
	}

	return tst
}

func (s *LvalueContext) Subscript(i int) ISubscriptContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISubscriptContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISubscriptContext)
}

func (s *LvalueContext) AllAnytag() []IAnytagContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAnytagContext)(nil)).Elem())
	var tst = make([]IAnytagContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAnytagContext)
		}
	}

	return tst
}

func (s *LvalueContext) Anytag(i int) IAnytagContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnytagContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAnytagContext)
}

func (s *LvalueContext) TAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserTAG, 0)
}

func (s *LvalueContext) MIXEDPTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserMIXEDPTAG, 0)
}

func (s *LvalueContext) PTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPTAG, 0)
}

func (s *LvalueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LvalueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LvalueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterLvalue(s)
	}
}

func (s *LvalueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitLvalue(s)
	}
}

func (p *PMMPStatemParser) Lvalue() (localctx ILvalueContext) {
	localctx = NewLvalueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, PMMPStatemParserRULE_lvalue)
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

	p.SetState(193)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPStatemParserMIXEDTAG:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(161)
			p.Match(PMMPStatemParserMIXEDTAG)
		}
		p.SetState(166)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == PMMPStatemParserLBRACKET || (((_la-47)&-(0x1f+1)) == 0 && ((1<<uint((_la-47)))&((1<<(PMMPStatemParserPTAG-47))|(1<<(PMMPStatemParserTAG-47))|(1<<(PMMPStatemParserMIXEDPTAG-47))|(1<<(PMMPStatemParserMIXEDTAG-47))|(1<<(PMMPStatemParserDECIMALTOKEN-47)))) != 0) {
			p.SetState(164)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case PMMPStatemParserLBRACKET, PMMPStatemParserDECIMALTOKEN:
				{
					p.SetState(162)
					p.Subscript()
				}

			case PMMPStatemParserPTAG, PMMPStatemParserTAG, PMMPStatemParserMIXEDPTAG, PMMPStatemParserMIXEDTAG:
				{
					p.SetState(163)
					p.Anytag()
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(168)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case PMMPStatemParserTAG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(169)
			p.Match(PMMPStatemParserTAG)
		}
		p.SetState(174)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == PMMPStatemParserLBRACKET || (((_la-47)&-(0x1f+1)) == 0 && ((1<<uint((_la-47)))&((1<<(PMMPStatemParserPTAG-47))|(1<<(PMMPStatemParserTAG-47))|(1<<(PMMPStatemParserMIXEDPTAG-47))|(1<<(PMMPStatemParserMIXEDTAG-47))|(1<<(PMMPStatemParserDECIMALTOKEN-47)))) != 0) {
			p.SetState(172)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case PMMPStatemParserLBRACKET, PMMPStatemParserDECIMALTOKEN:
				{
					p.SetState(170)
					p.Subscript()
				}

			case PMMPStatemParserPTAG, PMMPStatemParserTAG, PMMPStatemParserMIXEDPTAG, PMMPStatemParserMIXEDTAG:
				{
					p.SetState(171)
					p.Anytag()
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(176)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case PMMPStatemParserMIXEDPTAG:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(177)
			p.Match(PMMPStatemParserMIXEDPTAG)
		}
		p.SetState(182)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == PMMPStatemParserLBRACKET || (((_la-47)&-(0x1f+1)) == 0 && ((1<<uint((_la-47)))&((1<<(PMMPStatemParserPTAG-47))|(1<<(PMMPStatemParserTAG-47))|(1<<(PMMPStatemParserMIXEDPTAG-47))|(1<<(PMMPStatemParserMIXEDTAG-47))|(1<<(PMMPStatemParserDECIMALTOKEN-47)))) != 0) {
			p.SetState(180)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case PMMPStatemParserLBRACKET, PMMPStatemParserDECIMALTOKEN:
				{
					p.SetState(178)
					p.Subscript()
				}

			case PMMPStatemParserPTAG, PMMPStatemParserTAG, PMMPStatemParserMIXEDPTAG, PMMPStatemParserMIXEDTAG:
				{
					p.SetState(179)
					p.Anytag()
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(184)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case PMMPStatemParserPTAG:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(185)
			p.Match(PMMPStatemParserPTAG)
		}
		p.SetState(190)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == PMMPStatemParserLBRACKET || (((_la-47)&-(0x1f+1)) == 0 && ((1<<uint((_la-47)))&((1<<(PMMPStatemParserPTAG-47))|(1<<(PMMPStatemParserTAG-47))|(1<<(PMMPStatemParserMIXEDPTAG-47))|(1<<(PMMPStatemParserMIXEDTAG-47))|(1<<(PMMPStatemParserDECIMALTOKEN-47)))) != 0) {
			p.SetState(188)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case PMMPStatemParserLBRACKET, PMMPStatemParserDECIMALTOKEN:
				{
					p.SetState(186)
					p.Subscript()
				}

			case PMMPStatemParserPTAG, PMMPStatemParserTAG, PMMPStatemParserMIXEDPTAG, PMMPStatemParserMIXEDTAG:
				{
					p.SetState(187)
					p.Anytag()
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

			p.SetState(192)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.RuleIndex = PMMPStatemParserRULE_command
	return p
}

func (*CommandContext) IsCommandContext() {}

func NewCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandContext {
	var p = new(CommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_command

	return p
}

func (s *CommandContext) GetParser() antlr.Parser { return s.parser }

func (s *CommandContext) SaveStmt() ISaveStmtContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISaveStmtContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISaveStmtContext)
}

func (s *CommandContext) ShowvariableCmd() IShowvariableCmdContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IShowvariableCmdContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IShowvariableCmdContext)
}

func (s *CommandContext) DrawCmd() IDrawCmdContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDrawCmdContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDrawCmdContext)
}

func (s *CommandContext) FillCmd() IFillCmdContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFillCmdContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFillCmdContext)
}

func (s *CommandContext) PickupCmd() IPickupCmdContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPickupCmdContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPickupCmdContext)
}

func (s *CommandContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CommandContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CommandContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterCommand(s)
	}
}

func (s *CommandContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitCommand(s)
	}
}

func (p *PMMPStatemParser) Command() (localctx ICommandContext) {
	localctx = NewCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, PMMPStatemParserRULE_command)

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

	p.SetState(200)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPStatemParserT__2:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(195)
			p.SaveStmt()
		}

	case PMMPStatemParserT__3:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(196)
			p.ShowvariableCmd()
		}

	case PMMPStatemParserT__4:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(197)
			p.DrawCmd()
		}

	case PMMPStatemParserT__5:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(198)
			p.FillCmd()
		}

	case PMMPStatemParserT__6:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(199)
			p.PickupCmd()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ISaveStmtContext is an interface to support dynamic dispatch.
type ISaveStmtContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSaveStmtContext differentiates from other interfaces.
	IsSaveStmtContext()
}

type SaveStmtContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySaveStmtContext() *SaveStmtContext {
	var p = new(SaveStmtContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_saveStmt
	return p
}

func (*SaveStmtContext) IsSaveStmtContext() {}

func NewSaveStmtContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SaveStmtContext {
	var p = new(SaveStmtContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_saveStmt

	return p
}

func (s *SaveStmtContext) GetParser() antlr.Parser { return s.parser }

func (s *SaveStmtContext) AllTag() []ITagContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ITagContext)(nil)).Elem())
	var tst = make([]ITagContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ITagContext)
		}
	}

	return tst
}

func (s *SaveStmtContext) Tag(i int) ITagContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ITagContext)
}

func (s *SaveStmtContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PMMPStatemParserCOMMA)
}

func (s *SaveStmtContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserCOMMA, i)
}

func (s *SaveStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SaveStmtContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SaveStmtContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterSaveStmt(s)
	}
}

func (s *SaveStmtContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitSaveStmt(s)
	}
}

func (p *PMMPStatemParser) SaveStmt() (localctx ISaveStmtContext) {
	localctx = NewSaveStmtContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, PMMPStatemParserRULE_saveStmt)
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
		p.SetState(202)
		p.Match(PMMPStatemParserT__2)
	}
	{
		p.SetState(203)
		p.Tag()
	}
	p.SetState(208)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PMMPStatemParserCOMMA {
		{
			p.SetState(204)
			p.Match(PMMPStatemParserCOMMA)
		}
		{
			p.SetState(205)
			p.Tag()
		}

		p.SetState(210)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IShowvariableCmdContext is an interface to support dynamic dispatch.
type IShowvariableCmdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsShowvariableCmdContext differentiates from other interfaces.
	IsShowvariableCmdContext()
}

type ShowvariableCmdContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowvariableCmdContext() *ShowvariableCmdContext {
	var p = new(ShowvariableCmdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_showvariableCmd
	return p
}

func (*ShowvariableCmdContext) IsShowvariableCmdContext() {}

func NewShowvariableCmdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowvariableCmdContext {
	var p = new(ShowvariableCmdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_showvariableCmd

	return p
}

func (s *ShowvariableCmdContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowvariableCmdContext) Tag() ITagContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITagContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITagContext)
}

func (s *ShowvariableCmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowvariableCmdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowvariableCmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterShowvariableCmd(s)
	}
}

func (s *ShowvariableCmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitShowvariableCmd(s)
	}
}

func (p *PMMPStatemParser) ShowvariableCmd() (localctx IShowvariableCmdContext) {
	localctx = NewShowvariableCmdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, PMMPStatemParserRULE_showvariableCmd)

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
		p.SetState(211)
		p.Match(PMMPStatemParserT__3)
	}
	{
		p.SetState(212)
		p.Tag()
	}

	return localctx
}

// IDrawCmdContext is an interface to support dynamic dispatch.
type IDrawCmdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDrawCmdContext differentiates from other interfaces.
	IsDrawCmdContext()
}

type DrawCmdContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDrawCmdContext() *DrawCmdContext {
	var p = new(DrawCmdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_drawCmd
	return p
}

func (*DrawCmdContext) IsDrawCmdContext() {}

func NewDrawCmdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DrawCmdContext {
	var p = new(DrawCmdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_drawCmd

	return p
}

func (s *DrawCmdContext) GetParser() antlr.Parser { return s.parser }

func (s *DrawCmdContext) Pathexpression() IPathexpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathexpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathexpressionContext)
}

func (s *DrawCmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DrawCmdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DrawCmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterDrawCmd(s)
	}
}

func (s *DrawCmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitDrawCmd(s)
	}
}

func (p *PMMPStatemParser) DrawCmd() (localctx IDrawCmdContext) {
	localctx = NewDrawCmdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, PMMPStatemParserRULE_drawCmd)

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
		p.SetState(214)
		p.Match(PMMPStatemParserT__4)
	}
	{
		p.SetState(215)
		p.Pathexpression()
	}

	return localctx
}

// IFillCmdContext is an interface to support dynamic dispatch.
type IFillCmdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFillCmdContext differentiates from other interfaces.
	IsFillCmdContext()
}

type FillCmdContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFillCmdContext() *FillCmdContext {
	var p = new(FillCmdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_fillCmd
	return p
}

func (*FillCmdContext) IsFillCmdContext() {}

func NewFillCmdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FillCmdContext {
	var p = new(FillCmdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_fillCmd

	return p
}

func (s *FillCmdContext) GetParser() antlr.Parser { return s.parser }

func (s *FillCmdContext) Pathexpression() IPathexpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathexpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathexpressionContext)
}

func (s *FillCmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillCmdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FillCmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterFillCmd(s)
	}
}

func (s *FillCmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitFillCmd(s)
	}
}

func (p *PMMPStatemParser) FillCmd() (localctx IFillCmdContext) {
	localctx = NewFillCmdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, PMMPStatemParserRULE_fillCmd)

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
		p.SetState(217)
		p.Match(PMMPStatemParserT__5)
	}
	{
		p.SetState(218)
		p.Pathexpression()
	}

	return localctx
}

// IPickupCmdContext is an interface to support dynamic dispatch.
type IPickupCmdContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPickupCmdContext differentiates from other interfaces.
	IsPickupCmdContext()
}

type PickupCmdContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPickupCmdContext() *PickupCmdContext {
	var p = new(PickupCmdContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pickupCmd
	return p
}

func (*PickupCmdContext) IsPickupCmdContext() {}

func NewPickupCmdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PickupCmdContext {
	var p = new(PickupCmdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pickupCmd

	return p
}

func (s *PickupCmdContext) GetParser() antlr.Parser { return s.parser }

func (s *PickupCmdContext) PEN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPEN, 0)
}

func (s *PickupCmdContext) DECIMALTOKEN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserDECIMALTOKEN, 0)
}

func (s *PickupCmdContext) COLOR() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserCOLOR, 0)
}

func (s *PickupCmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PickupCmdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PickupCmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPickupCmd(s)
	}
}

func (s *PickupCmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPickupCmd(s)
	}
}

func (p *PMMPStatemParser) PickupCmd() (localctx IPickupCmdContext) {
	localctx = NewPickupCmdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, PMMPStatemParserRULE_pickupCmd)
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
		p.SetState(220)
		p.Match(PMMPStatemParserT__6)
	}
	{
		p.SetState(221)
		p.Match(PMMPStatemParserPEN)
	}
	p.SetState(224)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PMMPStatemParserSCALED {
		{
			p.SetState(222)
			p.Match(PMMPStatemParserSCALED)
		}
		{
			p.SetState(223)
			p.Match(PMMPStatemParserDECIMALTOKEN)
		}

	}
	p.SetState(228)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PMMPStatemParserT__7 {
		{
			p.SetState(226)
			p.Match(PMMPStatemParserT__7)
		}
		{
			p.SetState(227)
			p.Match(PMMPStatemParserCOLOR)
		}

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
	p.RuleIndex = PMMPStatemParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) Numtertiary() INumtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtertiaryContext)
}

func (s *ExpressionContext) Pairtertiary() IPairtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairtertiaryContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *PMMPStatemParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, PMMPStatemParserRULE_expression)

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

	p.SetState(232)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(230)
			p.numtertiary(0)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(231)
			p.pairtertiary(0)
		}

	}

	return localctx
}

// INumtertiaryContext is an interface to support dynamic dispatch.
type INumtertiaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNumtertiaryContext differentiates from other interfaces.
	IsNumtertiaryContext()
}

type NumtertiaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumtertiaryContext() *NumtertiaryContext {
	var p = new(NumtertiaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_numtertiary
	return p
}

func (*NumtertiaryContext) IsNumtertiaryContext() {}

func NewNumtertiaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumtertiaryContext {
	var p = new(NumtertiaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_numtertiary

	return p
}

func (s *NumtertiaryContext) GetParser() antlr.Parser { return s.parser }

func (s *NumtertiaryContext) Numsecondary() INumsecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumsecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumsecondaryContext)
}

func (s *NumtertiaryContext) Numtertiary() INumtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtertiaryContext)
}

func (s *NumtertiaryContext) PLUS() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPLUS, 0)
}

func (s *NumtertiaryContext) MINUS() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserMINUS, 0)
}

func (s *NumtertiaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumtertiaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NumtertiaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterNumtertiary(s)
	}
}

func (s *NumtertiaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitNumtertiary(s)
	}
}

func (p *PMMPStatemParser) Numtertiary() (localctx INumtertiaryContext) {
	return p.numtertiary(0)
}

func (p *PMMPStatemParser) numtertiary(_p int) (localctx INumtertiaryContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewNumtertiaryContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx INumtertiaryContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 38
	p.EnterRecursionRule(localctx, 38, PMMPStatemParserRULE_numtertiary, _p)
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
	{
		p.SetState(235)
		p.numsecondary(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(242)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewNumtertiaryContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, PMMPStatemParserRULE_numtertiary)
			p.SetState(237)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			p.SetState(238)
			_la = p.GetTokenStream().LA(1)

			if !(_la == PMMPStatemParserPLUS || _la == PMMPStatemParserMINUS) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
			{
				p.SetState(239)
				p.numsecondary(0)
			}

		}
		p.SetState(244)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext())
	}

	return localctx
}

// INumsecondaryContext is an interface to support dynamic dispatch.
type INumsecondaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNumsecondaryContext differentiates from other interfaces.
	IsNumsecondaryContext()
}

type NumsecondaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumsecondaryContext() *NumsecondaryContext {
	var p = new(NumsecondaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_numsecondary
	return p
}

func (*NumsecondaryContext) IsNumsecondaryContext() {}

func NewNumsecondaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumsecondaryContext {
	var p = new(NumsecondaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_numsecondary

	return p
}

func (s *NumsecondaryContext) GetParser() antlr.Parser { return s.parser }

func (s *NumsecondaryContext) Numprimary() INumprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumprimaryContext)
}

func (s *NumsecondaryContext) Numsecondary() INumsecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumsecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumsecondaryContext)
}

func (s *NumsecondaryContext) TIMES() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserTIMES, 0)
}

func (s *NumsecondaryContext) OVER() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserOVER, 0)
}

func (s *NumsecondaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumsecondaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NumsecondaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterNumsecondary(s)
	}
}

func (s *NumsecondaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitNumsecondary(s)
	}
}

func (p *PMMPStatemParser) Numsecondary() (localctx INumsecondaryContext) {
	return p.numsecondary(0)
}

func (p *PMMPStatemParser) numsecondary(_p int) (localctx INumsecondaryContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewNumsecondaryContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx INumsecondaryContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 40
	p.EnterRecursionRule(localctx, 40, PMMPStatemParserRULE_numsecondary, _p)
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
	{
		p.SetState(246)
		p.Numprimary()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(253)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewNumsecondaryContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, PMMPStatemParserRULE_numsecondary)
			p.SetState(248)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			p.SetState(249)
			_la = p.GetTokenStream().LA(1)

			if !(_la == PMMPStatemParserTIMES || _la == PMMPStatemParserOVER) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
			{
				p.SetState(250)
				p.Numprimary()
			}

		}
		p.SetState(255)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext())
	}

	return localctx
}

// INumprimaryContext is an interface to support dynamic dispatch.
type INumprimaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNumprimaryContext differentiates from other interfaces.
	IsNumprimaryContext()
}

type NumprimaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumprimaryContext() *NumprimaryContext {
	var p = new(NumprimaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_numprimary
	return p
}

func (*NumprimaryContext) IsNumprimaryContext() {}

func NewNumprimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumprimaryContext {
	var p = new(NumprimaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_numprimary

	return p
}

func (s *NumprimaryContext) GetParser() antlr.Parser { return s.parser }

func (s *NumprimaryContext) CopyFrom(ctx *NumprimaryContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *NumprimaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumprimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type InterpolationContext struct {
	*NumprimaryContext
}

func NewInterpolationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InterpolationContext {
	var p = new(InterpolationContext)

	p.NumprimaryContext = NewEmptyNumprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumprimaryContext))

	return p
}

func (s *InterpolationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InterpolationContext) Numatom() INumatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumatomContext)
}

func (s *InterpolationContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserLBRACKET, 0)
}

func (s *InterpolationContext) AllNumtertiary() []INumtertiaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem())
	var tst = make([]INumtertiaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(INumtertiaryContext)
		}
	}

	return tst
}

func (s *InterpolationContext) Numtertiary(i int) INumtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(INumtertiaryContext)
}

func (s *InterpolationContext) COMMA() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserCOMMA, 0)
}

func (s *InterpolationContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserRBRACKET, 0)
}

func (s *InterpolationContext) Numtokenatom() INumtokenatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtokenatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtokenatomContext)
}

func (s *InterpolationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterInterpolation(s)
	}
}

func (s *InterpolationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitInterpolation(s)
	}
}

type ScalarnumatomContext struct {
	*NumprimaryContext
}

func NewScalarnumatomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ScalarnumatomContext {
	var p = new(ScalarnumatomContext)

	p.NumprimaryContext = NewEmptyNumprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumprimaryContext))

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

func (s *ScalarnumatomContext) Numatom() INumatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumatomContext)
}

func (s *ScalarnumatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterScalarnumatom(s)
	}
}

func (s *ScalarnumatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitScalarnumatom(s)
	}
}

type SimplenumatomContext struct {
	*NumprimaryContext
}

func NewSimplenumatomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimplenumatomContext {
	var p = new(SimplenumatomContext)

	p.NumprimaryContext = NewEmptyNumprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumprimaryContext))

	return p
}

func (s *SimplenumatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimplenumatomContext) Numatom() INumatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumatomContext)
}

func (s *SimplenumatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterSimplenumatom(s)
	}
}

func (s *SimplenumatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitSimplenumatom(s)
	}
}

type PairpartContext struct {
	*NumprimaryContext
}

func NewPairpartContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PairpartContext {
	var p = new(PairpartContext)

	p.NumprimaryContext = NewEmptyNumprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumprimaryContext))

	return p
}

func (s *PairpartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairpartContext) PAIRPART() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPAIRPART, 0)
}

func (s *PairpartContext) Pairprimary() IPairprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairprimaryContext)
}

func (s *PairpartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPairpart(s)
	}
}

func (s *PairpartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPairpart(s)
	}
}

type PointdistanceContext struct {
	*NumprimaryContext
}

func NewPointdistanceContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PointdistanceContext {
	var p = new(PointdistanceContext)

	p.NumprimaryContext = NewEmptyNumprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumprimaryContext))

	return p
}

func (s *PointdistanceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PointdistanceContext) LENGTH() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserLENGTH, 0)
}

func (s *PointdistanceContext) Pairprimary() IPairprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairprimaryContext)
}

func (s *PointdistanceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPointdistance(s)
	}
}

func (s *PointdistanceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPointdistance(s)
	}
}

type FuncnumatomContext struct {
	*NumprimaryContext
}

func NewFuncnumatomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FuncnumatomContext {
	var p = new(FuncnumatomContext)

	p.NumprimaryContext = NewEmptyNumprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumprimaryContext))

	return p
}

func (s *FuncnumatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncnumatomContext) MATHFUNC() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserMATHFUNC, 0)
}

func (s *FuncnumatomContext) Numatom() INumatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumatomContext)
}

func (s *FuncnumatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterFuncnumatom(s)
	}
}

func (s *FuncnumatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitFuncnumatom(s)
	}
}

func (p *PMMPStatemParser) Numprimary() (localctx INumprimaryContext) {
	localctx = NewNumprimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, PMMPStatemParserRULE_numprimary)

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

	p.SetState(280)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 23, p.GetParserRuleContext()) {
	case 1:
		localctx = NewFuncnumatomContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(256)
			p.Match(PMMPStatemParserMATHFUNC)
		}
		{
			p.SetState(257)
			p.Numatom()
		}

	case 2:
		localctx = NewScalarnumatomContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(258)
			p.Scalarmulop()
		}
		{
			p.SetState(259)
			p.Numatom()
		}

	case 3:
		localctx = NewInterpolationContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(261)
			p.Numatom()
		}
		{
			p.SetState(262)
			p.Match(PMMPStatemParserLBRACKET)
		}
		{
			p.SetState(263)
			p.numtertiary(0)
		}
		{
			p.SetState(264)
			p.Match(PMMPStatemParserCOMMA)
		}
		{
			p.SetState(265)
			p.numtertiary(0)
		}
		{
			p.SetState(266)
			p.Match(PMMPStatemParserRBRACKET)
		}

	case 4:
		localctx = NewInterpolationContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(268)
			p.Numtokenatom()
		}
		{
			p.SetState(269)
			p.Match(PMMPStatemParserLBRACKET)
		}
		{
			p.SetState(270)
			p.numtertiary(0)
		}
		{
			p.SetState(271)
			p.Match(PMMPStatemParserCOMMA)
		}
		{
			p.SetState(272)
			p.numtertiary(0)
		}
		{
			p.SetState(273)
			p.Match(PMMPStatemParserRBRACKET)
		}

	case 5:
		localctx = NewSimplenumatomContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(275)
			p.Numatom()
		}

	case 6:
		localctx = NewPointdistanceContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(276)
			p.Match(PMMPStatemParserLENGTH)
		}
		{
			p.SetState(277)
			p.Pairprimary()
		}

	case 7:
		localctx = NewPairpartContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(278)
			p.Match(PMMPStatemParserPAIRPART)
		}
		{
			p.SetState(279)
			p.Pairprimary()
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
	p.RuleIndex = PMMPStatemParserRULE_scalarmulop
	return p
}

func (*ScalarmulopContext) IsScalarmulopContext() {}

func NewScalarmulopContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ScalarmulopContext {
	var p = new(ScalarmulopContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_scalarmulop

	return p
}

func (s *ScalarmulopContext) GetParser() antlr.Parser { return s.parser }

func (s *ScalarmulopContext) PLUS() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPLUS, 0)
}

func (s *ScalarmulopContext) MINUS() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserMINUS, 0)
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
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterScalarmulop(s)
	}
}

func (s *ScalarmulopContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitScalarmulop(s)
	}
}

func (p *PMMPStatemParser) Scalarmulop() (localctx IScalarmulopContext) {
	localctx = NewScalarmulopContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, PMMPStatemParserRULE_scalarmulop)
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

	p.SetState(284)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPStatemParserPLUS, PMMPStatemParserMINUS:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(282)
		_la = p.GetTokenStream().LA(1)

		if !(_la == PMMPStatemParserPLUS || _la == PMMPStatemParserMINUS) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}

	case PMMPStatemParserDECIMALTOKEN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(283)
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
	p.RuleIndex = PMMPStatemParserRULE_numtokenatom
	return p
}

func (*NumtokenatomContext) IsNumtokenatomContext() {}

func NewNumtokenatomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumtokenatomContext {
	var p = new(NumtokenatomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_numtokenatom

	return p
}

func (s *NumtokenatomContext) GetParser() antlr.Parser { return s.parser }

func (s *NumtokenatomContext) AllDECIMALTOKEN() []antlr.TerminalNode {
	return s.GetTokens(PMMPStatemParserDECIMALTOKEN)
}

func (s *NumtokenatomContext) DECIMALTOKEN(i int) antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserDECIMALTOKEN, i)
}

func (s *NumtokenatomContext) OVER() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserOVER, 0)
}

func (s *NumtokenatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumtokenatomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NumtokenatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterNumtokenatom(s)
	}
}

func (s *NumtokenatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitNumtokenatom(s)
	}
}

func (p *PMMPStatemParser) Numtokenatom() (localctx INumtokenatomContext) {
	localctx = NewNumtokenatomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, PMMPStatemParserRULE_numtokenatom)

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

	p.SetState(290)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(286)
			p.Match(PMMPStatemParserDECIMALTOKEN)
		}
		{
			p.SetState(287)
			p.Match(PMMPStatemParserOVER)
		}
		{
			p.SetState(288)
			p.Match(PMMPStatemParserDECIMALTOKEN)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(289)
			p.Match(PMMPStatemParserDECIMALTOKEN)
		}

	}

	return localctx
}

// INumatomContext is an interface to support dynamic dispatch.
type INumatomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNumatomContext differentiates from other interfaces.
	IsNumatomContext()
}

type NumatomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumatomContext() *NumatomContext {
	var p = new(NumatomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_numatom
	return p
}

func (*NumatomContext) IsNumatomContext() {}

func NewNumatomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumatomContext {
	var p = new(NumatomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_numatom

	return p
}

func (s *NumatomContext) GetParser() antlr.Parser { return s.parser }

func (s *NumatomContext) CopyFrom(ctx *NumatomContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *NumatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumatomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type InternalContext struct {
	*NumatomContext
}

func NewInternalContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InternalContext {
	var p = new(InternalContext)

	p.NumatomContext = NewEmptyNumatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumatomContext))

	return p
}

func (s *InternalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InternalContext) INTERNAL() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserINTERNAL, 0)
}

func (s *InternalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterInternal(s)
	}
}

func (s *InternalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitInternal(s)
	}
}

type VariableContext struct {
	*NumatomContext
}

func NewVariableContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *VariableContext {
	var p = new(VariableContext)

	p.NumatomContext = NewEmptyNumatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumatomContext))

	return p
}

func (s *VariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableContext) MIXEDTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserMIXEDTAG, 0)
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
	return s.GetToken(PMMPStatemParserTAG, 0)
}

func (s *VariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterVariable(s)
	}
}

func (s *VariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitVariable(s)
	}
}

type ExprgroupContext struct {
	*NumatomContext
}

func NewExprgroupContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExprgroupContext {
	var p = new(ExprgroupContext)

	p.NumatomContext = NewEmptyNumatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumatomContext))

	return p
}

func (s *ExprgroupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprgroupContext) BEGINGROUP() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserBEGINGROUP, 0)
}

func (s *ExprgroupContext) Statementlist() IStatementlistContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementlistContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStatementlistContext)
}

func (s *ExprgroupContext) Numtertiary() INumtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtertiaryContext)
}

func (s *ExprgroupContext) ENDGROUP() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserENDGROUP, 0)
}

func (s *ExprgroupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterExprgroup(s)
	}
}

func (s *ExprgroupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitExprgroup(s)
	}
}

type DecimalContext struct {
	*NumatomContext
}

func NewDecimalContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DecimalContext {
	var p = new(DecimalContext)

	p.NumatomContext = NewEmptyNumatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumatomContext))

	return p
}

func (s *DecimalContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DecimalContext) DECIMALTOKEN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserDECIMALTOKEN, 0)
}

func (s *DecimalContext) UNIT() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserUNIT, 0)
}

func (s *DecimalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterDecimal(s)
	}
}

func (s *DecimalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitDecimal(s)
	}
}

type SubexpressionContext struct {
	*NumatomContext
}

func NewSubexpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SubexpressionContext {
	var p = new(SubexpressionContext)

	p.NumatomContext = NewEmptyNumatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumatomContext))

	return p
}

func (s *SubexpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubexpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserLPAREN, 0)
}

func (s *SubexpressionContext) Numtertiary() INumtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtertiaryContext)
}

func (s *SubexpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserRPAREN, 0)
}

func (s *SubexpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterSubexpression(s)
	}
}

func (s *SubexpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitSubexpression(s)
	}
}

type WhateverContext struct {
	*NumatomContext
}

func NewWhateverContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *WhateverContext {
	var p = new(WhateverContext)

	p.NumatomContext = NewEmptyNumatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*NumatomContext))

	return p
}

func (s *WhateverContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhateverContext) WHATEVER() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserWHATEVER, 0)
}

func (s *WhateverContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterWhatever(s)
	}
}

func (s *WhateverContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitWhatever(s)
	}
}

func (p *PMMPStatemParser) Numatom() (localctx INumatomContext) {
	localctx = NewNumatomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, PMMPStatemParserRULE_numatom)

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

	p.SetState(323)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPStatemParserINTERNAL:
		localctx = NewInternalContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(292)
			p.Match(PMMPStatemParserINTERNAL)
		}

	case PMMPStatemParserWHATEVER:
		localctx = NewWhateverContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(293)
			p.Match(PMMPStatemParserWHATEVER)
		}

	case PMMPStatemParserMIXEDTAG:
		localctx = NewVariableContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(294)
			p.Match(PMMPStatemParserMIXEDTAG)
		}
		p.SetState(299)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(297)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case PMMPStatemParserLBRACKET, PMMPStatemParserDECIMALTOKEN:
					{
						p.SetState(295)
						p.Subscript()
					}

				case PMMPStatemParserPTAG, PMMPStatemParserTAG, PMMPStatemParserMIXEDPTAG, PMMPStatemParserMIXEDTAG:
					{
						p.SetState(296)
						p.Anytag()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}

			}
			p.SetState(301)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext())
		}

	case PMMPStatemParserTAG:
		localctx = NewVariableContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(302)
			p.Match(PMMPStatemParserTAG)
		}
		p.SetState(307)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(305)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case PMMPStatemParserLBRACKET, PMMPStatemParserDECIMALTOKEN:
					{
						p.SetState(303)
						p.Subscript()
					}

				case PMMPStatemParserPTAG, PMMPStatemParserTAG, PMMPStatemParserMIXEDPTAG, PMMPStatemParserMIXEDTAG:
					{
						p.SetState(304)
						p.Anytag()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}

			}
			p.SetState(309)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 29, p.GetParserRuleContext())
		}

	case PMMPStatemParserDECIMALTOKEN:
		localctx = NewDecimalContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(310)
			p.Match(PMMPStatemParserDECIMALTOKEN)
		}
		p.SetState(312)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(311)
				p.Match(PMMPStatemParserUNIT)
			}

		}

	case PMMPStatemParserLPAREN:
		localctx = NewSubexpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(314)
			p.Match(PMMPStatemParserLPAREN)
		}
		{
			p.SetState(315)
			p.numtertiary(0)
		}
		{
			p.SetState(316)
			p.Match(PMMPStatemParserRPAREN)
		}

	case PMMPStatemParserBEGINGROUP:
		localctx = NewExprgroupContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(318)
			p.Match(PMMPStatemParserBEGINGROUP)
		}
		{
			p.SetState(319)
			p.Statementlist()
		}
		{
			p.SetState(320)
			p.numtertiary(0)
		}
		{
			p.SetState(321)
			p.Match(PMMPStatemParserENDGROUP)
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
	p.RuleIndex = PMMPStatemParserRULE_subscript
	return p
}

func (*SubscriptContext) IsSubscriptContext() {}

func NewSubscriptContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubscriptContext {
	var p = new(SubscriptContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_subscript

	return p
}

func (s *SubscriptContext) GetParser() antlr.Parser { return s.parser }

func (s *SubscriptContext) DECIMALTOKEN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserDECIMALTOKEN, 0)
}

func (s *SubscriptContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserLBRACKET, 0)
}

func (s *SubscriptContext) Numtertiary() INumtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtertiaryContext)
}

func (s *SubscriptContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserRBRACKET, 0)
}

func (s *SubscriptContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubscriptContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubscriptContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterSubscript(s)
	}
}

func (s *SubscriptContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitSubscript(s)
	}
}

func (p *PMMPStatemParser) Subscript() (localctx ISubscriptContext) {
	localctx = NewSubscriptContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, PMMPStatemParserRULE_subscript)

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

	p.SetState(330)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPStatemParserDECIMALTOKEN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(325)
			p.Match(PMMPStatemParserDECIMALTOKEN)
		}

	case PMMPStatemParserLBRACKET:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(326)
			p.Match(PMMPStatemParserLBRACKET)
		}
		{
			p.SetState(327)
			p.numtertiary(0)
		}
		{
			p.SetState(328)
			p.Match(PMMPStatemParserRBRACKET)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IPairtertiaryContext is an interface to support dynamic dispatch.
type IPairtertiaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPairtertiaryContext differentiates from other interfaces.
	IsPairtertiaryContext()
}

type PairtertiaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPairtertiaryContext() *PairtertiaryContext {
	var p = new(PairtertiaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pairtertiary
	return p
}

func (*PairtertiaryContext) IsPairtertiaryContext() {}

func NewPairtertiaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PairtertiaryContext {
	var p = new(PairtertiaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pairtertiary

	return p
}

func (s *PairtertiaryContext) GetParser() antlr.Parser { return s.parser }

func (s *PairtertiaryContext) Pairsecondary() IPairsecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairsecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairsecondaryContext)
}

func (s *PairtertiaryContext) Pairtertiary() IPairtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairtertiaryContext)
}

func (s *PairtertiaryContext) PLUS() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPLUS, 0)
}

func (s *PairtertiaryContext) MINUS() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserMINUS, 0)
}

func (s *PairtertiaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairtertiaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PairtertiaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPairtertiary(s)
	}
}

func (s *PairtertiaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPairtertiary(s)
	}
}

func (p *PMMPStatemParser) Pairtertiary() (localctx IPairtertiaryContext) {
	return p.pairtertiary(0)
}

func (p *PMMPStatemParser) pairtertiary(_p int) (localctx IPairtertiaryContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewPairtertiaryContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IPairtertiaryContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 52
	p.EnterRecursionRule(localctx, 52, PMMPStatemParserRULE_pairtertiary, _p)
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
	{
		p.SetState(333)
		p.pairsecondary(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(340)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewPairtertiaryContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, PMMPStatemParserRULE_pairtertiary)
			p.SetState(335)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			p.SetState(336)
			_la = p.GetTokenStream().LA(1)

			if !(_la == PMMPStatemParserPLUS || _la == PMMPStatemParserMINUS) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
			{
				p.SetState(337)
				p.pairsecondary(0)
			}

		}
		p.SetState(342)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 33, p.GetParserRuleContext())
	}

	return localctx
}

// IPairsecondaryContext is an interface to support dynamic dispatch.
type IPairsecondaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPairsecondaryContext differentiates from other interfaces.
	IsPairsecondaryContext()
}

type PairsecondaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPairsecondaryContext() *PairsecondaryContext {
	var p = new(PairsecondaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pairsecondary
	return p
}

func (*PairsecondaryContext) IsPairsecondaryContext() {}

func NewPairsecondaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PairsecondaryContext {
	var p = new(PairsecondaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pairsecondary

	return p
}

func (s *PairsecondaryContext) GetParser() antlr.Parser { return s.parser }

func (s *PairsecondaryContext) CopyFrom(ctx *PairsecondaryContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *PairsecondaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairsecondaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type PairsecondContext struct {
	*PairsecondaryContext
}

func NewPairsecondContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PairsecondContext {
	var p = new(PairsecondContext)

	p.PairsecondaryContext = NewEmptyPairsecondaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PairsecondaryContext))

	return p
}

func (s *PairsecondContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairsecondContext) Pairprimary() IPairprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairprimaryContext)
}

func (s *PairsecondContext) Numsecondary() INumsecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumsecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumsecondaryContext)
}

func (s *PairsecondContext) TIMES() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserTIMES, 0)
}

func (s *PairsecondContext) Pairsecondary() IPairsecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairsecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairsecondaryContext)
}

func (s *PairsecondContext) Numprimary() INumprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumprimaryContext)
}

func (s *PairsecondContext) OVER() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserOVER, 0)
}

func (s *PairsecondContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPairsecond(s)
	}
}

func (s *PairsecondContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPairsecond(s)
	}
}

type TransformContext struct {
	*PairsecondaryContext
}

func NewTransformContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TransformContext {
	var p = new(TransformContext)

	p.PairsecondaryContext = NewEmptyPairsecondaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PairsecondaryContext))

	return p
}

func (s *TransformContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TransformContext) Pairsecondary() IPairsecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairsecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairsecondaryContext)
}

func (s *TransformContext) Transformer() ITransformerContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITransformerContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITransformerContext)
}

func (s *TransformContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterTransform(s)
	}
}

func (s *TransformContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitTransform(s)
	}
}

func (p *PMMPStatemParser) Pairsecondary() (localctx IPairsecondaryContext) {
	return p.pairsecondary(0)
}

func (p *PMMPStatemParser) pairsecondary(_p int) (localctx IPairsecondaryContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewPairsecondaryContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IPairsecondaryContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 54
	p.EnterRecursionRule(localctx, 54, PMMPStatemParserRULE_pairsecondary, _p)
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
	p.SetState(349)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext()) {
	case 1:
		localctx = NewPairsecondContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(344)
			p.Pairprimary()
		}

	case 2:
		localctx = NewPairsecondContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(345)
			p.numsecondary(0)
		}
		{
			p.SetState(346)
			p.Match(PMMPStatemParserTIMES)
		}
		{
			p.SetState(347)
			p.Pairprimary()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(358)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(356)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 35, p.GetParserRuleContext()) {
			case 1:
				localctx = NewPairsecondContext(p, NewPairsecondaryContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, PMMPStatemParserRULE_pairsecondary)
				p.SetState(351)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
				}
				p.SetState(352)
				_la = p.GetTokenStream().LA(1)

				if !(_la == PMMPStatemParserTIMES || _la == PMMPStatemParserOVER) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
				{
					p.SetState(353)
					p.Numprimary()
				}

			case 2:
				localctx = NewTransformContext(p, NewPairsecondaryContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, PMMPStatemParserRULE_pairsecondary)
				p.SetState(354)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				{
					p.SetState(355)
					p.Transformer()
				}

			}

		}
		p.SetState(360)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext())
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
	p.RuleIndex = PMMPStatemParserRULE_transformer
	return p
}

func (*TransformerContext) IsTransformerContext() {}

func NewTransformerContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TransformerContext {
	var p = new(TransformerContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_transformer

	return p
}

func (s *TransformerContext) GetParser() antlr.Parser { return s.parser }

func (s *TransformerContext) SCALED() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserSCALED, 0)
}

func (s *TransformerContext) Numprimary() INumprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumprimaryContext)
}

func (s *TransformerContext) ROTATED() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserROTATED, 0)
}

func (s *TransformerContext) SHIFTED() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserSHIFTED, 0)
}

func (s *TransformerContext) Pairprimary() IPairprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairprimaryContext)
}

func (s *TransformerContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TransformerContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TransformerContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterTransformer(s)
	}
}

func (s *TransformerContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitTransformer(s)
	}
}

func (p *PMMPStatemParser) Transformer() (localctx ITransformerContext) {
	localctx = NewTransformerContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, PMMPStatemParserRULE_transformer)

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

	p.SetState(367)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPStatemParserSCALED:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(361)
			p.Match(PMMPStatemParserSCALED)
		}
		{
			p.SetState(362)
			p.Numprimary()
		}

	case PMMPStatemParserROTATED:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(363)
			p.Match(PMMPStatemParserROTATED)
		}
		{
			p.SetState(364)
			p.Numprimary()
		}

	case PMMPStatemParserSHIFTED:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(365)
			p.Match(PMMPStatemParserSHIFTED)
		}
		{
			p.SetState(366)
			p.Pairprimary()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IPairprimaryContext is an interface to support dynamic dispatch.
type IPairprimaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPairprimaryContext differentiates from other interfaces.
	IsPairprimaryContext()
}

type PairprimaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPairprimaryContext() *PairprimaryContext {
	var p = new(PairprimaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pairprimary
	return p
}

func (*PairprimaryContext) IsPairprimaryContext() {}

func NewPairprimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PairprimaryContext {
	var p = new(PairprimaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pairprimary

	return p
}

func (s *PairprimaryContext) GetParser() antlr.Parser { return s.parser }

func (s *PairprimaryContext) CopyFrom(ctx *PairprimaryContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *PairprimaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairprimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SimplepairatomContext struct {
	*PairprimaryContext
}

func NewSimplepairatomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimplepairatomContext {
	var p = new(SimplepairatomContext)

	p.PairprimaryContext = NewEmptyPairprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PairprimaryContext))

	return p
}

func (s *SimplepairatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimplepairatomContext) Pairatom() IPairatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairatomContext)
}

func (s *SimplepairatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterSimplepairatom(s)
	}
}

func (s *SimplepairatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitSimplepairatom(s)
	}
}

type PathpointContext struct {
	*PairprimaryContext
}

func NewPathpointContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PathpointContext {
	var p = new(PathpointContext)

	p.PairprimaryContext = NewEmptyPairprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PairprimaryContext))

	return p
}

func (s *PathpointContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathpointContext) POINT() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPOINT, 0)
}

func (s *PathpointContext) Numtertiary() INumtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtertiaryContext)
}

func (s *PathpointContext) OF() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserOF, 0)
}

func (s *PathpointContext) Pathprimary() IPathprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathprimaryContext)
}

func (s *PathpointContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPathpoint(s)
	}
}

func (s *PathpointContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPathpoint(s)
	}
}

type PairinterpolationContext struct {
	*PairprimaryContext
}

func NewPairinterpolationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PairinterpolationContext {
	var p = new(PairinterpolationContext)

	p.PairprimaryContext = NewEmptyPairprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PairprimaryContext))

	return p
}

func (s *PairinterpolationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairinterpolationContext) Numatom() INumatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumatomContext)
}

func (s *PairinterpolationContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserLBRACKET, 0)
}

func (s *PairinterpolationContext) AllPairtertiary() []IPairtertiaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPairtertiaryContext)(nil)).Elem())
	var tst = make([]IPairtertiaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPairtertiaryContext)
		}
	}

	return tst
}

func (s *PairinterpolationContext) Pairtertiary(i int) IPairtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairtertiaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPairtertiaryContext)
}

func (s *PairinterpolationContext) COMMA() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserCOMMA, 0)
}

func (s *PairinterpolationContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserRBRACKET, 0)
}

func (s *PairinterpolationContext) Numtokenatom() INumtokenatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtokenatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtokenatomContext)
}

func (s *PairinterpolationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPairinterpolation(s)
	}
}

func (s *PairinterpolationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPairinterpolation(s)
	}
}

type ScalarmuloppairContext struct {
	*PairprimaryContext
}

func NewScalarmuloppairContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ScalarmuloppairContext {
	var p = new(ScalarmuloppairContext)

	p.PairprimaryContext = NewEmptyPairprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PairprimaryContext))

	return p
}

func (s *ScalarmuloppairContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScalarmuloppairContext) Scalarmulop() IScalarmulopContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IScalarmulopContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IScalarmulopContext)
}

func (s *ScalarmuloppairContext) Pairatom() IPairatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairatomContext)
}

func (s *ScalarmuloppairContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterScalarmuloppair(s)
	}
}

func (s *ScalarmuloppairContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitScalarmuloppair(s)
	}
}

func (p *PMMPStatemParser) Pairprimary() (localctx IPairprimaryContext) {
	localctx = NewPairprimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, PMMPStatemParserRULE_pairprimary)

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

	p.SetState(392)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 38, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSimplepairatomContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(369)
			p.Pairatom()
		}

	case 2:
		localctx = NewScalarmuloppairContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(370)
			p.Scalarmulop()
		}
		{
			p.SetState(371)
			p.Pairatom()
		}

	case 3:
		localctx = NewPathpointContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(373)
			p.Match(PMMPStatemParserPOINT)
		}
		{
			p.SetState(374)
			p.numtertiary(0)
		}
		{
			p.SetState(375)
			p.Match(PMMPStatemParserOF)
		}
		{
			p.SetState(376)
			p.Pathprimary()
		}

	case 4:
		localctx = NewPairinterpolationContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(378)
			p.Numatom()
		}
		{
			p.SetState(379)
			p.Match(PMMPStatemParserLBRACKET)
		}
		{
			p.SetState(380)
			p.pairtertiary(0)
		}
		{
			p.SetState(381)
			p.Match(PMMPStatemParserCOMMA)
		}
		{
			p.SetState(382)
			p.pairtertiary(0)
		}
		{
			p.SetState(383)
			p.Match(PMMPStatemParserRBRACKET)
		}

	case 5:
		localctx = NewPairinterpolationContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(385)
			p.Numtokenatom()
		}
		{
			p.SetState(386)
			p.Match(PMMPStatemParserLBRACKET)
		}
		{
			p.SetState(387)
			p.pairtertiary(0)
		}
		{
			p.SetState(388)
			p.Match(PMMPStatemParserCOMMA)
		}
		{
			p.SetState(389)
			p.pairtertiary(0)
		}
		{
			p.SetState(390)
			p.Match(PMMPStatemParserRBRACKET)
		}

	}

	return localctx
}

// IPairatomContext is an interface to support dynamic dispatch.
type IPairatomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPairatomContext differentiates from other interfaces.
	IsPairatomContext()
}

type PairatomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPairatomContext() *PairatomContext {
	var p = new(PairatomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pairatom
	return p
}

func (*PairatomContext) IsPairatomContext() {}

func NewPairatomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PairatomContext {
	var p = new(PairatomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pairatom

	return p
}

func (s *PairatomContext) GetParser() antlr.Parser { return s.parser }

func (s *PairatomContext) CopyFrom(ctx *PairatomContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *PairatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairatomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type PairexprgroupContext struct {
	*PairatomContext
}

func NewPairexprgroupContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PairexprgroupContext {
	var p = new(PairexprgroupContext)

	p.PairatomContext = NewEmptyPairatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PairatomContext))

	return p
}

func (s *PairexprgroupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairexprgroupContext) BEGINGROUP() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserBEGINGROUP, 0)
}

func (s *PairexprgroupContext) Statementlist() IStatementlistContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementlistContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStatementlistContext)
}

func (s *PairexprgroupContext) Pairtertiary() IPairtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairtertiaryContext)
}

func (s *PairexprgroupContext) ENDGROUP() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserENDGROUP, 0)
}

func (s *PairexprgroupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPairexprgroup(s)
	}
}

func (s *PairexprgroupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPairexprgroup(s)
	}
}

type SubpairexpressionContext struct {
	*PairatomContext
}

func NewSubpairexpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SubpairexpressionContext {
	var p = new(SubpairexpressionContext)

	p.PairatomContext = NewEmptyPairatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PairatomContext))

	return p
}

func (s *SubpairexpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubpairexpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserLPAREN, 0)
}

func (s *SubpairexpressionContext) Pairtertiary() IPairtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairtertiaryContext)
}

func (s *SubpairexpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserRPAREN, 0)
}

func (s *SubpairexpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterSubpairexpression(s)
	}
}

func (s *SubpairexpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitSubpairexpression(s)
	}
}

type PairvariableContext struct {
	*PairatomContext
}

func NewPairvariableContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PairvariableContext {
	var p = new(PairvariableContext)

	p.PairatomContext = NewEmptyPairatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PairatomContext))

	return p
}

func (s *PairvariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PairvariableContext) MIXEDPTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserMIXEDPTAG, 0)
}

func (s *PairvariableContext) AllSubscript() []ISubscriptContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISubscriptContext)(nil)).Elem())
	var tst = make([]ISubscriptContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISubscriptContext)
		}
	}

	return tst
}

func (s *PairvariableContext) Subscript(i int) ISubscriptContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISubscriptContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISubscriptContext)
}

func (s *PairvariableContext) AllAnytag() []IAnytagContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAnytagContext)(nil)).Elem())
	var tst = make([]IAnytagContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAnytagContext)
		}
	}

	return tst
}

func (s *PairvariableContext) Anytag(i int) IAnytagContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnytagContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAnytagContext)
}

func (s *PairvariableContext) PTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPTAG, 0)
}

func (s *PairvariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPairvariable(s)
	}
}

func (s *PairvariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPairvariable(s)
	}
}

type LiteralpairContext struct {
	*PairatomContext
}

func NewLiteralpairContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LiteralpairContext {
	var p = new(LiteralpairContext)

	p.PairatomContext = NewEmptyPairatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PairatomContext))

	return p
}

func (s *LiteralpairContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralpairContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserLPAREN, 0)
}

func (s *LiteralpairContext) AllNumtertiary() []INumtertiaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem())
	var tst = make([]INumtertiaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(INumtertiaryContext)
		}
	}

	return tst
}

func (s *LiteralpairContext) Numtertiary(i int) INumtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtertiaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(INumtertiaryContext)
}

func (s *LiteralpairContext) COMMA() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserCOMMA, 0)
}

func (s *LiteralpairContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserRPAREN, 0)
}

func (s *LiteralpairContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterLiteralpair(s)
	}
}

func (s *LiteralpairContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitLiteralpair(s)
	}
}

func (p *PMMPStatemParser) Pairatom() (localctx IPairatomContext) {
	localctx = NewPairatomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, PMMPStatemParserRULE_pairatom)

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

	p.SetState(425)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 43, p.GetParserRuleContext()) {
	case 1:
		localctx = NewLiteralpairContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(394)
			p.Match(PMMPStatemParserLPAREN)
		}
		{
			p.SetState(395)
			p.numtertiary(0)
		}
		{
			p.SetState(396)
			p.Match(PMMPStatemParserCOMMA)
		}
		{
			p.SetState(397)
			p.numtertiary(0)
		}
		{
			p.SetState(398)
			p.Match(PMMPStatemParserRPAREN)
		}

	case 2:
		localctx = NewPairvariableContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(400)
			p.Match(PMMPStatemParserMIXEDPTAG)
		}
		p.SetState(405)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(403)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case PMMPStatemParserLBRACKET, PMMPStatemParserDECIMALTOKEN:
					{
						p.SetState(401)
						p.Subscript()
					}

				case PMMPStatemParserPTAG, PMMPStatemParserTAG, PMMPStatemParserMIXEDPTAG, PMMPStatemParserMIXEDTAG:
					{
						p.SetState(402)
						p.Anytag()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}

			}
			p.SetState(407)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 40, p.GetParserRuleContext())
		}

	case 3:
		localctx = NewPairvariableContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(408)
			p.Match(PMMPStatemParserPTAG)
		}
		p.SetState(413)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(411)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case PMMPStatemParserLBRACKET, PMMPStatemParserDECIMALTOKEN:
					{
						p.SetState(409)
						p.Subscript()
					}

				case PMMPStatemParserPTAG, PMMPStatemParserTAG, PMMPStatemParserMIXEDPTAG, PMMPStatemParserMIXEDTAG:
					{
						p.SetState(410)
						p.Anytag()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}

			}
			p.SetState(415)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 42, p.GetParserRuleContext())
		}

	case 4:
		localctx = NewSubpairexpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(416)
			p.Match(PMMPStatemParserLPAREN)
		}
		{
			p.SetState(417)
			p.pairtertiary(0)
		}
		{
			p.SetState(418)
			p.Match(PMMPStatemParserRPAREN)
		}

	case 5:
		localctx = NewPairexprgroupContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(420)
			p.Match(PMMPStatemParserBEGINGROUP)
		}
		{
			p.SetState(421)
			p.Statementlist()
		}
		{
			p.SetState(422)
			p.pairtertiary(0)
		}
		{
			p.SetState(423)
			p.Match(PMMPStatemParserENDGROUP)
		}

	}

	return localctx
}

// IPathexpressionContext is an interface to support dynamic dispatch.
type IPathexpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPathexpressionContext differentiates from other interfaces.
	IsPathexpressionContext()
}

type PathexpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathexpressionContext() *PathexpressionContext {
	var p = new(PathexpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pathexpression
	return p
}

func (*PathexpressionContext) IsPathexpressionContext() {}

func NewPathexpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathexpressionContext {
	var p = new(PathexpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pathexpression

	return p
}

func (s *PathexpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *PathexpressionContext) AllPathtertiary() []IPathtertiaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPathtertiaryContext)(nil)).Elem())
	var tst = make([]IPathtertiaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPathtertiaryContext)
		}
	}

	return tst
}

func (s *PathexpressionContext) Pathtertiary(i int) IPathtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathtertiaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPathtertiaryContext)
}

func (s *PathexpressionContext) AllPATHCLIPOP() []antlr.TerminalNode {
	return s.GetTokens(PMMPStatemParserPATHCLIPOP)
}

func (s *PathexpressionContext) PATHCLIPOP(i int) antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPATHCLIPOP, i)
}

func (s *PathexpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathexpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PathexpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPathexpression(s)
	}
}

func (s *PathexpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPathexpression(s)
	}
}

func (p *PMMPStatemParser) Pathexpression() (localctx IPathexpressionContext) {
	localctx = NewPathexpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, PMMPStatemParserRULE_pathexpression)
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
		p.SetState(427)
		p.Pathtertiary()
	}
	p.SetState(432)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PMMPStatemParserPATHCLIPOP {
		{
			p.SetState(428)
			p.Match(PMMPStatemParserPATHCLIPOP)
		}
		{
			p.SetState(429)
			p.Pathtertiary()
		}

		p.SetState(434)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IPathtertiaryContext is an interface to support dynamic dispatch.
type IPathtertiaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPathtertiaryContext differentiates from other interfaces.
	IsPathtertiaryContext()
}

type PathtertiaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathtertiaryContext() *PathtertiaryContext {
	var p = new(PathtertiaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pathtertiary
	return p
}

func (*PathtertiaryContext) IsPathtertiaryContext() {}

func NewPathtertiaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathtertiaryContext {
	var p = new(PathtertiaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pathtertiary

	return p
}

func (s *PathtertiaryContext) GetParser() antlr.Parser { return s.parser }

func (s *PathtertiaryContext) AllPathfragm() []IPathfragmContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPathfragmContext)(nil)).Elem())
	var tst = make([]IPathfragmContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPathfragmContext)
		}
	}

	return tst
}

func (s *PathtertiaryContext) Pathfragm(i int) IPathfragmContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathfragmContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPathfragmContext)
}

func (s *PathtertiaryContext) AllPATHJOIN() []antlr.TerminalNode {
	return s.GetTokens(PMMPStatemParserPATHJOIN)
}

func (s *PathtertiaryContext) PATHJOIN(i int) antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPATHJOIN, i)
}

func (s *PathtertiaryContext) Cycle() ICycleContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICycleContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICycleContext)
}

func (s *PathtertiaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathtertiaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PathtertiaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPathtertiary(s)
	}
}

func (s *PathtertiaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPathtertiary(s)
	}
}

func (p *PMMPStatemParser) Pathtertiary() (localctx IPathtertiaryContext) {
	localctx = NewPathtertiaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, PMMPStatemParserRULE_pathtertiary)
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

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(435)
		p.Pathfragm()
	}
	p.SetState(440)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(436)
				p.Match(PMMPStatemParserPATHJOIN)
			}
			{
				p.SetState(437)
				p.Pathfragm()
			}

		}
		p.SetState(442)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 45, p.GetParserRuleContext())
	}
	p.SetState(444)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PMMPStatemParserPATHJOIN {
		{
			p.SetState(443)
			p.Cycle()
		}

	}

	return localctx
}

// IPathfragmContext is an interface to support dynamic dispatch.
type IPathfragmContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPathfragmContext differentiates from other interfaces.
	IsPathfragmContext()
}

type PathfragmContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathfragmContext() *PathfragmContext {
	var p = new(PathfragmContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pathfragm
	return p
}

func (*PathfragmContext) IsPathfragmContext() {}

func NewPathfragmContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathfragmContext {
	var p = new(PathfragmContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pathfragm

	return p
}

func (s *PathfragmContext) GetParser() antlr.Parser { return s.parser }

func (s *PathfragmContext) Pathsecondary() IPathsecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathsecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathsecondaryContext)
}

func (s *PathfragmContext) Pairtertiary() IPairtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairtertiaryContext)
}

func (s *PathfragmContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathfragmContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PathfragmContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPathfragm(s)
	}
}

func (s *PathfragmContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPathfragm(s)
	}
}

func (p *PMMPStatemParser) Pathfragm() (localctx IPathfragmContext) {
	localctx = NewPathfragmContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, PMMPStatemParserRULE_pathfragm)

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

	p.SetState(448)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPStatemParserSUBPATH, PMMPStatemParserREVERSE, PMMPStatemParserPATHTAG:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(446)
			p.pathsecondary(0)
		}

	case PMMPStatemParserLPAREN, PMMPStatemParserPLUS, PMMPStatemParserMINUS, PMMPStatemParserWHATEVER, PMMPStatemParserINTERNAL, PMMPStatemParserPAIRPART, PMMPStatemParserLENGTH, PMMPStatemParserMATHFUNC, PMMPStatemParserPOINT, PMMPStatemParserBEGINGROUP, PMMPStatemParserPTAG, PMMPStatemParserTAG, PMMPStatemParserMIXEDPTAG, PMMPStatemParserMIXEDTAG, PMMPStatemParserDECIMALTOKEN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(447)
			p.pairtertiary(0)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.RuleIndex = PMMPStatemParserRULE_cycle
	return p
}

func (*CycleContext) IsCycleContext() {}

func NewCycleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CycleContext {
	var p = new(CycleContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_cycle

	return p
}

func (s *CycleContext) GetParser() antlr.Parser { return s.parser }

func (s *CycleContext) PATHJOIN() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPATHJOIN, 0)
}

func (s *CycleContext) CYCLE() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserCYCLE, 0)
}

func (s *CycleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CycleContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CycleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterCycle(s)
	}
}

func (s *CycleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitCycle(s)
	}
}

func (p *PMMPStatemParser) Cycle() (localctx ICycleContext) {
	localctx = NewCycleContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, PMMPStatemParserRULE_cycle)

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
		p.SetState(450)
		p.Match(PMMPStatemParserPATHJOIN)
	}
	{
		p.SetState(451)
		p.Match(PMMPStatemParserCYCLE)
	}

	return localctx
}

// IPathsecondaryContext is an interface to support dynamic dispatch.
type IPathsecondaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPathsecondaryContext differentiates from other interfaces.
	IsPathsecondaryContext()
}

type PathsecondaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathsecondaryContext() *PathsecondaryContext {
	var p = new(PathsecondaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pathsecondary
	return p
}

func (*PathsecondaryContext) IsPathsecondaryContext() {}

func NewPathsecondaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathsecondaryContext {
	var p = new(PathsecondaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pathsecondary

	return p
}

func (s *PathsecondaryContext) GetParser() antlr.Parser { return s.parser }

func (s *PathsecondaryContext) Pathprimary() IPathprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathprimaryContext)
}

func (s *PathsecondaryContext) Pathsecondary() IPathsecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathsecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathsecondaryContext)
}

func (s *PathsecondaryContext) Transformer() ITransformerContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITransformerContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITransformerContext)
}

func (s *PathsecondaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathsecondaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PathsecondaryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPathsecondary(s)
	}
}

func (s *PathsecondaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPathsecondary(s)
	}
}

func (p *PMMPStatemParser) Pathsecondary() (localctx IPathsecondaryContext) {
	return p.pathsecondary(0)
}

func (p *PMMPStatemParser) pathsecondary(_p int) (localctx IPathsecondaryContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewPathsecondaryContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IPathsecondaryContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 70
	p.EnterRecursionRule(localctx, 70, PMMPStatemParserRULE_pathsecondary, _p)

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
		p.SetState(454)
		p.Pathprimary()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(460)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 48, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewPathsecondaryContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, PMMPStatemParserRULE_pathsecondary)
			p.SetState(456)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(457)
				p.Transformer()
			}

		}
		p.SetState(462)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 48, p.GetParserRuleContext())
	}

	return localctx
}

// IPathprimaryContext is an interface to support dynamic dispatch.
type IPathprimaryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPathprimaryContext differentiates from other interfaces.
	IsPathprimaryContext()
}

type PathprimaryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathprimaryContext() *PathprimaryContext {
	var p = new(PathprimaryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pathprimary
	return p
}

func (*PathprimaryContext) IsPathprimaryContext() {}

func NewPathprimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathprimaryContext {
	var p = new(PathprimaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pathprimary

	return p
}

func (s *PathprimaryContext) GetParser() antlr.Parser { return s.parser }

func (s *PathprimaryContext) CopyFrom(ctx *PathprimaryContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *PathprimaryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathprimaryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type AtomicpathContext struct {
	*PathprimaryContext
}

func NewAtomicpathContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AtomicpathContext {
	var p = new(AtomicpathContext)

	p.PathprimaryContext = NewEmptyPathprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PathprimaryContext))

	return p
}

func (s *AtomicpathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AtomicpathContext) Pathatom() IPathatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathatomContext)
}

func (s *AtomicpathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterAtomicpath(s)
	}
}

func (s *AtomicpathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitAtomicpath(s)
	}
}

type SubpathContext struct {
	*PathprimaryContext
}

func NewSubpathContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SubpathContext {
	var p = new(SubpathContext)

	p.PathprimaryContext = NewEmptyPathprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PathprimaryContext))

	return p
}

func (s *SubpathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubpathContext) SUBPATH() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserSUBPATH, 0)
}

func (s *SubpathContext) Pairtertiary() IPairtertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPairtertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPairtertiaryContext)
}

func (s *SubpathContext) OF() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserOF, 0)
}

func (s *SubpathContext) Pathprimary() IPathprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathprimaryContext)
}

func (s *SubpathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterSubpath(s)
	}
}

func (s *SubpathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitSubpath(s)
	}
}

type ReversepathContext struct {
	*PathprimaryContext
}

func NewReversepathContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ReversepathContext {
	var p = new(ReversepathContext)

	p.PathprimaryContext = NewEmptyPathprimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PathprimaryContext))

	return p
}

func (s *ReversepathContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ReversepathContext) REVERSE() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserREVERSE, 0)
}

func (s *ReversepathContext) Pathprimary() IPathprimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathprimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathprimaryContext)
}

func (s *ReversepathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterReversepath(s)
	}
}

func (s *ReversepathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitReversepath(s)
	}
}

func (p *PMMPStatemParser) Pathprimary() (localctx IPathprimaryContext) {
	localctx = NewPathprimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, PMMPStatemParserRULE_pathprimary)

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

	p.SetState(471)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPStatemParserPATHTAG:
		localctx = NewAtomicpathContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(463)
			p.Pathatom()
		}

	case PMMPStatemParserREVERSE:
		localctx = NewReversepathContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(464)
			p.Match(PMMPStatemParserREVERSE)
		}
		{
			p.SetState(465)
			p.Pathprimary()
		}

	case PMMPStatemParserSUBPATH:
		localctx = NewSubpathContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(466)
			p.Match(PMMPStatemParserSUBPATH)
		}
		{
			p.SetState(467)
			p.pairtertiary(0)
		}
		{
			p.SetState(468)
			p.Match(PMMPStatemParserOF)
		}
		{
			p.SetState(469)
			p.Pathprimary()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IPathatomContext is an interface to support dynamic dispatch.
type IPathatomContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPathatomContext differentiates from other interfaces.
	IsPathatomContext()
}

type PathatomContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathatomContext() *PathatomContext {
	var p = new(PathatomContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_pathatom
	return p
}

func (*PathatomContext) IsPathatomContext() {}

func NewPathatomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathatomContext {
	var p = new(PathatomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_pathatom

	return p
}

func (s *PathatomContext) GetParser() antlr.Parser { return s.parser }

func (s *PathatomContext) CopyFrom(ctx *PathatomContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *PathatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathatomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type PathvariableContext struct {
	*PathatomContext
}

func NewPathvariableContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PathvariableContext {
	var p = new(PathvariableContext)

	p.PathatomContext = NewEmptyPathatomContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PathatomContext))

	return p
}

func (s *PathvariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathvariableContext) PATHTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPATHTAG, 0)
}

func (s *PathvariableContext) AllSubscript() []ISubscriptContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISubscriptContext)(nil)).Elem())
	var tst = make([]ISubscriptContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISubscriptContext)
		}
	}

	return tst
}

func (s *PathvariableContext) Subscript(i int) ISubscriptContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISubscriptContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISubscriptContext)
}

func (s *PathvariableContext) AllAnytag() []IAnytagContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAnytagContext)(nil)).Elem())
	var tst = make([]IAnytagContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAnytagContext)
		}
	}

	return tst
}

func (s *PathvariableContext) Anytag(i int) IAnytagContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnytagContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAnytagContext)
}

func (s *PathvariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterPathvariable(s)
	}
}

func (s *PathvariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitPathvariable(s)
	}
}

func (p *PMMPStatemParser) Pathatom() (localctx IPathatomContext) {
	localctx = NewPathatomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, PMMPStatemParserRULE_pathatom)

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

	localctx = NewPathvariableContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(473)
		p.Match(PMMPStatemParserPATHTAG)
	}
	p.SetState(478)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(476)
			p.GetErrorHandler().Sync(p)

			switch p.GetTokenStream().LA(1) {
			case PMMPStatemParserLBRACKET, PMMPStatemParserDECIMALTOKEN:
				{
					p.SetState(474)
					p.Subscript()
				}

			case PMMPStatemParserPTAG, PMMPStatemParserTAG, PMMPStatemParserMIXEDPTAG, PMMPStatemParserMIXEDTAG:
				{
					p.SetState(475)
					p.Anytag()
				}

			default:
				panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			}

		}
		p.SetState(480)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 51, p.GetParserRuleContext())
	}

	return localctx
}

// ITagContext is an interface to support dynamic dispatch.
type ITagContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTagContext differentiates from other interfaces.
	IsTagContext()
}

type TagContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTagContext() *TagContext {
	var p = new(TagContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPStatemParserRULE_tag
	return p
}

func (*TagContext) IsTagContext() {}

func NewTagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TagContext {
	var p = new(TagContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_tag

	return p
}

func (s *TagContext) GetParser() antlr.Parser { return s.parser }

func (s *TagContext) TAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserTAG, 0)
}

func (s *TagContext) PTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPTAG, 0)
}

func (s *TagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterTag(s)
	}
}

func (s *TagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitTag(s)
	}
}

func (p *PMMPStatemParser) Tag() (localctx ITagContext) {
	localctx = NewTagContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, PMMPStatemParserRULE_tag)
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
	p.SetState(481)
	_la = p.GetTokenStream().LA(1)

	if !(_la == PMMPStatemParserPTAG || _la == PMMPStatemParserTAG) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
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
	p.RuleIndex = PMMPStatemParserRULE_anytag
	return p
}

func (*AnytagContext) IsAnytagContext() {}

func NewAnytagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnytagContext {
	var p = new(AnytagContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPStatemParserRULE_anytag

	return p
}

func (s *AnytagContext) GetParser() antlr.Parser { return s.parser }

func (s *AnytagContext) TAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserTAG, 0)
}

func (s *AnytagContext) MIXEDTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserMIXEDTAG, 0)
}

func (s *AnytagContext) PTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserPTAG, 0)
}

func (s *AnytagContext) MIXEDPTAG() antlr.TerminalNode {
	return s.GetToken(PMMPStatemParserMIXEDPTAG, 0)
}

func (s *AnytagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnytagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AnytagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.EnterAnytag(s)
	}
}

func (s *AnytagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPStatemListener); ok {
		listenerT.ExitAnytag(s)
	}
}

func (p *PMMPStatemParser) Anytag() (localctx IAnytagContext) {
	localctx = NewAnytagContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, PMMPStatemParserRULE_anytag)
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
	p.SetState(483)
	_la = p.GetTokenStream().LA(1)

	if !(((_la-47)&-(0x1f+1)) == 0 && ((1<<uint((_la-47)))&((1<<(PMMPStatemParserPTAG-47))|(1<<(PMMPStatemParserTAG-47))|(1<<(PMMPStatemParserMIXEDPTAG-47))|(1<<(PMMPStatemParserMIXEDTAG-47)))) != 0) {
		p.GetErrorHandler().RecoverInline(p)
	} else {
		p.GetErrorHandler().ReportMatch(p)
		p.Consume()
	}

	return localctx
}

func (p *PMMPStatemParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 19:
		var t *NumtertiaryContext = nil
		if localctx != nil {
			t = localctx.(*NumtertiaryContext)
		}
		return p.Numtertiary_Sempred(t, predIndex)

	case 20:
		var t *NumsecondaryContext = nil
		if localctx != nil {
			t = localctx.(*NumsecondaryContext)
		}
		return p.Numsecondary_Sempred(t, predIndex)

	case 26:
		var t *PairtertiaryContext = nil
		if localctx != nil {
			t = localctx.(*PairtertiaryContext)
		}
		return p.Pairtertiary_Sempred(t, predIndex)

	case 27:
		var t *PairsecondaryContext = nil
		if localctx != nil {
			t = localctx.(*PairsecondaryContext)
		}
		return p.Pairsecondary_Sempred(t, predIndex)

	case 35:
		var t *PathsecondaryContext = nil
		if localctx != nil {
			t = localctx.(*PathsecondaryContext)
		}
		return p.Pathsecondary_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *PMMPStatemParser) Numtertiary_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *PMMPStatemParser) Numsecondary_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 1:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *PMMPStatemParser) Pairtertiary_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *PMMPStatemParser) Pairsecondary_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 3:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *PMMPStatemParser) Pathsecondary_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 5:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
