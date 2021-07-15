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

type compInputs struct {
	comps  []string
	inputs []string
}

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

// returns componentNameWizard() output, a list of all components
func buildNameWizard(indent int, prefix string, compType string, compNames []string) string {
	const tabSize = 4
	space := indent * tabSize
	tab := strings.Repeat(" ", space)
	expected := fmt.Sprintf("%sAdd %s (enter to skip)\n", tab, compType)
	for i := 0; i < len(compNames); i++ {
		expected += fmt.Sprintf("%s%d: %s\n", tab, i, compNames[i])
	}
	expected += tab + "> "
	return prefix + expected
}

// returns componentListWizard() output
func buildListWizard(indent int, compGroup string, compNames []string, inputs []string) string {
	const tabSize = 4
	space := indent * tabSize
	tab := strings.Repeat(" ", space)
	expected := fmt.Sprintf("%sCurrent %ss: []\n", tab, compGroup)
	if len(inputs) == 0 {
		return buildNameWizard(1, expected, compGroup, compNames)
	}
	expected = buildNameWizard(1, expected, compGroup, compNames)
	for counter := 1; counter <= len(inputs); counter++ {
		theComp := strings.Split(inputs[counter-1], "/")[0]
		expected += fmt.Sprintf("%s%s %s extended name (optional) > ", tab, theComp, compGroup)
		expected += fmt.Sprintf("%sCurrent %ss: [", tab, compGroup)
		names := inputs[0:counter]
		for _, name := range names {
			expected += name + ", "
		}
		expected = expected[0 : len(expected)-2]
		expected += "]\n"
		expected = buildNameWizard(1, expected, compGroup, compNames)
	}
	return expected
}

// returns RpeWizard() output
func buildRpeWizard(
	recInput compInputs,
	procInput compInputs,
	expInput compInputs,
	extInput compInputs,
) (string, rpe) {
	expected := buildListWizard(1, "receiver", recInput.comps, recInput.inputs)
	expected += buildListWizard(1, "processor", procInput.comps, procInput.inputs)
	expected += buildListWizard(1, "exporter", expInput.comps, expInput.inputs)
	expected += buildListWizard(1, "extension", extInput.comps, extInput.inputs)
	expectedRPE := rpe{
		Receivers:  recInput.inputs,
		Processors: procInput.inputs,
		Exporters:  expInput.inputs,
		Extensions: extInput.inputs,
	}
	return expected, expectedRPE
}

// returns pipelineTypeWizard() output
func buildPipelineType(
	name string,
	recInput compInputs,
	procInput compInputs,
	expInput compInputs,
	extInput compInputs,
) (string, rpe) {
	wiz, rpe0 := buildRpeWizard(recInput, procInput, expInput, extInput)
	expected := fmt.Sprintf("%s pipeline extended name (optional) > ", strings.Split(strings.Title(name), "/")[0])
	expected += fmt.Sprintf("Pipeline \"%s\"\n", name)
	expected += wiz
	return expected, rpe0
}

func TestPipelineTypeWizard(t *testing.T) {
	testRecs := []string{"rec1", "rec2", "rec3"}
	testProcs := []string{"proc1", "proc2"}
	testExps := []string{"exps1", "exps2", "exps3", "exps4"}
	testExts := []string{"exts1", "exts2", "exts3"}

	// Test with no inputs
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{""}}
	io := clio{w.write, r.read}
	name, rpeOut := pipelineTypeWizard(io, "testing", testRecs, testProcs, testExps, testExts)
	expected0, rpe0 := buildPipelineType(
		name,
		compInputs{comps: testRecs},
		compInputs{comps: testProcs},
		compInputs{comps: testExps},
		compInputs{comps: testExts},
	)
	assert.Equal(t, "testing", name)
	assert.Equal(t, rpe0, rpeOut)
	assert.Equal(t, expected0, w.programOutput)

	// Testing w/ basic inputs
	w1 := fakeWriter{}
	r1 := fakeReaderPipe{userInput: []string{"", "0", "", "", "0", "", "", "0", "", "", "0", ""}}
	io1 := clio{w1.write, r1.read}
	name1, rpeOut1 := pipelineTypeWizard(io1, "testing1", testRecs, testProcs, testExps, testExts)
	expected1, rpe1 := buildPipelineType(
		name1,
		compInputs{comps: testRecs, inputs: []string{testRecs[0]}},
		compInputs{comps: testProcs, inputs: []string{testProcs[0]}},
		compInputs{comps: testExps, inputs: []string{testExps[0]}},
		compInputs{comps: testExts, inputs: []string{testExts[0]}},
	)
	assert.Equal(t, "testing1", name1)
	assert.Equal(t, rpe1, rpeOut1)
	assert.Equal(t, expected1, w1.programOutput)

	// Testing w/ basic inputs and extra names
	w2 := fakeWriter{}
	r2 := fakeReaderPipe{userInput: []string{"extpip", "0", "extr", "", "0", "extp", "", "0", "extexp", "", "0", "extext", ""}}
	io2 := clio{w2.write, r2.read}
	name2, rpeOut2 := pipelineTypeWizard(io2, "testingExt", testRecs, testProcs, testExps, testExts)
	expected2, rpe2 := buildPipelineType(
		name2,
		compInputs{comps: testRecs, inputs: []string{testRecs[0] + "/extr"}},
		compInputs{comps: testProcs, inputs: []string{testProcs[0] + "/extp"}},
		compInputs{comps: testExps, inputs: []string{testExps[0] + "/extexp"}},
		compInputs{comps: testExts, inputs: []string{testExts[0] + "/extext"}},
	)
	assert.Equal(t, "testingExt"+"/extpip", name2)
	assert.Equal(t, rpe2, rpeOut2)
	assert.Equal(t, expected2, w2.programOutput)
}

func TestRpeWizard(t *testing.T) {
	testRecs := []string{"rec1", "rec2", "rec3"}
	testProcs := []string{"proc1", "proc2"}
	testExps := []string{"exps1", "exps2", "exps3", "exps4"}
	testExts := []string{"exts1", "exts2", "exts3"}

	// Test with no inputs
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{""}}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	out := rpeWizard(io, pr, testRecs, testProcs, testExps, testExts)
	expected, expectedOut := buildRpeWizard(
		compInputs{comps: testRecs},
		compInputs{comps: testProcs},
		compInputs{comps: testExps},
		compInputs{comps: testExts},
	)
	assert.Equal(t, expectedOut, out)
	assert.Equal(t, expected, w.programOutput)

	// Test w/ user input for pipeline
	w2 := fakeWriter{}
	r2 := fakeReaderPipe{userInput: []string{"0", "", "", "0", "", "", "0", "", "", "0", ""}}
	io2 := clio{w2.write, r2.read}
	pr2 := io2.newIndentingPrinter(1)
	out2 := rpeWizard(io2, pr2, testRecs, testProcs, testExps, testExts)
	expected2, expectedOut2 := buildRpeWizard(
		compInputs{comps: testRecs, inputs: []string{testRecs[0]}},
		compInputs{comps: testProcs, inputs: []string{testProcs[0]}},
		compInputs{comps: testExps, inputs: []string{testExps[0]}},
		compInputs{comps: testExts, inputs: []string{testExts[0]}},
	)
	assert.Equal(t, expectedOut2, out2)
	assert.Equal(t, expected2, w2.programOutput)

	// Test w/ user input for pipeline w/ extended names
	w3 := fakeWriter{}
	r3 := fakeReaderPipe{userInput: []string{"0", "extr", "", "0", "extp", "", "0", "extexp", "", "0", "extext", ""}}
	io3 := clio{w3.write, r3.read}
	pr3 := io3.newIndentingPrinter(1)
	out3 := rpeWizard(io3, pr3, testRecs, testProcs, testExps, testExts)
	expected3, expectedOut3 := buildRpeWizard(
		compInputs{comps: testRecs, inputs: []string{testRecs[0] + "/extr"}},
		compInputs{comps: testProcs, inputs: []string{testProcs[0] + "/extp"}},
		compInputs{comps: testExps, inputs: []string{testExps[0] + "/extexp"}},
		compInputs{comps: testExts, inputs: []string{testExts[0] + "/extext"}},
	)

	assert.Equal(t, expectedOut3, out3)
	assert.Equal(t, expected3, w3.programOutput)

	// Test w/ mixture of inputs
	w4 := fakeWriter{}
	r4 := fakeReaderPipe{userInput: []string{"0", "", "1", "", "", "0", "", "", "1", "", "", "0", ""}}
	io4 := clio{w4.write, r4.read}
	pr4 := io4.newIndentingPrinter(1)
	out4 := rpeWizard(io4, pr4, testRecs, testProcs, testExps, testExts)

	expected4, expectedOut4 := buildRpeWizard(
		compInputs{comps: testRecs, inputs: []string{testRecs[0], testRecs[1]}},
		compInputs{comps: testProcs, inputs: []string{testProcs[0]}},
		compInputs{comps: testExps, inputs: []string{testExps[1]}},
		compInputs{comps: testExts, inputs: []string{testExts[0]}},
	)
	assert.Equal(t, expectedOut4, out4)
	assert.Equal(t, expected4, w4.programOutput)
}

func TestComponentListWizard(t *testing.T) {
	// Test with no extra inputs
	w := fakeWriter{}
	r := fakeReader{}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	compGroup := "test"
	compNames := []string{"comp1", "comp2", "comp3"}
	componentListWizard(io, pr, compGroup, compNames)
	expected := buildListWizard(1, compGroup, compNames, []string{})
	assert.Equal(t, expected, w.programOutput)

	// Testing with an input
	w2 := fakeWriter{}
	r2 := fakeReaderPipe{userInput: []string{"0", ""}, input: 0}
	io2 := clio{w2.write, r2.read}
	pr2 := io2.newIndentingPrinter(1)
	componentListWizard(io2, pr2, compGroup, compNames)
	expected2 := buildListWizard(
		1,
		compGroup,
		compNames,
		[]string{compNames[0]},
	)
	assert.Equal(t, expected2, w2.programOutput)

	// Testing extension and the input 1+ component
	w3 := fakeWriter{}
	r3 := fakeReaderPipe{userInput: []string{"0", "extension", "1", "", ""}, input: 0}
	io3 := clio{w3.write, r3.read}
	pr3 := io3.newIndentingPrinter(1)
	componentListWizard(io3, pr3, compGroup, compNames)
	expected3 := buildListWizard(
		1,
		compGroup,
		compNames,
		[]string{compNames[0] + "/extension", compNames[1]},
	)
	assert.Equal(t, expected3, w3.programOutput)

	// Extra test for buildListWizard on 2+
	w4 := fakeWriter{}
	r4 := fakeReaderPipe{userInput: []string{"0", "", "1", "extension", "2", "", ""}}
	io4 := clio{w4.write, r4.read}
	pr4 := io4.newIndentingPrinter(1)
	componentListWizard(io4, pr4, compGroup, compNames)
	expected4 := buildListWizard(
		1,
		compGroup,
		compNames,
		[]string{compNames[0], compNames[1] + "/extension", compNames[2]},
	)
	assert.Equal(t, expected4, w4.programOutput)

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
	expected := buildNameWizard(1, "", compType, compNames)
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
