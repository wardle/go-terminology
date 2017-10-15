// Package verhoeff provides an implementation of the Verhoeff check digit algorithm.
package verhoeff

import (
	"strconv"
)

// multiplication table
var multiplication = [][]int{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	{1, 2, 3, 4, 0, 6, 7, 8, 9, 5},
	{2, 3, 4, 0, 1, 7, 8, 9, 5, 6},
	{3, 4, 0, 1, 2, 8, 9, 5, 6, 7},
	{4, 0, 1, 2, 3, 9, 5, 6, 7, 8},
	{5, 9, 8, 7, 6, 0, 4, 3, 2, 1},
	{6, 5, 9, 8, 7, 1, 0, 4, 3, 2},
	{7, 6, 5, 9, 8, 2, 1, 0, 4, 3},
	{8, 7, 6, 5, 9, 3, 2, 1, 0, 4},
	{9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
}

// permutation table
var permutation = [][]int{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	{1, 5, 7, 6, 2, 8, 3, 0, 9, 4},
	{5, 8, 0, 3, 7, 9, 6, 1, 4, 2},
	{8, 9, 1, 6, 0, 4, 3, 5, 2, 7},
	{9, 4, 5, 3, 1, 2, 6, 8, 7, 0},
	{4, 2, 8, 6, 5, 7, 3, 9, 0, 1},
	{2, 7, 9, 3, 8, 0, 6, 4, 1, 5},
	{7, 0, 4, 6, 9, 1, 3, 2, 5, 8},
}

var digits = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

// inverse table
var inverse = []int{0, 4, 3, 2, 1, 5, 6, 7, 8, 9}

// CalculateVerhoeff generates a Verhoeff check digit
func CalculateVerhoeff(num string) int {
	c := 0
	ll := len(num)
	for i := 0; i < ll; i++ {
		c = multiplication[c][permutation[((i + 1) % 8)][num[ll-i-1]-'0']]
	}
	return inverse[c]
}

// AppendVerhoeff generates a new string appending a Verhoeff check digit
func AppendVerhoeff(s string) string {
	r := CalculateVerhoeff(s)
	return s + digits[r%10] //strconv.Itoa(r%10)
}

// ValidateVerhoeffString validates that the number is Verhoeff compliant with the last digit the correct check digit.
func ValidateVerhoeffString(num string) bool {
	c := 0
	ll := len(num)
	for i := 0; i < ll; i++ {
		c = multiplication[c][permutation[(i % 8)][num[ll-i-1]-'0']]
	}
	return (c == 0)
}

// ValidateVerhoeff validates that the number is Verhoeff compliant with the last digit the correct check digit.
func ValidateVerhoeff(num int) bool {
	return ValidateVerhoeffString(strconv.Itoa(num))
}
