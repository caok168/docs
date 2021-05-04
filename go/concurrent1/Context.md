# Context：信息穿透上下文

## 适用场景
* 上下文信息传递
* 控制子goroutine的运行
* 超时控制的方法调用
* 可以取消的方法调用

## 基本使用方法
### 4个实现方法
* Deadline方法：会返回这个Context会被取消的截止日期
* Done方法：返回一个Channel对象
* Err：如果Done没有被close，Err方法返回nil；如果Done被close，Err方法会返回Done被close的原因
* Value返回此ctx中和指定的key相关联的value

### 常用的生成顶层Context的方法
* context.Background():返回一个非nil、空的Context，没有任何值，不会被cancel，不会超时，没有截止日期
* context.TODO()：返回一个非nil、空的Context，没有任何值，不会被cancel，不会超时，没有截止日期

### 使用规则
* 一般函数使用Context的时候，会把这个参数放在第一个参数的位置
* 从来不把nil当做Context类型的参数值，可以使用context.Background()创建一个空的上下文对象，也不要使用nil
* Context只用来临时做函数之间的上下文透传，不能持久化Context或者Context长久保存
* key的类型不应该使用字符串类型或者其它内建类型，否则容易在包之间使用Context时候冲突
* 常常使用struct{}作为底层类型定义key的类型

### 创建特殊用途Context的方法
* WithValue：基于parent Context生成一个新的Context，保存了一个key-value键值对，常常用来传递上下文
* WithCancel：返回parent的副本，只是副本中的Done channel是新建的对象，它的类型是cancelCtx
* WithTimeout：其实和WithDeadline一样，只不过一个参数是超时时间，一个参数是截止时间
* WithDeadline：返回一个parent的副本，并且设置了一个不晚于参数d的截止时间，类型为timerCtx（或者是cancelCtx）
