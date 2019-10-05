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
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/golang/protobuf/ptypes/timestamp"
	context "golang.org/x/net/context"
)

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
	extendedMapRefsetFileType
	complexMapRefsetFileType
	attributeValueRefsetFileType
	associationRefsetFileType
	lastFileType
)

var fileTypeNames = [...]string{
	"Concepts",
	"Descriptions",
	"Relationships",
	"Refset Descriptor refset",
	"Language refset",
	"Simple refset",
	"Simple map refset",
	"Extended map refset",
	"Complex map refset",
	"Attribute value refset",
	"Association refset",
}
var columnNames = [...][]string{
	{"id", "effectiveTime", "active", "moduleId", "definitionStatusId"},
	{"id", "effectiveTime", "active", "moduleId", "conceptId", "languageCode", "typeId", "term", "caseSignificanceId"},
	{"id", "effectiveTime", "active", "moduleId", "sourceId", "destinationId", "relationshipGroup", "typeId", "characteristicTypeId", "modifierId"},
	{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "attributeDescription", "attributeType", "attributeOrder"},
	{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "acceptabilityId"},
	{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId"},
	{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "mapTarget"},
	{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "mapGroup", "mapPriority", "mapRule", "mapAdvice", "mapTarget", "correlationId", "mapCategoryId"},
	{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "mapGroup", "mapPriority", "mapRule", "mapAdvice", "mapTarget", "correlationId", "mapBlock"},
	{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "valueId"},
	{"id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "targetComponentId"},
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
	"der2_iisssccRefset_ExtendedMapSnapshot_\\S+_\\S+.txt", // extended
	"der2_iisssciRefset_ExtendedMapSnapshot_\\S+_\\S+.txt", // complex
	"der2_cRefset_AttributeValueSnapshot_\\S+_\\S+.txt",
	"der2_cRefset_AssociationSnapshot_\\S+_\\S+.txt",
}

// return the filename pattern for this file type
func (ft fileType) pattern() string {
	return fileTypeFilenamePatterns[ft]
}

// column names for this file type
func (ft fileType) cols() []string {
	return columnNames[ft]
}

func (ft fileType) String() string {
	return fileTypeNames[ft]
}

type task struct {
	filename  string
	batchSize int
	fileType  fileType
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

func parseReferenceSetItem(ft fileType, row []string, err *[]error) *ReferenceSetItem {
	switch ft {
	case refsetDescriptorRefsetFileType:
		return parseRefsetDescriptorRefset(row, err)
	case languageRefsetFileType:
		return parseLanguageRefset(row, err)
	case simpleRefsetFileType:
		return parseSimpleRefset(row, err)
	case simpleMapRefsetFileType:
		return parseSimpleMapRefset(row, err)
	case extendedMapRefsetFileType:
		return parseExtendedMapRefset(row, err)
	case complexMapRefsetFileType:
		return parseComplexMapRefset(row, err)
	case attributeValueRefsetFileType:
		return parseAttributeValueRefset(row, err)
	case associationRefsetFileType:
		return parseAssociationRefset(row, err)
	}
	*err = append(*err, fmt.Errorf("error: unable to process filetype %s", ft))
	return nil
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

// "id", "effectiveTime", "active", "moduleId", "refsetId", "referencedComponentId", "attributeDescription", "attributeType", "attributeOrder"},
func parseRefsetDescriptorRefset(row []string, errs *[]error) *ReferenceSetItem {
	item := parseReferenceSetHeader(row, errs)
	item.Body = &ReferenceSetItem_RefsetDescriptor{
		RefsetDescriptor: &RefSetDescriptorReferenceSet{
			AttributeDescriptionId: parseInt(row[6], errs),
			AttributeTypeId:        parseInt(row[7], errs),
			AttributeOrder:         uint32(parseInt(row[8], errs)),
		},
	}
	return item
}

func parseLanguageRefset(row []string, errs *[]error) *ReferenceSetItem {
	item := parseReferenceSetHeader(row, errs)
	item.Body = &ReferenceSetItem_Language{
		Language: &LanguageReferenceSet{
			AcceptabilityId: parseInt(row[6], errs),
		},
	}
	return item
}

func parseSimpleRefset(row []string, errs *[]error) *ReferenceSetItem {
	item := parseReferenceSetHeader(row, errs)
	item.Body = &ReferenceSetItem_Simple{
		Simple: &SimpleReferenceSet{},
	}
	return item
}

func parseSimpleMapRefset(row []string, errs *[]error) *ReferenceSetItem {
	item := parseReferenceSetHeader(row, errs)
	item.Body = &ReferenceSetItem_SimpleMap{
		SimpleMap: &SimpleMapReferenceSet{
			MapTarget: row[6],
		},
	}
	return item
}

func parseExtendedMapRefset(row []string, errs *[]error) *ReferenceSetItem {
	item := parseReferenceSetHeader(row, errs)
	item.Body = &ReferenceSetItem_ComplexMap{
		ComplexMap: &ComplexMapReferenceSet{
			MapGroup:    parseInt(row[6], errs),
			MapPriority: parseInt(row[7], errs),
			MapRule:     row[8],
			MapAdvice:   row[9],
			MapTarget:   strings.TrimSpace(row[10]),
			Correlation: parseInt(row[11], errs),
			MapCategory: parseInt(row[12], errs),
		},
	}
	return item
}
func parseComplexMapRefset(row []string, errs *[]error) *ReferenceSetItem {
	item := parseReferenceSetHeader(row, errs)
	item.Body = &ReferenceSetItem_ComplexMap{
		ComplexMap: &ComplexMapReferenceSet{
			MapGroup:    parseInt(row[6], errs),
			MapPriority: parseInt(row[7], errs),
			MapRule:     row[8],
			MapAdvice:   row[9],
			MapTarget:   strings.TrimSpace(row[10]),
			Correlation: parseInt(row[11], errs),
			MapBlock:    parseInt(row[12], errs),
		},
	}
	return item
}
func parseAttributeValueRefset(row []string, errs *[]error) *ReferenceSetItem {
	item := parseReferenceSetHeader(row, errs)
	item.Body = &ReferenceSetItem_AttributeValue{
		AttributeValue: &AttributeValueReferenceSet{
			ValueId: parseInt(row[6], errs),
		},
	}
	return item
}
func parseAssociationRefset(row []string, errs *[]error) *ReferenceSetItem {
	item := parseReferenceSetHeader(row, errs)
	item.Body = &ReferenceSetItem_Association{
		Association: &AssociationReferenceSet{
			TargetComponentId: parseInt(row[6], errs),
		},
	}
	return item
}

// ImportChannels defines the channels through which batches of data will be returned
type ImportChannels struct {
	Concepts      chan []*Concept
	Descriptions  chan []*Description
	Relationships chan []*Relationship
	Refsets       chan []*ReferenceSetItem
}

// Close all results channels
func (ir *ImportChannels) Close() {
	close(ir.Concepts)
	close(ir.Descriptions)
	close(ir.Relationships)
	close(ir.Refsets)
}

// FastImport imports all SNOMED datafiles from the specified root, returning data in batches through
// the returned channels.
func FastImport(ctx context.Context, root string, batchSize int) *ImportChannels {
	result := new(ImportChannels)
	result.Concepts = make(chan []*Concept)
	result.Descriptions = make(chan []*Description)
	result.Relationships = make(chan []*Relationship)
	result.Refsets = make(chan []*ReferenceSetItem)

	taskc := walkFiles(ctx, root, batchSize)

	// processFiles: takes tasks from walkfiles and turn into rows
	batchc := make(chan batch) // channel to handle batches of rows for processing
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			processFiles(ctx, taskc, batchc)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(batchc)
	}()

	// process batches: start some workers to process batches of data for import
	var batchWg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		batchWg.Add(1)
		go func() {
			processBatch(ctx, batchc, result)
			batchWg.Done()
		}()
	}
	go func() { // when all processBatch operations finish, close our output channels
		batchWg.Wait()
		result.Close()
	}()

	return result
}

type batch struct {
	task
	rows [][]string
}

// walkFiles walks the directory tree from the root specified and identifies
// SNOMED CT files and their type, emitting tasks on the created channel
func walkFiles(ctx context.Context, root string, batchSize int) <-chan task {
	tasks := make(chan task)
	go func() {
		defer close(tasks)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error processing %s : %s", path, err)
			}
			ft, success := calculateFileType(path)
			if !success {
				return nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case tasks <- task{filename: path, batchSize: batchSize, fileType: ft}: // when output channel free, send it a task
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}()
	return tasks
}

// processFiles will drain the tasks channel and then return, sending out batches of work to the batch channel
func processFiles(ctx context.Context, tasks <-chan task, batchc chan<- batch) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-tasks:
			if task.filename == "" {
				return
			}
			f, err := os.Open(task.filename)
			if err != nil {
				panic(fmt.Sprintf("unable to process file %s: %s", task.filename, err))
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			// read the first line and check that we have the right column names
			if scanner.Scan() == false {
				panic(fmt.Errorf("empty file %s", task.filename))
			}
			headings := strings.Split(scanner.Text(), "\t")
			if !reflect.DeepEqual(headings, task.fileType.cols()) {
				panic(fmt.Errorf("expecting column names: %v, got: %v", task.fileType.cols(), headings))
			}
			batch := batch{
				task: task,
				rows: make([][]string, 0, task.batchSize),
			}
			count := 0
			for scanner.Scan() {
				count++
				row := strings.Split(scanner.Text(), "\t")
				batch.rows = append(batch.rows, row)
				if count == task.batchSize {
					select {
					case batchc <- batch:
					case <-ctx.Done():
						return
					}
					batch.rows = nil
					count = 0
				}
			}
			select {
			case batchc <- batch:
			case <-ctx.Done():
				return
			}
		}
	}
}

func processBatch(ctx context.Context, batchc <-chan batch, results *ImportChannels) {
	for batch := range batchc {
		switch batch.fileType {
		case conceptsFileType:
			processConcepts(ctx, batch, results.Concepts)
		case descriptionsFileType:
			processDescriptions(ctx, batch, results.Descriptions)
		case relationshipsFileType:
			processRelationships(ctx, batch, results.Relationships)
		case refsetDescriptorRefsetFileType,
			languageRefsetFileType,
			simpleRefsetFileType,
			simpleMapRefsetFileType,
			complexMapRefsetFileType,
			extendedMapRefsetFileType,
			attributeValueRefsetFileType,
			associationRefsetFileType:
			processReferenceSetItems(ctx, batch, results.Refsets)
		default:
			panic(fmt.Errorf("unsupported file type: %s", batch.fileType))
		}
	}
}

func processConcepts(ctx context.Context, batch batch, concepts chan<- []*Concept) {
	result := make([]*Concept, 0, len(batch.rows))
	for _, row := range batch.rows {
		var errs []error
		c := parseConcept(row, &errs)
		if len(errs) > 0 {
			panic(fmt.Errorf("failed to parse concept %s : %v", row[0], errs))
		}
		result = append(result, c)
	}
	select {
	case concepts <- result:
	case <-ctx.Done():
		return
	}
}
func processDescriptions(ctx context.Context, batch batch, descriptions chan<- []*Description) {
	result := make([]*Description, 0, len(batch.rows))
	for _, row := range batch.rows {
		var errs []error
		c := parseDescription(row, &errs)
		if len(errs) > 0 {
			panic(fmt.Errorf("failed to parse description %s : %v", row[0], errs))
		}
		result = append(result, c)
	}
	select {
	case descriptions <- result:
	case <-ctx.Done():
		return
	}
}
func processRelationships(ctx context.Context, batch batch, outc chan<- []*Relationship) {
	result := make([]*Relationship, 0, len(batch.rows))
	for _, row := range batch.rows {
		var errs []error
		c := parseRelationship(row, &errs)
		if len(errs) > 0 {
			panic(fmt.Errorf("failed to parse description %s : %v", row[0], errs))
		}
		result = append(result, c)
	}
	select {
	case outc <- result:
	case <-ctx.Done():
		return
	}
}
func processReferenceSetItems(ctx context.Context, batch batch, outc chan<- []*ReferenceSetItem) {
	result := make([]*ReferenceSetItem, 0, len(batch.rows))
	for _, row := range batch.rows {
		var errs []error
		o := parseReferenceSetItem(batch.fileType, row, &errs)
		if len(errs) > 0 {
			panic(fmt.Errorf("failed to parse reference set item %s : %v", row[0], errs))
		}
		result = append(result, o)
	}
	select {
	case outc <- result:
	case <-ctx.Done():
		return
	}
}
