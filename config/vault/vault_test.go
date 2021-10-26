package vault_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	jsoncodec "go.unistack.org/micro-codec-json/v3"
	vault "go.unistack.org/micro-config-vault/v3"
	"go.unistack.org/micro/v3/config"
)

type Config struct {
	Name         string `json:"name"`
	ContactEmail string `json:"contact_email"`
}

func TestLoad(t *testing.T) {
	if len(os.Getenv("INTEGRATION_TESTS")) > 0 {
		t.Skip()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := &Config{}

	cfg := vault.NewConfig(config.Codec(jsoncodec.NewCodec()),
		config.Struct(conf), vault.Path("/secret/data/customer/acme"), vault.Token("s.3QKHWXe4VV7S0wqIZuKxuEv0"))
	if err := cfg.Init(); err != nil {
		t.Fatal(err)
	}

	if err := cfg.Load(ctx); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#+v\n", conf)
}
