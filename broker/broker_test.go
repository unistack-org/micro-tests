package broker_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	kg "github.com/twmb/franz-go/pkg/kgo"
	kgo "go.unistack.org/micro-broker-kgo/v3"
	jsoncodec "go.unistack.org/micro-codec-json/v3"
	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
)

var (
	msgcnt    = int64(60000)
	group     = "33"
	prefill   = true
	subtopic  = "subtest"
	pubtopic  = "pubtest"
	rateRecs  int64
	rateBytes int64
)

var subbm = &broker.Message{
	Header: map[string]string{"hkey": "hval", metadata.HeaderTopic: subtopic},
	Body:   []byte(`"bodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybody"`),
}

var pubbm = &broker.Message{
	Header: map[string]string{"hkey": "hval", metadata.HeaderTopic: pubtopic},
	Body:   []byte(`"bodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybodybody"`),
}

func printRate() {
	recs := atomic.SwapInt64(&rateRecs, 0)
	bytes := atomic.SwapInt64(&rateBytes, 0)
	fmt.Printf("%0.2f MiB/s; %0.2fk records/s\n", float64(bytes)/(1024*1024), float64(recs)/1000)
}

func TestOptionPassing(t *testing.T) {
	opts := []client.PublishOption{kgo.ClientPublishKey([]byte(`test`))}
	options := client.NewPublishOptions(opts...)

	if !strings.Contains(fmt.Sprintf("%#+v\n", options.Context), "kgo.publishKey") {
		t.Fatal("context does not have publish key")
	}
}

func TestKgo(t *testing.T) {
	if tr := os.Getenv("INTEGRATION_TESTS"); len(tr) > 0 {
		t.Skip()
	}

	_ = logger.DefaultLogger.Init(logger.WithLevel(logger.TraceLevel), logger.WithCallerSkipCount(3))
	ctx := context.Background()

	var addrs []string
	if addr := os.Getenv("BROKER_ADDRS"); len(addr) == 0 {
		addrs = []string{"172.18.0.201:29091", "172.18.0.202:29092", "172.18.0.203:29093"}
	} else {
		addrs = strings.Split(addr, ",")
	}

	b := kgo.NewBroker(
		broker.Codec(jsoncodec.NewCodec()),
		broker.Addrs(addrs...),
		kgo.CommitInterval(5*time.Second),
		kgo.Options(kg.ClientID("test"), kg.FetchMaxBytes(1*1024*1024)),
	)
	if err := b.Init(); err != nil {
		t.Fatal(err)
	}

	if err := b.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := b.Disconnect(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	if prefill {
		msgs := make([]*broker.Message, 0, msgcnt)
		for i := int64(0); i < msgcnt; i++ {
			msgs = append(msgs, subbm)
		}

		if err := b.BatchPublish(ctx, msgs); err != nil {
			t.Fatal(err)
		}
	}

	done := make(chan bool, 1)
	idx := int64(0)

	fn := func(m broker.Event) error {
		atomic.AddInt64(&idx, 1)
		time.Sleep(20 * time.Millisecond)
		if err := b.BatchPublish(ctx, []*broker.Message{pubbm}); err != nil {
			return err
		}
		atomic.AddInt64(&rateRecs, 1)
		atomic.AddInt64(&rateBytes, int64(len(m.Message().Body)))
		return m.Ack()
	}

	sub, err := b.Subscribe(ctx, subtopic, fn,
		broker.SubscribeAutoAck(true),
		broker.SubscribeGroup(group),
		broker.SubscribeBodyOnly(true))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := sub.Unsubscribe(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	// ticker := time.NewTicker(10 * time.Second)
	// defer ticker.Stop()

	pticker := time.NewTicker(1 * time.Second)
	defer pticker.Stop()

	go func() {
		for {
			select {
			case <-pticker.C:
				printRate()
				if prc := atomic.LoadInt64(&idx); prc == msgcnt {
					close(done)
				}
				//	case <-ticker.C:
				//		close(done)
			}
		}
	}()

	<-done
}
