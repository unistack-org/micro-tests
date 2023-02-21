package segmentio_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	segmentio "go.unistack.org/micro-broker-segmentio/v3"
	victoriameter "go.unistack.org/micro-meter-victoriametrics/v3"
	https "go.unistack.org/micro-server-http/v3"
	"go.unistack.org/micro/v3/broker"
	"go.unistack.org/micro/v3/codec"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/meter"
	meterhandler "go.unistack.org/micro/v3/meter/handler"
	"go.unistack.org/micro/v3/server"
)

type lg struct{}

func (l *lg) Printf(format string, args ...interface{}) {
	//	logger.Infof(context.Background(), format, args...)
}

var bm = &broker.Message{
	Header: map[string]string{"hkey": "hval"},
	Body:   []byte(`"body"`),
}

func TestConsumerGroup(t *testing.T) {
	topic := "test_topic"
	if tr := os.Getenv("INTEGRATION_TESTS"); len(tr) > 0 {
		return
		// t.Skip()
	}

	if err := logger.DefaultLogger.Init(logger.WithLevel(logger.TraceLevel)); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	var addrs []string
	if addr := os.Getenv("BROKER_ADDRS"); len(addr) == 0 {
		addrs = []string{"127.0.0.1:9092"}
	} else {
		addrs = strings.Split(addr, ",")
	}

	meter.DefaultMeter = victoriameter.NewMeter()

	s := https.NewServer(server.Context(ctx), server.Address("127.0.0.1:0"), server.Codec("text/plain", codec.NewCodec()))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	if err := meterhandler.RegisterMeterServiceServer(s, meterhandler.NewHandler()); err != nil {
		t.Fatal(err)
	}

	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	segmentio.DefaultWriterConfig.Async = true
	segmentio.DefaultWriterConfig.BatchTimeout = 1 * time.Second
	segmentio.DefaultWriterConfig.RequiredAcks = int(kafka.RequireAll)
	segmentio.DefaultReaderConfig.StartOffset = kafka.FirstOffset
	segmentio.DefaultReaderConfig.MinBytes = 1024 * 10        // 10 kb
	segmentio.DefaultReaderConfig.MaxBytes = 1024 * 1024 * 20 // 20 Mb
	segmentio.DefaultReaderConfig.MaxWait = 20 * time.Second  // 20s
	segmentio.DefaultReaderConfig.QueueCapacity = 500
	segmentio.DefaultReaderConfig.ReadBackoffMin = 2 * time.Second
	segmentio.DefaultReaderConfig.ReadBackoffMax = 5 * time.Second
	segmentio.DefaultReaderConfig.Logger = &lg{}
	segmentio.DefaultReaderConfig.CommitInterval = 1 * time.Second
	brk := segmentio.NewBroker(broker.Context(ctx), broker.Addrs(addrs...), segmentio.StatsInterval(5*time.Second),
		segmentio.ClientID("test_sub"),
	)
	t.Logf("init")
	if err := brk.Init(); err != nil {
		t.Fatal(err)
	}

	t.Logf("connect")
	if err := brk.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	defer func() {
		t.Logf("disconnect")
		if err := brk.Disconnect(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	fmt.Printf("prefill topic\n")
	go func() {
		for i := 0; i < 900000; i++ {
			if err := brk.Publish(ctx, topic, bm); err != nil {
				t.Fatal(err)
			}
			// log.Printf("publish %d", i)
			time.Sleep(600 * time.Millisecond)
		}
	}()
	fmt.Printf("prefill complete\n")

	var cnt uint64
	var wait atomic.Value
	wait.Store(true)

	fn := func(msg broker.Event) error {
		if wait.Load().(bool) {
			wait.Store(false)
		}
		atomic.AddUint64(&cnt, 1)
		fmt.Printf("processing mesage")
		return msg.Ack()
	}

	sub1, err := brk.Subscribe(ctx, topic, fn, broker.SubscribeGroup("test"), broker.SubscribeBodyOnly(true), broker.SubscribeAutoAck(true))
	if err != nil {
		t.Fatal(err)
	}
	defer sub1.Unsubscribe(ctx)

	sub2, err := brk.Subscribe(ctx, topic, fn, broker.SubscribeGroup("test"), broker.SubscribeBodyOnly(true), broker.SubscribeAutoAck(true))
	if err != nil {
		t.Fatal(err)
	}
	defer sub2.Unsubscribe(ctx)

	fmt.Printf("wait for ready\n")
	for {
		if !wait.Load().(bool) {
			break
		}
		time.Sleep(1 * time.Second)
	}
	/*
		time.Sleep(5 * time.Second)
		fmt.Printf("unsub\n")
		if err := sub1.Unsubscribe(ctx); err != nil {
			t.Fatal(err)
		}
		time.Sleep(9 * time.Second)
		sub1, err = brk.Subscribe(ctx, topic, fn, broker.SubscribeGroup("test"), broker.SubscribeBodyOnly(true), broker.SubscribeAutoAck(true))
		if err != nil {
			t.Fatal(err)
		}
	*/
	/*
		t1 := time.NewTicker(10 * time.Second)
		defer t1.Stop()
		t2 := time.NewTicker(30 * time.Second)
		defer t2.Stop()

			for {
				select {
				case <-t1.C:
					fmt.Printf("unsub from sub2\n")
					if err := sub2.Unsubscribe(ctx); err != nil {
						t.Fatal(err)
					}
					fmt.Printf("sub1\n")
					sub1, err = brk.Subscribe(ctx, topic, fn, broker.SubscribeGroup("test"), broker.SubscribeBodyOnly(true), broker.SubscribeAutoAck(true))
					if err != nil {
						t.Fatal(err)
					}
				case <-t2.C:
					fmt.Printf("unsub from sub1\n")
					if err := sub1.Unsubscribe(ctx); err != nil {
						t.Fatal(err)
					}
					fmt.Printf("sub2\n")
					sub2, err = brk.Subscribe(ctx, topic, fn, broker.SubscribeGroup("test"), broker.SubscribeBodyOnly(true), broker.SubscribeAutoAck(true))
					if err != nil {
						t.Fatal(err)
					}
				}
			}
	*/
	select {}
}

func TestSub(t *testing.T) {
	topic := "test_topic"
	if tr := os.Getenv("INTEGRATION_TESTS"); len(tr) > 0 {
		return
		// t.Skip()
	}

	if err := logger.DefaultLogger.Init(logger.WithLevel(logger.ErrorLevel)); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	var addrs []string
	if addr := os.Getenv("BROKER_ADDRS"); len(addr) == 0 {
		addrs = []string{"127.0.0.1:9092"}
	} else {
		addrs = strings.Split(addr, ",")
	}

	meter.DefaultMeter = victoriameter.NewMeter()

	s := https.NewServer(server.Context(ctx), server.Address("127.0.0.1:0"), server.Codec("text/plain", codec.NewCodec()))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	if err := meterhandler.RegisterMeterServiceServer(s, meterhandler.NewHandler()); err != nil {
		t.Fatal(err)
	}

	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	segmentio.DefaultWriterConfig.Async = true
	segmentio.DefaultWriterConfig.BatchTimeout = 1 * time.Second
	segmentio.DefaultWriterConfig.RequiredAcks = int(kafka.RequireAll)
	segmentio.DefaultReaderConfig.StartOffset = kafka.FirstOffset
	segmentio.DefaultReaderConfig.MinBytes = 1024 * 10        // 10 kb
	segmentio.DefaultReaderConfig.MaxBytes = 1024 * 1024 * 20 // 20 Mb
	segmentio.DefaultReaderConfig.MaxWait = 20 * time.Second  // 20s
	segmentio.DefaultReaderConfig.QueueCapacity = 500
	segmentio.DefaultReaderConfig.ReadBackoffMin = 2 * time.Second
	segmentio.DefaultReaderConfig.ReadBackoffMax = 5 * time.Second
	segmentio.DefaultReaderConfig.Logger = &lg{}
	segmentio.DefaultReaderConfig.CommitInterval = 1 * time.Second
	brk := segmentio.NewBroker(broker.Context(ctx), broker.Addrs(addrs...), segmentio.StatsInterval(5*time.Second),
		segmentio.ClientID("test_sub"),
	)
	t.Logf("init")
	if err := brk.Init(); err != nil {
		t.Fatal(err)
	}

	t.Logf("connect")
	if err := brk.Connect(ctx); err != nil {
		t.Fatal(err)
	}

	defer func() {
		t.Logf("disconnect")
		if err := brk.Disconnect(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	fmt.Printf("prefill topic\n")

	go func() {
		for i := 0; i < 900000; i++ {
			if err := brk.Publish(ctx, topic, bm); err != nil {
				t.Fatal(err)
			}
			log.Printf("publish %d", i)
			//		time.Sleep(1 * time.Second)
		}
	}()
	fmt.Printf("prefill complete\n")

	var cnt uint64
	var wait atomic.Value
	wait.Store(true)

	done := make(chan struct{})
	fn := func(msg broker.Event) error {
		if wait.Load().(bool) {
			wait.Store(false)
			fmt.Printf("done ready\n")
			close(done)
		}
		atomic.AddUint64(&cnt, 1)
		return msg.Ack()
	}

	sub, err := brk.Subscribe(ctx, topic, fn, broker.SubscribeGroup("test"), broker.SubscribeBodyOnly(true))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("wait for ready\n")
	<-done
	fmt.Printf("wait for bench\n")
	fmt.Printf("start %s\n", time.Now().String())
	<-time.After(20 * time.Second)
	fmt.Printf("stop %s\n", time.Now().String())
	rcnt := atomic.LoadUint64(&cnt)

	req, err := http.NewRequest(http.MethodGet, "http://"+s.Options().Address+"/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "text/plain")

	rsp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rsp.Body.Close()

	buf, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("unsub\n")
	if err := sub.Unsubscribe(ctx); err != nil {
		t.Fatal(err)
	}

	t.Logf("metrics: \n%s\n", buf)
	t.Logf("mesage count %d\n", rcnt)
}

func BenchmarkPub(b *testing.B) {
	if tr := os.Getenv("INTEGRATION_TESTS"); len(tr) > 0 {
		b.Skip()
	}

	if err := logger.DefaultLogger.Init(logger.WithLevel(logger.TraceLevel)); err != nil {
		b.Fatal(err)
	}
	ctx := context.Background()

	var addrs []string
	if addr := os.Getenv("BROKER_ADDRS"); len(addr) == 0 {
		addrs = []string{"127.0.0.1:9092"}
	} else {
		addrs = strings.Split(addr, ",")
	}

	meter.DefaultMeter = victoriameter.NewMeter()

	s := https.NewServer(server.Context(ctx), server.Address("127.0.0.1:0"), server.Codec("text/plain", codec.NewCodec()))
	if err := s.Init(); err != nil {
		b.Fatal(err)
	}
	if err := meterhandler.RegisterMeterServiceServer(s, meterhandler.NewHandler()); err != nil {
		b.Fatal(err)
	}

	if err := s.Start(); err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err := s.Stop(); err != nil {
			b.Fatal(err)
		}
	}()

	segmentio.DefaultWriterConfig.Async = false
	segmentio.DefaultWriterConfig.BatchTimeout = 1 * time.Second
	segmentio.DefaultWriterConfig.RequiredAcks = int(kafka.RequireAll)
	fn := func(msgs []kafka.Message, err error) {
		if err != nil {
			b.Logf("err %v", err)
		}
	}
	brk := segmentio.NewBroker(broker.Context(ctx), broker.Addrs(addrs...), segmentio.StatsInterval(1*time.Second),
		segmentio.WriterCompletionFunc(fn))
	b.Logf("init")
	if err := brk.Init(); err != nil {
		b.Fatal(err)
	}

	b.Logf("connect")
	if err := brk.Connect(ctx); err != nil {
		b.Fatal(err)
	}
	defer func() {
		b.Logf("disconnect")
		if err := brk.Disconnect(ctx); err != nil {
			b.Fatal(err)
		}
	}()

	cnt := 0
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := brk.Publish(ctx, "test_topic", bm); err != nil {
			b.Fatal(err)
		}
		cnt++
	}

	req, err := http.NewRequest(http.MethodGet, "http://"+s.Options().Address+"/metrics", nil)
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Add("Content-Type", "text/plain")

	rsp, err := (&http.Client{}).Do(req)
	if err != nil {
		b.Fatal(err)
	}

	defer rsp.Body.Close()

	buf, err := io.ReadAll(rsp.Body)
	if err != nil {
		b.Fatal(err)
	}

	b.Logf("metrics: \n%s\n", buf)
	b.Logf("mesage count %d\n", cnt)
}

func BenchmarkPubSub(b *testing.B) {
	b.Skip()
	ctx := context.Background()
	topic := "test_topic"
	var addrs []string
	if addr := os.Getenv("BROKER_ADDRS"); len(addr) == 0 {
		addrs = []string{"127.0.0.1:9092"}
	} else {
		addrs = strings.Split(addr, ",")
	}

	segmentio.DefaultWriterConfig.Async = true
	segmentio.DefaultWriterConfig.BatchTimeout = 1 * time.Second
	segmentio.DefaultReaderConfig.CommitInterval = 2 * time.Second
	brk := segmentio.NewBroker(broker.Context(ctx), broker.Addrs(addrs...), segmentio.StatsInterval(1*time.Minute))
	if err := brk.Init(); err != nil {
		b.Fatal(err)
	}

	if err := brk.Connect(ctx); err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err := brk.Disconnect(ctx); err != nil {
			b.Fatal(err)
		}
	}()

	wait := true
	var cnt uint64
	fn := func(msg broker.Event) error {
		if wait {
			wait = false
		}
		atomic.AddUint64(&cnt, 1)
		return msg.Ack()
	}

	if err := brk.Publish(ctx, topic, bm); err != nil {
		b.Fatal(err)
	}

	sub, err := brk.Subscribe(ctx, topic, fn, broker.SubscribeGroup("test"), broker.SubscribeBodyOnly(true))
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		if err := sub.Unsubscribe(ctx); err != nil {
			b.Fatal(err)
		}
	}()

	for {
		if !wait {
			break
		}
		time.Sleep(1 * time.Second)
	}
	b.ResetTimer()
	var result error
	sent := uint64(0)
	for n := 0; n < b.N; n++ {
		if err := brk.Publish(ctx, topic, bm); err != nil {
			b.Fatal(err)
		} else {
			result = err
		}
		sent++
	}

	b.Logf("publish done")
	for {
		c := atomic.LoadUint64(&cnt)
		if c >= sent {
			break
		}
		fmt.Printf("c %d seen %d\n", c, sent)
		time.Sleep(1 * time.Second)
	}
	_ = result
	fmt.Printf("c %d seen %d\n", atomic.LoadUint64(&cnt), sent)
}

func TestPubSub(t *testing.T) {
	if tr := os.Getenv("INTEGRATION_TESTS"); len(tr) > 0 {
		return
		// t.Skip()
	}

	if err := logger.DefaultLogger.Init(logger.WithLevel(logger.ErrorLevel)); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	var addrs []string
	if addr := os.Getenv("BROKER_ADDRS"); len(addr) == 0 {
		addrs = []string{"127.0.0.1:9092"}
	} else {
		addrs = strings.Split(addr, ",")
	}

	meter.DefaultMeter = victoriameter.NewMeter()

	s := https.NewServer(server.Context(ctx), server.Address("127.0.0.1:0"), server.Codec("text/plain", codec.NewCodec()))
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	if err := meterhandler.RegisterMeterServiceServer(s, meterhandler.NewHandler()); err != nil {
		t.Fatal(err)
	}

	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	b := segmentio.NewBroker(broker.Context(ctx), broker.Addrs(addrs...), segmentio.StatsInterval(500*time.Millisecond),
		segmentio.ClientID("test_pubsub"))
	t.Logf("init")
	if err := b.Init(); err != nil {
		t.Fatal(err)
	}

	t.Logf("connect")
	if err := b.Connect(ctx); err != nil {
		t.Fatal(err)
	}
	defer func() {
		t.Logf("disconnect")
		if err := b.Disconnect(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	wait := true
	fn := func(msg broker.Event) error {
		wait = false
		return msg.Ack()
	}

	t.Logf("subscribe")
	sub, err := b.Subscribe(ctx, "test_topic", fn, broker.SubscribeGroup("test"))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		t.Logf("unsubscribe")
		if err := sub.Unsubscribe(ctx); err != nil {
			t.Fatal(err)
		}
	}()

	if err := b.Publish(ctx, "test_topic", bm); err != nil {
		t.Fatal(err)
	}

	for {
		if !wait {
			break
		}
		time.Sleep(1 * time.Second)
	}

	req, err := http.NewRequest(http.MethodGet, "http://"+s.Options().Address+"/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "text/plain")

	rsp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rsp.Body.Close()

	buf, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("metrics: \n%s\n", buf)
}
