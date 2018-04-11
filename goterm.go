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
	"flag"
	"fmt"
	"github.com/wardle/go-terminology/server"
	"github.com/wardle/go-terminology/terminology"
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
var runserver = flag.Bool("server", false, "Run terminology server")
var runrpc = flag.Bool("rpc", false, "Run RPC service")
var port = flag.Int("port", 8080, "Port to use when running server")

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
	if *doImport != "" {
		sct.PerformImport(*doImport)
	}

	if *precompute {
		sct.PerformPrecomputations()
	}

	if *runserver {
		server.RunServer(sct, *port)
	}
	if *runrpc {
		server.RunRPCServer(sct, *port)
	}
}
