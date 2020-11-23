This repository has now been archived. 

My current SNOMED CT tooling is now available at [https://github.com/wardle/hermes](https://github.com/wardle/hermes).

-----


go-terminology
==============

[![Scc Count Badge](https://sloc.xyz/github/wardle/go-terminology)](https://github.com/wardle/go-terminology/)
[![Scc Cocomo Badge](https://sloc.xyz/github/wardle/go-terminology?category=cocomo)](https://github.com/wardle/go-terminology/)
[![Go Report Card](https://goreportcard.com/badge/github.com/wardle/go-terminology)](https://goreportcard.com/report/github.com/wardle/go-terminology)


Copyright 2018-2020 Mark Wardle and Eldrix Ltd

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

It replaces an older java-based microservice [rsterminology](https://github.com/wardle/rsterminology) and is, as far as I know, the first Golang implementation of SNOMED-CT. It uses protobuf as its persistence data structure and provides both a [gRPC](https://grpc.io) and REST API.

It is still in active development but is now in production-use. It can import a SNOMED-CT distribution and has no runtime dependencies except a filesystem. Importantly, it supports read-only operation from a shared filesystem making it ideal for use as a scalable microservice. 

It runs on Windows, Linux and Mac OS X and a range of other target architectures.

SNOMED CT is a medical ontology, and being able to process concepts and expressions, in the context of a wider information model, is critical for enabling the next generation of electronic health record systems to ensure that the right information is available at the right time, and that information is accessible and useful, both for direct care and the care of cohorts of patients, and for both professionals and patients themselves.

This tool is designed for use not only in operational clinical and research information systems, during the process of care, but also to support the rapid and straightforward use of SNOMED CT in analytics, by making it easier to parse and *understand* the meaning of a specific concept, such as "Does this patient have a type of motor neurone disease?". 

# Getting started

You will need a SNOMED CT distribution. For UK users, you can register and use [https://isd.digital.nhs.uk/trud3/user/guest/group/0/home](https://isd.digital.nhs.uk/trud3/user/guest/group/0/home). This example documents importing the following distributions:

* International release
* UK release
* UK dm+d (dictionary of medicines and devices) release

To get those from the NHS site, once registered and logged in, you need to seaarch for the approprite distribution, subscribe and then download the release - you may see more than one version, probably you want the latest.

The distributions are:

* [UK SNOMED CT Clinical Edition, RF2: Full, Snapshot & Delta](https://isd.digital.nhs.uk/trud3/user/authenticated/group/0/pack/26/subpack/101/releases) - contains the UK release and the International release on which it depends
* [UK SNOMED CT Drug Extension, RF2: Full, Snapshot & Delta](https://isd.digital.nhs.uk/trud3/user/authenticated/group/0/pack/26/subpack/105/releases)

Depending on your system, you may need to install some tool dependencies, particularly if you want to regenerate code from the protobuf/gRPC definitions.

```
brew install go protobuf protoc-gen-go

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
```

Then on all systems:

```
# If you're working within GOPATH you can get this repo (only needed the first time)
go get github.com/wardle/go-terminology

# Alternatively checkout the code using git
git clone https://github.com/wardle/go-terminology

# Get submodules (only needed if you are wanting to generate code from protobuf definitions)
git submodule update --init --recursive

# Go to the repo (assuming you have your go set up in the standard way)
cd ~/go/src/github.com/wardle/go-terminology

# Compile
make build

# Or run directly
go run goterm.go

# On linux, make sure there are enough file descriptors
ulimit -n 5000

# Import takes about 30 minutes for import, although it may take longer if you have a slow machine.
go run goterm.go -db ./snomed.db -v -import path/to/SNOMED-downloads/
```

> Import complete: : 28m45.6021888s: 958806 concepts, 2602531 descriptions, 6682738 relationships and 18503224 refset items...

```
# Before use, further precomputation is necessary. It now takes about 20 minutes to build the main indices and 10 minutes to precompute the search index. 
go run goterm.go -db ./snomed.db -precompute
```

> Indexing 21m5.803591824s: processed 2602531 descriptions, 6630693 relationships and 18503218 reference set items....
> Processed total: 2602531 descriptions in 10m41.973818972s.

Note: if you import from multiple distributions (such as the examples above in which I import the International, the UK and the UK dm+d distributions) there will be some duplicated components. Import will choose the version with the most recent "effective date". 

# And now you can run the terminology server 
go run goterm.go -db ./snomed.db -server
```
> 2019/08/22 20:30:54 gRPC Listening on [::]:8080
> 2019/08/22 20:30:54 HTTP Listening on :8081
```

# How to use the server

Full documentation of the API is available in [vendor/terminology/protos/server.proto](https://github.com/wardle/terminology/blob/master/protos/server.proto). In addition, Swagger documentation is generated as a part of the build.

# Example usage

Here are just a few examples of using the terminology server. They use [httpie](https://httpie.org) and some use [jq](https://stedolan.github.io/jq/) to extract parts of the result to make additional queries. You can obviously simply use a web browser instead.

For a short-time, I have created a small virtual server available at http://35.178.8.43

##### Fast free-text search:

Free-text search for ["MND"](http://35.178.8.43:8081/v1/snomed/search?s=mnd&is_a=64572001) in the 'disease' hierarchy. Note, multiple is_a parameters are supported. 
```
$ http get 'http://35.178.8.43:8081/v1/snomed/search?s=mnd&is_a=64572001'
```

The results are fast, and ideal for driving your autocompletion user interface. You can request to search in one or more hierarchies (via subsumption) or for items from specific reference sets. 

```
{
    "items": [
        {
            "concept_id": "37340000",
            "preferred_term": "Motor neuron disease",
            "term": "MND - Motor neurone disease"
        }
    ]
}

```
Get extended information about [laparoscopic cholecystectomy.](http://35.178.8.43:8081/v1/snomed/concepts/45595009/extended)
```
$ http get http://35.178.8.43:8081/v1/snomed/concepts/45595009/extended
```

Find out [how to refine a laparoscopic cholecystectomy](http://35.178.8.43:8081/v1/snomed/concepts/45595009/refinements), e.g. by access device, method and exact site(s). This can be used to drive interactive refinement, so that if a user chooses a procedure, you can then offer a choice to refine based on these characteristics.
```
$ http get http://35.178.8.43:8081/v1/snomed/concepts/45595009/refinements
```	

```
{
    "concept": {
        "active": true,
        "definition_status_id": "900000000000073002",
        "effective_time": "2002-01-31T00:00:00Z",
        "id": "45595009",
        "module_id": "900000000000207008"
    },
    "refinements": [
        {
            "attribute": {
                "concept_id": "425391005",
                "term": "Using access device"
            },
            "choices": [
                {
                    "concept_id": "701653007",
                    "term": "Externally-anchored laparoscopic retractor"
                },
                {
                    "concept_id": "462694004",
                    "term": "Neutral plasma surgical system control unit"
                },
                {
                    "concept_id": "465610003",
                    "term": "Vascular Doppler clamp"
                },
                {
                    "concept_id": "468274008",
                    "term": "Examination biliary catheter"
                },
                [...]
            ],
            "root_value": {
                "concept_id": "86174004",
                "term": "Laparoscope"
            }
        },
        {
            "attribute": {
                "concept_id": "405813007",
                "term": "Procedure site - Direct"
            },
            "choices": [
                {
                    "concept_id": "314739004",
                    "term": "Region of gallbladder"
                },
                {
                    "concept_id": "727273005",
                    "term": "Entire subserosa of gallbladder"
                },
                [...]
            ],
            "root_value": {
                "concept_id": "28231008",
                "term": "Gallbladder structure"
            }
        },
        {
            "attribute": {
                "concept_id": "260686004",
                "term": "Method"
            },
            "choices": [
                {
                    "concept_id": "289936007",
                    "term": "Shave excision"
                },
                {
                    "concept_id": "281838007",
                    "term": "Disarticulation - action"
                },
                [...]
            ],
            "root_value": {
                "concept_id": "129304002",
                "term": "Excision - action"
            }
        }
    ]
}
```
Get the descriptions (synonyms) for a ["surgical procedure".](http://35.178.8.43:8081/v1/snomed/concepts/387713003/descriptions)
```
$ http get http://35.178.8.43:8081/v1/snomed/concepts/387713003/descriptions
```

[Get multiple sclerosis concept](http://35.178.8.43:8081/v1/snomed/concepts/24700007)
```
$ http get http://35.178.8.43:8081/v1/snomed/concepts/24700007
```	

But that's only the raw data, and not terribly useful by itself. Instead, you want to understand this concept in the context of its relationships. 

Get [extended information about multiple sclerosis](http://35.178.8.43:8081/v1/snomed/concepts/24700007/extended), including a rapid way of determining subsumption. You can easily see that this is a "demyelinating disease of the CNS" (6118003) as it is listed in the "recursive_parent_ids" list. I prefer testing subsumption this way, at runtime in an EPR, and for analytics, because you can rapidly determine whether a patient has, say, a type of diabetes mellitus, while other servers have a subsumption endpoint that requires multiple round-trips.
```
$ http get http://35.178.8.43:8081/v1/snomed/concepts/24700007/extended
``` 

```
{
    "concept": {
        "active": true,
        "definition_status_id": "900000000000074008",
        "effective_time": "2002-01-31T00:00:00Z",
        "id": "24700007",
        "module_id": "900000000000207008"
    },
    "concept_refsets": [ list of reference sets to which this concept belongs
    ],
    "direct_parent_ids": [
        "6118003"
    ],
    "preferred_description": {
        "active": true,
        "case_significance": "900000000000448009",
        "concept_id": "24700007",
        "effective_time": "2017-07-31T00:00:00Z",
        "id": "41398015",
        "language_code": "en",
        "module_id": "900000000000207008",
        "term": "Multiple sclerosis",
        "type_id": "900000000000013009"
    },
    "recursive_parent_ids": [
        "6118003",
        "404684003",
        "138875005",
        "23853001",
        "362965005",
        "362975008",
        "64572001",
        "118940003",
        "123946008",
        "118234003",
        "246556002"
    ],
    "relationships": [ ... ]
}


```
To which [reference sets does multiple sclerosis belong](http://35.178.8.43:8081/v1/snomed/concepts/24700007/refsets)?
```
$ http get http://35.178.8.43:8081/v1/snomed/concepts/24700007/refsets
````

Parse a SNOMED expression
```
$ http get http://35.178.8.43:8081/v1/snomed/expression/parse?s="64572001 |disease|: 246454002 |occurrence| = 255407002 |neonatal|,  363698007 |finding site| = 113257007 |structure of cardiovascular system|"
```

```
{
    "clause": {
        "focus_concepts": [
            {
                "concept_id": "64572001",
                "term": "disease"
            }
        ],
        "refinements": [
            {
                "concept_value": {
                    "concept_id": "255407002",
                    "term": "neonatal"
                },
                "refinement_concept": {
                    "concept_id": "246454002",
                    "term": "occurrence"
                }
            },
            {
                "concept_value": {
                    "concept_id": "113257007",
                    "term": "structure of cardiovascular system"
                },
                "refinement_concept": {
                    "concept_id": "363698007",
                    "term": "finding site"
                }
            }
        ]
    }
}
```
Map "multiple sclerosis" into the [UK EU emergency care diagnostic subset](http://35.178.8.43:8081/v1/snomed/concepts/24700007/map?target_id=991411000000109) - and get 'multiple sclerosis', because it is in that subset.
```
$ http get http://35.178.8.43:8081/v1/snomed/concepts/24700007/map?target_id=991411000000109
```
 
 Now map a rare disorder ["ADEM" into the same diagnostic subset](http://35.178.8.43:8081/v1/snomed/concepts/83942000/map?target_id=991411000000109) - and get "demyelinating disease" (6118003) instead - useful for analytics to categorise highly granular or rarer diagnoses. It's really hard running analytics and reporting summaries unless you make it easier to categorise. 
```
$ http get http://35.178.8.43:8081/v1/snomed/concepts/83942000/map?target_id=991411000000109
```

```
{
    "active": true,
    "definition_status_id": "900000000000073002",
    "effective_time": "2002-01-31T00:00:00Z",
    "id": "6118003",
    "module_id": "900000000000207008"
}
```

Crossmap multiple sclerosis (24700007) to ICD-10 (G35X) [Go](http://35.178.8.43:8081/v1/snomed/concepts/24700007/crossmap?target_id=999002271000000101)
```
$ http get http://35.178.8.43:8081/v1/snomed/concepts/24700007/crossmap?target_id=999002271000000101
```

```
{
    "result": {
        "active": true,
        "complex_map": {
            "correlation": "447561005",
            "map_advice": "ALWAYS G35.X",
            "map_block": "1",
            "map_group": "1",
            "map_priority": "1",
            "map_target": "G35X"
        },
        "effective_time": "2018-04-01T00:00:00Z",
        "id": "57433204-2371-5c6f-855f-94ff9dad7ba6",
        "module_id": "999000031000000106",
        "referenced_component_id": "24700007",
        "refset_id": "999002271000000101"
    }
}
```
And let's get it as a [Read code](http://35.178.8.43:8081/v1/snomed/concepts/24700007/crossmap?target_id=900000000000497000)
```
$ http get http://35.178.8.43:8081/v1/snomed/concepts/24700007/crossmap?target_id=900000000000497000 | jq -r
```

```
F20..
```

Map it [back to SNOMED](http://35.178.8.43:8081/v1/snomed/crossmaps/900000000000497000/F20..)!
```
$ http get http://35.178.8.43:8081/v1/snomed/crossmaps/900000000000497000/F20.. | jq -r .translations[0].concept.id
```

```
24700007
```

This means we can process Read codes and make use of the SNOMED ontology e.g. when we're processing data from a GP system that doesn't yet use SNOMED, can we determine whether "F20" a type of demyelinating disease (6118003)?

```
$ id=`http get 35.178.8.43:8081/v1/snomed/crossmaps/900000000000497000/F20.. | jq -r .translations[0].concept.id`; http get "35.178.8.43:8081/v1/snomed/subsumes?code_a=$id&code_b=6118003" | jq -r .result
```
Yes!
```
SUBSUMED_BY
```
Is XU6qV (diabetes) a type of demyelinating disease? 
```
$ id=`http get 35.178.8.43:8081/v1/snomed/crossmaps/900000000000497000/XU6qV | jq -r .translations[0].concept.id`; http get "35.178.8.43:8081/v1/snomed/subsumes?code_a=$id&code_b=6118003" | jq -r .result
```
No!
```
NOT_SUBSUMED
```

Ok, so is XU6qV a disorder of carbohydrate metabolism (20957000)? 
```
$ id=`http get 35.178.8.43:8081/v1/snomed/crossmaps/900000000000497000/XU6qV | jq -r .translations[0].concept.id`; http get "35.178.8.43:8081/v1/snomed/subsumes?code_a=$id&code_b=20957000" | jq -r .result
```
Yes!
```
SUBSUMED_BY
```

Our user has searched for "heart attack" in their old unstructured letters. Can we help by also searching for [synonyms of this term](http://35.178.8.43:8081/v1/snomed/synonyms?s=heart%20attack&is_a=64572001)?
```
$ http get 'http://35.178.8.43:8081/v1/snomed/synonyms?s=heart%20attack&is_a=64572001'
```

Now we can help patients find information in their legacy information such as document archives. 

```
{"result":{"s":"Myocardial infarction"}}
{"result":{"s":"Infarction of heart"}}
{"result":{"s":"Cardiac infarction"}}
{"result":{"s":"Heart attack"}}
{"result":{"s":"MI - Myocardial infarction"}}
{"result":{"s":"Myocardial infarct"}}
```
It's also useful for case-finding for research, so let's find all of the terms that might have been used to [record a patient having a stroke](http://35.178.8.43:8081/v1/snomed/synonyms?s=stroke&is_a=64572001):
```
$ http get 'http://35.178.8.43:8081/v1/snomed/synonyms?s=stroke&is_a=64572001'
```
```
{"result":{"s":"Thalamic infarction"}}
{"result":{"s":"CVA - Cerebrovascular accident"}}
{"result":{"s":"Brain stem infarct"}}
{"result":{"s":"PACS - Partial anterior cerebral circulation stroke"}}
{"result":{"s":"TACI - Total anterior cerebral circulation infarction"}}
[...etc...]
```
We won't only get "stroke", but also "cerebral infarction", "thalamic infarction" and others and use those terms to search our legacy text archives.

*Mark*
