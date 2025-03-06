# local 保存方法

## Docker部署

```
docker run -d --name subs-check -p 8299:8199 -v ./subs-check/config:/app/config  -v ./subs-check/output:/app/output --restart always li5bo5/subs-check:latest
```

## 修改配置文件

- 将`save-method`配置为 `local`

- ./subs-check/config为配置文件config.yaml的目录

- ./subs-check/output为订阅保存的目录