
```

aibee@aibeedeMacBook-Pro conf % cat my.cnf 
[mysqld]
pid-file        = /var/run/mysqld/mysqld.pid
socket          = /var/run/mysqld/mysqld.sock
datadir         = /var/lib/mysql
secure-file-priv= NULL

log_bin             = /var/lib/mysql/mysql-bin.log
binlog_format = MIXED
expire_logs_days    = 7
max_binlog_size     = 1G
binlog_do_db        = car
sync_binlog = 10
innodb_flush_log_at_trx_commit = 10
binlog_cache_size = 1G

# Custom config should go here
!includedir /etc/mysql/conf.d/

```

docker-compose.yml
```

version: '3.1'

services:

  mysql:
    image: mysql:8.0
    restart: always
    container_name: mysql
    environment:
      MYSQL_ROOT_PASSWORD: car@2019
      MYSQL_DATABASE: car
    command:
      --default-authentication-plugin=mysql_native_password
      --character-set-server=utf8mb4
      --collation-server=utf8mb4_general_ci
      --explicit_defaults_for_timestamp=true
      --lower_case_table_names=1
    ports:
      - 3306:3306
    volumes:
      - ./data/mysql:/var/lib/mysql
      - ./conf/my.conf:/etc/mysql/my.conf

```

