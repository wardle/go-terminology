package server

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/wardle/go-terminology/terminology"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// RunServer runs a terminology server
// TODO:check precomputations have been run
func RunServer(sct *terminology.Svc, port int) {
	router := mux.NewRouter()
	router.Handle("/snomedct/concepts/{id}", &handler{getConcept, sct}).Methods("GET")
	router.Handle("/snomedct/concepts/{id}/descriptions", &handler{getConceptDescriptions, sct}).Methods("GET")
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
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		//grpcServer.ServeHTTP(w, r)
	} else {
		//otherHandler.ServeHTTP(w, r)
	}
	result := h.Handler(h.Svc, w, r)
	if result.hasError() {
		http.Error(w, result.error().Error(), result.status)
		return
	}
	if err := json.NewEncoder(w).Encode(result.v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
