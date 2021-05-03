# Cond:条件变量的实现机制及避坑指南

### Go标准库的Cond
Cond通常应用于等待某个条件的一组goroutine，等条件变为true的时候，其中一个goroutine或者所有的goroutine都会被唤醒执行。

### Cond的基本用法
三个方法
* Signal：允许调用者Caller唤醒一个等待此Cond的goroutine
* Broadcast：允许调用者Caller唤醒所有等待此Cond的goroutine
* Wait：会把调用者Caller放入Cond的等待队列中并阻塞，直到被Signal或者Broadcast的方法从等待队列中移除并唤醒

### 使用Cond的2个常见错误
* 1.调用Wait的时候没有加锁
* 2.只调用了一次Wait，没有检查等待条件是否满足，结果条件没满足，程序就继续执行了

### 知名项目中Cond的使用
* Kubernetes：定义了优先级队列PriorityQueue这样一个数据结构，用来实现Pod的调用。它内部有是三个Pod的队列，即activeQ、podBackoffQ和unschedulableQ
* Cond有3点特性，是Channel无法替代的
    - Cond和一个Locker关联，可以利用这个Locker对相关的依赖条件更改提供保护
    - Cond可以同时支持Signal和Broadcast方法，而Channel只能同时支持其中一种
    - Cond的Broadcast方法可以被重复调用

