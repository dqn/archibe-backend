package main

import (
	"log"
	"os"

	"github.com/dqn/archibe/app"
	"github.com/dqn/archibe/dbexec"
	"github.com/jmoiron/sqlx"
	"github.com/k0kubun/pp"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func run() error {
	if len(os.Args) != 3 {
		os.Exit(1)
	}

	address := os.Args[1]
	dns := os.Args[2]

	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		return err
	}
	defer db.Close()
	dbx := dbexec.NewExecutor(db)

	e := echo.New()

	a := app.App{DBX: dbx, Server: e}
	a.Start(address)
	pp.Print()

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
