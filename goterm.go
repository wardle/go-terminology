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
	"bufio"
	"flag"
	"fmt"
	"github.com/wardle/go-terminology/analysis"
	"github.com/wardle/go-terminology/server"
	"github.com/wardle/go-terminology/terminology"
	"log"
	"os"
	"runtime/pprof"
)

// automatically populated by linker flags
var version string
var build string

// commands and flags
var doVersion = flag.Bool("version", false, "Show version information")
var doImport = flag.Bool("import", false, "import SNOMED-CT data files from directories specified")
var precompute = flag.Bool("precompute", false, "perform precomputations and optimisations")
var reset = flag.Bool("reset", false, "clear precomputations and optimisations")
var database = flag.String("db", "", "filename of database to open or create (e.g. ./snomed.db)")
var lang = flag.String("lang", "en-GB", "language tags to be used, default 'en-GB'")

var verbose = flag.Bool("v", false, "verbose")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file specified")
var runserver = flag.Bool("server", false, "run terminology server")
var stats = flag.Bool("status", false, "get statistics")
var port = flag.Int("port", 8080, "port to use when running server")
var export = flag.Bool("export", false, "export expanded descriptions in delimited protobuf format")
var print = flag.Bool("print", false, "print information for each identifier in file specified")
var dof = flag.Bool("dof", false, "dimensionality analysis and reduction for file specified")

// flags for dof
var reduceDof = flag.Int("reduce", 0, "Reduce number of factors to specified number")
var minDistance = flag.Int("minimumDistance", 3, "Minimum distance from root")

func main() {
	flag.Parse()
	if *doVersion {
		fmt.Printf("%s v%s (%s)\n", os.Args[0], version, build)
		os.Exit(1)
	}
	if *database == "" {
		fmt.Fprint(os.Stderr, "error: missing mandatory database file\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	readOnly := true
	if *doImport || *precompute || *reset {
		readOnly = false
	}
	sct, err := terminology.NewService(*database, readOnly)
	if err != nil {
		log.Fatalf("couldn't open database: %v", err)
	}
	defer sct.Close()

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
		sct.ClearPrecomputations()
	}
	// perform import if an import root is specified
	if *doImport {
		if flag.NArg() == 0 {
			log.Fatalf("no input directories specified")
		}
		for _, filename := range flag.Args() {
			sct.PerformImport(filename, *verbose)
		}
	}

	// perform precomputations if requested
	if *precompute {
		sct.PerformPrecomputations()
	}

	// get statistics on store
	if *stats {
		s, err := sct.Statistics()
		if err != nil && err != terminology.ErrDatabaseNotInitialised {
			panic(err)
		}
		fmt.Printf("%v", s)
	}

	// export descriptions data in expanded denormalised format
	if *export {
		err := sct.Export(*lang)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *print {
		if flag.NArg() == 0 {
			log.Fatal("error: no input file(s) specified")
		}
		for _, filename := range flag.Args() {
			f, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			reader := bufio.NewReader(f)
			if err := analysis.Print(sct, reader); err != nil {
				log.Fatal(err)
			}
		}
	}

	// dimensionality analysis and reduction
	if *dof {
		if flag.NArg() == 0 {
			log.Fatal("error: no input file specified")
		}
		for _, filename := range flag.Args() {
			f, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			reader := bufio.NewReader(f)
			if *reduceDof > 0 {
				r := analysis.NewReducer(sct, *reduceDof, *minDistance)
				if err := r.ReduceCsv(reader, os.Stdout); err != nil {
					log.Fatal(err)
				}
			} else {
				factors, err := analysis.NumberFactors(reader)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(factors)
			}
		}
	}

	// optionally run a terminology server
	if *runserver {
		log.Fatal(server.RunServer(sct, *port, *lang))
	}
}
