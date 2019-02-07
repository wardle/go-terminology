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
	"bufio"
	"fmt"
	"golang.org/x/text/language"
	"io"
	"strconv"
	"strings"
)

// Print prints information about each concept
// TODO: move to different package
func Print(svc *Svc, reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		id, err := strconv.ParseInt(string(line), 10, 64)
		if err != nil {
			return err
		}
		c, err := svc.Concept(id)
		if err != nil {
			return err
		}
		fsn, found, err := svc.FullySpecifiedName(c, []language.Tag{BritishEnglish.Tag()})
		if err != nil {
			return err
		}
		if !found {
			// TODO: fallback?
		}
		fmt.Println(fsn.GetTerm())
	}
	return nil
}

// NumberFactors gives the number of unique factors in the data specified
func NumberFactors(reader io.Reader) (int, error) {
	scanner := bufio.NewScanner(reader)
	data := make(map[int64]struct{}, 0)
	for scanner.Scan() {
		line := scanner.Text()
		id, err := strconv.ParseInt(string(line), 10, 64)
		if err != nil {
			return 0, err
		}
		data[id] = struct{}{}
	}
	return len(data), nil
}

// Reducer performs dimensionality reduction
type Reducer struct {
	svc             *Svc
	maximumFactors  int // maximum number of factors permitted in your data
	minimumDistance int // minimum distance from root
	data            map[int64]*reducingConcept
	result          map[int64]int64
}

// NewReducer creates a new Reducer processor
func NewReducer(svc *Svc, maximumFactors int, minimumDistance int) *Reducer {
	if maximumFactors < 1 {
		maximumFactors = 1
	}
	return &Reducer{
		svc:             svc,
		maximumFactors:  maximumFactors,
		minimumDistance: minimumDistance,
		data:            make(map[int64]*reducingConcept),
		result:          make(map[int64]int64),
	}
}

type reducingConcept struct {
	id         int64
	pathToRoot []int64
	count      int
	mappedTo   *reducingConcept
	score      int
}

func (r *Reducer) mapped(id int64) *reducingConcept {
	if c, exists := r.data[id]; exists {
		if c.mappedTo != nil {
			return r.mapped(c.mappedTo.id)
		}
		return c
	}
	return nil
}

func (rc *reducingConcept) String() string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(rc.id, 10))
	sb.WriteString("[")
	for _, p := range rc.pathToRoot {
		sb.WriteString(strconv.FormatInt(p, 10))
		sb.WriteString("-")
	}
	sb.WriteString("](count:")
	sb.WriteString(strconv.Itoa(rc.count))
	sb.WriteString(",score:")
	sb.WriteString(strconv.Itoa(rc.score))
	sb.WriteString(")")
	return sb.String()
}

func (rc *reducingConcept) distanceFromRoot() int {
	return len(rc.pathToRoot)
}

// determine the number of unique factors in the data
func (r *Reducer) df() int {
	df := 0
	for _, c := range r.data {
		if c.count > 0 {
			df++
		}
	}
	return df
}

// Reduce reduces the list of concepts
func (r *Reducer) Reduce(concepts []int64) ([]int64, error) {
	for _, c := range concepts {
		if _, err := r.add(c); err != nil {
			return nil, err
		}
	}
	if err := r.reduce(); err != nil {
		return nil, err
	}
	result := make([]int64, len(concepts))
	for i, c := range concepts {
		result[i] = r.mapped(c).id
	}
	return result, nil
}

// ReduceCsv processes a csv file to reduce its dimensionality by genericising SNOMED-CT concepts
func (r *Reducer) ReduceCsv(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	source := make([]int64, 0)
	for scanner.Scan() {
		line := scanner.Text()
		id, err := strconv.ParseInt(string(line), 10, 64)
		if err != nil {
			return err
		}
		source = append(source, id)
		if _, err = r.add(id); err != nil {
			return err
		}
	}
	r.reduce()
	for _, id := range source {
		//fmt.Fprintf(writer, "%10d -- %10d\n", id, r.mapped(id))
		fmt.Fprintln(writer, strconv.FormatInt(r.mapped(id).id, 10))
	}
	return nil
}

// add adds an item of data to be reduced.
func (r *Reducer) add(id int64) (*reducingConcept, error) {
	if c, exists := r.data[id]; exists {
		c.count = c.count + 1
		return c, nil
	}
	path, err := r.shortestPathToRoot(id)
	if err != nil {
		return nil, err
	}
	c := &reducingConcept{id: id, pathToRoot: path, count: 1, mappedTo: nil}
	r.data[id] = c
	return c, nil
}

// reduce runs multiple passes of genericisation to reduce dimensionality of dataset.
func (r *Reducer) reduce() error {
	for r.df() > r.maximumFactors {
		//fmt.Printf("reduce(): degrees of freedom : %d, target: %d, so running reduction\n", r.df(), r.maximumFactors)
		if r.execute() == false {
			break
		}
	}
	if r.df() > r.maximumFactors {
		return fmt.Errorf("warning: reduced to %d factors, not %d due to minumum distance constraint %d", r.df(), r.maximumFactors, r.minimumDistance)
	}
	return nil
}

// execute executes a single pass of dimensionality reduction, returning whether it is possible to do more
func (r *Reducer) execute() bool {
	lowest := r.calculateScores()
	if lowest == -1 {
		return false
	}
	for id, c := range r.data {
		if c.count > 0 && c.score != -1 {
			if lowest == c.score && (c.distanceFromRoot()-1) > r.minimumDistance {
				r.genericise(id)
				return true // short-circuit
			}
		}
	}
	return true
}

// calculate a score to identify the concepts to be genericised
// score determined by the frequency principally but includes a score
// relating to the distance away from root, so that genericisation
// will occur in a concept furthest away from root, if frequencies are the same.
func (r *Reducer) calculateScores() int {
	lowestScore := -1
	maxDistance := r.maxDistance() + 1
	for _, c := range r.data {
		if c.count > 0 {
			l := len(c.pathToRoot)
			if c.distanceFromRoot() == r.minimumDistance {
				c.score = -1
			} else {
				c.score = c.count * (maxDistance - l)
			}
			if lowestScore == -1 || (lowestScore > c.score && c.score != -1) {
				lowestScore = c.score
			}
			//	fmt.Printf("%v\n", c)
		}
	}
	return lowestScore
}

// calculate the maximum distance from root for all data, not including those already mapped
func (r *Reducer) maxDistance() int {
	maxDistance := 0
	for _, c := range r.data {
		if c.count > 0 {
			d := len(c.pathToRoot)
			if d > maxDistance {
				maxDistance = d
			}
		}
	}
	return maxDistance
}

// genericise the specified concept turning it into a more generic one
func (r *Reducer) genericise(id int64) error {
	c := r.data[id]
	targetID := c.pathToRoot[1]
	target, err := r.add(targetID)
	if err != nil {
		return err
	}
	target.count += c.count
	c.count = 0
	c.score = -1
	c.mappedTo = target
	return nil
}

func (r *Reducer) shortestPathToRoot(conceptID int64) ([]int64, error) {
	c, err := r.svc.Concept(conceptID)
	if err != nil {
		return nil, err
	}
	path, err := r.svc.ShortestPathToRoot(c)
	if err != nil {
		return nil, err
	}
	distance := len(path)
	result := make([]int64, 0, distance)
	for _, p := range path {
		result = append(result, p.Id)
	}
	return result, nil
}
