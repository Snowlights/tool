package vsql

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	mysqlDriver    = "mysql"
	dataNameFormat = "%s:%s@tcp(%s)/%s?%s"

	dsnCharacter = "&"
	dsnFormat    = "%s=%s"

	defaultMaxIdleConn  = 128
	defaultMaxOpenConn  = 256
	defaultMaxLifetime  = time.Hour * 24
	defaultTimeout      = time.Second * 3
	defaultReadTimeout  = time.Second * 30
	defaultWriteTimeout = time.Second * 30
)

type DSN func() string // DSN is a function that returns a DSN string

var (
	utf8mb4DSN = DSN(func() string {
		return fmt.Sprintf(dsnFormat, "charset", "utf8mb4")
	})
	collationDSN = DSN(func() string {
		return fmt.Sprintf(dsnFormat, "collation", "utf8mb4_unicode_ci")
	})
	locDSN = DSN(func() string {
		return fmt.Sprintf(dsnFormat, "loc", url.QueryEscape("Asia/Shanghai"))
	})
	parseTimeDSN = DSN(func() string {
		return fmt.Sprintf(dsnFormat, "parseTime", "true")
	})
)

type Conf struct {
	driver   string
	dbName   string
	user     string
	password string
	host     string
	dsn      []DSN
}

func NewConf(dbName, user, pwd, host string) *Conf {
	return &Conf{
		driver:   mysqlDriver,
		dbName:   dbName,
		user:     user,
		password: pwd,
		host:     host,
	}
}

func (c *Conf) WithDSN(dsn ...DSN) *Conf {
	c.dsn = append(c.dsn, dsn...)
	return c
}

func (c *Conf) Open() (*sql.DB, error) {
	db, err := sql.Open(c.driver, c.formatDataSourceName())
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
	// like: root:123456@tcp(127.0.0.1:3312)/data?charset=utf8mb4&timeout=3000ms&readTimeout=3000ms&writeTimeout=3000ms
	return fmt.Sprintf(dataNameFormat, c.user, c.password, c.host, c.dbName, c.formatDSN())
}
