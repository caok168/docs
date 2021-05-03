# 如何实现线程安全的map

### Go内建的map类型
* map的类型是map[key]
* key类型的K必须是可比较的
* map[key]函数返回结果可以是一个值，也可以是两个值
* map是无序的，想要保证遍历map时元素有序，可以使用orderedmap

### 使用map的2种常见错误
* 未初始化
* 并发读写

### 如何实现线程安全的map类型
* 加读写锁：扩展map，支持并发读写
* 分片加锁：更高效的并发
    - GetShard是一个关键的方法，能够根据key计算出分片索引

### 应对特殊场景的sync.Map
* 应用场景不多
    - 只会增长的缓存系统中，一个key只写入一次而被读很多次；
    - 多个goroutine为不相交的键值读、写和重写键值对。
* 设计与实现
    - Store方法：设置一个键值对/更新一个键值对
    - Load方法：读取一个key对应的值
    - Delete方法：LoadAndDelete的实现
    - 辅助方法：
        * LoadAndDelete
        * LoadOrStore
        * Range
* 实现优化点
    - 空间换时间
    - 优先从read字段读取、更新、删除
    - 动态调整
    - double-checking
    - 延迟删除