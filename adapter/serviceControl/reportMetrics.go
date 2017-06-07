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

package serviceControl

import (
	"fmt"
	"math/rand"
	"time"

	"istio.io/mixer/adapter/serviceControl/config"
	"istio.io/mixer/pkg/adapter"
)

type (
	builder struct {
		adapter.DefaultBuilder
	}

	aspect struct {
		serviceName string
		service     *servicecontrol.Service
	}
)

var (
	name        = "service_control_metrics"
	desc        = "Pushes metrics to service controller"
	defaultConf = &config.Params{
		ClientId:     "mixc",
		ServiceName:  "chemistprober.googleprod.com",
		ClientSecret: "",
		Scope:        "",
		TokenFile:    "",
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

	return &aspect{params.ServiceName, ss}, err
}

func (a *aspect) Record(values []adapter.Value) error {
	vs := make([]*servicecontrol.MetricValueSet)
	for _, v := range values {
		var mv servicecontrol.MetricValue
		mv.Labels = mapLabels(v.Labels)
		mv.StatTime = v.StartTime
		mv.EndTime = v.EndTime
		i, _ := v.Int64()
		mv.Int64Value = &i

		ms := &servicecontrol.MericValueSet{
			MetricName:   v.Definition.Name,
			MetricValues: []*servicecontrol.MetricValue{&mv},
		}
		vs.append(ms)
	}

	op := &servicecontrol.Operation{
		OperationId:     fmt.Sprintf("%d", rand.Int()), // TODO use uuid
		OpeationName:    "reportMetrics",
		StartTime:       fmt.Sprintf("%d", time.Now()),
		EndTime:         fmt.Sprintf("%d", time.Now()),
		MetricValueSets: vs,
	}
	rq := &servicecontrol.ReportRequest{
		Operations: []*servicecontrol.Operation{op},
	}
	rp, err := a.service.Services.Report(a.serviceName, rq).Do()
	fmt.Printf("service control metric response for operation id %s: %v", op.OperationId, rp)
	return err
}

func mapLabels(labels map[string]interface{}) map[string]string {
	ml := make(map[string]string)
	for k, v := range labels {
		ml[k] = fmt.Sprintf("%v", v)
	}
	return ml
}

func (a *aspect) record(value adapter.Value) error {
	//TODO do not use
	return nil
}

func (a *aspect) Close() error {
	//TODO doesn't need?
	return nil
}
