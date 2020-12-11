Based on [goim v2.0](https://github.com/zqkgo/goim-enhanced), add stats ability to monitor runtime.

----

* [单机部署](#单机部署)
  * [启动kafka](#启动kafka)
  * [启动redis](#启动redis)
  * [启动注册中心](#启动注册中心)
  * [启动goim各服务](#启动goim各服务)
* [建立连接](#建立连接)
  * [启动模拟HTTP客户端](https://github.com/zqkgo/goim-enhanced#%E5%90%AF%E5%8A%A8%E6%A8%A1%E6%8B%9Fhttp%E5%AE%A2%E6%88%B7%E7%AB%AF)
  * [使用gRPC客户端](https://github.com/zqkgo/goim-enhanced#%E4%BD%BF%E7%94%A8grpc%E5%AE%A2%E6%88%B7%E7%AB%AF)
* [发送消息](#发送消息)
* [使用Prometheus](https://github.com/zqkgo/goim-enhanced#%E4%BD%BF%E7%94%A8prometheus)
* [使用Grafana](https://github.com/zqkgo/goim-enhanced#%E4%BD%BF%E7%94%A8grafana)



## 单机部署

按顺序启动以下组建或服务。

### 启动kafka

参考 [文档](https://kafka.apache.org/documentation/#quickstart)

```
bin/zookeeper-server-start.sh config/zookeeper.properties
bin/kafka-server-start.sh config/server.properties
```

### 启动redis

```
redis-server
```

### 启动注册中心

```
git clone https://github.com/bilibili/discovery
cd discovery/cmd/discovery
go run . -conf=discovery-example.toml
```

### 启动goim各服务

```
cd cmd/comet
go run . -region=sh -zone=sh001 -deploy.env=dev
cd cmd/logic
go run . -region=sh -zone=sh001 -deploy.env=dev
cd cmd/job
go run . -conf=job-example.toml -region=sh -zone=sh001 -deploy.env=dev
cd cmd/stat
go run . -conf=stat-example.toml -region=sh -zone=sh001 -deploy.env=dev -logtostderr=true
```

## 建立连接

### 启动模拟HTTP客户端

```
cd examples/javascript
go run .
```

访问 http://127.0.0.1:1999/

### 使用gRPC客户端

```
cd examples/grpc_clients
./run.sh
```

## 发送消息

```
curl -d '房间消息' 'http://127.0.0.1:3111/goim/push/room?operation=1000&type=live&room=1000'
```

## 使用Prometheus

*待补充*

## 使用Grafana

*待补充*
