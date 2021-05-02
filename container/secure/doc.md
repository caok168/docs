## 容器安全
### 一、两个方面需要注意
* 赋予容器合理的capabilities
* 在容器中以非root用户来运行程序

在容器中执行iptables -L 会提示Permission denied(you must be root)

这个时候在创建容器的时候加上privileged参数，就可以执行成功了

从Docker的代码里，如果配置上privileged的参数的话，就会获取所有的capabilities。

```
if ec.Privileged {
    p.Capabilities = caps.GetAllCapabilities()
}
```

### 二、基本概念
#### Linux capabilities
在Linux capabilities出现前，进程的权限可以简单分为两类：
* 第一类是特权用户的进程（root用户的进程）
* 第二类是非特权用户的进程（可以理解为非root用户进程）

从kernel2.2开始，Linux把特权用户所有的这些“特权”做了更详细的划分，这样被划分出来的每个单元就被称为capability。

对于任意一个进程，在做任意一个特权操作的时候，都需要有这个特权操作对应的capability。

比如：
* 运行iptables命令，对应的进程需要有CAP_NET_ADMIN这个capability。
* mount一个文件系统，那么对应的进程需要有CAP_SYS_ADMIN这个capability。

在普通Linux节点上，非root用户启动的进程缺省没有任何Linux capability，而root用户启动的进程缺省包含了所有的Linux capability。

* 使用capsh这个工具
* 可以通过/proc文件系统找到对应的status cat /proc/{processId}/status

### 解决问题
对于Privileged的容器也就是允许容器中的进程可以执行所有的特权操作。

考虑安全的因素，容器缺省启动的时候，哪怕是容器中的root用户的进程，系统也只允许了15个capability。

查看方式
* 查看runC spec文档中的security部分
* 查看容器init进程status里的Cap参数，看一下容器中缺省的capability
    * cat /proc/1/status |grep Cap

在容器平台上基本不允许把容器直接设置为"privileged"的，我们需要根据容器中进程需要的最少特权来赋予capabilities。

比如：
* docker run --name iptables --cap-add NET_ADMIN -it registry/iptables:v1 bash
    * root@cfedf124dcf1# iptables -L

### 重点小结
* 其实Linux capabilities就是把Linux root用户原来所有的特权做了细化，可以更加细粒度地给进程赋予不同权限。
* 每一个特权操作对应一个capability,对于一个capability，有的对应一个特权操作，有的可以对应很多个特权操作
* 每个Linux进程有5个capabilities集合参数，其中Effective集合里的capabilities决定了当前进程可以做哪些特权操作，而其他集合参数会和应用程序文件的capabilities集合参数一起来决定新启动程序的capabilities集合参数
* 对于容器的root用户，缺省只赋予了15个capabilities。如果发现容器中进程的权限不够，就需要分析它需要的最小capabilities集合，而不是直接赋予容器"privileged"。
* 因为"privileged"包含了所有的Linux capabilities，这样"privileged"就可以轻易获取宿主机上的所有资源，这会对宿主机的安全产生威胁。

### User Namespace
容器云平台Kubernetes上目前还不支持User Namespace

当然在云平台上，比如在Kubernetes里，我们可以限制容器去挂载宿主机的目录的

* 在没有User Namespace的情况下，容器中和宿主机上的用户都是root，也就是说，容器中用户的uid/gid和宿主机上的完全一样
* 虽然容器里root用户的capabilities被限制，但在容器中，对于被挂载上来的/etc目录下的文件，比如shadow文件，以这个root用户的权限还是可以做修改的

### 解决办法
* 方法一：Run as non-root user（给容器指定一个普通用户）
    * 在docker启动容器的时候加上"-u"参数，在参数中指定uid/gid
        * docker run -it --name root_example -u 6667:6667 -v /etc:/mnt centos bash
    * 在创建容器镜像的时候，用Dockerfile为容器镜像里创建一个用户
        ```
        # cat Dockerfile
        FROM centos
        RUN adduser -u 6667 nonroot
        USER nonroot
        ```
    * 存在的问题
        * 由于用户uid是整个节点中共享的，那么在容器中定义的uid，也就是宿主机上的uid，容易引起uid的冲突
        * 在一台Linux系统上，每个用户下的资源是有限制的，比如打开文件数目(open files)、最大进程数目(max user processs)等等。一旦有很多个容器共享一个uid，这些容器就很快消耗掉uid下的资源，很容易导致这些容器都不能再正常工作。
        * 解决问题，必须要有一个云平台级别的uid管理和分配，但选择这个方法也要付出代价。用户在定义自己容器中uid的时候，他们就需要有额外的操作，而且平台也需要新开发对uid平台级别的管理模块，完成这些事情要的工作量也不少。
* 方法二：User Namespace（用户隔离技术的支持）
    * User Namespace概念：
        * User Namespace隔离了一台Linux节点上的User ID(uid)和Group ID（gid），它给Namespace中uid/gid的值与宿主机上的uid/gid值建立了一个映射关系。
        * 经过User Namespace的隔离，我们在Namespace中看到的进程uid/gid，就和宿主机Namespace中看到的uid和gid不一样了。
            * 比如namespace_1里的uid值是0到999，但其实它在宿主机上对应的uid值是1000到1999。
            * User Namespace是可以嵌套的，namespace_2里可以再建立一个namespace_3，这个嵌套的特性是其他Namespace没有的。
    * Podman例子
        * podman run -ti -v /etc:/mnt --uidmap 0:2000:1000 centos bash
        * 0:2000:1000 -> 第一个0指在新的Namespace里uid从0开始，中间的那个2000指的是Host Namespace里被映射的uid从2000开始，最后一个1000是指总共需要连续映射1000个uid。
    * User Namespace有两个好处
        * 第一，它把容器中root用户（uid 0）映射成宿主机上的普通用户。
        * 第二，对于用户在容器中自己定义普通用户uid的情况，我们只要为每个容器在节点上分配一个uid范围，就不会出现在宿主机上uid冲突的问题了。
    * 目前Kubernetes还不支持User Namespace，相关工作进展可以看以下社区这个PR, https://github.com/kubernetes/enhancements/pull 2101

* 方法三：rootless container（以非root用户启动和管理容器）
    * 这里的rootless container中的"rootless"不仅仅指容器中以非root用户来运行进程，还指以非root用户来创建容器，管理容器。也就是说，启动容器的时候，Docker或者podman是以非root用户来执行的。
    * 目前Docker和podman都支持了rootless container，Kubernetes对rootless container支持的工作也在进行中。https://github.com/kubernetes/enhancements/issues/2033
    


