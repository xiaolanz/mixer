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

	servicecontrol "google.golang.org/api/servicecontrol/v1"

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
	name        = "serviceControl"
	desc        = "Pushes metrics to service controller"
	defaultConf = &config.Params{
		ServiceName:          "xiaolan-library-example.sandbox.googleapis.com",
		ClientCredentialPath: "/Users/xiaolan/credentials/",
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
	return nil
}

func (*builder) NewMetricsAspect(env adapter.Env, cfg adapter.Config, metrics map[string]*adapter.MetricDefinition) (adapter.MetricsAspect, error) {
	params := cfg.(*config.Params)

	ss, err := createAPIClient(env.Logger(), params.ClientCredentialPath)

	return &aspect{params.ServiceName, ss}, err
}

func (a *aspect) Record(values []adapter.Value) error {
	fmt.Printf("service control adaptor got metric values: %v\n", values)
	var vs []*servicecontrol.MetricValueSet
	for _, v := range values {
		// Only for request name.
		if v.Definition.Name != "request_count" {
			fmt.Printf("service control metric get other name: %s\n", v.Definition.Name)
			continue
		}
		var mv servicecontrol.MetricValue
		mv.Labels = fillLabels(v.Labels)
		//		mv.StartTime = v.StartTime.Format(time.RFC3339)
		//		mv.EndTime = v.EndTime.Format(time.RFC3339)
		i, _ := v.Int64()
		mv.Int64Value = &i

		ms := &servicecontrol.MetricValueSet{
			MetricName:   "serviceruntime.googleapis.com/api/producer/request_count",
			MetricValues: []*servicecontrol.MetricValue{&mv},
		}
		vs = append(vs, ms)
	}

	op := &servicecontrol.Operation{
		//	ConsumerId:      "project:xiaolan-api-codelab",
		OperationId:     fmt.Sprintf("mixer-metric-report-id-%d", rand.Int()), // TODO use uuid
		OperationName:   "reportMetrics",
		StartTime:       time.Now().Format(time.RFC3339),
		EndTime:         time.Now().Format(time.RFC3339),
		MetricValueSets: vs,
		Labels: map[string]string{"cloud.googleapis.com/location": "global",
			"serviceruntime.googleapis.com/api_method": "google.example.library.v1.LibraryService.MoveBook"},
	}
	rq := &servicecontrol.ReportRequest{
		Operations: []*servicecontrol.Operation{op},
	}

	// print out
	fmt.Printf("service control metric request: %v", len(rq.Operations[0].MetricValueSets))

	rp, err := a.service.Services.Report(a.serviceName, rq).Do()
	fmt.Printf("service control metric response for operation id %s: %v", op.OperationId, rp)
	return err
}

func fillLabels(labels map[string]interface{}) map[string]string {
	ml := make(map[string]string)
	for k, v := range labels {
		if k != "response_code" {
			continue
		}
		nk := fmt.Sprintf("/%s", k)
		ml[nk] = fmt.Sprintf("%v", v)
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
