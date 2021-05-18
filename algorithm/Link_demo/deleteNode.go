// 删除链表中的节点

public void deleteNode(ListNode node) {
    node.val = node.next.val;
    node.next = node.next.next;
}

func deleteNode(node ListNode) {
	node.val = node.next.val
	node.next = node.next.next
}

