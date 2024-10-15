package xmysql

import (
	"fmt"
	"io"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

// MySQLConf mysql配置结构
type MySQLXormConf struct {
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

func NewXorm(conf *MySQLXormConf, writer io.Writer) (*xorm.Engine, error) {
	var err error

	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", conf.User, conf.Password, conf.Host, conf.Port, conf.DbName, conf.Charset)
	engine, err := xorm.NewEngine("mysql", source)
	if err != nil {
		return nil, err
	}

	engine.SetMaxIdleConns(conf.MaxIdle)
	engine.SetMaxOpenConns(conf.MaxActive)
	engine.SetConnMaxLifetime(time.Second * time.Duration(conf.MaxLife))
	if writer != nil {
		engine.SetLogger(log.NewSimpleLogger(writer))
	}
	engine.ShowSQL(conf.ShowSQL)

	return engine, engine.Ping()
}
