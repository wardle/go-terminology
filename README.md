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


About go-terminology
===============
This is a SNOMED terminology server. It replaces an older java-based microservice [rsterminology](https://github.com/wardle/rsterminology) and is, 
as far as I know, the first Golang implementation of SNOMED-CT. It uses protobuf as its persistence data structure.

It is still in active development but is now in production-use. It can drive a [fast free-text search engine](https://github.com/wardle/rsterminology2) by exporting an optimised protobuf-based index. It can import a SNOMED-CT distribution and has no runtime dependencies except a filesystem. Importantly, it supports read-only operation from a shared filesystem making it ideal for use as a scalable microservice.

*Mark*