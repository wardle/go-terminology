// simple proof-of-concept SNOMED code
package main

import (
	"bitbucket.org/wardle/go-snomed/db"
	"bitbucket.org/wardle/go-snomed/rf2"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	bolt, err := db.NewBoltService("snomed.db")
	if err != nil {
		log.Fatal("Couldn't open database")
	}
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	importer := rf2.NewImporter(logger)
	concepts, descriptions, relationships := 0, 0, 0
	importer.SetConceptHandler(func(c []*rf2.Concept) {
		concepts = concepts + len(c)
		bolt.PutConcepts(c...)
	})
	importer.SetDescriptionHandler(func(d []*rf2.Description) {
		descriptions++
	})
	importer.SetRelationshipHandler(func(r []*rf2.Relationship) {
		relationships++
	})
	err = importer.ImportFiles("/Users/mark/Downloads/uk_sct2cl_24.0.0_20171001000001")
	if err != nil {
		log.Fatalf("Could not import files: %v", err)
	}
	fmt.Printf("Imported %d concepts, %d descriptions and %d relationships!", concepts, descriptions, relationships)
}
