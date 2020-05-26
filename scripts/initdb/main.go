package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const sqlsPath = "./sqls"

func getSQL(name string) (sql string, err error) {
	f, err := os.Open(path.Join(sqlsPath, name+".sql"))
	if err != nil {
		return
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	sql = string(b)

	return
}

func run() error {
	if len(os.Args) != 2 {
		os.Exit(1)
	}

	dsn := os.Args[1]

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return err
	}

	names := []string{
		"schemas",
		"tables",
		"indexes",
	}

	var sql string
	for _, v := range names {
		s, err := getSQL(v)
		if err != nil {
			return err
		}
		sql += "\n" + s
	}

	if _, err := db.Exec(sql); err != nil {
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
