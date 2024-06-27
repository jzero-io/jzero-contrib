# jzero-contrib

jzero contrib

## swaggerv2

在线展示 swagger ui 文档

![](https://oss.jaronnie.com/image-20240627175804999.png)

### Usage

将 swagger.json 放在 docs 文件夹下

```go
package main

import (
	"github.com/jzero-io/jzero-contrib/swaggerv2"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	server := rest.MustNewServer(rest.RestConf{
		Port: 8001,
	})
	swaggerv2.RegisterRoutes(server, swaggerv2.WithSwaggerPath("docs"))

	server.Start()
}
```

访问 localhost:8001/swagger

## logtoconsole

在 go-zero 中, 设置日志 mode 为 file 或者 volume 时, 无法在控制台上查看日志, 解决办法

```go
package main

import (
	"github.com/jzero-io/jzero-contrib/logtoconsole"
	"github.com/jzero-io/jzero-contrib/swaggerv2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	logConf := logx.LogConf{
		Mode:     "file",
		Path:     "logs",
		Encoding: "plain",
	}
	server := rest.MustNewServer(rest.RestConf{
		Port: 8001,
		ServiceConf: service.ServiceConf{
			Log: logConf,
		},
	})
	logtoconsole.Must(logConf)
	swaggerv2.RegisterRoutes(server, swaggerv2.WithSwaggerPath("docs"))

	logx.Info("starting server")
	server.Start()
}
```