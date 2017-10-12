// simple proof-of-concept SNOMED code
package main

import (
	"flag"

	"bitbucket.org/wardle/go-snomed/mcqs"
	_ "github.com/lib/pq"
)

// A simple proof-of-concept application to generate fake exam questions
func main() {
	var (
		number         int
		configFilename string
		precompute     bool
		truth          bool
		prevalence     bool
		diagnostic     bool
	)
	flag.IntVar(&number, "n", 0, "Number to generate. Default: all")
	flag.StringVar(&configFilename, "config", "config.yml", "Specifies a configuration file.")
	flag.BoolVar(&precompute, "precompute", false, "Generate a set of pre-computed SNOMED-CT data files.")
	flag.BoolVar(&truth, "truth", false, "Using precomputed SNOMED-CT, generate a fake truth dataset linking diagnostic concepts with clinical features.")
	flag.BoolVar(&prevalence, "prevalence", false, "Using fake prevalence figures, generate fake questions simply to model prevalence.")
	flag.BoolVar(&diagnostic, "diagnostic", false, "Using fake truth dataset, generate fake questions for machine learning proof-of-concept.")
	flag.Parse()
	if precompute || truth || prevalence || diagnostic {
		if precompute {
			mcqs.GenerateSnomedCT()
		}
		if truth {
			mcqs.GenerateTruth()
		}
		if prevalence {
			mcqs.GeneratePrevalence()
		}
		if diagnostic {
			mcqs.GenerateDiagnostic()
		}
	} else {
		flag.PrintDefaults()
	}
}
