// simple proof-of-concept SNOMED code
package main

import (
	"flag"

	_ "github.com/lib/pq"
	"github.com/wardle/go-snomed/snomed"
)

// A simple proof-of-concept application to generate fake exam questions
func main() {
	var (
		number         int
		configFilename string
		generate       bool
		truth          bool
		prevalence     bool
		questions      bool
	)
	flag.Int("n", number, "Number to generate. Default: all")
	flag.String("config", configFilename, "Specifies a configuration file.")
	flag.Bool("precompute", generate, "Generate a set of pre-computed SNOMED-CT data files.")
	flag.Bool("truth", truth, "Using precomputed SNOMED-CT, generate a fake truth dataset linking diagnostic concepts with clinical features.")
	flag.Bool("qprev", prevalence, "Using fake prevalence figures, generate fake questions simply to model prevalence.")
	flag.Bool("qdiag", questions, "Using fake truth dataset, generate fake questions for machine learning proof-of-concept.")
	flag.Parse()
	if generate || truth || prevalence || questions {
		if generate {
			snomed.GenerateSnomedCT()
		}
		if truth {
			snomed.GenerateTruth()
		}
		if prevalence {
			snomed.GeneratePrevalence()
		}
		if questions {
			snomed.GenerateDiagnostic()
		}
	} else {
		flag.PrintDefaults()
	}
}
