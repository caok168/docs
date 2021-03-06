# 基础篇 | TCP连接的建立和断开受哪些系统配置影响？

### 四次挥手
* TIME_WAIT状态存在的意义：最后发送的这个ACK包可能会被丢弃掉或者有延迟，这样对端就会再次发送FIN包。如果不维持TIME_WAIT这个状态，那么再次收到对端的FIN包后，本端就会回一个Reset包，这可能会产生一些异常。
* 所以维持TIME_WAIT状态一段时间，可以保障TCP连接正常断开。TIME_WAIT的默认存活时间在Linux上是60s（TCP_TIMEWAIT_LEN），这个时间对于数据中心而言可能还是有些长了，所以有的时候也会修改内核做些优化来减少该值，或者将该值设置为可通过sysctl来调节。


对于CLOSE_WAIT状态而言，系统中没有对应的配置项。但是该状态也是一个危险信号，如果这个状态的TCP连接较多，那往往意味着应用程序有Bug，在某些条件下没有调用close()来关闭连接。
