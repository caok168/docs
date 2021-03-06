# etcd的前世今生：为什么Kubernetes使用etcd？

一个协调服务，理想状态下大概需要满足以下五个目标：
* 可用性角度：高可用
* 数据一致性角度：提供读取“最新”数据的机制
* 容量角度：低容量、仅存储关键元数据配置
* 功能：增删该查，监听数据变化的机制
* 运维复杂度：可维护性

CoreOS为什么选择了从0到1开发一个新的协调服务？
* 当时其实是有ZooKeeper的，但是为什么没有不用ZooKeeper呢？
    * 从高可用、数据一致性、功能这三个角度来说，ZooKeeper是满足CoreOS诉求的，然而当时的ZooKeeper不支持通过API安全的变更成员，需要人工修改一个个节点的配置，并重启进程。
    * 若变更姿势不正确，可能出现脑裂等严重故障，维护成本相对较高。
    * ZooKeeper是用Java开发的，部署较繁琐，占用较多的内存资源，同时ZooKeeper RPC的序列化机制用的是Jute，自己实现的 RPC API。无法使用Curl之类的常用工具与之互动，CoreOS期望使用比较简单的HTTP + JSON。

### 1、高可用和数据一致性

* 单副本存在单点故障，而多副本又引入数据一致性问题
* 解决数据一致性问题，需要引入一个共识算法，确保各节点数据一致性，可以容忍一定节点故障。
    * 常见的共识算法有Paxos、ZAB、Raft等。
    * CoreOS选择易于理解的Raft算法，将复杂的一致性问题分解成Leader选举、日志同步、安全性三个相对独立的子问题。

### 2、数据模型（Data Model）和API
* 数据模型参考了ZooKeeper，使用的是基于目录的层次模式
* API相比ZooKeeper来说，使用了简单、易用的REST API。

### 3、key-value存储引擎上
* Zookeeper使用的是Concurrent HashMap
* etcd使用的则是简单内存树，它的节点数据精简后，含节点路径、值、孩子节点信息。
    * 这是一个典型的低容量设计，数据全放在内存，无需考虑数据分片，只能保存key的最新版本，简单易实现。

### 4、可维护性
* Raft算法提供了成员变更算法，可基于此实现成员在线、安全变更，同时此协调服务使用Go语言编写，无依赖，部署简单。

## etcd v2存在的问题

### 数据模型/Data Model
* 不支持范围查询
* 不支持分页

### HTTP/1.x JSON API
* JSON编解码耗CPU
* 不支持压缩、流量大
* 不支持连接多路复用

### TTL
* 大量key TTL相同时，续期开销大，扩展性较低

### Memory Tree
* 不支持多版本
* 内存开销大
* 快照备份开销较大

### 事务
* 不支持多key事务

### Watch
* 事件不可靠，可能会丢失

## etcd v3解决问题
### 内存开销、Watch事件可靠性、功能局限上
* 通过引入B-tree、boltdb实现一个MVCC数据库，数据模型从层次型目录结构改为扁平的key-value，提供稳定可靠的事件通知，实现了事务，支持多key原子操作，同时基于boltdb的持久化存储，显著降低了etcd的内存占用、避免了etcd v2定期生成快照时的昂贵的资源开销
### 性能上
* etcd v3使用了gRPC API，使用protobuf定义消息，消息编码性能相比JSON超过了2倍以上
* 通过HTTP/2.0多路复用机制，减少了大量watcher等场景下的连接数
* 使用Lease优化TTL机制，每个Lease具有一个TTL，相同的TTL的key关联一个Lease，Lease过期的时候自动删除相关联的所有key，不再需要为每个key单独续期。
* etcd v3支持范围、分页查询，可避免大包等expensive request。

## etcd与redis的区别
* 数据复制：
    * Redis是主备异步复制，etcd使用的是Raft，前者可能会丢数据，为了保证读写一致性，etcd读写性能相比redis差距比较大
* 数据分片：
    * Redis有各种集群版解决方案，可以承载上T数据，存储的一般是用户数据
    * etcd定位是个低容量的关键元数据存储，db大小一般不超过8g。
* 存储引擎和API：
    * Redis内存实现了各种丰富数据结构
    * etcd仅是kv API，使用的是持久化存储boltdb。



