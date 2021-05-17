# 基于角色的权限控制：RBAC

Kubernetes中所有的API对象，都保存在Etcd里。可是，对这些API对象的操作，却一定都是通过访问kube-apiserver实现的。其中一个非常重要的原因，就是需要APIServer来帮助你做授权工作。

在Kubernetes项目中，负责完成授权（Authorization）工作的机制，就是RBAC：基于角色的访问控制（Role-Based Access Control）。

### 三个最基本的概念
* Role：角色，它其实是一组规则，定义了一组对Kubernetes API对象的操作权限。
* Subject：被作用者，既可以是“人”，也可以是“机器”，也可以是你在Kubernetes里定义的“用户”。
* RoleBinding：定义了“被作用者”和“角色”的绑定关系。

### Role
```
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: mynamespace
  name: example-role
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list", "watch"]
```

### RoleBinding
```
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: example-rolebinding
  namespace: mynamespace
subjects:
- kind: User
  name: example-user
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: example-role
  apiGroup: rbac.authorization.k8s.io
```

通过roleRef这个字段，RoleBinding对象就可以直接通过这个名字，来引用我们前面定义的Role对象（example-role），从而定义了“被作用者（Subject）”和“角色（Role）”之间的绑定关系。

**注意：**
Role和RoleBinding对象都是Namespaced对象（Namespaced Object），它们对权限的限制规则仅在它们自己的Namespace内有效，roleRef也只能引用当前Namespace里的Role对象。

对于非Namespaced（Non-namespaced）对象（比如：Node），或者，某一个Role想要作用于所有的Namespace的时候，就需要使用ClusterRole和ClusterRoleBinding这两个组合了。

这两个API对象的用法和Role和RoleBinding完全一样，只不过，它们的定义里，没有了Namespace字段。


************************************

大多数时候，我们其实都不太使用“用户”这个功能，而是直接使用Kubernetes里的“内置用户”。

**ServiceAccount**

### 1. 创建一个ServiceAccount。
```
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: mynamespace
  name: example-sa
```
### 2. 编写RoleBinding的YAML文件，来为这个ServiceAccount分配权限：
```
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: example-rolebinding
  namespace: mynamespace
subjects:
- kind: ServiceAccount
  name: example-sa
  namespace: mynamespace
roleRef:
  kind: Role
  name: example-role
  apiGroup: rbac.authorization.k8s.io
```

### 3. 这个时候就可以使用这个创建pod了
```
apiVersion: v1
kind: Pod
metadata:
  namespace: mynamespace
  name: sa-token-test
spec:
  containers:
  - name: nginx
    image: nginx:1.7.9
  serviceAccountName: example-sa
```

* 等这个Pod运行起来之后，该ServiceAccount的token，也就是一个Secret对象，被Kubernetes自动挂载到了容器的/var/run/secrets/kubernetes.io/serviceaccount目录下

* 如果一个Pod没有声明serviceAccountName，Kubernetes会自动在它的Namespace下创建一个名叫default的默认ServiceAccount，然后分配给这个Pod。

* 但是这种情况下，这个默认的ServiceAccount并没有关联任何Role。也就是说，此时它有访问APIServer的绝大多数权限。当然，这个访问所需要的Token，还是默认ServiceAccount对应的Secret对象为它提供的。


一个ServiceAccount，在Kubernetes里对应的“用户”的名字是：
* system:serviceaccount:<Namespace名字>:<ServiceAccount名字>
它对应的内置“用户组”的名字，就是：
* system:serviceaccounts:<Namespace名字>


在Kubernetes中已经内置了很多个为系统保留的ClusterRole，它们的名字都以system:开头。

可以通过kubectl get clusterroles查看到它们。

Kubernetes还提供了四个预先定义好的ClusterRole来供用户直接使用：
* 1. cluster-admin
* 2. admin
* 3. edit
* 4. view
