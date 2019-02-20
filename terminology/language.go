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

package terminology

import (
	"golang.org/x/text/language"
	"log"
)

// Language defines a mapping between standard ISO language tags and the associated SNOMED-CT language reference sets
// TODO: add more supported languages
// TODO: check that the refset identifiers are correct
// TODO: add tests for other languages
type Language int

// Supported languages
const (
	AmericanEnglish Language = iota
	BritishEnglish
	French
	Spanish
	Danish
	lastLanguage
)

var tags = map[Language]language.Tag{
	BritishEnglish:  language.BritishEnglish,
	AmericanEnglish: language.AmericanEnglish,
	French:          language.French,
	Spanish:         language.Spanish,
	Danish:          language.Danish,
}

var identifiers = map[Language]int64{
	BritishEnglish:  999001261000000100,
	AmericanEnglish: 900000000000508004,
	French:          722131000,
	Spanish:         0,
	Danish:          554831000005107,
}

// Tag returns the language tag for this language
func (l Language) Tag() language.Tag {
	return tags[l]
}

// String returns the string representation of this language
func (l Language) String() string {
	return l.Tag().String()
}

// LanguageReferenceSetIdentifier returns the SNOMED-CT identifier for the language reference set for this language
func (l Language) LanguageReferenceSetIdentifier() int64 {
	return identifiers[l]
}

// newMatcher returns a language matcher that can be used to find the best service supported
// language given a user's requested preferences.
func newMatcher(svc *Svc) language.Matcher {
	allTags := make([]language.Tag, 0, len(tags))
	installed, err := svc.InstalledReferenceSets()
	if err != nil && err != ErrDatabaseNotInitialised {
		panic(err)
	}
	for l, v := range tags {
		refset := identifiers[l]
		if refset != 0 {
			for m := range installed {
				if _, ok := installed[m]; ok {
					allTags = append(allTags, v)
					break
				}
			}
		}
	}
	return language.NewMatcher(allTags)
}

// Match takes a list of requested languages and identifies the best supported match
func (svc *Svc) Match(preferred []language.Tag) Language {
	matchedTag, _, _ := svc.Matcher.Match(preferred...)
	for language, tag := range tags {
		if tag == matchedTag {
			return language
		}
	}
	log.Printf("failed to match language %s", matchedTag)
	return AmericanEnglish
}
