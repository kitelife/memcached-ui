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
