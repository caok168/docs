## 问题记录

### Buffer Pool相关
* 1. 对于Buffer Pool而言，里面会存放很多的缓存页以及对应的描述数据，那么假设Buffer Pool里的内存都用尽了，已经没有足够的剩余内存来存放缓存页和描述数据了，此时Buffer Pool里就一点内存都没有了吗？还是说Buffer Pool里会残留一些内存碎片呢？
    当然有。因为Buffer Pool大小是自己定的，很可能Buffer Pool划分完全部的缓存页和描述数据块之后，还剩一点点的内存，这一点的内存放不下任何一个缓存页了，所以这点内存就只能放着不能用，这就是内存碎片。
* 2. 如果你觉得Buffer Pool里会有内存碎片的话，怎么做才能尽可能减少Buffer Pool里的内存碎片呢？
    数据库在Buffer Pool中划分缓存页的时候，会让所有的缓存页和描述数据块都紧密的挨在一起，这样尽可能减少内存浪费，就可以尽可能的减少内存碎片的产生了。

### Buffer Pool中的三个链表
* free链表：磁盘上的数据是如何加载到缓存页里去的。
* flush链表：对缓存页修改之后，flush链表是如何用来记载脏数据页的。
* LRU链表：来判断哪些缓存页是不常用的。


### 线上数据库莫名其妙的随机性能抖动
#### 根本原因 2个
* buffer pool的缓存页都满了，此时执行一个SQL查询很多数据，一下子要把很多缓存页flush到磁盘上去，刷磁盘太慢了，就会导致查询语句执行的很慢；
* 执行更新语句的时候，redo log在磁盘上的所有文件都写满了，此时需要回到第一个redo log文件覆盖写，覆盖写的时候就涉及到第一个redo log文件里很多redo log日志对应的更新操作改动了缓存页，那些缓存还没flush到磁盘，此时就必须把那些缓存页flush到磁盘，才能执行后续的更新语句。

#### 解决问题核心点
* 尽量减少缓存页flush到磁盘的频率
* 尽量提升缓存页flush到磁盘的速度

#### 参数设置
* 参数：innodb_io_capacity：这个参数是告诉数据库采用多大的IO速率把缓存页flush到磁盘里去。
* 参数：innodb_flush_neighbors：在flush缓存页到磁盘的时候，会控制把缓存页临近的其他缓存页也刷到磁盘，但是这样有时候会导致flush的缓存页太多了；参数设置为0，禁止刷临近缓存页，这样就把每次刷新的缓存页数量降低到最少了。

### 快速复制一张表
#### 方式一
* select into from table2 from table1
* insert into table2(a, b, c) select a, b, c from table1

#### 方式二
* 使用mysqldump命令将数据导出成一组INSERT语句
```
mysqldump -h$host -P$port -u$user --add-locks=0 --no-create-info --single-transaction  --set-gtid-purged=OFF db1 t --where="a>900" --result-file=/client_tmp/t.sql
```

* --single-transaction：作用是，在导出数据的时候不需要对表db1.t加表锁，而是使用START TRANSACTION WITH CONSISTENT SNAPSHOT的方法；
* --add-locks设置为0，表示在输出的文件结果里，不增加"LOCK TABLES t WRITE;"
* --no-create-info：不需要导出表结构
* --set-gtid-purged=off：表示的是，不输出跟GTID相关的信息
* --result-file：指定了输出文件的路径，其中client表示生成的文件是在客户端机器上的。

```
mysql -h127.0.0.1 -P13000  -uroot db2 -e "source /client_tmp/t.sql"
```

#### 方式三
* 物理拷贝方法
```
1、create table r like t // 创建一个相同表结构的空表
2、alter table r discard tablespace // r.ibd文件会被删除
3、flush table t for export // 这时候db1 目录下会生成一个t.cfg文件
4、cp t.cfg r.cfg; cp t.ibd r.ibd // 注意读写权限
5、unlock tables // t.cfg文件会被删除
6、alter table r import tablespace // 将这个r.ibd文件作为表r新的表空间
```
