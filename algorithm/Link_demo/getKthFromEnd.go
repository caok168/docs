// 链表中倒数第k个节点

func getKthFromEnd(head *ListNode, k int) *ListNode {
	var len , root = 0, head

	for  {
		if root == nil{
			break
		}
		len ++
		root = root.Next
	}

	for i:=0; i<len-k; i++ {
		head = head.Next
	}
	return head
}

// 快慢指针也行

