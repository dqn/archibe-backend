package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dqn/tubekids/dbexec"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func run() error {
	if len(os.Args) != 3 {
		os.Exit(1)
	}

	dns := os.Args[1]
	channelID := os.Args[2]

	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		return err
	}
	defer db.Close()
	dbx := dbexec.NewExecutor(db)

	channel, err := dbx.Channels.Find(channelID)
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", channel)

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
