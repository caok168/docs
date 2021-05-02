# 容器CPU
## 容器CPU：怎么限制容器CPU使用？

CPU Usage一览表

 | 类型 | 具体含义 |
 | ----- | ----- |
 | us | User，用户态CPU时间，不包括低优先级进程的用户态时间(nice值 1-19) |
 | sys | System，内核态CPU时间 |
 | ni  | Nice，nice值 1 - 19的进程用户态CPU时间 |
 | id  | Idle，系统空闲CPU时间 |
 | wa  | Iowait，系统等待I/O的CPU时间，这个时间不计入进程CPU时间 |
 | hi  | Hardware irq,处理硬中断的时间，这个时间不计入进程CPU时间 |
 | si  | Softirq，处理软中断的时间，这个时间不计入进程CPU时间 |
 | st  | Steal，表示同一个宿主机上的其他虚拟机抢走的CPU时间 |

普通调度的算法在Linux中目前是CFS（即完全公平调度器）。
直接来看CPU Cgroup和CFS相关的参数，一共有三个。
* cpu.cfs_period_us：它是CFS算法的一个调度周期，一般它的值是100000，以microseconds为单位，也就100ms。
* cpu.cfs_quota_us：表示CFS算法中，在一个调度周期里这个控制组被允许的运行时间，比如这个值为50000时，就是50ms。
    * cpu.cfs_quota_us / cpu.cfs_period_us 比如 50ms/100ms = 0.5 个CPU。
* cpu.share：这个值是CPU Cgroup对于控制组之间的CPU分配比例，它的缺省值是1024。

对CPU Cgroup的参数做一个梳理
* cpu.cfs_quota_us和cpu.cfs_period_us这两个值决定了每个控制组中所有进程的可使用CPU资源的最大值
* cpu.shares这个值决定了CPU Cgroup子系统下控制组可用CPU的相对比例，不过只有当系统上CPU完全被占满的时候，这个比例才会在控制组间起作用。

