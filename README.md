Synthea Uploader
===================================================================================================================================================================

The Synthea uploader uploads Synthea-generated FHIR bundles to a FHIR server.

Usage
-----

```
$ ./uploader -help
Usage of ./uploader:
  -fhir string
    	Endpoint for the FHIR server
  -path string
    	Path to the folder containing records to upload
```

Example
-------

```
./uploader -fhir http://localhost:3001 -path ~/dev/synthetichealth/synthea/output/fhir/
```

License
-------

Copyright 2016 The MITRE Corporation

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
