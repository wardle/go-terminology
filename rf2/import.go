package rf2

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"bufio"
	"fmt"
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
// Clients do not have to register a handler for all filetypes, but only those
// depending on need.
//
// In this simple initial version, we don't handle batching or try concurrency.
type Importer struct {
	logger                        *log.Logger
	conceptHandler                func(*Concept)
	descriptionHandler            func(*Description)
	relationshipHandler           func(*Relationship)
	refsetDescriptorRefsetHandler func(*RefSetDescriptorReferenceSet)
	languageRefsetHandler         func(*LanguageReferenceSet)
	simpleRefsetHandler           func(*LanguageReferenceSet)
	simpleMapRefsetHandler        func(*SimpleMapReferenceSet)
	complexMapRefsetHandler       func(*ComplexMapReferenceSet)
}

// NewImporter creates a new importer on which you can register handlers
// to process different types of SNOMED-CT RF2 structure.
func NewImporter(logger *log.Logger) *Importer {
	return &Importer{logger: logger}
}

// SetConceptHandler defines a callback for handling concepts
func (im *Importer) SetConceptHandler(f func(*Concept)) {
	im.conceptHandler = f
}

// SetDescriptionHandler defines a callback for handling descriptions
func (im *Importer) SetDescriptionHandler(f func(*Description)) {
	im.descriptionHandler = f
}

// SetRelationshipHandler defines a callback for handling relationships
func (im *Importer) SetRelationshipHandler(f func(*Relationship)) {
	im.relationshipHandler = f
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
)

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
}

// Filename patterns for the supported file types
var fileTypeFilenamePatterns = [...]string{
	"sct2_Concept_Snapshot_\\S+_\\S+.txt",
	"sct2_Description_Snapshot-en_\\S+_\\S+.txt",
	"sct2_(Stated)*Relationship_Snapshot_\\S+_\\S+.txt",
	"der2_cciRefset_RefsetDescriptorSnapshot_\\S+_\\S+.txt",
	"der2_cRefset_LanguageSnapshot-\\S+_\\S+.txt",
	"der2_Refset_SimpleSnapshot_\\S+_\\S+.txt",
	"der2_sRefset_SimpleMapSnapshot_\\S+_\\S+.txt",
	"der2_iisssccRefset_ExtendedMapSnapshot_\\S+_\\S+.txt",
}

// return the filename pattern for this file type
func (ft fileType) filenamePattern() string {
	return fileTypeFilenamePatterns[ft]
}
func (ft fileType) columnNames() []string {
	return columnNames[ft]
}

func (ft fileType) String() string {
	return fileTypeNames[ft]
}

// calculateFileType determines the type of file from its filename, returning a
// boolean to indicate whether a file type was successfully determined.
func calculateFileType(path string) (fileType, bool) {
	filename := filepath.Base(path)
	for i, p := range fileTypeFilenamePatterns {
		matched, _ := regexp.MatchString(p, filename)
		if matched {
			return fileType(i), true
		}
	}
	return -1, false
}

// ImportFiles imports all SNOMED-CT files from a SNOMED-CT distribution
// See https://www.nhs-data.uk/Docs/SNOMEDCTFileSpec.pdf
// We must walk the directory tree and identify all of the different file types.
// We must then process those in turn, ensuring that concepts are imported before
// descriptions and relationships.
func (im Importer) ImportFiles(root string) error {
	files := make(map[fileType][]string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		ft, success := calculateFileType(path)
		if success {
			files[ft] = append(files[ft], path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, f := range files[conceptsFileType] {
		err = im.processConceptFile(f)
		if err != nil {
			return err
		}
	}
	for _, f := range files[descriptionsFileType] {
		err = im.processDescriptionFile(f)
		if err != nil {
			return err
		}
	}
	for _, f := range files[relationshipsFileType] {
		err = im.processRelationshipFile(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseIdentifier(s string, errs *[]error) snomed.Identifier {
	return snomed.Identifier(parseInt(s, errs))
}

func parseInt(s string, errs *[]error) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		*errs = append(*errs, err)
	}
	return i
}
func parseBoolean(s string, errs *[]error) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		*errs = append(*errs, err)
	}
	return b
}
func parseDate(s string, errs *[]error) time.Time {
	t, err := time.Parse("20060102", s)
	if err != nil {
		*errs = append(*errs, err)
	}
	return t
}

func (im Importer) processConceptFile(filename string) error {
	if im.conceptHandler == nil {
		im.logger.Printf("Ignoring concept file %s: no handler", filename)
		return nil
	}
	im.logger.Printf("Processing concept file %s\n", filename)
	return importFile(filename, conceptsFileType.columnNames(), im.logger, func(row []string) {
		var errs []error
		id := parseIdentifier(row[0], &errs)
		effectiveTime := parseDate(row[1], &errs)
		active := parseBoolean(row[2], &errs)
		moduleID := parseIdentifier(row[3], &errs)
		defnID := parseIdentifier(row[4], &errs)
		if len(errs) > 0 {
			im.logger.Printf("failed parsing concept %s : %v", row[0], errs)
		} else {
			concept := &Concept{ID: id, EffectiveTime: effectiveTime, Active: active, ModuleID: moduleID, DefinitionStatusID: defnID}
			im.conceptHandler(concept)
		}
	})
}

// id      effectiveTime   active  moduleId        conceptId       languageCode    typeId  term    caseSignificanceId
func (im Importer) processDescriptionFile(filename string) error {
	if im.descriptionHandler == nil {
		im.logger.Printf("Ignoring description file %s: no handler", filename)
		return nil
	}
	im.logger.Printf("Processing description file %s\n", filename)
	return importFile(filename, descriptionsFileType.columnNames(), im.logger, func(row []string) {
		var errs []error
		id := parseIdentifier(row[0], &errs)
		effectiveTime := parseDate(row[1], &errs)
		active := parseBoolean(row[2], &errs)
		moduleID := parseIdentifier(row[3], &errs)
		conceptID := parseIdentifier(row[4], &errs)
		languageCode := row[5]
		typeID := parseIdentifier(row[6], &errs)
		term := row[7]
		caseSigID := parseIdentifier(row[8], &errs)
		if len(errs) > 0 {
			im.logger.Printf("failed parsing description %s : %v", row[0], errs)
		} else {
			description := &Description{ID: id, EffectiveTime: effectiveTime, Active: active,
				ModuleID: moduleID, ConceptID: conceptID, LanguageCode: languageCode, TypeID: typeID, Term: term, CaseSignificance: caseSigID}
			im.descriptionHandler(description)
		}
	})
}

// id      effectiveTime   active  moduleId        sourceId        destinationId   relationshipGroup       typeId  characteristicTypeId    modifierId
func (im Importer) processRelationshipFile(filename string) error {
	if im.relationshipHandler == nil {
		im.logger.Printf("Ignoring relationship file %s: no handler", filename)
		return nil
	}
	im.logger.Printf("Processing relationship file %s\n", filename)
	return importFile(filename, relationshipsFileType.columnNames(), im.logger, func(row []string) {
		var errs []error
		id := parseIdentifier(row[0], &errs)
		effectiveTime := parseDate(row[1], &errs)
		active := parseBoolean(row[2], &errs)
		moduleID := parseIdentifier(row[3], &errs)
		sourceID := parseIdentifier(row[4], &errs)
		destinationID := parseIdentifier(row[5], &errs)
		relGroup := parseInt(row[6], &errs)
		typeID := parseIdentifier(row[7], &errs)
		charTypeID := parseIdentifier(row[8], &errs)
		modifierID := parseIdentifier(row[9], &errs)
		if len(errs) > 0 {
			im.logger.Printf("failed parsing relationship %s : %v", row[0], errs)
		} else {
			relationship := &Relationship{ID: id, EffectiveTime: effectiveTime, Active: active,
				ModuleID: moduleID, SourceID: sourceID, DestinationID: destinationID, RelationshipGroup: relGroup, TypeID: typeID, CharacteristicTypeID: charTypeID, ModifierID: modifierID}
			im.relationshipHandler(relationship)
		}
	})
}

// importFile reads a tab-delimited file and calls a handler for each row
func importFile(filename string, columnNames []string, logger *log.Logger, processFunc func(row []string)) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	nCols := len(columnNames)
	// read the first line and check that we have the right column names
	scanner.Scan()
	if err != nil {
		return err
	}
	headings := strings.Split(scanner.Text(), "\t")
	if !reflect.DeepEqual(headings, columnNames) {
		return fmt.Errorf("expecting column names: %v, got: %v", columnNames, headings)
	}
	// process each line
	for scanner.Scan() {
		record := strings.Split(scanner.Text(), "\t")
		l := len(record)
		if l > 0 {
			if l != nCols {
				logger.Fatalf("incorrect number of columns %v\n", record)
			} else {
				processFunc(record)
			}
		}
	}
	return nil
}