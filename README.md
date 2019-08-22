go-terminology
==============

Copyright 2018 Mark Wardle and Eldrix Ltd

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.


# About go-terminology

This is a SNOMED CT terminology server. 

It replaces an older java-based microservice [rsterminology](https://github.com/wardle/rsterminology) and is, as far as I know, the first Golang implementation of SNOMED-CT. It uses protobuf as its persistence data structure and provides both a GRPC and REST API.

It is still in active development but is now in production-use. It can import a SNOMED-CT distribution and has no runtime dependencies except a filesystem. Importantly, it supports read-only operation from a shared filesystem making it ideal for use as a scalable microservice. 

SNOMED CT is a medical ontology, and being able to process concepts and expressions, in the context of a wider information model, is critical for enabling the next generation of electronic health record systems to ensure that the right information is available at the right time, and that information is accessible and useful, both for direct care and the care of cohorts of patients, and for both professionals and patients themselves.

# Getting started
```
    # Fetch latest dependencies (currently no support for go modules)
	go get -u
	
	# Compile
	go build
	
    # Download a SNOMED distribution. For UK users, you can register and use https://isd.digital.nhs.uk/trud3/user/guest/group/0/home

	# Import (provide root of SNOMED extract - both International & UK data are imported, as well as dm+d). 
    # Takes about 30 minutes for import and indexing, although it may take longer if you have a slow machine.
	./gts -db ./snomed.db -v -import path/to/SNOMED-download/
	## Import complete: : 28m45.6021888s: 958806 concepts, 2602531 descriptions, 6682738 relationships and 18503224 refset items...

	# To use text search, further precomputation is necessary. Takes about 10 minutes to build the indexes.
	./gts -db ./snomed.db -precompute
	## Processed total: 2602531 descriptions in 10m41.973818972s.

    # And now you can run the terminology server 
    ./gts -db ./snomed.db -server

    #
    # Make some test calls (using httpie)
	#
    # Free-text search for "MND" in the 'disease' hierarchy. Note, multiple is_a parameters are supported.
    http get 'http://localhost:8081/v1/snomed/search?s=mnd&is_a=64572001'
	
	# get extended information about laparoscopic cholecystectomy
	http get http://localhost:8081/v1/snomed/concepts/45595009/extended

    # Find out how to refine a laparoscopic cholecystectomy, e.g. by access device, method and exact site(s)
    http get http://localhost:8081/v1/snomed/concepts/45595009/refinements
	
    # Get the descriptions (synonyms) for a "surgical procedure"
    http get http://localhost:8081/v1/snomed/concepts/387713003/descriptions
    
    # get multiple sclerosis
	http get http://localhost:8081/v1/snomed/concepts/24700007
	
    # get extended information about multiple sclerosis, including a rapid way of determining subsumption
    # you can easily see that this is a "demyelinating disease of the CNS" (6118003) as it is listed in the "recursive_parent_ids" list.
	http get http://localhost:8081/v1/snomed/concepts/24700007/extended
    
    # to which reference sets does multiple sclerosis belong?
   	http get http://localhost:8081/v1/snomed/concepts/24700007/refsets

    # parse a SNOMED expression
    http get http://localhost:8081/v1/snomed/expression/parse?s="64572001 |disease|: 246454002 |occurrence| = 255407002 |neonatal|,  363698007 |finding site| = 113257007 |structure of cardiovascular system|"

    # map "multiple sclerosis" into the UK EU emergency care diagnostic subset - and get 'multiple sclerosis'
    http get localhost:8081/v1/snomed/concepts/24700007/map?target_id=991411000000109

    # now map a rare disorder "ADEM" into the same diagnostic subset - and get "demyelinating disease" (6118003) instead - useful for analytics to categorise highly granular or rarer diagnoses.
    http get localhost:8081/v1/snomed/concepts/83942000/map?target_id=991411000000109

    # crossmap multiple sclerosis (24700007) to ICD-10 (G35X)
    http get http://localhost:8081/v1/snomed/concepts/24700007/crossmap?target_id=999002271000000101

    # and let's get it as a Read code
    http get localhost:8081/v1/snomed/concepts/24700007/crossmap?target_id=900000000000497000

    # crossmap ADEM to ICD-10
    http get localhost:8081/v1/snomed/concepts/83942000/crossmap?target_id=999002271000000101

    # if you have [jq(https://stedolan.github.io/jq/)] installed, you can parse JSON from the command line, like this
    
    # directly get Read code corresponding to 24700007 (multiple sclerosis)
    http get localhost:8081/v1/snomed/concepts/24700007/crossmap?target_id=900000000000497000 | jq .result.simple_map.map_target

    # map it back to SNOMED!
    http get localhost:8081/v1/snomed/crossmaps/900000000000497000/F20.. | jq -r .translations[0].concept.id

    # this means we can process Read codes and make use of the SNOMED ontology...
    # e.g. is F20 a type of demyelinating disease (6118003)? - Yes!
    id=`http get localhost:8081/v1/snomed/crossmaps/900000000000497000/F20.. | jq -r .translations[0].concept.id`; http get "localhost:8081/v1/snomed/subsumes?code_a=$id&code_b=6118003" | jq .result

    # e.g. is XU6qV (diabetes) a type of demyelinating disease? - No!
    id=`http get localhost:8081/v1/snomed/crossmaps/900000000000497000/XU6qV | jq -r .translations[0].concept.id`; http get "localhost:8081/v1/snomed/subsumes?code_a=$id&code_b=6118003" | jq .result

    # ok, so is XU6qV a disorder of carohydrateb metabolism (20957000)? Yes!
    id=`http get localhost:8081/v1/snomed/crossmaps/900000000000497000/XU6qV | jq -r .translations[0].concept.id`; http get "localhost:8081/v1/snomed/subsumes?code_a=$id&code_b=20957000" | jq .result

    # concepts that match ICD-10 code G040....
    http get localhost:8081/v1/snomed/crossmaps/999002271000000101/G040 | jq '[.translations[].concept.id]'

	# See server.proto for more details of the API
```

*Mark*