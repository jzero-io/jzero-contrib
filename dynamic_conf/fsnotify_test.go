package dynamic_conf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/a8m/envsubst"
	"github.com/jaronnie/genius"
	configurator "github.com/zeromicro/go-zero/core/configcenter"
	"github.com/zeromicro/go-zero/core/logx"
)

type TestSt struct {
	Name string `json:"name,"`
}

func TestLocalFsNotify(t *testing.T) {
	logx.MustSetup(logx.LogConf{
		Level:    "info",
		Encoding: "plain",
	})
	ss, err := NewFsNotify("testdata/etc.yaml")
	if err != nil {
		panic(err)
	}
	cc := configurator.MustNewConfigCenter[TestSt](configurator.Config{
		Type: "yaml", // 配置值类型：json,yaml,toml
	}, ss)

	v, err := cc.GetConfig()
	if err != nil {
		panic(err)
	}
	println(v.Name)
}

func TestEnvsubstYaml(t *testing.T) {
	os.Setenv("DatabaseType", "mysql")
	data, err := envsubst.ReadFile("testdata/etc.yaml")
	if err != nil {
		log.Fatalf("envsubst error: %v", err)
	}
	g, err := genius.NewFromType(data, filepath.Ext("testdata/etc.yaml"))
	if err != nil {
		panic(err)
	}
	fileBytes, err := g.EncodeToType(filepath.Ext("testdata/etc.yaml"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(fileBytes))
}

func TestEnvsubstJson(t *testing.T) {
	os.Setenv("DatabaseType", "mysql")
	data, err := envsubst.ReadFile("testdata/etc.json")
	if err != nil {
		log.Fatalf("envsubst error: %v", err)
	}
	g, err := genius.NewFromType(data, filepath.Ext("testdata/etc.json"))
	if err != nil {
		panic(err)
	}
	fileBytes, err := g.EncodeToType(filepath.Ext("testdata/etc.json"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(fileBytes))
}
