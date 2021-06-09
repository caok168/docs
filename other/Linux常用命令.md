# 常用命令

top
然后按 g 再输入3，进入内存模式

### 内存相关

在内存模式，可以看到各个进程内存的 %MEM、VIRT、RES、CODE、DATA、SHR、nMaj、nDRT，这些信息通过strace来跟踪top进程。

这些信息都是从 /proc/[pid]/statm 和 /proc/[pid]/stat 这个文件读取的。

```
strace -p `pidof top`

strace: Process 4115 attached
pselect6(1, [0], NULL, NULL, {1, 428569237}, {[], 8}) = 0 (Timeout)
lseek(5, 0, SEEK_SET)                   = 0
read(5, "MemTotal:        7862396 kB\nMemF"..., 8191) = 1307
lseek(4, 0, SEEK_SET)                   = 0
read(4, "188.54 669.86\n", 8191)        = 14
open("/proc", O_RDONLY|O_NONBLOCK|O_DIRECTORY|O_CLOEXEC) = 8
fstat(8, {st_mode=S_IFDIR|0555, st_size=0, ...}) = 0
getdents(8, /* 289 entries */, 32768)   = 7304
stat("/proc/1", {st_mode=S_IFDIR|0555, st_size=0, ...}) = 0
open("/proc/1/stat", O_RDONLY)          = 9
read(9, "1 (systemd) S 0 1 1 0 -1 4194560"..., 2048) = 185
close(9)                                = 0
open("/proc/1/statm", O_RDONLY)         = 9
read(9, "30020 1565 1014 352 0 4641 0\n", 2048) = 29
close(9)                                = 
```

#### vmstat
```
vmstat -w
```

```
ck@ck-ThinkPad-T430:~/docs/other$ vmstat -w
procs -----------------------memory---------------------- ---swap-- -----io---- -system-- --------cpu--------
 r  b         swpd         free         buff        cache   si   so    bi    bo   in   cs  us  sy  id  wa  st
 0  0            0      4143184       180024      2316332    0    0   134    41  247  819   5   2  93   0   0
```

#### pidstat
```
pidstat -r -p [pid]
```
```
ck@ck-ThinkPad-T430:~$ pidstat -r -p 4316
Linux 4.15.0-142-generic (ck-ThinkPad-T430) 	2021年05月30日 	_x86_64_	(4 CPU)

13时34分50秒   UID       PID  minflt/s  majflt/s     VSZ     RSS   %MEM  Command
13时34分50秒  1000      4316     22.31      0.21 4932312  137232   1.75  code
```

#### 继续查看某个进程的内存使用细节，可以使用pmap
```
pmap -x [pid]
```

```
ck@ck-ThinkPad-T430:~$ pmap -x 4316
4316:   /usr/share/code/code --no-sandbox --unity-launch /home/ck/docs
Address           Kbytes     RSS   Dirty Mode  Mapping
0000033400000000      48      48      48 rw---   [ anon ]
0000033400000000       0       0       0 rw---   [ anon ]
000003340000c000     208       0       0 -----   [ anon ]
000003340000c000       0       0       0 -----   [ anon ]
0000033400040000       4       4       4 rw---   [ anon ]
0000033400040000       0       0       0 rw---   [ anon ]
0000033400041000       4       0       0 -----   [ anon ]
```

pmap同样也是解析的/proc里的文件，具体文件是 /proc/[pid]/maps 和 /proc/[pid]/smaps,其中smaps文件相比maps的内容更详细，可以理解为是对maps的一个扩展。

除了观察进程自身的内存外，还可以观察进程分配的内存和系统指标的关联，我们就以常用的 /proc/meminfo为例。

提到的四种内存类型
* 私有匿名
* 私有文件
* 共享匿名
* 共享文件

凡是私有的内存都会体现在 /proc/meminfo中的AnonPages这一项，凡是共享的内存都会体现在Cached这一项，匿名共享的则还会体现在Shmem这一项。

总结：
* 进程直接读写的都是虚拟地址，虚拟地址最终会通过Paging（分页）来转换为物理内存的地址，Paging这个过程是由内核来完成的。
* 进程的内存类型可以从anon（匿名）与file（文件）、private（私有）与shared（共享）这四项来区分4种不同的类型，进程相关的所有内存都是这几种方式的不同组合。
* 查看进程内存时，可以先使用top来看系统中各个进程的内存使用情况，再使用pmap去观察某个进程的内存细节。

### 2.评测网络吞吐量
* Server端
    - iperf3 -s -i 1 -p 10000
* Client端
    - iperf3 -c 192.168.11.92 -b 1G -t 15 -P 2 -p 10000

### 3.评测网络性能
netperf：是一个衡量网络性能的工具，它可以提供单向吞吐量和端到端延迟的测试

测试宿主机和容器中网络的延迟情况
* 在宿主机上执行
    - netperf -H 192.168.1.1 -t TCP_RR
* 在容器中执行
    - netperf -H 192.168.1.1 -t TCP_RR



