// Code generated from CG.g4 by ANTLR 4.7.1. DO NOT EDIT.

package cg // CG
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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 3, 215, 423,
	4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9, 7,
	4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4, 13,
	9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 4, 16, 9, 16, 4, 17, 9, 17, 4, 18, 9,
	18, 4, 19, 9, 19, 4, 20, 9, 20, 4, 21, 9, 21, 4, 22, 9, 22, 4, 23, 9, 23,
	4, 24, 9, 24, 4, 25, 9, 25, 4, 26, 9, 26, 4, 27, 9, 27, 4, 28, 9, 28, 4,
	29, 9, 29, 4, 30, 9, 30, 4, 31, 9, 31, 4, 32, 9, 32, 4, 33, 9, 33, 4, 34,
	9, 34, 4, 35, 9, 35, 4, 36, 9, 36, 4, 37, 9, 37, 4, 38, 9, 38, 4, 39, 9,
	39, 3, 2, 3, 2, 3, 2, 3, 2, 5, 2, 83, 10, 2, 3, 2, 3, 2, 3, 2, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 5, 3, 94, 10, 3, 3, 4, 3, 4, 5, 4, 98, 10, 4,
	3, 5, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3, 7, 3, 7, 3, 7,
	3, 7, 3, 7, 7, 7, 114, 10, 7, 12, 7, 14, 7, 117, 11, 7, 3, 8, 3, 8, 3,
	8, 3, 8, 3, 8, 3, 8, 3, 8, 3, 8, 5, 8, 127, 10, 8, 3, 9, 3, 9, 3, 10, 3,
	10, 7, 10, 133, 10, 10, 12, 10, 14, 10, 136, 11, 10, 3, 10, 7, 10, 139,
	10, 10, 12, 10, 14, 10, 142, 11, 10, 3, 11, 3, 11, 5, 11, 146, 10, 11,
	3, 11, 3, 11, 3, 11, 5, 11, 151, 10, 11, 3, 11, 3, 11, 7, 11, 155, 10,
	11, 12, 11, 14, 11, 158, 11, 11, 3, 12, 3, 12, 3, 12, 3, 12, 3, 12, 3,
	12, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 3, 13, 7, 13, 172, 10, 13, 12, 13,
	14, 13, 175, 11, 13, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 14, 3, 15, 3,
	15, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 3, 16, 5, 16, 192, 10, 16,
	3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 3, 17, 5, 17, 201, 10, 17, 3,
	18, 3, 18, 6, 18, 205, 10, 18, 13, 18, 14, 18, 206, 3, 19, 5, 19, 210,
	10, 19, 3, 19, 3, 19, 5, 19, 214, 10, 19, 3, 20, 3, 20, 7, 20, 218, 10,
	20, 12, 20, 14, 20, 221, 11, 20, 3, 20, 5, 20, 224, 10, 20, 3, 21, 3, 21,
	3, 21, 6, 21, 229, 10, 21, 13, 21, 14, 21, 230, 3, 22, 3, 22, 3, 22, 3,
	22, 3, 22, 3, 22, 3, 22, 5, 22, 240, 10, 22, 3, 22, 3, 22, 3, 22, 3, 22,
	3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3,
	22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22,
	3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3,
	22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22,
	3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3,
	22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22,
	3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3,
	22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22, 3, 22,
	5, 22, 330, 10, 22, 3, 23, 3, 23, 3, 23, 3, 23, 7, 23, 336, 10, 23, 12,
	23, 14, 23, 339, 11, 23, 3, 24, 3, 24, 3, 25, 3, 25, 3, 26, 3, 26, 3, 27,
	3, 27, 3, 28, 3, 28, 3, 29, 3, 29, 3, 30, 3, 30, 3, 31, 3, 31, 3, 32, 3,
	32, 3, 33, 3, 33, 3, 33, 3, 33, 3, 33, 5, 33, 364, 10, 33, 3, 34, 3, 34,
	3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 3, 34, 5, 34, 375, 10, 34, 3,
	35, 3, 35, 3, 35, 3, 35, 3, 35, 3, 35, 5, 35, 383, 10, 35, 3, 36, 3, 36,
	3, 36, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3, 37, 3,
	37, 3, 37, 3, 37, 3, 37, 3, 37, 5, 37, 402, 10, 37, 3, 38, 3, 38, 3, 38,
	3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3, 38, 3,
	38, 3, 38, 5, 38, 419, 10, 38, 3, 39, 3, 39, 3, 39, 2, 2, 40, 2, 4, 6,
	8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42,
	44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 2,
	19, 4, 2, 17, 17, 19, 19, 3, 2, 22, 31, 3, 2, 23, 31, 3, 2, 7, 97, 3, 2,
	99, 100, 3, 2, 6, 7, 3, 2, 9, 65, 3, 2, 67, 100, 3, 2, 165, 194, 3, 2,
	133, 164, 3, 2, 196, 207, 3, 2, 101, 132, 3, 2, 209, 210, 3, 2, 117, 164,
	3, 2, 212, 214, 3, 2, 101, 116, 3, 2, 101, 164, 2, 439, 2, 78, 3, 2, 2,
	2, 4, 87, 3, 2, 2, 2, 6, 97, 3, 2, 2, 2, 8, 99, 3, 2, 2, 2, 10, 103, 3,
	2, 2, 2, 12, 107, 3, 2, 2, 2, 14, 118, 3, 2, 2, 2, 16, 128, 3, 2, 2, 2,
	18, 130, 3, 2, 2, 2, 20, 145, 3, 2, 2, 2, 22, 159, 3, 2, 2, 2, 24, 165,
	3, 2, 2, 2, 26, 176, 3, 2, 2, 2, 28, 182, 3, 2, 2, 2, 30, 191, 3, 2, 2,
	2, 32, 200, 3, 2, 2, 2, 34, 204, 3, 2, 2, 2, 36, 209, 3, 2, 2, 2, 38, 223,
	3, 2, 2, 2, 40, 225, 3, 2, 2, 2, 42, 232, 3, 2, 2, 2, 44, 337, 3, 2, 2,
	2, 46, 340, 3, 2, 2, 2, 48, 342, 3, 2, 2, 2, 50, 344, 3, 2, 2, 2, 52, 346,
	3, 2, 2, 2, 54, 348, 3, 2, 2, 2, 56, 350, 3, 2, 2, 2, 58, 352, 3, 2, 2,
	2, 60, 354, 3, 2, 2, 2, 62, 356, 3, 2, 2, 2, 64, 363, 3, 2, 2, 2, 66, 374,
	3, 2, 2, 2, 68, 382, 3, 2, 2, 2, 70, 384, 3, 2, 2, 2, 72, 401, 3, 2, 2,
	2, 74, 418, 3, 2, 2, 2, 76, 420, 3, 2, 2, 2, 78, 82, 5, 44, 23, 2, 79,
	80, 5, 6, 4, 2, 80, 81, 5, 44, 23, 2, 81, 83, 3, 2, 2, 2, 82, 79, 3, 2,
	2, 2, 82, 83, 3, 2, 2, 2, 83, 84, 3, 2, 2, 2, 84, 85, 5, 4, 3, 2, 85, 86,
	5, 44, 23, 2, 86, 3, 3, 2, 2, 2, 87, 93, 5, 12, 7, 2, 88, 89, 5, 44, 23,
	2, 89, 90, 7, 32, 2, 2, 90, 91, 5, 44, 23, 2, 91, 92, 5, 20, 11, 2, 92,
	94, 3, 2, 2, 2, 93, 88, 3, 2, 2, 2, 93, 94, 3, 2, 2, 2, 94, 5, 3, 2, 2,
	2, 95, 98, 5, 8, 5, 2, 96, 98, 5, 10, 6, 2, 97, 95, 3, 2, 2, 2, 97, 96,
	3, 2, 2, 2, 98, 7, 3, 2, 2, 2, 99, 100, 7, 35, 2, 2, 100, 101, 7, 35, 2,
	2, 101, 102, 7, 35, 2, 2, 102, 9, 3, 2, 2, 2, 103, 104, 7, 34, 2, 2, 104,
	105, 7, 34, 2, 2, 105, 106, 7, 34, 2, 2, 106, 11, 3, 2, 2, 2, 107, 115,
	5, 14, 8, 2, 108, 109, 5, 44, 23, 2, 109, 110, 7, 17, 2, 2, 110, 111, 5,
	44, 23, 2, 111, 112, 5, 14, 8, 2, 112, 114, 3, 2, 2, 2, 113, 108, 3, 2,
	2, 2, 114, 117, 3, 2, 2, 2, 115, 113, 3, 2, 2, 2, 115, 116, 3, 2, 2, 2,
	116, 13, 3, 2, 2, 2, 117, 115, 3, 2, 2, 2, 118, 126, 5, 16, 9, 2, 119,
	120, 5, 44, 23, 2, 120, 121, 7, 98, 2, 2, 121, 122, 5, 44, 23, 2, 122,
	123, 5, 18, 10, 2, 123, 124, 5, 44, 23, 2, 124, 125, 7, 98, 2, 2, 125,
	127, 3, 2, 2, 2, 126, 119, 3, 2, 2, 2, 126, 127, 3, 2, 2, 2, 127, 15, 3,
	2, 2, 2, 128, 129, 5, 42, 22, 2, 129, 17, 3, 2, 2, 2, 130, 140, 5, 64,
	33, 2, 131, 133, 5, 46, 24, 2, 132, 131, 3, 2, 2, 2, 133, 136, 3, 2, 2,
	2, 134, 132, 3, 2, 2, 2, 134, 135, 3, 2, 2, 2, 135, 137, 3, 2, 2, 2, 136,
	134, 3, 2, 2, 2, 137, 139, 5, 64, 33, 2, 138, 134, 3, 2, 2, 2, 139, 142,
	3, 2, 2, 2, 140, 138, 3, 2, 2, 2, 140, 141, 3, 2, 2, 2, 141, 19, 3, 2,
	2, 2, 142, 140, 3, 2, 2, 2, 143, 146, 5, 24, 13, 2, 144, 146, 5, 22, 12,
	2, 145, 143, 3, 2, 2, 2, 145, 144, 3, 2, 2, 2, 146, 156, 3, 2, 2, 2, 147,
	150, 5, 44, 23, 2, 148, 149, 7, 18, 2, 2, 149, 151, 5, 44, 23, 2, 150,
	148, 3, 2, 2, 2, 150, 151, 3, 2, 2, 2, 151, 152, 3, 2, 2, 2, 152, 153,
	5, 22, 12, 2, 153, 155, 3, 2, 2, 2, 154, 147, 3, 2, 2, 2, 155, 158, 3,
	2, 2, 2, 156, 154, 3, 2, 2, 2, 156, 157, 3, 2, 2, 2, 157, 21, 3, 2, 2,
	2, 158, 156, 3, 2, 2, 2, 159, 160, 7, 97, 2, 2, 160, 161, 5, 44, 23, 2,
	161, 162, 5, 24, 13, 2, 162, 163, 5, 44, 23, 2, 163, 164, 7, 99, 2, 2,
	164, 23, 3, 2, 2, 2, 165, 173, 5, 26, 14, 2, 166, 167, 5, 44, 23, 2, 167,
	168, 7, 18, 2, 2, 168, 169, 5, 44, 23, 2, 169, 170, 5, 26, 14, 2, 170,
	172, 3, 2, 2, 2, 171, 166, 3, 2, 2, 2, 172, 175, 3, 2, 2, 2, 173, 171,
	3, 2, 2, 2, 173, 174, 3, 2, 2, 2, 174, 25, 3, 2, 2, 2, 175, 173, 3, 2,
	2, 2, 176, 177, 5, 28, 15, 2, 177, 178, 5, 44, 23, 2, 178, 179, 7, 35,
	2, 2, 179, 180, 5, 44, 23, 2, 180, 181, 5, 30, 16, 2, 181, 27, 3, 2, 2,
	2, 182, 183, 5, 14, 8, 2, 183, 29, 3, 2, 2, 2, 184, 192, 5, 32, 17, 2,
	185, 186, 5, 54, 28, 2, 186, 187, 5, 34, 18, 2, 187, 188, 5, 54, 28, 2,
	188, 192, 3, 2, 2, 2, 189, 190, 7, 9, 2, 2, 190, 192, 5, 36, 19, 2, 191,
	184, 3, 2, 2, 2, 191, 185, 3, 2, 2, 2, 191, 189, 3, 2, 2, 2, 192, 31, 3,
	2, 2, 2, 193, 201, 5, 14, 8, 2, 194, 195, 7, 14, 2, 2, 195, 196, 5, 44,
	23, 2, 196, 197, 5, 4, 3, 2, 197, 198, 5, 44, 23, 2, 198, 199, 7, 15, 2,
	2, 199, 201, 3, 2, 2, 2, 200, 193, 3, 2, 2, 2, 200, 194, 3, 2, 2, 2, 201,
	33, 3, 2, 2, 2, 202, 205, 5, 66, 34, 2, 203, 205, 5, 68, 35, 2, 204, 202,
	3, 2, 2, 2, 204, 203, 3, 2, 2, 2, 205, 206, 3, 2, 2, 2, 206, 204, 3, 2,
	2, 2, 206, 207, 3, 2, 2, 2, 207, 35, 3, 2, 2, 2, 208, 210, 9, 2, 2, 2,
	209, 208, 3, 2, 2, 2, 209, 210, 3, 2, 2, 2, 210, 213, 3, 2, 2, 2, 211,
	214, 5, 40, 21, 2, 212, 214, 5, 38, 20, 2, 213, 211, 3, 2, 2, 2, 213, 212,
	3, 2, 2, 2, 214, 37, 3, 2, 2, 2, 215, 219, 5, 62, 32, 2, 216, 218, 5, 58,
	30, 2, 217, 216, 3, 2, 2, 2, 218, 221, 3, 2, 2, 2, 219, 217, 3, 2, 2, 2,
	219, 220, 3, 2, 2, 2, 220, 224, 3, 2, 2, 2, 221, 219, 3, 2, 2, 2, 222,
	224, 5, 60, 31, 2, 223, 215, 3, 2, 2, 2, 223, 222, 3, 2, 2, 2, 224, 39,
	3, 2, 2, 2, 225, 226, 5, 38, 20, 2, 226, 228, 7, 20, 2, 2, 227, 229, 5,
	58, 30, 2, 228, 227, 3, 2, 2, 2, 229, 230, 3, 2, 2, 2, 230, 228, 3, 2,
	2, 2, 230, 231, 3, 2, 2, 2, 231, 41, 3, 2, 2, 2, 232, 233, 5, 62, 32, 2,
	233, 234, 5, 58, 30, 2, 234, 235, 5, 58, 30, 2, 235, 236, 5, 58, 30, 2,
	236, 237, 5, 58, 30, 2, 237, 329, 5, 58, 30, 2, 238, 240, 5, 58, 30, 2,
	239, 238, 3, 2, 2, 2, 239, 240, 3, 2, 2, 2, 240, 330, 3, 2, 2, 2, 241,
	242, 5, 58, 30, 2, 242, 243, 5, 58, 30, 2, 243, 330, 3, 2, 2, 2, 244, 245,
	5, 58, 30, 2, 245, 246, 5, 58, 30, 2, 246, 247, 5, 58, 30, 2, 247, 330,
	3, 2, 2, 2, 248, 249, 5, 58, 30, 2, 249, 250, 5, 58, 30, 2, 250, 251, 5,
	58, 30, 2, 251, 252, 5, 58, 30, 2, 252, 330, 3, 2, 2, 2, 253, 254, 5, 58,
	30, 2, 254, 255, 5, 58, 30, 2, 255, 256, 5, 58, 30, 2, 256, 257, 5, 58,
	30, 2, 257, 258, 5, 58, 30, 2, 258, 330, 3, 2, 2, 2, 259, 260, 5, 58, 30,
	2, 260, 261, 5, 58, 30, 2, 261, 262, 5, 58, 30, 2, 262, 263, 5, 58, 30,
	2, 263, 264, 5, 58, 30, 2, 264, 265, 5, 58, 30, 2, 265, 330, 3, 2, 2, 2,
	266, 267, 5, 58, 30, 2, 267, 268, 5, 58, 30, 2, 268, 269, 5, 58, 30, 2,
	269, 270, 5, 58, 30, 2, 270, 271, 5, 58, 30, 2, 271, 272, 5, 58, 30, 2,
	272, 273, 5, 58, 30, 2, 273, 330, 3, 2, 2, 2, 274, 275, 5, 58, 30, 2, 275,
	276, 5, 58, 30, 2, 276, 277, 5, 58, 30, 2, 277, 278, 5, 58, 30, 2, 278,
	279, 5, 58, 30, 2, 279, 280, 5, 58, 30, 2, 280, 281, 5, 58, 30, 2, 281,
	282, 5, 58, 30, 2, 282, 330, 3, 2, 2, 2, 283, 284, 5, 58, 30, 2, 284, 285,
	5, 58, 30, 2, 285, 286, 5, 58, 30, 2, 286, 287, 5, 58, 30, 2, 287, 288,
	5, 58, 30, 2, 288, 289, 5, 58, 30, 2, 289, 290, 5, 58, 30, 2, 290, 291,
	5, 58, 30, 2, 291, 292, 5, 58, 30, 2, 292, 330, 3, 2, 2, 2, 293, 294, 5,
	58, 30, 2, 294, 295, 5, 58, 30, 2, 295, 296, 5, 58, 30, 2, 296, 297, 5,
	58, 30, 2, 297, 298, 5, 58, 30, 2, 298, 299, 5, 58, 30, 2, 299, 300, 5,
	58, 30, 2, 300, 301, 5, 58, 30, 2, 301, 302, 5, 58, 30, 2, 302, 303, 5,
	58, 30, 2, 303, 330, 3, 2, 2, 2, 304, 305, 5, 58, 30, 2, 305, 306, 5, 58,
	30, 2, 306, 307, 5, 58, 30, 2, 307, 308, 5, 58, 30, 2, 308, 309, 5, 58,
	30, 2, 309, 310, 5, 58, 30, 2, 310, 311, 5, 58, 30, 2, 311, 312, 5, 58,
	30, 2, 312, 313, 5, 58, 30, 2, 313, 314, 5, 58, 30, 2, 314, 315, 5, 58,
	30, 2, 315, 330, 3, 2, 2, 2, 316, 317, 5, 58, 30, 2, 317, 318, 5, 58, 30,
	2, 318, 319, 5, 58, 30, 2, 319, 320, 5, 58, 30, 2, 320, 321, 5, 58, 30,
	2, 321, 322, 5, 58, 30, 2, 322, 323, 5, 58, 30, 2, 323, 324, 5, 58, 30,
	2, 324, 325, 5, 58, 30, 2, 325, 326, 5, 58, 30, 2, 326, 327, 5, 58, 30,
	2, 327, 328, 5, 58, 30, 2, 328, 330, 3, 2, 2, 2, 329, 239, 3, 2, 2, 2,
	329, 241, 3, 2, 2, 2, 329, 244, 3, 2, 2, 2, 329, 248, 3, 2, 2, 2, 329,
	253, 3, 2, 2, 2, 329, 259, 3, 2, 2, 2, 329, 266, 3, 2, 2, 2, 329, 274,
	3, 2, 2, 2, 329, 283, 3, 2, 2, 2, 329, 293, 3, 2, 2, 2, 329, 304, 3, 2,
	2, 2, 329, 316, 3, 2, 2, 2, 330, 43, 3, 2, 2, 2, 331, 336, 5, 46, 24, 2,
	332, 336, 5, 48, 25, 2, 333, 336, 5, 50, 26, 2, 334, 336, 5, 52, 27, 2,
	335, 331, 3, 2, 2, 2, 335, 332, 3, 2, 2, 2, 335, 333, 3, 2, 2, 2, 335,
	334, 3, 2, 2, 2, 336, 339, 3, 2, 2, 2, 337, 335, 3, 2, 2, 2, 337, 338,
	3, 2, 2, 2, 338, 45, 3, 2, 2, 2, 339, 337, 3, 2, 2, 2, 340, 341, 7, 6,
	2, 2, 341, 47, 3, 2, 2, 2, 342, 343, 7, 3, 2, 2, 343, 49, 3, 2, 2, 2, 344,
	345, 7, 5, 2, 2, 345, 51, 3, 2, 2, 2, 346, 347, 7, 4, 2, 2, 347, 53, 3,
	2, 2, 2, 348, 349, 7, 8, 2, 2, 349, 55, 3, 2, 2, 2, 350, 351, 7, 66, 2,
	2, 351, 57, 3, 2, 2, 2, 352, 353, 9, 3, 2, 2, 353, 59, 3, 2, 2, 2, 354,
	355, 7, 22, 2, 2, 355, 61, 3, 2, 2, 2, 356, 357, 9, 4, 2, 2, 357, 63, 3,
	2, 2, 2, 358, 364, 9, 5, 2, 2, 359, 364, 9, 6, 2, 2, 360, 364, 5, 70, 36,
	2, 361, 364, 5, 72, 37, 2, 362, 364, 5, 74, 38, 2, 363, 358, 3, 2, 2, 2,
	363, 359, 3, 2, 2, 2, 363, 360, 3, 2, 2, 2, 363, 361, 3, 2, 2, 2, 363,
	362, 3, 2, 2, 2, 364, 65, 3, 2, 2, 2, 365, 375, 5, 48, 25, 2, 366, 375,
	5, 50, 26, 2, 367, 375, 5, 52, 27, 2, 368, 375, 9, 7, 2, 2, 369, 375, 9,
	8, 2, 2, 370, 375, 9, 9, 2, 2, 371, 375, 5, 70, 36, 2, 372, 375, 5, 72,
	37, 2, 373, 375, 5, 74, 38, 2, 374, 365, 3, 2, 2, 2, 374, 366, 3, 2, 2,
	2, 374, 367, 3, 2, 2, 2, 374, 368, 3, 2, 2, 2, 374, 369, 3, 2, 2, 2, 374,
	370, 3, 2, 2, 2, 374, 371, 3, 2, 2, 2, 374, 372, 3, 2, 2, 2, 374, 373,
	3, 2, 2, 2, 375, 67, 3, 2, 2, 2, 376, 377, 5, 56, 29, 2, 377, 378, 5, 54,
	28, 2, 378, 383, 3, 2, 2, 2, 379, 380, 5, 56, 29, 2, 380, 381, 5, 56, 29,
	2, 381, 383, 3, 2, 2, 2, 382, 376, 3, 2, 2, 2, 382, 379, 3, 2, 2, 2, 383,
	69, 3, 2, 2, 2, 384, 385, 9, 10, 2, 2, 385, 386, 5, 76, 39, 2, 386, 71,
	3, 2, 2, 2, 387, 388, 7, 195, 2, 2, 388, 389, 9, 11, 2, 2, 389, 402, 5,
	76, 39, 2, 390, 391, 9, 12, 2, 2, 391, 392, 5, 76, 39, 2, 392, 393, 5,
	76, 39, 2, 393, 402, 3, 2, 2, 2, 394, 395, 7, 208, 2, 2, 395, 396, 9, 13,
	2, 2, 396, 402, 5, 76, 39, 2, 397, 398, 9, 14, 2, 2, 398, 399, 5, 76, 39,
	2, 399, 400, 5, 76, 39, 2, 400, 402, 3, 2, 2, 2, 401, 387, 3, 2, 2, 2,
	401, 390, 3, 2, 2, 2, 401, 394, 3, 2, 2, 2, 401, 397, 3, 2, 2, 2, 402,
	73, 3, 2, 2, 2, 403, 404, 7, 211, 2, 2, 404, 405, 9, 15, 2, 2, 405, 406,
	5, 76, 39, 2, 406, 407, 5, 76, 39, 2, 407, 419, 3, 2, 2, 2, 408, 409, 9,
	16, 2, 2, 409, 410, 5, 76, 39, 2, 410, 411, 5, 76, 39, 2, 411, 412, 5,
	76, 39, 2, 412, 419, 3, 2, 2, 2, 413, 414, 7, 215, 2, 2, 414, 415, 9, 17,
	2, 2, 415, 416, 5, 76, 39, 2, 416, 417, 5, 76, 39, 2, 417, 419, 3, 2, 2,
	2, 418, 403, 3, 2, 2, 2, 418, 408, 3, 2, 2, 2, 418, 413, 3, 2, 2, 2, 419,
	75, 3, 2, 2, 2, 420, 421, 9, 18, 2, 2, 421, 77, 3, 2, 2, 2, 31, 82, 93,
	97, 115, 126, 134, 140, 145, 150, 156, 173, 191, 200, 204, 206, 209, 213,
	219, 223, 230, 239, 329, 335, 337, 363, 374, 382, 401, 418,
}
var deserializer = antlr.NewATNDeserializer(nil)
var deserializedATN = deserializer.DeserializeFromUInt16(parserATN)

var literalNames = []string{
	"", "'\u0009'", "'\u000A'", "'\u000D'", "' '", "'!'", "'\"'", "'#'", "'$'",
	"'%'", "'&'", "'''", "'('", "')'", "'*'", "'+'", "','", "'-'", "'.'", "'/'",
	"'0'", "'1'", "'2'", "'3'", "'4'", "'5'", "'6'", "'7'", "'8'", "'9'", "':'",
	"';'", "'<'", "'='", "'>'", "'?'", "'@'", "'A'", "'B'", "'C'", "'D'", "'E'",
	"'F'", "'G'", "'H'", "'I'", "'J'", "'K'", "'L'", "'M'", "'N'", "'O'", "'P'",
	"'Q'", "'R'", "'S'", "'T'", "'U'", "'V'", "'W'", "'X'", "'Y'", "'Z'", "'['",
	"'\\'", "']'", "'^'", "'_'", "'`'", "'a'", "'b'", "'c'", "'d'", "'e'",
	"'f'", "'g'", "'h'", "'i'", "'j'", "'k'", "'l'", "'m'", "'n'", "'o'", "'p'",
	"'q'", "'r'", "'s'", "'t'", "'u'", "'v'", "'w'", "'x'", "'y'", "'z'", "'{'",
	"'|'", "'}'", "'~'", "'\u0080'", "'\u0081'", "'\u0082'", "'\u0083'", "'\u0084'",
	"'\u0085'", "'\u0086'", "'\u0087'", "'\u0088'", "'\u0089'", "'\u008A'",
	"'\u008B'", "'\u008C'", "'\u008D'", "'\u008E'", "'\u008F'", "'\u0090'",
	"'\u0091'", "'\u0092'", "'\u0093'", "'\u0094'", "'\u0095'", "'\u0096'",
	"'\u0097'", "'\u0098'", "'\u0099'", "'\u009A'", "'\u009B'", "'\u009C'",
	"'\u009D'", "'\u009E'", "'\u009F'", "'\u00A0'", "'\u00A1'", "'\u00A2'",
	"'\u00A3'", "'\u00A4'", "'\u00A5'", "'\u00A6'", "'\u00A7'", "'\u00A8'",
	"'\u00A9'", "'\u00AA'", "'\u00AB'", "'\u00AC'", "'\u00AD'", "'\u00AE'",
	"'\u00AF'", "'\u00B0'", "'\u00B1'", "'\u00B2'", "'\u00B3'", "'\u00B4'",
	"'\u00B5'", "'\u00B6'", "'\u00B7'", "'\u00B8'", "'\u00B9'", "'\u00BA'",
	"'\u00BB'", "'\u00BC'", "'\u00BD'", "'\u00BE'", "'\u00BF'", "'\u00C2'",
	"'\u00C3'", "'\u00C4'", "'\u00C5'", "'\u00C6'", "'\u00C7'", "'\u00C8'",
	"'\u00C9'", "'\u00CA'", "'\u00CB'", "'\u00CC'", "'\u00CD'", "'\u00CE'",
	"'\u00CF'", "'\u00D0'", "'\u00D1'", "'\u00D2'", "'\u00D3'", "'\u00D4'",
	"'\u00D5'", "'\u00D6'", "'\u00D7'", "'\u00D8'", "'\u00D9'", "'\u00DA'",
	"'\u00DB'", "'\u00DC'", "'\u00DD'", "'\u00DE'", "'\u00DF'", "'\u00E0'",
	"'\u00E1'", "'\u00E2'", "'\u00E3'", "'\u00E4'", "'\u00E5'", "'\u00E6'",
	"'\u00E7'", "'\u00E8'", "'\u00E9'", "'\u00EA'", "'\u00EB'", "'\u00EC'",
	"'\u00ED'", "'\u00EE'", "'\u00EF'", "'\u00F0'", "'\u00F1'", "'\u00F2'",
	"'\u00F3'", "'\u00F4'",
}
var symbolicNames = []string{
	"", "TAB", "LF", "CR", "SPACE", "EXCLAMATION", "QUOTE", "POUND", "DOLLAR",
	"PERCENT", "AMPERSAND", "APOSTROPHE", "LEFT_PAREN", "RIGHT_PAREN", "ASTERISK",
	"PLUS", "COMMA", "DASH", "PERIOD", "SLASH", "ZERO", "ONE", "TWO", "THREE",
	"FOUR", "FIVE", "SIX", "SEVEN", "EIGHT", "NINE", "COLON", "SEMICOLON",
	"LESS_THAN", "EQUALS", "GREATER_THAN", "QUESTION", "AT", "CAP_A", "CAP_B",
	"CAP_C", "CAP_D", "CAP_E", "CAP_F", "CAP_G", "CAP_H", "CAP_I", "CAP_J",
	"CAP_K", "CAP_L", "CAP_M", "CAP_N", "CAP_O", "CAP_P", "CAP_Q", "CAP_R",
	"CAP_S", "CAP_T", "CAP_U", "CAP_V", "CAP_W", "CAP_X", "CAP_Y", "CAP_Z",
	"LEFT_BRACE", "BACKSLASH", "RIGHT_BRACE", "CARAT", "UNDERSCORE", "ACCENT",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O",
	"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "LEFT_CURLY_BRACE",
	"PIPE", "RIGHT_CURLY_BRACE", "TILDE", "U_0080", "U_0081", "U_0082", "U_0083",
	"U_0084", "U_0085", "U_0086", "U_0087", "U_0088", "U_0089", "U_008A", "U_008B",
	"U_008C", "U_008D", "U_008E", "U_008F", "U_0090", "U_0091", "U_0092", "U_0093",
	"U_0094", "U_0095", "U_0096", "U_0097", "U_0098", "U_0099", "U_009A", "U_009B",
	"U_009C", "U_009D", "U_009E", "U_009F", "U_00A0", "U_00A1", "U_00A2", "U_00A3",
	"U_00A4", "U_00A5", "U_00A6", "U_00A7", "U_00A8", "U_00A9", "U_00AA", "U_00AB",
	"U_00AC", "U_00AD", "U_00AE", "U_00AF", "U_00B0", "U_00B1", "U_00B2", "U_00B3",
	"U_00B4", "U_00B5", "U_00B6", "U_00B7", "U_00B8", "U_00B9", "U_00BA", "U_00BB",
	"U_00BC", "U_00BD", "U_00BE", "U_00BF", "U_00C2", "U_00C3", "U_00C4", "U_00C5",
	"U_00C6", "U_00C7", "U_00C8", "U_00C9", "U_00CA", "U_00CB", "U_00CC", "U_00CD",
	"U_00CE", "U_00CF", "U_00D0", "U_00D1", "U_00D2", "U_00D3", "U_00D4", "U_00D5",
	"U_00D6", "U_00D7", "U_00D8", "U_00D9", "U_00DA", "U_00DB", "U_00DC", "U_00DD",
	"U_00DE", "U_00DF", "U_00E0", "U_00E1", "U_00E2", "U_00E3", "U_00E4", "U_00E5",
	"U_00E6", "U_00E7", "U_00E8", "U_00E9", "U_00EA", "U_00EB", "U_00EC", "U_00ED",
	"U_00EE", "U_00EF", "U_00F0", "U_00F1", "U_00F2", "U_00F3", "U_00F4",
}

var ruleNames = []string{
	"expression", "subexpression", "definitionstatus", "equivalentto", "subtypeof",
	"focusconcept", "conceptreference", "conceptid", "term", "refinement",
	"attributegroup", "attributeset", "attribute", "attributename", "attributevalue",
	"expressionvalue", "stringvalue", "numericvalue", "integervalue", "decimalvalue",
	"sctid", "ws", "sp", "htab", "cr", "lf", "qm", "bs", "digit", "zero", "digitnonzero",
	"nonwsnonpipe", "anynonescapedchar", "escapedchar", "utf8_2", "utf8_3",
	"utf8_4", "utf8_tail",
}
var decisionToDFA = make([]*antlr.DFA, len(deserializedATN.DecisionToState))

func init() {
	for index, ds := range deserializedATN.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

type CGParser struct {
	*antlr.BaseParser
}

func NewCGParser(input antlr.TokenStream) *CGParser {
	this := new(CGParser)

	this.BaseParser = antlr.NewBaseParser(input)

	this.Interpreter = antlr.NewParserATNSimulator(this, deserializedATN, decisionToDFA, antlr.NewPredictionContextCache())
	this.RuleNames = ruleNames
	this.LiteralNames = literalNames
	this.SymbolicNames = symbolicNames
	this.GrammarFileName = "CG.g4"

	return this
}

// CGParser tokens.
const (
	CGParserEOF               = antlr.TokenEOF
	CGParserTAB               = 1
	CGParserLF                = 2
	CGParserCR                = 3
	CGParserSPACE             = 4
	CGParserEXCLAMATION       = 5
	CGParserQUOTE             = 6
	CGParserPOUND             = 7
	CGParserDOLLAR            = 8
	CGParserPERCENT           = 9
	CGParserAMPERSAND         = 10
	CGParserAPOSTROPHE        = 11
	CGParserLEFT_PAREN        = 12
	CGParserRIGHT_PAREN       = 13
	CGParserASTERISK          = 14
	CGParserPLUS              = 15
	CGParserCOMMA             = 16
	CGParserDASH              = 17
	CGParserPERIOD            = 18
	CGParserSLASH             = 19
	CGParserZERO              = 20
	CGParserONE               = 21
	CGParserTWO               = 22
	CGParserTHREE             = 23
	CGParserFOUR              = 24
	CGParserFIVE              = 25
	CGParserSIX               = 26
	CGParserSEVEN             = 27
	CGParserEIGHT             = 28
	CGParserNINE              = 29
	CGParserCOLON             = 30
	CGParserSEMICOLON         = 31
	CGParserLESS_THAN         = 32
	CGParserEQUALS            = 33
	CGParserGREATER_THAN      = 34
	CGParserQUESTION          = 35
	CGParserAT                = 36
	CGParserCAP_A             = 37
	CGParserCAP_B             = 38
	CGParserCAP_C             = 39
	CGParserCAP_D             = 40
	CGParserCAP_E             = 41
	CGParserCAP_F             = 42
	CGParserCAP_G             = 43
	CGParserCAP_H             = 44
	CGParserCAP_I             = 45
	CGParserCAP_J             = 46
	CGParserCAP_K             = 47
	CGParserCAP_L             = 48
	CGParserCAP_M             = 49
	CGParserCAP_N             = 50
	CGParserCAP_O             = 51
	CGParserCAP_P             = 52
	CGParserCAP_Q             = 53
	CGParserCAP_R             = 54
	CGParserCAP_S             = 55
	CGParserCAP_T             = 56
	CGParserCAP_U             = 57
	CGParserCAP_V             = 58
	CGParserCAP_W             = 59
	CGParserCAP_X             = 60
	CGParserCAP_Y             = 61
	CGParserCAP_Z             = 62
	CGParserLEFT_BRACE        = 63
	CGParserBACKSLASH         = 64
	CGParserRIGHT_BRACE       = 65
	CGParserCARAT             = 66
	CGParserUNDERSCORE        = 67
	CGParserACCENT            = 68
	CGParserA                 = 69
	CGParserB                 = 70
	CGParserC                 = 71
	CGParserD                 = 72
	CGParserE                 = 73
	CGParserF                 = 74
	CGParserG                 = 75
	CGParserH                 = 76
	CGParserI                 = 77
	CGParserJ                 = 78
	CGParserK                 = 79
	CGParserL                 = 80
	CGParserM                 = 81
	CGParserN                 = 82
	CGParserO                 = 83
	CGParserP                 = 84
	CGParserQ                 = 85
	CGParserR                 = 86
	CGParserS                 = 87
	CGParserT                 = 88
	CGParserU                 = 89
	CGParserV                 = 90
	CGParserW                 = 91
	CGParserX                 = 92
	CGParserY                 = 93
	CGParserZ                 = 94
	CGParserLEFT_CURLY_BRACE  = 95
	CGParserPIPE              = 96
	CGParserRIGHT_CURLY_BRACE = 97
	CGParserTILDE             = 98
	CGParserU_0080            = 99
	CGParserU_0081            = 100
	CGParserU_0082            = 101
	CGParserU_0083            = 102
	CGParserU_0084            = 103
	CGParserU_0085            = 104
	CGParserU_0086            = 105
	CGParserU_0087            = 106
	CGParserU_0088            = 107
	CGParserU_0089            = 108
	CGParserU_008A            = 109
	CGParserU_008B            = 110
	CGParserU_008C            = 111
	CGParserU_008D            = 112
	CGParserU_008E            = 113
	CGParserU_008F            = 114
	CGParserU_0090            = 115
	CGParserU_0091            = 116
	CGParserU_0092            = 117
	CGParserU_0093            = 118
	CGParserU_0094            = 119
	CGParserU_0095            = 120
	CGParserU_0096            = 121
	CGParserU_0097            = 122
	CGParserU_0098            = 123
	CGParserU_0099            = 124
	CGParserU_009A            = 125
	CGParserU_009B            = 126
	CGParserU_009C            = 127
	CGParserU_009D            = 128
	CGParserU_009E            = 129
	CGParserU_009F            = 130
	CGParserU_00A0            = 131
	CGParserU_00A1            = 132
	CGParserU_00A2            = 133
	CGParserU_00A3            = 134
	CGParserU_00A4            = 135
	CGParserU_00A5            = 136
	CGParserU_00A6            = 137
	CGParserU_00A7            = 138
	CGParserU_00A8            = 139
	CGParserU_00A9            = 140
	CGParserU_00AA            = 141
	CGParserU_00AB            = 142
	CGParserU_00AC            = 143
	CGParserU_00AD            = 144
	CGParserU_00AE            = 145
	CGParserU_00AF            = 146
	CGParserU_00B0            = 147
	CGParserU_00B1            = 148
	CGParserU_00B2            = 149
	CGParserU_00B3            = 150
	CGParserU_00B4            = 151
	CGParserU_00B5            = 152
	CGParserU_00B6            = 153
	CGParserU_00B7            = 154
	CGParserU_00B8            = 155
	CGParserU_00B9            = 156
	CGParserU_00BA            = 157
	CGParserU_00BB            = 158
	CGParserU_00BC            = 159
	CGParserU_00BD            = 160
	CGParserU_00BE            = 161
	CGParserU_00BF            = 162
	CGParserU_00C2            = 163
	CGParserU_00C3            = 164
	CGParserU_00C4            = 165
	CGParserU_00C5            = 166
	CGParserU_00C6            = 167
	CGParserU_00C7            = 168
	CGParserU_00C8            = 169
	CGParserU_00C9            = 170
	CGParserU_00CA            = 171
	CGParserU_00CB            = 172
	CGParserU_00CC            = 173
	CGParserU_00CD            = 174
	CGParserU_00CE            = 175
	CGParserU_00CF            = 176
	CGParserU_00D0            = 177
	CGParserU_00D1            = 178
	CGParserU_00D2            = 179
	CGParserU_00D3            = 180
	CGParserU_00D4            = 181
	CGParserU_00D5            = 182
	CGParserU_00D6            = 183
	CGParserU_00D7            = 184
	CGParserU_00D8            = 185
	CGParserU_00D9            = 186
	CGParserU_00DA            = 187
	CGParserU_00DB            = 188
	CGParserU_00DC            = 189
	CGParserU_00DD            = 190
	CGParserU_00DE            = 191
	CGParserU_00DF            = 192
	CGParserU_00E0            = 193
	CGParserU_00E1            = 194
	CGParserU_00E2            = 195
	CGParserU_00E3            = 196
	CGParserU_00E4            = 197
	CGParserU_00E5            = 198
	CGParserU_00E6            = 199
	CGParserU_00E7            = 200
	CGParserU_00E8            = 201
	CGParserU_00E9            = 202
	CGParserU_00EA            = 203
	CGParserU_00EB            = 204
	CGParserU_00EC            = 205
	CGParserU_00ED            = 206
	CGParserU_00EE            = 207
	CGParserU_00EF            = 208
	CGParserU_00F0            = 209
	CGParserU_00F1            = 210
	CGParserU_00F2            = 211
	CGParserU_00F3            = 212
	CGParserU_00F4            = 213
)

// CGParser rules.
const (
	CGParserRULE_expression        = 0
	CGParserRULE_subexpression     = 1
	CGParserRULE_definitionstatus  = 2
	CGParserRULE_equivalentto      = 3
	CGParserRULE_subtypeof         = 4
	CGParserRULE_focusconcept      = 5
	CGParserRULE_conceptreference  = 6
	CGParserRULE_conceptid         = 7
	CGParserRULE_term              = 8
	CGParserRULE_refinement        = 9
	CGParserRULE_attributegroup    = 10
	CGParserRULE_attributeset      = 11
	CGParserRULE_attribute         = 12
	CGParserRULE_attributename     = 13
	CGParserRULE_attributevalue    = 14
	CGParserRULE_expressionvalue   = 15
	CGParserRULE_stringvalue       = 16
	CGParserRULE_numericvalue      = 17
	CGParserRULE_integervalue      = 18
	CGParserRULE_decimalvalue      = 19
	CGParserRULE_sctid             = 20
	CGParserRULE_ws                = 21
	CGParserRULE_sp                = 22
	CGParserRULE_htab              = 23
	CGParserRULE_cr                = 24
	CGParserRULE_lf                = 25
	CGParserRULE_qm                = 26
	CGParserRULE_bs                = 27
	CGParserRULE_digit             = 28
	CGParserRULE_zero              = 29
	CGParserRULE_digitnonzero      = 30
	CGParserRULE_nonwsnonpipe      = 31
	CGParserRULE_anynonescapedchar = 32
	CGParserRULE_escapedchar       = 33
	CGParserRULE_utf8_2            = 34
	CGParserRULE_utf8_3            = 35
	CGParserRULE_utf8_4            = 36
	CGParserRULE_utf8_tail         = 37
)

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
	p.RuleIndex = CGParserRULE_expression
	return p
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) AllWs() []IWsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWsContext)(nil)).Elem())
	var tst = make([]IWsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWsContext)
		}
	}

	return tst
}

func (s *ExpressionContext) Ws(i int) IWsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWsContext)
}

func (s *ExpressionContext) Subexpression() ISubexpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISubexpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISubexpressionContext)
}

func (s *ExpressionContext) Definitionstatus() IDefinitionstatusContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDefinitionstatusContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDefinitionstatusContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *CGParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, CGParserRULE_expression)
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
		p.SetState(76)
		p.Ws()
	}
	p.SetState(80)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CGParserLESS_THAN || _la == CGParserEQUALS {
		{
			p.SetState(77)
			p.Definitionstatus()
		}
		{
			p.SetState(78)
			p.Ws()
		}

	}
	{
		p.SetState(82)
		p.Subexpression()
	}
	{
		p.SetState(83)
		p.Ws()
	}

	return localctx
}

// ISubexpressionContext is an interface to support dynamic dispatch.
type ISubexpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSubexpressionContext differentiates from other interfaces.
	IsSubexpressionContext()
}

type SubexpressionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySubexpressionContext() *SubexpressionContext {
	var p = new(SubexpressionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_subexpression
	return p
}

func (*SubexpressionContext) IsSubexpressionContext() {}

func NewSubexpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubexpressionContext {
	var p = new(SubexpressionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_subexpression

	return p
}

func (s *SubexpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *SubexpressionContext) Focusconcept() IFocusconceptContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IFocusconceptContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IFocusconceptContext)
}

func (s *SubexpressionContext) AllWs() []IWsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWsContext)(nil)).Elem())
	var tst = make([]IWsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWsContext)
		}
	}

	return tst
}

func (s *SubexpressionContext) Ws(i int) IWsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWsContext)
}

func (s *SubexpressionContext) COLON() antlr.TerminalNode {
	return s.GetToken(CGParserCOLON, 0)
}

func (s *SubexpressionContext) Refinement() IRefinementContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IRefinementContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IRefinementContext)
}

func (s *SubexpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubexpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubexpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterSubexpression(s)
	}
}

func (s *SubexpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitSubexpression(s)
	}
}

func (p *CGParser) Subexpression() (localctx ISubexpressionContext) {
	localctx = NewSubexpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, CGParserRULE_subexpression)

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
		p.SetState(85)
		p.Focusconcept()
	}
	p.SetState(91)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 1, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(86)
			p.Ws()
		}
		{
			p.SetState(87)
			p.Match(CGParserCOLON)
		}
		{
			p.SetState(88)
			p.Ws()
		}
		{
			p.SetState(89)
			p.Refinement()
		}

	}

	return localctx
}

// IDefinitionstatusContext is an interface to support dynamic dispatch.
type IDefinitionstatusContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDefinitionstatusContext differentiates from other interfaces.
	IsDefinitionstatusContext()
}

type DefinitionstatusContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDefinitionstatusContext() *DefinitionstatusContext {
	var p = new(DefinitionstatusContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_definitionstatus
	return p
}

func (*DefinitionstatusContext) IsDefinitionstatusContext() {}

func NewDefinitionstatusContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DefinitionstatusContext {
	var p = new(DefinitionstatusContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_definitionstatus

	return p
}

func (s *DefinitionstatusContext) GetParser() antlr.Parser { return s.parser }

func (s *DefinitionstatusContext) Equivalentto() IEquivalenttoContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEquivalenttoContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IEquivalenttoContext)
}

func (s *DefinitionstatusContext) Subtypeof() ISubtypeofContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISubtypeofContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISubtypeofContext)
}

func (s *DefinitionstatusContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DefinitionstatusContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DefinitionstatusContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterDefinitionstatus(s)
	}
}

func (s *DefinitionstatusContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitDefinitionstatus(s)
	}
}

func (p *CGParser) Definitionstatus() (localctx IDefinitionstatusContext) {
	localctx = NewDefinitionstatusContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, CGParserRULE_definitionstatus)

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

	p.SetState(95)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CGParserEQUALS:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(93)
			p.Equivalentto()
		}

	case CGParserLESS_THAN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(94)
			p.Subtypeof()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEquivalenttoContext is an interface to support dynamic dispatch.
type IEquivalenttoContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEquivalenttoContext differentiates from other interfaces.
	IsEquivalenttoContext()
}

type EquivalenttoContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEquivalenttoContext() *EquivalenttoContext {
	var p = new(EquivalenttoContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_equivalentto
	return p
}

func (*EquivalenttoContext) IsEquivalenttoContext() {}

func NewEquivalenttoContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EquivalenttoContext {
	var p = new(EquivalenttoContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_equivalentto

	return p
}

func (s *EquivalenttoContext) GetParser() antlr.Parser { return s.parser }

func (s *EquivalenttoContext) AllEQUALS() []antlr.TerminalNode {
	return s.GetTokens(CGParserEQUALS)
}

func (s *EquivalenttoContext) EQUALS(i int) antlr.TerminalNode {
	return s.GetToken(CGParserEQUALS, i)
}

func (s *EquivalenttoContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EquivalenttoContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EquivalenttoContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterEquivalentto(s)
	}
}

func (s *EquivalenttoContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitEquivalentto(s)
	}
}

func (p *CGParser) Equivalentto() (localctx IEquivalenttoContext) {
	localctx = NewEquivalenttoContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, CGParserRULE_equivalentto)

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
		p.SetState(97)
		p.Match(CGParserEQUALS)
	}
	{
		p.SetState(98)
		p.Match(CGParserEQUALS)
	}
	{
		p.SetState(99)
		p.Match(CGParserEQUALS)
	}

	return localctx
}

// ISubtypeofContext is an interface to support dynamic dispatch.
type ISubtypeofContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSubtypeofContext differentiates from other interfaces.
	IsSubtypeofContext()
}

type SubtypeofContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySubtypeofContext() *SubtypeofContext {
	var p = new(SubtypeofContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_subtypeof
	return p
}

func (*SubtypeofContext) IsSubtypeofContext() {}

func NewSubtypeofContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubtypeofContext {
	var p = new(SubtypeofContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_subtypeof

	return p
}

func (s *SubtypeofContext) GetParser() antlr.Parser { return s.parser }

func (s *SubtypeofContext) AllLESS_THAN() []antlr.TerminalNode {
	return s.GetTokens(CGParserLESS_THAN)
}

func (s *SubtypeofContext) LESS_THAN(i int) antlr.TerminalNode {
	return s.GetToken(CGParserLESS_THAN, i)
}

func (s *SubtypeofContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubtypeofContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubtypeofContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterSubtypeof(s)
	}
}

func (s *SubtypeofContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitSubtypeof(s)
	}
}

func (p *CGParser) Subtypeof() (localctx ISubtypeofContext) {
	localctx = NewSubtypeofContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, CGParserRULE_subtypeof)

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
		p.SetState(101)
		p.Match(CGParserLESS_THAN)
	}
	{
		p.SetState(102)
		p.Match(CGParserLESS_THAN)
	}
	{
		p.SetState(103)
		p.Match(CGParserLESS_THAN)
	}

	return localctx
}

// IFocusconceptContext is an interface to support dynamic dispatch.
type IFocusconceptContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFocusconceptContext differentiates from other interfaces.
	IsFocusconceptContext()
}

type FocusconceptContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFocusconceptContext() *FocusconceptContext {
	var p = new(FocusconceptContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_focusconcept
	return p
}

func (*FocusconceptContext) IsFocusconceptContext() {}

func NewFocusconceptContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FocusconceptContext {
	var p = new(FocusconceptContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_focusconcept

	return p
}

func (s *FocusconceptContext) GetParser() antlr.Parser { return s.parser }

func (s *FocusconceptContext) AllConceptreference() []IConceptreferenceContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IConceptreferenceContext)(nil)).Elem())
	var tst = make([]IConceptreferenceContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IConceptreferenceContext)
		}
	}

	return tst
}

func (s *FocusconceptContext) Conceptreference(i int) IConceptreferenceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IConceptreferenceContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IConceptreferenceContext)
}

func (s *FocusconceptContext) AllWs() []IWsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWsContext)(nil)).Elem())
	var tst = make([]IWsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWsContext)
		}
	}

	return tst
}

func (s *FocusconceptContext) Ws(i int) IWsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWsContext)
}

func (s *FocusconceptContext) AllPLUS() []antlr.TerminalNode {
	return s.GetTokens(CGParserPLUS)
}

func (s *FocusconceptContext) PLUS(i int) antlr.TerminalNode {
	return s.GetToken(CGParserPLUS, i)
}

func (s *FocusconceptContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FocusconceptContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FocusconceptContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterFocusconcept(s)
	}
}

func (s *FocusconceptContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitFocusconcept(s)
	}
}

func (p *CGParser) Focusconcept() (localctx IFocusconceptContext) {
	localctx = NewFocusconceptContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, CGParserRULE_focusconcept)

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
		p.SetState(105)
		p.Conceptreference()
	}
	p.SetState(113)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(106)
				p.Ws()
			}
			{
				p.SetState(107)
				p.Match(CGParserPLUS)
			}
			{
				p.SetState(108)
				p.Ws()
			}
			{
				p.SetState(109)
				p.Conceptreference()
			}

		}
		p.SetState(115)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext())
	}

	return localctx
}

// IConceptreferenceContext is an interface to support dynamic dispatch.
type IConceptreferenceContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsConceptreferenceContext differentiates from other interfaces.
	IsConceptreferenceContext()
}

type ConceptreferenceContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConceptreferenceContext() *ConceptreferenceContext {
	var p = new(ConceptreferenceContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_conceptreference
	return p
}

func (*ConceptreferenceContext) IsConceptreferenceContext() {}

func NewConceptreferenceContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConceptreferenceContext {
	var p = new(ConceptreferenceContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_conceptreference

	return p
}

func (s *ConceptreferenceContext) GetParser() antlr.Parser { return s.parser }

func (s *ConceptreferenceContext) Conceptid() IConceptidContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IConceptidContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IConceptidContext)
}

func (s *ConceptreferenceContext) AllWs() []IWsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWsContext)(nil)).Elem())
	var tst = make([]IWsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWsContext)
		}
	}

	return tst
}

func (s *ConceptreferenceContext) Ws(i int) IWsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWsContext)
}

func (s *ConceptreferenceContext) AllPIPE() []antlr.TerminalNode {
	return s.GetTokens(CGParserPIPE)
}

func (s *ConceptreferenceContext) PIPE(i int) antlr.TerminalNode {
	return s.GetToken(CGParserPIPE, i)
}

func (s *ConceptreferenceContext) Term() ITermContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ITermContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ITermContext)
}

func (s *ConceptreferenceContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConceptreferenceContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConceptreferenceContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterConceptreference(s)
	}
}

func (s *ConceptreferenceContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitConceptreference(s)
	}
}

func (p *CGParser) Conceptreference() (localctx IConceptreferenceContext) {
	localctx = NewConceptreferenceContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, CGParserRULE_conceptreference)

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
		p.SetState(116)
		p.Conceptid()
	}
	p.SetState(124)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(117)
			p.Ws()
		}
		{
			p.SetState(118)
			p.Match(CGParserPIPE)
		}
		{
			p.SetState(119)
			p.Ws()
		}
		{
			p.SetState(120)
			p.Term()
		}
		{
			p.SetState(121)
			p.Ws()
		}
		{
			p.SetState(122)
			p.Match(CGParserPIPE)
		}

	}

	return localctx
}

// IConceptidContext is an interface to support dynamic dispatch.
type IConceptidContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsConceptidContext differentiates from other interfaces.
	IsConceptidContext()
}

type ConceptidContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConceptidContext() *ConceptidContext {
	var p = new(ConceptidContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_conceptid
	return p
}

func (*ConceptidContext) IsConceptidContext() {}

func NewConceptidContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConceptidContext {
	var p = new(ConceptidContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_conceptid

	return p
}

func (s *ConceptidContext) GetParser() antlr.Parser { return s.parser }

func (s *ConceptidContext) Sctid() ISctidContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISctidContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISctidContext)
}

func (s *ConceptidContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConceptidContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConceptidContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterConceptid(s)
	}
}

func (s *ConceptidContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitConceptid(s)
	}
}

func (p *CGParser) Conceptid() (localctx IConceptidContext) {
	localctx = NewConceptidContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, CGParserRULE_conceptid)

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
		p.SetState(126)
		p.Sctid()
	}

	return localctx
}

// ITermContext is an interface to support dynamic dispatch.
type ITermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTermContext differentiates from other interfaces.
	IsTermContext()
}

type TermContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTermContext() *TermContext {
	var p = new(TermContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_term
	return p
}

func (*TermContext) IsTermContext() {}

func NewTermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TermContext {
	var p = new(TermContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_term

	return p
}

func (s *TermContext) GetParser() antlr.Parser { return s.parser }

func (s *TermContext) AllNonwsnonpipe() []INonwsnonpipeContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*INonwsnonpipeContext)(nil)).Elem())
	var tst = make([]INonwsnonpipeContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(INonwsnonpipeContext)
		}
	}

	return tst
}

func (s *TermContext) Nonwsnonpipe(i int) INonwsnonpipeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INonwsnonpipeContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(INonwsnonpipeContext)
}

func (s *TermContext) AllSp() []ISpContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISpContext)(nil)).Elem())
	var tst = make([]ISpContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISpContext)
		}
	}

	return tst
}

func (s *TermContext) Sp(i int) ISpContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISpContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISpContext)
}

func (s *TermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TermContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterTerm(s)
	}
}

func (s *TermContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitTerm(s)
	}
}

func (p *CGParser) Term() (localctx ITermContext) {
	localctx = NewTermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, CGParserRULE_term)
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
		p.SetState(128)
		p.Nonwsnonpipe()
	}
	p.SetState(138)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			p.SetState(132)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			for _la == CGParserSPACE {
				{
					p.SetState(129)
					p.Sp()
				}

				p.SetState(134)
				p.GetErrorHandler().Sync(p)
				_la = p.GetTokenStream().LA(1)
			}
			{
				p.SetState(135)
				p.Nonwsnonpipe()
			}

		}
		p.SetState(140)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext())
	}

	return localctx
}

// IRefinementContext is an interface to support dynamic dispatch.
type IRefinementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRefinementContext differentiates from other interfaces.
	IsRefinementContext()
}

type RefinementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRefinementContext() *RefinementContext {
	var p = new(RefinementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_refinement
	return p
}

func (*RefinementContext) IsRefinementContext() {}

func NewRefinementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RefinementContext {
	var p = new(RefinementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_refinement

	return p
}

func (s *RefinementContext) GetParser() antlr.Parser { return s.parser }

func (s *RefinementContext) Attributeset() IAttributesetContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttributesetContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttributesetContext)
}

func (s *RefinementContext) AllAttributegroup() []IAttributegroupContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAttributegroupContext)(nil)).Elem())
	var tst = make([]IAttributegroupContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAttributegroupContext)
		}
	}

	return tst
}

func (s *RefinementContext) Attributegroup(i int) IAttributegroupContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttributegroupContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAttributegroupContext)
}

func (s *RefinementContext) AllWs() []IWsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWsContext)(nil)).Elem())
	var tst = make([]IWsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWsContext)
		}
	}

	return tst
}

func (s *RefinementContext) Ws(i int) IWsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWsContext)
}

func (s *RefinementContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(CGParserCOMMA)
}

func (s *RefinementContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CGParserCOMMA, i)
}

func (s *RefinementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RefinementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *RefinementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterRefinement(s)
	}
}

func (s *RefinementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitRefinement(s)
	}
}

func (p *CGParser) Refinement() (localctx IRefinementContext) {
	localctx = NewRefinementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, CGParserRULE_refinement)
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
	p.SetState(143)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CGParserONE, CGParserTWO, CGParserTHREE, CGParserFOUR, CGParserFIVE, CGParserSIX, CGParserSEVEN, CGParserEIGHT, CGParserNINE:
		{
			p.SetState(141)
			p.Attributeset()
		}

	case CGParserLEFT_CURLY_BRACE:
		{
			p.SetState(142)
			p.Attributegroup()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(154)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(145)
				p.Ws()
			}
			p.SetState(148)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)

			if _la == CGParserCOMMA {
				{
					p.SetState(146)
					p.Match(CGParserCOMMA)
				}
				{
					p.SetState(147)
					p.Ws()
				}

			}
			{
				p.SetState(150)
				p.Attributegroup()
			}

		}
		p.SetState(156)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
	}

	return localctx
}

// IAttributegroupContext is an interface to support dynamic dispatch.
type IAttributegroupContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAttributegroupContext differentiates from other interfaces.
	IsAttributegroupContext()
}

type AttributegroupContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAttributegroupContext() *AttributegroupContext {
	var p = new(AttributegroupContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_attributegroup
	return p
}

func (*AttributegroupContext) IsAttributegroupContext() {}

func NewAttributegroupContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributegroupContext {
	var p = new(AttributegroupContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_attributegroup

	return p
}

func (s *AttributegroupContext) GetParser() antlr.Parser { return s.parser }

func (s *AttributegroupContext) LEFT_CURLY_BRACE() antlr.TerminalNode {
	return s.GetToken(CGParserLEFT_CURLY_BRACE, 0)
}

func (s *AttributegroupContext) AllWs() []IWsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWsContext)(nil)).Elem())
	var tst = make([]IWsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWsContext)
		}
	}

	return tst
}

func (s *AttributegroupContext) Ws(i int) IWsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWsContext)
}

func (s *AttributegroupContext) Attributeset() IAttributesetContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttributesetContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttributesetContext)
}

func (s *AttributegroupContext) RIGHT_CURLY_BRACE() antlr.TerminalNode {
	return s.GetToken(CGParserRIGHT_CURLY_BRACE, 0)
}

func (s *AttributegroupContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AttributegroupContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AttributegroupContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterAttributegroup(s)
	}
}

func (s *AttributegroupContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitAttributegroup(s)
	}
}

func (p *CGParser) Attributegroup() (localctx IAttributegroupContext) {
	localctx = NewAttributegroupContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, CGParserRULE_attributegroup)

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
		p.Match(CGParserLEFT_CURLY_BRACE)
	}
	{
		p.SetState(158)
		p.Ws()
	}
	{
		p.SetState(159)
		p.Attributeset()
	}
	{
		p.SetState(160)
		p.Ws()
	}
	{
		p.SetState(161)
		p.Match(CGParserRIGHT_CURLY_BRACE)
	}

	return localctx
}

// IAttributesetContext is an interface to support dynamic dispatch.
type IAttributesetContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAttributesetContext differentiates from other interfaces.
	IsAttributesetContext()
}

type AttributesetContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAttributesetContext() *AttributesetContext {
	var p = new(AttributesetContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_attributeset
	return p
}

func (*AttributesetContext) IsAttributesetContext() {}

func NewAttributesetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributesetContext {
	var p = new(AttributesetContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_attributeset

	return p
}

func (s *AttributesetContext) GetParser() antlr.Parser { return s.parser }

func (s *AttributesetContext) AllAttribute() []IAttributeContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAttributeContext)(nil)).Elem())
	var tst = make([]IAttributeContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAttributeContext)
		}
	}

	return tst
}

func (s *AttributesetContext) Attribute(i int) IAttributeContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttributeContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAttributeContext)
}

func (s *AttributesetContext) AllWs() []IWsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWsContext)(nil)).Elem())
	var tst = make([]IWsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWsContext)
		}
	}

	return tst
}

func (s *AttributesetContext) Ws(i int) IWsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWsContext)
}

func (s *AttributesetContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(CGParserCOMMA)
}

func (s *AttributesetContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(CGParserCOMMA, i)
}

func (s *AttributesetContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AttributesetContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AttributesetContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterAttributeset(s)
	}
}

func (s *AttributesetContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitAttributeset(s)
	}
}

func (p *CGParser) Attributeset() (localctx IAttributesetContext) {
	localctx = NewAttributesetContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, CGParserRULE_attributeset)

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
		p.SetState(163)
		p.Attribute()
	}
	p.SetState(171)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(164)
				p.Ws()
			}
			{
				p.SetState(165)
				p.Match(CGParserCOMMA)
			}
			{
				p.SetState(166)
				p.Ws()
			}
			{
				p.SetState(167)
				p.Attribute()
			}

		}
		p.SetState(173)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext())
	}

	return localctx
}

// IAttributeContext is an interface to support dynamic dispatch.
type IAttributeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAttributeContext differentiates from other interfaces.
	IsAttributeContext()
}

type AttributeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAttributeContext() *AttributeContext {
	var p = new(AttributeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_attribute
	return p
}

func (*AttributeContext) IsAttributeContext() {}

func NewAttributeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributeContext {
	var p = new(AttributeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_attribute

	return p
}

func (s *AttributeContext) GetParser() antlr.Parser { return s.parser }

func (s *AttributeContext) Attributename() IAttributenameContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttributenameContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttributenameContext)
}

func (s *AttributeContext) AllWs() []IWsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWsContext)(nil)).Elem())
	var tst = make([]IWsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWsContext)
		}
	}

	return tst
}

func (s *AttributeContext) Ws(i int) IWsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWsContext)
}

func (s *AttributeContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(CGParserEQUALS, 0)
}

func (s *AttributeContext) Attributevalue() IAttributevalueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAttributevalueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IAttributevalueContext)
}

func (s *AttributeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AttributeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AttributeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterAttribute(s)
	}
}

func (s *AttributeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitAttribute(s)
	}
}

func (p *CGParser) Attribute() (localctx IAttributeContext) {
	localctx = NewAttributeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, CGParserRULE_attribute)

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
		p.SetState(174)
		p.Attributename()
	}
	{
		p.SetState(175)
		p.Ws()
	}
	{
		p.SetState(176)
		p.Match(CGParserEQUALS)
	}
	{
		p.SetState(177)
		p.Ws()
	}
	{
		p.SetState(178)
		p.Attributevalue()
	}

	return localctx
}

// IAttributenameContext is an interface to support dynamic dispatch.
type IAttributenameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAttributenameContext differentiates from other interfaces.
	IsAttributenameContext()
}

type AttributenameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAttributenameContext() *AttributenameContext {
	var p = new(AttributenameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_attributename
	return p
}

func (*AttributenameContext) IsAttributenameContext() {}

func NewAttributenameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributenameContext {
	var p = new(AttributenameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_attributename

	return p
}

func (s *AttributenameContext) GetParser() antlr.Parser { return s.parser }

func (s *AttributenameContext) Conceptreference() IConceptreferenceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IConceptreferenceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IConceptreferenceContext)
}

func (s *AttributenameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AttributenameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AttributenameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterAttributename(s)
	}
}

func (s *AttributenameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitAttributename(s)
	}
}

func (p *CGParser) Attributename() (localctx IAttributenameContext) {
	localctx = NewAttributenameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, CGParserRULE_attributename)

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
		p.SetState(180)
		p.Conceptreference()
	}

	return localctx
}

// IAttributevalueContext is an interface to support dynamic dispatch.
type IAttributevalueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAttributevalueContext differentiates from other interfaces.
	IsAttributevalueContext()
}

type AttributevalueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAttributevalueContext() *AttributevalueContext {
	var p = new(AttributevalueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_attributevalue
	return p
}

func (*AttributevalueContext) IsAttributevalueContext() {}

func NewAttributevalueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AttributevalueContext {
	var p = new(AttributevalueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_attributevalue

	return p
}

func (s *AttributevalueContext) GetParser() antlr.Parser { return s.parser }

func (s *AttributevalueContext) Expressionvalue() IExpressionvalueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IExpressionvalueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IExpressionvalueContext)
}

func (s *AttributevalueContext) AllQm() []IQmContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IQmContext)(nil)).Elem())
	var tst = make([]IQmContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IQmContext)
		}
	}

	return tst
}

func (s *AttributevalueContext) Qm(i int) IQmContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQmContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IQmContext)
}

func (s *AttributevalueContext) Stringvalue() IStringvalueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IStringvalueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IStringvalueContext)
}

func (s *AttributevalueContext) POUND() antlr.TerminalNode {
	return s.GetToken(CGParserPOUND, 0)
}

func (s *AttributevalueContext) Numericvalue() INumericvalueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*INumericvalueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(INumericvalueContext)
}

func (s *AttributevalueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AttributevalueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AttributevalueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterAttributevalue(s)
	}
}

func (s *AttributevalueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitAttributevalue(s)
	}
}

func (p *CGParser) Attributevalue() (localctx IAttributevalueContext) {
	localctx = NewAttributevalueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, CGParserRULE_attributevalue)

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

	p.SetState(189)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CGParserLEFT_PAREN, CGParserONE, CGParserTWO, CGParserTHREE, CGParserFOUR, CGParserFIVE, CGParserSIX, CGParserSEVEN, CGParserEIGHT, CGParserNINE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(182)
			p.Expressionvalue()
		}

	case CGParserQUOTE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(183)
			p.Qm()
		}
		{
			p.SetState(184)
			p.Stringvalue()
		}
		{
			p.SetState(185)
			p.Qm()
		}

	case CGParserPOUND:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(187)
			p.Match(CGParserPOUND)
		}
		{
			p.SetState(188)
			p.Numericvalue()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IExpressionvalueContext is an interface to support dynamic dispatch.
type IExpressionvalueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExpressionvalueContext differentiates from other interfaces.
	IsExpressionvalueContext()
}

type ExpressionvalueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionvalueContext() *ExpressionvalueContext {
	var p = new(ExpressionvalueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_expressionvalue
	return p
}

func (*ExpressionvalueContext) IsExpressionvalueContext() {}

func NewExpressionvalueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionvalueContext {
	var p = new(ExpressionvalueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_expressionvalue

	return p
}

func (s *ExpressionvalueContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionvalueContext) Conceptreference() IConceptreferenceContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IConceptreferenceContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IConceptreferenceContext)
}

func (s *ExpressionvalueContext) LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(CGParserLEFT_PAREN, 0)
}

func (s *ExpressionvalueContext) AllWs() []IWsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IWsContext)(nil)).Elem())
	var tst = make([]IWsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IWsContext)
		}
	}

	return tst
}

func (s *ExpressionvalueContext) Ws(i int) IWsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IWsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IWsContext)
}

func (s *ExpressionvalueContext) Subexpression() ISubexpressionContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISubexpressionContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ISubexpressionContext)
}

func (s *ExpressionvalueContext) RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(CGParserRIGHT_PAREN, 0)
}

func (s *ExpressionvalueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionvalueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionvalueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterExpressionvalue(s)
	}
}

func (s *ExpressionvalueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitExpressionvalue(s)
	}
}

func (p *CGParser) Expressionvalue() (localctx IExpressionvalueContext) {
	localctx = NewExpressionvalueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, CGParserRULE_expressionvalue)

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

	p.SetState(198)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CGParserONE, CGParserTWO, CGParserTHREE, CGParserFOUR, CGParserFIVE, CGParserSIX, CGParserSEVEN, CGParserEIGHT, CGParserNINE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(191)
			p.Conceptreference()
		}

	case CGParserLEFT_PAREN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(192)
			p.Match(CGParserLEFT_PAREN)
		}
		{
			p.SetState(193)
			p.Ws()
		}
		{
			p.SetState(194)
			p.Subexpression()
		}
		{
			p.SetState(195)
			p.Ws()
		}
		{
			p.SetState(196)
			p.Match(CGParserRIGHT_PAREN)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IStringvalueContext is an interface to support dynamic dispatch.
type IStringvalueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsStringvalueContext differentiates from other interfaces.
	IsStringvalueContext()
}

type StringvalueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStringvalueContext() *StringvalueContext {
	var p = new(StringvalueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_stringvalue
	return p
}

func (*StringvalueContext) IsStringvalueContext() {}

func NewStringvalueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StringvalueContext {
	var p = new(StringvalueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_stringvalue

	return p
}

func (s *StringvalueContext) GetParser() antlr.Parser { return s.parser }

func (s *StringvalueContext) AllAnynonescapedchar() []IAnynonescapedcharContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IAnynonescapedcharContext)(nil)).Elem())
	var tst = make([]IAnynonescapedcharContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IAnynonescapedcharContext)
		}
	}

	return tst
}

func (s *StringvalueContext) Anynonescapedchar(i int) IAnynonescapedcharContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IAnynonescapedcharContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IAnynonescapedcharContext)
}

func (s *StringvalueContext) AllEscapedchar() []IEscapedcharContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IEscapedcharContext)(nil)).Elem())
	var tst = make([]IEscapedcharContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IEscapedcharContext)
		}
	}

	return tst
}

func (s *StringvalueContext) Escapedchar(i int) IEscapedcharContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IEscapedcharContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IEscapedcharContext)
}

func (s *StringvalueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringvalueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StringvalueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterStringvalue(s)
	}
}

func (s *StringvalueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitStringvalue(s)
	}
}

func (p *CGParser) Stringvalue() (localctx IStringvalueContext) {
	localctx = NewStringvalueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, CGParserRULE_stringvalue)
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
	p.SetState(202)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CGParserTAB)|(1<<CGParserLF)|(1<<CGParserCR)|(1<<CGParserSPACE)|(1<<CGParserEXCLAMATION)|(1<<CGParserPOUND)|(1<<CGParserDOLLAR)|(1<<CGParserPERCENT)|(1<<CGParserAMPERSAND)|(1<<CGParserAPOSTROPHE)|(1<<CGParserLEFT_PAREN)|(1<<CGParserRIGHT_PAREN)|(1<<CGParserASTERISK)|(1<<CGParserPLUS)|(1<<CGParserCOMMA)|(1<<CGParserDASH)|(1<<CGParserPERIOD)|(1<<CGParserSLASH)|(1<<CGParserZERO)|(1<<CGParserONE)|(1<<CGParserTWO)|(1<<CGParserTHREE)|(1<<CGParserFOUR)|(1<<CGParserFIVE)|(1<<CGParserSIX)|(1<<CGParserSEVEN)|(1<<CGParserEIGHT)|(1<<CGParserNINE)|(1<<CGParserCOLON)|(1<<CGParserSEMICOLON))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(CGParserLESS_THAN-32))|(1<<(CGParserEQUALS-32))|(1<<(CGParserGREATER_THAN-32))|(1<<(CGParserQUESTION-32))|(1<<(CGParserAT-32))|(1<<(CGParserCAP_A-32))|(1<<(CGParserCAP_B-32))|(1<<(CGParserCAP_C-32))|(1<<(CGParserCAP_D-32))|(1<<(CGParserCAP_E-32))|(1<<(CGParserCAP_F-32))|(1<<(CGParserCAP_G-32))|(1<<(CGParserCAP_H-32))|(1<<(CGParserCAP_I-32))|(1<<(CGParserCAP_J-32))|(1<<(CGParserCAP_K-32))|(1<<(CGParserCAP_L-32))|(1<<(CGParserCAP_M-32))|(1<<(CGParserCAP_N-32))|(1<<(CGParserCAP_O-32))|(1<<(CGParserCAP_P-32))|(1<<(CGParserCAP_Q-32))|(1<<(CGParserCAP_R-32))|(1<<(CGParserCAP_S-32))|(1<<(CGParserCAP_T-32))|(1<<(CGParserCAP_U-32))|(1<<(CGParserCAP_V-32))|(1<<(CGParserCAP_W-32))|(1<<(CGParserCAP_X-32))|(1<<(CGParserCAP_Y-32))|(1<<(CGParserCAP_Z-32))|(1<<(CGParserLEFT_BRACE-32)))) != 0) || (((_la-64)&-(0x1f+1)) == 0 && ((1<<uint((_la-64)))&((1<<(CGParserBACKSLASH-64))|(1<<(CGParserRIGHT_BRACE-64))|(1<<(CGParserCARAT-64))|(1<<(CGParserUNDERSCORE-64))|(1<<(CGParserACCENT-64))|(1<<(CGParserA-64))|(1<<(CGParserB-64))|(1<<(CGParserC-64))|(1<<(CGParserD-64))|(1<<(CGParserE-64))|(1<<(CGParserF-64))|(1<<(CGParserG-64))|(1<<(CGParserH-64))|(1<<(CGParserI-64))|(1<<(CGParserJ-64))|(1<<(CGParserK-64))|(1<<(CGParserL-64))|(1<<(CGParserM-64))|(1<<(CGParserN-64))|(1<<(CGParserO-64))|(1<<(CGParserP-64))|(1<<(CGParserQ-64))|(1<<(CGParserR-64))|(1<<(CGParserS-64))|(1<<(CGParserT-64))|(1<<(CGParserU-64))|(1<<(CGParserV-64))|(1<<(CGParserW-64))|(1<<(CGParserX-64))|(1<<(CGParserY-64))|(1<<(CGParserZ-64))|(1<<(CGParserLEFT_CURLY_BRACE-64)))) != 0) || (((_la-96)&-(0x1f+1)) == 0 && ((1<<uint((_la-96)))&((1<<(CGParserPIPE-96))|(1<<(CGParserRIGHT_CURLY_BRACE-96))|(1<<(CGParserTILDE-96)))) != 0) || (((_la-163)&-(0x1f+1)) == 0 && ((1<<uint((_la-163)))&((1<<(CGParserU_00C2-163))|(1<<(CGParserU_00C3-163))|(1<<(CGParserU_00C4-163))|(1<<(CGParserU_00C5-163))|(1<<(CGParserU_00C6-163))|(1<<(CGParserU_00C7-163))|(1<<(CGParserU_00C8-163))|(1<<(CGParserU_00C9-163))|(1<<(CGParserU_00CA-163))|(1<<(CGParserU_00CB-163))|(1<<(CGParserU_00CC-163))|(1<<(CGParserU_00CD-163))|(1<<(CGParserU_00CE-163))|(1<<(CGParserU_00CF-163))|(1<<(CGParserU_00D0-163))|(1<<(CGParserU_00D1-163))|(1<<(CGParserU_00D2-163))|(1<<(CGParserU_00D3-163))|(1<<(CGParserU_00D4-163))|(1<<(CGParserU_00D5-163))|(1<<(CGParserU_00D6-163))|(1<<(CGParserU_00D7-163))|(1<<(CGParserU_00D8-163))|(1<<(CGParserU_00D9-163))|(1<<(CGParserU_00DA-163))|(1<<(CGParserU_00DB-163))|(1<<(CGParserU_00DC-163))|(1<<(CGParserU_00DD-163))|(1<<(CGParserU_00DE-163))|(1<<(CGParserU_00DF-163))|(1<<(CGParserU_00E0-163))|(1<<(CGParserU_00E1-163)))) != 0) || (((_la-195)&-(0x1f+1)) == 0 && ((1<<uint((_la-195)))&((1<<(CGParserU_00E2-195))|(1<<(CGParserU_00E3-195))|(1<<(CGParserU_00E4-195))|(1<<(CGParserU_00E5-195))|(1<<(CGParserU_00E6-195))|(1<<(CGParserU_00E7-195))|(1<<(CGParserU_00E8-195))|(1<<(CGParserU_00E9-195))|(1<<(CGParserU_00EA-195))|(1<<(CGParserU_00EB-195))|(1<<(CGParserU_00EC-195))|(1<<(CGParserU_00ED-195))|(1<<(CGParserU_00EE-195))|(1<<(CGParserU_00EF-195))|(1<<(CGParserU_00F0-195))|(1<<(CGParserU_00F1-195))|(1<<(CGParserU_00F2-195))|(1<<(CGParserU_00F3-195))|(1<<(CGParserU_00F4-195)))) != 0) {
		p.SetState(202)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case CGParserTAB, CGParserLF, CGParserCR, CGParserSPACE, CGParserEXCLAMATION, CGParserPOUND, CGParserDOLLAR, CGParserPERCENT, CGParserAMPERSAND, CGParserAPOSTROPHE, CGParserLEFT_PAREN, CGParserRIGHT_PAREN, CGParserASTERISK, CGParserPLUS, CGParserCOMMA, CGParserDASH, CGParserPERIOD, CGParserSLASH, CGParserZERO, CGParserONE, CGParserTWO, CGParserTHREE, CGParserFOUR, CGParserFIVE, CGParserSIX, CGParserSEVEN, CGParserEIGHT, CGParserNINE, CGParserCOLON, CGParserSEMICOLON, CGParserLESS_THAN, CGParserEQUALS, CGParserGREATER_THAN, CGParserQUESTION, CGParserAT, CGParserCAP_A, CGParserCAP_B, CGParserCAP_C, CGParserCAP_D, CGParserCAP_E, CGParserCAP_F, CGParserCAP_G, CGParserCAP_H, CGParserCAP_I, CGParserCAP_J, CGParserCAP_K, CGParserCAP_L, CGParserCAP_M, CGParserCAP_N, CGParserCAP_O, CGParserCAP_P, CGParserCAP_Q, CGParserCAP_R, CGParserCAP_S, CGParserCAP_T, CGParserCAP_U, CGParserCAP_V, CGParserCAP_W, CGParserCAP_X, CGParserCAP_Y, CGParserCAP_Z, CGParserLEFT_BRACE, CGParserRIGHT_BRACE, CGParserCARAT, CGParserUNDERSCORE, CGParserACCENT, CGParserA, CGParserB, CGParserC, CGParserD, CGParserE, CGParserF, CGParserG, CGParserH, CGParserI, CGParserJ, CGParserK, CGParserL, CGParserM, CGParserN, CGParserO, CGParserP, CGParserQ, CGParserR, CGParserS, CGParserT, CGParserU, CGParserV, CGParserW, CGParserX, CGParserY, CGParserZ, CGParserLEFT_CURLY_BRACE, CGParserPIPE, CGParserRIGHT_CURLY_BRACE, CGParserTILDE, CGParserU_00C2, CGParserU_00C3, CGParserU_00C4, CGParserU_00C5, CGParserU_00C6, CGParserU_00C7, CGParserU_00C8, CGParserU_00C9, CGParserU_00CA, CGParserU_00CB, CGParserU_00CC, CGParserU_00CD, CGParserU_00CE, CGParserU_00CF, CGParserU_00D0, CGParserU_00D1, CGParserU_00D2, CGParserU_00D3, CGParserU_00D4, CGParserU_00D5, CGParserU_00D6, CGParserU_00D7, CGParserU_00D8, CGParserU_00D9, CGParserU_00DA, CGParserU_00DB, CGParserU_00DC, CGParserU_00DD, CGParserU_00DE, CGParserU_00DF, CGParserU_00E0, CGParserU_00E1, CGParserU_00E2, CGParserU_00E3, CGParserU_00E4, CGParserU_00E5, CGParserU_00E6, CGParserU_00E7, CGParserU_00E8, CGParserU_00E9, CGParserU_00EA, CGParserU_00EB, CGParserU_00EC, CGParserU_00ED, CGParserU_00EE, CGParserU_00EF, CGParserU_00F0, CGParserU_00F1, CGParserU_00F2, CGParserU_00F3, CGParserU_00F4:
			{
				p.SetState(200)
				p.Anynonescapedchar()
			}

		case CGParserBACKSLASH:
			{
				p.SetState(201)
				p.Escapedchar()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(204)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// INumericvalueContext is an interface to support dynamic dispatch.
type INumericvalueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNumericvalueContext differentiates from other interfaces.
	IsNumericvalueContext()
}

type NumericvalueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNumericvalueContext() *NumericvalueContext {
	var p = new(NumericvalueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_numericvalue
	return p
}

func (*NumericvalueContext) IsNumericvalueContext() {}

func NewNumericvalueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NumericvalueContext {
	var p = new(NumericvalueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_numericvalue

	return p
}

func (s *NumericvalueContext) GetParser() antlr.Parser { return s.parser }

func (s *NumericvalueContext) Decimalvalue() IDecimalvalueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDecimalvalueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDecimalvalueContext)
}

func (s *NumericvalueContext) Integervalue() IIntegervalueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegervalueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntegervalueContext)
}

func (s *NumericvalueContext) DASH() antlr.TerminalNode {
	return s.GetToken(CGParserDASH, 0)
}

func (s *NumericvalueContext) PLUS() antlr.TerminalNode {
	return s.GetToken(CGParserPLUS, 0)
}

func (s *NumericvalueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumericvalueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NumericvalueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterNumericvalue(s)
	}
}

func (s *NumericvalueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitNumericvalue(s)
	}
}

func (p *CGParser) Numericvalue() (localctx INumericvalueContext) {
	localctx = NewNumericvalueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, CGParserRULE_numericvalue)
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
	p.SetState(207)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == CGParserPLUS || _la == CGParserDASH {
		{
			p.SetState(206)
			_la = p.GetTokenStream().LA(1)

			if !(_la == CGParserPLUS || _la == CGParserDASH) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	p.SetState(211)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 16, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(209)
			p.Decimalvalue()
		}

	case 2:
		{
			p.SetState(210)
			p.Integervalue()
		}

	}

	return localctx
}

// IIntegervalueContext is an interface to support dynamic dispatch.
type IIntegervalueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIntegervalueContext differentiates from other interfaces.
	IsIntegervalueContext()
}

type IntegervalueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntegervalueContext() *IntegervalueContext {
	var p = new(IntegervalueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_integervalue
	return p
}

func (*IntegervalueContext) IsIntegervalueContext() {}

func NewIntegervalueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntegervalueContext {
	var p = new(IntegervalueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_integervalue

	return p
}

func (s *IntegervalueContext) GetParser() antlr.Parser { return s.parser }

func (s *IntegervalueContext) Digitnonzero() IDigitnonzeroContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDigitnonzeroContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDigitnonzeroContext)
}

func (s *IntegervalueContext) AllDigit() []IDigitContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IDigitContext)(nil)).Elem())
	var tst = make([]IDigitContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IDigitContext)
		}
	}

	return tst
}

func (s *IntegervalueContext) Digit(i int) IDigitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDigitContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IDigitContext)
}

func (s *IntegervalueContext) Zero() IZeroContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IZeroContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IZeroContext)
}

func (s *IntegervalueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntegervalueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntegervalueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterIntegervalue(s)
	}
}

func (s *IntegervalueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitIntegervalue(s)
	}
}

func (p *CGParser) Integervalue() (localctx IIntegervalueContext) {
	localctx = NewIntegervalueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, CGParserRULE_integervalue)
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

	p.SetState(221)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CGParserONE, CGParserTWO, CGParserTHREE, CGParserFOUR, CGParserFIVE, CGParserSIX, CGParserSEVEN, CGParserEIGHT, CGParserNINE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(213)
			p.Digitnonzero()
		}
		p.SetState(217)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CGParserZERO)|(1<<CGParserONE)|(1<<CGParserTWO)|(1<<CGParserTHREE)|(1<<CGParserFOUR)|(1<<CGParserFIVE)|(1<<CGParserSIX)|(1<<CGParserSEVEN)|(1<<CGParserEIGHT)|(1<<CGParserNINE))) != 0 {
			{
				p.SetState(214)
				p.Digit()
			}

			p.SetState(219)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case CGParserZERO:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(220)
			p.Zero()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IDecimalvalueContext is an interface to support dynamic dispatch.
type IDecimalvalueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDecimalvalueContext differentiates from other interfaces.
	IsDecimalvalueContext()
}

type DecimalvalueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDecimalvalueContext() *DecimalvalueContext {
	var p = new(DecimalvalueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_decimalvalue
	return p
}

func (*DecimalvalueContext) IsDecimalvalueContext() {}

func NewDecimalvalueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DecimalvalueContext {
	var p = new(DecimalvalueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_decimalvalue

	return p
}

func (s *DecimalvalueContext) GetParser() antlr.Parser { return s.parser }

func (s *DecimalvalueContext) Integervalue() IIntegervalueContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IIntegervalueContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IIntegervalueContext)
}

func (s *DecimalvalueContext) PERIOD() antlr.TerminalNode {
	return s.GetToken(CGParserPERIOD, 0)
}

func (s *DecimalvalueContext) AllDigit() []IDigitContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IDigitContext)(nil)).Elem())
	var tst = make([]IDigitContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IDigitContext)
		}
	}

	return tst
}

func (s *DecimalvalueContext) Digit(i int) IDigitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDigitContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IDigitContext)
}

func (s *DecimalvalueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DecimalvalueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DecimalvalueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterDecimalvalue(s)
	}
}

func (s *DecimalvalueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitDecimalvalue(s)
	}
}

func (p *CGParser) Decimalvalue() (localctx IDecimalvalueContext) {
	localctx = NewDecimalvalueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, CGParserRULE_decimalvalue)
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
		p.SetState(223)
		p.Integervalue()
	}
	{
		p.SetState(224)
		p.Match(CGParserPERIOD)
	}
	p.SetState(226)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = (((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CGParserZERO)|(1<<CGParserONE)|(1<<CGParserTWO)|(1<<CGParserTHREE)|(1<<CGParserFOUR)|(1<<CGParserFIVE)|(1<<CGParserSIX)|(1<<CGParserSEVEN)|(1<<CGParserEIGHT)|(1<<CGParserNINE))) != 0) {
		{
			p.SetState(225)
			p.Digit()
		}

		p.SetState(228)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ISctidContext is an interface to support dynamic dispatch.
type ISctidContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSctidContext differentiates from other interfaces.
	IsSctidContext()
}

type SctidContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySctidContext() *SctidContext {
	var p = new(SctidContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_sctid
	return p
}

func (*SctidContext) IsSctidContext() {}

func NewSctidContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SctidContext {
	var p = new(SctidContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_sctid

	return p
}

func (s *SctidContext) GetParser() antlr.Parser { return s.parser }

func (s *SctidContext) Digitnonzero() IDigitnonzeroContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDigitnonzeroContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IDigitnonzeroContext)
}

func (s *SctidContext) AllDigit() []IDigitContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IDigitContext)(nil)).Elem())
	var tst = make([]IDigitContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IDigitContext)
		}
	}

	return tst
}

func (s *SctidContext) Digit(i int) IDigitContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IDigitContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IDigitContext)
}

func (s *SctidContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SctidContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SctidContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterSctid(s)
	}
}

func (s *SctidContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitSctid(s)
	}
}

func (p *CGParser) Sctid() (localctx ISctidContext) {
	localctx = NewSctidContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, CGParserRULE_sctid)
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
		p.SetState(230)
		p.Digitnonzero()
	}

	{
		p.SetState(231)
		p.Digit()
	}

	{
		p.SetState(232)
		p.Digit()
	}

	{
		p.SetState(233)
		p.Digit()
	}

	{
		p.SetState(234)
		p.Digit()
	}

	{
		p.SetState(235)
		p.Digit()
	}

	p.SetState(327)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		p.SetState(237)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CGParserZERO)|(1<<CGParserONE)|(1<<CGParserTWO)|(1<<CGParserTHREE)|(1<<CGParserFOUR)|(1<<CGParserFIVE)|(1<<CGParserSIX)|(1<<CGParserSEVEN)|(1<<CGParserEIGHT)|(1<<CGParserNINE))) != 0 {
			{
				p.SetState(236)
				p.Digit()
			}

		}

	case 2:
		{
			p.SetState(239)
			p.Digit()
		}

		{
			p.SetState(240)
			p.Digit()
		}

	case 3:
		{
			p.SetState(242)
			p.Digit()
		}

		{
			p.SetState(243)
			p.Digit()
		}

		{
			p.SetState(244)
			p.Digit()
		}

	case 4:
		{
			p.SetState(246)
			p.Digit()
		}

		{
			p.SetState(247)
			p.Digit()
		}

		{
			p.SetState(248)
			p.Digit()
		}

		{
			p.SetState(249)
			p.Digit()
		}

	case 5:
		{
			p.SetState(251)
			p.Digit()
		}

		{
			p.SetState(252)
			p.Digit()
		}

		{
			p.SetState(253)
			p.Digit()
		}

		{
			p.SetState(254)
			p.Digit()
		}

		{
			p.SetState(255)
			p.Digit()
		}

	case 6:
		{
			p.SetState(257)
			p.Digit()
		}

		{
			p.SetState(258)
			p.Digit()
		}

		{
			p.SetState(259)
			p.Digit()
		}

		{
			p.SetState(260)
			p.Digit()
		}

		{
			p.SetState(261)
			p.Digit()
		}

		{
			p.SetState(262)
			p.Digit()
		}

	case 7:
		{
			p.SetState(264)
			p.Digit()
		}

		{
			p.SetState(265)
			p.Digit()
		}

		{
			p.SetState(266)
			p.Digit()
		}

		{
			p.SetState(267)
			p.Digit()
		}

		{
			p.SetState(268)
			p.Digit()
		}

		{
			p.SetState(269)
			p.Digit()
		}

		{
			p.SetState(270)
			p.Digit()
		}

	case 8:
		{
			p.SetState(272)
			p.Digit()
		}

		{
			p.SetState(273)
			p.Digit()
		}

		{
			p.SetState(274)
			p.Digit()
		}

		{
			p.SetState(275)
			p.Digit()
		}

		{
			p.SetState(276)
			p.Digit()
		}

		{
			p.SetState(277)
			p.Digit()
		}

		{
			p.SetState(278)
			p.Digit()
		}

		{
			p.SetState(279)
			p.Digit()
		}

	case 9:
		{
			p.SetState(281)
			p.Digit()
		}

		{
			p.SetState(282)
			p.Digit()
		}

		{
			p.SetState(283)
			p.Digit()
		}

		{
			p.SetState(284)
			p.Digit()
		}

		{
			p.SetState(285)
			p.Digit()
		}

		{
			p.SetState(286)
			p.Digit()
		}

		{
			p.SetState(287)
			p.Digit()
		}

		{
			p.SetState(288)
			p.Digit()
		}

		{
			p.SetState(289)
			p.Digit()
		}

	case 10:
		{
			p.SetState(291)
			p.Digit()
		}

		{
			p.SetState(292)
			p.Digit()
		}

		{
			p.SetState(293)
			p.Digit()
		}

		{
			p.SetState(294)
			p.Digit()
		}

		{
			p.SetState(295)
			p.Digit()
		}

		{
			p.SetState(296)
			p.Digit()
		}

		{
			p.SetState(297)
			p.Digit()
		}

		{
			p.SetState(298)
			p.Digit()
		}

		{
			p.SetState(299)
			p.Digit()
		}

		{
			p.SetState(300)
			p.Digit()
		}

	case 11:
		{
			p.SetState(302)
			p.Digit()
		}

		{
			p.SetState(303)
			p.Digit()
		}

		{
			p.SetState(304)
			p.Digit()
		}

		{
			p.SetState(305)
			p.Digit()
		}

		{
			p.SetState(306)
			p.Digit()
		}

		{
			p.SetState(307)
			p.Digit()
		}

		{
			p.SetState(308)
			p.Digit()
		}

		{
			p.SetState(309)
			p.Digit()
		}

		{
			p.SetState(310)
			p.Digit()
		}

		{
			p.SetState(311)
			p.Digit()
		}

		{
			p.SetState(312)
			p.Digit()
		}

	case 12:
		{
			p.SetState(314)
			p.Digit()
		}

		{
			p.SetState(315)
			p.Digit()
		}

		{
			p.SetState(316)
			p.Digit()
		}

		{
			p.SetState(317)
			p.Digit()
		}

		{
			p.SetState(318)
			p.Digit()
		}

		{
			p.SetState(319)
			p.Digit()
		}

		{
			p.SetState(320)
			p.Digit()
		}

		{
			p.SetState(321)
			p.Digit()
		}

		{
			p.SetState(322)
			p.Digit()
		}

		{
			p.SetState(323)
			p.Digit()
		}

		{
			p.SetState(324)
			p.Digit()
		}

		{
			p.SetState(325)
			p.Digit()
		}

	}

	return localctx
}

// IWsContext is an interface to support dynamic dispatch.
type IWsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWsContext differentiates from other interfaces.
	IsWsContext()
}

type WsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWsContext() *WsContext {
	var p = new(WsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_ws
	return p
}

func (*WsContext) IsWsContext() {}

func NewWsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WsContext {
	var p = new(WsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_ws

	return p
}

func (s *WsContext) GetParser() antlr.Parser { return s.parser }

func (s *WsContext) AllSp() []ISpContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ISpContext)(nil)).Elem())
	var tst = make([]ISpContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ISpContext)
		}
	}

	return tst
}

func (s *WsContext) Sp(i int) ISpContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ISpContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ISpContext)
}

func (s *WsContext) AllHtab() []IHtabContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IHtabContext)(nil)).Elem())
	var tst = make([]IHtabContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IHtabContext)
		}
	}

	return tst
}

func (s *WsContext) Htab(i int) IHtabContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHtabContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IHtabContext)
}

func (s *WsContext) AllCr() []ICrContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ICrContext)(nil)).Elem())
	var tst = make([]ICrContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ICrContext)
		}
	}

	return tst
}

func (s *WsContext) Cr(i int) ICrContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICrContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ICrContext)
}

func (s *WsContext) AllLf() []ILfContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*ILfContext)(nil)).Elem())
	var tst = make([]ILfContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(ILfContext)
		}
	}

	return tst
}

func (s *WsContext) Lf(i int) ILfContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILfContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(ILfContext)
}

func (s *WsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterWs(s)
	}
}

func (s *WsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitWs(s)
	}
}

func (p *CGParser) Ws() (localctx IWsContext) {
	localctx = NewWsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, CGParserRULE_ws)
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
	p.SetState(335)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CGParserTAB)|(1<<CGParserLF)|(1<<CGParserCR)|(1<<CGParserSPACE))) != 0 {
		p.SetState(333)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case CGParserSPACE:
			{
				p.SetState(329)
				p.Sp()
			}

		case CGParserTAB:
			{
				p.SetState(330)
				p.Htab()
			}

		case CGParserCR:
			{
				p.SetState(331)
				p.Cr()
			}

		case CGParserLF:
			{
				p.SetState(332)
				p.Lf()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(337)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// ISpContext is an interface to support dynamic dispatch.
type ISpContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSpContext differentiates from other interfaces.
	IsSpContext()
}

type SpContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySpContext() *SpContext {
	var p = new(SpContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_sp
	return p
}

func (*SpContext) IsSpContext() {}

func NewSpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SpContext {
	var p = new(SpContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_sp

	return p
}

func (s *SpContext) GetParser() antlr.Parser { return s.parser }

func (s *SpContext) SPACE() antlr.TerminalNode {
	return s.GetToken(CGParserSPACE, 0)
}

func (s *SpContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SpContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SpContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterSp(s)
	}
}

func (s *SpContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitSp(s)
	}
}

func (p *CGParser) Sp() (localctx ISpContext) {
	localctx = NewSpContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, CGParserRULE_sp)

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
		p.SetState(338)
		p.Match(CGParserSPACE)
	}

	return localctx
}

// IHtabContext is an interface to support dynamic dispatch.
type IHtabContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsHtabContext differentiates from other interfaces.
	IsHtabContext()
}

type HtabContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyHtabContext() *HtabContext {
	var p = new(HtabContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_htab
	return p
}

func (*HtabContext) IsHtabContext() {}

func NewHtabContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HtabContext {
	var p = new(HtabContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_htab

	return p
}

func (s *HtabContext) GetParser() antlr.Parser { return s.parser }

func (s *HtabContext) TAB() antlr.TerminalNode {
	return s.GetToken(CGParserTAB, 0)
}

func (s *HtabContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *HtabContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *HtabContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterHtab(s)
	}
}

func (s *HtabContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitHtab(s)
	}
}

func (p *CGParser) Htab() (localctx IHtabContext) {
	localctx = NewHtabContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, CGParserRULE_htab)

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
		p.SetState(340)
		p.Match(CGParserTAB)
	}

	return localctx
}

// ICrContext is an interface to support dynamic dispatch.
type ICrContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCrContext differentiates from other interfaces.
	IsCrContext()
}

type CrContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCrContext() *CrContext {
	var p = new(CrContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_cr
	return p
}

func (*CrContext) IsCrContext() {}

func NewCrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CrContext {
	var p = new(CrContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_cr

	return p
}

func (s *CrContext) GetParser() antlr.Parser { return s.parser }

func (s *CrContext) CR() antlr.TerminalNode {
	return s.GetToken(CGParserCR, 0)
}

func (s *CrContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CrContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CrContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterCr(s)
	}
}

func (s *CrContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitCr(s)
	}
}

func (p *CGParser) Cr() (localctx ICrContext) {
	localctx = NewCrContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, CGParserRULE_cr)

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
		p.SetState(342)
		p.Match(CGParserCR)
	}

	return localctx
}

// ILfContext is an interface to support dynamic dispatch.
type ILfContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLfContext differentiates from other interfaces.
	IsLfContext()
}

type LfContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLfContext() *LfContext {
	var p = new(LfContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_lf
	return p
}

func (*LfContext) IsLfContext() {}

func NewLfContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LfContext {
	var p = new(LfContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_lf

	return p
}

func (s *LfContext) GetParser() antlr.Parser { return s.parser }

func (s *LfContext) LF() antlr.TerminalNode {
	return s.GetToken(CGParserLF, 0)
}

func (s *LfContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LfContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LfContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterLf(s)
	}
}

func (s *LfContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitLf(s)
	}
}

func (p *CGParser) Lf() (localctx ILfContext) {
	localctx = NewLfContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, CGParserRULE_lf)

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
		p.SetState(344)
		p.Match(CGParserLF)
	}

	return localctx
}

// IQmContext is an interface to support dynamic dispatch.
type IQmContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQmContext differentiates from other interfaces.
	IsQmContext()
}

type QmContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQmContext() *QmContext {
	var p = new(QmContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_qm
	return p
}

func (*QmContext) IsQmContext() {}

func NewQmContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QmContext {
	var p = new(QmContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_qm

	return p
}

func (s *QmContext) GetParser() antlr.Parser { return s.parser }

func (s *QmContext) QUOTE() antlr.TerminalNode {
	return s.GetToken(CGParserQUOTE, 0)
}

func (s *QmContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QmContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QmContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterQm(s)
	}
}

func (s *QmContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitQm(s)
	}
}

func (p *CGParser) Qm() (localctx IQmContext) {
	localctx = NewQmContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, CGParserRULE_qm)

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
		p.SetState(346)
		p.Match(CGParserQUOTE)
	}

	return localctx
}

// IBsContext is an interface to support dynamic dispatch.
type IBsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBsContext differentiates from other interfaces.
	IsBsContext()
}

type BsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBsContext() *BsContext {
	var p = new(BsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_bs
	return p
}

func (*BsContext) IsBsContext() {}

func NewBsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BsContext {
	var p = new(BsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_bs

	return p
}

func (s *BsContext) GetParser() antlr.Parser { return s.parser }

func (s *BsContext) BACKSLASH() antlr.TerminalNode {
	return s.GetToken(CGParserBACKSLASH, 0)
}

func (s *BsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterBs(s)
	}
}

func (s *BsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitBs(s)
	}
}

func (p *CGParser) Bs() (localctx IBsContext) {
	localctx = NewBsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, CGParserRULE_bs)

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
		p.SetState(348)
		p.Match(CGParserBACKSLASH)
	}

	return localctx
}

// IDigitContext is an interface to support dynamic dispatch.
type IDigitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDigitContext differentiates from other interfaces.
	IsDigitContext()
}

type DigitContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDigitContext() *DigitContext {
	var p = new(DigitContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_digit
	return p
}

func (*DigitContext) IsDigitContext() {}

func NewDigitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DigitContext {
	var p = new(DigitContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_digit

	return p
}

func (s *DigitContext) GetParser() antlr.Parser { return s.parser }

func (s *DigitContext) ZERO() antlr.TerminalNode {
	return s.GetToken(CGParserZERO, 0)
}

func (s *DigitContext) ONE() antlr.TerminalNode {
	return s.GetToken(CGParserONE, 0)
}

func (s *DigitContext) TWO() antlr.TerminalNode {
	return s.GetToken(CGParserTWO, 0)
}

func (s *DigitContext) THREE() antlr.TerminalNode {
	return s.GetToken(CGParserTHREE, 0)
}

func (s *DigitContext) FOUR() antlr.TerminalNode {
	return s.GetToken(CGParserFOUR, 0)
}

func (s *DigitContext) FIVE() antlr.TerminalNode {
	return s.GetToken(CGParserFIVE, 0)
}

func (s *DigitContext) SIX() antlr.TerminalNode {
	return s.GetToken(CGParserSIX, 0)
}

func (s *DigitContext) SEVEN() antlr.TerminalNode {
	return s.GetToken(CGParserSEVEN, 0)
}

func (s *DigitContext) EIGHT() antlr.TerminalNode {
	return s.GetToken(CGParserEIGHT, 0)
}

func (s *DigitContext) NINE() antlr.TerminalNode {
	return s.GetToken(CGParserNINE, 0)
}

func (s *DigitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DigitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DigitContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterDigit(s)
	}
}

func (s *DigitContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitDigit(s)
	}
}

func (p *CGParser) Digit() (localctx IDigitContext) {
	localctx = NewDigitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, CGParserRULE_digit)
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
		p.SetState(350)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CGParserZERO)|(1<<CGParserONE)|(1<<CGParserTWO)|(1<<CGParserTHREE)|(1<<CGParserFOUR)|(1<<CGParserFIVE)|(1<<CGParserSIX)|(1<<CGParserSEVEN)|(1<<CGParserEIGHT)|(1<<CGParserNINE))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IZeroContext is an interface to support dynamic dispatch.
type IZeroContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsZeroContext differentiates from other interfaces.
	IsZeroContext()
}

type ZeroContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyZeroContext() *ZeroContext {
	var p = new(ZeroContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_zero
	return p
}

func (*ZeroContext) IsZeroContext() {}

func NewZeroContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ZeroContext {
	var p = new(ZeroContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_zero

	return p
}

func (s *ZeroContext) GetParser() antlr.Parser { return s.parser }

func (s *ZeroContext) ZERO() antlr.TerminalNode {
	return s.GetToken(CGParserZERO, 0)
}

func (s *ZeroContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ZeroContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ZeroContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterZero(s)
	}
}

func (s *ZeroContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitZero(s)
	}
}

func (p *CGParser) Zero() (localctx IZeroContext) {
	localctx = NewZeroContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, CGParserRULE_zero)

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
		p.SetState(352)
		p.Match(CGParserZERO)
	}

	return localctx
}

// IDigitnonzeroContext is an interface to support dynamic dispatch.
type IDigitnonzeroContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDigitnonzeroContext differentiates from other interfaces.
	IsDigitnonzeroContext()
}

type DigitnonzeroContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDigitnonzeroContext() *DigitnonzeroContext {
	var p = new(DigitnonzeroContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_digitnonzero
	return p
}

func (*DigitnonzeroContext) IsDigitnonzeroContext() {}

func NewDigitnonzeroContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DigitnonzeroContext {
	var p = new(DigitnonzeroContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_digitnonzero

	return p
}

func (s *DigitnonzeroContext) GetParser() antlr.Parser { return s.parser }

func (s *DigitnonzeroContext) ONE() antlr.TerminalNode {
	return s.GetToken(CGParserONE, 0)
}

func (s *DigitnonzeroContext) TWO() antlr.TerminalNode {
	return s.GetToken(CGParserTWO, 0)
}

func (s *DigitnonzeroContext) THREE() antlr.TerminalNode {
	return s.GetToken(CGParserTHREE, 0)
}

func (s *DigitnonzeroContext) FOUR() antlr.TerminalNode {
	return s.GetToken(CGParserFOUR, 0)
}

func (s *DigitnonzeroContext) FIVE() antlr.TerminalNode {
	return s.GetToken(CGParserFIVE, 0)
}

func (s *DigitnonzeroContext) SIX() antlr.TerminalNode {
	return s.GetToken(CGParserSIX, 0)
}

func (s *DigitnonzeroContext) SEVEN() antlr.TerminalNode {
	return s.GetToken(CGParserSEVEN, 0)
}

func (s *DigitnonzeroContext) EIGHT() antlr.TerminalNode {
	return s.GetToken(CGParserEIGHT, 0)
}

func (s *DigitnonzeroContext) NINE() antlr.TerminalNode {
	return s.GetToken(CGParserNINE, 0)
}

func (s *DigitnonzeroContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DigitnonzeroContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DigitnonzeroContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterDigitnonzero(s)
	}
}

func (s *DigitnonzeroContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitDigitnonzero(s)
	}
}

func (p *CGParser) Digitnonzero() (localctx IDigitnonzeroContext) {
	localctx = NewDigitnonzeroContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, CGParserRULE_digitnonzero)
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
		p.SetState(354)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CGParserONE)|(1<<CGParserTWO)|(1<<CGParserTHREE)|(1<<CGParserFOUR)|(1<<CGParserFIVE)|(1<<CGParserSIX)|(1<<CGParserSEVEN)|(1<<CGParserEIGHT)|(1<<CGParserNINE))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// INonwsnonpipeContext is an interface to support dynamic dispatch.
type INonwsnonpipeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsNonwsnonpipeContext differentiates from other interfaces.
	IsNonwsnonpipeContext()
}

type NonwsnonpipeContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNonwsnonpipeContext() *NonwsnonpipeContext {
	var p = new(NonwsnonpipeContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_nonwsnonpipe
	return p
}

func (*NonwsnonpipeContext) IsNonwsnonpipeContext() {}

func NewNonwsnonpipeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NonwsnonpipeContext {
	var p = new(NonwsnonpipeContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_nonwsnonpipe

	return p
}

func (s *NonwsnonpipeContext) GetParser() antlr.Parser { return s.parser }

func (s *NonwsnonpipeContext) EXCLAMATION() antlr.TerminalNode {
	return s.GetToken(CGParserEXCLAMATION, 0)
}

func (s *NonwsnonpipeContext) QUOTE() antlr.TerminalNode {
	return s.GetToken(CGParserQUOTE, 0)
}

func (s *NonwsnonpipeContext) POUND() antlr.TerminalNode {
	return s.GetToken(CGParserPOUND, 0)
}

func (s *NonwsnonpipeContext) DOLLAR() antlr.TerminalNode {
	return s.GetToken(CGParserDOLLAR, 0)
}

func (s *NonwsnonpipeContext) PERCENT() antlr.TerminalNode {
	return s.GetToken(CGParserPERCENT, 0)
}

func (s *NonwsnonpipeContext) AMPERSAND() antlr.TerminalNode {
	return s.GetToken(CGParserAMPERSAND, 0)
}

func (s *NonwsnonpipeContext) APOSTROPHE() antlr.TerminalNode {
	return s.GetToken(CGParserAPOSTROPHE, 0)
}

func (s *NonwsnonpipeContext) LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(CGParserLEFT_PAREN, 0)
}

func (s *NonwsnonpipeContext) RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(CGParserRIGHT_PAREN, 0)
}

func (s *NonwsnonpipeContext) ASTERISK() antlr.TerminalNode {
	return s.GetToken(CGParserASTERISK, 0)
}

func (s *NonwsnonpipeContext) PLUS() antlr.TerminalNode {
	return s.GetToken(CGParserPLUS, 0)
}

func (s *NonwsnonpipeContext) COMMA() antlr.TerminalNode {
	return s.GetToken(CGParserCOMMA, 0)
}

func (s *NonwsnonpipeContext) DASH() antlr.TerminalNode {
	return s.GetToken(CGParserDASH, 0)
}

func (s *NonwsnonpipeContext) PERIOD() antlr.TerminalNode {
	return s.GetToken(CGParserPERIOD, 0)
}

func (s *NonwsnonpipeContext) SLASH() antlr.TerminalNode {
	return s.GetToken(CGParserSLASH, 0)
}

func (s *NonwsnonpipeContext) ZERO() antlr.TerminalNode {
	return s.GetToken(CGParserZERO, 0)
}

func (s *NonwsnonpipeContext) ONE() antlr.TerminalNode {
	return s.GetToken(CGParserONE, 0)
}

func (s *NonwsnonpipeContext) TWO() antlr.TerminalNode {
	return s.GetToken(CGParserTWO, 0)
}

func (s *NonwsnonpipeContext) THREE() antlr.TerminalNode {
	return s.GetToken(CGParserTHREE, 0)
}

func (s *NonwsnonpipeContext) FOUR() antlr.TerminalNode {
	return s.GetToken(CGParserFOUR, 0)
}

func (s *NonwsnonpipeContext) FIVE() antlr.TerminalNode {
	return s.GetToken(CGParserFIVE, 0)
}

func (s *NonwsnonpipeContext) SIX() antlr.TerminalNode {
	return s.GetToken(CGParserSIX, 0)
}

func (s *NonwsnonpipeContext) SEVEN() antlr.TerminalNode {
	return s.GetToken(CGParserSEVEN, 0)
}

func (s *NonwsnonpipeContext) EIGHT() antlr.TerminalNode {
	return s.GetToken(CGParserEIGHT, 0)
}

func (s *NonwsnonpipeContext) NINE() antlr.TerminalNode {
	return s.GetToken(CGParserNINE, 0)
}

func (s *NonwsnonpipeContext) COLON() antlr.TerminalNode {
	return s.GetToken(CGParserCOLON, 0)
}

func (s *NonwsnonpipeContext) SEMICOLON() antlr.TerminalNode {
	return s.GetToken(CGParserSEMICOLON, 0)
}

func (s *NonwsnonpipeContext) LESS_THAN() antlr.TerminalNode {
	return s.GetToken(CGParserLESS_THAN, 0)
}

func (s *NonwsnonpipeContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(CGParserEQUALS, 0)
}

func (s *NonwsnonpipeContext) GREATER_THAN() antlr.TerminalNode {
	return s.GetToken(CGParserGREATER_THAN, 0)
}

func (s *NonwsnonpipeContext) QUESTION() antlr.TerminalNode {
	return s.GetToken(CGParserQUESTION, 0)
}

func (s *NonwsnonpipeContext) AT() antlr.TerminalNode {
	return s.GetToken(CGParserAT, 0)
}

func (s *NonwsnonpipeContext) CAP_A() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_A, 0)
}

func (s *NonwsnonpipeContext) CAP_B() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_B, 0)
}

func (s *NonwsnonpipeContext) CAP_C() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_C, 0)
}

func (s *NonwsnonpipeContext) CAP_D() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_D, 0)
}

func (s *NonwsnonpipeContext) CAP_E() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_E, 0)
}

func (s *NonwsnonpipeContext) CAP_F() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_F, 0)
}

func (s *NonwsnonpipeContext) CAP_G() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_G, 0)
}

func (s *NonwsnonpipeContext) CAP_H() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_H, 0)
}

func (s *NonwsnonpipeContext) CAP_I() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_I, 0)
}

func (s *NonwsnonpipeContext) CAP_J() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_J, 0)
}

func (s *NonwsnonpipeContext) CAP_K() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_K, 0)
}

func (s *NonwsnonpipeContext) CAP_L() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_L, 0)
}

func (s *NonwsnonpipeContext) CAP_M() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_M, 0)
}

func (s *NonwsnonpipeContext) CAP_N() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_N, 0)
}

func (s *NonwsnonpipeContext) CAP_O() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_O, 0)
}

func (s *NonwsnonpipeContext) CAP_P() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_P, 0)
}

func (s *NonwsnonpipeContext) CAP_Q() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_Q, 0)
}

func (s *NonwsnonpipeContext) CAP_R() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_R, 0)
}

func (s *NonwsnonpipeContext) CAP_S() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_S, 0)
}

func (s *NonwsnonpipeContext) CAP_T() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_T, 0)
}

func (s *NonwsnonpipeContext) CAP_U() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_U, 0)
}

func (s *NonwsnonpipeContext) CAP_V() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_V, 0)
}

func (s *NonwsnonpipeContext) CAP_W() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_W, 0)
}

func (s *NonwsnonpipeContext) CAP_X() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_X, 0)
}

func (s *NonwsnonpipeContext) CAP_Y() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_Y, 0)
}

func (s *NonwsnonpipeContext) CAP_Z() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_Z, 0)
}

func (s *NonwsnonpipeContext) LEFT_BRACE() antlr.TerminalNode {
	return s.GetToken(CGParserLEFT_BRACE, 0)
}

func (s *NonwsnonpipeContext) BACKSLASH() antlr.TerminalNode {
	return s.GetToken(CGParserBACKSLASH, 0)
}

func (s *NonwsnonpipeContext) RIGHT_BRACE() antlr.TerminalNode {
	return s.GetToken(CGParserRIGHT_BRACE, 0)
}

func (s *NonwsnonpipeContext) CARAT() antlr.TerminalNode {
	return s.GetToken(CGParserCARAT, 0)
}

func (s *NonwsnonpipeContext) UNDERSCORE() antlr.TerminalNode {
	return s.GetToken(CGParserUNDERSCORE, 0)
}

func (s *NonwsnonpipeContext) ACCENT() antlr.TerminalNode {
	return s.GetToken(CGParserACCENT, 0)
}

func (s *NonwsnonpipeContext) A() antlr.TerminalNode {
	return s.GetToken(CGParserA, 0)
}

func (s *NonwsnonpipeContext) B() antlr.TerminalNode {
	return s.GetToken(CGParserB, 0)
}

func (s *NonwsnonpipeContext) C() antlr.TerminalNode {
	return s.GetToken(CGParserC, 0)
}

func (s *NonwsnonpipeContext) D() antlr.TerminalNode {
	return s.GetToken(CGParserD, 0)
}

func (s *NonwsnonpipeContext) E() antlr.TerminalNode {
	return s.GetToken(CGParserE, 0)
}

func (s *NonwsnonpipeContext) F() antlr.TerminalNode {
	return s.GetToken(CGParserF, 0)
}

func (s *NonwsnonpipeContext) G() antlr.TerminalNode {
	return s.GetToken(CGParserG, 0)
}

func (s *NonwsnonpipeContext) H() antlr.TerminalNode {
	return s.GetToken(CGParserH, 0)
}

func (s *NonwsnonpipeContext) I() antlr.TerminalNode {
	return s.GetToken(CGParserI, 0)
}

func (s *NonwsnonpipeContext) J() antlr.TerminalNode {
	return s.GetToken(CGParserJ, 0)
}

func (s *NonwsnonpipeContext) K() antlr.TerminalNode {
	return s.GetToken(CGParserK, 0)
}

func (s *NonwsnonpipeContext) L() antlr.TerminalNode {
	return s.GetToken(CGParserL, 0)
}

func (s *NonwsnonpipeContext) M() antlr.TerminalNode {
	return s.GetToken(CGParserM, 0)
}

func (s *NonwsnonpipeContext) N() antlr.TerminalNode {
	return s.GetToken(CGParserN, 0)
}

func (s *NonwsnonpipeContext) O() antlr.TerminalNode {
	return s.GetToken(CGParserO, 0)
}

func (s *NonwsnonpipeContext) P() antlr.TerminalNode {
	return s.GetToken(CGParserP, 0)
}

func (s *NonwsnonpipeContext) Q() antlr.TerminalNode {
	return s.GetToken(CGParserQ, 0)
}

func (s *NonwsnonpipeContext) R() antlr.TerminalNode {
	return s.GetToken(CGParserR, 0)
}

func (s *NonwsnonpipeContext) S() antlr.TerminalNode {
	return s.GetToken(CGParserS, 0)
}

func (s *NonwsnonpipeContext) T() antlr.TerminalNode {
	return s.GetToken(CGParserT, 0)
}

func (s *NonwsnonpipeContext) U() antlr.TerminalNode {
	return s.GetToken(CGParserU, 0)
}

func (s *NonwsnonpipeContext) V() antlr.TerminalNode {
	return s.GetToken(CGParserV, 0)
}

func (s *NonwsnonpipeContext) W() antlr.TerminalNode {
	return s.GetToken(CGParserW, 0)
}

func (s *NonwsnonpipeContext) X() antlr.TerminalNode {
	return s.GetToken(CGParserX, 0)
}

func (s *NonwsnonpipeContext) Y() antlr.TerminalNode {
	return s.GetToken(CGParserY, 0)
}

func (s *NonwsnonpipeContext) Z() antlr.TerminalNode {
	return s.GetToken(CGParserZ, 0)
}

func (s *NonwsnonpipeContext) LEFT_CURLY_BRACE() antlr.TerminalNode {
	return s.GetToken(CGParserLEFT_CURLY_BRACE, 0)
}

func (s *NonwsnonpipeContext) RIGHT_CURLY_BRACE() antlr.TerminalNode {
	return s.GetToken(CGParserRIGHT_CURLY_BRACE, 0)
}

func (s *NonwsnonpipeContext) TILDE() antlr.TerminalNode {
	return s.GetToken(CGParserTILDE, 0)
}

func (s *NonwsnonpipeContext) Utf8_2() IUtf8_2Context {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUtf8_2Context)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUtf8_2Context)
}

func (s *NonwsnonpipeContext) Utf8_3() IUtf8_3Context {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUtf8_3Context)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUtf8_3Context)
}

func (s *NonwsnonpipeContext) Utf8_4() IUtf8_4Context {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUtf8_4Context)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUtf8_4Context)
}

func (s *NonwsnonpipeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NonwsnonpipeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NonwsnonpipeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterNonwsnonpipe(s)
	}
}

func (s *NonwsnonpipeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitNonwsnonpipe(s)
	}
}

func (p *CGParser) Nonwsnonpipe() (localctx INonwsnonpipeContext) {
	localctx = NewNonwsnonpipeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, CGParserRULE_nonwsnonpipe)
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

	p.SetState(361)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CGParserEXCLAMATION, CGParserQUOTE, CGParserPOUND, CGParserDOLLAR, CGParserPERCENT, CGParserAMPERSAND, CGParserAPOSTROPHE, CGParserLEFT_PAREN, CGParserRIGHT_PAREN, CGParserASTERISK, CGParserPLUS, CGParserCOMMA, CGParserDASH, CGParserPERIOD, CGParserSLASH, CGParserZERO, CGParserONE, CGParserTWO, CGParserTHREE, CGParserFOUR, CGParserFIVE, CGParserSIX, CGParserSEVEN, CGParserEIGHT, CGParserNINE, CGParserCOLON, CGParserSEMICOLON, CGParserLESS_THAN, CGParserEQUALS, CGParserGREATER_THAN, CGParserQUESTION, CGParserAT, CGParserCAP_A, CGParserCAP_B, CGParserCAP_C, CGParserCAP_D, CGParserCAP_E, CGParserCAP_F, CGParserCAP_G, CGParserCAP_H, CGParserCAP_I, CGParserCAP_J, CGParserCAP_K, CGParserCAP_L, CGParserCAP_M, CGParserCAP_N, CGParserCAP_O, CGParserCAP_P, CGParserCAP_Q, CGParserCAP_R, CGParserCAP_S, CGParserCAP_T, CGParserCAP_U, CGParserCAP_V, CGParserCAP_W, CGParserCAP_X, CGParserCAP_Y, CGParserCAP_Z, CGParserLEFT_BRACE, CGParserBACKSLASH, CGParserRIGHT_BRACE, CGParserCARAT, CGParserUNDERSCORE, CGParserACCENT, CGParserA, CGParserB, CGParserC, CGParserD, CGParserE, CGParserF, CGParserG, CGParserH, CGParserI, CGParserJ, CGParserK, CGParserL, CGParserM, CGParserN, CGParserO, CGParserP, CGParserQ, CGParserR, CGParserS, CGParserT, CGParserU, CGParserV, CGParserW, CGParserX, CGParserY, CGParserZ, CGParserLEFT_CURLY_BRACE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(356)
			_la = p.GetTokenStream().LA(1)

			if !((((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CGParserEXCLAMATION)|(1<<CGParserQUOTE)|(1<<CGParserPOUND)|(1<<CGParserDOLLAR)|(1<<CGParserPERCENT)|(1<<CGParserAMPERSAND)|(1<<CGParserAPOSTROPHE)|(1<<CGParserLEFT_PAREN)|(1<<CGParserRIGHT_PAREN)|(1<<CGParserASTERISK)|(1<<CGParserPLUS)|(1<<CGParserCOMMA)|(1<<CGParserDASH)|(1<<CGParserPERIOD)|(1<<CGParserSLASH)|(1<<CGParserZERO)|(1<<CGParserONE)|(1<<CGParserTWO)|(1<<CGParserTHREE)|(1<<CGParserFOUR)|(1<<CGParserFIVE)|(1<<CGParserSIX)|(1<<CGParserSEVEN)|(1<<CGParserEIGHT)|(1<<CGParserNINE)|(1<<CGParserCOLON)|(1<<CGParserSEMICOLON))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(CGParserLESS_THAN-32))|(1<<(CGParserEQUALS-32))|(1<<(CGParserGREATER_THAN-32))|(1<<(CGParserQUESTION-32))|(1<<(CGParserAT-32))|(1<<(CGParserCAP_A-32))|(1<<(CGParserCAP_B-32))|(1<<(CGParserCAP_C-32))|(1<<(CGParserCAP_D-32))|(1<<(CGParserCAP_E-32))|(1<<(CGParserCAP_F-32))|(1<<(CGParserCAP_G-32))|(1<<(CGParserCAP_H-32))|(1<<(CGParserCAP_I-32))|(1<<(CGParserCAP_J-32))|(1<<(CGParserCAP_K-32))|(1<<(CGParserCAP_L-32))|(1<<(CGParserCAP_M-32))|(1<<(CGParserCAP_N-32))|(1<<(CGParserCAP_O-32))|(1<<(CGParserCAP_P-32))|(1<<(CGParserCAP_Q-32))|(1<<(CGParserCAP_R-32))|(1<<(CGParserCAP_S-32))|(1<<(CGParserCAP_T-32))|(1<<(CGParserCAP_U-32))|(1<<(CGParserCAP_V-32))|(1<<(CGParserCAP_W-32))|(1<<(CGParserCAP_X-32))|(1<<(CGParserCAP_Y-32))|(1<<(CGParserCAP_Z-32))|(1<<(CGParserLEFT_BRACE-32)))) != 0) || (((_la-64)&-(0x1f+1)) == 0 && ((1<<uint((_la-64)))&((1<<(CGParserBACKSLASH-64))|(1<<(CGParserRIGHT_BRACE-64))|(1<<(CGParserCARAT-64))|(1<<(CGParserUNDERSCORE-64))|(1<<(CGParserACCENT-64))|(1<<(CGParserA-64))|(1<<(CGParserB-64))|(1<<(CGParserC-64))|(1<<(CGParserD-64))|(1<<(CGParserE-64))|(1<<(CGParserF-64))|(1<<(CGParserG-64))|(1<<(CGParserH-64))|(1<<(CGParserI-64))|(1<<(CGParserJ-64))|(1<<(CGParserK-64))|(1<<(CGParserL-64))|(1<<(CGParserM-64))|(1<<(CGParserN-64))|(1<<(CGParserO-64))|(1<<(CGParserP-64))|(1<<(CGParserQ-64))|(1<<(CGParserR-64))|(1<<(CGParserS-64))|(1<<(CGParserT-64))|(1<<(CGParserU-64))|(1<<(CGParserV-64))|(1<<(CGParserW-64))|(1<<(CGParserX-64))|(1<<(CGParserY-64))|(1<<(CGParserZ-64))|(1<<(CGParserLEFT_CURLY_BRACE-64)))) != 0)) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case CGParserRIGHT_CURLY_BRACE, CGParserTILDE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(357)
			_la = p.GetTokenStream().LA(1)

			if !(_la == CGParserRIGHT_CURLY_BRACE || _la == CGParserTILDE) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case CGParserU_00C2, CGParserU_00C3, CGParserU_00C4, CGParserU_00C5, CGParserU_00C6, CGParserU_00C7, CGParserU_00C8, CGParserU_00C9, CGParserU_00CA, CGParserU_00CB, CGParserU_00CC, CGParserU_00CD, CGParserU_00CE, CGParserU_00CF, CGParserU_00D0, CGParserU_00D1, CGParserU_00D2, CGParserU_00D3, CGParserU_00D4, CGParserU_00D5, CGParserU_00D6, CGParserU_00D7, CGParserU_00D8, CGParserU_00D9, CGParserU_00DA, CGParserU_00DB, CGParserU_00DC, CGParserU_00DD, CGParserU_00DE, CGParserU_00DF:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(358)
			p.Utf8_2()
		}

	case CGParserU_00E0, CGParserU_00E1, CGParserU_00E2, CGParserU_00E3, CGParserU_00E4, CGParserU_00E5, CGParserU_00E6, CGParserU_00E7, CGParserU_00E8, CGParserU_00E9, CGParserU_00EA, CGParserU_00EB, CGParserU_00EC, CGParserU_00ED, CGParserU_00EE, CGParserU_00EF:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(359)
			p.Utf8_3()
		}

	case CGParserU_00F0, CGParserU_00F1, CGParserU_00F2, CGParserU_00F3, CGParserU_00F4:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(360)
			p.Utf8_4()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IAnynonescapedcharContext is an interface to support dynamic dispatch.
type IAnynonescapedcharContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAnynonescapedcharContext differentiates from other interfaces.
	IsAnynonescapedcharContext()
}

type AnynonescapedcharContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAnynonescapedcharContext() *AnynonescapedcharContext {
	var p = new(AnynonescapedcharContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_anynonescapedchar
	return p
}

func (*AnynonescapedcharContext) IsAnynonescapedcharContext() {}

func NewAnynonescapedcharContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AnynonescapedcharContext {
	var p = new(AnynonescapedcharContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_anynonescapedchar

	return p
}

func (s *AnynonescapedcharContext) GetParser() antlr.Parser { return s.parser }

func (s *AnynonescapedcharContext) Htab() IHtabContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IHtabContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IHtabContext)
}

func (s *AnynonescapedcharContext) Cr() ICrContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ICrContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ICrContext)
}

func (s *AnynonescapedcharContext) Lf() ILfContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*ILfContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(ILfContext)
}

func (s *AnynonescapedcharContext) SPACE() antlr.TerminalNode {
	return s.GetToken(CGParserSPACE, 0)
}

func (s *AnynonescapedcharContext) EXCLAMATION() antlr.TerminalNode {
	return s.GetToken(CGParserEXCLAMATION, 0)
}

func (s *AnynonescapedcharContext) POUND() antlr.TerminalNode {
	return s.GetToken(CGParserPOUND, 0)
}

func (s *AnynonescapedcharContext) DOLLAR() antlr.TerminalNode {
	return s.GetToken(CGParserDOLLAR, 0)
}

func (s *AnynonescapedcharContext) PERCENT() antlr.TerminalNode {
	return s.GetToken(CGParserPERCENT, 0)
}

func (s *AnynonescapedcharContext) AMPERSAND() antlr.TerminalNode {
	return s.GetToken(CGParserAMPERSAND, 0)
}

func (s *AnynonescapedcharContext) APOSTROPHE() antlr.TerminalNode {
	return s.GetToken(CGParserAPOSTROPHE, 0)
}

func (s *AnynonescapedcharContext) LEFT_PAREN() antlr.TerminalNode {
	return s.GetToken(CGParserLEFT_PAREN, 0)
}

func (s *AnynonescapedcharContext) RIGHT_PAREN() antlr.TerminalNode {
	return s.GetToken(CGParserRIGHT_PAREN, 0)
}

func (s *AnynonescapedcharContext) ASTERISK() antlr.TerminalNode {
	return s.GetToken(CGParserASTERISK, 0)
}

func (s *AnynonescapedcharContext) PLUS() antlr.TerminalNode {
	return s.GetToken(CGParserPLUS, 0)
}

func (s *AnynonescapedcharContext) COMMA() antlr.TerminalNode {
	return s.GetToken(CGParserCOMMA, 0)
}

func (s *AnynonescapedcharContext) DASH() antlr.TerminalNode {
	return s.GetToken(CGParserDASH, 0)
}

func (s *AnynonescapedcharContext) PERIOD() antlr.TerminalNode {
	return s.GetToken(CGParserPERIOD, 0)
}

func (s *AnynonescapedcharContext) SLASH() antlr.TerminalNode {
	return s.GetToken(CGParserSLASH, 0)
}

func (s *AnynonescapedcharContext) ZERO() antlr.TerminalNode {
	return s.GetToken(CGParserZERO, 0)
}

func (s *AnynonescapedcharContext) ONE() antlr.TerminalNode {
	return s.GetToken(CGParserONE, 0)
}

func (s *AnynonescapedcharContext) TWO() antlr.TerminalNode {
	return s.GetToken(CGParserTWO, 0)
}

func (s *AnynonescapedcharContext) THREE() antlr.TerminalNode {
	return s.GetToken(CGParserTHREE, 0)
}

func (s *AnynonescapedcharContext) FOUR() antlr.TerminalNode {
	return s.GetToken(CGParserFOUR, 0)
}

func (s *AnynonescapedcharContext) FIVE() antlr.TerminalNode {
	return s.GetToken(CGParserFIVE, 0)
}

func (s *AnynonescapedcharContext) SIX() antlr.TerminalNode {
	return s.GetToken(CGParserSIX, 0)
}

func (s *AnynonescapedcharContext) SEVEN() antlr.TerminalNode {
	return s.GetToken(CGParserSEVEN, 0)
}

func (s *AnynonescapedcharContext) EIGHT() antlr.TerminalNode {
	return s.GetToken(CGParserEIGHT, 0)
}

func (s *AnynonescapedcharContext) NINE() antlr.TerminalNode {
	return s.GetToken(CGParserNINE, 0)
}

func (s *AnynonescapedcharContext) COLON() antlr.TerminalNode {
	return s.GetToken(CGParserCOLON, 0)
}

func (s *AnynonescapedcharContext) SEMICOLON() antlr.TerminalNode {
	return s.GetToken(CGParserSEMICOLON, 0)
}

func (s *AnynonescapedcharContext) LESS_THAN() antlr.TerminalNode {
	return s.GetToken(CGParserLESS_THAN, 0)
}

func (s *AnynonescapedcharContext) EQUALS() antlr.TerminalNode {
	return s.GetToken(CGParserEQUALS, 0)
}

func (s *AnynonescapedcharContext) GREATER_THAN() antlr.TerminalNode {
	return s.GetToken(CGParserGREATER_THAN, 0)
}

func (s *AnynonescapedcharContext) QUESTION() antlr.TerminalNode {
	return s.GetToken(CGParserQUESTION, 0)
}

func (s *AnynonescapedcharContext) AT() antlr.TerminalNode {
	return s.GetToken(CGParserAT, 0)
}

func (s *AnynonescapedcharContext) CAP_A() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_A, 0)
}

func (s *AnynonescapedcharContext) CAP_B() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_B, 0)
}

func (s *AnynonescapedcharContext) CAP_C() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_C, 0)
}

func (s *AnynonescapedcharContext) CAP_D() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_D, 0)
}

func (s *AnynonescapedcharContext) CAP_E() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_E, 0)
}

func (s *AnynonescapedcharContext) CAP_F() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_F, 0)
}

func (s *AnynonescapedcharContext) CAP_G() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_G, 0)
}

func (s *AnynonescapedcharContext) CAP_H() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_H, 0)
}

func (s *AnynonescapedcharContext) CAP_I() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_I, 0)
}

func (s *AnynonescapedcharContext) CAP_J() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_J, 0)
}

func (s *AnynonescapedcharContext) CAP_K() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_K, 0)
}

func (s *AnynonescapedcharContext) CAP_L() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_L, 0)
}

func (s *AnynonescapedcharContext) CAP_M() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_M, 0)
}

func (s *AnynonescapedcharContext) CAP_N() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_N, 0)
}

func (s *AnynonescapedcharContext) CAP_O() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_O, 0)
}

func (s *AnynonescapedcharContext) CAP_P() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_P, 0)
}

func (s *AnynonescapedcharContext) CAP_Q() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_Q, 0)
}

func (s *AnynonescapedcharContext) CAP_R() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_R, 0)
}

func (s *AnynonescapedcharContext) CAP_S() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_S, 0)
}

func (s *AnynonescapedcharContext) CAP_T() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_T, 0)
}

func (s *AnynonescapedcharContext) CAP_U() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_U, 0)
}

func (s *AnynonescapedcharContext) CAP_V() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_V, 0)
}

func (s *AnynonescapedcharContext) CAP_W() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_W, 0)
}

func (s *AnynonescapedcharContext) CAP_X() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_X, 0)
}

func (s *AnynonescapedcharContext) CAP_Y() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_Y, 0)
}

func (s *AnynonescapedcharContext) CAP_Z() antlr.TerminalNode {
	return s.GetToken(CGParserCAP_Z, 0)
}

func (s *AnynonescapedcharContext) LEFT_BRACE() antlr.TerminalNode {
	return s.GetToken(CGParserLEFT_BRACE, 0)
}

func (s *AnynonescapedcharContext) RIGHT_BRACE() antlr.TerminalNode {
	return s.GetToken(CGParserRIGHT_BRACE, 0)
}

func (s *AnynonescapedcharContext) CARAT() antlr.TerminalNode {
	return s.GetToken(CGParserCARAT, 0)
}

func (s *AnynonescapedcharContext) UNDERSCORE() antlr.TerminalNode {
	return s.GetToken(CGParserUNDERSCORE, 0)
}

func (s *AnynonescapedcharContext) ACCENT() antlr.TerminalNode {
	return s.GetToken(CGParserACCENT, 0)
}

func (s *AnynonescapedcharContext) A() antlr.TerminalNode {
	return s.GetToken(CGParserA, 0)
}

func (s *AnynonescapedcharContext) B() antlr.TerminalNode {
	return s.GetToken(CGParserB, 0)
}

func (s *AnynonescapedcharContext) C() antlr.TerminalNode {
	return s.GetToken(CGParserC, 0)
}

func (s *AnynonescapedcharContext) D() antlr.TerminalNode {
	return s.GetToken(CGParserD, 0)
}

func (s *AnynonescapedcharContext) E() antlr.TerminalNode {
	return s.GetToken(CGParserE, 0)
}

func (s *AnynonescapedcharContext) F() antlr.TerminalNode {
	return s.GetToken(CGParserF, 0)
}

func (s *AnynonescapedcharContext) G() antlr.TerminalNode {
	return s.GetToken(CGParserG, 0)
}

func (s *AnynonescapedcharContext) H() antlr.TerminalNode {
	return s.GetToken(CGParserH, 0)
}

func (s *AnynonescapedcharContext) I() antlr.TerminalNode {
	return s.GetToken(CGParserI, 0)
}

func (s *AnynonescapedcharContext) J() antlr.TerminalNode {
	return s.GetToken(CGParserJ, 0)
}

func (s *AnynonescapedcharContext) K() antlr.TerminalNode {
	return s.GetToken(CGParserK, 0)
}

func (s *AnynonescapedcharContext) L() antlr.TerminalNode {
	return s.GetToken(CGParserL, 0)
}

func (s *AnynonescapedcharContext) M() antlr.TerminalNode {
	return s.GetToken(CGParserM, 0)
}

func (s *AnynonescapedcharContext) N() antlr.TerminalNode {
	return s.GetToken(CGParserN, 0)
}

func (s *AnynonescapedcharContext) O() antlr.TerminalNode {
	return s.GetToken(CGParserO, 0)
}

func (s *AnynonescapedcharContext) P() antlr.TerminalNode {
	return s.GetToken(CGParserP, 0)
}

func (s *AnynonescapedcharContext) Q() antlr.TerminalNode {
	return s.GetToken(CGParserQ, 0)
}

func (s *AnynonescapedcharContext) R() antlr.TerminalNode {
	return s.GetToken(CGParserR, 0)
}

func (s *AnynonescapedcharContext) S() antlr.TerminalNode {
	return s.GetToken(CGParserS, 0)
}

func (s *AnynonescapedcharContext) T() antlr.TerminalNode {
	return s.GetToken(CGParserT, 0)
}

func (s *AnynonescapedcharContext) U() antlr.TerminalNode {
	return s.GetToken(CGParserU, 0)
}

func (s *AnynonescapedcharContext) V() antlr.TerminalNode {
	return s.GetToken(CGParserV, 0)
}

func (s *AnynonescapedcharContext) W() antlr.TerminalNode {
	return s.GetToken(CGParserW, 0)
}

func (s *AnynonescapedcharContext) X() antlr.TerminalNode {
	return s.GetToken(CGParserX, 0)
}

func (s *AnynonescapedcharContext) Y() antlr.TerminalNode {
	return s.GetToken(CGParserY, 0)
}

func (s *AnynonescapedcharContext) Z() antlr.TerminalNode {
	return s.GetToken(CGParserZ, 0)
}

func (s *AnynonescapedcharContext) LEFT_CURLY_BRACE() antlr.TerminalNode {
	return s.GetToken(CGParserLEFT_CURLY_BRACE, 0)
}

func (s *AnynonescapedcharContext) PIPE() antlr.TerminalNode {
	return s.GetToken(CGParserPIPE, 0)
}

func (s *AnynonescapedcharContext) RIGHT_CURLY_BRACE() antlr.TerminalNode {
	return s.GetToken(CGParserRIGHT_CURLY_BRACE, 0)
}

func (s *AnynonescapedcharContext) TILDE() antlr.TerminalNode {
	return s.GetToken(CGParserTILDE, 0)
}

func (s *AnynonescapedcharContext) Utf8_2() IUtf8_2Context {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUtf8_2Context)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUtf8_2Context)
}

func (s *AnynonescapedcharContext) Utf8_3() IUtf8_3Context {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUtf8_3Context)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUtf8_3Context)
}

func (s *AnynonescapedcharContext) Utf8_4() IUtf8_4Context {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUtf8_4Context)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUtf8_4Context)
}

func (s *AnynonescapedcharContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AnynonescapedcharContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AnynonescapedcharContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterAnynonescapedchar(s)
	}
}

func (s *AnynonescapedcharContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitAnynonescapedchar(s)
	}
}

func (p *CGParser) Anynonescapedchar() (localctx IAnynonescapedcharContext) {
	localctx = NewAnynonescapedcharContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, CGParserRULE_anynonescapedchar)
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

	p.SetState(372)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CGParserTAB:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(363)
			p.Htab()
		}

	case CGParserCR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(364)
			p.Cr()
		}

	case CGParserLF:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(365)
			p.Lf()
		}

	case CGParserSPACE, CGParserEXCLAMATION:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(366)
			_la = p.GetTokenStream().LA(1)

			if !(_la == CGParserSPACE || _la == CGParserEXCLAMATION) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case CGParserPOUND, CGParserDOLLAR, CGParserPERCENT, CGParserAMPERSAND, CGParserAPOSTROPHE, CGParserLEFT_PAREN, CGParserRIGHT_PAREN, CGParserASTERISK, CGParserPLUS, CGParserCOMMA, CGParserDASH, CGParserPERIOD, CGParserSLASH, CGParserZERO, CGParserONE, CGParserTWO, CGParserTHREE, CGParserFOUR, CGParserFIVE, CGParserSIX, CGParserSEVEN, CGParserEIGHT, CGParserNINE, CGParserCOLON, CGParserSEMICOLON, CGParserLESS_THAN, CGParserEQUALS, CGParserGREATER_THAN, CGParserQUESTION, CGParserAT, CGParserCAP_A, CGParserCAP_B, CGParserCAP_C, CGParserCAP_D, CGParserCAP_E, CGParserCAP_F, CGParserCAP_G, CGParserCAP_H, CGParserCAP_I, CGParserCAP_J, CGParserCAP_K, CGParserCAP_L, CGParserCAP_M, CGParserCAP_N, CGParserCAP_O, CGParserCAP_P, CGParserCAP_Q, CGParserCAP_R, CGParserCAP_S, CGParserCAP_T, CGParserCAP_U, CGParserCAP_V, CGParserCAP_W, CGParserCAP_X, CGParserCAP_Y, CGParserCAP_Z, CGParserLEFT_BRACE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(367)
			_la = p.GetTokenStream().LA(1)

			if !((((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<CGParserPOUND)|(1<<CGParserDOLLAR)|(1<<CGParserPERCENT)|(1<<CGParserAMPERSAND)|(1<<CGParserAPOSTROPHE)|(1<<CGParserLEFT_PAREN)|(1<<CGParserRIGHT_PAREN)|(1<<CGParserASTERISK)|(1<<CGParserPLUS)|(1<<CGParserCOMMA)|(1<<CGParserDASH)|(1<<CGParserPERIOD)|(1<<CGParserSLASH)|(1<<CGParserZERO)|(1<<CGParserONE)|(1<<CGParserTWO)|(1<<CGParserTHREE)|(1<<CGParserFOUR)|(1<<CGParserFIVE)|(1<<CGParserSIX)|(1<<CGParserSEVEN)|(1<<CGParserEIGHT)|(1<<CGParserNINE)|(1<<CGParserCOLON)|(1<<CGParserSEMICOLON))) != 0) || (((_la-32)&-(0x1f+1)) == 0 && ((1<<uint((_la-32)))&((1<<(CGParserLESS_THAN-32))|(1<<(CGParserEQUALS-32))|(1<<(CGParserGREATER_THAN-32))|(1<<(CGParserQUESTION-32))|(1<<(CGParserAT-32))|(1<<(CGParserCAP_A-32))|(1<<(CGParserCAP_B-32))|(1<<(CGParserCAP_C-32))|(1<<(CGParserCAP_D-32))|(1<<(CGParserCAP_E-32))|(1<<(CGParserCAP_F-32))|(1<<(CGParserCAP_G-32))|(1<<(CGParserCAP_H-32))|(1<<(CGParserCAP_I-32))|(1<<(CGParserCAP_J-32))|(1<<(CGParserCAP_K-32))|(1<<(CGParserCAP_L-32))|(1<<(CGParserCAP_M-32))|(1<<(CGParserCAP_N-32))|(1<<(CGParserCAP_O-32))|(1<<(CGParserCAP_P-32))|(1<<(CGParserCAP_Q-32))|(1<<(CGParserCAP_R-32))|(1<<(CGParserCAP_S-32))|(1<<(CGParserCAP_T-32))|(1<<(CGParserCAP_U-32))|(1<<(CGParserCAP_V-32))|(1<<(CGParserCAP_W-32))|(1<<(CGParserCAP_X-32))|(1<<(CGParserCAP_Y-32))|(1<<(CGParserCAP_Z-32))|(1<<(CGParserLEFT_BRACE-32)))) != 0)) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case CGParserRIGHT_BRACE, CGParserCARAT, CGParserUNDERSCORE, CGParserACCENT, CGParserA, CGParserB, CGParserC, CGParserD, CGParserE, CGParserF, CGParserG, CGParserH, CGParserI, CGParserJ, CGParserK, CGParserL, CGParserM, CGParserN, CGParserO, CGParserP, CGParserQ, CGParserR, CGParserS, CGParserT, CGParserU, CGParserV, CGParserW, CGParserX, CGParserY, CGParserZ, CGParserLEFT_CURLY_BRACE, CGParserPIPE, CGParserRIGHT_CURLY_BRACE, CGParserTILDE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(368)
			_la = p.GetTokenStream().LA(1)

			if !((((_la-65)&-(0x1f+1)) == 0 && ((1<<uint((_la-65)))&((1<<(CGParserRIGHT_BRACE-65))|(1<<(CGParserCARAT-65))|(1<<(CGParserUNDERSCORE-65))|(1<<(CGParserACCENT-65))|(1<<(CGParserA-65))|(1<<(CGParserB-65))|(1<<(CGParserC-65))|(1<<(CGParserD-65))|(1<<(CGParserE-65))|(1<<(CGParserF-65))|(1<<(CGParserG-65))|(1<<(CGParserH-65))|(1<<(CGParserI-65))|(1<<(CGParserJ-65))|(1<<(CGParserK-65))|(1<<(CGParserL-65))|(1<<(CGParserM-65))|(1<<(CGParserN-65))|(1<<(CGParserO-65))|(1<<(CGParserP-65))|(1<<(CGParserQ-65))|(1<<(CGParserR-65))|(1<<(CGParserS-65))|(1<<(CGParserT-65))|(1<<(CGParserU-65))|(1<<(CGParserV-65))|(1<<(CGParserW-65))|(1<<(CGParserX-65))|(1<<(CGParserY-65))|(1<<(CGParserZ-65))|(1<<(CGParserLEFT_CURLY_BRACE-65))|(1<<(CGParserPIPE-65)))) != 0) || _la == CGParserRIGHT_CURLY_BRACE || _la == CGParserTILDE) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case CGParserU_00C2, CGParserU_00C3, CGParserU_00C4, CGParserU_00C5, CGParserU_00C6, CGParserU_00C7, CGParserU_00C8, CGParserU_00C9, CGParserU_00CA, CGParserU_00CB, CGParserU_00CC, CGParserU_00CD, CGParserU_00CE, CGParserU_00CF, CGParserU_00D0, CGParserU_00D1, CGParserU_00D2, CGParserU_00D3, CGParserU_00D4, CGParserU_00D5, CGParserU_00D6, CGParserU_00D7, CGParserU_00D8, CGParserU_00D9, CGParserU_00DA, CGParserU_00DB, CGParserU_00DC, CGParserU_00DD, CGParserU_00DE, CGParserU_00DF:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(369)
			p.Utf8_2()
		}

	case CGParserU_00E0, CGParserU_00E1, CGParserU_00E2, CGParserU_00E3, CGParserU_00E4, CGParserU_00E5, CGParserU_00E6, CGParserU_00E7, CGParserU_00E8, CGParserU_00E9, CGParserU_00EA, CGParserU_00EB, CGParserU_00EC, CGParserU_00ED, CGParserU_00EE, CGParserU_00EF:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(370)
			p.Utf8_3()
		}

	case CGParserU_00F0, CGParserU_00F1, CGParserU_00F2, CGParserU_00F3, CGParserU_00F4:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(371)
			p.Utf8_4()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IEscapedcharContext is an interface to support dynamic dispatch.
type IEscapedcharContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsEscapedcharContext differentiates from other interfaces.
	IsEscapedcharContext()
}

type EscapedcharContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyEscapedcharContext() *EscapedcharContext {
	var p = new(EscapedcharContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_escapedchar
	return p
}

func (*EscapedcharContext) IsEscapedcharContext() {}

func NewEscapedcharContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EscapedcharContext {
	var p = new(EscapedcharContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_escapedchar

	return p
}

func (s *EscapedcharContext) GetParser() antlr.Parser { return s.parser }

func (s *EscapedcharContext) AllBs() []IBsContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IBsContext)(nil)).Elem())
	var tst = make([]IBsContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IBsContext)
		}
	}

	return tst
}

func (s *EscapedcharContext) Bs(i int) IBsContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IBsContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IBsContext)
}

func (s *EscapedcharContext) Qm() IQmContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IQmContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IQmContext)
}

func (s *EscapedcharContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EscapedcharContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *EscapedcharContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterEscapedchar(s)
	}
}

func (s *EscapedcharContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitEscapedchar(s)
	}
}

func (p *CGParser) Escapedchar() (localctx IEscapedcharContext) {
	localctx = NewEscapedcharContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, CGParserRULE_escapedchar)

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

	p.SetState(380)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(374)
			p.Bs()
		}
		{
			p.SetState(375)
			p.Qm()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(377)
			p.Bs()
		}
		{
			p.SetState(378)
			p.Bs()
		}

	}

	return localctx
}

// IUtf8_2Context is an interface to support dynamic dispatch.
type IUtf8_2Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUtf8_2Context differentiates from other interfaces.
	IsUtf8_2Context()
}

type Utf8_2Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUtf8_2Context() *Utf8_2Context {
	var p = new(Utf8_2Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_utf8_2
	return p
}

func (*Utf8_2Context) IsUtf8_2Context() {}

func NewUtf8_2Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Utf8_2Context {
	var p = new(Utf8_2Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_utf8_2

	return p
}

func (s *Utf8_2Context) GetParser() antlr.Parser { return s.parser }

func (s *Utf8_2Context) Utf8_tail() IUtf8_tailContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUtf8_tailContext)(nil)).Elem(), 0)

	if t == nil {
		return nil
	}

	return t.(IUtf8_tailContext)
}

func (s *Utf8_2Context) U_00C2() antlr.TerminalNode {
	return s.GetToken(CGParserU_00C2, 0)
}

func (s *Utf8_2Context) U_00C3() antlr.TerminalNode {
	return s.GetToken(CGParserU_00C3, 0)
}

func (s *Utf8_2Context) U_00C4() antlr.TerminalNode {
	return s.GetToken(CGParserU_00C4, 0)
}

func (s *Utf8_2Context) U_00C5() antlr.TerminalNode {
	return s.GetToken(CGParserU_00C5, 0)
}

func (s *Utf8_2Context) U_00C6() antlr.TerminalNode {
	return s.GetToken(CGParserU_00C6, 0)
}

func (s *Utf8_2Context) U_00C7() antlr.TerminalNode {
	return s.GetToken(CGParserU_00C7, 0)
}

func (s *Utf8_2Context) U_00C8() antlr.TerminalNode {
	return s.GetToken(CGParserU_00C8, 0)
}

func (s *Utf8_2Context) U_00C9() antlr.TerminalNode {
	return s.GetToken(CGParserU_00C9, 0)
}

func (s *Utf8_2Context) U_00CA() antlr.TerminalNode {
	return s.GetToken(CGParserU_00CA, 0)
}

func (s *Utf8_2Context) U_00CB() antlr.TerminalNode {
	return s.GetToken(CGParserU_00CB, 0)
}

func (s *Utf8_2Context) U_00CC() antlr.TerminalNode {
	return s.GetToken(CGParserU_00CC, 0)
}

func (s *Utf8_2Context) U_00CD() antlr.TerminalNode {
	return s.GetToken(CGParserU_00CD, 0)
}

func (s *Utf8_2Context) U_00CE() antlr.TerminalNode {
	return s.GetToken(CGParserU_00CE, 0)
}

func (s *Utf8_2Context) U_00CF() antlr.TerminalNode {
	return s.GetToken(CGParserU_00CF, 0)
}

func (s *Utf8_2Context) U_00D0() antlr.TerminalNode {
	return s.GetToken(CGParserU_00D0, 0)
}

func (s *Utf8_2Context) U_00D1() antlr.TerminalNode {
	return s.GetToken(CGParserU_00D1, 0)
}

func (s *Utf8_2Context) U_00D2() antlr.TerminalNode {
	return s.GetToken(CGParserU_00D2, 0)
}

func (s *Utf8_2Context) U_00D3() antlr.TerminalNode {
	return s.GetToken(CGParserU_00D3, 0)
}

func (s *Utf8_2Context) U_00D4() antlr.TerminalNode {
	return s.GetToken(CGParserU_00D4, 0)
}

func (s *Utf8_2Context) U_00D5() antlr.TerminalNode {
	return s.GetToken(CGParserU_00D5, 0)
}

func (s *Utf8_2Context) U_00D6() antlr.TerminalNode {
	return s.GetToken(CGParserU_00D6, 0)
}

func (s *Utf8_2Context) U_00D7() antlr.TerminalNode {
	return s.GetToken(CGParserU_00D7, 0)
}

func (s *Utf8_2Context) U_00D8() antlr.TerminalNode {
	return s.GetToken(CGParserU_00D8, 0)
}

func (s *Utf8_2Context) U_00D9() antlr.TerminalNode {
	return s.GetToken(CGParserU_00D9, 0)
}

func (s *Utf8_2Context) U_00DA() antlr.TerminalNode {
	return s.GetToken(CGParserU_00DA, 0)
}

func (s *Utf8_2Context) U_00DB() antlr.TerminalNode {
	return s.GetToken(CGParserU_00DB, 0)
}

func (s *Utf8_2Context) U_00DC() antlr.TerminalNode {
	return s.GetToken(CGParserU_00DC, 0)
}

func (s *Utf8_2Context) U_00DD() antlr.TerminalNode {
	return s.GetToken(CGParserU_00DD, 0)
}

func (s *Utf8_2Context) U_00DE() antlr.TerminalNode {
	return s.GetToken(CGParserU_00DE, 0)
}

func (s *Utf8_2Context) U_00DF() antlr.TerminalNode {
	return s.GetToken(CGParserU_00DF, 0)
}

func (s *Utf8_2Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Utf8_2Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Utf8_2Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterUtf8_2(s)
	}
}

func (s *Utf8_2Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitUtf8_2(s)
	}
}

func (p *CGParser) Utf8_2() (localctx IUtf8_2Context) {
	localctx = NewUtf8_2Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, CGParserRULE_utf8_2)
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
		p.SetState(382)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-163)&-(0x1f+1)) == 0 && ((1<<uint((_la-163)))&((1<<(CGParserU_00C2-163))|(1<<(CGParserU_00C3-163))|(1<<(CGParserU_00C4-163))|(1<<(CGParserU_00C5-163))|(1<<(CGParserU_00C6-163))|(1<<(CGParserU_00C7-163))|(1<<(CGParserU_00C8-163))|(1<<(CGParserU_00C9-163))|(1<<(CGParserU_00CA-163))|(1<<(CGParserU_00CB-163))|(1<<(CGParserU_00CC-163))|(1<<(CGParserU_00CD-163))|(1<<(CGParserU_00CE-163))|(1<<(CGParserU_00CF-163))|(1<<(CGParserU_00D0-163))|(1<<(CGParserU_00D1-163))|(1<<(CGParserU_00D2-163))|(1<<(CGParserU_00D3-163))|(1<<(CGParserU_00D4-163))|(1<<(CGParserU_00D5-163))|(1<<(CGParserU_00D6-163))|(1<<(CGParserU_00D7-163))|(1<<(CGParserU_00D8-163))|(1<<(CGParserU_00D9-163))|(1<<(CGParserU_00DA-163))|(1<<(CGParserU_00DB-163))|(1<<(CGParserU_00DC-163))|(1<<(CGParserU_00DD-163))|(1<<(CGParserU_00DE-163))|(1<<(CGParserU_00DF-163)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(383)
		p.Utf8_tail()
	}

	return localctx
}

// IUtf8_3Context is an interface to support dynamic dispatch.
type IUtf8_3Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUtf8_3Context differentiates from other interfaces.
	IsUtf8_3Context()
}

type Utf8_3Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUtf8_3Context() *Utf8_3Context {
	var p = new(Utf8_3Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_utf8_3
	return p
}

func (*Utf8_3Context) IsUtf8_3Context() {}

func NewUtf8_3Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Utf8_3Context {
	var p = new(Utf8_3Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_utf8_3

	return p
}

func (s *Utf8_3Context) GetParser() antlr.Parser { return s.parser }

func (s *Utf8_3Context) U_00E0() antlr.TerminalNode {
	return s.GetToken(CGParserU_00E0, 0)
}

func (s *Utf8_3Context) AllUtf8_tail() []IUtf8_tailContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IUtf8_tailContext)(nil)).Elem())
	var tst = make([]IUtf8_tailContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IUtf8_tailContext)
		}
	}

	return tst
}

func (s *Utf8_3Context) Utf8_tail(i int) IUtf8_tailContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUtf8_tailContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IUtf8_tailContext)
}

func (s *Utf8_3Context) U_00A0() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A0, 0)
}

func (s *Utf8_3Context) U_00A1() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A1, 0)
}

func (s *Utf8_3Context) U_00A2() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A2, 0)
}

func (s *Utf8_3Context) U_00A3() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A3, 0)
}

func (s *Utf8_3Context) U_00A4() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A4, 0)
}

func (s *Utf8_3Context) U_00A5() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A5, 0)
}

func (s *Utf8_3Context) U_00A6() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A6, 0)
}

func (s *Utf8_3Context) U_00A7() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A7, 0)
}

func (s *Utf8_3Context) U_00A8() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A8, 0)
}

func (s *Utf8_3Context) U_00A9() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A9, 0)
}

func (s *Utf8_3Context) U_00AA() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AA, 0)
}

func (s *Utf8_3Context) U_00AB() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AB, 0)
}

func (s *Utf8_3Context) U_00AC() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AC, 0)
}

func (s *Utf8_3Context) U_00AD() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AD, 0)
}

func (s *Utf8_3Context) U_00AE() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AE, 0)
}

func (s *Utf8_3Context) U_00AF() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AF, 0)
}

func (s *Utf8_3Context) U_00B0() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B0, 0)
}

func (s *Utf8_3Context) U_00B1() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B1, 0)
}

func (s *Utf8_3Context) U_00B2() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B2, 0)
}

func (s *Utf8_3Context) U_00B3() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B3, 0)
}

func (s *Utf8_3Context) U_00B4() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B4, 0)
}

func (s *Utf8_3Context) U_00B5() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B5, 0)
}

func (s *Utf8_3Context) U_00B6() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B6, 0)
}

func (s *Utf8_3Context) U_00B7() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B7, 0)
}

func (s *Utf8_3Context) U_00B8() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B8, 0)
}

func (s *Utf8_3Context) U_00B9() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B9, 0)
}

func (s *Utf8_3Context) U_00BA() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BA, 0)
}

func (s *Utf8_3Context) U_00BB() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BB, 0)
}

func (s *Utf8_3Context) U_00BC() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BC, 0)
}

func (s *Utf8_3Context) U_00BD() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BD, 0)
}

func (s *Utf8_3Context) U_00BE() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BE, 0)
}

func (s *Utf8_3Context) U_00BF() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BF, 0)
}

func (s *Utf8_3Context) U_00E1() antlr.TerminalNode {
	return s.GetToken(CGParserU_00E1, 0)
}

func (s *Utf8_3Context) U_00E2() antlr.TerminalNode {
	return s.GetToken(CGParserU_00E2, 0)
}

func (s *Utf8_3Context) U_00E3() antlr.TerminalNode {
	return s.GetToken(CGParserU_00E3, 0)
}

func (s *Utf8_3Context) U_00E4() antlr.TerminalNode {
	return s.GetToken(CGParserU_00E4, 0)
}

func (s *Utf8_3Context) U_00E5() antlr.TerminalNode {
	return s.GetToken(CGParserU_00E5, 0)
}

func (s *Utf8_3Context) U_00E6() antlr.TerminalNode {
	return s.GetToken(CGParserU_00E6, 0)
}

func (s *Utf8_3Context) U_00E7() antlr.TerminalNode {
	return s.GetToken(CGParserU_00E7, 0)
}

func (s *Utf8_3Context) U_00E8() antlr.TerminalNode {
	return s.GetToken(CGParserU_00E8, 0)
}

func (s *Utf8_3Context) U_00E9() antlr.TerminalNode {
	return s.GetToken(CGParserU_00E9, 0)
}

func (s *Utf8_3Context) U_00EA() antlr.TerminalNode {
	return s.GetToken(CGParserU_00EA, 0)
}

func (s *Utf8_3Context) U_00EB() antlr.TerminalNode {
	return s.GetToken(CGParserU_00EB, 0)
}

func (s *Utf8_3Context) U_00EC() antlr.TerminalNode {
	return s.GetToken(CGParserU_00EC, 0)
}

func (s *Utf8_3Context) U_00ED() antlr.TerminalNode {
	return s.GetToken(CGParserU_00ED, 0)
}

func (s *Utf8_3Context) U_0080() antlr.TerminalNode {
	return s.GetToken(CGParserU_0080, 0)
}

func (s *Utf8_3Context) U_0081() antlr.TerminalNode {
	return s.GetToken(CGParserU_0081, 0)
}

func (s *Utf8_3Context) U_0082() antlr.TerminalNode {
	return s.GetToken(CGParserU_0082, 0)
}

func (s *Utf8_3Context) U_0083() antlr.TerminalNode {
	return s.GetToken(CGParserU_0083, 0)
}

func (s *Utf8_3Context) U_0084() antlr.TerminalNode {
	return s.GetToken(CGParserU_0084, 0)
}

func (s *Utf8_3Context) U_0085() antlr.TerminalNode {
	return s.GetToken(CGParserU_0085, 0)
}

func (s *Utf8_3Context) U_0086() antlr.TerminalNode {
	return s.GetToken(CGParserU_0086, 0)
}

func (s *Utf8_3Context) U_0087() antlr.TerminalNode {
	return s.GetToken(CGParserU_0087, 0)
}

func (s *Utf8_3Context) U_0088() antlr.TerminalNode {
	return s.GetToken(CGParserU_0088, 0)
}

func (s *Utf8_3Context) U_0089() antlr.TerminalNode {
	return s.GetToken(CGParserU_0089, 0)
}

func (s *Utf8_3Context) U_008A() antlr.TerminalNode {
	return s.GetToken(CGParserU_008A, 0)
}

func (s *Utf8_3Context) U_008B() antlr.TerminalNode {
	return s.GetToken(CGParserU_008B, 0)
}

func (s *Utf8_3Context) U_008C() antlr.TerminalNode {
	return s.GetToken(CGParserU_008C, 0)
}

func (s *Utf8_3Context) U_008D() antlr.TerminalNode {
	return s.GetToken(CGParserU_008D, 0)
}

func (s *Utf8_3Context) U_008E() antlr.TerminalNode {
	return s.GetToken(CGParserU_008E, 0)
}

func (s *Utf8_3Context) U_008F() antlr.TerminalNode {
	return s.GetToken(CGParserU_008F, 0)
}

func (s *Utf8_3Context) U_0090() antlr.TerminalNode {
	return s.GetToken(CGParserU_0090, 0)
}

func (s *Utf8_3Context) U_0091() antlr.TerminalNode {
	return s.GetToken(CGParserU_0091, 0)
}

func (s *Utf8_3Context) U_0092() antlr.TerminalNode {
	return s.GetToken(CGParserU_0092, 0)
}

func (s *Utf8_3Context) U_0093() antlr.TerminalNode {
	return s.GetToken(CGParserU_0093, 0)
}

func (s *Utf8_3Context) U_0094() antlr.TerminalNode {
	return s.GetToken(CGParserU_0094, 0)
}

func (s *Utf8_3Context) U_0095() antlr.TerminalNode {
	return s.GetToken(CGParserU_0095, 0)
}

func (s *Utf8_3Context) U_0096() antlr.TerminalNode {
	return s.GetToken(CGParserU_0096, 0)
}

func (s *Utf8_3Context) U_0097() antlr.TerminalNode {
	return s.GetToken(CGParserU_0097, 0)
}

func (s *Utf8_3Context) U_0098() antlr.TerminalNode {
	return s.GetToken(CGParserU_0098, 0)
}

func (s *Utf8_3Context) U_0099() antlr.TerminalNode {
	return s.GetToken(CGParserU_0099, 0)
}

func (s *Utf8_3Context) U_009A() antlr.TerminalNode {
	return s.GetToken(CGParserU_009A, 0)
}

func (s *Utf8_3Context) U_009B() antlr.TerminalNode {
	return s.GetToken(CGParserU_009B, 0)
}

func (s *Utf8_3Context) U_009C() antlr.TerminalNode {
	return s.GetToken(CGParserU_009C, 0)
}

func (s *Utf8_3Context) U_009D() antlr.TerminalNode {
	return s.GetToken(CGParserU_009D, 0)
}

func (s *Utf8_3Context) U_009E() antlr.TerminalNode {
	return s.GetToken(CGParserU_009E, 0)
}

func (s *Utf8_3Context) U_009F() antlr.TerminalNode {
	return s.GetToken(CGParserU_009F, 0)
}

func (s *Utf8_3Context) U_00EE() antlr.TerminalNode {
	return s.GetToken(CGParserU_00EE, 0)
}

func (s *Utf8_3Context) U_00EF() antlr.TerminalNode {
	return s.GetToken(CGParserU_00EF, 0)
}

func (s *Utf8_3Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Utf8_3Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Utf8_3Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterUtf8_3(s)
	}
}

func (s *Utf8_3Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitUtf8_3(s)
	}
}

func (p *CGParser) Utf8_3() (localctx IUtf8_3Context) {
	localctx = NewUtf8_3Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, CGParserRULE_utf8_3)
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

	p.SetState(399)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CGParserU_00E0:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(385)
			p.Match(CGParserU_00E0)
		}
		{
			p.SetState(386)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-131)&-(0x1f+1)) == 0 && ((1<<uint((_la-131)))&((1<<(CGParserU_00A0-131))|(1<<(CGParserU_00A1-131))|(1<<(CGParserU_00A2-131))|(1<<(CGParserU_00A3-131))|(1<<(CGParserU_00A4-131))|(1<<(CGParserU_00A5-131))|(1<<(CGParserU_00A6-131))|(1<<(CGParserU_00A7-131))|(1<<(CGParserU_00A8-131))|(1<<(CGParserU_00A9-131))|(1<<(CGParserU_00AA-131))|(1<<(CGParserU_00AB-131))|(1<<(CGParserU_00AC-131))|(1<<(CGParserU_00AD-131))|(1<<(CGParserU_00AE-131))|(1<<(CGParserU_00AF-131))|(1<<(CGParserU_00B0-131))|(1<<(CGParserU_00B1-131))|(1<<(CGParserU_00B2-131))|(1<<(CGParserU_00B3-131))|(1<<(CGParserU_00B4-131))|(1<<(CGParserU_00B5-131))|(1<<(CGParserU_00B6-131))|(1<<(CGParserU_00B7-131))|(1<<(CGParserU_00B8-131))|(1<<(CGParserU_00B9-131))|(1<<(CGParserU_00BA-131))|(1<<(CGParserU_00BB-131))|(1<<(CGParserU_00BC-131))|(1<<(CGParserU_00BD-131))|(1<<(CGParserU_00BE-131))|(1<<(CGParserU_00BF-131)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(387)
			p.Utf8_tail()
		}

	case CGParserU_00E1, CGParserU_00E2, CGParserU_00E3, CGParserU_00E4, CGParserU_00E5, CGParserU_00E6, CGParserU_00E7, CGParserU_00E8, CGParserU_00E9, CGParserU_00EA, CGParserU_00EB, CGParserU_00EC:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(388)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-194)&-(0x1f+1)) == 0 && ((1<<uint((_la-194)))&((1<<(CGParserU_00E1-194))|(1<<(CGParserU_00E2-194))|(1<<(CGParserU_00E3-194))|(1<<(CGParserU_00E4-194))|(1<<(CGParserU_00E5-194))|(1<<(CGParserU_00E6-194))|(1<<(CGParserU_00E7-194))|(1<<(CGParserU_00E8-194))|(1<<(CGParserU_00E9-194))|(1<<(CGParserU_00EA-194))|(1<<(CGParserU_00EB-194))|(1<<(CGParserU_00EC-194)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		{
			p.SetState(389)
			p.Utf8_tail()
		}

		{
			p.SetState(390)
			p.Utf8_tail()
		}

	case CGParserU_00ED:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(392)
			p.Match(CGParserU_00ED)
		}
		{
			p.SetState(393)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-99)&-(0x1f+1)) == 0 && ((1<<uint((_la-99)))&((1<<(CGParserU_0080-99))|(1<<(CGParserU_0081-99))|(1<<(CGParserU_0082-99))|(1<<(CGParserU_0083-99))|(1<<(CGParserU_0084-99))|(1<<(CGParserU_0085-99))|(1<<(CGParserU_0086-99))|(1<<(CGParserU_0087-99))|(1<<(CGParserU_0088-99))|(1<<(CGParserU_0089-99))|(1<<(CGParserU_008A-99))|(1<<(CGParserU_008B-99))|(1<<(CGParserU_008C-99))|(1<<(CGParserU_008D-99))|(1<<(CGParserU_008E-99))|(1<<(CGParserU_008F-99))|(1<<(CGParserU_0090-99))|(1<<(CGParserU_0091-99))|(1<<(CGParserU_0092-99))|(1<<(CGParserU_0093-99))|(1<<(CGParserU_0094-99))|(1<<(CGParserU_0095-99))|(1<<(CGParserU_0096-99))|(1<<(CGParserU_0097-99))|(1<<(CGParserU_0098-99))|(1<<(CGParserU_0099-99))|(1<<(CGParserU_009A-99))|(1<<(CGParserU_009B-99))|(1<<(CGParserU_009C-99))|(1<<(CGParserU_009D-99))|(1<<(CGParserU_009E-99))|(1<<(CGParserU_009F-99)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(394)
			p.Utf8_tail()
		}

	case CGParserU_00EE, CGParserU_00EF:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(395)
			_la = p.GetTokenStream().LA(1)

			if !(_la == CGParserU_00EE || _la == CGParserU_00EF) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		{
			p.SetState(396)
			p.Utf8_tail()
		}

		{
			p.SetState(397)
			p.Utf8_tail()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IUtf8_4Context is an interface to support dynamic dispatch.
type IUtf8_4Context interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUtf8_4Context differentiates from other interfaces.
	IsUtf8_4Context()
}

type Utf8_4Context struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUtf8_4Context() *Utf8_4Context {
	var p = new(Utf8_4Context)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_utf8_4
	return p
}

func (*Utf8_4Context) IsUtf8_4Context() {}

func NewUtf8_4Context(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Utf8_4Context {
	var p = new(Utf8_4Context)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_utf8_4

	return p
}

func (s *Utf8_4Context) GetParser() antlr.Parser { return s.parser }

func (s *Utf8_4Context) U_00F0() antlr.TerminalNode {
	return s.GetToken(CGParserU_00F0, 0)
}

func (s *Utf8_4Context) U_0090() antlr.TerminalNode {
	return s.GetToken(CGParserU_0090, 0)
}

func (s *Utf8_4Context) U_0091() antlr.TerminalNode {
	return s.GetToken(CGParserU_0091, 0)
}

func (s *Utf8_4Context) U_0092() antlr.TerminalNode {
	return s.GetToken(CGParserU_0092, 0)
}

func (s *Utf8_4Context) U_0093() antlr.TerminalNode {
	return s.GetToken(CGParserU_0093, 0)
}

func (s *Utf8_4Context) U_0094() antlr.TerminalNode {
	return s.GetToken(CGParserU_0094, 0)
}

func (s *Utf8_4Context) U_0095() antlr.TerminalNode {
	return s.GetToken(CGParserU_0095, 0)
}

func (s *Utf8_4Context) U_0096() antlr.TerminalNode {
	return s.GetToken(CGParserU_0096, 0)
}

func (s *Utf8_4Context) U_0097() antlr.TerminalNode {
	return s.GetToken(CGParserU_0097, 0)
}

func (s *Utf8_4Context) U_0098() antlr.TerminalNode {
	return s.GetToken(CGParserU_0098, 0)
}

func (s *Utf8_4Context) U_0099() antlr.TerminalNode {
	return s.GetToken(CGParserU_0099, 0)
}

func (s *Utf8_4Context) U_009A() antlr.TerminalNode {
	return s.GetToken(CGParserU_009A, 0)
}

func (s *Utf8_4Context) U_009B() antlr.TerminalNode {
	return s.GetToken(CGParserU_009B, 0)
}

func (s *Utf8_4Context) U_009C() antlr.TerminalNode {
	return s.GetToken(CGParserU_009C, 0)
}

func (s *Utf8_4Context) U_009D() antlr.TerminalNode {
	return s.GetToken(CGParserU_009D, 0)
}

func (s *Utf8_4Context) U_009E() antlr.TerminalNode {
	return s.GetToken(CGParserU_009E, 0)
}

func (s *Utf8_4Context) U_009F() antlr.TerminalNode {
	return s.GetToken(CGParserU_009F, 0)
}

func (s *Utf8_4Context) U_00A0() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A0, 0)
}

func (s *Utf8_4Context) U_00A1() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A1, 0)
}

func (s *Utf8_4Context) U_00A2() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A2, 0)
}

func (s *Utf8_4Context) U_00A3() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A3, 0)
}

func (s *Utf8_4Context) U_00A4() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A4, 0)
}

func (s *Utf8_4Context) U_00A5() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A5, 0)
}

func (s *Utf8_4Context) U_00A6() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A6, 0)
}

func (s *Utf8_4Context) U_00A7() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A7, 0)
}

func (s *Utf8_4Context) U_00A8() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A8, 0)
}

func (s *Utf8_4Context) U_00A9() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A9, 0)
}

func (s *Utf8_4Context) U_00AA() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AA, 0)
}

func (s *Utf8_4Context) U_00AB() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AB, 0)
}

func (s *Utf8_4Context) U_00AC() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AC, 0)
}

func (s *Utf8_4Context) U_00AD() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AD, 0)
}

func (s *Utf8_4Context) U_00AE() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AE, 0)
}

func (s *Utf8_4Context) U_00AF() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AF, 0)
}

func (s *Utf8_4Context) U_00B0() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B0, 0)
}

func (s *Utf8_4Context) U_00B1() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B1, 0)
}

func (s *Utf8_4Context) U_00B2() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B2, 0)
}

func (s *Utf8_4Context) U_00B3() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B3, 0)
}

func (s *Utf8_4Context) U_00B4() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B4, 0)
}

func (s *Utf8_4Context) U_00B5() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B5, 0)
}

func (s *Utf8_4Context) U_00B6() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B6, 0)
}

func (s *Utf8_4Context) U_00B7() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B7, 0)
}

func (s *Utf8_4Context) U_00B8() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B8, 0)
}

func (s *Utf8_4Context) U_00B9() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B9, 0)
}

func (s *Utf8_4Context) U_00BA() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BA, 0)
}

func (s *Utf8_4Context) U_00BB() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BB, 0)
}

func (s *Utf8_4Context) U_00BC() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BC, 0)
}

func (s *Utf8_4Context) U_00BD() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BD, 0)
}

func (s *Utf8_4Context) U_00BE() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BE, 0)
}

func (s *Utf8_4Context) U_00BF() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BF, 0)
}

func (s *Utf8_4Context) AllUtf8_tail() []IUtf8_tailContext {
	var ts = s.GetTypedRuleContexts(reflect.TypeOf((*IUtf8_tailContext)(nil)).Elem())
	var tst = make([]IUtf8_tailContext, len(ts))

	for i, t := range ts {
		if t != nil {
			tst[i] = t.(IUtf8_tailContext)
		}
	}

	return tst
}

func (s *Utf8_4Context) Utf8_tail(i int) IUtf8_tailContext {
	var t = s.GetTypedRuleContext(reflect.TypeOf((*IUtf8_tailContext)(nil)).Elem(), i)

	if t == nil {
		return nil
	}

	return t.(IUtf8_tailContext)
}

func (s *Utf8_4Context) U_00F1() antlr.TerminalNode {
	return s.GetToken(CGParserU_00F1, 0)
}

func (s *Utf8_4Context) U_00F2() antlr.TerminalNode {
	return s.GetToken(CGParserU_00F2, 0)
}

func (s *Utf8_4Context) U_00F3() antlr.TerminalNode {
	return s.GetToken(CGParserU_00F3, 0)
}

func (s *Utf8_4Context) U_00F4() antlr.TerminalNode {
	return s.GetToken(CGParserU_00F4, 0)
}

func (s *Utf8_4Context) U_0080() antlr.TerminalNode {
	return s.GetToken(CGParserU_0080, 0)
}

func (s *Utf8_4Context) U_0081() antlr.TerminalNode {
	return s.GetToken(CGParserU_0081, 0)
}

func (s *Utf8_4Context) U_0082() antlr.TerminalNode {
	return s.GetToken(CGParserU_0082, 0)
}

func (s *Utf8_4Context) U_0083() antlr.TerminalNode {
	return s.GetToken(CGParserU_0083, 0)
}

func (s *Utf8_4Context) U_0084() antlr.TerminalNode {
	return s.GetToken(CGParserU_0084, 0)
}

func (s *Utf8_4Context) U_0085() antlr.TerminalNode {
	return s.GetToken(CGParserU_0085, 0)
}

func (s *Utf8_4Context) U_0086() antlr.TerminalNode {
	return s.GetToken(CGParserU_0086, 0)
}

func (s *Utf8_4Context) U_0087() antlr.TerminalNode {
	return s.GetToken(CGParserU_0087, 0)
}

func (s *Utf8_4Context) U_0088() antlr.TerminalNode {
	return s.GetToken(CGParserU_0088, 0)
}

func (s *Utf8_4Context) U_0089() antlr.TerminalNode {
	return s.GetToken(CGParserU_0089, 0)
}

func (s *Utf8_4Context) U_008A() antlr.TerminalNode {
	return s.GetToken(CGParserU_008A, 0)
}

func (s *Utf8_4Context) U_008B() antlr.TerminalNode {
	return s.GetToken(CGParserU_008B, 0)
}

func (s *Utf8_4Context) U_008C() antlr.TerminalNode {
	return s.GetToken(CGParserU_008C, 0)
}

func (s *Utf8_4Context) U_008D() antlr.TerminalNode {
	return s.GetToken(CGParserU_008D, 0)
}

func (s *Utf8_4Context) U_008E() antlr.TerminalNode {
	return s.GetToken(CGParserU_008E, 0)
}

func (s *Utf8_4Context) U_008F() antlr.TerminalNode {
	return s.GetToken(CGParserU_008F, 0)
}

func (s *Utf8_4Context) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Utf8_4Context) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Utf8_4Context) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterUtf8_4(s)
	}
}

func (s *Utf8_4Context) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitUtf8_4(s)
	}
}

func (p *CGParser) Utf8_4() (localctx IUtf8_4Context) {
	localctx = NewUtf8_4Context(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, CGParserRULE_utf8_4)
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

	p.SetState(416)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case CGParserU_00F0:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(401)
			p.Match(CGParserU_00F0)
		}
		{
			p.SetState(402)
			_la = p.GetTokenStream().LA(1)

			if !((((_la-115)&-(0x1f+1)) == 0 && ((1<<uint((_la-115)))&((1<<(CGParserU_0090-115))|(1<<(CGParserU_0091-115))|(1<<(CGParserU_0092-115))|(1<<(CGParserU_0093-115))|(1<<(CGParserU_0094-115))|(1<<(CGParserU_0095-115))|(1<<(CGParserU_0096-115))|(1<<(CGParserU_0097-115))|(1<<(CGParserU_0098-115))|(1<<(CGParserU_0099-115))|(1<<(CGParserU_009A-115))|(1<<(CGParserU_009B-115))|(1<<(CGParserU_009C-115))|(1<<(CGParserU_009D-115))|(1<<(CGParserU_009E-115))|(1<<(CGParserU_009F-115))|(1<<(CGParserU_00A0-115))|(1<<(CGParserU_00A1-115))|(1<<(CGParserU_00A2-115))|(1<<(CGParserU_00A3-115))|(1<<(CGParserU_00A4-115))|(1<<(CGParserU_00A5-115))|(1<<(CGParserU_00A6-115))|(1<<(CGParserU_00A7-115))|(1<<(CGParserU_00A8-115))|(1<<(CGParserU_00A9-115))|(1<<(CGParserU_00AA-115))|(1<<(CGParserU_00AB-115))|(1<<(CGParserU_00AC-115))|(1<<(CGParserU_00AD-115))|(1<<(CGParserU_00AE-115))|(1<<(CGParserU_00AF-115)))) != 0) || (((_la-147)&-(0x1f+1)) == 0 && ((1<<uint((_la-147)))&((1<<(CGParserU_00B0-147))|(1<<(CGParserU_00B1-147))|(1<<(CGParserU_00B2-147))|(1<<(CGParserU_00B3-147))|(1<<(CGParserU_00B4-147))|(1<<(CGParserU_00B5-147))|(1<<(CGParserU_00B6-147))|(1<<(CGParserU_00B7-147))|(1<<(CGParserU_00B8-147))|(1<<(CGParserU_00B9-147))|(1<<(CGParserU_00BA-147))|(1<<(CGParserU_00BB-147))|(1<<(CGParserU_00BC-147))|(1<<(CGParserU_00BD-147))|(1<<(CGParserU_00BE-147))|(1<<(CGParserU_00BF-147)))) != 0)) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		{
			p.SetState(403)
			p.Utf8_tail()
		}

		{
			p.SetState(404)
			p.Utf8_tail()
		}

	case CGParserU_00F1, CGParserU_00F2, CGParserU_00F3:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(406)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-210)&-(0x1f+1)) == 0 && ((1<<uint((_la-210)))&((1<<(CGParserU_00F1-210))|(1<<(CGParserU_00F2-210))|(1<<(CGParserU_00F3-210)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		{
			p.SetState(407)
			p.Utf8_tail()
		}

		{
			p.SetState(408)
			p.Utf8_tail()
		}

		{
			p.SetState(409)
			p.Utf8_tail()
		}

	case CGParserU_00F4:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(411)
			p.Match(CGParserU_00F4)
		}
		{
			p.SetState(412)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-99)&-(0x1f+1)) == 0 && ((1<<uint((_la-99)))&((1<<(CGParserU_0080-99))|(1<<(CGParserU_0081-99))|(1<<(CGParserU_0082-99))|(1<<(CGParserU_0083-99))|(1<<(CGParserU_0084-99))|(1<<(CGParserU_0085-99))|(1<<(CGParserU_0086-99))|(1<<(CGParserU_0087-99))|(1<<(CGParserU_0088-99))|(1<<(CGParserU_0089-99))|(1<<(CGParserU_008A-99))|(1<<(CGParserU_008B-99))|(1<<(CGParserU_008C-99))|(1<<(CGParserU_008D-99))|(1<<(CGParserU_008E-99))|(1<<(CGParserU_008F-99)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

		{
			p.SetState(413)
			p.Utf8_tail()
		}

		{
			p.SetState(414)
			p.Utf8_tail()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IUtf8_tailContext is an interface to support dynamic dispatch.
type IUtf8_tailContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsUtf8_tailContext differentiates from other interfaces.
	IsUtf8_tailContext()
}

type Utf8_tailContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUtf8_tailContext() *Utf8_tailContext {
	var p = new(Utf8_tailContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = CGParserRULE_utf8_tail
	return p
}

func (*Utf8_tailContext) IsUtf8_tailContext() {}

func NewUtf8_tailContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Utf8_tailContext {
	var p = new(Utf8_tailContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = CGParserRULE_utf8_tail

	return p
}

func (s *Utf8_tailContext) GetParser() antlr.Parser { return s.parser }

func (s *Utf8_tailContext) U_0080() antlr.TerminalNode {
	return s.GetToken(CGParserU_0080, 0)
}

func (s *Utf8_tailContext) U_0081() antlr.TerminalNode {
	return s.GetToken(CGParserU_0081, 0)
}

func (s *Utf8_tailContext) U_0082() antlr.TerminalNode {
	return s.GetToken(CGParserU_0082, 0)
}

func (s *Utf8_tailContext) U_0083() antlr.TerminalNode {
	return s.GetToken(CGParserU_0083, 0)
}

func (s *Utf8_tailContext) U_0084() antlr.TerminalNode {
	return s.GetToken(CGParserU_0084, 0)
}

func (s *Utf8_tailContext) U_0085() antlr.TerminalNode {
	return s.GetToken(CGParserU_0085, 0)
}

func (s *Utf8_tailContext) U_0086() antlr.TerminalNode {
	return s.GetToken(CGParserU_0086, 0)
}

func (s *Utf8_tailContext) U_0087() antlr.TerminalNode {
	return s.GetToken(CGParserU_0087, 0)
}

func (s *Utf8_tailContext) U_0088() antlr.TerminalNode {
	return s.GetToken(CGParserU_0088, 0)
}

func (s *Utf8_tailContext) U_0089() antlr.TerminalNode {
	return s.GetToken(CGParserU_0089, 0)
}

func (s *Utf8_tailContext) U_008A() antlr.TerminalNode {
	return s.GetToken(CGParserU_008A, 0)
}

func (s *Utf8_tailContext) U_008B() antlr.TerminalNode {
	return s.GetToken(CGParserU_008B, 0)
}

func (s *Utf8_tailContext) U_008C() antlr.TerminalNode {
	return s.GetToken(CGParserU_008C, 0)
}

func (s *Utf8_tailContext) U_008D() antlr.TerminalNode {
	return s.GetToken(CGParserU_008D, 0)
}

func (s *Utf8_tailContext) U_008E() antlr.TerminalNode {
	return s.GetToken(CGParserU_008E, 0)
}

func (s *Utf8_tailContext) U_008F() antlr.TerminalNode {
	return s.GetToken(CGParserU_008F, 0)
}

func (s *Utf8_tailContext) U_0090() antlr.TerminalNode {
	return s.GetToken(CGParserU_0090, 0)
}

func (s *Utf8_tailContext) U_0091() antlr.TerminalNode {
	return s.GetToken(CGParserU_0091, 0)
}

func (s *Utf8_tailContext) U_0092() antlr.TerminalNode {
	return s.GetToken(CGParserU_0092, 0)
}

func (s *Utf8_tailContext) U_0093() antlr.TerminalNode {
	return s.GetToken(CGParserU_0093, 0)
}

func (s *Utf8_tailContext) U_0094() antlr.TerminalNode {
	return s.GetToken(CGParserU_0094, 0)
}

func (s *Utf8_tailContext) U_0095() antlr.TerminalNode {
	return s.GetToken(CGParserU_0095, 0)
}

func (s *Utf8_tailContext) U_0096() antlr.TerminalNode {
	return s.GetToken(CGParserU_0096, 0)
}

func (s *Utf8_tailContext) U_0097() antlr.TerminalNode {
	return s.GetToken(CGParserU_0097, 0)
}

func (s *Utf8_tailContext) U_0098() antlr.TerminalNode {
	return s.GetToken(CGParserU_0098, 0)
}

func (s *Utf8_tailContext) U_0099() antlr.TerminalNode {
	return s.GetToken(CGParserU_0099, 0)
}

func (s *Utf8_tailContext) U_009A() antlr.TerminalNode {
	return s.GetToken(CGParserU_009A, 0)
}

func (s *Utf8_tailContext) U_009B() antlr.TerminalNode {
	return s.GetToken(CGParserU_009B, 0)
}

func (s *Utf8_tailContext) U_009C() antlr.TerminalNode {
	return s.GetToken(CGParserU_009C, 0)
}

func (s *Utf8_tailContext) U_009D() antlr.TerminalNode {
	return s.GetToken(CGParserU_009D, 0)
}

func (s *Utf8_tailContext) U_009E() antlr.TerminalNode {
	return s.GetToken(CGParserU_009E, 0)
}

func (s *Utf8_tailContext) U_009F() antlr.TerminalNode {
	return s.GetToken(CGParserU_009F, 0)
}

func (s *Utf8_tailContext) U_00A0() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A0, 0)
}

func (s *Utf8_tailContext) U_00A1() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A1, 0)
}

func (s *Utf8_tailContext) U_00A2() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A2, 0)
}

func (s *Utf8_tailContext) U_00A3() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A3, 0)
}

func (s *Utf8_tailContext) U_00A4() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A4, 0)
}

func (s *Utf8_tailContext) U_00A5() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A5, 0)
}

func (s *Utf8_tailContext) U_00A6() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A6, 0)
}

func (s *Utf8_tailContext) U_00A7() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A7, 0)
}

func (s *Utf8_tailContext) U_00A8() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A8, 0)
}

func (s *Utf8_tailContext) U_00A9() antlr.TerminalNode {
	return s.GetToken(CGParserU_00A9, 0)
}

func (s *Utf8_tailContext) U_00AA() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AA, 0)
}

func (s *Utf8_tailContext) U_00AB() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AB, 0)
}

func (s *Utf8_tailContext) U_00AC() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AC, 0)
}

func (s *Utf8_tailContext) U_00AD() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AD, 0)
}

func (s *Utf8_tailContext) U_00AE() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AE, 0)
}

func (s *Utf8_tailContext) U_00AF() antlr.TerminalNode {
	return s.GetToken(CGParserU_00AF, 0)
}

func (s *Utf8_tailContext) U_00B0() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B0, 0)
}

func (s *Utf8_tailContext) U_00B1() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B1, 0)
}

func (s *Utf8_tailContext) U_00B2() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B2, 0)
}

func (s *Utf8_tailContext) U_00B3() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B3, 0)
}

func (s *Utf8_tailContext) U_00B4() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B4, 0)
}

func (s *Utf8_tailContext) U_00B5() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B5, 0)
}

func (s *Utf8_tailContext) U_00B6() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B6, 0)
}

func (s *Utf8_tailContext) U_00B7() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B7, 0)
}

func (s *Utf8_tailContext) U_00B8() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B8, 0)
}

func (s *Utf8_tailContext) U_00B9() antlr.TerminalNode {
	return s.GetToken(CGParserU_00B9, 0)
}

func (s *Utf8_tailContext) U_00BA() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BA, 0)
}

func (s *Utf8_tailContext) U_00BB() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BB, 0)
}

func (s *Utf8_tailContext) U_00BC() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BC, 0)
}

func (s *Utf8_tailContext) U_00BD() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BD, 0)
}

func (s *Utf8_tailContext) U_00BE() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BE, 0)
}

func (s *Utf8_tailContext) U_00BF() antlr.TerminalNode {
	return s.GetToken(CGParserU_00BF, 0)
}

func (s *Utf8_tailContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Utf8_tailContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Utf8_tailContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.EnterUtf8_tail(s)
	}
}

func (s *Utf8_tailContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(CGListener); ok {
		listenerT.ExitUtf8_tail(s)
	}
}

func (p *CGParser) Utf8_tail() (localctx IUtf8_tailContext) {
	localctx = NewUtf8_tailContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, CGParserRULE_utf8_tail)
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
		p.SetState(418)
		_la = p.GetTokenStream().LA(1)

		if !((((_la-99)&-(0x1f+1)) == 0 && ((1<<uint((_la-99)))&((1<<(CGParserU_0080-99))|(1<<(CGParserU_0081-99))|(1<<(CGParserU_0082-99))|(1<<(CGParserU_0083-99))|(1<<(CGParserU_0084-99))|(1<<(CGParserU_0085-99))|(1<<(CGParserU_0086-99))|(1<<(CGParserU_0087-99))|(1<<(CGParserU_0088-99))|(1<<(CGParserU_0089-99))|(1<<(CGParserU_008A-99))|(1<<(CGParserU_008B-99))|(1<<(CGParserU_008C-99))|(1<<(CGParserU_008D-99))|(1<<(CGParserU_008E-99))|(1<<(CGParserU_008F-99))|(1<<(CGParserU_0090-99))|(1<<(CGParserU_0091-99))|(1<<(CGParserU_0092-99))|(1<<(CGParserU_0093-99))|(1<<(CGParserU_0094-99))|(1<<(CGParserU_0095-99))|(1<<(CGParserU_0096-99))|(1<<(CGParserU_0097-99))|(1<<(CGParserU_0098-99))|(1<<(CGParserU_0099-99))|(1<<(CGParserU_009A-99))|(1<<(CGParserU_009B-99))|(1<<(CGParserU_009C-99))|(1<<(CGParserU_009D-99))|(1<<(CGParserU_009E-99))|(1<<(CGParserU_009F-99)))) != 0) || (((_la-131)&-(0x1f+1)) == 0 && ((1<<uint((_la-131)))&((1<<(CGParserU_00A0-131))|(1<<(CGParserU_00A1-131))|(1<<(CGParserU_00A2-131))|(1<<(CGParserU_00A3-131))|(1<<(CGParserU_00A4-131))|(1<<(CGParserU_00A5-131))|(1<<(CGParserU_00A6-131))|(1<<(CGParserU_00A7-131))|(1<<(CGParserU_00A8-131))|(1<<(CGParserU_00A9-131))|(1<<(CGParserU_00AA-131))|(1<<(CGParserU_00AB-131))|(1<<(CGParserU_00AC-131))|(1<<(CGParserU_00AD-131))|(1<<(CGParserU_00AE-131))|(1<<(CGParserU_00AF-131))|(1<<(CGParserU_00B0-131))|(1<<(CGParserU_00B1-131))|(1<<(CGParserU_00B2-131))|(1<<(CGParserU_00B3-131))|(1<<(CGParserU_00B4-131))|(1<<(CGParserU_00B5-131))|(1<<(CGParserU_00B6-131))|(1<<(CGParserU_00B7-131))|(1<<(CGParserU_00B8-131))|(1<<(CGParserU_00B9-131))|(1<<(CGParserU_00BA-131))|(1<<(CGParserU_00BB-131))|(1<<(CGParserU_00BC-131))|(1<<(CGParserU_00BD-131))|(1<<(CGParserU_00BE-131))|(1<<(CGParserU_00BF-131)))) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}
