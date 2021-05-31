# kubelet　介绍

### kubelet简介
kubelet的主要功能就是定时从某个地方获取节点上的 pod/container的期望状态（运行什么容器、运行的副本数量、网络或者存储如何配置等等），并调用对应的容器平台接口达到这个状态。

kubelet组件运行在Node节点上，维持运行中的Pods以及提供Kubernetes运行时环境，主要完成以下使命：
* 1. 监视分配给该Node节点的pods
* 2. 挂载pod所需要的volumes
* 3. 下载pod的secret
* 4. 通过docker/rkt来运行pod中的容器
* 5. 周期的执行pod中为容器定义的liveness探针
* 6. 上报pod的状态给系统的其他组件
* 7. 上报Node的状态

### kubelet功能模块
在v1.12中，kubelet组件有18个主要的manager（kubernetes/pkg/kubelet/kubelet.go）:
* certificateManager
* cgroupManager
* containerManager
* cpuManager
* nodeContainreManager
* configmapManager
* containerReferenceManager
* evictionManager
* nvidiaGpuManager
* imageGCManager
* kuberuntimeManager
* hostportManager
* podManager
* probeManager
* secretManager
* statusManager
* volumeManager
* tokenManager

#### PLEG
PLEG全称PodLifecycleEvent，PLEG会一直调用container runtime获取本节点的pods，之后比较本模块中之前缓存的pods信息，比较最新的pods中的容器的状态是否发生改变，当状态发生改变的时候，生成一个eventRecord事件，输出到eventChannel中。syncPod模块会接收到eventChannel中的事件，来触发pod同步处理过程，调用container runtime来重建pod，保证pod工作正常。

参考链接：
cnblogs.com/liuxingxing/p/14392399.html


