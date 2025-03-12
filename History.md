# 更新历史

## 2024-03-12-1
### 修复构建错误
1. 修复函数命名冲突问题
   - 在 `utils/notify.go` 中修改常量名称，避免与函数名冲突
   - 将 `NotifyStart`、`NotifyResult` 和 `NotifyError` 常量重命名为 `TYPE_NotifyStart`、`TYPE_NotifyResult` 和 `TYPE_NotifyError`
   - 保持函数名称不变，确保API兼容性
2. 优化Dockerfile
   - 将CMD指令修改为JSON数组格式 `CMD ["/app/main"]`
   - 解决OS信号处理相关的潜在问题
   - 提高容器运行的稳定性

## 2024-03-11-1
### 添加通知渠道功能
1. 添加 `utils/notify.go` 文件实现通知功能
   - 使用 Apprise API 支持 100+ 种通知渠道
   - 支持三种通知类型：启动通知、结果通知、错误通知
   - 可配置不同情况下是否发送通知
2. 更新配置文件
   - 在 `config/config.go` 中添加通知相关配置字段
   - 在 `config/config.example.yaml` 中添加通知配置示例
3. 在 `main.go` 中集成通知功能
   - 程序启动时发送启动通知
   - 检测完成后发送结果通知
   - 出现错误时发送错误通知
4. 添加节点统计信息
   - 统计总节点数、可用节点数和可用率
   - 在日志和通知中显示统计信息

## 待完善：
- 记录变更内容
- 升级go.mod中组件版本
- 优化内存
- 添加通知渠道https://github.com/beck-8/subs-check/commit/4025634ce07b3cfef638587358a71dcf10fce339
- 添加sing-box生成格式

## 2024-03-08-3
### 移除 keep-success-proxies 功能
1. 移除配置项和相关实现
   - 从 `config/config.example.yaml` 移除 `keep-success-proxies` 配置项
   - 从 `config/config.go` 移除 `KeepSuccessProxies` 字段和 `GlobalProxies` 变量
   - 从 `check/check.go` 移除相关功能实现
   - 从 `main.go` 移除节点保存逻辑

## 2024-03-08-2

### 移除节点重命名功能

- 删除 `proxy/rename.go` 文件，移除节点重命名的主要实现
- 删除 `proxy/info.go` 文件，移除获取节点位置信息的功能
- 修改 `check/check.go` 中的 `updateProxyName` 方法，移除节点重命名相关代码，仅保留速度信息的添加
- 从配置文件中移除 `rename-node` 选项

## 2025-03-08-1
### 日志优化
1. 优化 utils/updatesubs.go 文件
   - 简化日志输出格式
   - 统一显示"更新完成，共 X 个节点"的提示
   - 移除中间状态的日志输出
   - 保持错误日志的完整性

## 2025-03-06-1
### 更换端口为8299
### 订阅获取优化
1. 优化 proxy/get.go 文件
   - https://github.com/bestruirui/BestSub/commit/09b569e4c22d3ec1890f9bf8185405d1c338f23b
   - 添加并发获取订阅功能
   - 使用 goroutine 池处理订阅链接
   - 通过 channel 进行任务分发和结果收集
   - 使用 sync.WaitGroup 确保任务完成
   - 动态计算最优并发数
   - 添加更详细的日志记录
   - 优化错误处理机制
   - 重构代码结构，提取 processSubscription 函数
   - 保持原有的调试日志和错误处理特性

2. 性能改进
   - 使用缓冲通道减少阻塞
   - 优化内存分配
   - 通过配置控制并发数量
   - 实际并发数不超过订阅链接数量

## 2025-03-04-8
### Dockerfile网络优化
1. 添加Go模块代理配置
   - 设置 GOPROXY 环境变量
   - 使用 goproxy.cn 作为国内代理
   - 优化依赖下载速度

## 2025-03-04-7
### Dockerfile构建修复
1. 修复Go版本兼容性问题
   - 确保使用 golang:1.24.0-alpine 基础镜像
   - 添加 GO111MODULE=on 环境变量
   - 移除不必要的 GOTOOLCHAIN 设置

## 2025-03-04-6
### Dockerfile版本修复
1. 修复Go版本兼容性问题
   - 确保使用 golang:1.24.0-alpine 基础镜像
   - 添加 GO111MODULE=on 环境变量
   - 移除不必要的 GOTOOLCHAIN 设置

## 2025-03-04-5
### Dockerfile CI/CD 适配
1. 优化构建环境
   - 添加包索引更新
   - 通过环境变量设置构建参数

## 2025-03-04-4
### Dockerfile构建修复
1. 添加必要的构建依赖
   - 增加 linux-headers 支持
   - 添加 git 和 make 工具
   - 优化依赖安装顺序

## 2025-03-04-3
### Dockerfile优化
1. 优化构建流程
   - 分离依赖安装和代码构建步骤
   - 添加 CGO_ENABLED=0 确保静态链接
   - 优化构建顺序，提高缓存利用率

2. 运行环境配置
   - 优化时区设置方式
   - 添加健康检查机制
   - 简化容器运行配置

## 2025-03-04-2
### 配置更新
1. 更新 Dockerfile 基础镜像
   - 升级 Go 版本到最新稳定版 1.24.0
   - 使用 golang:1.24.0-alpine 作为构建基础镜像
   - 移除旧版本注释

## 2025-03-04-1
### 代码优化
1. 优化 utils/updatesubs.go
   - 移除未使用的版本相关代码和结构体
   - 优化 UpdateSubs 函数的日志输出
   - 添加空订阅列表检查
   - 添加更新完成时的订阅数量统计
   - 更换代码中的仓库地址

### 配置文件更新
1. 更新 .goreleaser.yaml
   - 添加 version: 2 声明
   - 修复配置文件版本警告
   - 保持其他构建配置不变

2. 优化 Dockerfile
   - 指定具体的基础镜像版本：golang:1.21-alpine
   - 添加必要的构建依赖：gcc 和 musl-dev
   - 优化构建流程，分步执行依赖下载和验证
   - 使用 JSON 数组格式的 CMD 指令，提高稳定性

## 2025-03-03-5
### 小修
- 日志：未配置 MihomoApiUrl，跳过更新
- 移除updatesubs.go中getversion字段
- 移除MihomoApi更新订阅功能

## 2025-03-03-4
### 小修
- 移除config.example.yaml中worker-url和worker-token内容
- 移除config.example.yaml中mihomo-api-url和mihomo-api-secret内容
- 移除config.go中MihomoApiUrl和MihomoApiSecret字段

## 2025-03-03-2
### 修改
- 移除MihomoApi更新订阅功能
- CI/CD 配置调整
  - 移除 GitHub Container Registry 发布配置
  - 仅保留 Docker Hub 镜像发布

## 2025-03-03-1
### 移除功能
1. 移除解锁检测功能
   - 移除了check/check.go中的解锁检测相关代码和注释
   - 保留了基本的Google和Cloudflare检测功能
   - Result结构体中只保留了Proxy、Google和Cloudflare字段

2. 移除特定保存方式
   - 删除了save/method/webdav.go文件
   - 修改了save/save.go中的chooseSaveMethod函数，移除了WebDAV相关的保存方法
   - 现在只支持本地(local)和GitHub Gist两种保存方式

### 配置文件更新
1. 更新了config/config.example.yaml中的保存方法说明
2. 更新了README.md文件
   - 移除了解锁检测相关的功能说明
   - 更新了保存方法的说明，只保留本地和Gist两种方式
   - 移除了WebDAV相关的配置说明

### 代码优化
1. 简化了代码结构，移除了不必要的功能
2. 减少了项目依赖
3. 提高了代码的可维护性

### 保留的核心功能
- 检测节点可用性
- 合并多个订阅
- 转换订阅格式
- 节点去重和重命名
- 节点测速
- 对外提供HTTP服务