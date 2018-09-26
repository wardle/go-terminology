package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"golang.org/x/text/language"
	"net/http"
	"strconv"
)

// C represents a returned Concept including useful additional information
type C struct {
	*snomed.Concept
	IsA                  []int64               `json:"isA"`
	Descriptions         []*snomed.Description `json:"descriptions"`
	PreferredDescription *snomed.Description   `json:"preferredDescription"`
	PreferredFsn         *snomed.Description   `json:"preferredFsn"`
	ReferenceSets        []int64               `json:"referenceSets"`
}

type dFilter struct {
	includeInactive bool // whether to include inactive as well as active descriptions
	includeFsn      bool // whether to include FSN description
}

// create a description filter based on the HTTP request
func newDFilter(r *http.Request) *dFilter {
	filter := &dFilter{includeInactive: false, includeFsn: false}
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
	conceptID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return result{nil, err, http.StatusBadRequest}
	}
	concept, err := svc.GetConcept(conceptID)
	if err != nil {
		return result{nil, err, http.StatusNotFound}
	}
	return resultForConcept(svc, r, concept)
}

// TODO(MW): choose default language from system environment variables or command-line option
func resultForConcept(svc *terminology.Svc, r *http.Request, concept *snomed.Concept) result {
	tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil {
		tags = []language.Tag{language.BritishEnglish}
	}
	descriptions, err := svc.GetDescriptions(concept)
	if err != nil {
		return result{nil, err, http.StatusInternalServerError}
	}
	preferredDescription := svc.MustGetPreferredSynonym(concept, tags)
	preferredFsn := svc.MustGetFullySpecifiedName(concept, tags)
	referenceSets, err := svc.GetReferenceSets(concept.Id)
	if err != nil {
		return result{nil, err, http.StatusInternalServerError}
	}
	allParents, err := svc.GetAllParentIDs(concept)
	if err != nil {
		return result{nil, err, http.StatusInternalServerError}
	}
	return result{&C{
		Concept:              concept,
		IsA:                  allParents,
		Descriptions:         newDFilter(r).filter(descriptions),
		PreferredDescription: preferredDescription,
		PreferredFsn:         preferredFsn,
		ReferenceSets:        referenceSets,
	},
		nil, http.StatusOK}
}

func getConceptDescriptions(svc *terminology.Svc, w http.ResponseWriter, r *http.Request) result {

	params := mux.Vars(r)
	conceptID, err := strconv.ParseInt(params["id"], 10, 64)
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

func crossmap(svc *terminology.Svc, w http.ResponseWriter, r *http.Request) result {
	params := mux.Vars(r)
	componentID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return result{nil, err, http.StatusBadRequest}
	}
	refsetID := r.FormValue("refset")
	if refsetID == "" {
		return result{nil, fmt.Errorf("missing parameter: refset"), http.StatusBadRequest}
	}
	refset, err := strconv.ParseInt(refsetID, 10, 64)
	if err != nil {
		return result{nil, err, http.StatusBadRequest}
	}
	rsi, err := svc.GetFromReferenceSet(refset, componentID)
	if err != nil {
		return result{nil, err, http.StatusBadRequest}
	}
	if rsi == nil {
		return result{nil, err, http.StatusNotFound}
	}
	return result{rsi, nil, http.StatusOK}
}

// genericize maps a concept to an arbitrary root concept or to the best match in the specified refset
func genericize(svc *terminology.Svc, w http.ResponseWriter, r *http.Request) result {
	params := mux.Vars(r)
	conceptID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return result{nil, err, http.StatusBadRequest}
	}
	concept, err := svc.GetConcept(conceptID)
	if err != nil {
		return result{nil, err, http.StatusNotFound}
	}
	err = r.ParseForm()
	if err != nil {
		return result{nil, err, http.StatusBadRequest}
	}
	rootConceptIDs := r.Form["root"]
	if len(rootConceptIDs) > 0 {
		conceptIDs := make(map[int64]bool)
		for _, conceptID := range rootConceptIDs {
			root, err := strconv.ParseInt(conceptID, 10, 64)
			if err != nil {
				return result{nil, err, http.StatusBadRequest}
			}
			conceptIDs[root] = true
		}
		generic, ok := svc.GenericiseTo(concept, conceptIDs)
		if !ok {
			return result{nil, err, http.StatusNotFound}
		}
		return resultForConcept(svc, r, generic)
	}
	refsetID := r.FormValue("refset")
	if refsetID != "" {
		refset, err := strconv.ParseInt(refsetID, 10, 64)
		if err != nil {
			return result{nil, err, http.StatusBadRequest}
		}
		items, err := svc.GetReferenceSetItems(refset)
		if err != nil {
			return result{nil, err, http.StatusInternalServerError}
		}
		generic, ok := svc.GenericiseTo(concept, items)
		if !ok {
			return result{nil, fmt.Errorf("unable to genericise %d to a member of refset %d", conceptID, refset), http.StatusNotFound}
		}
		return resultForConcept(svc, r, generic)
	}
	return result{nil, fmt.Errorf("must specify either a root or refset"), http.StatusBadRequest}
}
