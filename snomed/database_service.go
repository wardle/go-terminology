package snomed

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq" // imported to nicely handle arrays with placeholders
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

// NaiveCache is an extraordinarily naive in-memory cache used for development
type NaiveCache struct {
	cache map[int]*Concept
}

// Put stores a concept in the cache
func (nc NaiveCache) Put(conceptID int, concept *Concept) {
	nc.cache[conceptID] = concept
}

// Get fetches a concept from the cache
func (nc NaiveCache) Get(conceptID int) *Concept {
	return nc.cache[conceptID]
}

// GetOrElse fetches a concept from the cache or performs the closure specified, caching the result
func (nc NaiveCache) GetOrElse(conceptID int, f func(conceptID int) (*Concept, error)) (*Concept, error) {
	concept := nc.cache[conceptID]
	if concept != nil {
		return concept, nil
	}
	concept, err := f(conceptID)
	if err == nil && concept != nil {
		nc.cache[conceptID] = concept
	}
	return concept, err
}

func (nc NaiveCache) String() string {
	return fmt.Sprintf("Cache with %d items.", len(nc.cache))
}

// DatabaseService is a concrete database-backed service for SNOMED-CT
type DatabaseService struct {
	db       *sql.DB
	language language.Tag
	cache    *NaiveCache
}

// NewDatabaseService creates a new database-backed service using the database specified.
// TODO: allow customisation of language preferences useful when getting preferred descriptions
// TODO: add more sophisticated caching
func NewDatabaseService(db *sql.DB) *DatabaseService {
	cache := &NaiveCache{make(map[int]*Concept)}
	return &DatabaseService{db, language.BritishEnglish, cache}
}

// SQL statements
const (
	// simple fetch of a concept and a list of recursive parents
	sqlFetchConcept = `select concept_id, fully_specified_name, concept_status_code,
	string_agg(parent_concept_id::text,',') as parents
	from t_concept left join t_cached_parent_concepts on 
	child_concept_id=concept_id 
	where concept_id=ANY($1) group by concept_id`
)

// FetchConcept fetches a concept with the given identifier
func (ds DatabaseService) FetchConcept(conceptID int) (*Concept, error) {
	return ds.cache.GetOrElse(conceptID, func(conceptID int) (*Concept, error) {
		fetched, err := ds.performFetchConcepts(conceptID)
		if err != nil {
			return nil, err
		}
		return fetched[conceptID], nil
	})
}

// FetchConcepts returns a list of concepts with the given identifiers
func (ds DatabaseService) FetchConcepts(conceptIDs ...int) ([]*Concept, error) {
	l := len(conceptIDs)
	result := make([]*Concept, l)
	fetch := make([]int, 0, l)
	for i, conceptID := range conceptIDs {
		cached := ds.cache.Get(conceptID)
		if cached != nil {
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
				ds.cache.Put(conceptID, concept)
				result[i] = concept
			} else {
				return nil, fmt.Errorf("Invalid concept identifier: %d", conceptID)
			}
		}
	}
	return result, nil
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
