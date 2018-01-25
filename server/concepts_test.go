package server

import (
	"github.com/wardle/go-terminology/snomed"
	"reflect"
	"testing"
)

type langTest struct {
	wanted       string // Accept-Languages format
	descriptions []*snomed.Description
	expected     []*snomed.Description
}

var en = &snomed.Description{Term: "english", LanguageCode: "en"}
var fr = &snomed.Description{Term: "francais", LanguageCode: "fr"}
var de = &snomed.Description{Term: "deutsch", LanguageCode: "de"}
var enUS = &snomed.Description{Term: "american", LanguageCode: "en-US"}
var enGB = &snomed.Description{Term: "british", LanguageCode: "en-GB"}

var langTests = []langTest{
	langTest{wanted: "en", descriptions: []*snomed.Description{en, fr, de}, expected: []*snomed.Description{en}},
	langTest{wanted: "en-US", descriptions: []*snomed.Description{en, fr, de, enGB, enUS}, expected: []*snomed.Description{en, enGB, enUS}},
	langTest{wanted: "fr", descriptions: []*snomed.Description{en, fr, de, enGB, enUS}, expected: []*snomed.Description{fr}},
	langTest{wanted: "de", descriptions: []*snomed.Description{en, fr, de, enGB, enUS}, expected: []*snomed.Description{de}},
	langTest{wanted: "en", descriptions: []*snomed.Description{en, fr, de, enUS, enGB}, expected: []*snomed.Description{en, enUS, enGB}},
}

func TestLanguage(t *testing.T) {
	for _, lt := range langTests {
		filter := &dFilter{langMatcher: parseLanguageMatcher(lt.wanted), includeFsn: false, includeInactive: true}
		ds := filter.filter(lt.descriptions)
		if reflect.DeepEqual(ds, lt.expected) == false {
			t.Errorf("expected %v, got %v", lt.expected, ds)
		}
	}
}
