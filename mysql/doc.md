# mysql explain之Extra字段讲解

* using index
* using index condition
* using filesort
* using where

### 上面的3者的性能开销是依次放大的。

## Using index
表示使用了覆盖索引，性能上会快很多
覆盖索引是指，索引上的信息足够满足查询请求，不需要再回到主键索引上去取数据。

## Using index condition
与 Using index 的区别在于，用上了索引（利用索引完成字段的筛选过滤），但是索引列不能够包含查询要求的所有字段，需要回表补全字段
回表是指，根据索引找到满足条件的id后，利用id回到主键索引上取出整行或者取出需要的字段

## Using filesort
表示的就是需要排序，MySQL 会给每个线程分配一块内存用于排序，称为 sort_buffer。
"排序"这个动作，可能在内存中完成，也可能需要使用外部排序，这取决于排序所需的内存和参数 sort_buffer_size。
sort_buffer_size，就是 MySQL 为排序开辟的内存（sort_buffer）的大小。如果要排序的数据量小于 sort_buffer_size，排序就在内存中完成。但如果排序数据量太大，内存放不下，则不得不利用磁盘临时文件辅助排序。

全字段排序
可以通过 OPTIMIZER_TRACE结果 查看
"sort_mode": "<fixed_sort_key, additional_fields>"

rowid排序
如果sort_buffer中存放的字段数太多，会造成内存里能够同时放下的的行数很少，就会使用临时文件来进行排序，那么排序的效率就会很差
控制排序的单行长度设置(就表t来说，max_length_for_sort_data值大于等于22时，均采用全字段排序) 具体语法为 SET max_length_for_sort_data = 21;
可以通过 OPTIMIZER_TRACE结果 查看
"sort_mode": "<fixed_sort_key, rowid>"

## Using where
首先你需要知道的是Mysql数据库包含sever层(连接器、查询缓存、分析器、优化器、执行器)与引擎层(innodb,myisam,memory)
Using where表示sever层在收到引擎层返回的行后会进行过滤(即应用WHERE过滤条件)。即会根据查询条件过滤结果集。

select * from t where a=? and b>? order by c limit 0,100
idx_acb (a,c,b)，这种方式才是真正最优的

前面就之前的测试得出的结论是 idx_cab (c,a,b) 最优，idx_ac (a,c)次之。


# 查询性能优化
## 1.如何定位慢的SQL语句
慢查询日志是将mysql服务器中影响数据库性能的相关SQL语句记录到日志文件，通过对这些特殊的SQL语句分析，改进以达到提高数据库性能的目的。

* 用查询缓存处理的查询不加到慢查询日志中，表有零行或一行而不能从索引中受益的查询也不写入慢查询日志。

## 2.如何对SQL进行优化和分析
### 2.1 Explain分析

* select_type : 查询类型，有简单查询、联合查询、子查询等
* type：扫描方式
* all：即全表扫描，如果是这个，尽量去优化
* index：按索引次序扫描，先读索引，再读实际的行，结果还是全表扫描，主要优点是避免了排序。因为索引是排好的
* range：以范围的形式扫描
* const 常量查询，查询过程中整个表最多只会有一条匹配的行
* eq_ref 使用唯一索引查找
* ref： 非唯一索引访问
* possible_keys：可能用到的索引
* key : 实际使用的索引
* rows : 扫描的行数
* key_len： 索引字段最大可能使用的长度
* Extra：其他信息
* Using index：此查询使用了覆盖索引，即通过索引就能返回结果，无须访问表
* Using where：表示 MySQL服务器从存储引擎收到行后再进行“后过滤”（Post-filter）。所谓“后过滤”，就是先读取整行数据，再检查此行是否符合where句的条件，符合就留下，不符合便丢弃。因为检查是在读取行后才进行的，所以称为“后过滤”。
* Using temporary：使用到了临时表，MySQL 使用临时表来实现 distinct 操作
* Using filesort：若查询所需的排序与使用的索引的排序一致，因为索引是已排序的，因此按索引的顺序读取结果返回；否则，在取得结果后，还需要按查询所需的顺序对结果进行排序，这时就会出现Using filesort 。

### 2.2 减少数据访问
* 减少请求的数据量
* 减少服务端扫描的行数

### 2.3 重构查询方式重构SQL
* 切分大查询
* 分解大连接查询



参考链接：
https://blog.csdn.net/wen_fei/article/details/89085813?utm_medium=distribute.pc_relevant_bbs_down.none-task-blog-baidujs-1.nonecase&depth_1-utm_source=distribute.pc_relevant_bbs_down.none-task-blog-baidujs-1.nonecase


## 3.分表

### 创建测试数据库
```
create database test charset=utf8;
```

### 创建测试表
```
use test;		// 先切换到test数据库
create table test_log
(
	time datetime,
	msg varchar(2000)
)
```

### 手动进行分区
#### 批量进行分区
```
alter table test_log partition by range columns(time)(
	partition p20191001 values less than('2019-10-01'),
	partition p20191002 values less than('2019-10-02'),
	partition p20191003 values less than('2019-10-03'),
	partition p20191003 values less than('2019-10-04')
);
```

#### 单条增加分区
```
alter table test_log add partition (partition p20191003 values less than('2019-10-03'));
```

#### 删除分区命令
```
alter table test_log drop partition p20191004;
```

### 插入数据
```
insert into test_log values('2019-10-01 10:11:13', 'hi');
insert into test_log values('2019-10-02 10:12:10', 'ni');
insert into test_log values('2019-10-03 10:12:10', 'hao');
```

### 查看分区表
```
select partition_name, partition_description as val from information_schema.partitions
where table_name='test_log' and table_schema='test';
```

参考链接：
https://www.cnblogs.com/wangsongbai/p/13444620.html
https://www.cnblogs.com/garfieldcgf/p/10143367.html

