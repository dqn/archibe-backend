package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func run() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("invalid arguments")
	}

	dsn := os.Args[1]

	f, err := os.Open(path.Join("sqls", "testdata.sql"))
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}

	if _, err = db.Exec(string(b)); err != nil {
		return err
	}

	fmt.Println("completed!")

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
