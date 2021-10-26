//go:build ignore
// +build ignore

package hystrix

import (
	"context"
	"testing"

	"github.com/afex/hystrix-go/hystrix"
	rrouter "go.unistack.org/micro-router-register/v3"
	"go.unistack.org/micro/register/memory"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/router"
)

func TestBreaker(t *testing.T) {
	// setup
	register := memory.NewRegister()

	c := client.NewClient(
		// set the selector
		client.Router(rrouter.NewRouter(router.Register(register))),
		// add the breaker wrapper
		client.Wrap(NewClientWrapper()),
	)

	req := c.NewRequest("test.service", "Test.Method", map[string]string{
		"foo": "bar",
	}, client.WithContentType("application/json"))

	var rsp map[string]interface{}

	// Force to point of trip
	for i := 0; i < (hystrix.DefaultVolumeThreshold * 3); i++ {
		c.Call(context.TODO(), req, rsp)
	}

	err := c.Call(context.TODO(), req, rsp)
	if err == nil {
		t.Error("Expecting tripped breaker, got nil error")
	}

	if err.Error() != "hystrix: circuit open" {
		t.Errorf("Expecting tripped breaker, got %v", err)
	}
}
