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

待完善：
- 添加sing-box生成格式