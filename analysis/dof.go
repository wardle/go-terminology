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

package analysis

import (
	"bufio"
	"fmt"
	"github.com/wardle/go-terminology/terminology"
	"io"
	"strconv"
)

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
	svc             *terminology.Svc
	maximumFactors  int // maximum number of factors permitted in your data
	minimumDistance int // minimum distance from root
	data            map[int64]*reducingConcept
	result          map[int64]int64
}

// NewReducer creates a new Reducer processor
func NewReducer(svc *terminology.Svc, maximumFactors int, minimumDistance int) *Reducer {
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
	mappedTo   int
	score      int
}

func (rc *reducingConcept) distanceFromRoot() int {
	return len(rc.pathToRoot) - rc.mappedTo
}

func (r *Reducer) df() int {
	df := 0
	for _, c := range r.data {
		if c.count > 0 {
			df++
		}
	}
	return df
}

func (r *Reducer) mapped(id int64) int64 {
	c := r.data[id]
	if c.mappedTo > 0 {
		return r.mapped(c.pathToRoot[c.mappedTo])
	}
	return id
}

// Reduce processes a csv file to reduce its dimensionality by genericising SNOMED-CT concepts
func (r *Reducer) Reduce(reader io.Reader, writer io.Writer) error {
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
	for r.df() > r.maximumFactors {
		if r.execute() == false {
			break
		}
	}
	for _, id := range source {
		fmt.Fprintln(writer, strconv.FormatInt(r.mapped(id), 10))
	}
	return nil
}

func (r *Reducer) add(id int64) (*reducingConcept, error) {
	c := r.data[id]
	if c == nil {
		path, err := r.longestPathToRoot(id)
		if err != nil {
			return nil, err
		}
		c = &reducingConcept{id: id, pathToRoot: path, count: 1, mappedTo: 0}
		r.data[id] = c
	} else {
		c.count = c.count + 1
	}
	return c, nil
}

// execute dimensionality reduction, returning whether it is possible to do more
func (r *Reducer) execute() bool {
	dirty := false
	lowest := r.calculateScores()
	for id, c := range r.data {
		if c.count > 0 {
			if lowest == c.score {
				if (len(c.pathToRoot) - c.mappedTo) >= r.minimumDistance {
					r.genericise(id)
					dirty = true
				}
			}
		}
	}
	return dirty
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
			c.score = (c.count * maxDistance) + (maxDistance - len(c.pathToRoot) + c.mappedTo)
			if lowestScore == -1 || lowestScore > c.score {
				lowestScore = c.score
			}
		}
	}
	return lowestScore
}

// calculate the maximum distance from root for all data, not including those already mapped
func (r *Reducer) maxDistance() int {
	maxDistance := 0
	for _, c := range r.data {
		if c.count > 0 {
			d := len(c.pathToRoot) - c.mappedTo
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
	c.mappedTo++
	targetID := c.pathToRoot[c.mappedTo]
	target := r.data[targetID]
	if target == nil {
		var err error
		if target, err = r.add(targetID); err != nil {
			return err
		}
	}
	target.count += c.count
	c.count = 0
	return nil
}

func (r *Reducer) longestPathToRoot(conceptID int64) ([]int64, error) {
	c, err := r.svc.GetConcept(conceptID)
	if err != nil {
		return nil, err
	}
	path, err := r.svc.LongestPathToRoot(c)
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
