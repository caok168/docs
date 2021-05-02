# 容器内存
## Swap:容器可以使用Swap空间吗？

Swap好处：可以应对一些瞬时突发的内存增大需求，不至于因为内存一时不够而触发OOM Killer，导致进程被杀死。

## 如何正确理解swappiness参数？
* /proc/sys/vm/swappiness （全局）
* swappiness可以决定系统将会多频繁地使用交换分区。
* 一个较高的值会使得内核更频繁地使用交换分区，而一个较低的取值，则代表着内核会尽量避免使用交换分区。
* swappiness的取值范围是 0-100,缺省值60。

在RSS里的内存，大部分都是没有对应磁盘文件的内存，比如用malloc（）申请得到的内存，这种内存也被称为**匿名内存**（Anonymous memory）。当Swap空间打开后，可以写入Swap空间的，就是这些匿名内存。

在Swap空间打开的时候，在内存紧张的时候，Linux系统怎么决定是先释放Page Cache，还是先把匿名内存释放并写入到Swap空间里呢？

这时候Swappiness的作用就起到了：

* 我们在释放内存的时候，需要平衡Page Cache的释放和匿名内存的释放，而swappiness，就是用来定义这个平衡的参数。
* swappiness的这个值的范围是0到100,它不是一个百分比，更像是一个权重。它是用来定义Page Cache内存和匿名内存的释放的一个比例。

swappness
* swappiness的取值范围在0到100,值为100的时候系统平等回收匿名内存;
* 一般缺省值为60,就是优先回收Page Cache;
* 即使swappiness为0,也不能完全禁止Swap分区的使用，就是说在内存紧张的时候，也会使用Swap来回收匿名内存。

swappiness参数除了在proc文件系统下有个全局的值外，在每个Memory Cgroup控制组里也有一个memory.swappiness，有什么不同？
* 每个Memory Cgroup控制组里的swappiness参数值为0的时候，就可以让控制组里的内存停止写入Swap。
* 这样，有了memory.swappiness这个参数后，需要使用Swap和不需要Swap的容器就可以在同一个宿主机上同时运行了，这样对于硬件的利用率也就更高了。


