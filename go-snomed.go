// simple proof-of-concept SNOMED code
package main

import (
	"bitbucket.org/wardle/go-snomed/snomed"
	"database/sql"
	"flag"
	"fmt"
	"os"

	"bitbucket.org/wardle/go-snomed/mcqs"
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

// A simple proof-of-concept application to generate fake exam questions
func main() {
	var (
		number     int
		truth      bool
		prevalence bool
		diagnostic bool
		all        bool
	)
	flag.IntVar(&number, "n", 0, "Number to generate. Default: all")
	flag.BoolVar(&truth, "truth", false, "Using precomputed SNOMED-CT, generate a fake truth dataset linking diagnostic concepts with clinical features.")
	flag.BoolVar(&prevalence, "prevalence", false, "Using fake prevalence figures, generate fake questions simply to model prevalence.")
	flag.BoolVar(&diagnostic, "diagnostic", false, "Using fake truth dataset, generate fake questions for machine learning proof-of-concept.")
	flag.BoolVar(&all, "all", false, "Build truth, prevalence and diagnostic data.")
	flag.Parse()
	if all {
		truth = true
		prevalence = true
		diagnostic = true
	}
	if truth || prevalence || diagnostic {
		dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
		db, err := sql.Open(dbDriver, dbinfo)
		if err != nil {
			fmt.Fprint(os.Stderr, "Error: could not open database connection.\n")
			os.Exit(1)
		}
		sct := snomed.NewDatabaseService(db)
		if truth {
			mcqs.GenerateFakeTruth(sct)
		}
		if prevalence {
			mcqs.GeneratePrevalence(sct, number)
		}
		if diagnostic {
			mcqs.GenerateDiagnostic(sct)
		}
	} else {
		flag.PrintDefaults()
	}
}
