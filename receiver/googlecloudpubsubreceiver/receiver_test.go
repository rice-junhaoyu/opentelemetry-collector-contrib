// Copyright The OpenTelemetry Authors
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

package googlecloudpubsubreceiver

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/pubsub/pstest"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/obsreport"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	pb "google.golang.org/genproto/googleapis/pubsub/v1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/googlecloudpubsubreceiver/testdata"
)

func TestStartReceiverNoSubscription(t *testing.T) {
	ctx := context.Background()
	// Start a fake server running locally.
	srv := pstest.NewServer()
	defer srv.Close()
	core, _ := observer.New(zap.WarnLevel)
	receiver := &pubsubReceiver{
		logger:    zap.New(core),
		userAgent: "test-user-agent",

		config: &Config{
			Endpoint:  srv.Addr,
			Insecure:  true,
			ProjectID: "my-project",
			TimeoutSettings: exporterhelper.TimeoutSettings{
				Timeout: 12 * time.Second,
			},
			Subscription: "projects/my-project/subscriptions/otlp",
		},
	}
	defer func() {
		assert.NoError(t, receiver.Shutdown(ctx))
	}()
	receiver.tracesConsumer = consumertest.NewNop()
	receiver.metricsConsumer = consumertest.NewNop()
	receiver.logsConsumer = consumertest.NewNop()
	// No error is thrown as the stream is handled async,
	// no locks should be kept though
	assert.NoError(t, receiver.Start(ctx, nil))
}

func TestReceiver(t *testing.T) {
	ctx := context.Background()
	// Start a fake server running locally.
	srv := pstest.NewServer()
	defer srv.Close()
	_, err := srv.GServer.CreateTopic(ctx, &pb.Topic{
		Name: "projects/my-project/topics/otlp",
	})
	assert.NoError(t, err)
	_, err = srv.GServer.CreateSubscription(ctx, &pb.Subscription{
		Topic:              "projects/my-project/topics/otlp",
		Name:               "projects/my-project/subscriptions/otlp",
		AckDeadlineSeconds: 10,
	})
	assert.NoError(t, err)

	core, _ := observer.New(zap.WarnLevel)
	params := componenttest.NewNopReceiverCreateSettings()
	traceSink := new(consumertest.TracesSink)
	metricSink := new(consumertest.MetricsSink)
	logSink := new(consumertest.LogsSink)
	receiver := &pubsubReceiver{
		logger: zap.New(core),
		obsrecv: obsreport.NewReceiver(obsreport.ReceiverSettings{
			ReceiverID:             config.NewComponentID(typeStr),
			Transport:              reportTransport,
			LongLivedCtx:           false,
			ReceiverCreateSettings: params,
		}),
		userAgent: "test-user-agent",

		config: &Config{
			Endpoint:  srv.Addr,
			Insecure:  true,
			ProjectID: "my-project",
			TimeoutSettings: exporterhelper.TimeoutSettings{
				Timeout: 1 * time.Second,
			},
			Subscription: "projects/my-project/subscriptions/otlp",
		},
		tracesConsumer:  traceSink,
		metricsConsumer: metricSink,
		logsConsumer:    logSink,
	}
	assert.NoError(t, receiver.Start(ctx, nil))

	receiver.tracesConsumer = traceSink
	receiver.metricsConsumer = metricSink
	receiver.logsConsumer = logSink
	// No error is thrown as the stream is handled async,
	// no locks should be kept though
	assert.NoError(t, receiver.Start(ctx, nil))

	time.Sleep(1 * time.Second)

	// Test an OTLP trace message
	traceSink.Reset()
	srv.Publish("projects/my-project/topics/otlp", testdata.CreateTraceExport(), map[string]string{
		"ce-type":      "org.opentelemetry.otlp.traces.v1",
		"content-type": "application/protobuf",
	})
	assert.Eventually(t, func() bool {
		return len(traceSink.AllTraces()) == 1
	}, 100*time.Second, 10*time.Millisecond)

	// Test an OTLP metric message
	metricSink.Reset()
	srv.Publish("projects/my-project/topics/otlp", testdata.CreateMetricExport(), map[string]string{
		"ce-type":      "org.opentelemetry.otlp.metrics.v1",
		"content-type": "application/protobuf",
	})
	assert.Eventually(t, func() bool {
		return len(metricSink.AllMetrics()) == 1
	}, time.Second, 10*time.Millisecond)

	// Test an OTLP log message
	logSink.Reset()
	srv.Publish("projects/my-project/topics/otlp", testdata.CreateLogExport(), map[string]string{
		"ce-type":      "org.opentelemetry.otlp.logs.v1",
		"content-type": "application/protobuf",
	})
	assert.Eventually(t, func() bool {
		return len(logSink.AllLogs()) == 1
	}, time.Second, 10*time.Millisecond)

	// Test a plain log message
	logSink.Reset()
	srv.Publish("projects/my-project/topics/otlp", testdata.CreateTextExport(), map[string]string{
		"content-type": "text/plain",
	})
	assert.Eventually(t, func() bool {
		return len(logSink.AllLogs()) == 1
	}, time.Second, 10*time.Millisecond)

	assert.Nil(t, receiver.Shutdown(ctx))
	assert.Nil(t, receiver.Shutdown(ctx))
}
