package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func run() error {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	dsn := os.Args[1]

	pool, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	rows, err := pool.Query("SELECT 1;")
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			return err
		}
		fmt.Println(id)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
