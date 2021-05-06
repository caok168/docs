# String

String有三种编码模式：
* int编码
* embstr编码
    - 字符串小于等于44字节
    - RedisObject中的元数据、指针和SDS是一块连续的内存区域，可以避免内存碎片
* raw编码
    - 字符串大于44字节
    - Redis不会把SDS和RedisObject布局在一起，而是会给SDS分配独立的空间，并用指针指向SDS结构