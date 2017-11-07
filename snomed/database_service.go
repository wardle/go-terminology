package snomed

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq" // imported to nicely handle arrays with placeholders
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

// DatabaseService is a concrete database-backed service for SNOMED-CT
type DatabaseService struct {
	db                      *sql.DB
	language                language.Tag
	cache                   *NaiveCache // cache for concepts, relationships and descriptions by id
	parentRelationshipCache *NaiveCache // cache for relationships by concept id
	childRelationshipCache  *NaiveCache // cache for relationships by concept id
	descriptionCache        *NaiveCache // cache for descriptions by concept id
}

// NewDatabaseService creates a new database-backed service using the database specified.
// TODO: allow customisation of language preferences, useful when getting preferred descriptions
// TODO: add more sophisticated caching
func NewDatabaseService(db *sql.DB) *DatabaseService {
	return &DatabaseService{db, language.BritishEnglish, NewCache(), NewCache(), NewCache(), NewCache()}
}

// SQL statements
const (
	// simple fetch of one or more concepts and a list of recursive parents
	sqlFetchConcept = `select concept_id, fully_specified_name, concept_status_code,
	string_agg(parent_concept_id::text,',') as parents
	from t_concept left join t_cached_parent_concepts on 
	child_concept_id=concept_id 
	where concept_id=ANY($1) group by concept_id`

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

	// fetch all descriptions for a given single concept
	sqlDescriptions = `select description_id, description_status_code, description_type_code, initial_capital_status, language_code, term
	from t_description 
	where concept_id=($1)`
)

// GetPreferredDescription returns the preferred description for this concept in the default language for this service.
func (ds DatabaseService) GetPreferredDescription(concept *Concept) (*Description, error) {
	return ds.GetPreferredDescriptionForLanguages(concept, []language.Tag{ds.language})
}

// GetPreferredDescriptionForLanguages returns the preferred description for this concept in the languages specified
func (ds DatabaseService) GetPreferredDescriptionForLanguages(concept *Concept, languages []language.Tag) (*Description, error) {
	preferred, err := ds.GetPreferredDescriptions(concept)
	if err != nil {
		return nil, err
	}
	matcher := language.NewMatcher(languages)
	tags := make([]language.Tag, 0, len(preferred))
	for _, d := range preferred {
		tags = append(tags, d.LanguageCode)
	}
	_, index, _ := matcher.Match(tags...)
	return preferred[index], nil
}

// GetPreferredDescriptions returns the preferred descriptions for the given concept
func (ds DatabaseService) GetPreferredDescriptions(concept *Concept) ([]*Description, error) {
	descriptions, err := ds.GetDescriptions(concept)
	if err != nil {
		return nil, err
	}
	preferred := make([]*Description, 0, len(descriptions))
	for _, description := range descriptions {
		if description.Type.IsPreferred() {
			preferred = append(preferred, description)
		}
	}
	return preferred, nil
}

// GetDescriptions returns all of the descriptions (synonyms) for the given concept
func (ds DatabaseService) GetDescriptions(concept *Concept) ([]*Description, error) {
	conceptID := int(concept.ConceptID)
	value, ok := ds.descriptionCache.Get(conceptID)
	if ok {
		return value.([]*Description), nil
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

// GetSiblings returns the siblings of this concept, ie: those who share the same parents
func (ds DatabaseService) GetSiblings(concept *Concept) ([]*Concept, error) {
	parents, err := ds.GetParents(concept)
	if err != nil {
		return nil, err
	}
	siblings := make([]*Concept, 0, 10)
	for _, parent := range parents {
		children, err := ds.GetChildren(parent)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			if child.ConceptID != concept.ConceptID {
				siblings = append(siblings, child)
			}
		}
	}
	return siblings, nil
}

// GetParents returns the direct IS-A relations of the specified concept.
func (ds DatabaseService) GetParents(concept *Concept) ([]*Concept, error) {
	return ds.GetParentsOfKind(concept, IsA)
}

// GetParentsOfKind returns the relations of the specified kind (type) of the specified concept.
func (ds DatabaseService) GetParentsOfKind(concept *Concept, kind Identifier) ([]*Concept, error) {
	relations, err := ds.FetchParentRelationships(concept)
	if err != nil {
		return nil, err
	}
	conceptIDs := make([]int, 0, len(relations))
	for _, relation := range relations {
		if relation.Type == kind {
			conceptIDs = append(conceptIDs, int(relation.Target))
		}
	}
	return ds.FetchConcepts(conceptIDs...)
}

// PathsToRoot returns the different possible paths to the root SNOMED-CT concept from this one.
func (ds DatabaseService) PathsToRoot(concept *Concept) ([][]*Concept, error) {
	return ds.pathsToRoot(nil, concept)
}

func debugPaths(paths [][]*Concept) {
	for i, path := range paths {
		fmt.Printf("Path %d: ", i)
		debugPath(path)
	}
}

func debugPath(path []*Concept) {
	for _, concept := range path {
		fmt.Printf("%d-", concept.ConceptID)
	}
	fmt.Print("\n")
}

// pathsToRoot recursively determines the paths from the concept to the root SNOMED-CT concept
func (ds DatabaseService) pathsToRoot(currentPath []*Concept, concept *Concept) ([][]*Concept, error) {
	parents, err := ds.GetParents(concept)
	if err != nil {
		return nil, err
	}
	if currentPath == nil {
		currentPath = make([]*Concept, 0, 1)
	}
	currentPath = append(currentPath, concept)
	results := make([][]*Concept, 0, len(parents)*2)
	if len(parents) == 0 { // if we're at the top of the hierarchy, add the current path
		results = append(results, currentPath)
	}
	for _, parent := range parents { // otherwise, recursively process parents
		parentResults, err := ds.pathsToRoot(currentPath, parent)
		if err != nil {
			return nil, err
		}
		results = append(results, parentResults...) // if not, keep going
	}
	return results, err
}

// FetchParentRelationships returns the relationships for a concept in which it is the source.
func (ds DatabaseService) FetchParentRelationships(concept *Concept) ([]*Relationship, error) {
	conceptID := int(concept.ConceptID)
	value, ok := ds.parentRelationshipCache.Get(conceptID)
	if ok {
		return value.([]*Relationship), nil
	}
	rows, err := ds.db.Query(sqlTargetRelationships, concept.ConceptID)
	if err != nil {
		return nil, err
	}
	relations, err := rowsToRelationships(rows)
	if err == nil {
		ds.parentRelationshipCache.Put(conceptID, relations)
	}
	return relations, err
}

// GetChildren returns the direct IS-A relations of the specified concept.
func (ds DatabaseService) GetChildren(concept *Concept) ([]*Concept, error) {
	return ds.GetChildrenOfKind(concept, IsA)
}

// GetChildrenOfKind returns the relations of the specified kind (type) of the specified concept.
func (ds DatabaseService) GetChildrenOfKind(concept *Concept, kind Identifier) ([]*Concept, error) {
	relations, err := ds.FetchChildRelationships(concept)
	if err != nil {
		return nil, err
	}
	conceptIDs := make([]int, 0, len(relations))
	for _, relation := range relations {
		if relation.Type == kind {
			conceptIDs = append(conceptIDs, int(relation.Source))
		}
	}
	return ds.FetchConcepts(conceptIDs...)
}

// FetchChildRelationships returns the relationships for a concept in which it is the target.
func (ds DatabaseService) FetchChildRelationships(concept *Concept) ([]*Relationship, error) {
	conceptID := int(concept.ConceptID)
	value, ok := ds.childRelationshipCache.Get(conceptID)
	if ok {
		return value.([]*Relationship), nil
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

// ConceptsForRelationship returns the concepts represented within a relationship
func (ds DatabaseService) ConceptsForRelationship(rel *Relationship) (source *Concept, kind *Concept, target *Concept, err error) {
	concepts, err := ds.FetchConcepts(int(rel.Source), int(rel.Type), int(rel.Target))
	if err != nil {
		return nil, nil, nil, err
	}
	return concepts[0], concepts[1], concepts[2], nil
}

// FetchConcept fetches a concept with the given identifier
func (ds DatabaseService) FetchConcept(conceptID int) (*Concept, error) {
	return ds.cache.GetConceptOrElse(conceptID, func(conceptID int) (*Concept, error) {
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

// FetchConcepts returns a list of concepts with the given identifiers
func (ds DatabaseService) FetchConcepts(conceptIDs ...int) ([]*Concept, error) {
	l := len(conceptIDs)
	result := make([]*Concept, l)
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

// FetchRecursiveChildrenIds fetches a list of identifiers representing all children of the given concept.
func (ds DatabaseService) FetchRecursiveChildrenIds(concept *Concept) ([]int, error) {
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

// FetchRecursiveChildren fetches all children of the given concept recursively.
// Use with caution with concepts at high levels of the hierarchy.
func (ds DatabaseService) FetchRecursiveChildren(concept *Concept) ([]*Concept, error) {
	children, err := ds.FetchRecursiveChildrenIds(concept)
	if err != nil {
		return nil, err
	}
	return ds.FetchConcepts(children...)
}

// GetAllParents returns all of the parents (recursively) for a given concept
func (ds DatabaseService) GetAllParents(concept *Concept) ([]*Concept, error) {
	return ds.FetchConcepts(concept.Parents...)
}

func (ds DatabaseService) performFetchConcepts(conceptIDs ...int) (map[int]*Concept, error) {
	rows, err := ds.db.Query(sqlFetchConcept, pq.Array(conceptIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	concepts, err := rowsToConcepts(rows)
	return concepts, nil
}

func rowsToConcepts(rows *sql.Rows) (map[int]*Concept, error) {
	concepts := make(map[int]*Concept)
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
		concept, err := NewConcept(Identifier(conceptID), fullySpecifiedName, conceptStatusCode, ListAtoi(parents.String))
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

func rowsToRelationships(rows *sql.Rows) ([]*Relationship, error) {
	relationships := make([]*Relationship, 0, 10)
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
		relationship := NewRelationship(Identifier(relationshipID), Identifier(sourceConceptID), Identifier(typeConceptID), Identifier(targetConceptID))
		relationships = append(relationships, relationship)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return relationships, nil
}

func rowsToDescriptions(rows *sql.Rows) ([]*Description, error) {
	descriptions := make([]*Description, 0, 10)
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
		description := &Description{Identifier(descriptionID), DescriptionStatus(descriptionStatusCode), DescriptionType(descriptionTypeCode), initialCapitalStatus, tag, term}
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
