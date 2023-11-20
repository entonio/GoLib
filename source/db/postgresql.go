package db

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"sync"

	"golib/lang"
	"golib/log"

	//_ "github.com/jackc/pgx/v5/stdlib" // 'pgx'
	_ "github.com/lib/pq" // 'postgres'
)

type PostgreSQL struct {
	Server   string
	Port     int
	Database string
	Username string
	Password string
	Encrypt  bool

	address string
}

var once sync.Once

func (self PostgreSQL) Connect() (db *sql.DB) {
	var err error

	enableOrDisable := func(b bool) string {
		if b {
			return "require"
		} else {
			return "disable"
		}
	}

	once.Do(func() {
		self.address = self.Server
		ips, err := net.DefaultResolver.LookupNetIP(context.TODO(), "ip4", self.Server)
		// log.Trace("ips: %v, err: %v", ips, err)
		if err == nil && len(ips) > 0 {
			self.address = ips[0].String()
		}
	})

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		self.Username, self.Password, self.address, self.Port, self.Database, enableOrDisable(self.Encrypt),
	)

	log.Debug("Connecting PostgreSQL... %s", lang.Mask(connectionString, self.Password))
	db, err = sql.Open("postgres", connectionString)
	assertNil(err)

	return
}
