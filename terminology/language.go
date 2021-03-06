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
	"sort"

	"golang.org/x/text/language"
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
	AmericanEnglish: language.AmericanEnglish,
	BritishEnglish:  language.BritishEnglish,
	French:          language.French,
	Spanish:         language.Spanish,
	Danish:          language.Danish,
}

var identifiers = map[Language]int64{
	AmericanEnglish: 900000000000509007,
	BritishEnglish:  999001261000000100,
	French:          722131000,
	Spanish:         450828004,
	Danish:          31000146106,
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

// LanguageForTag returns the language for the specified tag
func LanguageForTag(tag language.Tag) Language {
	for l, t := range tags {
		if t == tag {
			return l
		}
	}
	return AmericanEnglish
}

// AvailableLanguages returns the languages supported by the currently installed distribution
// Note: the sorting of the results is important for language matching, because the first matching language
// will be chosen, so we finish by sorting the result.
func (svc *Svc) AvailableLanguages() ([]language.Tag, error) {
	installed, err := svc.InstalledReferenceSets()
	if err != nil && err != ErrDatabaseNotInitialised {
		return nil, err
	}
	allTags := make([]language.Tag, 0)
	for l, t := range tags {
		if refsetID := identifiers[l]; refsetID != 0 {
			if _, ok := installed[refsetID]; ok {
				allTags = append(allTags, t)
			}
		}
	}
	sort.Slice(allTags, func(i, j int) bool {
		return LanguageForTag(allTags[i]) < LanguageForTag(allTags[j])
	})
	return allTags, nil
}
