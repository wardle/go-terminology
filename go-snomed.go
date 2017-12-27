// simple proof-of-concept SNOMED code
package main

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"database/sql"
	"fmt"
	"golang.org/x/text/language"
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
	_ = &snomed.Snomed{Service: snomed.NewDatabaseService(db), Language: language.BritishEnglish}

}
