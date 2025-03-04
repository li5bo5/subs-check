# 更新历史

## 2024-03-03-1

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

## 2024-03-03-2
### 修改
- 移除MihomoApi更新订阅功能
- CI/CD 配置调整
  - 移除 GitHub Container Registry 发布配置
  - 仅保留 Docker Hub 镜像发布
  - 保持多平台构建支持 (linux/amd64, linux/arm64, linux/arm/v7) 

## 2024-03-03-3
### 推送github
- 推送github

## 2024-03-03-4
### 小修
- 移除config.example.yaml中worker-url和worker-token内容
- 移除config.example.yaml中mihomo-api-url和mihomo-api-secret内容
- 移除config.go中MihomoApiUrl和MihomoApiSecret字段

## 2024-03-03-5
### 小修
- 日志：未配置 MihomoApiUrl，跳过更新
- 移除updatesubs.go中getversion字段
- 移除MihomoApi更新订阅功能

## 2024-03-04-1
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

## 2024-03-04-2
### 配置更新
1. 更新 Dockerfile 基础镜像
   - 升级 Go 版本到最新稳定版 1.24.0
   - 使用 golang:1.24.0-alpine 作为构建基础镜像
   - 移除旧版本注释

待完善：
- 添加sing-box生成格式