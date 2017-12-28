package database

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"database/sql"
	"fmt"
	"github.com/lib/pq" // imported to nicely handle arrays with placeholders
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

// SQLService is a concrete database-backed service for SNOMED-CT
// TODO: adopt a more sophisticated cache
type SQLService struct {
	db                      *sql.DB
	cache                   *NaiveCache // cache for concepts, relationships and descriptions by id
	parentRelationshipCache *NaiveCache // cache for relationships by concept id
	childRelationshipCache  *NaiveCache // cache for relationships by concept id
	descriptionCache        *NaiveCache // cache for descriptions by concept id
}

// this is to ensure that, at compile-time, our database service is a valid implementation of Service
var _ Service = (*SQLService)(nil)

// NewSQLService creates a new database-backed service using the database specified.
// TODO: allow customisation of language preferences, useful when getting preferred descriptions
// TODO: add more sophisticated caching
func NewSQLService(db *sql.DB) Service {
	return &SQLService{db, NewCache(), NewCache(), NewCache(), NewCache()}
}

// SQL statements
const (
	// simple fetch of one or more concepts and a list of recursive parents
	sqlFetchConcept = `select concept_id, fully_specified_name, concept_status_code,
	string_agg(parent_concept_id::text,',') as parents
	from t_concept left join t_cached_parent_concepts on 
	child_concept_id=concept_id 
	where concept_id=ANY($1) group by concept_id`

	sqlFetchAllConcepts = `select concept_id, fully_specified_name, concept_status_code,
	string_agg(parent_concept_id::text,',') as parents
	from t_concept left join t_cached_parent_concepts on 
	child_concept_id=concept_id 
	where parent_concept_id=ANY($1) group by concept_id`

	// fetch all recursive children for a given single concept
	sqlRecursiveChildren = `select child_concept_id from t_cached_parent_concepts where parent_concept_id=($1)`

	// fetch all parent relationships for a given single concept (relationships in which this is the source)
	sqlTargetRelationships = `select relationship_id, source_concept_id, relationship_type_concept_id, target_concept_id 
	from t_relationship
	where source_concept_id=($1)`

	// fetch all child relationships for a given single concept (relationships in which this is the target)
	sqlSourceRelationships = `select relationship_id, source_concept_id, relationship_type_concept_id, target_concept_id 
	from t_relationship
	where target_concept_id=($1)`

	// fetch all relationships for caching purposes
	sqlAllRelationships = `select relationship_id, source_concept_id, relationship_type_concept_id, target_concept_id 
	from t_relationship`

	// fetch all descriptions for a given single concept
	sqlDescriptions = `select description_id, description_status_code, description_type_code, initial_capital_status, language_code, term
	from t_description 
	where concept_id=($1)`
)

// GetDescriptions returns all of the descriptions (synonyms) for the given concept
func (ds *SQLService) GetDescriptions(concept *snomed.Concept) ([]*snomed.Description, error) {
	conceptID := int(concept.ConceptID)
	value, ok := ds.descriptionCache.Get(conceptID)
	if ok {
		return value.([]*snomed.Description), nil
	}
	rows, err := ds.db.Query(sqlDescriptions, concept.ConceptID)
	if err != nil {
		return nil, err
	}
	descriptions, err := rowsToDescriptions(rows)
	if err == nil {
		ds.descriptionCache.Put(conceptID, descriptions)
	}
	return descriptions, err
}

// PrecacheRelationships is a quick hack to preload all relationships for all concepts into in-memory cache.
// This caches only relationships for any concepts already cached.
func (ds *SQLService) PrecacheRelationships() error {
	rows, err := ds.db.Query(sqlAllRelationships)
	if err != nil {
		return err
	}
	defer rows.Close()
	relations, err := rowsToRelationships(rows)
	if err != nil {
		return err
	}
	for _, relation := range relations {
		precacheRelationship(ds.parentRelationshipCache, relation, relation.Target)
		precacheRelationship(ds.childRelationshipCache, relation, relation.Source)
	}
	return nil
}

func precacheRelationship(cache *NaiveCache, relation *snomed.Relationship, conceptID snomed.Identifier) {
	cached, ok := cache.Get(conceptID.AsInteger())
	if !ok {
		cached = make([]*snomed.Relationship, 0, 1)
	}
	cached = append(cached.([]*snomed.Relationship), relation)
	cache.Put(conceptID.AsInteger(), cached)
}

// GetParentRelationships returns the relationships for a concept in which it is the source.
func (ds *SQLService) GetParentRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error) {
	conceptID := int(concept.ConceptID)
	value, ok := ds.parentRelationshipCache.Get(conceptID)
	if ok {
		return value.([]*snomed.Relationship), nil
	}
	rows, err := ds.db.Query(sqlTargetRelationships, concept.ConceptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	relations, err := rowsToRelationships(rows)
	if err == nil {
		ds.parentRelationshipCache.Put(conceptID, relations)
	}
	return relations, err
}

// GetChildRelationships returns the relationships for a concept in which it is the target.
func (ds *SQLService) GetChildRelationships(concept *snomed.Concept) ([]*snomed.Relationship, error) {
	conceptID := int(concept.ConceptID)
	value, ok := ds.childRelationshipCache.Get(conceptID)
	if ok {
		return value.([]*snomed.Relationship), nil
	}
	rows, err := ds.db.Query(sqlSourceRelationships, concept.ConceptID)
	if err != nil {
		return nil, err
	}
	relations, err := rowsToRelationships(rows)
	if err == nil {
		ds.childRelationshipCache.Put(conceptID, relations)
	}
	return relations, err
}

// GetConcept fetches a concept with the given identifier
func (ds *SQLService) GetConcept(conceptID int) (*snomed.Concept, error) {
	return ds.cache.GetConceptOrElse(conceptID, func(conceptID int) (interface{}, error) {
		fetched, err := ds.performFetchConcepts(conceptID)
		if err != nil {
			return nil, err
		}
		concept := fetched[conceptID]
		if concept == nil {
			return nil, fmt.Errorf("No concept found with identifier %d", conceptID)
		}
		return concept, nil
	})
}

// GetConcepts returns a list of concepts with the given identifiers
func (ds *SQLService) GetConcepts(conceptIDs ...int) ([]*snomed.Concept, error) {
	l := len(conceptIDs)
	result := make([]*snomed.Concept, l)
	fetch := make([]int, 0, l)
	for i, conceptID := range conceptIDs {
		cached, ok := ds.cache.GetConcept(conceptID)
		if ok {
			result[i] = cached
		} else {
			fetch = append(fetch, conceptID)
		}
	}
	// perform fetch for concepts not in cache
	fetched, err := ds.performFetchConcepts(fetch...)
	if err != nil {
		return nil, err
	}
	// iterate through cached results and fill in blanks from fetched, populating cache as we go
	for i, concept := range result {
		if concept == nil {
			conceptID := conceptIDs[i]
			concept = fetched[conceptID]
			if concept != nil {
				ds.cache.PutConcept(conceptID, concept)
				result[i] = concept
			} else {
				return nil, fmt.Errorf("Invalid concept identifier: %d", conceptID)
			}
		}
	}
	return result, nil
}

// GetRecursiveChildrenIds fetches a list of identifiers representing all children of the given concept.
func (ds *SQLService) GetRecursiveChildrenIds(concept *snomed.Concept) ([]int, error) {
	rows, err := ds.db.Query(sqlRecursiveChildren, concept.ConceptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result = make([]int, 0, 10)
	for rows.Next() {
		var childConceptID int
		err = rows.Scan(&childConceptID)
		if err != nil {
			return nil, err
		}
		result = append(result, childConceptID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// Precache loads all of the SNOMED-CT dataset into the in-memory cache for speed.
func (ds *SQLService) Precache(rootConceptIDs ...int) error {
	err := ds.PrecacheConcepts(rootConceptIDs...)
	if err != nil {
		return err
	}
	err = ds.PrecacheRelationships()
	return err
}

// PrecacheConcepts loads all concepts into the cache
func (ds *SQLService) PrecacheConcepts(rootConceptIDs ...int) error {
	rows, err := ds.db.Query(sqlFetchAllConcepts, pq.Array(rootConceptIDs))
	if err != nil {
		return err
	}
	defer rows.Close()
	concepts, err := rowsToConcepts(rows)
	fmt.Printf("Fetched %d concepts\n", len(concepts))
	for conceptID, concept := range concepts {
		ds.cache.PutConcept(conceptID, concept)
	}
	return err
}

// Close closes the database
func (ds *SQLService) Close() error {
	return ds.db.Close()
}

// SliceToMap is a simple convenience method to convert a slice of concepts to a map
func SliceToMap(concepts []*snomed.Concept) map[snomed.Identifier]*snomed.Concept {
	l := len(concepts)
	r := make(map[snomed.Identifier]*snomed.Concept, l)
	for _, c := range concepts {
		r[c.ConceptID] = c
	}
	return r
}

// MapToSlice is a simple convenience method to convert a map of concepts to a slice
func MapToSlice(concepts map[snomed.Identifier]*snomed.Concept) []*snomed.Concept {
	l := len(concepts)
	r := make([]*snomed.Concept, 0, l)
	for _, c := range concepts {
		r = append(r, c)
	}
	return r

}

func (ds *SQLService) performFetchConcepts(conceptIDs ...int) (map[int]*snomed.Concept, error) {
	rows, err := ds.db.Query(sqlFetchConcept, pq.Array(conceptIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	concepts, err := rowsToConcepts(rows)
	return concepts, nil
}

func rowsToConcepts(rows *sql.Rows) (map[int]*snomed.Concept, error) {
	concepts := make(map[int]*snomed.Concept)
	var (
		conceptID          int
		fullySpecifiedName string
		conceptStatusCode  int
		parents            sql.NullString // may be null for root concept
	)
	for rows.Next() {
		err := rows.Scan(&conceptID, &fullySpecifiedName, &conceptStatusCode, &parents)
		if err != nil {
			return nil, err
		}
		concept, err := snomed.NewConcept(snomed.Identifier(conceptID), fullySpecifiedName, conceptStatusCode, ListAtoi(parents.String))
		if err != nil {
			return nil, err
		}
		concepts[conceptID] = concept
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return concepts, nil
}

func rowsToRelationships(rows *sql.Rows) ([]*snomed.Relationship, error) {
	relationships := make([]*snomed.Relationship, 0, 10)
	var (
		relationshipID  int
		sourceConceptID int
		typeConceptID   int
		targetConceptID int
	)
	for rows.Next() {
		err := rows.Scan(&relationshipID, &sourceConceptID, &typeConceptID, &targetConceptID)
		if err != nil {
			return nil, err
		}
		relationship := snomed.NewRelationship(snomed.Identifier(relationshipID), snomed.Identifier(sourceConceptID), snomed.Identifier(typeConceptID), snomed.Identifier(targetConceptID))
		relationships = append(relationships, relationship)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return relationships, nil
}

func rowsToDescriptions(rows *sql.Rows) ([]*snomed.Description, error) {
	descriptions := make([]*snomed.Description, 0, 10)
	var (
		descriptionID         int
		descriptionStatusCode int
		descriptionTypeCode   int
		initialCapitalStatus  int
		languageCode          string
		term                  string
	)
	for rows.Next() {
		err := rows.Scan(&descriptionID, &descriptionStatusCode, &descriptionTypeCode, &initialCapitalStatus, &languageCode, &term)
		if err != nil {
			return nil, err
		}
		tag, err := language.Parse(languageCode)
		if err != nil {
			return nil, err
		}
		description := &snomed.Description{DescriptionID: snomed.Identifier(descriptionID), Status: snomed.DescriptionStatus(descriptionStatusCode), Type: snomed.DescriptionType(descriptionTypeCode), InitialCapitalStatus: initialCapitalStatus, LanguageCode: tag, Term: term}
		descriptions = append(descriptions, description)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return descriptions, nil
}

// ListAtoi converts a comma-delimited string containing integers into a slice of integers
// TODO: move to utility package or find a way to make redundant
func ListAtoi(list string) []int {
	slist := strings.Split(strings.Replace(list, " ", "", -1), ",")
	r := make([]int, 0)
	for _, s := range slist {
		v, err := strconv.Atoi(s)
		if err == nil {
			r = append(r, v)
		}
	}
	return r
}

// ListItoA converts a slice of integers into a comma-delimited string
// TODO: move to utility package or find a way to make redundant
func ListItoA(list []int) string {
	r := make([]string, 0)
	for _, i := range list {
		s := strconv.Itoa(i)
		r = append(r, s)
	}
	return strings.Join(r, ",")
}
