# 如何提升HTTP1.1性能？

由于门槛低、易监控、自表达等特点，HTTP/1.1在互联网诞生之初就成为最广泛使用的应用层协议。然而，它的性能却很差，最为人诟病的是HTTP头部的传输占用了大量带宽。由于HTTP头部使用ASCII编码方式，这造成它往往达到几KB，而且滥用的Cookie头部进一步增大了体积。

在不升级协议的情况下，有3种优化思路：
* 首先是通过缓存避免发送HTTP请求
    - 客户端缓存响应，可以在有效期内避免发起HTTP请求。
    - 即使缓存过期后，如果服务器端资源未改变，仍然可以通过304响应避免发送包体资源。
    - 浏览器上的私有缓存、服务器上的共享缓存，都对HTTP协议的性能提升有很大意义。
* 其次，如果不得不发起请求，那么就得思考如何才能减少请求的个数；
    - 如将原本客户端处理的重定向请求，移至代理服务器处理可以减少重定向请求的数量
    - 从体验角度，使用懒加载技术延迟加载部分资源，也可以减少请求数量
    - 将多个文件合并后再传输，能够减少许多HTTP头部，而且减少TCP连接数量后也省去握手和慢启动的消耗。当然合并文件的副作用是小文件的更新，会导致整个合并后的大文件重传。
* 最后则是减少服务器响应的体积。
    - 通过亚索响应来降低传输的字节数，选择更优秀的压缩算法能够有效降低传输量，比如用Brotli无损压缩算法替换gzip，或者用WebP格式替换png等格式图片等。


