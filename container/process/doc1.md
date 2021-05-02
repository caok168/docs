### 进程
* Linux init进程，它最基本的功能都是创建出Linux系统中其他所有的进程，并且管理这些进程。
    ```
    ls -l /sbin/init
    lrwxrwxrwx 1 root root 20 Feb 5 01:07 /sbin/init -> /lib/systemd
    ```
### 信号
* 信号（Signal）其实就是Linux进程收到的一个通知。
* 对于每一个信号，进程对它的处理都有下面三个选择。
    * 忽略（Ignore）:就是对这个信号不做任何处理，但是有两个信号例外，对于SIGKILL和SIGSTOP这个两个信号，进程是不能忽略的。这是因为它们的主要作用是为Linux kernel和超级用户提供删除任意进程的特权。
    * 捕获（Catch）：这个是指用户进程可以注册自己针对这个信号的handler。
        * 对于捕获，SIGKILL和SIGSTOP这两个信号也同样例外，这两个信号不能有用户自己的处理代码，只能执行系统的缺省行为。
    * 缺省行为（Default）：Linux为每个信号都定义了一个缺省的行为，你可以在Linux系统中运行man 7 signal来查看每个信号的缺省行为。

### 总结
#### 两个基本概念
* Linux 1号进程：它是第一个用户态的进程。它直接或间接创建了Namespace中的其他进程。
* Linux信号：Linux有31个基本信号，进程在处理大部分信号时有三个选择：忽略、捕获和缺省行为。其中两个特权信号SIGKILL和SIGSTOP不能被忽略或者捕获。
#### 容器里1号进程对信号处理的两个要点：
* 在容器中，1号进程永远不会响应SIGKILL和SIGSTOP这两个特权信号；
* 对于其他的信号，如果用户自己注册了handler，1号进程可以响应。

### 僵尸进程
#### Linux的进程状态
无论进程还是线程，在Linux内核里其实都是用task_struct{}这个结构来表示的。它其实就是任务（task），也就是Linux里基本的调度单位。

在进程“活着”的时候只有两个状态：运行态（TASK_RUNNING）和睡眠态（TASK_INTERRUPTIBLE，TASK_UNINTERRUPTIBLE）。

运行态和睡眠态分别代表什么意思？
* 运行态：，无论进程是正在运行中（也就是获得了CPU资源），还是进程在run queue队列里随时可以运行，都处于这个状态。
* 睡眠态：进程需要等待某个资源而进入的状态，要等待的资源可以是一个信号量（Semaphore），或者是磁盘I/O，这个状态的进程会被放入到wait queue队列里。
    * 一个可以被打断的（TASK_INTERRUPTIBLE）,S stat。
    * 一个不可被打断的（TASK_UNINTERRUPTIBLE），D stat。

进程在调用do_exit()退出的时候，还有两个状态。
* EXIT_DEAD：也就是进程在真正结束退出的那一瞬间的状态；
* EXIT_ZOMBIE：这个进程在EXIT_DEAD前的一个状态，僵尸进程，就是处于这个状态中。

#### 限制容器中进程数目
在linux系统中如何限制进程数目
* 一台Linux机器上的进程总数目是有限制的。
* 这个最大值我们在 /proc/sys/kernel/pid_max这个参数中看到。
* Linux内核在初始化系统的时候，会根据机器CPU的数目来设置pid_max的值。
    * 如果机器中CPU数目小于32，那么pid_max就会被设置为32768（32K）；
    * 如果机器中CPU数目大于32，那么pid_max就被设置为N*1204（N就是CPU数目）。

对于每个容器来说，我们都需要限制它的最大进程数目，而这个功能由pids Cgroup这个子系统来完成。

这个功能的实现方法：
* pids Cgroup通过Cgroup文件系统的方式向用户提供操作接口，一般它的Cgroup文件系统挂载点在 /sys/fs/cgroup/pids。
* 在一个容器建立之后，创建容器的服务会在/sys/fs/cgroup/pids下建立一个子目录，就是一个控制组，控制组里最关键的一个文件就是pids.max。 我们可以向这个文件写入数值，而这个值就是这个容器中允许的最大进程数目。

父进程在创建完子进程之后就不管了，所以需要父进程调用wait()或者waitpid()系统调用来避免僵尸进程产生。




