# 套路篇：如何迅速分析出系统I/O的瓶颈在哪里？

* 通过df工具，既可以查看文件系统数据的空间容量，也可以查看索引节点的容量
* 文件系统缓存，通过/proc/meminfo、/proc/slabinfo以及slabtop等各种来源，观察页缓存、目录项缓存、索引节点缓存以及具体文件系统的缓存情况。
* iostat:可以得到磁盘I/O使用率、吞吐量、响应时间以及IOPS等性能指标
* pidstat:可以观察到进程的I/O吐吞量以及块设备I/O的延迟等

### 套路
* 通过top查看系统的CPU使用情况，发现iowait比较高
* 通过iostat发现磁盘的I/O使用率瓶颈
* 通过pidstat 找出了大量I/O的进程
* 最后通过strace和lsof，找出了问题进程正在读写的文件，并最终锁定性能问题的来源

### 块设备I/O事件跟踪

blkstrace：示例：blkstrace -d /dev/sda -o 



