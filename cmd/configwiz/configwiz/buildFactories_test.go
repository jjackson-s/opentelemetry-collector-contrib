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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/components"
)



func TestComps(t *testing.T) {
	out, _ := components.Components()
	//out := k8sprocessor.NewFactory()
	//out, _ := components()
	assert.Equal(t, "", out)
}

//func TestFindComps(t *testing.T) {
//	//"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver"
//
//	otelContrib := "opentelemetry-collector-contrib"
//	w := fakeWriter{}
//	r := fakeReader{}
//	io := clio{w.write, r.read}
//	out := findComps(io)
//	expectedOut := []string{
//		path.Join(otelContrib, "/receiver"),
//		path.Join(otelContrib, "/processor"),
//		path.Join(otelContrib, "/exporter"),
//		path.Join(otelContrib, "/extension"),
//	}
//	assert.Equal(t, expectedOut, out)
//}
