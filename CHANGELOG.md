# vNext
- 新增app配额查询接口

# Release 2.3.0
- SpecInfo 添加 Regions 字段
- 增加 list grants 接口

# Release 2.2.0
- vendorManaged 应用状态的response 改details为message
- ListRepoTags和GetImageConfig 支持获取imagesize
- spec和app添加Privileges字段
- 为各个 Client 接口添加 GetConfig 方法

# Release 2.1.0
- 新增应用平台user权限的接口
- 修改日志搜索返回结果字段及其Tag

# Release 2.0.0
- 新增app授权和撤销授权功能
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
