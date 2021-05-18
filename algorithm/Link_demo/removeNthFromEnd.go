// 删除链表的倒数第N个节点

func getLength(head *ListNode) (length int) {
    for ; head != nil; head = head.Next {
        length++
    }
    return
}

func removeNthFromEnd(head *ListNode, n int) *ListNode {
    length := getLength(head)
    dummy := &ListNode{0, head}
    cur := dummy
    for i := 0; i < length-n; i++ {
        cur = cur.Next
    }
    cur.Next = cur.Next.Next
    return dummy.Next
}


// 方法二

func removeNthFromEnd(head *ListNode, n int) *ListNode {
    if n <= 0 || head == nil {
        return head
    }
    fast := head
    for i := 1; i <= n && fast != nil; i++{
        fast = fast.Next
    }
    
    if fast == nil {
        return head.Next
    }
    
    
    slow := head
    for fast.Next != nil {
        slow = slow.Next
        fast = fast.Next
    }
    slow.Next = slow.Next.Next
    return head
}
