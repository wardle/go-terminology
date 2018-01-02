// simple proof-of-concept SNOMED code
package main

import (
	"bitbucket.org/wardle/go-snomed/db"
	"bitbucket.org/wardle/go-snomed/rf2"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

var doImport = flag.String("import", "", "import SNOMED-CT data files from directory specified")
var precompute = flag.Bool("precompute", false, "perform precomputations and optimisations")
var reset = flag.Bool("reset", false, "clear precomputations and optimisations")
var database = flag.String("db", "", "filename of database to open or create (e.g. ./snomed.db)")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file specified")
var server = flag.Bool("server", false, "Run terminology server")

func main() {
	flag.Parse()
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *database == "" {
		fmt.Fprint(os.Stderr, "error: missing mandatory database file\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	bolt, err := db.NewBoltService(*database)
	if err != nil {
		log.Fatal("Couldn't open database")
	}
	defer bolt.Close()
	// turn on CPU profiling if a profile file is specified
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// useful for user to be able to clear precomputations in case of wishing to share
	// a data file with another; the recipient can easily re-run precomputations
	if *reset {
		clearPrecomputations(bolt)
	}
	// perform import if an import root is specified
	if *doImport != "" {
		performImport(bolt, *doImport)
	}

	if *precompute {
		performPrecomputations(bolt)
	}

	if *server {
		runServer(bolt)
	}
}

// performs import from the root specified
// this automatically clears the precomputations, if they exist, but does
// not run precomputations at the end as the user may run multiple individual imports
// from multiple SNOMED-CT distributions before finally running precomputations
// at the end of multiple imports.
func performImport(bolt *db.BoltService, root string) {
	logger := log.New(os.Stdout, "logger: ", log.Lshortfile) // for future use
	importer := rf2.NewImporter(logger)
	concepts, descriptions, relationships := 0, 0, 0
	importer.SetConceptHandler(func(c []*rf2.Concept) {
		concepts = concepts + len(c)
		bolt.PutConcepts(c...)
	})
	importer.SetDescriptionHandler(func(d []*rf2.Description) {
		descriptions += len(d)
		bolt.PutDescriptions(d...)
	})
	importer.SetRelationshipHandler(func(r []*rf2.Relationship) {
		relationships += len(r)
		bolt.PutRelationships(r...)
	})
	clearPrecomputations(bolt)
	err := importer.ImportFiles(*doImport)
	if err != nil {
		log.Fatalf("Could not import files: %v", err)
	}
	fmt.Printf("Imported %d concepts, %d descriptions and %d relationships\n", concepts, descriptions, relationships)
}

// clear precomputations
func clearPrecomputations(bolt *db.BoltService) {
	// TODO(mw):implement
}

// perform precomputations
func performPrecomputations(bolt *db.BoltService) {
	// TODO(mw):implement
}

// run our terminology server
// TODO:check precomputations have been run
func runServer(bolt *db.BoltService) {
	// TODO(mw): implement
}
