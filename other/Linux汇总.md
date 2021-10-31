## 汇总

### 一. Crontab 无法执行脚本成功

* 1. 需要在脚本中加入
SHELL=/bin/bash
PATH=/sbin:/bin:/usr/sbin:/usr/bin

* 2. 权限问题
添加 +x执行权限

* 3. 路径问题

* 4. 时差问题

* 5. 变量问题
命令中含有变量，但crontab执行时没有


