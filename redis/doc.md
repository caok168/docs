# 

Redis有一种底层数据结构，叫压缩列表，是一种非常节省内存的结构。
* 压缩列表之所以节省内存，就在于它是用一系列连续的entry保存数据。
* 不需要额外的指针进行连接。

Redis基于压缩列表实现了List、Hash和Sorted Set这样的集合类型，这样做的最大好处就是节省了dictEntry的开销。

### Redis的底层数据结构
* 简单动态字符串
* 双向链表
* 压缩链表
* 哈希表
* 跳表
* 整数数组

### Redis数据类型和底层数据结构的对应关系
* String 
    - 简单动态字符串 SDS
* List
    - 双向链表
    - 压缩链表
* Hash
    - 压缩链表
    - 哈希表
* Sorted Set
    - 压缩链表
    - 跳表
* Set
    - 哈希表
    - 整形数组

### 组织结构
Redis使用了一个哈希表来保存所有键值对 （全局哈希表）

* 一个哈希表，其实就是一个数组，数组中的每个元素称为一个哈希桶。所以一个哈希表是由多个哈希桶组成的，每个哈希桶中保存了键值对数据。
* 哈希桶中的元素保存的并不是值本身，而是指向具体值的指针。也就是说，不管值是String，还是集合类型，哈希桶中的元素都是指向它们的指针。
* 哈希桶中的entry元素中保存了*key和 *value指针，分别指向了实际的键和值。


### 为什么哈希表操作变慢了

* Redis解决哈希冲突的方式，是链式哈希。指同一个哈希桶中的多个元素用一个链表来保存，它们之间依次用指针连接。

* 如果哈希冲突越来越多，这个链就越来越长，元素查找耗时长，效率降低。
* 所以Redis会对哈希表做rehash操作。rehash就是增加现有的哈希桶的数量，让逐渐增多的entry元素能在更多的桶之间分散保存，减少单个桶中的元素数量，从而减少单个桶中的冲突。
* rehash过程分三步
    - 给哈希表2分配更大的空间，例如是当前哈希表1大小的2倍；
    - 把哈希表1中的数据重新映射并拷贝到哈希表2中；
    - 释放哈希表1的空间。
* 避免大量数据拷贝，造成Redis线程阻塞，无法服务其他请求，Redis采用了渐进式rehash


查找的时间复杂度
|名称 | 时间复杂度|
|:-----|:------|
|哈希表 | O(1)|
|跳表 | O(logN)|
|双向链表 | O(n)|
|压缩链表 | O(n)|
|整数数组 | O(n)|


**问题：**

整数数组和压缩列表在查找时间复杂度方面并没有很大的优势，那为什么Redis还会把它们作为底层数据结构呢？
* 内存利用率，数组和压缩列表都是非常紧凑的数据结构，它比链表占用的内存要更少。Redis是内存数据库，大量的数据存到内存中，此时需要做尽可能的优化，提高内存的利用率
* 数组对CPU高速缓存支持更友好，所以Redis在设计时，集合数据元素较少情况下，默认采用内存紧凑排列的方式存储，同时利用CPU高速缓存不会降低访问速度。当数据元素超过设定阈值后，避免查询时间复杂度太高，转为哈希表和跳表数据结构存储，保证查询效率。
