# 容器网络配置
容器网络不通了要怎么调试

我们要让容器Network Namespace中的数据包最终发送到物理网卡上，需要完成哪些步骤呢？
* 第一步，就是要让数据包从容器的Network Namespace发送到Host Network Namespace上。
* 第二步，数据包发到了Host Network Namespace之后，还要解决数据包怎么从宿主机上的eth0发送出去的问题。

对于容器从自己的Network Namespace连接到Host Network Namespace的方法，一般来说就只有两类设备接口：
* 一类是veth（Docker启动的容器缺省的网络接口）
* 另一类是macvlan/ipvlan




