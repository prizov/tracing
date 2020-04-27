package tracing

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"net/http"
)

type JaegerTransport struct{}

func (JaegerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()
	host := r.URL.Hostname()
	method := r.Method
	url := r.URL.String()
	operationName := method + ":" + host
	span, _ := opentracing.StartSpanFromContext(ctx, operationName)
	if span == nil {
		span = opentracing.StartSpan(r.Method + "" + r.URL.Host)
		fmt.Println("Host", r.URL.Host)
	}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, r.Method)
	span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	res, err :=  http.DefaultTransport.RoundTrip(r)
	if err == nil {
		span.SetTag("http.status.code", res.StatusCode)
	}
	span.Finish()
	return res, err
}

func NewHTTPClient() *http.Client  {
	return &http.Client{Transport: JaegerTransport{}}
}