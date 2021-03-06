# 知识点汇总
## golang内存分配
golang运行时的内存分配算法主要源自Google为C语言开发的TCMalloc算法，全称Thread-Caching Malloc。

核心思想就是把内存分为多级管理，从而降低锁的粒度。它将可用的堆内存采用二级分配的方式进行管理：每个线程都会自行维护一个独立的内存池，进行内存分配时优先从该内存池中分配，当内存池不足时才会向全局内存池申请，以避免不同线程对全局内存池的频繁竞争。

### 基础概念
Go在程序启动的时候，会先向操作系统申请一块内存（注意这时还只是一段虚拟的地址空间，并不会真正地分配内存），切成小块后自己进行管理。

申请到的内存被分配了三个区域，在X64上分别是512MB，16GB，512GB大小。

* arena区域（512GB）
    - 就是我们所谓的堆区，Go动态分配的内存都是在这个区域，它把内存分割成8KB大小的页，一些页组合起来称为mspan。
* bitmap区域（16GB）
    - 标识arena区域哪些地址保存了对象，并且用4bit标志位标识对象是否包含指针、GC标记信息。
    - bitmap中一个byte大小的内存对应arena区域中4个指针大小（指针大小为8B）的内存，所以bitmap区域的大小是512GB/(4*8B)=15GB。
* spans区域（512MB）
    - 存放mspan的指针，每个指针对应一页，所以spans区域的大小就是512GB/8KB*8B=512MB。

### 内存管理单元
mspan：Go中内存管理的基本单元，是由一片连续的8KB的页组成的大块内存。

mspan是一个包含起始地址、mspan规格、页的数量等内容的双端链表。

### 内存管理组件
内存分配由内存分配器完成。分配器由3种组件构成：mcache，mcentral，mheap。

* mcache:每个工作线程都会绑定一个mcache，本地缓存可用的mspan资源，这样就可以直接给Goroutine分配，因为不存在多个Goroutine竞争的情况，所以不会消耗锁资源。
* mcentral：为所有mcache提供切分好的mspan资源。每个central保存一种特定大小的全局mspan列表，包括已分配出去的和未分配出去的。 每个mcentral对应一种mspan，而mspan的种类导致它分割的object大小不同。当工作线程的mcache中没有合适（也就是特定大小的）的mspan时就会从mcentral获取。
* mheap：代表Go程序持有的所有堆空间，Go程序使用一个mheap的全局对象_mheap来管理堆内存。

当mcentral没有空闲的mspan时，会向mheap申请。而mheap没有资源时，会向操作系统申请新内存。mheap主要用于大对象的内存分配，以及管理未切割的mspan，用于给mcentral切割成小对象。

### 分配流程
变量是在栈上分配还是在堆上分配，是由逃逸分析的结果决定的。通常情况下，编译器是倾向于将变量分配到栈上的，因为它的开销小，最极端的就是"zero garbage"，所有的变量都会在栈上分配，这样就不会存在内存碎片，垃圾回收之类的东西。

Go的内存分配器在分配对象时，根据对象的大小，分成三类：小对象（小于等于16B）、一般对象（大于16B，小于等于32KB）、大对象（大于32KB）。

#### 大体上的分配流程：
* 32KB 的对象，直接从mheap上分配；

* <=16B 的对象使用mcache的tiny分配器分配；
(16B,32KB] 的对象，首先计算对象的规格大小，然后使用mcache中相应规格大小的mspan分配；
* 如果mcache没有相应规格大小的mspan，则向mcentral申请
* 如果mcentral没有相应规格大小的mspan，则向mheap申请
* 如果mheap中也没有合适大小的mspan，则向操作系统申请

### 总结
* Go在程序启动时，会向操作系统申请一大块内存，之后自行管理。
* Go内存管理的基本单元是mspan，它由若干个页组成，每种mspan可以分配特定大小的object。
* mcache, mcentral, mheap是Go内存管理的三大组件，层层递进。mcache管理线程在本地缓存的mspan；mcentral管理全局的mspan供所有线程使用；mheap管理Go的所有动态分配内存。
* 极小对象会分配在一个object中，以节省资源，使用tiny分配器分配内存；一般小对象通过mspan分配内存；大对象则直接由mheap分配内存。

## golang垃圾回收
### Golang的三色标记法
golang 的垃圾回收(GC)是基于标记清扫算法，这种算法需要进行 STW（stop the world)，这个过程就会导致程序是卡顿的，频繁的 GC 会严重影响程序性能. golang 在此基础上进行了改进，通过三色标记清扫法与写屏障来减少 STW 的时间.

三色标记法的流程如下，它将对象通过白、灰、黑进行标记

* 1.所有对象最开始都是白色.
* 2.从 root 开始找到所有可达对象，标记为灰色，放入待处理队列。
* 3.历灰色对象队列，将其引用对象标记为灰色放入待处理队列，自身标记为黑色。
* 4.循环步骤3直到灰色队列为空为止，此时所有引用对象都被标记为黑色，所有不可达的对象依然为白色，白色的就是需要进行回收的对象。

三色标记法相对于普通标记清扫，减少了 STW 时间. 这主要得益于标记过程是 “on-the-fly” 的，在标记过程中是不需要 STW 的，它与程序是并发执行的，这就大大缩短了 STW 的时间.

### 写屏障
当标记和程序是并发执行的，这就会造成一个问题. 在标记过程中，有新的引用产生，可能会导致误清扫. 清扫开始前，标记为黑色的对象引用了一个新申请的对象，它肯定是白色的，而黑色对象不会被再次扫描，那么这个白色对象无法被扫描变成灰色、黑色，它就会最终被清扫，而实际它不应该被清扫. 这就需要用到屏障技术，golang 采用了写屏障，作用就是为了避免这类误清扫问题. 写屏障即在内存写操作前，维护一个约束，从而确保清扫开始前，黑色的对象不能引用白色对象.

### 触发条件
* 1> 当前内存分配达到一定比例则触发
* 2> 2 分钟没有触发过 GC 则触发 GC
* 3> 手动触发，调用 runtime.GC()

### 三色标记法，主要流程如下：
* 所有对象最开始都是白色。
* 从 root 开始找到所有可达对象，标记为灰色，放入待处理队列。
* 遍历灰色对象队列，将其引用对象标记为灰色放入待处理队列，自身标记为黑色。
* 处理完灰色对象队列，执行清扫工作。

* 1.首先从 root 开始遍历，root 包括全局指针和 goroutine 栈上的指针。
* 2.mark 有两个过程。
    - 从 root 开始遍历，标记为灰色。遍历灰色队列。
    - re-scan 全局指针和栈。因为 mark 和用户程序是并行的，所以在过程 1 的时候可能会有新的对象分配，这个时候就需要通过写屏障（write barrier）记录下来。re-scan 再完成检查一下。
* 3.Stop The World 有两个过程。
    - 第一个是 GC 将要开始的时候，这个时候主要是一些准备工作，比如 enable write barrier。
    - 第二个过程就是上面提到的 re-scan 过程。如果这个时候没有 stw，那么 mark 将无休止。

参考连接：
https://www.cnblogs.com/hezhixiong/p/9577199.html

## 3.Golang GC性能优化技巧
* (1)slice预先分配内存
* (2)新建的map也可以指定大小
    - 减少内存拷贝的开销，也可以减少rehash开销
    - map中保存值，而不是指针，使用分段map
    - 使用分段的，保存值的map的GC耗时最小。加上GODEBUG=gctrace=1分析GC轨迹
    - map保存值比保存指针的耗时少，主要是在GC的标记阶段耗时更少
* (3)string与[]byte的转换
    - string从设计上是不可变的。因此，string和[]byte的类型转化，都是产生一份新的副本
    - 如果确定转换的string/[]byte不会被修改，可以进行直接的转换，这样不会生成原有变量的副本。新的变量共享底层的数据指针
    -
    ```
    func String2Bytes(s string) []byte {
        stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
        bh := reflect.SliceHeader {
            Data: stringHeader.Data,
            Len: stringHeader.Len,
            Cap: stringHeader.Len,
        }
        return *(*[]byte)(unsafe.Pointer(&bh))
    }
    func Bytes2String(b []byte) string {
        sliceHeader
        sh := reflect.StringHeader {
            Data: sliceHeader.Data,
            Len: sliceHeader.Len,
        }
        return *(*string)(unsafe.Pointer(&sh))
    }
    ```
* (4)函数返回值使用值，不使用指针
* (5)使用struct{}优化
    - Golang中，没有集合set。如果要实现一个集合，可以使用struct{}作为值
    - struct{} 经过编译器特殊优化，指向同一个内存地址(runtime.zerobase)，不占用空间。

#### GC分析的工具
* go tool pprof
* go tool trace
* go build -gcflags="-m"
* GODEBUG="gctrace=1"

