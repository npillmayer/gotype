// Code generated from PMMPost.g4 by ANTLR 4.7.2. DO NOT EDIT.

package grammar // PMMPost
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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 67, 406,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2,
	3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4,
	3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4,
	3, 4, 3, 4, 5, 4, 106, 10, 4, 3, 5, 3, 5, 3, 5, 3, 5, 7, 5, 112, 10, 5,
	12, 5, 14, 5, 115, 11, 5, 3, 5, 3, 5, 5, 5, 119, 10, 5, 3, 5, 3, 5, 3,
	5, 7, 5, 124, 10, 5, 12, 5, 14, 5, 127, 11, 5, 5, 5, 129, 10, 5, 3, 6,
	3, 6, 3, 6, 3, 6, 7, 6, 135, 10, 6, 12, 6, 14, 6, 138, 11, 6, 3, 6, 3,
	6, 3, 6, 3, 6, 7, 6, 144, 10, 6, 12, 6, 14, 6, 147, 11, 6, 3, 6, 3, 6,
	3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 5, 6, 159, 10, 6, 3, 7,
	3, 7, 3, 7, 3, 8, 3, 8, 3, 8, 3, 9, 3, 9, 3, 9, 3, 9, 5, 9, 171, 10, 9,
	3, 9, 3, 9, 5, 9, 175, 10, 9, 3, 10, 5, 10, 178, 10, 10, 3, 10, 3, 10,
	5, 10, 182, 10, 10, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3, 11, 3,
	11, 3, 11, 5, 11, 193, 10, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 5, 12,
	200, 10, 12, 3, 13, 3, 13, 3, 13, 3, 13, 5, 13, 206, 10, 13, 3, 14, 6,
	14, 209, 10, 14, 13, 14, 14, 14, 210, 3, 14, 3, 14, 3, 15, 3, 15, 3, 15,
	3, 15, 7, 15, 219, 10, 15, 12, 15, 14, 15, 222, 11, 15, 3, 16, 3, 16, 3,
	16, 3, 16, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 3, 18, 3, 19, 3, 19, 5, 19,
	236, 10, 19, 3, 20, 3, 20, 3, 20, 6, 20, 241, 10, 20, 13, 20, 14, 20, 242,
	3, 21, 3, 21, 3, 21, 6, 21, 248, 10, 21, 13, 21, 14, 21, 249, 3, 22, 3,
	22, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 7, 23, 260, 10, 23, 12, 23,
	14, 23, 263, 11, 23, 3, 24, 3, 24, 3, 24, 5, 24, 268, 10, 24, 3, 24, 3,
	24, 3, 24, 7, 24, 273, 10, 24, 12, 24, 14, 24, 276, 11, 24, 3, 25, 3, 25,
	3, 25, 6, 25, 281, 10, 25, 13, 25, 14, 25, 282, 3, 25, 5, 25, 286, 10,
	25, 3, 26, 3, 26, 3, 26, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27,
	3, 27, 3, 27, 6, 27, 300, 10, 27, 13, 27, 14, 27, 301, 7, 27, 304, 10,
	27, 12, 27, 14, 27, 307, 11, 27, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3,
	28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28,
	3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3,
	28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 5, 28, 343, 10, 28,
	3, 29, 5, 29, 346, 10, 29, 3, 29, 3, 29, 3, 30, 3, 30, 3, 30, 3, 30, 5,
	30, 354, 10, 30, 3, 31, 3, 31, 5, 31, 358, 10, 31, 3, 31, 3, 31, 3, 31,
	3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3,
	31, 3, 31, 3, 31, 5, 31, 376, 10, 31, 3, 32, 3, 32, 3, 32, 7, 32, 381,
	10, 32, 12, 32, 14, 32, 384, 11, 32, 3, 32, 3, 32, 3, 32, 7, 32, 389, 10,
	32, 12, 32, 14, 32, 392, 11, 32, 3, 32, 5, 32, 395, 10, 32, 3, 33, 3, 33,
	3, 33, 3, 33, 3, 33, 5, 33, 402, 10, 33, 3, 34, 3, 34, 3, 34, 2, 5, 44,
	46, 52, 35, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32,
	34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 2,
	7, 3, 2, 34, 36, 11, 2, 8, 8, 21, 21, 30, 36, 39, 40, 44, 47, 49, 50, 55,
	55, 57, 60, 62, 62, 3, 2, 30, 31, 3, 2, 32, 33, 3, 2, 62, 63, 2, 433, 2,
	68, 3, 2, 2, 2, 4, 79, 3, 2, 2, 2, 6, 105, 3, 2, 2, 2, 8, 128, 3, 2, 2,
	2, 10, 158, 3, 2, 2, 2, 12, 160, 3, 2, 2, 2, 14, 163, 3, 2, 2, 2, 16, 166,
	3, 2, 2, 2, 18, 177, 3, 2, 2, 2, 20, 192, 3, 2, 2, 2, 22, 199, 3, 2, 2,
	2, 24, 201, 3, 2, 2, 2, 26, 208, 3, 2, 2, 2, 28, 214, 3, 2, 2, 2, 30, 223,
	3, 2, 2, 2, 32, 227, 3, 2, 2, 2, 34, 229, 3, 2, 2, 2, 36, 235, 3, 2, 2,
	2, 38, 237, 3, 2, 2, 2, 40, 244, 3, 2, 2, 2, 42, 251, 3, 2, 2, 2, 44, 253,
	3, 2, 2, 2, 46, 267, 3, 2, 2, 2, 48, 277, 3, 2, 2, 2, 50, 287, 3, 2, 2,
	2, 52, 290, 3, 2, 2, 2, 54, 342, 3, 2, 2, 2, 56, 345, 3, 2, 2, 2, 58, 353,
	3, 2, 2, 2, 60, 375, 3, 2, 2, 2, 62, 394, 3, 2, 2, 2, 64, 401, 3, 2, 2,
	2, 66, 403, 3, 2, 2, 2, 68, 69, 7, 3, 2, 2, 69, 70, 7, 26, 2, 2, 70, 71,
	7, 66, 2, 2, 71, 72, 7, 25, 2, 2, 72, 73, 7, 64, 2, 2, 73, 74, 7, 37, 2,
	2, 74, 75, 7, 25, 2, 2, 75, 76, 7, 64, 2, 2, 76, 77, 7, 37, 2, 2, 77, 78,
	7, 27, 2, 2, 78, 3, 3, 2, 2, 2, 79, 80, 7, 14, 2, 2, 80, 5, 3, 2, 2, 2,
	81, 82, 5, 30, 16, 2, 82, 83, 7, 24, 2, 2, 83, 106, 3, 2, 2, 2, 84, 85,
	5, 8, 5, 2, 85, 86, 7, 24, 2, 2, 86, 106, 3, 2, 2, 2, 87, 88, 5, 34, 18,
	2, 88, 89, 7, 24, 2, 2, 89, 106, 3, 2, 2, 2, 90, 91, 5, 36, 19, 2, 91,
	92, 7, 24, 2, 2, 92, 106, 3, 2, 2, 2, 93, 94, 5, 10, 6, 2, 94, 95, 7, 24,
	2, 2, 95, 106, 3, 2, 2, 2, 96, 97, 5, 32, 17, 2, 97, 98, 7, 24, 2, 2, 98,
	106, 3, 2, 2, 2, 99, 100, 5, 2, 2, 2, 100, 101, 7, 24, 2, 2, 101, 106,
	3, 2, 2, 2, 102, 103, 5, 4, 3, 2, 103, 104, 7, 24, 2, 2, 104, 106, 3, 2,
	2, 2, 105, 81, 3, 2, 2, 2, 105, 84, 3, 2, 2, 2, 105, 87, 3, 2, 2, 2, 105,
	90, 3, 2, 2, 2, 105, 93, 3, 2, 2, 2, 105, 96, 3, 2, 2, 2, 105, 99, 3, 2,
	2, 2, 105, 102, 3, 2, 2, 2, 106, 7, 3, 2, 2, 2, 107, 108, 7, 12, 2, 2,
	108, 113, 7, 62, 2, 2, 109, 110, 7, 25, 2, 2, 110, 112, 7, 62, 2, 2, 111,
	109, 3, 2, 2, 2, 112, 115, 3, 2, 2, 2, 113, 111, 3, 2, 2, 2, 113, 114,
	3, 2, 2, 2, 114, 129, 3, 2, 2, 2, 115, 113, 3, 2, 2, 2, 116, 118, 7, 41,
	2, 2, 117, 119, 7, 12, 2, 2, 118, 117, 3, 2, 2, 2, 118, 119, 3, 2, 2, 2,
	119, 120, 3, 2, 2, 2, 120, 125, 7, 62, 2, 2, 121, 122, 7, 25, 2, 2, 122,
	124, 7, 62, 2, 2, 123, 121, 3, 2, 2, 2, 124, 127, 3, 2, 2, 2, 125, 123,
	3, 2, 2, 2, 125, 126, 3, 2, 2, 2, 126, 129, 3, 2, 2, 2, 127, 125, 3, 2,
	2, 2, 128, 107, 3, 2, 2, 2, 128, 116, 3, 2, 2, 2, 129, 9, 3, 2, 2, 2, 130,
	131, 7, 59, 2, 2, 131, 136, 7, 62, 2, 2, 132, 133, 7, 25, 2, 2, 133, 135,
	7, 62, 2, 2, 134, 132, 3, 2, 2, 2, 135, 138, 3, 2, 2, 2, 136, 134, 3, 2,
	2, 2, 136, 137, 3, 2, 2, 2, 137, 159, 3, 2, 2, 2, 138, 136, 3, 2, 2, 2,
	139, 140, 7, 60, 2, 2, 140, 145, 7, 62, 2, 2, 141, 142, 7, 25, 2, 2, 142,
	144, 7, 62, 2, 2, 143, 141, 3, 2, 2, 2, 144, 147, 3, 2, 2, 2, 145, 143,
	3, 2, 2, 2, 145, 146, 3, 2, 2, 2, 146, 159, 3, 2, 2, 2, 147, 145, 3, 2,
	2, 2, 148, 149, 7, 58, 2, 2, 149, 159, 7, 66, 2, 2, 150, 151, 7, 61, 2,
	2, 151, 152, 5, 42, 22, 2, 152, 153, 7, 22, 2, 2, 153, 154, 7, 48, 2, 2,
	154, 159, 3, 2, 2, 2, 155, 159, 5, 16, 9, 2, 156, 159, 5, 12, 7, 2, 157,
	159, 5, 14, 8, 2, 158, 130, 3, 2, 2, 2, 158, 139, 3, 2, 2, 2, 158, 148,
	3, 2, 2, 2, 158, 150, 3, 2, 2, 2, 158, 155, 3, 2, 2, 2, 158, 156, 3, 2,
	2, 2, 158, 157, 3, 2, 2, 2, 159, 11, 3, 2, 2, 2, 160, 161, 7, 18, 2, 2,
	161, 162, 5, 44, 23, 2, 162, 13, 3, 2, 2, 2, 163, 164, 7, 17, 2, 2, 164,
	165, 5, 44, 23, 2, 165, 15, 3, 2, 2, 2, 166, 167, 7, 15, 2, 2, 167, 170,
	7, 16, 2, 2, 168, 169, 7, 4, 2, 2, 169, 171, 7, 64, 2, 2, 170, 168, 3,
	2, 2, 2, 170, 171, 3, 2, 2, 2, 171, 174, 3, 2, 2, 2, 172, 173, 7, 19, 2,
	2, 173, 175, 7, 20, 2, 2, 174, 172, 3, 2, 2, 2, 174, 175, 3, 2, 2, 2, 175,
	17, 3, 2, 2, 2, 176, 178, 5, 20, 11, 2, 177, 176, 3, 2, 2, 2, 177, 178,
	3, 2, 2, 2, 178, 179, 3, 2, 2, 2, 179, 181, 5, 22, 12, 2, 180, 182, 5,
	20, 11, 2, 181, 180, 3, 2, 2, 2, 181, 182, 3, 2, 2, 2, 182, 19, 3, 2, 2,
	2, 183, 184, 7, 5, 2, 2, 184, 185, 7, 11, 2, 2, 185, 186, 5, 44, 23, 2,
	186, 187, 7, 6, 2, 2, 187, 193, 3, 2, 2, 2, 188, 189, 7, 5, 2, 2, 189,
	190, 5, 44, 23, 2, 190, 191, 7, 6, 2, 2, 191, 193, 3, 2, 2, 2, 192, 183,
	3, 2, 2, 2, 192, 188, 3, 2, 2, 2, 193, 21, 3, 2, 2, 2, 194, 200, 7, 8,
	2, 2, 195, 196, 7, 7, 2, 2, 196, 197, 5, 24, 13, 2, 197, 198, 7, 7, 2,
	2, 198, 200, 3, 2, 2, 2, 199, 194, 3, 2, 2, 2, 199, 195, 3, 2, 2, 2, 200,
	23, 3, 2, 2, 2, 201, 202, 7, 9, 2, 2, 202, 205, 5, 44, 23, 2, 203, 204,
	7, 10, 2, 2, 204, 206, 5, 44, 23, 2, 205, 203, 3, 2, 2, 2, 205, 206, 3,
	2, 2, 2, 206, 25, 3, 2, 2, 2, 207, 209, 5, 6, 4, 2, 208, 207, 3, 2, 2,
	2, 209, 210, 3, 2, 2, 2, 210, 208, 3, 2, 2, 2, 210, 211, 3, 2, 2, 2, 211,
	212, 3, 2, 2, 2, 212, 213, 7, 2, 2, 3, 213, 27, 3, 2, 2, 2, 214, 215, 7,
	42, 2, 2, 215, 220, 7, 62, 2, 2, 216, 217, 7, 25, 2, 2, 217, 219, 7, 62,
	2, 2, 218, 216, 3, 2, 2, 2, 219, 222, 3, 2, 2, 2, 220, 218, 3, 2, 2, 2,
	220, 221, 3, 2, 2, 2, 221, 29, 3, 2, 2, 2, 222, 220, 3, 2, 2, 2, 223, 224,
	7, 39, 2, 2, 224, 225, 5, 26, 14, 2, 225, 226, 7, 40, 2, 2, 226, 31, 3,
	2, 2, 2, 227, 228, 3, 2, 2, 2, 228, 33, 3, 2, 2, 2, 229, 230, 5, 62, 32,
	2, 230, 231, 7, 21, 2, 2, 231, 232, 5, 48, 25, 2, 232, 35, 3, 2, 2, 2,
	233, 236, 5, 38, 20, 2, 234, 236, 5, 40, 21, 2, 235, 233, 3, 2, 2, 2, 235,
	234, 3, 2, 2, 2, 236, 37, 3, 2, 2, 2, 237, 240, 5, 44, 23, 2, 238, 239,
	7, 22, 2, 2, 239, 241, 5, 44, 23, 2, 240, 238, 3, 2, 2, 2, 241, 242, 3,
	2, 2, 2, 242, 240, 3, 2, 2, 2, 242, 243, 3, 2, 2, 2, 243, 39, 3, 2, 2,
	2, 244, 247, 5, 46, 24, 2, 245, 246, 9, 2, 2, 2, 246, 248, 5, 46, 24, 2,
	247, 245, 3, 2, 2, 2, 248, 249, 3, 2, 2, 2, 249, 247, 3, 2, 2, 2, 249,
	250, 3, 2, 2, 2, 250, 41, 3, 2, 2, 2, 251, 252, 9, 3, 2, 2, 252, 43, 3,
	2, 2, 2, 253, 254, 8, 23, 1, 2, 254, 255, 5, 46, 24, 2, 255, 261, 3, 2,
	2, 2, 256, 257, 12, 3, 2, 2, 257, 258, 7, 57, 2, 2, 258, 260, 5, 46, 24,
	2, 259, 256, 3, 2, 2, 2, 260, 263, 3, 2, 2, 2, 261, 259, 3, 2, 2, 2, 261,
	262, 3, 2, 2, 2, 262, 45, 3, 2, 2, 2, 263, 261, 3, 2, 2, 2, 264, 265, 8,
	24, 1, 2, 265, 268, 5, 48, 25, 2, 266, 268, 5, 52, 27, 2, 267, 264, 3,
	2, 2, 2, 267, 266, 3, 2, 2, 2, 268, 274, 3, 2, 2, 2, 269, 270, 12, 4, 2,
	2, 270, 271, 9, 4, 2, 2, 271, 273, 5, 52, 27, 2, 272, 269, 3, 2, 2, 2,
	273, 276, 3, 2, 2, 2, 274, 272, 3, 2, 2, 2, 274, 275, 3, 2, 2, 2, 275,
	47, 3, 2, 2, 2, 276, 274, 3, 2, 2, 2, 277, 280, 5, 52, 27, 2, 278, 279,
	7, 7, 2, 2, 279, 281, 5, 52, 27, 2, 280, 278, 3, 2, 2, 2, 281, 282, 3,
	2, 2, 2, 282, 280, 3, 2, 2, 2, 282, 283, 3, 2, 2, 2, 283, 285, 3, 2, 2,
	2, 284, 286, 5, 50, 26, 2, 285, 284, 3, 2, 2, 2, 285, 286, 3, 2, 2, 2,
	286, 49, 3, 2, 2, 2, 287, 288, 5, 18, 10, 2, 288, 289, 7, 56, 2, 2, 289,
	51, 3, 2, 2, 2, 290, 291, 8, 27, 1, 2, 291, 292, 5, 54, 28, 2, 292, 305,
	3, 2, 2, 2, 293, 294, 12, 4, 2, 2, 294, 295, 9, 5, 2, 2, 295, 304, 5, 54,
	28, 2, 296, 299, 12, 3, 2, 2, 297, 298, 7, 55, 2, 2, 298, 300, 5, 54, 28,
	2, 299, 297, 3, 2, 2, 2, 300, 301, 3, 2, 2, 2, 301, 299, 3, 2, 2, 2, 301,
	302, 3, 2, 2, 2, 302, 304, 3, 2, 2, 2, 303, 293, 3, 2, 2, 2, 303, 296,
	3, 2, 2, 2, 304, 307, 3, 2, 2, 2, 305, 303, 3, 2, 2, 2, 305, 306, 3, 2,
	2, 2, 306, 53, 3, 2, 2, 2, 307, 305, 3, 2, 2, 2, 308, 309, 7, 48, 2, 2,
	309, 343, 5, 60, 31, 2, 310, 311, 5, 56, 29, 2, 311, 312, 5, 60, 31, 2,
	312, 343, 3, 2, 2, 2, 313, 314, 5, 58, 30, 2, 314, 315, 7, 28, 2, 2, 315,
	316, 5, 46, 24, 2, 316, 317, 7, 25, 2, 2, 317, 318, 5, 46, 24, 2, 318,
	319, 7, 29, 2, 2, 319, 343, 3, 2, 2, 2, 320, 321, 5, 60, 31, 2, 321, 322,
	7, 28, 2, 2, 322, 323, 5, 46, 24, 2, 323, 324, 7, 25, 2, 2, 324, 325, 5,
	46, 24, 2, 325, 326, 7, 29, 2, 2, 326, 343, 3, 2, 2, 2, 327, 343, 5, 60,
	31, 2, 328, 329, 7, 43, 2, 2, 329, 343, 5, 54, 28, 2, 330, 331, 7, 52,
	2, 2, 331, 332, 5, 46, 24, 2, 332, 333, 7, 53, 2, 2, 333, 334, 5, 54, 28,
	2, 334, 343, 3, 2, 2, 2, 335, 336, 7, 50, 2, 2, 336, 343, 5, 54, 28, 2,
	337, 338, 7, 49, 2, 2, 338, 339, 5, 46, 24, 2, 339, 340, 7, 53, 2, 2, 340,
	341, 5, 54, 28, 2, 341, 343, 3, 2, 2, 2, 342, 308, 3, 2, 2, 2, 342, 310,
	3, 2, 2, 2, 342, 313, 3, 2, 2, 2, 342, 320, 3, 2, 2, 2, 342, 327, 3, 2,
	2, 2, 342, 328, 3, 2, 2, 2, 342, 330, 3, 2, 2, 2, 342, 335, 3, 2, 2, 2,
	342, 337, 3, 2, 2, 2, 343, 55, 3, 2, 2, 2, 344, 346, 9, 4, 2, 2, 345, 344,
	3, 2, 2, 2, 345, 346, 3, 2, 2, 2, 346, 347, 3, 2, 2, 2, 347, 348, 5, 58,
	30, 2, 348, 57, 3, 2, 2, 2, 349, 350, 7, 64, 2, 2, 350, 351, 7, 33, 2,
	2, 351, 354, 7, 64, 2, 2, 352, 354, 7, 64, 2, 2, 353, 349, 3, 2, 2, 2,
	353, 352, 3, 2, 2, 2, 354, 59, 3, 2, 2, 2, 355, 357, 7, 64, 2, 2, 356,
	358, 7, 37, 2, 2, 357, 356, 3, 2, 2, 2, 357, 358, 3, 2, 2, 2, 358, 376,
	3, 2, 2, 2, 359, 376, 5, 62, 32, 2, 360, 361, 7, 26, 2, 2, 361, 362, 5,
	46, 24, 2, 362, 363, 7, 25, 2, 2, 363, 364, 5, 46, 24, 2, 364, 365, 7,
	27, 2, 2, 365, 376, 3, 2, 2, 2, 366, 367, 7, 26, 2, 2, 367, 368, 5, 46,
	24, 2, 368, 369, 7, 27, 2, 2, 369, 376, 3, 2, 2, 2, 370, 371, 7, 39, 2,
	2, 371, 372, 5, 26, 14, 2, 372, 373, 5, 46, 24, 2, 373, 374, 7, 40, 2,
	2, 374, 376, 3, 2, 2, 2, 375, 355, 3, 2, 2, 2, 375, 359, 3, 2, 2, 2, 375,
	360, 3, 2, 2, 2, 375, 366, 3, 2, 2, 2, 375, 370, 3, 2, 2, 2, 376, 61, 3,
	2, 2, 2, 377, 382, 7, 63, 2, 2, 378, 381, 5, 64, 33, 2, 379, 381, 5, 66,
	34, 2, 380, 378, 3, 2, 2, 2, 380, 379, 3, 2, 2, 2, 381, 384, 3, 2, 2, 2,
	382, 380, 3, 2, 2, 2, 382, 383, 3, 2, 2, 2, 383, 395, 3, 2, 2, 2, 384,
	382, 3, 2, 2, 2, 385, 390, 7, 62, 2, 2, 386, 389, 5, 64, 33, 2, 387, 389,
	5, 66, 34, 2, 388, 386, 3, 2, 2, 2, 388, 387, 3, 2, 2, 2, 389, 392, 3,
	2, 2, 2, 390, 388, 3, 2, 2, 2, 390, 391, 3, 2, 2, 2, 391, 395, 3, 2, 2,
	2, 392, 390, 3, 2, 2, 2, 393, 395, 7, 38, 2, 2, 394, 377, 3, 2, 2, 2, 394,
	385, 3, 2, 2, 2, 394, 393, 3, 2, 2, 2, 395, 63, 3, 2, 2, 2, 396, 402, 7,
	64, 2, 2, 397, 398, 7, 28, 2, 2, 398, 399, 5, 46, 24, 2, 399, 400, 7, 29,
	2, 2, 400, 402, 3, 2, 2, 2, 401, 396, 3, 2, 2, 2, 401, 397, 3, 2, 2, 2,
	402, 65, 3, 2, 2, 2, 403, 404, 9, 6, 2, 2, 404, 67, 3, 2, 2, 2, 41, 105,
	113, 118, 125, 128, 136, 145, 158, 170, 174, 177, 181, 192, 199, 205, 210,
	220, 235, 242, 249, 261, 267, 274, 282, 285, 301, 303, 305, 342, 345, 353,
	357, 375, 380, 382, 388, 390, 394, 401,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'beginfig'", "'scaled'", "'{'", "'}'", "'..'", "", "'controls'", "'and'",
	"'curl'", "", "'figure'", "'endfig'", "'pickup'", "", "'fill'", "'draw'",
	"'withcolor'", "", "':='", "'='", "':'", "';'", "','", "'('", "')'", "'['",
	"']'", "'+'", "'-'", "'*'", "'/'", "'||'", "'|-'", "'~'", "", "'@'", "'begingroup'",
	"'endgroup'", "'local'", "'vardef'", "", "", "'edge'", "'frame'", "'box'",
	"", "'subpath'", "'reverse'", "'with'", "'point'", "'of'", "'to'", "",
	"'cycle'", "", "'proof'", "'save'", "'show'", "'let'", "", "", "", "'.'",
}
var symbolicNames = []string{
	"", "", "", "", "", "", "PATHJOIN", "CONTROLS", "AND", "CURL", "TYPE",
	"FIGURE", "ENDFIG", "PICKUP", "PEN", "FILL", "DRAW", "WITHCOLOR", "COLOR",
	"ASSIGN", "EQUALS", "COLON", "SEMIC", "COMMA", "LPAREN", "RPAREN", "LBRACKET",
	"RBRACKET", "PLUS", "MINUS", "TIMES", "OVER", "PARALLEL", "PERPENDIC",
	"CONGRUENT", "UNIT", "LAMBDAARG", "BEGINGROUP", "ENDGROUP", "LOCAL", "VARDEF",
	"PAIRPART", "EDGECONSTR", "EDGE", "FRAME", "BOX", "MATHFUNC", "SUBPATH",
	"REVERSE", "WITH", "POINT", "OF", "TO", "TRANSFORM", "CYCLE", "PATHCLIPOP",
	"PROOF", "SAVE", "SHOW", "LET", "TAG", "MIXEDTAG", "DECIMALTOKEN", "DOT",
	"LABEL", "WS",
}

var ruleNames = []string{
	"beginfig", "endfig", "statement", "declaration", "command", "drawCmd",
	"fillCmd", "pickupCmd", "pathjoin", "directionspec", "basicpathjoin", "controls",
	"statementlist", "vardef", "compound", "empty", "assignment", "constraint",
	"equation", "orientation", "token", "expression", "tertiary", "path", "cycle",
	"secondary", "primary", "scalarmulop", "numtokenatom", "atom", "variable",
	"subscript", "anytag",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type PMMPostParser struct {
	*antlr.BaseParser
}

func NewPMMPostParser(input antlr.TokenStream) *PMMPostParser {
	this := new(PMMPostParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "PMMPost.g4"

	return this
}

// PMMPostParser tokens.
const (
	PMMPostParserEOF          = antlr.TokenEOF
	PMMPostParserT__0         = 1
	PMMPostParserT__1         = 2
	PMMPostParserT__2         = 3
	PMMPostParserT__3         = 4
	PMMPostParserT__4         = 5
	PMMPostParserPATHJOIN     = 6
	PMMPostParserCONTROLS     = 7
	PMMPostParserAND          = 8
	PMMPostParserCURL         = 9
	PMMPostParserTYPE         = 10
	PMMPostParserFIGURE       = 11
	PMMPostParserENDFIG       = 12
	PMMPostParserPICKUP       = 13
	PMMPostParserPEN          = 14
	PMMPostParserFILL         = 15
	PMMPostParserDRAW         = 16
	PMMPostParserWITHCOLOR    = 17
	PMMPostParserCOLOR        = 18
	PMMPostParserASSIGN       = 19
	PMMPostParserEQUALS       = 20
	PMMPostParserCOLON        = 21
	PMMPostParserSEMIC        = 22
	PMMPostParserCOMMA        = 23
	PMMPostParserLPAREN       = 24
	PMMPostParserRPAREN       = 25
	PMMPostParserLBRACKET     = 26
	PMMPostParserRBRACKET     = 27
	PMMPostParserPLUS         = 28
	PMMPostParserMINUS        = 29
	PMMPostParserTIMES        = 30
	PMMPostParserOVER         = 31
	PMMPostParserPARALLEL     = 32
	PMMPostParserPERPENDIC    = 33
	PMMPostParserCONGRUENT    = 34
	PMMPostParserUNIT         = 35
	PMMPostParserLAMBDAARG    = 36
	PMMPostParserBEGINGROUP   = 37
	PMMPostParserENDGROUP     = 38
	PMMPostParserLOCAL        = 39
	PMMPostParserVARDEF       = 40
	PMMPostParserPAIRPART     = 41
	PMMPostParserEDGECONSTR   = 42
	PMMPostParserEDGE         = 43
	PMMPostParserFRAME        = 44
	PMMPostParserBOX          = 45
	PMMPostParserMATHFUNC     = 46
	PMMPostParserSUBPATH      = 47
	PMMPostParserREVERSE      = 48
	PMMPostParserWITH         = 49
	PMMPostParserPOINT        = 50
	PMMPostParserOF           = 51
	PMMPostParserTO           = 52
	PMMPostParserTRANSFORM    = 53
	PMMPostParserCYCLE        = 54
	PMMPostParserPATHCLIPOP   = 55
	PMMPostParserPROOF        = 56
	PMMPostParserSAVE         = 57
	PMMPostParserSHOW         = 58
	PMMPostParserLET          = 59
	PMMPostParserTAG          = 60
	PMMPostParserMIXEDTAG     = 61
	PMMPostParserDECIMALTOKEN = 62
	PMMPostParserDOT          = 63
	PMMPostParserLABEL        = 64
	PMMPostParserWS           = 65
)

// PMMPostParser rules.
const (
	PMMPostParserRULE_beginfig      = 0
	PMMPostParserRULE_endfig        = 1
	PMMPostParserRULE_statement     = 2
	PMMPostParserRULE_declaration   = 3
	PMMPostParserRULE_command       = 4
	PMMPostParserRULE_drawCmd       = 5
	PMMPostParserRULE_fillCmd       = 6
	PMMPostParserRULE_pickupCmd     = 7
	PMMPostParserRULE_pathjoin      = 8
	PMMPostParserRULE_directionspec = 9
	PMMPostParserRULE_basicpathjoin = 10
	PMMPostParserRULE_controls      = 11
	PMMPostParserRULE_statementlist = 12
	PMMPostParserRULE_vardef        = 13
	PMMPostParserRULE_compound      = 14
	PMMPostParserRULE_empty         = 15
	PMMPostParserRULE_assignment    = 16
	PMMPostParserRULE_constraint    = 17
	PMMPostParserRULE_equation      = 18
	PMMPostParserRULE_orientation   = 19
	PMMPostParserRULE_token         = 20
	PMMPostParserRULE_expression    = 21
	PMMPostParserRULE_tertiary      = 22
	PMMPostParserRULE_path          = 23
	PMMPostParserRULE_cycle         = 24
	PMMPostParserRULE_secondary     = 25
	PMMPostParserRULE_primary       = 26
	PMMPostParserRULE_scalarmulop   = 27
	PMMPostParserRULE_numtokenatom  = 28
	PMMPostParserRULE_atom          = 29
	PMMPostParserRULE_variable      = 30
	PMMPostParserRULE_subscript     = 31
	PMMPostParserRULE_anytag        = 32
)

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
	p.RuleIndex = PMMPostParserRULE_beginfig
	return p
}

func (*BeginfigContext) IsBeginfigContext() {}

func NewBeginfigContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BeginfigContext {
	var p = new(BeginfigContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_beginfig

	return p
}

func (s *BeginfigContext) GetParser() antlr.Parser { return s.parser }

func (s *BeginfigContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(PMMPostParserLPAREN, 0)
}

func (s *BeginfigContext) LABEL() antlr.TerminalNode {
	return s.GetToken(PMMPostParserLABEL, 0)
}

func (s *BeginfigContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserCOMMA)
}

func (s *BeginfigContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserCOMMA, i)
}

func (s *BeginfigContext) AllDECIMALTOKEN() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserDECIMALTOKEN)
}

func (s *BeginfigContext) DECIMALTOKEN(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserDECIMALTOKEN, i)
}

func (s *BeginfigContext) AllUNIT() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserUNIT)
}

func (s *BeginfigContext) UNIT(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserUNIT, i)
}

func (s *BeginfigContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PMMPostParserRPAREN, 0)
}

func (s *BeginfigContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BeginfigContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BeginfigContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterBeginfig(s)
	}
}

func (s *BeginfigContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitBeginfig(s)
	}
}

func (p *PMMPostParser) Beginfig() (localctx IBeginfigContext) {
	localctx = NewBeginfigContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, PMMPostParserRULE_beginfig)

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
		p.SetState(66)
		p.Match(PMMPostParserT__0)
	}
	{
		p.SetState(67)
		p.Match(PMMPostParserLPAREN)
	}
	{
		p.SetState(68)
		p.Match(PMMPostParserLABEL)
	}
	{
		p.SetState(69)
		p.Match(PMMPostParserCOMMA)
	}
	{
		p.SetState(70)
		p.Match(PMMPostParserDECIMALTOKEN)
	}
	{
		p.SetState(71)
		p.Match(PMMPostParserUNIT)
	}
	{
		p.SetState(72)
		p.Match(PMMPostParserCOMMA)
	}
	{
		p.SetState(73)
		p.Match(PMMPostParserDECIMALTOKEN)
	}
	{
		p.SetState(74)
		p.Match(PMMPostParserUNIT)
	}
	{
		p.SetState(75)
		p.Match(PMMPostParserRPAREN)
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
	p.RuleIndex = PMMPostParserRULE_endfig
	return p
}

func (*EndfigContext) IsEndfigContext() {}

func NewEndfigContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EndfigContext {
	var p = new(EndfigContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_endfig

	return p
}

func (s *EndfigContext) GetParser() antlr.Parser { return s.parser }

func (s *EndfigContext) ENDFIG() antlr.TerminalNode {
	return s.GetToken(PMMPostParserENDFIG, 0)
}

func (s *EndfigContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EndfigContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EndfigContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterEndfig(s)
	}
}

func (s *EndfigContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitEndfig(s)
	}
}

func (p *PMMPostParser) Endfig() (localctx IEndfigContext) {
	localctx = NewEndfigContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, PMMPostParserRULE_endfig)

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
		p.SetState(77)
		p.Match(PMMPostParserENDFIG)
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
	p.RuleIndex = PMMPostParserRULE_statement
	return p
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_statement

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

func (s *StatementContext) SEMIC() antlr.TerminalNode {
	return s.GetToken(PMMPostParserSEMIC, 0)
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

func (s *StatementContext) Empty() IEmptyContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEmptyContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEmptyContext)
}

func (s *StatementContext) Beginfig() IBeginfigContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBeginfigContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBeginfigContext)
}

func (s *StatementContext) Endfig() IEndfigContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEndfigContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEndfigContext)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (p *PMMPostParser) Statement() (localctx IStatementContext) {
	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, PMMPostParserRULE_statement)

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

	p.SetState(103)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(79)
			p.Compound()
		}
		{
			p.SetState(80)
			p.Match(PMMPostParserSEMIC)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(82)
			p.Declaration()
		}
		{
			p.SetState(83)
			p.Match(PMMPostParserSEMIC)
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(85)
			p.Assignment()
		}
		{
			p.SetState(86)
			p.Match(PMMPostParserSEMIC)
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(88)
			p.Constraint()
		}
		{
			p.SetState(89)
			p.Match(PMMPostParserSEMIC)
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(91)
			p.Command()
		}
		{
			p.SetState(92)
			p.Match(PMMPostParserSEMIC)
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(94)
			p.Empty()
		}
		{
			p.SetState(95)
			p.Match(PMMPostParserSEMIC)
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(97)
			p.Beginfig()
		}
		{
			p.SetState(98)
			p.Match(PMMPostParserSEMIC)
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(100)
			p.Endfig()
		}
		{
			p.SetState(101)
			p.Match(PMMPostParserSEMIC)
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
	p.RuleIndex = PMMPostParserRULE_declaration
	return p
}

func (*DeclarationContext) IsDeclarationContext() {}

func NewDeclarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclarationContext {
	var p = new(DeclarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_declaration

	return p
}

func (s *DeclarationContext) GetParser() antlr.Parser { return s.parser }

func (s *DeclarationContext) CopyFrom(ctx *DeclarationContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *DeclarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DeclarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type TypedeclContext struct {
	*DeclarationContext
}

func NewTypedeclContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TypedeclContext {
	var p = new(TypedeclContext)

	p.DeclarationContext = NewEmptyDeclarationContext()
	p.parser = parser
	p.CopyFrom(ctx.(*DeclarationContext))

	return p
}

func (s *TypedeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypedeclContext) TYPE() antlr.TerminalNode {
	return s.GetToken(PMMPostParserTYPE, 0)
}

func (s *TypedeclContext) AllTAG() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserTAG)
}

func (s *TypedeclContext) TAG(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserTAG, i)
}

func (s *TypedeclContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserCOMMA)
}

func (s *TypedeclContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserCOMMA, i)
}

func (s *TypedeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterTypedecl(s)
	}
}

func (s *TypedeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitTypedecl(s)
	}
}

type LocaldeclContext struct {
	*DeclarationContext
}

func NewLocaldeclContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LocaldeclContext {
	var p = new(LocaldeclContext)

	p.DeclarationContext = NewEmptyDeclarationContext()
	p.parser = parser
	p.CopyFrom(ctx.(*DeclarationContext))

	return p
}

func (s *LocaldeclContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LocaldeclContext) LOCAL() antlr.TerminalNode {
	return s.GetToken(PMMPostParserLOCAL, 0)
}

func (s *LocaldeclContext) AllTAG() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserTAG)
}

func (s *LocaldeclContext) TAG(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserTAG, i)
}

func (s *LocaldeclContext) TYPE() antlr.TerminalNode {
	return s.GetToken(PMMPostParserTYPE, 0)
}

func (s *LocaldeclContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserCOMMA)
}

func (s *LocaldeclContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserCOMMA, i)
}

func (s *LocaldeclContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterLocaldecl(s)
	}
}

func (s *LocaldeclContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitLocaldecl(s)
	}
}

func (p *PMMPostParser) Declaration() (localctx IDeclarationContext) {
	localctx = NewDeclarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, PMMPostParserRULE_declaration)
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

	p.SetState(126)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPostParserTYPE:
		localctx = NewTypedeclContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(105)
			p.Match(PMMPostParserTYPE)
		}
		{
			p.SetState(106)
			p.Match(PMMPostParserTAG)
		}
		p.SetState(111)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == PMMPostParserCOMMA {
			{
				p.SetState(107)
				p.Match(PMMPostParserCOMMA)
			}
			{
				p.SetState(108)
				p.Match(PMMPostParserTAG)
			}

			p.SetState(113)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case PMMPostParserLOCAL:
		localctx = NewLocaldeclContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(114)
			p.Match(PMMPostParserLOCAL)
		}
		p.SetState(116)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == PMMPostParserTYPE {
			{
				p.SetState(115)
				p.Match(PMMPostParserTYPE)
			}

		}
		{
			p.SetState(118)
			p.Match(PMMPostParserTAG)
		}
		p.SetState(123)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == PMMPostParserCOMMA {
			{
				p.SetState(119)
				p.Match(PMMPostParserCOMMA)
			}
			{
				p.SetState(120)
				p.Match(PMMPostParserTAG)
			}

			p.SetState(125)
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
	p.RuleIndex = PMMPostParserRULE_command
	return p
}

func (*CommandContext) IsCommandContext() {}

func NewCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandContext {
	var p = new(CommandContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_command

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
	return s.GetToken(PMMPostParserSHOW, 0)
}

func (s *ShowcmdContext) AllTAG() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserTAG)
}

func (s *ShowcmdContext) TAG(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserTAG, i)
}

func (s *ShowcmdContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserCOMMA)
}

func (s *ShowcmdContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserCOMMA, i)
}

func (s *ShowcmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterShowcmd(s)
	}
}

func (s *ShowcmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
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
	return s.GetToken(PMMPostParserPROOF, 0)
}

func (s *ProofcmdContext) LABEL() antlr.TerminalNode {
	return s.GetToken(PMMPostParserLABEL, 0)
}

func (s *ProofcmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterProofcmd(s)
	}
}

func (s *ProofcmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitProofcmd(s)
	}
}

type CmddrawContext struct {
	*CommandContext
}

func NewCmddrawContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CmddrawContext {
	var p = new(CmddrawContext)

	p.CommandContext = NewEmptyCommandContext()
	p.parser = parser
	p.CopyFrom(ctx.(*CommandContext))

	return p
}

func (s *CmddrawContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CmddrawContext) DrawCmd() IDrawCmdContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDrawCmdContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDrawCmdContext)
}

func (s *CmddrawContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterCmddraw(s)
	}
}

func (s *CmddrawContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitCmddraw(s)
	}
}

type CmdfillContext struct {
	*CommandContext
}

func NewCmdfillContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CmdfillContext {
	var p = new(CmdfillContext)

	p.CommandContext = NewEmptyCommandContext()
	p.parser = parser
	p.CopyFrom(ctx.(*CommandContext))

	return p
}

func (s *CmdfillContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CmdfillContext) FillCmd() IFillCmdContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFillCmdContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFillCmdContext)
}

func (s *CmdfillContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterCmdfill(s)
	}
}

func (s *CmdfillContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitCmdfill(s)
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
	return s.GetToken(PMMPostParserSAVE, 0)
}

func (s *SavecmdContext) AllTAG() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserTAG)
}

func (s *SavecmdContext) TAG(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserTAG, i)
}

func (s *SavecmdContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserCOMMA)
}

func (s *SavecmdContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserCOMMA, i)
}

func (s *SavecmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterSavecmd(s)
	}
}

func (s *SavecmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitSavecmd(s)
	}
}

type CmdpickupContext struct {
	*CommandContext
}

func NewCmdpickupContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CmdpickupContext {
	var p = new(CmdpickupContext)

	p.CommandContext = NewEmptyCommandContext()
	p.parser = parser
	p.CopyFrom(ctx.(*CommandContext))

	return p
}

func (s *CmdpickupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CmdpickupContext) PickupCmd() IPickupCmdContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPickupCmdContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPickupCmdContext)
}

func (s *CmdpickupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterCmdpickup(s)
	}
}

func (s *CmdpickupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitCmdpickup(s)
	}
}

type LetcmdContext struct {
	*CommandContext
}

func NewLetcmdContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LetcmdContext {
	var p = new(LetcmdContext)

	p.CommandContext = NewEmptyCommandContext()
	p.parser = parser
	p.CopyFrom(ctx.(*CommandContext))

	return p
}

func (s *LetcmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LetcmdContext) LET() antlr.TerminalNode {
	return s.GetToken(PMMPostParserLET, 0)
}

func (s *LetcmdContext) Token() ITokenContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITokenContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITokenContext)
}

func (s *LetcmdContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(PMMPostParserEQUALS, 0)
}

func (s *LetcmdContext) MATHFUNC() antlr.TerminalNode {
	return s.GetToken(PMMPostParserMATHFUNC, 0)
}

func (s *LetcmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterLetcmd(s)
	}
}

func (s *LetcmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitLetcmd(s)
	}
}

func (p *PMMPostParser) Command() (localctx ICommandContext) {
	localctx = NewCommandContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, PMMPostParserRULE_command)
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

	p.SetState(156)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPostParserSAVE:
		localctx = NewSavecmdContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(128)
			p.Match(PMMPostParserSAVE)
		}
		{
			p.SetState(129)
			p.Match(PMMPostParserTAG)
		}
		p.SetState(134)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == PMMPostParserCOMMA {
			{
				p.SetState(130)
				p.Match(PMMPostParserCOMMA)
			}
			{
				p.SetState(131)
				p.Match(PMMPostParserTAG)
			}

			p.SetState(136)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case PMMPostParserSHOW:
		localctx = NewShowcmdContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(137)
			p.Match(PMMPostParserSHOW)
		}
		{
			p.SetState(138)
			p.Match(PMMPostParserTAG)
		}
		p.SetState(143)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == PMMPostParserCOMMA {
			{
				p.SetState(139)
				p.Match(PMMPostParserCOMMA)
			}
			{
				p.SetState(140)
				p.Match(PMMPostParserTAG)
			}

			p.SetState(145)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case PMMPostParserPROOF:
		localctx = NewProofcmdContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(146)
			p.Match(PMMPostParserPROOF)
		}
		{
			p.SetState(147)
			p.Match(PMMPostParserLABEL)
		}

	case PMMPostParserLET:
		localctx = NewLetcmdContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(148)
			p.Match(PMMPostParserLET)
		}
		{
			p.SetState(149)
			p.Token()
		}
		{
			p.SetState(150)
			p.Match(PMMPostParserEQUALS)
		}
		{
			p.SetState(151)
			p.Match(PMMPostParserMATHFUNC)
		}

	case PMMPostParserPICKUP:
		localctx = NewCmdpickupContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(153)
			p.PickupCmd()
		}

	case PMMPostParserDRAW:
		localctx = NewCmddrawContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(154)
			p.DrawCmd()
		}

	case PMMPostParserFILL:
		localctx = NewCmdfillContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(155)
			p.FillCmd()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.RuleIndex = PMMPostParserRULE_drawCmd
	return p
}

func (*DrawCmdContext) IsDrawCmdContext() {}

func NewDrawCmdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DrawCmdContext {
	var p = new(DrawCmdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_drawCmd

	return p
}

func (s *DrawCmdContext) GetParser() antlr.Parser { return s.parser }

func (s *DrawCmdContext) DRAW() antlr.TerminalNode {
	return s.GetToken(PMMPostParserDRAW, 0)
}

func (s *DrawCmdContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *DrawCmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DrawCmdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DrawCmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterDrawCmd(s)
	}
}

func (s *DrawCmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitDrawCmd(s)
	}
}

func (p *PMMPostParser) DrawCmd() (localctx IDrawCmdContext) {
	localctx = NewDrawCmdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, PMMPostParserRULE_drawCmd)

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
		p.SetState(158)
		p.Match(PMMPostParserDRAW)
	}
	{
		p.SetState(159)
		p.expression(0)
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
	p.RuleIndex = PMMPostParserRULE_fillCmd
	return p
}

func (*FillCmdContext) IsFillCmdContext() {}

func NewFillCmdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FillCmdContext {
	var p = new(FillCmdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_fillCmd

	return p
}

func (s *FillCmdContext) GetParser() antlr.Parser { return s.parser }

func (s *FillCmdContext) FILL() antlr.TerminalNode {
	return s.GetToken(PMMPostParserFILL, 0)
}

func (s *FillCmdContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FillCmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FillCmdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FillCmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterFillCmd(s)
	}
}

func (s *FillCmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitFillCmd(s)
	}
}

func (p *PMMPostParser) FillCmd() (localctx IFillCmdContext) {
	localctx = NewFillCmdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, PMMPostParserRULE_fillCmd)

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
		p.SetState(161)
		p.Match(PMMPostParserFILL)
	}
	{
		p.SetState(162)
		p.expression(0)
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
	p.RuleIndex = PMMPostParserRULE_pickupCmd
	return p
}

func (*PickupCmdContext) IsPickupCmdContext() {}

func NewPickupCmdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PickupCmdContext {
	var p = new(PickupCmdContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_pickupCmd

	return p
}

func (s *PickupCmdContext) GetParser() antlr.Parser { return s.parser }

func (s *PickupCmdContext) PICKUP() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPICKUP, 0)
}

func (s *PickupCmdContext) PEN() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPEN, 0)
}

func (s *PickupCmdContext) DECIMALTOKEN() antlr.TerminalNode {
	return s.GetToken(PMMPostParserDECIMALTOKEN, 0)
}

func (s *PickupCmdContext) WITHCOLOR() antlr.TerminalNode {
	return s.GetToken(PMMPostParserWITHCOLOR, 0)
}

func (s *PickupCmdContext) COLOR() antlr.TerminalNode {
	return s.GetToken(PMMPostParserCOLOR, 0)
}

func (s *PickupCmdContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PickupCmdContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PickupCmdContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterPickupCmd(s)
	}
}

func (s *PickupCmdContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitPickupCmd(s)
	}
}

func (p *PMMPostParser) PickupCmd() (localctx IPickupCmdContext) {
	localctx = NewPickupCmdContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, PMMPostParserRULE_pickupCmd)
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
		p.SetState(164)
		p.Match(PMMPostParserPICKUP)
	}
	{
		p.SetState(165)
		p.Match(PMMPostParserPEN)
	}
	p.SetState(168)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PMMPostParserT__1 {
		{
			p.SetState(166)
			p.Match(PMMPostParserT__1)
		}
		{
			p.SetState(167)
			p.Match(PMMPostParserDECIMALTOKEN)
		}

	}
	p.SetState(172)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PMMPostParserWITHCOLOR {
		{
			p.SetState(170)
			p.Match(PMMPostParserWITHCOLOR)
		}
		{
			p.SetState(171)
			p.Match(PMMPostParserCOLOR)
		}

	}

	return localctx
}

// IPathjoinContext is an interface to support dynamic dispatch.
type IPathjoinContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsPathjoinContext differentiates from other interfaces.
	IsPathjoinContext()
}

type PathjoinContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPathjoinContext() *PathjoinContext {
	var p = new(PathjoinContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPostParserRULE_pathjoin
	return p
}

func (*PathjoinContext) IsPathjoinContext() {}

func NewPathjoinContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathjoinContext {
	var p = new(PathjoinContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_pathjoin

	return p
}

func (s *PathjoinContext) GetParser() antlr.Parser { return s.parser }

func (s *PathjoinContext) Basicpathjoin() IBasicpathjoinContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBasicpathjoinContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IBasicpathjoinContext)
}

func (s *PathjoinContext) AllDirectionspec() []IDirectionspecContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IDirectionspecContext)(nil)).Elem())
	var tst = make([]IDirectionspecContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IDirectionspecContext)
		}
	}

	return tst
}

func (s *PathjoinContext) Directionspec(i int) IDirectionspecContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDirectionspecContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IDirectionspecContext)
}

func (s *PathjoinContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PathjoinContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PathjoinContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterPathjoin(s)
	}
}

func (s *PathjoinContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitPathjoin(s)
	}
}

func (p *PMMPostParser) Pathjoin() (localctx IPathjoinContext) {
	localctx = NewPathjoinContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, PMMPostParserRULE_pathjoin)
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
	p.SetState(175)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PMMPostParserT__2 {
		{
			p.SetState(174)
			p.Directionspec()
		}

	}
	{
		p.SetState(177)
		p.Basicpathjoin()
	}
	p.SetState(179)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PMMPostParserT__2 {
		{
			p.SetState(178)
			p.Directionspec()
		}

	}

	return localctx
}

// IDirectionspecContext is an interface to support dynamic dispatch.
type IDirectionspecContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDirectionspecContext differentiates from other interfaces.
	IsDirectionspecContext()
}

type DirectionspecContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDirectionspecContext() *DirectionspecContext {
	var p = new(DirectionspecContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPostParserRULE_directionspec
	return p
}

func (*DirectionspecContext) IsDirectionspecContext() {}

func NewDirectionspecContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DirectionspecContext {
	var p = new(DirectionspecContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_directionspec

	return p
}

func (s *DirectionspecContext) GetParser() antlr.Parser { return s.parser }

func (s *DirectionspecContext) CopyFrom(ctx *DirectionspecContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *DirectionspecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DirectionspecContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type DirspecContext struct {
	*DirectionspecContext
}

func NewDirspecContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DirspecContext {
	var p = new(DirspecContext)

	p.DirectionspecContext = NewEmptyDirectionspecContext()
	p.parser = parser
	p.CopyFrom(ctx.(*DirectionspecContext))

	return p
}

func (s *DirspecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DirspecContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *DirspecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterDirspec(s)
	}
}

func (s *DirspecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitDirspec(s)
	}
}

type CurspecContext struct {
	*DirectionspecContext
}

func NewCurspecContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CurspecContext {
	var p = new(CurspecContext)

	p.DirectionspecContext = NewEmptyDirectionspecContext()
	p.parser = parser
	p.CopyFrom(ctx.(*DirectionspecContext))

	return p
}

func (s *CurspecContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CurspecContext) CURL() antlr.TerminalNode {
	return s.GetToken(PMMPostParserCURL, 0)
}

func (s *CurspecContext) Expression() IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *CurspecContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterCurspec(s)
	}
}

func (s *CurspecContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitCurspec(s)
	}
}

func (p *PMMPostParser) Directionspec() (localctx IDirectionspecContext) {
	localctx = NewDirectionspecContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, PMMPostParserRULE_directionspec)

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

	p.SetState(190)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
	case 1:
		localctx = NewCurspecContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(181)
			p.Match(PMMPostParserT__2)
		}
		{
			p.SetState(182)
			p.Match(PMMPostParserCURL)
		}
		{
			p.SetState(183)
			p.expression(0)
		}
		{
			p.SetState(184)
			p.Match(PMMPostParserT__3)
		}

	case 2:
		localctx = NewDirspecContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(186)
			p.Match(PMMPostParserT__2)
		}
		{
			p.SetState(187)
			p.expression(0)
		}
		{
			p.SetState(188)
			p.Match(PMMPostParserT__3)
		}

	}

	return localctx
}

// IBasicpathjoinContext is an interface to support dynamic dispatch.
type IBasicpathjoinContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBasicpathjoinContext differentiates from other interfaces.
	IsBasicpathjoinContext()
}

type BasicpathjoinContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBasicpathjoinContext() *BasicpathjoinContext {
	var p = new(BasicpathjoinContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPostParserRULE_basicpathjoin
	return p
}

func (*BasicpathjoinContext) IsBasicpathjoinContext() {}

func NewBasicpathjoinContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BasicpathjoinContext {
	var p = new(BasicpathjoinContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_basicpathjoin

	return p
}

func (s *BasicpathjoinContext) GetParser() antlr.Parser { return s.parser }

func (s *BasicpathjoinContext) PATHJOIN() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPATHJOIN, 0)
}

func (s *BasicpathjoinContext) Controls() IControlsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IControlsContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IControlsContext)
}

func (s *BasicpathjoinContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BasicpathjoinContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BasicpathjoinContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterBasicpathjoin(s)
	}
}

func (s *BasicpathjoinContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitBasicpathjoin(s)
	}
}

func (p *PMMPostParser) Basicpathjoin() (localctx IBasicpathjoinContext) {
	localctx = NewBasicpathjoinContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, PMMPostParserRULE_basicpathjoin)

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

	p.SetState(197)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPostParserPATHJOIN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(192)
			p.Match(PMMPostParserPATHJOIN)
		}

	case PMMPostParserT__4:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(193)
			p.Match(PMMPostParserT__4)
		}
		{
			p.SetState(194)
			p.Controls()
		}
		{
			p.SetState(195)
			p.Match(PMMPostParserT__4)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IControlsContext is an interface to support dynamic dispatch.
type IControlsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsControlsContext differentiates from other interfaces.
	IsControlsContext()
}

type ControlsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyControlsContext() *ControlsContext {
	var p = new(ControlsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPostParserRULE_controls
	return p
}

func (*ControlsContext) IsControlsContext() {}

func NewControlsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ControlsContext {
	var p = new(ControlsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_controls

	return p
}

func (s *ControlsContext) GetParser() antlr.Parser { return s.parser }

func (s *ControlsContext) CONTROLS() antlr.TerminalNode {
	return s.GetToken(PMMPostParserCONTROLS, 0)
}

func (s *ControlsContext) AllExpression() []IExpressionContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IExpressionContext)(nil)).Elem())
	var tst = make([]IExpressionContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IExpressionContext)
		}
	}

	return tst
}

func (s *ControlsContext) Expression(i int) IExpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ControlsContext) AND() antlr.TerminalNode {
	return s.GetToken(PMMPostParserAND, 0)
}

func (s *ControlsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ControlsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ControlsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterControls(s)
	}
}

func (s *ControlsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitControls(s)
	}
}

func (p *PMMPostParser) Controls() (localctx IControlsContext) {
	localctx = NewControlsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, PMMPostParserRULE_controls)
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
		p.SetState(199)
		p.Match(PMMPostParserCONTROLS)
	}
	{
		p.SetState(200)
		p.expression(0)
	}
	p.SetState(203)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PMMPostParserAND {
		{
			p.SetState(201)
			p.Match(PMMPostParserAND)
		}
		{
			p.SetState(202)
			p.expression(0)
		}

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
	p.RuleIndex = PMMPostParserRULE_statementlist
	return p
}

func (*StatementlistContext) IsStatementlistContext() {}

func NewStatementlistContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementlistContext {
	var p = new(StatementlistContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_statementlist

	return p
}

func (s *StatementlistContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementlistContext) EOF() antlr.TerminalNode {
	return s.GetToken(PMMPostParserEOF, 0)
}

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

func (s *StatementlistContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementlistContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementlistContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterStatementlist(s)
	}
}

func (s *StatementlistContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitStatementlist(s)
	}
}

func (p *PMMPostParser) Statementlist() (localctx IStatementlistContext) {
	localctx = NewStatementlistContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, PMMPostParserRULE_statementlist)
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
	p.SetState(206)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<PMMPostParserT__0)|(1<<PMMPostParserTYPE)|(1<<PMMPostParserENDFIG)|(1<<PMMPostParserPICKUP)|(1<<PMMPostParserFILL)|(1<<PMMPostParserDRAW)|(1<<PMMPostParserSEMIC)|(1<<PMMPostParserLPAREN)|(1<<PMMPostParserPLUS)|(1<<PMMPostParserMINUS))) != 0) || (((_la-36)&-(0x1f+1)) == 0 && ((1<<uint((_la-36)))&((1<<(PMMPostParserLAMBDAARG-36))|(1<<(PMMPostParserBEGINGROUP-36))|(1<<(PMMPostParserLOCAL-36))|(1<<(PMMPostParserPAIRPART-36))|(1<<(PMMPostParserMATHFUNC-36))|(1<<(PMMPostParserSUBPATH-36))|(1<<(PMMPostParserREVERSE-36))|(1<<(PMMPostParserPOINT-36))|(1<<(PMMPostParserPROOF-36))|(1<<(PMMPostParserSAVE-36))|(1<<(PMMPostParserSHOW-36))|(1<<(PMMPostParserLET-36))|(1<<(PMMPostParserTAG-36))|(1<<(PMMPostParserMIXEDTAG-36))|(1<<(PMMPostParserDECIMALTOKEN-36)))) != 0) {
		{
			p.SetState(205)
			p.Statement()
		}

		p.SetState(208)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(210)
		p.Match(PMMPostParserEOF)
	}

	return localctx
}

// IVardefContext is an interface to support dynamic dispatch.
type IVardefContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsVardefContext differentiates from other interfaces.
	IsVardefContext()
}

type VardefContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVardefContext() *VardefContext {
	var p = new(VardefContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPostParserRULE_vardef
	return p
}

func (*VardefContext) IsVardefContext() {}

func NewVardefContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VardefContext {
	var p = new(VardefContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_vardef

	return p
}

func (s *VardefContext) GetParser() antlr.Parser { return s.parser }

func (s *VardefContext) VARDEF() antlr.TerminalNode {
	return s.GetToken(PMMPostParserVARDEF, 0)
}

func (s *VardefContext) AllTAG() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserTAG)
}

func (s *VardefContext) TAG(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserTAG, i)
}

func (s *VardefContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserCOMMA)
}

func (s *VardefContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserCOMMA, i)
}

func (s *VardefContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VardefContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VardefContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterVardef(s)
	}
}

func (s *VardefContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitVardef(s)
	}
}

func (p *PMMPostParser) Vardef() (localctx IVardefContext) {
	localctx = NewVardefContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, PMMPostParserRULE_vardef)
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
		p.SetState(212)
		p.Match(PMMPostParserVARDEF)
	}
	{
		p.SetState(213)
		p.Match(PMMPostParserTAG)
	}
	p.SetState(218)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == PMMPostParserCOMMA {
		{
			p.SetState(214)
			p.Match(PMMPostParserCOMMA)
		}
		{
			p.SetState(215)
			p.Match(PMMPostParserTAG)
		}

		p.SetState(220)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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
	p.RuleIndex = PMMPostParserRULE_compound
	return p
}

func (*CompoundContext) IsCompoundContext() {}

func NewCompoundContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompoundContext {
	var p = new(CompoundContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_compound

	return p
}

func (s *CompoundContext) GetParser() antlr.Parser { return s.parser }

func (s *CompoundContext) BEGINGROUP() antlr.TerminalNode {
	return s.GetToken(PMMPostParserBEGINGROUP, 0)
}

func (s *CompoundContext) Statementlist() IStatementlistContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStatementlistContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStatementlistContext)
}

func (s *CompoundContext) ENDGROUP() antlr.TerminalNode {
	return s.GetToken(PMMPostParserENDGROUP, 0)
}

func (s *CompoundContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompoundContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompoundContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterCompound(s)
	}
}

func (s *CompoundContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitCompound(s)
	}
}

func (p *PMMPostParser) Compound() (localctx ICompoundContext) {
	localctx = NewCompoundContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, PMMPostParserRULE_compound)

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
		p.SetState(221)
		p.Match(PMMPostParserBEGINGROUP)
	}
	{
		p.SetState(222)
		p.Statementlist()
	}
	{
		p.SetState(223)
		p.Match(PMMPostParserENDGROUP)
	}

	return localctx
}

// IEmptyContext is an interface to support dynamic dispatch.
type IEmptyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEmptyContext differentiates from other interfaces.
	IsEmptyContext()
}

type EmptyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEmptyContext() *EmptyContext {
	var p = new(EmptyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPostParserRULE_empty
	return p
}

func (*EmptyContext) IsEmptyContext() {}

func NewEmptyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EmptyContext {
	var p = new(EmptyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_empty

	return p
}

func (s *EmptyContext) GetParser() antlr.Parser { return s.parser }
func (s *EmptyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EmptyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EmptyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterEmpty(s)
	}
}

func (s *EmptyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitEmpty(s)
	}
}

func (p *PMMPostParser) Empty() (localctx IEmptyContext) {
	localctx = NewEmptyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, PMMPostParserRULE_empty)

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
	p.RuleIndex = PMMPostParserRULE_assignment
	return p
}

func (*AssignmentContext) IsAssignmentContext() {}

func NewAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentContext {
	var p = new(AssignmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_assignment

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
	return s.GetToken(PMMPostParserASSIGN, 0)
}

func (s *AssignmentContext) Path() IPathContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathContext)
}

func (s *AssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterAssignment(s)
	}
}

func (s *AssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitAssignment(s)
	}
}

func (p *PMMPostParser) Assignment() (localctx IAssignmentContext) {
	localctx = NewAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, PMMPostParserRULE_assignment)

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
		p.SetState(227)
		p.Variable()
	}
	{
		p.SetState(228)
		p.Match(PMMPostParserASSIGN)
	}
	{
		p.SetState(229)
		p.Path()
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
	p.RuleIndex = PMMPostParserRULE_constraint
	return p
}

func (*ConstraintContext) IsConstraintContext() {}

func NewConstraintContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConstraintContext {
	var p = new(ConstraintContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_constraint

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
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterConstraint(s)
	}
}

func (s *ConstraintContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitConstraint(s)
	}
}

func (p *PMMPostParser) Constraint() (localctx IConstraintContext) {
	localctx = NewConstraintContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, PMMPostParserRULE_constraint)

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

	p.SetState(233)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(231)
			p.Equation()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(232)
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
	p.RuleIndex = PMMPostParserRULE_equation
	return p
}

func (*EquationContext) IsEquationContext() {}

func NewEquationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EquationContext {
	var p = new(EquationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_equation

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
	return s.GetTokens(PMMPostParserEQUALS)
}

func (s *EquationContext) EQUALS(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserEQUALS, i)
}

func (s *EquationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EquationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EquationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterEquation(s)
	}
}

func (s *EquationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitEquation(s)
	}
}

func (p *PMMPostParser) Equation() (localctx IEquationContext) {
	localctx = NewEquationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, PMMPostParserRULE_equation)
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
		p.SetState(235)
		p.expression(0)
	}
	p.SetState(238)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == PMMPostParserEQUALS {
		{
			p.SetState(236)
			p.Match(PMMPostParserEQUALS)
		}
		{
			p.SetState(237)
			p.expression(0)
		}

		p.SetState(240)
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
	p.RuleIndex = PMMPostParserRULE_orientation
	return p
}

func (*OrientationContext) IsOrientationContext() {}

func NewOrientationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrientationContext {
	var p = new(OrientationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_orientation

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
	return s.GetTokens(PMMPostParserPARALLEL)
}

func (s *OrientationContext) PARALLEL(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserPARALLEL, i)
}

func (s *OrientationContext) AllPERPENDIC() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserPERPENDIC)
}

func (s *OrientationContext) PERPENDIC(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserPERPENDIC, i)
}

func (s *OrientationContext) AllCONGRUENT() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserCONGRUENT)
}

func (s *OrientationContext) CONGRUENT(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserCONGRUENT, i)
}

func (s *OrientationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrientationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OrientationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterOrientation(s)
	}
}

func (s *OrientationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitOrientation(s)
	}
}

func (p *PMMPostParser) Orientation() (localctx IOrientationContext) {
	localctx = NewOrientationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, PMMPostParserRULE_orientation)
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
		p.SetState(242)
		p.tertiary(0)
	}
	p.SetState(245)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(PMMPostParserPARALLEL-32))|(1<<(PMMPostParserPERPENDIC-32))|(1<<(PMMPostParserCONGRUENT-32)))) != 0) {
		{
			p.SetState(243)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(PMMPostParserPARALLEL-32))|(1<<(PMMPostParserPERPENDIC-32))|(1<<(PMMPostParserCONGRUENT-32)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(244)
			p.tertiary(0)
		}

		p.SetState(247)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ITokenContext is an interface to support dynamic dispatch.
type ITokenContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTokenContext differentiates from other interfaces.
	IsTokenContext()
}

type TokenContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTokenContext() *TokenContext {
	var p = new(TokenContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = PMMPostParserRULE_token
	return p
}

func (*TokenContext) IsTokenContext() {}

func NewTokenContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TokenContext {
	var p = new(TokenContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_token

	return p
}

func (s *TokenContext) GetParser() antlr.Parser { return s.parser }

func (s *TokenContext) PLUS() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPLUS, 0)
}

func (s *TokenContext) MINUS() antlr.TerminalNode {
	return s.GetToken(PMMPostParserMINUS, 0)
}

func (s *TokenContext) TIMES() antlr.TerminalNode {
	return s.GetToken(PMMPostParserTIMES, 0)
}

func (s *TokenContext) OVER() antlr.TerminalNode {
	return s.GetToken(PMMPostParserOVER, 0)
}

func (s *TokenContext) ASSIGN() antlr.TerminalNode {
	return s.GetToken(PMMPostParserASSIGN, 0)
}

func (s *TokenContext) PARALLEL() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPARALLEL, 0)
}

func (s *TokenContext) PERPENDIC() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPERPENDIC, 0)
}

func (s *TokenContext) CONGRUENT() antlr.TerminalNode {
	return s.GetToken(PMMPostParserCONGRUENT, 0)
}

func (s *TokenContext) BEGINGROUP() antlr.TerminalNode {
	return s.GetToken(PMMPostParserBEGINGROUP, 0)
}

func (s *TokenContext) ENDGROUP() antlr.TerminalNode {
	return s.GetToken(PMMPostParserENDGROUP, 0)
}

func (s *TokenContext) EDGECONSTR() antlr.TerminalNode {
	return s.GetToken(PMMPostParserEDGECONSTR, 0)
}

func (s *TokenContext) PATHCLIPOP() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPATHCLIPOP, 0)
}

func (s *TokenContext) PATHJOIN() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPATHJOIN, 0)
}

func (s *TokenContext) EDGE() antlr.TerminalNode {
	return s.GetToken(PMMPostParserEDGE, 0)
}

func (s *TokenContext) FRAME() antlr.TerminalNode {
	return s.GetToken(PMMPostParserFRAME, 0)
}

func (s *TokenContext) BOX() antlr.TerminalNode {
	return s.GetToken(PMMPostParserBOX, 0)
}

func (s *TokenContext) REVERSE() antlr.TerminalNode {
	return s.GetToken(PMMPostParserREVERSE, 0)
}

func (s *TokenContext) SUBPATH() antlr.TerminalNode {
	return s.GetToken(PMMPostParserSUBPATH, 0)
}

func (s *TokenContext) PROOF() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPROOF, 0)
}

func (s *TokenContext) SAVE() antlr.TerminalNode {
	return s.GetToken(PMMPostParserSAVE, 0)
}

func (s *TokenContext) SHOW() antlr.TerminalNode {
	return s.GetToken(PMMPostParserSHOW, 0)
}

func (s *TokenContext) TRANSFORM() antlr.TerminalNode {
	return s.GetToken(PMMPostParserTRANSFORM, 0)
}

func (s *TokenContext) TAG() antlr.TerminalNode {
	return s.GetToken(PMMPostParserTAG, 0)
}

func (s *TokenContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TokenContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TokenContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterToken(s)
	}
}

func (s *TokenContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitToken(s)
	}
}

func (p *PMMPostParser) Token() (localctx ITokenContext) {
	localctx = NewTokenContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, PMMPostParserRULE_token)
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
		p.SetState(249)
		_la = p.GetTokenStream().LA(1)

		if !((((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<PMMPostParserPATHJOIN)|(1<<PMMPostParserASSIGN)|(1<<PMMPostParserPLUS)|(1<<PMMPostParserMINUS)|(1<<PMMPostParserTIMES)|(1<<PMMPostParserOVER))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(PMMPostParserPARALLEL-32))|(1<<(PMMPostParserPERPENDIC-32))|(1<<(PMMPostParserCONGRUENT-32))|(1<<(PMMPostParserBEGINGROUP-32))|(1<<(PMMPostParserENDGROUP-32))|(1<<(PMMPostParserEDGECONSTR-32))|(1<<(PMMPostParserEDGE-32))|(1<<(PMMPostParserFRAME-32))|(1<<(PMMPostParserBOX-32))|(1<<(PMMPostParserSUBPATH-32))|(1<<(PMMPostParserREVERSE-32))|(1<<(PMMPostParserTRANSFORM-32))|(1<<(PMMPostParserPATHCLIPOP-32))|(1<<(PMMPostParserPROOF-32))|(1<<(PMMPostParserSAVE-32))|(1<<(PMMPostParserSHOW-32))|(1<<(PMMPostParserTAG-32)))) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
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
	p.RuleIndex = PMMPostParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_expression

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
	return s.GetToken(PMMPostParserPATHCLIPOP, 0)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *PMMPostParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *PMMPostParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 42
	p.EnterRecursionRule(localctx, 42, PMMPostParserRULE_expression, _p)

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
		p.SetState(252)
		p.tertiary(0)
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(259)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewExpressionContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, PMMPostParserRULE_expression)
			p.SetState(254)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(255)
				p.Match(PMMPostParserPATHCLIPOP)
			}
			{
				p.SetState(256)
				p.tertiary(0)
			}

		}
		p.SetState(261)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext())
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
	p.RuleIndex = PMMPostParserRULE_tertiary
	return p
}

func (*TertiaryContext) IsTertiaryContext() {}

func NewTertiaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TertiaryContext {
	var p = new(TertiaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_tertiary

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
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterPathtertiary(s)
	}
}

func (s *PathtertiaryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitPathtertiary(s)
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

func (s *TermContext) Secondary() ISecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISecondaryContext)
}

func (s *TermContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *TermContext) PLUS() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPLUS, 0)
}

func (s *TermContext) MINUS() antlr.TerminalNode {
	return s.GetToken(PMMPostParserMINUS, 0)
}

func (s *TermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterTerm(s)
	}
}

func (s *TermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitTerm(s)
	}
}

func (p *PMMPostParser) Tertiary() (localctx ITertiaryContext) {
	return p.tertiary(0)
}

func (p *PMMPostParser) tertiary(_p int) (localctx ITertiaryContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewTertiaryContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ITertiaryContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 44
	p.EnterRecursionRule(localctx, 44, PMMPostParserRULE_tertiary, _p)
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
	p.SetState(265)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		localctx = NewPathtertiaryContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(263)
			p.Path()
		}

	case 2:
		localctx = NewTermContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(264)
			p.secondary(0)
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(272)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewTermContext(p, NewTertiaryContext(p, _parentctx, _parentState))
			p.PushNewRecursionContext(localctx, _startState, PMMPostParserRULE_tertiary)
			p.SetState(267)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(268)
				_la = p.GetTokenStream().LA(1)

				if !(_la == PMMPostParserPLUS || _la == PMMPostParserMINUS) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(269)
				p.secondary(0)
			}

		}
		p.SetState(274)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 22, p.GetParserRuleContext())
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
	p.RuleIndex = PMMPostParserRULE_path
	return p
}

func (*PathContext) IsPathContext() {}

func NewPathContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PathContext {
	var p = new(PathContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_path

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
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterPath(s)
	}
}

func (s *PathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitPath(s)
	}
}

func (p *PMMPostParser) Path() (localctx IPathContext) {
	localctx = NewPathContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, PMMPostParserRULE_path)

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
		p.SetState(275)
		p.secondary(0)
	}
	p.SetState(278)
	p.GetErrorHandler().Sync(p)
	_alt = 1
	for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		switch _alt {
		case 1:
			{
				p.SetState(276)
				p.Match(PMMPostParserT__4)
			}
			{
				p.SetState(277)
				p.secondary(0)
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(280)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 23, p.GetParserRuleContext())
	}
	p.SetState(283)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(282)
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
	p.RuleIndex = PMMPostParserRULE_cycle
	return p
}

func (*CycleContext) IsCycleContext() {}

func NewCycleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CycleContext {
	var p = new(CycleContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_cycle

	return p
}

func (s *CycleContext) GetParser() antlr.Parser { return s.parser }

func (s *CycleContext) Pathjoin() IPathjoinContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPathjoinContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPathjoinContext)
}

func (s *CycleContext) CYCLE() antlr.TerminalNode {
	return s.GetToken(PMMPostParserCYCLE, 0)
}

func (s *CycleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CycleContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CycleContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterCycle(s)
	}
}

func (s *CycleContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitCycle(s)
	}
}

func (p *PMMPostParser) Cycle() (localctx ICycleContext) {
	localctx = NewCycleContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, PMMPostParserRULE_cycle)

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
		p.SetState(285)
		p.Pathjoin()
	}
	{
		p.SetState(286)
		p.Match(PMMPostParserCYCLE)
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
	p.RuleIndex = PMMPostParserRULE_secondary
	return p
}

func (*SecondaryContext) IsSecondaryContext() {}

func NewSecondaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SecondaryContext {
	var p = new(SecondaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_secondary

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

func (s *TransformContext) AllTRANSFORM() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserTRANSFORM)
}

func (s *TransformContext) TRANSFORM(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserTRANSFORM, i)
}

func (s *TransformContext) AllPrimary() []IPrimaryContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IPrimaryContext)(nil)).Elem())
	var tst = make([]IPrimaryContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IPrimaryContext)
		}
	}

	return tst
}

func (s *TransformContext) Primary(i int) IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *TransformContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterTransform(s)
	}
}

func (s *TransformContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitTransform(s)
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

func (s *FactorContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *FactorContext) Secondary() ISecondaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISecondaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISecondaryContext)
}

func (s *FactorContext) TIMES() antlr.TerminalNode {
	return s.GetToken(PMMPostParserTIMES, 0)
}

func (s *FactorContext) OVER() antlr.TerminalNode {
	return s.GetToken(PMMPostParserOVER, 0)
}

func (s *FactorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterFactor(s)
	}
}

func (s *FactorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitFactor(s)
	}
}

func (p *PMMPostParser) Secondary() (localctx ISecondaryContext) {
	return p.secondary(0)
}

func (p *PMMPostParser) secondary(_p int) (localctx ISecondaryContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewSecondaryContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx ISecondaryContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 50
	p.EnterRecursionRule(localctx, 50, PMMPostParserRULE_secondary, _p)
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
	localctx = NewFactorContext(p, localctx)
	p.SetParserRuleContext(localctx)
	_prevctx = localctx

	{
		p.SetState(289)
		p.Primary()
	}

	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(303)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(301)
			p.GetErrorHandler().Sync(p)
			switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) {
			case 1:
				localctx = NewFactorContext(p, NewSecondaryContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, PMMPostParserRULE_secondary)
				p.SetState(291)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
				}
				{
					p.SetState(292)
					_la = p.GetTokenStream().LA(1)

					if !(_la == PMMPostParserTIMES || _la == PMMPostParserOVER) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(293)
					p.Primary()
				}

			case 2:
				localctx = NewTransformContext(p, NewSecondaryContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, PMMPostParserRULE_secondary)
				p.SetState(294)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
				}
				p.SetState(297)
				p.GetErrorHandler().Sync(p)
				_alt = 1
				for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
					switch _alt {
					case 1:
						{
							p.SetState(295)
							p.Match(PMMPostParserTRANSFORM)
						}
						{
							p.SetState(296)
							p.Primary()
						}

					default:
						panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
					}

					p.SetState(299)
					p.GetErrorHandler().Sync(p)
					_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())
				}

			}

		}
		p.SetState(305)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 27, p.GetParserRuleContext())
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
	p.RuleIndex = PMMPostParserRULE_primary
	return p
}

func (*PrimaryContext) IsPrimaryContext() {}

func NewPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryContext {
	var p = new(PrimaryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_primary

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
	return s.GetToken(PMMPostParserLBRACKET, 0)
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
	return s.GetToken(PMMPostParserCOMMA, 0)
}

func (s *InterpolationContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(PMMPostParserRBRACKET, 0)
}

func (s *InterpolationContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *InterpolationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterInterpolation(s)
	}
}

func (s *InterpolationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitInterpolation(s)
	}
}

type SimpleatomContext struct {
	*PrimaryContext
}

func NewSimpleatomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimpleatomContext {
	var p = new(SimpleatomContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *SimpleatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleatomContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *SimpleatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterSimpleatom(s)
	}
}

func (s *SimpleatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitSimpleatom(s)
	}
}

type FuncatomContext struct {
	*PrimaryContext
}

func NewFuncatomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FuncatomContext {
	var p = new(FuncatomContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *FuncatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FuncatomContext) MATHFUNC() antlr.TerminalNode {
	return s.GetToken(PMMPostParserMATHFUNC, 0)
}

func (s *FuncatomContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *FuncatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterFuncatom(s)
	}
}

func (s *FuncatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitFuncatom(s)
	}
}

type PointofContext struct {
	*PrimaryContext
}

func NewPointofContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PointofContext {
	var p = new(PointofContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *PointofContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PointofContext) POINT() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPOINT, 0)
}

func (s *PointofContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *PointofContext) OF() antlr.TerminalNode {
	return s.GetToken(PMMPostParserOF, 0)
}

func (s *PointofContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *PointofContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterPointof(s)
	}
}

func (s *PointofContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitPointof(s)
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
	return s.GetToken(PMMPostParserSUBPATH, 0)
}

func (s *SubpathContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *SubpathContext) OF() antlr.TerminalNode {
	return s.GetToken(PMMPostParserOF, 0)
}

func (s *SubpathContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *SubpathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterSubpath(s)
	}
}

func (s *SubpathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
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
	return s.GetToken(PMMPostParserPAIRPART, 0)
}

func (s *PairpartContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *PairpartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterPairpart(s)
	}
}

func (s *PairpartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitPairpart(s)
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
	return s.GetToken(PMMPostParserREVERSE, 0)
}

func (s *ReversepathContext) Primary() IPrimaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IPrimaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IPrimaryContext)
}

func (s *ReversepathContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterReversepath(s)
	}
}

func (s *ReversepathContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitReversepath(s)
	}
}

type ScalaratomContext struct {
	*PrimaryContext
}

func NewScalaratomContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ScalaratomContext {
	var p = new(ScalaratomContext)

	p.PrimaryContext = NewEmptyPrimaryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*PrimaryContext))

	return p
}

func (s *ScalaratomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScalaratomContext) Scalarmulop() IScalarmulopContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IScalarmulopContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IScalarmulopContext)
}

func (s *ScalaratomContext) Atom() IAtomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAtomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAtomContext)
}

func (s *ScalaratomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterScalaratom(s)
	}
}

func (s *ScalaratomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitScalaratom(s)
	}
}

func (p *PMMPostParser) Primary() (localctx IPrimaryContext) {
	localctx = NewPrimaryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, PMMPostParserRULE_primary)

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

	p.SetState(340)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 28, p.GetParserRuleContext()) {
	case 1:
		localctx = NewFuncatomContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(306)
			p.Match(PMMPostParserMATHFUNC)
		}
		{
			p.SetState(307)
			p.Atom()
		}

	case 2:
		localctx = NewScalaratomContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(308)
			p.Scalarmulop()
		}
		{
			p.SetState(309)
			p.Atom()
		}

	case 3:
		localctx = NewInterpolationContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(311)
			p.Numtokenatom()
		}
		{
			p.SetState(312)
			p.Match(PMMPostParserLBRACKET)
		}
		{
			p.SetState(313)
			p.tertiary(0)
		}
		{
			p.SetState(314)
			p.Match(PMMPostParserCOMMA)
		}
		{
			p.SetState(315)
			p.tertiary(0)
		}
		{
			p.SetState(316)
			p.Match(PMMPostParserRBRACKET)
		}

	case 4:
		localctx = NewInterpolationContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(318)
			p.Atom()
		}
		{
			p.SetState(319)
			p.Match(PMMPostParserLBRACKET)
		}
		{
			p.SetState(320)
			p.tertiary(0)
		}
		{
			p.SetState(321)
			p.Match(PMMPostParserCOMMA)
		}
		{
			p.SetState(322)
			p.tertiary(0)
		}
		{
			p.SetState(323)
			p.Match(PMMPostParserRBRACKET)
		}

	case 5:
		localctx = NewSimpleatomContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(325)
			p.Atom()
		}

	case 6:
		localctx = NewPairpartContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(326)
			p.Match(PMMPostParserPAIRPART)
		}
		{
			p.SetState(327)
			p.Primary()
		}

	case 7:
		localctx = NewPointofContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(328)
			p.Match(PMMPostParserPOINT)
		}
		{
			p.SetState(329)
			p.tertiary(0)
		}
		{
			p.SetState(330)
			p.Match(PMMPostParserOF)
		}
		{
			p.SetState(331)
			p.Primary()
		}

	case 8:
		localctx = NewReversepathContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(333)
			p.Match(PMMPostParserREVERSE)
		}
		{
			p.SetState(334)
			p.Primary()
		}

	case 9:
		localctx = NewSubpathContext(p, localctx)
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(335)
			p.Match(PMMPostParserSUBPATH)
		}
		{
			p.SetState(336)
			p.tertiary(0)
		}
		{
			p.SetState(337)
			p.Match(PMMPostParserOF)
		}
		{
			p.SetState(338)
			p.Primary()
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
	p.RuleIndex = PMMPostParserRULE_scalarmulop
	return p
}

func (*ScalarmulopContext) IsScalarmulopContext() {}

func NewScalarmulopContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ScalarmulopContext {
	var p = new(ScalarmulopContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_scalarmulop

	return p
}

func (s *ScalarmulopContext) GetParser() antlr.Parser { return s.parser }

func (s *ScalarmulopContext) Numtokenatom() INumtokenatomContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumtokenatomContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumtokenatomContext)
}

func (s *ScalarmulopContext) PLUS() antlr.TerminalNode {
	return s.GetToken(PMMPostParserPLUS, 0)
}

func (s *ScalarmulopContext) MINUS() antlr.TerminalNode {
	return s.GetToken(PMMPostParserMINUS, 0)
}

func (s *ScalarmulopContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ScalarmulopContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ScalarmulopContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterScalarmulop(s)
	}
}

func (s *ScalarmulopContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitScalarmulop(s)
	}
}

func (p *PMMPostParser) Scalarmulop() (localctx IScalarmulopContext) {
	localctx = NewScalarmulopContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, PMMPostParserRULE_scalarmulop)
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
	p.SetState(343)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == PMMPostParserPLUS || _la == PMMPostParserMINUS {
		{
			p.SetState(342)
			_la = p.GetTokenStream().LA(1)

			if !(_la == PMMPostParserPLUS || _la == PMMPostParserMINUS) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(345)
		p.Numtokenatom()
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
	p.RuleIndex = PMMPostParserRULE_numtokenatom
	return p
}

func (*NumtokenatomContext) IsNumtokenatomContext() {}

func NewNumtokenatomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumtokenatomContext {
	var p = new(NumtokenatomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_numtokenatom

	return p
}

func (s *NumtokenatomContext) GetParser() antlr.Parser { return s.parser }

func (s *NumtokenatomContext) AllDECIMALTOKEN() []antlr.TerminalNode {
	return s.GetTokens(PMMPostParserDECIMALTOKEN)
}

func (s *NumtokenatomContext) DECIMALTOKEN(i int) antlr.TerminalNode {
	return s.GetToken(PMMPostParserDECIMALTOKEN, i)
}

func (s *NumtokenatomContext) OVER() antlr.TerminalNode {
	return s.GetToken(PMMPostParserOVER, 0)
}

func (s *NumtokenatomContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumtokenatomContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NumtokenatomContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterNumtokenatom(s)
	}
}

func (s *NumtokenatomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitNumtokenatom(s)
	}
}

func (p *PMMPostParser) Numtokenatom() (localctx INumtokenatomContext) {
	localctx = NewNumtokenatomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, PMMPostParserRULE_numtokenatom)

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

	p.SetState(351)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 30, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(347)
			p.Match(PMMPostParserDECIMALTOKEN)
		}
		{
			p.SetState(348)
			p.Match(PMMPostParserOVER)
		}
		{
			p.SetState(349)
			p.Match(PMMPostParserDECIMALTOKEN)
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(350)
			p.Match(PMMPostParserDECIMALTOKEN)
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
	p.RuleIndex = PMMPostParserRULE_atom
	return p
}

func (*AtomContext) IsAtomContext() {}

func NewAtomContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AtomContext {
	var p = new(AtomContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_atom

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
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterVaratom(s)
	}
}

func (s *VaratomContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
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
	return s.GetToken(PMMPostParserBEGINGROUP, 0)
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
	return s.GetToken(PMMPostParserENDGROUP, 0)
}

func (s *ExprgroupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterExprgroup(s)
	}
}

func (s *ExprgroupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
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
	return s.GetToken(PMMPostParserDECIMALTOKEN, 0)
}

func (s *DecimalContext) UNIT() antlr.TerminalNode {
	return s.GetToken(PMMPostParserUNIT, 0)
}

func (s *DecimalContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterDecimal(s)
	}
}

func (s *DecimalContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
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
	return s.GetToken(PMMPostParserLPAREN, 0)
}

func (s *SubexpressionContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *SubexpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PMMPostParserRPAREN, 0)
}

func (s *SubexpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterSubexpression(s)
	}
}

func (s *SubexpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
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
	return s.GetToken(PMMPostParserLPAREN, 0)
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
	return s.GetToken(PMMPostParserCOMMA, 0)
}

func (s *LiteralpairContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(PMMPostParserRPAREN, 0)
}

func (s *LiteralpairContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterLiteralpair(s)
	}
}

func (s *LiteralpairContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitLiteralpair(s)
	}
}

func (p *PMMPostParser) Atom() (localctx IAtomContext) {
	localctx = NewAtomContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, PMMPostParserRULE_atom)

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

	p.SetState(373)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 32, p.GetParserRuleContext()) {
	case 1:
		localctx = NewDecimalContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(353)
			p.Match(PMMPostParserDECIMALTOKEN)
		}
		p.SetState(355)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 31, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(354)
				p.Match(PMMPostParserUNIT)
			}

		}

	case 2:
		localctx = NewVaratomContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(357)
			p.Variable()
		}

	case 3:
		localctx = NewLiteralpairContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(358)
			p.Match(PMMPostParserLPAREN)
		}
		{
			p.SetState(359)
			p.tertiary(0)
		}
		{
			p.SetState(360)
			p.Match(PMMPostParserCOMMA)
		}
		{
			p.SetState(361)
			p.tertiary(0)
		}
		{
			p.SetState(362)
			p.Match(PMMPostParserRPAREN)
		}

	case 4:
		localctx = NewSubexpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(364)
			p.Match(PMMPostParserLPAREN)
		}
		{
			p.SetState(365)
			p.tertiary(0)
		}
		{
			p.SetState(366)
			p.Match(PMMPostParserRPAREN)
		}

	case 5:
		localctx = NewExprgroupContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(368)
			p.Match(PMMPostParserBEGINGROUP)
		}
		{
			p.SetState(369)
			p.Statementlist()
		}
		{
			p.SetState(370)
			p.tertiary(0)
		}
		{
			p.SetState(371)
			p.Match(PMMPostParserENDGROUP)
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
	p.RuleIndex = PMMPostParserRULE_variable
	return p
}

func (*VariableContext) IsVariableContext() {}

func NewVariableContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VariableContext {
	var p = new(VariableContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_variable

	return p
}

func (s *VariableContext) GetParser() antlr.Parser { return s.parser }

func (s *VariableContext) MIXEDTAG() antlr.TerminalNode {
	return s.GetToken(PMMPostParserMIXEDTAG, 0)
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
	return s.GetToken(PMMPostParserTAG, 0)
}

func (s *VariableContext) LAMBDAARG() antlr.TerminalNode {
	return s.GetToken(PMMPostParserLAMBDAARG, 0)
}

func (s *VariableContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VariableContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VariableContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterVariable(s)
	}
}

func (s *VariableContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitVariable(s)
	}
}

func (p *PMMPostParser) Variable() (localctx IVariableContext) {
	localctx = NewVariableContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, PMMPostParserRULE_variable)

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

	p.SetState(392)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPostParserMIXEDTAG:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(375)
			p.Match(PMMPostParserMIXEDTAG)
		}
		p.SetState(380)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(378)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case PMMPostParserLBRACKET, PMMPostParserDECIMALTOKEN:
					{
						p.SetState(376)
						p.Subscript()
					}

				case PMMPostParserTAG, PMMPostParserMIXEDTAG:
					{
						p.SetState(377)
						p.Anytag()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}

			}
			p.SetState(382)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 34, p.GetParserRuleContext())
		}

	case PMMPostParserTAG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(383)
			p.Match(PMMPostParserTAG)
		}
		p.SetState(388)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext())

		for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			if _alt == 1 {
				p.SetState(386)
				p.GetErrorHandler().Sync(p)

				switch p.GetTokenStream().LA(1) {
				case PMMPostParserLBRACKET, PMMPostParserDECIMALTOKEN:
					{
						p.SetState(384)
						p.Subscript()
					}

				case PMMPostParserTAG, PMMPostParserMIXEDTAG:
					{
						p.SetState(385)
						p.Anytag()
					}

				default:
					panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				}

			}
			p.SetState(390)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 36, p.GetParserRuleContext())
		}

	case PMMPostParserLAMBDAARG:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(391)
			p.Match(PMMPostParserLAMBDAARG)
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
	p.RuleIndex = PMMPostParserRULE_subscript
	return p
}

func (*SubscriptContext) IsSubscriptContext() {}

func NewSubscriptContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubscriptContext {
	var p = new(SubscriptContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_subscript

	return p
}

func (s *SubscriptContext) GetParser() antlr.Parser { return s.parser }

func (s *SubscriptContext) DECIMALTOKEN() antlr.TerminalNode {
	return s.GetToken(PMMPostParserDECIMALTOKEN, 0)
}

func (s *SubscriptContext) LBRACKET() antlr.TerminalNode {
	return s.GetToken(PMMPostParserLBRACKET, 0)
}

func (s *SubscriptContext) Tertiary() ITertiaryContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITertiaryContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITertiaryContext)
}

func (s *SubscriptContext) RBRACKET() antlr.TerminalNode {
	return s.GetToken(PMMPostParserRBRACKET, 0)
}

func (s *SubscriptContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubscriptContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubscriptContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterSubscript(s)
	}
}

func (s *SubscriptContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitSubscript(s)
	}
}

func (p *PMMPostParser) Subscript() (localctx ISubscriptContext) {
	localctx = NewSubscriptContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, PMMPostParserRULE_subscript)

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

	p.SetState(399)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case PMMPostParserDECIMALTOKEN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(394)
			p.Match(PMMPostParserDECIMALTOKEN)
		}

	case PMMPostParserLBRACKET:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(395)
			p.Match(PMMPostParserLBRACKET)
		}
		{
			p.SetState(396)
			p.tertiary(0)
		}
		{
			p.SetState(397)
			p.Match(PMMPostParserRBRACKET)
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
	p.RuleIndex = PMMPostParserRULE_anytag
	return p
}

func (*AnytagContext) IsAnytagContext() {}

func NewAnytagContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnytagContext {
	var p = new(AnytagContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = PMMPostParserRULE_anytag

	return p
}

func (s *AnytagContext) GetParser() antlr.Parser { return s.parser }

func (s *AnytagContext) TAG() antlr.TerminalNode {
	return s.GetToken(PMMPostParserTAG, 0)
}

func (s *AnytagContext) MIXEDTAG() antlr.TerminalNode {
	return s.GetToken(PMMPostParserMIXEDTAG, 0)
}

func (s *AnytagContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnytagContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AnytagContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.EnterAnytag(s)
	}
}

func (s *AnytagContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(PMMPostListener); ok {
		listenerT.ExitAnytag(s)
	}
}

func (p *PMMPostParser) Anytag() (localctx IAnytagContext) {
	localctx = NewAnytagContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, PMMPostParserRULE_anytag)
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
		p.SetState(401)
		_la = p.GetTokenStream().LA(1)

		if !(_la == PMMPostParserTAG || _la == PMMPostParserMIXEDTAG) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

func (p *PMMPostParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 21:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	case 22:
		var t *TertiaryContext = nil
		if localctx != nil {
			t = localctx.(*TertiaryContext)
		}
		return p.Tertiary_Sempred(t, predIndex)

	case 25:
		var t *SecondaryContext = nil
		if localctx != nil {
			t = localctx.(*SecondaryContext)
		}
		return p.Secondary_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *PMMPostParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *PMMPostParser) Tertiary_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 1:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}

func (p *PMMPostParser) Secondary_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 2:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
