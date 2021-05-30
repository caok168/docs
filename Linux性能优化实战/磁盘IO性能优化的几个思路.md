# 套路篇：磁盘I/O性能优化的几个思路

### I/O基准测试

fio 是最常用的文件系统和磁盘I/O性能基准测试工具。

```
sudo apt install -y fio
```

fio支持I/O的重放。借助blktrace，再配合上fio，就可以实现对应用程序I/O模式的基准测试。

需要先用blktrace，记录磁盘设备的I/O访问情况；然后使用fio，重放blktrace的记录。

```
#使用blktrace跟踪磁盘I/O，注意指定应用程序正在操作的磁盘
blktrace /dev/sdb

# 查看blktrace记录的结果
ls
sdb.blktrace.0 sdb.blktrace.1

# 将结果转化为二进制文件
blktrace sdb -d sdb.bin

# 使用fio重放日志
fio --name=replay --filename=/dev/sdb --direct=1 --read_iolog=sdb.bin
```


