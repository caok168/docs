# 容器文件系统
## 一、如何理解容器文件系统
在容器里运行df命令，可以在容器中的根目录（/）的文件系统类型是"overlay"

为了有效的减少磁盘上冗余的镜像数据，同时减少冗余的镜像数据在网络上的传输，选择一种针对于容器的文件系统是很有必要的，而这类的文件系统被称为UnionFS。

- UnionFS这类文件系统实现的主要功能是把多个目录（处于不同的分区）一起挂载（mount）在一个目录下。
- 可以把ubuntu18.04 和 app_1 两个文件夹，mount到container_1文件夹下；

## OverlayFS
UnionFS类似的有很多种实现，包括在Docker里最早使用的AUFS，还有目前我们使用的OverlayFS。

在Linux内核3.18版本中，OverlayFS代码正式合入Linux内核的主分支。在这之后，OverlayFS也就逐渐称为各个主流Linux发型版本里缺省使用的容器文件系统了。

* OverlayFS的一个mount命令牵涉到的四类目录
    * lower
    * upper
    * merged
    * work

* lower/，也就是被mount两层目录中底下的这层（lowerdir）；最底下这一层里的文件是不会被修改的，可以认为它是只读的； OverlayFS是支持多个lowerdir的。
* upper/，是被mount两层目录中上面的这层（upperdir）；在OverlayFS中，如果有文件的创建，修改，删除操作，那么都会在这一层反映出来，它是可读写的。
* merged：它是挂载点（mount point）目录，也是用户看到的目录，用户的实际文件操作在这里
* work/：只是一个存放临时文件的目录，OverlayFS中如果有文件修改，就会在中间过程中临时存放文件到这里。

容器启动后，对镜像文件中修改就会被保存在upperdir里了。
```
root@8124ddb8bc62:/# cat /proc/mounts |grep overlay
overlay / overlay rw,relatime,lowerdir=/var/lib/docker/overlay2/l/S6JDSAMUF4ASJCCKVFNBAM4GRW:/var/lib/docker/overlay2/l/KHDYLTOF54RRRKESGRESH7KRVW:/var/lib/docker/overlay2/l/FMXXVEJI6FZI5TCPLJ7BNNPZEE:/var/lib/docker/overlay2/l/W7P4KQFHS7OXG5GS5ETFM4KC4Y,upperdir=/var/lib/docker/overlay2/9adda346fd4bebbc9652f9cc98df95f6cdf18a4540805b5edf72f16ea32f1409/diff,workdir=/var/lib/docker/overlay2/9adda346fd4bebbc9652f9cc98df95f6cdf18a4540805b5edf72f16ea32f1409/work 0 0
```

## **重点总结**
为什么要有容器自己的文件系统？
* 很重要的一点是减少相同镜像文件在同一个节点上的数据冗余，可以节省磁盘空间，也可以减少镜像文件下载占用的网络资源。
* OverlayFS也是把多个目录合并挂载，被挂载的目录分为两大类：lowerdir和upperdir。
* lowerdir允许有多个目录，在被挂载后，这些目录里的文件都是不会被修改或者删除的，也就是只读的；
* upperdir只有一个，不过这个目录是可读写的，挂载点目录中的所有文件修改都会在upperdir中反映出来。
* 容器的镜像文件中各层作为OverlFS的lowerdir的目录，加上一个空的upperdir一起挂载好后，就组成了容器的文件系统。



