package db

import (
	"database/sql"
	"fmt"

	"golib/lang"
	"golib/log"

	_ "github.com/denisenkom/go-mssqldb"
)

type SQLServer struct {
	Server   string
	Port     int
	Database string
	Username string
	Password string
	Encrypt  bool
}

func (self SQLServer) Connect() (db *sql.DB) {
	var err error

	enableOrDisable := func(b bool) string {
		if b {
			return "enable"
		} else {
			return "disable"
		}
	}

	connectionString := fmt.Sprintf(
		"server=%s;port=%d;database=%s;encrypt=%s",
		self.Server, self.Port, self.Database, enableOrDisable(self.Encrypt),
	)
	if len(self.Username) > 0 {
		connectionString += fmt.Sprintf(
			";user id=%s;password=%s",
			self.Username, self.Password,
		)
	}
	log.Debug("Connecting SQLServer... %s", lang.Mask(connectionString, self.Password))
	db, err = sql.Open("mssql", connectionString)
	assertNil(err)

	return
}
