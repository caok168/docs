## sysbench 

### 1、安装
```
sudo apt update
sudo apt install sysbench -y
sysbench --version
```
### 2、基于sysbench构造测试表和测试数据
```
sysbench --db-driver=mysql --time=300 --threads=10 --report-interval=1 --mysql-host=127.0.0.1 --mysql-port=3306 --mysql-user=root --mysql-password=car@2019 --mysql-db=test_db --tables=20 --table-size=1000000 oltp_read_write --db-ps-mode=disable prepare
```

```
root@aibee243:~# sysbench --db-driver=mysql --time=300 --threads=10 --report-interval=1 --mysql-host=127.0.0.1 --mysql-port=3306 --mysql-user=root --mysql-password=car@2019 --mysql-db=test_db --tables=20 --table-size=1000000 oltp_read_write --db-ps-mode=disable prepare
sysbench 1.0.11 (using system LuaJIT 2.1.0-beta3)

Initializing worker threads...

Creating table 'sbtest2'...
Creating table 'sbtest3'...
Creating table 'sbtest5'...
Creating table 'sbtest6'...
Creating table 'sbtest1'...
Creating table 'sbtest10'...
Creating table 'sbtest7'...
Creating table 'sbtest8'...
Creating table 'sbtest9'...
Creating table 'sbtest4'...
Inserting 1000000 records into 'sbtest9'
Inserting 1000000 records into 'sbtest5'
Inserting 1000000 records into 'sbtest10'
Inserting 1000000 records into 'sbtest2'
Inserting 1000000 records into 'sbtest7'
Inserting 1000000 records into 'sbtest1'
Inserting 1000000 records into 'sbtest8'
Inserting 1000000 records into 'sbtest3'
Inserting 1000000 records into 'sbtest4'
Inserting 1000000 records into 'sbtest6'
```

#### 参数介绍
* --db-driver=mysql：基于mysql的驱动去连接mysql数据库
* --time=300：连续访问300秒
* --threads=10：用10个线程模拟并发访问
* --report-interval=1：每隔1秒输出一下压测情况
* --mysql-host=127.0.0.1
* --mysql-port=3306
* --mysql-user=root
* --mysql-password=pwd
* --mysql-db=test_db --tables=20 --table_size=100000 ：在test_db这个库里，构造20个测试表，每个测试表构造100万条测试数据
* oltp_read_write：执行oltp数据库的读写测试
* --db-ps-mode=disable：禁止ps模式
* prepare：参照这个命令的设置去构造出来我们需要的数据库里的数据，他会自动创建20个测试表，每个表里创建100万条测试数据。

### 对数据库进行360的全方位测试
#### 测试数据库的综合读写TPS，使用的是oltp_read_write模式
```
sysbench --db-driver=mysql --time=300 --threads=10 --report-interval=1 --mysql-host=127.0.0.1 --mysql-port=3306 --mysql-user=root --mysql_password=car@2019 --mysql-db=test_db --tables=20 --table-size=1000000 oltp_read_write --db-ps-mode=disable run
```

#### 测试数据库的只读性能，使用的是oltp_read_only模式
```
sysbench --db-driver=mysql --time=300 --threads=10 --report-interval=1 --mysql-host=127.0.0.1 --mysql-port=3306 --mysql-user=root --mysql_password=car@2019 --mysql-db=test_db --tables=20 --table-size=1000000 oltp_read_only --db-ps-mode=disable run
```

#### 测试数据库的删除性能，使用的是oltp_delete模式：
```
sysbench --db-driver=mysql --time=300 --threads=10 --report-interval=1 --mysql-host=127.0.0.1 --mysql-port=3306 --mysql-user=root --mysql_password=car@2019 --mysql-db=test_db --tables=20 --table-size=1000000 oltp_delete --db-ps-mode=disable run
```

#### 测试数据库的更新索引字段的性能，使用的是oltp_update_index模式：
```
sysbench --db-driver=mysql --time=300 --threads=10 --report-interval=1 --mysql-host=127.0.0.1 --mysql-port=3306 --mysql-user=root --mysql_password=car@2019 --mysql-db=test_db --tables=20 --table-size=1000000 oltp_update_index --db-ps-mode=disable run
```

#### 测试数据库的更新非索引字段的性能，使用的是oltp_update_non_index模式：
```
sysbench --db-driver=mysql --time=300 --threads=10 --report-interval=1 --mysql-host=127.0.0.1 --mysql-port=3306 --mysql-user=root --mysql_password=car@2019 --mysql-db=test_db --tables=20 --table-size=1000000 oltp_update_non_index --db-ps-mode=disable run
```

#### 测试数据库的插入性能，使用的是oltp_insert模式：
```
sysbench --db-driver=mysql --time=300 --threads=10 --report-interval=1 --mysql-host=127.0.0.1 --mysql-port=3306 --mysql-user=root --mysql_password=car@2019 --mysql-db=test_db --tables=20 --table-size=1000000 oltp_insert --db-ps-mode=disable run
```

#### 测试数据库的写入性能，使用的是oltp_write_only模式：
```
sysbench --db-driver=mysql --time=300 --threads=10 --report-interval=1 --mysql-host=127.0.0.1 --mysql-port=3306 --mysql-user=root --mysql_password=car@2019 --mysql-db=test_db --tables=20 --table-size=1000000 oltp_write_only --db-ps-mode=disable run
```

#### 压测完成之后，执行cleanup命令，清理数据
```
sysbench --db-driver=mysql --time=300 --threads=10 --report-interval=1 --mysql-host=127.0.0.1 --mysql-port=3306 --mysql-user=root --mysql_password=car@2019 --mysql-db=test_db --tables=20 --table-size=1000000 oltp_write_only --db-ps-mode=disable cleanup
```
