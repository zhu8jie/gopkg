package xconfig

import (
	"net"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

// MySQLConf mysql配置结构
type MySQLConf struct {
	User      string `json:"user" yaml:"user"`
	Password  string `json:"password" yaml:"password"`
	Protol    string `json:"protol" yaml:"protol"`
	Host      string `json:"host" yaml:"host"`
	Port      int    `json:"port" yaml:"port"`
	DbName    string `json:"db_name" yaml:"dbname"`
	Charset   string `json:"charset" yaml:"charset"`
	MaxIdle   int    `json:"max_idle" yaml:"max_idle"`
	MaxActive int    `json:"max_active" yaml:"max_active"`
	MaxLife   int    `json:"max_life" yaml:"max_life"`
	ShowSQL   bool   `json:"show_sql" yaml:"show_sql"`
}

func (c *MySQLConf) DSN() string {
	cfg := mysql.NewConfig()
	cfg.User = c.User
	cfg.Passwd = c.Password
	if c.Protol == "" {
		c.Protol = "tcp"
	}
	cfg.Net = c.Protol
	cfg.Addr = net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	cfg.DBName = c.DbName
	cfg.ParseTime = true
	cfg.Loc = time.Local
	return cfg.FormatDSN()
}

// SetMysqlConfig 使用环境变量替换mysql配置参数
func (c *MySQLConf) SetMysqlConfig() {
	host := os.Getenv("DATABASE_SERVER_HOST")
	if host != "" {
		c.Host = host
	}

	port := os.Getenv("DATABASE_SERVER_PORT")
	if p, err := strconv.Atoi(port); err == nil && p > 0 {
		c.Port = p
	}

	dbName := os.Getenv("DATABASE_SERVER_DBNAME")
	if dbName != "" {
		c.DbName = dbName
	}

	user := os.Getenv("DATABASE_SERVER_USERNAME")
	if user != "" {
		c.User = user
	}

	password := os.Getenv("DATABASE_SERVER_PASSWORD")
	if password != "" {
		c.Password = password
	}

	charset := os.Getenv("DATABASE_SERVER_CHARSET")
	if charset != "" {
		c.Charset = charset
	}

	maxActive := os.Getenv("DATABASE_SERVER_MAX_ACTIVE")
	if p, err := strconv.Atoi(maxActive); err == nil && p > 0 {
		c.MaxActive = p
	}

	maxIdle := os.Getenv("DATABASE_SERVER_MAX_IDLE")
	if p, err := strconv.Atoi(maxIdle); err == nil && p > 0 {
		c.MaxIdle = p
	}

	showSQL := os.Getenv("DATABASE_SERVER_SHOW_SQL")
	if showSQL == "true" {
		c.ShowSQL = true
	}
}
