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

// Package stdioLogger provides an implementation of Mixer's logger aspect that
// writes logs (serialized as JSON) to a standard stream (stdout | stderr).
package serviceControl

import (
	"encoding/json"
	"io"

	servicecontrol "google.golang.org/api/servicecontrol/v1"
	"istio.io/mixer/adapter/serviceControl/config"
	"istio.io/mixer/pkg/adapter"
)

type (
	builder struct{ adapter.DefaultBuilder }

	logger struct {
		serviceName string
		service     *servicecontrol.Service
	}
)

// Register records the builders exposed by this adapter.
func Register(r adapter.Registrar) {
	b := builder{adapter.NewDefaultBuilder(
		"serviceControl",
		"Writes log entries to service controller",
		&config.Params{
			ServiceName:      "xiaolan-library-example.sandbox.googleapis.com",
			ClientCredentialPath: "/Users/xiaolan/credentials/",
		},
	)}

	r.RegisterApplicationLogsBuilder(b)
	r.RegisterAccessLogsBuilder(b)
}

func (builder) NewApplicationLogsAspect(env adapter.Env, cfg adapter.Config) (adapter.ApplicationLogsAspect, error) {
	return newLogger(cfg)
}

func (builder) NewAccessLogsAspect(env adapter.Env, cfg adapter.Config) (adapter.AccessLogsAspect, error) {
	return newLogger(cfg)
}

func newLogger(cfg adapter.Config) (*logger, error) {
	params := cfg.(*config.Params)

	ss, err := createAPIClient(env.Logger(), params.ClientCredentialPath)

	return &logger{params.ServiceName, ss}, err
}

func (l *logger) Log(entries []adapter.LogEntry) error {
	fmt.Printf("service control adaptor got log entriess: %v\n", entries)
	var ls []*servicecontrol.LogEntries
	for _, e := range entries {
		l := &servicecontrol.LogEntry {
			Name: e.LogName,
			Labels: e.Labels,
			Severity: e.Severity.String(),
			TextPayload: e.TextPlayload,
			StructPayload: e.StructPayload,
			Timestamp: e.Timestamp,
		}
		ls = append(ls, l)
	}

	op := &servicecontrol.Operation{
		ConsumerId:      "project:xiaolan-api-codelab",
		OperationId:     fmt.Sprintf("mixer-test-report-id-%d", rand.Int()), // TODO use uuid
		OperationName:   "reportLogs",
		StartTime:       time.Now().Format(time.RFC3339),
		EndTime:         time.Now().Format(time.RFC3339),
		LogEntries: ls,
		Labels:          map[string]string{"cloud.googleapis.com/location": "global"},
	}

	rq := &servicecontrol.ReportRequest{
		Operations: []*servicecontrol.Operation{op},
	}

	fmt.Printf("service control metric request: %v\n", len(rq.Operations[0].MetricValueSets))

	rp, err := l.service.Services.Report(l.serviceName, rq).Do()
	fmt.Printf("service control metric response for operation id %s: %v", op.OperationId, rp)
	return err
}

func (l *logger) LogAccess(entries []adapter.LogEntry) error {
	// call check api?
	return nil
}

func (l *logger) Close() error { return nil }

func writeJSON(w io.Writer, le interface{}) error {
	return json.NewEncoder(w).Encode(le)
}
