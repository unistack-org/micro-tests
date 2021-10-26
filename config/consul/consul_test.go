package consul_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	jsoncodec "go.unistack.org/micro-codec-json/v3"
	consul "go.unistack.org/micro-config-consul/v3"
	"go.unistack.org/micro/v3/config"
)

type Cfg struct {
	StringValue string `default:"string_value"`
	IgnoreValue string `json:"-"`
	StructValue struct {
		StringValue string `default:"string_value"`
	}
	IntValue int `default:"99"`
}

func TestWatch(t *testing.T) {
	if tr := os.Getenv("INTEGRATION_TESTS"); len(tr) > 0 {
		t.Skip()
	}

	addrs := ""
	if addr := os.Getenv("CONSUL_ADDRS"); len(addr) == 0 {
		addrs = "127.0.0.1:8500"
	} else {
		addrs = addr
	}

	ctx := context.Background()

	conf := &Cfg{IntValue: 10}

	cfg := consul.NewConfig(config.Struct(conf), consul.Address(addrs), consul.Path("test/consul"), config.Codec(jsoncodec.NewCodec()))
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}
	if err := cfg.Load(ctx); err != nil {
		t.Fatal(err)
	}

	w, err := cfg.Watch(ctx, config.WatchInterval(700*time.Millisecond, 2000*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = w.Stop()
	}()

	done := make(chan struct{})

	go func() {
		for {
			mp, err := w.Next()
			if err != nil {
				t.Fatal(err)
			}
			if len(mp) != 1 {
				close(done)
				t.Fatal(fmt.Errorf("consul watcher err: %v", mp))
				return
			}
			v, ok := mp["IntValue"]
			if !ok {
				close(done)
				t.Fatal(fmt.Errorf("consul watcher err: %v", v))
				return
			}
			if nv, ok := v.(int); !ok || nv != 5 {
				close(done)
				t.Fatal(fmt.Errorf("consul watcher err: %v", v))
				return
			}
			close(done)
			return
		}
	}()

	<-done
}
