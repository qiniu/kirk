# vNext
- 为 Service 以及 Job 增加 Confs 域
- 增加 ConfigService 相关 V3 接口
- 为indexd列举tag增加分页接口

# Release 1.1.0
- AccountClient 相关 API 使用 appd V3 接口
- 为异步 API 添加同步版本
- 为 port 添加创建和更新时间
- 为 job 添加 deps 域
- 支持 domain 鉴权
- 为 ServiceInspect 添加 progress 域
- 修复BUG: 非 nil 空 slice 未被包含在发送 request body 中

# Release 1.0.0
- Initial release
