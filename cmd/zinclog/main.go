package main

import (
	"path/filepath"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"

	"github.com/jzero-io/jzero-contrib/zinclog"
)

var (
	cfgFile = filepath.Join("etc", "etc.yaml")
	cfg     = zinclog.ZincLogstash{}
)

func main() {
	group := service.NewServiceGroup()

	err := conf.Load(cfgFile, &cfg)
	logx.Must(err)

	logstash := zinclog.NewZincLogstash(&cfg)
	group.Add(logstash)

	group.Start()

	select {}
}
