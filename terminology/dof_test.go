package terminology

import (
	"testing"
)

var testData = []int64{
	22298006,  // myocardial infarction
	230690007, // stroke
	73211009,  // diabetes
	73211009,
	73211009,
	230690007,
	335621000000101, // maternally inherited diabetes mellitus // should go to generic diabetes mellitus
	230690007,
	230690007,
	22298006,
	22298006,
	22298006,
}

var expected3Data = []int64{
	22298006,
	230690007,
	73211009,
	73211009,
	73211009,
	230690007,
	73211009,
	230690007,
	230690007,
	22298006,
	22298006,
	22298006,
}

func TestReduce(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()

	r1 := NewReducer(svc, 3, 1)
	r1r, err := r1.Reduce(testData)
	if err != nil {
		t.Fatal(err)
	}
	for i, id := range r1r {
		if id != expected3Data[i] {
			t.Errorf("expected: %v, got: %v", expected3Data, r1r)
		}
	}

	r2 := NewReducer(svc, 1, 0)
	r2r, err := r2.Reduce([]int64{73211009, 335621000000101}) // diabetes and maternally inherited diabetes
	if err != nil {
		t.Fatal(err)
	}
	if r2r[0] != 73211009 && r2r[1] != 73211009 {
		t.Errorf("expected: 73211009+73211009, got %v", r2r)
	}

	r3 := NewReducer(svc, 0, 0)
	r3r, err := r3.Reduce(testData)
	if r3r[0] != 64572001 {
		t.Errorf("expected test data to be mapped to Disease (64572001). got %v", r3r)
	}
}
