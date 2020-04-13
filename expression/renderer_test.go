package expression

import (
	"strings"
	"testing"

	"golang.org/x/text/language"
	"google.golang.org/protobuf/proto"
)

func TestRoundtrip(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	defaultRenderer := NewDefaultRenderer()
	canonicalRenderer := NewCanonicalRenderer()
	for _, roundtrip := range roundtrips {
		s1 := roundtrip.expression
		e1, err := Parse(s1)
		if err != nil {
			t.Fatal(err)
		}
		s2, err := defaultRenderer.Render(e1)
		if err != nil {
			t.Fatal(err)
		}
		e2, err := Parse(s2)
		if err != nil {
			t.Fatal(err)
		}
		if proto.Equal(e1, e2) == false {
			t.Fatalf("failed to roundtrip. input:\n%s\noutput:\n%s", s1, s2)
		}
		s3, err := canonicalRenderer.Render(e1)
		if err != nil {
			t.Fatal(err)
		}
		e3, err := Parse(s3)
		if err != nil {
			t.Fatal(err)
		}
		if roundtrip.canonical != "" {
			if strings.EqualFold(s3, roundtrip.canonical) == false {
				t.Errorf("incorrect canonical form, expected:%s got:%s", roundtrip.canonical, s3)
			}
			e4, err := Parse(roundtrip.canonical)
			if err != nil {
				t.Fatal(err)
			}
			if Equal(e4, e1) == false {
				t.Fatalf("failed to roundtrip. input:\n%s\noutput:\n%s", roundtrip.canonical, roundtrip.expression)
			}
		}
		if Equal(e1, e3) == false {
			t.Fatalf("failed to roundtrip. input:\n%s\noutput:\n%s", s1, s3)
		}
	}
}

func TestUpdatingExpressions(t *testing.T) {
	svc := setUp(t)
	defer svc.Close()
	tags, _, _ := language.ParseAcceptLanguage("en-GB") // TODO(mw): better language support
	updatingRenderer := NewUpdatingRenderer(svc, tags)
	canonicalRenderer := NewCanonicalRenderer()
	s1 := "80146002 | Appendectomy |" // original - using term from international release, not GB
	s2 := "80146002|Appendicectomy|"  // expected result from updating renderer - using en-GB
	s3 := "80146002"                  // expected result from canonical renderer
	e1, err := Parse(s1)
	if err != nil {
		t.Fatal(err)
	}
	updated, err := updatingRenderer.Render(e1)
	if err != nil {
		t.Fatal(err)
	}
	if strings.EqualFold(s2, updated) == false {
		t.Errorf("failed to update term from:%s to:%s, got:%s", s1, s2, updated)
	}
	canonical, err := canonicalRenderer.Render(e1)
	if err != nil {
		t.Fatal(err)
	}
	if strings.EqualFold(s3, canonical) == false {
		t.Errorf("failed to generate canonical expression from:%s to:%s, got:%s", s1, s3, canonical)
	}
}
