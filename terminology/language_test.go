package terminology_test

import (
	"testing"

	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
)

func TestSimpleMatch(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	preferred := "en-gb"
	wanted, _, err := language.ParseAcceptLanguage(preferred)
	if err != nil {
		t.Fatal(err)
	}
	available, err := svc.AvailableLanguages()
	if err != nil {
		t.Fatal(err)
	}
	matcher := language.NewMatcher(available)
	_, i, _ := matcher.Match(wanted...)
	best := available[i]
	if best != terminology.BritishEnglish.Tag() {
		t.Fatalf("Didn't correctly match British English. Matched %v", best)
	}
}
