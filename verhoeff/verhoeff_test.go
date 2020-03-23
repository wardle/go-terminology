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
	"testing"
)

func TestVerhoeff(t *testing.T) {
	test(t, "311220190006", true)
	test(t, "1234567890", true)
	test(t, "24700007", true)
	test(t, "1334567890", false)
	test(t, "1234567892", false)
	test(t, "14567894", true)
	test(t, "14567895", false)

}

func test(t *testing.T, input string, expected bool) {
	n, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		t.Error(err)
	}
	if Validate(n) != expected || ValidateString(input) != expected {
		t.Errorf("Failed to validate Verhoeff check digit for: %s", input)
	}
}
