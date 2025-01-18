package wrapper_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"

	//"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	vmeter "go.unistack.org/micro-meter-victoriametrics/v3"
	ot "go.unistack.org/micro-tracer-opentracing/v3"
	wrapper "go.unistack.org/micro-wrapper-sql/v3"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/meter"
	"go.unistack.org/micro/v3/tracer"
	"modernc.org/sqlite"
)

var schema = `
  CREATE TABLE IF NOT EXISTS person (
	  first_name text,
	  last_name text,
	  email text
	);`

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
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

func TestSqliteWrapper(t *testing.T) {
	t.Skip()
	ctx := context.Background()
	wrapper.DefaultMeterStatsInterval = 100 * time.Millisecond
	meter.DefaultMeter = vmeter.NewMeter()
	buf := bytes.NewBuffer(nil)
	logger.DefaultLogger = logger.NewLogger(logger.WithLevel(logger.DebugLevel), logger.WithOutput(buf))

	if err := logger.DefaultLogger.Init(); err != nil {
		t.Fatal(err)
	}

	tr, c := initJaeger(fmt.Sprintf("Test tracing %s", time.Now().Format(time.RFC1123Z)))
	defer c.Close()
	opentracing.SetGlobalTracer(tr)
	tracer.DefaultTracer = ot.NewTracer(ot.Tracer(tr))
	if err := tracer.DefaultTracer.Init(); err != nil {
		t.Fatal(ctx, err)
	}

	sql.Register("micro-wrapper-sql", wrapper.NewWrapper(&sqlite.Driver{},
		wrapper.DatabaseHost("localhost"),
		wrapper.DatabaseName("memory"),
	))
	wdb, err := sql.Open("micro-wrapper-sql", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err = wdb.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	db := sqlx.NewDb(wdb, "sqlite")
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	wrapper.NewStatsMeter(ctx, db, wrapper.DatabaseHost("localhost"), wrapper.DatabaseName("memory"))

	tx, err := wdb.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.ExecContext(wrapper.QueryName(ctx, "schema create"), schema); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	tx, err = wdb.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.ExecContext(wrapper.QueryName(ctx, "insert one"), "INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Fist1", "Last1", "Email1"); err != nil {
		t.Fatal(err)
	}

	if _, err := wdb.ExecContext(wrapper.QueryName(ctx, "double schema"), schema); err != nil {
		t.Fatal(err)
	}

	if _, err := tx.ExecContext(wrapper.QueryName(ctx, "insert two"), "INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Fist2", "Last2", "Email2"); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	var peoples []*Person
	if err := sqlx.SelectContext(wrapper.QueryName(ctx, "get_all_person"), db, &peoples, "SELECT * FROM person limit 2"); err != nil {
		t.Fatal(err)
	}

	_ = peoples
	time.Sleep(1 * time.Second)

	mbuf := bytes.NewBuffer(nil)
	_ = meter.DefaultMeter.Write(mbuf, meter.WriteProcessMetrics(true))

	if !bytes.Contains(mbuf.Bytes(), []byte(`micro_sql_idle_connections`)) {
		t.Fatalf("micro-wrapper-sql meter output contains invalid output: %s", buf.Bytes())
	}

	for _, tcase := range [][]byte{
		[]byte(`"method":"ExecContext"`),
		[]byte(`"method":"Open"`),
		[]byte(`"method":"BeginTx"`),
		[]byte(`"method":"Commit"`),
		[]byte(`"method":"QueryContext"`),
		[]byte(`"query":"get_all_person"`),
		[]byte(`"took":`),
	} {
		if !bytes.Contains(buf.Bytes(), tcase) {
			t.Fatalf("micro-wrapper-sql logger output miss %s in output: %s", tcase, buf.Bytes())
		}
	}

	t.Logf("%s", buf.Bytes())
}

/*
func TestPostgresWrapper(t *testing.T) {
	ctx := context.Background()
	wrapper.DefaultMeterStatsInterval = 100 * time.Millisecond
	meter.DefaultMeter = vmeter.NewMeter()
	buf := bytes.NewBuffer(nil)
	logger.DefaultLogger = logger.NewLogger(logger.WithLevel(logger.DebugLevel), logger.WithOutput(buf))

	if err := logger.DefaultLogger.Init(); err != nil {
		t.Fatal(err)
	}

	tr, c := initJaeger(fmt.Sprintf("Test tracing %s", time.Now().Format(time.RFC1123Z)))
	defer c.Close()
	opentracing.SetGlobalTracer(tr)
	tracer.DefaultTracer = ot.NewTracer(ot.Tracer(tr))
	if err := tracer.DefaultTracer.Init(); err != nil {
		logger.Fatal(ctx, err)
	}

	sql.Register("micro-wrapper-sql", wrapper.NewWrapper(&sqlite.Driver{},
		wrapper.DatabaseHost("localhost"),
		wrapper.DatabaseName("memory"),
		wrapper.LoggerLevel(logger.DebugLevel),
		wrapper.LoggerEnabled(true),
	))
	wdb, err := sql.Open("micro-wrapper-sql", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err = wdb.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	db := sqlx.NewDb(wdb, "sqlite")
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	wrapper.NewStatsMeter(ctx, db, wrapper.DatabaseHost("localhost"), wrapper.DatabaseName("memory"))

	tx, err := wdb.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.ExecContext(wrapper.QueryName(ctx, "schema create"), schema); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	tx, err = wdb.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.ExecContext(wrapper.QueryName(ctx, "insert one"), "INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Fist1", "Last1", "Email1"); err != nil {
		t.Fatal(err)
	}

	if _, err := wdb.ExecContext(wrapper.QueryName(ctx, "double schema"), schema); err != nil {
		t.Fatal(err)
	}

	if _, err := tx.ExecContext(wrapper.QueryName(ctx, "insert two"), "INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Fist2", "Last2", "Email2"); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	var peoples []*Person
	if err := sqlx.SelectContext(wrapper.QueryName(ctx, "get_all_person"), db, &peoples, "SELECT * FROM person limit 2"); err != nil {
		t.Fatal(err)
	}

	_ = peoples
	time.Sleep(1 * time.Second)

	mbuf := bytes.NewBuffer(nil)
	_ = meter.DefaultMeter.Write(mbuf, meter.WriteProcessMetrics(true))

	if !bytes.Contains(mbuf.Bytes(), []byte(`micro_sql_idle_connections`)) {
		t.Fatalf("micro-wrapper-sql meter output contains invalid output: %s", buf.Bytes())
	}

	for _, tcase := range [][]byte{
		[]byte(`"method":"ExecContext"`),
		[]byte(`"method":"Open"`),
		[]byte(`"method":"BeginTx"`),
		[]byte(`"method":"Commit"`),
		[]byte(`"method":"QueryContext"`),
		[]byte(`"query":"get_all_person"`),
		[]byte(`"took":`),
	} {
		if !bytes.Contains(buf.Bytes(), tcase) {
			t.Fatalf("micro-wrapper-sql logger output miss %s in output: %s", tcase, buf.Bytes())
		}
	}

	t.Logf("%s", buf.Bytes())
}
*/
