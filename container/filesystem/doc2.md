# 容器存储
## 容器磁盘限速
### 场景
如果多个容器同时读写节点上的同一块磁盘，那么它们的磁盘读写相互之间影响吗？如果容器之间读写磁盘相互影响，我们有什么办法解决吗？

**场景再现**
Dockerfile
```
FROM centos:8.1.1911

RUN yum install -y fio
```
脚本
```
#!/bin/bash

mkdir -p /tmp/test1
docker stop fio_test1;docker rm fio_test1

docker run --name fio_test1 --volume /tmp/test1:/tmp  registry/fio:v1 fio -direct=1 -rw=write -ioengine=libaio -bs=4k -size=1G -numjobs=1  -name=/tmp/fio_test1.log
```

脚本2
```
#!/bin/bash

mkdir -p /tmp/test1
mkdir -p /tmp/test2

docker stop fio_test1;docker rm fio_test1
docker stop fio_test2;docker rm fio_test2

docker run --name fio_test1 --volume /tmp/test1:/tmp  registry/fio:v1 fio -direct=1 -rw=write -ioengine=libaio -bs=4k -size=1G -numjobs=1  -name=/tmp/fio_test1.log &
docker run --name fio_test2 --volume /tmp/test2:/tmp  registry/fio:v1 fio -direct=1 -rw=write -ioengine=libaio -bs=4k -size=1G -numjobs=1  -name=/tmp/fio_test2.log &
```
通过例子可以观察到，多个容器同时写一块磁盘的时候，它的性能会收到干扰。

在Cgroup v1中有blkio子系统，它可以来限制磁盘的I/O。

## 知识详解
### Blkio Cgroup
衡量磁盘性能的两个常见指标：
* IOPS：Input/Output Operations Per Second的简称，也就是每秒钟磁盘读写的次数，这个数值越大，当然也就表示性能越好
* 吞吐量：Throughput，是指每秒钟磁盘中数据的读取量，一般以MB/s为单位。这个读取量可以叫作吞吐量，有时候也被叫作带宽。
* 在IOPS固定的情况下，如果读写的每一个数据块越大，那么吞吐量也越大，它们的关系大概是这样的：吞吐量=数据块大小*IOPS 。

blkio Cgroup也是Cgroups里的一个子系统。在Cgroups v1里，blkio Cgroup的虚拟文件系统挂载点一般在 "/sys/fs/cgroup/blkio/"。

在blkio Cgroup中，有四个最主要的参数，它们可以用来限制磁盘I/O性能：
* blkio.throttle.read_iops_device [磁盘读取IOPS限制]
* blkio.throttle.read_bps_device   【磁盘读取吞吐量限制】
* blkio.throttle.write_iops_device [磁盘写入IOPS限制]
* blkio.throttle.write_bps_device 【磁盘写入吞吐量限制】

```
echo "253:0" 10485760 > $CGROUP_CONTAINER_PATH1/blkio.throttle.read_bps_device
```

这时候，两个容器里执行的结构都是10MB/s了。

### 注意
如果将fio命令里的 "-direct=1"给去掉，这个时候即使我们设置了blkio Cgroup，也根本不能限制磁盘的吞吐量了。

### Direct I/O和Buffered I/O
Linux系统有两种文件I/O模式：
* Direct I/O
* Buffered I/O

**Direct I/O**

用户进程如果要写磁盘文件，就会通过Linux内核的文件系统层(filesystem)->块设备层(block layer)->磁盘驱动->磁盘硬件，这样一路下去写入磁盘

**Buffered I/O**

* 用户进程只是把文件数据写到内存中(Page Cache)就返回了，而Linux内核自己有线程会把内存中的数据再写入到磁盘中。

* 在linux里，由于考虑到性能问题，绝大多数的应用都会使用Buffered I/O模式。

v1的CPU Cgroup，memory Cgroup和blkio Cgroup，那么Cgroup v1的一个整体结构，它的每一个子系统都是独立的，资源的限制只能在子系统中发生。

比如进程pid_y，可以分别属于memory Cgroup和blkio Cgroup。但是在blkio Cgroup对进程pid_y做磁盘I/O限制的时候，bklio子系统是不会关心pid_y用了哪些内存，哪些内存是不是属于Page Cache，而这些Page Cache的页面在刷入磁盘的时候，产生的I/O也不会计算到进程pid_y上面。

就是这个原因，导致了blkio在Cgroup v1里不能限制Buffered I/O。

### Cgroup V2
Cgroup v2相比Cgroup v1做的最大的变动就是一个进程属于一个控制组，而每个控制组里可以定义自己需要的多个子系统。

比如：进程pid_y属于控制组group2，而在group2里同时打开了io和memory子系统。

在Cgroup v2里，尝试以下设置 blkio Cgroup + Memory Cgroup
* 第一步：在Linux系统中打开Cgroup v2的功能。
    * 配置一个kernel参数"cgroup_no_v1=blkio,memory"，这表示把Cgroup v1的blkio和Memory两个子系统给禁止，这样Cgroup v2的io和Memory这两个子系统就打开了。
    * 可以把这个参数配置到grub中，然后重启机器
* 对Cgroup v2 io的限速配置
    ```

    ```




