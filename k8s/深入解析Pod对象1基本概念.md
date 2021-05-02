# 深入解析Pod对象（一）：基本概念

Pod扮演的是传统部署环境里“虚拟机”的角色。

凡是调度、网络、存储，以及安全相关的属性，基本上是Pod级别的。

### Pod中几个重要字段的含义和用法
* NodeSelector：是一个供用户将Pod与Node进行绑定的字段，用法如下所示：
    ```
    apiVersion: v1
    kind: Pod
    ...
    spec:
      nodeSelector:
        disktype: ssd
    ```
* NodeName:一旦Pod的这个字段被赋值，Kubernetes项目就会被认为这个Pod已经经过了调度，调度的结果就是赋值的节点名字。 所以这个字段一般由调度器负责设置，但用户也可以设置它来“骗过”调度器，当然这个做法一般是在测试或者调试的时候才会用到。
* HostAliases：定义了Pod的hosts文件（比如/etc/hosts）里的内容，用法如下：
    ```
    apiVersion: v1
    kind: Pod
    ...
    spec:
      hostAliases:
      - ip: "10.1.2.3"
        hostname:
        - "foo.remote"
        - "bar.remote"
    ...
    ```

除了上述跟“机器”相关的配置外，你可能也会发现，凡是跟容器的Linux Namespace相关的属性，也一定是Pod级别的。
* Pod的设计，就是要让它里面的容器尽可能多地共享Linux Namespace,仅保留必要的隔离和限制能力。
* 这样，Pod模拟出的效果，就跟虚拟机里程序间的关系非常类似了。

举个例子，在下面这个Pod的YAML文件中，我定义了shareProcessNamespace=true:
```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  shareProcessNamespace: true
  containers:
  - name: nginx
    image: nginx
  - name: shell
    image: busybox
    stdin: true
    tty: true
```
这就意味着这个Pod里的容器要共享PID Namespace。

凡是Pod中的容器要共享宿主机的Namespace，也一定是Pod级别的定义，比如：
```
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  hostNetwork: true
  hostIPC: true
  hostPID: true
  containers:
  - name: nginx
    image: nginx
  - name: shell
    image: busybox
    stdin: true
    tty: true
```

### Container的主要字段
* ImagePullPolicy字段
    * 默认是Always
    * Never
    * IfNotPresent
* Lifecycle字段
    * 定义的是Container Lifecycle Hooks。
    * 作用：是在容器状态发生变化时触发一系列"钩子"

### Pod对象在Kubernetes中的生命周期
Pod生命周期的变化，主要体现在Pod API对象的Status部分，这是它除了Metadata和Spec之外的第三个重要字段。

其中，pod.status.phase，就是Pod的当前状态，它有如下几种可能的情况：
* 1.Pending。这个状态意味着，Pod的YAML文件已经提交给了Kubernetes，API对象已经被创建并保存在Etcd当中。但是，这个Pod里有些容器因为某种原因而不能被顺利创建。比如，调度不成功。
* 2.Running。
* 3.Succeeded。
* 4.Failed。
* 5.Unknown。这是一个异常状态，意味着Pod的状态不能持续地被kubelet汇报给kube-apiserver，这很有可能是主从节点（Master和Kubulet）间的通信出现了问题。

更进一步，Pod对象的Status字段，还可以再细分出一组Conditions。这些细分状态的值包括：PodScheduled、Ready、Initialized，以及Unschedulable。它们主要用于描述造成当前Status的具体原因是什么。

