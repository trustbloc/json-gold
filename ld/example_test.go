// Copyright 2015-2017 Piprate Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ld_test

import (
	"fmt"
	"log"

	"github.com/piprate/json-gold/ld"
)

func ExampleJsonLdProcessor_Expand_online() {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	// expanding remote document

	expanded, err := proc.Expand("http://json-ld.org/test-suite/tests/expand-0002-in.jsonld", options)
	if err != nil {
		log.Println("Error when expanding JSON-LD document:", err)
		return
	}

	ld.PrintDocument("JSON-LD expansion succeeded", expanded)

	// Output:
	// JSON-LD expansion succeeded
	// [
	//   {
	//     "@id": "http://example.com/id1",
	//     "@type": [
	//       "http://example.com/t1"
	//     ],
	//     "http://example.com/term1": [
	//       {
	//         "@value": "v1"
	//       }
	//     ],
	//     "http://example.com/term2": [
	//       {
	//         "@type": "http://example.com/t2",
	//         "@value": "v2"
	//       }
	//     ],
	//     "http://example.com/term3": [
	//       {
	//         "@language": "en",
	//         "@value": "v3"
	//       }
	//     ],
	//     "http://example.com/term4": [
	//       {
	//         "@value": 4
	//       }
	//     ],
	//     "http://example.com/term5": [
	//       {
	//         "@value": 50
	//       },
	//       {
	//         "@value": 51
	//       }
	//     ]
	//   }
	// ]
}

func ExampleJsonLdProcessor_Expand_inmemory() {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	// expanding in-memory document

	doc := map[string]interface{}{
		"@context":  "http://schema.org/",
		"@type":     "Person",
		"name":      "Jane Doe",
		"jobTitle":  "Professor",
		"telephone": "(425) 123-4567",
		"url":       "http://www.janedoe.com",
	}

	expanded, err := proc.Expand(doc, options)
	if err != nil {
		log.Println("Error when expanding JSON-LD document:", err)
		return
	}

	ld.PrintDocument("JSON-LD expansion succeeded", expanded)

	// Output:
	// JSON-LD expansion succeeded
	// [
	//   {
	//     "@type": [
	//       "http://schema.org/Person"
	//     ],
	//     "http://schema.org/jobTitle": [
	//       {
	//         "@value": "Professor"
	//       }
	//     ],
	//     "http://schema.org/name": [
	//       {
	//         "@value": "Jane Doe"
	//       }
	//     ],
	//     "http://schema.org/telephone": [
	//       {
	//         "@value": "(425) 123-4567"
	//       }
	//     ],
	//     "http://schema.org/url": [
	//       {
	//         "@id": "http://www.janedoe.com"
	//       }
	//     ]
	//   }
	// ]
}

func ExampleJsonLdProcessor_Compact() {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	doc := map[string]interface{}{
		"@id": "http://example.org/test#book",
		"http://example.org/vocab#contains": map[string]interface{}{
			"@id": "http://example.org/test#chapter",
		},
		"http://purl.org/dc/elements/1.1/title": "Title",
	}

	context := map[string]interface{}{
		"@context": map[string]interface{}{
			"dc": "http://purl.org/dc/elements/1.1/",
			"ex": "http://example.org/vocab#",
			"ex:contains": map[string]interface{}{
				"@type": "@id",
			},
		},
	}

	compactedDoc, err := proc.Compact(doc, context, options)
	if err != nil {
		log.Println("Error when compacting JSON-LD document:", err)
		return
	}

	ld.PrintDocument("JSON-LD compact doc", compactedDoc)

	// Output:
	// JSON-LD compact doc
	// {
	//   "@context": {
	//     "dc": "http://purl.org/dc/elements/1.1/",
	//     "ex": "http://example.org/vocab#",
	//     "ex:contains": {
	//       "@type": "@id"
	//     }
	//   },
	//   "@id": "http://example.org/test#book",
	//   "dc:title": "Title",
	//   "ex:contains": "http://example.org/test#chapter"
	// }
}

func ExampleJsonLdProcessor_Flatten() {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	doc := map[string]interface{}{
		"@context": []interface{}{
			map[string]interface{}{
				"name": "http://xmlns.com/foaf/0.1/name",
				"homepage": map[string]interface{}{
					"@id":   "http://xmlns.com/foaf/0.1/homepage",
					"@type": "@id",
				},
			},
			map[string]interface{}{
				"ical": "http://www.w3.org/2002/12/cal/ical#",
			},
		},
		"@id":           "http://example.com/speakers#Alice",
		"name":          "Alice",
		"homepage":      "http://xkcd.com/177/",
		"ical:summary":  "Alice Talk",
		"ical:location": "Lyon Convention Centre, Lyon, France",
	}

	flattenedDoc, err := proc.Flatten(doc, nil, options)
	if err != nil {
		log.Println("Error when flattening JSON-LD document:", err)
		return
	}

	ld.PrintDocument("JSON-LD flattened doc", flattenedDoc)

	// Output:
	// JSON-LD flattened doc
	// [
	//   {
	//     "@id": "http://example.com/speakers#Alice",
	//     "http://www.w3.org/2002/12/cal/ical#location": [
	//       {
	//         "@value": "Lyon Convention Centre, Lyon, France"
	//       }
	//     ],
	//     "http://www.w3.org/2002/12/cal/ical#summary": [
	//       {
	//         "@value": "Alice Talk"
	//       }
	//     ],
	//     "http://xmlns.com/foaf/0.1/homepage": [
	//       {
	//         "@id": "http://xkcd.com/177/"
	//       }
	//     ],
	//     "http://xmlns.com/foaf/0.1/name": [
	//       {
	//         "@value": "Alice"
	//       }
	//     ]
	//   }
	// ]
}

func ExampleJsonLdProcessor_Frame() {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	doc := map[string]interface{}{
		"@context": map[string]interface{}{
			"dc":          "http://purl.org/dc/elements/1.1/",
			"ex":          "http://example.org/vocab#",
			"ex:contains": map[string]interface{}{"@type": "@id"},
		},
		"@graph": []interface{}{
			map[string]interface{}{
				"@id":         "http://example.org/test/#library",
				"@type":       "ex:Library",
				"ex:contains": "http://example.org/test#book",
			},
			map[string]interface{}{
				"@id":            "http://example.org/test#book",
				"@type":          "ex:Book",
				"dc:contributor": "Writer",
				"dc:title":       "My Book",
				"ex:contains":    "http://example.org/test#chapter",
			},
			map[string]interface{}{
				"@id":            "http://example.org/test#chapter",
				"@type":          "ex:Chapter",
				"dc:description": "Fun",
				"dc:title":       "Chapter One",
			},
		},
	}

	frame := map[string]interface{}{
		"@context": map[string]interface{}{
			"dc": "http://purl.org/dc/elements/1.1/",
			"ex": "http://example.org/vocab#",
		},
		"@type": "ex:Library",
		"ex:contains": map[string]interface{}{
			"@type": "ex:Book",
			"ex:contains": map[string]interface{}{
				"@type": "ex:Chapter",
			},
		},
	}

	framedDoc, err := proc.Frame(doc, frame, options)
	if err != nil {
		log.Println("Error when framing JSON-LD document:", err)
		return
	}

	ld.PrintDocument("JSON-LD framed doc", framedDoc)

	// Output:
	// JSON-LD framed doc
	// {
	//   "@context": {
	//     "dc": "http://purl.org/dc/elements/1.1/",
	//     "ex": "http://example.org/vocab#"
	//   },
	//   "@graph": [
	//     {
	//       "@id": "http://example.org/test/#library",
	//       "@type": "ex:Library",
	//       "ex:contains": {
	//         "@id": "http://example.org/test#book",
	//         "@type": "ex:Book",
	//         "dc:contributor": "Writer",
	//         "dc:title": "My Book",
	//         "ex:contains": {
	//           "@id": "http://example.org/test#chapter",
	//           "@type": "ex:Chapter",
	//           "dc:description": "Fun",
	//           "dc:title": "Chapter One"
	//         }
	//       }
	//     }
	//   ]
	// }
}

func ExampleJsonLdProcessor_ToRDF() {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.Format = "application/nquads"

	// this JSON-LD document was taken from http://json-ld.org/test-suite/tests/toRdf-0028-in.jsonld
	doc := map[string]interface{}{
		"@context": map[string]interface{}{
			"sec":        "http://purl.org/security#",
			"xsd":        "http://www.w3.org/2001/XMLSchema#",
			"rdf":        "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
			"dc":         "http://purl.org/dc/terms/",
			"sec:signer": map[string]interface{}{"@type": "@id"},
			"dc:created": map[string]interface{}{"@type": "xsd:dateTime"},
		},
		"@id":                "http://example.org/sig1",
		"@type":              []interface{}{"rdf:Graph", "sec:SignedGraph"},
		"dc:created":         "2011-09-23T20:21:34Z",
		"sec:signer":         "http://payswarm.example.com/i/john/keys/5",
		"sec:signatureValue": "OGQzNGVkMzVm4NTIyZTkZDYMmMzQzNmExMgoYzI43Q3ODIyOWM32NjI=",
		"@graph": map[string]interface{}{
			"@id":      "http://example.org/fact1",
			"dc:title": "Hello World!",
		},
	}
	triples, err := proc.ToRDF(doc, options)
	if err != nil {
		log.Println("Error running ToRDF:", err)
		return
	}

	fmt.Printf("%s\n", triples)

	// Output:
	// <http://example.org/fact1> <http://purl.org/dc/terms/title> "Hello World!" <http://example.org/sig1> .
	// <http://example.org/sig1> <http://purl.org/dc/terms/created> "2011-09-23T20:21:34Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
	// <http://example.org/sig1> <http://purl.org/security#signatureValue> "OGQzNGVkMzVm4NTIyZTkZDYMmMzQzNmExMgoYzI43Q3ODIyOWM32NjI=" .
	// <http://example.org/sig1> <http://purl.org/security#signer> <http://payswarm.example.com/i/john/keys/5> .
	// <http://example.org/sig1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://purl.org/security#SignedGraph> .
	// <http://example.org/sig1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://www.w3.org/1999/02/22-rdf-syntax-ns#Graph> .
	//
}

func ExampleJsonLdProcessor_FromRDF() {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	triples := `
	<http://example.com/Subj1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.com/Type> .
	<http://example.com/Subj1> <http://example.com/prop1> <http://example.com/Obj1> .
	<http://example.com/Subj1> <http://example.com/prop2> "Plain" .
	<http://example.com/Subj1> <http://example.com/prop2> "2012-05-12"^^<http://www.w3.org/2001/XMLSchema#date> .
	<http://example.com/Subj1> <http://example.com/prop2> "English"@en .
`

	doc, err := proc.FromRDF(triples, options)
	if err != nil {
		log.Println("Error running FromRDF:", err)
		return
	}

	ld.PrintDocument("JSON-LD doc from RDF", doc)

	// Output:
	// JSON-LD doc from RDF
	// [
	//   {
	//     "@id": "http://example.com/Subj1",
	//     "@type": [
	//       "http://example.com/Type"
	//     ],
	//     "http://example.com/prop1": [
	//       {
	//         "@id": "http://example.com/Obj1"
	//       }
	//     ],
	//     "http://example.com/prop2": [
	//       {
	//         "@value": "Plain"
	//       },
	//       {
	//         "@type": "http://www.w3.org/2001/XMLSchema#date",
	//         "@value": "2012-05-12"
	//       },
	//       {
	//         "@language": "en",
	//         "@value": "English"
	//       }
	//     ]
	//   }
	// ]
}

func ExampleJsonLdProcessor_Normalize() {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	options.Format = "application/nquads"
	options.Algorithm = "URDNA2015"

	doc := map[string]interface{}{
		"@context": map[string]interface{}{
			"ex": "http://example.org/vocab#",
		},
		"@id":   "http://example.org/test#example",
		"@type": "ex:Foo",
		"ex:embed": map[string]interface{}{
			"@type": "ex:Bar",
		},
	}

	normalizedTriples, err := proc.Normalize(doc, options)
	if err != nil {
		log.Println("Error running Normalize:", err)
		return
	}

	fmt.Printf("%s\n", normalizedTriples)

	// Output:
	// <http://example.org/test#example> <http://example.org/vocab#embed> _:c14n0 .
	// <http://example.org/test#example> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/vocab#Foo> .
	// _:c14n0 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/vocab#Bar> .
}