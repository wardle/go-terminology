package terminology_test

import (
	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
	"testing"
)

func TestSimpleMatch(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	preferred := "en-gb"
	wanted, _, err := language.ParseAcceptLanguage(preferred)
	if err != nil {
		t.Fatal(err)
	}
	best := svc.Match(wanted)
	if best != terminology.BritishEnglish {
		t.Fatalf("Didn't correctly match British English. Matched %v", best)
	}
}
