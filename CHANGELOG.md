# vNext
- Service SCALING 状态拆分为 SCALING-UP SCALING-DOWN
- 日志搜索结果添加CollectedAtNano字段 
- 添加 GetWebProxy 方法

# Release 1.2.0
- 添加禁用/启用AP端口的API，并在查看/搜索AP的API返回的端口信息中返回端口的启用状态（启用/禁用）。
- 查看服务时，给出与该服务关联的AP端口信息。
- 为 Service 以及 Job 增加 Confs 域
- 增加 ConfigService 相关 V3 接口
- 为 index 列取tag 接口增加排序，分页和时间参数

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
