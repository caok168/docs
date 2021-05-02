# 容器存储
## 容器中的内存与IO
### 现象
当使用Buffered I/O的应用程序从虚拟机迁移到容器，我们会发现多了Memory Cgroup的限制之后，write()写相同大小的数据块花费的时间，延时波动会比较大

### 时间波动是因为Dirty Pages的影响吗？
对于Buffer I/O，用户的数据是先写入到Page Cache里的。而这些写入了数据的内存页面，在它们没有被写入到磁盘文件之前，就被叫作dirty pages。

Linux内核会有专门的内核线程（每个磁盘设备对应的kworker/flush线程）把dirty pages写入到磁盘中。

比如A等于dirty pages的内存/节点可用内存*100%

**dirty pages相关的参数**
* 在/proc/sys/vm里和dirty page相关的内核参数
* dirty_background_ratio：这个参数里的数值是一个百分比值，缺省是10%。如果比值A大于这个参数的值的话，内核flush线程就会把dirty pages刷到磁盘里。
* dirty_background_bytes：和dirty_background_ratio作用相同。区别只是表示具体的字节数。
* dirty_ratio：缺省值为20%.如果比值A大于默认值，这时候正在执行Buffered I/O写文件的进程就会被阻塞住，直到它写的数据页面都写到磁盘为止。
* dirty_bytes：与dirty_ratio作用相同。
* dirty_writeback_centisecs：这个参数的值是个时间值，以百分之一秒为单位，缺省值为500,5秒钟。表示每5秒中会唤醒内核的flush线程来处理dirty pages。
* dirty_expire_centisecs：这个参数的值是个时间值，以百分之一秒为单位，缺省值是3000，也就是30秒钟。它定义了dirty page在内存中存放的最长时间，如果一个dirty page超过这里定义的时间，那么内核的flush线程也会把这个页面写入磁盘。

**注意**
* dirty_background_ratio与dirty_background_bytes只能使用一个
* dirty_ratio与dirty_bytes只能使用一个

查看dirty pages的实时数目：

```
watch -n 1 "cat /proc/vmstat |grep dirty"
```

****



