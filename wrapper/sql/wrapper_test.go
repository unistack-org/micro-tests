package wrapper_test

import (
	"bytes"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	vmeter "go.unistack.org/micro-meter-victoriametrics/v3"
	wrapper "go.unistack.org/micro-wrapper-sql/v3"
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
	wrapper.DefaultMeterStatsInterval = 100 * time.Millisecond
	meter.DefaultMeter = vmeter.NewMeter()

	sql.Register("micro-wrapper-sql", wrapper.NewWrapper(&sqlite.Driver{}))
	wdb, err := sql.Open("micro-wrapper-sql", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	db := sqlx.NewDb(wdb, "sqlite")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wrapper.NewStatsMeter(ctx, db, wrapper.DatabaseHost("localhost"), wrapper.DatabaseName("memory"))
	if _, err := db.Exec(schema); err != nil {
		t.Fatal(err)
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tx.NamedExec("INSERT OR REPLACE INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Fist1", "Last1", "Email1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := tx.NamedExec("INSERT OR REPLACE INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Fist2", "Last2", "Email2"}); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	var peoples []*Person
	if err := sqlx.SelectContext(ctx, db, &peoples, "SELECT * FROM person limit 2"); err != nil {
		t.Fatal(err)
	}

	_ = peoples
	time.Sleep(1 * time.Second)
	buf := bytes.NewBuffer(nil)
	_ = meter.DefaultMeter.Write(buf, meter.WriteProcessMetrics(true))

	if !bytes.Contains(buf.Bytes(), []byte(`micro_sql_idle_connections`)) {
		t.Fatalf("micro-wrapper-sql output contains invalid output: %s", buf.Bytes())
	}
}
