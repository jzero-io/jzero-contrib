package modelx

import (
	"testing"

	"github.com/huandu/go-assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestConfig(t *testing.T) {
	config := &ModelxConfig{}
	conf.MustLoad("./config.yaml", config)

	assert.Equal(t, config.DatabaseType, "mysql11111")
	assert.Equal(t, config.Mysql.Host, "127.0.0.1")
	assert.Equal(t, config.Mysql.Username, "hhh")
	assert.Equal(t, config.Sqlite.Path, "testPath")
}
