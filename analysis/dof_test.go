package analysis

import (
	"github.com/wardle/go-terminology/terminology"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	dbFilename = "../snomed.db" // real, live database
)

func setUp(tb testing.TB) *terminology.Svc {
	if _, err := os.Stat(dbFilename); os.IsNotExist(err) { // skip these tests if no working live snomed db
		tb.Skipf("Skipping tests against a live database. To run, create a database named %s", dbFilename)
	}
	svc, err := terminology.NewService(dbFilename, true)
	if err != nil {
		tb.Fatal(err)
	}
	return svc
}

// TestReduce tests dimensionality reduction.
//
// test data contains:
// 22298006 myocardial infarction
// 230690007 stroke
// 73211009 Diabetes mellitus
// 335621000000101 - maternally inherited diabetes mellitus // should go to generic diabetes mellitus
func TestReduce(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	r := NewReducer(svc, 3, 0)
	output := &strings.Builder{}
	f, err := os.Open("testdata.txt")
	if err != nil {
		t.Fatal(err)
	}
	if err := r.ReduceCsv(f, output); err != nil {
		t.Fatal(err)
	}
	expected, err := ioutil.ReadFile("testdata-expected.txt")
	if err != nil {
		t.Fatal(err)
	}
	if output.String() != string(expected) {
		t.Fatalf("Result not as expected. Got: %s", output.String())
	}
}
