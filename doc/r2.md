# r2 保存方法

## 部署

- 创建R2存储桶

- 将 [worker](./cloudflare/worker.js) 部署到 cloudflare workers

- `变量和机密` 设置`AUTH_TOKEN`为访问密钥

- `绑定` 在worker中绑定R2储存桶,并将`变量名称`设置为`SUB_BUCKET`

## 修改配置文件

- 将`save-method`配置为 `r2`

- 将`worker-url`设置为你的 worker 地址

- 将`worker-token`设置为你的 worker token

## 获取订阅

- 全部订阅

```
https://your-worker-url/storage?key=all.yaml&token=AUTH_TOKEN
```

- 解锁openai的节点

```
https://your-worker-url/storage?key=openai.yaml&token=AUTH_TOKEN
```

- 解锁netflix的节点

```
https://your-worker-url/storage?key=netflix.yaml&token=AUTH_TOKEN
```

- 解锁disney的节点

```
https://your-worker-url/storage?key=disney.yaml&token=AUTH_TOKEN
```

- 解锁youtube的节点

```
https://your-worker-url/storage?key=youtube.yaml&token=AUTH_TOKEN
```
