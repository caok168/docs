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
