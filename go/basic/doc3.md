# 知识点汇总

## Go语言如何深度拷贝对象
深度复制可以基于reflect包的反射机制完成，但是全部重头手写的话会很繁琐。

最简单的方式是基于序列化和反序列化来实现对象的深度复制：
```
func deepCopy(dst, src interface{}) error {
    var buf bytes.Buffer
    if err := gob.NewEncoder(&buf).Encode(src); err != nil{
        return err
    }
    return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
```
Gob和bytes.Buffer简单组合就搞定了。当然，Gob的底层也是基于reflect包实现的。

## go里面一个协程能保证绑定在一个内核线程上面的
golang的runtime提供了一个LockOSThread的函数，该方法的作用是可以让当前协程绑定并独立一个线程M。

绑定线程的哪个协程new出来的子协程不会继承lockOSThread特性。

#### 什么场景下用到runtime.LockOSThread?
* 我们知道golang的schedule可以理解为公平协作调度和抢占的综合体，不支持优先级调度。
* 当开了几十万个goroutine，并且大多数协程已经在runq等待调度了，那么如果有一个重要的周期性的协程需要优先执行该怎么办？
* 可以借助runtime.LockOSThread()方法来绑定线程，绑定线程M后的好处在于，可以由system kernel内核来调度，因为它本质是线程了。
#### 总结
runtime.LockOSThread会锁定当前协程只跑在一个系统线程上，这个线程也只跑该协程。他们是互相牵制的。


