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
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/cmd/configschema/configschema"
	"go.opentelemetry.io/collector/component"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/redisreceiver"
)

func TestGetComponents(t *testing.T) {
	m := map[string]interface{}{
		"pipelines": map[string]interface{}{
			"metrics": rpe{
				Exporters: []string{"ccc", "bbb", "aaa"},
			},
		},
	}
	require.Equal(t, map[string][]string{
		"exporter": {"ccc", "bbb", "aaa"},
	}, serviceToComponentNames(m))
}

func TestCli(t *testing.T) {
	factories := component.Factories{}
	receivers := []component.ReceiverFactory{
		redisreceiver.NewFactory(),
	}
	factories.Receivers, _ = component.MakeReceiverFactoryMap(receivers...)
	compGroup := "receiver"
	name := "redis"
	cfgInfo, err := configschema.GetCfgInfo(factories, compGroup, name)
	if err != nil {
		panic(err)
	}
	dr := configschema.NewDirResolver("../../..", "github.com/open-telemetry/opentelemetry-collector-contrib")
	f, err:= configschema.ReadFields(reflect.ValueOf(cfgInfo.CfgInstance), dr)
	if err != nil {
		panic(err)
	}
	fmt.Print(f)

}

