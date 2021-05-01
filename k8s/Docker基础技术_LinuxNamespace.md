# Docker基础技术：Linux Namespace
Linux Namespace是Linux提供的一种内核级别环境隔离的方法。

很早以前有一个叫chroot的系统调用（通过修改根目录把用户jail到一个特定目录下），chroot提供了一种简单的隔离模式：chroot内部的文件系统无法访问外部的内容。

Linux Namespace在此基础上，提供了对UTS、IPC、mount、PID、network、User等的隔离机制。

|分类 |系统调用参数 |相关内核版本|
|----- |:----- |:----- |
|Mount namespaces | CLONE_NEWNS| Linux 2.4.19|
|UTS namespaces | CLONE_NEWUTS | Linux 2.6.19|
|IPC namespaces | CLONE_NEWIPC | Linux 2.6.19|
|PID namespaces | CLONE_NEWPID| Linux 2.6.24|
|Network namespaces | CLONE_NEWNET|始于Linux 2.6.24完成于Linux 2.6.29|
|User Namespaces | CLONE_NEWUSER|始于Linux 2.6.23完成于Linux 3.8|

### 主要是三个系统调用
* clone()- 实现线程的系统调用，用来创建一个新的进程，并可以通过设计上述参数达到隔离
* unshare()- 使某进程脱离某个namespace
* setns()- 把某进程加入到某个namespace

### UTS Namespace
子进程的hostname变成了container

### IPC Namespace
* IPC全程Inter-Process Communication，是Unix/Linux下进程间通信的一种方式，IPC有共享内存、信号量、消息队列等方法。

* 所以，为了隔离，我们也需要把IPC给隔离开来，这样，只有在同一个Namespace下的进程才能互相通信。
* IPC需要有一个全局的ID，既然是全局的，那么就意味着我们的Namespace需要对这个ID隔离，不能让别的Namespace的进程看到。

### PID Namespace
子进程的pid是1了

* 在传统的UNIX系统中，PID为1的进程是init，地位非常特殊。
* 作为所有进程的父进程，有很多特权（比如：屏蔽信号等），另外，还会为检查所有进程的状态，我们知道，如果某个子进程脱离了父进程（父进程没有wait它），那么init就会负责回收资源并结束这个子进程。
* 要做到进程空间的隔离，首先创建出PID为1的进程，最好就行chroot那样，把子进程的PID在容器内变成1。

**注意：**
* 但是，我们会发现，在子进程的shell里输入ps，top等命令，我们还是会看到所有进程。 说明并没有完全隔离。
* 这是因为，像ps、top这些命令会去读/proc文件系统，所以，因为/proc文件系统在父进程和子进程都是一样的，所以这些命令显示的东西都是一样的。
* 所以我们还需要对文件系统进行隔离。

### Mount Namespace
我们在启用了mount namespace并在子进程中重新mount了/proc文件系统。

ps top后就只有两个了，/proc目录下也干净了很多。

多说一下：
* 在通过CLONE_NEWNS创建mount namespace后，父进程会把自己的文件结构复制给子进程中
* 而子进程中新的namespace中的所有mount操作都只影响自身的文件系统，而不对外界产生任何影响。这样可以做到比较严格的隔离。

### User Namespace
User Namespace主要是用了CLONE_NEWUSER的参数。使用了这个参数后，内部看到的UID和GID已经与外部不同了，默认显示为65534 。

那是因为容器找不到其真正的UID所以，设置了最大的UID（其设置定义在/proc/sys/kernel/overflowuid）。

### Network Namespace




