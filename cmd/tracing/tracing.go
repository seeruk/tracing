package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"

	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func main() {
	jaegerAgentHost := os.Getenv("JAEGER_AGENT_HOST")
	jaegerAgentPort := os.Getenv("JAEGER_AGENT_PORT")

	if jaegerAgentHost != "" && jaegerAgentPort != "" {
		cfg := jaegercfg.Configuration{
			Sampler: &jaegercfg.SamplerConfig{
				Type:  "const",
				Param: 1,
			},
			Reporter: &jaegercfg.ReporterConfig{
				LogSpans:            true,
				BufferFlushInterval: time.Second,
				LocalAgentHostPort:  fmt.Sprintf("%s:%s", jaegerAgentHost, jaegerAgentPort),
			},
		}

		tracer, _, err := cfg.New("seeruk-tracing", jaegercfg.Logger(jaeger.StdLogger))
		if err != nil {
			log.Fatalf("failed to initiate tracer: %v\n", err)
		}

		log.Printf("Set tracer: %s\n", fmt.Sprintf("%s:%s", jaegerAgentHost, jaegerAgentPort))

		opentracing.SetGlobalTracer(tracer)
	}

	span := opentracing.StartSpan("cmd/tracing:main")

	time.Sleep(500 * time.Millisecond)

	span.Finish()

	// Actually allow the spans to be flushed...
	time.Sleep(1.5 * time.Second)
}
