package mcqs

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

// define simple connection settings statically for this hack. TODO: move to configuration file
const (
	dbDriver   = "postgres"
	dbUser     = "mark"
	dbPassword = ""
	dbName     = "rsdb"
)

// GenerateSnomedCT generates a set of useful intermediate files containing pre-computed and cached
// SNOMED-CT data, useful for sharing with others without depending on a dedicated terminology server.
func GenerateSnomedCT(path string) (SnomedDataset, error) {
	spinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	fmt.Print("Building pre-computed optimised SNOMED-CT data files... ")
	spinner.Start()
	defer spinner.Stop()
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbUser, dbName)
	db, err := sql.Open(dbDriver, dbinfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	defer fmt.Print("Done.\n")
	return NewSnomedDataset(db, path)
}
