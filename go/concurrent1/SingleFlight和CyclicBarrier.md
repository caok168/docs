# SingleFlight和CyclicBarrier：请求合并和循环栅栏该怎么用？

* SingleFlight：将并发请求合并成一个请求，以减少对下层服务的压力；
* CyclicBarrier：是一个可重用的栅栏并发原语，用来控制一组请求同时执行的数据结构。

SingleFlight和sync.Once有什么区别？
* sync.Once不是只在并发的时候保证只有一个goroutine执行函数f，而是会保证永远只执行一次，而SingleFlight是每次调用都重新执行，并且在多个请求同时调用的时候只有一个执行。
* 它们两个面对的场景是不同的，sync.Once主要是用在单次初始化场景中，而SingleFlight主要用在合并并发请求的场景中，尤其是缓存场景。

### 实现原理
SingleFlight使用互斥锁Mutex和Map来实现。Mutex提供并发时的读写保护，Map用来保存同一个key的正在处理的请求。

