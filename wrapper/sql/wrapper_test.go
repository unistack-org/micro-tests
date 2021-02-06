// +build ignore

package wrapper_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	//_ "github.com/mattn/go-sqlite3"
	wrapper "github.com/unistack-org/micro-wrapper-sql"
)

var (
	schema = `
  CREATE TABLE IF NOT EXISTS person (
	  first_name text,
	  last_name text,
	  email text
	);`
)

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

func TestWrapper(t *testing.T) {
	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		t.Fatal(err)
	}

	w := wrapper.NewWrapper(db)
	defer w.Close()
	ctx := context.Background()

	db.MustExec(schema)

	tx := db.MustBegin()
	//tx.NamedExec("INSERT OR REPLACE INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com"})
	tx.Commit()

	people := &Person{}
	if err := sqlx.GetContext(ctx, w, people, "SELECT * FROM person limit 1"); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#+v\n", people)
}
