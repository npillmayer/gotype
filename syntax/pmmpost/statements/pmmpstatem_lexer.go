// Generated from PMMPStatem.g4 by ANTLR 4.7.

package statements

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 58, 566,
	8, 1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7,
	9, 7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12,
	4, 13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4,
	18, 9, 18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23,
	9, 23, 4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9,
	28, 4, 29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33,
	4, 34, 9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4,
	39, 9, 39, 4, 40, 9, 40, 4, 41, 9, 41, 4, 42, 9, 42, 4, 43, 9, 43, 4, 44,
	9, 44, 4, 45, 9, 45, 4, 46, 9, 46, 4, 47, 9, 47, 4, 48, 9, 48, 4, 49, 9,
	49, 4, 50, 9, 50, 4, 51, 9, 51, 4, 52, 9, 52, 4, 53, 9, 53, 4, 54, 9, 54,
	4, 55, 9, 55, 4, 56, 9, 56, 4, 57, 9, 57, 4, 58, 9, 58, 3, 2, 3, 2, 3,
	2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3,
	5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3,
	6, 3, 7, 3, 7, 3, 7, 3, 7, 3, 7, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 3,
	8, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 9, 3, 10, 3,
	10, 3, 10, 3, 11, 3, 11, 3, 12, 3, 12, 3, 13, 3, 13, 3, 14, 3, 14, 3, 15,
	3, 15, 3, 16, 3, 16, 3, 17, 3, 17, 3, 18, 3, 18, 3, 19, 3, 19, 3, 19, 3,
	19, 5, 19, 202, 10, 19, 3, 19, 3, 19, 3, 19, 3, 19, 5, 19, 208, 10, 19,
	3, 19, 5, 19, 211, 10, 19, 3, 20, 3, 20, 3, 21, 3, 21, 3, 22, 3, 22, 3,
	23, 3, 23, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24, 3, 24,
	5, 24, 230, 10, 24, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3, 25, 3,
	25, 3, 25, 3, 25, 5, 25, 242, 10, 25, 3, 26, 3, 26, 3, 26, 3, 26, 3, 26,
	3, 26, 3, 26, 3, 26, 3, 27, 3, 27, 3, 27, 3, 27, 3, 27, 3, 28, 3, 28, 3,
	28, 3, 28, 3, 28, 3, 29, 3, 29, 3, 29, 3, 29, 3, 29, 3, 29, 3, 29, 3, 29,
	3, 29, 3, 29, 3, 29, 5, 29, 273, 10, 29, 3, 30, 3, 30, 3, 30, 3, 30, 3,
	30, 3, 30, 3, 30, 3, 30, 3, 30, 3, 30, 5, 30, 285, 10, 30, 3, 31, 3, 31,
	3, 31, 3, 31, 3, 31, 3, 31, 3, 31, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3,
	32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 3, 32, 5, 32, 307, 10, 32,
	3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3,
	34, 3, 35, 3, 35, 3, 35, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36, 3, 36,
	3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 38, 3, 38, 3,
	38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 39, 3, 39, 3, 39, 3, 39, 3, 39,
	3, 39, 3, 39, 3, 39, 3, 39, 3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 3, 40, 3,
	40, 3, 40, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 41, 3, 42,
	3, 42, 3, 42, 3, 42, 3, 42, 3, 42, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43, 3,
	43, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43, 3, 43,
	3, 43, 3, 43, 3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 3,
	44, 3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 3, 44, 5, 44,
	413, 10, 44, 3, 45, 3, 45, 3, 45, 3, 45, 3, 45, 3, 45, 3, 45, 5, 45, 422,
	10, 45, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46,
	3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3,
	46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 3, 46, 5, 46, 451, 10, 46,
	3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3, 47, 3,
	47, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 48, 3, 49,
	3, 49, 7, 49, 475, 10, 49, 12, 49, 14, 49, 478, 11, 49, 3, 50, 3, 50, 7,
	50, 482, 10, 50, 12, 50, 14, 50, 485, 11, 50, 3, 51, 3, 51, 3, 51, 5, 51,
	490, 10, 51, 3, 51, 7, 51, 493, 10, 51, 12, 51, 14, 51, 496, 11, 51, 3,
	52, 3, 52, 3, 52, 5, 52, 501, 10, 52, 3, 52, 7, 52, 504, 10, 52, 12, 52,
	14, 52, 507, 11, 52, 3, 53, 3, 53, 7, 53, 511, 10, 53, 12, 53, 14, 53,
	514, 11, 53, 3, 54, 6, 54, 517, 10, 54, 13, 54, 14, 54, 518, 3, 54, 3,
	54, 6, 54, 523, 10, 54, 13, 54, 14, 54, 524, 5, 54, 527, 10, 54, 3, 54,
	5, 54, 530, 10, 54, 3, 55, 3, 55, 3, 56, 3, 56, 3, 56, 6, 56, 537, 10,
	56, 13, 56, 14, 56, 538, 3, 56, 3, 56, 3, 57, 3, 57, 3, 57, 3, 57, 7, 57,
	547, 10, 57, 12, 57, 14, 57, 550, 11, 57, 3, 57, 5, 57, 553, 10, 57, 3,
	57, 3, 57, 3, 58, 6, 58, 558, 10, 58, 13, 58, 14, 58, 559, 3, 58, 5, 58,
	563, 10, 58, 3, 58, 3, 58, 3, 548, 2, 59, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7,
	13, 8, 15, 9, 17, 10, 19, 11, 21, 12, 23, 13, 25, 14, 27, 15, 29, 16, 31,
	17, 33, 18, 35, 19, 37, 20, 39, 21, 41, 22, 43, 23, 45, 24, 47, 25, 49,
	26, 51, 27, 53, 28, 55, 29, 57, 30, 59, 31, 61, 32, 63, 33, 65, 34, 67,
	35, 69, 36, 71, 37, 73, 38, 75, 39, 77, 40, 79, 41, 81, 42, 83, 43, 85,
	44, 87, 45, 89, 46, 91, 47, 93, 48, 95, 49, 97, 50, 99, 51, 101, 52, 103,
	53, 105, 54, 107, 55, 109, 56, 111, 57, 113, 2, 115, 58, 3, 2, 12, 4, 2,
	50, 59, 99, 104, 4, 2, 114, 115, 124, 124, 3, 2, 99, 124, 4, 2, 99, 113,
	116, 123, 5, 2, 48, 48, 50, 59, 99, 124, 3, 2, 67, 92, 6, 2, 48, 48, 50,
	59, 67, 92, 99, 124, 4, 2, 67, 92, 99, 124, 8, 2, 34, 34, 47, 47, 50, 59,
	67, 92, 97, 97, 99, 124, 5, 2, 11, 12, 15, 15, 34, 34, 2, 597, 2, 3, 3,
	2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2, 11, 3,
	2, 2, 2, 2, 13, 3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 2, 17, 3, 2, 2, 2, 2, 19,
	3, 2, 2, 2, 2, 21, 3, 2, 2, 2, 2, 23, 3, 2, 2, 2, 2, 25, 3, 2, 2, 2, 2,
	27, 3, 2, 2, 2, 2, 29, 3, 2, 2, 2, 2, 31, 3, 2, 2, 2, 2, 33, 3, 2, 2, 2,
	2, 35, 3, 2, 2, 2, 2, 37, 3, 2, 2, 2, 2, 39, 3, 2, 2, 2, 2, 41, 3, 2, 2,
	2, 2, 43, 3, 2, 2, 2, 2, 45, 3, 2, 2, 2, 2, 47, 3, 2, 2, 2, 2, 49, 3, 2,
	2, 2, 2, 51, 3, 2, 2, 2, 2, 53, 3, 2, 2, 2, 2, 55, 3, 2, 2, 2, 2, 57, 3,
	2, 2, 2, 2, 59, 3, 2, 2, 2, 2, 61, 3, 2, 2, 2, 2, 63, 3, 2, 2, 2, 2, 65,
	3, 2, 2, 2, 2, 67, 3, 2, 2, 2, 2, 69, 3, 2, 2, 2, 2, 71, 3, 2, 2, 2, 2,
	73, 3, 2, 2, 2, 2, 75, 3, 2, 2, 2, 2, 77, 3, 2, 2, 2, 2, 79, 3, 2, 2, 2,
	2, 81, 3, 2, 2, 2, 2, 83, 3, 2, 2, 2, 2, 85, 3, 2, 2, 2, 2, 87, 3, 2, 2,
	2, 2, 89, 3, 2, 2, 2, 2, 91, 3, 2, 2, 2, 2, 93, 3, 2, 2, 2, 2, 95, 3, 2,
	2, 2, 2, 97, 3, 2, 2, 2, 2, 99, 3, 2, 2, 2, 2, 101, 3, 2, 2, 2, 2, 103,
	3, 2, 2, 2, 2, 105, 3, 2, 2, 2, 2, 107, 3, 2, 2, 2, 2, 109, 3, 2, 2, 2,
	2, 111, 3, 2, 2, 2, 2, 115, 3, 2, 2, 2, 3, 117, 3, 2, 2, 2, 5, 126, 3,
	2, 2, 2, 7, 133, 3, 2, 2, 2, 9, 138, 3, 2, 2, 2, 11, 151, 3, 2, 2, 2, 13,
	156, 3, 2, 2, 2, 15, 161, 3, 2, 2, 2, 17, 168, 3, 2, 2, 2, 19, 178, 3,
	2, 2, 2, 21, 181, 3, 2, 2, 2, 23, 183, 3, 2, 2, 2, 25, 185, 3, 2, 2, 2,
	27, 187, 3, 2, 2, 2, 29, 189, 3, 2, 2, 2, 31, 191, 3, 2, 2, 2, 33, 193,
	3, 2, 2, 2, 35, 195, 3, 2, 2, 2, 37, 210, 3, 2, 2, 2, 39, 212, 3, 2, 2,
	2, 41, 214, 3, 2, 2, 2, 43, 216, 3, 2, 2, 2, 45, 218, 3, 2, 2, 2, 47, 229,
	3, 2, 2, 2, 49, 241, 3, 2, 2, 2, 51, 243, 3, 2, 2, 2, 53, 251, 3, 2, 2,
	2, 55, 256, 3, 2, 2, 2, 57, 272, 3, 2, 2, 2, 59, 284, 3, 2, 2, 2, 61, 286,
	3, 2, 2, 2, 63, 306, 3, 2, 2, 2, 65, 308, 3, 2, 2, 2, 67, 313, 3, 2, 2,
	2, 69, 319, 3, 2, 2, 2, 71, 322, 3, 2, 2, 2, 73, 329, 3, 2, 2, 2, 75, 337,
	3, 2, 2, 2, 77, 345, 3, 2, 2, 2, 79, 354, 3, 2, 2, 2, 81, 362, 3, 2, 2,
	2, 83, 370, 3, 2, 2, 2, 85, 376, 3, 2, 2, 2, 87, 412, 3, 2, 2, 2, 89, 414,
	3, 2, 2, 2, 91, 450, 3, 2, 2, 2, 93, 452, 3, 2, 2, 2, 95, 463, 3, 2, 2,
	2, 97, 472, 3, 2, 2, 2, 99, 479, 3, 2, 2, 2, 101, 489, 3, 2, 2, 2, 103,
	500, 3, 2, 2, 2, 105, 508, 3, 2, 2, 2, 107, 516, 3, 2, 2, 2, 109, 531,
	3, 2, 2, 2, 111, 533, 3, 2, 2, 2, 113, 542, 3, 2, 2, 2, 115, 562, 3, 2,
	2, 2, 117, 118, 7, 100, 2, 2, 118, 119, 7, 103, 2, 2, 119, 120, 7, 105,
	2, 2, 120, 121, 7, 107, 2, 2, 121, 122, 7, 112, 2, 2, 122, 123, 7, 104,
	2, 2, 123, 124, 7, 107, 2, 2, 124, 125, 7, 105, 2, 2, 125, 4, 3, 2, 2,
	2, 126, 127, 7, 103, 2, 2, 127, 128, 7, 112, 2, 2, 128, 129, 7, 102, 2,
	2, 129, 130, 7, 104, 2, 2, 130, 131, 7, 107, 2, 2, 131, 132, 7, 105, 2,
	2, 132, 6, 3, 2, 2, 2, 133, 134, 7, 117, 2, 2, 134, 135, 7, 99, 2, 2, 135,
	136, 7, 120, 2, 2, 136, 137, 7, 103, 2, 2, 137, 8, 3, 2, 2, 2, 138, 139,
	7, 117, 2, 2, 139, 140, 7, 106, 2, 2, 140, 141, 7, 113, 2, 2, 141, 142,
	7, 121, 2, 2, 142, 143, 7, 120, 2, 2, 143, 144, 7, 99, 2, 2, 144, 145,
	7, 116, 2, 2, 145, 146, 7, 107, 2, 2, 146, 147, 7, 99, 2, 2, 147, 148,
	7, 100, 2, 2, 148, 149, 7, 110, 2, 2, 149, 150, 7, 103, 2, 2, 150, 10,
	3, 2, 2, 2, 151, 152, 7, 102, 2, 2, 152, 153, 7, 116, 2, 2, 153, 154, 7,
	99, 2, 2, 154, 155, 7, 121, 2, 2, 155, 12, 3, 2, 2, 2, 156, 157, 7, 104,
	2, 2, 157, 158, 7, 107, 2, 2, 158, 159, 7, 110, 2, 2, 159, 160, 7, 110,
	2, 2, 160, 14, 3, 2, 2, 2, 161, 162, 7, 114, 2, 2, 162, 163, 7, 107, 2,
	2, 163, 164, 7, 101, 2, 2, 164, 165, 7, 109, 2, 2, 165, 166, 7, 119, 2,
	2, 166, 167, 7, 114, 2, 2, 167, 16, 3, 2, 2, 2, 168, 169, 7, 121, 2, 2,
	169, 170, 7, 107, 2, 2, 170, 171, 7, 118, 2, 2, 171, 172, 7, 106, 2, 2,
	172, 173, 7, 101, 2, 2, 173, 174, 7, 113, 2, 2, 174, 175, 7, 110, 2, 2,
	175, 176, 7, 113, 2, 2, 176, 177, 7, 116, 2, 2, 177, 18, 3, 2, 2, 2, 178,
	179, 7, 60, 2, 2, 179, 180, 7, 63, 2, 2, 180, 20, 3, 2, 2, 2, 181, 182,
	7, 63, 2, 2, 182, 22, 3, 2, 2, 2, 183, 184, 7, 60, 2, 2, 184, 24, 3, 2,
	2, 2, 185, 186, 7, 61, 2, 2, 186, 26, 3, 2, 2, 2, 187, 188, 7, 46, 2, 2,
	188, 28, 3, 2, 2, 2, 189, 190, 7, 42, 2, 2, 190, 30, 3, 2, 2, 2, 191, 192,
	7, 43, 2, 2, 192, 32, 3, 2, 2, 2, 193, 194, 7, 93, 2, 2, 194, 34, 3, 2,
	2, 2, 195, 196, 7, 95, 2, 2, 196, 36, 3, 2, 2, 2, 197, 198, 7, 47, 2, 2,
	198, 199, 7, 47, 2, 2, 199, 201, 3, 2, 2, 2, 200, 202, 7, 47, 2, 2, 201,
	200, 3, 2, 2, 2, 201, 202, 3, 2, 2, 2, 202, 211, 3, 2, 2, 2, 203, 204,
	7, 48, 2, 2, 204, 205, 7, 48, 2, 2, 205, 207, 3, 2, 2, 2, 206, 208, 7,
	48, 2, 2, 207, 206, 3, 2, 2, 2, 207, 208, 3, 2, 2, 2, 208, 211, 3, 2, 2,
	2, 209, 211, 7, 40, 2, 2, 210, 197, 3, 2, 2, 2, 210, 203, 3, 2, 2, 2, 210,
	209, 3, 2, 2, 2, 211, 38, 3, 2, 2, 2, 212, 213, 7, 45, 2, 2, 213, 40, 3,
	2, 2, 2, 214, 215, 7, 47, 2, 2, 215, 42, 3, 2, 2, 2, 216, 217, 7, 44, 2,
	2, 217, 44, 3, 2, 2, 2, 218, 219, 7, 49, 2, 2, 219, 46, 3, 2, 2, 2, 220,
	230, 7, 65, 2, 2, 221, 222, 7, 121, 2, 2, 222, 223, 7, 106, 2, 2, 223,
	224, 7, 99, 2, 2, 224, 225, 7, 118, 2, 2, 225, 226, 7, 103, 2, 2, 226,
	227, 7, 120, 2, 2, 227, 228, 7, 103, 2, 2, 228, 230, 7, 116, 2, 2, 229,
	220, 3, 2, 2, 2, 229, 221, 3, 2, 2, 2, 230, 48, 3, 2, 2, 2, 231, 232, 7,
	111, 2, 2, 232, 242, 7, 111, 2, 2, 233, 234, 7, 101, 2, 2, 234, 242, 7,
	111, 2, 2, 235, 236, 7, 107, 2, 2, 236, 242, 7, 112, 2, 2, 237, 238, 7,
	114, 2, 2, 238, 242, 7, 118, 2, 2, 239, 240, 7, 114, 2, 2, 240, 242, 7,
	101, 2, 2, 241, 231, 3, 2, 2, 2, 241, 233, 3, 2, 2, 2, 241, 235, 3, 2,
	2, 2, 241, 237, 3, 2, 2, 2, 241, 239, 3, 2, 2, 2, 242, 50, 3, 2, 2, 2,
	243, 244, 7, 112, 2, 2, 244, 245, 7, 119, 2, 2, 245, 246, 7, 111, 2, 2,
	246, 247, 7, 103, 2, 2, 247, 248, 7, 116, 2, 2, 248, 249, 7, 107, 2, 2,
	249, 250, 7, 101, 2, 2, 250, 52, 3, 2, 2, 2, 251, 252, 7, 114, 2, 2, 252,
	253, 7, 99, 2, 2, 253, 254, 7, 107, 2, 2, 254, 255, 7, 116, 2, 2, 255,
	54, 3, 2, 2, 2, 256, 257, 7, 114, 2, 2, 257, 258, 7, 99, 2, 2, 258, 259,
	7, 118, 2, 2, 259, 260, 7, 106, 2, 2, 260, 56, 3, 2, 2, 2, 261, 262, 7,
	121, 2, 2, 262, 263, 7, 107, 2, 2, 263, 264, 7, 102, 2, 2, 264, 265, 7,
	118, 2, 2, 265, 273, 7, 106, 2, 2, 266, 267, 7, 106, 2, 2, 267, 268, 7,
	103, 2, 2, 268, 269, 7, 107, 2, 2, 269, 270, 7, 105, 2, 2, 270, 271, 7,
	106, 2, 2, 271, 273, 7, 118, 2, 2, 272, 261, 3, 2, 2, 2, 272, 266, 3, 2,
	2, 2, 273, 58, 3, 2, 2, 2, 274, 275, 7, 122, 2, 2, 275, 276, 7, 114, 2,
	2, 276, 277, 7, 99, 2, 2, 277, 278, 7, 116, 2, 2, 278, 285, 7, 118, 2,
	2, 279, 280, 7, 123, 2, 2, 280, 281, 7, 114, 2, 2, 281, 282, 7, 99, 2,
	2, 282, 283, 7, 116, 2, 2, 283, 285, 7, 118, 2, 2, 284, 274, 3, 2, 2, 2,
	284, 279, 3, 2, 2, 2, 285, 60, 3, 2, 2, 2, 286, 287, 7, 110, 2, 2, 287,
	288, 7, 103, 2, 2, 288, 289, 7, 112, 2, 2, 289, 290, 7, 105, 2, 2, 290,
	291, 7, 118, 2, 2, 291, 292, 7, 106, 2, 2, 292, 62, 3, 2, 2, 2, 293, 294,
	7, 104, 2, 2, 294, 295, 7, 110, 2, 2, 295, 296, 7, 113, 2, 2, 296, 297,
	7, 113, 2, 2, 297, 307, 7, 116, 2, 2, 298, 299, 7, 101, 2, 2, 299, 300,
	7, 103, 2, 2, 300, 301, 7, 107, 2, 2, 301, 307, 7, 110, 2, 2, 302, 303,
	7, 117, 2, 2, 303, 304, 7, 115, 2, 2, 304, 305, 7, 116, 2, 2, 305, 307,
	7, 118, 2, 2, 306, 293, 3, 2, 2, 2, 306, 298, 3, 2, 2, 2, 306, 302, 3,
	2, 2, 2, 307, 64, 3, 2, 2, 2, 308, 309, 7, 121, 2, 2, 309, 310, 7, 107,
	2, 2, 310, 311, 7, 118, 2, 2, 311, 312, 7, 106, 2, 2, 312, 66, 3, 2, 2,
	2, 313, 314, 7, 114, 2, 2, 314, 315, 7, 113, 2, 2, 315, 316, 7, 107, 2,
	2, 316, 317, 7, 112, 2, 2, 317, 318, 7, 118, 2, 2, 318, 68, 3, 2, 2, 2,
	319, 320, 7, 113, 2, 2, 320, 321, 7, 104, 2, 2, 321, 70, 3, 2, 2, 2, 322,
	323, 7, 117, 2, 2, 323, 324, 7, 101, 2, 2, 324, 325, 7, 99, 2, 2, 325,
	326, 7, 110, 2, 2, 326, 327, 7, 103, 2, 2, 327, 328, 7, 102, 2, 2, 328,
	72, 3, 2, 2, 2, 329, 330, 7, 117, 2, 2, 330, 331, 7, 106, 2, 2, 331, 332,
	7, 107, 2, 2, 332, 333, 7, 104, 2, 2, 333, 334, 7, 118, 2, 2, 334, 335,
	7, 103, 2, 2, 335, 336, 7, 102, 2, 2, 336, 74, 3, 2, 2, 2, 337, 338, 7,
	116, 2, 2, 338, 339, 7, 113, 2, 2, 339, 340, 7, 118, 2, 2, 340, 341, 7,
	99, 2, 2, 341, 342, 7, 118, 2, 2, 342, 343, 7, 103, 2, 2, 343, 344, 7,
	102, 2, 2, 344, 76, 3, 2, 2, 2, 345, 346, 7, 112, 2, 2, 346, 347, 7, 119,
	2, 2, 347, 348, 7, 110, 2, 2, 348, 349, 7, 110, 2, 2, 349, 350, 7, 114,
	2, 2, 350, 351, 7, 99, 2, 2, 351, 352, 7, 118, 2, 2, 352, 353, 7, 106,
	2, 2, 353, 78, 3, 2, 2, 2, 354, 355, 7, 117, 2, 2, 355, 356, 7, 119, 2,
	2, 356, 357, 7, 100, 2, 2, 357, 358, 7, 114, 2, 2, 358, 359, 7, 99, 2,
	2, 359, 360, 7, 118, 2, 2, 360, 361, 7, 106, 2, 2, 361, 80, 3, 2, 2, 2,
	362, 363, 7, 116, 2, 2, 363, 364, 7, 103, 2, 2, 364, 365, 7, 120, 2, 2,
	365, 366, 7, 103, 2, 2, 366, 367, 7, 116, 2, 2, 367, 368, 7, 117, 2, 2,
	368, 369, 7, 103, 2, 2, 369, 82, 3, 2, 2, 2, 370, 371, 7, 101, 2, 2, 371,
	372, 7, 123, 2, 2, 372, 373, 7, 101, 2, 2, 373, 374, 7, 110, 2, 2, 374,
	375, 7, 103, 2, 2, 375, 84, 3, 2, 2, 2, 376, 377, 7, 107, 2, 2, 377, 378,
	7, 112, 2, 2, 378, 379, 7, 118, 2, 2, 379, 380, 7, 103, 2, 2, 380, 381,
	7, 116, 2, 2, 381, 382, 7, 117, 2, 2, 382, 383, 7, 103, 2, 2, 383, 384,
	7, 101, 2, 2, 384, 385, 7, 118, 2, 2, 385, 386, 7, 107, 2, 2, 386, 387,
	7, 113, 2, 2, 387, 388, 7, 112, 2, 2, 388, 389, 7, 114, 2, 2, 389, 390,
	7, 113, 2, 2, 390, 391, 7, 107, 2, 2, 391, 392, 7, 112, 2, 2, 392, 393,
	7, 118, 2, 2, 393, 86, 3, 2, 2, 2, 394, 395, 7, 114, 2, 2, 395, 396, 7,
	103, 2, 2, 396, 397, 7, 112, 2, 2, 397, 398, 7, 101, 2, 2, 398, 399, 7,
	107, 2, 2, 399, 400, 7, 116, 2, 2, 400, 401, 7, 101, 2, 2, 401, 402, 7,
	110, 2, 2, 402, 413, 7, 103, 2, 2, 403, 404, 7, 114, 2, 2, 404, 405, 7,
	103, 2, 2, 405, 406, 7, 112, 2, 2, 406, 407, 7, 117, 2, 2, 407, 408, 7,
	115, 2, 2, 408, 409, 7, 119, 2, 2, 409, 410, 7, 99, 2, 2, 410, 411, 7,
	116, 2, 2, 411, 413, 7, 103, 2, 2, 412, 394, 3, 2, 2, 2, 412, 403, 3, 2,
	2, 2, 413, 88, 3, 2, 2, 2, 414, 415, 7, 37, 2, 2, 415, 416, 9, 2, 2, 2,
	416, 417, 9, 2, 2, 2, 417, 421, 9, 2, 2, 2, 418, 419, 9, 2, 2, 2, 419,
	420, 9, 2, 2, 2, 420, 422, 9, 2, 2, 2, 421, 418, 3, 2, 2, 2, 421, 422,
	3, 2, 2, 2, 422, 90, 3, 2, 2, 2, 423, 424, 7, 119, 2, 2, 424, 425, 7, 112,
	2, 2, 425, 426, 7, 107, 2, 2, 426, 427, 7, 113, 2, 2, 427, 451, 7, 112,
	2, 2, 428, 429, 7, 107, 2, 2, 429, 430, 7, 112, 2, 2, 430, 431, 7, 118,
	2, 2, 431, 432, 7, 103, 2, 2, 432, 433, 7, 116, 2, 2, 433, 434, 7, 117,
	2, 2, 434, 435, 7, 103, 2, 2, 435, 436, 7, 101, 2, 2, 436, 437, 7, 118,
	2, 2, 437, 438, 7, 107, 2, 2, 438, 439, 7, 113, 2, 2, 439, 451, 7, 112,
	2, 2, 440, 441, 7, 102, 2, 2, 441, 442, 7, 107, 2, 2, 442, 443, 7, 104,
	2, 2, 443, 444, 7, 104, 2, 2, 444, 445, 7, 103, 2, 2, 445, 446, 7, 116,
	2, 2, 446, 447, 7, 103, 2, 2, 447, 448, 7, 112, 2, 2, 448, 449, 7, 101,
	2, 2, 449, 451, 7, 103, 2, 2, 450, 423, 3, 2, 2, 2, 450, 428, 3, 2, 2,
	2, 450, 440, 3, 2, 2, 2, 451, 92, 3, 2, 2, 2, 452, 453, 7, 100, 2, 2, 453,
	454, 7, 103, 2, 2, 454, 455, 7, 105, 2, 2, 455, 456, 7, 107, 2, 2, 456,
	457, 7, 112, 2, 2, 457, 458, 7, 105, 2, 2, 458, 459, 7, 116, 2, 2, 459,
	460, 7, 113, 2, 2, 460, 461, 7, 119, 2, 2, 461, 462, 7, 114, 2, 2, 462,
	94, 3, 2, 2, 2, 463, 464, 7, 103, 2, 2, 464, 465, 7, 112, 2, 2, 465, 466,
	7, 102, 2, 2, 466, 467, 7, 105, 2, 2, 467, 468, 7, 116, 2, 2, 468, 469,
	7, 113, 2, 2, 469, 470, 7, 119, 2, 2, 470, 471, 7, 114, 2, 2, 471, 96,
	3, 2, 2, 2, 472, 476, 9, 3, 2, 2, 473, 475, 9, 4, 2, 2, 474, 473, 3, 2,
	2, 2, 475, 478, 3, 2, 2, 2, 476, 474, 3, 2, 2, 2, 476, 477, 3, 2, 2, 2,
	477, 98, 3, 2, 2, 2, 478, 476, 3, 2, 2, 2, 479, 483, 9, 5, 2, 2, 480, 482,
	9, 4, 2, 2, 481, 480, 3, 2, 2, 2, 482, 485, 3, 2, 2, 2, 483, 481, 3, 2,
	2, 2, 483, 484, 3, 2, 2, 2, 484, 100, 3, 2, 2, 2, 485, 483, 3, 2, 2, 2,
	486, 487, 7, 48, 2, 2, 487, 490, 9, 4, 2, 2, 488, 490, 9, 3, 2, 2, 489,
	486, 3, 2, 2, 2, 489, 488, 3, 2, 2, 2, 490, 494, 3, 2, 2, 2, 491, 493,
	9, 6, 2, 2, 492, 491, 3, 2, 2, 2, 493, 496, 3, 2, 2, 2, 494, 492, 3, 2,
	2, 2, 494, 495, 3, 2, 2, 2, 495, 102, 3, 2, 2, 2, 496, 494, 3, 2, 2, 2,
	497, 498, 7, 48, 2, 2, 498, 501, 9, 4, 2, 2, 499, 501, 9, 5, 2, 2, 500,
	497, 3, 2, 2, 2, 500, 499, 3, 2, 2, 2, 501, 505, 3, 2, 2, 2, 502, 504,
	9, 6, 2, 2, 503, 502, 3, 2, 2, 2, 504, 507, 3, 2, 2, 2, 505, 503, 3, 2,
	2, 2, 505, 506, 3, 2, 2, 2, 506, 104, 3, 2, 2, 2, 507, 505, 3, 2, 2, 2,
	508, 512, 9, 7, 2, 2, 509, 511, 9, 8, 2, 2, 510, 509, 3, 2, 2, 2, 511,
	514, 3, 2, 2, 2, 512, 510, 3, 2, 2, 2, 512, 513, 3, 2, 2, 2, 513, 106,
	3, 2, 2, 2, 514, 512, 3, 2, 2, 2, 515, 517, 4, 50, 59, 2, 516, 515, 3,
	2, 2, 2, 517, 518, 3, 2, 2, 2, 518, 516, 3, 2, 2, 2, 518, 519, 3, 2, 2,
	2, 519, 526, 3, 2, 2, 2, 520, 522, 5, 109, 55, 2, 521, 523, 4, 50, 59,
	2, 522, 521, 3, 2, 2, 2, 523, 524, 3, 2, 2, 2, 524, 522, 3, 2, 2, 2, 524,
	525, 3, 2, 2, 2, 525, 527, 3, 2, 2, 2, 526, 520, 3, 2, 2, 2, 526, 527,
	3, 2, 2, 2, 527, 529, 3, 2, 2, 2, 528, 530, 7, 39, 2, 2, 529, 528, 3, 2,
	2, 2, 529, 530, 3, 2, 2, 2, 530, 108, 3, 2, 2, 2, 531, 532, 7, 48, 2, 2,
	532, 110, 3, 2, 2, 2, 533, 534, 7, 36, 2, 2, 534, 536, 9, 9, 2, 2, 535,
	537, 9, 10, 2, 2, 536, 535, 3, 2, 2, 2, 537, 538, 3, 2, 2, 2, 538, 536,
	3, 2, 2, 2, 538, 539, 3, 2, 2, 2, 539, 540, 3, 2, 2, 2, 540, 541, 7, 36,
	2, 2, 541, 112, 3, 2, 2, 2, 542, 543, 7, 49, 2, 2, 543, 544, 7, 49, 2,
	2, 544, 548, 3, 2, 2, 2, 545, 547, 11, 2, 2, 2, 546, 545, 3, 2, 2, 2, 547,
	550, 3, 2, 2, 2, 548, 549, 3, 2, 2, 2, 548, 546, 3, 2, 2, 2, 549, 552,
	3, 2, 2, 2, 550, 548, 3, 2, 2, 2, 551, 553, 7, 15, 2, 2, 552, 551, 3, 2,
	2, 2, 552, 553, 3, 2, 2, 2, 553, 554, 3, 2, 2, 2, 554, 555, 7, 12, 2, 2,
	555, 114, 3, 2, 2, 2, 556, 558, 9, 11, 2, 2, 557, 556, 3, 2, 2, 2, 558,
	559, 3, 2, 2, 2, 559, 557, 3, 2, 2, 2, 559, 560, 3, 2, 2, 2, 560, 563,
	3, 2, 2, 2, 561, 563, 5, 113, 57, 2, 562, 557, 3, 2, 2, 2, 562, 561, 3,
	2, 2, 2, 563, 564, 3, 2, 2, 2, 564, 565, 8, 58, 2, 2, 565, 116, 3, 2, 2,
	2, 32, 2, 201, 207, 210, 229, 241, 272, 284, 306, 412, 421, 450, 476, 483,
	489, 492, 494, 500, 503, 505, 512, 518, 524, 526, 529, 538, 548, 552, 559,
	562, 3, 8, 2, 2,
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
	"", "'beginfig'", "'endfig'", "'save'", "'showvariable'", "'draw'", "'fill'",
	"'pickup'", "'withcolor'", "':='", "'='", "':'", "';'", "','", "'('", "')'",
	"'['", "']'", "", "'+'", "'-'", "'*'", "'/'", "", "", "'numeric'", "'pair'",
	"'path'", "", "", "'length'", "", "'with'", "'point'", "'of'", "'scaled'",
	"'shifted'", "'rotated'", "'nullpath'", "'subpath'", "'reverse'", "'cycle'",
	"'intersectionpoint'", "", "", "", "'begingroup'", "'endgroup'", "", "",
	"", "", "", "", "'.'",
}

var lexerSymbolicNames = []string{
	"", "", "", "", "", "", "", "", "", "ASSIGN", "EQUALS", "COLON", "SEMIC",
	"COMMA", "LPAREN", "RPAREN", "LBRACKET", "RBRACKET", "PATHJOIN", "PLUS",
	"MINUS", "TIMES", "OVER", "WHATEVER", "UNIT", "NUMERIC", "PAIR", "PATH",
	"INTERNAL", "PAIRPART", "LENGTH", "MATHFUNC", "WITH", "POINT", "OF", "SCALED",
	"SHIFTED", "ROTATED", "NULLPATH", "SUBPATH", "REVERSE", "CYCLE", "INTERSECTIONPOINT",
	"PEN", "COLOR", "PATHCLIPOP", "BEGINGROUP", "ENDGROUP", "PTAG", "TAG",
	"MIXEDPTAG", "MIXEDTAG", "PATHTAG", "DECIMALTOKEN", "DOT", "LABEL", "WS",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "T__2", "T__3", "T__4", "T__5", "T__6", "T__7", "ASSIGN",
	"EQUALS", "COLON", "SEMIC", "COMMA", "LPAREN", "RPAREN", "LBRACKET", "RBRACKET",
	"PATHJOIN", "PLUS", "MINUS", "TIMES", "OVER", "WHATEVER", "UNIT", "NUMERIC",
	"PAIR", "PATH", "INTERNAL", "PAIRPART", "LENGTH", "MATHFUNC", "WITH", "POINT",
	"OF", "SCALED", "SHIFTED", "ROTATED", "NULLPATH", "SUBPATH", "REVERSE",
	"CYCLE", "INTERSECTIONPOINT", "PEN", "COLOR", "PATHCLIPOP", "BEGINGROUP",
	"ENDGROUP", "PTAG", "TAG", "MIXEDPTAG", "MIXEDTAG", "PATHTAG", "DECIMALTOKEN",
	"DOT", "LABEL", "LINECOMMENT", "WS",
}

type PMMPStatemLexer struct {
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

func NewPMMPStatemLexer(input antlr.CharStream) *PMMPStatemLexer {

	l := new(PMMPStatemLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "PMMPStatem.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// PMMPStatemLexer tokens.
const (
	PMMPStatemLexerT__0              = 1
	PMMPStatemLexerT__1              = 2
	PMMPStatemLexerT__2              = 3
	PMMPStatemLexerT__3              = 4
	PMMPStatemLexerT__4              = 5
	PMMPStatemLexerT__5              = 6
	PMMPStatemLexerT__6              = 7
	PMMPStatemLexerT__7              = 8
	PMMPStatemLexerASSIGN            = 9
	PMMPStatemLexerEQUALS            = 10
	PMMPStatemLexerCOLON             = 11
	PMMPStatemLexerSEMIC             = 12
	PMMPStatemLexerCOMMA             = 13
	PMMPStatemLexerLPAREN            = 14
	PMMPStatemLexerRPAREN            = 15
	PMMPStatemLexerLBRACKET          = 16
	PMMPStatemLexerRBRACKET          = 17
	PMMPStatemLexerPATHJOIN          = 18
	PMMPStatemLexerPLUS              = 19
	PMMPStatemLexerMINUS             = 20
	PMMPStatemLexerTIMES             = 21
	PMMPStatemLexerOVER              = 22
	PMMPStatemLexerWHATEVER          = 23
	PMMPStatemLexerUNIT              = 24
	PMMPStatemLexerNUMERIC           = 25
	PMMPStatemLexerPAIR              = 26
	PMMPStatemLexerPATH              = 27
	PMMPStatemLexerINTERNAL          = 28
	PMMPStatemLexerPAIRPART          = 29
	PMMPStatemLexerLENGTH            = 30
	PMMPStatemLexerMATHFUNC          = 31
	PMMPStatemLexerWITH              = 32
	PMMPStatemLexerPOINT             = 33
	PMMPStatemLexerOF                = 34
	PMMPStatemLexerSCALED            = 35
	PMMPStatemLexerSHIFTED           = 36
	PMMPStatemLexerROTATED           = 37
	PMMPStatemLexerNULLPATH          = 38
	PMMPStatemLexerSUBPATH           = 39
	PMMPStatemLexerREVERSE           = 40
	PMMPStatemLexerCYCLE             = 41
	PMMPStatemLexerINTERSECTIONPOINT = 42
	PMMPStatemLexerPEN               = 43
	PMMPStatemLexerCOLOR             = 44
	PMMPStatemLexerPATHCLIPOP        = 45
	PMMPStatemLexerBEGINGROUP        = 46
	PMMPStatemLexerENDGROUP          = 47
	PMMPStatemLexerPTAG              = 48
	PMMPStatemLexerTAG               = 49
	PMMPStatemLexerMIXEDPTAG         = 50
	PMMPStatemLexerMIXEDTAG          = 51
	PMMPStatemLexerPATHTAG           = 52
	PMMPStatemLexerDECIMALTOKEN      = 53
	PMMPStatemLexerDOT               = 54
	PMMPStatemLexerLABEL             = 55
	PMMPStatemLexerWS                = 56
)
