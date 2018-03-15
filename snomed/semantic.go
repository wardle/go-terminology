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

// Common concept identifiers, predominantly for the SNOMED-CT metadata model.
const (
	// The SNOMED CT root concept
	RootConceptID Identifier = 138875005

	// IsAConceptID represents the relationship type, IS-A; the commonest type of relationship
	IsAConceptID int64 = 116680003
)
