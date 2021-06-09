# 案例篇：Redis响应严重延迟，如何解决？

响应慢，这种情况下，我们先会怀疑，是不是系统资源出现了瓶颈。所以，先观察CPU、内存和磁盘I/O等情况肯定不会错。

### 一、执行top命令，分析系统的CPU使用情况
* 发现CPU0的iowait比较高

### 二、执行iostat，查看有没有I/O性能问题
* iostat -d -x 1
* 发现磁盘sda每秒写数据(wkB/s)为2.5MB，I/O使用率是０，虽然有些I/O操作，但并没有导致磁盘的I/O瓶颈

###　三、使用pidstat，查看哪些进程进行I/O操作
* pidstat -d 1
* 发现redis-server在向磁盘写数据

### 四、执行strace命令，并且指定redis-server的进程号9085
* strace -f -T -tt -p 9085
  - -f表示跟踪子进程和子线程
  - -T表示显示系统调用的时长
  - －tt表示显示跟踪时间
  - 
  ```
  # -f表示跟踪子进程和子线程，-T表示显示系统调用的时长，-tt表示显示跟踪时间
  $ strace -f -T -tt -p 9085
  [pid 9085] 14:20:16.826131 epoll_pwait(5, [{EPOLLIN, {u32=8, u64=8}}], 10128, 65, NULL, 8) = 1 <0.000055>
  [pid 9085] 14:20:16.826301 read(8, "*2\r\n$3\r\nGET\r\n$41\r\nuuid:5b2e76cc-"..., 16384) = 61 <0.000071>
  [pid 9085] 14:20:16.826477 read(3, 0x7fff366a5747, 1) = -1 EAGAIN (Resource temporarily unavailable) <0.000063>
  [pid 9085] 14:20:16.826645 write(8, "$3\r\nbad\r\n", 9) = 9 <0.000173>
  [pid 9085] 14:20:16.826907 epoll_pwait(5, [{EPOLLIN, {u32=8, u64=8}}], 10128, 65, NULL, 8) = 1 <0.000032>
  [pid 9085] 14:20:16.827030 read(8, "*2\r\n$3\r\nGET\r\n$41\r\nuuid:55862ada-"..., 16384) = 61 <0.000044>
  [pid 9085] 14:20:16.827149 read(3, 0x7fff366a5747, 1) = -1 EAGAIN (Resource temporarily unavailable) <0.000043>
  [pid 9085] 14:20:16.827285 write(8, "$3\r\nbad\r\n", 9) = 9 <0.000141>
  [pid 9085] 14:20:16.827514 epoll_pwait(5, [{EPOLLIN, {u32=8, u64=8}}], 10128, 64, NULL, 8) = 1 <0.000049>
  [pid 9085] 14:20:16.827641 read(8, "*2
  ```
* 从系统调用来看，epoll_pwait、read、write、fdatasync这些系统调用都比较频繁，刚才观察到的写磁盘，应该就是write或者fdatasync导致的了。

**也可以用strace，观察这个系统调用的执行情况。比如通过-e选项指定fdatasync后**

* strace -f -p 9085 -T -tt -e fdatasync
```
$ strace -f -p 9085 -T -tt -e fdatasync
strace: Process 9085 attached with 4 threads
[pid  9085] 14:22:52.013547 fdatasync(7) = 0 <0.007112>
[pid  9085] 14:22:52.022467 fdatasync(7) = 0 <0.008572>
[pid  9085] 14:22:52.032223 fdatasync(7) = 0 <0.006769>
...
[pid  9085] 14:22:52.139629 fdatasync(7) = 0 <0.008183>
```

### 五、运行lsof命令，找出系统调用的操作对象
* lsof -p 9085
```
$ lsof -p 9085
redis-ser 9085 systemd-network 3r FIFO 0,12 0t0 15447970 pipe
redis-ser 9085 systemd-network 4w FIFO 0,12 0t0 15447970 pipe
redis-ser 9085 systemd-network 5u a_inode 0,13 0 10179 [eventpoll]
redis-ser 9085 systemd-network 6u sock 0,9 0t0 15447972 protocol: TCP
redis-ser 9085 systemd-network 7w REG 8,1 8830146 2838532 /data/appendonly.aof
redis-ser 9085 systemd-network 8u sock 0,9 0t0 15448709 protocol: TCP
```


### 其他
* lsof -i
  - -i表示显示网络套接字信息