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

	multierror "github.com/hashicorp/go-multierror"

	"istio.io/mixer/adapter/statsd/config"
	"istio.io/mixer/pkg/adapter"
	"istio.io/mixer/pkg/pool"
)

type (
	builder struct {
		adapter.DefaultBuilder
	}

	aspect struct {
		clientState *clientState
	}
)

var (
	name        = "service_control_metrics"
	desc        = "Pushes metrics to service controller"
	defaultConf = &config.Params{
		Address: "chemistprober.googleprod.com",
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

	//TODO error handling
	cs, err := createAPIClient(params.Address)

	return &aspect{cs}, err
}

func (a *aspect) Record(values []adapter.Value) error {
	var result *multierror.Error
	for _, v := range values {
		if err := a.record(v); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result.ErrorOrNil()
}

func (a *aspect) record(value adapter.Value) error {
	//TODO mapping metrics to chemist proto
	mname := value.Definition.Name
	if t, found := a.templates[mname]; found {
		buf := pool.GetBuffer()

		// We don't check the error here because Execute should only fail when the template is invalid; since
		// we check that the templates are parsable in ValidateConfig and further check that they can be executed
		// with the metric's labels in NewMetricsAspect, this should never fail.
		_ = t.Execute(buf, value.Labels)
		mname = buf.String()
		pool.PutBuffer(buf)
	}

	var result error
	switch value.Definition.Kind {
	case adapter.Gauge:
		v, err := value.Int64()
		if err != nil {
			return fmt.Errorf("could not record gauge '%s': %v", mname, err)
		}
		result = a.client.Gauge(mname, v, a.rate)

	case adapter.Counter:
		v, err := value.Int64()
		if err != nil {
			return fmt.Errorf("could not record counter '%s': %v", mname, err)
		}
		result = a.client.Inc(mname, v, a.rate)
	case adapter.Distribution:
		// TODO: figure out how to program histograms via config.*
		// updates
		v, err := value.Duration()
		if err == nil {
			return a.client.TimingDuration(mname, v, a.rate)
		}
		// TODO: figure out support for non-duration distributions.
		vint, err := value.Int64()
		if err == nil {
			return a.client.Inc(mname, vint, a.rate)
		}
		return fmt.Errorf("could not record distribution '%s': %v", mname, err)
	}

	return result
}

func (a *aspect) Close() error {
	return deleteAPIClient(a.clientState)
}
