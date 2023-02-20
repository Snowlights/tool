# tool
工具


##整体架构

-----------------

<details>
<summary>展开查看</summary>
<pre><code>.
│—— breaker: 熔断部分
    │—— 计数熔断(后续新增滑动窗口或其他方式)
│—— cache: 缓存部分
    │—— redis
│—— common 公共依赖
    │—— idl 接口文件
│—— load_balance 负载均衡
     │—— consistent_hash 均衡算法
│—— parse 解析&反射
    │—— parse 数据解析
│—— stat 统计
    │—— vcpu cpu指标
    │—— vload 负载
    │—— vmemory 内存
│—— vconfig 配置
    │—— apollo 配置中心
│—— vfile 文件
    │—— file 文件基础操作
│—— vlog 日志
│—— vmongo mongodb
│—— vmq 消息队列
    │—— vkafka kafka
│—— vnet 网络
│—— vprometheus prometheus
    │—— metric 指标
    │—— vvollector 采集器
    │—— vmetric prometheus指标
│—— vservice 服务注册发现
    │—— client 服务发现
    │—— common 公用信息
    │—— server 服务注册
    │—— test 测试用例
│—— vsql 数据库
    │—— builder 数据库构建
    │—— scan 数据库反射
│—— vtrace 跟踪
    │—— trace 跟踪
</code></pre>
</details>

### 整体依赖
* redis
  * redis(sdk版本为v6.15.9),仅支持基础命令集和pipeline动作
  * 文档请参考[链接](https://redis.io/docs/)
* mongodb
  * mongodb(驱动采用go.mongodb.org/mongo-driver v1.9.0)
  * 文档请参考[链接](https://www.mongodb.com/docs/)
* kafka
    * 文档请参考[链接](https://kafka.apache.org/documentation/)
* prometheus
    * 文档请参考[链接](https://prometheus.io/docs/introduction/overview/)
* etcd 
    * 文档请参考[链接](https://etcd.io/docs/v3.5/)
* zookeeper
  * 文档请参考[链接](https://zookeeper.apache.org/doc/r3.1.2/index.html)
* consul
  * 文档请参考[链接](https://www.consul.io/docs)
* jaeger
    * 文档请参考[链接](https://www.jaegertracing.io/docs/1.34/)
* apollo
    * 文档请参考[链接](https://www.bookstack.cn/read/ctripcorp-apollo/66fd39d228fadcad.md)


### 模块功能

主入口从service进入，其他模块可单独调用或单独使用

<details>
<summary>展开查看</summary>
<pre><code>.
│—— vservice
  │—— client
      │—— 通过sdk的方式获取到其他服务的信息，并进行调用
      │—— 支持http,grpc,thrift
  │—— common
      │—— 公共信息&配置信息
  │—— server
        │—— 初始化服务注册的地址(etcd,zookeeper)在机器或者镜像的环境变量中
        │—— 服务注册，通过指定服务注册类型(etcd、zookeeper)，将会自动获取该服务的注册地址并进行注册(并自动维护心跳)
        │—— 服务注册时会从apollo拉取(redis、mongo、kafka、jaeger、consul、prometheus等配置信息)，并在服务注册之前初始化这些信息
        │—— 服务注册需要用户提供自己的processor方法来实现，用于实现用户自定义实现的业务逻辑
        │—— 可以通过设置服务注册的配置信息，来实现服务注册的配置信息的更新
  │—— test
        │—— 测试用例
│—— 其他模块
    │—— 服务注册之后，用户可以通过各个模块暴露的方法来调用其他服务(redis等)而无需关心服务依赖的节点
    │—— prometheus及jaeger也可在用户服务中进行使用
</code></pre>
</details>