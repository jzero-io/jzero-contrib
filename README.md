# jzero-contrib

jzero contrib是一个 jzero 框架的扩展库，包含了一些 jzero 的扩展功能, go-zero 用户也可以利用该生态增强开发体验和效率。

## cache

统一封装的 cache 缓存层, 建立在 `github.com/zeromicro/go-zero/core/stores/cache` 之上, 提供更多功能.

## condition

条件构造器, 基于 `sqlbuilder` 封装了常用条件构造, 并结合 `go-zero` 的 `model` 层 和 `sqlx` 框架, 针对常用业务进一步封装, 强化开发体验.

## dynamic_conf

动态配置, 基于 `github.com/zeromicro/go-zero/core/configcenter` 和 `github.com/fsnotify/fsnotify`, 提供一种本地配置文件动态更新机制.

## swaggerv2

一键为你的服务提供 swagger 在线访问 UI 页面.

