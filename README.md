## Memcached UI

一个简单实用的Memcached Web客户端。

- 以插件方式支持多种Web框架使用Memached的方式，比如Yii框架会对原始的键值进行处理，存入Memcached的键是一个哈希值，而value只是序列化后的结果。
- 不依赖第三方Memcached客户端库
- 支持Basic Auth身份认证方式

------

插件的实现见代码文件：`middleman/middleman/default.go` 和 `middleman/middleman/yii.go`。

------

界面：

![memcached-ui](https://raw.github.com/youngsterxyf/memcached-ui/master/sample.png)

部署：

1. `go get github.com/youngsterxyf/memcached-ui`
2. `cd $GOPATH/src/github.com/youngsterxyf/memcached-ui`
3. `go build mu.go`
4. `cp app.json.example app.json`，并根据需求修改配置信息
5. `nohup ./mu > mu.log 2>&1 &`
6. 访问 http://127.0.0.1:8080
