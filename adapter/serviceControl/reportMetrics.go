// Copyright 2017 Istio Authors
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

package servicecontrol

import (
	"fmt"

	"github.com/google/google-api-go-client/servicecontrol/v1"

	"istio.io/mixer/adapter/serviceControl/config"
	"istio.io/mixer/pkg/adapter"
)

type (
	builder struct {
		adapter.DefaultBuilder
	}

	aspect struct {
		service *v1.Service
	}
)

var (
	name        = "service_control_metrics"
	desc        = "Pushes metrics to service controller"
	defaultConf = &config.Params{
		//Address: "chemistprober.googleprod.com",
		client_id: "mixc",
	}
)

// Register records the builders exposed by this adapter.
func Register(r adapter.Registrar) {
	r.RegisterMetricsBuilder(newBuilder())
}

func newBuilder() *builder {
	return &builder{adapter.NewDefaultBuilder(name, desc, defaultConf)}
}

func (b *builder) ValidateConfig(c adapter.Config) (ce *adapter.ConfigErrors) {
	return
}

func (*builder) NewMetricsAspect(env adapter.Env, cfg adapter.Config, metrics map[string]*adapter.MetricDefinition) (adapter.MetricsAspect, error) {
	params := cfg.(*config.Params)

	ss, err := createAPIClient(params.clientId, params.clientSecret, params.scope, params.tokenFile)

	return &aspect{ss}, err
}

func (a *aspect) Record(values []adapter.Value) error {
	// create proto
	// create operation
	for v := range values {
		var mv v1.MetricValue
		mv.labels = mapLabels(v.Labels, v.Definition.Labels)
		mv.StatTime = v.StartTime
		mv.EndTime = v.EndTime
		mv.value, err = v.Int64()

		ms := &v1.MericValueSet{
			metricName:   v.Definition.name,
			metricValues: []v.MetricValue{mv},
		}
		// load values into operation
	}

	a.service.report()
	return nil
}

func mapLabels(labels map[string]interface{}, labelType map[string]adapter.LabelType) (map[string]string, error) {
	ml := make(map[string]string)
	for k, v := range labels {
		if labelType[k] != adapter.String {
			return nil, fmt.Errorf("Only support string labels")
		}
		ml[k] = string(v)
	}
	return ml
}

func (a *aspect) record(value adapter.Value) error {
	//TODO do not use
	return nil
}

func (a *aspect) Close() error {
	return deleteAPIClient(a.cs)
}
