# Pool:性能提升大杀器

sync.Pool数据类型用来保存一组可独立访问的**临时**对象。

临时说明了，它池化的对象会在未来的某个时候被毫无预兆地移除掉。而且，如果没有别的对象引用这个被移除的对象的话，这个被移除的对象就会被垃圾回收掉。

因为Pool可以有效地减少新对象的申请，从而提高程序性能，所以Go内部库也用到了sync.Pool，比如fmt包，它会使用一个动态大小的buffer池做输出缓存，当大量的goroutine并发输出的时候，就会创建比较多的buffer，并且在不需要的时候回收掉。

### sync.Pool的特点
* sync.Pool本身就是线程安全的，多个goroutine可以并发地调用它的方法存取对象
* sync.Pool不可在使用之后再复制使用

### 使用方法
* New：当调用Pool的Get方法从池中获取元素，没有更多的空闲元素可返回时，就会调用这个New方法来创建新的元素。如果你没有设置New字段，没有更多的空闲元素可返回时，Get方法将返回nil，表明当前没有可用的元素
* Get：如果调用这个方法，就会从Pool取走一个元素
* Put：这个方法用于将一个元素返还给Pool，Pool会把这个元素保存到池中，并且可以复用

### sync.Pool的坑
* 内存泄露
* 内存浪费

### 第三方库
* bytebufferpool
* oxtoacart/bpool

### 常用场景
* 连接池
    - 标准库中的http client池
    - TCP连接池
    - 数据库连接池
    - memcached client连接池
* Worker Pool
    - gammazero/workerpool
    - ivpusic/grpool
    - dpaks/goworkers


