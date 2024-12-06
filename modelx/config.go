package modelx

import (
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ModelxConfig struct {
	DatabaseType string `json:"databaseType,default=mysql"`

	Mysql  MysqlConf  `json:"mysql,"`
	Sqlite SqliteConf `json:"sqlite,"`
}

type MysqlConf struct {
	DatabaseConf
}

type SqliteConf struct {
	Path string `json:"path,default=data.db"`
}

type DatabaseConf struct {
	Host     string `json:"host,default=localhost"`
	Port     int    `json:"port,default=3306"`
	Username string `json:"username,default=root"`
	Password string `json:"password,default=123456"`
	DbName   string `json:"dbName,default=test"`
}

func BuildDataSource(c ModelxConfig) string {
	switch c.DatabaseType {
	case "mysql":
		sqlbuilder.DefaultFlavor = sqlbuilder.MySQL
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.Mysql.Username,
			c.Mysql.Password,
			c.Mysql.Host+":"+cast.ToString(c.Mysql.Port),
			c.Mysql.DbName)
	case "sqlite":
		sqlbuilder.DefaultFlavor = sqlbuilder.SQLite
		return c.Sqlite.Path
	}
	return ""
}

func MustSqlxConn(c ModelxConfig) sqlx.SqlConn {
	sqlConn := sqlx.NewSqlConn(c.DatabaseType, BuildDataSource(c))
	_, err := sqlConn.Exec("select 1")
	if err != nil {
		panic(err)
	}
	return sqlConn
}
