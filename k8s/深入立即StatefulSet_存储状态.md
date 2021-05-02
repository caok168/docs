# 深入理解StatefulSet（二）：存储状态

解读StatefulSet对存储状态的管理机制。这个机制，主要使用的是一个叫作Persistent Volume Claim的功能。

### StatefulSet的工作原理
* 首先，StatefulSet的控制器直接管理的是Pod。
    * 这是因为，StatefulSet里的不同Pod实例，不再像ReplicaSet中那样都是完全一样的，有了细微的区别的。
    * 比如，每个Pod的hostname、名字等都是不同的，携带了编号的。
    * 而StatefulSet区分这些实例的方式，就是通过在Pod的名字里加上事先约定好的编号。
* 其次，Kubernetes通过Headless Service，为这些有编号的Pod，在DNS服务器中生成带有同样编号的DNS记录。
* 最后，StatefulSet还为每一个Pod分配并创建一个同样编号的PVC。


如何对StatefulSet进行“滚动更新”（rolling update）？

* 只要修改StatefulSet的Pod模板，就会自动触发“滚动更新”。
* StatefulSet Controller就会按照与Pod编号相反的顺序，从最后一个Pod开始，逐一更新这个StatefulSet管理的每个Pod。而如果更新发生了错误，这次“滚动更新”就会停止。
* StatefulSet的“滚动更新”还允许我们进行更精细的控制，比如金丝雀发布或者灰度发布，这意味着应用的多个实例中被指定的一部分不会被更新到最新的版本。
    * 这个字段，正是StatefulSet的spec.updateStrategy.rollingUpdate的partition字段。

StatefulSet可以说是Kubernetes项目中最为复杂的编排对象。

