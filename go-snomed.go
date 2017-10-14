// simple proof-of-concept SNOMED code
package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/wardle/go-snomed/mcqs"
	_ "github.com/lib/pq"
)

// A simple proof-of-concept application to generate fake exam questions
func main() {
	var (
		number     int
		path       string
		precompute bool
		truth      bool
		prevalence bool
		diagnostic bool
	)
	flag.IntVar(&number, "n", 0, "Number to generate. Default: all")
	flag.StringVar(&path, "path", "./", "Location of data files")
	flag.BoolVar(&precompute, "precompute", false, "Generate a set of pre-computed SNOMED-CT data files.")
	flag.BoolVar(&truth, "truth", false, "Using precomputed SNOMED-CT, generate a fake truth dataset linking diagnostic concepts with clinical features.")
	flag.BoolVar(&prevalence, "prevalence", false, "Using fake prevalence figures, generate fake questions simply to model prevalence.")
	flag.BoolVar(&diagnostic, "diagnostic", false, "Using fake truth dataset, generate fake questions for machine learning proof-of-concept.")
	flag.Parse()
	if precompute || truth || prevalence || diagnostic {
		var dataset mcqs.SnomedDataset
		var err error
		if precompute {
			dataset, err = mcqs.GenerateSnomedCT(path)
		}
		if err == nil && dataset == nil {
			dataset, err = mcqs.OpenSnomedDataset(path)
		}
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprint(os.Stderr, "Error: did not find precompute files.\n")
			} else {
				fmt.Fprintf(os.Stderr, "Error opening SNOMED dataset: %s\n", err)
			}
			os.Exit(1)
		}
		if truth {
			mcqs.GenerateFakeTruth(dataset)
		}
		if prevalence {
			mcqs.GeneratePrevalence(dataset)
		}
		if diagnostic {
			mcqs.GenerateDiagnostic(dataset)
		}
	} else {
		flag.PrintDefaults()
	}
}
