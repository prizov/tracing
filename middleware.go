package tracing

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

func TraceHandler(handler http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		tracer := opentracing.GlobalTracer()
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))
		span := tracer.StartSpan(request.Method + " " + request.URL.Host, ext.RPCServerOption(spanCtx))
		ctx := opentracing.ContextWithSpan(request.Context(), span)
		request = request.WithContext(ctx)
		defer span.Finish()
		handler.ServeHTTP(writer, request)
	}
}
