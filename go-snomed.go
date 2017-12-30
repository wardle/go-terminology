// simple proof-of-concept SNOMED code
package main

import (
	"bitbucket.org/wardle/go-snomed/database"
	"bitbucket.org/wardle/go-snomed/rf2"
	"database/sql"
	"fmt"
	"golang.org/x/text/language"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// database connection parameters
// TODO: permit configuration at runtime
const (
	dbDriver   = "postgres"
	dbUser     = "mark"
	dbPassword = ""
	dbName     = "rsdb"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open(dbDriver, dbinfo)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: could not open database connection.\n")
		os.Exit(1)
	}
	_ = &database.Snomed{Service: database.NewSQLService(db), Language: language.BritishEnglish}

	logger := log.New(os.Stdout, "logger: ", log.Lshortfile)
	importer := rf2.NewImporter(logger)
	concepts, descriptions, relationships := 0, 0, 0
	importer.SetConceptHandler(func(c *rf2.Concept) {
		concepts++
	})
	importer.SetDescriptionHandler(func(d *rf2.Description) {
		descriptions++
	})
	importer.SetRelationshipHandler(func(r *rf2.Relationship) {
		relationships++
	})
	err = importer.ImportFiles("/Users/mark/Downloads/uk_sct2cl_24.0.0_20171001000001")
	if err != nil {
		log.Fatalf("Could not import files: %v", err)
	}
	fmt.Printf("Imported %d concepts, %d descriptions and %d relationships!", concepts, descriptions, relationships)
}
