package opentracing_test

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	cli "go.unistack.org/micro-client-grpc/v3"
	jsoncodec "go.unistack.org/micro-codec-json/v3"
	rrouter "go.unistack.org/micro-router-register/v3"
	srv "go.unistack.org/micro-server-grpc/v3"
	ot "go.unistack.org/micro-tracer-opentracing/v3"
	"go.unistack.org/micro/v3/api"
	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/errors"
	"go.unistack.org/micro/v3/register"
	"go.unistack.org/micro/v3/router"
	"go.unistack.org/micro/v3/server"
	mt "go.unistack.org/micro/v3/tracer"
	otwrapper "go.unistack.org/micro/v3/tracer/wrapper"
)

type Test interface {
	Method(ctx context.Context, in *TestRequest, opts ...client.CallOption) (*TestResponse, error)
}

type TestRequest struct {
	IsError bool
}

type TestResponse struct {
	Message string
}

type testHandler struct{}

func (t *testHandler) Method(ctx context.Context, req *TestRequest, rsp *TestResponse) error {
	var span mt.Span
	ctx, span = mt.DefaultTracer.Start(ctx, "internal", mt.WithSpanKind(mt.SpanKindInternal))
	defer span.Finish()
	span.AddLabels("some key", "some val")
	if req.IsError {
		return errors.BadRequest("bad", "test error")
	}

	rsp.Message = "passed"

	return nil
}

func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		ServiceName: service,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}

func TestClient(t *testing.T) {
	// setup
	assert := assert.New(t)
	for name, tt := range map[string]struct {
		message     string
		isError     bool
		wantMessage string
		wantStatus  string
	}{
		"OK": {
			message:     "passed",
			isError:     false,
			wantMessage: "passed",
			wantStatus:  "OK",
		},
		"Invalid": {
			message:     "",
			isError:     true,
			wantMessage: "",
			wantStatus:  "InvalidArgument",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var tr opentracing.Tracer
			tr = mocktracer.New()

			_ = tr
			var cl io.Closer
			tr, cl = initJaeger(fmt.Sprintf("Test tracing %s", time.Now().Format(time.RFC1123Z)))
			defer cl.Close()
			opentracing.SetGlobalTracer(tr)

			reg := register.NewRegister()
			brk := broker.NewBroker(broker.Register(reg))

			serverName := "service"
			serverID := "1234567890"
			serverVersion := "1.0.0"

			rt := rrouter.NewRouter(router.Register(reg))

			mtr := ot.NewTracer(ot.Tracer(tr))
			if err := mtr.Init(); err != nil {
				t.Fatal(err)
			}
			mt.DefaultTracer = mtr

			c := cli.NewClient(
				client.Codec("application/grpc+json", jsoncodec.NewCodec()),
				client.Codec("application/json", jsoncodec.NewCodec()),
				client.Router(rt),
				client.Wrap(otwrapper.NewClientWrapper(otwrapper.WithTracer(mtr))),
			)

			if err := c.Init(); err != nil {
				t.Fatal(err)
			}

			s := srv.NewServer(
				server.Codec("application/grpc+json", jsoncodec.NewCodec()),
				server.Codec("application/json", jsoncodec.NewCodec()),
				server.Name(serverName),
				server.Version(serverVersion),
				server.ID(serverID),
				server.Register(reg),
				server.Broker(brk),
				server.WrapSubscriber(otwrapper.NewServerSubscriberWrapper(otwrapper.WithTracer(mtr))),
				server.WrapHandler(otwrapper.NewServerHandlerWrapper(otwrapper.WithTracer(mtr))),
				server.Address("127.0.0.1:0"),
			)
			if err := s.Init(); err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = s.Stop()
			}()

			type Test struct {
				*testHandler
			}

			nopts := []server.HandlerOption{
				api.WithEndpoint(&api.Endpoint{
					Name:    "Test.Method",
					Method:  []string{"POST"},
					Handler: "rpc",
				}),
			}

			if err := s.Handle(s.NewHandler(&Test{new(testHandler)}, nopts...)); err != nil {
				t.Fatal(err)
			}

			if err := s.Start(); err != nil {
				t.Fatalf("Unexpected error starting server: %v", err)
			}

			ctx, span := mtr.Start(context.Background(), "root", mt.WithSpanKind(mt.SpanKindClient))
			var err error
			req := c.NewRequest("service", "Test.Method", &TestRequest{IsError: tt.isError}, client.RequestContentType("application/json"))
			fmt.Printf("%s.%s\n", req.Service(), req.Endpoint())
			rsp := TestResponse{}
			err = c.Call(ctx, req, &rsp)
			if tt.isError {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			assert.Equal(rsp.Message, tt.message)

			span.Finish()

			return
			/*
				spans := tr.FinishedSpans()
				assert.Len(spans, 3)

				var rootSpan opentracing.Span
				for _, s := range spans {
					// order of traces in buffer is not garanteed
					switch s.OperationName {
					case "root":
						rootSpan = s
					}
				}

				for _, s := range spans {
					fmt.Printf("root %#+v\ncheck span %#+v\n", rootSpan, s)
					assert.Equal(rootSpan.Context().(mocktracer.MockSpanContext).TraceID, s.Context().(mocktracer.MockSpanContext).TraceID)
				}
			*/
		})

	}
}
