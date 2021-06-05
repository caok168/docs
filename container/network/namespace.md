

Ping(ICMP)：验证IP的可达性

telnet：验证服务的可用性

sudo ip netns list

sudo ip netns delete test1

sudo ip netns add test1

sudo ip netns exec test1 ip a
sudo ip netns exec test1 ip link
sudo ip netns exec test1 ip link set dev lo up

sudo ip link add veth-test1 type veth peer name veth-test2