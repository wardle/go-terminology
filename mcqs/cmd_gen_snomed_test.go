package mcqs

import "testing"
import "reflect"

func TestListAtoi(t *testing.T) {
	testAtoi(t, "123,456,789,123456789", []int{123, 456, 789, 123456789})
	testAtoi(t, "123,456,789, 123456789", []int{123, 456, 789, 123456789})
	testAtoi(t, "aaa,vbb,cc,123", []int{123})
}

func testAtoi(t *testing.T, input string, expected []int) {
	r := listAtoi(input)
	if reflect.DeepEqual(r, expected) == false {
		t.Errorf("Failed to parse: %s. Parsed to %v", input, r)
	}
}
