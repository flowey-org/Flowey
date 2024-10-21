package db

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"slices"
)

func Occupied(path string) (bool, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return true, nil
		} else {
			return false, err
		}
	}
	file.Close()
	return false, nil
}

type tableInfoRow struct {
	cid        int
	name       string
	typeDef    string
	notnull    int
	dflt_value interface{}
	pk         int
}

type tableInfo = []tableInfoRow

func Validate(path string) error {
	expectedTableInfos := []tableInfo{
		{
			{cid: 0, name: "id", typeDef: "INTEGER", notnull: 1, dflt_value: nil, pk: 1},
			{cid: 1, name: "username", typeDef: "TEXT", notnull: 1, dflt_value: nil, pk: 0},
			{cid: 2, name: "password", typeDef: "TEXT", notnull: 1, dflt_value: nil, pk: 0},
		},
		{
			{cid: 0, name: "session_key", typeDef: "TEXT", notnull: 1, dflt_value: nil, pk: 1},
			{cid: 1, name: "user_id", typeDef: "INTEGER", notnull: 1, dflt_value: nil, pk: 0},
		},
	}

	tableInfos := []tableInfo{}
	for _, tableName := range []string{"users", "sessions"} {
		rows, err := db.Query(`PRAGMA table_info(` + tableName + `)`)
		if err != nil {
			return err
		}

		tableInfo := tableInfo{}
		for rows.Next() {
			var r tableInfoRow
			if err := rows.Scan(
				&r.cid, &r.name, &r.typeDef, &r.notnull, &r.dflt_value, &r.pk,
			); err != nil {
				return err
			}
			tableInfo = append(tableInfo, r)
		}

		tableInfos = append(tableInfos, tableInfo)
		rows.Close()
	}

	for i, expectedTableInfo := range expectedTableInfos {
		if !slices.Equal(expectedTableInfo, tableInfos[i]) {
			return fmt.Errorf("failed to validate the database")
		}
	}

	log.Printf("validated the database")
	return nil
}

func Init(path string) error {
	query := `CREATE TABLE users(
  id INTEGER NOT NULL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);
CREATE TABLE sessions(
  session_key TEXT NOT NULL PRIMARY KEY,
  user_id INTEGER NOT NULL
)`
	if _, err := db.Exec(query); err != nil {
		log.Println(err)
		return nil
	}

	log.Printf("initialized the database")
	return nil
}

func Open(path string) error {
	var err error
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	log.Printf("opened a connection to the database at %s", path)
	return nil
}

func Prepare(path string) error {
	occupied, err := Occupied(path)
	if err != nil {
		return err
	}

	if err := Open(path); err != nil {
		return err
	}

	if occupied {
		err = Validate(path)
	} else {
		err = Init(path)
	}
	if err != nil {
		Close()
		return err
	}

	return nil
}

func PrepareCmd(args []string, path string) error {
	flagSet := flag.NewFlagSet("flowey db prepare", flag.ExitOnError)

	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: flowey db prepare")
	}

	flagSet.Parse(args)

	if flagSet.NArg() > 0 {
		flagSet.Usage()
		return nil
	}

	if err := Prepare(path); err != nil {
		return err
	}
	defer Close()

	return nil
}
