package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

func main() {
	// jaeger config
	cfg := jaegercfg.Configuration{
		ServiceName: "test-app",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	// end of jaeger setting

	r := gin.Default()
	r.GET("/hi", func(ctx *gin.Context) {
		// create the span
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
		serverSpan := tracer.StartSpan("server say hi", ext.RPCServerOption(spanCtx))
		// this finish() is critical
		defer serverSpan.Finish()
		// end

		ctx.JSON(http.StatusOK, gin.H{"msg": "Up and Running..."})
	})
	r.Run(":8080")
}
