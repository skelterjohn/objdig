/*
Copyright 2015 Google Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
Usage: objdig [-j|-y] KEY+

objdig takes a document on stdin and prints the values for one or more keys
on stdout.

If the doc is JSON, use -j. For YAML, use -y.

A key is as defined for github.com/skelterjohn/overwrite.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/skelterjohn/overwrite"
	"gopkg.in/yaml.v2"
)

func main() {
	typ := ""
	var args []string
	for _, arg := range os.Args[1:] {
		switch arg {
		case "-j", "-y":
			if typ != "" {
				fmt.Fprintln(os.Stderr, "only one doc type may be declared")
				os.Exit(1)
			}
			typ = arg
		case "-r":
			fmt.Fprintln(os.Stderr, "rjson support is not implemented")
			os.Exit(1)
		default:
			args = append(args, arg)
		}
	}

	doc := map[string]interface{}{}
	switch typ {
	case "-j":
		if err := json.NewDecoder(os.Stdin).Decode(&doc); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "-y":
		if data, err := ioutil.ReadAll(os.Stdin); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		} else {
			if err := yaml.Unmarshal(data, doc); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	default:
		fmt.Fprintln(os.Stderr, "must provide a doc type (-j or -y)")
		os.Exit(1)
	}

	for _, arg := range args {
		obj, err := overwrite.Fetch(doc, arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "for %q: %s\n", arg, err)
		} else {
			fmt.Println(obj)
		}
	}
}
