/*
包 qiniupkg.com/kirk 是七牛容器云产品 KIRK 相关的 Golang 开源代码库，包含对接产品 API 接口的 kirksdk 等。

和七牛云其他产品相同，管理容器云资源需要密钥 AccessKey/SecretKey，请注意保证账号密钥的安全。
*/
package kirk

import (
	_ "qiniupkg.com/kirk/kirksdk"
)
