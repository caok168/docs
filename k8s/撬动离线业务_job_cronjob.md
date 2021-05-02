# 撬动离线业务：Job与CronJob

例子：
```
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: resouer/ubuntu-bc
        command: ["sh", "-c", "echo 'scale=10000; 4*a(1)' | bc -l "]
      restartPolicy: Never
  backoffLimit: 4
```

* 跟其他控制器不同的是，Job对象并不要求你定义一个spec.selector来描述要控制哪些Pod。
* 这个Job对象在创建后，它的Pod模板，被自动加上了一个controller-uid=<一个随机字符串> 这样的Label。而这个Job对象本身，则被自动加上了这个Label对应的Selector，从而保证了Job与它所管理的Pod之间的匹配关系。
* Job Controller之所以要使用这种携带了UID的label，就是为了避免不同Job对象所管理的Pod发生重合。



