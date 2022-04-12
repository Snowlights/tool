package vsql

import (
	"database/sql"
	"fmt"
	"strings"
)

const (
	mysqlDriver    = "mysql"
	dataNameFormat = "%s:%s@tcp(%s:%d)/%s?%s"

	dsnCharacter = "&"
	dsnFormat    = "%s=%s"
)

type DSN func() string // DSN is a function that returns a DSN string

var (
	utf8mb4DSN = DSN(func() string {
		return fmt.Sprintf(dsnFormat, "charset", "utf8mb4")
	})
)

type Conf struct {
	driver   string
	dbName   string
	user     string
	password string
	host     string
	port     int
	dsn      []DSN
}

func NewConf(dbName, user, pwd, host string, port int) *Conf {
	return &Conf{
		driver:   mysqlDriver,
		dbName:   dbName,
		user:     user,
		password: pwd,
		host:     host,
		port:     port,
	}
}

func (c *Conf) WithDSN(dsn ...DSN) *Conf {
	c.dsn = append(c.dsn, dsn...)
	return c
}

func (c *Conf) Open() (*sql.DB, error) {
	db, err := c.openDB()
	if nil != err {
		return nil, err
	}

	return db, db.Ping()
}

func (c *Conf) formatDSN() string {
	var dsn []string
	for _, d := range c.dsn {
		dsn = append(dsn, d())
	}
	return strings.Join(dsn, dsnCharacter)
}

func (c *Conf) formatDataSourceName() string {
	return fmt.Sprintf(dataNameFormat, c.user, c.password, c.host, c.port, c.dbName, c.formatDSN())
}

func (c *Conf) openDB() (*sql.DB, error) {
	return sql.Open(c.driver, c.formatDataSourceName())
}
