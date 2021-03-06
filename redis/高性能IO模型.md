# 高性能IO模型：为什么单线程Redis能那么快？

Redis是单线程，主要是指Redis的网络IO和键值对读写是由一个线程来完成的，这也是Redis对外提供键值存储服务的主要流程。

但Redis的其他功能，比如持久化、异步删除、集群数据同步等，其实是由额外的线程执行的。

### Redis为什么使用单线程
* 多线程的开销
  - 当有多个线程要修改这个共享资源时，为了保证共享资源的正确性，就需要有额外的机制进行保证，而这个额外的机制，就会带来额外的开销。
  - 多线程编程模式面临的共享资源的并发访问控制问题。
  - 采用多线程开发一般会引入同步原语来保护共享资源的并发访问，这也会降低代码的易调试性和可维护性。


### 单线程Redis为什么那么快？
* 一方面，Redis的大部分操作在内存上完成，再加上采用了高效的数据结构
* 另一方面，就是Redis采用了多路复用机制，使其在网络IO操作中能并发处理大量的客户端请求，实现高吞吐率。



### 问题：
在Redis基本IO模型中，你觉得还有哪些潜在的性能瓶颈吗？

Redis单线程处理IO请求性能瓶颈主要包括2个方面：
* 任意一个请求在server中一旦发生耗时，都会影响整个server的性能，也就是说后面的请求都要等前面的耗时请求处理完成，自己才能被处理到。耗时的操作包括以下几种：
  - 操作bigkey：写入一个bigkey在分配内存时需要消耗更多的时间，同样，删除bigkey释放内存同样会产生耗时;
  - 使用复杂度过高的命令：
  - 大量key集中过期：Redis的过期机制也是在主线程中执行的，大量key集中过期会导致一个请求时，耗时都在删除过期key，耗时变长;
  - 淘汰策略：淘汰策略也是在主线程执行的，当内存超过Redis内存上限后，每次写入都需要淘汰一些key，也会造成耗时变长;
  - AOF刷盘开启always机制：每次写入都需要把这个操作刷到磁盘，写磁盘的速度远比写内存慢，会托慢Redis的性能;
  - 主从全量同步生成RDB：虽然采用fork子进程生成数据快照，但fork这一瞬间也是会阻塞整个线程的，实例越大，阻塞时间越久;
* 并发量非常大时，单线程读写客户端IO数据存在性能瓶颈，虽然采用IO多路复用机制，但是读写客户端数据依旧是同步IO，只能单线程依次读取客户端的数据，无法利用到CPU多核。

针对问题1,一方面需要业务人员去规避，一方面Redis在4.0推出了lazy-free机制，把bigkey释放内存的耗时操作放在了异步线程中执行，降低对主线程的影响。

针对问题2,Redis在6.0推出了多线程，可以在高并发场景下利用CPU多核多线程读写客户端数据，进一步提升server性能，当然，只是针对客户端的读写是并行的，每个命令的真正操作依旧是单线程的。

