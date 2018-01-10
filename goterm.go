// SNOMED-CT command line utility and terminology server
//
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
package main

import (
	"github.com/wardle/go-terminology/db"
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
var index = flag.String("index", "", "filename of index to open or create (e.g. ./snomed.index). ")
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
	readOnly := true
	if *doImport != "" || *precompute || *reset {
		readOnly = false
	}
	bolt, err := db.NewBoltService(*database, readOnly)
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
		db.ClearPrecomputations(bolt)
	}
	// perform import if an import root is specified
	if *doImport != "" {
		db.PerformImport(bolt, *doImport)
	}

	if *precompute {
		db.PerformPrecomputations(bolt)
	}

	if *server {
		runServer(bolt)
	}
}

// run our terminology server
// TODO:check precomputations have been run
func runServer(bolt *db.BoltService) {
	// TODO(mw): implement
}
