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

package snomed

import (
	"fmt"
	"strconv"

	"github.com/wardle/go-terminology/verhoeff"
)

// Identifier (SCTID) is a checksummed (Verhoeff) globally unique persistent identifier
// See https://confluence.ihtsdotools.org/display/DOCTIG/3.1.4.2.+Component+features+-+Identifiers
// The SCTID data type is 64-bit Integer which is allocated and represented in accordance with a set of rules.
// These rules enable each Identifier to refer unambiguously to a unique component.
// They also support separate partitions for allocation of Identifiers for particular types of component and
// namespaces that distinguish between different issuing organizations.
//
// A valid identifier can be represented either as a uint64 or an int64. See https://confluence.ihtsdotools.org/display/DOCRELFMT/6.3+SCTID+Constraints
//
type Identifier int64

// ParseIdentifier converts a string into an identifier
func ParseIdentifier(s string) (Identifier, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return Identifier(id), nil
}

// ParseAndValidate converts a string into an identifier and validates
func ParseAndValidate(s string) (Identifier, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	id2 := Identifier(id)
	if !id2.IsValid() {
		return 0, fmt.Errorf("invalid identifier '%s'", s)
	}
	return id2, nil
}

// Integer is a convenience method to convert to integer
func (id Identifier) Integer() int64 {
	return int64(id)
}

// String returns a string representation of this identifier
func (id Identifier) String() string {
	return strconv.FormatInt(int64(id), 10)
}

// IsConcept will return true if this identifier refers to a concept
func (id Identifier) IsConcept() bool {
	pid := id.partitionIdentifier()
	return pid == "00" || pid == "10"
}

// IsDescription will return true if this identifier refers to a description.
func (id Identifier) IsDescription() bool {
	pid := id.partitionIdentifier()
	return pid == "01" || pid == "11"
}

// IsRelationship will return true if this identifier refers to a relationship.
func (id Identifier) IsRelationship() bool {
	pid := id.partitionIdentifier()
	return pid == "02" || pid == "12"
}

// IsValid will return true if this is a valid SNOMED CT identifier
func (id Identifier) IsValid() bool {
	return verhoeff.Validate(int64(id))
}

// partitionIdentifier returns the penultimate last digit digits
// see https://confluence.ihtsdotools.org/display/DOCRELFMT/5.5.+Partition+Identifier
// 0123456789
// xxxxxxxppc
func (id Identifier) partitionIdentifier() string {
	s := strconv.FormatInt(int64(id), 10)
	l := len(s)
	return s[l-3 : l-1]
}
