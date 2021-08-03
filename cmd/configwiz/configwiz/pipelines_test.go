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
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/extension/extensionhelper"

	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.opentelemetry.io/collector/receiver/receiverhelper"

	testcomponents2 "github.com/open-telemetry/opentelemetry-collector-contrib/internal/testcomponents"
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

var testRecs = []string{"rec1", "rec2", "rec3"}
var testProcs = []string{"proc1", "proc2"}
var testExps = []string{"exps1", "exps2", "exps3", "exps4"}
var testExts = []string{"exts1", "exts2", "exts3"}
var compNames = []string{"comp1", "comp2", "comp3"}
const compType = "test"

// PipelineWizard() Tests
func TestPipelineWizardTraces (t *testing.T){
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"2", ""}}
	io := clio{w.write, r.read}
	testFact := createTestFactories()
	out := pipelinesWizard(io, testFact)
	expected, _ := buildPipelineWizard(testFact, []string{"traces"})
	assert.Equal(t, map[string]interface{}(map[string]interface{}{
		"traces": rpe{},
	}), out)
	assert.Equal(t, expected, w.programOutput)
}

func TestPipelineWizardMetric(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"1", ""}}
	io := clio{w.write, r.read}
	testFact := createTestFactories()
	out := pipelinesWizard(io, testFact)
	expected, _ := buildPipelineWizard(testFact, []string{"metrics"})
	assert.Equal(t, map[string]interface{}(map[string]interface{}{
		"metrics": rpe{},
	}), out)
	assert.Equal(t, expected, w.programOutput)
}

func TestPipelineWizardEmpty(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{""}}
	io := clio{w.write, r.read}
	testFact := createTestFactories()
	out := pipelinesWizard(io, testFact)
	expected, _ := buildPipelineWizard(testFact, []string{})
	assert.Equal(t, map[string]interface{}(map[string]interface{}{}), out)
	assert.Equal(t, expected, w.programOutput)
}

//singlePipelineWizardTests
func TestSinglePipelineWizardFail(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"-1", ""}}
	io := clio{w.write, r.read}
	testFact := createTestFactories()
	_, rpeOut := singlePipelineWizard(io, testFact)
	expected := "Add pipeline (enter to skip)\n1: Metrics\n2: Traces\n> "
	expected += "Invalid input. Try again.\n" + expected
	assert.Equal(t, rpe{}, rpeOut)
	assert.Equal(t, expected, w.programOutput)
}

func TestSinglePipelineWizardEmpty(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{""}}
	io := clio{w.write, r.read}
	testFact := createTestFactories()
	_, rpeOut := singlePipelineWizard(io, testFact)
	expected := "Add pipeline (enter to skip)\n1: Metrics\n2: Traces\n> "
	assert.Equal(t, rpe{}, rpeOut)
	assert.Equal(t, expected, w.programOutput)
}

func TestSinglePipelineWizardTraces(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"2", ""}}
	io := clio{w.write, r.read}
	testFact := createTestFactories()
	name, rpeOut := singlePipelineWizard(io, testFact)
	expectedOut, rpe0 := buildSinglePipelineWiz(testFact, name)
	assert.Equal(t, rpe0, rpeOut)
	assert.Equal(t, expectedOut, w.programOutput)
}

func TestSinglePipelineWizardMetrics(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"1", ""}}
	io := clio{w.write, r.read}
	testFact := createTestFactories()
	name, rpeOut := singlePipelineWizard(io, testFact)
	expectedOut, rpe0 := buildSinglePipelineWiz(testFact, name)
	assert.Equal(t, rpe0, rpeOut)
	assert.Equal(t, expectedOut, w.programOutput)
}

// pipeline wizard Tests
func TestPipelineTypeWizardEmpty(t *testing.T) {
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
}

func TestPipelineTypeWizardBasicInp(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"", "0", "", "", "0", "", "", "0", "", "", "0", ""}}
	io := clio{w.write, r.read}
	name, rpeOut := pipelineTypeWizard(io, "testing1", testRecs, testProcs, testExps, testExts)
	expected, rpe0 := buildPipelineType(
		name,
		compInputs{comps: testRecs, inputs: []string{testRecs[0]}},
		compInputs{comps: testProcs, inputs: []string{testProcs[0]}},
		compInputs{comps: testExps, inputs: []string{testExps[0]}},
		compInputs{comps: testExts, inputs: []string{testExts[0]}},
	)
	assert.Equal(t, "testing1", name)
	assert.Equal(t, rpe0, rpeOut)
	assert.Equal(t, expected, w.programOutput)
}

func TestPipelineTypeWizardExtendedNames(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"extpip", "0", "extr", "", "0", "extp", "", "0", "extexp", "", "0", "extext", ""}}
	io := clio{w.write, r.read}
	name, rpeOut := pipelineTypeWizard(io, "testingExt", testRecs, testProcs, testExps, testExts)
	expected, rpe0 := buildPipelineType(
		name,
		compInputs{comps: testRecs, inputs: []string{testRecs[0] + "/extr"}},
		compInputs{comps: testProcs, inputs: []string{testProcs[0] + "/extp"}},
		compInputs{comps: testExps, inputs: []string{testExps[0] + "/extexp"}},
		compInputs{comps: testExts, inputs: []string{testExts[0] + "/extext"}},
	)
	assert.Equal(t, "testingExt"+"/extpip", name)
	assert.Equal(t, rpe0, rpeOut)
	assert.Equal(t, expected, w.programOutput)
}

// RpeWizard tests
func TestRpeWizardEmpty(t *testing.T) {
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
}

func TestRpeWizardBasicInp(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"0", "", "", "0", "", "", "0", "", "", "0", ""}}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	out := rpeWizard(io, pr, testRecs, testProcs, testExps, testExts)
	expected, expectedOut := buildRpeWizard(
		compInputs{comps: testRecs, inputs: []string{testRecs[0]}},
		compInputs{comps: testProcs, inputs: []string{testProcs[0]}},
		compInputs{comps: testExps, inputs: []string{testExps[0]}},
		compInputs{comps: testExts, inputs: []string{testExts[0]}},
	)
	assert.Equal(t, expectedOut, out)
	assert.Equal(t, expected, w.programOutput)
}

func TestRpeWizardMultipleInputs(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"0", "", "1", "extr", "", "0", "", "", "1", "", "", "0", ""}}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	out := rpeWizard(io, pr, testRecs, testProcs, testExps, testExts)
	expected, expectedOut := buildRpeWizard(
		compInputs{comps: testRecs, inputs: []string{testRecs[0], testRecs[1] + "/extr"}},
		compInputs{comps: testProcs, inputs: []string{testProcs[0]}},
		compInputs{comps: testExps, inputs: []string{testExps[1]}},
		compInputs{comps: testExts, inputs: []string{testExts[0]}},
	)
	assert.Equal(t, expectedOut, out)
	assert.Equal(t, expected, w.programOutput)
}

//ComponentListWizardTest()
func TestComponentListWizardEmpty(t *testing.T) {
	w := fakeWriter{}
	r := fakeReader{}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	componentListWizard(io, pr, compType, compNames)
	expected := buildListWizard(1, compType, compNames, []string{})
	assert.Equal(t, expected, w.programOutput)
}

func TestComponentListWizardSingleInp(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"0", ""}, input: 0}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	componentListWizard(io, pr, compType, compNames)
	expected := buildListWizard(
		1,
		compType,
		compNames,
		[]string{compNames[0]},
	)
	assert.Equal(t, expected, w.programOutput)
}

func TestComponentListWizardMultipleInp(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{userInput: []string{"0", "", "1", "extension", "2", "", ""}}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	componentListWizard(io, pr, compType, compNames)
	expected := buildListWizard(
		1,
		compType,
		compNames,
		[]string{compNames[0], compNames[1] + "/extension", compNames[2]},
	)
	assert.Equal(t, expected, w.programOutput)
}

// Test ComponentNameWizard()
func TestComponentNameWizardEmpty(t *testing.T) {
	w := fakeWriter{}
	r := fakeReader{}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	componentNameWizard(io, pr, compType, compNames)
	expected := buildNameWizard(1, "", compType, compNames)
	assert.Equal(t, expected, w.programOutput)
}

func TestComponentNameWizardExtended(t *testing.T) {
	w := fakeWriter{}
	r := fakeReader{"0"}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	out, val := componentNameWizard(io, pr, compType, compNames)
	expected := buildNameWizard(1, "", compType, compNames)
	tab := strings.Repeat(" ", 4)
	expected += fmt.Sprintf("%s%s %s extended name (optional) > ", tab, out, compType)
	assert.Equal(t, compNames[0], out)
	assert.Equal(t, val, "0")
	assert.Equal(t, expected, w.programOutput)
}


func TestComponentNameWizardError(t *testing.T) {
	w := fakeWriter{}
	r := fakeReaderPipe{[]string{"-1", ""}, 0}
	io := clio{w.write, r.read}
	pr := io.newIndentingPrinter(1)
	componentNameWizard(io, pr, compType, compNames)
	expected := buildNameWizard(1, "", compType, compNames)
	expected += "Invalid input. Try again.\n"
	expected += buildNameWizard(1, "", compType, compNames)
	assert.Equal(t, expected, w.programOutput)
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

func buildCompInputs(testFactory component.Factories, pipeType bool, recInp []string, procInp []string, expInp []string, extInp []string) []compInputs {
	if pipeType {
		return []compInputs{
			{receiverNames(testFactory, isMetricsReceiver), recInp},
			{processorNames(testFactory, isMetricProcessor), procInp},
			{exporterNames(testFactory, isMetricsExporter), expInp},
			{extensionNames(testFactory, isExtension), extInp},
		}
	}
	return []compInputs{
		{receiverNames(testFactory, isTracesReceiver), recInp},
		{processorNames(testFactory, isTracesProcessor), procInp},
		{exporterNames(testFactory, isTracesExporter), expInp},
		{extensionNames(testFactory, isExtension), extInp},
	}
}

func buildSinglePipelineWiz(testFact component.Factories, name string) (string, rpe) {
	expected := "Add pipeline (enter to skip)\n1: Metrics\n2: Traces\n> "
	comps := buildCompInputs(testFact, false, nil, nil, nil, nil)
	expectedOut, rpe0 := buildPipelineType(name, comps[0], comps[1], comps[2], comps[3])
	return expected + expectedOut, rpe0

}

func buildPipelineWizard(testFact component.Factories, inputs []string) (string, rpe) {
	expected := "Current pipelines: []\n"
	addPipe := "Add pipeline (enter to skip)\n1: Metrics\n2: Traces\n> "
	if len(inputs) == 0 {
		return expected + addPipe, rpe{}
	}
	expectedOut, rpe0 := buildSinglePipelineWiz(testFact, inputs[0])
	expectedOut = expected + expectedOut
	for i := range inputs {
		expectedOut += "Current pipelines: ["
		currPipes := inputs[:i+1]
		for _, pipe := range currPipes {
			expectedOut += pipe + ", "
		}
		expectedOut = expectedOut[0 : len(expectedOut)-2]
		expectedOut += "]\n" + addPipe
	}
	return expectedOut, rpe0
}


func createTestFactories() component.Factories {
	exampleReceiverFactory := testcomponents2.ExampleReceiverFactory
	exampleProcessorFactory := testcomponents2.ExampleProcessorFactory
	exampleExporterFactory := testcomponents2.ExampleExporterFactory
	badExtensionFactory := newBadExtensionFactory()
	badReceiverFactory := newBadReceiverFactory()
	badProcessorFactory := newBadProcessorFactory()
	badExporterFactory := newBadExporterFactory()

	factories := component.Factories{
		Extensions: map[config.Type]component.ExtensionFactory{
			badExtensionFactory.Type(): badExtensionFactory,
		},
		Receivers: map[config.Type]component.ReceiverFactory{
			exampleReceiverFactory.Type(): exampleReceiverFactory,
			badReceiverFactory.Type():     badReceiverFactory,
		},
		Processors: map[config.Type]component.ProcessorFactory{
			exampleProcessorFactory.Type(): exampleProcessorFactory,
			badProcessorFactory.Type():     badProcessorFactory,
		},
		Exporters: map[config.Type]component.ExporterFactory{
			exampleExporterFactory.Type(): exampleExporterFactory,
			badExporterFactory.Type():     badExporterFactory,
		},
	}

	return factories
}

func newBadReceiverFactory() component.ReceiverFactory {
	return receiverhelper.NewFactory("bf", func() config.Receiver {
		return &struct {
			config.ReceiverSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
		}{
			ReceiverSettings: config.NewReceiverSettings(config.NewID("bf")),
		}
	})
}

func newBadProcessorFactory() component.ProcessorFactory {
	return processorhelper.NewFactory("bf", func() config.Processor {
		return &struct {
			config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
		}{
			ProcessorSettings: config.NewProcessorSettings(config.NewID("bf")),
		}
	})
}

func newBadExporterFactory() component.ExporterFactory {
	return exporterhelper.NewFactory("bf", func() config.Exporter {
		return &struct {
			config.ExporterSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
		}{
			ExporterSettings: config.NewExporterSettings(config.NewID("bf")),
		}
	})
}

func newBadExtensionFactory() component.ExtensionFactory {
	return extensionhelper.NewFactory(
		"bf",
		func() config.Extension {
			return &struct {
				config.ExtensionSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
			}{
				ExtensionSettings: config.NewExtensionSettings(config.NewID("bf")),
			}
		},
		func(ctx context.Context, params component.ExtensionCreateSettings, extension config.Extension) (component.Extension, error) {
			return nil, nil
		},
	)
}
