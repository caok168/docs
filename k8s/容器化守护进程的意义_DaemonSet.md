# 容器化守护进程的意义：DaemonSet

DaemonSet的主要作用，是让你在Kubernetes集群里，运行一个Daemon Pod。
所以这个Pod有如下三个特征：
* 这个Pod运行在Kubernetes集群里的每一个节点（Node）上；
* 每个节点上只有一个这样的Pod实例；
* 当有新的节点加入Kubernetes集群后，该Pod会自动地在新节点上被创建出来；而当旧节点被删除后，它上面的Pod也相应地会被回收掉。

列举几个例子：
* 各种网络插件的Agent组件，都必须运行在每一个节点上，用来处理这个节点上的容器网络；
* 各种存储插件的Agent组件，也必须运行在每一个节点上，用来在这个节点上挂载远程存储目录，操作容器的Volume目录；
* 各种监控组件和日志组件，也必须运行在每一个节点上，负责这个节点上的监控信息和日志搜索。

更重要的是，跟其他编排对象不一样，DaemonSet开始运行的时机，很多时候比整个Kubernetes集群出现的时机都要早。

DaemonSet其实是一个非常简单的控制器。
* 在它的控制循环中，只需要遍历所有节点，然后根据节点上是否有被管理Pod的情况，来决定是否要创建或删除一个Pod。
* 在创建每个Pod的时候，DaemonSet会自动给这个Pod加上一个nodeAffinity，从而保证这个Pod只会在指定节点上启动。
* 同时，它还会自动给这个Pod加上一个Toleration，从而忽略节点的unschedulable”污点“。
* 当然，也可以在Pod模板里加上更多种类的Toleration，从而利用DaemonSet达到自己的目的。

Deployment管理这些版本，靠的是”一个版本对应一个ReplicaSet对象“。

DaemonSet控制器操作的直接就是Pod，不可能有ReplicaSet这样的对象参与其中。那么，它的这些版本又是如何维护的呢？

所谓，一切皆对象！

**ControllerRevision**
专门用来记录某种Controller对象的版本。

* DaemonSet使用ControllerRevision，来保存和管理自己对应的”版本“。
* StatefulSet也是直接控制Pod对象的。
* 在Kubernetes项目里，ControllerRevision其实是一个通用的版本管理对象。Kubernetes项目就巧妙地避免了每种控制器都要维护一套冗余的代码和逻辑的问题。



