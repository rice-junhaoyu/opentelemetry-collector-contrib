// Copyright 2019, OpenTelemetry Authors
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

package googlecloudexporter

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configtest"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, configtest.CheckConfigStruct(cfg))
}

func TestCreateExporter(t *testing.T) {
	defer setPdataFeatureGateForTest(t, true)()
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("Default credentials not set, skip creating Google Cloud exporter")
	}
	ctx := context.Background()
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	eCfg := cfg.(*Config)
	eCfg.ProjectID = "test"

	te, err := factory.CreateTracesExporter(ctx, componenttest.NewNopExporterCreateSettings(), eCfg)
	assert.NoError(t, err)
	assert.NotNil(t, te, "failed to create trace exporter")

	me, err := factory.CreateMetricsExporter(ctx, componenttest.NewNopExporterCreateSettings(), eCfg)
	assert.NoError(t, err)
	assert.NotNil(t, me, "failed to create metrics exporter")
}

func TestCreateLegacyExporter(t *testing.T) {
	defer setPdataFeatureGateForTest(t, false)()
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		t.Skip("Default credentials not set, skip creating Google Cloud exporter")
	}
	ctx := context.Background()
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	eCfg := cfg.(*Config)
	eCfg.ProjectID = "test"

	te, err := factory.CreateTracesExporter(ctx, componenttest.NewNopExporterCreateSettings(), eCfg)
	assert.NoError(t, err)
	assert.NotNil(t, te, "failed to create trace exporter")

	me, err := factory.CreateMetricsExporter(ctx, componenttest.NewNopExporterCreateSettings(), eCfg)
	assert.NoError(t, err)
	assert.NotNil(t, me, "failed to create metrics exporter")
}
