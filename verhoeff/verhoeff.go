// Package verhoeff provides an implementation of the Verhoeff check digit algorithm.
//
// Copyright 2018 Mark Wardle / Eldrix Ltd
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//
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

// Calculate generates a Verhoeff check digit
func Calculate(num string) int {
	c := 0
	ll := len(num)
	for i := 0; i < ll; i++ {
		c = multiplication[c][permutation[((i + 1) % 8)][num[ll-i-1]-'0']]
	}
	return inverse[c]
}

// AppendCheckDigit generates a new string appending a Verhoeff check digit
func AppendCheckDigit(s string) string {
	r := Calculate(s)
	return s + digits[r%10] //strconv.Itoa(r%10)
}

// ValidateString validates that the number is Verhoeff compliant with the last digit the correct check digit.
func ValidateString(num string) bool {
	c := 0
	ll := len(num)
	for i := 0; i < ll; i++ {
		c = multiplication[c][permutation[(i % 8)][num[ll-i-1]-'0']]
	}
	return c == 0
}

// Validate validates that the number is Verhoeff compliant with the last digit the correct check digit.
func Validate(i int64) bool {
	return ValidateString(strconv.FormatInt(i, 10))
}
