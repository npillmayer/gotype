// Generated from Gallery.g4 by ANTLR 4.7.

package grammar

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 52, 526,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23,
	9, 23, 4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9,
	28, 4, 29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33,
	4, 34, 9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4,
	39, 9, 39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44,
	9, 44, 4, 45, 9, 45, 4, 46, 9, 46, 4, 47, 9, 47, 4, 48, 9, 48, 4, 49, 9,
	49, 4, 50, 9, 50, 4, 51, 9, 51, 4, 52, 9, 52, 4, 53, 9, 53, 3, 2, 3, 2,
	3, 2, 5, 2, 111, 10, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4,
	3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4,
	3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4,
	3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 5, 4, 166, 10, 4,
	3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 7, 3, 7, 3, 8, 3, 8, 3, 9, 3, 9, 3, 10,
	3, 10, 3, 11, 3, 11, 3, 12, 3, 12, 3, 13, 3, 13, 3, 14, 3, 14, 3, 15, 3,
	15, 3, 16, 3, 16, 3, 17, 3, 17, 3, 18, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19,
	3, 20, 3, 20, 3, 21, 3, 21, 3, 21, 3, 21, 3, 21, 3, 21, 3, 21, 3, 21, 3,
	21, 3, 21, 3, 21, 3, 21, 5, 21, 215, 10, 21, 3, 22, 3, 22, 3, 23, 3, 23,
	3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 23, 3, 24, 3,
	24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 25, 3, 25, 3, 25,
	3, 25, 3, 25, 3, 25, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26, 3,
	27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 5, 27,
	262, 10, 27, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3,
	28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 3, 28, 5, 28,
	282, 10, 28, 3, 29, 3, 29, 3, 29, 3, 29, 3, 29, 3, 30, 3, 30, 3, 30, 3,
	30, 3, 30, 3, 30, 3, 31, 3, 31, 3, 31, 3, 31, 3, 32, 3, 32, 3, 32, 3, 32,
	3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3,
	32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32,
	3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 5, 32, 329, 10, 32, 3, 32, 5, 32, 332,
	10, 32, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 34,
	3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 35, 3, 35, 3, 35, 3,
	35, 3, 35, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36, 3, 37, 3, 37, 3, 37,
	3, 38, 3, 38, 3, 38, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3,
	39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39,
	3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39, 3,
	39, 3, 39, 3, 39, 5, 39, 398, 10, 39, 3, 40, 3, 40, 3, 40, 3, 40, 3, 40,
	3, 40, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3,
	41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41,
	3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 5, 41, 433, 10, 41, 3,
	42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43,
	3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 3, 45, 3, 45, 3, 45, 3, 45, 3, 46, 3,
	46, 6, 46, 457, 10, 46, 13, 46, 14, 46, 458, 3, 47, 6, 47, 462, 10, 47,
	13, 47, 14, 47, 463, 3, 48, 5, 48, 467, 10, 48, 3, 48, 3, 48, 7, 48, 471,
	10, 48, 12, 48, 14, 48, 474, 11, 48, 3, 49, 6, 49, 477, 10, 49, 13, 49,
	14, 49, 478, 3, 49, 3, 49, 6, 49, 483, 10, 49, 13, 49, 14, 49, 484, 5,
	49, 487, 10, 49, 3, 49, 5, 49, 490, 10, 49, 3, 50, 3, 50, 3, 51, 3, 51,
	3, 51, 6, 51, 497, 10, 51, 13, 51, 14, 51, 498, 3, 51, 3, 51, 3, 52, 3,
	52, 3, 52, 3, 52, 7, 52, 507, 10, 52, 12, 52, 14, 52, 510, 11, 52, 3, 52,
	5, 52, 513, 10, 52, 3, 52, 3, 52, 3, 53, 6, 53, 518, 10, 53, 13, 53, 14,
	53, 519, 3, 53, 5, 53, 523, 10, 53, 3, 53, 3, 53, 3, 508, 2, 54, 3, 3,
	5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 15, 9, 17, 10, 19, 11, 21, 12, 23, 13,
	25, 14, 27, 15, 29, 16, 31, 17, 33, 18, 35, 19, 37, 20, 39, 21, 41, 22,
	43, 23, 45, 24, 47, 25, 49, 26, 51, 27, 53, 28, 55, 29, 57, 30, 59, 31,
	61, 32, 63, 33, 65, 34, 67, 35, 69, 36, 71, 37, 73, 38, 75, 39, 77, 40,
	79, 41, 81, 42, 83, 43, 85, 44, 87, 45, 89, 46, 91, 2, 93, 47, 95, 48,
	97, 49, 99, 50, 101, 51, 103, 2, 105, 52, 3, 2, 6, 4, 2, 67, 92, 99, 124,
	6, 2, 48, 48, 50, 59, 67, 92, 99, 124, 8, 2, 34, 34, 47, 47, 50, 59, 67,
	92, 97, 97, 99, 124, 5, 2, 11, 12, 15, 15, 34, 34, 2, 562, 2, 3, 3, 2,
	2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2, 11, 3, 2,
	2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2, 17, 3, 2, 2, 2, 2, 19, 3,
	2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 23, 3, 2, 2, 2, 2, 25, 3, 2, 2, 2, 2, 27,
	3, 2, 2, 2, 2, 29, 3, 2, 2, 2, 2, 31, 3, 2, 2, 2, 2, 33, 3, 2, 2, 2, 2,
	35, 3, 2, 2, 2, 2, 37, 3, 2, 2, 2, 2, 39, 3, 2, 2, 2, 2, 41, 3, 2, 2, 2,
	2, 43, 3, 2, 2, 2, 2, 45, 3, 2, 2, 2, 2, 47, 3, 2, 2, 2, 2, 49, 3, 2, 2,
	2, 2, 51, 3, 2, 2, 2, 2, 53, 3, 2, 2, 2, 2, 55, 3, 2, 2, 2, 2, 57, 3, 2,
	2, 2, 2, 59, 3, 2, 2, 2, 2, 61, 3, 2, 2, 2, 2, 63, 3, 2, 2, 2, 2, 65, 3,
	2, 2, 2, 2, 67, 3, 2, 2, 2, 2, 69, 3, 2, 2, 2, 2, 71, 3, 2, 2, 2, 2, 73,
	3, 2, 2, 2, 2, 75, 3, 2, 2, 2, 2, 77, 3, 2, 2, 2, 2, 79, 3, 2, 2, 2, 2,
	81, 3, 2, 2, 2, 2, 83, 3, 2, 2, 2, 2, 85, 3, 2, 2, 2, 2, 87, 3, 2, 2, 2,
	2, 89, 3, 2, 2, 2, 2, 93, 3, 2, 2, 2, 2, 95, 3, 2, 2, 2, 2, 97, 3, 2, 2,
	2, 2, 99, 3, 2, 2, 2, 2, 101, 3, 2, 2, 2, 2, 105, 3, 2, 2, 2, 3, 110, 3,
	2, 2, 2, 5, 112, 3, 2, 2, 2, 7, 165, 3, 2, 2, 2, 9, 167, 3, 2, 2, 2, 11,
	170, 3, 2, 2, 2, 13, 172, 3, 2, 2, 2, 15, 174, 3, 2, 2, 2, 17, 176, 3,
	2, 2, 2, 19, 178, 3, 2, 2, 2, 21, 180, 3, 2, 2, 2, 23, 182, 3, 2, 2, 2,
	25, 184, 3, 2, 2, 2, 27, 186, 3, 2, 2, 2, 29, 188, 3, 2, 2, 2, 31, 190,
	3, 2, 2, 2, 33, 192, 3, 2, 2, 2, 35, 194, 3, 2, 2, 2, 37, 197, 3, 2, 2,
	2, 39, 200, 3, 2, 2, 2, 41, 214, 3, 2, 2, 2, 43, 216, 3, 2, 2, 2, 45, 218,
	3, 2, 2, 2, 47, 229, 3, 2, 2, 2, 49, 238, 3, 2, 2, 2, 51, 244, 3, 2, 2,
	2, 53, 261, 3, 2, 2, 2, 55, 281, 3, 2, 2, 2, 57, 283, 3, 2, 2, 2, 59, 288,
	3, 2, 2, 2, 61, 294, 3, 2, 2, 2, 63, 331, 3, 2, 2, 2, 65, 333, 3, 2, 2,
	2, 67, 341, 3, 2, 2, 2, 69, 349, 3, 2, 2, 2, 71, 354, 3, 2, 2, 2, 73, 360,
	3, 2, 2, 2, 75, 363, 3, 2, 2, 2, 77, 397, 3, 2, 2, 2, 79, 399, 3, 2, 2,
	2, 81, 432, 3, 2, 2, 2, 83, 434, 3, 2, 2, 2, 85, 440, 3, 2, 2, 2, 87, 445,
	3, 2, 2, 2, 89, 450, 3, 2, 2, 2, 91, 454, 3, 2, 2, 2, 93, 461, 3, 2, 2,
	2, 95, 466, 3, 2, 2, 2, 97, 476, 3, 2, 2, 2, 99, 491, 3, 2, 2, 2, 101,
	493, 3, 2, 2, 2, 103, 502, 3, 2, 2, 2, 105, 522, 3, 2, 2, 2, 107, 108,
	7, 47, 2, 2, 108, 111, 7, 47, 2, 2, 109, 111, 7, 40, 2, 2, 110, 107, 3,
	2, 2, 2, 110, 109, 3, 2, 2, 2, 111, 4, 3, 2, 2, 2, 112, 113, 7, 114, 2,
	2, 113, 114, 7, 99, 2, 2, 114, 115, 7, 116, 2, 2, 115, 116, 7, 99, 2, 2,
	116, 117, 7, 111, 2, 2, 117, 118, 7, 103, 2, 2, 118, 119, 7, 118, 2, 2,
	119, 120, 7, 103, 2, 2, 120, 121, 7, 116, 2, 2, 121, 6, 3, 2, 2, 2, 122,
	123, 7, 112, 2, 2, 123, 124, 7, 119, 2, 2, 124, 125, 7, 111, 2, 2, 125,
	126, 7, 103, 2, 2, 126, 127, 7, 116, 2, 2, 127, 128, 7, 107, 2, 2, 128,
	166, 7, 101, 2, 2, 129, 130, 7, 114, 2, 2, 130, 131, 7, 99, 2, 2, 131,
	132, 7, 107, 2, 2, 132, 166, 7, 116, 2, 2, 133, 134, 7, 114, 2, 2, 134,
	135, 7, 99, 2, 2, 135, 136, 7, 118, 2, 2, 136, 166, 7, 106, 2, 2, 137,
	138, 7, 118, 2, 2, 138, 139, 7, 116, 2, 2, 139, 140, 7, 99, 2, 2, 140,
	141, 7, 112, 2, 2, 141, 142, 7, 117, 2, 2, 142, 143, 7, 104, 2, 2, 143,
	144, 7, 113, 2, 2, 144, 145, 7, 116, 2, 2, 145, 166, 7, 111, 2, 2, 146,
	147, 7, 104, 2, 2, 147, 148, 7, 116, 2, 2, 148, 149, 7, 99, 2, 2, 149,
	150, 7, 111, 2, 2, 150, 151, 7, 103, 2, 2, 151, 152, 7, 102, 2, 2, 152,
	153, 7, 100, 2, 2, 153, 154, 7, 113, 2, 2, 154, 166, 7, 122, 2, 2, 155,
	156, 7, 104, 2, 2, 156, 157, 7, 116, 2, 2, 157, 158, 7, 99, 2, 2, 158,
	159, 7, 111, 2, 2, 159, 160, 7, 103, 2, 2, 160, 161, 7, 102, 2, 2, 161,
	162, 7, 114, 2, 2, 162, 163, 7, 99, 2, 2, 163, 164, 7, 118, 2, 2, 164,
	166, 7, 106, 2, 2, 165, 122, 3, 2, 2, 2, 165, 129, 3, 2, 2, 2, 165, 133,
	3, 2, 2, 2, 165, 137, 3, 2, 2, 2, 165, 146, 3, 2, 2, 2, 165, 155, 3, 2,
	2, 2, 166, 8, 3, 2, 2, 2, 167, 168, 7, 60, 2, 2, 168, 169, 7, 63, 2, 2,
	169, 10, 3, 2, 2, 2, 170, 171, 7, 63, 2, 2, 171, 12, 3, 2, 2, 2, 172, 173,
	7, 60, 2, 2, 173, 14, 3, 2, 2, 2, 174, 175, 7, 61, 2, 2, 175, 16, 3, 2,
	2, 2, 176, 177, 7, 46, 2, 2, 177, 18, 3, 2, 2, 2, 178, 179, 7, 42, 2, 2,
	179, 20, 3, 2, 2, 2, 180, 181, 7, 43, 2, 2, 181, 22, 3, 2, 2, 2, 182, 183,
	7, 93, 2, 2, 183, 24, 3, 2, 2, 2, 184, 185, 7, 95, 2, 2, 185, 26, 3, 2,
	2, 2, 186, 187, 7, 45, 2, 2, 187, 28, 3, 2, 2, 2, 188, 189, 7, 47, 2, 2,
	189, 30, 3, 2, 2, 2, 190, 191, 7, 44, 2, 2, 191, 32, 3, 2, 2, 2, 192, 193,
	7, 49, 2, 2, 193, 34, 3, 2, 2, 2, 194, 195, 7, 126, 2, 2, 195, 196, 7,
	126, 2, 2, 196, 36, 3, 2, 2, 2, 197, 198, 7, 126, 2, 2, 198, 199, 7, 47,
	2, 2, 199, 38, 3, 2, 2, 2, 200, 201, 7, 128, 2, 2, 201, 40, 3, 2, 2, 2,
	202, 203, 7, 100, 2, 2, 203, 215, 7, 114, 2, 2, 204, 205, 7, 111, 2, 2,
	205, 215, 7, 111, 2, 2, 206, 207, 7, 101, 2, 2, 207, 215, 7, 111, 2, 2,
	208, 209, 7, 107, 2, 2, 209, 215, 7, 112, 2, 2, 210, 211, 7, 114, 2, 2,
	211, 215, 7, 118, 2, 2, 212, 213, 7, 114, 2, 2, 213, 215, 7, 101, 2, 2,
	214, 202, 3, 2, 2, 2, 214, 204, 3, 2, 2, 2, 214, 206, 3, 2, 2, 2, 214,
	208, 3, 2, 2, 2, 214, 210, 3, 2, 2, 2, 214, 212, 3, 2, 2, 2, 215, 42, 3,
	2, 2, 2, 216, 217, 7, 66, 2, 2, 217, 44, 3, 2, 2, 2, 218, 219, 7, 100,
	2, 2, 219, 220, 7, 103, 2, 2, 220, 221, 7, 105, 2, 2, 221, 222, 7, 107,
	2, 2, 222, 223, 7, 112, 2, 2, 223, 224, 7, 105, 2, 2, 224, 225, 7, 116,
	2, 2, 225, 226, 7, 113, 2, 2, 226, 227, 7, 119, 2, 2, 227, 228, 7, 114,
	2, 2, 228, 46, 3, 2, 2, 2, 229, 230, 7, 103, 2, 2, 230, 231, 7, 112, 2,
	2, 231, 232, 7, 102, 2, 2, 232, 233, 7, 105, 2, 2, 233, 234, 7, 116, 2,
	2, 234, 235, 7, 113, 2, 2, 235, 236, 7, 119, 2, 2, 236, 237, 7, 114, 2,
	2, 237, 48, 3, 2, 2, 2, 238, 239, 7, 110, 2, 2, 239, 240, 7, 113, 2, 2,
	240, 241, 7, 101, 2, 2, 241, 242, 7, 99, 2, 2, 242, 243, 7, 110, 2, 2,
	243, 50, 3, 2, 2, 2, 244, 245, 7, 120, 2, 2, 245, 246, 7, 99, 2, 2, 246,
	247, 7, 116, 2, 2, 247, 248, 7, 102, 2, 2, 248, 249, 7, 103, 2, 2, 249,
	250, 7, 104, 2, 2, 250, 52, 3, 2, 2, 2, 251, 252, 7, 122, 2, 2, 252, 253,
	7, 114, 2, 2, 253, 254, 7, 99, 2, 2, 254, 255, 7, 116, 2, 2, 255, 262,
	7, 118, 2, 2, 256, 257, 7, 123, 2, 2, 257, 258, 7, 114, 2, 2, 258, 259,
	7, 99, 2, 2, 259, 260, 7, 116, 2, 2, 260, 262, 7, 118, 2, 2, 261, 251,
	3, 2, 2, 2, 261, 256, 3, 2, 2, 2, 262, 54, 3, 2, 2, 2, 263, 264, 7, 118,
	2, 2, 264, 265, 7, 113, 2, 2, 265, 282, 7, 114, 2, 2, 266, 267, 7, 110,
	2, 2, 267, 268, 7, 103, 2, 2, 268, 269, 7, 104, 2, 2, 269, 282, 7, 118,
	2, 2, 270, 271, 7, 116, 2, 2, 271, 272, 7, 107, 2, 2, 272, 273, 7, 105,
	2, 2, 273, 274, 7, 106, 2, 2, 274, 282, 7, 118, 2, 2, 275, 276, 7, 100,
	2, 2, 276, 277, 7, 113, 2, 2, 277, 278, 7, 118, 2, 2, 278, 279, 7, 118,
	2, 2, 279, 280, 7, 113, 2, 2, 280, 282, 7, 111, 2, 2, 281, 263, 3, 2, 2,
	2, 281, 266, 3, 2, 2, 2, 281, 270, 3, 2, 2, 2, 281, 275, 3, 2, 2, 2, 282,
	56, 3, 2, 2, 2, 283, 284, 7, 103, 2, 2, 284, 285, 7, 102, 2, 2, 285, 286,
	7, 105, 2, 2, 286, 287, 7, 103, 2, 2, 287, 58, 3, 2, 2, 2, 288, 289, 7,
	104, 2, 2, 289, 290, 7, 116, 2, 2, 290, 291, 7, 99, 2, 2, 291, 292, 7,
	111, 2, 2, 292, 293, 7, 103, 2, 2, 293, 60, 3, 2, 2, 2, 294, 295, 7, 100,
	2, 2, 295, 296, 7, 113, 2, 2, 296, 297, 7, 122, 2, 2, 297, 62, 3, 2, 2,
	2, 298, 299, 7, 102, 2, 2, 299, 300, 7, 107, 2, 2, 300, 301, 7, 117, 2,
	2, 301, 302, 7, 118, 2, 2, 302, 303, 7, 99, 2, 2, 303, 304, 7, 112, 2,
	2, 304, 305, 7, 101, 2, 2, 305, 329, 7, 103, 2, 2, 306, 307, 7, 110, 2,
	2, 307, 308, 7, 103, 2, 2, 308, 309, 7, 112, 2, 2, 309, 310, 7, 105, 2,
	2, 310, 311, 7, 118, 2, 2, 311, 329, 7, 106, 2, 2, 312, 313, 7, 104, 2,
	2, 313, 314, 7, 110, 2, 2, 314, 315, 7, 113, 2, 2, 315, 316, 7, 113, 2,
	2, 316, 329, 7, 116, 2, 2, 317, 318, 7, 101, 2, 2, 318, 319, 7, 103, 2,
	2, 319, 320, 7, 107, 2, 2, 320, 329, 7, 110, 2, 2, 321, 322, 7, 117, 2,
	2, 322, 323, 7, 115, 2, 2, 323, 324, 7, 116, 2, 2, 324, 329, 7, 118, 2,
	2, 325, 326, 7, 102, 2, 2, 326, 327, 7, 107, 2, 2, 327, 329, 7, 116, 2,
	2, 328, 298, 3, 2, 2, 2, 328, 306, 3, 2, 2, 2, 328, 312, 3, 2, 2, 2, 328,
	317, 3, 2, 2, 2, 328, 321, 3, 2, 2, 2, 328, 325, 3, 2, 2, 2, 329, 332,
	3, 2, 2, 2, 330, 332, 5, 91, 46, 2, 331, 328, 3, 2, 2, 2, 331, 330, 3,
	2, 2, 2, 332, 64, 3, 2, 2, 2, 333, 334, 7, 117, 2, 2, 334, 335, 7, 119,
	2, 2, 335, 336, 7, 100, 2, 2, 336, 337, 7, 114, 2, 2, 337, 338, 7, 99,
	2, 2, 338, 339, 7, 118, 2, 2, 339, 340, 7, 106, 2, 2, 340, 66, 3, 2, 2,
	2, 341, 342, 7, 116, 2, 2, 342, 343, 7, 103, 2, 2, 343, 344, 7, 120, 2,
	2, 344, 345, 7, 103, 2, 2, 345, 346, 7, 116, 2, 2, 346, 347, 7, 117, 2,
	2, 347, 348, 7, 103, 2, 2, 348, 68, 3, 2, 2, 2, 349, 350, 7, 121, 2, 2,
	350, 351, 7, 107, 2, 2, 351, 352, 7, 118, 2, 2, 352, 353, 7, 106, 2, 2,
	353, 70, 3, 2, 2, 2, 354, 355, 7, 114, 2, 2, 355, 356, 7, 113, 2, 2, 356,
	357, 7, 107, 2, 2, 357, 358, 7, 112, 2, 2, 358, 359, 7, 118, 2, 2, 359,
	72, 3, 2, 2, 2, 360, 361, 7, 113, 2, 2, 361, 362, 7, 104, 2, 2, 362, 74,
	3, 2, 2, 2, 363, 364, 7, 118, 2, 2, 364, 365, 7, 113, 2, 2, 365, 76, 3,
	2, 2, 2, 366, 367, 7, 117, 2, 2, 367, 368, 7, 101, 2, 2, 368, 369, 7, 99,
	2, 2, 369, 370, 7, 110, 2, 2, 370, 371, 7, 103, 2, 2, 371, 398, 7, 102,
	2, 2, 372, 373, 7, 117, 2, 2, 373, 374, 7, 106, 2, 2, 374, 375, 7, 107,
	2, 2, 375, 376, 7, 104, 2, 2, 376, 377, 7, 118, 2, 2, 377, 378, 7, 103,
	2, 2, 378, 398, 7, 102, 2, 2, 379, 380, 7, 116, 2, 2, 380, 381, 7, 113,
	2, 2, 381, 382, 7, 118, 2, 2, 382, 383, 7, 99, 2, 2, 383, 384, 7, 118,
	2, 2, 384, 385, 7, 103, 2, 2, 385, 398, 7, 102, 2, 2, 386, 387, 7, 118,
	2, 2, 387, 388, 7, 116, 2, 2, 388, 389, 7, 99, 2, 2, 389, 390, 7, 112,
	2, 2, 390, 391, 7, 117, 2, 2, 391, 392, 7, 104, 2, 2, 392, 393, 7, 113,
	2, 2, 393, 394, 7, 116, 2, 2, 394, 395, 7, 111, 2, 2, 395, 396, 7, 103,
	2, 2, 396, 398, 7, 102, 2, 2, 397, 366, 3, 2, 2, 2, 397, 372, 3, 2, 2,
	2, 397, 379, 3, 2, 2, 2, 397, 386, 3, 2, 2, 2, 398, 78, 3, 2, 2, 2, 399,
	400, 7, 101, 2, 2, 400, 401, 7, 123, 2, 2, 401, 402, 7, 101, 2, 2, 402,
	403, 7, 110, 2, 2, 403, 404, 7, 103, 2, 2, 404, 80, 3, 2, 2, 2, 405, 406,
	7, 119, 2, 2, 406, 407, 7, 112, 2, 2, 407, 408, 7, 107, 2, 2, 408, 409,
	7, 113, 2, 2, 409, 433, 7, 112, 2, 2, 410, 411, 7, 107, 2, 2, 411, 412,
	7, 112, 2, 2, 412, 413, 7, 118, 2, 2, 413, 414, 7, 103, 2, 2, 414, 415,
	7, 116, 2, 2, 415, 416, 7, 117, 2, 2, 416, 417, 7, 103, 2, 2, 417, 418,
	7, 101, 2, 2, 418, 419, 7, 118, 2, 2, 419, 420, 7, 107, 2, 2, 420, 421,
	7, 113, 2, 2, 421, 433, 7, 112, 2, 2, 422, 423, 7, 102, 2, 2, 423, 424,
	7, 107, 2, 2, 424, 425, 7, 104, 2, 2, 425, 426, 7, 104, 2, 2, 426, 427,
	7, 103, 2, 2, 427, 428, 7, 116, 2, 2, 428, 429, 7, 103, 2, 2, 429, 430,
	7, 112, 2, 2, 430, 431, 7, 101, 2, 2, 431, 433, 7, 103, 2, 2, 432, 405,
	3, 2, 2, 2, 432, 410, 3, 2, 2, 2, 432, 422, 3, 2, 2, 2, 433, 82, 3, 2,
	2, 2, 434, 435, 7, 114, 2, 2, 435, 436, 7, 116, 2, 2, 436, 437, 7, 113,
	2, 2, 437, 438, 7, 113, 2, 2, 438, 439, 7, 104, 2, 2, 439, 84, 3, 2, 2,
	2, 440, 441, 7, 117, 2, 2, 441, 442, 7, 99, 2, 2, 442, 443, 7, 120, 2,
	2, 443, 444, 7, 103, 2, 2, 444, 86, 3, 2, 2, 2, 445, 446, 7, 117, 2, 2,
	446, 447, 7, 106, 2, 2, 447, 448, 7, 113, 2, 2, 448, 449, 7, 121, 2, 2,
	449, 88, 3, 2, 2, 2, 450, 451, 7, 110, 2, 2, 451, 452, 7, 103, 2, 2, 452,
	453, 7, 118, 2, 2, 453, 90, 3, 2, 2, 2, 454, 456, 7, 66, 2, 2, 455, 457,
	9, 2, 2, 2, 456, 455, 3, 2, 2, 2, 457, 458, 3, 2, 2, 2, 458, 456, 3, 2,
	2, 2, 458, 459, 3, 2, 2, 2, 459, 92, 3, 2, 2, 2, 460, 462, 9, 2, 2, 2,
	461, 460, 3, 2, 2, 2, 462, 463, 3, 2, 2, 2, 463, 461, 3, 2, 2, 2, 463,
	464, 3, 2, 2, 2, 464, 94, 3, 2, 2, 2, 465, 467, 7, 48, 2, 2, 466, 465,
	3, 2, 2, 2, 466, 467, 3, 2, 2, 2, 467, 468, 3, 2, 2, 2, 468, 472, 9, 2,
	2, 2, 469, 471, 9, 3, 2, 2, 470, 469, 3, 2, 2, 2, 471, 474, 3, 2, 2, 2,
	472, 470, 3, 2, 2, 2, 472, 473, 3, 2, 2, 2, 473, 96, 3, 2, 2, 2, 474, 472,
	3, 2, 2, 2, 475, 477, 4, 50, 59, 2, 476, 475, 3, 2, 2, 2, 477, 478, 3,
	2, 2, 2, 478, 476, 3, 2, 2, 2, 478, 479, 3, 2, 2, 2, 479, 486, 3, 2, 2,
	2, 480, 482, 5, 99, 50, 2, 481, 483, 4, 50, 59, 2, 482, 481, 3, 2, 2, 2,
	483, 484, 3, 2, 2, 2, 484, 482, 3, 2, 2, 2, 484, 485, 3, 2, 2, 2, 485,
	487, 3, 2, 2, 2, 486, 480, 3, 2, 2, 2, 486, 487, 3, 2, 2, 2, 487, 489,
	3, 2, 2, 2, 488, 490, 7, 39, 2, 2, 489, 488, 3, 2, 2, 2, 489, 490, 3, 2,
	2, 2, 490, 98, 3, 2, 2, 2, 491, 492, 7, 48, 2, 2, 492, 100, 3, 2, 2, 2,
	493, 494, 7, 36, 2, 2, 494, 496, 9, 2, 2, 2, 495, 497, 9, 4, 2, 2, 496,
	495, 3, 2, 2, 2, 497, 498, 3, 2, 2, 2, 498, 496, 3, 2, 2, 2, 498, 499,
	3, 2, 2, 2, 499, 500, 3, 2, 2, 2, 500, 501, 7, 36, 2, 2, 501, 102, 3, 2,
	2, 2, 502, 503, 7, 49, 2, 2, 503, 504, 7, 49, 2, 2, 504, 508, 3, 2, 2,
	2, 505, 507, 11, 2, 2, 2, 506, 505, 3, 2, 2, 2, 507, 510, 3, 2, 2, 2, 508,
	509, 3, 2, 2, 2, 508, 506, 3, 2, 2, 2, 509, 512, 3, 2, 2, 2, 510, 508,
	3, 2, 2, 2, 511, 513, 7, 15, 2, 2, 512, 511, 3, 2, 2, 2, 512, 513, 3, 2,
	2, 2, 513, 514, 3, 2, 2, 2, 514, 515, 7, 12, 2, 2, 515, 104, 3, 2, 2, 2,
	516, 518, 9, 5, 2, 2, 517, 516, 3, 2, 2, 2, 518, 519, 3, 2, 2, 2, 519,
	517, 3, 2, 2, 2, 519, 520, 3, 2, 2, 2, 520, 523, 3, 2, 2, 2, 521, 523,
	5, 103, 52, 2, 522, 517, 3, 2, 2, 2, 522, 521, 3, 2, 2, 2, 523, 524, 3,
	2, 2, 2, 524, 525, 8, 53, 2, 2, 525, 106, 3, 2, 2, 2, 25, 2, 110, 165,
	214, 261, 281, 328, 331, 397, 432, 458, 463, 466, 472, 478, 484, 486, 489,
	498, 508, 512, 519, 522, 3, 8, 2, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "", "'parameter'", "", "':='", "'='", "':'", "';'", "','", "'('", "')'",
	"'['", "']'", "'+'", "'-'", "'*'", "'/'", "'||'", "'|-'", "'~'", "", "'@'",
	"'begingroup'", "'endgroup'", "'local'", "'vardef'", "", "", "'edge'",
	"'frame'", "'box'", "", "'subpath'", "'reverse'", "'with'", "'point'",
	"'of'", "'to'", "", "'cycle'", "", "'proof'", "'save'", "'show'", "'let'",
	"", "", "", "'.'",
}

var lexerSymbolicNames = []string{
	"", "PATHJOIN", "PARAMETER", "TYPE", "ASSIGN", "EQUALS", "COLON", "SEMIC",
	"COMMA", "LPAREN", "RPAREN", "LBRACKET", "RBRACKET", "PLUS", "MINUS", "TIMES",
	"OVER", "PARALLEL", "PERPENDIC", "CONGRUENT", "UNIT", "LAMBDAARG", "BEGINGROUP",
	"ENDGROUP", "LOCAL", "VARDEF", "PAIRPART", "EDGECONSTR", "EDGE", "FRAME",
	"BOX", "MATHFUNC", "SUBPATH", "REVERSE", "WITH", "POINT", "OF", "TO", "TRANSFORM",
	"CYCLE", "PATHCLIPOP", "PROOF", "SAVE", "SHOW", "LET", "TAG", "MIXEDTAG",
	"DECIMALTOKEN", "DOT", "LABEL", "WS",
}

var lexerRuleNames = []string{
	"PATHJOIN", "PARAMETER", "TYPE", "ASSIGN", "EQUALS", "COLON", "SEMIC",
	"COMMA", "LPAREN", "RPAREN", "LBRACKET", "RBRACKET", "PLUS", "MINUS", "TIMES",
	"OVER", "PARALLEL", "PERPENDIC", "CONGRUENT", "UNIT", "LAMBDAARG", "BEGINGROUP",
	"ENDGROUP", "LOCAL", "VARDEF", "PAIRPART", "EDGECONSTR", "EDGE", "FRAME",
	"BOX", "MATHFUNC", "SUBPATH", "REVERSE", "WITH", "POINT", "OF", "TO", "TRANSFORM",
	"CYCLE", "PATHCLIPOP", "PROOF", "SAVE", "SHOW", "LET", "LUAFUNC", "TAG",
	"MIXEDTAG", "DECIMALTOKEN", "DOT", "LABEL", "LINECOMMENT", "WS",
}

type GalleryLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewGalleryLexer(input antlr.CharStream) *GalleryLexer {

	l := new(GalleryLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "Gallery.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// GalleryLexer tokens.
const (
	GalleryLexerPATHJOIN     = 1
	GalleryLexerPARAMETER    = 2
	GalleryLexerTYPE         = 3
	GalleryLexerASSIGN       = 4
	GalleryLexerEQUALS       = 5
	GalleryLexerCOLON        = 6
	GalleryLexerSEMIC        = 7
	GalleryLexerCOMMA        = 8
	GalleryLexerLPAREN       = 9
	GalleryLexerRPAREN       = 10
	GalleryLexerLBRACKET     = 11
	GalleryLexerRBRACKET     = 12
	GalleryLexerPLUS         = 13
	GalleryLexerMINUS        = 14
	GalleryLexerTIMES        = 15
	GalleryLexerOVER         = 16
	GalleryLexerPARALLEL     = 17
	GalleryLexerPERPENDIC    = 18
	GalleryLexerCONGRUENT    = 19
	GalleryLexerUNIT         = 20
	GalleryLexerLAMBDAARG    = 21
	GalleryLexerBEGINGROUP   = 22
	GalleryLexerENDGROUP     = 23
	GalleryLexerLOCAL        = 24
	GalleryLexerVARDEF       = 25
	GalleryLexerPAIRPART     = 26
	GalleryLexerEDGECONSTR   = 27
	GalleryLexerEDGE         = 28
	GalleryLexerFRAME        = 29
	GalleryLexerBOX          = 30
	GalleryLexerMATHFUNC     = 31
	GalleryLexerSUBPATH      = 32
	GalleryLexerREVERSE      = 33
	GalleryLexerWITH         = 34
	GalleryLexerPOINT        = 35
	GalleryLexerOF           = 36
	GalleryLexerTO           = 37
	GalleryLexerTRANSFORM    = 38
	GalleryLexerCYCLE        = 39
	GalleryLexerPATHCLIPOP   = 40
	GalleryLexerPROOF        = 41
	GalleryLexerSAVE         = 42
	GalleryLexerSHOW         = 43
	GalleryLexerLET          = 44
	GalleryLexerTAG          = 45
	GalleryLexerMIXEDTAG     = 46
	GalleryLexerDECIMALTOKEN = 47
	GalleryLexerDOT          = 48
	GalleryLexerLABEL        = 49
	GalleryLexerWS           = 50
)