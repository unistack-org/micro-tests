// +build ignore

package gobreaker

import (
	"context"
	"testing"

	"github.com/sony/gobreaker"
	"github.com/unistack-org/micro/register/memory"
	"github.com/unistack-org/micro/v3/client"
	"github.com/unistack-org/micro/v3/errors"
	"github.com/unistack-org/micro/v3/router"
)

func TestBreaker(t *testing.T) {
	// setup
	r := memory.NewRegister()

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
	for i := 0; i < 6; i++ {
		c.Call(context.TODO(), req, rsp)
	}

	err := c.Call(context.TODO(), req, rsp)
	if err == nil {
		t.Error("Expecting tripped breaker, got nil error")
	}

	merr := err.(*errors.Error)
	if merr.Code != 502 {
		t.Errorf("Expecting tripped breaker, got %v", err)
	}
}

func TestCustomBreaker(t *testing.T) {
	// setup
	r := memory.NewRegister()

	c := client.NewClient(
		// set the selector
		client.Router(rrouter.NewRouter(router.Register(register))),
		// add the breaker wrapper
		client.Wrap(NewCustomClientWrapper(
			gobreaker.Settings{},
			BreakService,
		)),
	)

	req := c.NewRequest("test.service", "Test.Method", map[string]string{
		"foo": "bar",
	}, client.WithContentType("application/json"))

	var rsp map[string]interface{}

	// Force to point of trip
	for i := 0; i < 6; i++ {
		c.Call(context.TODO(), req, rsp)
	}

	err := c.Call(context.TODO(), req, rsp)
	if err == nil {
		t.Error("Expecting tripped breaker, got nil error")
	}

	merr := err.(*errors.Error)
	if merr.Code != 502 {
		t.Errorf("Expecting tripped breaker, got %v", err)
	}
}
