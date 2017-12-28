// simple proof-of-concept SNOMED code
package main

import (
	"bitbucket.org/wardle/go-snomed/database"
	"bitbucket.org/wardle/go-snomed/rf2"
	"bytes"
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

	var (
		buf    bytes.Buffer
		logger = log.New(&buf, "logger: ", log.Lshortfile)
	)
	rf2.ImportFiles("/Users/mark/Downloads/SnomedCT_InternationalRF2_PRODUCTION_20170731T150000Z", logger)
}
