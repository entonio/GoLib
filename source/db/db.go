package db

import (
	"database/sql"
	"reflect"
	"strings"
	"time"

	"golib/arrays"
	"golib/lang"
	"golib/log"
)

type Connection struct {
	db *sql.DB
}

var openConnections []*Connection

func NewConnection(db *sql.DB) *Connection {
	self := &Connection{db: db}
	openConnections = append(openConnections, self)
	return self
}

func (self *Connection) Close() {
	self.db.Close()
	openConnections = arrays.Without(openConnections, self)
}

func Close() {
	for _, c := range openConnections {
		c.Close()
	}
}

type ConnectionInfo interface {
	Connect() *sql.DB
}

type ValuePreprocessor func(value any) (bool, any)

type DBMS struct {
	Preprocessor ValuePreprocessor
	Info         ConnectionInfo

	connection *Connection
}

func (self *DBMS) OpenConnection() *Connection {
	self.connection = NewConnection(self.Info.Connect())
	return self.connection
}

func (self *DBMS) ConnectAndRunQuery(query string, values []any, eachRow func(rows *sql.Rows)) {
	defer self.OpenConnection().Close()
	self.RunQuery(query, values, eachRow)
}

func (self *DBMS) ConnectAndRunDML(query string, values []any) sql.Result {
	defer self.OpenConnection().Close()
	return self.RunDML(query, values)
}

func (self *DBMS) RunQuery(query string, values []any, eachRow func(rows *sql.Rows)) {
	self.ensureConnection()
	runQuery(self.Preprocessor, self.connection.db, query, values, true, eachRow)
}

func (self *DBMS) RunDML(query string, values []any) sql.Result {
	self.ensureConnection()
	return runDML(self.Preprocessor, self.connection.db, query, values, true)
}

func (self *DBMS) ensureConnection() {
	if self.connection == nil {
		self.OpenConnection()
	}
}

func runDML(vp ValuePreprocessor, connection *sql.DB, dml string, values []any, printQuery bool) sql.Result {

	var result sql.Result
	defer func() {
		var count int64
		var err error
		if result != nil {
			count, err = result.RowsAffected()
		}
		// for DML this can be done after, it's not supposed to be slow
		if printQuery {
			if count > 0 {
				log.Info(stringFromSQL(dml, values))
			} else {
				log.Trace(stringFromSQL(dml, values))
			}
		}
		if err != nil {
			log.Error("DML: %s", err.Error())
		} else {
			log.Trace("DML: Rows count: %d", count)
		}
	}()

	values = preprocessValues(vp, values)

	statement, err := connection.Prepare(dml)
	assertNil(err)

	defer statement.Close()

	result, err = statement.Exec(values...)
	assertNil(err)

	return result
}

func runQuery(vp ValuePreprocessor, connection *sql.DB, query string, values []any, printQuery bool, eachRow func(rows *sql.Rows)) {

	var count int64
	var err error
	defer func() {
		if err == nil {
			log.Trace("SELECT: Rows count: %d", count)
		}
	}()

	values = preprocessValues(vp, values)

	statement, err := connection.Prepare(query)
	assertNil(err)

	if printQuery {
		log.Trace(stringFromSQL(query, values))
		//log.Trace("Q: %s", query)
		//log.Trace("V: %v", values)
	}

	defer statement.Close()
	rows, err := statement.Query(values...)
	assertNil(err)
	for rows.Next() {
		count += 1
		eachRow(rows)
	}

	err = rows.Err()
	assertNil(err)
}

func preprocessValues(vp ValuePreprocessor, values []any) []any {
	if vp != nil {
		for i, v := range values {
			converted, newValue := vp(v)
			if converted {
				values[i] = newValue
			}
		}
	}
	/*
		if quoteStrings {
			for i, v := range values {
				if reflect.ValueOf(v).Kind() == reflect.String {
					values[i] = quotedString(v.(fmt.Stringer).String())
				}
				/*
					switch v.(type) {
					case string:
						values[i] = quoteOrNot(v.(string))
					}
				* /
				//		log.Debug("ref: %#v", ref)
			}
		}
	*/
	return values
}

func quotedString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func stringFromSQL(sql string, values []any) string {
	var lines []string
	for _, line := range strings.Split(sql, "\n") {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 && !strings.HasPrefix(trimmed, "--") {
			lines = append(lines, line)
		}
	}
	sql = strings.Join(lines, "\n")
	for i, _ := range values {
		i = len(values) - i - 1
		n := print(i + 1)
		p := "$" + n
		k := "/*" + n + "*/ "
		v := values[i]
		switch reflect.ValueOf(v).Kind() {
		case reflect.String:
			sql = strings.ReplaceAll(sql, p, k+quotedString(print(v)))
		default:
			sql = strings.ReplaceAll(sql, p, k+print(v))
		}
	}
	if len(lines) <= 3 {
		lines = nil
		for _, line := range strings.Split(sql, "\n") {
			trimmed := strings.TrimSpace(line)
			lines = append(lines, trimmed)
		}
		sql = strings.Join(lines, " ")
	}
	return sql
}

type EnsuredTime time.Time

func (t *EnsuredTime) Scan(v any) error {
	if v == nil {
		*t = EnsuredTime(EnsuredTime{})
		return nil
	}
	switch v.(type) {
	case []byte:
		vt, err := lang.LocalTime("2006-01-02 15:04:05", string(v.([]byte)))
		if err != nil {
			return err
		}
		*t = EnsuredTime(vt)
		return nil
	default:
		*t = EnsuredTime(v.(time.Time))
		return nil
	}
}
