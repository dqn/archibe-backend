package main

import (
	"log"
	"os"

	_ "github.com/lib/pq"
)

func run() error {
	if len(os.Args) != 3 {
		os.Exit(1)
	}

	// dns := os.Args[1]
	// channelID := os.Args[2]

	// db, err := sqlx.Open("postgres", dns)
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()
	// dbx := dbexec.NewExecutor(db)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
