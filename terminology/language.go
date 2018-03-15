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
)

// Language defines a mapping between standard ISO language tags and the associated SNOMED-CT language reference sets
type Language int

// Supported languages
const (
	BritishEnglish Language = iota
	AmericanEnglish
	lastLanguage
)

var tags = map[Language]language.Tag{
	BritishEnglish:  language.BritishEnglish,
	AmericanEnglish: language.AmericanEnglish,
}

var identifiers = map[Language]int64{
	BritishEnglish:  999001261000000100,
	AmericanEnglish: 900000000000508004,
}

// Tag returns the language tag for this language
func (l Language) Tag() language.Tag {
	return tags[l]
}

// LanguageReferenceSetIdentifier returns the SNOMED-CT identifier for the language reference set for this language
func (l Language) LanguageReferenceSetIdentifier() int64 {
	return identifiers[l]
}

// NewMatcher returns a language matcher that can be used to find the best service supported
// language given a user's requested preferences.
func NewMatcher(svc Svc) language.Matcher {
	allTags := make([]language.Tag, 0, len(tags))
	installed, err := svc.GetReferenceSets()
	if err != nil {
		panic(err)
	}
	for l, v := range tags {
		refset := identifiers[l]
		for m := range installed {
			if installed[m] == refset {
				allTags = append(allTags, v)
				break
			}
		}
	}
	return language.NewMatcher(allTags)
}

// Match takes a list of requested languages and identifies the best supported match
func Match(svc Svc, preferred []language.Tag) int64 {
	_, index, _ := NewMatcher(svc).Match(preferred...)
	return Language(index).LanguageReferenceSetIdentifier()
}
