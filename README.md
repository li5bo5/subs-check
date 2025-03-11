# 订阅合并转换检测工具

兴趣爱好，目前不建议使用本项目

## 项目介绍

 - 详情看[项目概述.md](./doc/项目概述.md)

## 原项目

- https://github.com/bestruirui/BestSub
- https://github.com/beck-8/subs-check


## 修改内容

- 移除MihomoApi更新订阅功能
- 移除 GitHub Container Registry 发布配置
- Docker Hub 镜像发布
- 移除config.example.yaml中worker-url和worker-token内容
- 移除config.example.yaml中mihomo-api-url和mihomo-api-secret内容
- 移除config.go中MihomoApiUrl和MihomoApiSecret字段
- 移除updatesubs.go中getversion字段
- 抄作业[多线程获取订阅]https://github.com/bestruirui/BestSub/commit/09b569e4c22d3ec1890f9bf8185405d1c338f23b

## 保存方式

- 本地
- Gist

## 生成订阅格式

- mihomo
- base64

## TODO

- 生成sing-box订阅
- 移除重命名
- 其他

## 使用方法

- 1.Docker部署
- 2.进入主机文件夹./subs-check/config，找到config.yaml
- 3.将config.yaml中的sub-urls下的example.com替换为你的订阅链接
- 4.保存config.yaml文件
- 5.重新运行docker容器
- 6.在主机8299端口查看订阅
```bash
curl http://localhost:8299/all.yaml
```
> 如果拉取订阅速度慢，可使用通用的 `HTTP_PROXY` `HTTPS_PROXY` 环境变量加快速度；此变量不会影响节点测试速度

### docker运行

```bash
docker run -d --name subs-check -p 8299:8299 -v ./subs-check/config:/app/config  -v ./subs-check/output:/app/output --restart always li5bo5/subs-check:latest
```

## 订阅使用方法

原作者推荐直接裸核运行 tun 模式 

原作者写的Windows下的裸核运行应用 [minihomo](https://github.com/bestruirui/minihomo)

- 下载[base.yaml](./doc/base.yaml)
- 将文件中对应的链接改为自己的即可
- 只修改一个url的内容，然后这个base.yaml当作本地订阅文件导入OpenClash，就可以用了

例如:

```yaml
proxy-providers:
  ProviderALL:
    url: _http://127.0.0.1:8299/all.yaml_ #将此处替换为自己的链接
    type: http
    interval: 600
    proxy: DIRECT
    health-check:
      enable: true
      url: http://www.google.com/generate_204
      interval: 60
    path: ./proxy_provider/ALL.yaml
```