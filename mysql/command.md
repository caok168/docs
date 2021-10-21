## 导出数据库

### 导出所有数据库
* mysqldump -uroot -proot --all-databases > all.sql

### 导出db1、db2两个数据库的所有数据
* mysqldump -uroot -proot --databases db1 db2 > user.sql

### 导出db1中的a1、a2表
* mysqldump -uroot -proot --databases db1 --tables a1 a2 > db1.sql

### 条件导出，导出db1表a1中id=1的数据
* mysqldump -uroot -proot --databases db1 --tables a1 --where='id=1' > a1.sql
* 字段是字符串，并且导出的sql中不包含drop table，create table
* mysqldump -uroot -proot --no-create-info --databases db1 --tables a1 --where="id='1'" > a1.sql

#### 参考链接
https://www.cnblogs.com/chenmh/p/5300370.html

## 查看数据库情况

### 查看目前最大使用连接数和线程数
* mycli -h172.30.62.66 -uroot -pcar@2019 -P32512 -e "show global status" |grep -E 'Max_used_connections|Threads_connected'

### 查看mysql活动情况
* mycli -h172.30.62.66 -uroot -pcar@2019 -P32512 -e "show FULL PROCESSLIST;"

### 查看最大连接数
* mycli -h172.30.62.66 -uroot -pcar@2019 -P32512 -e "show variables " |grep max_connections 

### 设置最大连接数
* mycli -h172.30.62.66 -uroot -pcar@2019 -P32512 -e "set GLOBAL max_connections=1000"

### 查看binlog是否开启
* show variables like 'log_%';

### 查看所有binlog日志列表
* show master logs; 
* 在mysql的数据保存目录下，可以看到

### 修改日志大小和保存时间
* set global max_binlog_size=1073741824; # 1G
* set global expire_logs_days=7

### 查看日志大小和保存时间
* show variables like '%max_binlog_size%';
* show variables like '%expire_logs_days%'

### 查看慢查询日志
* show variables like '%slow_query_log%'
* show variables like '%long_query_time%'
* show variables like 'log_queries_not_using_indexes'

### 设置慢查询日志
* set global slow_query_log = [ON|OFF]
* set global slow_query_log_file = /sql_log/slowlog.log
* set global long_query_time = xx.xxx秒
* set global log_queries_not_using_indexes = [ON|OFF]

### 分析MySql慢查询日志
* mysqldumpslow
* pt-query-digest

### 监控长时间运行的SQL
* select id, `user`,`host`,DB,command,`time`,state,info from information_schema.PROCESSLIST WHERE TIME>=60

## ibd文件相关
* ALTER TABLE positions DISCARD TABLESPACE; 
* ALTER TABLE positions IMPORT TABLESPACE; SHOW WARNINGS;
如果import tablespace 的时候，报错 ibd文件与表的 ROW_TYPE_COMPACT 不兼容，则需要在建表语句最后 加上 ROW_FORMAT=COMPACT保持一致！

mysql ibd 文件过大问题
* 解决方案：

* 第一种，删除表，然后重新建。drop table 操作自动回收表空间

* 第二种，alter table tablename engine=innodb 。搞定

参考链接
https://www.cnblogs.com/zgq123456/p/9956820.html



## 锁
### 共享锁
* select * from table lock in share mode;

### 互斥锁
* select * from table for update;

### 元数据锁
* metadata locks 在进行DDL的时候

### 表锁
#### 表锁分为表锁和表级的意向锁
* Lock tables xxx READ 这是加表级共享锁
* Lock tables xxx WRITE 这是加表级独占锁
一般没有人用这种方法对表进行加锁



