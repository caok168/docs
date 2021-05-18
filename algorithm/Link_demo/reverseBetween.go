// 反转链表II
func reverseBetween(head *ListNode, m int, n int) *ListNode {
    if head == nil || m == n {
		return head
	}
	guard := &ListNode{Val: -1}
	guard.Next = head
	cur := guard
	pre := guard
	tail := guard
	var newNode *ListNode
	index := 0
	for cur != nil && index <= n{
		curNext := cur.Next
		if index < m{
			pre = cur
		}
		if index >= m {
			if index == m {
				tail = cur
			}
			cur.Next = newNode
			newNode = cur
		}
		index++
		cur = curNext
	}
	tail.Next = cur
	pre.Next = newNode
	return guard.Next
}
