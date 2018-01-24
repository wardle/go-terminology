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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wardle/go-terminology/snomed"
	"github.com/wardle/go-terminology/terminology"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strconv"
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

	if *server {
		runServer(sct, 8080)
	}
}

// run our terminology server
// TODO:check precomputations have been run
func runServer(sct *terminology.Svc, port int) {
	router := mux.NewRouter()
	router.Handle("/snomedct/concepts/{id}", &handler{getConcept, sct}).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

// Result represents the result from a handler
type result struct {
	v      interface{}
	err    error
	status int
}

// HasError returns whether the result is a failure or not
func (r result) hasError() bool {
	return r.status >= 400
}

// Error returns the underlying error, falling back to generic error based on status code if necessary
func (r result) error() error {
	if r.err != nil {
		return r.err
	}
	if r.hasError() {
		return errors.New(http.StatusText(r.status))
	}
	return nil
}

type handler struct {
	Handler func(svc *terminology.Svc, w http.ResponseWriter, r *http.Request) result
	Svc     *terminology.Svc
}

// ServeHTTP allows your type to satisfy the http.Handler interface.
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := h.Handler(h.Svc, w, r)
	if result.hasError() {
		http.Error(w, result.error().Error(), result.status)
		return
	}
	if err := json.NewEncoder(w).Encode(result.v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getConcept(svc *terminology.Svc, w http.ResponseWriter, r *http.Request) result {
	params := mux.Vars(r)
	conceptID, err := strconv.Atoi(params["id"])
	if err != nil {
		return result{nil, err, http.StatusBadRequest}
	}
	concept, err := svc.GetConcept(conceptID)
	if err != nil {
		return result{nil, err, http.StatusNotFound}
	}
	descriptions, err := svc.GetDescriptions(concept)
	if err != nil {
		return result{nil, err, http.StatusInternalServerError}
	}
	allParents, err := svc.GetAllParentIDs(concept)
	if err != nil {
		return result{nil, err, http.StatusInternalServerError}
	}
	return result{&C{concept, allParents, descriptions}, nil, http.StatusOK}
}

// C represents a returned Concept including useful additional information
// TODO: include derivation of preferredTerm for the locale requested
type C struct {
	*snomed.Concept
	IsA          []int                 `json:"isA"`
	Descriptions []*snomed.Description `json:"descriptions"`
}
