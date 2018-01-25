package server

import (
	"github.com/gorilla/mux"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
	"net/http"
	"strconv"
)

// C represents a returned Concept including useful additional information
// TODO: include derivation of preferredTerm for the locale requested
type C struct {
	*snomed.Concept
	IsA          []int                 `json:"isA"`
	Descriptions []*snomed.Description `json:"descriptions"`
}

type dFilter struct {
	langMatcher     language.Matcher // user accepted languages, may be nil
	includeInactive bool             // whether to include inactive as well as active descriptions
	includeFsn      bool             // whether to include FSN description
}

func parseLanguageMatcher(acceptedLanguage string) language.Matcher {
	if desired, _, err := language.ParseAcceptLanguage(acceptedLanguage); err == nil {
		return language.NewMatcher(desired)
	}
	return nil
}

// create a description filter based on the HTTP request
func newDFilter(r *http.Request) *dFilter {
	filter := &dFilter{langMatcher: nil, includeInactive: false, includeFsn: false}
	if accept := r.Header.Get("Accept-Language"); accept != "" {
		filter.langMatcher = parseLanguageMatcher(accept)
	}
	if includeInactive, err := strconv.ParseBool(r.FormValue("includeInactive")); err == nil {
		filter.includeInactive = includeInactive
	}
	if includeFsn, err := strconv.ParseBool(r.FormValue("includeFsn")); err == nil {
		filter.includeFsn = includeFsn
	}
	return filter
}

// filter a slice of descriptions
func (df *dFilter) filter(descriptions []*snomed.Description) []*snomed.Description {
	ds := make([]*snomed.Description, 0)
	for _, d := range descriptions {
		if df.test(d) {
			ds = append(ds, d)
		}
	}
	return ds
}

// test whether an individual description should be included or not
func (df *dFilter) test(d *snomed.Description) bool {
	if df.langMatcher != nil {
		_, _, conf := df.langMatcher.Match(d.LanguageTag())
		if conf < language.High {
			return false
		}
	}
	if d.Active == false && df.includeInactive == false {
		return false
	}
	if d.IsFullySpecifiedName() && df.includeFsn == false {
		return false
	}
	return true
}

// return a single concept
func getConcept(svc *terminology.Svc, w http.ResponseWriter, r *http.Request) result {
	params := mux.Vars(r)
	conceptID, err := strconv.Atoi(params["id"])
	if err != nil {
		return result{nil, err, http.StatusBadRequest}
	}
	concept, err := svc.GetConcept(conceptID)
	if err != nil {
		return result{nil, err, http.StatusNotFound}
	}
	descriptions, err := svc.GetDescriptions(concept)
	if err != nil {
		return result{nil, err, http.StatusInternalServerError}
	}
	allParents, err := svc.GetAllParentIDs(concept)
	if err != nil {
		return result{nil, err, http.StatusInternalServerError}
	}
	return result{&C{concept, allParents, newDFilter(r).filter(descriptions)}, nil, http.StatusOK}
}

func getConceptDescriptions(svc *terminology.Svc, w http.ResponseWriter, r *http.Request) result {
	params := mux.Vars(r)
	conceptID, err := strconv.Atoi(params["id"])
	if err != nil {
		return result{nil, err, http.StatusBadRequest}
	}
	concept, err := svc.GetConcept(conceptID)
	if err != nil {
		return result{nil, err, http.StatusNotFound}
	}
	descriptions, err := svc.GetDescriptions(concept)
	if err != nil {
		return result{nil, err, http.StatusInternalServerError}
	}
	return result{newDFilter(r).filter(descriptions), nil, http.StatusOK}
}
