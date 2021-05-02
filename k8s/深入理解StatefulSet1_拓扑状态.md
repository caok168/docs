# 深入理解StatefulSet（一）：拓扑状态

Kubernetes项目很早就在Deployment的基础上，扩展出了对“有状态应用”的初步支持。这个编排功能，就是：StatefulSet。

StatefulSet的设计其实非常容易理解。它把真实世界里的应用状态，抽象成了两种情况：
* 拓扑状态
* 存储状态

StatefulSet的核心功能，就是通过某种方式记录这些状态，然后在Pod被重新创建时，能够为新Pod恢复这些状态。

### 非常实用的概念：Headless Service
Service是Kubernetes项目中用来将一组Pod暴露给外界访问的一种机制。

Service是如何被访问的呢？
* 第一种方式，是以Service的VIP（Virtual IP，即：虚拟IP）方式
* 第二种方式，是以Service的DNS方式。
    * （1）是Normal Service。这种情况下，你访问"my-svc.my-namespace.svc.cluster.local"解析到的，正式my-svc这个Service的VIP，后面的流程和VIP方式一致了。
    * （2）是Headless Service。这种情况下，你访问"my-svc.my-namespace.svc.cluster.local"解析到的，直接就是my-svc代理的某一个Pod的IP地址。**这里的区别在于，Headless Service不需要分配一个VIP，而是可以直接以DNS记录的方式解析出被代理Pod的IP地址。**

### 这样的设计有什么作用呢？
从Headless Service的定义方式看起
```
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
  spec:
    ports:
    - port: 80
      name: web
    clusterIP: none
    selector:
      app: nginx
```
所谓的Headless Service，其实仍是一个标准的Service的YAML文件。只不过，它的clusterIP字段的值是：None

* 当创建了一个Headless Service之后，它所代理的所有Pod的IP地址，都会被绑定一个这样格式的DNS记录，如下所示：
```
<pod-name>.<svc-name>.<namespace>.svc.cluster.local
```
* 这个DNS记录，正式Kubernetes项目为Pod分配的唯一的“可解析身份”(Resolvable Identity)。
* 有了这个“可解析身份”,只要知道了一个Pod的名字，以及它对应的Service的名字，就可以非常确定地通过这条DNS记录访问到Pod的IP地址。

#### StatefulSet又是如何使用这个DNS记录来维持Pod的拓扑状态的呢？
```
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  serviceName: "nginx"
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.9.1
        ports:
        - containerPort: 80
          name: web
```
* 这个YAML文件，和前面nginx-deployment的唯一区别，就是多了一个serviceName=nginx字段。
* 这个字段的作用，就是告诉StatefulSet控制器，在执行控制循环（Control Loop）的时候，请使用nginx这个Headless Service来保证Pod的“可解析身份”。

```
kubectl get pods -w -l app=nginx
```

### 总结
* StatefulSet这个控制器的主要作用之一，就是使用Pod模板创建Pod的时候，对它们进行编号，并且按照编号顺序逐一完成创建工作。而当StatefulSet的“控制循环”发现Pod的“实际状态”与“期望状态”不一致，需要新建或者删除Pod进行“调谐”的时候，它会严格按照这些Pod编号的顺序，逐一完成这些操作。
* 通过Headless Service的方式，StatefulSet为每个Pod创建了一个固定并且稳定的DNS记录，来作为它的访问入口。

