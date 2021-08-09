// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configwiz

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	"go.opentelemetry.io/collector/cmd/configschema/configschema"
	"go.opentelemetry.io/collector/component"
)

const defaultFileName = "compose.yaml"

func CLI(io Clio, factories component.Factories) {
	fileName := getFileName(io)
	service := map[string]interface{}{
		// this is the overview (top-level) part of the wizard, where the user just creates the pipelines
		"pipelines": pipelinesWizard(io, factories),
	}
	m := map[string]interface{}{
		"service": service,
	}
	dr := configschema.NewDirResolver(".", "github.com/open-telemetry/opentelemetry-collector-contrib")
	// build each individual component that the user chose.
	for componentGroup, names := range serviceToComponentNames(service) {
		handleComponent(factories, m, componentGroup, names, dr)
	}
	out := buildYamlFile(m)
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Println(out)
	rv := buildYamlFile(m)
	writeFile(fileName, []byte(rv))
}

// getFileName prompts the user to input a fileName, and returns the value that they chose.
// defaults to out.yaml
func getFileName(io Clio) string {
	pr := io.newIndentingPrinter(0)
	pr.println("Name of file (default out.yaml):")
	pr.print("> ")
	fileName := io.Read("")
	if fileName == "" {
		fileName = defaultFileName
	}
	if !strings.HasSuffix(fileName, ".yaml") {
		fileName += ".yaml"
	}
	return fileName
}

// buildYamlFile outputs a .yaml file based on the configuration we created
func buildYamlFile(m map[string]interface{}) string {
	rv := ""
	out := outYaml{
		map[string]interface{}{
			"receivers": m["receivers"],
		},
		map[string]interface{}{
			"processors": m["processors"],
		},
		map[string]interface{}{
			"exporters" : m["exporters"],
		},
		map[string]interface{}{
			"extensions": m["extensions"],
		},
		map[string]interface{}{
			"service" : m["service"],
		},
	}
	rv += marshal(out.receivers)
	rv += marshal(out.processors)
	rv += marshal(out.exporters)
	rv += marshal(out.extensions)
	rv += marshal(out.service)
	return rv
}

func marshal(comp interface{}) string {
	if bytes, err := yaml.Marshal(comp); err != nil {
		panic(err)
	} else {
		return string(bytes)
	}
}


type outYaml struct {
	receivers interface{}
	processors interface{}
	exporters interface{}
	extensions interface{}
	service interface{}
}
