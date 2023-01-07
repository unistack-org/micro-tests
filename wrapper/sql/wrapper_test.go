package wrapper_test

import (
	"bytes"
	"context"
	"database/sql"
	//"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	vmeter "go.unistack.org/micro-meter-victoriametrics/v3"
	wrapper "go.unistack.org/micro-wrapper-sql/v3"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/meter"
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

func TestWrapper(t *testing.T) {
	ctx := context.Background()
	wrapper.DefaultMeterStatsInterval = 100 * time.Millisecond
	meter.DefaultMeter = vmeter.NewMeter()
	buf := bytes.NewBuffer(nil)
	logger.DefaultLogger = logger.NewLogger(logger.WithLevel(logger.DebugLevel), logger.WithOutput(buf))

	if err := logger.DefaultLogger.Init(); err != nil {
		t.Fatal(err)
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

	db := sqlx.NewDb(wdb, "sqlite")
	var cancel func()
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	wrapper.NewStatsMeter(ctx, db, wrapper.DatabaseHost("localhost"), wrapper.DatabaseName("memory"))
	if _, err := db.ExecContext(ctx, schema); err != nil {
		t.Fatal(err)
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.NamedExecContext(ctx, "INSERT OR REPLACE INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Fist1", "Last1", "Email1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.NamedExecContext(ctx, "INSERT OR REPLACE INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Fist2", "Last2", "Email2"}); err != nil {
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
			t.Fatalf("micro-wrapper-sql logger output contains invalid output: %s", buf.Bytes())
		}
	}
}
