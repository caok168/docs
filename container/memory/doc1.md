# 容器内存
## 如何理解OOM Killer？
OOM是Out of Memory的缩写，顾名思义就是内存不足的意思，而Killer在这里指需要杀死某个进程。

OOM Killer就是在Linux系统里如果内存不足时，就需要杀死一个正在运行的进程来释放一些内存。

在发生OOM的时候，Linux到底是根据什么标准来选择被杀的进程呢？
* 在Linux内核里有一个oom_badness()函数，就是它定义了选择进程的标准。
* 判断标准也很简单，函数中涉及两个条件：
    * 1.进程已经使用的物理内存页面数；
    * 2.每个进程的OOM校准值oom_score_adj。在/proc文件系统中，每个进程都有一个/proc//oom_score_adj的接口文件。我们可以在这个文件中输入-1000到1000之间的任意一个数值，调整进程被OOM Kill的几率。
* 函数oom_badness()的最终计算方法：
    * 用系统总的可用页面数，去乘以OOM校准值oom_score_adj，再加上进程已经使用的物理页面数，计算出来的数值越大，那么这个进程被OOM Kill的几率也就越大。

## 如何理解Memory Cgroup
### 基本概念
* 1. Memory Cgroup中每一个控制组可以为一组进程限制内存使用量，一旦所有进程使用内存的总量达到限制值，缺省情况下，就会触发OOM Killer。这样一来，控制组里的“某个进程”就会被杀死。
* 2. 杀死某个进程的选择标准是，控制组中总的可用页面乘以进程的oom_score_adj，加上进程已经使用的物理内存页面，所得值最大的进程，就会被系统选中杀死。

### 三个基本参数
* memory.limit_in_bytes：直接限制控制组里所有进程可使用内存的最大值
* memory.oom_control：当控制组中的进程内存使用达到上限时，这个参数能够决定会不会触发OOM Killer，默认会触发
* memory.usage_in_bytes：只读参数，里面的数值是当前控制组里所有进程实际使用的内存总和。数值越接近参数1，OOM的风险越高

### 解决问题
怎么才能快速确定容器发生了OOM呢？这个可以通过查看内核日志及时地发现。

使用journal -k命令，或者直接查看日志文件/var/log/message。
* 第一部分就是容器里每一个进程使用的内存页面数量。
* 第二部分oom-kill，这行列出了发生OOM的Memory Cgroup的控制组，可以从控制组的信息中知道OOM是在哪个容器发生的。
* 第三部分显示了最终被OOM Killer杀死的进程。

**分析**
* 第一种情况是这个进程本身的确需要很大的内存，这说明我们给memory.limit_in_bytes里的内存上限值设置小了，那么就需要增大内存的上限值。
* 第二种情况是进程的代码中有Bug，会导致内存泄露，进程内存使用到达了Memory Cgroup中的上限。如果是这种情况，就需要我们具体去解决代码里的问题了。



