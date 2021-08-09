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

const defaultFileName = "out.yaml"

func CLI(io Clio, factories component.Factories) {
	fileName := getFileName(io)
	service := map[string]interface{}{
		// this is the overview (top-level) part of the wizard, where the user just creates the pipelines
		"pipelines": pipelinesWizard(io, factories),
	}
	m := map[string]interface{}{
		"service": service,
	}
	dr := configschema.NewDefaultDirResolver()
	// build each individual component that the user chose.
	for componentGroup, names := range serviceToComponentNames(service) {
		handleComponent(factories, m, componentGroup, names, dr)
	}
	bytes, _ := yaml.Marshal(m)
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Println(string(bytes))
	writeFile(fileName, bytes)
}

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
