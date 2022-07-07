package telemetry

import (
	"context"
	"fmt"
	"log"
	"time"

	stdout "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

func RunMetricsCollection(ctx context.Context) error {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatalln("failed to initialize metric stdout exporter:", err)
	}
	cont := controller.New(
		processor.NewFactory(
			simple.NewWithInexpensiveDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(3*time.Second),
	)
	if err = cont.Start(context.Background()); err != nil {
		log.Fatalln("failed to start the metric controller:", err)
	}
	global.SetMeterProvider(cont)

	// Собираем стандартные метрики рантайма -
	// стандартная библиотека runtime.
	// Пример - кол-во горутин: https://pkg.go.dev/runtime#NumGoroutine
	if err = runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(time.Second),
	); err != nil {
		log.Fatalln("failed to start runtime instrumentation:", err)
	}

	<-ctx.Done()

	if err = cont.Stop(context.Background()); err != nil {
		return fmt.Errorf("failed to stop the metric controller: %v", err)
	}

	return nil
}
