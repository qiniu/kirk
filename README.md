# 简介
本 SDK 基于 Golang 语言， 用于与 QINIU KIRK 通用计算服务（容器云）开放 API 的编程交互，提供了在开发者业务服务器（服务端或客户端）管理七牛容器云资源的能力。

[![Qiniu Logo](http://open.qiniudn.com/logo.png)](http://qiniu.com/)

| branch | usage | status |
| --- | --- | --- |
| master | (default) stable version| [![Build Status](https://travis-ci.org/qiniu/kirk.svg?branch=master)](https://travis-ci.org/qiniu/kirk) [![GoDoc](https://godoc.org/qiniupkg.com/kirk?status.svg)](https://godoc.org/qiniupkg.com/kirk) |
| develop | develop version | [![Build Status](https://travis-ci.org/qiniu/kirk.svg?branch=develop)](https://travis-ci.org/qiniu/kirk) |

# 安装
> 本 SDK 需要 [go 1.7](https://golang.org/dl/) 以上版本

## 使用 go get 安装
```
$ go get -u qiniupkg.com/kirk/kirksdk
```

## 使用 [glide](https://glide.sh) 安装
- 安装 [glide](https://glide.sh) 包管理工具
- 在项目有中添加一个 import “qiniupkg.com/kirk/kirksdk” 的 .go 源文件，并执行如下命令。glide会自动扫描代码并下载需要的包

```bash
$ cd your_project_dir
$ glide init
$ glide install
```

# 示例
## 创建 App
```golang
import "qiniupkg.com/kirk/kirksdk"

...

cfg := kirksdk.AccountConfig{
	AccessKey: ACCESS_KEY,
	SecretKey: SECRET_KEY,
	Host:      kirksdk.DefaultAccountHost,
}

accountClient := kirksdk.NewAccountClient(cfg)

createdApp, err := accountClient.CreateApp(context.TODO(), "myapp", kirksdk.CreateAppArgs{
	Title:  "title",
	Region: "nq",
})

if err != nil {
// 错误处理
}

fmt.Println(createdApp.URI)
```

## 在 App 下创建 Service
```golang
import "qiniupkg.com/kirk/kirksdk"

...

qcosClient, err := accountClient.GetQcosClient(context.TODO(), createdApp.URI)
if err != nil {
// 错误处理
}

err = qcosClient.CreateService(context.TODO(), "mystack", kirksdk.CreateServiceArgs{
	Name: "myservice",
})
if err != nil {
// 错误处理
}
```

## 列出账号下所有镜像仓库
```golang
import "qiniupkg.com/kirk/kirksdk"

...

accountInfo, err := accountClient.GetAccountInfo(context.TODO())
if err != nil {
// 错误处理
}

indexClient, err := accountClient.GetIndexClient(context.TODO())
if err != nil {
// 错误处理
}

repos, err := indexClient.ListRepo(context.TODO(), accountInfo.Name)
if err != nil {
// 错误处理
}

for _, repo := range repos {
	fmt.Println(repo.Name)
}
```

# 参考
- [qiniupkg.com/kirk](https://godoc.org/qiniupkg.com/kirk)
- [文档](https://developer.qiniu.com/kirk)
