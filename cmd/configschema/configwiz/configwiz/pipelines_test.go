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
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeReaderPipe struct {
	userInput []string
	input     int
}

func (r *fakeReaderPipe) read(defaultVal string) string {
	out := r.userInput[r.input]
	if r.input < len(r.userInput)-1 {
		r.input++
	}
	return out
}

func buildCompExpected(indent int, prefix string, compType string, compNames []string) string {
	const tabSize = 4
	space := indent * tabSize
	tab := strings.Repeat(" ", space)
	expected := fmt.Sprintf("%sAdd %s (enter to skip)\n", tab, compType)
	for i := 0; i < 3; i++ {
		expected += fmt.Sprintf("%s%d: %s\n", tab, i, compNames[i])
	}
	expected += tab + "> "
	return prefix + expected
}

func TestComponentListWizard(t *testing.T) {
	tab := strings.Repeat(" ", 4)

	w := fakeWriter{}
	r := fakeReader{}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	compGroup := "test"
	compNames := []string{"comp1", "comp2", "comp3"}
	componentListWizard(io, pr, compGroup, compNames)
	expected := fmt.Sprintf("%sCurrent %ss: []\n", tab, compGroup)
	expected = buildCompExpected(1, expected, compGroup, compNames)
	assert.Equal(t, expected, w.programOutput)

	//Testing inputting a value inside
	w2 := fakeWriter{}
	r2 := fakeReaderPipe{userInput: []string{"0", ""}, input: 0}
	io2 := clio{w2.write, r2.read}
	pr2 := io2.newIndentingPrinter(1)
	componentListWizard(io2, pr2, compGroup, compNames)
	expected2 := expected + fmt.Sprintf("%s%s %s extended name (optional) > ", tab, compNames[0], compGroup)
	expected2 += fmt.Sprintf("%sCurrent tests: [%s]\n", tab, compNames[0])
	expected2 = buildCompExpected(1, expected2, compGroup, compNames)
	assert.Equal(t, expected2, w2.programOutput)

	// Testing extension and the input of another value
	w3 := fakeWriter{}
	r3 := fakeReaderPipe{userInput: []string{"0", "extension", "1", "", ""}, input: 0}
	io3 := clio{w3.write, r3.read}
	pr3 := io3.newIndentingPrinter(1)
	componentListWizard(io3, pr3, compGroup, compNames)
	expected3 := expected + fmt.Sprintf("%s%s %s extended name (optional) > ", tab, compNames[0], compGroup)
	expected3 += fmt.Sprintf("%sCurrent tests: [%s/extension]\n", tab, compNames[0])
	expected3 = buildCompExpected(1, expected3, compGroup, compNames)
	expected3 += fmt.Sprintf("%s%s %s extended name (optional) > ", tab, compNames[1], compGroup)
	expected3 += fmt.Sprintf("%sCurrent tests: [%s/extension, %s]\n", tab, compNames[0], compNames[1])
	expected3 = buildCompExpected(1, expected3, compGroup, compNames)
	assert.Equal(t, expected3, w3.programOutput)
}

func TestComponentNameWizard(t *testing.T) {
	// Test components get printed out
	w := fakeWriter{}
	r := fakeReader{}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	compType := "test"
	compNames := []string{"comp1", "comp2", "comp3"}
	componentNameWizard(io, pr, compType, compNames)
	tab := strings.Repeat(" ", 4)
	expected := buildCompExpected(1, "", compType, compNames)
	assert.Equal(t, expected, w.programOutput)

	// Test extended name
	w2 := fakeWriter{}
	r2 := fakeReader{"0"}
	io2 := clio{w2.write, r2.read}
	pr2 := io2.newIndentingPrinter(1)
	out, val := componentNameWizard(io2, pr2, compType, compNames)
	assert.Equal(t, compNames[0], out)
	assert.Equal(t, val, "0")
	expected2 := expected + fmt.Sprintf("%s%s %s extended name (optional) > ", tab, out, compType)
	assert.Equal(t, expected2, w2.programOutput)

	// Test error
	w3 := fakeWriter{}
	r3 := fakeReaderPipe{[]string{"-1", ""}, 0}
	io3 := clio{w3.write, r3.read}
	pr3 := io3.newIndentingPrinter(1)
	componentNameWizard(io3, pr3, compType, compNames)
	expected3 := expected + "Invalid input. Try again.\n"
	expected3 += expected
	assert.Equal(t, expected3, w3.programOutput)
}
