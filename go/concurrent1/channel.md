# Channel:另辟蹊径，解决并发问题

## Channel的发展
Channel这种go语言中特有的数据结构，要追溯到CSP模型。

CSP是Communicating Sequential Process的简称，中文直译为通信顺序进程，或者叫作交换信息的循序进程，是用来描述并发系统中进行交互的一种模式。

CSP允许使用进程组件来描述系统，它们独立运行，并且只通过消息传递的方式通信。

## Channel的应用场景
不要通过共享内存的方式来通信，而是要通过Channel通信的方式分享数据。

* communicate by sharing memory 和 share memory by communicating 是两种不同的并发处理模式。
* communicate by sharing memory是传统的并发编程处理方式，就是指，共享的数据需要用锁进行保护，goroutine需要获取到锁，才能并发访问数据。
* share memory by communicating 则是类似于CSP模型的方式，通过通信的方式，一个goroutine可以把数据的“所有权”交给另外一个goroutine（虽然go中没有“所有权”的概念，但是从逻辑上说，可以把它理解为是所有权的转移）。

Channel的应用场景分为5种类型
* 数据交流：当作并发的buffer或者queue，解决生产者-消费者问题。多个goroutine可以并发当作生产者（Producer）和消费者（Consumer）。
* 数据传递：一个goroutine将数据交给另一个goroutine，相当于把数据的拥有权（引用）托付出去。
* 信号通知：一个goroutine可以将信号（closing、closed、data ready等）传递给另一个或者另一组goroutine。
* 任务编排：可以让一组goroutine按照一定的顺序并发或者串行的执行，这就是编排的功能。
    - Or-Done模式
    - 扇入模式
    - 扇出模式
    - Stream
    - Map-Reduce
* 锁：利用Channel也可以实现互斥锁的机制。

## Channel的基本用法
Channel类型分为只能接收、只能发送、既可以接收又可以发送三种类型。


## Channel的实现原理
主要介绍chan的数据结构、初始化的方法以及三个重要的操作方法，分别是send、recv和close。

**数据结构**
chan类型的数据结构如下所示，它的数据类型是runtime.hchan。
|字段|说明|
|:-----|:-----|
|qcount uint|循环队列元素的个数|
|dataqsiz uiint |循环队列的大小|
|buf unsafe.Pointer|循环队列的指针|
|elemsize uint16|chan中元素的大小|
|closed uint32|是否已close|
|elemtype *_type|chan中元素个数|
|sendx uint16|send在buf中的索引|
|recvx uint16|recv在buf中的索引|
|recvq waitq|receiver的等待队列 |
|sendq waitq|sender的等待队列 |
|lock mutex | 互斥锁，保护所有字段|

* 初始化：go在编译的时候，会根据容量的大小选择调用makechan64，还是makechan
* send：Go在编译发送数据给chan的时候，会把send语句转换成chansend1函数，chansend1函数会调用chansend
* recv：在处理从chan中接收数据时，Go会把代码转换成chanrecv1函数
* close：通过close函数，可以把chan关闭，编译器会替换成closechan方法的调用

## 使用Channel容易犯的错误
使用Channel最常见的错误是panic和goroutine泄露。
### 3种会panic的情况
* close为nil的chan；
* send已经close的chan；
* close已经close的chan；

### goroutine泄露的问题
```
func process(timeout time.Duration) bool {
    ch := make(chan bool)

    go func(){
        // 模拟处理耗时的业务
        time.Sleep((timeout + time.Second))
        ch <- true // block
        fmt.Println("exit goroutine")
    }()
    select {
    case result := <-ch:
        return result
    case <-time.After(timeout):
        return false
    }
}
```
发生超时，process函数就会返回了，这就会导致unbuffered的chan从来就没有被读取。

我们知道，unbuffered chan必须等reader和writer都准备好了才能交流，否则就会阻塞。超时导致未读，结果就是子goroutine就阻塞在第7行永远结束不了，进而导致goroutine泄露。

解决这个bug的办法很简单，就是将unbuffered chan改成容量为1的chan

**一套选择的方法**
* 共享资源的并发访问使用传统并发原语
* 复杂的任务编排和消息传递使用Channel
* 消息通知机制使用Channel，除非只想signal一个goroutine，才使用Cond
* 简单等待所有任务的完成用WaitGroup，也有Channel的推崇着用Channel，都可以
* 需要和Select语句结合，使用Channel
* 需要和超时配合时，使用Channel和Context

