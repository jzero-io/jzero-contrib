package modelx

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/huandu/go-assert"
)

func TestConfig(t *testing.T) {
	configJson := `{
  "databasetype": "mysql11111",
  "Mysql": {
    "Host": "127.0.0.1",
	"USernaMe": "hhh"
  },
  "sqlite": {
	"path": "testPath",
  }
}`
	config := &ModelxConfig{}

	err := json.Unmarshal([]byte(configJson), config)
	if err != nil {
		fmt.Println(err)
		return
	}
	assert.Equal(t, config.DatabaseType, "mysql11111")
	assert.Equal(t, config.Mysql.Host, "127.0.0.1")
	assert.Equal(t, config.Mysql.Username, "hhh")
	assert.Equal(t, config.Sqlite.Path, "testPath")
}
