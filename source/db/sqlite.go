package db

import (
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"golib/log"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	File       string
	Connection *Connection
}

type SQLiteRowId int64

func AsRowId(s string) SQLiteRowId {
	rowId, err := strconv.ParseInt(s, 10, 64)
	assertNil(err)
	return SQLiteRowId(rowId)
}

func AssertSQLiteAvailable() {
	if !strings.Contains(runtime.GOARCH, "64") {
		log.Debug("Para coisas que envolvem o SQLite tem de se usar a versão x64 do .exe")
		os.Exit(1)
	}
}

func (self SQLite) Connect() (db *sql.DB) {
	var err error

	connectionString := fmt.Sprintf(
		"file:%s?parseTime=true",
		self.File,
	)
	log.Trace("Connecting SQLite... %s", connectionString)
	db, err = sql.Open("sqlite3", connectionString)
	assertNil(err)

	return
}
