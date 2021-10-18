## 问题记录

### Buffer Pool相关
* 1. 对于Buffer Pool而言，里面会存放很多的缓存页以及对应的描述数据，那么假设Buffer Pool里的内存都用尽了，已经没有足够的剩余内存来存放缓存页和描述数据了，此时Buffer Pool里就一点内存都没有了吗？还是说Buffer Pool里会残留一些内存碎片呢？
* 2. 如果你觉得Buffer Pool里会有内存碎片的话，怎么做才能尽可能减少Buffer Pool里的内存碎片呢？

