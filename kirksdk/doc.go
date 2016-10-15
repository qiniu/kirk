/*
包 qiniupkg.com/kirk/kirksdk 提供了在开发者业务服务器（服务端或客户端）管理七牛容器云资源的能力。

首先，开发者需要配置七牛云账户的 AccessKey/SecretKey（在 https://portal.qiniu.com/ 中可以查看）。

    import (
        "golang.org/x/net/context"
        "qiniupkg.com/kirk/kirksdk"
    )

    cfg := kirksdk.AppConfig{
        AccessKey: "UserAccountAccessKey",
        SecretKey: "UserAccountSecretKey",
        Host: kirksdk.DefaultAppHost
    }
    appClient := kirksdk.NewAppClient(cfg)

    indexClient, err := appClient.GetIndexClient(context.TODO())

    subAppClient, err := appClient.GetSubAppClient(context.TODO(), "kirk-test.new-app")

在构建了 appClient/indexClient/subAppClient 后，开发者可通过 client 调用账号、镜像空间以及各区域应用下容器、网络资源管理接口。


前往 http://kirk-docs.qiniu.com/apidocs/?go 可以查看更多 API 和 SDK 文档。
*/
package kirksdk
