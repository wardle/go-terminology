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
	"bufio"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	proto "github.com/golang/protobuf/ptypes/timestamp"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Importer manages the handling of different types of SNOMED-CT data structure
//
type Importer struct {
	logger    *log.Logger
	batchSize int
	handler   func(interface{})
}

// NewImporter creates a new importer on which you can register a handler
// to process different types of SNOMED-CT RF2 structure.
func NewImporter(logger *log.Logger, handler func(interface{})) *Importer {
	return &Importer{logger: logger, batchSize: 5000, handler: handler}
}

// fileType represents a type of SNOMED-CT distribution file
type fileType int

// Supported file types
// These are listed in order of importance for import
const (
	conceptsFileType fileType = iota
	descriptionsFileType
	relationshipsFileType
	refsetDescriptorRefsetFileType
	languageRefsetFileType
	simpleRefsetFileType
	simpleMapRefsetFileType
	complexMapRefsetFileType
	lastFileType
)

type task struct {
	filename  string
	batchSize int
	fileType  fileType
}

var fileTypeNames = [...]string{
	"Concepts",
	"Descriptions",
	"Relationships",
	"Refset Descriptor refset",
	"Language refset",
	"Simple refset",
	"Simple map refset",
	"Complex / extended map refset",
}
var columnNames = [...][]string{
	[]string{"id", "effectiveTime", "active", "moduleId", "definitionStatusId"},
	[]string{"id", "effectiveTime", "active", "moduleId", "conceptId", "languageCode", "typeId", "term", "caseSignificanceId"},
	[]string{"id", "effectiveTime", "active", "moduleId", "sourceId", "destinationId", "relationshipGroup", "typeId", "characteristicTypeId", "modifierId"},
	[]string{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "attributeDescription", "attributeType", "attributeOrder"},
	[]string{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "acceptabilityId"},
	[]string{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId"},
	[]string{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "mapTarget"},
	[]string{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "mapGroup", "mapPriority", "mapRule", "mapAdvice", "mapTarget", "correlationId", "mapBlock"},
}

// Filename patterns for the supported file types
var fileTypeFilenamePatterns = [...]string{
	"sct2_Concept_Snapshot_\\S+_\\S+.txt",
	"sct2_Description_Snapshot-en\\S+_\\S+.txt",
	"sct2_(Stated)*Relationship_Snapshot_\\S+_\\S+.txt",
	"der2_cciRefset_RefsetDescriptorSnapshot_\\S+_\\S+.txt",
	"der2_cRefset_LanguageSnapshot-\\S+_\\S+.txt",
	"der2_Refset_SimpleSnapshot_\\S+_\\S+.txt",
	"der2_sRefset_SimpleMapSnapshot_\\S+_\\S+.txt",
	"der2_iisssciRefset_ExtendedMapSnapshot_\\S+_\\S+.txt",
}

// Processors for each file type
var processors = [...]func(im *Importer, task *task) error{
	processConceptFile,
	processDescriptionFile,
	processRelationshipFile,
	nil,
	processLanguageRefsetFile,
	processSimpleRefsetFile,
	processSimpleMapRefsetFile,
	processComplexMapRefsetFile,
}

// return the filename pattern for this file type
func (ft fileType) pattern() string {
	return fileTypeFilenamePatterns[ft]
}

// column names for this file type
func (ft fileType) cols() []string {
	return columnNames[ft]
}
func (ft fileType) processor() func(im *Importer, task *task) error {
	return processors[ft]
}

func (ft fileType) String() string {
	return fileTypeNames[ft]
}

// calculateFileType determines the type of file from its filename, returning a
// boolean to indicate whether a file type was successfully determined.
func calculateFileType(path string) (fileType, bool) {
	filename := filepath.Base(path)
	for ft := conceptsFileType; ft < lastFileType; ft++ {
		matched, _ := regexp.MatchString(ft.pattern(), filename)
		if matched {
			return ft, true
		}
	}
	return -1, false
}

// ImportFiles imports all SNOMED-CT files from a SNOMED-CT distribution
// See https://www.nhs-data.uk/Docs/SNOMEDCTFileSpec.pdf
// We must walk the directory tree and identify all of the different file types.
// We must then process those in turn, ensuring that concepts are imported before
// descriptions and relationships.
func (im *Importer) ImportFiles(root string) error {
	tasks := make(map[int][]*task)
	maxRank := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if ft, success := calculateFileType(path); success {
			task := &task{filename: path, batchSize: im.batchSize, fileType: ft}
			rank := int(ft)
			tasks[rank] = append(tasks[rank], task)
			if rank > maxRank {
				maxRank = rank
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		return fmt.Errorf("error: found 0 datafiles at path '%s'", root)
	}
	// execute each task, but in rank order so that concepts come before descriptions and relationships
	for rank := 0; rank <= maxRank; rank++ {
		rankedTasks := tasks[rank]
		for _, task := range rankedTasks {
			if task.fileType.processor() != nil {
				if err = task.fileType.processor()(im, task); err != nil {
					return nil
				}

			}
		}
	}
	return nil
}

func parseIdentifier(s string, errs *[]error) int64 {
	return parseInt(s, errs)
}

func parseInt(s string, errs *[]error) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		*errs = append(*errs, err)
	}
	return int64(i)
}
func parseBoolean(s string, errs *[]error) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		*errs = append(*errs, err)
	}
	return b
}
func parseDate(s string, errs *[]error) *proto.Timestamp {
	t, err := time.Parse("20060102", s)
	if err != nil {
		*errs = append(*errs, err)
	}
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		*errs = append(*errs, err)
	}
	return ts
}

func processConceptFile(im *Importer, task *task) error {
	im.logger.Printf("Processing concept file %s\n", task.filename)
	return importFile(task, im.logger, func(rows [][]string) {
		result := make([]*Concept, 0, len(rows))
		for _, row := range rows {
			var errs []error
			concept := parseConcept(row, &errs)
			if len(errs) > 0 {
				im.logger.Printf("failed to parse concept %s : %v", row[0], errs)
			} else {
				result = append(result, concept)
			}
		}
		im.handler(result)
	})
}

func processDescriptionFile(im *Importer, task *task) error {
	im.logger.Printf("Processing description file %s\n", task.filename)
	return importFile(task, im.logger, func(rows [][]string) {
		result := make([]*Description, 0, len(rows))
		for _, row := range rows {
			var errs []error
			description := parseDescription(row, &errs)
			if len(errs) > 0 {
				im.logger.Printf("failed to parse description %s : %v", row[0], errs)
			} else {
				result = append(result, description)
			}
		}
		im.handler(result)
	})
}

func processRelationshipFile(im *Importer, task *task) error {
	im.logger.Printf("Processing relationship file %s\n", task.filename)
	return importFile(task, im.logger, func(rows [][]string) {
		var result = make([]*Relationship, 0, len(rows))
		for _, row := range rows {
			var errs []error
			relationship := parseRelationship(row, &errs)
			if len(errs) > 0 {
				im.logger.Printf("failed to parse relationship %s : %v", row[0], errs)
			} else {
				result = append(result, relationship)
			}
		}
		im.handler(result)
	})
}

// id      effectiveTime   active  moduleId        refsetId        referencedComponentId   acceptabilityId
// bba5806d-8d8e-5295-ac6a-962b67c8ed50    20040131        1       999000011000000103      900000000000508004      999002221000000116      900000000000548007
func processLanguageRefsetFile(im *Importer, task *task) error {
	im.logger.Printf("Processing language refset file %s\n", task.filename)
	return importFile(task, im.logger, func(rows [][]string) {
		var result = make([]*ReferenceSetItem, 0, len(rows))
		for _, row := range rows {
			var errs []error
			item := parseReferenceSetHeader(row, &errs)
			item.Body = &ReferenceSetItem_Language{
				Language: &LanguageReferenceSet{
					AcceptabilityId: parseInt(row[6], &errs),
				},
			}
			if len(errs) > 0 {
				im.logger.Printf("failed to parse language refset %s : %v", row[0], errs)
			} else {
				result = append(result, item)
			}
		}
		im.handler(result)
	})
}

// id      effectiveTime   active  moduleId        refsetId        referencedComponentId
func processSimpleRefsetFile(im *Importer, task *task) error {
	im.logger.Printf("Processing simple refset file %s\n", task.filename)
	return importFile(task, im.logger, func(rows [][]string) {
		var result = make([]*ReferenceSetItem, 0, len(rows))
		for _, row := range rows {
			var errs []error
			item := parseReferenceSetHeader(row, &errs)
			item.Body = &ReferenceSetItem_Simple{Simple: &SimpleReferenceSet{}}
			if len(errs) > 0 {
				im.logger.Printf("failed to parse simple refset %s : %v", row[0], errs)
			} else {
				result = append(result, item)
			}
		}
		im.handler(result)
	})
}

func processSimpleMapRefsetFile(im *Importer, task *task) error {
	im.logger.Printf("Processing simple map refset file %s\n", task.filename)
	return importFile(task, im.logger, func(rows [][]string) {
		var result = make([]*ReferenceSetItem, 0, len(rows))
		for _, row := range rows {
			var errs []error
			item := parseReferenceSetHeader(row, &errs)
			item.Body = &ReferenceSetItem_SimpleMap{
				SimpleMap: &SimpleMapReferenceSet{
					MapTarget: row[6],
				},
			}
			if len(errs) > 0 {
				im.logger.Printf("failed to parse simple map refset %s : %v", row[0], errs)
			} else {
				result = append(result, item)
			}
		}
		im.handler(result)
	})
}

func processComplexMapRefsetFile(im *Importer, task *task) error {
	im.logger.Printf("Processing complex map refset file %s\n", task.filename)
	return importFile(task, im.logger, func(rows [][]string) {
		var result = make([]*ReferenceSetItem, 0, len(rows))
		for _, row := range rows {
			var errs []error
			item := parseReferenceSetHeader(row, &errs)
			item.Body = &ReferenceSetItem_ComplexMap{
				ComplexMap: &ComplexMapReferenceSet{
					MapGroup:    parseInt(row[6], &errs),
					MapPriority: parseInt(row[7], &errs),
					MapRule:     row[8],
					MapAdvice:   row[9],
					MapTarget:   row[10],
					Correlation: parseInt(row[11], &errs),
					MapCategory: parseInt(row[12], &errs),
				},
			}
			if len(errs) > 0 {
				im.logger.Printf("failed to parse complex map refset %s : %v", row[0], errs)
			} else {
				result = append(result, item)
			}
		}
		im.handler(result)
	})
}

func parseConcept(row []string, errs *[]error) *Concept {
	return &Concept{
		Id:                 parseIdentifier(row[0], errs),
		EffectiveTime:      parseDate(row[1], errs),
		Active:             parseBoolean(row[2], errs),
		ModuleId:           parseIdentifier(row[3], errs),
		DefinitionStatusId: parseIdentifier(row[4], errs)}
}

// id      effectiveTime   active  moduleId        conceptId       languageCode    typeId  term    caseSignificanceId
func parseDescription(row []string, errs *[]error) *Description {
	return &Description{
		Id:               parseIdentifier(row[0], errs),
		EffectiveTime:    parseDate(row[1], errs),
		Active:           parseBoolean(row[2], errs),
		ModuleId:         parseIdentifier(row[3], errs),
		ConceptId:        parseIdentifier(row[4], errs),
		LanguageCode:     row[5],
		TypeId:           parseIdentifier(row[6], errs),
		Term:             row[7],
		CaseSignificance: parseIdentifier(row[8], errs)}
}

// id      effectiveTime   active  moduleId        sourceId        destinationId   relationshipGroup       typeId  characteristicTypeId    modifierId
func parseRelationship(row []string, errs *[]error) *Relationship {
	return &Relationship{
		Id:                   parseIdentifier(row[0], errs),
		EffectiveTime:        parseDate(row[1], errs),
		Active:               parseBoolean(row[2], errs),
		ModuleId:             parseIdentifier(row[3], errs),
		SourceId:             parseIdentifier(row[4], errs),
		DestinationId:        parseIdentifier(row[5], errs),
		RelationshipGroup:    parseInt(row[6], errs),
		TypeId:               parseIdentifier(row[7], errs),
		CharacteristicTypeId: parseIdentifier(row[8], errs),
		ModifierId:           parseIdentifier(row[9], errs)}

}

// parse a reference set from the row
func parseReferenceSetHeader(row []string, errs *[]error) *ReferenceSetItem {
	return &ReferenceSetItem{
		Id:                    row[0], // identifier is a long unique uuid string,
		EffectiveTime:         parseDate(row[1], errs),
		Active:                parseBoolean(row[2], errs),
		ModuleId:              parseIdentifier(row[3], errs),
		RefsetId:              parseIdentifier(row[4], errs),
		ReferencedComponentId: parseIdentifier(row[5], errs),
	}
}

// importFile reads a tab-delimited file and calls a handler for a batch of rows
func importFile(task *task, logger *log.Logger, processFunc func(rows [][]string)) error {
	f, err := os.Open(task.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	// read the first line and check that we have the right column names
	if scanner.Scan() == false {
		return fmt.Errorf("empty file %s", task.filename)
	}
	headings := strings.Split(scanner.Text(), "\t")
	if !reflect.DeepEqual(headings, task.fileType.cols()) {
		return fmt.Errorf("expecting column names: %v, got: %v", task.fileType.cols(), headings)
	}
	batch := make([][]string, 0, task.batchSize)
	for scanner.Scan() {
		record := strings.Split(scanner.Text(), "\t")
		batch = append(batch, record)
		if len(batch) == task.batchSize {
			processFunc(batch)
			batch = nil
		}
	}
	if len(batch) > 0 {
		processFunc(batch)
	}
	return nil
}
