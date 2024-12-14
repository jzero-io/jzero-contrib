package dynamic_conf

import (
	"testing"

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

	// 如果想监听配置变化，可以添加 listener
	cc.AddListener(func() {
		v, err := cc.GetConfig()
		if err != nil {
			panic(err)
		}
		println(v.Name)
	})

	select {}
}
