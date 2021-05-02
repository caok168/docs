# 深入理解Pod对象（二）：使用进阶
先从一种特殊的Volume开始

这种特殊的Volume，叫作Projected Volume，可以翻译成“投射数据卷”。

在Kubernetes中，有几种特殊的Volume，它们存在的意义不是为了存放容器里的数据，也不是用来进行容器和宿主机之间的数据交换。

这些特殊的Volume的作用，是为容器提供预先定义好的数据。所以从容器的角度来看，这些Volume里的信息就是仿佛是被Kubernetes“投射（Project）”进入容器当中的。

到目前为止，Kubernetes支持的Projected Volume一共有四种：
* Secret
* ConfigMap
* Downward API
* ServiceAccountToken



