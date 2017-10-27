package verhoeff

import (
	"strconv"
	"testing"
)

func TestVerhoeff(t *testing.T) {
	test(t, "1234567890", true)
	test(t, "24700007", true)
	test(t, "1334567890", false)
	test(t, "1234567892", false)
	test(t, "14567894", true)
	test(t, "14567895", false)

}

func test(t *testing.T, input string, expected bool) {
	n, err := strconv.Atoi(input)
	if err != nil {
		t.Error(err)
	}
	if ValidateVerhoeff(n) != expected || ValidateVerhoeffString(input) != expected {
		t.Errorf("Failed to validate Verhoeff check digit for: %s", input)
	}
}
